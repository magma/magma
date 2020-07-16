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
    <path d="M15.41 7.41L10.83 12l4.58 4.59L14 18l-6-6 6-6z" />
  </SvgIcon>
);

export default ArrowLeftIcon;
