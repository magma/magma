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
import type {KPIRows} from '../../components/KPIGrid';
import type {subscriber} from '../../../../../fbcnms-packages/fbcnms-magma-api';

import Button from '@material-ui/core/Button';
import CardHeader from '@material-ui/core/CardHeader';
import Collapse from '@material-ui/core/Collapse';
import Divider from '@material-ui/core/Divider';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import KPIGrid from '../../components/KPIGrid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';

import {CardTitleFilterRow} from '../../components/layout/CardTitleRow';
import {EditSubscriberButton} from './SubscriberAddDialog';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  list: {
    padding: 0,
  },
  kpiLabel: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  kpiValue: {
    color: colors.primary.brightGray,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    width: '100%',
  },
  kpiBox: {
    width: '100%',
    padding: 0,
    '& > div': {
      width: '100%',
    },
  },
}));

export default function SubscriberDetailConfig({
  subscriberInfo,
}: {
  subscriberInfo: subscriber,
}) {
  const classes = useStyles();

  function ConfigFilter() {
    return <Button variant="text">Edit</Button>;
  }

  function TrafficFilter() {
    return <Button variant="text">Edit</Button>;
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6} alignItems="center">
              <Grid container>
                <Grid item xs={6}>
                  <CardTitleFilterRow
                    icon={SettingsIcon}
                    label="Config"
                    filter={ConfigFilter}
                  />
                </Grid>

                <Grid container item xs={6} justify="flex-end">
                  <EditSubscriberButton />
                </Grid>
              </Grid>

              <SubscriberInfoConfig
                readOnly={true}
                subscriberInfo={subscriberInfo}
              />
            </Grid>

            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleFilterRow
                icon={GraphicEqIcon}
                label="Traffic Policy"
                filter={TrafficFilter}
              />
              <SubscriberConfigTrafficPolicy
                readOnly={true}
                subscriberInfo={subscriberInfo}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function SubscriberConfigTrafficPolicy({
  subscriberInfo,
}: {
  subscriberInfo: subscriber,
}) {
  const [open, setOpen] = useState({
    activeAPN: false,
    baseNames: false,
    activePolicies: false,
  });
  const handleCollapse = (config: string) => {
    setOpen({
      ...open,
      [config]: !open[config],
    });
  };
  const classes = useStyles();

  function ListItems(props) {
    return (
      <>
        <ListItem>
          <ListItemText primary={props.data} />
        </ListItem>
        <Divider />
      </>
    );
  }

  function ListNull() {
    return (
      <>
        <ListItem>
          <ListItemText primary="-" />
        </ListItem>
        <Divider />
      </>
    );
  }

  return (
    <List component={Paper} elevation={0} className={classes.list}>
      <ListItem button onClick={() => handleCollapse('activeAPN')}>
        <CardHeader
          title="Active APNs"
          className={classes.kpiBox}
          subheader={subscriberInfo.active_apns?.length || 0}
          titleTypographyProps={{
            variant: 'body3',
            className: classes.kpiLabel,
            title: 'Active APNs',
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.kpiValue,
            title: subscriberInfo.active_apns?.length || 0,
          }}
        />
        {open['activeAPN'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse
        key="activeAPN"
        in={open['activeAPN']}
        timeout="auto"
        unmountOnExit>
        {subscriberInfo.active_apns?.map(data => <ListItems data={data} />) || (
          <ListNull />
        )}
      </Collapse>
      <ListItem button onClick={() => handleCollapse('baseNames')}>
        <CardHeader
          title="Base Names"
          className={classes.kpiBox}
          subheader={subscriberInfo.active_base_names?.length || 0}
          titleTypographyProps={{
            variant: 'body3',
            className: classes.kpiLabel,
            title: 'Base Names',
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.kpiValue,
            title: subscriberInfo.active_base_names?.length || 0,
          }}
        />
        {open['baseNames'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse
        key="baseNames"
        in={open['baseNames']}
        timeout="auto"
        unmountOnExit>
        {subscriberInfo.active_base_names?.map(data => (
          <ListItems data={data} />
        )) || <ListNull />}
      </Collapse>
      <ListItem button onClick={() => handleCollapse('activePolicies')}>
        <CardHeader
          title="Active Policies"
          className={classes.kpiBox}
          subheader={subscriberInfo.active_policies?.length || 0}
          titleTypographyProps={{
            variant: 'body3',
            className: classes.kpiLabel,
            title: 'Active Policies',
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.kpiValue,
            title: subscriberInfo.active_policies?.length || 0,
          }}
        />
        {open['activePolicies'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse
        key="activePolicies"
        in={open['activePolicies']}
        timeout="auto"
        unmountOnExit>
        {subscriberInfo.active_policies?.map(data => (
          <ListItems data={data} />
        )) || <ListNull />}
      </Collapse>
    </List>
  );
}

function SubscriberInfoConfig({subscriberInfo}: {subscriberInfo: subscriber}) {
  const [authKey, _setAuthKey] = useState(subscriberInfo.lte.auth_key);
  const [authOPC, _setAuthOPC] = useState(subscriberInfo.lte.auth_opc ?? false);
  const [dataPlan, _setDataPlan] = useState(subscriberInfo.lte.sub_profile);

  const kpiData: KPIRows[] = [
    [
      {
        category: 'LTE Network Access',
        value: subscriberInfo.lte.state,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Data plan',
        value: dataPlan,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Auth Key',
        value: authKey,
        statusCircle: false,
      },
    ],
  ];

  if (authOPC) {
    kpiData.push([
      {
        category: 'Auth OPC',
        value: authOPC,
        statusCircle: false,
      },
    ]);
  }

  return <KPIGrid data={kpiData} />;
}
