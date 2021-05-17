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

package models

import (
	"github.com/hashicorp/go-multierror"

	"magma/orc8r/cloud/go/obsidian/models"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	lteAuthKeyLength = 16
	lteAuthOpcLength = 16
)

func (m *LteSubscription) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	authKeyLen := len([]byte(m.AuthKey))
	if authKeyLen != lteAuthKeyLength {
		return models.ValidateErrorf("expected lte auth key to be %d bytes but got %d bytes", lteAuthKeyLength, authKeyLen)
	}

	// OPc is optional, but if it's provided it should be 16 bytes
	authOpcLen := len([]byte(m.AuthOpc))
	if authOpcLen > 0 && authOpcLen != lteAuthOpcLength {
		return models.ValidateErrorf("expected lte auth opc to be %d bytes but got %d bytes", lteAuthOpcLength, authOpcLen)
	}

	return nil
}

func (m *MutableSubscriber) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if err := m.Lte.ValidateModel(); err != nil {
		return err
	}

	// You can't assign a static IP allocation if the subscriber doesn't have
	// the APN active
	apnSet := funk.Map(m.ActiveApns, func(apn string) (string, bool) { return apn, true }).(map[string]bool)
	for apn := range m.StaticIps {
		if _, exists := apnSet[apn]; !exists {
			return errors.Errorf("static IP assigned to APN %s which is not active for the subscriber", apn)
		}
	}
	return nil
}

func (m MutableSubscribers) ValidateModel() error {
	errs := &multierror.Error{}
	for _, s := range m {
		if err := s.ValidateModel(); err != nil {
			multierror.Append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}

func (m *IcmpStatus) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *MsisdnAssignment) ValidateModel() error {
	return m.Validate(strfmt.Default)
}
