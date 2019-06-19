package coa

import (
	"context"
	"testing"

	"fbc/cwf/radius/modules/coa/protos"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestAuthorizationChange(t *testing.T) {

	// Set up a connection to the Server.
	const address = "localhost:3798"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	assert.True(t, err == nil)

	ac := protos.NewAuthorizationClient(conn)
	context := context.Background()

	t.Run("Change", func(t *testing.T) {

		req := &protos.ChangeRequest{}

		r, err := ac.Change(context, req)
		assert.Error(t, err)
		assert.True(t, r == nil)
	})
}

func TestAuthorizationDissconnect(t *testing.T) {

	// Set up a connection to the Server.
	const address = "localhost:3798"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	assert.True(t, err == nil)

	ac := protos.NewAuthorizationClient(conn)
	context := context.Background()

	t.Run("Change", func(t *testing.T) {

		req := &protos.DisconnectRequest{}

		r, err := ac.Disconnect(context, req)
		assert.Error(t, err)
		assert.True(t, r == nil)

	})
}
