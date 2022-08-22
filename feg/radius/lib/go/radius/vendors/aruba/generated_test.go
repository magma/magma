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

package aruba_test

import (
	"testing"

	"fbc/lib/go/radius"
	. "fbc/lib/go/radius/vendors/aruba"
)

func TestLookup(t *testing.T) {
	p := radius.New(radius.CodeAccessRequest, []byte(`12345`))
	ArubaUserRole_SetString(p, "Admin")
	ArubaDeviceType_SetString(p, "Desktop")

	if dt := ArubaDeviceType_GetString(p); dt != "Desktop" {
		t.Fatalf("ArubaDeviceType = %v; expecting %v", dt, "Desktop")
	}
}
