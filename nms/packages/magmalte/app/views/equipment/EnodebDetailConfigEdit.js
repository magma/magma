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
  enodeb,
  enodeb_configuration,
  network_ran_configs,
  unmanaged_enodeb_configuration,
} from '../../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import EnodebContext from '../../components/context/EnodebContext';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import {
  EnodebBandwidthOption,
  EnodebDeviceClass,
} from '../../components/lte/EnodebUtils';

import EnodeConfigEditFdd from './EnodebDetailConfigFdd';
import EnodeConfigEditTdd from './EnodebDetailConfigTdd';

import {AltFormField} from '../../components/FormField';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const CONFIG_TITLE = 'Config';
const RAN_TITLE = 'Ran';
const DEFAULT_ENB_CONFIG = {
  name: '',
  serial: '',
  description: '',
  config: {
    cell_id: 0,
    device_class: 'Baicells Nova-233 G2 OD FDD',
    transmit_enabled: false,
  },
};

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
}));

const EditTableType = {
  config: 0,
  ran: 1,
};

type EditProps = {
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

export default function AddEditEnodeButton(props: ButtonProps) {
  const classes = useStyles();
  const [open, setOpen] = useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <>
      <EnodeEditDialog
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

function EnodeEditDialog(props: DialogProps) {
  const {open, editProps} = props;
  const classes = useStyles();
  const [enb, setEnb] = useState<enodeb>({});
  const {match} = useRouter();
  const ctx = useContext(EnodebContext);
  const enodebSerial: string = match.params.enodebSerial;
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const lteRanConfigs = ctx.lteRanConfigs;

  const [tabPos, setTabPos] = useState(
    editProps ? EditTableType[editProps.editTable] : 0,
  );

  const onClose = () => {
    // clear existing state
    props.onClose();
  };

  useEffect(() => {
    setTabPos(editProps ? EditTableType[editProps.editTable] : 0);
    setEnb({});
  }, [editProps, open]);

  return (
    <Dialog data-testid="editDialog" open={open} fullWidth={true} maxWidth="sm">
      <DialogTitle
        label={editProps ? 'Edit eNodeB' : 'Add New eNodeB'}
        onClose={onClose}
      />
      <Tabs
        value={tabPos}
        onChange={(_, v) => setTabPos(v)}
        indicatorColor="primary"
        className={classes.tabBar}>
        <Tab key="config" data-testid="configTab" label={CONFIG_TITLE} />; ;
        <Tab
          key="ran"
          data-testid="ranTab"
          disabled={editProps ? false : true}
          label={RAN_TITLE}
        />
      </Tabs>
      {tabPos === 0 && (
        <ConfigEdit
          isAdd={!editProps}
          enb={Object.keys(enb).length != 0 ? enb : enbInfo?.enb}
          lteRanConfigs={lteRanConfigs}
          onClose={onClose}
          onSave={(enb: enodeb) => {
            setEnb(enb);
            if (editProps) {
              onClose();
            } else {
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === 1 && (
        <RanEdit
          isAdd={!editProps}
          enb={Object.keys(enb).length != 0 ? enb : enbInfo?.enb}
          lteRanConfigs={lteRanConfigs}
          onClose={onClose}
          onSave={onClose}
        />
      )}
    </Dialog>
  );
}

type Props = {
  isAdd: boolean,
  enb?: enodeb,
  lteRanConfigs: ?network_ran_configs,
  onClose: () => void,
  onSave: enodeb => void,
};

type BandwidthMhzType = $PropertyType<enodeb_configuration, 'bandwidth_mhz'>;

type OptConfig = {
  earfcndl: string,
  bandwidthMhz: BandwidthMhzType,
  specialSubframePattern: string,
  subframeAssignment: string,
  pci: string,
  tac: string,
  a1_threshold_rsrp: string,
  lte_a1_threshold_rsrq: string,
  hysteresis: string,
  time_to_trigger: string,
  a2_threshold_rsrp: string,
  lte_a2_threshold_rsrq: string,
  lte_a2_threshold_rsrp_irat_volte: string,
  lte_a2_threshold_rsrq_irat_volte: string,
  a3_offset: string,
  a3_offset_anr: string,
  a4_threshold_rsrp: string,
  lte_intra_a5_threshold_1_rsrp: string,
  lte_intra_a5_threshold_2_rsrp: string,
  lte_inter_anr_a5_threshold_1_rsrp: string,
  lte_inter_anr_a5_threshold_2_rsrp: string,
  b2_threshold1_rsrp: string,
  b2_threshold2_rsrp: string,
  b2_geran_irat_threshold: string,
  qrxlevmin_sib1: string,
  qrxlevminoffset: string,
  s_intrasearch: string,
  s_nonintrasearch: string,
  qrxlevmin_sib3: string,
  reselection_priority: string,
  threshservinglow: string,
  ciphering_algorithm: string,
  integrity_algorithm: string,
};
type OptKey = $Keys<OptConfig>;

export function RanEdit(props: Props) {
  const {match} = useRouter();
  const ctx = useContext(EnodebContext);
  const enodebSerial: string = match.params.enodebSerial;
  const enbInfo = ctx.state.enbInfo[enodebSerial];

  const handleEnbChange = (key: string, val) =>
    setConfig({...config, [key]: val});

  const handleUnmanagedEnbChange = (key: string, val) =>
    setUnmanagedConfig({...unmanagedConfig, [key]: val});

  const handleOptChange = (key: OptKey, val) =>
    setOptConfig({...optConfig, [(key: string)]: val});

  const [error, setError] = useState('');

  const [
    unmanagedConfig,
    setUnmanagedConfig,
  ] = useState<unmanaged_enodeb_configuration>(
    props.enb?.enodeb_config?.unmanaged_config || {
      cell_id: 0,
      ip_address: '',
      tac: 0,
    },
  );

  const [config, setConfig] = useState<enodeb_configuration>(
    props.enb?.enodeb_config?.managed_config || DEFAULT_ENB_CONFIG.config,
  );

  const [enbConfigType, setEnbConfigType] = useState<'MANAGED' | 'UNMANAGED'>(
    props.enb?.enodeb_config?.config_type ?? 'MANAGED',
  );

  const [optConfig, setOptConfig] = useState<OptConfig>({
    earfcndl: String(config.earfcndl ?? ''),
    bandwidthMhz: config.bandwidth_mhz ?? EnodebBandwidthOption['20'],
    specialSubframePattern: String(config.special_subframe_pattern ?? ''),
    subframeAssignment: String(config.subframe_assignment ?? ''),
    pci: String(config.pci ?? ''),
    tac: String(config.tac ?? ''),
    a1_threshold_rsrp: String(
      config.ho_algorithm_config.a1_threshold_rsrp ?? '',
    ),
    lte_a1_threshold_rsrq: String(
      config.ho_algorithm_config.lte_a1_threshold_rsrq ?? '',
    ),
    hysteresis: String(
      config.ho_algorithm_config.hysteresis ?? '',
    ),
    time_to_trigger: String(
      config.ho_algorithm_config.time_to_trigger ?? '',
    ),
    a2_threshold_rsrp: String(
      config.ho_algorithm_config.a2_threshold_rsrp ?? '',
    ),
    lte_a2_threshold_rsrq: String(
      config.ho_algorithm_config.lte_a2_threshold_rsrq ?? '',
    ),
    lte_a2_threshold_rsrp_irat_volte: String(
      config.ho_algorithm_config.lte_a2_threshold_rsrp_irat_volte ?? '',
    ),
    lte_a2_threshold_rsrq_irat_volte: String(
      config.ho_algorithm_config.lte_a2_threshold_rsrq_irat_volte ?? '',
    ),
    a3_offset: String(
      config.ho_algorithm_config.a3_offset ?? '',
    ),
    a3_offset_anr: String(
      config.ho_algorithm_config.a3_offset_anr ?? '',
    ),
    a4_threshold_rsrp: String(
      config.ho_algorithm_config.a4_threshold_rsrp ?? '',
    ),
    lte_intra_a5_threshold_1_rsrp: String(
      config.ho_algorithm_config.lte_intra_a5_threshold_1_rsrp ?? '',
    ),
    lte_intra_a5_threshold_2_rsrp: String(
      config.ho_algorithm_config.lte_intra_a5_threshold_2_rsrp ?? '',
    ),
      lte_inter_anr_a5_threshold_1_rsrp: String(
      config.ho_algorithm_config.lte_inter_anr_a5_threshold_1_rsrp ?? '',
    ),
      lte_inter_anr_a5_threshold_2_rsrp: String(
      config.ho_algorithm_config.lte_inter_anr_a5_threshold_2_rsrp ?? '',
    ),
    b2_threshold1_rsrp: String(
      config.ho_algorithm_config.b2_threshold1_rsrp ?? '',
    ),
    b2_threshold2_rsrp: String(
      config.ho_algorithm_config.b2_threshold2_rsrp ?? '',
    ),
    b2_geran_irat_threshold: String(
      config.ho_algorithm_config.b2_geran_irat_threshold ?? '',
    ),
    qrxlevmin_sib1: String(
      config.ho_algorithm_config.qrxlevmin_sib1 ?? '',
    ),
    qrxlevminoffset: String(
      config.ho_algorithm_config.qrxlevminoffset ?? '',
    ),
    s_intrasearch: String(
      config.ho_algorithm_config.s_intrasearch ?? '',
    ),
    s_nonintrasearch: String(
      config.ho_algorithm_config.s_nonintrasearch ?? '',
    ),
    qrxlevmin_sib3: String(
      config.ho_algorithm_config.qrxlevmin_sib3 ?? '',
    ),
    reselection_priority: String(
      config.ho_algorithm_config.reselection_priority ?? '',
    ),
    threshservinglow: String(
      config.ho_algorithm_config.threshservinglow ?? '',
    ),
    ciphering_algorithm: String(
      config.ho_algorithm_config.ciphering_algorithm ?? '',
    ),
    integrity_algorithm: String(
      config.ho_algorithm_config.integrity_algorithm ?? '',
    ),
  });

  const enqueueSnackbar = useEnqueueSnackbar();

  const onSave = async () => {
    try {
      const enb: enodeb = {
        ...(props.enb || DEFAULT_ENB_CONFIG),
        config:
          enbConfigType === 'MANAGED'
            ? buildRanConfig(config, optConfig)
            : DEFAULT_ENB_CONFIG.config,
        enodeb_config: {
          config_type: enbConfigType,
          managed_config:
            enbConfigType === 'MANAGED'
              ? buildRanConfig(config, optConfig)
              : undefined,
          unmanaged_config:
            enbConfigType === 'UNMANAGED' ? unmanagedConfig : undefined,
        },
      };

      await ctx.setState(enb.serial, {
        enb_state: enbInfo?.enb_state ?? {},
        enb: enb,
      });

      enqueueSnackbar('eNodeb saved successfully', {
        variant: 'success',
      });
      props.onSave(enb);
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
          <AltFormField label={'eNodeB Managed Externally'}>
            <Switch
              data-testid="enodeb_config"
              onChange={({target}) =>
                setEnbConfigType(target.checked ? 'UNMANAGED' : 'MANAGED')
              }
              checked={enbConfigType === 'UNMANAGED'}
            />
          </AltFormField>
          {enbConfigType === 'UNMANAGED' ? (
            <>
              <AltFormField label={'Cell ID'}>
                <OutlinedInput
                  data-testid="cellId"
                  type="number"
                  min={0}
                  max={Math.pow(2, 28) - 1}
                  fullWidth={true}
                  value={unmanagedConfig.cell_id}
                  onChange={({target}) =>
                    handleUnmanagedEnbChange('cell_id', parseInt(target.value))
                  }
                />
              </AltFormField>
              <AltFormField label={'TAC'}>
                <OutlinedInput
                  data-testid="tac"
                  type="number"
                  min={0}
                  max={65535}
                  fullWidth={true}
                  value={unmanagedConfig.tac}
                  onChange={({target}) =>
                    handleUnmanagedEnbChange('tac', parseInt(target.value))
                  }
                />
              </AltFormField>
              <AltFormField label={'IP Address'}>
                <OutlinedInput
                  data-testid="ipAddress"
                  fullWidth={true}
                  placeholder="192.168.0.1/24"
                  value={unmanagedConfig.ip_address}
                  onChange={({target}) =>
                    handleUnmanagedEnbChange('ip_address', target.value)
                  }
                />
              </AltFormField>
            </>
          ) : (
            <>
              <AltFormField label={'Device Class'}>
                <FormControl>
                  <Select
                    value={config.device_class}
                    onChange={({target}) =>
                      handleEnbChange(
                        'device_class',
                        coerceValue(target.value, EnodebDeviceClass),
                      )
                    }
                    input={<OutlinedInput id="deviceClass" />}>
                    {Object.keys(EnodebDeviceClass).map(
                      (k: string, idx: number) => (
                        <MenuItem key={idx} value={EnodebDeviceClass[k]}>
                          {EnodebDeviceClass[k]}
                        </MenuItem>
                      ),
                    )}
                  </Select>
                </FormControl>
              </AltFormField>
              <AltFormField label={'Cell ID'}>
                <OutlinedInput
                  data-testid="cellId"
                  type="number"
                  min={0}
                  max={Math.pow(2, 28) - 1}
                  fullWidth={true}
                  value={config.cell_id}
                  onChange={({target}) =>
                    handleEnbChange('cell_id', parseInt(target.value))
                  }
                />
              </AltFormField>
              <AltFormField label={'Bandwidth'}>
                <FormControl>
                  <Select
                    value={optConfig.bandwidthMhz}
                    onChange={({target}) =>
                      handleOptChange(
                        'bandwidthMhz',
                        coerceValue(target.value, EnodebBandwidthOption),
                      )
                    }
                    input={<OutlinedInput id="bandwidth" />}>
                    {Object.keys(EnodebBandwidthOption).map(
                      (k: string, idx: number) => (
                        <MenuItem key={idx} value={EnodebBandwidthOption[k]}>
                          {EnodebBandwidthOption[k]}
                        </MenuItem>
                      ),
                    )}
                  </Select>
                </FormControl>
              </AltFormField>
              {props.lteRanConfigs?.tdd_config && (
                <EnodeConfigEditTdd
                  earfcndl={optConfig.earfcndl}
                  specialSubframePattern={optConfig.specialSubframePattern}
                  subframeAssignment={optConfig.subframeAssignment}
                  setEarfcndl={v => handleOptChange('earfcndl', v)}
                  setSubframeAssignment={v =>
                    handleOptChange('subframeAssignment', v)
                  }
                  setSpecialSubframePattern={v =>
                    handleOptChange('specialSubframePattern', v)
                  }
                />
              )}
              {props.lteRanConfigs?.fdd_config && (
                <EnodeConfigEditFdd
                  earfcndl={optConfig.earfcndl}
                  earfcnul={props.lteRanConfigs.fdd_config.earfcnul.toString()}
                  setEarfcndl={v => handleOptChange('earfcndl', v)}
                />
              )}
              <AltFormField label={'PCI'}>
                <OutlinedInput
                  data-testid="pci"
                  placeholder="Enter PCI"
                  fullWidth={true}
                  value={optConfig.pci}
                  onChange={({target}) => handleOptChange('pci', target.value)}
                />
              </AltFormField>

              <AltFormField label={'TAC'}>
                <OutlinedInput
                  data-testid="tac"
                  placeholder="Enter TAC"
                  fullWidth={true}
                  value={optConfig.tac}
                  onChange={({target}) => handleOptChange('tac', target.value)}
                />
              </AltFormField>

              <AltFormField label={'A1 Thraeshold Rsrp'}>
                <OutlinedInput
                  data-testid="a1_threshold_rsrp"
                  placeholder="Enter A1THRESHOLDRSRP"
                  fullWidth={true}
                  value={optConfig.a1_threshold_rsrp}
                  onChange={({target}) =>
                    handleOptChange('a1_threshold_rsrp', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'lte a1 threshold rsrq'}>
                <OutlinedInput
                  data-testid="lte_a1_threshold_rsrq"
                  placeholder="Enter LTEA1THRESHOLDRSRQ"
                  fullWidth={true}
                  value={optConfig.lte_a1_threshold_rsrq}
                  onChange={({target}) =>
                    handleOptChange('lte_a1_threshold_rsrq', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'hysteresis'}>
                <OutlinedInput
                  data-testid="hysteresis"
                  placeholder="Enter HYSTERESIS"
                  fullWidth={true}
                  value={optConfig.hysteresis}
                  onChange={({target}) =>
                    handleOptChange('hysteresis', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'time to trigger'}>
                <OutlinedInput
                  data-testid="time_to_trigger"
                  placeholder="Enter TIMETOTRIGGER"
                  fullWidth={true}
                  value={optConfig.time_to_trigger.toString()}
                  onChange={({target}) =>
                    handleOptChange('time_to_trigger', target.value.toString())
                  }
                />
              </AltFormField>
              <AltFormField label={'a2 threshold rsrp'}>
                <OutlinedInput
                  data-testid="a2_threshold_rsrp"
                  placeholder="Enter A2THRESHOLDRSRP"
                  fullWidth={true}
                  value={optConfig.a2_threshold_rsrp}
                  onChange={({target}) =>
                    handleOptChange('a2_threshold_rsrp', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'lte a2 threshold rsrq'}>
                <OutlinedInput
                  data-testid="lte_a2_threshold_rsrq"
                  placeholder="Enter LTEA2THRESHOLDRSRQ"
                  fullWidth={true}
                  value={optConfig.lte_a2_threshold_rsrq}
                  onChange={({target}) =>
                    handleOptChange('lte_a2_threshold_rsrq', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'lte a2 threshold rsrp irat volte'}>
                <OutlinedInput
                  data-testid="lte_a2_threshold_rsrp_irat_volte"
                  placeholder="Enter LTEA2THRESHOLDRSRPIRATVOLTE"
                  fullWidth={true}
                  value={optConfig.lte_a2_threshold_rsrp_irat_volte}
                  onChange={({target}) =>
                    handleOptChange('lte_a2_threshold_rsrp_irat_volte', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'lte a2 threshold rsrq irat volte'}>
                <OutlinedInput
                  data-testid="lte_a2_threshold_rsrq_irat_volte"
                  placeholder="Enter LTEA2THRESHOLDRSRQIRATVOLTE"
                  fullWidth={true}
                  value={optConfig.lte_a2_threshold_rsrq_irat_volte}
                  onChange={({target}) =>
                    handleOptChange('lte_a2_threshold_rsrq_irat_volte', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'a3 offset'}>
                <OutlinedInput
                  data-testid="a3_offset"
                  placeholder="Enter A3OFFSET"
                  fullWidth={true}
                  value={optConfig.a3_offset}
                  onChange={({target}) =>
                    handleOptChange('a3_offset', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'a3 offset anr'}>
                <OutlinedInput
                  data-testid="a3_offset_anr"
                  placeholder="Enter A3OFFSETANR"
                  fullWidth={true}
                  value={optConfig.a3_offset_anr}
                  onChange={({target}) =>
                    handleOptChange('a3_offset_anr', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'a4 threshold rsrp'}>
                <OutlinedInput
                  data-testid="a4_threshold_rsrp"
                  placeholder="Enter A4THRESHOLDRSRP"
                  fullWidth={true}
                  value={optConfig.a4_threshold_rsrp}
                  onChange={({target}) =>
                    handleOptChange('a4_threshold_rsrp', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'lte intraa5 threshold1 rsrp'}>
                <OutlinedInput
                  data-testid="lte_intra_a5_threshold_1_rsrp"
                  placeholder="Enter LTEINTRAA5THRESHOLD1RSRP"
                  fullWidth={true}
                  value={optConfig.lte_intra_a5_threshold_1_rsrp}
                  onChange={({target}) =>
                    handleOptChange('lte_intra_a5_threshold_1_rsrp', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'lte intraa5 threshold2 rsrp'}>
                <OutlinedInput
                  data-testid="lte_intra_a5_threshold_2_rsrp"
                  placeholder="Enter LTEINTRAA5THRESHOLD2RSRP"
                  fullWidth={true}
                  value={optConfig.lte_intra_a5_threshold_2_rsrp}
                  onChange={({target}) =>
                    handleOptChange('lte_intra_a5_threshold_2_rsrp', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'lte interanra5 threshold1 rsrp'}>
                <OutlinedInput
                  data-testid="lte_inter_anr_a5_threshold_1_rsrp"
                  placeholder="Enter LTEINTERANRA5THRESHOLD1RSRP"
                  fullWidth={true}
                  value={optConfig.lte_inter_anr_a5_threshold_1_rsrp}
                  onChange={({target}) =>
                    handleOptChange('lte_inter_anr_a5_threshold_1_rsrp', target.value)
                }
                />
              </AltFormField>
              <AltFormField label={'lte interanra5 threshold2 rsrp'}>
                <OutlinedInput
                  data-testid="lte_inter_anr_a5_threshold_2_rsrp"
                  placeholder="Enter LTEINTERANRA5THRESHOLD2RSRP"
                  fullWidth={true}
                  value={optConfig.lte_inter_anr_a5_threshold_2_rsrp}
                  onChange={({target}) =>
                    handleOptChange('lte_inter_anr_a5_threshold_2_rsrp', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'b2 threshold1 rsrp'}>
                <OutlinedInput
                  data-testid="b2_threshold1_rsrp"
                  placeholder="Enter B2THRESHOLD1RSRP"
                  fullWidth={true}
                  value={optConfig.b2_threshold1_rsrp}
                  onChange={({target}) =>
                    handleOptChange('b2_threshold1_rsrp', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'b2 threshold2 rsrp'}>
                <OutlinedInput
                  data-testid="b2_threshold2_rsrp"
                  placeholder="Enter B2THRESHOLD2RSRP"
                  fullWidth={true}
                  value={optConfig.b2_threshold2_rsrp}
                  onChange={({target}) =>
                    handleOptChange('b2_threshold2_rsrp', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'b2 geran irat threshold'}>
                <OutlinedInput
                  data-testid="b2_geran_irat_threshold"
                  placeholder="Enter B2GERANIRATTHRESHOLD"
                  fullWidth={true}
                  value={optConfig.b2_geran_irat_threshold}
                  onChange={({target}) =>
                    handleOptChange('b2_geran_irat_threshold', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'qrxlevmin sib1'}>
                <OutlinedInput
                  data-testid="qrxlevmin_sib1"
                  placeholder="Enter QRXLEVMINSIB1"
                  fullWidth={true}
                  value={optConfig.qrxlevmin_sib1}
                  onChange={({target}) =>
                    handleOptChange('qrxlevmin_sib1', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'qrxlevminoffset'}>
                <OutlinedInput
                  data-testid="qrxlevminoffset"
                  placeholder="Enter QRXLEVMINOFFSET"
                  fullWidth={true}
                  value={optConfig.qrxlevminoffset}
                  onChange={({target}) =>
                    handleOptChange('qrxlevminoffset', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'s intrasearch'}>
                <OutlinedInput
                  data-testid="s_intrasearch"
                  placeholder="Enter SINTRASEARCH"
                  fullWidth={true}
                  value={optConfig.s_intrasearch}
                  onChange={({target}) =>
                    handleOptChange('s_intrasearch', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'s nonintrasearch'}>
                <OutlinedInput
                  data-testid="s_nonintrasearch"
                  placeholder="Enter SNONINTRASEARCH"
                  fullWidth={true}
                  value={optConfig.s_nonintrasearch}
                  onChange={({target}) =>
                    handleOptChange('s_nonintrasearch', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'qrxlevmin sib3'}>
                <OutlinedInput
                  data-testid="qrxlevmin_sib3"
                  placeholder="Enter QRXLEVMINSIB3"
                  fullWidth={true}
                  value={optConfig.qrxlevmin_sib3}
                  onChange={({target}) =>
                    handleOptChange('qrxlevmin_sib3', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'reselection priority'}>
                <OutlinedInput
                  data-testid="reselection_priority"
                  placeholder="Enter RESELECTIONPRIORITY"
                  fullWidth={true}
                  value={optConfig.reselection_priority}
                  onChange={({target}) =>
                    handleOptChange('reselection_priority', target.value)
                  }
                />
              </AltFormField>
              <AltFormField label={'threshservinglow'}>
                <OutlinedInput
                  data-testid="threshservinglow"
                  placeholder="Enter THRESHSERVINGLOW"
                  fullWidth={true}
                  value={optConfig.threshservinglow}
                  onChange={({target}) =>
                    handleOptChange('threshservinglow', target.value)
                  }
                />
              </AltFormField>

              <AltFormField label={'ciphering algorithm'}>
                <OutlinedInput
                  data-testid="ciphering_algorithm"
                  placeholder="Enter CIPHERINGALGORITHM"
                  fullWidth={true}
                  value={optConfig.ciphering_algorithm.toString()}
                  onChange={({target}) =>
                    handleOptChange('ciphering_algorithm', target.value.toString())
                  }
                />
              </AltFormField>
              <AltFormField label={'integrity algorithm'}>
                <OutlinedInput
                  data-testid="integrity_algorithm"
                  placeholder="Enter INTEGRITYALGORITHM"
                  fullWidth={true}
                  value={optConfig.integrity_algorithm.toString()}
                  onChange={({target}) =>
                    handleOptChange('integrity_algorithm', target.value.toString())
                  }
                />
              </AltFormField>

              <AltFormField label={'Transmit'}>
                <FormControl variant={'outlined'}>
                  <Select
                    value={config.transmit_enabled ? 1 : 0}
                    onChange={({target}) =>
                      handleEnbChange('transmit_enabled', target.value === 1)
                    }
                    input={<OutlinedInput id="transmitEnabled" />}>
                    <MenuItem value={0}>Disabled</MenuItem>
                    <MenuItem value={1}>Enabled</MenuItem>
                  </Select>
                </FormControl>
              </AltFormField>
            </>
          )}
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.isAdd ? 'Save And Add eNodeB' : 'Save'}
        </Button>
      </DialogActions>
    </>
  );
}

export function ConfigEdit(props: Props) {
  const [error, setError] = useState('');
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(EnodebContext);
  const enodebSerial: string = match.params.enodebSerial;
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const [enb, setEnb] = useState<enodeb>(props.enb || DEFAULT_ENB_CONFIG);
  const onSave = async () => {
    try {
      if (props.isAdd) {
        // check if it is not a modify during add i.e we aren't switching tabs back
        // during add and modifying the information other than the serial number
        if (
          enb.serial in ctx.state.enbInfo &&
          enb.serial !== props.enb?.serial
        ) {
          setError(`eNodeB ${enb.serial} already exists`);
          return;
        }
      }

      if (enb.config == null) {
        enb.config = DEFAULT_ENB_CONFIG.config;
      }
      if (enb.enodeb_config == null || enb.enodeb_config.config_type == '') {
        enb.enodeb_config = {
          config_type: 'MANAGED',
          managed_config: DEFAULT_ENB_CONFIG.config,
        };
      }
      await ctx.setState(enb.serial, {
        enb_state: enbInfo?.enb_state ?? {},
        enb: enb,
      });
      if (props.enb) {
        enqueueSnackbar('eNodeb saved successfully', {
          variant: 'success',
        });
      }
      props.onSave(enb);
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
          <AltFormField label={'Name'}>
            <OutlinedInput
              data-testid="name"
              placeholder="Enter Name"
              fullWidth={true}
              value={enb.name}
              onChange={({target}) => setEnb({...enb, name: target.value})}
            />
          </AltFormField>
          <AltFormField label={'Serial Number'}>
            <OutlinedInput
              data-testid="serial"
              placeholder="Ex: 12020000261814C0021"
              fullWidth={true}
              value={enb.serial}
              readOnly={props.enb ? true : false}
              onChange={({target}) => setEnb({...enb, serial: target.value})}
            />
          </AltFormField>
          <AltFormField label={'Description'}>
            <OutlinedInput
              data-testid="description"
              placeholder="Enter Description"
              fullWidth={true}
              multiline
              rows={4}
              value={enb.description}
              onChange={({target}) =>
                setEnb({...enb, description: target.value})
              }
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.isAdd ? 'Save And Continue' : 'Save'}
        </Button>
      </DialogActions>
    </>
  );
}

function coerceValue<T>(value: string, options: {[string]: T}): T {
  const values = Object.values(options);
  const keys = Object.keys(options);
  const optionKey = values.indexOf(value);
  if (optionKey > -1) {
    return options[keys[optionKey]];
  } else {
    throw Error('Expected a valid selection.');
  }
}

function isNumberInRange(value: string | number, lower: number, upper: number) {
  const val = parseInt(value, 10);
  if (isNaN(val)) {
    return false;
  }
  return val >= lower && val <= upper;
}

function buildRanConfig(config: enodeb_configuration, optConfig: OptConfig) {
  const response = {
    ...config,
    bandwidth_mhz: optConfig.bandwidthMhz,
    ho_algorithm_config: {
      a1_threshold_rsrp: null,
      lte_a1_threshold_rsrq: null,
      hysteresis: null,
      time_to_trigger: null,
      a2_threshold_rsrp: null,
      lte_a2_threshold_rsrq: null,
      lte_a2_threshold_rsrp_irat_volte: null,
      lte_a2_threshold_rsrq_irat_volte: null,
      a3_offset: null,
      a3_offset_anr: null,
      a4_threshold_rsrp: null,
      lte_intra_a5_threshold_1_rsrp: null,
      lte_intra_a5_threshold_2_rsrp: null,
      lte_inter_anr_a5_threshold_1_rsrp: null,
      lte_inter_anr_a5_threshold_2_rsrp: null,
      b2_threshold1_rsrp: null,
      b2_threshold2_rsrp: null,
      b2_geran_irat_threshold: null,
      qrxlevmin_sib1: null,
      qrxlevminoffset: null,
      s_intrasearch: null,
      s_nonintrasearch: null,
      qrxlevmin_sib3: null,
      reselection_priority: null,
      threshservinglow: null,
      ciphering_algorithm: null,
      integrity_algorithm: null,
    },
  };

  if (!isNumberInRange(config.cell_id, 0, Math.pow(2, 28) - 1)) {
    throw Error('Invalid Configuration Cell ID. Valid range 0 - (2^28) - 1');
  }
  if (optConfig.earfcndl !== '') {
    if (!isNumberInRange(optConfig.earfcndl, 0, 65535)) {
      throw Error('Invalid EARFCNDL. Valid range 0 - 645535');
    }
    response['earfcndl'] = parseInt(optConfig.earfcndl);
  }

  if (optConfig.pci !== '') {
    if (!isNumberInRange(optConfig.pci, 0, 504)) {
      throw Error('Invalid PCI. Valid range 0 - 504');
    }
    response['pci'] = parseInt(optConfig.pci);
  }

  if (optConfig.specialSubframePattern !== '') {
    if (!isNumberInRange(optConfig.specialSubframePattern, 0, 9)) {
      throw Error('Invalid Special SubFrame Pattern, Valid range 0 - 9');
    }
    response['special_subframe_pattern'] = parseInt(
      optConfig.specialSubframePattern,
    );
  }

  if (optConfig.subframeAssignment !== '') {
    if (!isNumberInRange(optConfig.subframeAssignment, 0, 6)) {
      throw Error('Invalid SubFrame Assignment, Valid range 0 - 6');
    }
    response['subframe_assignment'] = parseInt(optConfig.subframeAssignment);
  }

  if (optConfig.tac !== '') {
    if (!isNumberInRange(optConfig.tac, 0, 65535)) {
      throw Error('Invalid TAC, Valid Range 0 - 65535');
    }
    response['tac'] = parseInt(optConfig.tac);
  }
  if (optConfig.a1_threshold_rsrp !== '') {
    response.ho_algorithm_config.a1_threshold_rsrp = parseInt(
      optConfig.a1_threshold_rsrp,
    );
  }
  if (optConfig.lte_a1_threshold_rsrq !== '') {
    response.ho_algorithm_config.lte_a1_threshold_rsrq = parseInt(
      optConfig.lte_a1_threshold_rsrq,
    );
  }
  if (optConfig.hysteresis !== '') {
    response.ho_algorithm_config.hysteresis = parseInt(
      optConfig.hysteresis,
    );
  }
  if (optConfig.time_to_trigger !== '') {
    response.ho_algorithm_config.time_to_trigger = optConfig.time_to_trigger;
  }
  if (optConfig.a2_threshold_rsrp !== '') {
    response.ho_algorithm_config.a2_threshold_rsrp = parseInt(
      optConfig.a2_threshold_rsrp,
    );
  }
  if (optConfig.lte_a2_threshold_rsrq !== '') {
    response.ho_algorithm_config.lte_a2_threshold_rsrq = parseInt(
      optConfig.lte_a2_threshold_rsrq,
    );
  }
  if (optConfig.lte_a2_threshold_rsrq !== '') {
    response.ho_algorithm_config.lte_a2_threshold_rsrq = parseInt(
      optConfig.lte_a2_threshold_rsrq,
    );
  }
  if (optConfig.lte_a2_threshold_rsrp_irat_volte !== '') {
    response.ho_algorithm_config.lte_a2_threshold_rsrp_irat_volte = parseInt(
      optConfig.lte_a2_threshold_rsrp_irat_volte,
    );
  }if (optConfig.lte_a2_threshold_rsrq_irat_volte !== '') {
    response.ho_algorithm_config.lte_a2_threshold_rsrq_irat_volte = parseInt(
      optConfig.lte_a2_threshold_rsrq_irat_volte,
    );
  }
  if (optConfig.a3_offset !== '') {
    response.ho_algorithm_config.a3_offset = parseInt(
      optConfig.a3_offset,
    );
  }
  if (optConfig.a3_offset_anr !== '') {
    response.ho_algorithm_config.a3_offset_anr = parseInt(
      optConfig.a3_offset_anr,
    );
  }if (optConfig.a4_threshold_rsrp !== '') {
    response.ho_algorithm_config.a4_threshold_rsrp = parseInt(
      optConfig.a4_threshold_rsrp,
    );
  }
  if (optConfig.lte_intra_a5_threshold_1_rsrp !== '') {
    response.ho_algorithm_config.lte_intra_a5_threshold_1_rsrp = parseInt(
      optConfig.lte_intra_a5_threshold_1_rsrp,
    );
  }
  if (optConfig.lte_intra_a5_threshold_2_rsrp !== '') {
    response.ho_algorithm_config.lte_intra_a5_threshold_2_rsrp = parseInt(
      optConfig.lte_intra_a5_threshold_2_rsrp,
    );
  }
  if (optConfig.lte_inter_anr_a5_threshold_1_rsrp !== '') {
    response.ho_algorithm_config.lte_inter_anr_a5_threshold_1_rsrp = parseInt(
      optConfig.lte_inter_anr_a5_threshold_1_rsrp,
    );
  }
  if (optConfig.lte_inter_anr_a5_threshold_2_rsrp !== '') {
    response.ho_algorithm_config.lte_inter_anr_a5_threshold_2_rsrp = parseInt(
      optConfig.lte_inter_anr_a5_threshold_2_rsrp,
    );
  }
  if (optConfig.b2_threshold1_rsrp !== '') {
    response.ho_algorithm_config.b2_threshold1_rsrp = parseInt(
      optConfig.b2_threshold1_rsrp,
    );
  }
  if (optConfig.b2_threshold2_rsrp !== '') {
    response.ho_algorithm_config.b2_threshold2_rsrp = parseInt(
      optConfig.b2_threshold2_rsrp,
    );
  }
  if (optConfig.b2_geran_irat_threshold !== '') {
    response.ho_algorithm_config.b2_geran_irat_threshold = parseInt(
      optConfig.b2_geran_irat_threshold,
    );
  }
  if (optConfig.qrxlevmin_sib1 !== '') {
    response.ho_algorithm_config.qrxlevmin_sib1 = parseInt(
      optConfig.qrxlevmin_sib1,
    );
  }
  if (optConfig.qrxlevminoffset !== '') {
    response.ho_algorithm_config.qrxlevminoffset = parseInt(
      optConfig.qrxlevminoffset,
    );
  }
  if (optConfig.s_intrasearch !== '') {
    response.ho_algorithm_config.s_intrasearch = parseInt(
      optConfig.s_intrasearch,
    );
  }
  if (optConfig.s_nonintrasearch !== '') {
    response.ho_algorithm_config.s_nonintrasearch = parseInt(
      optConfig.s_nonintrasearch,
    );
  }
  if (optConfig.qrxlevmin_sib3 !== '') {
    response.ho_algorithm_config.qrxlevmin_sib3 = parseInt(
      optConfig.qrxlevmin_sib3,
    );
  }
  if (optConfig.reselection_priority !== '') {
    response.ho_algorithm_config.reselection_priority = parseInt(
      optConfig.reselection_priority,
    );
  }
  /*if (optConfig.a1_threshold_rsrp !== '') {
    response.ho_algorithm_config.a1_threshold_rsrp = parseInt(
      optConfig.a1_threshold_rsrp,
    );
  }*/
  if (optConfig.threshservinglow !== '') {
    response.ho_algorithm_config.threshservinglow = parseInt(
      optConfig.threshservinglow,
    );
  }
  if (optConfig.ciphering_algorithm !== '') {
    response.ho_algorithm_config.ciphering_algorithm = optConfig.ciphering_algorithm;
  }
  if (optConfig.integrity_algorithm !== '') {
    response.ho_algorithm_config.integrity_algorithm = optConfig.integrity_algorithm;
  }




  return response;
}
