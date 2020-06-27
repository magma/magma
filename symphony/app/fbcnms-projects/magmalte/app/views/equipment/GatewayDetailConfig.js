/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {lte_gateway} from '@fbcnms/magma-api';

import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import Text from '@fbcnms/ui/components/design-system/Text';

import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
  },
  itemTitle: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  itemValue: {
    color: colors.primary.brightGray,
  },
}));

export default function GatewayConfig({gwInfo}: {gwInfo: lte_gateway}) {
  const classes = useStyles();
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3}>
        <Grid container spacing={3} item xs={12}>
          <Text>
            <SettingsIcon /> Config
          </Text>
        </Grid>
        <Grid container spacing={3} item xs={6}>
          <Grid item xs={12}>
            <Grid container>
              <Grid item xs={6}>
                <Text>Gateway</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <GatewayInfoConfig readOnly={true} gwInfo={gwInfo} />
          </Grid>
          <Grid item xs={12}>
            <Grid container>
              <Grid item xs={6}>
                <Text>Aggregations</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <Paper className={classes.paper} />
          </Grid>
        </Grid>
        <Grid container spacing={3} item xs={6}>
          <Grid item xs={12}>
            <Grid container>
              <Grid item xs={6}>
                <Text>EPC</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <Paper className={classes.paper} />
          </Grid>
          <Grid item xs={12}>
            <Grid container>
              <Grid item xs={6}>
                <Text>Ran</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <Paper className={classes.paper} />
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function GatewayInfoConfig({gwInfo}: {gwInfo: lte_gateway}) {
  const classes = useStyles();
  const typographyProps = {
    primaryTypographyProps: {
      variant: 'caption',
      className: classes.itemTitle,
    },
    secondaryTypographyProps: {
      variant: 'h6',
      className: classes.itemValue,
    },
  };
  return (
    <List component={Paper}>
      <ListItem>
        <ListItemText
          primary="Name"
          secondary={gwInfo.name}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={gwInfo.id}
          primary="Gateway ID"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={gwInfo.device.hardware_id}
          primary="Hardware UUID"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={
            gwInfo.status?.platform_info?.packages?.[0]?.version ?? 'null'
          }
          primary="Version"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={gwInfo.description}
          primary="Description"
          {...typographyProps}
        />
      </ListItem>
    </List>
  );
}
