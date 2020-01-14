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

const NextArrowIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path
      d="M8 0L6.59 1.41 10.173 5H0v2h10.173L6.59 10.59 8 12l5.293-5.293a1 1 0 000-1.414L8 0z"
      fillRule="evenodd"
    />
  </SvgIcon>
);

export default NextArrowIcon;
