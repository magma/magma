/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
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
