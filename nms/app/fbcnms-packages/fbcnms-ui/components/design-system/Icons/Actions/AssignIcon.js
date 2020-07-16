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

const AddIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <defs>
      <path
        d="M13 7v4h4v2h-4.001L13 17h-2l-.001-4H7v-2h4V7h2zM5 5h3l2-2H5a2 2 0 00-2 2v5l2-2V5zm16 9v5a2 2 0 01-2 2h-5l2-2h3v-3l2-2zm-2-9v3l2 2V5a2 2 0 00-2-2h-5l2 2h3zM8 19l2 2H5a2 2 0 01-2-2v-5l2 2v3h3z"
        id="assignIcon"
      />
    </defs>
    <use xlinkHref="#assignIcon" />
  </SvgIcon>
);

export default AddIcon;
