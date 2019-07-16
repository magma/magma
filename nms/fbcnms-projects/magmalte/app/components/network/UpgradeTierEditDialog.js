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
import type {NetworkUpgradeTier} from '../../common/MagmaAPIType';
import type {WithStyles} from '@material-ui/core';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormGroup from '@material-ui/core/FormGroup';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';

import nullthrows from '@fbcnms/util/nullthrows';
import {MagmaAPIUrls, fetchNetworkUpgradeTier} from '../../common/MagmaAPI';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = {};

type Props = ContextRouter &
  WithStyles & {
    onSave: (config: NetworkUpgradeTier) => void,
    onCancel: () => void,
    tierId: ?string,
  };

type State = {
  isLoading: boolean,
  tier: NetworkUpgradeTier,
};

class UpgradeTierEditDialog extends React.Component<Props, State> {
  state = {
    isLoading: !this.isNewTier(),
    tier: {
      id: this.props.tierId || '',
      name: '',
      version: '',
      images: [],
    },
  };

  componentDidMount() {
    const {tierId, match} = this.props;
    if (tierId) {
      fetchNetworkUpgradeTier(nullthrows(match.params.networkId), tierId).then(
        tier =>
          this.setState({
            isLoading: false,
            tier,
          }),
      );
    }
  }

  _handleFieldChange = (evt: SyntheticInputEvent<*>, field: string): void => {
    this.setState({
      tier: {
        ...this.state.tier,
        [field]: evt.target.value,
      },
    });
  };

  handleIdChanged = evt => this._handleFieldChange(evt, 'id');
  handleNameChanged = evt => this._handleFieldChange(evt, 'name');
  handleVersionChanged = evt => this._handleFieldChange(evt, 'version');

  isNewTier() {
    return !this.props.tierId;
  }

  render() {
    const {classes, tierId: initialTierId} = this.props;
    const {isLoading, tier} = this.state;
    return (
      <Dialog open={!isLoading} onClose={this.props.onCancel} scroll="body">
        <DialogTitle>
          {this.isNewTier() ? 'Add Upgrade Tier' : 'Edit Upgrade Tier'}
        </DialogTitle>
        <DialogContent>
          <FormGroup row className={classes.formGroup}>
            <TextField
              required
              label="Tier ID"
              placeholder="E.g. t1"
              margin="normal"
              disabled={!!initialTierId}
              className={classes.textField}
              value={tier.id}
              onChange={this.handleIdChanged}
            />
          </FormGroup>
          <FormGroup row className={classes.formGroup}>
            <TextField
              required
              label="Tier Name"
              placeholder="E.g. Example Tier"
              margin="normal"
              className={classes.textField}
              value={tier.name}
              onChange={this.handleNameChanged}
            />
          </FormGroup>
          <FormGroup row className={classes.formGroup}>
            <TextField
              required
              label="Tier Version"
              placeholder="E.g. 1.0.0-0"
              margin="normal"
              className={classes.textField}
              value={tier.version}
              onChange={this.handleVersionChanged}
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

  onSave = () => {
    const {tier} = this.state;
    if (this.isNewTier()) {
      axios
        .post(MagmaAPIUrls.networkTiers(this.props.match), tier)
        .then(resp => {
          const newTier: NetworkUpgradeTier = {
            ...tier,
            id: resp.data,
          };
          this.props.onSave(newTier);
        })
        .catch(console.error);
    } else {
      axios
        .put(MagmaAPIUrls.networkTier(this.props.match, tier.id), tier)
        .then(_resp => this.props.onSave(tier))
        .catch(console.error);
    }
  };
}

export default withStyles(styles)(withRouter(UpgradeTierEditDialog));
