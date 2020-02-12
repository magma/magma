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

const UploadIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(2,3.2)">
      <path
        d="M18 16v-4h2v6H0v-6h2v4h16zm-7-3V3.775L14.313 7l1.437-1.4L10.698.68a1 1 0 00-1.396 0L4.25 5.6 5.688 7 9 3.775V13h2z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default UploadIcon;
