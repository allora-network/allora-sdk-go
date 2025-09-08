#!/usr/bin/env bash

set -e

echo "🤖 Generating gogo proto code"
cd proto-deps
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    # Only generate for the specific modules we care about:
    # - Allora emissions v9
    # - Allora mint v5  
    # - Cosmos SDK auth, bank, staking
    # - Tendermint/CometBFT service
    if [[ $file == *"emissions/v9/"* ]] || \
       [[ $file == *"mint/v5/"* ]] || \
       [[ $file == *"cosmos/auth/"* ]] || \
       [[ $file == *"cosmos/bank/"* ]] || \
       [[ $file == *"cosmos/staking/"* ]] || \
       [[ $file == *"tendermint/types/"* ]] || \
       [[ $file == *"tendermint/version/"* ]] || \
       [[ $file == *"tendermint/p2p/"* ]] || \
       [[ $file == *"client/grpc/tmservice/"* ]]; then
      buf generate --template buf.gen.yaml $file
    fi
  done
done

cd ..

# Copy generated files to the right locations
echo "🤖 Moving generated files"

# Create directories
mkdir -p types/emissions/v9
mkdir -p types/mint/v5
mkdir -p types/cosmos/auth/v1beta1
mkdir -p types/cosmos/bank/v1beta1
mkdir -p types/cosmos/staking/v1beta1
mkdir -p types/tendermint/types
mkdir -p types/tendermint/version
mkdir -p types/tendermint/p2p
mkdir -p types/cosmos/base/tendermint/v1beta1

# Move generated files
find proto-deps -name "*.pb.go" -exec cp {} types/ \;
find proto-deps -name "*_grpc.pb.go" -exec cp {} types/ \;

echo "🤖 Proto code generation complete!"