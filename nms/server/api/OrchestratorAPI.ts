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

import https from 'https';
import nullthrows from '../../shared/util/nullthrows';
import {API_HOST, apiCredentials} from '../../config/config';
import {setUpApi} from '../../api/API';

const httpsAgent = new https.Agent({
  cert: apiCredentials().cert,
  key: apiCredentials().key,
  rejectUnauthorized: false,
});

const host = nullthrows(API_HOST);

const orchestratorUrl = !/^https?\:\/\//.test(host)
  ? `https://${host}/magma/v1`
  : `${host}/magma/v1`;

export default setUpApi(orchestratorUrl, httpsAgent);
