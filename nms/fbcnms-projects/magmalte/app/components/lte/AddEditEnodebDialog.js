/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {WithStyles} from '@material-ui/core';
import type {
  enodeb,
  enodeb_configuration,
} from '../../common/__generated__/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import EnodebPropertySelector from './EnodebPropertySelector';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormLabel from '@material-ui/core/FormLabel';
import MagmaV1API from '../../common/MagmaV1API';
import React from 'react';
import Switch from '@material-ui/core/Switch';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import {EnodebBandwidthOption, EnodebDeviceClass} from './EnodebUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
};

type Props = ContextRouter &
  WithStyles<typeof styles> & {
    // Only set if we are editing an eNodeB configuration
    editingEnodeb: ?enodeb,
    onClose: () => void,
    onSave: enodeb => void,
  };

type EditingEnodeb = {
  name: string,
  serialId: string,
  deviceClass: $PropertyType<enodeb_configuration, 'device_class'>,
  earfcndl: string,
  subframeAssignment: string,
  specialSubframePattern: string,
  pci: string,
  bandwidthMhz: $PropertyType<enodeb_configuration, 'bandwidth_mhz'>,
  tac: string,
  enodebId: string, // 20 bit number, user chosen
  cellNumber: string, // 8-bit number, static
  transmitEnabled: boolean,
};

type State = {
  error: string,
  editingEnodeb: EditingEnodeb,
};

class AddEditEnodebDialog extends React.Component<Props, State> {
  state = {
    error: '',
    editingEnodeb: this.getEditingEnodeb(),
  };

  getEditingEnodeb(): EditingEnodeb {
    const editingEnodeb = this.props.editingEnodeb;
    if (editingEnodeb == null) {
      return {
        name: '',
        serialId: '',
        deviceClass: EnodebDeviceClass['BAICELLS_ID'],
        earfcndl: '0',
        subframeAssignment: '0',
        specialSubframePattern: '0',
        pci: '0',
        bandwidthMhz: EnodebBandwidthOption['20'],
        tac: '0',
        enodebId: '0',
        cellNumber: '1',
        transmitEnabled: false,
      };
    }

    const cellIdBits = editingEnodeb.config.cell_id
      .toString(2)
      .padStart(28, '0');
    const enodebId = parseInt(cellIdBits.substring(0, 20), 2).toString();
    const cellNumber = parseInt(cellIdBits.substring(20, 28), 2).toString();
    return {
      name: editingEnodeb.name,
      serialId: editingEnodeb.serial,
      deviceClass: editingEnodeb.config.device_class,
      earfcndl: String(editingEnodeb.config.earfcndl || 0),
      subframeAssignment: String(editingEnodeb.config.subframe_assignment || 0),
      specialSubframePattern: String(
        editingEnodeb.config.special_subframe_pattern || 0,
      ),
      pci: String(editingEnodeb.config.pci || 0),
      bandwidthMhz: editingEnodeb.config.bandwidth_mhz,
      tac: String(editingEnodeb.config.tac || 0),
      enodebId: enodebId,
      cellNumber: cellNumber,
      transmitEnabled: editingEnodeb.config.transmit_enabled,
    };
  }

  render() {
    const {classes} = this.props;
    const error = this.state.error ? (
      <FormLabel error>{this.state.error}</FormLabel>
    ) : null;

    return (
      <Dialog open={true} onClose={this.props.onClose}>
        <DialogTitle>
          {this.props.editingEnodeb ? 'Edit eNodeB' : 'Add eNodeB'}
        </DialogTitle>
        <DialogContent>
          {error}
          <TextField
            label="eNodeB Name"
            className={classes.input}
            value={this.state.editingEnodeb.name}
            onChange={this.onNameChange}
            placeholder="Name of eNodeB, eg. 'Market NW Corner'"
          />
          <TextField
            label="eNodeB Serial ID"
            className={classes.input}
            disabled={this.props.editingEnodeb != null}
            value={this.state.editingEnodeb.serialId}
            onChange={this.onSerialIdChange}
            placeholder="Unique Serial ID of eNodeB, eg. 120200002618AGP0003"
          />
          <EnodebPropertySelector
            titleLabel="eNodeB Device Class"
            value={this.state.editingEnodeb.deviceClass}
            valueOptionsByKey={EnodebDeviceClass}
            onChange={this.onDeviceClassChange}
            className={classes.input}
          />
          <TextField
            label="EARFCNDL"
            className={classes.input}
            value={this.state.editingEnodeb.earfcndl}
            onChange={this.onEarfcndlChange}
            placeholder="0-65535"
            error={!this.isEarfcndlValid()}
          />
          <TextField
            label="Subframe Assignment"
            className={classes.input}
            value={this.state.editingEnodeb.subframeAssignment}
            onChange={this.onSubframeAssignmentChange}
            placeholder="0-6"
            error={!this.isSubframeAssignmentValid()}
          />
          <TextField
            label="Special Subframe Pattern"
            className={classes.input}
            value={this.state.editingEnodeb.specialSubframePattern}
            onChange={this.onSpecialSubframePatternChange}
            inputProps={{min: 0, max: 9}}
            placeholder="0-9"
            error={!this.isSpecialSubframePatternValid()}
          />
          <TextField
            label="Physical Cell Identifier"
            className={classes.input}
            value={this.state.editingEnodeb.pci}
            onChange={this.onPciChange}
            placeholder="0-504"
            error={!this.isPciValid()}
          />
          <EnodebPropertySelector
            titleLabel="eNodeB DL/UL Bandwidth (MHz)"
            value={this.state.editingEnodeb.bandwidthMhz || ''}
            valueOptionsByKey={EnodebBandwidthOption}
            onChange={this.onBandwidthMhzChange}
            className={classes.input}
          />
          <TextField
            label="Tracking Area Code"
            className={classes.input}
            value={this.state.editingEnodeb.tac}
            onChange={this.onTacChange}
            placeholder="0-65535"
            error={!this.isTacValid()}
          />
          <TextField
            label="Enodeb ID"
            className={classes.input}
            value={this.state.editingEnodeb.enodebId}
            onChange={this.onEnodebIdChange}
            placeholder="0-1048576"
            error={!this.isEnodebIdValid()}
          />
          <TextField
            disabled
            label="Cell Number"
            className={classes.input}
            value={this.state.editingEnodeb.cellNumber}
            error={false}
          />
          <FormControl className={classes.input}>
            <FormControlLabel
              control={
                <Switch
                  checked={this.state.editingEnodeb.transmitEnabled}
                  onChange={this.onTransmitEnabledChange}
                  color="primary"
                />
              }
              label="Transmit Enabled"
            />
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onClose} color="primary">
            Cancel
          </Button>
          <Button onClick={this.onSave} color="primary" variant="contained">
            Save
          </Button>
        </DialogActions>
      </Dialog>
    );
  }

  fieldChangedHandler = (
    field:
      | 'name'
      | 'serialId'
      | 'earfcndl'
      | 'subframeAssignment'
      | 'specialSubframePattern'
      | 'pci'
      | 'tac'
      | 'enodebId'
      | 'cellNumber',
  ) => (event: SyntheticEvent<HTMLInputElement>) => {
    const {target} = event;
    if (target instanceof HTMLInputElement) {
      this.setState({
        editingEnodeb: {
          ...this.state.editingEnodeb,
          [field]: target.value,
        },
      });
    } else {
      throw Error('Expected event to be a SyntheticEvent<HTMLInputElement>');
    }
  };

  onNameChange = this.fieldChangedHandler('name');
  onSerialIdChange = this.fieldChangedHandler('serialId');
  onDeviceClassChange = event => {
    const optionKey = Object.values(EnodebDeviceClass).indexOf(
      event.target.value,
    );
    if (optionKey > -1) {
      const value =
        EnodebDeviceClass[Object.keys(EnodebDeviceClass)[optionKey]];
      this.setState({
        editingEnodeb: {
          ...this.state.editingEnodeb,
          deviceClass: value,
        },
      });
    } else {
      throw Error('Expected a valid eNodeB device class selection.');
    }
  };
  onEarfcndlChange = this.fieldChangedHandler('earfcndl');
  onSubframeAssignmentChange = this.fieldChangedHandler('subframeAssignment');
  onSpecialSubframePatternChange = this.fieldChangedHandler(
    'specialSubframePattern',
  );
  onPciChange = this.fieldChangedHandler('pci');
  onBandwidthMhzChange = event => {
    const optionKey = Object.values(EnodebBandwidthOption).indexOf(
      event.target.value,
    );
    if (optionKey > -1) {
      const value =
        EnodebBandwidthOption[Object.keys(EnodebBandwidthOption)[optionKey]];
      this.setState({
        editingEnodeb: {
          ...this.state.editingEnodeb,
          bandwidthMhz: value,
        },
      });
    } else {
      throw Error('Expected a valid bandwidth (MHz) selection.');
    }
  };
  onTacChange = this.fieldChangedHandler('tac');
  onEnodebIdChange = this.fieldChangedHandler('enodebId');
  onCellNumberChange = this.fieldChangedHandler('cellNumber');
  onTransmitEnabledChange = () => {
    this.setState({
      editingEnodeb: {
        ...this.state.editingEnodeb,
        transmitEnabled: !this.state.editingEnodeb.transmitEnabled,
      },
    });
  };

  isSerialIdValid = () => this.state.editingEnodeb.serialId.length > 0;
  isEarfcndlValid = () => {
    const val = parseInt(this.state.editingEnodeb.earfcndl, 10);
    if (isNaN(val)) {
      return false;
    }
    return val >= 0 && val <= 65535;
  };
  isSubframeAssignmentValid = () => {
    const val = parseInt(this.state.editingEnodeb.subframeAssignment, 10);
    if (isNaN(val)) {
      return false;
    }
    return val >= 0 && val <= 6;
  };
  isSpecialSubframePatternValid = () => {
    const val = parseInt(this.state.editingEnodeb.specialSubframePattern, 10);
    if (isNaN(val)) {
      return false;
    }
    return val >= 0 && val <= 9;
  };
  isPciValid = () => {
    const val = parseInt(this.state.editingEnodeb.pci, 10);
    if (isNaN(val)) {
      return false;
    }
    return val >= 0 && val <= 504;
  };
  isBandwidthMhzValid = () => {
    const val = parseInt(this.state.editingEnodeb.bandwidthMhz, 10);
    if (isNaN(val)) {
      return false;
    }
    return val >= 0 && val <= 20;
  };
  isTacValid = () => {
    const val = parseInt(this.state.editingEnodeb.tac, 10);
    if (isNaN(val)) {
      return false;
    }
    return val >= 0 && val <= 65535;
  };
  isEnodebIdValid = () => {
    const val = parseInt(this.state.editingEnodeb.enodebId);
    if (isNaN(val)) {
      return false;
    }
    // Maximum value is 2^20 - 1
    return val >= 0 && val <= 1048575;
  };
  isCellNumberValid = () => {
    const val = parseInt(this.state.editingEnodeb.cellNumber);
    if (isNaN(val)) {
      return false;
    }
    // Maximum value is 2^8 - 1
    return val >= 0 && val <= 255;
  };
  isCellIdValid = () => {
    if (!this.isEnodebIdValid() || !this.isCellNumberValid) {
      return false;
    }
    const cellId = this.getCellId();
    if (isNaN(cellId)) {
      return false;
    }
    // Maximum value is 2^28 - 1
    return cellId >= 0 && cellId <= 268435455;
  };
  isTransmitEnabledValid = () => {
    return typeof this.state.editingEnodeb.transmitEnabled == 'boolean';
  };

  getCellId = (): number => {
    const enodebId = parseInt(this.state.editingEnodeb.enodebId);
    const cellNumber = parseInt(this.state.editingEnodeb.cellNumber);
    return 256 * enodebId + cellNumber;
  };

  onSave = async () => {
    if (
      !this.isSerialIdValid() ||
      !this.isEarfcndlValid() ||
      !this.isSubframeAssignmentValid() ||
      !this.isSpecialSubframePatternValid() ||
      !this.isPciValid() ||
      !this.isBandwidthMhzValid() ||
      !this.isTacValid() ||
      !this.isCellIdValid() ||
      !this.isTransmitEnabledValid()
    ) {
      this.setState({error: 'Please complete all fields with valid values'});
      return;
    }
    const enb = {
      name: this.state.editingEnodeb.name,
      serial: this.state.editingEnodeb.serialId,
      config: {
        device_class: this.state.editingEnodeb.deviceClass,
        earfcndl: parseInt(this.state.editingEnodeb.earfcndl),
        subframe_assignment: parseInt(
          this.state.editingEnodeb.subframeAssignment,
        ),
        special_subframe_pattern: parseInt(
          this.state.editingEnodeb.specialSubframePattern,
        ),
        pci: parseInt(this.state.editingEnodeb.pci),
        bandwidth_mhz: this.state.editingEnodeb.bandwidthMhz,
        tac: parseInt(this.state.editingEnodeb.tac),
        cell_id: this.getCellId(),
        transmit_enabled: this.state.editingEnodeb.transmitEnabled,
      },
    };
    const networkId = nullthrows(this.props.match.params.networkId);
    try {
      if (this.props.editingEnodeb != null) {
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
      this.props.onSave(data);
    } catch (e) {
      this.setState({error: e?.response?.data?.message || e?.message || e});
    }
  };
}

export default withStyles(styles)(withRouter(AddEditEnodebDialog));
