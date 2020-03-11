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

const EmailIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M21 5a1 1 0 011 1v12a1 1 0 01-1 1H3a1 1 0 01-1-1V6a1 1 0 011-1h18zm-6.419 8.06L12 15.314 9.418 13.06 6.135 17h11.729l-3.283-3.94zM20 8.327l-3.912 3.416L20 16.437v-8.11zM4 8.328v8.108l3.911-4.693L4 8.328zM18.479 7H5.52L12 12.659 18.479 7z" />
  </SvgIcon>
);

export default EmailIcon;
