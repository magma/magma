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
    <g transform="translate(2,3)">
      <path
        d="M18 16v-4h2v6H0v-6h2v4h16zM11 0v9.225L14.313 6l1.437 1.4-5.052 4.92a1 1 0 01-1.396 0L4.25 7.4 5.688 6 9 9.225V0h2z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default DownloadIcon;
