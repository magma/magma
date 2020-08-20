/**
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
import type {EnodebContextType} from '../../components/context/EnodebContext';
import type {EnodebInfo} from '../lte/EnodebUtils';
import type {GatewayContextType} from '../../components/context/GatewayContext';
import type {GatewayTierContextType} from '../../components/context/GatewayTierContext';
import type {NetworkType} from '@fbcnms/types/network';
import type {SectionsConfigs} from '../layout/Section';
import type {SubscriberContextType} from '../../components/context/SubscriberContext';
import type {
  lte_gateway,
  mutable_subscriber,
  network_id,
  network_ran_configs,
  subscriber_id,
  tier,
} from '@fbcnms/magma-api';

import AlarmIcon from '@material-ui/icons/Alarm';
import Alarms from '@fbcnms/ui/insights/Alarms/Alarms';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import DashboardIcon from '@material-ui/icons/Dashboard';
import EnodebContext from '../../components/context/EnodebContext';
import Enodebs from './Enodebs';
import EquipmentDashboard from '../../views/equipment/EquipmentDashboard';
import GatewayContext from '../../components/context/GatewayContext';
import GatewayTierContext from '../../components/context/GatewayTierContext';
import Gateways from '../Gateways';
import InitSubscriberState from '../../state/SubscriberState';
import Insights from '@fbcnms/ui/insights/Insights';
import ListIcon from '@material-ui/icons/List';
import Logs from '@fbcnms/ui/insights/Logs/Logs';
import LteConfigure from '../LteConfigure';
import LteDashboard from './LteDashboard';
import LteMetrics from './LteMetrics';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NetworkCheckIcon from '@material-ui/icons/NetworkCheck';
import NetworkDashboard from '../../views/network/NetworkDashboard';
import PeopleIcon from '@material-ui/icons/People';
import PublicIcon from '@material-ui/icons/Public';
import React from 'react';
import RouterIcon from '@material-ui/icons/Router';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import SubscriberContext from '../context/SubscriberContext';
import SubscriberDashboard from '../../views/subscriber/SubscriberOverview';
import Subscribers from '../Subscribers';
import TrafficDashboard from '../../views/traffic/TrafficOverview';
import WifiTetheringIcon from '@material-ui/icons/WifiTethering';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import {
  InitEnodeState,
  InitTierState,
  SetEnodebState,
  SetGatewayState,
  SetTierState,
  UpdateGateway,
} from '../../state/EquipmentState';
import {LTE} from '@fbcnms/types/network';
import {
  getSubscriberGatewayMap,
  setSubscriberState,
} from '../../state/SubscriberState';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

type LteContextType = {
  subscriberCtx: SubscriberContextType,
  enodebCtx: EnodebContextType,
  gatewayCtx: GatewayContextType,
  gatewayTierCtx: GatewayTierContextType,
};

export function useLteContext(
  networkId: ?network_id,
  networkType: ?NetworkType,
  skipContext: boolean,
): ?LteContextType {
  const [subscriberMap, setSubscriberMap] = useState({});
  const [subscriberMetrics, setSubscriberMetrics] = useState({});
  const [enbInfo, setEnbInfo] = useState<{[string]: EnodebInfo}>({});
  const [lteGateways, setLteGateways] = useState<{[string]: lte_gateway}>({});
  const [lteRanConfigs, setLteRanConfigs] = useState<network_ran_configs>({});
  const [tiers, setTiers] = useState<{[string]: tier}>({});
  const [supportedVersions, setSupportedVersions] = useState<Array<string>>([]);
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchLteState = async () => {
      if (networkId == null) {
        return;
      }
      const [
        lteGatewaysResp,
        lteRanConfigsResp,
        stableChannelResp,
      ] = await Promise.allSettled([
        MagmaV1API.getLteByNetworkIdGateways({networkId}),
        MagmaV1API.getLteByNetworkIdCellularRan({networkId}),
        MagmaV1API.getChannelsByChannelId({channelId: 'stable'}),
        InitTierState({networkId, setTiers, enqueueSnackbar}),
        InitEnodeState({networkId, setEnbInfo, enqueueSnackbar}),
        InitSubscriberState({
          networkId,
          setSubscriberMap,
          setSubscriberMetrics,
          enqueueSnackbar,
        }),
      ]);
      if (stableChannelResp.value) {
        setSupportedVersions(
          stableChannelResp.value.supported_versions.reverse(),
        );
      }
      if (lteGatewaysResp.value) {
        setLteGateways(lteGatewaysResp.value);
      }
      if (lteRanConfigsResp.value) {
        setLteRanConfigs(lteRanConfigsResp.value);
      }
      setIsLoading(false);
    };
    if (!skipContext && networkType === LTE) {
      fetchLteState();
    }
  }, [networkId, networkType, skipContext]);

  if (skipContext || networkId == null || isLoading) {
    return null;
  }

  const subscriberCtx = {
    state: subscriberMap,
    metrics: subscriberMetrics,
    gwSubscriberMap: getSubscriberGatewayMap(subscriberMap),
    setState: (key: subscriber_id, value?: mutable_subscriber) =>
      setSubscriberState({
        networkId,
        subscriberMap,
        setSubscriberMap,
        key,
        value,
      }),
  };

  const gatewayCtx = {
    state: lteGateways,
    setState: (key, value?) =>
      SetGatewayState({lteGateways, setLteGateways, networkId, key, value}),
    updateGateway: props =>
      UpdateGateway({networkId, setLteGateways, ...props}),
  };

  const enodebCtx = {
    state: {enbInfo},
    lteRanConfigs: lteRanConfigs,
    setState: (key, value?) =>
      SetEnodebState({enbInfo, setEnbInfo, networkId, key, value}),
  };

  const gatewayTierCtx = {
    state: {supportedVersions, tiers},
    setState: (key, value?) =>
      SetTierState({tiers, setTiers, networkId, key, value}),
  };

  const lteContext = {subscriberCtx, enodebCtx, gatewayCtx, gatewayTierCtx};
  return lteContext;
}

export function getLteSections(
  alertsEnabled: boolean,
  logsEnabled: boolean,
): SectionsConfigs {
  const sections = [
    'map', // landing path
    [
      {
        path: 'map',
        label: 'Map',
        icon: <PublicIcon />,
        component: Insights,
      },
      {
        path: 'metrics',
        label: 'Metrics',
        icon: <ShowChartIcon />,
        component: LteMetrics,
      },
      {
        path: 'subscribers',
        label: 'Subscribers',
        icon: <PeopleIcon />,
        component: Subscribers,
      },
      {
        path: 'gateways',
        label: 'Gateways',
        icon: <CellWifiIcon />,
        component: Gateways,
      },
      {
        path: 'enodebs',
        label: 'eNodeB Devices',
        icon: <SettingsInputAntennaIcon />,
        component: Enodebs,
      },
      {
        path: 'configure',
        label: 'Configure',
        icon: <SettingsCellIcon />,
        component: LteConfigure,
      },
      {
        path: 'alerts',
        label: 'Alerts',
        icon: <AlarmIcon />,
        component: Alarms,
      },
    ],
  ];
  if (logsEnabled) {
    sections[1].splice(2, 0, {
      path: 'logs',
      label: 'Logs',
      icon: <ListIcon />,
      component: Logs,
    });
  }
  if (alertsEnabled) {
    sections[1].splice(2, 0, {
      path: 'alerts',
      label: 'Alerts',
      icon: <AlarmIcon />,
      component: Alarms,
    });
  }
  return sections;
}

export function getLteSectionsV2(
  alertsEnabled: boolean,
  lteContext: ?LteContextType,
): SectionsConfigs {
  const sections = [
    'dashboard', // landing path
    [
      {
        path: 'dashboard',
        label: 'Dashboard',
        icon: <DashboardIcon />,
        component: lteContext
          ? () => (
              <EnodebContext.Provider value={lteContext.enodebCtx}>
                <GatewayContext.Provider value={lteContext.gatewayCtx}>
                  <LteDashboard />
                </GatewayContext.Provider>
              </EnodebContext.Provider>
            )
          : LoadingFiller,
      },
      {
        path: 'equipment',
        label: 'Equipment',
        icon: <RouterIcon />,
        component: lteContext
          ? () => (
              <EnodebContext.Provider value={lteContext.enodebCtx}>
                <GatewayContext.Provider value={lteContext.gatewayCtx}>
                  <GatewayTierContext.Provider
                    value={lteContext.gatewayTierCtx}>
                    <SubscriberContext.Provider
                      value={lteContext.subscriberCtx}>
                      <EquipmentDashboard />
                    </SubscriberContext.Provider>
                  </GatewayTierContext.Provider>
                </GatewayContext.Provider>
              </EnodebContext.Provider>
            )
          : LoadingFiller,
      },
      {
        path: 'network',
        label: 'Network',
        icon: <NetworkCheckIcon />,
        component: lteContext
          ? () => (
              <EnodebContext.Provider value={lteContext.enodebCtx}>
                <GatewayContext.Provider value={lteContext.gatewayCtx}>
                  <SubscriberContext.Provider value={lteContext.subscriberCtx}>
                    <NetworkDashboard />
                  </SubscriberContext.Provider>
                </GatewayContext.Provider>
              </EnodebContext.Provider>
            )
          : LoadingFiller,
      },
      {
        path: 'subscribers',
        label: 'Subscriber',
        icon: <PeopleIcon />,
        component: lteContext
          ? () => (
              <SubscriberContext.Provider value={lteContext.subscriberCtx}>
                <SubscriberDashboard />
              </SubscriberContext.Provider>
            )
          : LoadingFiller,
      },
      {
        path: 'traffic',
        label: 'Traffic',
        icon: <WifiTetheringIcon />,
        component: TrafficDashboard,
      },
      {
        path: 'metrics',
        label: 'Metrics',
        icon: <ShowChartIcon />,
        component: LteMetrics,
      },
    ],
  ];
  if (alertsEnabled) {
    sections[1].splice(2, 0, {
      path: 'alerts',
      label: 'Alerts',
      icon: <AlarmIcon />,
      component: Alarms,
    });
  }
  return sections;
}
