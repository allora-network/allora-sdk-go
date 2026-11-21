package allora

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testProtoDepsDir = "proto-deps"
	testRepoMarker   = "revision"
)

type protoFixture struct {
	dir      string
	makefile string
	names    []string
}

func TestProtoDepsTrackVersions(t *testing.T) {
	fixture := newProtoFixture(t)

	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)
	fixture.assertRevision(t, "allora-chain", "v1")

	marker := filepath.Join(fixture.dir, testProtoDepsDir, "gogoproto", "untouched")
	if err := os.WriteFile(marker, []byte("unchanged checkout\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("unchanged version refetched gogoproto: %v", err)
	}

	fixture.make(t, append([]string{"proto-deps"}, append(fixture.versions("v1"), "ALLORA_CHAIN_VERSION=v2")...)...)
	fixture.assertRevision(t, "allora-chain", "v2")
	fixture.assertRevision(t, "gogoproto", "v1")
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("allora-chain version change refetched gogoproto: %v", err)
	}
}

func TestProtoDepsUpdateWorksWithoutCheckouts(t *testing.T) {
	fixture := newProtoFixture(t)

	fixture.make(t, append([]string{"proto-deps-update"}, fixture.versions("v1")...)...)
	fixture.make(t, append([]string{"proto-deps"}, fixture.versions("v1")...)...)
	for _, name := range fixture.names {
		fixture.assertRevision(t, name, "v1")
	}
}

func TestGoogleapisUpdateFetchesMovingBranch(t *testing.T) {
	fixture := newProtoFixture(t)
	versions := fixture.movingBranchVersions()

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

func newProtoFixture(t *testing.T) *protoFixture {
	t.Helper()
	dir := t.TempDir()
	fixture := &protoFixture{
		dir:      dir,
		makefile: filepath.Join(dir, "Makefile"),
		names: []string{
			"gogoproto",
			"cosmos-proto",
			"cosmos-sdk",
			"googleapis",
			"allora-chain",
		},
	}

	data, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatal(err)
	}
	makefile := string(data)
	for _, name := range fixture.names {
		for _, owner := range []string{"cosmos", "googleapis", "allora-network"} {
			url := fmt.Sprintf("https://github.com/%s/%s", owner, name)
			makefile = strings.ReplaceAll(makefile, url, "file://"+fixture.remote(name))
		}
	}
	if err := os.WriteFile(fixture.makefile, []byte(makefile), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, name := range fixture.names {
		fixture.initRemote(t, name)
	}
	return fixture
}

func newProtoFixtureCommand(t *testing.T, fixture *protoFixture, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = fixture.dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}

func (f *protoFixture) remote(name string) string {
	return filepath.Join(f.dir, "remotes", name+".git")
}

func (f *protoFixture) work(name string) string {
	return filepath.Join(f.dir, "work", name)
}

func (f *protoFixture) initRemote(t *testing.T, name string) {
	t.Helper()
	runGit(t, "", "init", "--bare", f.remote(name))
	runGit(t, "", "init", "-b", "master", f.work(name))
	runGit(t, f.work(name), "remote", "add", "origin", f.remote(name))
	runGit(t, f.work(name), "config", "user.name", "Codegen Test")
	runGit(t, f.work(name), "config", "user.email", "codegen@example.invalid")
	f.commit(t, name, "first")
	runGit(t, f.work(name), "tag", "v1")
	runGit(t, f.work(name), "push", "origin", "v1")
	f.commit(t, name, "second")
	runGit(t, f.work(name), "tag", "v2")
	runGit(t, f.work(name), "push", "origin", "v2")
}

func (f *protoFixture) commit(t *testing.T, name, contents string) {
	t.Helper()
	work := f.work(name)
	if err := os.WriteFile(filepath.Join(work, testRepoMarker), []byte(contents+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGit(t, work, "add", testRepoMarker)
	runGit(t, work, "commit", "-m", contents)
	runGit(t, work, "push", "--force", "origin", "master")
	runGit(t, work, "push", "--force", "origin", "--tags")
}

func (f *protoFixture) versions(version string) []string {
	variables := []string{
		"GOGOPROTO_VERSION",
		"COSMOS_PROTO_VERSION",
		"COSMOS_SDK_VERSION",
		"GOOGLEAPIS_VERSION",
		"ALLORA_CHAIN_VERSION",
	}
	args := []string{"PROTO_DEPS=" + filepath.Join(f.dir, testProtoDepsDir)}
	for _, variable := range variables {
		args = append(args, variable+"="+version)
	}
	return args
}

func (f *protoFixture) movingBranchVersions() []string {
	args := f.versions("v1")
	for i, arg := range args {
		if arg == "GOOGLEAPIS_VERSION=v1" {
			args[i] = "GOOGLEAPIS_VERSION=master"
		}
	}
	return args
}

func (f *protoFixture) make(t *testing.T, args ...string) {
	t.Helper()
	args = append([]string{"--no-print-directory", "--no-builtin-rules", "-f", f.makefile}, args...)
	newProtoFixtureCommand(t, f, "make", args...)
}

func (f *protoFixture) revision(t *testing.T, name, revision string) string {
	t.Helper()
	return runGit(t, "", "--git-dir", f.remote(name), "rev-parse", revision+"^{commit}")
}

func (f *protoFixture) checkout(t *testing.T, name string) string {
	t.Helper()
	return runGit(t, filepath.Join(f.dir, testProtoDepsDir, name), "rev-parse", "HEAD")
}

func (f *protoFixture) assertRevision(t *testing.T, name, revision string) string {
	t.Helper()
	want := f.revision(t, name, revision)
	if got := f.checkout(t, name); got != want {
		t.Fatalf("%s checkout = %s, want %s (%s)", name, got, want, revision)
	}
	return want
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}
