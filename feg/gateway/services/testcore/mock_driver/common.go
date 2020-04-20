/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mock_driver

import (
	"fmt"

	"magma/feg/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/protobuf/ptypes/wrappers"
)

type CreditControlRequestPK struct {
	imsi        string
	requestType protos.CCRequestType
}

func NewCCRequestPK(imsi string, requestType protos.CCRequestType) CreditControlRequestPK {
	return CreditControlRequestPK{
		imsi:        imsi,
		requestType: requestType,
	}
}

func CompareRequestNumber(pk CreditControlRequestPK, expected *wrappers.Int32Value, actual datatype.Unsigned32) error {
	if expected == nil {
		return nil
	}
	expectedRN := expected.GetValue()
	if expectedRN != int32(actual) {
		return fmt.Errorf("For Request=%v, Expected Request Number: %v, Received Request Number: %v", pk, expectedRN, actual)
	}
	return nil
}

func (r CreditControlRequestPK) String() string {
	return fmt.Sprintf("Imsi: %v, Type: %v", r.imsi, r.requestType)
}

func EqualWithinDelta(a, b, delta uint64) bool {
	if b >= a && b-a <= delta {
		return true
	}
	if a >= b && a-b <= delta {
		return true
	}
	return false
}
