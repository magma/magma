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
import ActionTable from '../../components/ActionTable';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const POLICY_TITLE = 'Policies';
const APN_TITLE = 'APNs';
const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
    color: colors.primary.white,
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: `0 ${theme.spacing(5)}px`,
  },
  tabs: {
    color: colors.primary.white,
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '16px 0 16px 0',
    display: 'flex',
    alignItems: 'center',
  },
  tabIconLabel: {
    marginRight: '8px',
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
  appBarBtnSecondary: {
    color: colors.primary.white,
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));

export default function PolicyDashboard() {
  const classes = useStyles();
  const {relativePath, relativeUrl} = useRouter();

  return (
    <>
      <div className={classes.topBar}>
        <Text color="light" weight="medium">
          Traffic
        </Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container>
          <Grid item xs={6}>
            <Tabs
              value={0}
              indicatorColor="primary"
              TabIndicatorProps={{style: {height: '5px'}}}
              textColor="inherit"
              className={classes.tabs}>
              <Tab
                key="Policy"
                component={NestedRouteLink}
                label={<PolicyDashboardTabLabel />}
                to="/policy"
                className={classes.tab}
              />
              <Tab
                key="APNs"
                component={NestedRouteLink}
                label={<APNDashboardTabLabel />}
                to="/policy"
                className={classes.tab}
              />
            </Tabs>
          </Grid>
          <Grid
            container
            item
            xs={6}
            justify="flex-end"
            alignItems="center"
            spacing={2}>
            <Grid item>
              <Button color="secondary" variant="contained">
                Download Template
              </Button>
            </Grid>

            <Grid item>
              <Button color="secondary" variant="contained">
                Upload CSV
              </Button>
            </Grid>

            <Grid item>
              <Button color="primary" variant="contained">
                Create New Policy
              </Button>
            </Grid>
          </Grid>
        </Grid>
      </AppBar>

      <Switch>
        <Route
          path={relativePath('/policy')}
          component={PolicyDashboardInternal}
        />
        <Redirect to={relativeUrl('/policy')} />
      </Switch>
    </>
  );
}

type PolicyRowType = {
  policyID: string,
  numFlows: number,
  priority: number,
  numSubscribers: number,
  monitoringKey: string,
  rating: string,
  trackingType: string,
};

function PolicyDashboardInternal() {
  const classes = useStyles();
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);

  // this for enabling edit, deactivate actions
  const [_, setCurrRow] = useState<PolicyRowType>({});

  const {response, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPoliciesRulesViewFull,
    {
      networkId: networkId,
    },
  );

  if (isLoading) {
    return <LoadingFiller />;
  }
  const policyRows: Array<PolicyRowType> = response
    ? Object.keys(response).map((policyID: string) => {
        const policyRule = response[policyID];
        return {
          policyID: policyRule.id,
          numFlows: policyRule.flow_list.length,
          priority: policyRule.priority,
          numSubscribers: policyRule.assigned_subscribers?.length ?? 0,
          monitoringKey: policyRule.monitoring_key ?? '',
          rating: policyRule.rating_group?.toString() ?? 'not found',
          trackingType: policyRule.tracking_type ?? 'NO_TRACKING',
        };
      })
    : [];
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <Text key="title">
            <LibraryBooksIcon /> {POLICY_TITLE}
          </Text>
        </Grid>

        <Grid item xs={12}>
          <ActionTable
            data={policyRows}
            columns={[
              {title: 'Policy ID', field: 'policyID'},
              {title: 'Flows', field: 'numFlows', type: 'numeric'},
              {title: 'Priority', field: 'priority', type: 'numeric'},
              {title: 'Subscribers', field: 'numSubscribers', type: 'numeric'},
              {
                title: 'Monitoring Key',
                field: 'monitoringKey',
                render: rowData => {
                  return (
                    <TextField
                      type="password"
                      value={rowData.monitoringKey}
                      InputProps={{
                        disableUnderline: true,
                        readOnly: true,
                      }}
                    />
                  );
                },
              },
              {title: 'Rating', field: 'rating'},
              {title: 'Tracking Type', field: 'trackingType'},
            ]}
            handleCurrRow={(row: PolicyRowType) => setCurrRow(row)}
            menuItems={[{name: 'Edit'}, {name: 'Deactivate'}, {name: 'Remove'}]}
            options={{
              actionsColumnIndex: -1,
              pageSizeOptions: [5, 10],
            }}
          />
        </Grid>
      </Grid>
    </div>
  );
}

function PolicyDashboardTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <LibraryBooksIcon className={classes.tabIconLabel} /> {POLICY_TITLE}
    </div>
  );
}

function APNDashboardTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <RssFeedIcon className={classes.tabIconLabel} /> {APN_TITLE}
    </div>
  );
}
