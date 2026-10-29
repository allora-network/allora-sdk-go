package allora

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListDatasetsUsesDataAPIKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "/api/datasets/", r.URL.Path) ||
			!assert.Equal(t, "forge_sk_test", r.Header.Get("X-API-Key")) ||
			!assert.Empty(t, r.Header.Get("Authorization")) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"count":    1,
			"next":     nil,
			"previous": nil,
			"results": []map[string]any{
				{
					"id":           3,
					"name":         "tiingo_btcusd_1min",
					"description":  "BTC",
					"metadata":     map[string]any{"ticker": "btcusd", "frequency": "1min"},
					"owner_key_id": 2,
					"created_at":   "2026-09-30T13:15:49.236047Z",
					"updated_at":   "2026-09-30T14:28:19.955361Z",
				},
			},
		})
	}))
	defer srv.Close()

	client := testClient(t, "http://gateway.example", srv.URL+"/api")

	page, err := client.ListDatasets(context.Background(), WithDatasetLimit(10))
	require.NoError(t, err)
	require.Equal(t, 1, page.Count)
	require.Len(t, page.Results, 1)
	require.Equal(t, "tiingo_btcusd_1min", page.Results[0].Name)
	require.Equal(t, "btcusd", page.Results[0].Metadata["ticker"])
}

func TestFindDatasetsByTicker(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "/api/datasets/search/", r.URL.Path) ||
			!assert.Equal(t, "btcusd", r.URL.Query().Get("ticker")) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"count": 1,
			"results": []map[string]any{
				{"id": 3, "name": "tiingo_btcusd_1min", "metadata": map[string]any{"ticker": "btcusd"},
					"created_at": "2026-09-30T13:15:49Z", "updated_at": "2026-09-30T14:28:19Z"},
			},
		})
	}))
	defer srv.Close()

	client := testClient(t, "http://gateway.example", srv.URL+"/api")

	datasets, err := client.FindDatasetsByTicker(context.Background(), "BTCUSD")
	require.NoError(t, err)
	require.Len(t, datasets, 1)
	require.Equal(t, 3, datasets[0].ID)
}

func TestGetCandlesParsesRawAndAggregate(t *testing.T) {
	var paths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path+"?"+r.URL.RawQuery)
		if !assert.Equal(t, "forge_sk_test", r.Header.Get("X-API-Key")) {
			http.Error(w, "bad key", http.StatusUnauthorized)
			return
		}

		switch {
		case r.URL.Path == "/api/datasets/search/":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"count": 1,
				"results": []map[string]any{
					{"id": 3, "name": "tiingo_btcusd_1min", "metadata": map[string]any{"ticker": "btcusd"},
						"created_at": "2026-09-30T13:15:49Z", "updated_at": "2026-09-30T14:28:19Z"},
				},
			})
		case r.URL.Query().Get("resolution") == "1h":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"count": 1,
				"results": []map[string]any{
					{
						"dataset_id": 3,
						"timestamp":  "2024-01-01T00:00:00Z",
						"open":       100.0,
						"high":       110.0,
						"low":        90.0,
						"close":      105.0,
						"volume":     50.0,
						"count":      60,
					},
				},
			})
		default:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"count": 1,
				"next":  "/api/rows/?dataset_name=tiingo_btcusd_1min&limit=1&offset=1",
				"results": []map[string]any{
					{
						"id":         1,
						"dataset_id": 3,
						"timestamp":  "2024-01-01T00:00:00Z",
						"values": map[string]any{
							"open": 42286.5, "high": 42345.2, "low": 42260.9, "close": 42298.4,
							"volume": 58.6, "volume_notional": 2482757.8, "trades_done": 633,
						},
					},
				},
			})
		}
	}))
	defer srv.Close()

	client := testClient(t, "http://gateway.example", srv.URL+"/api")

	raw, err := client.GetCandles(context.Background(), CandleQuery{
		Ticker: "btcusd",
		Start:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
		Limit:  1,
	})
	require.NoError(t, err)
	require.Len(t, raw.Results, 1)
	require.InDelta(t, 42298.4, raw.Results[0].Close, 0.01)
	require.Equal(t, 633, raw.Results[0].TradesDone)
	require.NotNil(t, raw.Next)

	agg, err := client.GetCandles(context.Background(), CandleQuery{
		DatasetName: "tiingo_btcusd_1min",
		Resolution:  CandleResolution1h,
		Start:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		End:         time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		Limit:       1,
	})
	require.NoError(t, err)
	require.Len(t, agg.Results, 1)
	require.Equal(t, 105.0, agg.Results[0].Close)
	require.Equal(t, 60, agg.Results[0].Count)

	require.NotEmpty(t, paths)
}

func TestGetCandlesRequiresExactlyOneSelector(t *testing.T) {
	client := testClient(t, "http://gateway.example", "http://example.invalid/api")

	_, err := client.GetCandles(context.Background(), CandleQuery{})
	require.Error(t, err)

	_, err = client.GetCandles(context.Background(), CandleQuery{
		DatasetID: 1, DatasetName: "x",
	})
	require.Error(t, err)
}

func TestGetCandlesRejectsAmbiguousTicker(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"count": 2,
			"results": []map[string]any{
				{"id": 1, "name": "tiingo_btcusd_5min", "metadata": map[string]any{"ticker": "btcusd"},
					"created_at": "2026-09-30T13:15:49Z", "updated_at": "2026-09-30T14:28:19Z"},
				{"id": 2, "name": "tiingo_btcusd_15min", "metadata": map[string]any{"ticker": "btcusd"},
					"created_at": "2026-09-30T13:15:49Z", "updated_at": "2026-09-30T14:28:19Z"},
			},
		})
	}))
	defer srv.Close()

	client := testClient(t, "http://gateway.example", srv.URL+"/api")
	_, err := client.GetCandles(context.Background(), CandleQuery{Ticker: "btcusd"})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAmbiguousTicker)
	require.Contains(t, err.Error(), "DatasetID or DatasetName")
}

func TestIterCandlesFollowsNext(t *testing.T) {
	var queries []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "forge_sk_test", r.Header.Get("X-API-Key")) {
			http.Error(w, "bad key", http.StatusUnauthorized)
			return
		}
		queries = append(queries, r.URL.RawQuery)
		if r.URL.Query().Get("offset") == "" {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"count": 2,
				"next":  "/api/rows/?dataset_name=tiingo_btcusd_1min&limit=1&offset=1",
				"results": []map[string]any{
					{"dataset_id": 3, "timestamp": "2024-01-01T00:00:00Z",
						"values": map[string]any{"open": 1, "high": 2, "low": 0.5, "close": 1.5, "volume": 10}},
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
				{"dataset_id": 3, "timestamp": "2024-01-01T00:01:00Z",
					"values": map[string]any{"open": 1.5, "high": 3, "low": 1, "close": 2.5, "volume": 11}},
			},
		})
	}))
	defer srv.Close()

	client := testClient(t, "http://gateway.example", srv.URL+"/api")

	var closes []float64
	for candles, err := range client.IterCandles(context.Background(), CandleQuery{
		DatasetName: "tiingo_btcusd_1min",
		Limit:       1,
	}) {
		require.NoError(t, err)
		for _, c := range candles {
			closes = append(closes, c.Close)
		}
	}
	require.Equal(t, []float64{1.5, 2.5}, closes)
	require.Len(t, queries, 2)
	require.NotContains(t, queries[0], "offset=")
	require.Contains(t, queries[1], "offset=1")
}

func TestIterCandlesDetectsNextCycle(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		// Point next back at the same query the client just issued (no offset).
		next := "/api/rows/?dataset_name=tiingo_btcusd_1min&limit=1"
		_ = json.NewEncoder(w).Encode(map[string]any{
			"count": 1,
			"next":  next,
			"results": []map[string]any{
				{"dataset_id": 3, "timestamp": "2024-01-01T00:00:00Z",
					"values": map[string]any{"open": 1, "high": 2, "low": 0.5, "close": 1.5, "volume": 10}},
			},
		})
	}))
	defer srv.Close()

	client := testClient(t, "http://gateway.example", srv.URL+"/api")
	var sawErr error
	var pages int
	for _, err := range client.IterCandles(context.Background(), CandleQuery{
		DatasetName: "tiingo_btcusd_1min",
		Limit:       1,
	}) {
		if err != nil {
			sawErr = err
			break
		}
		pages++
	}
	require.Error(t, sawErr)
	require.Contains(t, sawErr.Error(), "cycle detected")
	require.Equal(t, 1, pages, "must not re-yield the initial page when next cycles back to it")
	require.Equal(t, 1, hits)
}

func TestIterCandlesDetectsReorderedQueryCycle(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		// Same logical query as the initial request but with reordered params.
		next := "/api/rows/?limit=1&dataset_name=tiingo_btcusd_1min"
		_ = json.NewEncoder(w).Encode(map[string]any{
			"count": 1,
			"next":  next,
			"results": []map[string]any{
				{"dataset_id": 3, "timestamp": "2024-01-01T00:00:00Z",
					"values": map[string]any{"open": 1, "high": 2, "low": 0.5, "close": 1.5, "volume": 10}},
			},
		})
	}))
	defer srv.Close()

	client := testClient(t, "http://gateway.example", srv.URL+"/api")
	var sawErr error
	for _, err := range client.IterCandles(context.Background(), CandleQuery{
		DatasetName: "tiingo_btcusd_1min",
		Limit:       1,
	}) {
		if err != nil {
			sawErr = err
			break
		}
	}
	require.Error(t, sawErr)
	require.Contains(t, sawErr.Error(), "cycle detected")
	require.Equal(t, 1, hits, "reordered next must be caught on the first hop")
}

func TestResolveDataNextURLRejectsForeignHost(t *testing.T) {
	client := testClient(t, "http://gateway.example", "https://forge-data.allora.network/api")
	_, err := client.resolveDataNextURL("https://evil.example/api/rows/?offset=1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "does not match data API")
}

func TestResolveDataNextURLKeepsAPIPrefix(t *testing.T) {
	client := testClient(t, "http://gateway.example", "https://forge-data.allora.network/api")
	got, err := client.resolveDataNextURL("rows/?offset=1")
	require.NoError(t, err)
	require.Equal(t, "https://forge-data.allora.network/api/rows/?offset=1", got)
}
