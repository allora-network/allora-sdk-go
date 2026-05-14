package allora

import (
	"os"
	"regexp"
	"testing"
)

func TestCosmosSDKCodegenVersionMatchesRuntimeReplace(t *testing.T) {
	makefileBytes, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}

	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}

	codegenVersion := mustFindSubmatch(t, string(makefileBytes), `(?m)^COSMOS_SDK_VERSION := (\S+)$`)
	runtimeVersion := mustFindSubmatch(t, string(goModBytes), `(?m)^replace github\.com/cosmos/cosmos-sdk => github\.com/cosmos/cosmos-sdk (\S+)$`)

	if codegenVersion != runtimeVersion {
		t.Fatalf("COSMOS_SDK_VERSION %s must match go.mod replace version %s", codegenVersion, runtimeVersion)
	}
}

func mustFindSubmatch(t *testing.T, text string, pattern string) string {
	t.Helper()

	matches := regexp.MustCompile(pattern).FindStringSubmatch(text)
	if len(matches) != 2 {
		t.Fatalf("expected exactly one submatch for pattern %q", pattern)
	}

	return matches[1]
}
