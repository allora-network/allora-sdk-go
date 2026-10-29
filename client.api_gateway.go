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

// NetworkInference is the Forge gateway inference payload.
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
	Network  string    `json:"network"`
	Count    int       `json:"count"`
	CachedAt time.Time `json:"cached_at"`
	Topics   []Topic   `json:"topics"`
}

// Topic is a Forge gateway catalog entry (native wire shape).
// LatestNetworkInference is attached client-side by GetTopic (best-effort) and
// is not part of the gateway topics list payload.
type Topic struct {
	ID                     int64      `json:"id"`
	Creator                string     `json:"creator"`
	Metadata               string     `json:"metadata"`
	LossMethod             string     `json:"loss_method"`
	EpochLastEnded         int64      `json:"epoch_last_ended"`
	EpochLength            int64      `json:"epoch_length"`
	GroundTruthLag         int64      `json:"ground_truth_lag"`
	WorkerSubmissionWindow int64      `json:"worker_submission_window"`
	CreatedAt              *time.Time `json:"created_at,omitempty"`

	LatestNetworkInference *NetworkInference `json:"latest_network_inference,omitempty"`
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

// GetTopic returns one topic from the gateway catalog and best-effort attaches
// the latest network inference. The gateway has no single-topic route yet, so
// this scans the catalog list. Missing ids return ErrTopicNotFound. Inference
// failures leave LatestNetworkInference nil so catalog metadata stays
// independently usable (inference is a separately billed route).
func (c *apiClient) GetTopic(topicID uint64) (*Topic, error) {
	topics, err := c.fetchTopics(context.Background())
	if err != nil {
		return nil, err
	}
	for i := range topics {
		if uint64(topics[i].ID) != topicID {
			continue
		}
		topic := topics[i]
		inf, err := c.GetNetworkInference(context.Background(), topicID)
		if err != nil {
			c.logger.Warn().Err(err).Uint64("topic_id", topicID).Msg("latest inference unavailable")
		} else {
			topic.LatestNetworkInference = inf
		}
		return &topic, nil
	}
	return nil, errors.Wrap(ErrTopicNotFound, "get topic", "topic_id", topicID)
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
	return resp.Topics, nil
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
