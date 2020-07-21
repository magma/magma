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

const GroupIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M16.89 4C19.152 4 21 5.785 21 8a3.936 3.936 0 01-1.283 2.903C21.095 11.803 22 13.296 22 15v4h-5v2H7v-2H2v-4c0-1.704.906-3.197 2.282-4.098A3.933 3.933 0 013 8c0-2.215 1.847-4 4.11-4 1.239 0 2.374.531 3.137 1.403A3.987 3.987 0 0112 5c.632 0 1.229.146 1.76.407A4.15 4.15 0 0116.89 4zM12 14a3 3 0 00-3 3v2h6v-2a3 3 0 00-3-3zm4.654-2c-.797 0-1.542.251-2.127.685A4.996 4.996 0 0117 17h3v-2c0-1.637-1.48-3-3.346-3zm-9.308 0c-1.801 0-3.243 1.27-3.34 2.832L4 15v2h3c0-1.839.993-3.446 2.472-4.315a3.54 3.54 0 00-1.889-.678L7.346 12zM12 7a2 2 0 100 4 2 2 0 000-4zM7.11 6c-1.163 0-2.096.902-2.096 2 0 1.098.933 2 2.097 2 .343 0 .672-.079.965-.224a3.982 3.982 0 01.64-3.061A2.13 2.13 0 007.111 6zm9.78 0c-.639 0-1.22.274-1.608.714a3.978 3.978 0 01.643 3.06c.292.147.62.226.964.226 1.164 0 2.097-.902 2.097-2 0-1.098-.933-2-2.097-2z" />
  </SvgIcon>
);

export default GroupIcon;
