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

const MandatoryIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M13 3v7.267l6.294-3.633 1 1.732L14 12l6.294 3.634-1 1.732L13 13.732V21h-2v-7.268l-6.294 3.634-1-1.732 6.293-3.635-6.293-3.633 1-1.732L11 10.267V3h2z" />
  </SvgIcon>
);

export default MandatoryIcon;
