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

import Sequelize from 'sequelize';
import https from 'https';

import {API_HOST, apiCredentials} from '@fbcnms/platform-server/config';
import {Organization} from '@fbcnms/sequelize-models';

async function enableNetworks(
  organizationName: string,
  networkIDs: Array<string>,
) {
  const organization = await Organization.findOne({
    where: {
      name: Sequelize.where(
        Sequelize.fn('lower', Sequelize.col('name')),
        Sequelize.fn('lower', organizationName),
      ),
    },
  });
  if (organization) {
    await organization.update({networkIDs});
  }
}

function main() {
  const args = process.argv.slice(2);
  if (args.length !== 1) {
    console.log('Usage: enableAllNetworkIDsForOrganization.js <orgName>');
    process.exit(1);
  }

  const options = {
    hostname: API_HOST,
    path: '/magma/v1/networks',
    cert: apiCredentials().cert,
    key: apiCredentials().key,
  };

  https.get(options, res => {
    res.on('data', networkIDs => {
      enableNetworks(args[0], JSON.parse(networkIDs.toString()))
        .then(_res => {
          console.log('Success');
          process.exit();
        })
        .catch(err => {
          console.error(err);
          process.exit(1);
        });
    });
  });
}

main();
