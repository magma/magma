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
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  header: {
    backgroundColor: theme.palette.common.white,
    borderBottom: `1px solid ${theme.palette.divider}`,
  },
  paper: {
    margin: theme.spacing(3),
  },
  searchBar: {
    marginLeft: theme.spacing(1),
  },
}));

export default function Logs() {
  const classes = useStyles();
  return (
    <>
      <div className={classes.header}>
        <TextField
          id="outlined-search"
          placeholder="Search logs"
          type="search"
          margin="normal"
          variant="outlined"
          className={classes.searchBar}
        />
      </div>
      <div className={classes.paper} />
    </>
  );
}
