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

import type {MetricGraphConfig} from '@fbcnms/ui/insights/Metrics';

import AppBar from '@material-ui/core/AppBar';
import MultiMetrics from './MultiMetrics';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import {Redirect, Route, Switch} from 'react-router-dom';
import {colors} from '../../theme/default';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(_ => ({
  bar: {
    backgroundColor: colors.primary.brightGray,
  },
  tabs: {
    flex: 1,
    color: colors.primary.white,
  },
}));

const CONFIGS: Array<MetricGraphConfig> = [
  {
    basicQueryConfigs: [
      {
        metric: 'solchrg_solar_output_power_192_168_1_12',
        filters: [{name: 'service', value: 'solchrgd'}],
      },
      {
        metric: 'solchrg_solar_input_power_192_168_1_12',
        filters: [{name: 'service', value: 'solchrgd'}],
      },
    ],
    label: 'Solar Power',
    unit: 'Watts',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'solchrg_battery_voltage_filtered_192_168_1_12',
        filters: [{name: 'service', value: 'solchrgd'}],
      },
      {
        metric: 'solchrg_array_voltage_filtered_192_168_1_12',
        filters: [{name: 'service', value: 'solchrgd'}],
      },
      {
        metric: 'solchrg_12V_supply_filtered_192_168_1_12',
        filters: [{name: 'service', value: 'solchrgd'}],
      },
    ],
    label: 'Voltage (filtered)',
    unit: 'Volts',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'solchrg_battery_current_filtered_192_168_1_12',
        filters: [{name: 'service', value: 'solchrgd'}],
      },
      {
        metric: 'solchrg_array_current_filtered_192_168_1_12',
        filters: [{name: 'service', value: 'solchrgd'}],
      },
    ],
    label: 'Current (filtered)',
    unit: 'Amperes',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'solchrg_heatsink_temperature_192_168_1_12',
        filters: [{name: 'service', value: 'solchrgd'}],
      },
      {
        metric: 'solchrg_RTS_temperature_192_168_1_12',
        filters: [{name: 'service', value: 'solchrgd'}],
      },
    ],
    label: 'Temperature',
    unit: 'Celsius',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'measurementsInletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '1'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '1'},
          {name: 'outletId', value: '1'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '1'},
          {name: 'outletId', value: '2'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '1'},
          {name: 'outletId', value: '3'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '1'},
          {name: 'outletId', value: '4'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '1'},
          {name: 'outletId', value: '5'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '1'},
          {name: 'outletId', value: '6'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '1'},
          {name: 'outletId', value: '7'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '1'},
          {name: 'outletId', value: '8'},
        ],
      },
    ],
    label: 'PDU Current Usage',
    unit: 'Ampere',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'measurementsInletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '5'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '5'},
          {name: 'outletId', value: '1'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '5'},
          {name: 'outletId', value: '2'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '5'},
          {name: 'outletId', value: '3'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '5'},
          {name: 'outletId', value: '4'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '5'},
          {name: 'outletId', value: '5'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '5'},
          {name: 'outletId', value: '6'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '5'},
          {name: 'outletId', value: '7'},
        ],
      },
      {
        metric: 'measurementsOutletSensorValue',
        filters: [
          {name: 'service', value: 'snmp'},
          {name: 'pduId', value: '192.168.1.10'},
          {name: 'sensorType', value: '5'},
          {name: 'outletId', value: '8'},
        ],
      },
    ],
    label: 'PDU Active Power',
    unit: 'Watts',
  },
];

function RhinoMetrics() {
  return <MultiMetrics configs={CONFIGS} />;
}

export default function () {
  const classes = useStyles();
  const {match, relativePath, relativeUrl, location} = useRouter();

  const currentTab = findIndex(['gateways', 'network'], route =>
    location.pathname.startsWith(match.url + '/' + route),
  );

  return (
    <>
      <AppBar position="static" color="default" className={classes.bar}>
        <Tabs
          value={currentTab !== -1 ? currentTab : 0}
          indicatorColor="primary"
          textColor="inherit"
          className={classes.tabs}>
          <Tab component={NestedRouteLink} label="Gateways" to="/gateways" />
        </Tabs>
      </AppBar>
      <Switch>
        <Route path={relativePath('/gateways')} component={RhinoMetrics} />
        <Redirect to={relativeUrl('/gateways')} />
      </Switch>
    </>
  );
}
