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
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import WACDeviceDialog from './WACDeviceDialog';
import WifiMeshDialog from '../wifi/WifiMeshDialog';
import WifiMeshRow from '../wifi/WifiMeshRow';
import axios from 'axios';
import {Route} from 'react-router-dom';

import nullthrows from '@fbcnms/util/nullthrows';
import {MagmaAPIUrls} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {buildWifiGatewayFromPayload, meshesURL} from '../wifi/WifiUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing(3),
  },
});

type Props = ContextRouter & WithStyles<typeof styles> & {};

type State = {
  groups: Map<string, any>,
  errorMessage: ?string,
};

class WACDevices extends React.Component<Props, State> {
  state = {
    groups: new Map(),
    errorMessage: null,
  };

  componentDidMount() {
    this.loadData();
  }

  async loadData() {
    try {
      const [devicesResponse, groupsResponse] = await axios.all([
        axios.get(MagmaAPIUrls.gateways(this.props.match, true)),
        axios.get(meshesURL(this.props.match)),
      ]);

      const groups = new Map();
      groupsResponse.data.forEach(groupID => groups.set(groupID, []));

      const now = new Date().getTime();
      devicesResponse.data
        // TODO: skip filter when magma API bug fixed t34643616
        .filter(device => device.record && device.config)
        .forEach(devicePayload => {
          const device = buildWifiGatewayFromPayload(devicePayload, now);
          groups.set(device.meshid, groups.get(device.meshid) || []);
          nullthrows(groups.get(device.meshid)).push(device);
        });

      this.setState({groups});
    } catch (error) {
      this.setState({errorMessage: error.toString()});
    }
  }

  render() {
    const rows = [];
    this.state.groups.forEach((devices, groupID) =>
      rows.push(
        <WifiMeshRow
          key={groupID}
          gateways={devices}
          meshID={groupID}
          onDeleteMesh={this.onDeleteGroup}
          // $FlowFixMe: Return types don't match. Please fix.
          onDeleteDevice={this.onDeleteDevice}
        />,
      ),
    );

    return (
      <>
        <div className={this.props.classes.paper}>
          <div className={this.props.classes.header}>
            <Text variant="h5">Devices</Text>
            <NestedRouteLink to="/add_group/">
              <Button>Add Group</Button>
            </NestedRouteLink>
          </div>
          <Paper elevation={2}>
            {this.state.errorMessage}
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Info</TableCell>
                  <TableCell>ID</TableCell>
                  <TableCell>Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>{rows}</TableBody>
            </Table>
          </Paper>
        </div>
        <Route
          path={`${this.props.match.path}/add_group`}
          component={this.renderAddGroupDialog}
        />
        <Route
          path={`${this.props.match.path}/edit_mesh/:meshID`}
          component={this.renderEditGroupDialog}
        />
        <Route
          path={`${this.props.match.path}/add_device/:meshID`}
          component={this.renderAddDeviceDialog}
        />
      </>
    );
  }

  renderAddGroupDialog = () => (
    <WifiMeshDialog onSave={this.onAddGroup} onCancel={this.onCancelDialog} />
  );

  renderEditGroupDialog = () => (
    <WifiMeshDialog
      onSave={this.onCancelDialog}
      onCancel={this.onCancelDialog}
    />
  );

  renderAddDeviceDialog = () => (
    <WACDeviceDialog onSave={this.onAddDevice} onCancel={this.onCancelDialog} />
  );

  onAddGroup = (groupID: string) => {
    const {groups} = this.state;
    groups.set(groupID, []);
    this.setState({groups});
    this.onCancelDialog();
  };

  onDeleteGroup = groupID => {
    const {groups} = this.state;
    groups.delete(groupID);
    this.setState({groups});
  };

  onAddDevice = device => {
    const {groups} = this.state;
    const devices = nullthrows(groups.get(device.meshid)).slice();
    devices.push(device);
    groups.set(device.meshid, devices);
    this.setState({groups});
    this.onCancelDialog();
  };

  onDeleteDevice = async oldDevice => {
    const {groups} = this.state;
    const devices = nullthrows(groups.get(oldDevice.meshid)).filter(
      device => oldDevice.id !== device.id,
    );
    groups.set(oldDevice.meshid, devices);
    this.setState({groups});

    const groupURL = meshesURL(this.props.match) + '/' + oldDevice.meshid;
    const groupResult = await axios.get(groupURL);
    const devicesPayload = (groupResult.data || []).filter(
      id => id !== oldDevice.id,
    );
    await axios.put(groupURL, devicesPayload);
  };

  onCancelDialog = () => this.props.history.push(`${this.props.match.url}/`);
}

export default withStyles(styles)(withRouter(WACDevices));
