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

const DownloadAttachmentIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <g transform="translate(2,2)">
      <path
        d="M10 18a8 8 0 100-16 8 8 0 000 16zm0 2C4.477 20 0 15.523 0 10S4.477 0 10 0s10 4.477 10 10-4.477 10-10 10zm1-15v7.589L13.563 10l1.187 1.2L11 14.988V15h-.012l-.277.282a1 1 0 01-1.422 0L9.011 15H9v-.011L5.25 11.2 6.438 10 9 12.589V5h2z"
        fillRule="evenodd"
      />
    </g>
  </SvgIcon>
);

export default DownloadAttachmentIcon;
