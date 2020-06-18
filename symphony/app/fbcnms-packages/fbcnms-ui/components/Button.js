/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';

type Props = {
  onClick?: () => void,
  error?: boolean,
  children: React.Node,
};

export default function Button(props: Props) {
  const styles = {
    border: '1px solid #bbb',
    borderRadius: 6,
    cursor: 'pointer',
    fontSize: 15,
    padding: '3px 10px',
  };
  if (props.error != null) {
    styles['border'] = '1px solid red';
  }
  return (
    <button style={styles} onClick={props.onClick}>
      {props.children}
    </button>
  );
}
