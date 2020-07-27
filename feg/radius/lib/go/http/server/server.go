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

package server

import (
	"context"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"fbc/lib/go/http/middleware"
	"fbc/lib/go/log"

	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type (
	// Server provides a http server implementation.
	Server struct {
		http.Server
		Mux      *http.ServeMux
		Logger   log.Factory
		closers  sync.WaitGroup
		listener net.Listener
	}

	// An Option configures the server.
	Option struct {
		apply  func(*Server) error
		weight int
	}
)

// OptionFunc is a helper to allow the use of
// ordinary function as Option.
func OptionFunc(f func(*Server) error) Option {
	return Option{apply: f}
}

// Logger can be provided to override server logger.
func Logger(logger *zap.Logger) Option {
	return Option{
		apply: func(srv *Server) error {
			srv.Logger = log.NewFactory(logger)
			srv.ErrorLog = zap.NewStdLog(logger)
			return nil
		},
		weight: math.MinInt32,
	}
}

// Closer can be provided and are executes on server termination.
func Closer(closer io.Closer) Option {
	return OptionFunc(func(srv *Server) error {
		srv.RegisterOnShutdown(func() {
			if err := closer.Close(); err != nil {
				srv.Logger.Bg().Warn("running closer", zap.Error(err))
			}
			srv.closers.Done()
		})
		srv.closers.Add(1)
		return nil
	})
}

// Listener can be provided to override server listener.
func Listener(listener net.Listener) Option {
	return OptionFunc(func(srv *Server) error {
		srv.Addr = listener.Addr().String()
		srv.listener = listener
		return nil
	})
}

// New creates a new http server.
func New(config Config, options ...Option) (*Server, error) {
	sort.Slice(options, func(i, j int) bool {
		return options[i].weight < options[j].weight
	})
	idx := sort.Search(len(options), func(i int) bool {
		return options[i].weight >= 0
	})

	srv := &Server{Server: http.Server{Addr: config.Addr}}
	err := srv.Apply(options[:idx]...)
	if err != nil {
		return nil, errors.Wrap(err, "applying early options")
	}

	if srv.Logger == nil {
		logger, err := config.createLogger()
		if err != nil {
			return nil, err
		}
		_ = srv.Apply(Logger(logger))
	}

	srv.Mux, srv.Handler = config.createServeMux(srv.Logger)
	err = srv.Apply(options[idx:]...)
	if err != nil {
		return nil, errors.Wrap(err, "applying late options")
	}

	return srv, nil
}

// Apply applied server options.
func (srv *Server) Apply(options ...Option) error {
	var err error
	for _, option := range options {
		err = multierr.Append(err, option.apply(srv))
	}
	return err
}

// ServeFiles serves files from the given file system root.
func (srv *Server) ServeFiles(path string, root http.FileSystem) {
	fs := http.FileServer(root)
	srv.Mux.Handle(path, http.StripPrefix(path, fs))
}

// Handle registers the handler for the given pattern.
func (srv *Server) Handle(pattern string, handler http.Handler) {
	chain := alice.New(
		middleware.Tracing(),
		middleware.Logger(srv.Logger.Bg()),
		middleware.Recovery(
			middleware.RecoveryLogger(srv.Logger),
		),
	)
	srv.Mux.Handle(pattern, chain.Then(handler))
}

// HandleFunc registers the handler function for the given pattern.
func (srv *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	srv.Handle(pattern, http.HandlerFunc(handler))
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.Handler.ServeHTTP(w, r)
}

func (srv *Server) serve() error {
	if srv.listener != nil {
		return srv.Serve(srv.listener)
	}
	return srv.ListenAndServe()
}

func (srv *Server) run(done <-chan os.Signal) error {
	logger := srv.Logger.Bg()
	defer logger.Sync()

	defer srv.closers.Wait()

	go func() {
		<-done
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	logger.Info("starting http server", zap.String("address", srv.Addr))
	err := srv.serve()
	logger.Info("terminating http server", zap.Error(err))

	return err
}

// Run starts the http server.
func (srv *Server) Run() error {
	done := make(chan os.Signal, 1)
	defer close(done)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	return srv.run(done)
}
