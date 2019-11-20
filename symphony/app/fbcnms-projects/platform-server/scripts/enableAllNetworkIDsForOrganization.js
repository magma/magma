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

import https from 'https';
import {API_HOST, apiCredentials} from '@fbcnms/platform-server/config';
import {Organization} from '@fbcnms/sequelize-models';

async function enableNetworks(
  organizationName: string,
  networkIDs: Array<string>,
) {
  const organization = await Organization.findOne({
    where: {name: organizationName},
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
