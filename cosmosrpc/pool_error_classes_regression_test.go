package cosmosrpc

import (
	"context"
	"net"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/allora-network/allora-sdk-go/config"
	sdkgrpc "github.com/allora-network/allora-sdk-go/gen/grpc"
)

type voteStateServer struct {
	govv1.UnimplementedQueryServer
	hasVote bool
}

func (s *voteStateServer) Vote(context.Context, *govv1.QueryVoteRequest) (*govv1.QueryVoteResponse, error) {
	if !s.hasVote {
		return nil, status.Error(codes.InvalidArgument, "vote not found in this state")
	}
	voter := types.AccAddress(make([]byte, 20))
	return &govv1.QueryVoteResponse{Vote: &govv1.Vote{
		ProposalId: 1,
		Voter:      voter.String(),
		Options:    govv1.NewNonSplitVoteOption(govv1.OptionYes),
	}}, nil
}

func startVoteStateServer(t *testing.T, hasVote bool) string {
	t.Helper()
	server := grpcpkg.NewServer()
	govv1.RegisterQueryServer(server, &voteStateServer{hasVote: hasVote})
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	return "grpc://" + listener.Addr().String()
}

// TestPoolGRPCInvalidArgumentFailsOverToState uses the actual Gov wrapper and
// real gRPC servers, one without the vote and one that holds it.
func TestPoolGRPCInvalidArgumentFailsOverToState(t *testing.T) {
	logger := zerolog.Nop()
	urlMissing := startVoteStateServer(t, false)
	urlHeld := startVoteStateServer(t, true)
	clientMissing, err := sdkgrpc.NewGRPCClient(config.EndpointConfig{URL: urlMissing, Protocol: config.ProtocolGRPC}, logger)
	require.NoError(t, err)
	clientHeld, err := sdkgrpc.NewGRPCClient(config.EndpointConfig{URL: urlHeld, Protocol: config.ProtocolGRPC}, logger)
	require.NoError(t, err)

	p := NewClientPool([]Client{clientMissing, clientHeld}, logger, 0)
	t.Cleanup(func() { _ = p.Close() })

	response, err := p.Gov().Vote(context.Background(), &govv1.QueryVoteRequest{ProposalId: 1})

	require.NoError(t, err)
	require.Equal(t, uint64(1), response.Vote.ProposalId)
	require.Equal(t, 2, p.GetHealthStatus()["active_clients"])
}
