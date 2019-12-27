/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 *
 * Wrappper component for editors such as AddEditRule, AddEditReceiver, etc
 */

import * as React from 'react';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import {makeStyles} from '@material-ui/styles';

type Props = {
  children: React.Node,
  onExit: () => void,
  onSave: () => void,
  isNew: boolean,
};
const useStyles = makeStyles(theme => ({
  gridContainer: {
    flexGrow: 1,
  },
  editingSpace: {
    height: '100%',
    padding: theme.spacing(3),
  },
}));
export default function Editor({
  children,
  isNew,
  onExit,
  onSave,
  ...props
}: Props) {
  const classes = useStyles();

  return (
    <Grid {...props} className={classes.gridContainer} container spacing={0}>
      <Grid className={classes.editingSpace} item xs>
        <form
          onSubmit={e => {
            e.preventDefault();
            onSave();
          }}>
          <Grid container spacing={3}>
            <Grid
              container
              item
              direction="column"
              spacing={2}
              wrap="nowrap"
              xs={12}
              sm={4}>
              {children}
            </Grid>
            <Grid container item spacing={1} xs={12}>
              <Grid item>
                <Button
                  variant="outlined"
                  onClick={() => onExit()}
                  className={classes.button}>
                  Close
                </Button>
              </Grid>
              <Grid item>
                <Button
                  variant="contained"
                  color="primary"
                  type="submit"
                  className={classes.button}
                  data-testid="editor-submit-button">
                  {isNew ? 'Add' : 'Save'}
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </form>
      </Grid>
    </Grid>
  );
}
