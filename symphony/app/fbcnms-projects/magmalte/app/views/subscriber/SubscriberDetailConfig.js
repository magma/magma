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
import type {subscriber} from '../../../../../fbcnms-packages/fbcnms-magma-api';

import Collapse from '@material-ui/core/Collapse';
import Divider from '@material-ui/core/Divider';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
}));

export default function SubscriberDetailConfig({
  subscriberInfo,
}: {
  subscriberInfo: subscriber,
}) {
  const classes = useStyles();
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3} alignItems="stretch">
        <Grid container spacing={3} alignItems="stretch" item xs={12}>
          <Grid item xs={6}>
            <Grid container>
              <Grid item xs={6}>
                <Text>
                  <SettingsIcon /> Config
                </Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <SubscriberInfoConfig
              readOnly={true}
              subscriberInfo={subscriberInfo}
            />
          </Grid>

          <Grid item xs={6}>
            <Grid container>
              <Grid item xs={6}>
                <Text>
                  <GraphicEqIcon />
                  Traffic Policy
                </Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <SubscriberConfigTrafficPolicy
              readOnly={true}
              subscriberInfo={subscriberInfo}
            />
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function SubscriberConfigTrafficPolicy({
  subscriberInfo,
  readOnly,
}: {
  subscriberInfo: subscriber,
  readOnly: boolean,
}) {
  const [open, setOpen] = useState({
    activeAPN: true,
    baseNames: true,
    activePolicies: true,
  });
  const handleCollapse = (config: string) => {
    setOpen({
      ...open,
      [config]: !open[config],
    });
  };
  return (
    <List component={Paper}>
      <ListItem button onClick={() => handleCollapse('activeAPN')}>
        <TextField
          fullWidth={true}
          value={subscriberInfo.active_apns?.length || 0}
          label="Active APNs"
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
        {open['activeAPN'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse
        key="activeAPN"
        in={open['activeAPN']}
        timeout="auto"
        unmountOnExit>
        <ListItem>
          <TextField
            fullWidth={true}
            value={subscriberInfo.active_apns?.join(', ') || 0}
            InputProps={{disableUnderline: true, readOnly: readOnly}}
          />
        </ListItem>
        <Divider />
      </Collapse>
      <ListItem button onClick={() => handleCollapse('baseNames')}>
        <TextField
          fullWidth={true}
          value={subscriberInfo.active_base_names?.length || 0}
          label="Base Names"
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
        {open['baseNames'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse
        key="baseNames"
        in={open['baseNames']}
        timeout="auto"
        unmountOnExit>
        <ListItem>
          <TextField
            fullWidth={true}
            value={subscriberInfo.active_base_names?.join(', ') || 0}
            InputProps={{disableUnderline: true, readOnly: readOnly}}
          />
        </ListItem>
        <Divider />
      </Collapse>
      <ListItem button onClick={() => handleCollapse('activePolicies')}>
        <TextField
          fullWidth={true}
          value={subscriberInfo.active_policies?.length || 0}
          label="Active Policies"
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
        {open['activePolicies'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Collapse
        key="activePolicies"
        in={open['activePolicies']}
        timeout="auto"
        unmountOnExit>
        <ListItem>
          <TextField
            fullWidth={true}
            value={subscriberInfo.active_policies?.join(', ') || 0}
            InputProps={{disableUnderline: true, readOnly: readOnly}}
          />
        </ListItem>
      </Collapse>
    </List>
  );
}

function SubscriberInfoConfig({
  subscriberInfo,
  readOnly,
}: {
  subscriberInfo: subscriber,
  readOnly: boolean,
}) {
  const [authKey, setAuthKey] = useState(subscriberInfo.lte.auth_key);
  const [authOPC, setAuthOPC] = useState(subscriberInfo.lte.auth_opc ?? false);
  const [dataPlan, setDataPlan] = useState(subscriberInfo.lte.sub_profile);

  return (
    <List component={Paper}>
      <ListItem>
        <TextField
          fullWidth={true}
          value={subscriberInfo.lte.state}
          label="LTE Network Access"
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <TextField
          fullWidth={true}
          value={dataPlan}
          label="Data plan"
          onChange={({target}) => setDataPlan(target.value)}
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <TextField
          type="password"
          fullWidth={true}
          value={authKey}
          label="Auth Key"
          onChange={({target}) => setAuthKey(target.value)}
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
        <Divider />
        {authOPC && (
          <TextField
            type="password"
            fullWidth={true}
            value={authOPC}
            label="Auth OPC"
            onChange={({target}) => setAuthOPC(target.value)}
            InputProps={{disableUnderline: true, readOnly: readOnly}}
          />
        )}
      </ListItem>
    </List>
  );
}
