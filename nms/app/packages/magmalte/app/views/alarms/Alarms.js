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
import type {ApiUtil} from '@fbcnms/alarms/components/AlarmsApi';
import type {FiringAlarm} from '@fbcnms/alarms/components/AlarmAPIType';
import type {Labels} from '@fbcnms/alarms/components/AlarmAPIType';
import type {RuleInterfaceMap} from '@fbcnms/alarms/components/rules/RuleInterface';

import AddAlertIcon from '@material-ui/icons/AddAlert';
import AlarmContext from '@fbcnms/alarms/components/AlarmContext';
import AlarmIcon from '@material-ui/icons/Alarm';
import AlertRules from './AlertRules';
import ContactMailIcon from '@material-ui/icons/ContactMail';
import FiringAlerts from '@fbcnms/alarms/components/alertmanager/FiringAlerts';
import React from 'react';
import Receivers from '@fbcnms/alarms/components/alertmanager/Receivers/Receivers';
import TopBar from '../../components/TopBar';
import getPrometheusRuleInterface from '@fbcnms/alarms/components/rules/PrometheusEditor/getRuleInterface';

import {MagmaAlarmsApiUtil} from '../../state/AlarmsApiUtil';
import {Redirect, Route, Switch} from 'react-router-dom';
import {useRouter} from '@fbcnms/ui/hooks';

const DEFAULT_TAB_NAME = 'alerts';

type Props<TRuleUnion> = {
  //props specific to this component
  // context props
  ruleMap?: ?RuleInterfaceMap<TRuleUnion>,
  thresholdEditorEnabled?: boolean,
  alertManagerGlobalConfigEnabled?: boolean,
  filterLabels?: (labels: Labels) => Labels,
  getAlertType?: (alert: FiringAlarm) => string,
};

export default function AlarmsDashboard<TRuleUnion>(props: Props<TRuleUnion>) {
  const {
    filterLabels,
    alertManagerGlobalConfigEnabled,
    ruleMap,
    getAlertType,
  } = props;
  const {match} = useRouter();
  const apiUtil = MagmaAlarmsApiUtil;
  const thresholdEditorEnabled = true;
  const mergedRuleMap = useMergedRuleMap<TRuleUnion>({ruleMap, apiUtil});

  const tabs = [
    {
      to: '/alerts',
      label: 'Alerts',
      icon: AlarmIcon,
    },
    {
      to: '/alert_rules',
      label: 'Alert Rules',
      icon: AddAlertIcon,
    },
    {
      to: 'receivers',
      label: 'Receivers',
      icon: ContactMailIcon,
    },
  ];

  return (
    <AlarmContext.Provider
      value={{
        apiUtil,
        thresholdEditorEnabled,
        alertManagerGlobalConfigEnabled,
        filterLabels,
        ruleMap: mergedRuleMap,
        getAlertType: getAlertType,
      }}>
      <TopBar header={'Alarms'} tabs={tabs} />

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
