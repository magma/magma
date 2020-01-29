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

const DuplicateIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(2,2)">
      <path
        d="M6 6v12h12V6H6zM5 4h14a1 1 0 011 1v14a1 1 0 01-1 1H5a1 1 0 01-1-1V5a1 1 0 011-1zm-3 9h2v2H1a1 1 0 01-1-1V1a1 1 0 011-1h13a1 1 0 011 1v3h-2V2H2v11z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default DuplicateIcon;
