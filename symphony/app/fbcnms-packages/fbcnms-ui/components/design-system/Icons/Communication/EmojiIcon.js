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

const EmojiIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(2,2)">
      <path
        d="M10 18a8 8 0 100-16 8 8 0 000 16zm0 2C4.477 20 0 15.523 0 10S4.477 0 10 0s10 4.477 10 10-4.477 10-10 10zM6.5 9a1.5 1.5 0 100-3 1.5 1.5 0 000 3zM6 11c0 1.713 2.008 3.226 4 3.226s4-1.513 4-3.226h2c0 3.235-3.312 5.135-5.818 5.223l-.182.003c-2.524 0-6-1.912-6-5.226h2zm7.5-5a1.5 1.5 0 110 3 1.5 1.5 0 010-3z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default EmojiIcon;
