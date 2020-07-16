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

const BackArrowIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M11.978 5.5L13.5 7.027 9.549 11H19v2H9.55l3.95 3.973-1.522 1.527-5.775-5.794a1 1 0 010-1.412L11.978 5.5z" />
  </SvgIcon>
);

export default BackArrowIcon;
