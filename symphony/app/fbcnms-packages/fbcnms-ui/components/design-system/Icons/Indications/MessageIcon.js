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
    <path d="M4 19.263L7.394 17H20V7H4v12.263zM3 5h18a1 1 0 011 1v12a1 1 0 01-1 1H8l-6 4V6a1 1 0 011-1zm3 4h12v2H6V9zm0 4h12v2H6v-2z" />
  </SvgIcon>
);

export default MessageIcon;
