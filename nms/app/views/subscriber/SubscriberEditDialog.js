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
import type {SubscriberInfo} from './SubscriberUtils';
import type {subscriber} from '../../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DialogTitle from '../../theme/design-system/DialogTitle';
import EditSubscriberApnStaticIps from './SubscriberApnStaticIpsEdit';
import EditSubscriberTrafficPolicy from './SubscriberTrafficPolicyEdit';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import ListItemText from '@material-ui/core/ListItemText';
// $FlowFixMe migrated to typescript
import LteNetworkContext from '../../components/context/LteNetworkContext';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
// $FlowFixMe migrated to typescript
import SubscriberContext from '../../components/context/SubscriberContext';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import TypedSelect from '../../components/TypedSelect';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {AltFormField, PasswordInput} from '../../components/FormField';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {CoreNetworkTypes} from './SubscriberUtils';
// $FlowFixMe[cannot-resolve-module]
import {base64ToHex, hexToBase64, isValidHex} from '../../util/strings';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(() => ({
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
  dialog: {
    height: '750px',
  },
}));

const SUBSCRIBER_TITLE = 'Subscriber';
const TRAFFIC_TITLE = 'Traffic Policy';
const STATIC_IPS_TITLE = 'APNs Static IPs';
const forbiddenNetworkTypes = Object.keys(CoreNetworkTypes).map(
  key => CoreNetworkTypes[key],
);

export const EditTableType = {
  subscriber: 0,
  trafficPolicy: 1,
  staticIps: 2,
};

export type EditProps = {
  editTable: $Keys<typeof EditTableType>,
};
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

export function SubscriberEditDialog(props: DialogProps) {
  const {editProps} = props;
  const enqueueSnackbar = useEnqueueSnackbar();
  const [tabPos, setTabPos] = useState(
    editProps ? EditTableType[editProps.editTable] : 0,
  );
  const ctx = useContext(SubscriberContext);
  const lteCtx = useContext(LteNetworkContext);
  const classes = useStyles();
  const params = useParams();
  const subscriberId = nullthrows(params.subscriberId);
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
        <Button onClick={props.onClose}>Close</Button>
        <Button
          data-testid={`${props.editProps?.editTable || ''}-saveButton`}
          variant="contained"
          color="primary"
          onClick={onSave}>
          Save
        </Button>
      </DialogActions>
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
  forbiddenNetworkTypes: Array<subscriberForbiddenNetworkTypes>,
  authKey: string,
  authOpc: string,
  setAuthKey: (key: string) => void,
  setAuthOpc: (key: string) => void,
};

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
