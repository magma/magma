// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/symphony/pkg/ent"

	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/pkg/viewer"
)

func newTestServer(t *testing.T, client *ent.Client, subscriber pubsub.Subscriber, handlers []Handler) *Server {
	return &Server{
		tenancy:    viewer.NewFixedTenancy(client),
		logger:     logtest.NewTestLogger(t),
		subscriber: subscriber,
		handlers:   handlers,
	}
}

func getLogEntry() pubsub.LogEntry {
	return pubsub.LogEntry{
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
	emitter, subscriber := pubsub.Pipe()
	defer func() {
		ctx := context.Background()
		_ = emitter.Shutdown(ctx)
		_ = subscriber.Shutdown(ctx)
	}()
	logEntry := getLogEntry()
	data, err := pubsub.Marshal(logEntry)
	require.NoError(t, err)
	client := viewertest.NewTestClient(t)
	ctx, cancel := context.WithCancel(context.Background())
	h := Func(func(ctx context.Context, entry pubsub.LogEntry) error {
		v := viewer.FromContext(ctx)
		require.Equal(t, tenantName, v.Tenant())
		require.Equal(t, serviceName, v.Name())
		require.Equal(t, user.RoleOWNER, v.Role())
		require.EqualValues(t, logEntry, entry)
		cancel()
		return nil
	})
	server := newTestServer(t, client, subscriber, []Handler{h})
	var wg sync.WaitGroup
	listener, err := server.Subscribe(ctx, &wg)
	require.NoError(t, err)
	defer listener.Shutdown(ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := listener.Listen(ctx)
		require.True(t, errors.Is(err, context.Canceled))
	}()
	err = emitter.Emit(ctx, tenantName, pubsub.EntMutation, data)
	require.NoError(t, err)
	wg.Wait()
}

func TestServerBadData(t *testing.T) {
	emitter, subscriber := pubsub.Pipe()
	defer func() {
		ctx := context.Background()
		_ = emitter.Shutdown(ctx)
		_ = subscriber.Shutdown(ctx)
	}()
	client := viewertest.NewTestClient(t)
	ctx, cancel := context.WithCancel(context.Background())
	h := Func(func(context.Context, pubsub.LogEntry) error {
		cancel()
		return nil
	})
	server := newTestServer(t, client, subscriber, []Handler{h})
	var wg sync.WaitGroup
	listener, err := server.Subscribe(ctx, &wg)
	require.NoError(t, err)
	defer listener.Shutdown(ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := listener.Listen(ctx)
		require.Error(t, err)
		require.False(t, errors.Is(err, context.Canceled))
	}()
	err = emitter.Emit(ctx, viewertest.DefaultTenant, pubsub.EntMutation, []byte(""))
	require.NoError(t, err)
	wg.Wait()
}

func TestServerHandlerError(t *testing.T) {
	tenantName := "Random"
	emitter, subscriber := pubsub.Pipe()
	defer func() {
		ctx := context.Background()
		_ = emitter.Shutdown(ctx)
		_ = subscriber.Shutdown(ctx)
	}()
	logEntry := getLogEntry()
	data, err := pubsub.Marshal(logEntry)
	require.NoError(t, err)
	client := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), client)
	cancelledCtx, cancel := context.WithCancel(ctx)

	h := Func(func(ctx context.Context, entry pubsub.LogEntry) error {
		client := ent.FromContext(ctx)
		client.LocationType.Create().
			SetName("LocationType").
			SaveX(ctx)
		cancel()
		return errors.New("operation failed")
	})
	server := newTestServer(t, client, subscriber, []Handler{h})
	var wg sync.WaitGroup
	listener, err := server.Subscribe(cancelledCtx, &wg)
	require.NoError(t, err)
	defer listener.Shutdown(ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := listener.Listen(cancelledCtx)
		require.True(t, errors.Is(err, context.Canceled))
	}()
	err = emitter.Emit(cancelledCtx, tenantName, pubsub.EntMutation, data)
	require.NoError(t, err)
	wg.Wait()
	require.False(t, client.LocationType.Query().Where().ExistX(ctx))
}

func TestServerHandlerNoError(t *testing.T) {
	tenantName := "Random"
	emitter, subscriber := pubsub.Pipe()
	defer func() {
		ctx := context.Background()
		_ = emitter.Shutdown(ctx)
		_ = subscriber.Shutdown(ctx)
	}()
	logEntry := getLogEntry()
	data, err := pubsub.Marshal(logEntry)
	require.NoError(t, err)
	client := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), client)
	cancelledCtx, cancel := context.WithCancel(ctx)

	h := Func(func(ctx context.Context, entry pubsub.LogEntry) error {
		client := ent.FromContext(ctx)
		client.LocationType.Create().
			SetName("LocationType").
			SaveX(ctx)
		cancel()
		return nil
	})
	server := newTestServer(t, client, subscriber, []Handler{h})
	var wg sync.WaitGroup
	listener, err := server.Subscribe(cancelledCtx, &wg)
	require.NoError(t, err)
	defer listener.Shutdown(ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := listener.Listen(cancelledCtx)
		require.True(t, errors.Is(err, context.Canceled))
	}()
	err = emitter.Emit(cancelledCtx, tenantName, pubsub.EntMutation, data)
	require.NoError(t, err)
	wg.Wait()
	require.True(t, client.LocationType.Query().Where().ExistX(ctx))
}
