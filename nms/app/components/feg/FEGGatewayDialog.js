/**
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
  csfb,
  diameter_client_configs,
  federation_gateway,
  gateway_federation_configs,
  gx,
  s8,
  virtual_apn_rule,
} from '../../../generated/MagmaAPIBindings';

import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FEGGatewayContext from '../context/FEGGatewayContext';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import InputLabel from '@material-ui/core/InputLabel';
import KeyValueFields from '../KeyValueFields';
import LoadingFillerBackdrop from '../LoadingFillerBackdrop';
import MagmaV1API from '../../../generated/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useContext, useState} from 'react';
import Select from '@material-ui/core/Select';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import TextField from '@material-ui/core/TextField';

import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPIFlow';

import {
  AddGatewayFields,
  EMPTY_GATEWAY_FIELDS,
  MAGMAD_DEFAULT_CONFIGS,
  // $FlowFixMe migrated to typescript
} from '../AddGatewayDialog';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(() => ({
  appBar: {
    backgroundColor: '#f5f5f5',
    marginBottom: '20px',
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {|
  onClose: () => void,
  onSave: federation_gateway => void,
  editingGateway?: federation_gateway,
  tabOption?: TabOption,
|};

type SCTPValues = {
  server_address: string,
  local_address: string,
};

function getSCTPConfigs(cfg: SCTPValues): csfb {
  return {
    client: {...cfg},
  };
}

function getInitialSCTPConfigs(cfg: ?csfb): SCTPValues {
  return {
    server_address: cfg?.client?.server_address || '',
    local_address: cfg?.client?.local_address || '',
  };
}

type S8Values = {
  local_address: string,
  pgw_address: string,
  apn_operator_suffix: string,
};

function getS8Configs(cfg: S8Values): s8 {
  return {...cfg};
}

function getInitialS8Configs(cfg: ?s8): S8Values {
  return {
    local_address: cfg?.local_address || '',
    pgw_address: cfg?.pgw_address || '',
    apn_operator_suffix: cfg?.apn_operator_suffix || '',
  };
}

type DiameterValues = {
  address: string,
  dest_host: string,
  dest_realm: string,
  host: string,
  realm: string,
  local_address: string,
  product_name: string,
  protocol: $PropertyType<diameter_client_configs, 'protocol'>,
  disable_dest_host: boolean,
};

export const TAB_OPTIONS = Object.freeze({
  GENERAL: 'general',
  GX: 'gx',
  GY: 'gy',
  SWX: 'swx',
  S6A: 's6a',
  S8: 's8',
  CSFB: 'csfb',
});

export type TabOption = $Values<typeof TAB_OPTIONS>;

function getDiameterConfigs(cfg: DiameterValues): gx {
  return {
    server: {...cfg},
  };
}

function getDiameterServerConfig(
  server: ?diameter_client_configs,
): DiameterValues {
  return {
    address: server?.address || '',
    dest_host: server?.dest_host || '',
    dest_realm: server?.dest_realm || '',
    host: server?.host || '',
    realm: server?.realm || '',
    local_address: server?.local_address || '',
    product_name: server?.product_name || '',
    protocol: server?.protocol || 'tcp',
    disable_dest_host: server?.disable_dest_host || false,
  };
}

function getVirtualApnRules(
  rules: ?Array<virtual_apn_rule>,
): ?Array<[string, string]> {
  return rules?.map(entry => {
    return [entry.apn_filter || '', entry.apn_overwrite || ''];
  });
}

function virtualApnRulesToObject(
  props: ?Array<[string, string]>,
): ?Array<virtual_apn_rule> {
  return props
    ?.filter(p => p[0])
    .map(pair => {
      return {apn_filter: pair[0], apn_overwrite: pair[1]};
    });
}

export default function FEGGatewayDialog(props: Props) {
  const classes = useStyles();
  const params = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();

  const ctx = useContext(FEGGatewayContext);
  const {editingGateway, tabOption} = props;
  const [tab, setTab] = useState(
    editingGateway ? tabOption ?? TAB_OPTIONS.GX : TAB_OPTIONS.GENERAL,
  );
  const [generalFields, setGeneralFields] = useState(EMPTY_GATEWAY_FIELDS);
  const [gx, setGx] = useState<DiameterValues>(
    getDiameterServerConfig(editingGateway?.federation?.gx?.server),
  );
  const [gy, setGy] = useState<DiameterValues>(
    getDiameterServerConfig(editingGateway?.federation?.gy?.server),
  );
  const [swx, setSWx] = useState<DiameterValues>(
    getDiameterServerConfig(editingGateway?.federation?.swx?.server),
  );
  const [s6a, setS6A] = useState<DiameterValues>(
    getDiameterServerConfig(editingGateway?.federation?.s6a?.server),
  );
  const [s8, setS8] = useState<S8Values>(
    getInitialS8Configs(editingGateway?.federation?.s8),
  );

  const [csfb, setCSFB] = useState<SCTPValues>(
    getInitialSCTPConfigs(editingGateway?.federation?.csfb),
  );

  const [gxVirtualApnRules, setGxVirtualApnRules] = useState<?Array<
    [string, string],
  >>(getVirtualApnRules(editingGateway?.federation?.gx?.virtual_apn_rules));

  const [gyVirtualApnRules, setGyVirtualApnRules] = useState<?Array<
    [string, string],
  >>(getVirtualApnRules(editingGateway?.federation?.gy?.virtual_apn_rules));

  const networkID = nullthrows(params.networkId);
  const {response: tiers, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdTiers,
    {
      networkId: networkID,
    },
  );

  if (isLoading || !tiers) {
    return <LoadingFillerBackdrop />;
  }

  const getFederationConfigs = (): gateway_federation_configs => ({
    aaa_server: {},
    eap_aka: {},
    gx: {
      server: getDiameterConfigs(gx).server,
      virtual_apn_rules: virtualApnRulesToObject(gxVirtualApnRules) || [],
    },
    gy: {
      server: getDiameterConfigs(gy).server,
      init_method: 2,
      virtual_apn_rules: virtualApnRulesToObject(gyVirtualApnRules) || [],
    },
    health: {},
    hss: {},
    s6a: {...getDiameterConfigs(s6a), plmn_ids: []},
    s8: {...getS8Configs(s8)},
    served_network_ids: [],
    swx: {...getDiameterConfigs(swx)},
    csfb: {...getSCTPConfigs(csfb)},
  });

  const onSave = async () => {
    try {
      if (editingGateway) {
        const editedGateway = {
          ...editingGateway,
          federation: getFederationConfigs(),
        };
        await ctx.setState(editingGateway.id, editedGateway);
      } else {
        const newGateway = {
          device: {
            hardware_id: generalFields.hardwareID,
            key: {
              key: generalFields.challengeKey,
              key_type: 'SOFTWARE_ECDSA_SHA256', // default key type should be ECDSA (do not use ECHO for prod)
            },
          },
          federation: getFederationConfigs(),
          magmad: MAGMAD_DEFAULT_CONFIGS,
          id: generalFields.gatewayID,
          description: generalFields.description,
          name: generalFields.name,
          tier: generalFields.tier,
        };
        await ctx.setState(newGateway.id, newGateway);
      }

      const gateway = await MagmaV1API.getFegByNetworkIdGatewaysByGatewayId({
        networkId: networkID,
        gatewayId: editingGateway?.id || generalFields.gatewayID,
      });
      props.onSave(gateway);
    } catch (e) {
      enqueueSnackbar(e?.response?.data?.message || e?.message || e, {
        variant: 'error',
      });
    }
  };

  let content;
  let contentOverwriteAPN;
  switch (tab) {
    case TAB_OPTIONS.GENERAL:
      content = (
        <AddGatewayFields
          onChange={setGeneralFields}
          values={generalFields}
          tiers={tiers}
        />
      );
      break;
    case TAB_OPTIONS.GX:
      content = (
        <DiameterFields
          onChange={setGx}
          values={gx}
          supportedProtocols={['tcp']}
        />
      );
      contentOverwriteAPN = (
        <KeyValueFields
          key_label="APN Filter"
          value_label="APN Overwrite"
          onChange={setGxVirtualApnRules}
          keyValuePairs={
            gxVirtualApnRules && gxVirtualApnRules.length
              ? gxVirtualApnRules
              : [['', '']]
          }
        />
      );
      break;
    case TAB_OPTIONS.GY:
      content = (
        <DiameterFields
          onChange={setGy}
          values={gy}
          supportedProtocols={['tcp']}
        />
      );
      contentOverwriteAPN = (
        <KeyValueFields
          key_label="APN Filter"
          value_label="APN Overwrite"
          onChange={setGyVirtualApnRules}
          keyValuePairs={
            gyVirtualApnRules && gyVirtualApnRules.length
              ? gyVirtualApnRules
              : [['', '']]
          }
        />
      );
      break;
    case TAB_OPTIONS.SWX:
      content = (
        <DiameterFields
          onChange={setSWx}
          values={swx}
          supportedProtocols={['tcp', 'sctp']}
        />
      );
      break;
    case TAB_OPTIONS.S6A:
      content = (
        <DiameterFields
          onChange={setS6A}
          values={s6a}
          supportedProtocols={['tcp', 'sctp']}
        />
      );
      break;
    case TAB_OPTIONS.S8:
      content = <S8Fields onChange={setS8} values={s8} />;
      break;
    case TAB_OPTIONS.CSFB:
      content = <SCTPFields onChange={setCSFB} values={csfb} />;
      break;
  }

  return (
    <Dialog open={true} onClose={props.onClose} maxWidth="md" scroll="body">
      <AppBar position="static" className={classes.appBar}>
        <Tabs
          indicatorColor="primary"
          textColor="primary"
          value={tab}
          onChange={(event, tab) => setTab(tab)}>
          {!editingGateway && <Tab label="General" value="general" />}
          <Tab label="Gx" value="gx" />
          <Tab label="Gy" value="gy" />
          <Tab label="SWx" value="swx" />
          <Tab label="S6A" value="s6a" />
          <Tab label="S8" value="s8" />
          <Tab label="CSFB" value="csfb" />
        </Tabs>
      </AppBar>
      <DialogContent>
        {content}
        {contentOverwriteAPN}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} color="primary">
          Cancel
        </Button>
        <Button onClick={onSave} color="primary" variant="contained">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}

function S8Fields(props: {values: S8Values, onChange: S8Values => void}) {
  const classes = useStyles();
  const {values} = props;
  const onChange = field => event =>
    // $FlowFixMe Set state for each field
    props.onChange({...values, [field]: event.target.value});

  return (
    <>
      <TextField
        label="Local Address"
        className={classes.input}
        value={values.local_address}
        onChange={onChange('local_address')}
        placeholder="example.magma.com:5555"
        inputProps={{'data-testid': 'localAddress'}}
      />
      <TextField
        label="PGW Address"
        className={classes.input}
        value={values.pgw_address}
        onChange={onChange('pgw_address')}
        placeholder="pgw.magma.com:5555"
        inputProps={{'data-testid': 'pgwAddress'}}
      />
      <TextField
        label="APN Operator Suffix"
        className={classes.input}
        value={values.apn_operator_suffix}
        onChange={onChange('apn_operator_suffix')}
        placeholder=".operator.com"
        inputProps={{'data-testid': 'apnOperatorSuffix'}}
      />
    </>
  );
}

function SCTPFields(props: {values: SCTPValues, onChange: SCTPValues => void}) {
  const classes = useStyles();
  const {values} = props;
  const onChange = field => event =>
    // $FlowFixMe Set state for each field
    props.onChange({...values, [field]: event.target.value});

  return (
    <>
      <TextField
        label="Server Address"
        className={classes.input}
        value={values.server_address}
        onChange={onChange('server_address')}
        placeholder="example.magma.com:5555"
        inputProps={{'data-testid': 'serverAddress'}}
      />
      <TextField
        label="Local Address"
        className={classes.input}
        value={values.local_address}
        onChange={onChange('local_address')}
        placeholder="example.magma.com:5555"
        inputProps={{'data-testid': 'localAddress'}}
      />
    </>
  );
}

function DiameterFields(props: {
  values: DiameterValues,
  onChange: DiameterValues => void,
  supportedProtocols: Array<
    $NonMaybeType<$PropertyType<diameter_client_configs, 'protocol'>>,
  >,
}) {
  const classes = useStyles();
  const {values, supportedProtocols} = props;
  const onChange = field => event =>
    // $FlowFixMe Set state for each field
    props.onChange({...values, [field]: event.target.value});

  return (
    <>
      <TextField
        label="Address"
        className={classes.input}
        value={values.address}
        onChange={onChange('address')}
        placeholder="example.magma.com:5555"
        inputProps={{'data-testid': 'address'}}
      />
      <TextField
        label="Destination Host"
        className={classes.input}
        value={values.dest_host}
        onChange={onChange('dest_host')}
        placeholder="magma-fedgw.magma.com"
        inputProps={{'data-testid': 'destinationHost'}}
      />
      <TextField
        label="Dest Realm"
        className={classes.input}
        value={values.dest_realm}
        onChange={onChange('dest_realm')}
        placeholder="magma.com"
        inputProps={{'data-testid': 'destRealm'}}
      />
      <TextField
        label="Host"
        className={classes.input}
        value={values.host}
        onChange={onChange('host')}
        placeholder="magma.com"
        inputProps={{'data-testid': 'host'}}
      />
      <TextField
        label="Realm"
        className={classes.input}
        value={values.realm}
        onChange={onChange('realm')}
        placeholder="realm"
        inputProps={{'data-testid': 'realm'}}
      />
      <TextField
        label="Local Address"
        className={classes.input}
        value={values.local_address}
        onChange={onChange('local_address')}
        placeholder=":56789"
        inputProps={{'data-testid': 'localAddress'}}
      />
      <TextField
        label="Product Name"
        className={classes.input}
        value={values.product_name}
        onChange={onChange('product_name')}
        placeholder="Magma"
        inputProps={{'data-testid': 'productName'}}
      />
      <FormControl className={classes.input}>
        <InputLabel htmlFor="protocol">Protocol</InputLabel>
        <Select
          inputProps={{id: 'protocol', 'data-testid': 'protocol'}}
          value={values.protocol}
          onChange={({target}) => {
            switch (target.value) {
              case 'tcp':
              case 'tcp4':
              case 'tcp6':
              case 'sctp':
              case 'sctp4':
              case 'sctp6':
                props.onChange({...values, protocol: target.value});
            }
          }}>
          {supportedProtocols.map(item => (
            <MenuItem value={item} key={item}>
              {item.toUpperCase()}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
      <FormControlLabel
        control={
          <Checkbox
            checked={values.disable_dest_host}
            onChange={({target}) =>
              props.onChange({...values, disable_dest_host: target.checked})
            }
            color="primary"
            inputProps={{'data-testid': 'disableDestinationHost'}}
          />
        }
        label="Disable Destination Host"
      />
    </>
  );
}
