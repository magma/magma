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

import AddIcon from '@material-ui/icons/Add';
import Fab from '@material-ui/core/Fab';
import React from 'react';
import RemoveCircleOutlineIcon from '@material-ui/icons/RemoveCircleOutline';
import classNames from 'classnames';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  action: 'add' | 'remove',
  onClick: () => void,
} & WithStyles<typeof styles>;

const styles = theme => ({
  actionButton: {
    width: '20px',
    height: '20px',
    minWidth: '20px',
    minHeight: '20px',
    boxShadow: 'none',
  },
  removeFab: {
    backgroundColor: 'transparent',
    marginTop: '-4px',
  },
  removeFabIcon: {
    fill: theme.palette.grey[600],
    fontSize: '28px',
  },
  addFabIcon: {
    fontSize: '20px',
  },
});

const ActionButton = (props: Props) => {
  const {action, classes, onClick} = props;
  switch (action) {
    case 'add':
      return (
        <Fab
          className={classes.actionButton}
          color="primary"
          size="small"
          onClick={onClick}>
          <AddIcon className={classes.addFabIcon} />
        </Fab>
      );
    case 'remove':
      return (
        <Fab
          className={classNames(classes.actionButton, classes.removeFab)}
          size="small"
          onClick={onClick}>
          <RemoveCircleOutlineIcon className={classes.removeFabIcon} />
        </Fab>
      );
    default:
      return null;
  }
};

export default withStyles(styles)(ActionButton);
