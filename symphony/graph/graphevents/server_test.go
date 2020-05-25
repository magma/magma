// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphevents

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/event"

	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
)

func newTestServer(t *testing.T, client *ent.Client, subscriber event.Subscriber, handlers []Handler) *Server {
	return &Server{
		tenancy:    viewer.NewFixedTenancy(client),
		logger:     logtest.NewTestLogger(t),
		subscriber: subscriber,
		handlers:   handlers,
	}
}

func getLogEntry() event.LogEntry {
	return event.LogEntry{
		UserName:  "",
		UserID:    nil,
		Time:      time.Time{},
		Operation: ent.OpCreate,
		PrevState: nil,
		CurrState: &ent.Node{
			ID:   rand.Int(),
			Type: "Dog",
			Fields: []*ent.Field{
				{
					Type:  "string",
					Name:  "Name",
					Value: "Lassie",
				},
			},
			Edges: nil,
		},
	}
}

func TestServer(t *testing.T) {
	tenantName := "Random"
	emitter, subscriber := event.Pipe()
	defer func() {
		ctx := context.Background()
		_ = emitter.Shutdown(ctx)
		_ = subscriber.Shutdown(ctx)
	}()
	logEntry := getLogEntry()
	data, err := event.Marshal(logEntry)
	require.NoError(t, err)
	client := viewertest.NewTestClient(t)
	ctx, cancel := context.WithCancel(context.Background())
	h := HandlerFunc(func(ctx context.Context, entry event.LogEntry) error {
		v := viewer.FromContext(ctx)
		require.Equal(t, tenantName, v.Tenant())
		require.Equal(t, serviceName, v.Name())
		require.Equal(t, user.RoleOWNER, v.Role())
		require.EqualValues(t, logEntry, entry)
		cancel()
		return nil
	})
	server := newTestServer(t, client, subscriber, []Handler{h})
	listener, err := server.Subscribe(ctx)
	require.NoError(t, err)
	defer listener.Shutdown(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := listener.Listen(ctx)
		require.True(t, errors.Is(err, context.Canceled))
	}()
	err = emitter.Emit(ctx, tenantName, event.EntMutation, data)
	require.NoError(t, err)
	wg.Wait()
}

func TestServerBadData(t *testing.T) {
	emitter, subscriber := event.Pipe()
	defer func() {
		ctx := context.Background()
		_ = emitter.Shutdown(ctx)
		_ = subscriber.Shutdown(ctx)
	}()
	client := viewertest.NewTestClient(t)
	ctx, cancel := context.WithCancel(context.Background())
	h := HandlerFunc(func(context.Context, event.LogEntry) error {
		cancel()
		return nil
	})
	server := newTestServer(t, client, subscriber, []Handler{h})
	listener, err := server.Subscribe(ctx)
	require.NoError(t, err)
	defer listener.Shutdown(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := listener.Listen(ctx)
		require.Error(t, err)
		require.False(t, errors.Is(err, context.Canceled))
	}()
	err = emitter.Emit(ctx, viewertest.DefaultTenant, event.EntMutation, []byte(""))
	require.NoError(t, err)
	wg.Wait()
}

func TestServerHandlerError(t *testing.T) {
	tenantName := "Random"
	emitter, subscriber := event.Pipe()
	defer func() {
		ctx := context.Background()
		_ = emitter.Shutdown(ctx)
		_ = subscriber.Shutdown(ctx)
	}()
	logEntry := getLogEntry()
	data, err := event.Marshal(logEntry)
	require.NoError(t, err)
	client := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), client)
	h := HandlerFunc(func(ctx context.Context, entry event.LogEntry) error {
		client := ent.FromContext(ctx)
		client.LocationType.Create().
			SetName("LocationType").
			SaveX(ctx)
		return errors.New("operation failed")
	})
	server := newTestServer(t, client, subscriber, []Handler{h})
	listener, err := server.Subscribe(ctx)
	require.NoError(t, err)
	defer listener.Shutdown(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := listener.Listen(ctx)
		require.Error(t, err)
		require.False(t, client.LocationType.Query().Where().ExistX(ctx))
	}()
	err = emitter.Emit(ctx, tenantName, event.EntMutation, data)
	require.NoError(t, err)
	wg.Wait()
}
