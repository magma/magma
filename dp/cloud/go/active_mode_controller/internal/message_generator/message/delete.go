package message

import (
	"context"
	"fmt"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func NewDeleteMessage(id int64) *deleteMessage {
	return &deleteMessage{id: id}
}

type deleteMessage struct {
	id int64
}

func (d *deleteMessage) Send(ctx context.Context, provider ClientProvider) error {
	req := &active_mode.DeleteCbsdRequest{Id: d.id}
	client := provider.GetActiveModeClient()
	_, err := client.DeleteCbsd(ctx, req)
	return err
}

func (d *deleteMessage) String() string {
	return fmt.Sprintf("delete: %d", d.id)
}
