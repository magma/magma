/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';

type Props = {
  onClick?: () => void,
  error?: boolean,
  children: any,
};

type State = {};

export default class Button extends React.Component<Props, State> {
  render() {
    const styles = {
      border: '1px solid #bbb',
      borderRadius: 6,
      cursor: 'pointer',
      fontSize: 15,
      padding: '3px 10px',
    };
    if (this.props.error) {
      styles['border'] = '1px solid red';
    }
    return (
      <button style={styles} onClick={this.props.onClick}>
        {this.props.children}
      </button>
    );
  }
}
