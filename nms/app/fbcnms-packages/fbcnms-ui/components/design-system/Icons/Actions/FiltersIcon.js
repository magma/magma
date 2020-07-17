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

const FiltersIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M9 12a4.002 4.002 0 013.874 3H21v2h-8.126a4.002 4.002 0 01-7.748 0H3v-2h2.126c.444-1.725 2.01-3 3.874-3zm0 2a2 2 0 100 4 2 2 0 000-4zm6-10a4.002 4.002 0 013.874 3H21v2h-2.126a4.002 4.002 0 01-7.748 0H3V7h8.126c.444-1.725 2.01-3 3.874-3zm0 2a2 2 0 100 4 2 2 0 000-4z" />
  </SvgIcon>
);

export default FiltersIcon;
