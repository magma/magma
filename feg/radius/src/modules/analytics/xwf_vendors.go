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

package analytics

// Vendor is the GraphQL enum for representing an AP vendor.
// It's used as a translation layer between meter's models and www's models.
type Vendor int

// Enum list of vendors supported by Express-Wifi program in its wider meaning (Wifi, Carrier Wifi, Standalone, ...).
const (
	Cambium Vendor = iota
	Ruckus
	Mojo
	CoovaChilli
	NonCertCambium
	ChilliSpot
	IPNet
	HP
)

// Vendors are the text representation of the Vendor enums.
var Vendors = [...]string{
	Cambium:        "CAMBIUM",
	Ruckus:         "RUCKUS",
	Mojo:           "MOJO",
	CoovaChilli:    "COOVACHILLI",
	NonCertCambium: "NON_CERT_CAMBIUM",
	ChilliSpot:     "CHILLISPOT",
	IPNet:          "IPNET",
	HP:             "HP",
}

// String implements the fmt.Stringer interface.
func (v Vendor) String() string { return Vendors[v] }
