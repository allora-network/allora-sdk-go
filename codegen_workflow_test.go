package allora

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const (
	testProtoDepsDir = "proto-deps"
	testRepoMarker   = "revision"
	testBranchPin    = "release/v1"
)

var testDepNames = []string{"gogoproto", "cosmos-proto", "cosmos-sdk", "googleapis", "allora-chain"}

// urlPattern matches the remote URLs a Makefile can name.
var urlPattern = regexp.MustCompile(`(?:https?|git|ssh|file)://[^\s"',)]+`)

// protoFixture is a throwaway SDK checkout with local bare remotes standing in
// for the five proto dependencies. Every git and make command runs with an
// isolated git configuration so the developer's global settings cannot change
// the result.
type protoFixture struct {
	root     string
	sdk      string
	makefile string
	env      []string
}

func TestProtoDepsReplacesPlainDirectoryWithoutTouchingEnclosingRepo(t *testing.T) {
	fixture := newProtoFixture(t)
	before := fixture.enclosingState(t)

	plain := filepath.Join(fixture.sdk, testProtoDepsDir, "gogoproto")
	if err := os.MkdirAll(plain, 0o755); err != nil {
		t.Fatal(err)
	}
	placeholder := filepath.Join(plain, "placeholder.txt")
	if err := os.WriteFile(placeholder, []byte("not a checkout\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)

	if _, err := os.Stat(placeholder); !os.IsNotExist(err) {
		t.Fatalf("plain directory survived the sync: %v", err)
	}
	fixture.assertRevision(t, "gogoproto", "v1")

	if after := fixture.enclosingState(t); after != before {
		t.Fatalf("enclosing repo changed:\nbefore: %+v\nafter:  %+v", before, after)
	}
}

func TestProtoDepsReplacesSymlinkToEnclosingRepo(t *testing.T) {
	fixture := newProtoFixture(t)
	before := fixture.enclosingState(t)

	link := filepath.Join(fixture.sdk, testProtoDepsDir, "gogoproto")
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(fixture.sdk, link); err != nil {
		t.Fatal(err)
	}

	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)

	if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("symlink survived the sync: %v", err)
	}
	fixture.assertRevision(t, "gogoproto", "v1")

	if after := fixture.enclosingState(t); after != before {
		t.Fatalf("enclosing repo changed:\nbefore: %+v\nafter:  %+v", before, after)
	}
}

func TestProtoDepsUnchangedPinIsNoOp(t *testing.T) {
	fixture := newProtoFixture(t)
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)
	head := fixture.checkout(t, "gogoproto")

	tracked := filepath.Join(fixture.sdk, testProtoDepsDir, "gogoproto", testRepoMarker)
	if err := os.WriteFile(tracked, []byte("local edit\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Removing the remote makes any refetch fail, so a passing run proves the
	// unchanged pin never reaches the network.
	if err := os.RemoveAll(fixture.remote("gogoproto")); err != nil {
		t.Fatal(err)
	}

	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)

	if got := fixture.checkout(t, "gogoproto"); got != head {
		t.Fatalf("unchanged pin moved the checkout from %s to %s", head, got)
	}
	data, err := os.ReadFile(tracked)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "local edit\n" {
		t.Fatalf("unchanged pin reset a tracked file to %q", data)
	}
}

func TestProtoDepsRepairsDriftedCheckout(t *testing.T) {
	fixture := newProtoFixture(t)
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)
	want := fixture.assertRevision(t, "gogoproto", "v1")

	dep := filepath.Join(fixture.sdk, testProtoDepsDir, "gogoproto")
	fixture.git(t, dep, "fetch", "--depth", "1", "origin", "v2")
	fixture.git(t, dep, "checkout", "--force", "--detach", "FETCH_HEAD")
	if got := fixture.checkout(t, "gogoproto"); got == want {
		t.Fatal("drift setup did not move the checkout")
	}

	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)

	if got := fixture.checkout(t, "gogoproto"); got != want {
		t.Fatalf("drifted checkout = %s, want %s", got, want)
	}
}

func TestProtoDepsRevertsToEarlierPin(t *testing.T) {
	fixture := newProtoFixture(t)
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v2")...)...)
	fixture.assertRevision(t, "gogoproto", "v2")

	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)

	fixture.assertRevision(t, "gogoproto", "v1")
}

func TestProtoDepsRepairsStampThatDisagreesWithCheckout(t *testing.T) {
	fixture := newProtoFixture(t)
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)
	head := fixture.checkout(t, "gogoproto")

	stamp := fixture.stampPath("gogoproto")
	stale := fmt.Sprintf("file://%s\nv1\n%s\n", fixture.remote("gogoproto"), strings.Repeat("0", 40))
	if err := os.WriteFile(stamp, []byte(stale), 0o600); err != nil {
		t.Fatal(err)
	}

	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)

	if got := fixture.checkout(t, "gogoproto"); got != head {
		t.Fatalf("checkout = %s, want %s", got, head)
	}
	if got := fixture.stampCommit(t, "gogoproto"); got != head {
		t.Fatalf("stamp records %s, want the checkout commit %s", got, head)
	}
}

// TestProtoDepsInterruptedRunLeavesTruthfulStamp makes the checkout fail after
// the fetch succeeds. The stamp must not claim the new pin, because the
// checkout did not complete; the next run must finish the move.
func TestProtoDepsInterruptedRunLeavesTruthfulStamp(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores the read-only work tree this test relies on")
	}

	fixture := newProtoFixture(t)
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)

	dep := filepath.Join(fixture.sdk, testProtoDepsDir, "gogoproto")
	if err := os.Chmod(dep, 0o500); err != nil {
		t.Fatal(err)
	}
	fixture.makeExpectFailure(t, append([]string{"proto-deps"}, fixture.versions("v2")...)...)
	if err := os.Chmod(dep, 0o700); err != nil {
		t.Fatal(err)
	}

	if got := fixture.stampPin(t, "gogoproto"); got != "v1" {
		t.Fatalf("stamp claims pin %q after a failed checkout, want v1", got)
	}

	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v2")...)...)
	want := fixture.assertRevision(t, "gogoproto", "v2")
	if got := fixture.stampCommit(t, "gogoproto"); got != want {
		t.Fatalf("stamp records %s after recovery, want %s", got, want)
	}
}

func TestProtoDepsStampRecordsResolvedCommit(t *testing.T) {
	fixture := newProtoFixture(t)
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)

	data, err := os.ReadFile(fixture.stampPath("gogoproto"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("stamp has %d lines, want 3:\n%s", len(lines), data)
	}
	if want := "file://" + fixture.remote("gogoproto"); lines[0] != want {
		t.Fatalf("stamp url = %q, want %q", lines[0], want)
	}
	if lines[1] != "v1" {
		t.Fatalf("stamp pin = %q, want v1", lines[1])
	}
	if want := fixture.checkout(t, "gogoproto"); lines[2] != want {
		t.Fatalf("stamp commit = %q, want %q", lines[2], want)
	}
}

func TestProtoDepsAcceptsBranchPin(t *testing.T) {
	fixture := newProtoFixture(t)
	args := replaceVersion(fixture.versions("v1"), "ALLORA_CHAIN_VERSION", testBranchPin)

	fixture.make(t, append([]string{"proto-deps"}, args...)...)

	fixture.assertRevision(t, "allora-chain", testBranchPin)
}

func TestProtoDepsUpdateWorksWithoutCheckouts(t *testing.T) {
	fixture := newProtoFixture(t)

	fixture.make(t, append([]string{"proto-deps-update"}, fixture.versions("v1")...)...)
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)

	for _, name := range testDepNames {
		fixture.assertRevision(t, name, "v1")
	}
}

func TestProtoDepsUpdateForcesDirtyCheckout(t *testing.T) {
	fixture := newProtoFixture(t)
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)

	dep := filepath.Join(fixture.sdk, testProtoDepsDir, "gogoproto")
	tracked := filepath.Join(dep, testRepoMarker)
	if err := os.WriteFile(tracked, []byte("dirty\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dep, "untracked.txt"), []byte("junk\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	fixture.make(t, append([]string{"proto-deps-update"}, fixture.versions("v2")...)...)

	fixture.assertRevision(t, "gogoproto", "v2")
	data, err := os.ReadFile(tracked)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "second\n" {
		t.Fatalf("tracked file = %q, want the v2 contents", data)
	}
	if _, err := os.Stat(filepath.Join(dep, "untracked.txt")); !os.IsNotExist(err) {
		t.Fatalf("untracked file survived the forced update: %v", err)
	}
}

func TestGoogleapisUpdateFetchesMovingBranch(t *testing.T) {
	fixture := newProtoFixture(t)
	versions := replaceVersion(fixture.versions("v1"), "GOOGLEAPIS_VERSION", "master")

	fixture.make(t, append([]string{"proto-deps-update"}, versions...)...)
	first := fixture.assertRevision(t, "googleapis", "master")
	fixture.commit(t, "googleapis", "third")
	third := fixture.revision(t, "googleapis", "master")

	fixture.make(t, append([]string{"proto-deps"}, versions...)...)
	if got := fixture.checkout(t, "googleapis"); got != first {
		t.Fatalf("proto-deps moved an unchanged branch from %s to %s", first, got)
	}

	fixture.make(t, append([]string{"proto-deps-update"}, versions...)...)
	if got := fixture.checkout(t, "googleapis"); got != third {
		t.Fatalf("proto-deps-update checkout = %s, want origin/master %s", got, third)
	}
}

// TestDevOrdersProtoDepsBeforeCodegen asserts the dependency graph makes
// codegen wait for the proto checkouts, then runs `make -j4 dev` on a fresh
// fixture to show the ordering holds in practice.
func TestDevOrdersProtoDepsBeforeCodegen(t *testing.T) {
	fixture := newProtoFixture(t)

	graph := fixture.makeGraph(t, "dev")
	if !strings.Contains(graph, "\ncodegen: proto-deps") {
		t.Fatalf("codegen does not depend on proto-deps, so `make -j` can run them concurrently:\n%s", graph)
	}

	bin := filepath.Join(fixture.root, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(fixture.root, "go.log")
	checkout := filepath.Join(fixture.sdk, testProtoDepsDir, "gogoproto", ".git")
	script := fmt.Sprintf("#!/bin/sh\nif [ -d '%s' ]; then echo present >> '%s'; else echo missing >> '%s'; fi\n", checkout, log, log)
	if err := os.WriteFile(filepath.Join(bin, "go"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	env := withEnv(fixture.env, "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	fixture.makeWithEnv(t, env, append([]string{"-j4", "dev"}, fixture.versions("v1")...)...)

	data, err := os.ReadFile(log)
	if err != nil {
		t.Fatalf("codegen did not run: %v", err)
	}
	if got := strings.TrimSpace(string(data)); got != "present" {
		t.Fatalf("codegen saw %q, want the proto checkouts in place", got)
	}
}

func newProtoFixture(t *testing.T) *protoFixture {
	t.Helper()

	root := t.TempDir()
	home := filepath.Join(root, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}

	fixture := &protoFixture{
		root:     root,
		sdk:      filepath.Join(root, "sdk"),
		makefile: filepath.Join(root, "sdk", "Makefile"),
		env: withEnv(os.Environ(),
			"HOME="+home,
			"XDG_CONFIG_HOME="+filepath.Join(home, ".config"),
			"GIT_CONFIG_GLOBAL=/dev/null",
			"GIT_CONFIG_NOSYSTEM=1",
			"GIT_CONFIG_COUNT=0",
			"GIT_TERMINAL_PROMPT=0",
		),
	}
	if err := os.MkdirAll(fixture.sdk, 0o755); err != nil {
		t.Fatal(err)
	}

	for _, name := range testDepNames {
		fixture.initRemote(t, name)
	}
	fixture.initEnclosingRepo(t)
	return fixture
}

// initEnclosingRepo builds the SDK checkout the Makefile lives in. The tests
// assert this repository is untouched by a proto dependency sync.
func (f *protoFixture) initEnclosingRepo(t *testing.T) {
	t.Helper()

	f.git(t, f.sdk, "init", "-b", "master")
	f.git(t, f.sdk, "config", "user.name", "Codegen Test")
	f.git(t, f.sdk, "config", "user.email", "codegen@example.invalid")
	f.git(t, f.sdk, "remote", "add", "origin", "https://example.invalid/allora-sdk-go.git")

	if err := os.WriteFile(filepath.Join(f.sdk, "tracked.txt"), []byte("sdk tracked content\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	// The real repository ignores proto-deps/, so the fixture does too. That
	// keeps `git status` meaningful as an assertion about the enclosing repo.
	if err := os.WriteFile(filepath.Join(f.sdk, ".gitignore"), []byte("proto-deps/\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	f.writeMakefile(t)
	f.git(t, f.sdk, "add", "tracked.txt", ".gitignore", "Makefile")
	f.git(t, f.sdk, "commit", "-m", "sdk initial")
}

// writeMakefile copies the repository Makefile into the fixture and points
// every dependency URL at a local bare remote. It fails if any URL in the
// result is not one of the fixture's remotes, so a future URL change cannot
// silently reach the real network.
func (f *protoFixture) writeMakefile(t *testing.T) {
	t.Helper()

	data, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatal(err)
	}
	makefile := string(data)
	for _, name := range testDepNames {
		for _, owner := range []string{"cosmos", "googleapis", "allora-network"} {
			url := fmt.Sprintf("https://github.com/%s/%s", owner, name)
			makefile = strings.ReplaceAll(makefile, url, "file://"+f.remote(name))
		}
	}

	allowed := map[string]bool{}
	for _, name := range testDepNames {
		allowed["file://"+f.remote(name)] = true
	}
	for _, url := range urlPattern.FindAllString(makefile, -1) {
		if !allowed[url] {
			t.Fatalf("fixture Makefile contains %q, which is not a fixture remote; the test would reach the network:\n%s", url, makefile)
		}
	}
	for _, name := range testDepNames {
		if !strings.Contains(makefile, "file://"+f.remote(name)) {
			t.Fatalf("fixture Makefile does not point %s at %s", name, f.remote(name))
		}
	}
	if err := os.WriteFile(f.makefile, []byte(makefile), 0o600); err != nil {
		t.Fatal(err)
	}
}

func (f *protoFixture) remote(name string) string {
	return filepath.Join(f.root, "remotes", name+".git")
}

func (f *protoFixture) work(name string) string {
	return filepath.Join(f.root, "work", name)
}

func (f *protoFixture) stampPath(name string) string {
	return filepath.Join(f.sdk, testProtoDepsDir, name, ".git", "codegen.version")
}

func (f *protoFixture) initRemote(t *testing.T, name string) {
	t.Helper()

	f.git(t, "", "init", "--bare", f.remote(name))
	f.git(t, "", "init", "-b", "master", f.work(name))
	f.git(t, f.work(name), "remote", "add", "origin", f.remote(name))
	f.git(t, f.work(name), "config", "user.name", "Codegen Test")
	f.git(t, f.work(name), "config", "user.email", "codegen@example.invalid")

	f.commit(t, name, "first")
	f.git(t, f.work(name), "tag", "v1")
	f.git(t, f.work(name), "push", "origin", "v1")
	f.git(t, f.work(name), "branch", testBranchPin)
	f.git(t, f.work(name), "push", "origin", testBranchPin)

	f.commit(t, name, "second")
	f.git(t, f.work(name), "tag", "v2")
	f.git(t, f.work(name), "push", "origin", "v2")
}

func (f *protoFixture) commit(t *testing.T, name, contents string) {
	t.Helper()

	work := f.work(name)
	if err := os.WriteFile(filepath.Join(work, testRepoMarker), []byte(contents+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	f.git(t, work, "add", testRepoMarker)
	f.git(t, work, "commit", "-m", contents)
	f.git(t, work, "push", "--force", "origin", "master")
	f.git(t, work, "push", "--force", "origin", "--tags")
}

func (f *protoFixture) versions(version string) []string {
	args := []string{}
	for _, variable := range []string{
		"GOGOPROTO_VERSION",
		"COSMOS_PROTO_VERSION",
		"COSMOS_SDK_VERSION",
		"GOOGLEAPIS_VERSION",
		"ALLORA_CHAIN_VERSION",
	} {
		args = append(args, variable+"="+version)
	}
	return args
}

func replaceVersion(args []string, variable, version string) []string {
	replaced := make([]string, len(args))
	for i, arg := range args {
		if strings.HasPrefix(arg, variable+"=") {
			replaced[i] = variable + "=" + version
			continue
		}
		replaced[i] = arg
	}
	return replaced
}

func (f *protoFixture) make(t *testing.T, args ...string) {
	t.Helper()
	f.makeWithEnv(t, f.env, args...)
}

func (f *protoFixture) makeWithEnv(t *testing.T, env []string, args ...string) {
	t.Helper()

	cmd := exec.Command("make", append([]string{"--no-print-directory", "-f", f.makefile}, args...)...)
	cmd.Dir = f.sdk
	cmd.Env = env
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("make %s: %v\n%s", strings.Join(args, " "), err, output)
	}
}

// makeGraph returns make's target database for the given targets, which shows
// each target's prerequisites.
func (f *protoFixture) makeGraph(t *testing.T, args ...string) string {
	t.Helper()

	cmd := exec.Command("make", append([]string{"--no-print-directory", "-p", "-n", "-f", f.makefile}, args...)...)
	cmd.Dir = f.sdk
	cmd.Env = f.env
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make -p %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

// makeExpectFailure runs make and requires a non-zero exit, returning the
// combined output for assertions.
func (f *protoFixture) makeExpectFailure(t *testing.T, args ...string) string {
	t.Helper()

	cmd := exec.Command("make", append([]string{"--no-print-directory", "-f", f.makefile}, args...)...)
	cmd.Dir = f.sdk
	cmd.Env = f.env
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("make %s succeeded, want a failure\n%s", strings.Join(args, " "), output)
	}
	return string(output)
}

func (f *protoFixture) git(t *testing.T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = f.env
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}

func (f *protoFixture) revision(t *testing.T, name, revision string) string {
	t.Helper()
	return f.git(t, "", "--git-dir", f.remote(name), "rev-parse", revision+"^{commit}")
}

func (f *protoFixture) checkout(t *testing.T, name string) string {
	t.Helper()
	return f.git(t, filepath.Join(f.sdk, testProtoDepsDir, name), "rev-parse", "HEAD")
}

func (f *protoFixture) stampCommit(t *testing.T, name string) string {
	t.Helper()
	return f.stampLines(t, name)[2]
}

func (f *protoFixture) stampPin(t *testing.T, name string) string {
	t.Helper()
	return f.stampLines(t, name)[1]
}

func (f *protoFixture) stampLines(t *testing.T, name string) []string {
	t.Helper()

	data, err := os.ReadFile(f.stampPath(name))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("stamp has %d lines, want 3:\n%s", len(lines), data)
	}
	return lines
}

func (f *protoFixture) assertRevision(t *testing.T, name, revision string) string {
	t.Helper()

	want := f.revision(t, name, revision)
	if got := f.checkout(t, name); got != want {
		t.Fatalf("%s checkout = %s, want %s (%s)", name, got, want, revision)
	}
	return want
}

// repoState captures the parts of the enclosing repository a stray git command
// would disturb: the commit, the branch, the origin URL, and every tracked
// path with its content hash.
type repoState struct {
	head    string
	branch  string
	origin  string
	tracked string
	status  string
}

func (f *protoFixture) enclosingState(t *testing.T) repoState {
	t.Helper()

	return repoState{
		head:    f.git(t, f.sdk, "rev-parse", "HEAD"),
		branch:  f.git(t, f.sdk, "symbolic-ref", "--short", "HEAD"),
		origin:  f.git(t, f.sdk, "remote", "get-url", "origin"),
		tracked: f.git(t, f.sdk, "ls-files", "-s"),
		status:  f.git(t, f.sdk, "status", "--porcelain"),
	}
}

// withEnv returns base with each key=value applied, replacing any existing
// entry for the same key so the last value wins deterministically.
func withEnv(base []string, kv ...string) []string {
	env := make([]string, 0, len(base)+len(kv))
	for _, entry := range base {
		replaced := false
		for _, override := range kv {
			if strings.HasPrefix(entry, strings.SplitN(override, "=", 2)[0]+"=") {
				replaced = true
				break
			}
		}
		if !replaced {
			env = append(env, entry)
		}
	}
	return append(env, kv...)
}
