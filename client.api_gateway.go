package allora

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brynbellomy/go-utils/errors"
)

type inferenceOptions struct {
	evmChainID string
}

// InferenceOption configures GetNetworkInference (e.g. optional consumer signatures).
type InferenceOption func(*inferenceOptions)

// WithEVMChainID requests AlloraConsumer signatures for the given EVM chain id.
func WithEVMChainID(chainID uint64) InferenceOption {
	return func(o *inferenceOptions) {
		o.evmChainID = strconv.FormatUint(chainID, 10)
	}
}

// WithEVMChainIDSlug requests signatures using a configured chain slug (e.g. "base").
func WithEVMChainIDSlug(slug string) InferenceOption {
	return func(o *inferenceOptions) {
		o.evmChainID = strings.TrimSpace(slug)
	}
}

// LabeledValue is one label entry in a network-inference bundle.
type LabeledValue struct {
	LabelName string `json:"label_name"`
	Value     string `json:"value"`
}

// WorkerInference is one worker's labeled values (inferer or forecaster).
type WorkerInference struct {
	Worker string         `json:"worker"`
	Values []LabeledValue `json:"values"`
}

// NetworkInference is the full Forge gateway inference payload (extended helper).
type NetworkInference struct {
	Network                                 string            `json:"network"`
	TopicID                                 int64             `json:"topic_id"`
	CombinedValue                           []LabeledValue    `json:"combined_value"`
	NaiveValue                              []LabeledValue    `json:"naive_value,omitempty"`
	LabelNames                              []string          `json:"label_names,omitempty"`
	Nonce                                   int64             `json:"nonce,omitempty"`
	ExtraData                               string            `json:"extra_data,omitempty"`
	TimestampMs                             int64             `json:"timestamp_ms"`
	ConsensusTime                           time.Time         `json:"consensus_time"`
	HeightTx                                int64             `json:"height_tx"`
	InferenceBlockHeight                    string            `json:"inference_block_height"`
	InfererValues                           []WorkerInference `json:"inferer_values"`
	ForecasterValues                        []WorkerInference `json:"forecaster_values"`
	CachedAt                                time.Time         `json:"cached_at"`
	NextPredictionExpectedAt                time.Time         `json:"next_prediction_expected_at"`
	IsStale                                 bool              `json:"is_stale"`
	SecondsLate                             *float64          `json:"seconds_late"`
	NetworkInference                        string            `json:"network_inference"`
	NetworkInferenceNormalized              string            `json:"network_inference_normalized"`
	TokenDecimals                           int               `json:"token_decimals"`
	Timestamp                               int64             `json:"timestamp"`
	Signature                               string            `json:"signature"`
	ConfidenceIntervalPercentiles           []string          `json:"confidence_interval_percentiles"`
	ConfidenceIntervalPercentilesNormalized []string          `json:"confidence_interval_percentiles_normalized"`
	ConfidenceIntervalValues                []string          `json:"confidence_interval_values"`
	ConfidenceIntervalValuesNormalized      []string          `json:"confidence_interval_values_normalized"`
	EVMChainID                              uint64            `json:"evm_chain_id,omitempty"`
	ConsumerAddress                         string            `json:"consumer_address,omitempty"`
	VectorSignature                         string            `json:"vector_signature,omitempty"`
	VectorSignaturePublic                   string            `json:"vector_signature_public,omitempty"`
	SignatureExtraData                      string            `json:"signature_extra_data,omitempty"`
}

type gatewayTopicsResponse struct {
	Network  string         `json:"network"`
	Count    int            `json:"count"`
	CachedAt time.Time      `json:"cached_at"`
	Topics   []gatewayTopic `json:"topics"`
}

type gatewayTopic struct {
	ID                     int64      `json:"id"`
	Creator                string     `json:"creator"`
	Metadata               string     `json:"metadata"`
	LossMethod             string     `json:"loss_method"`
	EpochLastEnded         int64      `json:"epoch_last_ended"`
	EpochLength            int64      `json:"epoch_length"`
	GroundTruthLag         int64      `json:"ground_truth_lag"`
	WorkerSubmissionWindow int64      `json:"worker_submission_window"`
	CreatedAt              *time.Time `json:"created_at,omitempty"`
}

func (c *apiClient) GetTopics() iter.Seq2[*Topic, error] {
	return func(yield func(*Topic, error) bool) {
		topics, err := c.fetchTopics(context.Background())
		if err != nil {
			c.logger.Error().Err(err).Msg("failed to get topics")
			yield(nil, err)
			return
		}
		for i := range topics {
			if !yield(&topics[i], nil) {
				return
			}
		}
	}
}

// GetTopic returns one topic from the gateway catalog and attaches latest inference.
// The gateway has no single-topic route yet, so this scans the catalog list.
// Missing ids return ErrTopicNotFound. Inference failures are returned to the caller.
func (c *apiClient) GetTopic(topicID uint64) (*Topic, error) {
	topics, err := c.fetchTopics(context.Background())
	if err != nil {
		return nil, err
	}
	for i := range topics {
		if uint64(topics[i].TopicID) != topicID {
			continue
		}
		topic := topics[i]
		inf, err := c.GetInference(context.Background(), topicID)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to load latest inference")
		}
		topic.LatestNetworkInference = inf
		return &topic, nil
	}
	return nil, fmt.Errorf("%w: %d", ErrTopicNotFound, topicID)
}

func (c *apiClient) fetchTopics(ctx context.Context) ([]Topic, error) {
	var resp gatewayTopicsResponse
	reqURL, err := c.gatewayURL(string(c.network)+"/topics", nil)
	if err != nil {
		return nil, errors.WithMessage(err, "could not parse url")
	}
	if err := c.doJSON(ctx, "GetTopics", http.MethodGet, reqURL, &resp); err != nil {
		return nil, err
	}
	out := make([]Topic, 0, len(resp.Topics))
	for _, gt := range resp.Topics {
		out = append(out, mapGatewayTopic(gt))
	}
	return out, nil
}

func mapGatewayTopic(gt gatewayTopic) Topic {
	// The gateway's free-form `metadata` string (e.g. "BTC 5min") is the closest
	// equivalent to the previous SDK's TopicName. Other Upshot-only fields
	// (Description, WorkerCount, IsActive, stakes, UpdatedAt, …) stay zero/unset —
	// the gateway has no update timestamp, so we do not invent one from cache/created_at.
	return Topic{
		TopicID:                int(gt.ID),
		TopicName:              gt.Metadata,
		EpochLength:            int(gt.EpochLength),
		GroundTruthLag:         int(gt.GroundTruthLag),
		LossMethod:             gt.LossMethod,
		WorkerSubmissionWindow: int(gt.WorkerSubmissionWindow),
	}
}

// GetInference returns the previous SDK Inference shape (combined/naive/timestamp).
func (c *apiClient) GetInference(ctx context.Context, topicID uint64) (*Inference, error) {
	ni, err := c.GetNetworkInference(ctx, topicID)
	if err != nil {
		return nil, err
	}
	return mapNetworkInference(ni), nil
}

// GetNetworkInference returns the full Forge gateway inference payload.
func (c *apiClient) GetNetworkInference(ctx context.Context, topicID uint64, opts ...InferenceOption) (*NetworkInference, error) {
	options := &inferenceOptions{}
	for _, opt := range opts {
		opt(options)
	}

	query := map[string]string{}
	if options.evmChainID != "" {
		query["evm_chain_id"] = options.evmChainID
	}

	path := fmt.Sprintf("%s/inference/%d", c.network, topicID)
	reqURL, err := c.gatewayURL(path, query)
	if err != nil {
		return nil, errors.WithMessage(err, "could not parse url")
	}

	var inf NetworkInference
	if err := c.doJSON(ctx, "GetNetworkInference", http.MethodGet, reqURL, &inf); err != nil {
		return nil, err
	}
	return &inf, nil
}

func mapNetworkInference(ni *NetworkInference) *Inference {
	if ni == nil {
		return nil
	}
	// Prefer the single normalized decimal string. Avoid raw network_inference
	// (base-unit integer, ~10^24 larger) and avoid picking an arbitrary label
	// from multi-label combined_value bundles.
	combined := strings.TrimSpace(ni.NetworkInferenceNormalized)
	if combined == "" {
		combined = firstLabeledValue(ni.CombinedValue)
	}
	naive := firstLabeledValue(ni.NaiveValue)
	ts := ni.TimestampMs
	if ts == 0 && ni.Timestamp > 0 {
		ts = ni.Timestamp * 1000
	}
	return &Inference{
		CombinedValue: combined,
		NaiveValue:    naive,
		Timestamp:     ts,
	}
}

func firstLabeledValue(values []LabeledValue) string {
	if len(values) == 0 {
		return ""
	}
	return values[0].Value
}
