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

const NewTabIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M13.25 5v2H7v10h10v-6.25h2v6.5a1.75 1.75 0 01-1.606 1.744L17.25 19H6.75a1.75 1.75 0 01-1.744-1.606L5 17.25V6.75a1.75 1.75 0 011.606-1.744L6.75 5h6.5zM20 3a1 1 0 011 1v4h-2V6.414l-6 6L11.586 11l5.999-6H16V3h4z" />
  </SvgIcon>
);

export default NewTabIcon;
