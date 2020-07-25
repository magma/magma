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

import "magma/lte/cloud/go/protos"

// validateSubscriberData ensures that a subscriber data proto is not nil and
// that it contains a valid subscriber id.
func validateSubscriberData(subscriber *protos.SubscriberData) error {
	if subscriber == nil {
		return NewInvalidArgumentError("Subscriber data cannot be nil")
	}
	if subscriber.Sid == nil {
		return NewInvalidArgumentError("Subscriber data must contain a subscriber id")
	}
	return validateSubscriberID(subscriber.Sid.Id)
}

// validateSubscriberID ensures that a subscriber ID can be stored
func validateSubscriberID(id string) error {
	if len(id) == 0 {
		return NewInvalidArgumentError("Subscriber id cannot be the empty string")
	}
	return nil
}
