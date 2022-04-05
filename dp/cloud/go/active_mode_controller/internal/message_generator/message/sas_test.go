package message_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/message"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"
)

const data = "some data"

func TestSasMessageString(t *testing.T) {
	m := message.NewSasMessage(data)
	expected := fmt.Sprintf("request: %s", data)
	assert.Equal(t, expected, m.String())
}

func TestSasMessageSend(t *testing.T) {
	client := &stubRequestsClient{}
	provider := &stubRequestsClientProvider{client: client}

	m := message.NewSasMessage(data)
	require.NoError(t, m.Send(context.Background(), provider))

	expected := &requests.RequestPayload{Payload: data}
	assert.Equal(t, expected, client.req)
}

type stubRequestsClientProvider struct {
	client *stubRequestsClient
}

func (s *stubRequestsClientProvider) GetRequestsClient() requests.RadioControllerClient {
	return s.client
}

func (s *stubRequestsClientProvider) GetActiveModeClient() active_mode.ActiveModeControllerClient {
	panic("not implemented")
}

type stubRequestsClient struct {
	req *requests.RequestPayload
}

func (s *stubRequestsClient) UploadRequests(_ context.Context, in *requests.RequestPayload, _ ...grpc.CallOption) (*requests.RequestDbIds, error) {
	s.req = in
	return &requests.RequestDbIds{}, nil
}
