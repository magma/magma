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

const EditIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M20.5 6.121a1 1 0 010 1.415l-2.122 2.12-1.414 1.415-9.193 9.192-3.63.519a1 1 0 01-1.13-1.131l.518-3.63L16.257 3.293a1 1 0 011.414 0l2.828 2.828zm-4.951 3.535l-1.415-1.414-8.72 8.722-.235 1.65 1.65-.236 8.72-8.722zm2.829-2.828l-1.414-1.414-1.415 1.414 1.414 1.414 1.415-1.414z" />
  </SvgIcon>
);

export default EditIcon;
