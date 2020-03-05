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

const ChecklistCheckIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M21 2a1 1 0 011 1v18a1 1 0 01-1 1H3a1 1 0 01-1-1V3a1 1 0 011-1h18zm-1 2H4v16h16V4zm-4.207 4.4a1 1 0 011.414 1.413l-5.894 5.894a1 1 0 01-1.414 0l-2.754-2.754A1 1 0 118.56 11.54l2.047 2.047z" />
  </SvgIcon>
);

export default ChecklistCheckIcon;
