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

const ThreeDotsIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <defs>
      <path
        d="M12 7a2 2 0 100-4 2 2 0 000 4zm0 7a2 2 0 100-4 2 2 0 000 4zm0 7a2 2 0 100-4 2 2 0 000 4z"
        id="treeDotsIcon"
      />
    </defs>
    <use xlinkHref="#treeDotsIcon" fillRule="evenodd" />
  </SvgIcon>
);

export default ThreeDotsIcon;
