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

package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"magma/dp/cloud/go/active_mode_controller/config"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"
)

type App struct {
	additionalGrpcOpts []grpc.DialOption
	clock              Clock
	rng                message_generator.RNG
	cfg                *config.Config
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
	return func(a *App) {
		a.additionalGrpcOpts = append(a.additionalGrpcOpts, grpc.WithContextDialer(dialer))
	}
}

func WithRNG(rng message_generator.RNG) Option {
	return func(a *App) {
		a.rng = rng
	}
}

func WithClock(clock Clock) Option {
	return func(a *App) {
		a.clock = clock
	}
}

func WithConfig(cfg *config.Config) Option {
	return func(a *App) {
		a.cfg = cfg
	}
}

func (a *App) Run(ctx context.Context) error {
	conn, err := a.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	provider := &clientProvider{
		amcClient: active_mode.NewActiveModeControllerClient(conn),
		rcClient:  requests.NewRadioControllerClient(conn),
	}
	ticker := a.clock.Tick(a.cfg.PollingInterval)
	defer ticker.Stop()
	generator := newGenerator(a.cfg, a.rng)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			state, err := a.getState(ctx, provider.amcClient)
			if err != nil {
				log.Printf("failed to get state: %s", err)
				continue
			}
			messages := generator.GenerateMessages(state, a.clock.Now())
			for _, msg := range messages {
				if err := a.sendMessage(ctx, provider, msg); err != nil {
					log.Printf("failed to send message '%s': %s", msg, err)
				}
			}
		}
	}
}

func newGenerator(cfg *config.Config, rng message_generator.RNG) messageGenerator {
	return message_generator.NewMessageGenerator(
		cfg.HeartbeatSendTimeout+cfg.PollingInterval+cfg.RequestProcessingInterval,
		cfg.CbsdInactivityTimeout,
		rng,
	)
}

type messageGenerator interface {
	GenerateMessages(*active_mode.State, time.Time) []message_generator.Message
}

type clientProvider struct {
	amcClient active_mode.ActiveModeControllerClient
	rcClient  requests.RadioControllerClient
}

func (c *clientProvider) GetActiveModeClient() active_mode.ActiveModeControllerClient {
	return c.amcClient
}

func (c *clientProvider) GetRequestsClient() requests.RadioControllerClient {
	return c.rcClient
}

func (a *App) connect(ctx context.Context) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	opts = append(opts, a.additionalGrpcOpts...)
	dialCtx, cancel := context.WithTimeout(ctx, a.cfg.DialTimeout)
	defer cancel()
	addr := fmt.Sprintf("%s:%d", a.cfg.GrpcService, a.cfg.GrpcPort)
	return grpc.DialContext(dialCtx, addr, opts...)
}

func (a *App) getState(ctx context.Context, c active_mode.ActiveModeControllerClient) (*active_mode.State, error) {
	log.Printf("getting state")
	reqCtx, cancel := context.WithTimeout(ctx, a.cfg.RequestTimeout)
	defer cancel()
	return c.GetState(reqCtx, &active_mode.GetStateRequest{})
}

func (a *App) sendMessage(ctx context.Context, provider *clientProvider, msg message_generator.Message) error {
	log.Printf("sending message: %s", msg)
	reqCtx, cancel := context.WithTimeout(ctx, a.cfg.RequestTimeout)
	defer cancel()
	return msg.Send(reqCtx, provider)
}
