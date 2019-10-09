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
import type {WithStyles} from '@material-ui/core';

import Button from '@material-ui/core/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Divider from '@material-ui/core/Divider';
import FormField from './FormField';
import Input from '@material-ui/core/Input';
import MagmaV1API from '../common/MagmaV1API';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import Typography from '@material-ui/core/Typography';

import nullthrows from '@fbcnms/util/nullthrows';
import {toString} from './GatewayUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  input: {
    margin: '10px 0',
    width: '100%',
  },
  title: {
    fontSize: '15px',
  },
  divider: {
    margin: '10px 0',
  },
});

type Props = ContextRouter &
  WithStyles<typeof styles> & {
    onClose: () => void,
    onSave: (gatewayID: string) => void,
    gateway: GatewayV1,
  };

type State = {
  natEnabled: boolean,
  ipBlock: string,
  attachedEnodebSerials: Array<string>,
  pci: string,
  transmitEnabled: boolean,
  nonEPSServiceControl: number,
  csfbRAT: number,
  mcc: string,
  mnc: string,
  lac: string,
};

class GatewayCellularFields extends React.Component<Props, State> {
  state = {
    natEnabled: this.props.gateway.epc.natEnabled,
    ipBlock: this.props.gateway.epc.ipBlock,
    attachedEnodebSerials: this.props.gateway.attachedEnodebSerials,
    pci: toString(this.props.gateway.ran.pci),
    transmitEnabled: this.props.gateway.ran.transmitEnabled,
    nonEPSServiceControl: this.props.gateway.nonEPSService.control,
    csfbRAT: this.props.gateway.nonEPSService.csfbRAT,
    mcc: toString(this.props.gateway.nonEPSService.csfbMCC),
    mnc: toString(this.props.gateway.nonEPSService.csfbMNC),
    lac: toString(this.props.gateway.nonEPSService.lac),
  };

  render() {
    const nonEPSServiceControlOff = this.state.nonEPSServiceControl == 0;
    return (
      <>
        <DialogContent>
          <Typography className={this.props.classes.title} variant="h6">
            EPC Configs
          </Typography>
          <FormField label="NAT Enabled">
            <Select
              className={this.props.classes.input}
              value={this.state.natEnabled ? 1 : 0}
              onChange={this.natEnabledChanged}>
              <MenuItem value={1}>Enabled</MenuItem>
              <MenuItem value={0}>Disabled</MenuItem>
            </Select>
          </FormField>
          <FormField label="IP Block">
            <Input
              className={this.props.classes.input}
              value={this.state.ipBlock}
              onChange={this.ipBlockChanged}
              placeholder="E.g. 20.20.20.0/24"
              disabled={this.state.natEnabled}
            />
          </FormField>
          <Divider className={this.props.classes.divider} />
          <Typography className={this.props.classes.title} variant="h6">
            RAN Configs
          </Typography>
          <FormField label="Registered eNB IDs">
            <Input
              className={this.props.classes.input}
              value={this.state.attachedEnodebSerials.toString()}
              onChange={this.attachedEnodebSerialsChanged}
              placeholder="E.g. 123, 456"
            />
          </FormField>
          <FormField label="PCI">
            <Input
              className={this.props.classes.input}
              value={this.state.pci}
              onChange={this.pciChanged}
              placeholder="E.g. 123"
            />
          </FormField>
          <FormField label="ENODEB Transmit Enabled">
            <Select
              className={this.props.classes.input}
              value={this.state.transmitEnabled ? 1 : 0}
              onChange={this.transmitEnabledChanged}>
              <MenuItem value={1}>Enabled</MenuItem>
              <MenuItem value={0}>Disabled</MenuItem>
            </Select>
          </FormField>
          <Divider className={this.props.classes.divider} />
          <Typography className={this.props.classes.title} variant="h6">
            NonEPS Configs
          </Typography>
          <FormField label="NonEPS Service Control">
            <Select
              className={this.props.classes.input}
              value={this.state.nonEPSServiceControl}
              onChange={this.nonEPSServiceControlChanged}>
              <MenuItem value={0}>Off</MenuItem>
              <MenuItem value={1}>CSFB SMS</MenuItem>
              <MenuItem value={2}>SMS</MenuItem>
            </Select>
          </FormField>
          <FormField label="CSFB RAT Type">
            <Select
              disabled={nonEPSServiceControlOff}
              className={this.props.classes.input}
              value={this.state.csfbRAT}
              onChange={this.csfbRATChanged}>
              <MenuItem value={0}>2G</MenuItem>
              <MenuItem value={1}>3G</MenuItem>
            </Select>
          </FormField>
          <FormField label="CSFB MCC">
            <Input
              disabled={nonEPSServiceControlOff}
              className={this.props.classes.input}
              value={this.state.mcc}
              onChange={this.mccChanged}
              placeholder="E.g. 01"
            />
          </FormField>
          <FormField label="CSFB MNC">
            <Input
              disabled={nonEPSServiceControlOff}
              className={this.props.classes.input}
              value={this.state.mnc}
              onChange={this.mncChanged}
              placeholder="E.g. 01"
            />
          </FormField>
          <FormField label="LAC">
            <Input
              disabled={nonEPSServiceControlOff}
              className={this.props.classes.input}
              value={this.state.lac}
              onChange={this.lacChanged}
              placeholder="E.g. 01"
            />
          </FormField>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onClose} color="primary">
            Cancel
          </Button>
          <Button onClick={this.onSave} color="primary" variant="contained">
            Save
          </Button>
        </DialogActions>
      </>
    );
  }

  onSave = () => {
    const id = this.props.gateway.logicalID;
    const {cellular} = this.props.gateway.rawGateway;
    const {nonEPSServiceControl, csfbRAT} = this.state;

    // these conditions should never be true since these values are coming from
    // a selector, but they're needed for Flow
    if (
      nonEPSServiceControl !== 0 &&
      nonEPSServiceControl !== 1 &&
      nonEPSServiceControl !== 2
    ) {
      return;
    }

    if (csfbRAT !== 1 && csfbRAT !== 0) {
      return;
    }

    const config = {
      ...cellular,
      epc: {
        ...cellular.epc,
        nat_enabled: this.state.natEnabled,
        ip_block: this.state.ipBlock,
      },
      ran: {
        ...cellular.ran,
        pci: parseInt(this.state.pci),
        transmit_enabled: this.state.transmitEnabled,
      },
      non_eps_service: {
        ...cellular.non_eps_service,
        non_eps_service_control: nonEPSServiceControl,
        csfb_rat: csfbRAT,
        csfb_mcc: this.state.mcc,
        csfb_mnc: this.state.mnc,
        lac: parseInt(this.state.lac),
      },
      // Override the registered eNodeB devices with new values
      attached_enodeb_serials: this.state.attachedEnodebSerials,
    };

    const {match} = this.props;
    MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdCellular({
      networkId: nullthrows(match.params.networkId),
      gatewayId: id,
      config,
    }).then(() => this.props.onSave(id));
  };

  natEnabledChanged = ({target}) => this.setState({natEnabled: !!target.value});
  ipBlockChanged = ({target}) => this.setState({ipBlock: target.value});
  pciChanged = ({target}) => this.setState({pci: target.value});
  transmitEnabledChanged = ({target}) =>
    this.setState({transmitEnabled: !!target.value});
  nonEPSServiceControlChanged = ({target}) =>
    this.setState({nonEPSServiceControl: parseInt(target.value)});
  csfbRATChanged = ({target}) =>
    this.setState({csfbRAT: parseInt(target.value)});
  mccChanged = ({target}) => this.setState({mcc: target.value});
  mncChanged = ({target}) => this.setState({mnc: target.value});
  lacChanged = ({target}) => this.setState({lac: target.value});
  attachedEnodebSerialsChanged = ({target}) => {
    this.setState({
      attachedEnodebSerials: target.value.replace(' ', '').split(','),
    });
  };
}

export default withStyles(styles)(withRouter(GatewayCellularFields));
