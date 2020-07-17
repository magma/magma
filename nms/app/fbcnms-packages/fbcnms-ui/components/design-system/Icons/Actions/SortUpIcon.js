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

const SortUpIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g>
      <path
        fill="#9DA9BE"
        d="M8.94 14L12 17.09 15.06 14l.94.951L12 19l-4-4.049z"
      />
      <path d="M15.06 10L12 6.91 8.94 10 8 9.049 12 5l4 4.049z" />
    </g>
  </SvgIcon>
);

export default SortUpIcon;
