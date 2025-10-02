#!/usr/bin/make -f

# --- Version Variables
GOGOPROTO_VERSION := v1.7.0
COSMOS_PROTO_VERSION := v1.0.0-beta.5
COSMOS_SDK_VERSION := v0.50.13
GOOGLEAPIS_VERSION := master
ALLORA_CHAIN_VERSION := v0.12.3

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
  $(ALLORA_CHAIN_DIR)/x/emissions/proto/emissions/v9/query.proto \
  $(ALLORA_CHAIN_DIR)/x/mint/proto/mint/v5/query.proto


# --- Git dependencies
$(GOGOPROTO_DIR)/.git:
	rm -rf "$(GOGOPROTO_DIR)"
	git clone --depth 1 --single-branch --branch $(GOGOPROTO_VERSION) \
	  https://github.com/cosmos/gogoproto "$(GOGOPROTO_DIR)"

$(COSMOS_PROTO_DIR)/.git:
	rm -rf "$(COSMOS_PROTO_DIR)"
	git clone --depth 1 --single-branch --branch $(COSMOS_PROTO_VERSION) \
	  https://github.com/cosmos/cosmos-proto "$(COSMOS_PROTO_DIR)"

$(COSMOS_SDK_DIR)/.git:
	rm -rf "$(COSMOS_SDK_DIR)"
	git clone --depth 1 --single-branch --branch $(COSMOS_SDK_VERSION) \
	  https://github.com/cosmos/cosmos-sdk "$(COSMOS_SDK_DIR)"

$(GOOGLEAPIS_DIR)/.git:
	rm -rf "$(GOOGLEAPIS_DIR)"
	git clone --depth 1 --single-branch --branch $(GOOGLEAPIS_VERSION) \
	  https://github.com/googleapis/googleapis "$(GOOGLEAPIS_DIR)"

$(ALLORA_CHAIN_DIR)/.git:
	rm -rf "$(ALLORA_CHAIN_DIR)"
	git clone --depth 1 --single-branch --branch $(ALLORA_CHAIN_VERSION) \
	  https://github.com/allora-network/allora-chain "$(ALLORA_CHAIN_DIR)"

.PHONY: proto-deps
proto-deps: \
  $(GOGOPROTO_DIR)/.git \
  $(COSMOS_PROTO_DIR)/.git \
  $(COSMOS_SDK_DIR)/.git \
  $(GOOGLEAPIS_DIR)/.git \
  $(ALLORA_CHAIN_DIR)/.git

.PHONY: proto-deps-update
proto-deps-update:
	git -C "$(GOGOPROTO_DIR)" fetch --depth 1 origin $(GOGOPROTO_VERSION) && git -C "$(GOGOPROTO_DIR)" reset --hard FETCH_HEAD
	git -C "$(COSMOS_PROTO_DIR)" fetch --depth 1 origin $(COSMOS_PROTO_VERSION) && git -C "$(COSMOS_PROTO_DIR)" reset --hard FETCH_HEAD
	git -C "$(COSMOS_SDK_DIR)" fetch --depth 1 origin $(COSMOS_SDK_VERSION) && git -C "$(COSMOS_SDK_DIR)" reset --hard FETCH_HEAD
	git -C "$(GOOGLEAPIS_DIR)" fetch --depth 1 origin $(GOOGLEAPIS_VERSION) && git -C "$(GOOGLEAPIS_DIR)" reset --hard FETCH_HEAD
	git -C "$(ALLORA_CHAIN_DIR)" fetch --depth 1 origin $(ALLORA_CHAIN_VERSION) && git -C "$(ALLORA_CHAIN_DIR)" reset --hard FETCH_HEAD



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

