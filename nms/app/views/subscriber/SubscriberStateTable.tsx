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
 */
import type {Subscriber, SubscriberState} from '../../../generated';

import ActionTable from '../../components/ActionTable';
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import CardTitleRow from '../../components/layout/CardTitleRow';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Link from '@mui/material/Link';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import PeopleIcon from '@mui/icons-material/People';
import React from 'react';
import SubscriberContext from '../../context/SubscriberContext';
import Text from '../../theme/design-system/Text';

import {Column} from '@material-table/core';
import {JsonDialog, RenderLink} from './SubscriberUtils';
import {REFRESH_INTERVAL} from '../../context/AppContext';
import {Theme} from '@mui/material/styles';
import {colors} from '../../theme/default';
import {makeStyles} from '@mui/styles';
import {useContext, useState} from 'react';
import {useInterval} from '../../hooks';

const useStyles = makeStyles<Theme>(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
}));
export type SubscriberRowType = {
  name: string;
  imsi: string;
  activeApns?: string;
  ipAddresses?: string;
  activeSessions?: number;
  service: string;
  currentUsage: string;
  dailyAvg: string;
  lastReportedTime: Date | string;
};

type SubscriberSessionRowType = {
  apnName: string;
  sessionId: string;
  ipAddr: string;
  state: string;
  activeDuration: string;
  activePolicies: Array<{id: string}>;
};
type SubscriberStateDetailPanelProps = {
  rowData: SubscriberRowType;
};
function SubscriberStateDetailPanel(props: SubscriberStateDetailPanelProps) {
  const ctx = useContext(SubscriberContext);
  const sessionState: Record<string, SubscriberState> = ctx.sessionState;
  const subscriber: Record<string, any> =
    sessionState[props.rowData.imsi]?.subscriber_state || {};
  const subscriberSessionRows: Array<SubscriberSessionRowType> = [];
  Object.keys(subscriber).map((apn: string) => {
    /* eslint-disable @typescript-eslint/restrict-template-expressions,@typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-argument,@typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access */
    subscriber[apn].map((infos: any) => {
      subscriberSessionRows.push({
        apnName: apn,
        sessionId: infos.session_id,
        ipAddr: infos.ipv4 ?? '-',
        state: infos.lifecycle_state,
        activeDuration: `${infos.active_duration_sec} sec`,
        activePolicies: infos.active_policy_rules,
      });
      /* eslint-enable @typescript-eslint/restrict-template-expressions,@typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-argument,@typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access */
    });
  });

  return (
    <ActionTable
      data-testid="detailPanel"
      title=""
      data={subscriberSessionRows}
      columns={[
        {title: 'APN Name', field: 'apnName'},
        {title: 'Session ID', field: 'sessionId'},
        {title: 'State', field: 'state'},
        {title: 'IP Address', field: 'ipAddr'},
        {
          title: 'Active Duration',
          field: 'activeDuration',
        },
        {
          title: 'Active Policy IDs',
          field: 'activePolicies',
          render: currRow =>
            currRow.activePolicies.length ? (
              <List>
                {currRow.activePolicies.map(policy => (
                  <ListItem key={policy.id}>
                    <Link underline="hover">{policy.id} </Link>
                  </ListItem>
                ))}
              </List>
            ) : (
              <Text>{'-'}</Text>
            ),
        },
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5],
        toolbar: false,
        paging: false,
        rowStyle: {background: '#f7f7f7'},
        headerStyle: {
          background: '#f7f7f7',
          color: colors.primary.comet,
        },
      }}
    />
  );
}

export default function SubscriberStateTable() {
  const [currRow, setCurrRow] = useState<SubscriberRowType>(
    {} as SubscriberRowType,
  );
  const classes = useStyles();
  const ctx = useContext(SubscriberContext);
  const subscriberMap: Record<string, Subscriber> = ctx.state;
  const sessionState: Record<string, SubscriberState> = ctx.sessionState;
  const subscriberMetrics = ctx.metrics;
  const [jsonDialog, setJsonDialog] = useState(false);
  const subscribersIds = Object.keys(sessionState);
  const [refresh, setRefresh] = useState(true);
  // Auto refresh subscribers sessions every 30 seconds
  useInterval(
    () => ctx.refetchSessionState(),
    refresh ? REFRESH_INTERVAL : null,
  );

  const tableData: Array<SubscriberRowType> = subscribersIds.map(
    (imsi: string) => {
      const subscriberInfo = subscriberMap[imsi] || {};
      const metrics = subscriberMetrics?.[`${imsi}`];
      const subscriber: Record<string, any> =
        ctx.sessionState?.[imsi]?.subscriber_state || {};
      const ipAddresses: Array<string> = [];
      const activeApns: Array<string> = [];
      let activeSessions = 0;
      Object.keys(subscriber || {}).forEach((apn: string) => {
        /* eslint-disable @typescript-eslint/no-unsafe-argument,@typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access */
        subscriber[apn].forEach((session: any) => {
          if (session.lifecycle_state === 'SESSION_ACTIVE') {
            ipAddresses.push(session?.ipv4);
            activeSessions++;
          }
          /* eslint-enable @typescript-eslint/no-unsafe-argument,@typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access */
        });
        activeApns.push(apn);
      });
      return {
        name: subscriberInfo.name ?? imsi,
        imsi: imsi,
        service: subscriberInfo.lte?.state || '',
        currentUsage: metrics?.currentUsage ?? '0',
        activeApns: activeApns.length > 0 ? activeApns.join() : '-',
        activeSessions: activeSessions,
        ipAddress: ipAddresses.length > 0 ? ipAddresses.join() : '-',
        dailyAvg: metrics?.dailyAvg ?? '0',
        lastReportedTime:
          subscriberInfo.monitoring?.icmp?.last_reported_time === 0
            ? new Date(subscriberInfo.monitoring?.icmp?.last_reported_time)
            : '-',
      };
    },
  );

  const tableColumns: Array<Column<SubscriberRowType>> = [
    {
      title: 'Name',
      field: 'name',
    },
    {
      title: 'IMSI',
      field: 'imsi',
      render: (currRow: SubscriberRowType) => {
        const subscriberConfig = subscriberMap[currRow.imsi];
        return (
          <RenderLink subscriberConfig={subscriberConfig} currRow={currRow} />
        );
      },
    },
    {
      title: 'Service',
      field: 'service',
      width: 100,
    },
    {
      title: 'Current Usage',
      field: 'currentUsage',
      width: 175,
    },
    {
      title: 'Daily Average',
      field: 'dailyAvg',
      width: 175,
    },
    {
      title: 'Last Reported Time',
      field: 'lastReportedTime',
      type: 'datetime',
      width: 200,
    },
    {
      title: 'Active Sessions',
      field: 'activeSessions',
      width: 175,
    },
    {
      title: 'Active APNs',
      field: 'activeApns',
    },
    {
      title: 'Session IP Address',
      field: 'ipAddress',
    },
  ];
  const onClose = () => setJsonDialog(false);
  return (
    <>
      <div className={classes.dashboardRoot}>
        <CardTitleRow
          key="title"
          icon={PeopleIcon}
          label={'Subscriber Sessions'}
          filter={() => (
            <AutorefreshCheckbox
              autorefreshEnabled={refresh}
              onToggle={() => setRefresh(current => !current)}
            />
          )}
        />
        <JsonDialog open={jsonDialog} onClose={onClose} imsi={currRow.imsi} />
        <ActionTable
          data={tableData}
          columns={tableColumns}
          handleCurrRow={(row: SubscriberRowType) => setCurrRow(row)}
          menuItems={[
            {
              name: 'View JSON',
              handleFunc: () => {
                setJsonDialog(true);
              },
            },
          ]}
          options={{
            actionsColumnIndex: -1,
            pageSize: 10,
            pageSizeOptions: [10, 20],
          }}
          detailPanel={[
            {
              icon: () => {
                return <ExpandMore data-testid="details" />;
              },
              openIcon: ExpandLess,
              render: rowData => <SubscriberStateDetailPanel {...rowData} />,
            },
          ]}
        />
      </div>
    </>
  );
}
