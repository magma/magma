// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipe(t *testing.T) {
	emitter, subscriber := Pipe()
	require.NotNil(t, emitter)
	require.NotNil(t, subscriber)
	subscription, err := subscriber.Subscribe(context.Background())
	require.NoError(t, err)

	err = emitter.Emit(context.Background(), t.Name(), t.Name(), nil)
	require.NoError(t, err)
	msg, err := subscription.Receive(context.Background())
	require.NoError(t, err)
	require.Equal(t, t.Name(), msg.Metadata[TenantHeader])
	require.Equal(t, t.Name(), msg.Metadata[NameHeader])
	require.Empty(t, msg.Body)
}
