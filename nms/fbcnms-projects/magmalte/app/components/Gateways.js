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
import type {Gateway, GatewayPayload} from './GatewayUtils';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import AddGatewayDialog from './AddGatewayDialog';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditGatewayDialog from './EditGatewayDialog';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {GatewayStatus} from './GatewayUtils';
import {MagmaAPIUrls} from '../common/MagmaAPI';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const myInt = (n: ?(string | number)): ?number => {
  return n ? parseInt(n) : null;
};

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
  gateways: ?(Gateway[]),
  editingGateway: ?any,
};

class Gateways extends React.Component<Props, State> {
  state = {
    showDialog: false,
    gateways: null,
    editingGateway: null,
  };

  componentDidMount() {
    const {match} = this.props;
    axios
      .get(MagmaAPIUrls.gateways(match, true))
      .then(response => {
        const gateways = response.data
          .filter(g => g.record && g.config)
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
          <GatewayStatus
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
          <Typography variant="h5">Configure Gateways</Typography>
          <Button variant="contained" color="primary" onClick={this.showDialog}>
            Add Gateway
          </Button>
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
        <AddGatewayDialog
          open={this.state.showDialog}
          onClose={this.hideDialog}
          onSave={this.onSave}
        />
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
        axios
          .delete(MagmaAPIUrls.gateway(match, gateway.logicalID))
          .then(_resp =>
            this.setState({
              gateways: gateways.filter(
                gw => gw.logicalID != gateway.logicalID,
              ),
            }),
          );
      });
  };

  _buildGatewayFromPayload(gateway: GatewayPayload): Gateway {
    if (!gateway.record || !gateway.config) {
      throw Error('Cannot read gateway without `record` or `config`');
    }

    let enodebRFTXOn = false;
    let enodebConnected = false;
    let gpsConnected = false;
    let mmeConnected = false;
    let version = 'Not Reported';
    let vpnIP = 'Not Reported';
    let lastCheckin = 'Not Reported';
    let isBackhaulDown = true;
    const latLon = {lat: 0, lon: 0};
    const {status} = gateway;
    if (status) {
      version = status.version || version;
      vpnIP = status.vpn_ip || vpnIP;
      lastCheckin = status.checkin_time
        ? status.checkin_time.toString()
        : lastCheckin;

      // if the last check-in time is more than 5 minutes
      // we treat it as backhaul is down
      const dutation = Math.max(0, Date.now() - status.checkin_time);
      isBackhaulDown = dutation > 1000 * 5 * 60;

      if (status.meta) {
        if (!isBackhaulDown) {
          enodebRFTXOn = status.meta && status.meta.rf_tx_on;
        }

        latLon.lat = status.meta.gps_latitude;
        latLon.lon = status.meta.gps_longitude;
        gpsConnected = status.meta.gps_connected == 1;
        enodebConnected = status.meta.enodeb_connected == 1;
        mmeConnected = status.meta.mme_connected == 1;
      }
    }

    let autoupgradePollInterval,
      checkinInterval,
      checkinTimeout,
      tier,
      autoupgradeEnabled;
    if (gateway.config && gateway.config.magmad_gateway) {
      const {magmad_gateway: magmadGateway} = gateway.config;
      autoupgradePollInterval = myInt(magmadGateway.autoupgrade_poll_interval);
      checkinInterval = myInt(magmadGateway.checkin_interval);
      checkinTimeout = myInt(magmadGateway.checkin_timeout);
      tier = magmadGateway.tier;
      autoupgradeEnabled = magmadGateway.autoupgrade_enabled;
    }

    let ipBlock, natEnabled;
    let attachedEnodebSerials = [];
    let pci, transmitEnabled;
    let control, csfbRAT, csfbMCC, csfbMNC, lac;
    if (gateway.config && gateway.config.cellular_gateway) {
      const {cellular_gateway: cellularGateway} = gateway.config;
      ipBlock = (cellularGateway.epc || {}).ip_block;
      natEnabled = (cellularGateway.epc || {}).nat_enabled;

      attachedEnodebSerials = cellularGateway.attached_enodeb_serials || [];

      pci = myInt((cellularGateway.ran || {}).pci);
      transmitEnabled = (cellularGateway.ran || {}).transmit_enabled || false;

      const nonEPSService = cellularGateway.non_eps_service || {};
      control = myInt(nonEPSService.non_eps_service_control);
      csfbRAT = myInt(nonEPSService.csfb_rat);
      csfbMNC = myInt(nonEPSService.csfb_mnc);
      csfbMCC = myInt(nonEPSService.csfb_mcc);
      lac = myInt(nonEPSService.lac);
    }

    return {
      hardware_id: gateway.record.hardware_id,
      name: gateway.name || 'N/A',
      logicalID: gateway.gateway_id,
      challengeType: gateway.record.key.key_type,
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
      tier,
      autoupgradeEnabled: !!autoupgradeEnabled,
      attachedEnodebSerials,
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
