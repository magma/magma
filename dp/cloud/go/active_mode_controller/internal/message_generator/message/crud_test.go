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

const serialNumber = "some_serial_number"

func TestDeleteMessageString(t *testing.T) {
	m := message.NewDeleteMessage(serialNumber)
	expected := fmt.Sprintf("delete: %s", serialNumber)
	assert.Equal(t, expected, m.String())
}

func TestDeleteMessageSend(t *testing.T) {
	client := &stubActiveModeClient{}
	provider := &stubActiveModeClientProvider{client: client}

	m := message.NewDeleteMessage(serialNumber)
	require.NoError(t, m.Send(context.Background(), provider))

	expected := &active_mode.DeleteCbsdRequest{SerialNumber: serialNumber}
	assert.Equal(t, expected, client.req)
}

type stubActiveModeClientProvider struct {
	client *stubActiveModeClient
}

func (s *stubActiveModeClientProvider) GetRequestsClient() requests.RadioControllerClient {
	panic("not implemented")
}

func (s *stubActiveModeClientProvider) GetActiveModeClient() active_mode.ActiveModeControllerClient {
	return s.client
}

type stubActiveModeClient struct {
	req *active_mode.DeleteCbsdRequest
}

func (s *stubActiveModeClient) DeleteCbsd(_ context.Context, in *active_mode.DeleteCbsdRequest, _ ...grpc.CallOption) (*active_mode.DeleteCbsdResponse, error) {
	s.req = in
	return &active_mode.DeleteCbsdResponse{}, nil
}

func (s *stubActiveModeClient) GetState(_ context.Context, _ *active_mode.GetStateRequest, _ ...grpc.CallOption) (*active_mode.State, error) {
	panic("not implemented")
}
