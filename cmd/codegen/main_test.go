package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

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
	dir := t.TempDir()
	proto := filepath.Join(dir, "query.proto")
	writeFile(t, proto, `syntax = "proto3";
package test.v1;
option go_package = "example.com/test/v1";
service Query {
  rpc Ping(PingRequest) returns (PingResponse);
}
message PingRequest {}
message PingResponse {}
`)

	work := t.TempDir()
	restore := chdir(t, work)
	defer restore()

	// protoc and gofmt are present so the run would otherwise succeed; only
	// goimports is missing. A run that writes gen/ therefore proves the
	// preflight did not run first.
	bin := t.TempDir()
	linkTool(t, bin, "protoc")
	linkTool(t, bin, "gofmt")
	restorePath := setPath(t, bin)
	defer restorePath()

	err := generateFromProtos([]string{"-I", dir, "-f", proto})
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
	dir := t.TempDir()
	proto := filepath.Join(dir, "query.proto")
	writeFile(t, proto, `syntax = "proto3";
package test.v1;
option go_package = "example.com/test/v1";
service Query {
  rpc Ping(PingRequest) returns (PingResponse);
}
message PingRequest {}
message PingResponse {}
`)

	bin := t.TempDir()
	writeFile(t, filepath.Join(bin, "goimports"), "#!/bin/sh\necho 'goimports exploded' >&2\nexit 1\n")
	if err := os.Chmod(filepath.Join(bin, "goimports"), 0o755); err != nil {
		t.Fatal(err)
	}
	linkTool(t, bin, "protoc")
	linkTool(t, bin, "gofmt")

	work := t.TempDir()
	restore := chdir(t, work)
	defer restore()

	restorePath := setPath(t, bin)
	defer restorePath()

	err := generateFromProtos([]string{"-I", dir, "-f", proto})
	if err == nil {
		t.Fatal("expected a goimports runtime failure to fail generation")
	}
	if !strings.Contains(err.Error(), "goimports") {
		t.Fatalf("error %q does not mention goimports", err)
	}
}

// linkTool symlinks a tool from the ambient PATH into dir so a test can build a
// PATH that contains exactly the tools it wants.
func linkTool(t *testing.T, dir, name string) {
	t.Helper()

	path, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s is not installed: %v", name, err)
	}
	if err := os.Symlink(path, filepath.Join(dir, name)); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
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

// setPath replaces PATH with dir plus the directories holding the tools the
// generator shells out to, so a test can control which goimports is found.
func setPath(t *testing.T, dir string) func() {
	t.Helper()

	previous, had := os.LookupEnv("PATH")
	if err := os.Setenv("PATH", dir); err != nil {
		t.Fatal(err)
	}
	return func() {
		if had {
			_ = os.Setenv("PATH", previous)
			return
		}
		_ = os.Unsetenv("PATH")
	}
}

type execError struct {
	name string
}

func (e *execError) Error() string {
	return "executable file not found: " + e.name
}
