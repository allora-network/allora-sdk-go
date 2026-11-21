package mocks_test

import (
	allora "github.com/allora-network/allora-sdk-go"
	"github.com/allora-network/allora-sdk-go/mocks"
)

var (
	_ allora.Client    = (*mocks.MockClient)(nil)
	_ allora.APIClient = (*mocks.MockAPIClient)(nil)
)
