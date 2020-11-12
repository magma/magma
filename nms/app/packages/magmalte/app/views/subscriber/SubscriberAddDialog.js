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
import type {subscriber} from '@fbcnms/magma-api';

import Accordion from '@material-ui/core/Accordion';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import ActionTable from '../../components/ActionTable';
import AddIcon from '@material-ui/icons/Add';
import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import DeleteIcon from '@material-ui/icons/Delete';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import LteNetworkContext from '../../components/context/LteNetworkContext';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Select from '@material-ui/core/Select';
import SubscriberContext from '../../components/context/SubscriberContext';
import Text from '@fbcnms/ui/components/design-system/Text';
import TypedSelect from '@fbcnms/ui/components/TypedSelect';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {AltFormField, PasswordInput} from '../../components/FormField';
import {SelectEditComponent} from '../../components/ActionTable';
import {base64ToHex, hexToBase64, isValidHex} from '@fbcnms/util/strings';
import {colors, typography} from '../../theme/default';
import {forwardRef} from 'react';
import {makeStyles} from '@material-ui/styles';
import {useContext, useRef, useState} from 'react';
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
    padding: `0 ${theme.spacing(5)}px`,
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
  appBarBtnSecondary: {
    color: colors.primary.white,
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
  accordionList: {
    width: '100%',
  },
  placeholder: {
    opacity: 0.5,
  },
  apnButton: {
    margin: '20px 0',
  },
}));

const MAX_UPLOAD_FILE_SZ_BYTES = 10 * 1024 * 1024;

type SubscriberInfo = {
  name: string,
  imsi: string,
  authKey: string,
  authOpc: string,
  state: 'INACTIVE' | 'ACTIVE',
  dataPlan: string,
  apns: Array<string>,
};

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
          const items = line.split(',').map(item => item.trim());
          if (items.length != 7) {
            reject('failed parsing ' + line);
            return;
          }
          subscribers.push({
            name: items[0],
            imsi: items[1],
            authKey: items[2],
            authOpc: items[3],
            state: items[4] === 'ACTIVE' ? 'ACTIVE' : 'INACTIVE',
            dataPlan: items[5],
            apns: items[6]
              .split('|')
              .map(item => item.trim())
              .filter(Boolean),
          });
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

export default function AddSubscriberButton() {
  const classes = useStyles();
  const [open, setOpen] = useState(false);

  return (
    <>
      <AddSubscriberDialog open={open} onClose={() => setOpen(false)} />
      <Button onClick={() => setOpen(true)} className={classes.appBarBtn}>
        {'Add Subscriber'}
      </Button>
    </>
  );
}

export function EditSubscriberButton() {
  const [open, setOpen] = useState(false);

  return (
    <>
      <EditSubscriberDialog open={open} onClose={() => setOpen(false)} />
      <Button component="button" variant="text" onClick={() => setOpen(true)}>
        {'Edit'}
      </Button>
    </>
  );
}

type DialogProps = {
  open: boolean,
  onClose: () => void,
};

function AddSubscriberDialog(props: DialogProps) {
  const classes = useStyles();
  return (
    <Dialog
      data-testid="editDialog"
      open={props.open}
      fullWidth={true}
      maxWidth="lg">
      <DialogTitle className={classes.topBar}>
        <Text color="light" weight="medium">
          {'Add Subscribers'}
        </Text>
      </DialogTitle>

      <AddSubscriberDetails {...props} />
    </Dialog>
  );
}

function EditSubscriberDialog(props: DialogProps) {
  const classes = useStyles();
  return (
    <Dialog
      data-testid="editDialog"
      open={props.open}
      fullWidth={true}
      maxWidth="md">
      <DialogTitle className={classes.topBar}>
        <Text color="light" weight="medium">
          {'Edit Subscriber Settings'}
        </Text>
      </DialogTitle>

      <EditSubscriberDetails {...props} />
    </Dialog>
  );
}

function AddSubscriberDetails(props: DialogProps) {
  const classes = useStyles();
  const ctx = useContext(SubscriberContext);
  const {match} = useRouter();

  const [error, setError] = useState('');
  const [subscribers, setSubscribers] = useState<Array<SubscriberInfo>>([]);
  const fileInput = useRef(null);
  const enqueueSnackbar = useEnqueueSnackbar();

  const {isLoading: subProfilesLoading, response: epcConfigs} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularEpc,
    {
      networkId: nullthrows(match.params.networkId),
    },
  );

  const {isLoading: apnsLoading, response: networkAPNs} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdApns,
    {
      networkId: nullthrows(match.params.networkId),
    },
  );

  if (subProfilesLoading || apnsLoading) {
    return <LoadingFiller />;
  }

  const subProfiles = Array.from(
    new Set(Object.keys(epcConfigs?.sub_profiles || {})).add('default'),
  );
  const apns = Array.from(new Set(Object.keys(networkAPNs || {})));

  const saveSubscribers = async () => {
    for (const subscriber of subscribers) {
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
        await ctx.setState?.(subscriber.imsi, {
          active_apns: subscriber.apns,
          id: subscriber.imsi,
          name: subscriber.name,
          lte: {
            auth_algo: 'MILENAGE',
            auth_key: authKey,
            auth_opc: authOpc,
            state: subscriber.state,
            sub_profile: subscriber.dataPlan,
          },
        });
      } catch (e) {
        const errMsg = e.response?.data?.message ?? e.message ?? e;
        setError('error saving ' + subscriber.imsi + ' : ' + errMsg);
        return;
      }
    }
    props.onClose();
  };

  return (
    <>
      <DialogContent>
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
                  variant="outlined"
                  placeholder="Enter Name"
                  type="text"
                  value={props.value}
                  onChange={e => props.onChange(e.target.value)}
                />
              ),
            },
            {
              title: 'IMSI',
              field: 'imsi',
              editComponent: props => (
                <OutlinedInput
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
        <Button onClick={saveSubscribers}> Save and Add Subscribers </Button>
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
function EditSubscriberDetails(props: DialogProps) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(SubscriberContext);
  const lteCtx = useContext(LteNetworkContext);
  const staticIpAssignments =
    lteCtx.state.cellular?.epc?.mobility?.enable_static_ip_assignments;
  const classes = useStyles();
  const {match} = useRouter();
  const subscriberId = nullthrows(match.params.subscriberId);
  const [subscriberState, setSubscriberState] = useState(
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
    Object.keys(ctx.state[subscriberId].config?.static_ips || {}).map(
      (apn: string) => {
        return {
          apnName: apn,
          staticIp: ctx.state[subscriberId].config?.static_ips?.[apn] || '',
        };
      },
    ),
  );
  const [error, setError] = useState('');
  const {isLoading: subProfilesLoading, response: epcConfigs} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularEpc,
    {
      networkId: nullthrows(match.params.networkId),
    },
  );

  const {isLoading: apnsLoading, response: networkAPNs} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdApns,
    {
      networkId: nullthrows(match.params.networkId),
    },
  );

  const saveSubscriber = async () => {
    try {
      if (authOpc !== '') {
        if (isValidHex(authOpc)) {
          subscriberState.lte.auth_opc = hexToBase64(authOpc);
        } else {
          setError('auth_opc is not a valid hex ');
          return;
        }
      }

      if (authKey !== '') {
        if (isValidHex(authKey)) {
          subscriberState.lte.auth_key = hexToBase64(authKey);
        } else {
          setError('auth_key is not a valid hex ');
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
      const errMsg = e.response.data?.message ?? e.message;
      setError('error saving ' + subscriberState.id + ' : ' + errMsg);
      return;
    }
    props.onClose();
  };

  if (subProfilesLoading || apnsLoading) {
    return <LoadingFiller />;
  }

  // we are doing this to ensure we can map subprofiles from an array
  // for e.g. ['foo', 'default'] -> {foo: 'foo', default: 'default}
  // this is done because TypedSelect expects items in this form to verify
  // if the passed in input is of expected type
  const subProfiles = Array.from(
    new Set(Object.keys(epcConfigs?.sub_profiles || {})).add('default'),
  ).reduce(function (o, v) {
    o[v] = v;
    return o;
  }, {});
  const apns = Array.from(new Set(Object.keys(networkAPNs || {})));
  const handleSubscriberChange = (key: string, val) =>
    setSubscriberState({...subscriberState, [key]: val});

  const handleSubscriberStaticIpsChange = (key: string, val, index: number) => {
    const rows = subscriberStaticIPRows;
    rows[index][key] = val;
    setSubscriberStaticIPRows([...rows]);
  };
  const deleteApn = apn => {
    setSubscriberStaticIPRows([
      ...subscriberStaticIPRows.filter(
        (deletedApn: subscriberStaticIpsRowType) => apn !== deletedApn,
      ),
    ]);
  };
  const addApnStaticIP = () => {
    setSubscriberStaticIPRows([
      ...subscriberStaticIPRows,
      {apnName: '', staticIp: ''},
    ]);
  };
  return (
    <>
      <DialogContent>
        {error !== '' && <FormLabel error>{error}</FormLabel>}
        <List component={Paper}>
          <AltFormField label={'Subscriber Name'}>
            <OutlinedInput
              className={classes.input}
              placeholder="Enter Name"
              fullWidth={true}
              value={subscriberState.name}
              onChange={({target}) =>
                handleSubscriberChange('name', target.value)
              }
            />
          </AltFormField>
          <AltFormField label={'Service State'}>
            <TypedSelect
              className={classes.input}
              input={<OutlinedInput />}
              value={subscriberState.lte.state}
              items={{
                ACTIVE: 'Active',
                INACTIVE: 'Inactive',
              }}
              onChange={value => {
                handleSubscriberChange('lte', {
                  ...subscriberState.lte,
                  state: value,
                });
              }}
            />
          </AltFormField>
          <AltFormField label={'Data Plan'}>
            <TypedSelect
              className={classes.input}
              input={<OutlinedInput />}
              value={subscriberState.lte.sub_profile}
              items={subProfiles}
              onChange={value => {
                handleSubscriberChange('lte', {
                  ...subscriberState.lte,
                  sub_profile: value,
                });
              }}
            />
          </AltFormField>
          <AltFormField label={'Auth Key'}>
            <PasswordInput
              className={classes.input}
              placeholder="Enter Auth Key"
              value={authKey}
              error={authKey && !isValidHex(authKey)}
              onChange={v => setAuthKey(v)}
            />
          </AltFormField>
          <AltFormField label={'Auth OPC'}>
            <PasswordInput
              value={authOpc}
              placeholder="Enter Auth OPC"
              className={classes.input}
              error={authOpc && !isValidHex(authOpc)}
              onChange={v => setAuthOpc(v)}
            />
          </AltFormField>
          <AltFormField label={'Active APNs'}>
            <FormControl className={classes.input}>
              <Select
                multiple
                value={subscriberState.active_apns ?? []}
                onChange={({target}) => {
                  handleSubscriberChange('active_apns', target.value);
                }}
                renderValue={selected => selected.join(', ')}
                input={<OutlinedInput />}>
                {apns.map((k: string, idx: number) => (
                  <MenuItem key={idx} value={k}>
                    <Checkbox
                      checked={
                        subscriberState.active_apns != null
                          ? subscriberState.active_apns.indexOf(k) > -1
                          : false
                      }
                    />
                    <ListItemText primary={k} />
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </AltFormField>
        </List>
        <Button
          onClick={addApnStaticIP}
          disabled={!staticIpAssignments ?? false}
          className={classes.apnButton}>
          Add New APN Static IP
          <AddIcon />
        </Button>
        {subscriberStaticIPRows.map((apn, index) => (
          <Accordion>
            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
              <List className={classes.accordionList}>
                <ListItem>
                  <ListItemText
                    primary={
                      apn.apnName || (
                        <Text className={classes.placeholder}>{'APN'}</Text>
                      )
                    }
                  />
                  <ListItemSecondaryAction>
                    <IconButton
                      edge="end"
                      aria-label="delete"
                      onClick={event => {
                        event.stopPropagation();
                        deleteApn(apn);
                      }}>
                      <DeleteIcon />
                    </IconButton>
                  </ListItemSecondaryAction>
                </ListItem>
              </List>
            </AccordionSummary>
            <AccordionDetails>
              <AltFormField label={'APN name'}>
                <FormControl className={classes.input}>
                  <Select
                    value={apn.apnName}
                    onChange={({target}) => {
                      const staticIpApn = subscriberStaticIPRows.map(
                        apn => apn.apnName,
                      );
                      if (!staticIpApn.includes(target.value)) {
                        handleSubscriberStaticIpsChange(
                          'apnName',
                          target.value,
                          index,
                        );
                      }
                    }}
                    input={<OutlinedInput />}>
                    {(subscriberState.active_apns || []).map(apn => (
                      <MenuItem value={apn}>
                        <ListItemText primary={apn} />
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </AltFormField>
              <AltFormField label={'APN Static IP'}>
                <OutlinedInput
                  className={classes.input}
                  placeholder="Enter IP"
                  fullWidth={true}
                  value={apn.staticIp}
                  onChange={({target}) => {
                    handleSubscriberStaticIpsChange(
                      'staticIp',
                      target.value,
                      index,
                    );
                  }}
                />
              </AltFormField>
            </AccordionDetails>
          </Accordion>
        ))}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}> Cancel </Button>
        <Button onClick={saveSubscriber}> Save </Button>
      </DialogActions>
    </>
  );
}
