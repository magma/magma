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

const HierarchyArrowIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(6.5,7)">
      <path
        d="M2 0v5l3.999-.001L6 2l5 4-5 4-.001-3.001L0 7V0h2z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default HierarchyArrowIcon;
