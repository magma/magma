package message

import (
	"context"
	"fmt"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func NewUpdateMessage(id int64) *updateMessage {
	return &updateMessage{id: id}
}

type updateMessage struct {
	id    int64
	delta int64
}

func (u *updateMessage) Send(ctx context.Context, provider ClientProvider) error {
	req := &active_mode.AcknowledgeCbsdUpdateRequest{Id: u.id}
	client := provider.GetActiveModeClient()
	_, err := client.AcknowledgeCbsdUpdate(ctx, req)
	return err
}

func (u *updateMessage) String() string {
	return fmt.Sprintf("update: %d", u.id)
}
