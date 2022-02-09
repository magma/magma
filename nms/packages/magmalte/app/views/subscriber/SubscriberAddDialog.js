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
import type {SubscriberActionType} from './SubscriberUtils';
import type {
  core_network_types,
  subscriber,
} from '../../../generated/MagmaAPIBindings';

import ActionTable from '../../components/ActionTable';
import Alert from '@material-ui/lab/Alert';
import AlertTitle from '@material-ui/lab/AlertTitle';
import ApnContext from '../../components/context/ApnContext';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import Checkbox from '@material-ui/core/Checkbox';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import EditSubscriberApnStaticIps from './SubscriberApnStaticIpsEdit';
import EditSubscriberTrafficPolicy from './SubscriberTrafficPolicyEdit';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import List from '@material-ui/core/List';
import ListItemText from '@material-ui/core/ListItemText';
import LteNetworkContext from '../../components/context/LteNetworkContext';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import PolicyContext from '../../components/context/PolicyContext';
import React from 'react';
import Select from '@material-ui/core/Select';
import SubscriberContext from '../../components/context/SubscriberContext';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../../theme/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import TypedSelect from '@fbcnms/ui/components/TypedSelect';
import nullthrows from '@fbcnms/util/nullthrows';

import {AltFormField, PasswordInput} from '../../components/FormField';
import {CoreNetworkTypes, SUBSCRIBER_ADD_ERRORS} from './SubscriberUtils';
import {DropzoneArea} from 'material-ui-dropzone';
import {SelectEditComponent} from '../../components/ActionTable';
import {base64ToHex, hexToBase64, isValidHex} from '@fbcnms/util/strings';
import {colors, typography} from '../../theme/default';
import {forwardRef} from 'react';
import {handleSubscriberQuery} from '../../state/lte/SubscriberState';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  alert: {
    backgroundColor: colors.primary.white,
    padding: '20px 40px 20px 40px',
  },
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    color: colors.primary.white,
  },
  tabs: {
    color: colors.primary.white,
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '16px 0 16px 0',
    display: 'flex',
    alignItems: 'center',
  },
  tabIconLabel: {
    marginRight: '8px',
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
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

const MAX_UPLOAD_FILE_SZ_BYTES = 10 * 1024 * 1024;
const SUBSCRIBER_TITLE = 'Subscriber';
const TRAFFIC_TITLE = 'Traffic Policy';
const STATIC_IPS_TITLE = 'APNs Static IPs';
const UPLOAD_DOC_LINK =
  'https://docs.magmacore.org/docs/nms/subscriber#uploading-a-subscriber-csv-file';
const ADD_INSTRUCTIONS =
  'You can download this template that automatically maps the fields. Find more instruction in ';
const DELETE_INSTRUCTIONS =
  'You can export all subscribers and select the subscribers you want to delete. Find more instruction in ';
const EDIT_INSTRUCTIONS =
  'You can export all subscribers to edit and upload the file. Find more instruction in ';

export type SubscriberInfo = {
  name: string,
  imsi: string,
  authKey: string,
  authOpc: string,
  state: 'INACTIVE' | 'ACTIVE',
  forbiddenNetworkTypes: core_network_types,
  dataPlan: string,
  apns: Array<string>,
  policies?: Array<string>,
};

const SUB_NAME_OFFSET = 0;
const SUB_IMSI_OFFSET = 1;
const SUB_AUTH_KEY_OFFSET = 2;
const SUB_AUTH_OPC_OFFSET = 3;
const SUB_STATE_OFFSET = 4;
const SUB_FORBIDDEN_NETWORK_TYPE_OFFSET = 5;
const SUB_DATAPLAN_OFFSET = 6;
const SUB_APN_OFFSET = 7;
const SUB_POLICY_OFFSET = 8;
const SUB_MAX_FIELDS = 9;
const forbiddenNetworkTypes = Object.keys(CoreNetworkTypes).map(
  key => CoreNetworkTypes[key],
);

function parseSubscriber(line: string) {
  const items = line.split(',').map(item => item.trim());
  if (items.length > SUB_MAX_FIELDS) {
    throw new Error(
      `Too many fields to parse, expected ${SUB_MAX_FIELDS} fields, received ${items.length} fields`,
    );
  }
  return {
    name: items[SUB_NAME_OFFSET],
    imsi: items[SUB_IMSI_OFFSET],
    authKey: items[SUB_AUTH_KEY_OFFSET],
    authOpc: items[SUB_AUTH_OPC_OFFSET],
    state: items[SUB_STATE_OFFSET] === 'ACTIVE' ? 'ACTIVE' : 'INACTIVE',
    forbiddenNetworkTypes: forbiddenNetworkTypes.filter(value =>
      items[SUB_FORBIDDEN_NETWORK_TYPE_OFFSET]?.split('|')
        .map(item => item.trim())
        .filter(Boolean)
        .includes(value),
    ),
    dataPlan: items[SUB_DATAPLAN_OFFSET],
    apns: items[SUB_APN_OFFSET]?.split('|')
      .map(item => item.trim())
      .filter(Boolean),
    policies: items?.[SUB_POLICY_OFFSET]?.split('|')
      .map(item => item.trim())
      .filter(Boolean),
  };
}

function parseSubscriberFile(fileObj: File) {
  const reader = new FileReader();
  const subscribers = [];
  return new Promise((resolve, reject) => {
    if (fileObj.size > MAX_UPLOAD_FILE_SZ_BYTES) {
      reject(
        'file size exceeds max upload size of 10MB, please upload smaller file',
      );
      return;
    }

    reader.onload = async e => {
      try {
        if (!(e.target instanceof FileReader)) {
          reject('invalid target type');
          return;
        }

        const text = e.target.result;
        if (typeof text !== 'string') {
          reject('invalid file content');
          return;
        }

        for (const line of text
          .split('\n')
          .map(item => item.trim())
          .filter(Boolean)) {
          subscribers.push(parseSubscriber(line));
        }
      } catch (e) {
        reject('Failed parsing the file ' + fileObj.name + '. ' + e?.message);
        return;
      }
      resolve(subscribers);
    };
    reader.readAsText(fileObj);
  });
}

export function EditSubscriberButton(props: EditProps) {
  const [open, setOpen] = useState(false);
  return (
    <>
      <SubscriberEditDialog
        editProps={props}
        open={open}
        onClose={() => setOpen(false)}
        subscriberAction={''}
      />
      <Button
        component="button"
        data-testid={props.editTable}
        variant="text"
        onClick={() => setOpen(true)}>
        {'Edit'}
      </Button>
    </>
  );
}

const EditTableType = {
  subscriber: 0,
  trafficPolicy: 1,
  staticIps: 2,
};

type EditProps = {
  editTable: $Keys<typeof EditTableType>,
};

type ActionDialogProps = {
  open: boolean,
  onClose: () => void,
  editProps?: EditProps,
  onSave?: (
    subscribers: Array<SubscriberInfo>,
    selectedSubscribers?: Array<string>,
  ) => void,
  error?: string,
  subscriberAction: SubscriberActionType,
};

type DialogProps = {
  open: boolean,
  onClose: () => void,
  editProps?: EditProps,
  onSave?: (
    subscribers: Array<SubscriberInfo>,
    selectedSubscribers?: Array<string>,
  ) => void,
  error?: string,
};

/**
 * Dialog used to Add/Delete/Update subscribers
 */
export function AddSubscriberDialog(props: ActionDialogProps) {
  return (
    <>
      <Dialog
        data-testid="addSubscriberDialog"
        open={props.open}
        onSave={(subscribers, selectedSubscribers) => {
          props.onSave?.(subscribers || [], selectedSubscribers);
        }}
        maxWidth="xl">
        <DialogTitle
          onClose={props.onClose}
          label={`${props.subscriberAction} Subscriber(s)`}
        />

        <SubscriberDetailsDialogContent {...props} />
      </Dialog>
    </>
  );
}

export type EditSubscriberProps = {
  subscriberState: subscriber,
  onSubscriberChange: (key: string, val: string | number | {}) => void,
  inputClass: string,
  onTrafficPolicyChange: (
    key: string,
    val: string | number | {},
    index: number,
  ) => void,
  onDeleteApn: (apn: {}) => void,
  onAddApnStaticIP: () => void,
  subProfiles: {},
  subscriberStaticIPRows: Array<subscriberStaticIpsRowType>,
  forbiddenNetworkTypes: Array<subscriberForbiddenNetworkTypes>,
  authKey: string,
  authOpc: string,
  setAuthKey: (key: string) => void,
  setAuthOpc: (key: string) => void,
};

export function SubscriberEditDialog(props: DialogProps) {
  const {editProps} = props;
  const enqueueSnackbar = useEnqueueSnackbar();
  const [tabPos, setTabPos] = useState(
    editProps ? EditTableType[editProps.editTable] : 0,
  );
  const ctx = useContext(SubscriberContext);
  const lteCtx = useContext(LteNetworkContext);
  const classes = useStyles();
  const {match} = useRouter();
  const subscriberId = nullthrows(match.params.subscriberId);
  const [subscriberState, setSubscriberState] = useState<subscriber>(
    ctx.state[subscriberId],
  );

  const [authKey, setAuthKey] = useState(
    subscriberState.lte.auth_key
      ? base64ToHex(subscriberState.lte.auth_key)
      : '',
  );
  const [authOpc, setAuthOpc] = useState(
    subscriberState.lte.auth_opc != null
      ? base64ToHex(subscriberState.lte.auth_opc)
      : '',
  );
  const [subscriberStaticIPRows, setSubscriberStaticIPRows] = useState<
    Array<subscriberStaticIpsRowType>,
  >(
    Object.keys(ctx.state[subscriberId].config.static_ips || {}).map(
      (apn: string) => {
        return {
          apnName: apn,
          staticIp: ctx.state[subscriberId].config.static_ips?.[apn] || '',
        };
      },
    ),
  );

  const subscriberCoreNetwork = Array<subscriberForbiddenNetworkTypes>(
    Object.keys(CoreNetworkTypes).map((key: string) => {
      return {
        nwTypes: key,
      };
    }),
  );

  const [error, setError] = useState('');
  useEffect(() => {
    setTabPos(props.editProps ? EditTableType[props.editProps.editTable] : 0);
  }, [props.editProps]);

  const onClose = () => {
    props.onClose();
  };

  // we are doing this to ensure we can map subprofiles from an array
  // for e.g. ['foo', 'default'] -> {foo: 'foo', default: 'default}
  // this is done because TypedSelect expects items in this form to verify
  // if the passed in input is of expected type
  const subProfiles = Array.from(
    new Set(Object.keys(lteCtx.state.cellular?.epc?.sub_profiles || {})).add(
      'default',
    ),
  ).reduce(function (o, v) {
    o[v] = v;
    return o;
  }, {});

  const subscriberProps: EditSubscriberProps = {
    subscriberState: subscriberState,
    onSubscriberChange: (key: string, val) => {
      setSubscriberState({...subscriberState, [key]: val});
    },
    onTrafficPolicyChange: (key: string, val, index: number) => {
      const rows = subscriberStaticIPRows;
      rows[index][key] = val;
      setSubscriberStaticIPRows([...rows]);
    },
    onDeleteApn: (apn: {}) => {
      setSubscriberStaticIPRows([
        ...subscriberStaticIPRows.filter(
          (deletedApn: subscriberStaticIpsRowType) => apn !== deletedApn,
        ),
      ]);
    },
    onAddApnStaticIP: () => {
      setSubscriberStaticIPRows([
        ...subscriberStaticIPRows,
        {apnName: '', staticIp: ''},
      ]);
    },
    subProfiles: subProfiles,
    subscriberStaticIPRows: subscriberStaticIPRows,
    forbiddenNetworkTypes: subscriberCoreNetwork,
    authKey: authKey,
    authOpc: authOpc,
    setAuthKey: (key: string) => setAuthKey(key),
    setAuthOpc: (key: string) => setAuthOpc(key),
    inputClass: classes.input,
  };

  const onSave = async () => {
    try {
      if (authOpc !== '') {
        if (isValidHex(authOpc)) {
          subscriberState.lte.auth_opc = hexToBase64(authOpc);
        } else {
          setError('auth_opc is not a valid hex');
          return;
        }
      }

      if (authKey !== '') {
        if (isValidHex(authKey)) {
          subscriberState.lte.auth_key = hexToBase64(authKey);
        } else {
          setError('auth_key is not a valid hex');
          return;
        }
      }
      const {config: _, ...mutableSubscriber} = {...subscriberState};
      const staticIps = {};
      subscriberStaticIPRows.forEach(
        apn => (staticIps[apn.apnName] = apn.staticIp),
      );
      await ctx.setState?.(subscriberState.id, {
        ...mutableSubscriber,
        static_ips: staticIps,
      });
      enqueueSnackbar('Subscriber saved successfully', {
        variant: 'success',
      });
    } catch (e) {
      const errMsg = e.response?.data?.message ?? e.message;
      setError('error saving ' + subscriberState.id + ' : ' + errMsg);
      return;
    }
    props.onClose();
  };

  return (
    <Dialog
      classes={{paper: classes.dialog}}
      data-testid="editDialog"
      open={props.open}
      fullWidth={true}
      maxWidth="sm">
      <DialogTitle label={'Edit Subscriber Settings'} onClose={onClose} />
      <Tabs
        value={tabPos}
        onChange={(_, v) => setTabPos(v)}
        indicatorColor="primary"
        className={classes.tabBar}>
        <Tab
          key="subscriber"
          data-testid="subscriberTab"
          label={SUBSCRIBER_TITLE}
        />
        ;
        <Tab
          key="trafficPolicy"
          data-testid="trafficPolicyTab"
          label={TRAFFIC_TITLE}
        />
        <Tab
          key="apnStaticIps"
          data-testid="staticIpsTab"
          label={STATIC_IPS_TITLE}
        />
        ;
      </Tabs>
      <DialogContent>
        <List>
          {error !== '' && (
            <AltFormField disableGutters label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}

          {tabPos === 0 && (
            <div>
              <EditSubscriberDetails {...subscriberProps} />
            </div>
          )}
          {tabPos === 1 && <EditSubscriberTrafficPolicy {...subscriberProps} />}
          {tabPos === 2 && (
            <div>
              <EditSubscriberApnStaticIps {...subscriberProps} />
            </div>
          )}
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          {'Close'}
        </Button>
        <Button
          data-testid={`${props.editProps?.editTable || ''}-saveButton`}
          variant="contained"
          color="primary"
          onClick={onSave}>
          {'Save'}
        </Button>
      </DialogActions>
    </Dialog>
  );
}

type AddSubscribersProps = {
  // Subscribers to add, edit or delete
  setSubscribers: (Array<SubscriberInfo>) => void,
  subscribers: Array<SubscriberInfo>,
  // Formatting error (eg: field missing, wrong IMSI format)
  setAddError: (Array<string>) => void,
  addError: Array<string>,
  // Display dropzone if set to true
  setUpload: boolean => void,
  upload: boolean,
  onClose: () => void,
  // Add, edit or delete subscribers
  onSave?: (Array<SubscriberInfo>, selectedSubscribers?: Array<string>) => void,
  error?: string,
  // Row added with the Add New Row button
  rowAdd: boolean,
  setRowAdd: boolean => void,
  // Delete subscribers dialog
  subscriberAction: SubscriberActionType,
};

function SubscriberDetailsUpload(props: AddSubscribersProps) {
  const {
    setSubscribers,
    setAddError,
    setUpload,
    upload,
    subscribers,
    subscriberAction,
  } = props;
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [fileName, setFileName] = useState('');
  const DropzoneText = () => (
    <div>
      Drag and drop or <Link>browse files</Link>
    </div>
  );

  return (
    <>
      <DialogContent classes={{root: classes.uploadDialog}}>
        <CardTitleRow label={'Upload CSV'} />
        <Grid container>
          <Grid item xs={12}>
            <Alert severity="warning">
              This will replace the subscribers you entered on the previous
              page.
            </Alert>
          </Grid>
          {!fileName ? (
            <Grid item xs={12}>
              <DropzoneArea
                dropzoneText={<DropzoneText />}
                useChipsForPreview
                showPreviewsInDropzone={false}
                filesLimit={1}
                showAlerts={false}
                onChange={async files => {
                  if (files.length) {
                    try {
                      const newSubscribers = await parseSubscriberFile(
                        files[0],
                      );
                      if (newSubscribers) {
                        setSubscribers([...newSubscribers]);
                        const errors = validateSubscribers(
                          newSubscribers,
                          subscriberAction === 'delete',
                        );
                        setFileName(files[0].name);
                        if (!(subscriberAction === 'delete')) {
                          setUpload(false);
                          setAddError(errors);
                        }
                      }
                    } catch (e) {
                      enqueueSnackbar(e, {
                        variant: 'error',
                      });
                    }
                  }
                }}
              />
              <Text variant="body2" className={classes.uploadInstructions}>
                {`Accepted file type: .csv (<10 MB).  ${
                  subscriberAction === 'delete'
                    ? DELETE_INSTRUCTIONS
                    : subscriberAction === 'add'
                    ? ADD_INSTRUCTIONS
                    : EDIT_INSTRUCTIONS
                }`}
                <Link href={UPLOAD_DOC_LINK}>documentation</Link>
              </Text>
            </Grid>
          ) : (
            <Grid item xs={12}>
              <Alert severity="success">{`${fileName} is uploaded`}</Alert>
            </Grid>
          )}
        </Grid>
      </DialogContent>
      <DialogActions>
        <Grid container justify="space-between">
          <Grid item>
            {upload && (
              <Button
                onClick={() => {
                  setUpload(false);
                  if (subscriberAction === 'delete' && subscribers.length > 0) {
                    setSubscribers([]);
                  }
                }}>
                Back
              </Button>
            )}
          </Grid>
          <Grid item>
            <Button onClick={props.onClose}> Cancel </Button>
            <Button
              data-testid="saveSubscriber"
              onClick={() => {
                props.onSave?.(subscribers);
              }}>
              {subscriberAction === 'delete' ?? false
                ? 'Delete Subcribers'
                : 'Save and Add Subscribers'}
            </Button>
          </Grid>
        </Grid>
      </DialogActions>
    </>
  );
}

function SubscriberDetailsTable(props: AddSubscribersProps) {
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
  const tableActions = {
    onRowUpdate: (newData, oldData) => {
      return new Promise((resolve, reject) => {
        const err = validateSubscribers([newData]);
        setAddError(err);
        if (err.length > 0) {
          return reject();
        }
        const dataUpdate = [...subscribers];
        const index = oldData.tableData.id;
        dataUpdate[index] = newData;
        setSubscribers([...dataUpdate]);
        resolve();
      });
    },
    onRowDelete: oldData =>
      new Promise(resolve => {
        const dataDelete = [...subscribers];
        const index = oldData.tableData.id;
        dataDelete.splice(index, 1);
        setSubscribers([...dataDelete]);
        resolve();
      }),
  };
  const [selectedSubscribers, setSelectedSubscribers] = useState<Array<string>>(
    [],
  );
  const {match} = useRouter();
  const [maxPageRowCount, setMaxPageRowCount] = useState(0);
  const [tokenList, setTokenList] = useState(['']);

  const networkId: string = nullthrows(match.params.networkId);
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

  const columns = [
    {
      title: 'IMSI',
      field: 'imsi',
      editComponent: props => (
        <OutlinedInput
          data-testid="IMSI"
          type="text"
          placeholder="Enter IMSI"
          variant="outlined"
          value={props.value}
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
          variant="outlined"
          placeholder="Enter Name"
          type="text"
          value={props.value}
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
          value={props.value || ''}
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
          value={props.value}
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
            value={props.value ?? []}
            onChange={({target}) => props.onChange(target.value)}
            displayEmpty={true}
            renderValue={selected => {
              if (!selected.length) {
                return 'Select Forbidden Network Types';
              }
              return selected.join(', ');
            }}
            input={
              <OutlinedInput
                className={props.value ? '' : classes.placeholder}
              />
            }>
            {forbiddenNetworkTypes.map((k, idx) => (
              <MenuItem key={idx} value={k}>
                <Checkbox
                  checked={props.value ? props.value.indexOf(k) > -1 : false}
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
            value={props.value ?? []}
            onChange={({target}) => props.onChange(target.value)}
            displayEmpty={true}
            renderValue={selected => {
              if (!selected.length) {
                return 'Select APNs';
              }
              return selected.join(', ');
            }}
            input={
              <OutlinedInput
                className={props.value ? '' : classes.placeholder}
              />
            }>
            {apns.map((k: string, idx: number) => (
              <MenuItem key={idx} value={k}>
                <Checkbox
                  checked={props.value ? props.value.indexOf(k) > -1 : false}
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
            value={props.value ?? []}
            onChange={({target}) => props.onChange(target.value)}
            displayEmpty={true}
            renderValue={selected => {
              if (!selected.length) {
                return 'Select Policies';
              }
              return selected.join(', ');
            }}
            input={
              <OutlinedInput
                className={props.value ? '' : classes.placeholder}
              />
            }>
            {policies.map((k: string, idx: number) => (
              <MenuItem key={idx} value={k}>
                <Checkbox
                  checked={props.value ? props.value.indexOf(k) > -1 : false}
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
            subscriberAction === 'delete' && !subscribers.length
              ? getSubscribers
              : subscribers
          }
          columns={
            !(subscriberAction === 'delete')
              ? [
                  {
                    title: '',
                    field: '',
                    width: '5%',
                    editable: 'never',
                    render: rowData => (
                      <Text variant="subtitle3">
                        {rowData.tableData?.id + 1 || ''}
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
            selectionProps: rowData => {
              return {
                checked: selectedSubscribers.includes(rowData.imsi),
                value: rowData.imsi,
                onClick: event => {
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
                    return new Promise((resolve, reject) => {
                      const err = validateSubscribers([newData]);
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
                    icon: forwardRef((props, ref) => (
                      <Button
                        startIcon={<CloudUploadIcon {...props} ref={ref} />}
                        variant="outlined"
                        color="primary">
                        {subscriberAction === 'delete'
                          ? 'Delete from CSV'
                          : 'Upload CSV'}
                      </Button>
                    )),
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
        <Grid container justify="space-between">
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
                const err = validateSubscribers(subscribers);
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

/**
 * Dialog content used to Add/Delete/Update subscribers
 * Displays upload subscriber dropzone or subscriber table
 */
function SubscriberDetailsDialogContent(props: ActionDialogProps) {
  const [addError, setAddError] = useState([]);
  const [subscribers, setSubscribers] = useState<Array<SubscriberInfo>>([]);
  const [upload, setUpload] = useState(false);
  const [rowAdd, setRowAdd] = useState(false);

  const addSubscriberProps = {
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
  };

  return (
    <>
      {!upload ? (
        <SubscriberDetailsTable {...addSubscriberProps} />
      ) : (
        <SubscriberDetailsUpload {...addSubscriberProps} />
      )}
    </>
  );
}

export function validateSubscriberInfo(
  info: SubscriberInfo,
  subscribers: {[string]: subscriber},
  edit?: boolean,
) {
  if (!info.imsi.match(/^(IMSI\d{10,15})$/)) {
    return "imsi invalid, should match '^(IMSId{10,15})$'";
  }
  if (info.imsi in subscribers && !(edit ?? false)) {
    return 'imsi already exists';
  }
  if (info.authKey && !isValidHex(info.authKey)) {
    return 'auth key is not a valid hex';
  }
  if (info.authOpc && !isValidHex(info.authOpc)) {
    return 'auth opc is not a valid hex';
  }
  return '';
}

type SubscriberError = $Values<typeof SUBSCRIBER_ADD_ERRORS>;

/**
 * Checks subscriber fields format
 *
 * @param {Array<SubscriberInfo>} subscribers Array of subcribers to validate
 * @returns {Array<string>} Returns array of errors
 */
export function validateSubscribers(
  subscribers: Array<SubscriberInfo>,
  isDelete?: boolean,
) {
  const errors: {
    [error: SubscriberError]: Array<number>,
  } = {};
  const imsiList = [];

  Object.keys(SUBSCRIBER_ADD_ERRORS).map(error => {
    const subscriberError = SUBSCRIBER_ADD_ERRORS[error];
    errors[subscriberError] = [];
  });
  subscribers.forEach((info, i) => {
    if (!info.authKey && !(isDelete ?? false)) {
      errors[SUBSCRIBER_ADD_ERRORS['REQUIRED_AUTH_KEY']].push(i + 1);
    }
    if (!info?.imsi?.match(/^(IMSI\d{10,15})$/)) {
      errors[SUBSCRIBER_ADD_ERRORS['INVALID_IMSI']].push(i + 1);
    }
    if (info.authKey && !isValidHex(info.authKey)) {
      errors[SUBSCRIBER_ADD_ERRORS['INVALID_AUTH_KEY']].push(i + 1);
    }
    if (info.authOpc && !isValidHex(info.authOpc)) {
      errors[SUBSCRIBER_ADD_ERRORS['INVALID_AUTH_OPC']].push(i + 1);
    }
    if (!info.dataPlan && !(isDelete ?? false)) {
      errors[SUBSCRIBER_ADD_ERRORS['REQUIRED_SUB_PROFILE']].push(i + 1);
    }
    if (imsiList.includes(info.imsi) && !(isDelete ?? false)) {
      errors[SUBSCRIBER_ADD_ERRORS['DUPLICATE_IMSI']].push(i + 1);
    } else {
      imsiList.push(info.imsi);
    }
  });

  const errorList: Array<string> = Object.keys(SUBSCRIBER_ADD_ERRORS)
    .map(error => SUBSCRIBER_ADD_ERRORS[error])
    .reduce((res, errorMessage) => {
      if (errors[errorMessage].length > 0) {
        res.push(
          `${errorMessage} : Row ${errors[errorMessage].sort().join(', ')}`,
        );
      }
      return res;
    }, []);

  return errorList;
}

type subscriberStaticIpsRowType = {
  apnName: string,
  staticIp: string,
};

type subscriberForbiddenNetworkTypes = {
  nwTypes: string,
};

function EditSubscriberDetails(props: EditSubscriberProps) {
  const classes = useStyles();
  return (
    <div>
      <List>
        <AltFormField label={'Subscriber Name'}>
          <OutlinedInput
            data-testid="name"
            className={classes.input}
            placeholder="Enter Name"
            fullWidth={true}
            value={props.subscriberState.name}
            onChange={({target}) =>
              props.onSubscriberChange('name', target.value)
            }
          />
        </AltFormField>
        <AltFormField label={'Service State'}>
          <TypedSelect
            className={classes.input}
            input={<OutlinedInput />}
            value={props.subscriberState.lte.state}
            items={{
              ACTIVE: 'Active',
              INACTIVE: 'Inactive',
            }}
            onChange={value => {
              props.onSubscriberChange('lte', {
                ...props.subscriberState.lte,
                state: value,
              });
            }}
          />
        </AltFormField>
        <AltFormField label={'Data Plan'}>
          <TypedSelect
            className={classes.input}
            input={<OutlinedInput />}
            value={props.subscriberState.lte.sub_profile}
            items={props.subProfiles}
            onChange={value => {
              props.onSubscriberChange('lte', {
                ...props.subscriberState.lte,
                sub_profile: value,
              });
            }}
          />
        </AltFormField>
        <AltFormField label={'Forbidden Network Types'}>
          <FormControl className={classes.input}>
            <Select
              multiple
              value={props.subscriberState.forbidden_network_types ?? []}
              onChange={({target}) => {
                props.onSubscriberChange(
                  'forbidden_network_types',
                  target.value,
                );
              }}
              renderValue={selected => selected.join(', ')}
              input={<OutlinedInput />}>
              {forbiddenNetworkTypes.map((k: string, idx: number) => (
                <MenuItem key={idx} value={k}>
                  <Checkbox
                    checked={
                      props.subscriberState.forbidden_network_types != null
                        ? props.subscriberState.forbidden_network_types.indexOf(
                            k,
                          ) > -1
                        : false
                    }
                  />
                  <ListItemText primary={k} />
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </AltFormField>

        <AltFormField label={'Auth Key'}>
          <PasswordInput
            data-testid="authKey"
            className={classes.input}
            placeholder="Eg. 8baf473f2f8fd09487cccbd7097c6862"
            value={props.authKey}
            error={props.authKey && !isValidHex(props.authKey) ? true : false}
            onChange={v => props.setAuthKey(v)}
          />
        </AltFormField>
        <AltFormField label={'Auth OPC'}>
          <PasswordInput
            data-testid="authOPC"
            value={props.authOpc}
            placeholder="Eg. 8e27b6af0e692e750f32667a3b14605d"
            className={classes.input}
            error={props.authOpc && !isValidHex(props.authOpc) ? true : false}
            onChange={v => props.setAuthOpc(v)}
          />
        </AltFormField>
      </List>
    </div>
  );
}
