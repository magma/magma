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
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';
import type {Gateway} from './GatewayUtils';

import axios from 'axios';
import {MagmaAPIUrls} from '../common/MagmaAPI';
import Button from '@material-ui/core/Button';
import Check from '@material-ui/icons/Check';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import React from 'react';
import Divider from '@material-ui/core/Divider';
import Fade from '@material-ui/core/Fade';
import Typography from '@material-ui/core/Typography';
import FormField from './FormField';
import Input from '@material-ui/core/Input';
import grey from '@material-ui/core/colors/grey';
import LinearProgress from '@material-ui/core/LinearProgress';

import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  input: {
    margin: '10px 0',
    width: '100%',
  },
  divider: {
    margin: '10px 0',
  },
});

type Props = ContextRouter &
  WithAlert &
  WithStyles & {
    onClose: () => void,
    onSave: (gatewayID: string) => void,
    gateway: Gateway,
    showRestartCommand: boolean,
    showPingCommand: boolean,
    showGenericCommand: boolean,
  };

type State = {
  showRebootCheck: boolean,
  showRestartCheck: boolean,
  pingHosts: string,
  pingPackets: string,
  pingResponse: string,
  showPingProgress: boolean,
  genericCommandName: string,
  genericParams: string,
  genericResponse: string,
  showGenericProgress: boolean,
};

function CommandResponse(props) {
  return (
    <pre
      style={{
        backgroundColor: grey[100],
        fontSize: '12px',
        color: grey[900],
      }}>
      {props.showProgressBar && <LinearProgress />}
      <code>{props.response}</code>
    </pre>
  );
}

class GatewayCommandFields extends React.Component<Props, State> {
  state = {
    showRebootCheck: false,
    showRestartCheck: false,
    pingHosts: '',
    pingPackets: '',
    pingResponse: '',
    showPingProgress: false,
    genericCommandName: '',
    genericParams: '{\n}',
    genericResponse: '',
    showGenericProgress: false,
  };

  render() {
    return (
      <>
        <DialogContent>
          <Typography className={this.props.classes.title} variant="subtitle1">
            Reboot
          </Typography>
          <FormField label="Reboot Device">
            <Button onClick={this.handleRebootGateway} color="primary">
              Reboot
            </Button>
            <Fade in={this.state.showRebootCheck} timeout={500}>
              <Check style={{verticalAlign: 'middle'}} nativeColor="green" />
            </Fade>
          </FormField>
          <div style={this.props.showRestartCommand ? {} : {display: 'none'}}>
            <FormField label="Restart Services">
              <Button onClick={this.handleRestartServices} color="primary">
                Restart Services
              </Button>
              <Fade in={this.state.showRestartCheck} timeout={500}>
                <Check style={{verticalAlign: 'middle'}} nativeColor="green" />
              </Fade>
            </FormField>
          </div>
          <div style={this.props.showPingCommand ? {} : {display: 'none'}}>
            <Divider className={this.props.classes.divider} />
            <Typography
              className={this.props.classes.title}
              variant="subtitle1">
              Ping
            </Typography>
            <FormField label="Host(s) (one per line)">
              <Input
                className={this.props.classes.input}
                value={this.state.pingHosts}
                onChange={this.pingHostsChanged}
                placeholder="E.g. example.com"
                multiline={true}
              />
            </FormField>
            <FormField label="Packets (default 4)">
              <Input
                className={this.props.classes.input}
                value={this.state.pingPackets}
                onChange={this.pingPacketsChanged}
                placeholder="E.g. 4"
                type="number"
              />
            </FormField>
            <FormField label="">
              <Button onClick={this.handlePing} color="primary">
                Ping
              </Button>
            </FormField>
            <FormField label="">
              <CommandResponse
                response={this.state.pingResponse}
                showProgressBar={this.state.showPingProgress}
              />
            </FormField>
          </div>
          <div style={this.props.showGenericCommand ? {} : {display: 'none'}}>
            <Divider className={this.props.classes.divider} />
            <Typography
              className={this.props.classes.title}
              variant="subtitle1">
              Generic
            </Typography>
            <FormField label="Command">
              <Input
                className={this.props.classes.input}
                value={this.state.genericCommandName}
                onChange={this.genericCommandNameChanged}
                placeholder="Command name"
              />
            </FormField>
            <FormField label="Parameters">
              <Input
                className={this.props.classes.input}
                value={this.state.genericParams}
                onChange={this.genericParamsChanged}
                multiline={true}
                style={{fontFamily: 'monospace', fontSize: '14px'}}
              />
            </FormField>
            <FormField label="">
              <Button onClick={this.handleGeneric} color="primary">
                Execute
              </Button>
            </FormField>
            <FormField label="">
              <CommandResponse
                response={this.state.genericResponse}
                showProgressBar={this.state.showGenericProgress}
              />
            </FormField>
          </div>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onClose} color="primary">
            Close
          </Button>
        </DialogActions>
      </>
    );
  }

  handleRebootGateway = () => {
    const {match, gateway} = this.props;
    const id = gateway.logicalID;
    const commandName = 'reboot';

    axios
      .post(MagmaAPIUrls.command(match, id, commandName))
      .then(_resp => {
        this.props.alert('Successfully initiated reboot');
        this.setState({showRebootCheck: true}, () => {
          setTimeout(() => this.setState({showRebootCheck: false}), 5000);
        });
      })
      .catch(error =>
        this.props.alert('Reboot failed: ' + error.response.data.message),
      );
  };

  handleRestartServices = () => {
    const {match, gateway} = this.props;
    const id = gateway.logicalID;
    const commandName = 'restart_services';

    axios
      .post(MagmaAPIUrls.command(match, id, commandName), [])
      .then(_resp => {
        this.props.alert('Successfully initiated service restart');
        this.setState({showRestartCheck: true}, () => {
          setTimeout(() => this.setState({showRestartCheck: false}), 5000);
        });
      })
      .catch(error =>
        this.props.alert(
          'Restart services failed: ' + error.response.data.message,
        ),
      );
  };

  handlePing = () => {
    const {match, gateway} = this.props;
    const id = gateway.logicalID;
    const commandName = 'ping';

    const hosts = this.state.pingHosts.split('\n').filter(host => host);
    const packets = parseInt(this.state.pingPackets);
    const params = {
      hosts,
      packets,
    };

    this.setState({showPingProgress: true});
    axios
      .post(MagmaAPIUrls.command(match, id, commandName), params)
      .then(resp => {
        this.setState({pingResponse: JSON.stringify(resp.data, null, 2)});
      })
      .catch(error =>
        this.props.alert('Ping failed: ' + error.response.data.message),
      )
      .finally(() => this.setState({showPingProgress: false}));
  };

  handleGeneric = () => {
    const {match, gateway} = this.props;
    const id = gateway.logicalID;
    const commandName = 'generic';

    const genericCommandName = this.state.genericCommandName;
    let genericCommandParams = {};
    try {
      genericCommandParams = JSON.parse(this.state.genericParams);
    } catch (e) {
      this.props.alert('Error parsing params: ' + e);
      return;
    }
    const params = {
      command: genericCommandName,
      params: genericCommandParams,
    };

    this.setState({showGenericProgress: true});
    axios
      .post(MagmaAPIUrls.command(match, id, commandName), params)
      .then(resp => {
        this.setState({genericResponse: JSON.stringify(resp.data, null, 2)});
      })
      .catch(error =>
        this.props.alert(
          'Generic command failed: ' + error.response.data.message,
        ),
      )
      .finally(() => this.setState({showGenericProgress: false}));
  };

  pingHostsChanged = ({target}) => this.setState({pingHosts: target.value});
  pingPacketsChanged = ({target}) => this.setState({pingPackets: target.value});

  genericCommandNameChanged = ({target}) =>
    this.setState({genericCommandName: target.value});
  genericParamsChanged = ({target}) => {
    this.setState({genericParams: target.value});
  };
}

export default withStyles(styles)(withRouter(withAlert(GatewayCommandFields)));
