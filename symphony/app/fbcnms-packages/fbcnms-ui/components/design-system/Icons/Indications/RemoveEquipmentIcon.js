/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {SvgIconStyleProps} from '../SvgIcon';

import React from 'react';
import SvgIcon from '../SvgIcon';

const RemoveEquipmentIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M20.9 7.343a1 1 0 010 1.414l-2.122 2.122.707.707a2 2 0 010 2.828l-2.828 2.829a7.002 7.002 0 01-7.302 1.645L8.171 20.07l-1.414-1.414-2.828 2.828-1.414-1.414 2.828-2.828-1.414-1.415 1.183-1.183a7.002 7.002 0 011.645-7.302l2.829-2.828a2 2 0 012.828 0l.708.706 2.12-2.12a1 1 0 011.328-.078l.087.078a1 1 0 01.078 1.326l-.078.088-2.12 2.12 2.827 2.83 2.121-2.122a1 1 0 011.414 0zM11 5.93L8.172 8.757a5.002 5.002 0 00-.978 5.692l.305.638-.741.742 1.414 1.414.742-.741.638.304a5.002 5.002 0 005.69-.978L18.073 13 11 5.929zm3.536 7.778l-1.415 1.414L8.88 10.88l1.414-1.415 4.243 4.243z" />
  </SvgIcon>
);

export default RemoveEquipmentIcon;
