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
	"context"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// A Factory can create logger instances either for a given context or context-less.
	Factory interface {
		Bg() *zap.Logger
		For(context.Context) *zap.Logger
		With(...zap.Field) Factory
	}

	factory struct {
		background *zap.Logger
		contextual spanCore
	}

	nopFactory struct {
		*zap.Logger
	}
)

// NewFactory creates a new factory
func NewFactory(logger *zap.Logger) Factory {
	return factory{
		background: logger,
	}
}

// Bg returns a context-unaware logger
func (f factory) Bg() *zap.Logger {
	return f.background
}

// For returns a context-aware logger
func (f factory) For(ctx context.Context) *zap.Logger {
	logger := f.background
	if span := trace.FromContext(ctx); span != nil {
		logger = f.background.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, f.contextual.withSpan(span))
		}))
	}
	return logger.With(FieldsFromContext(ctx)...)
}

// With creates a child factory and optionally adds some context fields to that factory
func (f factory) With(fields ...zap.Field) Factory {
	if len(fields) > 0 {
		return factory{
			background: f.background.With(fields...),
			contextual: f.contextual.with(fields),
		}
	}
	return f
}

// NewNopFactory creates a new nop factory
func NewNopFactory() Factory {
	return nopFactory{zap.NewNop()}
}

func (f nopFactory) Bg() *zap.Logger                 { return f.Logger }
func (f nopFactory) For(context.Context) *zap.Logger { return f.Logger }
func (f nopFactory) With(...zap.Field) Factory       { return f }
