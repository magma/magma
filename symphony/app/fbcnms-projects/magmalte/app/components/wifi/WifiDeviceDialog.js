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
  gateway_device,
  gateway_status,
  gateway_wifi_configs,
  magmad_gateway_configs,
  wifi_gateway,
} from '@fbcnms/magma-api';

import type {WifiGateway} from './WifiUtils';
import type {WithStyles} from '@material-ui/core';

import AppBar from '@material-ui/core/AppBar';
import Button from '@fbcnms/ui/components/design-system/Button';
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

import GatewayCommandFields from '@fbcnms/magmalte/app/components/GatewayCommandFields';
import nullthrows from '@fbcnms/util/nullthrows';
import {
  DEFAULT_HW_ID_PREFIX,
  DEFAULT_WIFI_GATEWAY_CONFIGS,
  additionalPropsToArray,
  additionalPropsToObject,
  buildWifiGatewayFromPayloadV1,
} from './WifiUtils';
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
  record: ?gateway_device,
  macAddress: string,
  magmaConfigs: ?magmad_gateway_configs,
  magmaConfigsChanged: boolean,
  error: string,
  status: ?gateway_status,
  wifiConfigs: ?gateway_wifi_configs,
  additionalProps: Array<[string, string]>,
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
            client_channel: '11',
            is_production: false,
          },
      additionalProps: [['', '']],
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
        const device = await MagmaV1API.getWifiByNetworkIdGatewaysByGatewayId({
          networkId: nullthrows(this.props.match.params.networkId),
          gatewayId: deviceID,
        });
        this.setState({
          record: device.device,
          magmaConfigs: device.magmad,
          status: device.status,
          wifiConfigs: device.wifi,
          additionalProps:
            additionalPropsToArray(device.wifi.additional_props) || [],
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
            additionalProps={this.state.additionalProps}
            handleMACAddressChange={
              deviceID ? undefined : this.handleMACAddressChange
            }
            configChangeHandler={this.wifiConfigChangeHandler}
            additionalPropsChangeHandler={this.wifiPropsConfigChangeHandler}
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
        content =
          (deviceID && (
            <GatewayCommandFields
              gatewayID={deviceID}
              showRestartCommand={true}
              showRebootEnodebCommand={false}
              showPingCommand={true}
              showGenericCommand={true}
            />
          )) ||
          null;
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
          <Button onClick={this.props.onCancel} skin="regular">
            Cancel
          </Button>
          <Button onClick={deviceID ? this.onEdit : this.onCreate}>Save</Button>
        </DialogActions>
      </Dialog>
    );
  }

  handleMACAddressChange = macAddress =>
    this.setState({macAddress: macAddress.trim()});
  wifiConfigChangeHandler = (fieldName, value) => {
    this.setState({
      wifiConfigsChanged: true,
      wifiConfigs: {
        ...nullthrows(this.state.wifiConfigs),
        [fieldName]: value,
      },
    });
  };
  wifiPropsConfigChangeHandler = additionalProps => {
    this.setState({wifiConfigsChanged: true, additionalProps});
  };

  magmaConfigChangeHandler = (fieldName, value) => {
    this.setState({
      magmaConfigsChanged: true,
      magmaConfigs: {
        ...nullthrows(this.state.magmaConfigs),
        [fieldName]: value,
      },
    });
  };

  onEdit = async () => {
    const networkId = nullthrows(this.props.match.params.networkId);
    const gatewayId = nullthrows(this.props.match.params.deviceID);

    try {
      const requests = [];
      if (this.state.wifiConfigsChanged) {
        requests.push(
          MagmaV1API.putWifiByNetworkIdGatewaysByGatewayIdWifi({
            networkId,
            gatewayId,
            config: this.getWifiConfigs(),
          }),
        );
      }

      if (this.state.magmaConfigsChanged) {
        requests.push(
          MagmaV1API.putWifiByNetworkIdGatewaysByGatewayIdMagmad({
            networkId,
            gatewayId,
            magmad: this.getMagmaConfigs(),
          }),
        );
      }
      await Promise.all(requests);

      const result = await MagmaV1API.getWifiByNetworkIdGatewaysByGatewayId({
        networkId,
        gatewayId,
      });

      this.props.onSave(buildWifiGatewayFromPayloadV1(result));
    } catch (e) {
      this.setState({error: e?.response?.data?.message || e.message});
    }
  };

  onCreate = async () => {
    if (!this.state.macAddress || !nullthrows(this.state.wifiConfigs).info) {
      this.setState({error: 'MAC Address and Info fields cannot be empty'});
      return;
    }

    const meshID = nullthrows(this.props.match.params.meshID);
    const macAddress = this.state.macAddress.replace(/[:]/g, '').toLowerCase();

    const data = {
      hardware_id: DEFAULT_HW_ID_PREFIX + macAddress,
      key: {key_type: 'ECHO'},
    };

    try {
      const gateway: wifi_gateway = {
        id: `${meshID}_id_${macAddress}`,
        name: macAddress,
        description: macAddress,
        device: data,
        magmad: DEFAULT_WIFI_GATEWAY_CONFIGS,
        wifi: this.getWifiConfigs(),
        tier: 'default',
      };

      // Workaround(1): there's a bug when creating gateways
      //     where meshId needs to be updated to correctly store the association
      // 1) create gateway with no meshID, then
      // 2) update gateway with meshID - association will be created.
      const {mesh_id: _, ...wifiConfigsNoMesh} = gateway.wifi;
      await MagmaV1API.postWifiByNetworkIdGateways({
        networkId: nullthrows(this.props.match.params.networkId),
        // TODO: replace with "gateway" after workaround(1) is unneeded
        gateway: {...gateway, wifi: wifiConfigsNoMesh},
      });

      // TODO: remove this section after workaround(1) is unneeded
      await MagmaV1API.putWifiByNetworkIdGatewaysByGatewayIdWifi({
        networkId: nullthrows(this.props.match.params.networkId),
        gatewayId: gateway.id,
        config: gateway.wifi,
      });

      this.props.onSave(buildWifiGatewayFromPayloadV1(gateway));
    } catch (e) {
      this.setState({error: e.response.data.message || e.message});
    }
  };

  getWifiConfigs(): gateway_wifi_configs {
    const {latitude, longitude, ...otherFields} = nullthrows(
      this.state.wifiConfigs,
    );

    if (latitude && Number.isNaN(parseFloat(latitude))) {
      throw Error('Latitude invalid');
    }
    if (longitude && Number.isNaN(parseFloat(longitude))) {
      throw Error('Longitude invalid');
    }

    const configs = {
      ...otherFields,
      latitude: latitude ? parseFloat(latitude) : undefined,
      longitude: longitude ? parseFloat(longitude) : undefined,
      additional_props:
        additionalPropsToObject(this.state.additionalProps) || {},
    };

    return configs;
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
