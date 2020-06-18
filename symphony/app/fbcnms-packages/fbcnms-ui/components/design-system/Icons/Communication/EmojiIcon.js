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
    <path d="M12 20a8 8 0 100-16 8 8 0 000 16zm0 2C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zM8.5 11a1.5 1.5 0 100-3 1.5 1.5 0 000 3zM8 13c0 1.713 2.008 3.226 4 3.226s4-1.513 4-3.226h2c0 3.235-3.312 5.135-5.818 5.223l-.182.003c-2.524 0-6-1.912-6-5.226h2zm7.5-5a1.5 1.5 0 110 3 1.5 1.5 0 010-3z" />
  </SvgIcon>
);

export default EmojiIcon;
