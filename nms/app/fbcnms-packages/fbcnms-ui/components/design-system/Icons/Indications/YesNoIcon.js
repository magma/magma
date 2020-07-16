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

const YesNoIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 2c5.523 0 10 4.477 10 10s-4.477 10-10 10S2 17.523 2 12 6.477 2 12 2zm0 2a8 8 0 100 16 8 8 0 000-16zm5.355 8.293a1 1 0 01.083 1.32l-.083.094-3.894 3.894a1 1 0 01-1.32.083l-.094-.083-1.754-1.754a1 1 0 011.32-1.497l.094.083 1.047 1.047 3.187-3.187a1 1 0 011.414 0zm-4.657-5.019a.95.95 0 01.101 1.28l-.081.092-1.337 1.353 1.337 1.355a.95.95 0 01-.02 1.372 1.018 1.018 0 01-1.41-.02L10 11.4l-1.288 1.306c-.354.355-.92.39-1.315.099l-.095-.08a.95.95 0 01-.101-1.278l.081-.093L8.618 10 7.282 8.646a.95.95 0 01.02-1.372 1.018 1.018 0 011.41.02L10 8.599l1.288-1.305a1.018 1.018 0 011.41-.02z" />
  </SvgIcon>
);

export default YesNoIcon;
