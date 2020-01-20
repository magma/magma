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
import type {GatewayV1} from './GatewayUtils';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';
import type {lte_gateway} from '@fbcnms/magma-api';

import AddGatewayDialog from './AddGatewayDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import EditGatewayDialog from './EditGatewayDialog';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {MAGMAD_DEFAULT_CONFIGS} from './AddGatewayDialog';
import {find} from 'lodash';
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

type Props = ContextRouter & WithAlert & WithStyles<typeof styles> & {};

type State = {
  showDialog: boolean,
  gateways: ?(GatewayV1[]),
  editingGateway: ?GatewayV1,
};

class Gateways extends React.Component<Props, State> {
  state = {
    showDialog: false,
    gateways: null,
    editingGateway: null,
  };

  componentDidMount() {
    const {match} = this.props;
    MagmaV1API.getLteByNetworkIdGateways({
      networkId: nullthrows(match.params.networkId),
    })
      .then(response => {
        const gateways = Object.keys(response)
          .map(k => response[k])
          .filter(g => g.cellular)
          .map(this._buildGatewayFromPayload);
        this.setState({gateways});
      })
      .catch(error => {
        this.props.alert(error);
      });
  }

  render() {
    const {gateways} = this.state;
    const rows = (gateways || []).map(gateway => (
      <TableRow key={gateway.logicalID}>
        <TableCell>
          <DeviceStatusCircle
            isGrey={!gateway.enodebRFTXOn}
            isActive={gateway.enodebRFTXOn === gateway.enodebRFTXEnabled}
          />
          {gateway.name}
        </TableCell>
        <TableCell>{gateway.hardware_id}</TableCell>
        <TableCell>
          <IconButton
            data-testid="edit-gateway-icon"
            color="primary"
            onClick={this.editGateway.bind(this, gateway)}>
            <EditIcon />
          </IconButton>
          <IconButton
            data-testid="delete-gateway-icon"
            color="primary"
            onClick={this.deleteGateway.bind(this, gateway)}>
            <DeleteIcon />
          </IconButton>
        </TableCell>
      </TableRow>
    ));

    return (
      <div className={this.props.classes.paper}>
        <div className={this.props.classes.header}>
          <Text variant="h5">Configure Gateways</Text>
          <Button onClick={this.showDialog}>Add Gateway</Button>
        </div>
        <Paper elevation={2}>
          {gateways ? (
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Name</TableCell>
                  <TableCell>Hardware UUID</TableCell>
                  <TableCell />
                </TableRow>
              </TableHead>
              <TableBody>{rows}</TableBody>
            </Table>
          ) : (
            <LoadingFiller />
          )}
        </Paper>
        {this.state.showDialog && (
          <AddGatewayDialog
            onClose={this.hideDialog}
            onSave={this.onGatewayAdd}
          />
        )}
        <EditGatewayDialog
          key={this.state.editingGateway && this.state.editingGateway.logicalID}
          gateway={this.state.editingGateway}
          onClose={() => this.setState({editingGateway: null})}
          onSave={this.onSave}
        />
      </div>
    );
  }

  showDialog = () => this.setState({showDialog: true});
  hideDialog = () => this.setState({showDialog: false});
  editGateway = editingGateway => this.setState({editingGateway});

  onGatewayAdd = async ({
    gatewayID,
    name,
    description,
    hardwareID,
    challengeKey,
    tier,
  }) => {
    const networkID = nullthrows(this.props.match.params.networkId);
    await MagmaV1API.postLteByNetworkIdGateways({
      networkId: networkID,
      gateway: {
        id: gatewayID,
        name,
        description,
        cellular: {
          epc: {nat_enabled: true, ip_block: '192.168.128.0/24'},
          ran: {pci: 260, transmit_enabled: false},
          non_eps_service: undefined,
        },
        magmad: MAGMAD_DEFAULT_CONFIGS,
        device: {
          hardware_id: hardwareID,
          key: {
            key: challengeKey,
            key_type: 'SOFTWARE_ECDSA_SHA256', // default key/challenge type
          },
        },
        connected_enodeb_serials: [],
        tier,
      },
    });
    const gatewayPayload = await MagmaV1API.getLteByNetworkIdGatewaysByGatewayId(
      {
        networkId: networkID,
        gatewayId: gatewayID,
      },
    );
    this.onSave(gatewayPayload);
  };

  onSave = gatewayPayload => {
    const gateway = this._buildGatewayFromPayload(gatewayPayload);
    const newGateways = nullthrows(this.state.gateways).slice(0);
    if (this.state.editingGateway) {
      newGateways[newGateways.indexOf(this.state.editingGateway)] = gateway;
    } else {
      newGateways.push(gateway);
    }
    this.setState({
      gateways: newGateways,
      showDialog: false,
      editingGateway: null,
    });
  };

  deleteGateway = gateway => {
    const gateways = nullthrows(this.state.gateways);
    const {match} = this.props;
    this.props
      .confirm(`Are you sure you want to delete ${gateway.name}?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        MagmaV1API.deleteLteByNetworkIdGatewaysByGatewayId({
          networkId: nullthrows(match.params.networkId),
          gatewayId: gateway.logicalID,
        }).then(() =>
          this.setState({
            gateways: gateways.filter(gw => gw.logicalID != gateway.logicalID),
          }),
        );
      });
  };

  _buildGatewayFromPayload(gateway: lte_gateway): GatewayV1 {
    let enodebRFTXOn = false;
    let enodebConnected = false;
    let gpsConnected = false;
    let mmeConnected = false;
    let version = 'Not Reported';
    let vpnIP = 'Not Reported';
    let lastCheckin = 'Not Reported';
    let hardwareID = 'Not reported';
    let isBackhaulDown = true;
    const latLon = {lat: 0, lon: 0};
    const {status} = gateway;
    if (status) {
      vpnIP = status.platform_info?.vpn_ip || vpnIP;
      const packages = find(status.platform_info?.packages || [], {
        name: 'magma',
      });
      version = packages?.version || '';
      // if the last check-in time is more than 5 minutes
      // we treat it as backhaul is down
      const checkin = status.checkin_time;
      if (checkin != null) {
        const duration = Math.max(0, Date.now() - checkin);
        isBackhaulDown = duration > 1000 * 5 * 60;
        lastCheckin = checkin.toString();
      }

      const {meta} = status;
      if (meta) {
        if (!isBackhaulDown) {
          enodebRFTXOn = status.meta && status.meta.rf_tx_on;
        }

        latLon.lat = parseFloat(meta.gps_latitude);
        latLon.lon = parseFloat(meta.gps_longitude);
        gpsConnected = meta.gps_connected == '1';
        enodebConnected = meta.enodeb_connected == '1';
        mmeConnected = meta.mme_connected == '1';
      }

      if (status.hardware_id) {
        hardwareID = status.hardware_id;
      }
    }

    let autoupgradePollInterval,
      checkinInterval,
      checkinTimeout,
      autoupgradeEnabled;
    if (gateway.magmad) {
      const {magmad} = gateway;
      autoupgradePollInterval = magmad.autoupgrade_poll_interval;
      checkinInterval = magmad.checkin_interval;
      checkinTimeout = magmad.checkin_timeout;
      autoupgradeEnabled = magmad.autoupgrade_enabled;
    }

    let ipBlock, natEnabled;
    let pci, transmitEnabled;
    let control, csfbRAT, csfbMCC, csfbMNC, lac;
    if (gateway.cellular) {
      const {cellular} = gateway;
      ipBlock = cellular?.epc?.ip_block;
      natEnabled = cellular?.epc?.nat_enabled;
      pci = cellular?.ran?.pci;
      transmitEnabled = cellular?.ran?.transmit_enabled ?? false;

      const nonEPSService = cellular.non_eps_service || {};
      control = nonEPSService.non_eps_service_control;
      csfbRAT = nonEPSService.csfb_rat;
      csfbMNC = nonEPSService.csfb_mnc;
      csfbMCC = nonEPSService.csfb_mcc;
      lac = nonEPSService.lac;
    }

    return {
      hardware_id: hardwareID,
      name: gateway.name || 'N/A',
      logicalID: gateway.id,
      challengeType: gateway?.device?.key?.key_type || '',
      enodebRFTXEnabled: !!transmitEnabled,
      enodebRFTXOn: !!enodebRFTXOn,
      enodebConnected,
      gpsConnected,
      isBackhaulDown,
      lastCheckin,
      latLon,
      mmeConnected,
      version,
      vpnIP,
      autoupgradePollInterval,
      checkinInterval,
      checkinTimeout,
      tier: gateway.tier,
      autoupgradeEnabled: !!autoupgradeEnabled,
      attachedEnodebSerials: gateway.connected_enodeb_serials || [],
      ran: {
        pci,
        transmitEnabled: !!transmitEnabled,
      },
      epc: {
        ipBlock: ipBlock || '',
        natEnabled: natEnabled || false,
      },
      nonEPSService: {
        control: control || 0,
        csfbRAT: csfbRAT || 0,
        csfbMCC,
        csfbMNC,
        lac,
      },
      rawGateway: gateway,
    };
  }
}

export default withStyles(styles)(withAlert(withRouter(Gateways)));
