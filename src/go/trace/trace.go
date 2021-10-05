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

package trace

import (
	"time"
)

//go:generate go run github.com/golang/mock/mockgen -destination mock_trace/mock_trace.go . Trace

// TraceID identifies a trace.
type TraceID [16]byte

// SpanID identifies a span.
type SpanID [8]byte

// SpanStatus is the status of a span.
type SpanStatus uint8

// Sampled signifies a sampling decision.
type Sampled int8

// A Span is the building block of a Sentry transaction. Spans build up a tree
// structure of timed operations. The span tree makes up a transaction event
// that is sent to Sentry when the root span is finished.
//
// Spans must be started with either StartSpan or Span.StartChild.
type Span struct {
	TraceID      TraceID
	SpanID       SpanID
	ParentSpanID SpanID
	Op           string
	Description  string
	Status       SpanStatus
	Tags         map[string]string
	StartTime    time.Time
	EndTime      time.Time
	Data         map[string]interface{}

	Sampled Sampled
}

type Trace interface {
	Add(operation string) Span
	SetTag(name, value string)
	Finish()
}
