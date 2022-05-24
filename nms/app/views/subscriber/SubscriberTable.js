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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ActionTable from '../../components/ActionTable';
import Button from '@material-ui/core/Button';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
import LaunchIcon from '@material-ui/icons/Launch';
import MenuItem from '@material-ui/core/MenuItem';
// $FlowFixMe migrated to typescript
import NetworkContext from '../../components/context/NetworkContext';
import React, {useContext, useEffect, useRef, useState} from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SubscriberContext from '../../components/context/SubscriberContext';
import Text from '../../theme/design-system/Text';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import withAlert from '../../components/Alert/withAlert';
import {
  DEFAULT_PAGE_SIZE,
  REFRESH_TIMEOUT,
  SUBSCRIBER_EXPORT_COLUMNS,
  // $FlowFixMe[cannot-resolve-module] for TypeScript migration
} from './SubscriberUtils';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {ActionQuery} from '../../components/ActionTable';
import type {EnqueueSnackbarOptions} from 'notistack';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {SubscriberActionType, SubscriberInfo} from './SubscriberUtils';
import type {WithAlert} from '../../components/Alert/withAlert';
import type {
  lte_subscription,
  mutable_subscriber,
  mutable_subscribers,
  paginated_subscribers,
  subscriber,
} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe migrated to typescript
import MenuButton from '../../components/MenuButton';
import {AddSubscriberDialog} from './SubscriberAddDialog';
import {CsvBuilder} from 'filefy';
import {
  FetchSubscribers,
  handleSubscriberQuery,
} from '../../state/lte/SubscriberState';
import {JsonDialog, RenderLink} from './SubscriberOverview';
// $FlowFixMe[cannot-resolve-module]
import {base64ToHex, hexToBase64, isValidHex} from '../../util/strings';
import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';

// number of subscriber in a chunk
const SUBSCRIBERS_CHUNK_SIZE = 1000;

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
  const params = useParams();
  const networkId = nullthrows(params.networkId);
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
              case 'forbidden_network_types':
              case 'name':
                return typeof subscriberInfo[columnDef.field] === 'object'
                  ? subscriberInfo[columnDef.field].join(', ')
                  : subscriberInfo[columnDef.field];
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

function SubscriberActionsMenu(props: {onClose: () => void}) {
  const [open, setOpen] = React.useState(false);
  const [error, setError] = React.useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(SubscriberContext);
  const successCountRef = useRef(0);
  const [
    subscriberAction,
    setSubscriberAction,
  ] = useState<SubscriberActionType>('add');

  /**
   * Delete array of subscriber IMSIs.
   *
   * @param {Array<string>} subscribers Array of subscriber IMSI to delete
   */
  const deleteSubscribers = async (subscribers: Array<string>) => {
    try {
      // Delete subscribers
      subscribers.map(imsi => ctx.setState?.(imsi));
      enqueueSnackbar(`${subscribers.length} subscriber(s) deleted`, {
        variant: 'success',
      });
    } catch (e) {
      enqueueSnackbar('Deleting subscribers failed', {
        variant: 'error',
      });
    }
    props.onClose();
    setOpen(false);
  };

  const addSubscriberChunk = async (addedSubscribers: mutable_subscribers) => {
    try {
      await ctx.setState?.('', addedSubscribers);
      return true;
    } catch (e) {
      const errMsg = e.response?.data?.message ?? e.message ?? e;
      setError('Error saving subscribers: ' + errMsg);
      return false;
    }
  };

  /**
   * Add or update subscriber chunks
   *
   * @param {Array<SubscriberInfo>} subscribers Array of subscribers to Add or Update
   * @param {SubscriberActionType} subscriberAction Add or Update subscribers
   */
  const handleSubscribers = async (subscribers: Array<SubscriberInfo>) => {
    // Create array of subscriber chunk
    const subscriberChunks = subscribers.reduce((chunks, subscriber, index) => {
      const chunkIndex = Math.floor(index / SUBSCRIBERS_CHUNK_SIZE);
      if (!chunks[chunkIndex]) {
        chunks[chunkIndex] = [];
      }
      const authKey =
        subscriber.authKey && isValidHex(subscriber.authKey)
          ? hexToBase64(subscriber.authKey)
          : '';

      const authOpc =
        subscriber.authOpc !== undefined && isValidHex(subscriber.authOpc)
          ? hexToBase64(subscriber.authOpc)
          : '';
      const newSubscriber: mutable_subscriber = {
        active_apns: subscriber.apns,
        active_policies: subscriber.policies,
        forbidden_network_types: subscriber.forbiddenNetworkTypes,

        id: subscriber.imsi,
        name: subscriber.name,
        lte: {
          auth_algo: 'MILENAGE',
          auth_key: authKey,
          auth_opc: authOpc,
          state: subscriber.state,
          sub_profile: subscriber.dataPlan,
        },
      };

      chunks[chunkIndex].push(newSubscriber);
      return chunks;
    }, []);

    for (let i = 0; i < subscriberChunks.length; i++) {
      const subscriberChunk = subscriberChunks[i];
      try {
        if (subscriberAction === 'edit') {
          // Update subscribers
          subscriberChunk.map(subscriber => {
            ctx.setState?.(subscriber.id, subscriber);
          });
        } else {
          // Add subscribers

          const success = await addSubscriberChunk(subscriberChunk);
          if (success) {
            successCountRef.current =
              successCountRef.current + subscriberChunk.length;
          } else {
            enqueueSnackbar('Saving subscribers failed', {
              variant: 'error',
            });
            return;
          }
        }
      } catch (e) {
        const errMsg = e.response?.data?.message ?? e.message ?? e;
        enqueueSnackbar('Saving subscribers failed : ' + errMsg, {
          variant: 'error',
        });
      }
    }
    enqueueSnackbar('Subscriber(s) saved successfully', {
      variant: 'success',
    });
    props.onClose();
    setOpen(false);
  };

  return (
    <div>
      <AddSubscriberDialog
        error={error}
        subscriberAction={subscriberAction}
        open={open}
        onSave={(subscribers: Array<SubscriberInfo>, selectedSubscribers) => {
          if (subscriberAction === 'delete') {
            deleteSubscribers(
              selectedSubscribers?.length
                ? selectedSubscribers
                : subscribers.map(subscriber => subscriber.imsi),
            );
          } else {
            handleSubscribers(subscribers);
          }
        }}
        onClose={() => {
          setOpen(false);
          props.onClose();
        }}
      />
      <MenuButton label="Manage Subscribers">
        <MenuItem
          data-testid=""
          onClick={() => {
            setSubscriberAction('add');
            setOpen(true);
          }}>
          <Text variant="body2">Add Subscribers</Text>
        </MenuItem>
        <MenuItem>
          <Text
            variant="body2"
            onClick={() => {
              setSubscriberAction('edit');
              setOpen(true);
            }}>
            Update Subscribers
          </Text>
        </MenuItem>
        <MenuItem
          onClick={() => {
            setSubscriberAction('delete');
            setOpen(true);
          }}>
          <Text variant="body2">Delete Subscribers</Text>
        </MenuItem>
      </MenuButton>
    </div>
  );
}
function SubscribersTable(props: WithAlert) {
  const navigate = useNavigate();
  const params = useParams();
  const [currRow, setCurrRow] = useState<SubscriberRowType>({});
  const classes = useStyles();
  const networkId: string = nullthrows(params.networkId);
  const networkCtx = useContext(NetworkContext);
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(SubscriberContext);
  const subscriberMetrics = ctx.metrics;
  const [jsonDialog, setJsonDialog] = useState(false);
  const [maxPageRowCount, setMaxPageRowCount] = useState(0);
  // first token (page 1) is an empty string
  const [tokenList, setTokenList] = useState(['']);
  const onClose = () => setJsonDialog(false);
  const tableRef = React.useRef();
  const subscriberMap = ctx.state;
  const [refresh, setRefresh] = useState(false);

  const tableColumns = [
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
      title: 'Name',
      field: 'name',
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
  }, [refresh]);
  return (
    <>
      <div className={classes.dashboardRoot}>
        <CardTitleRow
          key="title"
          icon={SettingsIcon}
          label={'Subscribers'}
          filter={() => (
            <Grid
              container
              justifyContent="flex-end"
              alignItems="center"
              spacing={2}>
              <Grid item>
                <ExportSubscribersButton />
              </Grid>
              <Grid item>
                <SubscriberActionsMenu
                  onClose={() => {
                    setTimeout(() => {
                      setRefresh(!refresh);
                    }, REFRESH_TIMEOUT);
                  }}
                />
              </Grid>
            </Grid>
          )}
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
              deleteTable: false,
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
                navigate(currRow.imsi);
              },
            },
            {
              name: 'Edit',
              handleFunc: () => {
                navigate(currRow.imsi + '/config');
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
