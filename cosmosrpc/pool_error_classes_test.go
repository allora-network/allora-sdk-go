package cosmosrpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/allora-network/allora-sdk-go/gen/rest"
)

// restEndpoint starts a REST endpoint that answers every request with the given
// status. It returns the endpoint URL and a counter of the requests it
// received.
func restEndpoint(t *testing.T, httpStatus int) (string, *atomic.Int64) {
	t.Helper()

	var hits atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		w.WriteHeader(httpStatus)
	}))
	t.Cleanup(server.Close)

	return server.URL, &hits
}

// TestPoolRESTStatusClasses pins, end to end over HTTP, which statuses end the
// request on the first endpoint and which fail over. The pool decides from the
// gRPC code the REST client derives from the status, so this test covers the
// two halves together.
func TestPoolRESTStatusClasses(t *testing.T) {
	for _, tc := range []struct {
		name string
		// httpStatus is what the first endpoint answers.
		httpStatus int
		// wantAttempts is how many endpoints the pool calls.
		wantAttempts int64
		// wantCooled is how many endpoints end up in the cooling pool.
		wantCooled int
	}{
		// Definitive answers: the node is healthy and no other node can
		// answer differently.
		{"501 unimplemented", http.StatusNotImplemented, 1, 0},
		// A provider's gateway refuses us; another endpoint may not.
		{"401 unauthenticated", http.StatusUnauthorized, 2, 0},
		{"403 permission denied", http.StatusForbidden, 2, 0},
		// The node or its proxy is rate limiting us.
		{"429 rate limited", http.StatusTooManyRequests, 2, 2},
		// The node is down or too slow.
		{"503 unavailable", http.StatusServiceUnavailable, 2, 2},
		{"504 gateway timeout", http.StatusGatewayTimeout, 2, 2},
		// Ambiguous statuses stay uncoded and count as endpoint faults.
		{"500 internal", http.StatusInternalServerError, 2, 2},
		{"502 bad gateway", http.StatusBadGateway, 2, 2},
		{"404 not found", http.StatusNotFound, 2, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := zerolog.Nop()
			urlA, hitsA := restEndpoint(t, tc.httpStatus)
			urlB, hitsB := restEndpoint(t, tc.httpStatus)

			p := NewClientPool([]Client{
				rest.NewRESTClient(urlA, logger),
				rest.NewRESTClient(urlB, logger),
			}, logger)
			t.Cleanup(func() { _ = p.Close() })

			_, err := p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
				&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})

			require.Error(t, err)
			require.Equal(t, tc.wantAttempts, hitsA.Load()+hitsB.Load(),
				"attempts for HTTP %d", tc.httpStatus)

			health := p.GetHealthStatus()
			require.Equal(t, tc.wantCooled, health["cooling_clients"],
				"cooling pool for HTTP %d", tc.httpStatus)
			require.Equal(t, 2-tc.wantCooled, health["active_clients"],
				"active pool for HTTP %d", tc.httpStatus)
		})
	}
}

// TestPoolRESTRateLimitFailsOverToHealthyEndpoint pins the useful half of the
// rate-limit rule: when one endpoint rate limits us and another serves us, the
// caller gets an answer.
func TestPoolRESTRateLimitFailsOverToHealthyEndpoint(t *testing.T) {
	logger := zerolog.Nop()

	limitedURL, limitedHits := restEndpoint(t, http.StatusTooManyRequests)

	var healthyHits atomic.Int64
	healthy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		healthyHits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"is_topic_worker_whitelist_enabled":true}`))
	}))
	t.Cleanup(healthy.Close)

	p := NewClientPool([]Client{
		rest.NewRESTClient(limitedURL, logger),
		rest.NewRESTClient(healthy.URL, logger),
	}, logger)
	t.Cleanup(func() { _ = p.Close() })

	resp, err := p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
		&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})

	require.NoError(t, err, "a rate-limited endpoint must not fail the request when another endpoint serves it")
	require.True(t, resp.IsTopicWorkerWhitelistEnabled)
	require.Equal(t, int64(1), limitedHits.Load())
	require.Equal(t, int64(1), healthyHits.Load())
}

// TestPoolRESTUnimplementedIsTerminal pins that a 501 ends the request on the
// first endpoint with the code the node meant, so the caller can tell "this
// node does not serve this method" from "this node is broken".
func TestPoolRESTUnimplementedIsTerminal(t *testing.T) {
	logger := zerolog.Nop()
	urlA, hitsA := restEndpoint(t, http.StatusNotImplemented)
	urlB, hitsB := restEndpoint(t, http.StatusNotImplemented)

	p := NewClientPool([]Client{
		rest.NewRESTClient(urlA, logger),
		rest.NewRESTClient(urlB, logger),
	}, logger)
	t.Cleanup(func() { _ = p.Close() })

	_, err := p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
		&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})

	require.Error(t, err)
	require.Equal(t, codes.Unimplemented, status.Code(err), "got %v", err)
	require.Equal(t, int64(1), hitsA.Load()+hitsB.Load(),
		"an unimplemented method must not be retried on the other endpoints")
	require.Equal(t, 2, p.GetHealthStatus()["active_clients"],
		"a 501 says nothing about endpoint health")
}
