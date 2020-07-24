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
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggerFieldContext(t *testing.T) {
	core, o := observer.New(zap.InfoLevel)
	logger := NewFactory(zap.New(core))

	ctx := NewFieldsContext(context.Background(), zap.String("name", "test"))
	ctx = NewFieldsContext(ctx, zap.String("lang", "go"))
	logger.For(ctx).Info("test message", zap.Int("speed", 42))

	assert.Equal(t, 1, o.
		FilterMessage("test message").
		FilterField(zap.String("name", "test")).
		FilterField(zap.String("lang", "go")).
		Len(),
	)
}
