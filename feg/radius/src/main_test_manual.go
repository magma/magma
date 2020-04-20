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

package main

import (
	"context"
	"log"
	"testing"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
)

// TestManual Manually test server
//nolint:deadcode
func TestManual(_ *testing.T) {
	// Act
	packet := radius.New(radius.CodeAccessRequest, []byte(`123456`))
	rfc2865.UserName_SetString(packet, "tim")
	rfc2865.UserPassword_SetString(packet, "12345")
	response, err := radius.Exchange(context.Background(), packet, "localhost:1812")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Code:", response.Code)
}
