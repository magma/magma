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

const DeleteIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(4,2)">
      <path
        d="M10 0a1 1 0 011 1v1h4a1 1 0 011 1v5h-1v11a1 1 0 01-1 1H2a1 1 0 01-1-1V8H0V3a1 1 0 011-1h4V1a1 1 0 011-1h4zm2.999 8h-10L3 18h10l-.001-10zM7 10v6H5v-6h2zm4 0v6H9v-6h2zm3-4V4H2v2h12z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default DeleteIcon;
