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

const MessageIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(2,3)">
      <path
        d="M2 14.263L5.394 12H18V2H2v12.263zM1 0h18a1 1 0 011 1v12a1 1 0 01-1 1H6l-6 4V1a1 1 0 011-1zm3 4h12v2H4V4zm0 4h12v2H4V8z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default MessageIcon;
