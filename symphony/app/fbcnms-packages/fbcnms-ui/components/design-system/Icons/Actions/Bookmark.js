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

const BookmarkIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M18 6a2 2 0 012 2v11.586A2 2 0 0116.586 21L14 18.414 11.414 21A2 2 0 018 19.586V8a2 2 0 012-2h8zm0 2h-8v11.586L12.586 17a2 2 0 012.828 0L18 19.586V8zm-4-6a1 1 0 010 2H6v12.586a1 1 0 11-2 0V4a2 2 0 012-2z" />
  </SvgIcon>
);

export default BookmarkIcon;
