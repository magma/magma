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
	"magma/lte/cloud/go/protos"
)

// SubscriberStore is an interface for storing and accessing subscriber data.
type SubscriberStore interface {
	// AddSubscriber tries to add this subscriber to the server.
	// This function returns an AlreadyExists error if the subscriber has already
	// been added.
	// Input: The subscriber data which will be added.
	AddSubscriber(data *protos.SubscriberData) error

	// GetSubscriberData looks up a subscriber by their Id.
	// If the subscriber cannot be found, an error is returned instead.
	// Input: The id of the subscriber to be looked up.
	// Output: The data of the corresponding subscriber or an error.
	GetSubscriberData(id string) (*protos.SubscriberData, error)

	// UpdateSubscriber changes the data stored for an existing subscriber.
	// If the subscriber cannot be found, an error is returned instead.
	// Input: The new subscriber data to store
	UpdateSubscriber(data *protos.SubscriberData) error

	// DeleteSubscriber deletes a subscriber by their Id.
	// If the subscriber is not found, then this call is ignored.
	// Input: The id of the subscriber to be deleted.
	DeleteSubscriber(id string) error

	// DeleteAllSubscribers deletes all the data from the store.
	DeleteAllSubscribers() error
}
