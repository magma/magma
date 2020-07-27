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
	"io"
	"net/http"

	// register pprof with default http mux
	_ "net/http/pprof"

	"fbc/lib/go/http/middleware"
	"fbc/lib/go/log"

	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	// DiscardLogger is a nop logger.
	DiscardLogger = "discard"
	// DevelopmentLogger is a logger that should be used during development.
	DevelopmentLogger = "development"
	// ProductionLogger is a logger that should be used in production.
	ProductionLogger = "production"
)

type (
	// Config defines the server config.
	Config struct {
		Addr   string
		Logger *LoggerConfig
	}

	// LoggerConfig defines logger config.
	LoggerConfig struct {
		Type    string `envconfig:"TYPE" env:"TYPE" long:"type" default:"development" description:"runtime logger type"`
		Verbose bool   `envconfig:"VERBOSE" env:"VERBOSE" long:"verbose" description:"enable verbose logging"`
	}
)

// DefaultLoggerConfig is the default logger configuration.
var DefaultLoggerConfig = &LoggerConfig{
	Type: DiscardLogger,
}

func (config *Config) createLogger() (*zap.Logger, error) {
	lc := config.Logger
	if lc == nil || lc.Type == "" {
		lc = DefaultLoggerConfig
	}

	var (
		logger *zap.Logger
		err    error
	)
	switch lc.Type {
	case DiscardLogger:
		logger = zap.NewNop()
	case DevelopmentLogger:
		logger, err = zap.NewDevelopment()
	case ProductionLogger:
		config := zap.NewProductionConfig()
		if lc.Verbose {
			config.Level.SetLevel(zap.DebugLevel)
		}
		logger, err = config.Build()
	default:
		err = errors.Errorf("unknown logger type: %q", lc.Type)
	}

	if err != nil {
		return nil, errors.Wrap(err, "creating logger")
	}

	return logger.WithOptions(zap.AddStacktrace(zap.DPanicLevel)), nil
}

func (Config) createServeMux(logger log.Factory) (*http.ServeMux, http.Handler) {
	mux := http.NewServeMux()

	mux.Handle("/debug/pprof/", http.DefaultServeMux)

	health := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "OK")
	})
	mux.Handle("/health", health)
	mux.Handle("/healthz", health)

	chain := alice.New(
		middleware.Recovery(
			middleware.RecoveryLogger(logger),
		),
		middleware.RequestID,
	)

	return mux, chain.Then(mux)
}
