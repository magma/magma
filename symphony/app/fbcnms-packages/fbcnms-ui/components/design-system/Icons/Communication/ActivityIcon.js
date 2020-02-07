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

const ActivityIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(2,6)">
      <path
        d="M10 10v2H0v-2h10zm10-5v2H0V5h20zm0-5v2H0V0h20z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default ActivityIcon;
