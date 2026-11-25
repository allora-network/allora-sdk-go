package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// testProto is a minimal proto file the generator can parse. The fake protoc
// ignores its contents and emits a descriptor set describing the same package.
const testProto = `syntax = "proto3";
package test.v1;
option go_package = "example.com/test/v1";
service Query {
  rpc Ping(PingRequest) returns (PingResponse);
}
message PingRequest {}
message PingResponse {}
`

func TestRequireGoimports(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		lookPath := func(file string) (string, error) {
			return "", &execError{name: file}
		}

		err := requireGoimports(lookPath)
		if err == nil {
			t.Fatal("expected missing goimports to fail")
		}
		for _, want := range []string{"goimports", "go install golang.org/x/tools/cmd/goimports@latest"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not contain %q", err, want)
			}
		}
	})

	t.Run("present", func(t *testing.T) {
		called := false
		lookPath := func(file string) (string, error) {
			called = true
			if file != "goimports" {
				t.Fatalf("looked up %q instead of goimports", file)
			}
			return "/test/bin/goimports", nil
		}

		if err := requireGoimports(lookPath); err != nil {
			t.Fatalf("requireGoimports returned an error: %v", err)
		}
		if !called {
			t.Fatal("requireGoimports did not look up goimports")
		}
	})
}

// TestGenerateFromProtosFailsWhenGoimportsIsMissing proves the preflight runs
// before any file is written: the run reports the missing tool and leaves the
// output directory untouched.
func TestGenerateFromProtosFailsWhenGoimportsIsMissing(t *testing.T) {
	work := t.TempDir()
	restore := chdir(t, work)
	defer restore()

	// protoc and gofmt are present so the run would otherwise succeed; only
	// goimports is missing. A run that writes gen/ therefore proves the
	// preflight did not run first.
	bin := t.TempDir()
	writeFakeProtoc(t, bin)
	writeFakeGofmt(t, bin)
	restorePath := setPath(t, bin)
	defer restorePath()

	err := generateFromProtos(testGenerateArgs(t))
	if err == nil {
		t.Fatal("expected generation to fail without goimports")
	}
	if !strings.Contains(err.Error(), "goimports") {
		t.Fatalf("error %q does not mention goimports", err)
	}
	if _, statErr := os.Stat(filepath.Join(work, "gen")); !os.IsNotExist(statErr) {
		t.Fatalf("generation wrote output before the preflight: %v", statErr)
	}
}

// TestGenerateFromProtosFailsWhenGoimportsFails proves a goimports runtime
// failure is fatal rather than a warning.
func TestGenerateFromProtosFailsWhenGoimportsFails(t *testing.T) {
	work := t.TempDir()
	restore := chdir(t, work)
	defer restore()

	bin := t.TempDir()
	writeFakeProtoc(t, bin)
	writeFakeGofmt(t, bin)
	writeExecutable(t, filepath.Join(bin, "goimports"), "#!/bin/sh\necho 'goimports exploded' >&2\nexit 1\n")
	restorePath := setPath(t, bin)
	defer restorePath()

	err := generateFromProtos(testGenerateArgs(t))
	if err == nil {
		t.Fatal("expected a goimports runtime failure to fail generation")
	}
	if !strings.Contains(err.Error(), "goimports") {
		t.Fatalf("error %q does not mention goimports", err)
	}
}

// TestGenerateFromProtosSucceedsWithFakeTools proves the whole pipeline runs
// end to end against the fake tools, so the failure tests above are not passing
// only because generation stops early for an unrelated reason.
func TestGenerateFromProtosSucceedsWithFakeTools(t *testing.T) {
	work := t.TempDir()
	restore := chdir(t, work)
	defer restore()

	bin := t.TempDir()
	writeFakeProtoc(t, bin)
	writeFakeGofmt(t, bin)
	writeExecutable(t, filepath.Join(bin, "goimports"), "#!/bin/sh\nexit 0\n")
	restorePath := setPath(t, bin)
	defer restorePath()

	if err := generateFromProtos(testGenerateArgs(t)); err != nil {
		t.Fatalf("generation failed: %v", err)
	}
	for _, name := range []string{
		"gen/wrapper/client.test.go",
		"gen/rest/client.test.go",
		"gen/grpc/client.test.go",
		"gen/interfaces/client.test.go",
		"gen/interfaces/client.go",
		"gen/rest/client.go",
		"gen/grpc/client.go",
		"gen/wrapper/client.go",
	} {
		if _, err := os.Stat(filepath.Join(work, name)); err != nil {
			t.Errorf("expected %s: %v", name, err)
		}
	}
}

func TestGenerateFromProtosFailsWithoutResolvedInputs(t *testing.T) {
	work := t.TempDir()
	restore := chdir(t, work)
	defer restore()

	bin := t.TempDir()
	writeFakeProtoc(t, bin)
	writeFakeGofmt(t, bin)
	writeExecutable(t, filepath.Join(bin, "goimports"), "#!/bin/sh\nexit 0\n")
	restorePath := setPath(t, bin)
	defer restorePath()

	missing := filepath.Join(t.TempDir(), "missing.proto")
	include := t.TempDir()
	err := generateFromProtos([]string{"-I", include, "-f", missing, "-I", include, "-f", missing})
	if err == nil {
		t.Fatal("expected an unresolved proto file to fail generation")
	}
	for _, want := range []string{missing, "include paths"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not contain %q", err, want)
		}
	}
	if _, statErr := os.Stat(filepath.Join(work, "gen")); !os.IsNotExist(statErr) {
		t.Fatalf("generation wrote output for an unresolved input: %v", statErr)
	}
}

func TestGenerateFromProtosFailsWithoutServices(t *testing.T) {
	work := t.TempDir()
	restore := chdir(t, work)
	defer restore()

	bin := t.TempDir()
	writeFakeProtoc(t, bin)
	writeFakeGofmt(t, bin)
	writeExecutable(t, filepath.Join(bin, "goimports"), "#!/bin/sh\nexit 0\n")
	restorePath := setPath(t, bin)
	defer restorePath()

	include := filepath.Join(t.TempDir(), "include")
	if err := os.Mkdir(include, 0o755); err != nil {
		t.Fatal(err)
	}
	proto := filepath.Join(include, "query.proto")
	if err := os.WriteFile(proto, []byte(testProto), 0o600); err != nil {
		t.Fatal(err)
	}
	err := generateFromProtos([]string{"-I", include, "-I", fakeProtocNoServicesArg, "-f", proto})
	if err == nil {
		t.Fatal("expected zero services to fail generation")
	}
	if !strings.Contains(err.Error(), "no services") {
		t.Fatalf("error %q does not mention missing services", err)
	}
	if _, statErr := os.Stat(filepath.Join(work, "gen")); !os.IsNotExist(statErr) {
		t.Fatalf("generation wrote output with zero services: %v", statErr)
	}
}

func TestBuildFileDescriptorSetForFilesUsesPrivateTemporaryDirectory(t *testing.T) {
	bin := t.TempDir()
	writeFakeProtoc(t, bin)
	restorePath := setPath(t, bin)
	defer restorePath()

	t.Setenv("TMPDIR", t.TempDir())
	stalePath := filepath.Join(os.TempDir(), "descriptor-set.pb")
	stale := []byte("stale descriptor data")
	if err := os.WriteFile(stalePath, stale, 0o600); err != nil {
		t.Fatal(err)
	}
	fds, err := buildFileDescriptorSetForFiles([]string{"query.proto"}, nil)
	if err != nil {
		t.Fatalf("building descriptor set: %v", err)
	}
	if len(fds.File) != 1 || len(fds.File[0].Service) != 1 {
		t.Fatalf("descriptor set was not read from the current run: %v", fds)
	}
	contents, err := os.ReadFile(stalePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != string(stale) {
		t.Fatal("generator used the shared descriptor-set path")
	}
	if matches, err := filepath.Glob(os.Getenv("TMPDIR") + "/allora-codegen-*"); err != nil {
		t.Fatal(err)
	} else if len(matches) != 0 {
		t.Fatalf("descriptor temp directory remains: %v", matches)
	}
}

func TestGenerateFromProtosFailsWithoutGoPackage(t *testing.T) {
	work := t.TempDir()
	restore := chdir(t, work)
	defer restore()

	bin := t.TempDir()
	writeFakeProtoc(t, bin)
	writeFakeGofmt(t, bin)
	writeExecutable(t, filepath.Join(bin, "goimports"), "#!/bin/sh\nexit 0\n")
	restorePath := setPath(t, bin)
	defer restorePath()

	args := append(testGenerateArgs(t), "-I", fakeProtocWithoutGoPackageArg)
	err := generateFromProtos(args)
	if err == nil {
		t.Fatal("expected a missing go_package to fail generation")
	}
	for _, want := range []string{"go_package", "query.proto"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not contain %q", err, want)
		}
	}
	if _, statErr := os.Stat(filepath.Join(work, "gen")); !os.IsNotExist(statErr) {
		t.Fatalf("generation wrote output without go_package: %v", statErr)
	}
}

func TestGenerateFromProtosReportsAllFilesWithoutModules(t *testing.T) {
	work := t.TempDir()
	restore := chdir(t, work)
	defer restore()

	bin := t.TempDir()
	writeFakeProtoc(t, bin)
	writeFakeGofmt(t, bin)
	writeExecutable(t, filepath.Join(bin, "goimports"), "#!/bin/sh\nexit 0\n")
	restorePath := setPath(t, bin)
	defer restorePath()

	args := append(testGenerateArgs(t), "-I", fakeProtocEmptyServiceArg)
	err := generateFromProtos(args)
	if err == nil {
		t.Fatal("expected an empty service to leave no modules")
	}
	proto := args[3]
	for _, want := range []string{"no modules", proto} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not contain %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "no services") {
		t.Fatalf("error %q incorrectly blames the absence of services", err)
	}
	if _, statErr := os.Stat(filepath.Join(work, "gen")); !os.IsNotExist(statErr) {
		t.Fatalf("generation wrote output without modules: %v", statErr)
	}
}

// TestGenerateFromProtosReturnsParseError proves an unparseable flag is
// reported to the caller instead of exiting the process.
func TestGenerateFromProtosReturnsParseError(t *testing.T) {
	err := generateFromProtos([]string{"-not-a-flag"})
	if err == nil {
		t.Fatal("expected an unknown flag to return an error")
	}
	if !strings.Contains(err.Error(), "parsing generate flags") {
		t.Fatalf("error %q does not mention flag parsing", err)
	}
}

// TestGenerateFromProtosReturnsHelp proves -h is reported as flag.ErrHelp so
// main can print usage and exit zero.
func TestGenerateFromProtosReturnsHelp(t *testing.T) {
	err := generateFromProtos([]string{"-h"})
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("error = %v, want flag.ErrHelp", err)
	}
}

// testGenerateArgs writes the test proto into a fresh include directory and
// returns the flags the generator needs to process it.
func testGenerateArgs(t *testing.T) []string {
	t.Helper()

	include := t.TempDir()
	proto := filepath.Join(include, "query.proto")
	if err := os.WriteFile(proto, []byte(testProto), 0o600); err != nil {
		t.Fatal(err)
	}
	return []string{"-I", include, "-f", proto}
}

// writeFakeGofmt installs a gofmt that accepts the arguments the generator
// passes and leaves the files alone.
func writeFakeGofmt(t *testing.T, dir string) {
	t.Helper()

	writeExecutable(t, filepath.Join(dir, "gofmt"), "#!/bin/sh\nexit 0\n")
}

func writeExecutable(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(contents), 0o755); err != nil {
		t.Fatal(err)
	}
}

func chdir(t *testing.T, dir string) func() {
	t.Helper()

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatal(err)
		}
	}
}

// setPath replaces PATH with dir, so a test controls exactly which tools the
// generator can find. The fake protoc is the test binary re-executed, so it
// needs the marker variable in its environment.
func setPath(t *testing.T, dir string) func() {
	t.Helper()

	previousPath, hadPath := os.LookupEnv("PATH")
	previousMarker, hadMarker := os.LookupEnv(fakeProtocEnv)
	if err := os.Setenv("PATH", dir); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv(fakeProtocEnv, "1"); err != nil {
		t.Fatal(err)
	}
	return func() {
		restoreEnv("PATH", previousPath, hadPath)
		restoreEnv(fakeProtocEnv, previousMarker, hadMarker)
	}
}

func restoreEnv(key, value string, had bool) {
	if had {
		_ = os.Setenv(key, value)
		return
	}
	_ = os.Unsetenv(key)
}

type execError struct {
	name string
}

func (e *execError) Error() string {
	return "executable file not found: " + e.name
}
