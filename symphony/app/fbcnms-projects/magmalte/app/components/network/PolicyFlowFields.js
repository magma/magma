/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';
import type {flow_description} from '@fbcnms/magma-api';

import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import FormControl from '@material-ui/core/FormControl';
import IconButton from '@material-ui/core/IconButton';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import RemoveCircleOutline from '@material-ui/icons/RemoveCircleOutline';
import Select from '@material-ui/core/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';

import {ACTION, DIRECTION, PROTOCOL} from './PolicyTypes';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  root: {
    '&$expanded': {
      minHeight: 'auto',
    },
  },
  expanded: {},
  block: {
    display: 'block',
  },
  flex: {display: 'flex'},
  panel: {flexGrow: 1},
  removeIcon: {alignSelf: 'baseline'},
};

type ActionType = $Keys<typeof ACTION>;
type Props = WithStyles<typeof styles> & {
  index: number,
  flow: flow_description,
  handleActionChange: (number, ActionType) => void,
  handleFieldChange: (number, string, string | number) => void,
  handleDelete: number => void,
};

class PolicyFlowFields extends React.Component<Props> {
  render() {
    const {classes, flow} = this.props;
    return (
      <div className={classes.flex}>
        <ExpansionPanel className={classes.panel}>
          <ExpansionPanelSummary
            classes={{root: classes.root, expanded: classes.expanded}}
            expandIcon={<ExpandMoreIcon />}>
            <Text variant="body2">Flow {this.props.index + 1}</Text>
          </ExpansionPanelSummary>
          <ExpansionPanelDetails classes={{root: classes.block}}>
            <div className={classes.flex}>
              <FormControl className={classes.input}>
                <InputLabel htmlFor="action">Action</InputLabel>
                <Select
                  value={flow.action}
                  onChange={this.handleActionChange}
                  input={<Input id="action" />}>
                  <MenuItem value={ACTION.PERMIT}>Permit</MenuItem>
                  <MenuItem value={ACTION.DENY}>Deny</MenuItem>
                </Select>
              </FormControl>
              <FormControl className={classes.input}>
                <InputLabel htmlFor="direction">Direction</InputLabel>
                <Select
                  value={flow.match.direction}
                  onChange={this.handleDirectionChange}
                  input={<Input id="direction" />}>
                  <MenuItem value={DIRECTION.UPLINK}>Uplink</MenuItem>
                  <MenuItem value={DIRECTION.DOWNLINK}>Downlink</MenuItem>
                </Select>
              </FormControl>
              <FormControl className={classes.input}>
                <InputLabel htmlFor="protocol">Protocol</InputLabel>
                <Select
                  value={flow.match.ip_proto}
                  onChange={this.handleProtocolChange}
                  input={<Input id="protocol" />}>
                  <MenuItem value={PROTOCOL.IPPROTO_IP}>IP</MenuItem>
                  <MenuItem value={PROTOCOL.IPPROTO_UDP}>UDP</MenuItem>
                  <MenuItem value={PROTOCOL.IPPROTO_TCP}>TCP</MenuItem>
                  <MenuItem value={PROTOCOL.IPPROTO_ICMP}>ICMP</MenuItem>
                </Select>
              </FormControl>
            </div>
            {flow.match.ip_proto !== PROTOCOL.IPPROTO_ICMP && (
              <div className={classes.flex}>
                <TextField
                  className={classes.input}
                  label="IPv4 Source"
                  margin="normal"
                  value={flow.match.ipv4_src}
                  onChange={this.handleIPv4SourceChange}
                />
                <TextField
                  className={classes.input}
                  label="IPv4 Destination"
                  margin="normal"
                  value={flow.match.ipv4_dst}
                  onChange={this.handleIPv4DestinationChange}
                />
              </div>
            )}
            {flow.match.ip_proto === PROTOCOL.IPPROTO_UDP && (
              <div className={classes.flex}>
                <TextField
                  className={classes.input}
                  label="UDP Source Port"
                  margin="normal"
                  value={flow.match.udp_src}
                  onChange={this.handleUDPSourceChange}
                />
                <TextField
                  className={classes.input}
                  label="UDP Destination Port"
                  margin="normal"
                  value={flow.match.udp_dst}
                  onChange={this.handleUDPDestinationChange}
                />
              </div>
            )}
            {flow.match.ip_proto === PROTOCOL.IPPROTO_TCP && (
              <div className={classes.flex}>
                <TextField
                  className={classes.input}
                  label="TCP Source Port"
                  margin="normal"
                  value={flow.match.tcp_src}
                  onChange={this.handleTCPSourceChange}
                />
                <TextField
                  className={classes.input}
                  label="TCP Destination Port"
                  margin="normal"
                  value={flow.match.tcp_dst}
                  onChange={this.handleTCPDestinationChange}
                />
              </div>
            )}
          </ExpansionPanelDetails>
        </ExpansionPanel>
        <IconButton className={classes.removeIcon} onClick={this.handleDelete}>
          <RemoveCircleOutline />
        </IconButton>
      </div>
    );
  }

  handleActionChange = ({target}) =>
    this.props.handleActionChange(
      this.props.index,
      // eslint-disable-next-line flowtype/no-weak-types
      ((target.value: any): ActionType),
    );
  handleDirectionChange = ({target}) =>
    this.props.handleFieldChange(this.props.index, 'direction', target.value);
  handleProtocolChange = ({target}) =>
    this.props.handleFieldChange(this.props.index, 'ip_proto', target.value);
  handleIPv4SourceChange = ({target}) =>
    this.props.handleFieldChange(this.props.index, 'ipv4_src', target.value);
  handleIPv4DestinationChange = ({target}) =>
    this.props.handleFieldChange(this.props.index, 'ipv4_dst', target.value);
  handleUDPSourceChange = ({target}) =>
    this.props.handleFieldChange(
      this.props.index,
      'udp_src',
      parseInt(target.value),
    );
  handleUDPDestinationChange = ({target}) =>
    this.props.handleFieldChange(
      this.props.index,
      'udp_dst',
      parseInt(target.value),
    );
  handleTCPSourceChange = ({target}) =>
    this.props.handleFieldChange(
      this.props.index,
      'tcp_src',
      parseInt(target.value),
    );
  handleTCPDestinationChange = ({target}) =>
    this.props.handleFieldChange(
      this.props.index,
      'tcp_dst',
      parseInt(target.value),
    );
  handleDelete = () => this.props.handleDelete(this.props.index);
}

export default withStyles(styles)(PolicyFlowFields);
