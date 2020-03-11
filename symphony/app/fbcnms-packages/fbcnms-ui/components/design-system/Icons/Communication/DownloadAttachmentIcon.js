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

const DownloadAttachmentIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 20a8 8 0 100-16 8 8 0 000 16zm0 2C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm1-15v7.589L15.563 12l1.187 1.2L13 16.988V17h-.012l-.277.282a1 1 0 01-1.422 0L11.011 17H11v-.011L7.25 13.2 8.438 12 11 14.589V7h2z" />
  </SvgIcon>
);

export default DownloadAttachmentIcon;
