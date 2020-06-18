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
import SvgIcon from '@material-ui/core/SvgIcon';

type Props = {
  className?: string,
};

const ProfileIcon = (props: Props) => (
  <SvgIcon
    color="inherit"
    viewBox="0 0 12 15"
    width="12px"
    height="15px"
    className={props.className}>
    <g id="baseline-account_circle-24px" transform="translate(-0.2, 0)">
      <path d="m6 0c1.7 0 3 1.3 3 3s-1.3 3-3 3-3-1.3-3-3 1.3-3 3-3zm0 14.2c-2.5 0-4.7-1.3-6-3.2 0-2 4-3.1 6-3.1s6 1.1 6 3.1c-1.3 1.9-3.5 3.2-6 3.2z" />
    </g>
  </SvgIcon>
);

export default ProfileIcon;
