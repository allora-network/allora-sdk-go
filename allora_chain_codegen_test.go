package allora

import (
	"os"
	"regexp"
	"testing"
)

func TestAlloraChainCodegenVersionIsPinnedReleaseTagOrCommit(t *testing.T) {
	makefileBytes, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}

	version := mustFindSubmatch(t, string(makefileBytes), `(?m)^ALLORA_CHAIN_VERSION := (\S+)$`)
	if !regexp.MustCompile(`^(v[0-9]+\.[0-9]+\.[0-9]+(?:[-+][A-Za-z0-9.-]+)?|[0-9a-f]{40})$`).MatchString(version) {
		t.Fatalf("ALLORA_CHAIN_VERSION %q must be a release tag or full commit SHA", version)
	}
}
