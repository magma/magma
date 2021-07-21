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

package subscriberdb

type Config struct {
	// DigestsEnabled is a feature flag for the flat digest functionality.
	DigestsEnabled bool `yaml:"digestsEnabled"`
	// ChangesetSizeThreshold specifies the max size of the cloud-AGW changeset
	// past which a resync signal will be sent back to the AGW.
	ChangesetSizeThreshold int `yaml:"changesetSizeThreshold"`
	// MaxProtosLoadSize specifies the max size of cached subscriber protos that
	// can be loaded for a page.
	MaxProtosLoadSize uint64 `yaml:"maxProtosLoadSize"`
}
