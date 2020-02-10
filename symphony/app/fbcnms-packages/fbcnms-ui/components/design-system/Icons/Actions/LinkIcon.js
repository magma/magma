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

const LinkIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(3.456,3.456)">
      <path
        d="M9 1.929L6.172 4.757l1.414 1.415 2.828-2.829a3.009 3.009 0 014.243 0 3.009 3.009 0 010 4.243l-2.829 2.828 1.415 1.414L16.07 9a5.002 5.002 0 000-7.071 5.002 5.002 0 00-7.071 0zm1.414 9.9l-2.828 2.828a3.009 3.009 0 01-4.243 0 3.009 3.009 0 010-4.243l2.829-2.828-1.415-1.414L1.93 9a5.002 5.002 0 000 7.071 5.002 5.002 0 007.071 0l2.828-2.828-1.414-1.415zm-4.95-.708l5.657-5.657 1.415 1.415-5.657 5.657-1.415-1.415z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default LinkIcon;
