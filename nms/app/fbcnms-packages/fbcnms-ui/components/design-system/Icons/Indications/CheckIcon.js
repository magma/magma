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

const CheckIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M10.292 16.293a.996.996 0 01-1.41 0l-3.59-3.59a.996.996 0 111.41-1.41l2.88 2.88 6.88-6.88a.996.996 0 111.41 1.41l-7.58 7.59z" />
  </SvgIcon>
);

export default CheckIcon;
