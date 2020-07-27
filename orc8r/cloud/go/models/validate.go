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
	"github.com/go-openapi/strfmt"
)

func (m *NetworkName) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkType) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkDescription) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayName) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayDescription) ValidateModel() error {
	return m.Validate(strfmt.Default)
}
