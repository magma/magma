/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import AddIcon from '@material-ui/icons/Add';
import Fab from '@material-ui/core/Fab';
import React from 'react';
import RemoveCircleOutlineIcon from '@material-ui/icons/RemoveCircleOutline';
import classNames from 'classnames';
import symphony from '../theme/symphony';
import {makeStyles} from '@material-ui/styles';

type Props = {
  action: 'add' | 'remove',
  onClick: () => void,
};

const useStyles = makeStyles(() => ({
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
    fill: symphony.palette.D600,
    fontSize: '28px',
  },
  addFabIcon: {
    fontSize: '20px',
  },
}));

const ActionButton = (props: Props) => {
  const {action, onClick} = props;
  const classes = useStyles();

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

export default ActionButton;
