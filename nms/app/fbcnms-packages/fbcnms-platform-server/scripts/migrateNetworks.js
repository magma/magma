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
 * @flow strict-local
 * @format
 */
'use strict';

import MagmaV1API from '../magma';
import axios from 'axios';
import https from 'https';
import nullthrows from '@fbcnms/util/nullthrows';
import {API_HOST, apiCredentials} from '@fbcnms/platform-server/config';

const httpsAgent = new https.Agent({
  cert: apiCredentials().cert,
  key: apiCredentials().key,
  rejectUnauthorized: false,
});

function apiUrl(): string {
  return !/^https?\:\/\//.test(nullthrows(API_HOST))
    ? `https://${nullthrows(API_HOST)}/magma`
    : `${nullthrows(API_HOST)}/magma`;
}

async function main() {
  const networks = await MagmaV1API.getNetworks();
  await Promise.all(
    networks.map(async networkId => {
      try {
        const network = await MagmaV1API.getNetworksByNetworkId({networkId});
        const {data} = await axios({
          baseURL: apiUrl(),
          url: `/networks/${networkId}`,
          method: 'GET',
          httpsAgent,
        });

        const networkType = data.features?.networkType;
        if (networkType && networkType !== network.type) {
          console.log(
            `Mismatch \t${networkType}\t${network.type || 'null'}\t`,
            networkId,
          );

          // uncomment this to write the changes
          // const newType = networkType === 'cellular' ? 'lte' : networkType;
          // await MagmaV1API.putNetworksByNetworkIdType({
          //   networkId,
          //   type: JSON.stringify(`${newType}`),
          // });
        }
      } catch (e) {
        console.log('Error: ', networkId, e);
      }
    }),
  );
}

main();
