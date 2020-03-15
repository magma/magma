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

const CloseIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M16.588 6L18 7.412 13.411 12 18 16.588 16.588 18 12 13.411 7.412 18 6 16.588 10.588 12 6 7.412 7.412 6 12 10.588 16.588 6z" />
  </SvgIcon>
);

export default CloseIcon;
