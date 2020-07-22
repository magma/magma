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

const ArrowUpIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M16.59 15.41L12 10.83l-4.59 4.58L6 14l6-6 6 6z" />
  </SvgIcon>
);

export default ArrowUpIcon;
