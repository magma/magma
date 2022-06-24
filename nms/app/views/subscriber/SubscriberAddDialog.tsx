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
import ActionTable, {SelectEditComponent} from '../../components/ActionTable';
import Alert from '@material-ui/lab/Alert';
import AlertTitle from '@material-ui/lab/AlertTitle';
import ApnContext from '../../components/context/ApnContext';
import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import Grid from '@material-ui/core/Grid';
import ListItemText from '@material-ui/core/ListItemText';
import LteNetworkContext from '../../components/context/LteNetworkContext';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import PolicyContext from '../../components/context/PolicyContext';
import React, {forwardRef, useContext, useState} from 'react';
import Select from '@material-ui/core/Select';
import SubscriberContext from '../../components/context/SubscriberContext';
import Text from '../../theme/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import nullthrows from '../../../shared/util/nullthrows';
import {Column, MaterialTableProps} from '@material-table/core';
import {PasswordInput} from '../../components/FormField';
import {SubscriberDetailsUpload} from './SubscriberUpload';
import {
  SubscriberRowType,
  handleSubscriberQuery,
} from '../../state/lte/SubscriberState';
import {colors} from '../../theme/default';
import {forbiddenNetworkTypes, validateSubscribers} from './SubscriberUtils';
import {makeStyles} from '@material-ui/styles';
import {useParams} from 'react-router-dom';
import type {ActionQuery} from '../../components/ActionTable';
import type {EditProps} from './SubscriberEditDialog';
import type {
  SubscriberActionType,
  SubscriberInfo,
  SubscribersDialogDetailProps,
} from './SubscriberUtils';

const useStyles = makeStyles(() => ({
  dialogTitle: {
    textTransform: 'capitalize',
    backgroundColor: colors.primary.brightGray,
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    color: colors.primary.white,
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
  placeholder: {
    opacity: 0.5,
  },
  dialog: {
    height: '750px',
  },
  ellipsis: {
    textOverflow: 'ellipsis',
    overflow: 'hidden',
    width: '160px',
    whiteSpace: 'nowrap',
  },
  rowId: {
    color: colors.primary.comet,
  },
  uploadDialog: {
    width: '800px',
  },
  uploadInstructions: {
    marginTop: '16px',
    color: colors.primary.comet,
  },
}));

type ActionDialogProps = {
  open: boolean;
  onClose: () => void;
  editProps?: EditProps;
  onSave: (
    subscribers: Array<SubscriberInfo>,
    selectedSubscribers?: Array<string>,
  ) => void;
  error?: string;
  subscriberAction: SubscriberActionType;
};

/**
 * Dialog used to Add/Delete/Update subscribers
 */
export function AddSubscriberDialog(props: ActionDialogProps) {
  const classes = useStyles();

  return (
    <>
      <Dialog data-testid="addSubscriberDialog" open={props.open} maxWidth="xl">
        <DialogTitle
          classes={{root: classes.dialogTitle}}
          onClose={props.onClose}
          label={`${props.subscriberAction} Subscriber(s)`}
        />

        <SubscriberDetailsDialogContent {...props} />
      </Dialog>
    </>
  );
}

/**
 * Dialog content used to Add/Delete/Update subscribers
 * Displays upload subscriber dropzone or subscriber table
 */
function SubscriberDetailsDialogContent(props: ActionDialogProps) {
  const [addError, setAddError] = useState<Array<string>>([]);
  const [subscribers, setSubscribers] = useState<Array<SubscriberInfo>>([]);
  const [upload, setUpload] = useState(false);
  const [rowAdd, setRowAdd] = useState(false);

  const subscriberProps = {
    upload,
    setUpload,
    subscribers,
    setSubscribers,
    addError,
    setAddError,
    error: props.error,
    onClose: props.onClose,
    onSave: props.onSave,
    rowAdd,
    setRowAdd,
    subscriberAction: props.subscriberAction,
  } as const;

  return (
    <>
      {!(upload || props.subscriberAction === 'edit') ? (
        <SubscriberDetailsTable {...subscriberProps} />
      ) : (
        <SubscriberDetailsUpload {...subscriberProps} />
      )}
    </>
  );
}

function SubscriberDetailsTable(props: SubscribersDialogDetailProps) {
  const {
    setSubscribers,
    setAddError,
    setUpload,
    subscribers,
    addError,
    rowAdd,
    setRowAdd,
    subscriberAction,
  } = props;
  const classes = useStyles();
  const ctx = useContext(SubscriberContext);
  const apnCtx = useContext(ApnContext);
  const lteCtx = useContext(LteNetworkContext);
  const policyCtx = useContext(PolicyContext);
  const apns = Array.from(new Set(Object.keys(apnCtx.state || {})));
  const subProfiles = Array.from(
    new Set(Object.keys(lteCtx.state.cellular?.epc?.sub_profiles || {})).add(
      'default',
    ),
  );
  const policies = Array.from(
    new Set(Object.keys(policyCtx.state || {})).add('default'),
  );
  const tableActions: MaterialTableProps<SubscriberInfo>['editable'] = {
    onRowUpdate: (newData, oldData) => {
      return new Promise<void>((resolve, reject) => {
        const err = validateSubscribers([newData], subscriberAction);
        setAddError(err);
        if (err.length > 0) {
          return reject();
        }
        const dataUpdate = [...subscribers];
        const index = ((oldData as unknown) as {tableData: {id: number}})
          .tableData.id;
        dataUpdate[index] = newData;
        setSubscribers([...dataUpdate]);
        resolve();
      });
    },
    onRowDelete: oldData =>
      new Promise<void>(resolve => {
        const dataDelete = [...subscribers];
        const index = ((oldData as unknown) as {tableData: {id: number}})
          .tableData.id;
        dataDelete.splice(index, 1);
        setSubscribers([...dataDelete]);
        resolve();
      }),
  };
  const [selectedSubscribers, setSelectedSubscribers] = useState<Array<string>>(
    [],
  );
  const params = useParams();
  const [maxPageRowCount, setMaxPageRowCount] = useState(0);
  const [tokenList, setTokenList] = useState(['']);
  const networkId: string = nullthrows(params.networkId);
  const subscriberMetrics = ctx.metrics;
  const getSubscribers = (query: ActionQuery) =>
    handleSubscriberQuery({
      networkId,
      query,
      ctx,
      maxPageRowCount,
      setMaxPageRowCount,
      tokenList,
      setTokenList,
      pageSize: 100,
      subscriberMetrics,
      deleteTable: true,
    });

  const columns: Array<Column<SubscriberInfo>> = [
    {
      title: 'IMSI',
      field: 'imsi',
      editComponent: props => (
        <OutlinedInput
          data-testid="IMSI"
          type="text"
          placeholder="Enter IMSI"
          value={props.value as string}
          onChange={e => props.onChange(e.target.value)}
        />
      ),
    },
    {
      title: 'Subscriber Name',
      field: 'name',
      editComponent: props => (
        <OutlinedInput
          data-testid="name"
          placeholder="Enter Name"
          type="text"
          value={props.value as string}
          onChange={e => {
            props.onChange(e.target.value);
          }}
        />
      ),
    },
    {
      title: 'Auth Key',
      field: 'authKey',
      editComponent: props => (
        <PasswordInput
          data-testid="authKey"
          placeholder="Key"
          value={(props.value as string) || ''}
          onChange={v => props.onChange(v)}
        />
      ),
      render: rowData => {
        return (
          <Tooltip title={rowData.authKey} placement="top">
            <div className={classes.ellipsis}>{rowData.authKey}</div>
          </Tooltip>
        );
      },
    },
    {
      title: 'Auth OPC',
      field: 'authOpc',
      editComponent: props => (
        <PasswordInput
          data-testid="authOpc"
          placeholder="OPC"
          value={props.value as string}
          onChange={v => props.onChange(v)}
        />
      ),
      render: rowData => {
        return (
          <Tooltip title={rowData.authOpc} placement="top">
            <div className={classes.ellipsis}>{rowData.authOpc}</div>
          </Tooltip>
        );
      },
    },
    {
      title: 'Service',
      field: 'state',
      editComponent: props => {
        return (
          <SelectEditComponent
            {...props}
            testId="service"
            defaultValue={'ACTIVE'}
            content={['ACTIVE', 'INACTIVE']}
            onChange={value => props.onChange(value)}
          />
        );
      },
    },
    {
      title: 'Forbidden Network Types',
      field: 'forbiddenNetworkTypes',
      editComponent: props => (
        <FormControl>
          <Select
            data-testid="forbiddenNetworkTypes"
            multiple
            value={(props.value ?? []) as Array<string>}
            onChange={({target}) => props.onChange(target.value)}
            displayEmpty={true}
            renderValue={selected => {
              if (!(selected as Array<string>).length) {
                return 'Select Forbidden Network Types';
              }
              return (selected as Array<string>).join(', ');
            }}
            input={
              <OutlinedInput
                className={props.value ? '' : classes.placeholder}
              />
            }>
            {forbiddenNetworkTypes.map((k, idx) => (
              <MenuItem key={idx} value={k}>
                <Checkbox
                  checked={
                    props.value
                      ? (props.value as Array<string>).indexOf(k) > -1
                      : false
                  }
                />
                <ListItemText primary={k} />
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      ),
    },
    {
      title: 'Data Plan',
      field: 'dataPlan',
      editComponent: props => (
        <SelectEditComponent
          {...props}
          testId="dataPlan"
          defaultValue={'default'}
          content={subProfiles}
          onChange={value => props.onChange(value)}
        />
      ),
    },
    {
      title: 'Active APNs',
      field: 'apns',
      editComponent: props => (
        <FormControl>
          <Select
            data-testid="activeApns"
            multiple
            value={(props.value ?? []) as Array<string>}
            onChange={({target}) => props.onChange(target.value)}
            displayEmpty={true}
            renderValue={selected => {
              if (!(selected as Array<string>).length) {
                return 'Select APNs';
              }
              return (selected as Array<string>).join(', ');
            }}
            input={
              <OutlinedInput
                className={props.value ? '' : classes.placeholder}
              />
            }>
            {apns.map((k: string, idx: number) => (
              <MenuItem key={idx} value={k}>
                <Checkbox
                  checked={
                    props.value
                      ? (props.value as Array<string>).indexOf(k) > -1
                      : false
                  }
                />
                <ListItemText primary={k} />
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      ),
    },
    {
      title: 'Active Policies',
      field: 'policies',
      editComponent: props => (
        <FormControl>
          <Select
            data-testid="activePolicies"
            multiple
            value={(props.value ?? []) as Array<string>}
            onChange={({target}) => props.onChange(target.value)}
            displayEmpty={true}
            renderValue={selected => {
              if (!(selected as Array<string>).length) {
                return 'Select Policies';
              }
              return (selected as Array<string>).join(', ');
            }}
            input={
              <OutlinedInput
                className={props.value ? '' : classes.placeholder}
              />
            }>
            {policies.map((k: string, idx: number) => (
              <MenuItem key={idx} value={k}>
                <Checkbox
                  checked={
                    props.value
                      ? (props.value as Array<string>).indexOf(k) > -1
                      : false
                  }
                />
                <ListItemText primary={k} />
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      ),
    },
  ];

  return (
    <>
      <DialogContent>
        {(addError.length > 0 || props.error) && (
          <Grid item>
            <Alert severity="error">
              <AlertTitle>Error Adding Subscriber(s)</AlertTitle>
              {addError.length > 0 ? (
                <ul>
                  {addError.map(e => (
                    <li>{e}</li>
                  ))}
                </ul>
              ) : (
                <> {props.error} </>
              )}
            </Alert>
          </Grid>
        )}
        <Grid>
          <Text>
            {subscriberAction === 'delete' && selectedSubscribers.length
              ? `Select Subscribers (${selectedSubscribers.length} Selected)`
              : `${subscriberAction === 'delete' ? 'Deleting' : 'Adding'} ${
                  subscribers.length
                } subscriber(s)`}
          </Text>
        </Grid>

        <ActionTable
          data={
            // The table is rendered with two different data types and the columns and action ar changed depending on the type.
            // To not complete break the typing we pretend that the data is always the same.
            subscriberAction === 'delete' && !subscribers.length
              ? ((getSubscribers as unknown) as Array<SubscriberInfo>)
              : subscribers
          }
          columns={
            !(subscriberAction === 'delete')
              ? [
                  {
                    title: '',
                    field: '',
                    width: '70px',
                    editable: 'never',
                    align: 'center',
                    render: rowData => (
                      <Text variant="subtitle3">
                        {((rowData as unknown) as {tableData: {id: number}})
                          .tableData?.id + 1 || ''}
                      </Text>
                    ),
                  },
                  ...columns,
                ]
              : columns
          }
          onSelectionChange={(rows: Array<SubscriberInfo>) => {
            const newSubscribers = rows.map(r => r.imsi);
            setSelectedSubscribers(oldSubscribers => {
              return [...new Set([...newSubscribers, ...oldSubscribers])];
            });
          }}
          options={{
            actionsColumnIndex: -1,
            pageSize: 100,
            pageSizeOptions: [100],
            tableLayout: 'fixed',
            fixedColumns: {
              left: 1,
            },
            showTextRowsSelected: false,
            selection: subscriberAction === 'delete' && !subscribers.length,
            selectionProps: (rowData: SubscriberInfo | SubscriberRowType) => {
              return {
                checked: selectedSubscribers.includes(rowData.imsi),
                value: rowData.imsi,
                onClick: (event: React.ChangeEvent<HTMLInputElement>) => {
                  if (selectedSubscribers.includes(event.target.value)) {
                    const newSubscribers = selectedSubscribers.filter(
                      imsi => imsi !== event.target.value,
                    );
                    setSelectedSubscribers([...newSubscribers]);
                  }
                },
              };
            },
          }}
          editable={
            // Hide 'Upload CSV' and 'Add New Row' button if subscribers are uploaded
            // or if subscribers are added one by one

            subscriberAction === 'delete'
              ? {}
              : subscribers.length > 0 && !rowAdd
              ? tableActions
              : {
                  ...tableActions,
                  onRowAdd: newData => {
                    setRowAdd(true);
                    return new Promise<void>((resolve, reject) => {
                      const err = validateSubscribers(
                        [newData],
                        subscriberAction,
                      );
                      setAddError(err);
                      if (err.length > 0) {
                        return reject();
                      }
                      setSubscribers([...subscribers, newData]);
                      resolve();
                    });
                  },
                }
          }
          actions={
            subscribers.length > 0 && !rowAdd
              ? []
              : [
                  {
                    icon: (forwardRef<SVGSVGElement>((props, ref) => (
                      <Button
                        startIcon={<CloudUploadIcon {...props} ref={ref} />}
                        variant="outlined"
                        color="primary">
                        {subscriberAction === 'delete'
                          ? 'Delete from CSV'
                          : 'Upload CSV'}
                      </Button>
                    )) as unknown) as () => React.ReactElement<any>,
                    tooltip: 'Upload',
                    isFreeAction: true,
                    onClick: () => {
                      setUpload(true);
                    },
                  },
                ]
          }
        />
      </DialogContent>
      <DialogActions>
        <Grid container justifyContent="space-between">
          <Grid item>
            <Button
              disabled={!(subscribers.length > 0) || rowAdd}
              onClick={() => {
                setUpload(true);
              }}>
              Back
            </Button>
          </Grid>
          <Grid item>
            <Button onClick={props.onClose}> Cancel </Button>
            <Button
              variant="contained"
              color="primary"
              data-testid="saveSubscriber"
              onClick={() => {
                const err = validateSubscribers(subscribers, subscriberAction);
                setAddError(err);
                if (!err.length) {
                  props.onSave?.(subscribers, selectedSubscribers);
                  setSelectedSubscribers([]);
                }
              }}>
              {subscriberAction === 'delete'
                ? 'Delete Subcribers'
                : 'Save and Add Subscribers'}
            </Button>
          </Grid>
        </Grid>
      </DialogActions>
    </>
  );
}
