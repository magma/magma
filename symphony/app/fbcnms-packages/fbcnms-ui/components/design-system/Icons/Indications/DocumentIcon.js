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

const DocumentIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M15 2l5 6v13a1 1 0 01-1 1H5a1 1 0 01-1-1V3a1 1 0 011-1h10zm-2 2H6v16h12V9h-5V4zm3 10v2H8v-2h8zm0-3v2H8v-2h8zm-1-5.875V7h1.563L15 5.125z" />
  </SvgIcon>
);

export default DocumentIcon;
