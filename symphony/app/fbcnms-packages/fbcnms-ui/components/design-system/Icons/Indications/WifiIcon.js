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

const WifiIcon = (props: SvgIconStyleProps) => (
  <SvgIcon {...props}>
    <path d="M12 17a1 1 0 110 2 1 1 0 010-2zm2.76-3.136l.188.138.736.57L14.506 16l-.737-.57c-.98-.76-2.35-.8-3.371-.12l-.167.12-.737.57-1.178-1.427.736-.571a4.84 4.84 0 015.708-.138zm3.014-3.152l.266.18.802.562L17.675 13l-.802-.562c-2.82-1.976-6.61-2.034-9.487-.174l-.259.174-.802.562-1.167-1.546.802-.562c3.513-2.461 8.242-2.521 11.814-.18zm3.093-3.089l.307.206.826.566L20.834 10l-.826-.566c-4.703-3.226-10.954-3.292-15.72-.198l-.296.198-.826.566L2 8.395l.826-.566c5.398-3.702 12.577-3.77 18.04-.206z" />
  </SvgIcon>
);

export default WifiIcon;
