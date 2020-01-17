/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AlertRules from './AlertRules';
import AppBar from '@material-ui/core/AppBar';
import FiringAlerts from './prometheus/FiringAlerts';
import React from 'react';
import Receivers from './prometheus/Receivers/Receivers';
import Routes from './prometheus/Routes';
import Suppressions from './prometheus/Suppressions';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import {Link, Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {matchPath} from 'react-router';
import {useRouter} from '@fbcnms/ui/hooks';

import type {ApiUtil} from './AlarmsApi';
import type {FiringAlarm, Labels} from './AlarmAPIType';
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
  apiUtil: ApiUtil,
  makeTabLink: ({match: Match, keyName: string}) => string,
  disabledTabs?: Array<string>,
  thresholdEditorEnabled?: boolean,
  filterLabels?: (labels: Labels, alarm: FiringAlarm) => Labels,
  ruleMap?: ?RuleInterfaceMap<TRuleUnion>,
};

export default function Alarms<TRuleUnion>(props: Props<TRuleUnion>) {
  const {
    apiUtil,
    filterLabels,
    makeTabLink,
    disabledTabs,
    thresholdEditorEnabled,
    ruleMap,
  } = props;
  const classes = useStyles();
  const {match, location} = useRouter();

  const currentTabMatch = matchPath(location.pathname, {
    path: `${match.path}/:tabName`,
  });

  const disabledTabSet = React.useMemo(() => {
    return new Set(disabledTabs ?? []);
  }, [disabledTabs]);

  const alarmProps = {apiUtil};
  return (
    <>
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
          render={() => (
            <FiringAlerts {...alarmProps} filterLabels={filterLabels} />
          )}
        />
        <Route
          path={`${match.path}/alert_rules`}
          render={() => (
            <AlertRules
              {...alarmProps}
              ruleMap={ruleMap}
              thresholdEditorEnabled={thresholdEditorEnabled}
            />
          )}
        />
        <Route
          path={`${match.path}/suppressions`}
          render={() => <Suppressions {...alarmProps} />}
        />
        <Route
          path={`${match.path}/routes`}
          render={() => <Routes {...alarmProps} />}
        />
        <Route
          path={`${match.path}/receivers`}
          render={() => <Receivers {...alarmProps} />}
        />
        <Redirect to={`${match.path}/${DEFAULT_TAB_NAME}`} />
      </Switch>
    </>
  );
}
