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

package alwaysaccept

import (
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

func TestAccessRequest(t *testing.T) {
	// Arrange
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	ctx, err := Init(logger, modules.ModuleConfig{})
	require.NoError(t, err, "failed to init")

	// Act
	res, err := Handle(
		ctx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: nil,
		},
		&radius.Request{Packet: &radius.Packet{
			Code: radius.CodeAccessRequest,
		}},
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			require.Fail(t, "next method is called bu not expected to")
			return nil, nil
		},
	)

	// Act and Assert
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, radius.CodeAccessAccept, res.Code)
}

func TestNotAccessRequest(t *testing.T) {
	// Arrange
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	ctx, err := Init(logger, modules.ModuleConfig{})
	require.NoError(t, err, "failed to init")

	// Act
	res, err := Handle(
		ctx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: nil,
		},
		&radius.Request{Packet: &radius.Packet{
			Code: radius.CodeAccountingRequest,
		}},
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			require.Fail(t, "next method is called bu not expected to")
			return nil, nil
		},
	)

	// Act and Assert
	require.NotNil(t, err)
	require.Nil(t, res)
}
