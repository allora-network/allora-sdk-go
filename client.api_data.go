package allora

import (
	"context"
	"encoding/json"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/brynbellomy/go-utils/errors"
)

// CandleResolution selects Atlas row aggregation.
type CandleResolution string

const (
	CandleResolutionRaw CandleResolution = "raw"
	CandleResolution5m  CandleResolution = "5m"
	CandleResolution1h  CandleResolution = "1h"
	CandleResolution1d  CandleResolution = "1d"
)

// Dataset is a Forge Data (Atlas) dataset catalog entry.
type Dataset struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Metadata    map[string]any `json:"metadata"`
	OwnerKeyID  *int           `json:"owner_key_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// DatasetPage is one page of GET /api/datasets/ (or search).
type DatasetPage struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []Dataset `json:"results"`
}

type datasetListOptions struct {
	limit  int
	offset int
	search string
}

// DatasetListOption configures ListDatasets.
type DatasetListOption func(*datasetListOptions)

// WithDatasetLimit sets page size for ListDatasets.
func WithDatasetLimit(limit int) DatasetListOption {
	return func(o *datasetListOptions) { o.limit = limit }
}

// WithDatasetOffset sets the offset for ListDatasets.
func WithDatasetOffset(offset int) DatasetListOption {
	return func(o *datasetListOptions) { o.offset = offset }
}

// WithDatasetSearch sets the text search for ListDatasets.
func WithDatasetSearch(search string) DatasetListOption {
	return func(o *datasetListOptions) { o.search = strings.TrimSpace(search) }
}

// CandleQuery selects OHLCV rows from Forge Data.
// Provide exactly one of DatasetID, DatasetName, or Ticker.
type CandleQuery struct {
	DatasetID   int
	DatasetName string
	Ticker      string

	Start time.Time
	End   time.Time

	Resolution CandleResolution
	Limit      int
	Offset     int
	Ordering   string

	// MaxPages caps how many Atlas pages a single iteration will follow.
	// Zero means DefaultMaxCandlePages. Exceeding the cap returns an error.
	MaxPages int
}

// Candle is one OHLCV bar from Atlas.
type Candle struct {
	DatasetID      int       `json:"dataset_id"`
	Timestamp      time.Time `json:"timestamp"`
	Open           float64   `json:"open"`
	High           float64   `json:"high"`
	Low            float64   `json:"low"`
	Close          float64   `json:"close"`
	Volume         float64   `json:"volume"`
	VolumeNotional float64   `json:"volume_notional,omitempty"`
	TradesDone     int       `json:"trades_done,omitempty"`
	Count          int       `json:"count,omitempty"`
}

// CandlePage is one page of GET /api/rows/.
type CandlePage struct {
	Count    int      `json:"count"`
	Next     *string  `json:"next"`
	Previous *string  `json:"previous"`
	Results  []Candle `json:"results"`
}

type atlasDatasetWire struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Metadata    json.RawMessage `json:"metadata"`
	OwnerKeyID  *int            `json:"owner_key_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type atlasDatasetPageWire struct {
	Count    int                `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []atlasDatasetWire `json:"results"`
}

type atlasRowWire struct {
	ID        int64           `json:"id"`
	DatasetID int             `json:"dataset_id"`
	Timestamp time.Time       `json:"timestamp"`
	Values    json.RawMessage `json:"values"`
	CreatedAt time.Time       `json:"created_at"`
	Open      *float64        `json:"open"`
	High      *float64        `json:"high"`
	Low       *float64        `json:"low"`
	Close     *float64        `json:"close"`
	Volume    *float64        `json:"volume"`
	Count     *int            `json:"count"`
}

type atlasCandleValues struct {
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Close          float64 `json:"close"`
	Volume         float64 `json:"volume"`
	VolumeNotional float64 `json:"volume_notional"`
	TradesDone     int     `json:"trades_done"`
}

type atlasRowsPageWire struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []atlasRowWire `json:"results"`
}

// ListDatasets calls GET /api/datasets/.
func (c *apiClient) ListDatasets(ctx context.Context, opts ...DatasetListOption) (*DatasetPage, error) {
	options := &datasetListOptions{}
	for _, opt := range opts {
		opt(options)
	}

	query := map[string]string{}
	if options.limit > 0 {
		query["limit"] = strconv.Itoa(options.limit)
	}
	if options.offset > 0 {
		query["offset"] = strconv.Itoa(options.offset)
	}
	if options.search != "" {
		query["search"] = options.search
	}

	reqURL, err := c.dataURLPath("datasets/", query)
	if err != nil {
		return nil, errors.WithMessage(err, "could not parse url")
	}

	var wire atlasDatasetPageWire
	if err := c.doDataJSON(ctx, "ListDatasets", http.MethodGet, reqURL, &wire); err != nil {
		return nil, err
	}
	return decodeDatasetPage(wire)
}

// FindDatasetsByTicker calls GET /api/datasets/search/?ticker=….
func (c *apiClient) FindDatasetsByTicker(ctx context.Context, ticker string) ([]Dataset, error) {
	ticker = strings.ToLower(strings.TrimSpace(ticker))
	if ticker == "" {
		return nil, errors.New("ticker is required")
	}

	reqURL, err := c.dataURLPath("datasets/search/", map[string]string{"ticker": ticker})
	if err != nil {
		return nil, errors.WithMessage(err, "could not parse url")
	}

	var wire atlasDatasetPageWire
	if err := c.doDataJSON(ctx, "FindDatasetsByTicker", http.MethodGet, reqURL, &wire); err != nil {
		return nil, err
	}
	page, err := decodeDatasetPage(wire)
	if err != nil {
		return nil, err
	}
	return page.Results, nil
}

// DefaultMaxCandlePages caps Atlas pagination for IterCandles / GetCandles.
// At Limit=250 this is ~2.5M candles; raise via CandleQuery.MaxPages for longer pulls.
// Exceeding the cap returns an error (tail is not silently truncated).
const DefaultMaxCandlePages = 10_000

// IterCandles walks every page of GetCandles, following Atlas `next` links.
// Empty pages are skipped; an all-empty result set ends the iterator without error.
func (c *apiClient) IterCandles(ctx context.Context, q CandleQuery) iter.Seq2[[]Candle, error] {
	return func(yield func([]Candle, error) bool) {
		for page, err := range c.iterCandlePages(ctx, q) {
			if err != nil {
				yield(nil, err)
				return
			}
			if len(page.Results) == 0 {
				continue
			}
			if !yield(page.Results, nil) {
				return
			}
		}
	}
}

// iterCandlePages is the single pagination loop for Atlas candle queries.
// It caps pages and aborts if a resolved next URL repeats (cycle guard),
// including cycles that return to the initial request URL. Cycle keys are
// canonicalized so reordered/re-encoded query parameters are still detected.
func (c *apiClient) iterCandlePages(ctx context.Context, q CandleQuery) iter.Seq2[*CandlePage, error] {
	return func(yield func(*CandlePage, error) bool) {
		pageCap := q.MaxPages
		if pageCap <= 0 {
			pageCap = DefaultMaxCandlePages
		}
		page, reqURL, err := c.fetchCandlePage(ctx, q)
		seen := map[string]struct{}{}
		if reqURL != "" {
			if k, kerr := canonicalNextKey(reqURL); kerr == nil {
				seen[k] = struct{}{}
			}
		}
		pages := 0
		for {
			if err != nil {
				yield(nil, err)
				return
			}
			if page == nil {
				return
			}
			pages++
			if pages > pageCap {
				yield(nil, errors.Errorf("candle pagination exceeded %d pages (set CandleQuery.MaxPages to raise; default %d)", pageCap, DefaultMaxCandlePages))
				return
			}
			if !yield(page, nil) {
				return
			}
			if page.Next == nil || *page.Next == "" {
				return
			}
			nextURL, resolveErr := c.resolveDataNextURL(*page.Next)
			if resolveErr != nil {
				yield(nil, errors.WithMessage(resolveErr, "could not resolve next url"))
				return
			}
			key, kerr := canonicalNextKey(nextURL)
			if kerr != nil {
				key = nextURL
			}
			if _, dup := seen[key]; dup {
				yield(nil, errors.Errorf("candle pagination cycle detected at %s", nextURL))
				return
			}
			seen[key] = struct{}{}
			page, err = c.getCandlesNext(ctx, *page.Next)
		}
	}
}

// canonicalNextKey normalizes a URL so that two logically identical pages with
// reordered or differently-encoded query parameters compare equal.
func canonicalNextKey(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	u.Fragment = ""
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}
	u.RawQuery = query.Encode() // sorted keys, canonical escaping
	return u.String(), nil
}

// GetCandles calls GET /api/rows/ for one page of OHLCV data.
func (c *apiClient) GetCandles(ctx context.Context, q CandleQuery) (*CandlePage, error) {
	page, _, err := c.fetchCandlePage(ctx, q)
	return page, err
}

func (c *apiClient) fetchCandlePage(ctx context.Context, q CandleQuery) (*CandlePage, string, error) {
	if err := q.validateSelectors(); err != nil {
		return nil, "", err
	}

	datasetName := strings.TrimSpace(q.DatasetName)
	datasetID := q.DatasetID
	if ticker := strings.TrimSpace(q.Ticker); ticker != "" {
		datasets, err := c.FindDatasetsByTicker(ctx, ticker)
		if err != nil {
			return nil, "", err
		}
		picked, err := pickDatasetForTicker(ticker, datasets)
		if err != nil {
			return nil, "", err
		}
		datasetName = picked.Name
		datasetID = 0
	}

	query := map[string]string{}
	if datasetID > 0 {
		query["dataset"] = strconv.Itoa(datasetID)
	}
	if datasetName != "" {
		query["dataset_name"] = datasetName
	}
	if !q.Start.IsZero() {
		query["start"] = q.Start.UTC().Format(time.RFC3339)
	}
	if !q.End.IsZero() {
		query["end"] = q.End.UTC().Format(time.RFC3339)
	}
	if q.Resolution != "" && q.Resolution != CandleResolutionRaw {
		query["resolution"] = string(q.Resolution)
	}
	if q.Limit > 0 {
		query["limit"] = strconv.Itoa(q.Limit)
	}
	if q.Offset > 0 {
		query["offset"] = strconv.Itoa(q.Offset)
	}
	if q.Ordering != "" {
		query["ordering"] = q.Ordering
	}

	reqURL, err := c.dataURLPath("rows/", query)
	if err != nil {
		return nil, "", errors.WithMessage(err, "could not parse url")
	}

	var wire atlasRowsPageWire
	if err := c.doDataJSON(ctx, "GetCandles", http.MethodGet, reqURL, &wire); err != nil {
		return nil, reqURL, err
	}
	page, err := decodeCandlePage(wire)
	return page, reqURL, err
}

func (c *apiClient) getCandlesNext(ctx context.Context, next string) (*CandlePage, error) {
	reqURL, err := c.resolveDataNextURL(next)
	if err != nil {
		return nil, errors.WithMessage(err, "could not resolve next url")
	}
	var wire atlasRowsPageWire
	if err := c.doDataJSON(ctx, "GetCandles", http.MethodGet, reqURL, &wire); err != nil {
		return nil, err
	}
	return decodeCandlePage(wire)
}

func (q CandleQuery) validateSelectors() error {
	n := 0
	if q.DatasetID > 0 {
		n++
	}
	if strings.TrimSpace(q.DatasetName) != "" {
		n++
	}
	if strings.TrimSpace(q.Ticker) != "" {
		n++
	}
	if n != 1 {
		return errors.New("CandleQuery requires exactly one of DatasetID, DatasetName, or Ticker")
	}
	return nil
}

func decodeDatasetPage(wire atlasDatasetPageWire) (*DatasetPage, error) {
	out := &DatasetPage{
		Count:    wire.Count,
		Next:     wire.Next,
		Previous: wire.Previous,
		Results:  make([]Dataset, 0, len(wire.Results)),
	}
	for _, w := range wire.Results {
		ds, err := decodeDataset(w)
		if err != nil {
			return nil, err
		}
		out.Results = append(out.Results, ds)
	}
	return out, nil
}

func decodeDataset(w atlasDatasetWire) (Dataset, error) {
	meta := map[string]any{}
	if len(w.Metadata) > 0 && string(w.Metadata) != "null" {
		if err := json.Unmarshal(w.Metadata, &meta); err != nil {
			return Dataset{}, errors.WithMessage(err, "failed to decode dataset metadata")
		}
	}
	return Dataset{
		ID:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		Metadata:    meta,
		OwnerKeyID:  w.OwnerKeyID,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}, nil
}

func decodeCandlePage(wire atlasRowsPageWire) (*CandlePage, error) {
	out := &CandlePage{
		Count:    wire.Count,
		Next:     wire.Next,
		Previous: wire.Previous,
		Results:  make([]Candle, 0, len(wire.Results)),
	}
	for _, row := range wire.Results {
		c, err := decodeCandle(row)
		if err != nil {
			return nil, err
		}
		out.Results = append(out.Results, c)
	}
	return out, nil
}

func decodeCandle(row atlasRowWire) (Candle, error) {
	c := Candle{
		DatasetID: row.DatasetID,
		Timestamp: row.Timestamp,
	}

	fullAggregate := row.Open != nil && row.High != nil && row.Low != nil && row.Close != nil
	if fullAggregate {
		c.Open = *row.Open
		c.High = *row.High
		c.Low = *row.Low
		c.Close = *row.Close
		if row.Volume != nil {
			c.Volume = *row.Volume
		}
		if row.Count != nil {
			c.Count = *row.Count
		}
		// Fill volume/trades from values when aggregate omitted them.
		if row.Volume == nil && len(row.Values) > 0 && string(row.Values) != "null" {
			var v atlasCandleValues
			if err := json.Unmarshal(row.Values, &v); err != nil {
				return c, errors.WithMessage(err, "failed to decode candle values")
			}
			c.Volume = v.Volume
			c.VolumeNotional = v.VolumeNotional
			c.TradesDone = v.TradesDone
		}
		return c, nil
	}

	if len(row.Values) == 0 || string(row.Values) == "null" {
		return c, errors.New("candle row missing values and aggregate fields")
	}
	var v atlasCandleValues
	if err := json.Unmarshal(row.Values, &v); err != nil {
		return c, errors.WithMessage(err, "failed to decode candle values")
	}
	c.Open = v.Open
	c.High = v.High
	c.Low = v.Low
	c.Close = v.Close
	c.Volume = v.Volume
	c.VolumeNotional = v.VolumeNotional
	c.TradesDone = v.TradesDone
	return c, nil
}

func pickDatasetForTicker(ticker string, datasets []Dataset) (Dataset, error) {
	if len(datasets) == 0 {
		return Dataset{}, errors.Wrap(ErrDatasetNotFound, "resolve ticker", "ticker", ticker)
	}
	if len(datasets) == 1 {
		return datasets[0], nil
	}
	// Explicit preference: unique 1-minute series when several datasets share a ticker.
	var preferred []Dataset
	for _, d := range datasets {
		name := strings.ToLower(d.Name)
		if strings.Contains(name, "1min") || strings.HasSuffix(name, "_1m") {
			preferred = append(preferred, d)
		}
	}
	if len(preferred) == 1 {
		return preferred[0], nil
	}
	return Dataset{}, errors.Wrap(ErrAmbiguousTicker, "resolve ticker: use DatasetID or DatasetName",
		"ticker", ticker, "matches", len(datasets))
}

func (c *apiClient) resolveDataNextURL(next string) (string, error) {
	next = strings.TrimSpace(next)
	if next == "" {
		return "", errors.New("empty next url")
	}
	u, err := url.Parse(next)
	if err != nil {
		return "", err
	}
	base, err := url.Parse(c.dataURL)
	if err != nil {
		return "", err
	}
	if u.IsAbs() {
		if u.Scheme != base.Scheme || u.Host != base.Host {
			return "", errors.New("next url host does not match data API")
		}
		return u.String(), nil
	}
	if strings.HasPrefix(u.Path, "/") {
		return (&url.URL{
			Scheme:   base.Scheme,
			Host:     base.Host,
			Path:     u.Path,
			RawQuery: u.RawQuery,
			Fragment: u.Fragment,
		}).String(), nil
	}
	// Resolve relative to dataURL with a trailing slash so "/api" is not dropped.
	baseWithSlash, err := url.Parse(strings.TrimRight(c.dataURL, "/") + "/")
	if err != nil {
		return "", err
	}
	return baseWithSlash.ResolveReference(u).String(), nil
}
