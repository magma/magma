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

const ArrowRightIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M8 16.59L12.58 12 8 7.41 9.41 6l6 6-6 6z" />
  </SvgIcon>
);

export default ArrowRightIcon;
