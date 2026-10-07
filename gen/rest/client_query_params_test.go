package rest

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// A pagination key is raw bytes, often not valid UTF-8. The node's REST
// gateway reads a bytes query parameter as base64, so the key must arrive
// base64-encoded for the next page to start where the last one ended.
func TestRESTClientSendsBytesQueryParamsAsBase64(t *testing.T) {
	key := []byte{0x14, 0x00, 0x52, 0x24, 0x68, 0xf3, 0xd9, 0xa8, 0xff}

	var got url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Query()
		_, _ = w.Write([]byte("{}"))
	}))
	defer server.Close()

	client := NewBankRESTClient(NewRESTClientCore(server.URL, zerolog.Nop()), zerolog.Nop())
	_, err := client.DenomOwners(context.Background(), &banktypes.QueryDenomOwnersRequest{
		Denom:      "uallo",
		Pagination: &query.PageRequest{Key: key, Limit: 1000},
	})
	require.NoError(t, err)

	decoded, err := base64.StdEncoding.DecodeString(got.Get("pagination.key"))
	require.NoError(t, err)
	require.Equal(t, key, decoded)
	require.Equal(t, "1000", got.Get("pagination.limit"))
}
