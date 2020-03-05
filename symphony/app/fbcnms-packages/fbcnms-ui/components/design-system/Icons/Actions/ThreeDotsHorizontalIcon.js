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

const ThreeDotsHorizontalIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <defs>
      <path
        d="M19 14a2 2 0 100-4 2 2 0 000 4zm-7 0a2 2 0 100-4 2 2 0 000 4zm-7 0a2 2 0 100-4 2 2 0 000 4z"
        id="threeDotsHorizontalIcon"
      />
    </defs>
    <use xlinkHref="#threeDotsHorizontalIcon" />
  </SvgIcon>
);

export default ThreeDotsHorizontalIcon;
