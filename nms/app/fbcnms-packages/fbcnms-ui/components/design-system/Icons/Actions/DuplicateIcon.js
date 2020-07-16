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

const DuplicateIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M8 8v12h12V8H8zM7 6h14a1 1 0 011 1v14a1 1 0 01-1 1H7a1 1 0 01-1-1V7a1 1 0 011-1zm-3 9h2v2H3a1 1 0 01-1-1V3a1 1 0 011-1h13a1 1 0 011 1v3h-2V4H4v11z" />
  </SvgIcon>
);

export default DuplicateIcon;
