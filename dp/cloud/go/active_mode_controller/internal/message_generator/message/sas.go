package message

import (
	"context"
	"fmt"

	"magma/dp/cloud/go/active_mode_controller/protos/requests"
)

func NewSasMessage(data string) *sasMessage {
	return &sasMessage{data: data}
}

type sasMessage struct {
	data string
}

func (s *sasMessage) Send(ctx context.Context, provider ClientProvider) error {
	payload := &requests.RequestPayload{Payload: s.data}
	client := provider.GetRequestsClient()
	_, err := client.UploadRequests(ctx, payload)
	return err
}

func (s *sasMessage) String() string {
	return fmt.Sprintf("request: %s", s.data)
}
