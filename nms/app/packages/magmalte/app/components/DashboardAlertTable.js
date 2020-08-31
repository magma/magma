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
import type {prom_firing_alert} from '@fbcnms/magma-api';

import ActionTable from './ActionTable';
import CardTitleRow from './layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../theme/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {Alarm} from '@material-ui/icons';
import {colors, typography} from '../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';
import {withStyles} from '@material-ui/core/styles';

const useStyles = makeStyles(theme => ({
  tab: {
    backgroundColor: colors.primary.white,
    borderRadius: '4px 4px 0 0',
    boxShadow: `inset 0 -2px 0 0 ${colors.primary.concrete}`,
    '& + &': {
      marginLeft: '4px',
    },
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
  warn: 'Major',
  major: 'Major',
  minor: 'Minor',
};

type AlertRowType = {
  label: string,
  labelInfo: string,
  annotations: string,
  status: string,
  timingInfo: Date,
};

type AlertTable = {[Severity]: Array<AlertRowType>};

export default function () {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const {isLoading, response} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdAlerts,
    {
      networkId,
    },
  );
  if (isLoading) {
    return <LoadingFiller />;
  }

  const alerts: Array<prom_firing_alert> = response ?? [];
  const data: AlertTable = {Critical: [], Major: [], Minor: [], Other: []};

  alerts.forEach(alert => {
    const sev: Severity = severityMap[alert.labels['severity']] || 'Other';
    data[sev].push({
      label: alert.labels.alertname,
      labelInfo: `${alert.labels.job} - ${alert.labels.instance}`,
      annotations: `${alert.annotations.description} - ${alert.annotations.summary}`,
      status: alert.status.state,
      timingInfo: new Date(alert.startsAt),
    });
  });
  return (
    <>
      <CardTitleRow icon={Alarm} label={`Alerts (${alerts.length})`} />
      <AlertsTabbedTable alerts={data} />
    </>
  );
}
type TabPanelProps = {
  alerts: Array<AlertRowType>,
  label: string,
};

function TabPanel(props: TabPanelProps) {
  const classes = useStyles();
  const {history, match} = useRouter();

  if (props.alerts.length === 0) {
    return (
      <Paper elevation={0}>
        <Grid
          container
          alignItems="center"
          justify="center"
          className={classes.emptyTable}>
          <Grid item xs={12} className={classes.emptyTableContent}>
            <Text variant="body2">You have 0 {props.label} Alerts</Text>
            <Text variant="body3">
              To add alert triggers click
              <Link
                onClick={() => {
                  history.push(
                    match.url.replace(`dashboard/network`, `alerts/alerts`),
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

  return (
    <ActionTable
      data={props.alerts}
      columns={[
        {title: 'Label', field: 'label'},
        {title: 'Label Info', field: 'labelInfo'},
        {title: 'Annotations', field: 'annotations'},
        {title: 'Status', field: 'status'},
        {title: 'Date', field: 'timingInfo', type: 'datetime'},
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5, 10],
        toolbar: false,
        header: false,
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
        {Object.keys(props.alerts).map((k: Severity, idx: number) => {
          return (
            <MagmaTab
              key={idx}
              label={`${props.alerts[k].length} ${k}`}
              className={classes.tab}
            />
          );
        })}
      </MagmaTabs>
      <TabPanel
        label={severityTabs[currTabIndex]}
        alerts={props.alerts[severityTabs[currTabIndex]]}
      />
    </>
  );
}
