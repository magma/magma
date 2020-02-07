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

const RemoveIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <defs>
      <path
        d="M12 20a8 8 0 100-16 8 8 0 000 16zm0 2C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm-5-9v-2h10v2H7z"
        id="removeIcon"
      />
    </defs>
    <use xlinkHref="#removeIcon" fillRule="evenodd" />
  </SvgIcon>
);

export default RemoveIcon;
