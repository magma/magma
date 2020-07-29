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

package logtest

import (
	"context"
	"testing"

	"fbc/lib/go/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type testFactory struct {
	*zap.Logger
}

// NewFactory creates a new testing logger factory
func NewFactory(t *testing.T) log.Factory {
	return testFactory{zaptest.NewLogger(t)}
}

func (f testFactory) Bg() *zap.Logger                 { return f.Logger }
func (f testFactory) For(context.Context) *zap.Logger { return f.Logger }

func (f testFactory) With(fields ...zap.Field) log.Factory {
	return testFactory{f.Logger.With(fields...)}
}
