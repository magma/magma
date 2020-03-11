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

const InfoTinyIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 18a6 6 0 100-12 6 6 0 000 12zm-1-7h2v5h-2v-5zm0-3h2v2h-2V8z" />
  </SvgIcon>
);

export default InfoTinyIcon;
