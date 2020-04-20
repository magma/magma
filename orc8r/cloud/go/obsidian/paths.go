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

package obsidian

const (
	UrlSep = "/"

	MagmaNetworksUrlPart  = "networks"
	MagmaOperatorsUrlPart = "operators"

	// "/magma"
	RestRoot = UrlSep + "magma"
	// "/magma/networks"
	NetworksRoot = RestRoot + UrlSep + MagmaNetworksUrlPart
	// "/magma/operators"
	OperatorsRoot = RestRoot + UrlSep + MagmaOperatorsUrlPart

	// Supported API versions
	V0 = ""
	V1 = "v1"
	// Note the trailing slash (this is actually important for apidocs to render properly)
	V1Root = RestRoot + UrlSep + V1 + UrlSep
)
