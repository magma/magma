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

const ArrowRightIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M0 10.59L4.58 6 0 1.41 1.41 0l6 6-6 6z" fillRule="evenodd" />
  </SvgIcon>
);

export default ArrowRightIcon;
