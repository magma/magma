/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package streamer provides streamer client Go implementation for golang based gateways
package streamer

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"

	"magma/feg/gateway/registry"
	"magma/orc8r/lib/go/protos"
)

const (
	StreamerServiceName = "streamer"
	StreamingInterval   = time.Second * 20
)

type listener struct {
	Listener
	done int32
}

type streamerClient struct {
	listeners       map[string]*listener
	listenersMu     sync.Mutex
	serviceRegistry registry.CloudRegistry
}

// NewStreamerClient creates new streamer client with an empty listeners list and provided cloud registry
// The created streamer is ready to serve new listeners after they are added via AddListener() call
func NewStreamerClient(cr registry.CloudRegistry) Client {
	if cr == nil {
		cr = registry.NewCloudRegistry()
	}
	return &streamerClient{listeners: map[string]*listener{}, serviceRegistry: cr}
}

// AddListener registers a new streaming updates listener for the
// listener.GetName() stream and starts stream loop routine for it.
// The stream name must be unique and AddListener will error out if a listener
// for the same stream is already registered.
func (cl *streamerClient) AddListener(l Listener) error {
	if l == nil {
		return fmt.Errorf("Invalid (nil) Listener")
	}
	stream := l.GetName()
	cl.listenersMu.Lock()
	defer cl.listenersMu.Unlock()
	if lis, exist := cl.listeners[stream]; exist && lis != nil {
		return fmt.Errorf("Listener for stream '%s' already exist", stream)
	}

	lis := &listener{Listener: l, done: 0}
	cl.listeners[stream] = lis
	go cl.streamUpdates(lis)

	return nil
}

// RemoveListener removes currently registered listener. It returns true is the
// listener with provided l.GetName() exists and was unregistered successfully
// RemoveListener is the only way to terminate streaming loop
func (cl *streamerClient) RemoveListener(l Listener) bool {
	if l != nil {
		stream := l.GetName()
		cl.listenersMu.Lock()
		defer cl.listenersMu.Unlock()
		if lis, exist := cl.listeners[stream]; exist {
			lis.setDone()
			delete(cl.listeners, stream)
			return true
		}
	}
	return false
}

func (cl *streamerClient) streamUpdates(l *listener) {
	for {
		conn, grpcStreamClient, err := cl.startStreaming(l)
		if err != nil {
			// Notify Listener & Check if it wants to continue
			l.notifyError(fmt.Errorf("Failed to create stream for %s: %v", StreamerServiceName, err))
		} else {
			for !l.isDone() {
				updatesBatch, err := grpcStreamClient.Recv()
				if err != nil {
					// Don't notify on EOF, streaming service may have closed stream due to
					// being idle 4 too long/restart/etc. In this case we'll just reconnect
					if err != io.EOF {
						l.notifyError(fmt.Errorf("Stream %s receive error: %v", StreamerServiceName, err))
					}
					break // reconnect and continue or exit
				}
				if !l.Update(updatesBatch) {
					l.setDone() // Listener indicated not to continue streaming, cleanup and return
				} else {
					time.Sleep(StreamingInterval)
				}
			}
			conn.Close()
		}
		if l.isDone() {
			break
		}
		time.Sleep(StreamingInterval)
	}
	cl.RemoveListener(l)
}

func (cl *streamerClient) startStreaming(l *listener) (*grpc.ClientConn, protos.Streamer_GetUpdatesClient, error) {
	conn, err := cl.serviceRegistry.GetCloudConnection(StreamerServiceName)
	if err != nil {
		return nil, nil, err
	}
	grpcStreamerClient, err := protos.NewStreamerClient(conn).GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: "", StreamName: l.GetName()},
	)
	if err != nil {
		conn.Close()
		return nil, nil, err
	}
	return conn, grpcStreamerClient, err
}

func (l *listener) clearDone() {
	atomic.StoreInt32(&l.done, 0)
}

func (l *listener) setDone() {
	atomic.StoreInt32(&l.done, 1)
}

func (l *listener) isDone() bool {
	return atomic.LoadInt32(&l.done) != 0
}

func (l *listener) notifyError(err error) bool {
	if l.ReportError(err) == nil { // Notify Listener & Check if it wants to continue
		return true
	}
	l.setDone()
	return false
}
