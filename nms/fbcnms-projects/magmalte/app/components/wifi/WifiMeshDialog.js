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

import Button from '@fbcnms/ui/components/design-system/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import FormLabel from '@material-ui/core/FormLabel';
import KeyValueFields from '@fbcnms/magmalte/app/components/KeyValueFields';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';

import {
  additionalPropsToArray,
  additionalPropsToObject,
  meshesURL,
} from './WifiUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  backdrop: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    position: 'fixed',
    zIndex: '13000',
  },
};

type Props = ContextRouter &
  WithStyles<typeof styles> & {
    onCancel: () => void,
    onSave: string => void,
  };

type State = {
  meshID: string,
  configs: {[string]: any},
  error?: string,
};

class WifiMeshDialog extends React.Component<Props, State> {
  constructor(props: Props) {
    super();

    const {meshID} = props.match.params;
    this.state = {
      meshID: meshID || '',
      configs: {},
    };

    if (meshID) {
      axios
        .get(`${meshesURL(props.match)}/${meshID}/configs`)
        .then(response =>
          this.setState({
            configs: {
              ...response.data,
              additional_props: additionalPropsToArray(
                response.data.additional_props,
              ),
            },
          }),
        )
        .catch(() => this.props.onCancel());
    }
  }

  render() {
    const {meshID} = this.props.match.params;
    if (meshID && Object.keys(this.state.configs).length === 0) {
      return <LoadingFillerBackdrop />;
    }

    return (
      <Dialog open={true} onClose={this.props.onCancel}>
        <DialogTitle>{meshID ? 'Edit Mesh' : 'New Mesh'}</DialogTitle>
        <DialogContent>
          {this.state.error ? (
            <FormLabel error>{this.state.error}</FormLabel>
          ) : null}
          <FormGroup row>
            <TextField
              required
              className={this.props.classes.input}
              label="Mesh Name"
              margin="normal"
              onChange={this.handlemeshIDChange}
              value={this.state.meshID}
              disabled={!!meshID}
            />
            <TextField
              required
              className={this.props.classes.input}
              label="SSID"
              margin="normal"
              value={this.state.configs.ssid}
              onChange={this.handleSSIDChange}
            />
            <TextField
              className={this.props.classes.input}
              label="Password"
              margin="normal"
              value={this.state.configs.password}
              onChange={this.handlePasswordChange}
            />
            <FormControlLabel
              control={
                <Checkbox
                  checked={this.state.configs.xwf_enabled}
                  onChange={this.handledEnableXWFChange}
                  color="primary"
                />
              }
              label="Enable XWF"
            />
            <KeyValueFields
              keyValuePairs={this.state.configs.additional_props || [['', '']]}
              onChange={this.handleAdditionalPropsChange}
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

  handlemeshIDChange = event => this.setState({meshID: event.target.value});
  handleSSIDChange = ({target}) =>
    this.setState({
      configs: {
        ...this.state.configs,
        ssid: target.value,
      },
    });
  handlePasswordChange = ({target}) =>
    this.setState({
      configs: {
        ...this.state.configs,
        password: target.value,
      },
    });
  handledEnableXWFChange = ({target}) =>
    this.setState({
      configs: {
        ...this.state.configs,
        xwf_enabled: target.checked,
      },
    });
  handleAdditionalPropsChange = value =>
    this.setState({
      configs: {
        ...this.state.configs,
        additional_props: value,
      },
    });

  onSave = async () => {
    try {
      const editingMeshID = this.props.match.params.meshID;
      const url = meshesURL(this.props.match);
      if (editingMeshID) {
        await axios.put(`${url}/${editingMeshID}/configs`, this.getConfigs());
        this.props.onSave(editingMeshID);
        return;
      }

      // creating a mesh requires two steps:
      // 1st Step: create the mesh object
      await axios.post(url, {mesh_id: this.state.meshID, gateway_ids: []});

      // 2nd Step: create the config object
      await axios.post(
        `${url}/${this.state.meshID}/configs`,
        this.getConfigs(),
      );

      this.props.onSave(this.state.meshID);
    } catch (e) {
      this.setState({error: e.response.data.message || e.message});
    }
  };

  getConfigs = () => {
    return {
      ...this.state.configs,
      mesh_frequency: parseInt(this.state.configs.mesh_frequency),
      additional_props: additionalPropsToObject(
        this.state.configs.additional_props,
      ),
    };
  };
}

export default withStyles(styles)(withRouter(WifiMeshDialog));
