//go:build integration

package allora

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Run with:
//
//	ALLORA_API_KEY=forge_sk_… go test -tags=integration -run TestIntegrationForgeData -count=1 -v .
//
// Optional overrides (defaults: DefaultAPIBaseURL / DefaultDataAPIBaseURL):
//
//	ALLORA_API_BASE_URL=http://localhost:8080 \
//	ALLORA_DATA_API_BASE_URL=https://gcp-forge-data-service-stage.tail110056.ts.net/api \
//	ALLORA_API_KEY=forge_sk_… \
//	go test -tags=integration -run TestIntegrationForgeData -count=1 -v .
func TestIntegrationForgeDataCandles(t *testing.T) {
	apiKey := os.Getenv("ALLORA_API_KEY")
	if apiKey == "" {
		t.Skip("ALLORA_API_KEY not set")
	}

	var opts []APIClientOption
	if gatewayURL := os.Getenv("ALLORA_API_BASE_URL"); gatewayURL != "" {
		opts = append(opts, WithBaseURL(gatewayURL))
	}
	if dataURL := os.Getenv("ALLORA_DATA_API_BASE_URL"); dataURL != "" {
		opts = append(opts, WithDataBaseURL(dataURL))
	}

	client, err := NewAPIClient(apiKey, opts...)
	require.NoError(t, err)
	t.Logf("gateway=%s data=%s", client.url, client.dataURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	page, err := client.ListDatasets(ctx, WithDatasetLimit(5))
	require.NoError(t, err)
	require.Greater(t, page.Count, 0)
	require.NotEmpty(t, page.Results)
	t.Logf("datasets: count=%d page=%d first=%q", page.Count, len(page.Results), page.Results[0].Name)

	datasets, err := client.FindDatasetsByTicker(ctx, "btcusd")
	require.NoError(t, err)
	require.NotEmpty(t, datasets)
	t.Logf("btcusd datasets: %d first=%q", len(datasets), datasets[0].Name)

	candles, err := client.GetCandles(ctx, CandleQuery{
		Ticker: "btcusd",
		Start:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
		Limit:  3,
	})
	require.NoError(t, err)
	require.NotEmpty(t, candles.Results)
	require.Greater(t, candles.Results[0].Close, 0.0)
	c0 := candles.Results[0]
	t.Logf("candles: n=%d first={ts=%s o=%.2f h=%.2f l=%.2f c=%.2f v=%.4f}",
		len(candles.Results), c0.Timestamp.UTC().Format(time.RFC3339),
		c0.Open, c0.High, c0.Low, c0.Close, c0.Volume)

	var iterPages int
	for batch, err := range client.IterCandles(ctx, CandleQuery{
		Ticker: "btcusd",
		Start:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Limit:  250,
	}) {
		require.NoError(t, err)
		require.NotEmpty(t, batch)
		require.Greater(t, batch[0].Close, 0.0)
		iterPages++
		d0 := batch[0]
		t.Logf("iter candles: n=%d first={ts=%s o=%.2f h=%.2f l=%.2f c=%.2f v=%.4f}",
			len(batch), d0.Timestamp.UTC().Format(time.RFC3339),
			d0.Open, d0.High, d0.Low, d0.Close, d0.Volume)
		break
	}
	require.Equal(t, 1, iterPages)

	hourly, err := client.GetCandles(ctx, CandleQuery{
		DatasetName: datasets[0].Name,
		Start:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		End:         time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		Resolution:  CandleResolution1h,
		Limit:       2,
	})
	require.NoError(t, err)
	require.NotEmpty(t, hourly.Results)
	h0 := hourly.Results[0]
	t.Logf("hourly candles: n=%d dataset=%q first={ts=%s c=%.2f}",
		len(hourly.Results), datasets[0].Name, h0.Timestamp.UTC().Format(time.RFC3339), h0.Close)
}
