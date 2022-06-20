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

import {Strategy} from 'passport';

//use this in place of the real openid strategies until discovery finishes
export default class StubStrategy extends Strategy {
  constructor() {
    super();
    this.name = 'stub';
  }
  authenticate() {
    return this.fail('No implementation found for strategy');
  }
}
