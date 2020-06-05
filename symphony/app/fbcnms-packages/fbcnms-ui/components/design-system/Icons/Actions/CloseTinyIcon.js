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

const CloseTinyIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M16 15.2l-.8.8-3.2-3.201L8.8 16l-.8-.8 3.201-3.2L8 8.8l.8-.8 3.2 3.2L15.2 8l.8.8-3.2 3.199L16 15.2z" />
  </SvgIcon>
);

export default CloseTinyIcon;
