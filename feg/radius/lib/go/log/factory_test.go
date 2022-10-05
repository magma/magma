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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestFactory(t *testing.T) {
	core, o := observer.New(zap.InfoLevel)
	logger := NewFactory(zap.New(core))

	msg := "context-less"
	logger.Bg().Info(msg)
	assert.Len(t, o.TakeAll(), 1)

	ctx := context.Background()
	logger.With().For(ctx).Warn(msg)
	assert.Equal(t, 1, o.FilterMessage(msg).Len())

	exporter := &mockExporter{}
	trace.RegisterExporter(exporter)

	ctx, span := trace.StartSpan(context.Background(), "test",
		trace.WithSampler(trace.AlwaysSample()))
	field, msg := zap.Int("result", 42), "context-aware"
	logger.With(field).For(ctx).Info(msg)
	span.End()

	assert.Equal(t, 1, o.FilterField(field).FilterMessage(msg).Len())
	spans := exporter.spans
	assert.Len(t, spans, 1)
	annotations := spans[0].Annotations
	assert.Len(t, annotations, 1)
	assert.Equal(t, "context-aware", annotations[0].Message)
	assert.EqualValues(t, 42, annotations[0].Attributes["result"])
}

func TestNopFactory(t *testing.T) {
	logger := NewNopFactory()
	assert.EqualValues(t, zap.NewNop(), logger.Bg())
	assert.EqualValues(t, zap.NewNop(), logger.For(context.Background()))
	assert.EqualValues(t, logger, logger.With(zap.Int("id", 42)))
}
