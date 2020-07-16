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
  variant: 'small' | 'large',
  className?: string,
};

const ActiveProviderEndpointIcon = (props: Props) =>
  props.variant === 'large' ? (
    <SvgIcon
      color="inherit"
      width="40px"
      height="44px"
      viewBox="0 0 40 44"
      className={props.className}>
      <g
        fill="none"
        fillRule="evenodd"
        transform="translate(-1 1)"
        stroke="#9E43DF"
        strokeWidth="2">
        <path
          d="m24 1.7321 12.187 7.0359c1.8564 1.0718 3 3.0526 3 5.1962v14.072c0 2.1436-1.1436 4.1244-3 5.1962l-12.187 7.0359c-1.8564 1.0718-4.1436 1.0718-6 0l-12.187-7.0359c-1.8564-1.0718-3-3.0526-3-5.1962v-14.072c0-2.1436 1.1436-4.1244 3-5.1962l12.187-7.0359c1.8564-1.0718 4.1436-1.0718 6 0z"
          fill="#E0B8FC"
        />
        <polyline points="39 11 21 21 3 11" />
        <path d="m21 21v20" strokeLinecap="square" />
      </g>
    </SvgIcon>
  ) : (
    <SvgIcon
      color="inherit"
      width="16px"
      height="18px"
      viewBox="0 0 16 18"
      className={props.className}>
      <g
        fill="none"
        fillRule="evenodd"
        transform="translate(0 1)"
        stroke="#9E43DF"
        strokeWidth="2">
        <path
          d="m9.5 0.86603 3.9282 2.2679c0.9282 0.5359 1.5 1.5263 1.5 2.5981v4.5359c0 1.0718-0.5718 2.0622-1.5 2.5981l-3.9282 2.2679c-0.9282 0.5359-2.0718 0.5359-3 0l-3.9282-2.2679c-0.9282-0.5359-1.5-1.5263-1.5-2.5981v-4.5359c0-1.0718 0.5718-2.0622 1.5-2.5981l3.9282-2.2679c0.9282-0.5359 2.0718-0.5359 3 0z"
          fill="#E0B8FC"
        />
        <polyline points="14.857 4.1905 8 8 1.1429 4.1905" />
        <path d="m8 8v7" strokeLinecap="square" />
      </g>
    </SvgIcon>
  );

export default ActiveProviderEndpointIcon;
