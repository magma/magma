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
import type {EnqueueSnackbarOptions} from 'notistack';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {
  lte_subscription,
  paginated_subscribers,
  subscriber,
} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
import LaunchIcon from '@material-ui/icons/Launch';
import NetworkContext from '../../components/context/NetworkContext';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SubscriberContext from '../../components/context/SubscriberContext';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import {CsvBuilder} from 'filefy';
import {
  DEFAULT_PAGE_SIZE,
  SUBSCRIBER_EXPORT_COLUMNS,
} from '../../views/subscriber/SubscriberUtils';
import {
  FetchSubscribers,
  handleSubscriberQuery,
} from '../../state/lte/SubscriberState';
import {JsonDialog} from './SubscriberOverview';
import {RenderLink} from './SubscriberOverview';
import {base64ToHex} from '@fbcnms/util/strings';
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

function ExportSubscribersButton() {
  const {match} = useRouter();
  const networkId = nullthrows(match.params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();

  return (
    <Grid item>
      <Button
        variant="contained"
        color="primary"
        startIcon={<LaunchIcon />}
        onClick={() =>
          exportSubscribers({
            networkId,
            enqueueSnackbar,
          })
        }>
        Export
      </Button>
    </Grid>
  );
}
type ExportProps = {
  networkId: string,
  enqueueSnackbar: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
};
/**
 * Export subscribers in csv format.
 * Iterates over paginated subscribers.
 *
 * @param {string} networkId ID of the network.
 * @param {ActionQuery} enqueueSnackbar Snackbar to display error message
 */
async function exportSubscribers(props: ExportProps) {
  const {networkId, enqueueSnackbar} = props;
  let page = 1;
  let token = undefined;
  const currTs = Date.now();
  const subscriberExport = new CsvBuilder(`subscribers_${currTs}.csv`)
    .setDelimeter(',')
    .setColumns(SUBSCRIBER_EXPORT_COLUMNS.map(columnDef => columnDef.title));
  try {
    // last page next_page_token is an empty string
    while (token !== '') {
      // $FlowIgnore
      const subscriberRows: paginated_subscribers = await FetchSubscribers({
        networkId,
        token,
      });
      if (subscriberRows) {
        page = page + 1;
        token = subscriberRows.next_page_token;
      }
      const subscriberData = Object.keys(subscriberRows.subscribers).map(
        rowData =>
          SUBSCRIBER_EXPORT_COLUMNS.map(columnDef => {
            const subscriberConfig: lte_subscription =
              subscriberRows.subscribers[rowData].config.lte;
            const subscriberInfo: subscriber =
              subscriberRows.subscribers[rowData];
            switch (columnDef.field) {
              case 'auth_opc':
              case 'auth_key':
                return base64ToHex(subscriberConfig[columnDef.field] ?? '');
              case 'state':
              case 'sub_profile':
                return subscriberConfig[columnDef.field];
              case 'id':
              case 'active_apns':
              case 'name':
                return typeof subscriberInfo[columnDef.field] === 'object'
                  ? subscriberInfo[columnDef.field].join('|')
                  : subscriberInfo[columnDef.field];
              default:
                console.log('invalid field not found', columnDef.field);
            }
          }),
      );
      subscriberExport.addRows(subscriberData);
    }
    if (subscriberExport) {
      subscriberExport.exportFile();
    }
  } catch (e) {
    enqueueSnackbar(e?.message ?? 'error retrieving subscribers', {
      variant: 'error',
    });
  }
}

function SubscribersTable(props: WithAlert & {refresh: boolean}) {
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
  // refresh data on subscriber add
  useEffect(() => {
    tableRef.current?.onQueryChange();
  }, [props.refresh]);
  return (
    <>
      <div className={classes.dashboardRoot}>
        <CardTitleRow
          key="title"
          icon={SettingsIcon}
          label={'Subscribers'}
          filter={() => <ExportSubscribersButton />}
        />
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
                      // refresh table data
                      tableRef.current?.onQueryChange();
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
