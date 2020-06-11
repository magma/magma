// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gocloud.dev/gcerrors"
)

func TestPipe(t *testing.T) {
	emitter, subscriber := Pipe()
	require.NotNil(t, emitter)
	require.NotNil(t, subscriber)
	ctx := context.Background()
	subscription, err := subscriber.Subscribe(ctx)
	require.NoError(t, err)

	err = emitter.Emit(ctx, t.Name(), t.Name(), nil)
	require.NoError(t, err)
	msg, err := subscription.Receive(ctx)
	require.NoError(t, err)
	require.Equal(t, t.Name(), msg.Metadata[TenantHeader])
	require.Equal(t, t.Name(), msg.Metadata[NameHeader])
	require.Empty(t, msg.Body)

	err = emitter.Shutdown(ctx)
	require.NoError(t, err)
	err = emitter.Emit(ctx, t.Name(), t.Name(), nil)
	require.Error(t, err)

	err = subscriber.Shutdown(ctx)
	require.NoError(t, err)
	_, err = subscriber.Subscribe(ctx)
	require.Error(t, err)
	err = subscription.Shutdown(ctx)
	require.Equal(t, gcerrors.FailedPrecondition, gcerrors.Code(err))
}
