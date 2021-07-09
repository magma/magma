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

package subscriberdb

const (
	ServiceName = "subscriberdb"

	EntityType = "subscriber"

	LookupTableBlobstore       = "subscriber_lookup_blobstore"
	PerSubDigestTableBlobstore = "per_sub_digest_blobstore"

	// MinimumSyncInterval is the the minimum interval in seconds between
	// gateway requests to sync its subscriberdb with the cloud.
	// orc8r should never send a value lower than MinimumSyncInterval.
	MinimumSyncInterval = 60
)
