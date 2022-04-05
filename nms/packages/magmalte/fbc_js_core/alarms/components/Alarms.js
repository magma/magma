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
 * @flow
 * @format
 */

import AccountTreeIcon from '@material-ui/icons/AccountTree';
import AlarmContext from './AlarmContext';
import AlertRules from './AlertRules';
import FiringAlerts from './alertmanager/FiringAlerts';
import Grid from '@material-ui/core/Grid';
import GroupIcon from '@material-ui/icons/Group';
import NotificationsActiveIcon from '@material-ui/icons/NotificationsActive';
import React from 'react';
import Receivers from './alertmanager/Receivers/Receivers';
import Routes from './alertmanager/Routes';
import Suppressions from './alertmanager/Suppressions';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import getPrometheusRuleInterface from './rules/PrometheusEditor/getRuleInterface';
import useRouter from '../../../fbc_js_core/ui/hooks/useRouter';
import {Link, Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {matchPath} from 'react-router';

import type {ApiUtil} from './AlarmsApi';
import type {Element} from 'react';
import type {FiringAlarm} from './AlarmAPIType';
import type {Labels} from './AlarmAPIType';
import type {Match} from 'react-router-dom';
import type {RuleInterfaceMap} from './rules/RuleInterface';

const useTabStyles = makeStyles(theme => ({
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

type TabData = {
  icon: Element<*>,
  name: string,
};

type TabMap = {
  [string]: TabData,
};

const TABS: TabMap = {
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
};

const DEFAULT_TAB_NAME = 'alerts';

type Props<TRuleUnion> = {
  //props specific to this component
  makeTabLink: ({match: Match, keyName: string}) => string,
  disabledTabs?: Array<string>,
  // context props
  apiUtil: ApiUtil,
  getNetworkId?: () => string,
  ruleMap?: ?RuleInterfaceMap<TRuleUnion>,
  thresholdEditorEnabled?: boolean,
  alertManagerGlobalConfigEnabled?: boolean,
  filterLabels?: (labels: Labels) => Labels,
  getAlertType?: (alert: FiringAlarm) => string,
  emptyAlerts?: React$Node,
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
  const {match, location} = useRouter();

  const currentTabMatch = matchPath(
    location.pathname,
    `${match.path}/:tabName`,
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
      <Grid container spacing={2} justify="space-between">
        <Grid item xs={12}>
          <Tabs
            value={currentTabMatch?.params?.tabName || DEFAULT_TAB_NAME}
            indicatorColor="primary"
            textColor="primary">
            {Object.keys(TABS).map(keyName => {
              if (disabledTabSet.has(keyName)) {
                return null;
              }
              const {icon, name} = TABS[keyName];
              return (
                <Tab
                  classes={tabStyles}
                  component={Link}
                  to={makeTabLink({keyName, match})}
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
      <Switch>
        <Route
          path={`${match.path}/alerts`}
          render={() => (
            <FiringAlerts
              emptyAlerts={emptyAlerts}
              filterLabels={filterLabels}
            />
          )}
        />
        <Route
          path={`${match.path}/rules`}
          render={() => (
            <AlertRules
              ruleMap={ruleMap}
              thresholdEditorEnabled={thresholdEditorEnabled}
            />
          )}
        />
        <Route
          path={`${match.path}/suppressions`}
          render={() => <Suppressions />}
        />
        <Route path={`${match.path}/routes`} render={() => <Routes />} />
        <Route path={`${match.path}/teams`} render={() => <Receivers />} />
        <Redirect to={`${match.path}/${DEFAULT_TAB_NAME}`} />
      </Switch>
    </AlarmContext.Provider>
  );
}

// merge custom ruleMap with default prometheus rule map
function useMergedRuleMap<TRuleUnion>({
  ruleMap,
  apiUtil,
}: {
  ruleMap: ?RuleInterfaceMap<TRuleUnion>,
  apiUtil: ApiUtil,
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
