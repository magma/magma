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

const BackArrowIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path
      d="M6.978.5L8.5 2.027 4.549 6H14v2H4.55l3.95 3.973L6.978 13.5 1.203 7.706a1 1 0 010-1.412L6.978.5z"
      fillRule="evenodd"
    />
  </SvgIcon>
);

export default BackArrowIcon;
