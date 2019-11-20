/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import Paper from '@material-ui/core/Paper';
import Text from '@fbcnms/ui/components/design-system/Text';

import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  paper: {
    padding: theme.spacing(2),
    textAlign: 'center',
    '&:hover': {
      background: symphony.palette.D100,
    },
    cursor: 'pointer',
  },
}));

export default function ActionsCard(props: {
  icon: React.Node,
  message: React.Node,
  onClick: () => void,
}) {
  const classes = useStyles();
  return (
    <Paper className={classes.paper} onClick={props.onClick}>
      <div>{props.icon}</div>
      <Text variant="subtitle1">{props.message}</Text>
    </Paper>
  );
}
