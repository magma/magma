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
 */

import AccountTreeIcon from '@material-ui/icons/AccountTree';
import AlarmContext from './AlarmContext';
import AlertRules from './AlertRules';
import AlertmanagerRoutes from './alertmanager/Routes';
import FiringAlerts from './alertmanager/FiringAlerts';
import Grid from '@material-ui/core/Grid';
import GroupIcon from '@material-ui/icons/Group';
import NotificationsActiveIcon from '@material-ui/icons/NotificationsActive';
import React from 'react';
import Receivers from './alertmanager/Receivers/Receivers';
import Suppressions from './alertmanager/Suppressions';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import getPrometheusRuleInterface from './rules/PrometheusEditor/getRuleInterface';
import {
  Link,
  Navigate,
  Route,
  Routes,
  matchPath,
  useLocation,
  useParams,
  useResolvedPath,
} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';

import {Theme} from '@material-ui/core/styles';
import type {ApiUtil} from './AlarmsApi';
import type {FiringAlarm} from './AlarmAPIType';
import type {Labels} from './AlarmAPIType';
import type {RuleInterfaceMap} from './rules/RuleInterface';

const useTabStyles = makeStyles<Theme>(theme => ({
  root: {
    minWidth: 'auto',
    minHeight: theme.spacing(4),
  },
  wrapper: {
    flexDirection: 'row',
    textTransform: 'capitalize',
    '& svg, .material-icons': {
      marginRight: theme.spacing(1),
    },
  },
}));

const TABS = {
  alerts: {
    name: 'Alerts',
    icon: <NotificationsActiveIcon />,
  },
  rules: {
    name: 'Rules',
    icon: <AccountTreeIcon />,
  },
  suppressions: {
    name: 'Suppressions',
    icon: <React.Fragment />,
  },
  routes: {
    name: 'Routes',
    icon: <React.Fragment />,
  },
  teams: {
    name: 'Teams',
    icon: <GroupIcon />,
  },
} as const;

const DEFAULT_TAB_NAME = 'alerts';

type Props<TRuleUnion> = {
  //props specific to this component
  makeTabLink: (tabData: {networkId?: string; keyName: string}) => string;
  disabledTabs?: Array<string>;
  // context props
  apiUtil: ApiUtil;
  getNetworkId?: () => string;
  ruleMap?: RuleInterfaceMap<TRuleUnion> | null | undefined;
  thresholdEditorEnabled?: boolean;
  alertManagerGlobalConfigEnabled?: boolean;
  filterLabels?: (labels: Labels) => Labels;
  getAlertType?: (alert: FiringAlarm) => string;
  emptyAlerts?: React.ReactNode;
};

export default function Alarms<TRuleUnion>(props: Props<TRuleUnion>) {
  const {
    apiUtil,
    filterLabels,
    makeTabLink,
    getNetworkId,
    disabledTabs,
    thresholdEditorEnabled,
    alertManagerGlobalConfigEnabled,
    ruleMap,
    getAlertType,
    emptyAlerts,
  } = props;
  const tabStyles = useTabStyles();
  const location = useLocation();
  const resolvedPath = useResolvedPath('');
  const params = useParams();

  const currentTabMatch = matchPath(
    location.pathname,
    `${resolvedPath.pathname}/:tabName`,
  );
  const mergedRuleMap = useMergedRuleMap<TRuleUnion>({ruleMap, apiUtil});

  const disabledTabSet = React.useMemo(() => {
    return new Set(disabledTabs ?? []);
  }, [disabledTabs]);

  return (
    <AlarmContext.Provider
      value={{
        apiUtil,
        thresholdEditorEnabled,
        alertManagerGlobalConfigEnabled,
        filterLabels,
        getNetworkId,
        ruleMap: mergedRuleMap,
        getAlertType: getAlertType,
      }}>
      <Grid container spacing={2} justifyContent="space-between">
        <Grid item xs={12}>
          <Tabs
            value={currentTabMatch?.params?.tabName || DEFAULT_TAB_NAME}
            indicatorColor="primary"
            textColor="primary">
            {(Object.keys(TABS) as Array<keyof typeof TABS>).map(keyName => {
              if (disabledTabSet.has(keyName)) {
                return null;
              }
              const {icon, name} = TABS[keyName];
              return (
                <Tab
                  classes={tabStyles}
                  component={Link}
                  to={makeTabLink({keyName, networkId: params.networkId})}
                  key={keyName}
                  icon={icon}
                  label={name}
                  value={keyName}
                />
              );
            })}
          </Tabs>
        </Grid>
      </Grid>
      <Routes>
        <Route
          path="/alerts"
          element={<FiringAlerts emptyAlerts={emptyAlerts} />}
        />
        <Route path="/rules" element={<AlertRules />} />
        <Route path="/suppressions" element={<Suppressions />} />
        <Route path="/routes" element={<AlertmanagerRoutes />} />
        <Route path="/teams" element={<Receivers />} />
        <Route index element={<Navigate to={DEFAULT_TAB_NAME} replace />} />
      </Routes>
    </AlarmContext.Provider>
  );
}

// merge custom ruleMap with default prometheus rule map
function useMergedRuleMap<TRuleUnion>({
  ruleMap,
  apiUtil,
}: {
  ruleMap: RuleInterfaceMap<TRuleUnion> | null | undefined;
  apiUtil: ApiUtil;
}): RuleInterfaceMap<TRuleUnion> {
  const mergedRuleMap = React.useMemo<RuleInterfaceMap<TRuleUnion>>(
    () =>
      Object.assign(
        {},
        getPrometheusRuleInterface({apiUtil: apiUtil}),
        ruleMap || {},
      ),
    [ruleMap, apiUtil],
  );
  return mergedRuleMap;
}
