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

import type {enodeb, enodeb_configuration} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import EnodebPropertySelector from './EnodebPropertySelector';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React, {useState} from 'react';
import Switch from '@material-ui/core/Switch';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import {EnodebBandwidthOption, EnodebDeviceClass} from './EnodebUtils';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type BandwidthMhzType = $PropertyType<enodeb_configuration, 'bandwidth_mhz'>;
type DeviceClassType = $PropertyType<enodeb_configuration, 'device_class'>;

type Props = {
  // Only set if we are editing an eNodeB configuration
  editingEnodeb: ?enodeb,
  onClose: () => void,
  onSave: enodeb => void,
};

export default function AddEditEnodebDialog(props: Props) {
  const {editingEnodeb} = props;
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const networkId = nullthrows(match.params.networkId);

  let defaultEnodebId = '0';
  let defaultCellNumber = '1';
  if (editingEnodeb) {
    [defaultEnodebId, defaultCellNumber] = unpackCellId(
      editingEnodeb.config.cell_id,
    );
  }

  const [deviceClass, setDeviceClass] = useState<DeviceClassType>(
    editingEnodeb?.config.device_class ?? EnodebDeviceClass['BAICELLS_ID'],
  );
  const [bandwidthMhz, setBandwidthMhz] = useState<BandwidthMhzType>(
    editingEnodeb?.config.bandwidth_mhz ?? EnodebBandwidthOption['20'],
  );
  const [specialSubframePattern, setSpecialSubframePattern] = useState<string>(
    String(editingEnodeb?.config.special_subframe_pattern || ''),
  );
  const [subframeAssignment, setSubframeAssignment] = useState<string>(
    String(editingEnodeb?.config.subframe_assignment || ''),
  );
  const [serialId, setSerialId] = useState<string>(editingEnodeb?.serial || '');
  const [tac, setTac] = useState(String(editingEnodeb?.config.tac || ''));
  const [transmitEnabled, setTransmitEnabled] = useState<boolean>(
    editingEnodeb?.config.transmit_enabled ?? false,
  );
  const [earfcndl, setEarfcndl] = useState<string>(
    String(editingEnodeb?.config.earfcndl || ''),
  );
  const [pci, setPci] = useState<string>(
    String(editingEnodeb?.config.pci || ''),
  );
  const [enodebId, setEnodebId] = useState<string>(defaultEnodebId);
  const [cellNumber, setCellNumber] = useState<string>(defaultCellNumber);
  const [name, setName] = useState<string>(editingEnodeb?.name || '');

  const cellId = packCellId(enodebId, cellNumber);

  const isBandwidthMhzValid = isNumberInRange(bandwidthMhz || 0, 0, 20);
  const isCellNumberValid = isNumberInRange(cellNumber, 0, Math.pow(2, 8) - 1);
  const isCellIdValid = isNumberInRange(cellId, 0, Math.pow(2, 28) - 1);
  const isEarfcndlValid =
    earfcndl === '' || isNumberInRange(earfcndl, 0, 65535);
  const isEnodebIdValid = isNumberInRange(enodebId, 0, Math.pow(2, 20) - 1);
  const isNameValid = name !== '';
  const isPciValid = pci === '' || isNumberInRange(pci, 0, 504);
  const isSerialIdValid = serialId.length > 0;
  const isSpecialSubframePatternValid =
    specialSubframePattern === '' ||
    isNumberInRange(specialSubframePattern, 0, 9);
  const isSubframeAssignmentValid =
    subframeAssignment === '' || isNumberInRange(subframeAssignment, 0, 6);
  const isTacValid = tac === '' || isNumberInRange(tac, 0, 65535);
  const isTransmitEnabledValid = typeof transmitEnabled === 'boolean';

  const isFormValid =
    isNameValid &&
    isSerialIdValid &&
    isEarfcndlValid &&
    isSubframeAssignmentValid &&
    isSpecialSubframePatternValid &&
    isPciValid &&
    isBandwidthMhzValid &&
    isTacValid &&
    isEnodebIdValid &&
    isCellNumberValid &&
    isCellIdValid &&
    isTransmitEnabledValid;

  const onSave = async () => {
    if (!isFormValid) {
      enqueueSnackbar('Please complete all fields with valid values', {
        variant: 'error',
      });
      return;
    }
    const enb: enodeb = {
      name: name,
      serial: serialId,
      config: {
        device_class: deviceClass,
        bandwidth_mhz: bandwidthMhz,
        cell_id: cellId,
        transmit_enabled: transmitEnabled,
      },
    };

    if (earfcndl !== '') {
      enb.config.earfcndl = parseInt(earfcndl);
    }
    if (subframeAssignment !== '') {
      enb.config.subframe_assignment = parseInt(subframeAssignment);
    }
    if (specialSubframePattern !== '') {
      enb.config.special_subframe_pattern = parseInt(specialSubframePattern);
    }
    if (pci !== '') {
      enb.config.pci = parseInt(pci);
    }
    if (tac !== '') {
      enb.config.tac = parseInt(tac);
    }

    try {
      if (props.editingEnodeb != null) {
        await MagmaV1API.putLteByNetworkIdEnodebsByEnodebSerial({
          networkId,
          enodebSerial: enb.serial,
          enodeb: enb,
        });
      } else {
        await MagmaV1API.postLteByNetworkIdEnodebs({
          networkId,
          enodeb: enb,
        });
      }
      const data = await MagmaV1API.getLteByNetworkIdEnodebsByEnodebSerial({
        networkId,
        enodebSerial: enb.serial,
      });
      props.onSave(data);
    } catch (e) {
      enqueueSnackbar(e?.response?.data?.message || e?.message || e, {
        variant: 'error',
      });
    }
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>{editingEnodeb ? 'Edit eNodeB' : 'Add eNodeB'}</DialogTitle>
      <DialogContent>
        <TextField
          label="eNodeB Name"
          className={classes.input}
          value={name}
          onChange={({target}) => setName(target.value)}
          placeholder="Name of eNodeB, eg. 'Market NW Corner'"
          error={!isNameValid}
        />
        <TextField
          label="eNodeB Serial ID"
          className={classes.input}
          disabled={props.editingEnodeb != null}
          value={serialId}
          onChange={({target}) => setSerialId(target.value)}
          placeholder="Unique Serial ID of eNodeB, eg. 120200002618AGP0003"
          error={!isSerialIdValid}
        />
        <EnodebPropertySelector
          titleLabel="eNodeB Device Class"
          className={classes.input}
          value={deviceClass}
          valueOptionsByKey={EnodebDeviceClass}
          onChange={({target}) =>
            setDeviceClass(coerceValue(target.value, EnodebDeviceClass))
          }
        />
        <TextField
          label="EARFCNDL"
          className={classes.input}
          value={earfcndl}
          onChange={({target}) => setEarfcndl(target.value)}
          placeholder="0-65535"
          error={!isEarfcndlValid}
        />
        <TextField
          label="Subframe Assignment"
          className={classes.input}
          value={subframeAssignment}
          onChange={({target}) => setSubframeAssignment(target.value)}
          placeholder="0-6"
          error={!isSubframeAssignmentValid}
        />
        <TextField
          label="Special Subframe Pattern"
          className={classes.input}
          value={specialSubframePattern}
          onChange={({target}) => setSpecialSubframePattern(target.value)}
          inputProps={{min: 0, max: 9}}
          placeholder="0-9"
          error={!isSpecialSubframePatternValid}
        />
        <TextField
          label="Physical Cell Identifier"
          className={classes.input}
          value={pci}
          onChange={({target}) => setPci(target.value)}
          placeholder="0-504"
          error={!isPciValid}
        />
        <EnodebPropertySelector
          titleLabel="eNodeB DL/UL Bandwidth (MHz)"
          value={bandwidthMhz || ''}
          valueOptionsByKey={EnodebBandwidthOption}
          onChange={({target}) =>
            setBandwidthMhz(coerceValue(target.value, EnodebBandwidthOption))
          }
          className={classes.input}
          error={!isBandwidthMhzValid}
        />
        <TextField
          label="Tracking Area Code"
          className={classes.input}
          value={tac}
          onChange={({target}) => setTac(target.value)}
          placeholder="0-65535"
          error={!isTacValid}
        />
        <TextField
          label="Enodeb ID"
          className={classes.input}
          value={enodebId}
          onChange={({target}) => setEnodebId(target.value)}
          placeholder="0-1048576"
          error={!isEnodebIdValid || !isCellIdValid}
        />
        <TextField
          label="Cell Number"
          className={classes.input}
          value={cellNumber}
          onChange={({target}) => setCellNumber(target.value)}
          error={!isCellNumberValid}
        />
        <FormControl className={classes.input}>
          <FormControlLabel
            control={
              <Switch
                checked={transmitEnabled}
                onChange={() => setTransmitEnabled(!transmitEnabled)}
                color="primary"
              />
            }
            label="Transmit Enabled"
          />
        </FormControl>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button disabled={!isFormValid} onClick={onSave}>
          Save
        </Button>
      </DialogActions>
    </Dialog>
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

function unpackCellId(cellId: number) {
  const cellIdBits = cellId.toString(2).padStart(28, '0');
  return [
    parseInt(cellIdBits.substring(0, 20), 2).toString(),
    parseInt(cellIdBits.substring(20, 28), 2).toString(),
  ];
}

function packCellId(enodebId: string, cellNumber: string) {
  return 256 * parseInt(enodebId) + parseInt(cellNumber);
}
