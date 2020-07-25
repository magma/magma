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

package storage

import (
	"magma/orc8r/cloud/go/services/certifier/protos"
)

// CertifierStorage provides storage functionality for mapping serial numbers to certificate information.
type CertifierStorage interface {
	// ListSerialNumbers returns all tracked serial numbers.
	ListSerialNumbers() ([]string, error)

	// GetCertInfo returns the certificate info associated with the serial number.
	// If not found, returns ErrNotFound from magma/orc8r/lib/go/errors.
	GetCertInfo(serialNumber string) (*protos.CertificateInfo, error)

	// GetManyCertInfo maps the passed serial numbers to their associated certificate info.
	GetManyCertInfo(serialNumbers []string) (map[string]*protos.CertificateInfo, error)

	// GetAllCertInfo returns a map of all serial numbers to their associated certificate info.
	GetAllCertInfo() (map[string]*protos.CertificateInfo, error)

	// PutCertInfo associates certificate info with the passed serial number.
	PutCertInfo(serialNumber string, certInfo *protos.CertificateInfo) error

	// DeleteCertInfo removes the serial number and its certificate info.
	// Returns success even when nothing is deleted (i.e. serial number not found).
	DeleteCertInfo(serialNumber string) error
}
