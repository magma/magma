/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import CellWifiIcon from '@material-ui/icons/CellWifi';
import KPITray from '../../components/KPITray';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {useRouter} from '@fbcnms/ui/hooks';

export default function NetworkKPI() {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const {response: lteGatwayResp, isLoading: isLteRespLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdGateways,
    {
      networkId: networkId,
    },
  );

  const {response: enb, isLoading: isEnbRespLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdEnodebs,
    {
      networkId: networkId,
    },
  );

  const {response: policyRules, isLoading: isPolicyLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPoliciesRules,
    {
      networkId: networkId,
    },
  );

  const {response: subscriber, isLoading: isSubscriberLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdSubscribers,
    {
      networkId: networkId,
    },
  );

  const {response: apns, isLoading: isAPNsLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdApns,
    {
      networkId: networkId,
    },
  );
  if (
    isLteRespLoading ||
    isEnbRespLoading ||
    isPolicyLoading ||
    isSubscriberLoading ||
    isAPNsLoading
  ) {
    return <LoadingFiller />;
  }
  return (
    <KPITray
      data={[
        {
          icon: CellWifiIcon,
          category: 'Gateways',
          value: lteGatwayResp ? Object.keys(lteGatwayResp).length : 0,
        },
        {
          icon: SettingsInputAntennaIcon,
          category: 'eNodeBs',
          value: enb ? Object.keys(enb).length : 0,
        },
        {
          icon: PeopleIcon,
          category: 'Subscribers',
          value: subscriber ? Object.keys(subscriber).length : 0,
        },
        {
          icon: LibraryBooksIcon,
          category: 'Policies',
          value: policyRules ? policyRules.length : 0,
        },
        {
          icon: RssFeedIcon,
          category: 'APNs',
          value: apns ? Object.keys(apns).length : 0,
        },
      ]}
    />
  );
}
