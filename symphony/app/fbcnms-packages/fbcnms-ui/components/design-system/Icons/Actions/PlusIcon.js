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

const IconIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M13 5v6h6v2h-6v6h-2v-6H5v-2h6V5h2z" />
  </SvgIcon>
);

export default IconIcon;
