/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {lte_gateway, network_ran_configs, tier} from '@fbcnms/magma-api';

import AddEditEnodeButton from './EnodebDetailConfigEdit';
import AddEditGatewayButton from './GatewayDetailConfigEdit';
import AppBar from '@material-ui/core/AppBar';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import Enodeb from './EquipmentEnodeb';
import EnodebContext from '../../components/context/EnodebContext';
import EnodebDetail from './EnodebDetailMain';
import Gateway from './EquipmentGateway';
import GatewayContext from '../../components/context/GatewayContext';
import GatewayDetail from './GatewayDetailMain';
import GatewayTierContext from '../../components/context/GatewayTierContext';
import Grid from '@material-ui/core/Grid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../../theme/design-system/Text';
import UpgradeButton from './UpgradeTiersDialog';
import nullthrows from '@fbcnms/util/nullthrows';

import {GetCurrentTabPos} from '../../components/TabUtils.js';
import {
  InitEnodeState,
  InitTierState,
  SetEnodebState,
  SetGatewayState,
  SetTierState,
  UpdateGateway,
} from '../../state/EquipmentState';
import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
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
}));

function EquipmentDashboard() {
  const {match, relativePath, relativeUrl} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const [enbInfo, setEnbInfo] = useState<{[string]: EnodebInfo}>({});
  const [isLoading, setIsLoading] = useState(true);
  const [lteGateways, setLteGateways] = useState<{[string]: lte_gateway}>({});
  const [lteRanConfigs, setLteRanConfigs] = useState<network_ran_configs>({});
  const [tiers, setTiers] = useState<{[string]: tier}>({});
  const [supportedVersions, setSupportedVersions] = useState<Array<string>>([]);
  const enqueueSnackbar = useEnqueueSnackbar();
  const enodebCtx = {
    state: {enbInfo, lteRanConfigs},
    setState: (key, value?) =>
      SetEnodebState({enbInfo, setEnbInfo, networkId, key, value}),
  };
  const tierCtx = {
    state: {supportedVersions, tiers},
    setState: (key, value?) =>
      SetTierState({tiers, setTiers, networkId, key, value}),
  };
  const gatewayCtx = {
    state: lteGateways,
    setState: (key, value?) =>
      SetGatewayState({lteGateways, setLteGateways, networkId, key, value}),
    updateGateway: props =>
      UpdateGateway({networkId, setLteGateways, ...props}),
  };

  useEffect(() => {
    const fetchAllData = async () => {
      const [
        lteGatewaysResp,
        lteRanConfigsResp,
        stableChannelResp,
      ] = await Promise.allSettled([
        MagmaV1API.getLteByNetworkIdGateways({networkId}),
        MagmaV1API.getLteByNetworkIdCellularRan({networkId}),
        MagmaV1API.getChannelsByChannelId({channelId: 'stable'}),
        InitEnodeState({networkId, setEnbInfo, enqueueSnackbar}),
        InitTierState({networkId, setTiers, enqueueSnackbar}),
      ]);

      if (lteGatewaysResp.value) {
        setLteGateways(lteGatewaysResp.value);
      }
      if (lteRanConfigsResp.value) {
        setLteRanConfigs(lteRanConfigsResp.value);
      }
      if (stableChannelResp.value) {
        setSupportedVersions(
          stableChannelResp.value.supported_versions.reverse(),
        );
      }
      setIsLoading(false);
    };

    fetchAllData();
  }, [networkId, isLoading, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }
  return (
    <>
      <Switch>
        <Route
          path={relativePath('/overview/gateway/:gatewayId')}
          render={() => (
            <EnodebContext.Provider value={enodebCtx}>
              <GatewayContext.Provider value={gatewayCtx}>
                <GatewayTierContext.Provider value={tierCtx}>
                  <GatewayDetail />
                </GatewayTierContext.Provider>
              </GatewayContext.Provider>
            </EnodebContext.Provider>
          )}
        />
        <Route
          path={relativePath('/overview/enodeb/:enodebSerial')}
          render={() => (
            <EnodebContext.Provider value={enodebCtx}>
              <EnodebDetail />
            </EnodebContext.Provider>
          )}
        />
        <Route
          path={relativePath('/overview')}
          render={() => (
            <EnodebContext.Provider value={enodebCtx}>
              <GatewayContext.Provider value={gatewayCtx}>
                <GatewayTierContext.Provider value={tierCtx}>
                  <EquipmentDashboardInternal />
                </GatewayTierContext.Provider>
              </GatewayContext.Provider>
            </EnodebContext.Provider>
          )}
        />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}

function EquipmentDashboardInternal() {
  const classes = useStyles();
  const {relativePath, relativeUrl, match} = useRouter();
  const tabPos = GetCurrentTabPos(match.url, ['gateway', 'enodeb']);

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
          <Grid item xs={6}>
            <Grid container justify="flex-end" alignItems="center" spacing={2}>
              <Grid item>
                {/* TODO: these button styles need to be localized */}
                {tabPos == 0 && <UpgradeButton />}
              </Grid>
              <Grid item>
                {tabPos == 0 && (
                  <AddEditGatewayButton title="Add New" isLink={false} />
                )}
                {tabPos == 1 && (
                  <AddEditEnodeButton title="Add New" isLink={false} />
                )}
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </AppBar>
      <Switch>
        <Route path={relativePath('/gateway')} render={() => <Gateway />} />
        <Route path={relativePath('/enodeb')} render={() => <Enodeb />} />
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
