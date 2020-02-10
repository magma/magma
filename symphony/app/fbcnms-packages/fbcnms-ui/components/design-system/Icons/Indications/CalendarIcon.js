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
    <g transform="translate(2,2)">
      <path
        d="M19 0a1 1 0 011 1l-.001 5H20v2h-.001L20 19a1 1 0 01-1 1H1a1 1 0 01-1-1V1a1 1 0 011-1h18zm-1 8H2v10h16V8zM7 14v2H5v-2h2zm4 0v2H9v-2h2zm4 0v2h-2v-2h2zm-8-4v2H5v-2h2zm4 0v2H9v-2h2zm4 0v2h-2v-2h2zm0-5h-2V2H7v3H5V2H2v4h16V2h-3v3z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default CalendarIcon;
