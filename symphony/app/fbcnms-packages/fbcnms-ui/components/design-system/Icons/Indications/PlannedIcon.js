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

const PlannedIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 20a8 8 0 100-16 8 8 0 000 16zm0 2C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm1-15l-.001 5.184 3.658 3.659-1.414 1.414L11 13.014l.014-.014H11V7h2z" />
  </SvgIcon>
);

export default PlannedIcon;
