/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package active_mode_controller

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"magma/dp/cloud/go/services/dp/active_mode_controller/message_generator"
	"magma/dp/cloud/go/services/dp/active_mode_controller/protos/active_mode"
)

type App struct {
	dialer                Dialer
	clock                 Clock
	rng                   message_generator.RNG
	dialTimeout           time.Duration
	heartbeatSendTimeout  time.Duration
	requestTimeout        time.Duration
	pollingInterval       time.Duration
	grpcService           string
	cbsdInactivityTimeout time.Duration
}

func NewApp(options ...Option) *App {
	a := &App{}
	for _, o := range options {
		o(a)
	}
	return a
}

type Clock interface {
	Now() time.Time
	Tick(duration time.Duration) *time.Ticker
}

type Option func(*App)

type Dialer func(context.Context, string) (net.Conn, error)

func WithDialer(dialer Dialer) Option {
	return func(a *App) { a.dialer = dialer }
}

func WithRNG(rng message_generator.RNG) Option {
	return func(a *App) { a.rng = rng }
}

func WithClock(clock Clock) Option {
	return func(a *App) { a.clock = clock }
}

func WithDialTimeout(timeout time.Duration) Option {
	return func(a *App) { a.dialTimeout = timeout }
}

func WithHeartbeatSendTimeout(sendTimeout time.Duration, sendInterval time.Duration) Option {
	return func(a *App) { a.heartbeatSendTimeout = sendTimeout + sendInterval }
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(a *App) { a.requestTimeout = timeout }
}

func WithPollingInterval(interval time.Duration) Option {
	return func(a *App) { a.pollingInterval = interval }
}

func WithCbsdInactivityTimeout(timeout time.Duration) Option {
	return func(a *App) { a.cbsdInactivityTimeout = timeout }
}

func WithGrpcService(service string, port int) Option {
	return func(a *App) {
		a.grpcService = fmt.Sprintf("%s:%d", service, port)
	}
}

func (a *App) Run(ctx context.Context) error {
	conn, err := a.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := active_mode.NewActiveModeControllerClient(conn)
	ticker := a.clock.Tick(a.pollingInterval)
	defer ticker.Stop()
	generator := a.newGenerator()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			state, err := a.getState(ctx, client)
			if err != nil {
				log.Printf("failed to get state: %s", err)
				continue
			}
			messages := generator.GenerateMessages(state, a.clock.Now())
			for _, msg := range messages {
				if err := a.sendMessage(ctx, client, msg); err != nil {
					log.Printf("failed to send message '%s': %s", msg, err)
				}
			}
		}
	}
}

func (a *App) newGenerator() messageGenerator {
	return message_generator.NewMessageGenerator(
		a.heartbeatSendTimeout+a.pollingInterval,
		a.cbsdInactivityTimeout,
		a.rng,
	)
}

type messageGenerator interface {
	GenerateMessages(*active_mode.State, time.Time) []message_generator.Message
}

func (a *App) connect(ctx context.Context) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	if a.dialer != nil {
		opts = append(opts, grpc.WithContextDialer(a.dialer))
	}
	dialCtx, cancel := context.WithTimeout(ctx, a.dialTimeout)
	defer cancel()
	return grpc.DialContext(dialCtx, a.grpcService, opts...)
}

func (a *App) getState(ctx context.Context, c active_mode.ActiveModeControllerClient) (*active_mode.State, error) {
	glog.Infof("getting state")
	reqCtx, cancel := context.WithTimeout(ctx, a.requestTimeout)
	defer cancel()
	return c.GetState(reqCtx, &active_mode.GetStateRequest{})
}

func (a *App) sendMessage(ctx context.Context, client active_mode.ActiveModeControllerClient, msg message_generator.Message) error {
	glog.Infof("sending message: %s", msg)
	reqCtx, cancel := context.WithTimeout(ctx, a.requestTimeout)
	defer cancel()
	return msg.Send(reqCtx, client)
}
