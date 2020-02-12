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

const InfoTinyIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(6,6)">
      <path
        d="M6 12A6 6 0 106 0a6 6 0 000 12zM5 5h2v5H5V5zm0-3h2v2H5V2z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default InfoTinyIcon;
