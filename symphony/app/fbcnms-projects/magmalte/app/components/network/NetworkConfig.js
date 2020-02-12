/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';
import type {
  network_cellular_configs,
  network_ran_configs,
} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import FormHelperText from '@material-ui/core/FormHelperText';
import IconButton from '@material-ui/core/IconButton';
import Input from '@material-ui/core/Input';
import InputAdornment from '@material-ui/core/InputAdornment';
import InputLabel from '@material-ui/core/InputLabel';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import VisibilityIcon from '@material-ui/icons/Visibility';
import VisibilityOffIcon from '@material-ui/icons/VisibilityOff';

import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {base64ToHex, hexToBase64, isValidHex} from '@fbcnms/util/strings';
import {get, merge} from 'lodash';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  formContainer: {
    paddingBottom: theme.spacing(2),
  },
  formGroup: {
    marginLeft: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  select: {
    marginRight: theme.spacing(),
    minWidth: 200,
  },
  saveButton: {
    marginTop: theme.spacing(2),
  },
  textField: {
    marginRight: theme.spacing(),
  },
});

type State = {
  config: ?network_cellular_configs,
  isLoading: boolean,
  lteAuthOpHex: string,
  showLteAuthOP: boolean,
  bandSelection: string,
  tddConfig: ?TDDConfig,
  fddConfig: ?FDDConfig,
};

type Props = ContextRouter & WithAlert & WithStyles<typeof styles> & {};

const TDD = 'tdd';
const FDD = 'fdd';
type TDDConfig = $PropertyType<network_ran_configs, 'tdd_config'>;
type FDDConfig = $PropertyType<network_ran_configs, 'fdd_config'>;

class NetworkConfig extends React.Component<Props, State> {
  state = {
    config: null,
    isLoading: true,
    lteAuthOpHex: '',
    showLteAuthOP: false,
    bandSelection: '',
    fddConfig: null,
    tddConfig: null,
  };

  componentDidMount() {
    MagmaV1API.getLteByNetworkIdCellular({
      networkId: nullthrows(this.props.match.params.networkId),
    })
      .then(response =>
        this.setState({
          config: response,
          isLoading: false,
          lteAuthOpHex: base64ToHex(response.epc.lte_auth_op),
          bandSelection: response.ran.fdd_config ? FDD : TDD,
          tddConfig: response.ran.tdd_config || {
            earfcndl: 0,
            special_subframe_pattern: 0,
            subframe_assignment: 0,
          },
          fddConfig: response.ran.fdd_config || {
            earfcndl: 0,
            earfcnul: 0,
          },
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
      const newConfig = {
        ...this.state.config,
        // $FlowFixMe Set state for each field
        [epcOrRan]: {
          ...this.state.config[epcOrRan],
          // $FlowFixMe Set state for each field
          [field]: evt.target.value,
        },
      };
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
  handleMccChanged = this.updateNetworkConfigField('epc', 'mcc');
  handleMncChanged = this.updateNetworkConfigField('epc', 'mnc');
  handleTacChanged = this.updateNetworkConfigField('epc', 'tac');

  handleTddOrFddConfigChanged = (
    tddOrFdd: 'tddConfig' | 'fddConfig',
    field: string,
  ) => {
    return evt =>
      this.setState({
        // $FlowFixMe Set state for each field
        [tddOrFdd]: {
          ...this.state[tddOrFdd],
          [field]: evt.target.value,
        },
      });
  };

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

    let bandeSelectionFields;
    if (this.state.bandSelection === FDD) {
      bandeSelectionFields = (
        <FormGroup row className={classes.formGroup}>
          <TextField
            required
            label="EARFCNDL"
            margin="normal"
            className={classes.textField}
            value={this.state.fddConfig?.earfcndl}
            onChange={this.handleTddOrFddConfigChanged('fddConfig', 'earfcndl')}
          />
          <TextField
            required
            label="EARFCNUL"
            margin="normal"
            className={classes.textField}
            value={this.state.fddConfig?.earfcnul}
            onChange={this.handleTddOrFddConfigChanged('fddConfig', 'earfcnul')}
          />
        </FormGroup>
      );
    } else {
      bandeSelectionFields = (
        <FormGroup row className={classes.formGroup}>
          <TextField
            required
            label="EARFCNDL"
            margin="normal"
            className={classes.textField}
            value={this.state.tddConfig?.earfcndl}
            onChange={this.handleTddOrFddConfigChanged('tddConfig', 'earfcndl')}
          />
          <TextField
            required
            label="Special Subframe Pattern"
            margin="normal"
            className={classes.textField}
            value={this.state.tddConfig?.special_subframe_pattern}
            onChange={this.handleTddOrFddConfigChanged(
              'tddConfig',
              'special_subframe_pattern',
            )}
          />
          <TextField
            required
            label="Subframe Assignment"
            margin="normal"
            className={classes.textField}
            value={this.state.tddConfig?.subframe_assignment}
            onChange={this.handleTddOrFddConfigChanged(
              'tddConfig',
              'subframe_assignment',
            )}
          />
        </FormGroup>
      );
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
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <FormControl className={classes.select}>
            <InputLabel htmlFor="band_selection">Band Selection</InputLabel>
            <Select
              inputProps={{id: 'bend_selection'}}
              value={this.state.bandSelection}
              onChange={({target}) => {
                this.setState({bandSelection: target.value});
              }}>
              <MenuItem value={TDD}>TDD</MenuItem>
              <MenuItem value={FDD}>FDD</MenuItem>
            </Select>
          </FormControl>
        </FormGroup>
        {bandeSelectionFields}
        <FormGroup row className={classes.formGroup}>
          <Button
            disabled={!this.canSubmitForm()}
            className={classes.saveButton}
            onClick={this.handleSave}>
            Save
          </Button>
        </FormGroup>
      </div>
    );
  }

  handleSave = () => {
    const config = nullthrows(this.state.config);
    const bandSeletionConfig: {|
      tdd_config?: TDDConfig,
      fdd_config?: FDDConfig,
    |} = {tdd_config: undefined, fdd_config: undefined};
    if (this.state.bandSelection === TDD) {
      const tddConfig = nullthrows(this.state.tddConfig);
      bandSeletionConfig.tdd_config = {
        earfcndl: parseInt(tddConfig.earfcndl),
        special_subframe_pattern: parseInt(tddConfig.special_subframe_pattern),
        subframe_assignment: parseInt(tddConfig.subframe_assignment),
      };
    } else {
      const fddConfig = nullthrows(this.state.fddConfig);
      bandSeletionConfig.fdd_config = {
        earfcndl: parseInt(fddConfig.earfcndl),
        earfcnul: parseInt(fddConfig.earfcnul),
      };
    }

    MagmaV1API.putLteByNetworkIdCellular({
      networkId: nullthrows(this.props.match.params.networkId),
      config: {
        ...config,
        ran: {
          ...config.ran,
          ...bandSeletionConfig,
        },
        epc: {
          ...config.epc,
          tac: parseInt(config.epc.tac),
        },
      },
    })
      .then(_resp => {
        this.props.alert('Saved successfully');
      })
      .catch(this.props.alert);
  };
}

export default withStyles(styles)(withAlert(withRouter(NetworkConfig)));
