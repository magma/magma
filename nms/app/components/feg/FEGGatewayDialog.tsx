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
 */

import type {
  Csfb,
  DiameterClientConfigs,
  FederationGateway,
  GatewayFederationConfigs,
  Gx,
  MutableFederationGateway,
  S8,
  VirtualApnRule,
} from '../../../generated';

import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import FEGGatewayContext from '../../context/FEGGatewayContext';
import KeyValueFields from '../KeyValueFields';
import LoadingFillerBackdrop from '../LoadingFillerBackdrop';
import MagmaAPI from '../../api/MagmaAPI';
import MenuItem from '@mui/material/MenuItem';
import OutlinedInput from '@mui/material/OutlinedInput';
import React, {ChangeEvent, useContext, useState} from 'react';
import Select from '@mui/material/Select';
import Switch from '@mui/material/Switch';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPI';
import {
  AddGatewayFields,
  EMPTY_GATEWAY_FIELDS,
  MAGMAD_DEFAULT_CONFIGS,
} from '../AddGatewayDialog';
import {AltFormField} from '../FormField';
import {colors} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(() => ({
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    color: colors.primary.white,
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  onClose: () => void;
  onSave: (gateway: FederationGateway) => void;
  editingGateway?: FederationGateway;
  tabOption?: TabOption;
};

type SCTPValues = {
  server_address: string;
  local_address: string;
};

function getSCTPConfigs(cfg: SCTPValues): Csfb {
  return {
    client: {...cfg},
  };
}

function getInitialSCTPConfigs(cfg: Csfb | null | undefined): SCTPValues {
  return {
    server_address: cfg?.client?.server_address || '',
    local_address: cfg?.client?.local_address || '',
  };
}

type S8Values = {
  local_address: string;
  pgw_address: string;
  apn_operator_suffix: string;
};

function getS8Configs(cfg: S8Values): S8 {
  return {...cfg};
}

function getInitialS8Configs(cfg: S8 | null | undefined): S8Values {
  return {
    local_address: cfg?.local_address || '',
    pgw_address: cfg?.pgw_address || '',
    apn_operator_suffix: cfg?.apn_operator_suffix || '',
  };
}

type DiameterValues = {
  address: string;
  dest_host: string;
  dest_realm: string;
  host: string;
  realm: string;
  local_address: string;
  product_name: string;
  protocol: DiameterClientConfigs['protocol'];
  disable_dest_host: boolean;
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

export type TabOption = typeof TAB_OPTIONS[keyof typeof TAB_OPTIONS];

function getDiameterConfigs(cfg: DiameterValues): Gx {
  return {
    server: {...cfg},
  };
}

function getDiameterServerConfig(
  server: DiameterClientConfigs | null | undefined,
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
  rules: Array<VirtualApnRule> | null | undefined,
): Array<[string, string]> | null | undefined {
  return rules?.map(entry => {
    return [entry.apn_filter || '', entry.apn_overwrite || ''];
  });
}

function virtualApnRulesToObject(
  props: Array<[string, string]> | null | undefined,
): Array<VirtualApnRule> | null | undefined {
  return props
    ?.filter(p => p[0])
    .map(pair => {
      return {apn_filter: pair[0], apn_overwrite: pair[1]};
    });
}

export function FEGAddGatewayButton() {
  const [isVisible, setIsVisible] = useState(false);

  const handleClose = () => setIsVisible(false);

  return (
    <>
      <Button
        onClick={() => setIsVisible(true)}
        color="primary"
        size="small"
        variant="contained">
        Add Gateway
      </Button>
      {isVisible && (
        <FEGGatewayDialog onClose={handleClose} onSave={handleClose} />
      )}
    </>
  );
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

  const [gxVirtualApnRules, setGxVirtualApnRules] = useState<
    Array<[string, string]> | undefined | null
  >(getVirtualApnRules(editingGateway?.federation?.gx?.virtual_apn_rules));

  const [gyVirtualApnRules, setGyVirtualApnRules] = useState<
    Array<[string, string]> | undefined | null
  >(getVirtualApnRules(editingGateway?.federation?.gy?.virtual_apn_rules));

  const networkID = nullthrows(params.networkId);
  const {response: tiers, isLoading} = useMagmaAPI(
    MagmaAPI.upgrades.networksNetworkIdTiersGet,
    {
      networkId: networkID,
    },
  );

  if (isLoading || !tiers) {
    return <LoadingFillerBackdrop />;
  }

  const getFederationConfigs = (): GatewayFederationConfigs => ({
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
        const newGateway: MutableFederationGateway = {
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

      const gateway = (
        await MagmaAPI.federationGateways.fegNetworkIdGatewaysGatewayIdGet({
          networkId: networkID,
          gatewayId: editingGateway?.id || generalFields.gatewayID,
        })
      ).data;
      props.onSave(gateway);
    } catch (e) {
      enqueueSnackbar(getErrorMessage(e), {
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
    <Dialog
      open={true}
      onClose={props.onClose}
      maxWidth="md"
      scroll="body"
      data-testid="FEGGatewayDialog">
      <DialogTitle
        label={editingGateway ? 'Edit Gateway' : 'Add New Gateway'}
        onClose={props.onClose}
      />
      <Tabs
        value={tab}
        className={classes.tabBar}
        indicatorColor="primary"
        onChange={(event, tab) => setTab(tab as string)}>
        {!editingGateway && <Tab label="General" value="general" />}
        <Tab label="Gx" value="gx" />
        <Tab label="Gy" value="gy" />
        <Tab label="SWx" value="swx" />
        <Tab label="S6A" value="s6a" />
        <Tab label="S8" value="s8" />
        <Tab label="CSFB" value="csfb" />
      </Tabs>
      <DialogContent>
        {content}
        {contentOverwriteAPN}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button
          onClick={() => void onSave()}
          color="primary"
          variant="contained">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}

function S8Fields(props: {
  values: S8Values;
  onChange: (values: S8Values) => void;
}) {
  const {values} = props;
  const onChange = (field: keyof S8Values) => (
    event: ChangeEvent<HTMLInputElement>,
  ) => props.onChange({...values, [field]: event.target.value});

  return (
    <>
      <AltFormField label="Local Address">
        <OutlinedInput
          fullWidth={true}
          value={values.local_address}
          onChange={onChange('local_address')}
          placeholder="example.magma.com:5555"
          inputProps={{'data-testid': 'localAddress'}}
        />
      </AltFormField>
      <AltFormField label="PGW Address">
        <OutlinedInput
          fullWidth={true}
          value={values.pgw_address}
          onChange={onChange('pgw_address')}
          placeholder="pgw.magma.com:5555"
          inputProps={{'data-testid': 'pgwAddress'}}
        />
      </AltFormField>
      <AltFormField label="APN Operator Suffix">
        <OutlinedInput
          fullWidth={true}
          value={values.apn_operator_suffix}
          onChange={onChange('apn_operator_suffix')}
          placeholder=".operator.com"
          inputProps={{'data-testid': 'apnOperatorSuffix'}}
        />
      </AltFormField>
    </>
  );
}

function SCTPFields(props: {
  values: SCTPValues;
  onChange: (sctpValues: SCTPValues) => void;
}) {
  const {values} = props;
  const onChange = (field: keyof SCTPValues) => (
    event: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>,
  ) => props.onChange({...values, [field]: event.target.value});

  return (
    <>
      <AltFormField label="Server Address">
        <OutlinedInput
          fullWidth={true}
          value={values.server_address}
          onChange={onChange('server_address')}
          placeholder="example.magma.com:5555"
          inputProps={{'data-testid': 'serverAddress'}}
        />
      </AltFormField>
      <AltFormField label="Local Address">
        <OutlinedInput
          fullWidth={true}
          value={values.local_address}
          onChange={onChange('local_address')}
          placeholder="example.magma.com:5555"
          inputProps={{'data-testid': 'localAddress'}}
        />
      </AltFormField>
    </>
  );
}

function DiameterFields(props: {
  values: DiameterValues;
  onChange: (diameterValues: DiameterValues) => void;
  supportedProtocols: Array<Required<DiameterClientConfigs['protocol']>>;
}) {
  const {values, supportedProtocols} = props;
  const onChange = (field: keyof DiameterValues) => (
    event: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>,
  ) => props.onChange({...values, [field]: event.target.value});

  return (
    <>
      <AltFormField label="Address">
        <OutlinedInput
          fullWidth={true}
          placeholder="example.magma.com:5555"
          value={values.address}
          onChange={onChange('address')}
          inputProps={{'data-testid': 'address'}}
        />
      </AltFormField>
      <AltFormField label="Destination Host">
        <OutlinedInput
          fullWidth={true}
          onChange={onChange('dest_host')}
          placeholder="magma-fedgw.magma.com"
          value={values.dest_host}
          inputProps={{'data-testid': 'destinationHost'}}
        />
      </AltFormField>
      <AltFormField label="Dest Realm">
        <OutlinedInput
          fullWidth={true}
          onChange={onChange('dest_realm')}
          placeholder="magma.com"
          value={values.dest_realm}
          inputProps={{'data-testid': 'destRealm'}}
        />
      </AltFormField>
      <AltFormField label="Host">
        <OutlinedInput
          fullWidth={true}
          onChange={onChange('host')}
          placeholder="magma.com"
          value={values.host}
          inputProps={{'data-testid': 'host'}}
        />
      </AltFormField>
      <AltFormField label="Realm">
        <OutlinedInput
          fullWidth={true}
          onChange={onChange('realm')}
          placeholder="realm"
          value={values.realm}
          inputProps={{'data-testid': 'realm'}}
        />
      </AltFormField>

      <AltFormField label="Local Address">
        <OutlinedInput
          fullWidth={true}
          onChange={onChange('local_address')}
          placeholder=":56789"
          value={values.local_address}
          inputProps={{'data-testid': 'localAddress'}}
        />
      </AltFormField>

      <AltFormField label="Product Name">
        <OutlinedInput
          fullWidth={true}
          onChange={onChange('product_name')}
          placeholder="Magma"
          value={values.product_name}
          inputProps={{'data-testid': 'productName'}}
        />
      </AltFormField>

      <AltFormField label="Protocol">
        <Select
          fullWidth={true}
          variant={'outlined'}
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
              {item!.toUpperCase()}
            </MenuItem>
          ))}
        </Select>
      </AltFormField>

      <AltFormField label="Disable Destination Host">
        <Switch
          checked={values.disable_dest_host}
          onChange={({target}) =>
            props.onChange({...values, disable_dest_host: target.checked})
          }
        />
      </AltFormField>
    </>
  );
}
