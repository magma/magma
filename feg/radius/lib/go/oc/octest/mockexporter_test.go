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

package octest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
)

func TestMockExporter(t *testing.T) {
	exporter := &MockExporter{}
	trace.RegisterExporter(exporter)
	defer trace.UnregisterExporter(exporter)

	_, span := trace.StartSpan(context.Background(), "test",
		trace.WithSampler(trace.AlwaysSample()))

	span.AddAttributes([]trace.Attribute{
		trace.StringAttribute("message", "hello"),
		trace.Int64Attribute("rank", 42),
	}...)
	span.End()

	spans := exporter.ExportedSpans()
	require.Len(t, spans, 1)
	assert.Equal(t, "test", spans[0].Name)
	assert.Equal(t, "hello", spans[0].Attributes["message"])
	assert.EqualValues(t, 42, spans[0].Attributes["rank"])

	exporter.Reset()
	assert.Empty(t, exporter.ExportedSpans())
}
