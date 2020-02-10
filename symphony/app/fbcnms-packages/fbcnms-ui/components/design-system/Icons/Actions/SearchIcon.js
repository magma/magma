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

const SearchIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(3.41,3.4)">
      <path
        d="M2.575 2.075a7 7 0 0110.557 9.142l4.293 4.293a1 1 0 01-1.415 1.415l-4.292-4.293A7.002 7.002 0 012.575 2.075zM3.99 3.49a5 5 0 107.07 7.07 5 5 0 00-7.07-7.07z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default SearchIcon;
