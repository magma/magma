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

const DownloadIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M20 19v-4h2v6H2v-6h2v4h16zM13 3v9.225L16.313 9l1.437 1.4-5.052 4.92a1 1 0 01-1.396 0L6.25 10.4 7.688 9 11 12.225V3h2z" />
  </SvgIcon>
);

export default DownloadIcon;
