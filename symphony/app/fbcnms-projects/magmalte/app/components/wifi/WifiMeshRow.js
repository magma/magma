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
import type {WifiGateway} from './WifiUtils';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import AddIcon from '@material-ui/icons/Add';
import Button from '@material-ui/core/Button';
import ChevronRight from '@material-ui/icons/ChevronRight';
import ClipboardLink from '@fbcnms/ui/components/ClipboardLink';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import EditIcon from '@material-ui/icons/Edit';
import ExpandMore from '@material-ui/icons/ExpandMore';
import IconButton from '@material-ui/core/IconButton';
import InfoIcon from '@material-ui/icons/Info';
import LinkIcon from '@material-ui/icons/Link';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import Tooltip from '@material-ui/core/Tooltip';
import url from 'url';
import {groupBy} from 'lodash';

import WifiDeviceDetails, {InfoRow} from './WifiDeviceDetails';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  actionsCell: {
    textAlign: 'right',
  },
  gatewayCell: {
    paddingBottom: '15px',
    paddingLeft: '75px',
    paddingRight: '15px',
    paddingTop: '15px',
  },
  deviceWarning: {
    color: 'red',
    paddingLeft: 40,
  },
  iconButton: {
    color: theme.palette.secondary.light,
    padding: '5px',
  },
  meshButton: {
    margin: 0,
    textTransform: 'none',
  },
  meshCell: {
    padding: '5px',
  },
  meshID: {
    color: theme.palette.primary.dark,
    fontWeight: 'bolder',
  },
  meshIconButton: {
    color: theme.palette.primary.dark,
    padding: '5px',
  },
  tableCell: {
    padding: '15px',
  },
  tableRow: {
    height: 'auto',
    whiteSpace: 'nowrap',
    verticalAlign: 'top',
  },
});

type Props = ContextRouter &
  WithStyles<typeof styles> &
  WithAlert & {
    enableDeviceEditing?: boolean,
    meshID: string,
    gateways: WifiGateway[],
    onDeleteMesh: string => void,
    onDeleteDevice: WifiGateway => void,
  };

const EXPANDED_STATE_TYPES = {
  none: 0,
  device: 1,
  neighbors: 2,
  fullDump: 3,
  configs: 4,
};

const MESH_ID_PARAM = 'meshID';
const DEVICE_ID_PARAM = 'deviceID';
const EXPANDED_STATE_PARAM = 'expandedState';

type State = {
  expanded: boolean,
  expandedGateways: {[key: string]: 0 | 1 | 2 | 3 | 4},
};

class WifiMeshRow extends React.Component<Props, State> {
  constructor(props) {
    super(props);
    const queryParams = new URLSearchParams(this.props.location.search);
    const initialState = {
      expanded: false,
      expandedGateways: {},
    };
    if (queryParams.get(MESH_ID_PARAM) === this.props.meshID) {
      initialState.expanded = true;
      const deviceID = queryParams.get(DEVICE_ID_PARAM);
      if (deviceID != null) {
        const expandedState = parseInt(queryParams.get(EXPANDED_STATE_PARAM));
        initialState.expandedGateways = {[deviceID]: expandedState};
      }
    }
    this.state = initialState;
  }

  handleToggleAllDevices = () => {
    const {gateways} = this.props;

    // determine old max state by getting max() of all the states
    const maxState = gateways
      .map(gateway => this.state.expandedGateways[gateway.id])
      .reduce((max, state) => (state ? Math.max(max, state) : max), 0);

    // calculate next state
    const nextState = (maxState + 1) % Object.keys(EXPANDED_STATE_TYPES).length;

    // assign same next state to all gateways
    if (nextState === 0) {
      // no need to set any gateways for 0/unexpanded state
      this.setState({expandedGateways: {}});
    } else {
      const expandedGateways = gateways
        .map(gateway => gateway.id)
        .reduce((expandedGateways, id) => {
          expandedGateways[id] = nextState;
          return expandedGateways;
        }, {});
      this.setState({expandedGateways});
    }
  };

  render() {
    const {meshID, gateways, classes} = this.props;

    let gatewayRows;
    if (this.state.expanded) {
      gatewayRows = gateways.map(gateway => (
        <TableRow className={this.props.classes.tableRow} key={gateway.id}>
          <TableCell className={classes.gatewayCell}>{gateway.info}</TableCell>
          <TableCell className={classes.tableCell}>
            {status}
            <DeviceStatusCircle
              isGrey={!gateway.status}
              isActive={!!gateway.up}
            />
            <Tooltip
              title="Click to toggle device info"
              enterDelay={400}
              placement={'right'}>
              <span onClick={() => this.expandGateway(gateway.id)}>
                {gateway.id}
              </span>
            </Tooltip>
            {gateway.coordinates.includes(NaN) && (
              <span className={classes.deviceWarning}>
                {' '}
                Please configure Lat/Lng
              </span>
            )}
            {gateway.status &&
              gateway.status.meta &&
              gateway.status.meta['validation_status'] !== 'passed' && (
                <span className={classes.deviceWarning}>
                  {' '}
                  Please check image validation status
                </span>
              )}

            {!!this.state.expandedGateways[gateway.id] && gateway.status && (
              <WifiDeviceDetails
                device={gateway}
                hideHeader={true}
                showConfigs={
                  this.state.expandedGateways[gateway.id] ===
                  EXPANDED_STATE_TYPES.configs
                }
                showDevice={
                  this.state.expandedGateways[gateway.id] ===
                  EXPANDED_STATE_TYPES.device
                }
                showNeighbors={
                  this.state.expandedGateways[gateway.id] ===
                  EXPANDED_STATE_TYPES.neighbors
                }
                showFullDump={
                  this.state.expandedGateways[gateway.id] ===
                  EXPANDED_STATE_TYPES.fullDump
                }
              />
            )}
          </TableCell>
          <TableCell className={classes.actionsCell}>
            <ClipboardLink title="Copy link to this device">
              {({copyString}) => (
                <IconButton
                  className={classes.iconButton}
                  onClick={() =>
                    copyString(this.buildLinkURL(meshID, gateway.id))
                  }>
                  <LinkIcon />
                </IconButton>
              )}
            </ClipboardLink>
            <Tooltip title="Click to toggle device info" enterDelay={400}>
              <IconButton
                className={classes.iconButton}
                onClick={() => this.expandGateway(gateway.id)}>
                <InfoIcon />
              </IconButton>
            </Tooltip>
            {this.props.enableDeviceEditing && (
              <NestedRouteLink to={`/${meshID}/edit_device/${gateway.id}`}>
                <IconButton className={classes.iconButton}>
                  <EditIcon />
                </IconButton>
              </NestedRouteLink>
            )}
            <IconButton
              className={classes.iconButton}
              onClick={() => this.showDeviceDeleteDialog(gateway)}>
              <DeleteIcon />
            </IconButton>
          </TableCell>
        </TableRow>
      ));
    }

    // construct version list per mesh
    const versionGroups: {string: Array<WifiGateway>} = groupBy(
      gateways,
      device => {
        if (device.versionParsed) {
          if (device.versionParsed.fbpkg !== 'none') {
            return device.versionParsed.fbpkg;
          } else {
            return device.versionParsed.hash;
          }
        }
        return device.version || 'UNKNOWN';
      },
    );

    // sort by device count, then version string
    const sortedVersions: Array<string> = Object.keys(versionGroups);
    sortedVersions.sort((a, b) => {
      // keep "Not Reported at the bottom"
      if (a === 'Not Reported') {
        return 1;
      } else if (b === 'Not Reported') {
        return -1;
      } else if (versionGroups[a].length === versionGroups[b].length) {
        // if device counts are equal, then use version string
        return a.localeCompare(b);
      } else {
        // sort by device count
        return versionGroups[b].length - versionGroups[a].length;
      }
    });

    const gatewayVersions = sortedVersions.map(version => (
      <div key={version}>
        <Tooltip
          title={`${versionGroups[version].length} device(s) with ${versionGroups[version][0].version}`}
          enterDelay={100}
          key={version}>
          <span style={{fontFamily: 'monospace'}}>{version}</span>
        </Tooltip>
        :{' '}
        <span style={{fontSize: '88%', fontWeight: 'bold'}}>
          {versionGroups[version].length}
        </span>
      </div>
    ));

    return (
      <>
        <TableRow className={this.props.classes.tableRow}>
          <TableCell className={this.props.classes.meshCell}>
            <IconButton
              className={classes.meshIconButton}
              onClick={gateways.length == 0 ? null : this.onToggleExpand}>
              {this.state.expanded ? <ExpandMore /> : <ChevronRight />}
            </IconButton>
            <span className={classes.meshID}>{meshID}</span>
          </TableCell>
          <TableCell className={this.props.classes.meshCell}>
            {gateways.length > 0 && (
              <>
                <InfoRow
                  label="Up"
                  data={`${gateways.filter(gateway => gateway.up).length} of ${
                    gateways.length
                  }`}
                />
                {gatewayVersions}
                {this.state.expanded && (
                  <>
                    <Tooltip
                      title="Click to toggle device info"
                      enterDelay={400}
                      placement={'right'}>
                      <Button
                        size="small"
                        className={classes.meshButton}
                        onClick={this.handleToggleAllDevices}>
                        toggle info
                      </Button>
                    </Tooltip>
                  </>
                )}
              </>
            )}
          </TableCell>
          <TableCell className={this.props.classes.actionsCell}>
            <ClipboardLink title="Copy link to this mesh">
              {({copyString}) => (
                <IconButton
                  className={classes.iconButton}
                  onClick={() => copyString(this.buildLinkURL(meshID))}>
                  <LinkIcon />
                </IconButton>
              )}
            </ClipboardLink>
            <NestedRouteLink to={`/add_device/${meshID}`}>
              <IconButton className={classes.meshIconButton}>
                <AddIcon />
              </IconButton>
            </NestedRouteLink>
            <NestedRouteLink to={`/edit_mesh/${meshID}`}>
              <IconButton className={classes.meshIconButton}>
                <EditIcon />
              </IconButton>
            </NestedRouteLink>
            <IconButton
              className={classes.meshIconButton}
              onClick={this.showMeshDeleteDialog}>
              <DeleteIcon />
            </IconButton>
          </TableCell>
        </TableRow>
        {gatewayRows}
      </>
    );
  }

  onToggleExpand = () => this.setState({expanded: !this.state.expanded});
  expandGateway = id => {
    const expandedGateways = {
      ...this.state.expandedGateways,
      [id]:
        ((this.state.expandedGateways[id] | 0) + 1) %
        Object.keys(EXPANDED_STATE_TYPES).length,
    };

    this.setState({expandedGateways});
  };

  showMeshDeleteDialog = () => {
    this.props
      .confirm(
        `Are you sure you want to delete mesh "${this.props.meshID}" and all its devices (count: ${this.props.gateways.length})?`,
      )
      .then(async confirmed => {
        if (!confirmed) {
          return;
        }

        await Promise.all(
          this.props.gateways.map(device =>
            MagmaV1API.deleteWifiByNetworkIdGatewaysByGatewayId({
              networkId: nullthrows(this.props.match.params.networkId),
              gatewayId: device.id,
            }),
          ),
        );

        await MagmaV1API.deleteWifiByNetworkIdMeshesByMeshId({
          networkId: nullthrows(this.props.match.params.networkId),
          meshId: this.props.meshID,
        });

        this.props.onDeleteMesh(this.props.meshID);
      });
  };

  showDeviceDeleteDialog = (device: WifiGateway) => {
    this.props
      .confirm(`Are you sure you want to delete "${device.id}"?`)
      .then(async confirmed => {
        if (!confirmed) {
          return;
        }

        // V1 API call will delete all parts of the device
        await MagmaV1API.deleteWifiByNetworkIdGatewaysByGatewayId({
          networkId: nullthrows(this.props.match.params.networkId),
          gatewayId: device.id,
        });

        this.props.onDeleteDevice(device);
      });
  };

  buildLinkURL = (meshID: string, deviceID: ?string = null): string => {
    const query: {[string]: string | number} = {[MESH_ID_PARAM]: meshID};
    if (deviceID) {
      query[DEVICE_ID_PARAM] = deviceID;
      query[EXPANDED_STATE_PARAM] = this.state.expandedGateways[deviceID] ?? 1;
    }
    const {protocol, host, pathname} = window.location;
    return url.format({protocol, host, pathname, query});
  };
}

export default withStyles(styles)(withRouter(withAlert(WifiMeshRow)));
