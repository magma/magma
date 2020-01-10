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

const FiltersIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(3,4)">
      <path
        d="M6 8a4.002 4.002 0 013.874 3H18v2H9.874a4.002 4.002 0 01-7.748 0H0v-2h2.126C2.57 9.275 4.136 8 6 8zm0 2a2 2 0 100 4 2 2 0 000-4zm6-10a4.002 4.002 0 013.874 3H18v2h-2.126a4.002 4.002 0 01-7.748 0H0V3h8.126C8.57 1.275 10.136 0 12 0zm0 2a2 2 0 100 4 2 2 0 000-4z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default FiltersIcon;
