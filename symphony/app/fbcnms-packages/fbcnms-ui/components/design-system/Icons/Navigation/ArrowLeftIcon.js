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

const ArrowLeftIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M7.41 1.41L2.83 6l4.58 4.59L6 12 0 6l6-6z" fillRule="evenodd" />
  </SvgIcon>
);

export default ArrowLeftIcon;
