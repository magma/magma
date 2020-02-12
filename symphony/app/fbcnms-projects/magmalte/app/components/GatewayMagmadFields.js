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
import type {GatewayV1} from './GatewayUtils';
import type {WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import {toString} from './GatewayUtils';
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
    onClose: () => void,
    onSave: (gatewayID: string) => void,
    gateway: GatewayV1,
  };

type State = {
  autoupgradeEnabled: boolean,
  autoupgradePollInterval: string,
  checkinInterval: string,
  checkinTimeout: string,
};

class GatewayMagmadFields extends React.Component<Props, State> {
  state = {
    autoupgradeEnabled: this.props.gateway.autoupgradeEnabled,
    autoupgradePollInterval: toString(
      this.props.gateway.autoupgradePollInterval,
    ),
    checkinInterval: toString(this.props.gateway.checkinInterval),
    checkinTimeout: toString(this.props.gateway.checkinTimeout),
  };

  render() {
    return (
      <>
        <DialogContent>
          <FormControl className={this.props.classes.input}>
            <InputLabel htmlFor="autoupgradeEnabled">
              Autoupgrade Enabled
            </InputLabel>
            <Select
              inputProps={{id: 'autoupgradeEnabled'}}
              value={this.state.autoupgradeEnabled ? 1 : 0}
              onChange={this.autoupgradeEnabledChanged}>
              <MenuItem value={1}>Enabled</MenuItem>
              <MenuItem value={0}>Disabled</MenuItem>
            </Select>
          </FormControl>
          <TextField
            label="Autoupgrade Poll Interval (seconds)"
            className={this.props.classes.input}
            value={this.state.autoupgradePollInterval}
            onChange={this.autoupgradePollIntervalChanged}
            placeholder="E.g. 300"
          />
          <TextField
            label="Checkin Interval (seconds)"
            className={this.props.classes.input}
            value={this.state.checkinInterval}
            onChange={this.checkinIntervalChanged}
            placeholder="E.g. 60"
          />
          <TextField
            label="Checkin Timeout (seconds)"
            className={this.props.classes.input}
            value={this.state.checkinTimeout}
            onChange={this.checkinTimeoutChanged}
            placeholder="E.g. 5"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onClose} skin="regular">
            Cancel
          </Button>
          <Button onClick={this.onSave}>Save</Button>
        </DialogActions>
      </>
    );
  }

  onSave = () => {
    const magmad = {
      autoupgrade_enabled: this.state.autoupgradeEnabled,
      autoupgrade_poll_interval: parseInt(this.state.autoupgradePollInterval),
      checkin_interval: parseInt(this.state.checkinInterval),
      checkin_timeout: parseInt(this.state.checkinTimeout),
      tier: this.props.gateway.tier,
    };

    MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdMagmad({
      networkId: nullthrows(this.props.match.params.networkId),
      gatewayId: this.props.gateway.logicalID,
      magmad,
    }).then(() => this.props.onSave(this.props.gateway.logicalID));
  };

  autoupgradeEnabledChanged = ({target}) =>
    this.setState({autoupgradeEnabled: !!target.value});

  autoupgradePollIntervalChanged = ({target}) =>
    this.setState({autoupgradePollInterval: target.value});

  checkinIntervalChanged = ({target}) =>
    this.setState({checkinInterval: target.value});

  checkinTimeoutChanged = ({target}) =>
    this.setState({checkinTimeout: target.value});
}

export default withStyles(styles)(withRouter(GatewayMagmadFields));
