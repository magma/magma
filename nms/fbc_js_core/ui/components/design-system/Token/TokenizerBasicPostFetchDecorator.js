/**
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
 *
 * @flow
 * @format
 */

import type {Entries} from './Tokenizer';

export default function TokenizerBasicPostFetchDecorator<TEntry>(
  response: Entries<TEntry>,
  queryString: string,
  currentTokens: Entries<TEntry>,
): Entries<TEntry> {
  return response.filter(
    entry =>
      entry.label.toLowerCase().includes(queryString.toLowerCase()) &&
      !currentTokens.some(token => token.key === entry.key),
  );
}
