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
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {network_epc_configs} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormGroup,
  InputAdornment,
  TextField,
} from '@material-ui/core';

import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withRouter} from 'react-router-dom';

import {BITRATE_MULTIPLIER, DATA_PLAN_UNLIMITED_RATES} from './DataPlanConst';

type CellularNetworkProfile = $Values<
  $NonMaybeType<$PropertyType<network_epc_configs, 'sub_profiles'>>,
>;

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
  WithAlert & {
    onSave: (dataPlanId: string, newEPCConfig: network_epc_configs) => void,
    onCancel: () => void,
    epcConfig: network_epc_configs,
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
    const {dataPlanId} = this.props;
    const dataPlan = this._getDataPlan();
    return (
      <Dialog open={true} onClose={this.props.onCancel} scroll="body">
        <DialogTitle>{dataPlanId ? 'Edit' : 'Add'} Data Plan</DialogTitle>
        <DialogContent>
          <FormGroup row>
            <TextField
              required
              label="Name"
              margin="normal"
              disabled={!!dataPlanId}
              value={this._getDataPlanIdField()}
              onChange={this.handleNameChanged}
            />
          </FormGroup>
          <FormGroup row>
            <MegabyteTextField
              label="Download Limit"
              dataPlan={dataPlan}
              field="max_dl_bit_rate"
              value={this.state.editedMaxDlBitRate}
              onChange={this.handleDownloadLimitChanged}
            />
          </FormGroup>
          <FormGroup row>
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
          <Button onClick={this.props.onCancel} skin="regular">
            Cancel
          </Button>
          <Button onClick={this.onSave}>Save</Button>
        </DialogActions>
      </Dialog>
    );
  }

  _getDataPlan(): ?CellularNetworkProfile {
    const {dataPlanId, epcConfig} = this.props;
    return (dataPlanId && epcConfig.sub_profiles?.[dataPlanId]) || null;
  }

  _getDataPlanIdField(): string {
    const dataPlanId =
      this.state.editedName !== null
        ? this.state.editedName
        : this.props.dataPlanId;
    return dataPlanId || '';
  }

  onSave = () => {
    const {match} = this.props;
    const epcConfig = this.props.epcConfig;
    const dataPlanId = this._getDataPlanIdField();
    const dataPlan = this._getDataPlan();
    if (!this.props.dataPlanId && (epcConfig.sub_profiles || {})[dataPlanId]) {
      this.props.alert(
        'Data plan name is already used. Please use a different name',
      );
      return;
    }

    const newConfig = {
      ...epcConfig,
      sub_profiles: {
        ...(epcConfig.sub_profiles || {}),
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
    };

    MagmaV1API.putLteByNetworkIdCellularEpc({
      networkId: nullthrows(match.params.networkId),
      config: newConfig,
    })
      .then(_resp => this.props.onSave(dataPlanId, newConfig))
      .catch(error => this.props.alert(error.response?.data?.message || error));
  };
}

export default withAlert(withRouter(DataPlanEditDialog));
