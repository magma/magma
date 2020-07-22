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

const CalendarIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M21 2a1 1 0 011 1l-.001 5H22v2h-.001L22 21a1 1 0 01-1 1H3a1 1 0 01-1-1V3a1 1 0 011-1h18zm-1 8H4v10h16V10zM9 16v2H7v-2h2zm4 0v2h-2v-2h2zm4 0v2h-2v-2h2zm-8-4v2H7v-2h2zm4 0v2h-2v-2h2zm4 0v2h-2v-2h2zm0-5h-2V4H9v3H7V4H4v4h16V4h-3v3z" />
  </SvgIcon>
);

export default CalendarIcon;
