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

package packet

// Code as defined in RFC3748 section 4
type Code int

// Code values
const (
	CodeREQUEST  Code = 1
	CodeRESPONSE Code = 2
	CodeSUCCESS  Code = 3
	CodeFAILURE  Code = 4
)

// IsValid Verify if the value is a valid Code
// (may be coming from external source like incoming EAP packet)
func (c Code) IsValid() bool {
	switch c {
	case
		CodeREQUEST,
		CodeRESPONSE,
		CodeSUCCESS,
		CodeFAILURE:
		return true
	}
	return false
}

// IsRequestOrResponse helper method
func (c Code) IsRequestOrResponse() bool {
	return (c == CodeREQUEST) || (c == CodeRESPONSE)
}
