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
	"errors"
)

func (m *GetHostnameForHWIDRequest) Validate() error {
	if m == nil {
		return errors.New("request cannot be nil")
	}
	if m.Hwid == "" {
		return errors.New("request hwid cannot be empty")
	}
	return nil
}

func (m *MapHWIDToHostnameRequest) Validate() error {
	if m == nil {
		return errors.New("request cannot be nil")
	}
	if m.HwidToHostname == nil {
		return errors.New("request hwidToHostname cannot be empty")
	}
	return nil
}

func (m *GetIMSIForSessionIDRequest) Validate() error {
	if m == nil {
		return errors.New("request cannot be nil")
	}
	if m.SessionID == "" {
		return errors.New("request sessionID cannot be empty")
	}
	return nil
}

func (m *MapSessionIDToIMSIRequest) Validate() error {
	if m == nil {
		return errors.New("request cannot be nil")
	}
	if m.NetworkID == "" {
		return errors.New("network ID cannot be empty")
	}
	if m.SessionIDToIMSI == nil {
		return errors.New("request sessionIDToIMSI cannot be empty")
	}
	return nil
}

func (m *GetHWIDForSgwCTeidRequest) Validate() error {
	if m == nil {
		return errors.New("request cannot be nil")
	}
	if m.Teid == "" {
		return errors.New("request teid cannot be empty")
	}
	return nil
}

func (m *MapSgwCTeidToHWIDRequest) Validate() error {
	if m == nil {
		return errors.New("request cannot be nil")
	}
	if m.NetworkID == "" {
		return errors.New("network ID cannot be empty")
	}
	if m.TeidToHwid == nil {
		return errors.New("request TeidToHwid cannot be empty")
	}
	return nil
}
