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

import React from 'react';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  onClick: () => void,
  children: any,
} & WithStyles<typeof styles>;

const styles = theme => ({
  root: {
    cursor: 'pointer',
    textDecoration: 'underline',
    color: theme.palette.primary.main,
  },
});

// TODO(T38660666) - style according to design
class Link extends React.Component<Props> {
  render() {
    const {classes, onClick, children} = this.props;
    return (
      <a className={classes.root} onClick={onClick}>
        {children}
      </a>
    );
  }
}

export default withStyles(styles)(Link);
