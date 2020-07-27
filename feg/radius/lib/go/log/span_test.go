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
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type mockExporter struct {
	spans []*trace.SpanData
}

func (e *mockExporter) ExportSpan(s *trace.SpanData) {
	e.spans = append(e.spans, s)
}

func TestSpanCoreCheck(t *testing.T) {
	tests := []struct {
		name    string
		sampler trace.Sampler
		level   zapcore.Level
		expect  func(*testing.T, []*trace.SpanData)
	}{
		{
			name:    "suppressed-level",
			sampler: trace.AlwaysSample(),
			level:   zap.DebugLevel,
			expect: func(t *testing.T, spans []*trace.SpanData) {
				require.Len(t, spans, 1)
				assert.Empty(t, spans[0].Attributes)
				assert.Empty(t, spans[0].Annotations)
			},
		},
		{
			name:    "suppressed-sampler",
			sampler: trace.NeverSample(),
			level:   zap.ErrorLevel,
			expect: func(t *testing.T, spans []*trace.SpanData) {
				assert.Empty(t, spans)
			},
		},
		{
			name:    "emitted",
			sampler: trace.AlwaysSample(),
			level:   zap.InfoLevel,
			expect: func(t *testing.T, spans []*trace.SpanData) {
				require.Len(t, spans, 1)
				assert.Len(t, spans[0].Annotations, 1)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			exporter := &mockExporter{}
			trace.RegisterExporter(exporter)

			_, span := trace.StartSpan(context.Background(), "test",
				trace.WithSampler(tc.sampler))
			logger := zap.New(spanCore{span: span})
			if ce := logger.Check(tc.level, tc.level.String()+" message"); ce != nil {
				ce.Write()
			}
			span.End()

			tc.expect(t, exporter.spans)
		})
	}
}

func TestSpanCoreWith(t *testing.T) {
	exporter := &mockExporter{}
	trace.RegisterExporter(exporter)
	_, span := trace.StartSpan(context.Background(), "test",
		trace.WithSampler(trace.AlwaysSample()))

	root := zap.New(spanCore{span: span})
	root = root.With(zap.String("root", "root"))
	left := root.With(zap.String("left", "left"))
	right := root.With(zap.String("right", "right"))
	leaf := left.With(zap.String("leaf", "leaf"))

	loggers := []*zap.Logger{root, left, right, leaf}
	for _, logger := range loggers {
		logger.Info("")
	}

	span.End()
	spans := exporter.spans
	assert.Len(t, spans, 1)

	annotations := spans[0].Annotations
	assert.Len(t, annotations, len(loggers))

	assert.Len(t, annotations[0].Attributes, 2)
	assert.Equal(t, "root", annotations[0].Attributes["root"])

	assert.Len(t, annotations[1].Attributes, 3)
	assert.Equal(t, "root", annotations[1].Attributes["root"])
	assert.Equal(t, "left", annotations[1].Attributes["left"])

	assert.Len(t, annotations[2].Attributes, 3)
	assert.Equal(t, "root", annotations[2].Attributes["root"])
	assert.Equal(t, "right", annotations[2].Attributes["right"])

	assert.Len(t, annotations[3].Attributes, 4)
	assert.Equal(t, "root", annotations[3].Attributes["root"])
	assert.Equal(t, "left", annotations[3].Attributes["left"])
	assert.Equal(t, "leaf", annotations[3].Attributes["leaf"])
}

type loggable struct{ bool }

func (l loggable) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if !l.bool {
		return errors.New("can't marshal")
	}
	enc.AddString("loggable", "yes")
	return nil
}

func TestSpanCoreWrite(t *testing.T) {
	exporter := &mockExporter{}
	trace.RegisterExporter(exporter)

	_, span := trace.StartSpan(context.Background(), "test",
		trace.WithSampler(trace.AlwaysSample()))
	logger := zap.New(spanCore{span: span})
	logger.Info("field dump",
		zap.Bool("b", true),
		zap.Float32("f32", math.Pi),
		zap.Float64("f64", math.E),
		zap.Int("i", 0),
		zap.Int8("i8", -8),
		zap.Int16("i16", -16),
		zap.Int32("i32", -32),
		zap.Int64("i64", -64),
		zap.Uintptr("ptr", 0xbadbeef),
		zap.Uint("u", 0),
		zap.Uint8("u8", 8),
		zap.Uint16("u16", 16),
		zap.Uint32("u32", 32),
		zap.Uint64("u64", 64),
		zap.Complex64("c64", 1+1i),
		zap.Complex128("c128", 2+2i),
		zap.Duration("duration", time.Second+time.Second/2),
		zap.Time("date", time.Now()),
		zap.Binary("bin", []byte{5, 4, 3}),
		zap.ByteString("bytes", []byte{1, 2, 3}),
		zap.Reflect("numbers", []int{1, 2, 3}),
		zap.Bools("bools", []bool{true, true, false}),
		zap.Object("obj", loggable{true}),
	)
	span.End()

	spans := exporter.spans
	require.Len(t, spans, 1)
	annotations := spans[0].Annotations
	require.Len(t, annotations, 1)
	annotation := annotations[0]
	assert.Equal(t, true, annotation.Attributes["b"])
	assert.EqualValues(t, math.Float32bits(math.Pi), annotation.Attributes["f32"])
	assert.EqualValues(t, math.Float64bits(math.E), annotation.Attributes["f64"])
	assert.EqualValues(t, 0, annotation.Attributes["i"])
	assert.EqualValues(t, -8, annotation.Attributes["i8"])
	assert.EqualValues(t, -16, annotation.Attributes["i16"])
	assert.EqualValues(t, -32, annotation.Attributes["i32"])
	assert.EqualValues(t, -64, annotation.Attributes["i64"])
	assert.EqualValues(t, 0xbadbeef, annotation.Attributes["ptr"])
	assert.EqualValues(t, 0, annotation.Attributes["u"])
	assert.EqualValues(t, 8, annotation.Attributes["u8"])
	assert.EqualValues(t, 16, annotation.Attributes["u16"])
	assert.EqualValues(t, 32, annotation.Attributes["u32"])
	assert.EqualValues(t, 64, annotation.Attributes["u64"])
	assert.Equal(t, "(1+1i)", annotation.Attributes["c64"])
	assert.Equal(t, "(2+2i)", annotation.Attributes["c128"])
	assert.Equal(t, "1.5s", annotation.Attributes["duration"])
	assert.NotEmpty(t, annotation.Attributes["date"])
	assert.Equal(t, "BQQD", annotation.Attributes["bin"])
	assert.Equal(t, "\x01\x02\x03", annotation.Attributes["bytes"])
	assert.Equal(t, "([]int) (len=3) {\n (int) 1,\n (int) 2,\n (int) 3\n}\n", annotation.Attributes["numbers"])
	assert.Equal(t, "[true true false]", annotation.Attributes["bools"])
	assert.Equal(t, "(map[string]interface {}) (len=1) {\n (string) (len=8) \"loggable\": (string) (len=3) \"yes\"\n}\n", annotation.Attributes["obj"])
}

func TestSpanCoreOnPanic(t *testing.T) {
	exporter := &mockExporter{}
	trace.RegisterExporter(exporter)
	_, span := trace.StartSpan(context.Background(), "test",
		trace.WithSampler(trace.AlwaysSample()))

	logger := zap.New(spanCore{span: span}, zap.AddStacktrace(zap.PanicLevel))
	assert.Panics(t, func() { logger.Panic("oh no!") })
	span.End()

	spans := exporter.spans
	require.Len(t, spans, 1)
	assert.Equal(t, true, spans[0].Attributes["error"])
	annotations := spans[0].Annotations
	require.Len(t, annotations, 1)
	var annotation *trace.Annotation
	for i := range annotations {
		if annotations[i].Message == "oh no!" {
			annotation = &annotations[i]
			break
		}
	}
	require.NotNil(t, annotation)
	assert.Equal(t, "panic", annotation.Attributes["level"])
	assert.NotEmpty(t, annotation.Attributes["stack"])
}

func TestSpanCoreSync(t *testing.T) {
	core := spanCore{}
	assert.NoError(t, core.Sync())
}
