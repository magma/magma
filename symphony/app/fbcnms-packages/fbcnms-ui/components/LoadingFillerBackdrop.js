/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';

import LoadingFiller from './LoadingFiller';
import React from 'react';

import {withStyles} from '@material-ui/core/styles';

const styles = {
  backdrop: {
    alignItems: 'center',
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    bottom: 0,
    display: 'flex',
    justifyContent: 'center',
    left: 0,
    position: 'fixed',
    right: 0,
    top: 0,
    zIndex: '13000',
  },
};

type Props = WithStyles<typeof styles> & {};

class LoadingFillerBackdrop extends React.Component<Props> {
  render() {
    return (
      <div className={this.props.classes.backdrop}>
        <LoadingFiller />
      </div>
    );
  }
}

export default withStyles(styles)(LoadingFillerBackdrop);
