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

package protos

import (
	"github.com/pkg/errors"
)

func (m *GetMSISDNsRequest) Validate() error {
	if m.NetworkId == "" {
		return errors.New("network ID cannot be empty")
	}
	return nil
}

func (m *SetMSISDNRequest) Validate() error {
	if m.NetworkId == "" {
		return errors.New("network ID cannot be empty")
	}
	if m.Msisdn == "" {
		return errors.New("msisdn cannot be empty")
	}
	if m.Imsi == "" {
		return errors.New("imsi cannot be empty")
	}
	return nil
}

func (m *DeleteMSISDNRequest) Validate() error {
	if m.NetworkId == "" {
		return errors.New("network ID cannot be empty")
	}
	if m.Msisdn == "" {
		return errors.New("msisdn cannot be empty")
	}
	return nil
}

func (m *GetIPsRequest) Validate() error {
	if m.NetworkId == "" {
		return errors.New("network ID cannot be empty")
	}
	return nil
}

func (m *SetIPsRequest) Validate() error {
	if m.NetworkId == "" {
		return errors.New("network ID cannot be empty")
	}
	for _, mapping := range m.IpMappings {
		if mapping.Ip == "" {
			return errors.Errorf("ip cannot be empty in mapping %v", mapping)
		}
		if mapping.Imsi == "" {
			return errors.Errorf("imsi cannot be empty in mapping %v", mapping)
		}
		if mapping.Apn == "" {
			return errors.Errorf("apn cannot be empty in mapping %v", mapping)
		}
	}
	return nil
}
