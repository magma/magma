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

const InfoSmallIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 20a8 8 0 100-16 8 8 0 000 16zm0 2C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm-1.001-11H10V9h3v6h1v2h-4v-2h1l-.001-4zM11 6h2v2h-2V6z" />
  </SvgIcon>
);

export default InfoSmallIcon;
