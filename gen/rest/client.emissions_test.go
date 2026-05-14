package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
)

func TestEmissionsRESTClientUsesV10LatestInputInferenceRoute(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte("{}"))
	}))
	defer server.Close()

	client := NewEmissionsRESTClient(NewRESTClientCore(server.URL, zerolog.Nop()), zerolog.Nop())
	_, err := client.GetWorkerLatestInputInferenceByTopicId(context.Background(), &emissionstypes.GetWorkerLatestInputInferenceByTopicIdRequest{
		TopicId:       7,
		WorkerAddress: "allo1worker",
	})
	require.NoError(t, err)
	require.Equal(t, "/emissions/v10/topics/7/workers/allo1worker/latest_input_inference", gotPath)
}

func TestEmissionsRESTClientUsesV10OpenReputerSubmissionWindowsRoute(t *testing.T) {
	gotPath := captureEmissionsRESTPath(t, func(ctx context.Context, client *EmissionsRESTClient) error {
		_, err := client.GetOpenReputerSubmissionWindows(ctx, &emissionstypes.GetOpenReputerSubmissionWindowsRequest{
			TopicId: 42,
		})
		return err
	})

	require.Equal(t, "/emissions/v10/open_reputer_submission_windows/42", gotPath)
}

func TestEmissionsRESTClientUsesV10OpenWorkerSubmissionWindowsRoute(t *testing.T) {
	gotPath := captureEmissionsRESTPath(t, func(ctx context.Context, client *EmissionsRESTClient) error {
		_, err := client.GetOpenWorkerSubmissionWindows(ctx, &emissionstypes.GetOpenWorkerSubmissionWindowsRequest{
			TopicId: 42,
		})
		return err
	})

	require.Equal(t, "/emissions/v10/open_worker_submission_windows/42", gotPath)
}

func captureEmissionsRESTPath(t *testing.T, call func(context.Context, *EmissionsRESTClient) error) string {
	t.Helper()

	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte("{}"))
	}))
	defer server.Close()

	client := NewEmissionsRESTClient(NewRESTClientCore(server.URL, zerolog.Nop()), zerolog.Nop())
	require.NoError(t, call(context.Background(), client))

	return gotPath
}
