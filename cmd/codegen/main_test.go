package main

import (
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

type execError struct {
	name string
}

func (e *execError) Error() string {
	return "executable file not found: " + e.name
}
