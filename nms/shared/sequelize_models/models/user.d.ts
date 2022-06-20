/**
 * Copyright 2022 The Magma Authors.
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

// This is the type required for creation
type UserRawInitType = {
  email: string;
  organization?: string;
  password: string;
  role: number;
  networkIDs?: Array<string>;
};

// This is the type read back
export type UserRawType = {
  id: number;
  networkIDs: Array<string>;
  isSuperUser: boolean;
  isReadOnlyUser: boolean;
  role: number;
} & UserRawInitType;
