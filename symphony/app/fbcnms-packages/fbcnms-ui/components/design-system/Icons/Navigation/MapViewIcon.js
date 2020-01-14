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

const IconIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path
      d="M19.5 0c.28 0 .5.22.5.5v17.12c0 .23-.15.41-.36.48L13 20l-6-2.1-6.34 2.07L.5 20c-.28 0-.5-.22-.5-.5V2.38c0-.23.15-.41.36-.48L7 0l6 2.1L19.34.03 19.5 0zM18 2.694l-4 1.238v13.687l4-1.077V2.694zM8 2.47v13.662l4 1.4V3.869l-4-1.4zM6 2.38L2 3.458v13.848l4-1.24V2.38z"
      fillRule="evenodd"
    />
  </SvgIcon>
);

export default IconIcon;
