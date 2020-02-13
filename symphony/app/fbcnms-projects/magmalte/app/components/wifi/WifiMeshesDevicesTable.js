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
import type {WifiGateway} from './WifiUtils';
import type {WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import IconButton from '@material-ui/core/IconButton';
import LinearProgress from '@material-ui/core/LinearProgress';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import RefreshIcon from '@material-ui/icons/Refresh';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import WifiDeviceDialog from './WifiDeviceDialog';
import WifiMeshDialog from './WifiMeshDialog';
import WifiMeshRow from './WifiMeshRow';
import nullthrows from '@fbcnms/util/nullthrows';

import {Route, withRouter} from 'react-router-dom';
import {buildWifiGatewayFromPayloadV1} from './WifiUtils';
import {map} from 'lodash';
import {sortBy} from 'lodash';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  actionsColumn: {
    width: '160px',
  },
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  infoColumn: {
    width: '400px',
  },
  paper: {
    margin: theme.spacing(3),
  },
});

type Props = ContextRouter & WithStyles<typeof styles> & {};

type State = {
  isLoading: boolean,
  meshes: Map<string, WifiGateway[]>,
  errorMessage: ?string,
  lastRefreshTime: string,
};

class WifiMeshesDevicesTable extends React.Component<Props, State> {
  state = {
    isLoading: false,
    meshes: new Map(),
    errorMessage: null,
    lastRefreshTime: new Date().toLocaleString(),
  };

  componentDidMount() {
    this.fetchMeshes();
  }

  render() {
    const meshIDs: Array<string> = sortBy(
      [...this.state.meshes.keys()], // sortBy can't sort a MapIterator
      [m => m.toLowerCase()],
    );
    const rows = meshIDs.map(meshID => (
      <WifiMeshRow
        enableDeviceEditing={true}
        key={meshID}
        gateways={this.state.meshes.get(meshID) || []}
        meshID={meshID}
        onDeleteMesh={this.onDeleteMesh}
        onDeleteDevice={this.onDeleteDevice}
      />
    ));

    return (
      <>
        <div className={this.props.classes.paper}>
          <div className={this.props.classes.header}>
            <Text variant="h5">Devices</Text>
            <div>
              <Tooltip title={'Last refreshed: ' + this.state.lastRefreshTime}>
                <span>
                  <IconButton
                    color="inherit"
                    onClick={this.fetchMeshes}
                    disabled={this.state.isLoading}>
                    <RefreshIcon />
                  </IconButton>
                </span>
              </Tooltip>
              <NestedRouteLink to="/add_mesh/">
                <Button>New Mesh</Button>
              </NestedRouteLink>
            </div>
          </div>
          <Text color="error">{this.state.errorMessage}</Text>
          <Paper elevation={2}>
            {this.state.isLoading ? <LinearProgress /> : null}
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell className={this.props.classes.infoColumn}>
                    Info
                  </TableCell>
                  <TableCell>ID</TableCell>
                  <TableCell className={this.props.classes.actionsColumn}>
                    Actions
                  </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>{rows}</TableBody>
            </Table>
          </Paper>
          <Route
            path={`${this.props.match.path}/add_mesh`}
            component={this.renderAddMeshDialog}
          />
          <Route
            path={`${this.props.match.path}/edit_mesh/:meshID`}
            component={this.renderEditMeshDialog}
          />
          <Route
            path={`${this.props.match.path}/add_device/:meshID`}
            component={this.renderAddDeviceDialog}
          />
          <Route
            path={`${this.props.match.path}/:meshID/edit_device/:deviceID`}
            component={this.renderEditDeviceDialog}
          />
        </div>
      </>
    );
  }

  fetchMeshes = () => {
    this.setState({isLoading: true});
    const networkId = nullthrows(this.props.match.params.networkId);
    Promise.all([
      MagmaV1API.getWifiByNetworkIdGateways({networkId}),
      MagmaV1API.getWifiByNetworkIdMeshes({networkId}),
    ])
      .then(([gatewaysResponse, meshesResponse]) => {
        const meshes = new Map();
        meshesResponse.forEach(meshID => meshes.set(meshID, []));

        const now = new Date().getTime();
        map(gatewaysResponse) // turn id->gateway map into gateway list
          .filter(gateway => gateway.device)
          .forEach(gatewayPayload => {
            const gateway = buildWifiGatewayFromPayloadV1(gatewayPayload, now);
            meshes.set(gateway.meshid, meshes.get(gateway.meshid) || []);
            nullthrows(meshes.get(gateway.meshid)).push(gateway);
          });

        meshes.forEach(gateways => gateways.sort(this.sortDevices));
        this.setState({
          isLoading: false,
          meshes: meshes,
          lastRefreshTime: new Date().toLocaleString(),
          errorMessage: null,
        });
      })
      .catch((error, _) =>
        this.setState({
          errorMessage: error.toString(),
          isLoading: false,
        }),
      );
  };

  renderAddMeshDialog = () => {
    return (
      <WifiMeshDialog onSave={this.onAddMesh} onCancel={this.onCancelDialog} />
    );
  };

  renderEditMeshDialog = () => {
    return (
      <WifiMeshDialog
        onSave={this.onCancelDialog}
        onCancel={this.onCancelDialog}
      />
    );
  };

  renderAddDeviceDialog = () => {
    return (
      <WifiDeviceDialog
        title="Add"
        onSave={this.onAddDevice}
        onCancel={this.onCancelDialog}
      />
    );
  };

  renderEditDeviceDialog = () => {
    return (
      <WifiDeviceDialog
        title="Edit"
        onSave={this.onEditDevice}
        onCancel={this.onCancelDialog}
      />
    );
  };

  onAddMesh = meshID => {
    if (meshID) {
      const {meshes} = this.state;
      meshes.set(meshID, []);
      this.setState({meshes});
      this.onCancelDialog();
    }
  };

  onDeleteMesh = meshID => {
    const {meshes} = this.state;
    meshes.delete(meshID);
    this.setState({meshes});
  };

  onAddDevice = device => {
    const {meshes} = this.state;
    const devices = nullthrows(meshes.get(device.meshid)).slice();
    devices.push(device);
    devices.sort(this.sortDevices);
    meshes.set(device.meshid, devices);
    this.setState({meshes});
    this.onCancelDialog();
  };

  onEditDevice = newDevice => {
    const {meshes} = this.state;
    const devices = nullthrows(meshes.get(newDevice.meshid)).map(oldDevice => {
      return oldDevice.id === newDevice.id ? newDevice : oldDevice;
    });
    devices.sort(this.sortDevices);
    meshes.set(newDevice.meshid, devices);
    this.setState({meshes});
    this.onCancelDialog();
  };

  onDeleteDevice = oldDevice => {
    const {meshes} = this.state;
    const devices = nullthrows(meshes.get(oldDevice.meshid)).filter(
      device => oldDevice.id !== device.id,
    );
    meshes.set(oldDevice.meshid, devices);
    this.setState({meshes});
  };

  onCancelDialog = () => this.props.history.push(`${this.props.match.url}/`);

  sortDevices = (d1, d2) =>
    d1.info.toLowerCase() > d2.info.toLowerCase() ? 1 : -1;
}

export default withStyles(styles)(withRouter(WifiMeshesDevicesTable));
