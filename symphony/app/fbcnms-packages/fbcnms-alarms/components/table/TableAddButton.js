/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 *
 * The add button at the bottom right of the tables
 */

import * as React from 'react';
import AddIcon from '@material-ui/icons/Add';
import Fab from '@material-ui/core/Fab';
import {makeStyles} from '@material-ui/styles';

type Props = {
  label: string,
  onClick: () => void,
};

const useStyles = makeStyles(theme => ({
  addButton: {
    position: 'fixed',
    bottom: 0,
    right: 0,
    margin: theme.spacing(2),
  },
}));
export default function TableAddButton({label, onClick, ...props}: Props) {
  const classes = useStyles();
  return (
    <Fab
      {...props}
      className={classes.addButton}
      color="primary"
      onClick={onClick}
      aria-label={label}>
      <AddIcon />
    </Fab>
  );
}
