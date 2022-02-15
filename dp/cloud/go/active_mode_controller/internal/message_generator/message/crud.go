package message

import (
	"context"
	"fmt"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func NewDeleteMessage(serialNumber string) *deleteMessage {
	return &deleteMessage{serialNumber: serialNumber}
}

type deleteMessage struct {
	serialNumber string
}

func (d *deleteMessage) Send(ctx context.Context, provider ClientProvider) error {
	req := &active_mode.DeleteCbsdRequest{SerialNumber: d.serialNumber}
	client := provider.GetActiveModeClient()
	_, err := client.DeleteCbsd(ctx, req)
	return err
}

func (d *deleteMessage) String() string {
	return fmt.Sprintf("delete: %s", d.serialNumber)
}
