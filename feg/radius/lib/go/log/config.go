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

package log

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config offers a declarative way to construct a logger.
type Config struct {
	// Level is the minimum enabled logging level.
	Level string `env:"LEVEL" long:"level" default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" description:"Only log messages with the given severity or above."`
	// Format sets the logging format. Valid values are "json" and "console".
	Format string `env:"FORMAT" long:"format" default:"console" choice:"console" choice:"json" description:"Output format of log messages."`
}

// Build constructs a logger from Config.
func (cfg Config) Build() (Factory, error) {
	if cfg == (Config{}) {
		return NewNopFactory(), nil
	}

	var c zap.Config
	switch cfg.Format {
	case "console":
		c = zap.NewDevelopmentConfig()
	case "json":
		c = zap.NewProductionConfig()
	default:
		return nil, errors.Errorf("unsupported logging format: %q", cfg.Format)
	}

	var level zapcore.Level
	if err := level.Set(cfg.Level); err != nil {
		return nil, errors.Wrap(err, "setting log level")
	}
	c.Level = zap.NewAtomicLevelAt(level)

	logger, err := c.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return nil, errors.Wrap(err, "creating logger")
	}
	return NewFactory(logger), nil
}
