// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/magma/magma/src/go/trace"
)

type Trace struct {
	sentry.Span
}

func NewTrace(ctx context.Context, operation string) trace.Trace {
	span := sentry.StartSpan(ctx, operation)

	return &Trace{
		Span: *span,
	}
}

func (t *Trace) Add(operation string) trace.Span {
	return convertSpan(*t.Span.StartChild(operation))
}

func (t *Trace) SetTag(name, value string) {
	t.Span.SetTag(name, value)
}

func (t *Trace) Finish() {
	t.Span.Finish()
}

func convertSpan(sentrySpan sentry.Span) trace.Span {
	return trace.Span{
		TraceID:      trace.TraceID(sentrySpan.TraceID),
		SpanID:       trace.SpanID(sentrySpan.SpanID),
		ParentSpanID: trace.SpanID(sentrySpan.ParentSpanID),
		Op:           sentrySpan.Op,
		Description:  sentrySpan.Description,
		Status:       trace.SpanStatus(sentrySpan.Status),
		Tags:         sentrySpan.Tags,
		StartTime:    sentrySpan.StartTime,
		EndTime:      sentrySpan.EndTime,
		Data:         sentrySpan.Data,
		Sampled:      trace.Sampled(sentrySpan.Sampled),
	}
}
