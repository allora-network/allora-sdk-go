#!/usr/bin/make -f

# --- Version Variables
GOGOPROTO_VERSION := v1.7.0
COSMOS_PROTO_VERSION := v1.0.0-beta.5
COSMOS_SDK_VERSION := v0.50.14
GOOGLEAPIS_VERSION := master
ALLORA_CHAIN_VERSION := v0.17.0

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


# Each stamp records one version. A different version refetches only that checkout.
.PHONY: proto-deps-update-gogoproto proto-deps-update-cosmos-proto proto-deps-update-cosmos-sdk proto-deps-update-googleapis proto-deps-update-allora-chain

proto-deps-update-gogoproto:
	rm -f "$(GOGOPROTO_DIR)"/.git/codegen.version.*
	$(MAKE) --no-print-directory "$(GOGOPROTO_DIR)/.git/codegen.version.$(GOGOPROTO_VERSION)"

proto-deps-update-cosmos-proto:
	rm -f "$(COSMOS_PROTO_DIR)"/.git/codegen.version.*
	$(MAKE) --no-print-directory "$(COSMOS_PROTO_DIR)/.git/codegen.version.$(COSMOS_PROTO_VERSION)"

proto-deps-update-cosmos-sdk:
	rm -f "$(COSMOS_SDK_DIR)"/.git/codegen.version.*
	$(MAKE) --no-print-directory "$(COSMOS_SDK_DIR)/.git/codegen.version.$(COSMOS_SDK_VERSION)"

proto-deps-update-googleapis:
	rm -f "$(GOOGLEAPIS_DIR)"/.git/codegen.version.*
	$(MAKE) --no-print-directory "$(GOOGLEAPIS_DIR)/.git/codegen.version.$(GOOGLEAPIS_VERSION)"

proto-deps-update-allora-chain:
	rm -f "$(ALLORA_CHAIN_DIR)"/.git/codegen.version.*
	$(MAKE) --no-print-directory "$(ALLORA_CHAIN_DIR)/.git/codegen.version.$(ALLORA_CHAIN_VERSION)"

define fetch_proto_dep
	@mkdir -p "$(dir $(1))"
	@if ! git -C "$(1)" rev-parse --git-dir >/dev/null 2>&1; then \
		rm -rf "$(1)"; \
		git init "$(1)"; \
	fi
	@if ! git -C "$(1)" remote get-url origin >/dev/null 2>&1; then \
		git -C "$(1)" remote add origin "$(2)"; \
	fi
	git -C "$(1)" remote set-url origin "$(2)"
	git -C "$(1)" fetch --depth 1 origin "$(3)"
	git -C "$(1)" checkout --detach FETCH_HEAD
	git -C "$(1)" reset --hard FETCH_HEAD
	rm -f "$(4)"/codegen.version.*
	touch "$@"
endef

$(GOGOPROTO_DIR)/.git/codegen.version.$(GOGOPROTO_VERSION):
	$(call fetch_proto_dep,$(GOGOPROTO_DIR),https://github.com/cosmos/gogoproto,$(GOGOPROTO_VERSION),$(GOGOPROTO_DIR)/.git)

$(COSMOS_PROTO_DIR)/.git/codegen.version.$(COSMOS_PROTO_VERSION):
	$(call fetch_proto_dep,$(COSMOS_PROTO_DIR),https://github.com/cosmos/cosmos-proto,$(COSMOS_PROTO_VERSION),$(COSMOS_PROTO_DIR)/.git)

$(COSMOS_SDK_DIR)/.git/codegen.version.$(COSMOS_SDK_VERSION):
	$(call fetch_proto_dep,$(COSMOS_SDK_DIR),https://github.com/cosmos/cosmos-sdk,$(COSMOS_SDK_VERSION),$(COSMOS_SDK_DIR)/.git)

$(GOOGLEAPIS_DIR)/.git/codegen.version.$(GOOGLEAPIS_VERSION):
	$(call fetch_proto_dep,$(GOOGLEAPIS_DIR),https://github.com/googleapis/googleapis,$(GOOGLEAPIS_VERSION),$(GOOGLEAPIS_DIR)/.git)

$(ALLORA_CHAIN_DIR)/.git/codegen.version.$(ALLORA_CHAIN_VERSION):
	$(call fetch_proto_dep,$(ALLORA_CHAIN_DIR),https://github.com/allora-network/allora-chain,$(ALLORA_CHAIN_VERSION),$(ALLORA_CHAIN_DIR)/.git)

.PHONY: proto-deps
proto-deps: \
  $(GOGOPROTO_DIR)/.git/codegen.version.$(GOGOPROTO_VERSION) \
  $(COSMOS_PROTO_DIR)/.git/codegen.version.$(COSMOS_PROTO_VERSION) \
  $(COSMOS_SDK_DIR)/.git/codegen.version.$(COSMOS_SDK_VERSION) \
  $(GOOGLEAPIS_DIR)/.git/codegen.version.$(GOOGLEAPIS_VERSION) \
  $(ALLORA_CHAIN_DIR)/.git/codegen.version.$(ALLORA_CHAIN_VERSION)

.PHONY: proto-deps-update
proto-deps-update: \
  proto-deps-update-gogoproto \
  proto-deps-update-cosmos-proto \
  proto-deps-update-cosmos-sdk \
  proto-deps-update-googleapis \
  proto-deps-update-allora-chain



.PHONY: dev
dev: proto-deps codegen
	@echo "✅ Ready for development."

.PHONY: codegen
codegen:
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
