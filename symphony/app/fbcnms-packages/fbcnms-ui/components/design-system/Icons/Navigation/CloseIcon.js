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

const CloseIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path
      d="M10.588 0L12 1.412 7.411 6 12 10.588 10.588 12 6 7.411 1.412 12 0 10.588 4.588 6 0 1.412 1.412 0 6 4.588 10.588 0z"
      fillRule="evenodd"
    />
  </SvgIcon>
);

export default CloseIcon;
