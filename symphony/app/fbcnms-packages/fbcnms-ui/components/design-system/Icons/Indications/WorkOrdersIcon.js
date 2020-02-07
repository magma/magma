/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {SvgIconStyleProps} from '../SvgIcon';

import React from 'react';
import SvgIcon from '../SvgIcon';

const WorkOrdersIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(2,2)">
      <path
        d="M12 4V2H8v2h4zM2 6v12h16V6H2zm16-2c1.11 0 2 .89 2 2v12c0 1.11-.89 2-2 2H2c-1.11 0-2-.89-2-2L.01 6C.01 4.89.89 4 2 4h4V2c0-1.11.89-2 2-2h4c1.11 0 2 .89 2 2v2h4z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default WorkOrdersIcon;
