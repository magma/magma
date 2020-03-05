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

const FilterIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M6 6v1l3.791 2.844L11.608 18h.792l1.809-8.156L18 7V6H6zM4 4h16v3.5a1 1 0 01-.4.8L16 11l-1.995 9h-4L8 11 4.4 8.3a1 1 0 01-.4-.8V4z" />
  </SvgIcon>
);

export default FilterIcon;
