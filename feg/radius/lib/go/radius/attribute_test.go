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

package radius_test

import (
	"testing"

	"fbc/lib/go/radius"
)

func TestNewUserPassword_length(t *testing.T) {
	tbl := []struct {
		Password      string
		EncodedLength int
	}{
		{"", 16},
		{"abc", 16},
		{"0123456789abcde", 16},
		{"0123456789abcdef", 16},
		{"0123456789abcdef0", 16 * 2},
		{"0123456789abcdef0123456789abcdef0123456789abcdef", 16 * 3},
	}

	secret := []byte(`12345`)
	ra := []byte(`0123456789abcdef`)

	for _, x := range tbl {
		attr, err := radius.NewUserPassword([]byte(x.Password), secret, ra)
		if err != nil {
			t.Fatal(err)
		}
		if len(attr) != x.EncodedLength {
			t.Fatalf("expected encoded length of %#v = %d, got %d", x.Password, x.EncodedLength, len(attr))
		}
	}
}
