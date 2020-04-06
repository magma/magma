/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import symphony from '../theme/symphony';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  onClick: () => void,
  children: React.Node,
} & WithStyles<typeof styles>;

const styles = () => ({
  root: {
    cursor: 'pointer',
    textDecoration: 'underline',
    color: symphony.palette.primary,
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
