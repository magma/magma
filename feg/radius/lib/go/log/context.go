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

	"go.uber.org/zap"
)

type contextKey struct{}

// NewFieldsContext returns a new context with the given fields attached.
func NewFieldsContext(parent context.Context, fields ...zap.Field) context.Context {
	f := FieldsFromContext(parent)
	return context.WithValue(parent, contextKey{}, append(f[:len(f):len(f)], fields...))
}

// FieldsFromContext returns the fields stored in a context, or nil if there isn't one.
func FieldsFromContext(ctx context.Context) []zap.Field {
	f, _ := ctx.Value(contextKey{}).([]zap.Field)
	return f
}
