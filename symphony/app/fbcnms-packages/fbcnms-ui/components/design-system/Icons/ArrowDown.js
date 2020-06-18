/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import SymphonyTheme from '../../../theme/symphony';

type Props = {
  className?: string,
};

const ArrowDown = ({className}: Props) => (
  <svg width="10" height="10" xmlns="http://www.w3.org/2000/svg">
    <path
      className={className}
      d="M4.444 7.916V0h1.112v7.916l3.632-3.472.812.777L5 10 0 5.22l.812-.776 3.632 3.472z"
      fill={SymphonyTheme.palette.primary}
    />
  </svg>
);

export default ArrowDown;
