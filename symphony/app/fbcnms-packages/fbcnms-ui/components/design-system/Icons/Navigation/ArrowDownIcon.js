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
    <path d="M7.41 8L12 12.58 16.59 8 18 9.41l-6 6-6-6z" />
  </SvgIcon>
);

export default ArrowDownIcon;
