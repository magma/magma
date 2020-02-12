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
    <path d="M10.59 7.41L6 2.83 1.41 7.41 0 6l6-6 6 6z" fillRule="evenodd" />
  </SvgIcon>
);

export default ArrowUpIcon;
