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
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import LinearProgress from '@material-ui/core/LinearProgress';
import RefreshIcon from '@material-ui/icons/Refresh';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  header: {
    padding: theme.spacing(3),
    display: 'flex',
    justifyContent: 'space-between',
    backgroundColor: 'white',
    borderBottom: `1px solid ${theme.palette.divider}`,
  },
}));

type Props = {|
  title: string,
  isLoading: boolean,
  lastRefreshTime: string,
  onRefreshClick: string => void,
  children?: React.Node,
  'data-testid'?: string,
|};

export default function AlarmsHeader(props: Props) {
  const classes = useStyles();
  const {
    title,
    isLoading,
    lastRefreshTime,
    onRefreshClick,
    children,
    ...divProps
  } = props;

  return (
    <>
      <div className={classes.header} {...divProps}>
        <Typography variant="h5">{title}</Typography>
        <div>
          <Grid container spacing={1} justify="flex-end" alignItems="center">
            <Grid item>
              <Tooltip title={'Last refreshed: ' + lastRefreshTime}>
                <div>
                  <IconButton
                    color="inherit"
                    onClick={() => onRefreshClick(new Date().toLocaleString())}
                    disabled={isLoading}>
                    <RefreshIcon />
                  </IconButton>
                </div>
              </Tooltip>
            </Grid>
            {React.Children.map(children, child => (
              <Grid item>{child}</Grid>
            ))}
          </Grid>
        </div>
      </div>
      {isLoading ? <LinearProgress /> : null}
    </>
  );
}
