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
import type {gateway_id, lte_gateway, tier, tier_id} from '@fbcnms/magma-api';

import AddEditEnodeButton from './EnodebDetailConfigEdit';
import AddEditGatewayButton from './GatewayDetailConfigEdit';
import AppBar from '@material-ui/core/AppBar';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import Enodeb from './EquipmentEnodeb';
import EnodebDetail from './EnodebDetailMain';
import Gateway from './EquipmentGateway';
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
  const [lteGateways, setLteGatways] = useState<{[string]: lte_gateway}>({});
  const [tiers, setTiers] = useState<{[string]: tier}>({});
  const [supportedVersions, setSupportedVersions] = useState<Array<string>>([]);
  const enqueueSnackbar = useEnqueueSnackbar();

  const updateTier = async (key: tier_id, value: tier) => {
    if (key in tiers) {
      await MagmaV1API.putNetworksByNetworkIdTiersByTierId({
        networkId: networkId,
        tierId: key,
        tier: value,
      });
    } else {
      await MagmaV1API.postNetworksByNetworkIdTiers({
        networkId: networkId,
        tier: value,
      });
    }
    setTiers({...tiers, [key]: value});
  };

  const removeTier = async (key: tier_id) => {
    await MagmaV1API.deleteNetworksByNetworkIdTiersByTierId({
      networkId: networkId,
      tierId: key,
    });
    const {key: _, ...newTiers} = tiers;
    setTiers(newTiers);
  };

  const updateGatewayTier = async (gatewayId: gateway_id, tierId: tier_id) => {
    await MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdTier({
      networkId,
      gatewayId: gatewayId,
      tierId: JSON.stringify(`"${tierId}"`),
    });
    const gateways = await MagmaV1API.getLteByNetworkIdGateways({
      networkId,
    });
    setLteGatways(gateways);
  };

  useEffect(() => {
    const fetchAllNetworkUpgradeTiers = async () => {
      let tierIdList = [];
      try {
        tierIdList = await MagmaV1API.getNetworksByNetworkIdTiers({networkId});
      } catch (e) {
        enqueueSnackbar('failed fetching tier information', {variant: 'error'});
      }

      const requests = tierIdList.map(tierId => {
        try {
          return MagmaV1API.getNetworksByNetworkIdTiersByTierId({
            networkId,
            tierId,
          });
        } catch (e) {
          enqueueSnackbar('failed fetching tier information for ' + tierId, {
            variant: 'error',
          });
          return;
        }
      });

      return await Promise.all(requests);
    };

    const fetchEnodebState = async () => {
      let enb = {};
      try {
        enb = await MagmaV1API.getLteByNetworkIdEnodebs({networkId});
      } catch (e) {
        enqueueSnackbar('failed fetching enodeb information', {
          variant: 'error',
        });
        return [];
      }

      let err = false;
      const requests = Object.keys(enb).map(async k => {
        try {
          const {serial} = enb[k];
          // eslint-disable-next-line max-len
          const enbSt = await MagmaV1API.getLteByNetworkIdEnodebsByEnodebSerialState(
            {
              networkId: networkId,
              enodebSerial: serial,
            },
          );
          return [enb[k], enbSt];
        } catch (e) {
          err = true;
          return [enb[k], {}];
        }
      });
      if (err) {
        enqueueSnackbar('failed fetching enodeb state information', {
          variant: 'error',
        });
      }
      return await Promise.all(requests);
    };

    const fetchAllData = async () => {
      const [
        lteGateways,
        enbResp,
        tierResponse,
        stableChannel,
      ] = await Promise.all([
        MagmaV1API.getLteByNetworkIdGateways({networkId}),
        fetchEnodebState(),
        fetchAllNetworkUpgradeTiers(),
        MagmaV1API.getChannelsByChannelId({channelId: 'stable'}),
      ]);
      const enbInfoLocal = {};
      enbResp.filter(Boolean).forEach(r => {
        if (r.length > 0) {
          const [enb, enbSt] = r;
          if (enb != null && enbSt != null) {
            enbInfoLocal[enb.serial] = {
              enb: enb,
              enb_state: enbSt,
            };
          }
        }
      });

      const tiers = {};
      // reduce function gives a flow lint, hence using forEach instead
      tierResponse.filter(Boolean).forEach(item => {
        tiers[item.id] = item;
      });

      setTiers(tiers);
      setLteGatways(lteGateways);
      setEnbInfo(enbInfoLocal);
      setSupportedVersions(stableChannel.supported_versions.reverse());
      setIsLoading(false);
    };
    fetchAllData().catch(e => {
      enqueueSnackbar(e?.message, {variant: 'error'});
      setIsLoading(false);
    });
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
            <GatewayTierContext.Provider
              value={{
                supportedVersions: supportedVersions,
                tiers: tiers,
                updateTier: updateTier,
                removeTier: removeTier,
                updateGatewayTier: updateGatewayTier,
              }}>
              <EquipmentDashboardInternal
                enbInfo={enbInfo}
                lteGateways={lteGateways}
              />
            </GatewayTierContext.Provider>
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
