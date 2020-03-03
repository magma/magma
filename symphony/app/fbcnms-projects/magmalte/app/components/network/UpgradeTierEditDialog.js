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
import type {tier} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormGroup from '@material-ui/core/FormGroup';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import {withRouter} from 'react-router-dom';

type Props = ContextRouter & {
  onSave: tier => void,
  onCancel: () => void,
  tierId?: string,
};

type State = {
  isLoading: boolean,
  tier: tier,
};

class UpgradeTierEditDialog extends React.Component<Props, State> {
  state = {
    isLoading: !this.isNewTier(),
    tier: {
      id: this.props.tierId || '',
      name: '',
      version: '',
      images: [],
      gateways: [],
    },
  };

  componentDidMount() {
    const {tierId, match} = this.props;
    if (tierId) {
      MagmaV1API.getNetworksByNetworkIdTiersByTierId({
        networkId: nullthrows(match.params.networkId),
        tierId,
      }).then(tier =>
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
    const {tierId: initialTierId} = this.props;
    const {isLoading, tier} = this.state;
    return (
      <Dialog open={!isLoading} onClose={this.props.onCancel} scroll="body">
        <DialogTitle>
          {this.isNewTier() ? 'Add Upgrade Tier' : 'Edit Upgrade Tier'}
        </DialogTitle>
        <DialogContent>
          <FormGroup row>
            <TextField
              required
              label="Tier ID"
              placeholder="E.g. t1"
              margin="normal"
              disabled={!!initialTierId}
              value={tier.id}
              onChange={this.handleIdChanged}
            />
          </FormGroup>
          <FormGroup row>
            <TextField
              required
              label="Tier Name"
              placeholder="E.g. Example Tier"
              margin="normal"
              value={tier.name}
              onChange={this.handleNameChanged}
            />
          </FormGroup>
          <FormGroup row>
            <TextField
              required
              label="Tier Version"
              placeholder="E.g. 1.0.0-0"
              margin="normal"
              value={tier.version}
              onChange={this.handleVersionChanged}
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

  onSave = () => {
    const {tier} = this.state;
    if (this.isNewTier()) {
      MagmaV1API.postNetworksByNetworkIdTiers({
        networkId: nullthrows(this.props.match.params.networkId),
        tier,
      })
        .then(() => this.props.onSave(tier))
        .catch(console.error);
    } else {
      MagmaV1API.putNetworksByNetworkIdTiersByTierId({
        networkId: nullthrows(this.props.match.params.networkId),
        tierId: tier.id,
        tier,
      })
        .then(_resp => this.props.onSave(tier))
        .catch(console.error);
    }
  };
}

export default withRouter(UpgradeTierEditDialog);
