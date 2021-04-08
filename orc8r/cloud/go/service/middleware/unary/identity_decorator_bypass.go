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

package unary

// identityDecoratorBypassList is a map of RPC methods which are allowed to
// bypass Identity verification checks.
// For now, the only service in this category is Bootstrapper.
var identityDecoratorBypassList = map[string]struct{}{
	"/magma.orc8r.Bootstrapper/GetChallenge": {},
	"/magma.orc8r.Bootstrapper/RequestSign":  {},
}
