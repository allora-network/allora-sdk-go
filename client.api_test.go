package allora

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T, gatewayURL, dataURL string, opts ...APIClientOption) *apiClient {
	t.Helper()
	base := []APIClientOption{
		WithBaseURL(gatewayURL),
		WithDataBaseURL(dataURL),
		WithNetwork(NetworkTestnet),
	}
	base = append(base, opts...)
	client, err := NewAPIClient("forge_sk_test", base...)
	require.NoError(t, err)
	return client
}

func TestNewAPIClientDefaultsAndOverrides(t *testing.T) {
	_, err := NewAPIClient("")
	require.Error(t, err)

	_, err = NewAPIClient("UP-deadbeef")
	require.Error(t, err)
	require.Contains(t, err.Error(), forgeAPIKeyPrefix)

	client, err := NewAPIClient("forge_sk_test")
	require.NoError(t, err)
	require.Equal(t, DefaultAPIBaseURL, client.url)
	require.Equal(t, DefaultDataAPIBaseURL, client.dataURL)
	require.Equal(t, NetworkMainnet, client.network)

	client, err = NewAPIClient("forge_sk_test",
		WithBaseURL("http://gateway.example"),
		WithDataBaseURL("http://data.example/api"),
		WithNetwork(NetworkTestnet),
	)
	require.NoError(t, err)
	require.Equal(t, NetworkTestnet, client.network)
	require.Equal(t, "http://gateway.example", client.url)
	require.Equal(t, "http://data.example/api", client.dataURL)

	_, err = NewAPIClient("forge_sk_test", WithNetwork(Network("bogus")))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported network")
}

func TestGetTopicsCompatShape(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/testnet/topics", r.URL.Path)
		require.Equal(t, "Bearer forge_sk_test", r.Header.Get("Authorization"))

		_ = json.NewEncoder(w).Encode(map[string]any{
			"network":   "testnet",
			"count":     1,
			"cached_at": "2026-07-14T10:40:00Z",
			"topics": []map[string]any{
				{
					"id":                       72,
					"creator":                  "allo1creator",
					"metadata":                 "BTC 5min",
					"loss_method":              "mse",
					"epoch_length":             720,
					"ground_truth_lag":         1,
					"worker_submission_window": 10,
					"created_at":               "2026-01-01T00:00:00Z",
				},
			},
		})
	}))
	defer srv.Close()

	client := testClient(t, srv.URL, "http://data.example/api")

	var got []*Topic
	for topic, err := range client.GetTopics() {
		require.NoError(t, err)
		got = append(got, topic)
	}
	require.Len(t, got, 1)
	require.Equal(t, 72, got[0].TopicID)
	require.Equal(t, "BTC 5min", got[0].TopicName)
	require.Equal(t, 720, got[0].EpochLength)
	require.False(t, got[0].IsActive) // not invented
	require.Zero(t, got[0].WorkerCount)
	require.Nil(t, got[0].LatestNetworkInference)
	require.Nil(t, got[0].Description)
}

func TestGetTopicAttachesCompatInference(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/testnet/topics":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"network": "testnet",
				"topics": []map[string]any{
					{"id": 72, "metadata": "BTC 5min", "epoch_length": 720, "loss_method": "mse"},
				},
			})
		case "/testnet/inference/72":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"topic_id": 72,
				"combined_value": []map[string]string{
					{"label_name": "y", "value": "0.00006678"},
				},
				"naive_value": []map[string]string{
					{"label_name": "y", "value": "0.00006"},
				},
				"timestamp_ms":                 1784025644244,
				"network_inference_normalized": "0.00006678",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	client := testClient(t, srv.URL, "http://data.example/api")

	topic, err := client.GetTopic(72)
	require.NoError(t, err)
	require.Equal(t, 72, topic.TopicID)
	require.NotNil(t, topic.LatestNetworkInference)
	require.Equal(t, "0.00006678", topic.LatestNetworkInference.CombinedValue)
	require.Equal(t, "0.00006", topic.LatestNetworkInference.NaiveValue)
	require.Equal(t, int64(1784025644244), topic.LatestNetworkInference.Timestamp)

	inf, err := client.GetInference(context.Background(), 72)
	require.NoError(t, err)
	require.Equal(t, "0.00006678", inf.CombinedValue)

	_, err = client.GetTopic(999)
	require.ErrorIs(t, err, ErrTopicNotFound)
}

func TestGetTopicBestEffortOnInferenceError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/testnet/topics":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"topics": []map[string]any{{"id": 72, "metadata": "BTC"}},
			})
		case "/testnet/inference/72":
			w.WriteHeader(http.StatusPaymentRequired)
			_, _ = w.Write([]byte(`{"error":"insufficient balance"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	client := testClient(t, srv.URL, "http://data.example/api")
	topic, err := client.GetTopic(72)
	require.NoError(t, err)
	require.Equal(t, 72, topic.TopicID)
	require.Equal(t, "BTC", topic.TopicName)
	require.Nil(t, topic.LatestNetworkInference)
}

func TestGetNetworkInferenceParsesFlatGatewayJSON(t *testing.T) {
	var sawQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/testnet/inference/72", r.URL.Path)
		sawQuery = r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(map[string]any{
			"network":  "testnet",
			"topic_id": 72,
			"combined_value": []map[string]string{
				{"label_name": "y", "value": "0.00006678"},
			},
			"timestamp_ms":                 1784025644244,
			"network_inference":            "66781380725120668592",
			"network_inference_normalized": "0.00006678",
			"signature":                    "0xabc",
			"evm_chain_id":                 8453,
			"consumer_address":             "0xconsumer",
			"vector_signature":             "0xvec",
			"cached_at":                    "2026-07-14T10:40:45Z",
			"is_stale":                     false,
		})
	}))
	defer srv.Close()

	client := testClient(t, srv.URL, "http://data.example/api")

	ni, err := client.GetNetworkInference(context.Background(), 72, WithEVMChainID(8453))
	require.NoError(t, err)
	require.Equal(t, "evm_chain_id=8453", sawQuery)
	require.Equal(t, int64(72), ni.TopicID)
	require.Equal(t, "0xabc", ni.Signature)
	require.Equal(t, uint64(8453), ni.EVMChainID)
	require.WithinDuration(t, time.Date(2026, 7, 14, 10, 40, 45, 0, time.UTC), ni.CachedAt.UTC(), time.Second)
}

func TestGetInferenceMapsHTTPErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
		_, _ = w.Write([]byte(`{"error":"insufficient balance"}`))
	}))
	defer srv.Close()

	client := testClient(t, srv.URL, "http://data.example/api")

	_, err := client.GetInference(context.Background(), 1)
	require.Error(t, err)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusPaymentRequired, apiErr.StatusCode)
	require.True(t, strings.Contains(apiErr.Error(), "insufficient balance"))
}

func TestGetOHLCDataCompatShape(t *testing.T) {
	var rowQueries []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "forge_sk_test", r.Header.Get("X-API-Key")) {
			http.Error(w, "bad key", http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/api/datasets/search/":
			if !assert.Equal(t, "btcusd", r.URL.Query().Get("ticker")) {
				http.Error(w, "bad ticker", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"count": 1,
				"results": []map[string]any{
					{"id": 3, "name": "tiingo_btcusd_1min", "metadata": map[string]any{"ticker": "btcusd"},
						"created_at": "2026-09-30T13:15:49Z", "updated_at": "2026-09-30T14:28:19Z"},
				},
			})
		case "/api/rows/":
			rowQueries = append(rowQueries, r.URL.RawQuery)
			if r.URL.Query().Get("offset") == "" {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"count": 2,
					"next":  "/api/rows/?dataset_name=tiingo_btcusd_1min&limit=1&offset=1",
					"results": []map[string]any{
						{
							"dataset_id": 3,
							"timestamp":  "2024-01-01T00:00:00Z",
							"values": map[string]any{
								"open": 100.5, "high": 110.0, "low": 90.25, "close": 105.0,
								"volume": 12.5, "volume_notional": 1300.0, "trades_done": 42,
							},
						},
					},
				})
				return
			}
			if !assert.Equal(t, "1", r.URL.Query().Get("offset")) {
				http.Error(w, "expected offset=1", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"count": 2,
				"next":  nil,
				"results": []map[string]any{
					{
						"dataset_id": 3,
						"timestamp":  "2024-01-01T00:01:00Z",
						"values": map[string]any{
							"open": 105.0, "high": 106.0, "low": 104.0, "close": 105.5,
							"volume": 1.0, "volume_notional": 105.5, "trades_done": 3,
						},
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	client := testClient(t, "http://gateway.example", srv.URL+"/api")

	var pages []*OHLCResponse
	for resp, err := range client.GetOHLCData(context.Background(), "BTCUSD", "2024-01-01") {
		require.NoError(t, err)
		require.NotNil(t, resp)
		pages = append(pages, resp)
	}
	require.Len(t, pages, 2)
	require.NotNil(t, pages[0].ContinuationToken)
	require.Nil(t, pages[1].ContinuationToken)
	require.Len(t, rowQueries, 2)
	require.NotContains(t, rowQueries[0], "offset=")
	require.Contains(t, rowQueries[1], "offset=1")

	require.Equal(t, "BTCUSD", pages[0].Data[0].Ticker)
	require.Equal(t, "2024-01-01T00:00:00Z", pages[0].Data[0].Date)
	require.Equal(t, "100.5", pages[0].Data[0].Open)
	require.Equal(t, "105", pages[0].Data[0].Close)
	require.Equal(t, 42, pages[0].Data[0].TradesDone)
}
