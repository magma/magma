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
import ActionTable, {TableRef} from '../../components/ActionTable';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import EmptyState from '../../components/EmptyState';
import Grid from '@material-ui/core/Grid';
import LaunchIcon from '@material-ui/icons/Launch';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React, {
  SyntheticEvent,
  useContext,
  useEffect,
  useRef,
  useState,
} from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SubscriberContext from '../../context/SubscriberContext';
import Text from '../../theme/design-system/Text';
import nullthrows from '../../../shared/util/nullthrows';
import withAlert from '../../components/Alert/withAlert';

import {AddSubscriberDialog} from './SubscriberAddDialog';
import {CsvBuilder} from 'filefy';
import {
  DEFAULT_PAGE_SIZE,
  JsonDialog,
  REFRESH_TIMEOUT,
  RenderLink,
  SUBSCRIBER_EXPORT_COLUMNS,
} from './SubscriberUtils';
import type {ActionQuery} from '../../components/ActionTable';
import type {
  LteSubscription,
  MutableSubscriber,
  PaginatedSubscribers,
  Subscriber,
} from '../../../generated';
import type {OptionsObject} from 'notistack';
import type {SubscriberActionType, SubscriberInfo} from './SubscriberUtils';
import type {WithAlert} from '../../components/Alert/withAlert';

import {Column} from '@material-table/core';
import {MenuProps} from '@material-ui/core/Menu/Menu';
import {Theme} from '@material-ui/core/styles';
import {base64ToHex, hexToBase64, isValidHex} from '../../util/strings';
import {
  fetchSubscribers,
  handleSubscriberQuery,
} from '../../util/SubscriberState';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

// number of subscriber in a chunk
const SUBSCRIBERS_CHUNK_SIZE = 1000;
const EMPTY_STATE_OVERVIEW =
  'The subscriber page allows you to add, edit, and delete your subscribers. You’ll be able to view current data ' +
  'usage, average data usage, last reported time (displayed if subscriber monitoring is enabled), and other status information from the subscriber table.';
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
          void exportSubscribers({
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
  networkId: string;
  enqueueSnackbar: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
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
      const subscriberRows = (await fetchSubscribers({
        networkId,
        token,
      })) as PaginatedSubscribers;
      if (subscriberRows) {
        page = page + 1;
        token = subscriberRows.next_page_token;
      }
      const subscriberData = Object.keys(subscriberRows.subscribers).map(
        rowData =>
          SUBSCRIBER_EXPORT_COLUMNS.map(columnDef => {
            const subscriberConfig: LteSubscription =
              subscriberRows.subscribers[rowData].config.lte;
            const subscriberInfo: Subscriber =
              subscriberRows.subscribers[rowData];
            switch (columnDef.field) {
              case 'auth_opc':
              case 'auth_key':
                return base64ToHex(subscriberConfig[columnDef.field] ?? '');
              case 'state':
              case 'sub_profile':
                return subscriberConfig[columnDef.field];
              case 'forbidden_network_types':
              case 'name': {
                const field = subscriberInfo[columnDef.field]!;
                return typeof field === 'object' ? field.join(', ') : field;
              }
              case 'id':
              case 'active_apns': {
                const field = subscriberInfo[columnDef.field]!;
                return typeof field === 'object' ? field.join('|') : field;
              }
              default:
                // @ts-ignore
                console.error('invalid field not found', columnDef.field);
                return '';
            }
          }),
      );
      subscriberExport.addRows(subscriberData);
    }
    if (subscriberExport) {
      subscriberExport.exportFile();
    }
  } catch (e) {
    enqueueSnackbar(getErrorMessage(e, 'error retrieving subscribers'), {
      variant: 'error',
    });
  }
}

const StyledMenu = withStyles({
  paper: {
    border: '1px solid #d3d4d5',
  },
})((props: MenuProps) => (
  <Menu
    elevation={0}
    getContentAnchorEl={null}
    anchorOrigin={{
      vertical: 'bottom',
      horizontal: 'center',
    }}
    transformOrigin={{
      vertical: 'top',
      horizontal: 'center',
    }}
    {...props}
  />
));

function SubscriberActionsMenu(props: {
  onClose: () => void;
  // used for empty state to only show add subscriber dialog
  addDialog: boolean;
  // hide manage button
  hideButton: boolean;
}) {
  const [anchorEl, setAnchorEl] = React.useState<Element | null>(null);
  const [open, setOpen] = React.useState(false);
  const [error, setError] = React.useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(SubscriberContext);
  const successCountRef = useRef(0);
  const [subscriberAction, setSubscriberAction] = useState<
    SubscriberActionType
  >('add');
  const handleClick = (event: SyntheticEvent<Element>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };
  useEffect(() => {
    if (props.addDialog) {
      setOpen(true);
      setSubscriberAction('add');
    }
  }, [props.addDialog]);

  /**
   * Delete array of subscriber IMSIs.
   *
   * @param {Array<string>} subscribers Array of subscriber IMSI to delete
   */
  const deleteSubscribers = (subscribers: Array<string>) => {
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

  const addSubscriberChunk = async (
    addedSubscribers: Array<MutableSubscriber>,
  ) => {
    try {
      await ctx.setState?.('', addedSubscribers);
      return true;
    } catch (e) {
      setError('Error saving subscribers: ' + getErrorMessage(e));
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
      const newSubscriber: MutableSubscriber = {
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
    }, [] as Array<Array<MutableSubscriber>>);

    for (let i = 0; i < subscriberChunks.length; i++) {
      const subscriberChunk = subscriberChunks[i];
      try {
        if (subscriberAction === 'edit') {
          // Update subscribers
          subscriberChunk.forEach(subscriber => {
            void ctx.setState?.(subscriber.id, subscriber);
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
        const errMsg = getErrorMessage(e);
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
            void handleSubscribers(subscribers);
          }
        }}
        onClose={() => {
          setOpen(false);
          props.onClose();
        }}
      />
      {!props.hideButton && (
        <>
          <Button
            variant="contained"
            color="primary"
            onClick={handleClick}
            endIcon={<ArrowDropDownIcon />}>
            {'Manage Subscribers'}
          </Button>
          <StyledMenu
            anchorEl={anchorEl}
            keepMounted
            open={Boolean(anchorEl)}
            onClose={handleClose}>
            <MenuItem
              data-testid=""
              onClick={() => {
                setSubscriberAction('add');
                setOpen(true);
              }}>
              <Text variant="subtitle2">Add Subscribers</Text>
            </MenuItem>
            <MenuItem>
              <Text
                variant="subtitle2"
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
              <Text variant="subtitle2">Delete Subscribers</Text>
            </MenuItem>
          </StyledMenu>
        </>
      )}
    </div>
  );
}
function SubscribersTable(props: WithAlert) {
  const navigate = useNavigate();
  const params = useParams();
  const [currRow, setCurrRow] = useState<SubscriberRowType>(
    {} as SubscriberRowType,
  );
  const classes = useStyles();
  const networkId: string = nullthrows(params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(SubscriberContext);
  const subscriberMetrics = ctx.metrics;
  const [jsonDialog, setJsonDialog] = useState(false);
  const [maxPageRowCount, setMaxPageRowCount] = useState(0);
  // first token (page 1) is an empty string
  const [tokenList, setTokenList] = useState(['']);
  const onClose = () => setJsonDialog(false);
  const tableRef: TableRef = React.useRef();
  const subscriberMap = ctx.state;
  const [refresh, setRefresh] = useState(false);
  const [addDialog, setAddDialog] = useState(false);

  const tableColumns: Array<Column<SubscriberRowType>> = [
    {
      title: 'IMSI',
      field: 'imsi',
      render: currRow => {
        const subscriberConfig = subscriberMap[currRow.imsi];
        return (
          <RenderLink subscriberConfig={subscriberConfig} currRow={currRow} />
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

  const cardActions = {
    buttonText: 'Add Subscribers',
    onClick: () => setAddDialog(true),
    linkText: 'Learn more about Subscribers',
    link:
      'https://docs.magmacore.org/docs/next/nms/subscriber#subscriber-dashboard',
  };
  return (
    <>
      <div className={classes.dashboardRoot}>
        {Object.keys(subscriberMap).length > 0 ? (
          <>
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
                      addDialog={false}
                      hideButton={false}
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
            <JsonDialog
              open={jsonDialog}
              onClose={onClose}
              imsi={currRow.imsi}
            />
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
                    void props
                      .confirm(
                        `Are you sure you want to delete ${currRow.imsi}?`,
                      )
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
                idSynonym: 'imsi',
                sorting: false,
                actionsColumnIndex: -1,
                pageSize: DEFAULT_PAGE_SIZE,
                pageSizeOptions: [],
                showFirstLastPageButtons: false,
              }}
            />
          </>
        ) : (
          <Grid container justifyContent="space-between" spacing={3}>
            <SubscriberActionsMenu
              addDialog={addDialog}
              hideButton={true}
              onClose={() => {
                setTimeout(() => {
                  setAddDialog(false);
                  setRefresh(!refresh);
                }, REFRESH_TIMEOUT);
              }}
            />
            <EmptyState
              title={'Set up Subscribers'}
              instructions={
                'Add subscriber by manually entering subscriber information, or uploading a CSV file.'
              }
              cardActions={cardActions}
              overviewTitle={'Subscribers Overview'}
              overviewDescription={EMPTY_STATE_OVERVIEW}
            />
          </Grid>
        )}
      </div>
    </>
  );
}

const SubscriberTable = withAlert(SubscribersTable);
export default SubscriberTable;
