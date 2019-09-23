/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {NetworkContextType} from '../context/NetworkContext';
import type {NetworkType} from '@fbcnms/types/network';
import type {SectionsConfigs} from '../layout/Section';

import AppContext from '@fbcnms/ui/context/AppContext';
import NetworkContext from '../context/NetworkContext';
import axios from 'axios';
import {useContext, useEffect, useState} from 'react';

import {CELLULAR, RHINO} from '@fbcnms/types/network';
import {MagmaAPIUrls} from '../../common/MagmaAPI';
import {getLteSections} from '../lte/LteSections';
import {getRhinoSections} from '../rhino/RhinoSections';
import {useFeatureFlag} from '@fbcnms/ui/hooks';

export default function useSections(): SectionsConfigs {
  const {networkId} = useContext<NetworkContextType>(NetworkContext);
  const [networkType, setNetworkType] = useState<?NetworkType>(null);
  const alertsEnabled = useFeatureFlag(AppContext, 'alerts');

  useEffect(() => {
    if (networkId) {
      axios
        .get(MagmaAPIUrls.network(networkId))
        .then(({data}) => setNetworkType(data?.features?.networkType || ''));
    }
  }, [networkId]);

  if (!networkId || networkType === null) {
    return [null, []];
  }

  switch (networkType) {
    case RHINO:
      return getRhinoSections();
    case CELLULAR:
    default:
      return getLteSections(alertsEnabled);
  }
}
