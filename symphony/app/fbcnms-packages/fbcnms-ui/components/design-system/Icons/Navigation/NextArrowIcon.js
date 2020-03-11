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

const NextArrowIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M13 6l-1.41 1.41L15.173 11H5v2h10.173l-3.583 3.59L13 18l5.293-5.293a1 1 0 000-1.414L13 6z" />
  </SvgIcon>
);

export default NextArrowIcon;
