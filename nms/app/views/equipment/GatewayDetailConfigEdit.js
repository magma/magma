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
import type {
  apn_resources,
  challenge_key,
  distribution_package,
  enodeb_serials,
  gateway_device,
  gateway_dns_configs,
  gateway_epc_configs,
  gateway_he_config,
  gateway_logging_configs,
  gateway_ran_configs,
  lte_gateway,
  magmad_gateway_configs,
} from '../../../generated/MagmaAPIBindings';

import Accordion from '@material-ui/core/Accordion';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AddIcon from '@material-ui/icons/Add';
// $FlowFixMe migrated to typescript
import ApnContext from '../../components/context/ApnContext';
import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import DeleteIcon from '@material-ui/icons/Delete';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
// $FlowFixMe migrated to typescript
import EnodebContext from '../../components/context/EnodebContext';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import GatewayContext from '../../components/context/GatewayContext';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import LteNetworkContext from '../../components/context/LteNetworkContext';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../../theme/design-system/Text';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {AltFormField} from '../../components/FormField';
import {
  DEFAULT_DNS_CONFIG,
  DEFAULT_GATEWAY_CONFIG,
  DEFAULT_HE_CONFIG,
  DynamicServices,
} from '../../components/GatewayUtils';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const GATEWAY_TITLE = 'Gateway';
const RAN_TITLE = 'Ran';
const AGGREGATION_TITLE = 'Aggregation';
const EPC_TITLE = 'Epc';
const APN_RESOURCES_TITLE = 'APN Resources';
const HEADER_ENRICHMENT_TITLE = 'Header Enrichment';

const useStyles = makeStyles(_ => ({
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
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    color: colors.primary.white,
  },
  selectMenu: {
    maxHeight: '200px',
  },
  placeholder: {
    opacity: 0.5,
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
}));

const EditTableType = {
  info: 0,
  aggregation: 1,
  epc: 2,
  ran: 3,
  apnResources: 4,
  headerEnrichment: 5,
};

export type EditProps = {
  editTable: $Keys<typeof EditTableType>,
};

type DialogProps = {
  open: boolean,
  onClose: () => void,
  editProps?: EditProps,
};

type ButtonProps = {
  title: string,
  isLink: boolean,
  editProps?: EditProps,
};

export default function AddEditGatewayButton(props: ButtonProps) {
  const classes = useStyles();
  const [open, setOpen] = React.useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <>
      <GatewayEditDialog
        open={open}
        onClose={handleClose}
        editProps={props.editProps}
      />
      {props.isLink ? (
        <Button
          data-testid={(props.editProps?.editTable ?? '') + 'EditButton'}
          component="button"
          variant="text"
          onClick={handleClickOpen}>
          {props.title}
        </Button>
      ) : (
        <Button
          variant="text"
          className={classes.appBarBtn}
          onClick={handleClickOpen}>
          {props.title}
        </Button>
      )}
    </>
  );
}

function GatewayEditDialog(props: DialogProps) {
  const {open, editProps} = props;
  const classes = useStyles();
  const params = useParams();
  const [gateway, setGateway] = useState<lte_gateway>(DEFAULT_GATEWAY_CONFIG);
  const gatewayId: string = params.gatewayId;
  const [tabPos, setTabPos] = useState(
    editProps ? EditTableType[editProps.editTable] : 0,
  );
  const ctx = useContext(GatewayContext);
  const onClose = () => {
    props.onClose();
  };

  useEffect(() => {
    setTabPos(editProps ? EditTableType[editProps.editTable] : 0);
    setGateway(DEFAULT_GATEWAY_CONFIG);
  }, [editProps, open]);

  return (
    <Dialog data-testid="editDialog" open={open} fullWidth={true} maxWidth="md">
      <DialogTitle
        label={editProps ? 'Edit Gateway' : 'Add New Gateway'}
        onClose={onClose}
      />
      <Tabs
        value={tabPos}
        onChange={(_, v) => setTabPos(v)}
        indicatorColor="primary"
        className={classes.tabBar}>
        <Tab key="gateway" data-testid="gatewayTab" label={GATEWAY_TITLE} />;
        <Tab
          key="aggregations"
          data-testid="aggregationsTab"
          disabled={editProps ? false : true}
          label={AGGREGATION_TITLE}
        />
        <Tab
          key="epc"
          data-testid="EPCTab"
          disabled={editProps ? false : true}
          label={EPC_TITLE}
        />
        <Tab
          key="ran"
          data-testid="ranTab"
          disabled={editProps ? false : true}
          label={RAN_TITLE}
        />
        <Tab
          key="apnResources"
          data-testid="apnResourcesTab"
          disabled={editProps ? false : true}
          label={APN_RESOURCES_TITLE}
        />
        <Tab
          key="headerEnrichment"
          data-testid="headerEnrichmentTab"
          disabled={editProps ? false : true}
          label={HEADER_ENRICHMENT_TITLE}
        />
        ;
      </Tabs>
      {tabPos === 0 && (
        <ConfigEdit
          isAdd={!editProps}
          gateway={!editProps ? gateway : ctx.state[gatewayId]}
          onClose={onClose}
          onSave={(gateway: lte_gateway) => {
            setGateway(gateway);
            if (editProps) {
              onClose();
            } else {
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === 1 && (
        <DynamicServicesEdit
          isAdd={!editProps}
          gateway={!editProps ? gateway : ctx.state[gatewayId]}
          onClose={onClose}
          onSave={(gateway: lte_gateway) => {
            setGateway(gateway);
            if (editProps) {
              onClose();
            } else {
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === 2 && (
        <EPCEdit
          isAdd={!editProps}
          gateway={!editProps ? gateway : ctx.state[gatewayId]}
          onClose={onClose}
          onSave={(gateway: lte_gateway) => {
            setGateway(gateway);
            if (editProps) {
              onClose();
            } else {
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === 3 && (
        <RanEdit
          isAdd={!editProps}
          gateway={!editProps ? gateway : ctx.state[gatewayId]}
          onClose={onClose}
          onSave={(gateway: lte_gateway) => {
            setGateway(gateway);
            if (editProps) {
              onClose();
            } else {
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === 4 && (
        <ApnResourcesEdit
          isAdd={!editProps}
          gateway={!editProps ? gateway : ctx.state[gatewayId]}
          onClose={onClose}
          onSave={(gateway: lte_gateway) => {
            setGateway(gateway);
            if (editProps) {
              onClose();
            } else {
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === 5 && (
        <HeaderEnrichmentConfig
          isAdd={!editProps}
          gateway={!editProps ? gateway : ctx.state[gatewayId]}
          onClose={onClose}
          onSave={(gateway: lte_gateway) => {
            setGateway(gateway);
            if (editProps) {
              onClose();
            }
            onClose();
          }}
        />
      )}
    </Dialog>
  );
}

type Props = {
  isAdd: boolean,
  gateway: lte_gateway,
  onClose: () => void,
  onSave: lte_gateway => void,
};

type VersionType = $PropertyType<distribution_package, 'version'>;

export function ConfigEdit(props: Props) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const ctx = useContext(GatewayContext);

  const [gateway, setGateway] = useState<lte_gateway>({
    ...(props.gateway || DEFAULT_GATEWAY_CONFIG),
    connected_enodeb_serials:
      props.gateway?.connected_enodeb_serials ||
      DEFAULT_GATEWAY_CONFIG.connected_enodeb_serials,
  });

  const [gatewayDevice, SetGatewayDevice] = useState<gateway_device>(
    props.gateway?.device || (DEFAULT_GATEWAY_CONFIG.device ?? {}),
  );

  const [challengeKey, setChallengeKey] = useState<challenge_key>(
    props.gateway?.device?.key || (DEFAULT_GATEWAY_CONFIG.device?.key ?? {}),
  );

  const [gatewayVersion, setGatewayVersion] = useState<VersionType>(
    props.gateway?.status?.platform_info?.packages?.[0].version ||
      DEFAULT_GATEWAY_CONFIG.status?.platform_info?.packages?.[0]?.version,
  );

  const onSave = async () => {
    try {
      const gatewayInfos = {
        ...gateway,
        connected_enodeb_serials:
          props.gateway?.connected_enodeb_serials ||
          DEFAULT_GATEWAY_CONFIG.connected_enodeb_serials,
        status: {
          platform_info: {
            packages: [{version: gatewayVersion}],
          },
        },
        device: {...gatewayDevice, key: challengeKey},
      };
      if (props.isAdd) {
        // check if it is not a modify during add i.e we aren't switching tabs
        // back during add and modifying the information other than the serial
        // number
        if (gateway.id in ctx.state && gateway.id !== props.gateway?.id) {
          setError(`Gateway ${gateway.id} already exists`);
          return;
        }
      }
      await ctx.setState(gateway.id, gatewayInfos);
      enqueueSnackbar('Gateway saved successfully', {
        variant: 'success',
      });
      props.onSave(gatewayInfos);
    } catch (e) {
      setError(e.response?.data?.message ?? e.message);
    }
  };
  return (
    <>
      <DialogContent data-testid="configEdit">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}
          <AltFormField label={'Gateway Name'}>
            <OutlinedInput
              data-testid="name"
              placeholder="Enter Name"
              fullWidth={true}
              value={gateway.name}
              onChange={({target}) => {
                setGateway({...gateway, name: target.value});
              }}
            />
          </AltFormField>
          <AltFormField label={'Gateway ID'}>
            <OutlinedInput
              data-testid="id"
              placeholder="Enter ID"
              fullWidth={true}
              value={gateway.id}
              readOnly={props.gateway.id !== '' ? true : false}
              onChange={({target}) =>
                setGateway({...gateway, id: target.value})
              }
            />
          </AltFormField>
          <AltFormField label={'Hardware UUID'}>
            <OutlinedInput
              data-testid="hardwareId"
              placeholder="Eg. 4dfe212f-df33-4cd2-910c-41892a042fee"
              fullWidth={true}
              value={gatewayDevice.hardware_id}
              onChange={({target}) =>
                SetGatewayDevice({
                  ...gatewayDevice,
                  ['hardware_id']: target.value,
                })
              }
            />
          </AltFormField>
          <AltFormField label={'Version'}>
            <OutlinedInput
              data-testid="version"
              placeholder="Enter Version"
              fullWidth={true}
              value={gatewayVersion}
              readOnly={false}
              onChange={({target}) => setGatewayVersion(target.value)}
            />
          </AltFormField>
          <AltFormField label={'Gateway Description'}>
            <OutlinedInput
              data-testid="description"
              placeholder="Enter Description"
              fullWidth={true}
              value={gateway.description}
              onChange={({target}) =>
                setGateway({...gateway, description: target.value})
              }
            />
          </AltFormField>
          <AltFormField label={'Challenge Key'}>
            <OutlinedInput
              data-testid="challengeKey"
              placeholder="A base64 bytestring of the key in DER format"
              fullWidth={true}
              value={challengeKey.key}
              onChange={({target}) =>
                setChallengeKey({...challengeKey, key: target.value})
              }
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.isAdd ? 'Save And Continue' : 'Save'}
        </Button>
      </DialogActions>
    </>
  );
}

export function DynamicServicesEdit(props: Props) {
  const [error, setError] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(GatewayContext);
  const params = useParams();
  const gatewayId: string = props.gateway.id || nullthrows(params.gatewayId);
  const [config, setConfig] = useState<magmad_gateway_configs>(
    props.gateway.magmad,
  );

  const handleChange = (val: boolean, key: string) => {
    const dynamicServices = [...(config.dynamic_services || [])];
    if (val) {
      dynamicServices.push(key);
      setConfig({
        ...config,
        ['dynamic_services']: dynamicServices,
      });
    } else {
      const index = dynamicServices.indexOf(key);
      if (index !== -1) {
        dynamicServices.splice(index, 1);
        setConfig({
          ...config,
          ['dynamic_services']: dynamicServices,
        });
      }
    }
  };

  const onSave = async () => {
    try {
      if (config.dynamic_services?.includes(DynamicServices.TD_AGENT_BIT)) {
        const logging: gateway_logging_configs = {
          aggregation: {
            target_files_by_tag: {
              mme: 'var/log/mme.log',
            },
          },
          log_level: 'DEBUG',
        };
        config.logging = logging;
      } else {
        if (config.logging) {
          delete config.logging;
        }
      }

      const gateway = {
        ...props.gateway,
        magmad: config,
      };
      await ctx.updateGateway({gatewayId, magmadConfigs: config});
      enqueueSnackbar('Gateway saved successfully', {
        variant: 'success',
      });
      props.onSave(gateway);
    } catch (e) {
      setError(e.response?.data?.message ?? e.message);
    }
  };

  return (
    <>
      <DialogContent data-testid="dynamicServicesEdit">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel error>{error}</FormLabel>
            </AltFormField>
          )}
          <AltFormField label={'Event Aggregation'}>
            <Switch
              data-testid="eventdService"
              onChange={({target}) =>
                handleChange(target.checked, DynamicServices.EVENTD)
              }
              checked={config.dynamic_services?.includes(
                DynamicServices.EVENTD,
              )}
            />
          </AltFormField>
          <AltFormField label={'Log Aggregation'}>
            <Switch
              data-testid="tdAgentService"
              onChange={({target}) =>
                handleChange(target.checked, DynamicServices.TD_AGENT_BIT)
              }
              checked={config.dynamic_services?.includes(
                DynamicServices.TD_AGENT_BIT,
              )}
            />
          </AltFormField>
          <AltFormField label={'CPE Monitoring'}>
            <Switch
              data-testid="monitordService"
              onChange={({target}) =>
                handleChange(target.checked, DynamicServices.MONITORD)
              }
              checked={config.dynamic_services?.includes(
                DynamicServices.MONITORD,
              )}
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.isAdd ? 'Save And Continue' : 'Save'}
        </Button>
      </DialogActions>
    </>
  );
}

export function EPCEdit(props: Props) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const ctx = useContext(GatewayContext);

  const handleEPCChange = (key: string, val) => {
    setEPCConfig({...EPCConfig, [key]: val});
  };

  const [EPCConfig, setEPCConfig] = useState<gateway_epc_configs>(
    props.gateway.cellular.epc,
  );

  useEffect(() => {
    setEPCConfig(props.gateway.cellular.epc);
    setError('');
  }, [props.gateway.cellular.epc]);

  const onSave = async () => {
    try {
      const gateway = {
        ...props.gateway,
        cellular: {
          ...props.gateway.cellular,
          epc: EPCConfig,
        },
      };
      await ctx.updateGateway({gatewayId: gateway.id, epcConfigs: EPCConfig});

      enqueueSnackbar('Gateway saved successfully', {
        variant: 'success',
      });
      props.onSave(gateway);
    } catch (e) {
      setError(e.response?.data?.message ?? e.message);
    }
  };
  return (
    <>
      <DialogContent data-testid="epcEdit">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel error>{error}</FormLabel>
            </AltFormField>
          )}
          <AltFormField label={'Nat Enabled'}>
            <Switch
              data-testid="natEnabled"
              onChange={() =>
                handleEPCChange('nat_enabled', !EPCConfig.nat_enabled)
              }
              checked={EPCConfig.nat_enabled}
            />
          </AltFormField>
          <AltFormField label={'IP Block'}>
            <OutlinedInput
              data-testid="ipBlock"
              placeholder="192.168.128.0/24"
              type="string"
              fullWidth={true}
              value={EPCConfig.ip_block}
              onChange={({target}) => handleEPCChange('ip_block', target.value)}
            />
          </AltFormField>
          <AltFormField label={'IPv6 Block'}>
            <OutlinedInput
              data-testid="ipv6Block"
              placeholder="fdee:5:6c::/48"
              type="string"
              fullWidth={true}
              value={EPCConfig.ipv6_block}
              onChange={({target}) =>
                handleEPCChange('ipv6_block', target.value)
              }
            />
          </AltFormField>
          <AltFormField label={'DNS Primary'}>
            <OutlinedInput
              data-testid="dnsPrimary"
              placeholder="8.8.8.8"
              type="string"
              fullWidth={true}
              value={EPCConfig.dns_primary ?? ''}
              onChange={({target}) =>
                handleEPCChange('dns_primary', target.value)
              }
            />
          </AltFormField>
          <AltFormField label={'DNS Secondary'}>
            <OutlinedInput
              data-testid="dnsSecondary"
              placeholder="8.8.4.4"
              type="string"
              fullWidth={true}
              value={EPCConfig.dns_secondary ?? ''}
              onChange={({target}) =>
                handleEPCChange('dns_secondary', target.value)
              }
            />
          </AltFormField>
          <AltFormField label={'SGi network Gateway IP address'}>
            <OutlinedInput
              data-testid="gwSgiIP"
              placeholder="1.1.1.1"
              required={
                EPCConfig.sgi_management_iface_static_ip ?? false ? true : false
              }
              type="string"
              fullWidth={true}
              value={EPCConfig.sgi_management_iface_gw ?? ''}
              onChange={({target}) =>
                handleEPCChange('sgi_management_iface_gw', target.value)
              }
            />
          </AltFormField>
          <AltFormField label={'SGi management interface IP address'}>
            <OutlinedInput
              data-testid="sgiStaticIP"
              placeholder="1.1.1.1/24"
              type="string"
              fullWidth={true}
              value={EPCConfig.sgi_management_iface_static_ip ?? ''}
              onChange={({target}) =>
                handleEPCChange('sgi_management_iface_static_ip', target.value)
              }
            />
          </AltFormField>
          {!EPCConfig.nat_enabled && (
            <AltFormField label={'SGi management network VLAN id'}>
              <OutlinedInput
                data-testid="sgiVlanID"
                placeholder="100"
                type="string"
                fullWidth={true}
                value={EPCConfig.sgi_management_iface_vlan ?? ''}
                onChange={({target}) =>
                  handleEPCChange('sgi_management_iface_vlan', target.value)
                }
              />
            </AltFormField>
          )}
          <AltFormField label={'SGi management Gateway IPv6 address'}>
            <OutlinedInput
              data-testid="gwSgiIpv6"
              placeholder="2001:4860:4860:0:0:0:0:1"
              type="string"
              fullWidth={true}
              value={EPCConfig.sgi_management_iface_ipv6_gw ?? ''}
              onChange={({target}) =>
                handleEPCChange('sgi_management_iface_ipv6_gw', target.value)
              }
            />
          </AltFormField>
          <AltFormField label={'SGi management interface IPv6 address'}>
            <OutlinedInput
              data-testid="sgiStaticIpv6"
              placeholder="2001:4860:4860:0:0:0:0:8888/64"
              type="string"
              fullWidth={true}
              value={EPCConfig.sgi_management_iface_ipv6_addr ?? ''}
              onChange={({target}) =>
                handleEPCChange('sgi_management_iface_ipv6_addr', target.value)
              }
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.isAdd ? 'Save And Continue' : 'Save'}
        </Button>
      </DialogActions>
    </>
  );
}

export function RanEdit(props: Props) {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const ctx = useContext(GatewayContext);
  const enbsCtx = useContext(EnodebContext);
  const [ranConfig, setRanConfig] = useState<gateway_ran_configs>(
    props.gateway.cellular.ran,
  );
  const [dnsConfig, setDnsConfig] = useState<gateway_dns_configs>(
    props.gateway.cellular.dns ?? {},
  );
  const [connectedEnodebs, setConnectedEnodebs] = useState<enodeb_serials>(
    props.gateway.connected_enodeb_serials,
  );
  const handleRanChange = (key: string, val) => {
    setRanConfig({...ranConfig, [key]: val});
  };
  const handleDnsChange = (key: string, val) => {
    setDnsConfig({...dnsConfig, [key]: val});
  };
  const onSave = async () => {
    try {
      const gateway = {
        ...props.gateway,
        cellular: {
          ...props.gateway.cellular,
          ran: ranConfig,
          dns: {...DEFAULT_DNS_CONFIG, ...dnsConfig},
        },
        connected_enodeb_serials: connectedEnodebs,
      };
      await ctx.updateGateway({
        gatewayId: gateway.id,
        enbs: connectedEnodebs,
        ranConfigs: ranConfig,
        dnsConfig: Object.keys(dnsConfig).length
          ? {...DEFAULT_DNS_CONFIG, ...dnsConfig}
          : undefined,
      });
      enqueueSnackbar('Gateway saved successfully', {
        variant: 'success',
      });
      props.onSave(gateway);
    } catch (e) {
      setError(e.response?.data?.message ?? e.message);
    }
  };
  return (
    <>
      <DialogContent data-testid="ranEdit">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel error>{error}</FormLabel>
            </AltFormField>
          )}
          <AltFormField label={'PCI'}>
            <OutlinedInput
              disabled={!(dnsConfig?.dhcp_server_enabled ?? true)}
              data-testid="pci"
              placeholder="Enter PCI"
              type="number"
              fullWidth={true}
              value={ranConfig.pci}
              onChange={({target}) =>
                handleRanChange('pci', parseInt(target.value))
              }
            />
          </AltFormField>
          <AltFormField label={'Registered eNodeBs'}>
            <Select
              multiple
              variant={'outlined'}
              fullWidth={true}
              displayEmpty={true}
              value={connectedEnodebs}
              onChange={({target}) => {
                setConnectedEnodebs(Array.from(target.value));
              }}
              MenuProps={{classes: {paper: classes.selectMenu}}}
              renderValue={selected => {
                if (!selected.length) {
                  return 'Select eNodeBs';
                }
                return selected.join(', ');
              }}
              input={
                <OutlinedInput
                  data-testid="registeredEnodeb"
                  className={connectedEnodebs.length ? '' : classes.placeholder}
                />
              }>
              {enbsCtx?.state &&
                Object.keys(enbsCtx.state.enbInfo).map(enbSerial => (
                  <MenuItem key={enbSerial} value={enbSerial}>
                    <Checkbox checked={connectedEnodebs.includes(enbSerial)} />
                    <ListItemText
                      primary={enbsCtx.state.enbInfo[enbSerial].enb.name}
                      secondary={enbSerial}
                    />
                  </MenuItem>
                ))}
            </Select>
          </AltFormField>
          <AltFormField label={'Transmit Enabled'}>
            <Switch
              disabled={!(dnsConfig?.dhcp_server_enabled ?? true)}
              onChange={() =>
                handleRanChange('transmit_enabled', !ranConfig.transmit_enabled)
              }
              checked={ranConfig.transmit_enabled}
            />
          </AltFormField>
          <AltFormField label={'eNodeB DHCP Service'}>
            <Switch
              data-testid="enbDhcpService"
              onChange={() =>
                handleDnsChange(
                  'dhcp_server_enabled',
                  !(dnsConfig?.dhcp_server_enabled ?? true),
                )
              }
              checked={dnsConfig?.dhcp_server_enabled ?? true}
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.isAdd ? 'Save And Continue' : 'Save'}
        </Button>
      </DialogActions>
    </>
  );
}

export function ApnResourcesEdit(props: Props) {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const ctx = useContext(GatewayContext);
  const apnCtx = useContext(ApnContext);
  const lteCtx = useContext(LteNetworkContext);
  const apnResources: apn_resources = props.gateway.apn_resources ?? {};
  const [apnResourcesRows, setApnResourcesRows] = useState(
    Object.keys(apnResources).map(apn => apnResources[apn]),
  );
  const handleApnResourcesChange = (key: string, val, index: number) => {
    const rows = apnResourcesRows;
    rows[index][key] = val;
    setApnResourcesRows([...rows]);
  };
  const deleteApn = deletedApn =>
    setApnResourcesRows([
      ...apnResourcesRows.filter(apn => apn !== deletedApn),
    ]);

  const addApnResource = () => {
    setApnResourcesRows([
      ...apnResourcesRows,
      {apn_name: '', id: '', vlan_id: null},
    ]);
  };

  const onSave = async () => {
    try {
      const gatewayApnResources = {};
      apnResourcesRows.forEach(
        apn => (gatewayApnResources[apn.apn_name] = apn),
      );
      const gateway = {
        ...props.gateway,
        apn_resources: gatewayApnResources,
      };
      await ctx.setState(gateway.id, gateway);

      enqueueSnackbar('Gateway saved successfully', {
        variant: 'success',
      });
      props.onSave(gateway);
    } catch (e) {
      setError(e.response?.data?.message ?? e.message);
    }
  };

  return (
    <>
      <DialogContent data-testid="apnResourcesEdit">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel error>{error}</FormLabel>
            </AltFormField>
          )}
          <Button
            data-testid="apnResourcesAdd"
            onClick={addApnResource}
            disabled={
              !lteCtx.state.cellular.epc.mobility
                ?.enable_multi_apn_ip_allocation ?? false
            }>
            Add New APN Resource
            <AddIcon />
          </Button>
          {apnResourcesRows.map((apn, index) => (
            <Accordion key={index}>
              <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <List className={classes.accordionList}>
                  <ListItem>
                    <ListItemText
                      primary={
                        apn.apn_name || (
                          <Text
                            className={
                              apn.apn_name.length ? '' : classes.placeholder
                            }>
                            {'APN'}
                          </Text>
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
                      data-testid="apnName"
                      value={apn.apn_name}
                      onChange={({target}) => {
                        const apns = apnResourcesRows;
                        apns[index].apn_name = target.value;
                        setApnResourcesRows([...apns]);
                      }}
                      input={<OutlinedInput />}>
                      {(Object.keys(apnCtx.state) || []).map(apn => (
                        <MenuItem key={apn} value={apn}>
                          <ListItemText primary={apn} />
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </AltFormField>
                <AltFormField label={'APN Resource ID'}>
                  <OutlinedInput
                    data-testid="apnID"
                    className={classes.input}
                    placeholder="Enter ID"
                    fullWidth={true}
                    value={apn.id}
                    onChange={({target}) => {
                      const apns = apnResourcesRows;
                      apns[index].id = target.value;
                      setApnResourcesRows([...apns]);
                    }}
                  />
                </AltFormField>
                <AltFormField label={'VLAN ID'}>
                  <OutlinedInput
                    data-testid="vlanID"
                    className={classes.input}
                    type="number"
                    placeholder="Enter number"
                    fullWidth={true}
                    value={apn.vlan_id}
                    onChange={({target}) => {
                      handleApnResourcesChange(
                        'vlan_id',
                        parseInt(target.value),
                        index,
                      );
                    }}
                  />
                </AltFormField>
              </AccordionDetails>
            </Accordion>
          ))}
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.isAdd ? 'Save And Continue' : 'Save'}
        </Button>
      </DialogActions>
    </>
  );
}

export function HeaderEnrichmentConfig(props: Props) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const ctx = useContext(GatewayContext);
  const [heConfig, setHeConfig] = useState<gateway_he_config>(
    props.gateway.cellular.he_config || DEFAULT_HE_CONFIG,
  );

  const handleHEChange = (key: string, val) => {
    setHeConfig({...heConfig, [key]: val});
  };
  const heEncodingTypes = ['BASE64', 'HEX2BIN'];
  const heEncryptionAlgorithmTypes = [
    'RC4',
    'AES256_CBC_HMAC_MD5',
    'AES256_ECB_HMAC_MD5',
    'GZIPPED_AES256_ECB_SHA1',
  ];
  const heHashFunctionTypes = ['MD5', 'HEX', 'SHA256'];

  const onSave = async () => {
    try {
      const gateway = {
        ...props.gateway,
        cellular: {
          ...props.gateway.cellular,
          he_config: heConfig.enable_header_enrichment
            ? {
                ...heConfig,
                encryption_key: heConfig.enable_encryption
                  ? heConfig.encryption_key
                  : '',
              }
            : undefined,
        },
      };
      await ctx.updateGateway({
        gatewayId: gateway.id,
        cellularConfigs: gateway.cellular,
      });
      enqueueSnackbar('Gateway saved successfully', {
        variant: 'success',
      });
      props.onSave(gateway);
    } catch (e) {
      setError(e.response?.data?.message ?? e.message);
    }
  };

  return (
    <>
      <DialogContent data-testid="headerEnrichmentEdit">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel error>{error}</FormLabel>
            </AltFormField>
          )}
          <AltFormField label={'Enable Header Enrichment'}>
            <Switch
              data-testid="enableHE"
              onChange={() =>
                handleHEChange(
                  'enable_header_enrichment',
                  !(heConfig?.enable_header_enrichment ?? false),
                )
              }
              checked={heConfig?.enable_header_enrichment ?? false}
            />
          </AltFormField>
          <AltFormField label={'Enable Encryption'}>
            <Switch
              data-testid="enableEncryption"
              disabled={!heConfig.enable_header_enrichment}
              onChange={() =>
                handleHEChange(
                  'enable_encryption',
                  !(heConfig?.enable_encryption ?? false),
                )
              }
              checked={heConfig?.enable_encryption ?? false}
            />
          </AltFormField>

          {heConfig.enable_encryption && (
            <Grid data-testid="encryptionEdit">
              <AltFormField label={'Encryption Key'}>
                <OutlinedInput
                  disabled={!heConfig.enable_header_enrichment}
                  data-testid="encryptionKey"
                  type="string"
                  fullWidth={true}
                  value={
                    heConfig.encryption_key ?? DEFAULT_HE_CONFIG.encryption_key
                  }
                  onChange={({target}) =>
                    handleHEChange('encryption_key', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'Encoding Type'}>
                <Select
                  disabled={!heConfig.enable_header_enrichment}
                  fullWidth={true}
                  variant={'outlined'}
                  value={
                    heConfig.he_encoding_type ??
                    DEFAULT_HE_CONFIG.he_encoding_type
                  }
                  onChange={({target}) => {
                    handleHEChange('he_encoding_type', target.value);
                  }}
                  input={<OutlinedInput id="encodingType" />}>
                  {heEncodingTypes.map(type => (
                    <MenuItem key={type} value={type}>
                      <ListItemText primary={type} />
                    </MenuItem>
                  ))}
                </Select>
              </AltFormField>
              <AltFormField label={'Encryption Algorithm'}>
                <Select
                  disabled={!heConfig.enable_header_enrichment}
                  fullWidth={true}
                  variant={'outlined'}
                  value={
                    heConfig.he_encryption_algorithm ??
                    DEFAULT_HE_CONFIG.he_encoding_type
                  }
                  onChange={({target}) => {
                    handleHEChange('he_encryption_algorithm', target.value);
                  }}
                  input={<OutlinedInput id="encryptionAlgorithm" />}>
                  {heEncryptionAlgorithmTypes.map(type => (
                    <MenuItem key={type} value={type}>
                      <ListItemText primary={type} />
                    </MenuItem>
                  ))}
                </Select>
              </AltFormField>
              <AltFormField label={'Hash Function'}>
                <Select
                  disabled={!heConfig.enable_header_enrichment}
                  fullWidth={true}
                  variant={'outlined'}
                  value={
                    heConfig.he_hash_function ??
                    DEFAULT_HE_CONFIG.he_encoding_type
                  }
                  onChange={({target}) => {
                    handleHEChange('he_hash_function', target.value);
                  }}
                  input={<OutlinedInput id="hashFunction" />}>
                  {heHashFunctionTypes.map(type => (
                    <MenuItem key={type} value={type}>
                      <ListItemText primary={type} />
                    </MenuItem>
                  ))}
                </Select>
              </AltFormField>
            </Grid>
          )}
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.isAdd ? 'Save And Close' : 'Save'}
        </Button>
      </DialogActions>
    </>
  );
}
