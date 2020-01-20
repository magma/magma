/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Configure from '../network/Configure';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import PoliciesConfig from '../network/PoliciesConfig';
import React from 'react';
import UpgradeConfig from '../network/UpgradeConfig';

import useMagmaAPI from '../../common/useMagmaAPI';
import {useRouter} from '@fbcnms/ui/hooks';

export default function CWFConfigure() {
  const tabs = [
    {
      component: UpgradeConfig,
      label: 'Upgrades',
      path: 'upgrades',
    },
    {
      component: CWFPolicies,
      label: 'Policies',
      path: 'policies',
    },
  ];
  return <Configure tabRoutes={tabs} />;
}

function CWFPolicies() {
  const {match} = useRouter();

  const {response, isLoading} = useMagmaAPI(MagmaV1API.getCwfByNetworkId, {
    networkId: match.params.networkId,
  });

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <PoliciesConfig mirrorNetwork={response?.federation?.feg_network_id} />
  );
}
