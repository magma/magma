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

const NumbersIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M16.959 2.79a1 1 0 01.707 1.224L16.597 8 19 8a1 1 0 010 2h-2.939l-1.072 4H18a1 1 0 010 2h-3.547l-1.187 4.435a1 1 0 11-1.932-.518L12.384 16H9.195l-1.188 4.436a1 1 0 11-1.932-.518L7.125 16 5 16a1 1 0 010-2h2.661l1.072-4H6a1 1 0 110-2h3.268l1.207-4.503a1 1 0 011.932.517L11.34 8h3.188l1.207-4.502a1 1 0 011.225-.707zM13.992 10h-3.189l-1.072 4h3.189l1.072-4z" />
  </SvgIcon>
);

export default NumbersIcon;
