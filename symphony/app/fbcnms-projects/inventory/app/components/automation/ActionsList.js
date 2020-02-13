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
import ActionsHead from './ActionsHead';
import ActionsListCard from './ActionsListCard';
import Button from '@fbcnms/ui/components/design-system/Button';
import Grid from '@material-ui/core/Grid';
import Text from '@fbcnms/ui/components/design-system/Text';

import useRouter from '@fbcnms/ui/hooks/useRouter';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
  },
  spacer: {
    flexGrow: 1,
  },
  main: {
    margin: theme.spacing(2),
  },
}));

export default function ActionsList() {
  const classes = useStyles();
  const {history} = useRouter();

  return (
    <div className={classes.root}>
      <ActionsHead>
        <div className={classes.spacer} />
        <div>
          <Button
            variant="text"
            onClick={() => history.push('/automation/actions')}>
            Create Rule
          </Button>
        </div>
      </ActionsHead>
      <div className={classes.main}>
        <Text variant="h6">Existing Rules</Text>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <ActionsListCard />
          </Grid>
        </Grid>
      </div>
    </div>
  );
}
