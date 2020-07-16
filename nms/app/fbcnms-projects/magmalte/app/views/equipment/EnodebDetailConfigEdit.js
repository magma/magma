/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {
  enodeb,
  enodeb_configuration,
  network_ran_configs,
} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import Link from '@material-ui/core/Link';
import List from '@material-ui/core/List';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Select from '@material-ui/core/Select';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import {
  EnodebBandwidthOption,
  EnodebDeviceClass,
} from '../../components/lte/EnodebUtils';

import EnodeConfigEditFdd from './EnodebDetailConfigFdd';
import EnodeConfigEditTdd from './EnodebDetailConfigTdd';

import {AltFormField} from '../../components/FormField';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
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
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
    color: colors.primary.white,
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
    width: '80%',
  },
}));

const EditTableType = {
  config: 0,
  ran: 1,
};

type EditProps = {
  editTable: $Keys<typeof EditTableType>,
  enb: enodeb,
  lteRanConfigs?: ?network_ran_configs,
  onSave: enodeb => void,
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
  const [open, setOpen] = React.useState(false);

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
        <Link
          data-testid={(props.editProps?.editTable ?? '') + 'EditButton'}
          component="button"
          variant="body2"
          onClick={handleClickOpen}>
          {props.title}
        </Link>
      ) : (
        <Button
          variant="contained"
          className={classes.appBarBtn}
          onClick={handleClickOpen}>
          {props.title}
        </Button>
      )}
    </>
  );
}

function EnodeEditDialog({open, onClose, editProps}: DialogProps) {
  const classes = useStyles();
  const [enb, setEnb] = useState<enodeb>({});
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const [lteRanConfigs, setLteRanConfigs] = useState<network_ran_configs>(
    editProps?.lteRanConfigs ?? {},
  );
  const [isLoading, setLoading] = useState<boolean>(true);

  const [tabPos, setTabPos] = useState(
    editProps ? EditTableType[editProps.editTable] : 0,
  );

  useEffect(() => {
    const fetchSt = async () => {
      try {
        const lteRanConfigs = await MagmaV1API.getLteByNetworkIdCellularRan({
          networkId: networkId,
        });
        setLteRanConfigs(lteRanConfigs);
        setLoading(false);
      } catch (error) {
        setLoading(false);
      }
    };
    if (Object.keys(editProps?.lteRanConfigs ?? {}).length > 0) {
      setLoading(false);
      return;
    }
    fetchSt();
  }, [networkId, editProps]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <Dialog data-testid="editDialog" open={open} fullWidth={true} maxWidth="sm">
      <DialogTitle className={classes.topBar}>
        <Text color="light" weight="medium">
          {editProps ? 'Edit eNodeB' : 'Add New eNodeB'}
        </Text>
      </DialogTitle>
      <Tabs
        value={tabPos}
        onChange={(_, v) => setTabPos(v)}
        indicatorColor="primary"
        textColor="primary"
        variant="fullWidth">
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
          saveButtonTitle={editProps ? 'Save' : 'Save And Continue'}
          enb={Object.keys(enb).length != 0 ? enb : editProps?.enb}
          lteRanConfigs={lteRanConfigs}
          onClose={onClose}
          onSave={(enb: enodeb) => {
            setEnb(enb);
            if (editProps) {
              editProps.onSave(enb);
              onClose();
            } else {
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === 1 && (
        <RanEdit
          saveButtonTitle={editProps ? 'Save' : 'Save And Add eNodeB'}
          enb={Object.keys(enb).length != 0 ? enb : editProps?.enb}
          lteRanConfigs={lteRanConfigs}
          onClose={onClose}
          onSave={(enb: enodeb) => {
            setEnb(enb);
            if (editProps) {
              editProps.onSave(enb);
            }
            onClose();
          }}
        />
      )}
    </Dialog>
  );
}

type Props = {
  saveButtonTitle: string,
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
};
type OptKey = $Keys<OptConfig>;

export function RanEdit(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);

  const handleEnbChange = (key: string, val) =>
    setConfig({...config, [key]: val});

  const handleOptChange = (key: OptKey, val) =>
    setOptConfig({...optConfig, [(key: string)]: val});

  const [error, setError] = useState('');
  const [config, setConfig] = useState<enodeb_configuration>(
    props.enb?.config || DEFAULT_ENB_CONFIG.config,
  );

  const [optConfig, setOptConfig] = useState<OptConfig>({
    earfcndl: String(config.earfcndl ?? ''),
    bandwidthMhz: config.bandwidth_mhz ?? EnodebBandwidthOption['20'],
    specialSubframePattern: String(config.special_subframe_pattern ?? ''),
    subframeAssignment: String(config.subframe_assignment ?? ''),
    pci: String(config.pci ?? ''),
    tac: String(config.tac ?? ''),
  });

  const enqueueSnackbar = useEnqueueSnackbar();

  const onSave = async () => {
    try {
      const enb = {
        ...(props.enb || DEFAULT_ENB_CONFIG),
        config: buildRanConfig(config, optConfig),
      };
      await MagmaV1API.putLteByNetworkIdEnodebsByEnodebSerial({
        networkId: networkId,
        enodebSerial: props.enb?.serial ?? '',
        enodeb: enb,
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
        {error !== '' && <FormLabel error>{error}</FormLabel>}
        <List component={Paper}>
          <AltFormField label={'Device Class'}>
            <FormControl className={classes.input}>
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
              className={classes.input}
              fullWidth={true}
              value={config.cell_id}
              onChange={({target}) =>
                handleEnbChange('cell_id', parseInt(target.value))
              }
            />
          </AltFormField>
          <AltFormField label={'Bandwidth'}>
            <FormControl className={classes.input}>
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
              className={classes.input}
              fullWidth={true}
              value={optConfig.pci}
              onChange={({target}) => handleOptChange('pci', target.value)}
            />
          </AltFormField>

          <AltFormField label={'TAC'}>
            <OutlinedInput
              data-testid="tac"
              className={classes.input}
              fullWidth={true}
              value={optConfig.tac}
              onChange={({target}) => handleOptChange('tac', target.value)}
            />
          </AltFormField>

          <AltFormField label={'Transmit'}>
            <FormControl variant={'outlined'} className={classes.input}>
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
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>{props.saveButtonTitle}</Button>
      </DialogActions>
    </>
  );
}

export function ConfigEdit(props: Props) {
  const classes = useStyles();
  const [error, setError] = useState('');
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();

  const [enb, setEnb] = useState<enodeb>(props.enb || DEFAULT_ENB_CONFIG);

  const onSave = async () => {
    try {
      if (props.enb) {
        await MagmaV1API.putLteByNetworkIdEnodebsByEnodebSerial({
          networkId: networkId,
          enodebSerial: enb.serial,
          enodeb: enb,
        });
        enqueueSnackbar('eNodeb saved successfully', {
          variant: 'success',
        });
      } else {
        await MagmaV1API.postLteByNetworkIdEnodebs({
          networkId: networkId,
          enodeb: enb,
        });
        enqueueSnackbar('eNodeb added successfully', {
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
        {error !== '' && <FormLabel error>{error}</FormLabel>}
        <List component={Paper}>
          <AltFormField label={'Name'}>
            <OutlinedInput
              data-testid="name"
              fullWidth={true}
              className={classes.input}
              value={enb.name}
              onChange={({target}) => setEnb({...enb, name: target.value})}
            />
          </AltFormField>
          <AltFormField label={'Serial Number'}>
            <OutlinedInput
              data-testid="serial"
              fullWidth={true}
              className={classes.input}
              value={enb.serial}
              readOnly={props.enb ? true : false}
              onChange={({target}) => setEnb({...enb, serial: target.value})}
            />
          </AltFormField>
          <AltFormField label={'Description'}>
            <OutlinedInput
              data-testid="description"
              fullWidth={true}
              multiline
              className={classes.input}
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
        <Button onClick={onSave}>{props.saveButtonTitle}</Button>
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
  const response = {...config, bandwidth_mhz: optConfig.bandwidthMhz};

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

  return response;
}
