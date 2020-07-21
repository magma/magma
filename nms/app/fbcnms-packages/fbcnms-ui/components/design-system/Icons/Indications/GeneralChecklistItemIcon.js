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

const GeneralChecklistItemIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 3v2H5v14h14v-8h2v8a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h7zm6.67-.707L21.5 5.12a1 1 0 010 1.415l-2.122 2.12-1.414 1.415-6.193 6.192-3.63.519a1 1 0 01-1.13-1.131l.518-3.63 9.728-9.728a1 1 0 011.414 0zM15.136 7.24l-5.72 5.723-.236 1.65 1.65-.236 5.72-5.723-1.414-1.414zm2.829-2.827l-1.415 1.413 1.415 1.414 1.414-1.413-1.414-1.414z" />
  </SvgIcon>
);

export default GeneralChecklistItemIcon;
