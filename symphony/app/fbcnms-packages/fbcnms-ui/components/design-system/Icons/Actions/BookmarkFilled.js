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

const BookmarkFilledIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M14 4a2 2 0 012 2v11.586A2 2 0 0112.586 19L10 16.414 7.414 19A2 2 0 014 17.586V6a2 2 0 012-2h8zm-4-4a1 1 0 010 2H2v12.586a1 1 0 11-2 0V2a2 2 0 012-2z" />
  </SvgIcon>
);

export default BookmarkFilledIcon;
