package app

import (
	"context"
	"fmt"
	"log"
	"magma/dp/cloud/go/active_mode_controller/config"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"
	"net"
	"time"

	"google.golang.org/grpc"
)

type App struct {
	additionalGrpcOpts []grpc.DialOption
	clock              Clock
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
	stateGetter := active_mode.NewActiveModeControllerClient(conn)
	requestSender := requests.NewRadioControllerClient(conn)
	ticker := a.clock.Tick(a.cfg.PollingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			state, err := a.getState(ctx, stateGetter)
			if err != nil {
				log.Printf("failed to get state: %s", err)
				continue
			}
			messages := message_generator.GenerateMessages(a.clock.Now, state)
			for _, request := range messages {
				if _, err := a.uploadRequests(ctx, requestSender, request); err != nil {
					log.Printf("failed to send request '%s': %s", request.Payload, err)
				}
			}
		}
	}
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

func (a *App) uploadRequests(ctx context.Context, c requests.RadioControllerClient, req *requests.RequestPayload) (*requests.RequestDbIds, error) {
	log.Printf("uploading request: '%s'", req.Payload)
	reqCtx, cancel := context.WithTimeout(ctx, a.cfg.RequestTimeout)
	defer cancel()
	return c.UploadRequests(reqCtx, req)
}
