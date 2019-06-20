/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CellularNetworkConfig} from '../../common/MagmaAPIType';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {ContextRouter} from 'react-router-dom';
import type {WithStyles} from '@material-ui/core';

import React from 'react';
import axios from 'axios';
import FormGroup from '@material-ui/core/FormGroup';
import FormHelperText from '@material-ui/core/FormHelperText';
import Input from '@material-ui/core/Input';
import InputAdornment from '@material-ui/core/InputAdornment';
import IconButton from '@material-ui/core/IconButton';
import InputLabel from '@material-ui/core/InputLabel';
import LoadingFiller from '../LoadingFiller';
import MenuItem from '@material-ui/core/MenuItem';
import Button from '@material-ui/core/Button';
import Select from '@material-ui/core/Select';
import FormControl from '@material-ui/core/FormControl';
import TextField from '@material-ui/core/TextField';
import VisibilityIcon from '@material-ui/icons/Visibility';
import VisibilityOffIcon from '@material-ui/icons/VisibilityOff';

import {MagmaAPIUrls} from '../../common/MagmaAPI';
import {get, merge} from 'lodash';
import {base64ToHex, hexToBase64, isValidHex} from '@fbcnms/util/strings';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  formContainer: {
    paddingBottom: theme.spacing.unit * 2,
  },
  formGroup: {
    marginLeft: theme.spacing.unit * 2,
    marginBottom: theme.spacing.unit * 2,
  },
  select: {
    marginRight: theme.spacing.unit,
    minWidth: 200,
  },
  saveButton: {
    marginTop: theme.spacing.unit * 2,
  },
  textField: {
    marginRight: theme.spacing.unit,
  },
});

type State = {
  config: ?CellularNetworkConfig,
  isLoading: boolean,
  lteAuthOpHex: string,
  showLteAuthOP: boolean,
};

type Props = ContextRouter & WithAlert & WithStyles & {};

class DataPlanConfig extends React.Component<Props, State> {
  state = {
    config: null,
    isLoading: true,
    lteAuthOpHex: '',
    showLteAuthOP: false,
  };

  componentDidMount() {
    axios
      .get(MagmaAPIUrls.networkConfigsForType(this.props.match, 'cellular'))
      .then(response =>
        this.setState({
          config: response.data,
          isLoading: false,
          lteAuthOpHex: base64ToHex(response.data.epc.lte_auth_op),
        }),
      )
      .catch(error => {
        this.props.alert(get(error, 'response.data.message', error));
        this.setState({
          isLoading: false,
        });
      });
  }

  updateNetworkConfigField = (
    epcOrRan: 'epc' | 'ran',
    field: string | number,
  ) => {
    return evt => {
      if (!this.state.config) {
        return;
      }
      const newConfig = merge({}, this.state.config, {
        [epcOrRan]: {
          [field]: evt.target.value,
        },
      });
      this.setState({config: newConfig});
    };
  };

  handleBandwidthChanged = this.updateNetworkConfigField(
    'ran',
    'bandwidth_mhz',
  );
  handleLteAuthOpChanged = evt => {
    this.setState({lteAuthOpHex: evt.target.value});
    this.setState({
      config: merge({}, this.state.config, {
        epc: {
          lte_auth_op: hexToBase64(evt.target.value),
        },
      }),
    });
  };
  handleEarfcndlChanged = this.updateNetworkConfigField('ran', 'earfcndl');
  handleMccChanged = this.updateNetworkConfigField('epc', 'mcc');
  handleMncChanged = this.updateNetworkConfigField('epc', 'mnc');
  handleTacChanged = this.updateNetworkConfigField('epc', 'tac');
  handleUploadDownloadRatioChanged = this.updateNetworkConfigField(
    'ran',
    'ul_dl_ratio',
  );

  handleMouseDownPassword = event => {
    event.preventDefault();
  };

  handleClickShowPassword = () => {
    this.setState(state => ({showLteAuthOP: !state.showLteAuthOP}));
  };

  canSubmitForm(): boolean {
    return isValidHex(this.state.lteAuthOpHex);
  }

  render() {
    const {classes} = this.props;
    const {config, lteAuthOpHex, showLteAuthOP} = this.state;
    if (!config) {
      return <LoadingFiller />;
    }

    return (
      <div className={classes.formContainer}>
        <FormGroup row className={classes.formGroup}>
          <TextField
            required
            label="MCC"
            margin="normal"
            className={classes.textField}
            value={config.epc.mcc}
            onChange={this.handleMccChanged}
          />
          <TextField
            required
            label="MNC"
            margin="normal"
            className={classes.textField}
            value={config.epc.mnc}
            onChange={this.handleMncChanged}
          />
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <TextField
            required
            label="TAC"
            margin="normal"
            className={classes.textField}
            value={config.epc.tac}
            onChange={this.handleTacChanged}
          />
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <TextField
            required
            label="EARFCNDL"
            margin="normal"
            className={classes.textField}
            value={config.ran.earfcndl}
            onChange={this.handleEarfcndlChanged}
          />
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <FormControl
            className={classes.textField}
            error={!isValidHex(lteAuthOpHex)}>
            <InputLabel htmlFor="lte_auth_op">Auth OP</InputLabel>
            <Input
              id="lte_auth_op"
              type={showLteAuthOP ? 'text' : 'password'}
              value={lteAuthOpHex}
              onChange={this.handleLteAuthOpChanged}
              endAdornment={
                <InputAdornment position="end">
                  <IconButton
                    onClick={this.handleClickShowPassword}
                    onMouseDown={this.handleMouseDownPassword}>
                    {showLteAuthOP ? <VisibilityOffIcon /> : <VisibilityIcon />}
                  </IconButton>
                </InputAdornment>
              }
            />
            {!isValidHex(lteAuthOpHex) && (
              <FormHelperText>Invalid hex value</FormHelperText>
            )}
          </FormControl>
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <FormControl className={classes.select}>
            <InputLabel htmlFor="">Bandwidth (Mhz)</InputLabel>
            <Select
              value={config.ran.bandwidth_mhz}
              onChange={this.handleBandwidthChanged}>
              <MenuItem value={3}>3</MenuItem>
              <MenuItem value={5}>5</MenuItem>
              <MenuItem value={10}>10</MenuItem>
              <MenuItem value={15}>15</MenuItem>
              <MenuItem value={20}>20</MenuItem>
            </Select>
          </FormControl>
          <FormControl className={classes.select}>
            <InputLabel htmlFor="">DL:UL Ratio</InputLabel>
            <Select
              value={config.ran.ul_dl_ratio || 2}
              onChange={this.handleUploadDownloadRatioChanged}>
              <MenuItem value={2}>3:1</MenuItem>
              <MenuItem value={1}>1:1</MenuItem>
            </Select>
          </FormControl>
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <Button
            disabled={!this.canSubmitForm()}
            className={classes.saveButton}
            variant="contained"
            color="primary"
            onClick={this.handleSave}>
            Save
          </Button>
        </FormGroup>
      </div>
    );
  }

  handleSave = () => {
    const {config} = this.state;
    axios
      .put(
        MagmaAPIUrls.networkConfigsForType(this.props.match, 'cellular'),
        config,
      )
      .then(_resp => {
        this.props.alert('Saved successfully');
      })
      .catch(this.props.alert);
  };
}

export default withStyles(styles)(withAlert(withRouter(DataPlanConfig)));
