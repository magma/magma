/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AlarmContext from './AlarmContext';
import AlertRules from './AlertRules';
import AppBar from '@material-ui/core/AppBar';
import FiringAlerts from './prometheus/FiringAlerts';
import React from 'react';
import Receivers from './prometheus/Receivers/Receivers';
import Routes from './prometheus/Routes';
import Suppressions from './prometheus/Suppressions';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import getPrometheusRuleInterface from './rules/PrometheusEditor/getRuleInterface';
import {Link, Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {matchPath} from 'react-router';
import {useRouter} from '@fbcnms/ui/hooks';

import type {ApiUtil} from './AlarmsApi';
import type {FiringAlarm} from './AlarmAPIType';
import type {Labels} from './AlarmAPIType';
import type {Match} from 'react-router-dom';
import type {RuleInterfaceMap} from './rules/RuleInterface';

const useStyles = makeStyles(_theme => ({
  appBar: {
    position: 'inherit',
  },
}));

type TabMap = {
  [string]: {name: string},
};

const TABS: TabMap = {
  alerts: {
    name: 'Alerts',
  },
  alert_rules: {
    name: 'Alert Rules',
  },
  suppressions: {
    name: 'Suppressions',
  },
  routes: {
    name: 'Routes',
  },
  receivers: {
    name: 'Receivers',
  },
};

const DEFAULT_TAB_NAME = 'alerts';

type Props<TRuleUnion> = {
  //props specific to this component
  makeTabLink: ({match: Match, keyName: string}) => string,
  disabledTabs?: Array<string>,
  // context props
  apiUtil: ApiUtil,
  ruleMap?: ?RuleInterfaceMap<TRuleUnion>,
  thresholdEditorEnabled?: boolean,
  filterLabels?: (labels: Labels) => Labels,
  getAlertType?: (alert: FiringAlarm) => string,
};

export default function Alarms<TRuleUnion>(props: Props<TRuleUnion>) {
  const {
    apiUtil,
    filterLabels,
    makeTabLink,
    disabledTabs,
    thresholdEditorEnabled,
    ruleMap,
    getAlertType,
  } = props;
  const classes = useStyles();
  const {match, location} = useRouter();

  const currentTabMatch = matchPath(location.pathname, {
    path: `${match.path}/:tabName`,
  });
  const mergedRuleMap = useMergedRuleMap<TRuleUnion>({ruleMap, apiUtil});

  const disabledTabSet = React.useMemo(() => {
    return new Set(disabledTabs ?? []);
  }, [disabledTabs]);

  return (
    <AlarmContext.Provider
      value={{
        apiUtil,
        thresholdEditorEnabled,
        filterLabels,
        ruleMap: mergedRuleMap,
        getAlertType: getAlertType,
      }}>
      <AppBar className={classes.appBar} color="default">
        <Tabs
          value={currentTabMatch?.params?.tabName || 'alerts'}
          indicatorColor="primary"
          textColor="primary">
          {Object.keys(TABS).map(keyName => {
            if (disabledTabSet.has(keyName)) {
              return null;
            }
            return (
              <Tab
                component={Link}
                to={makeTabLink({keyName, match})}
                key={keyName}
                className={classes.selectedTab}
                label={TABS[keyName].name}
                value={keyName}
              />
            );
          })}
        </Tabs>
      </AppBar>

      <Switch>
        <Route
          path={`${match.path}/alerts`}
          render={() => <FiringAlerts filterLabels={filterLabels} />}
        />
        <Route
          path={`${match.path}/alert_rules`}
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
        <Route path={`${match.path}/receivers`} render={() => <Receivers />} />
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
