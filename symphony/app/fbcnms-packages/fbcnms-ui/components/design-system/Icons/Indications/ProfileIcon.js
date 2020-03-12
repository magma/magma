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

const ProfileIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 2c5.523 0 10 4.477 10 10a9.982 9.982 0 01-3.804 7.85L18 20a9.952 9.952 0 01-6 2C6.477 22 2 17.523 2 12S6.477 2 12 2zm0 12c-2.198 0-4 1.892-4 4.25v.68A7.963 7.963 0 0012 20c1.458 0 2.824-.39 4.001-1.07L16 18.25c0-2.358-1.802-4.25-4-4.25zm0-10a8 8 0 00-5.94 13.36c.314-2.282 1.81-4.173 3.83-4.963a4 4 0 114.221.001c2.019.789 3.515 2.68 3.828 4.96A8 8 0 0012 4zm0 3a2 2 0 100 4 2 2 0 000-4z" />
  </SvgIcon>
);

export default ProfileIcon;
