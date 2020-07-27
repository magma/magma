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

package gy_test

import (
	"testing"

	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestRedirectServer_ToProto(t *testing.T) {
	var convertedRedirectServer *protos.RedirectServer = nil
	convertedRedirectServer = (&gy.RedirectServer{
		RedirectAddressType:   gy.IPV4Address,
		RedirectServerAddress: "www.magma.com",
	}).ToProto()

	assert.Equal(t, protos.RedirectServer_IPV4, convertedRedirectServer.RedirectAddressType)
	assert.Equal(t, "www.magma.com", convertedRedirectServer.RedirectServerAddress)

	var nilRedirectServer *gy.RedirectServer = nil
	convertedRedirectServer = nilRedirectServer.ToProto()

	assert.Equal(t, &protos.RedirectServer{}, convertedRedirectServer)
}
