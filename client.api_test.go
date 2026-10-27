package allora

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brynbellomy/go-utils/errors"
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

	// Prefix is no longer validated client-side; server decides.
	client, err := NewAPIClient("UP-deadbeef")
	require.NoError(t, err)
	require.Equal(t, "UP-deadbeef", client.apiKey)

	client, err = NewAPIClient("forge_sk_test")
	require.NoError(t, err)
	require.Equal(t, DefaultAPIBaseURL, client.url)
	require.Equal(t, DefaultDataAPIBaseURL, client.dataURL)
	require.Equal(t, "https://forge-data.allora.run/api", DefaultDataAPIBaseURL)
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

func TestWithHTTPClientAndTimeout(t *testing.T) {
	hc := &http.Client{Timeout: 5 * time.Second}
	client, err := NewAPIClient("k", WithHTTPClient(hc), WithTimeout(2*time.Second))
	require.NoError(t, err)
	require.NotSame(t, hc, client.httpClient, "WithHTTPClient must copy so the caller's client is not mutated")
	require.Equal(t, 2*time.Second, client.httpClient.Timeout)
	require.Equal(t, 5*time.Second, hc.Timeout, "caller Timeout must stay untouched")
	require.NotNil(t, client.httpClient.CheckRedirect)
	require.Nil(t, hc.CheckRedirect, "caller CheckRedirect must stay untouched")
}

func TestGetTopicsNativeShape(t *testing.T) {
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
	require.Equal(t, int64(72), got[0].ID)
	require.Equal(t, "BTC 5min", got[0].Metadata)
	require.Equal(t, int64(720), got[0].EpochLength)
	require.Equal(t, "allo1creator", got[0].Creator)
	require.Nil(t, got[0].LatestNetworkInference)
}

func TestGetTopicAttachesNetworkInference(t *testing.T) {
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
	require.Equal(t, int64(72), topic.ID)
	require.NotNil(t, topic.LatestNetworkInference)
	require.Equal(t, "0.00006678", topic.LatestNetworkInference.NetworkInferenceNormalized)
	require.Equal(t, int64(1784025644244), topic.LatestNetworkInference.TimestampMs)
	require.Equal(t, []LabeledValue{{LabelName: "y", Value: "0.00006678"}}, topic.LatestNetworkInference.CombinedValue)
	require.Equal(t, []LabeledValue{{LabelName: "y", Value: "0.00006"}}, topic.LatestNetworkInference.NaiveValue)

	inf, err := client.GetNetworkInference(context.Background(), 72)
	require.NoError(t, err)
	require.Equal(t, "0.00006678", inf.NetworkInferenceNormalized)

	_, err = client.GetTopic(999)
	require.ErrorIs(t, err, ErrTopicNotFound)
}

func TestGetTopicBestEffortOnInferenceError(t *testing.T) {
	var inferenceHit bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/testnet/topics":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"topics": []map[string]any{{"id": 72, "metadata": "BTC"}},
			})
		case "/testnet/inference/72":
			inferenceHit = true
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
	require.True(t, inferenceHit, "GetTopic must attempt inference before tolerating failure")
	require.Equal(t, int64(72), topic.ID)
	require.Equal(t, "BTC", topic.Metadata)
	require.Nil(t, topic.LatestNetworkInference)
}

func TestDataAPIRefusesCrossHostRedirect(t *testing.T) {
	var foreignHit bool
	foreign := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		foreignHit = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"count":0,"results":[]}`))
	}))
	defer foreign.Close()

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "forge_sk_test", r.Header.Get("X-API-Key"))
		http.Redirect(w, r, foreign.URL+"/api/datasets/", http.StatusFound)
	}))
	defer origin.Close()

	client := testClient(t, "http://gateway.example", origin.URL+"/api")
	_, err := client.ListDatasets(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "cross-host redirect")
	require.False(t, foreignHit, "X-API-Key must not be forwarded to a foreign host")
}

func TestDataAPIAllowsSameHostRedirect(t *testing.T) {
	var finalHit bool
	var srvURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery == "" {
			http.Redirect(w, r, srvURL+"/api/datasets/?redirected=1", http.StatusFound)
			return
		}
		finalHit = true
		require.Equal(t, "forge_sk_test", r.Header.Get("X-API-Key"))
		_ = json.NewEncoder(w).Encode(map[string]any{"count": 0, "results": []any{}})
	}))
	defer srv.Close()
	srvURL = srv.URL

	client := testClient(t, "http://gateway.example", srv.URL+"/api")
	page, err := client.ListDatasets(context.Background())
	require.NoError(t, err)
	require.True(t, finalHit)
	require.Equal(t, 0, page.Count)
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

func TestGetNetworkInferenceMapsHTTPErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
		_, _ = w.Write([]byte(`{"error":"insufficient balance"}`))
	}))
	defer srv.Close()

	client := testClient(t, srv.URL, "http://data.example/api")

	_, err := client.GetNetworkInference(context.Background(), 1)
	require.Error(t, err)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusPaymentRequired, apiErr.StatusCode)
	require.True(t, strings.Contains(apiErr.Error(), "insufficient balance"))
	require.Equal(t, errors.FaultCaller, errors.GetFault(err))
	require.False(t, errors.IsRetryable(err))
}

func TestIterCandlesNativeShape(t *testing.T) {
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

	var pages [][]Candle
	for candles, err := range client.IterCandles(context.Background(), CandleQuery{
		Ticker:   "btcusd",
		Start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Ordering: "timestamp",
		Limit:    1,
	}) {
		require.NoError(t, err)
		pages = append(pages, candles)
	}
	require.Len(t, pages, 2)
	require.Len(t, rowQueries, 2)
	require.NotContains(t, rowQueries[0], "offset=")
	require.Contains(t, rowQueries[1], "offset=1")

	require.Equal(t, 3, pages[0][0].DatasetID)
	require.Equal(t, 100.5, pages[0][0].Open)
	require.Equal(t, 105.0, pages[0][0].Close)
	require.Equal(t, 42, pages[0][0].TradesDone)
}
