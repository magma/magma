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
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {lte_gateway} from '@fbcnms/magma-api';

import CellWifiIcon from '@material-ui/icons/CellWifi';
import Enodeb from './EquipmentEnodeb';
import EnodebDetail from './EnodebDetailMain';
import Gateway from './EquipmentGateway';
import GatewayDetail from './GatewayDetailMain';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import nullthrows from '@fbcnms/util/nullthrows';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import Text from '../../theme/design-system/Text';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {AppBar, Button, Grid, Tab, Tabs} from '@material-ui/core';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {Redirect, Route, Switch} from 'react-router-dom';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

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

function EquipmentDashboard() {
  const {match, relativePath, relativeUrl} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const [enbInfo, setEnbInfo] = useState<{[string]: EnodebInfo}>({});
  const [isEnbStLoading, setIsEnbStLoading] = useState(true);

  const enqueueSnackbar = useEnqueueSnackbar();

  const {response: lteGatwayResp, isLoading: isLteRespLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdGateways,
    {
      networkId: networkId,
    },
  );

  const {response: enb, isLoading: isEnbRespLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdEnodebs,
    {
      networkId: networkId,
    },
  );

  useEffect(() => {
    const fetchEnodebState = async () => {
      let err = false;
      if (!enb) {
        return;
      }
      const requests = Object.keys(enb).map(async k => {
        const {serial} = enb[k];
        try {
          // eslint-disable-next-line max-len
          const enbSt = await MagmaV1API.getLteByNetworkIdEnodebsByEnodebSerialState(
            {
              networkId: networkId,
              enodebSerial: serial,
            },
          );
          return {serial, enbSt};
        } catch (error) {
          err = true;
          console.error('error getting enodeb status for ' + serial);
          return null;
        }
      });
      if (err) {
        enqueueSnackbar(
          'There was a problem fetching enodeb state from the server',
          {variant: 'error'},
        );
      }
      Promise.all(requests).then(allResponses => {
        const enbInfoLocal = {};
        allResponses.filter(Boolean).forEach(r => {
          enbInfoLocal[r.serial] = {
            enb: enb[r.serial],
            enb_state: r.enbSt,
          };
        });
        setEnbInfo(enbInfoLocal);
        setIsEnbStLoading(false);
      });
    };
    if (!enb && !isEnbRespLoading) {
      setIsEnbStLoading(false);
      return;
    }
    fetchEnodebState();
  }, [networkId, enb, isEnbRespLoading, enqueueSnackbar]);

  if (isLteRespLoading || isEnbStLoading) {
    return <LoadingFiller />;
  }
  const lteGateways: {[string]: lte_gateway} = lteGatwayResp ?? {};
  return (
    <>
      <Switch>
        <Route
          path={relativePath('/overview/gateway/:gatewayId')}
          render={() => (
            <GatewayDetail lteGateways={lteGateways} enbInfo={enbInfo} />
          )}
        />
        <Route
          path={relativePath('/overview/enodeb/:enodebSerial')}
          render={() => <EnodebDetail enbInfo={enbInfo} />}
        />
        <Route
          path={relativePath('/overview')}
          render={() => (
            <EquipmentDashboardInternal
              enbInfo={enbInfo}
              lteGateways={lteGateways}
            />
          )}
        />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}

function EquipmentDashboardInternal({
  lteGateways,
  enbInfo,
}: {
  lteGateways: {[string]: lte_gateway},
  enbInfo: {[string]: EnodebInfo},
}) {
  const classes = useStyles();
  const {relativePath, relativeUrl} = useRouter();
  const [tabPos, setTabPos] = React.useState(0);
  return (
    <>
      <div className={classes.topBar}>
        <Text variant="body2">Equipment</Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container direction="row" justify="flex-end" alignItems="center">
          <Grid item xs={6}>
            <Tabs
              value={tabPos}
              onChange={(_, v) => setTabPos(v)}
              indicatorColor="primary"
              TabIndicatorProps={{style: {height: '5px'}}}
              textColor="inherit"
              className={classes.tabs}>
              <Tab
                key="Gateways"
                component={NestedRouteLink}
                label={<GatewayTabLabel />}
                to="/gateway"
                className={classes.tab}
              />
              <Tab
                key="EnodeBs"
                component={NestedRouteLink}
                label={<EnodebTabLabel />}
                to="/enodeb"
                className={classes.tab}
              />
            </Tabs>
          </Grid>
          <Grid
            item
            xs={6}
            direction="row"
            justify="flex-end"
            alignItems="center">
            <Grid container justify="flex-end" alignItems="center" spacing={2}>
              <Grid item>
                {/* TODO: these button styles need to be localized */}
                <Button variant="text" className={classes.appBarBtnSecondary}>
                  Secondary Action
                </Button>
              </Grid>
              <Grid item>
                <Button variant="contained" className={classes.appBarBtn}>
                  Create New
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </AppBar>
      <Switch>
        <Route
          path={relativePath('/gateway')}
          render={() => <Gateway lteGateways={lteGateways} />}
        />
        <Route
          path={relativePath('/enodeb')}
          render={() => <Enodeb enbInfo={enbInfo} />}
        />
        <Redirect to={relativeUrl('/gateway')} />
      </Switch>
    </>
  );
}

function GatewayTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <CellWifiIcon className={classes.tabIconLabel} /> Gateway
    </div>
  );
}

function EnodebTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <SettingsInputAntennaIcon className={classes.tabIconLabel} /> eNodeB
    </div>
  );
}

export default EquipmentDashboard;
