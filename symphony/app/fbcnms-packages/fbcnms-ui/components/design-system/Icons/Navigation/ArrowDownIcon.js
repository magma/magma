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

const ArrowDownIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M1.41 0L6 4.58 10.59 0 12 1.41l-6 6-6-6z" fillRule="evenodd" />
  </SvgIcon>
);

export default ArrowDownIcon;
