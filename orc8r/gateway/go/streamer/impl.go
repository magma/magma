/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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

	"magma/gateway/service_registry"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
)

const (
	StreamingInterval = time.Second * 20
)

type listener struct {
	Listener
	done int32
}

type streamerClient struct {
	listenersMu     sync.Mutex
	listeners       map[string]*listener
	serviceRegistry service_registry.GatewayRegistry
}

// NewStreamerClient creates new streamer client with an empty listeners list
// The created streamer is ready to serve new listeners after they are added via AddListener() call
func NewStreamerClient(reg service_registry.GatewayRegistry) Client {
	if reg == nil {
		reg = service_registry.Get()
	}
	return &streamerClient{listeners: map[string]*listener{}, serviceRegistry: reg}
}

// AddListener registers a new streaming updates listener for the
// listener.GetName() stream
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

// Stream starts streaming loop for a registered by AddListener listener
// If successful, Stream never return and should be called in it's own go routine or main()
// If the provided Listener is not registered, Stream will try to register it prior to starting streaming
func (cl *streamerClient) Stream(l Listener) error {
	if cl == nil {
		return fmt.Errorf("Invalid (nil) Stream Client")
	}
	if l == nil {
		return fmt.Errorf("Invalid (nil) Listener")
	}
	stream := l.GetName()
	cl.listenersMu.Lock()
	lis, exist := cl.listeners[stream]
	if !exist && lis == nil {
		lis = &listener{Listener: l, done: 0}
		cl.listeners[stream] = lis
	}
	cl.listenersMu.Unlock()
	cl.streamUpdates(lis)
	return nil
}

func (cl *streamerClient) streamUpdates(l *listener) {
	for {
		conn, grpcStreamClient, err := cl.startStreaming(l)
		if err != nil {
			// Notify Listener & Check if it wants to continue
			l.notifyError(fmt.Errorf("Failed to create stream for %s: %v", definitions.StreamerServiceName, err))
		} else {
			for !l.isDone() {
				updatesBatch, err := grpcStreamClient.Recv()
				if err != nil {
					// Don't notify on EOF, streaming service may have closed stream due to
					// being idle 4 too long/restart/etc. In this case we'll just reconnect
					if err != io.EOF {
						l.notifyError(fmt.Errorf("Stream %s receive error: %v", definitions.StreamerServiceName, err))
					}
					break // reconnect and continue or exit
				}
				if !l.Update(updatesBatch) {
					// Listener indicated not to continue streaming
					// send io.EOF to Listener's ReportError receiver to give it an option to terminate streaming
					// break & cleanup and depending on the result of ReportError reopen stream or terminate
					l.notifyError(io.EOF)
					break
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
	conn, err := cl.serviceRegistry.GetCloudConnection(definitions.StreamerServiceName)
	if err != nil {
		return nil, nil, err
	}
	req := &protos.StreamRequest{GatewayId: "", StreamName: l.GetName(), ExtraArgs: l.GetExtraArgs()}
	grpcStreamerClient, err := protos.NewStreamerClient(conn).GetUpdates(context.Background(), req)
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
