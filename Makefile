#!/usr/bin/make -f

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
	git clone --depth 1 --single-branch --branch v1.7.0 \
	  https://github.com/cosmos/gogoproto "$(GOGOPROTO_DIR)"

$(COSMOS_PROTO_DIR)/.git:
	rm -rf "$(COSMOS_PROTO_DIR)"
	git clone --depth 1 --single-branch --branch v1.0.0-beta.5 \
	  https://github.com/cosmos/cosmos-proto "$(COSMOS_PROTO_DIR)"

$(COSMOS_SDK_DIR)/.git:
	rm -rf "$(COSMOS_SDK_DIR)"
	git clone --depth 1 --single-branch --branch v0.50.13 \
	  https://github.com/cosmos/cosmos-sdk "$(COSMOS_SDK_DIR)"

$(GOOGLEAPIS_DIR)/.git:
	rm -rf "$(GOOGLEAPIS_DIR)"
	git clone --depth 1 --single-branch --branch master \
	  https://github.com/googleapis/googleapis "$(GOOGLEAPIS_DIR)"

$(ALLORA_CHAIN_DIR)/.git:
	rm -rf "$(ALLORA_CHAIN_DIR)"
	git clone --depth 1 --single-branch --branch v0.12.1 \
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
	git -C "$(GOGOPROTO_DIR)" fetch --depth 1 origin v1.7.0 && git -C "$(GOGOPROTO_DIR)" reset --hard FETCH_HEAD
	git -C "$(COSMOS_PROTO_DIR)" fetch --depth 1 origin v1.0.0-beta.5 && git -C "$(COSMOS_PROTO_DIR)" reset --hard FETCH_HEAD
	git -C "$(COSMOS_SDK_DIR)" fetch --depth 1 origin v0.50.13 && git -C "$(COSMOS_SDK_DIR)" reset --hard FETCH_HEAD
	git -C "$(GOOGLEAPIS_DIR)" fetch --depth 1 origin master && git -C "$(GOOGLEAPIS_DIR)" reset --hard FETCH_HEAD
	git -C "$(ALLORA_CHAIN_DIR)" fetch --depth 1 origin v0.12.1 && git -C "$(ALLORA_CHAIN_DIR)" reset --hard FETCH_HEAD



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

