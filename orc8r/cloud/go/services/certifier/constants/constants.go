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

package constants

const (
	// CertInfoType is the type of CertInfo used in blobstore type fields.
	CertInfoType = "certificate_info"

	// UserType is the type of CertInfo used in blobstore type fields.
	UserType = "user"

	// PolicyType is the type of policy used in blobstore type fileds
	PolicyType = "policy"
)

type ResourceType string

const (
	Path      ResourceType = "path"
	NetworkID ResourceType = "network_id"
	TenantID  ResourceType = "tenant_id"
)
