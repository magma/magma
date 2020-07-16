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

const LockIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M8 0a4 4 0 014 4v2h2a2 2 0 012 2v8a2 2 0 01-2 2H2a2 2 0 01-2-2V8a2 2 0 012-2h2V4a4 4 0 014-4zm6 8H2v8h12V8zm-4 2v2H9v3H7v-3H6v-2h4zM8 2a2 2 0 00-2 2v2h4V4a2 2 0 00-2-2z" />
  </SvgIcon>
);

export default LockIcon;
