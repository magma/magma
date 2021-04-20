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
import type {ActionQuery} from '../../components/ActionTable';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import CardTitleRow from '../../components/layout/CardTitleRow';
import NetworkContext from '../../components/context/NetworkContext';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SubscriberContext from '../../components/context/SubscriberContext';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import {DEFAULT_PAGE_SIZE} from '../../views/subscriber/SubscriberUtils';
import {JsonDialog} from './SubscriberOverview';
import {RenderLink} from './SubscriberOverview';
import {handleSubscriberQuery} from '../../state/lte/SubscriberState';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

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

function SubscribersTable(props: WithAlert) {
  const {history, match, relativeUrl} = useRouter();
  const [currRow, setCurrRow] = useState<SubscriberRowType>({});
  const classes = useStyles();
  const networkId: string = nullthrows(match.params.networkId);
  const networkCtx = useContext(NetworkContext);
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(SubscriberContext);
  const subscriberMetrics = ctx.metrics;
  const [jsonDialog, setJsonDialog] = useState(false);
  // first token (page 1) is an empty string
  const [maxPageRowCount, setMaxPageRowCount] = useState(0);
  const [tokenList, setTokenList] = useState(['']);
  const onClose = () => setJsonDialog(false);
  const tableRef = React.useRef();
  const subscriberMap = ctx.state;
  const subscriberCount = Object.keys(subscriberMap).length;

  useEffect(() => {
    tableRef.current?.onQueryChange();
  }, [subscriberCount]);

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
  ];

  return (
    <>
      <div className={classes.dashboardRoot}>
        <CardTitleRow key="title" icon={SettingsIcon} label={'Subscribers'} />
        <JsonDialog open={jsonDialog} onClose={onClose} imsi={currRow.imsi} />
        <ActionTable
          tableRef={tableRef}
          localization={{
            toolbar: {
              searchPlaceholder: 'Search IMSI001011234560000',
            },
          }}
          data={(query: ActionQuery) => {
            return handleSubscriberQuery({
              networkId,
              query,
              ctx,
              maxPageRowCount,
              setMaxPageRowCount,
              tokenList,
              setTokenList,
              pageSize: DEFAULT_PAGE_SIZE,
              subscriberMetrics,
            });
          }}
          columns={tableColumns}
          handleCurrRow={(row: SubscriberRowType) => setCurrRow(row)}
          menuItems={[
            {
              name: 'View JSON',
              handleFunc: () => {
                setJsonDialog(true);
              },
            },
            {
              name: 'View',
              handleFunc: () => {
                history.push(relativeUrl('/' + currRow.imsi));
              },
            },
            {
              name: 'Edit',
              handleFunc: () => {
                history.push(relativeUrl('/' + currRow.imsi + '/config'));
              },
            },
            {
              name: 'Remove',
              handleFunc: () => {
                props
                  .confirm(`Are you sure you want to delete ${currRow.imsi}?`)
                  .then(async confirmed => {
                    if (!confirmed) {
                      return;
                    }

                    try {
                      await ctx.setState?.(currRow.imsi);
                    } catch (e) {
                      enqueueSnackbar(
                        'failed deleting subscriber ' + currRow.imsi,
                        {
                          variant: 'error',
                        },
                      );
                    }
                  });
              },
            },
          ]}
          options={{
            actionsColumnIndex: -1,
            pageSize: DEFAULT_PAGE_SIZE,
            pageSizeOptions: [],
            showFirstLastPageButtons: false,
          }}
        />
      </div>
    </>
  );
}

const SubscriberTable = withAlert(SubscribersTable);
export default SubscriberTable;
