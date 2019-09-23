/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {MagmaAPIUrls} from '@fbcnms/magmalte/app/common/MagmaAPI';

import axios from 'axios';

export async function isWACNetwork(networkID: string) {
  try {
    const networkConfigs = await axios.get(
      MagmaAPIUrls.networkConfigsForType(networkID, 'wifi'),
    );
    return networkConfigs.data.additional_props?.wac_type === 'aruba';
  } catch (e) {}
  return false;
}
