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

const PhoneIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M21.394 16.895a2 2 0 01-.186.214l-2.693 2.694a3 3 0 01-3.651.459l-2.435-1.444a20.194 20.194 0 01-7.071-7.071L3.914 9.312a3 3 0 01.459-3.652l2.693-2.693a2 2 0 013.014.215L12.42 6.3a1 1 0 01-.093 1.307L10.594 9.34c1.242 1.576 2.666 3 4.242 4.243l1.733-1.734a1 1 0 011.307-.092l3.118 2.338a2 2 0 01.4 2.8zM8.48 4.382L5.787 7.075a1 1 0 00-.153 1.217l1.444 2.435a18.194 18.194 0 006.37 6.37l2.436 1.444a1 1 0 001.217-.153l2.693-2.693-2.424-1.818-2.376 2.377-1.396-1.1a27.495 27.495 0 01-4.575-4.576l-1.1-1.395 2.376-2.377L8.48 4.382z" />
  </SvgIcon>
);

export default PhoneIcon;
