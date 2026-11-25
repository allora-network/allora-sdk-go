#!/usr/bin/make -f

# --- Version Variables
GOGOPROTO_VERSION := v1.7.2
COSMOS_PROTO_VERSION := v1.0.0-beta.5
COSMOS_SDK_VERSION := v0.50.15
GOOGLEAPIS_VERSION := master
ALLORA_CHAIN_VERSION := a67f2c37667c04ef9515d79e180b8cc9208e905a

# --- Paths
PROTO_DEPS := ./proto-deps
COSMOS_SDK_DIR := $(PROTO_DEPS)/cosmos-sdk
COSMOS_PROTO_DIR := $(PROTO_DEPS)/cosmos-proto
GOGOPROTO_DIR := $(PROTO_DEPS)/gogoproto
GOOGLEAPIS_DIR := $(PROTO_DEPS)/googleapis
ALLORA_CHAIN_DIR := $(PROTO_DEPS)/allora-chain

# --- Codegen inputs
CODEGEN_INCLUDES := \
  -I $(COSMOS_SDK_DIR)/proto \
  -I $(COSMOS_PROTO_DIR)/proto \
  -I $(GOOGLEAPIS_DIR) \
  -I $(GOGOPROTO_DIR) \
  -I $(ALLORA_CHAIN_DIR)/x/emissions/proto \
  -I $(ALLORA_CHAIN_DIR)/x/mint/proto

CODEGEN_FILES := \
  $(COSMOS_SDK_DIR)/proto/cosmos/auth/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/authz/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/bank/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/base/node/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/base/tendermint/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/consensus/v1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/distribution/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/evidence/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/feegrant/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/gov/v1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/params/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/slashing/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/staking/v1beta1/query.proto \
  $(COSMOS_SDK_DIR)/proto/cosmos/tx/v1beta1/service.proto \
  $(ALLORA_CHAIN_DIR)/x/emissions/proto/emissions/v10/query.proto \
  $(ALLORA_CHAIN_DIR)/x/mint/proto/mint/v5/query.proto


# --- Proto dependency checkouts
#
# Each dependency is a standalone clone under $(PROTO_DEPS). The stamp file
# <checkout>/.git/codegen.version holds three lines: the remote URL, the pinned
# revision, and the commit the checkout holds. A run refetches when any of the
# three disagrees with the current configuration or with the checkout's HEAD,
# so an interrupted run, a manual checkout, or a reverted pin all repair
# themselves. An unchanged pin costs a few local git calls and no network.
#
# The guard requires <checkout>/.git to be a directory and git to resolve the
# checkout as its own work tree root. The expected root is built from the
# physical parent directory plus the literal final component, so a symlink
# pointing at the enclosing repository does not pass. Without both checks, a
# plain directory or a symlink inside the enclosing repository would pass, and
# the commands below would rewrite that repository's origin, fetch into it, and
# reset its HEAD.

PROTO_DEPS_FORCE ?=

define sync_proto_dep
	@set -eu; \
	dep_dir='$(1)'; dep_url='$(2)'; dep_pin='$(3)'; \
	stamp="$$dep_dir/.git/codegen.version"; \
	mkdir -p "$$(dirname "$$dep_dir")"; \
	dep_root="$$(cd "$$(dirname "$$dep_dir")" && pwd -P)/$$(basename "$$dep_dir")"; \
	if [ -d "$$dep_dir/.git" ] && \
	   [ "$$(git -C "$$dep_dir" rev-parse --show-toplevel 2>/dev/null)" = "$$dep_root" ]; then \
		:; \
	else \
		rm -rf "$$dep_dir"; \
		git init -q "$$dep_dir"; \
	fi; \
	if git -C "$$dep_dir" remote get-url origin >/dev/null 2>&1; then \
		git -C "$$dep_dir" remote set-url origin "$$dep_url"; \
	else \
		git -C "$$dep_dir" remote add origin "$$dep_url"; \
	fi; \
	head=$$(git -C "$$dep_dir" rev-parse --verify HEAD 2>/dev/null || true); \
	if [ -z '$(PROTO_DEPS_FORCE)' ] && [ -f "$$stamp" ] && \
	   [ "$$(sed -n 1p "$$stamp")" = "$$dep_url" ] && \
	   [ "$$(sed -n 2p "$$stamp")" = "$$dep_pin" ] && \
	   [ "$$(sed -n 3p "$$stamp")" = "$$head" ]; then \
		echo "proto-deps: $$dep_dir is at $$dep_pin ($$head)"; \
		exit 0; \
	fi; \
	echo "proto-deps: fetching $$dep_url $$dep_pin into $$dep_dir"; \
	git -C "$$dep_dir" fetch --depth 1 origin "$$dep_pin"; \
	git -C "$$dep_dir" clean -ffdxq; \
	git -C "$$dep_dir" checkout --force --detach FETCH_HEAD; \
	commit=$$(git -C "$$dep_dir" rev-parse HEAD); \
	printf '%s\n%s\n%s\n' "$$dep_url" "$$dep_pin" "$$commit" > "$$stamp.tmp"; \
	mv "$$stamp.tmp" "$$stamp"; \
	echo "proto-deps: $$dep_dir is at $$dep_pin ($$commit)"
endef

.PHONY: proto-deps-gogoproto proto-deps-cosmos-proto proto-deps-cosmos-sdk proto-deps-googleapis proto-deps-allora-chain

proto-deps-gogoproto:
	$(call sync_proto_dep,$(GOGOPROTO_DIR),https://github.com/cosmos/gogoproto,$(GOGOPROTO_VERSION))

proto-deps-cosmos-proto:
	$(call sync_proto_dep,$(COSMOS_PROTO_DIR),https://github.com/cosmos/cosmos-proto,$(COSMOS_PROTO_VERSION))

proto-deps-cosmos-sdk:
	$(call sync_proto_dep,$(COSMOS_SDK_DIR),https://github.com/cosmos/cosmos-sdk,$(COSMOS_SDK_VERSION))

proto-deps-googleapis:
	$(call sync_proto_dep,$(GOOGLEAPIS_DIR),https://github.com/googleapis/googleapis,$(GOOGLEAPIS_VERSION))

proto-deps-allora-chain:
	$(call sync_proto_dep,$(ALLORA_CHAIN_DIR),https://github.com/allora-network/allora-chain,$(ALLORA_CHAIN_VERSION))

.PHONY: proto-deps
proto-deps: \
  proto-deps-gogoproto \
  proto-deps-cosmos-proto \
  proto-deps-cosmos-sdk \
  proto-deps-googleapis \
  proto-deps-allora-chain

.PHONY: proto-deps-update-gogoproto proto-deps-update-cosmos-proto proto-deps-update-cosmos-sdk proto-deps-update-googleapis proto-deps-update-allora-chain

proto-deps-update-gogoproto: PROTO_DEPS_FORCE := 1
proto-deps-update-gogoproto: proto-deps-gogoproto

proto-deps-update-cosmos-proto: PROTO_DEPS_FORCE := 1
proto-deps-update-cosmos-proto: proto-deps-cosmos-proto

proto-deps-update-cosmos-sdk: PROTO_DEPS_FORCE := 1
proto-deps-update-cosmos-sdk: proto-deps-cosmos-sdk

proto-deps-update-googleapis: PROTO_DEPS_FORCE := 1
proto-deps-update-googleapis: proto-deps-googleapis

proto-deps-update-allora-chain: PROTO_DEPS_FORCE := 1
proto-deps-update-allora-chain: proto-deps-allora-chain

.PHONY: proto-deps-update
proto-deps-update: PROTO_DEPS_FORCE := 1
proto-deps-update: proto-deps


.PHONY: dev
dev: codegen
	@echo "✅ Ready for development."

.PHONY: codegen
codegen: proto-deps
	@echo "🤖 Generating clients into ./gen"
	@go run cmd/codegen/main.go $(CODEGEN_INCLUDES) $(foreach f,$(CODEGEN_FILES),-f $(f))

.PHONY: clean
clean:
	@rm -rf gen/

.PHONY: lint
lint:
	@echo "🤖 Running golangci-lint"
	@golangci-lint run --timeout=10m

.PHONY: lint-fix
lint-fix:
	@echo "🤖 Running golangci-lint (fix)"
	@golangci-lint run --fix --timeout=10m
