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
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';
import type {network_wifi_configs} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import FormGroup from '@material-ui/core/FormGroup';
import KeyValueFields from '@fbcnms/magmalte/app/components/KeyValueFields';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {
  additionalPropsToArray,
  additionalPropsToObject,
} from '../wifi/WifiUtils';
import {get} from 'lodash';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  block: {
    display: 'block',
    marginRight: theme.spacing(),
    width: '245px',
  },
  formContainer: {
    paddingBottom: theme.spacing(2),
  },
  formGroup: {
    marginLeft: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  keyValueFieldsInputValue: {
    width: '585px',
  },
  saveButton: {
    marginTop: theme.spacing(2),
  },
  textArea: {
    width: '600px',
  },
  textField: {
    marginRight: theme.spacing(),
    width: '245px',
  },
});

type State = {
  config: network_wifi_configs,
  additionalProps: Array<[string, string]>,
  isLoading: boolean,
};

type Props = ContextRouter & WithAlert & WithStyles<typeof styles> & {};

class WifiNetworkConfig extends React.Component<Props, State> {
  state = {
    config: {},
    additionalProps: [],
    isLoading: true,
  };

  componentDidMount() {
    MagmaV1API.getWifiByNetworkIdWifi({
      networkId: nullthrows(this.props.match.params.networkId),
    })
      .then(data =>
        this.setState({
          config: {...data},
          additionalProps: additionalPropsToArray(data.additional_props) || [],
          isLoading: false,
        }),
      )
      .catch(error => {
        this.props.alert(get(error, 'response.data.message', error));
        this.setState({
          isLoading: false,
        });
      });
  }

  render() {
    if (this.state.isLoading) {
      return <LoadingFiller />;
    }

    const {classes} = this.props;
    const {config} = this.state;
    return (
      <div className={classes.formContainer}>
        <FormGroup row className={classes.formGroup}>
          <TextField
            required
            label="VPN Proto"
            margin="normal"
            className={classes.textField}
            value={config.mgmt_vpn_proto}
            onChange={this.handleVPNProtoChanged}
          />
          <TextField
            required
            label="VPN Remote"
            margin="normal"
            className={classes.textField}
            value={config.mgmt_vpn_remote}
            onChange={this.handleVPNRemoteChanged}
          />
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <TextField
            required
            label="Ping Host List"
            margin="normal"
            className={classes.textField}
            value={(config.ping_host_list || []).join(',')}
            onChange={this.handlePingHostListChanged}
          />
          <TextField
            required
            label="Ping Number of Packets"
            margin="normal"
            className={classes.textField}
            value={config.ping_num_packets}
            onChange={this.handlePingNumPacketsChanged}
          />
          <TextField
            required
            label="Ping Timeout (s)"
            margin="normal"
            className={classes.textField}
            value={config.ping_timeout_secs}
            onChange={this.handlePingTimeoutChanged}
          />
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <div>
            <TextField
              required
              label="XWF Radius Server"
              margin="normal"
              className={classes.block}
              value={config.xwf_radius_server}
              onChange={this.handleXWFRadiusServerChanged}
            />
            <TextField
              required
              label="XWF Radius Auth Port"
              margin="normal"
              className={classes.block}
              value={config.xwf_radius_auth_port}
              onChange={this.handleXWFRadiusAuthPortChanged}
            />
            <TextField
              required
              label="XWF Radius Acct Port"
              margin="normal"
              className={classes.block}
              value={config.xwf_radius_acct_port}
              onChange={this.handleXWFRadiusAcctPortChanged}
            />
          </div>
          <div>
            <TextField
              required
              label="XWF Radius Shared Secret"
              margin="normal"
              className={classes.block}
              value={config.xwf_radius_shared_secret}
              onChange={this.handleXWFRadiusSharedSecretChanged}
            />
            <TextField
              required
              label="XWF UAM Secret"
              margin="normal"
              className={classes.block}
              value={config.xwf_uam_secret}
              onChange={this.handleXWFUAMSecretChanged}
            />
          </div>
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <div>
            <TextField
              required
              label="XWF Partner Name"
              margin="normal"
              className={classes.block}
              value={config.xwf_partner_name}
              onChange={this.handleXWFPartnerNameChanged}
            />
            <TextField
              required
              label="XWF DHCP DNS 1"
              margin="normal"
              className={classes.block}
              value={config.xwf_dhcp_dns1}
              onChange={this.handleXWFDHCPDNS1Changed}
            />
            <TextField
              required
              label="XWF DHCP DNS 2"
              margin="normal"
              className={classes.block}
              value={config.xwf_dhcp_dns2}
              onChange={this.handleXWFDHCPDNS2Changed}
            />
          </div>
          <TextField
            multiline
            rowsMax="8"
            label="XWF Config"
            margin="normal"
            className={classes.textArea}
            value={config.xwf_config}
            onChange={this.handleXWFConfigChanged}
          />
        </FormGroup>
        <FormGroup className={classes.formGroup}>
          <KeyValueFields
            keyValuePairs={this.state.additionalProps || [['', '']]}
            onChange={this.handleAdditionalPropsChange}
            classes={{inputValue: classes.keyValueFieldsInputValue}}
          />
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <Button className={classes.saveButton} onClick={this.handleSave}>
            Save
          </Button>
        </FormGroup>
      </div>
    );
  }

  handlePingHostListChanged = ({target}) =>
    this.handleConfigChange('ping_host_list', target.value.split(','));
  handlePingNumPacketsChanged = ({target}) =>
    this.handleConfigChange('ping_num_packets', target.value);
  handlePingTimeoutChanged = ({target}) =>
    this.handleConfigChange('ping_timeout_secs', target.value);
  handleVPNProtoChanged = ({target}) =>
    this.handleConfigChange('mgmt_vpn_proto', target.value);
  handleVPNRemoteChanged = ({target}) =>
    this.handleConfigChange('mgmt_vpn_remote', target.value);
  handleXWFPartnerNameChanged = ({target}) =>
    this.handleConfigChange('xwf_partner_name', target.value);
  handleXWFConfigChanged = ({target}) =>
    this.handleConfigChange('xwf_config', target.value);
  handleXWFDHCPDNS1Changed = ({target}) =>
    this.handleConfigChange('xwf_dhcp_dns1', target.value);
  handleXWFDHCPDNS2Changed = ({target}) =>
    this.handleConfigChange('xwf_dhcp_dns2', target.value);
  handleXWFRadiusServerChanged = ({target}) =>
    this.handleConfigChange('xwf_radius_server', target.value);
  handleXWFRadiusSharedSecretChanged = ({target}) =>
    this.handleConfigChange('xwf_radius_shared_secret', target.value);
  handleXWFRadiusAuthPortChanged = ({target}) =>
    this.handleConfigChange('xwf_radius_auth_port', target.value);
  handleXWFRadiusAcctPortChanged = ({target}) =>
    this.handleConfigChange('xwf_radius_acct_port', target.value);
  handleXWFUAMSecretChanged = ({target}) =>
    this.handleConfigChange('xwf_uam_secret', target.value);
  handleAdditionalPropsChange = value =>
    this.setState({additionalProps: value});

  handleConfigChange = (
    field: string,
    value: string | number | Array<string>,
  ) => {
    this.setState({
      config: {
        ...this.state.config,
        [field]: value,
      },
    });
  };

  handleSave = () => {
    MagmaV1API.putWifiByNetworkIdWifi({
      networkId: nullthrows(this.props.match.params.networkId),
      config: this.getConfigs(),
    })
      .then(_resp => this.props.alert('Saved successfully'))
      .catch(this.props.alert);
  };

  getConfigs(): network_wifi_configs {
    return {
      ...this.state.config,
      ping_num_packets: parseInt(this.state.config.ping_num_packets),
      ping_timeout_secs: parseInt(this.state.config.ping_timeout_secs),
      xwf_radius_auth_port: parseInt(this.state.config.xwf_radius_auth_port),
      xwf_radius_acct_port: parseInt(this.state.config.xwf_radius_acct_port),
      additional_props:
        additionalPropsToObject(this.state.additionalProps) || {},
    };
  }
}

export default withStyles(styles)(withAlert(withRouter(WifiNetworkConfig)));
