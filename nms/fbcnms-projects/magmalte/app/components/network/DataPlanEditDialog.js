/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  CellularNetworkConfig,
  CellularNetworkProfile,
} from '../../common/MagmaAPIType';
import type {ContextRouter} from 'react-router-dom';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import React from 'react';
import axios from 'axios';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormGroup,
  InputAdornment,
  TextField,
  withStyles,
} from '@material-ui/core';

import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {MagmaAPIUrls} from '../../common/MagmaAPI';
import {merge} from 'lodash';
import {withRouter} from 'react-router-dom';

import {BITRATE_MULTIPLIER, DATA_PLAN_UNLIMITED_RATES} from './DataPlanConst';

const styles = {};

// Calculates the bitrate value in bps based on the default
// value input in the form, and the previous (default) bitrate
function _getBitRateValue(props: {
  dataPlan: ?CellularNetworkProfile,
  editedValue: ?string,
  field: 'max_ul_bit_rate' | 'max_dl_bit_rate',
}): number {
  const {dataPlan, editedValue, field} = props;
  if (editedValue !== null && editedValue === '') {
    return DATA_PLAN_UNLIMITED_RATES[field];
  } else if (editedValue !== null) {
    return parseFloat(editedValue) * BITRATE_MULTIPLIER;
  }
  return dataPlan ? dataPlan[field] : DATA_PLAN_UNLIMITED_RATES[field];
}

const MegabyteTextField = (props: {
  label: string,
  field: 'max_ul_bit_rate' | 'max_dl_bit_rate',
  value: ?string,
  dataPlan: ?CellularNetworkProfile,
  onChange: (event: SyntheticInputEvent<*>) => void,
}) => {
  const {dataPlan, field, label, onChange, value} = props;
  const defaultMaxBitRate = dataPlan ? dataPlan[field] : null;
  let fieldValue;
  if (value !== null) {
    fieldValue = value || '';
  } else if (defaultMaxBitRate === null) {
    fieldValue = '';
  } else if (defaultMaxBitRate === DATA_PLAN_UNLIMITED_RATES[field]) {
    fieldValue = '';
  } else {
    fieldValue = (defaultMaxBitRate || 0) / BITRATE_MULTIPLIER + '';
  }

  return (
    <TextField
      required
      type="number"
      label={label}
      margin="normal"
      placeholder="Unlimited"
      value={fieldValue}
      InputProps={{
        endAdornment: <InputAdornment position="end">Mbps</InputAdornment>,
      }}
      InputLabelProps={{
        shrink: true,
      }}
      onChange={onChange}
    />
  );
};

type Props = ContextRouter &
  WithStyles &
  WithAlert & {
    onSave: (
      dataPlanId: string,
      newNetworkConfig: CellularNetworkConfig,
    ) => void,
    onCancel: () => void,
    networkConfig: ?CellularNetworkConfig,
    dataPlanId: ?string,
  };

type State = {
  editedName: ?string,
  editedMaxUlBitRate: ?string,
  editedMaxDlBitRate: ?string,
};

class DataPlanEditDialog extends React.Component<Props, State> {
  state = {
    editedName: null,
    editedMaxUlBitRate: null,
    editedMaxDlBitRate: null,
  };

  handleDownloadLimitChanged = evt =>
    this.setState({editedMaxDlBitRate: evt.target.value});
  handleNameChanged = evt => this.setState({editedName: evt.target.value});
  handleUploadLimitChanged = evt =>
    this.setState({editedMaxUlBitRate: evt.target.value});

  render() {
    const {classes, dataPlanId} = this.props;
    const dataPlan = this._getDataPlan();
    return (
      <Dialog open={true} onClose={this.props.onCancel} scroll="body">
        <DialogTitle>{dataPlanId ? 'Edit' : 'Add'} Data Plan</DialogTitle>
        <DialogContent>
          <FormGroup row className={classes.formGroup}>
            <TextField
              required
              label="Name"
              margin="normal"
              disabled={!!dataPlanId}
              className={classes.textField}
              value={this._getDataPlanIdField()}
              onChange={this.handleNameChanged}
            />
          </FormGroup>
          <FormGroup row className={classes.formGroup}>
            <MegabyteTextField
              label="Download Limit"
              dataPlan={dataPlan}
              field="max_dl_bit_rate"
              value={this.state.editedMaxDlBitRate}
              onChange={this.handleDownloadLimitChanged}
            />
          </FormGroup>
          <FormGroup row className={classes.formGroup}>
            <MegabyteTextField
              label="Upload Limit"
              dataPlan={dataPlan}
              field="max_ul_bit_rate"
              value={this.state.editedMaxUlBitRate}
              onChange={this.handleUploadLimitChanged}
            />
          </FormGroup>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onCancel} color="primary">
            Cancel
          </Button>
          <Button onClick={this.onSave} color="primary" variant="contained">
            Save
          </Button>
        </DialogActions>
      </Dialog>
    );
  }

  _getDataPlan(): ?CellularNetworkProfile {
    const {dataPlanId, networkConfig} = this.props;
    return dataPlanId && networkConfig && networkConfig.epc
      ? networkConfig.epc.sub_profiles[dataPlanId]
      : null;
  }

  _getDataPlanIdField(): string {
    const dataPlanId =
      this.state.editedName !== null
        ? this.state.editedName
        : this.props.dataPlanId;
    return dataPlanId || '';
  }

  onSave = () => {
    const {match, networkConfig} = this.props;
    const dataPlanId = this._getDataPlanIdField();
    const dataPlan = this._getDataPlan();
    const newConfig = merge({}, networkConfig, {
      epc: {
        sub_profiles: {
          [dataPlanId]: {
            max_dl_bit_rate: _getBitRateValue({
              dataPlan,
              field: 'max_dl_bit_rate',
              editedValue: this.state.editedMaxDlBitRate,
            }),
            max_ul_bit_rate: _getBitRateValue({
              dataPlan,
              field: 'max_ul_bit_rate',
              editedValue: this.state.editedMaxUlBitRate,
            }),
          },
        },
      },
    });
    axios
      .put(MagmaAPIUrls.networkConfigsForType(match, 'cellular'), newConfig)
      .then(_resp => this.props.onSave(dataPlanId, newConfig))
      .catch(error => this.props.alert(error.response?.data?.message || error));
  };
}

export default withStyles(styles)(withAlert(withRouter(DataPlanEditDialog)));
