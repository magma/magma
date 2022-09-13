/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import "magma/dp/cloud/go/services/dp/storage/db"

// AmcManager is supposed to be a library that will replace radio controller
// it is not supposed to be a service but rather an interface to database
// could be implemented in this file as separate struct or combined with cbsd manager
// also its methods are supposed to be used in transaction (they should start a new one)
type AmcManager interface {
	// GetState is equivalent to GetState grpc method
	// it should return list of all feasible cbsd with grants
	// cbsd is considered feasible if and only if
	// - it has no pending requests
	// - one of the following conditions is satisfied
	//	 - it has all necessary parameters to perform sas requests (registration/grant)
	//   - it has some pending db action (e.g. it needs to be deleted)
	GetState() []*DetailedCbsd
	// CreateRequest should just store given request in the database (no data processing/marshaling)
	CreateRequest(*DBRequest)
	// DeleteCbsd should just delete cbsd (no need to check if it exists)
	DeleteCbsd(*DBCbsd)
	// UpdateCbsd should replace AcknowledgeCbsdUpdate, AcknowledgeCbsdRelinquish
	// and StoreAvailableFrequencies
	// it should just update cbsd (no need to lock)
	UpdateCbsd(*DBCbsd, db.FieldMask)
}
