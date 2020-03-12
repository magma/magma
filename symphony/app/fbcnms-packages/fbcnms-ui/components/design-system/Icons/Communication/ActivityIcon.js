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

const ActivityIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 16v2H2v-2h10zm10-5v2H2v-2h20zm0-5v2H2V6h20z" />
  </SvgIcon>
);

export default ActivityIcon;
