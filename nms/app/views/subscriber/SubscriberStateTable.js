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
import type {
  subscriber,
  subscriber_state,
} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ActionTable from '../../components/ActionTable';
// $FlowFixMe migrated to typescript
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import Link from '@material-ui/core/Link';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
// $FlowFixMe migrated to typescript
import NetworkContext from '../../components/context/NetworkContext';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
// $FlowFixMe migrated to typescript
import SubscriberContext from '../../components/context/SubscriberContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {
  REFRESH_INTERVAL,
  useRefreshingContext,
  // $FlowFixMe[cannot-resolve-module] for TypeScript migration
} from '../../components/context/RefreshContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {JsonDialog, RenderLink} from './SubscriberTypes';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
}));
export type SubscriberRowType = {
  name: string,
  imsi: string,
  activeApns?: string,
  ipAddresses?: string,
  activeSessions?: number,
  service: string,
  currentUsage: string,
  dailyAvg: string,
  lastReportedTime: Date | string,
};

type SubscriberSessionRowType = {
  apnName: string,
  sessionId: string,
  ipAddr: string,
  state: string,
  activeDuration: string,
  activePolicies: Array<string>,
};
type SubscriberStateDetailPanelProps = {
  rowData: SubscriberRowType,
};
function SubscriberStateDetailPanel(props: SubscriberStateDetailPanelProps) {
  const ctx = useContext(SubscriberContext);
  const sessionState: {[string]: subscriber_state} = ctx.sessionState;
  const subscriber = sessionState[props.rowData.imsi]?.subscriber_state || {};
  const subscriberSessionRows: Array<SubscriberSessionRowType> = [];
  Object.keys(subscriber).map((apn: string) => {
    subscriber[apn].map(infos => {
      subscriberSessionRows.push({
        apnName: apn,
        sessionId: infos.session_id,
        ipAddr: infos.ipv4 ?? '-',
        state: infos.lifecycle_state,
        activeDuration: `${infos.active_duration_sec} sec`,
        activePolicies: infos.active_policy_rules,
      });
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
                    <Link>{policy.id} </Link>
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
  const params = useParams();
  const [currRow, setCurrRow] = useState<SubscriberRowType>({});
  const classes = useStyles();
  const networkId: string = nullthrows(params.networkId);
  const ctx = useContext(SubscriberContext);
  const subscriberMap: {[string]: subscriber} = ctx.state;
  const sessionState: {[string]: subscriber_state} = ctx.sessionState;
  const subscriberMetrics = ctx.metrics;
  const [jsonDialog, setJsonDialog] = useState(false);
  const subscribersIds = Object.keys(sessionState);
  const networkCtx = useContext(NetworkContext);
  const [refresh, setRefresh] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  // Auto refresh subscribers sessions every 30 seconds
  const refreshingSessionState = useRefreshingContext({
    context: SubscriberContext,
    networkId,
    type: 'subscriber',
    interval: REFRESH_INTERVAL,
    enqueueSnackbar,
    refresh,
  });

  const tableData: Array<SubscriberRowType> = subscribersIds.map(
    (imsi: string) => {
      const subscriberInfo = subscriberMap[imsi] || {};
      const metrics = subscriberMetrics?.[`${imsi}`];
      const subscriber =
        // $FlowIgnore
        refreshingSessionState.sessionState?.[imsi]?.subscriber_state || {};
      const ipAddresses = [];
      const activeApns = [];
      let activeSessions = 0;
      Object.keys(subscriber || {}).forEach(apn => {
        subscriber[apn].forEach(session => {
          if (session.lifecycle_state === 'SESSION_ACTIVE') {
            ipAddresses.push(session?.ipv4);
            activeSessions++;
          }
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

  const tableColumns = [
    {
      title: 'Name',
      field: 'name',
    },
    {
      title: 'IMSI',
      field: 'imsi',
      render: currRow => {
        const subscriberConfig = subscriberMap[currRow.imsi];
        return (
          <RenderLink
            subscriberConfig={subscriberConfig}
            currRow={currRow}
            networkCtx={networkCtx}
          />
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
