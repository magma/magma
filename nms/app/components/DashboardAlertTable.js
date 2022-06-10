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
import type {
  prom_alert_labels,
  prom_firing_alert,
} from '../../generated/MagmaAPIBindings';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ActionTable from './ActionTable';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from './layout/CardTitleRow';
import Chip from '@material-ui/core/Chip';
import ErrorIcon from '@material-ui/icons/Error';
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import Grid from '@material-ui/core/Grid';
import InfoIcon from '@material-ui/icons/Info';
import Link from '@material-ui/core/Link';
// $FlowFixMe migrated to typescript
import LoadingFiller from './LoadingFiller';
import MagmaV1API from '../../generated/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../theme/design-system/Text';
import WarningIcon from '@material-ui/icons/Warning';
// $FlowFixMe migrated to typescript
import nullthrows from '../../shared/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPIFlow';

import {Alarm} from '@material-ui/icons';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {REFRESH_INTERVAL} from './context/RefreshContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors, typography} from '../theme/default';
import {intersection} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useNavigate, useParams, useResolvedPath} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  tab: {
    backgroundColor: colors.primary.white,
    borderRadius: '4px 4px 0 0',
    boxShadow: `inset 0 -2px 0 0 ${colors.primary.concrete}`,
    '& + &': {
      marginLeft: '4px',
    },
  },
  tabLabel: {
    padding: '4px 0 4px 0',
    display: 'flex',
    alignItems: 'center',
  },
  emptyTable: {
    backgroundColor: colors.primary.white,
    padding: theme.spacing(4),
    minHeight: '96px',
  },
  emptyTableContent: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    color: colors.primary.comet,
  },
  rowTitle: {
    color: colors.primary.brightGray,
  },
  rowText: {
    color: colors.primary.comet,
  },
  tabIconLabel: {
    marginRight: '8px',
  },
  labelChip: {
    backgroundColor: colors.state.errorFill,
    color: colors.state.error,
    margin: '5px',
  },
}));

const MagmaTabs = withStyles({
  indicator: {
    backgroundColor: colors.secondary.dodgerBlue,
  },
})(Tabs);

const MagmaTab = withStyles({
  root: {
    fontFamily: typography.body1.fontFamily,
    fontWeight: typography.body1.fontWeight,
    fontSize: typography.body1.fontSize,
    lineHeight: typography.body1.lineHeight,
    letterSpacing: typography.body1.letterSpacing,
    color: colors.primary.brightGray,
    textTransform: 'none',
  },
})(Tab);

type Severity = 'Critical' | 'Major' | 'Minor' | 'Other';

const severityMap: {[string]: Severity} = {
  critical: 'Critical',
  page: 'Critical',
  warning: 'Major',
  major: 'Major',
  minor: 'Minor',
  info: 'Other',
  notice: 'Other',
};

type AlertRowType = {
  alertName: string,
  labels: prom_alert_labels,
  status: string,
  service: string,
  gatewayId: string,
  date: Date,
};

type AlertTable = {[Severity]: Array<AlertRowType>};

type DashboardAlertTableProps = {
  labelFilters?: {[string]: string},
};

function checkFilter(
  alert: prom_firing_alert,
  labelFilters?: {[string]: string},
) {
  if (labelFilters) {
    const labels = intersection(
      Object.keys(labelFilters),
      Object.keys(alert.labels),
    );
    if (!labels.length) {
      return false;
    }

    let filtersMatch = true;
    labels.forEach(k => {
      if (alert.labels[k] !== labelFilters[k]) {
        filtersMatch = false;
      }
    });
    return filtersMatch;
  }
  return true;
}

export default function DashboardAlertTable(props: DashboardAlertTableProps) {
  const classes = useStyles();
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const [lastRefreshTime, setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );

  const {isLoading, response} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdAlerts,
    {
      networkId,
    },
    undefined,
    lastRefreshTime,
  );
  useEffect(() => {
    const intervalId = setInterval(
      () => setLastRefreshTime(new Date().toLocaleString()),
      REFRESH_INTERVAL,
    );

    return () => {
      clearInterval(intervalId);
    };
  }, []);

  if (isLoading) {
    return <LoadingFiller />;
  }

  let alerts: Array<prom_firing_alert> = response ?? [];
  const data: AlertTable = {Critical: [], Major: [], Minor: [], Other: []};

  alerts = alerts.filter(alert => checkFilter(alert, props.labelFilters));
  alerts.forEach(alert => {
    const severity = alert.labels?.['severity']?.toLowerCase();
    const sev: Severity = severityMap?.[severity] ?? 'Other';
    data[sev].push({
      alertName: alert.labels?.['alertname'] ?? '-',
      labels: alert.labels,
      gatewayId: alert.labels?.['gatewayID'] ?? '-',
      service: alert.labels?.['service'] ?? '-',
      status: alert.status.state,
      date: new Date(alert.startsAt),
    });
  });

  return (
    <div className={props.labelFilters && classes.dashboardRoot}>
      <CardTitleRow icon={Alarm} label={`Alerts (${alerts.length})`} />
      <AlertsTabbedTable alerts={data} />
    </div>
  );
}

type TabPanelProps = {
  alerts: Array<AlertRowType>,
  label: string,
};

function TabPanel(props: TabPanelProps) {
  const classes = useStyles();
  const resolvedPath = useResolvedPath('');
  const navigate = useNavigate();

  if (props.alerts.length === 0) {
    return (
      <Paper elevation={0}>
        <Grid
          container
          alignItems="center"
          justifyContent="center"
          className={classes.emptyTable}>
          <Grid item xs={12} className={classes.emptyTableContent}>
            <Text variant="body2">You have 0 {props.label} Alerts</Text>
            <Text variant="body3">
              To add alert triggers click
              <Link
                onClick={() => {
                  navigate(
                    resolvedPath.pathname.replace(
                      `dashboard/network`,
                      `alerts/alerts`,
                    ),
                  );
                }}>
                {' '}
                alert settings
              </Link>
            </Text>
          </Grid>
        </Grid>
      </Paper>
    );
  }

  const ignoreLabelList = [
    'networkID',
    'gatewayID',
    'monitor',
    'severity',
    'alertname',
    'service',
  ];
  return (
    <ActionTable
      data={props.alerts}
      columns={[
        {title: 'Date', field: 'date', type: 'datetime', defaultSort: 'desc'},
        {title: 'Status', field: 'status', width: 200},
        {title: 'Alert Name', field: 'alertName', width: 200},
        {title: 'Service', field: 'service', width: 200},
        {title: 'Gateway', field: 'gatewayId', width: 200},
        {
          title: 'Labels',
          field: 'labels',
          render: (currRow: AlertRowType) => (
            <div>
              {Object.keys(currRow.labels)
                .filter(k => !ignoreLabelList.includes(k))
                .map(k => (
                  <Chip
                    key={k}
                    className={classes.labelChip}
                    label={
                      <span>
                        <em>{k}</em>={currRow.labels[k]}
                      </span>
                    }
                    size="small"
                  />
                ))}
            </div>
          ),
        },
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5, 10],
        toolbar: false,
      }}
      localization={{
        header: {actions: ''},
      }}
    />
  );
}

type Props = {
  alerts: AlertTable,
};

function AlertsTabbedTable(props: Props) {
  const classes = useStyles();
  const [currTabIndex, setCurrTabIndex] = useState<number>(0);
  const severityTabs: Array<Severity> = ['Critical', 'Major', 'Minor', 'Other'];

  return (
    <>
      <MagmaTabs
        value={currTabIndex}
        onChange={(_, newIndex: number) => setCurrTabIndex(newIndex)}
        variant="fullWidth">
        <MagmaTab
          key={'severe'}
          label={
            <div className={classes.tabLabel}>
              <ErrorIcon
                style={{color: colors.alerts.severe}}
                className={classes.tabIconLabel}
              />
              {`Critical(${props.alerts['Critical'].length})`}
            </div>
          }
          className={classes.tab}
        />
        <MagmaTab
          key={'major'}
          label={
            <div className={classes.tabLabel}>
              <WarningIcon
                style={{color: colors.alerts.major}}
                className={classes.tabIconLabel}
              />
              {`Major(${props.alerts['Major'].length})`}
            </div>
          }
          className={classes.tab}
        />
        <MagmaTab
          key={'minor'}
          label={
            <div className={classes.tabLabel}>
              <ErrorOutlineIcon
                style={{color: colors.alerts.minor}}
                className={classes.tabIconLabel}
              />
              {`Minor(${props.alerts['Minor'].length})`}
            </div>
          }
          className={classes.tab}
        />
        <MagmaTab
          key={'other'}
          label={
            <div className={classes.tabLabel}>
              <InfoIcon
                style={{color: colors.alerts.other}}
                className={classes.tabIconLabel}
              />
              {`Other(${props.alerts['Other'].length})`}
            </div>
          }
          className={classes.tab}
        />
      </MagmaTabs>
      <TabPanel
        label={severityTabs[currTabIndex]}
        alerts={props.alerts[severityTabs[currTabIndex]]}
      />
    </>
  );
}
