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

const ListViewIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path
      d="M19 0a1 1 0 011 1v18a1 1 0 01-1 1H1a1 1 0 01-1-1V1a1 1 0 011-1h18zm-1 2H2v16h16V2zM7 13v2H5v-2h2zm8 0v2H9v-2h6zM7 9v2H5V9h2zm8 0v2H9V9h6zM7 5v2H5V5h2zm8 0v2H9V5h6z"
      fillRule="evenodd"
    />
  </SvgIcon>
);

export default ListViewIcon;
