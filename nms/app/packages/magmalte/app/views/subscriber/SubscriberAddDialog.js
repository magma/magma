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
import type {mutable_subscribers, subscriber} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import ApnContext from '../../components/context/ApnContext';
import Button from '@material-ui/core/Button';
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
import TypedSelect from '@fbcnms/ui/components/TypedSelect';
import nullthrows from '@fbcnms/util/nullthrows';

import {
  AltFormField,
  LinearProgressWithLabel,
  PasswordInput,
} from '../../components/FormField';
import {SelectEditComponent} from '../../components/ActionTable';
import {base64ToHex, hexToBase64, isValidHex} from '@fbcnms/util/strings';
import {colors, typography} from '../../theme/default';
import {forwardRef} from 'react';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useRef, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
    color: colors.primary.white,
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
}));

const MAX_UPLOAD_FILE_SZ_BYTES = 10 * 1024 * 1024;
const SUBSCRIBER_TITLE = 'Subscriber';
const TRAFFIC_TITLE = 'Traffic Policy';
const STATIC_IPS_TITLE = 'APNs Static IPs';
const SUBSCRIBERS_CHUNK_SIZE = 1000;

type SubscriberInfo = {
  name: string,
  imsi: string,
  authKey: string,
  authOpc: string,
  state: 'INACTIVE' | 'ACTIVE',
  dataPlan: string,
  apns: Array<string>,
  policies?: Array<string>,
};

const SUB_NAME_OFFSET = 0;
const SUB_IMSI_OFFSET = 1;
const SUB_AUTH_KEY_OFFSET = 2;
const SUB_AUTH_OPC_OFFSET = 3;
const SUB_STATE_OFFSET = 4;
const SUB_DATAPLAN_OFFSET = 5;
const SUB_APN_OFFSET = 6;
const SUB_POLICY_OFFSET = 7;
const SUB_MAX_FIELDS = 8;

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
    dataPlan: items[SUB_DATAPLAN_OFFSET],
    apns: items[SUB_APN_OFFSET].split('|')
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

export default function AddSubscriberButton(props: {onClose: () => void}) {
  const classes = useStyles();
  const [open, setOpen] = useState(false);

  return (
    <>
      <AddSubscriberDialog
        open={open}
        onClose={() => {
          setOpen(false);
          props.onClose();
        }}
      />
      <Button onClick={() => setOpen(true)} className={classes.appBarBtn}>
        {'Add Subscriber'}
      </Button>
    </>
  );
}

export function EditSubscriberButton(props: EditProps) {
  const [open, setOpen] = useState(false);
  return (
    <>
      <SubscriberEditDialog
        editProps={props}
        open={open}
        onClose={() => setOpen(false)}
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

type DialogProps = {
  open: boolean,
  onClose: () => void,
  editProps?: EditProps,
};

function AddSubscriberDialog(props: DialogProps) {
  const classes = useStyles();
  return (
    <Dialog
      data-testid="addSubscriberDialog"
      open={props.open}
      fullWidth={true}
      maxWidth="lg">
      <DialogTitle
        className={classes.topBar}
        onClose={props.onClose}
        label={'Add Subscribers'}
      />

      <AddSubscriberDetails {...props} />
    </Dialog>
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
function AddSubscriberDetails(props: DialogProps) {
  const classes = useStyles();
  const ctx = useContext(SubscriberContext);
  const apnCtx = useContext(ApnContext);
  const lteCtx = useContext(LteNetworkContext);
  const policyCtx = useContext(PolicyContext);
  const [error, setError] = useState('');
  const [subscribers, setSubscribers] = useState<Array<SubscriberInfo>>([]);
  const successCountRef = useRef(0);
  const fileInput = useRef(null);
  const enqueueSnackbar = useEnqueueSnackbar();

  const apns = Array.from(new Set(Object.keys(apnCtx.state || {})));

  const subProfiles = Array.from(
    new Set(Object.keys(lteCtx.state.cellular?.epc?.sub_profiles || {})).add(
      'default',
    ),
  );
  const policies = Array.from(
    new Set(Object.keys(policyCtx.state || {})).add('default'),
  );
  const bulkAdd = async (
    addedSubscribers: mutable_subscribers,
    subscriberErrors: string,
  ) => {
    let success = true;
    try {
      if (subscriberErrors.length > 0) {
        setError(subscriberErrors);
        return false;
      }
      await ctx.setState?.('', addedSubscribers);
      return success;
    } catch (e) {
      const errMsg = e.response?.data?.message ?? e.message ?? e;
      setError('error saving subscribers: ' + errMsg);
      success = false;
      return success;
    }
  };

  const saveSubscribers = async () => {
    let addedSubscribers = [];
    let subscriberErrors = '';
    for (const [i, subscriber] of subscribers.entries()) {
      try {
        const err = validateSubscriberInfo(subscriber, ctx.state);
        if (err.length > 0) {
          throw err;
        }
        const authKey =
          subscriber.authKey && isValidHex(subscriber.authKey)
            ? hexToBase64(subscriber.authKey)
            : '';

        const authOpc =
          subscriber.authOpc !== undefined && isValidHex(subscriber.authOpc)
            ? hexToBase64(subscriber.authOpc)
            : '';
        const newSubscriber = {
          active_apns: subscriber.apns,
          active_policies: subscriber.policies,
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
        addedSubscribers.push(newSubscriber);
        // bulk add chunked subscribers
        if (
          addedSubscribers.length == SUBSCRIBERS_CHUNK_SIZE ||
          i == subscribers.length - 1
        ) {
          const success = await bulkAdd(addedSubscribers, subscriberErrors);
          if (success) {
            successCountRef.current =
              successCountRef.current + addedSubscribers.length;
            addedSubscribers = [];
          } else {
            enqueueSnackbar('Saving subscribers to the api failed: ', {
              variant: 'error',
            });
            return;
          }
        }
      } catch (e) {
        const errMsg = e.response?.data?.message ?? e.message ?? e;
        subscriberErrors +=
          'error saving ' + subscriber.imsi + ' : ' + errMsg + '\n';
        //report saved errors if we reach end of loop without calling bulkadd.
        if (i == subscribers.length - 1) {
          setError(subscriberErrors);
          enqueueSnackbar('Saving subscribers to the api failed: ', {
            variant: 'error',
          });
          return;
        }
      }
    }
    enqueueSnackbar(
      ` Subscriber${
        successCountRef.current > 0 ? 's ' : ''
      } saved successfully`,
      {
        variant: 'success',
      },
    );
    props.onClose();
  };

  return (
    <>
      <DialogContent>
        {successCountRef.current > 0 && subscribers.length > 0 && (
          <LinearProgressWithLabel
            value={Math.round(
              (successCountRef.current * 100) / subscribers.length,
            )}
            text={`${successCountRef.current}/${subscribers.length}`}
          />
        )}
        {error !== '' && <FormLabel error>{error}</FormLabel>}
        <input
          type="file"
          ref={fileInput}
          accept={'.csv'}
          style={{display: 'none'}}
          onChange={async () => {
            if (fileInput.current) {
              try {
                const newSubscribers = await parseSubscriberFile(
                  fileInput.current.files[0],
                );
                setSubscribers([...subscribers, ...newSubscribers]);
              } catch (e) {
                enqueueSnackbar(e, {
                  variant: 'error',
                });
              }
            }
          }}
        />
        <ActionTable
          data={subscribers}
          columns={[
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
              title: 'Auth Key',
              field: 'authKey',
              editComponent: props => (
                <PasswordInput
                  data-testid="authKey"
                  placeholder="Key"
                  value={props.value}
                  onChange={v => props.onChange(v)}
                />
              ),
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
                          checked={
                            props.value ? props.value.indexOf(k) > -1 : false
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
                          checked={
                            props.value ? props.value.indexOf(k) > -1 : false
                          }
                        />
                        <ListItemText primary={k} />
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              ),
            },
          ]}
          options={{
            actionsColumnIndex: -1,
            pageSizeOptions: [5, 10],
          }}
          editable={{
            onRowAdd: newData =>
              new Promise((resolve, reject) => {
                const err = validateSubscriberInfo(newData, ctx.state);
                setError(err);
                if (err.length > 0) {
                  return reject();
                }
                setSubscribers([...subscribers, newData]);
                resolve();
              }),
            onRowUpdate: (newData, oldData) =>
              new Promise((resolve, reject) => {
                const err = validateSubscriberInfo(newData, ctx.state);
                setError(err);
                if (err.length > 0) {
                  return reject();
                }
                const dataUpdate = [...subscribers];
                const index = oldData.tableData.id;
                dataUpdate[index] = newData;
                setSubscribers([...dataUpdate]);
                resolve();
              }),
            onRowDelete: oldData =>
              new Promise(resolve => {
                const dataDelete = [...subscribers];
                const index = oldData.tableData.id;
                dataDelete.splice(index, 1);
                setSubscribers([...dataDelete]);
                resolve();
              }),
          }}
          actions={[
            {
              icon: forwardRef((props, ref) => (
                <CloudUploadIcon {...props} ref={ref} />
              )),
              tooltip: 'Upload',
              isFreeAction: true,
              onClick: () => {
                if (fileInput.current) {
                  fileInput.current.click();
                }
              },
            },
          ]}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}> Cancel </Button>
        <Button data-testid="saveSubscriber" onClick={saveSubscribers}>
          {'Save and Add Subscribers'}
        </Button>
      </DialogActions>
    </>
  );
}

function validateSubscriberInfo(
  info: SubscriberInfo,
  subscribers: {[string]: subscriber},
) {
  if (!info.imsi.match(/^(IMSI\d{10,15})$/)) {
    return "imsi invalid, should match '^(IMSId{10,15})$'";
  }
  if (info.imsi in subscribers) {
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
type subscriberStaticIpsRowType = {
  apnName: string,
  staticIp: string,
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
