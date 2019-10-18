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
import type {
  Record,
  WifiConfig,
} from '@fbcnms/magmalte/app/common/MagmaAPIType';
import type {magmad_gateway_configs} from '@fbcnms/magma-api';

import type {WifiGateway} from './WifiUtils';
import type {WithStyles} from '@material-ui/core';

import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormLabel from '@material-ui/core/FormLabel';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaDeviceFields from '@fbcnms/magmalte/app/components/MagmaDeviceFields';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import WifiDeviceFields from './WifiDeviceFields';
import WifiDeviceHardwareFields from './WifiDeviceHardwareFields';
import axios from 'axios';

import GatewayCommandFields from '@fbcnms/magmalte/app/components/GatewayCommandFields';
import nullthrows from '@fbcnms/util/nullthrows';
import {
  DEFAULT_HW_ID_PREFIX,
  DEFAULT_WIFI_GATEWAY_CONFIGS,
  additionalPropsToArray,
  additionalPropsToObject,
  buildWifiGatewayFromPayload,
} from './WifiUtils';
import {
  MagmaAPIUrls,
  createDevice,
  fetchDevice,
} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  appBar: {
    backgroundColor: '#f5f5f5',
    marginBottom: '20px',
  },
};

type Props = ContextRouter &
  WithStyles<typeof styles> & {
    title: string,
    onCancel: () => void,
    onSave: WifiGateway => void,
  };

type State = {
  record: ?Record,
  macAddress: string,
  magmaConfigs: ?magmad_gateway_configs,
  magmaConfigsChanged: boolean,
  error: string,
  status: ?{
    [key: string]: string,
    meta: {
      [key: string]: string,
    },
  },
  wifiConfigs: ?{
    ...WifiConfig,
    latitude: string,
    longitude: string,
    additional_props: Array<[string, string]>,
  },
  wifiConfigsChanged: boolean,
  tab: number,
};

class WifiDeviceDialog extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    const {deviceID, meshID} = props.match.params;

    this.state = ({
      record: null,
      macAddress: '',
      magmaConfigs: null,
      magmaConfigsChanged: false,
      error: '',
      status: null,
      tab: 0,
      wifiConfigs: deviceID
        ? null
        : {
            mesh_id: nullthrows(meshID),
            info: '',
            latitude: '',
            longitude: '',
            client_channel: '11',
            is_production: false,
            additional_props: [['', '']],
          },
      wifiConfigsChanged: false,
    }: State);
  }

  componentDidMount() {
    const {deviceID} = this.props.match.params;
    if (!deviceID) {
      return;
    }

    (async () => {
      try {
        const device = await fetchDevice(this.props.match, deviceID);
        this.setState({
          record: {...device.record},
          magmaConfigs: {
            ...device.config.magmad_gateway,
          },
          status: {...device.status},
          wifiConfigs: {
            ...device.config.wifi_gateway,
            additional_props: additionalPropsToArray(
              device.config.wifi_gateway.additional_props,
            ),
          },
        });
      } catch (error) {
        this.props.onCancel();
      }
    })();
  }

  render() {
    const {deviceID} = this.props.match.params;
    if (!this.state.wifiConfigs) {
      return <LoadingFillerBackdrop />;
    }

    const header = deviceID ? (
      <AppBar position="static" className={this.props.classes.appBar}>
        <Tabs
          indicatorColor="primary"
          textColor="primary"
          value={this.state.tab}
          onChange={this.onTabChange}>
          <Tab label="Wi-Fi" />
          <Tab label="Controller" />
          <Tab label="Hardware" />
          <Tab label="Command" />
        </Tabs>
      </AppBar>
    ) : (
      <DialogTitle>Add Device</DialogTitle>
    );

    let content;
    switch (this.state.tab) {
      case 0:
        content = (
          <WifiDeviceFields
            macAddress={this.state.macAddress}
            status={this.state.status}
            configs={nullthrows(this.state.wifiConfigs)}
            handleMACAddressChange={
              deviceID ? undefined : this.handleMACAddressChange
            }
            configChangeHandler={this.wifiConfigChangeHandler}
          />
        );
        break;
      case 1:
        content = (
          <MagmaDeviceFields
            configs={nullthrows(this.state.magmaConfigs)}
            configChangeHandler={this.magmaConfigChangeHandler}
          />
        );
        break;
      case 2:
        content = (
          <WifiDeviceHardwareFields record={nullthrows(this.state.record)} />
        );
        break;
      case 3:
        content = (
          <GatewayCommandFields
            // $FlowFixMe: deviceID is nullable. Please fix.
            gatewayID={deviceID}
            showRestartCommand={true}
            showRebootEnodebCommand={false}
            showPingCommand={true}
            showGenericCommand={true}
          />
        );
        break;
    }

    return (
      <Dialog
        open={true}
        onClose={this.props.onCancel}
        maxWidth="md"
        scroll="body">
        {header}
        <DialogContent>
          {this.state.error ? (
            <FormLabel error>{this.state.error}</FormLabel>
          ) : null}
          {content}
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onCancel} color="primary">
            Cancel
          </Button>
          <Button
            onClick={deviceID ? this.onEdit : this.onCreate}
            color="primary"
            variant="contained">
            Save
          </Button>
        </DialogActions>
      </Dialog>
    );
  }

  handleMACAddressChange = macAddress => this.setState({macAddress});
  wifiConfigChangeHandler = (fieldName, value) => {
    this.setState({
      wifiConfigsChanged: true,
      wifiConfigs: {
        ...this.state.wifiConfigs,
        [fieldName]: value,
      },
    });
  };

  magmaConfigChangeHandler = (fieldName, value) => {
    this.setState({
      magmaConfigsChanged: true,
      magmaConfigs: {
        ...this.state.magmaConfigs,
        [fieldName]: value,
      },
    });
  };

  onEdit = async () => {
    const {match} = this.props;
    const deviceID = nullthrows(match.params.deviceID);

    const requests = [];
    if (this.state.wifiConfigsChanged) {
      requests.push(
        axios.put(
          MagmaAPIUrls.gatewayConfigsForType(match, deviceID, 'wifi'),
          this.getWifiConfigs(),
        ),
      );
    }

    if (this.state.magmaConfigsChanged) {
      requests.push(
        MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdMagmad({
          networkId: nullthrows(match.params.networkId),
          gatewayId: deviceID,
          magmad: this.getMagmaConfigs(),
        }),
      );
    }

    await axios.all(requests);
    const result = await fetchDevice(match, deviceID);
    this.props.onSave(buildWifiGatewayFromPayload(result));
  };

  onCreate = async () => {
    const {match} = this.props;
    if (!this.state.macAddress || !nullthrows(this.state.wifiConfigs).info) {
      this.setState({error: 'MAC Address and Info fields cannot be empty'});
      return;
    }

    const meshID = nullthrows(match.params.meshID);
    const macAddress = this.state.macAddress.replace(/[:]/g, '').toLowerCase();

    const data = {
      hardware_id: DEFAULT_HW_ID_PREFIX + macAddress,
      key: {key_type: 'ECHO'},
    };

    try {
      const result = await createDevice(
        meshID + '_id_' + macAddress,
        data,
        'wifi',
        DEFAULT_WIFI_GATEWAY_CONFIGS,
        this.getWifiConfigs(),
        match,
      );
      this.props.onSave(buildWifiGatewayFromPayload(result));
    } catch (e) {
      this.setState({error: e.response.data.message || e.message});
    }
  };

  getWifiConfigs(): WifiConfig {
    const wifiConfigs = nullthrows(this.state.wifiConfigs);
    const {additional_props, latitude, longitude, ...otherFields} = wifiConfigs;
    return {
      ...otherFields,
      latitude: parseFloat(latitude),
      longitude: parseFloat(longitude),
      additional_props: additionalPropsToObject(additional_props),
    };
  }

  getMagmaConfigs(): magmad_gateway_configs {
    const magmaConfigs = nullthrows(this.state.magmaConfigs);
    return {
      ...magmaConfigs,
      autoupgrade_poll_interval: parseInt(
        magmaConfigs.autoupgrade_poll_interval,
      ),
      checkin_interval: parseInt(magmaConfigs.checkin_interval),
      checkin_timeout: parseInt(magmaConfigs.checkin_timeout),
    };
  }

  onTabChange = (event, tab) => this.setState({tab});
}

export default withStyles(styles)(withRouter(WifiDeviceDialog));
