/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */

import type {WifiGateway} from './WifiUtils';
import type {
  gateway_device,
  gateway_status,
  gateway_wifi_configs,
  magmad_gateway_configs,
  wifi_gateway,
} from '@fbcnms/magma-api';

import AppBar from '@material-ui/core/AppBar';
import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import GatewayCommandFields from '@fbcnms/magmalte/app/components/GatewayCommandFields';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaDeviceFields from '@fbcnms/magmalte/app/components/MagmaDeviceFields';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import WifiDeviceFields from './WifiDeviceFields';
import WifiDeviceHardwareFields from './WifiDeviceHardwareFields';

import nullthrows from '@fbcnms/util/nullthrows';
import {
  DEFAULT_HW_ID_PREFIX,
  DEFAULT_WIFI_GATEWAY_CONFIGS,
  additionalPropsToArray,
  additionalPropsToObject,
  buildWifiGatewayFromPayloadV1,
} from './WifiUtils';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(() => ({
  appBar: {
    backgroundColor: '#f5f5f5',
    marginBottom: '20px',
  },
}));

type Props = {
  title: string,
  onCancel: () => void,
  onSave: WifiGateway => void,
};

export default function WifiDeviceDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const {deviceID, meshID} = match.params;

  const [record, setRecord] = useState<?gateway_device>();
  const [macAddress, setMacAddress] = useState('');
  const [magmaConfigs, setMagmaConfigs] = useState<?magmad_gateway_configs>();
  const [magmaConfigsChanged, setMmagmaConfigsChanged] = useState(false);
  const [status, setStatus] = useState<?gateway_status>();
  const [wifiConfigs, setWifiConfigs] = useState<?gateway_wifi_configs>(
    deviceID
      ? null
      : {
          mesh_id: nullthrows(meshID),
          info: '',
          client_channel: '11',
          is_production: false,
        },
  );
  const [wifiConfigsChanged, setWifiConfigsChanged] = useState(false);
  const [additionalProps, setAdditionalProps] = useState([['', '']]);
  const [tab, setTab] = useState(0);

  useEffect(() => {
    if (deviceID) {
      MagmaV1API.getWifiByNetworkIdGatewaysByGatewayId({
        networkId: nullthrows(match.params.networkId),
        gatewayId: deviceID,
      })
        .then(device => {
          setRecord(device.device);
          setMagmaConfigs(device.magmad);
          setStatus(device.status);
          setWifiConfigs(device.wifi);
          setAdditionalProps(
            additionalPropsToArray(device.wifi.additional_props) || [],
          );
        })
        .catch(e => {
          enqueueSnackbar(e?.response?.data?.message || e.message, {
            variant: 'error',
          });
          props.onCancel();
        });
    }
  }, [deviceID, enqueueSnackbar, match.params.networkId, props]);

  if (!wifiConfigs) {
    return <LoadingFillerBackdrop />;
  }

  const wifiConfigChangeHandler = (fieldName, value) => {
    setWifiConfigsChanged(true);
    setWifiConfigs({
      ...nullthrows(wifiConfigs),
      [fieldName]: value,
    });
  };

  const magmaConfigChangeHandler = (fieldName, value) => {
    setMmagmaConfigsChanged(true);
    setMagmaConfigs({
      ...nullthrows(magmaConfigs),
      [fieldName]: value,
    });
  };

  const onEdit = async () => {
    const networkId = nullthrows(match.params.networkId);
    const gatewayId = nullthrows(match.params.deviceID);

    try {
      const requests = [];
      if (wifiConfigsChanged) {
        requests.push(
          MagmaV1API.putWifiByNetworkIdGatewaysByGatewayIdWifi({
            networkId,
            gatewayId,
            config: getWifiConfigs(wifiConfigs, additionalProps),
          }),
        );
      }

      if (magmaConfigsChanged) {
        requests.push(
          MagmaV1API.putWifiByNetworkIdGatewaysByGatewayIdMagmad({
            networkId,
            gatewayId,
            magmad: getMagmaConfigs(magmaConfigs),
          }),
        );
      }
      await Promise.all(requests);

      const result = await MagmaV1API.getWifiByNetworkIdGatewaysByGatewayId({
        networkId,
        gatewayId,
      });

      props.onSave(buildWifiGatewayFromPayloadV1(result));
    } catch (e) {
      enqueueSnackbar(e?.response?.data?.message || e.message, {
        variant: 'error',
      });
    }
  };

  const onCreate = async () => {
    if (!macAddress || !nullthrows(wifiConfigs).info) {
      this.setState({error: 'MAC Address and Info fields cannot be empty'});
      return;
    }

    const meshID = nullthrows(match.params.meshID);
    const sanitizedMac = macAddress.replace(/[:]/g, '').toLowerCase();

    const data = {
      hardware_id: DEFAULT_HW_ID_PREFIX + sanitizedMac,
      key: {key_type: 'ECHO'},
    };

    try {
      const gateway: wifi_gateway = {
        id: `${meshID}_id_${sanitizedMac}`,
        name: sanitizedMac,
        description: sanitizedMac,
        device: data,
        magmad: DEFAULT_WIFI_GATEWAY_CONFIGS,
        wifi: getWifiConfigs(wifiConfigs, additionalProps),
        tier: 'default',
      };

      // Workaround(1): there's a bug when creating gateways
      //     where meshId needs to be updated to correctly store the association
      // 1) create gateway with no meshID, then
      // 2) update gateway with meshID - association will be created.
      const {mesh_id: _, ...wifiConfigsNoMesh} = gateway.wifi;
      await MagmaV1API.postWifiByNetworkIdGateways({
        networkId: nullthrows(match.params.networkId),
        // TODO: replace with "gateway" after workaround(1) is unneeded
        gateway: {...gateway, wifi: wifiConfigsNoMesh},
      });

      // TODO: remove this section after workaround(1) is unneeded
      await MagmaV1API.putWifiByNetworkIdGatewaysByGatewayIdWifi({
        networkId: nullthrows(match.params.networkId),
        gatewayId: gateway.id,
        config: gateway.wifi,
      });

      props.onSave(buildWifiGatewayFromPayloadV1(gateway));
    } catch (e) {
      enqueueSnackbar(e?.response?.data?.message || e.message, {
        variant: 'error',
      });
    }
  };

  const header = deviceID ? (
    <AppBar position="static" className={classes.appBar}>
      <Tabs
        indicatorColor="primary"
        textColor="primary"
        value={tab}
        onChange={(_, newTab) => setTab(newTab)}>
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
  switch (tab) {
    case 0:
      content = (
        <WifiDeviceFields
          macAddress={macAddress}
          status={status}
          configs={nullthrows(wifiConfigs)}
          additionalProps={additionalProps}
          handleMACAddressChange={
            deviceID ? undefined : m => setMacAddress(m.trim())
          }
          configChangeHandler={wifiConfigChangeHandler}
          additionalPropsChangeHandler={newValue => {
            setAdditionalProps(newValue);
            setWifiConfigsChanged(true);
          }}
        />
      );
      break;
    case 1:
      content = (
        <MagmaDeviceFields
          configs={nullthrows(magmaConfigs)}
          configChangeHandler={magmaConfigChangeHandler}
        />
      );
      break;
    case 2:
      content = <WifiDeviceHardwareFields record={nullthrows(record)} />;
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
    <Dialog open={true} onClose={props.onCancel} maxWidth="md" scroll="body">
      {header}
      <DialogContent>{content}</DialogContent>
      <DialogActions>
        <Button onClick={props.onCancel} skin="regular">
          Cancel
        </Button>
        <Button onClick={deviceID ? onEdit : onCreate}>Save</Button>
      </DialogActions>
    </Dialog>
  );
}

function getWifiConfigs(wifiConfigs, additionalProps): gateway_wifi_configs {
  const {latitude, longitude, ...otherFields} = nullthrows(wifiConfigs);

  if (latitude != null && Number.isNaN(parseFloat(latitude))) {
    throw Error('Latitude invalid');
  }
  if (longitude != null && Number.isNaN(parseFloat(longitude))) {
    throw Error('Longitude invalid');
  }

  const configs = {
    ...otherFields,
    latitude: latitude != null ? parseFloat(latitude) : undefined,
    longitude: longitude != null ? parseFloat(longitude) : undefined,
    additional_props: additionalPropsToObject(additionalProps) || {},
  };

  return configs;
}

function getMagmaConfigs(configs): magmad_gateway_configs {
  const magmaConfigs = nullthrows(configs);
  return {
    ...magmaConfigs,
    autoupgrade_poll_interval: parseInt(magmaConfigs.autoupgrade_poll_interval),
    checkin_interval: parseInt(magmaConfigs.checkin_interval),
    checkin_timeout: parseInt(magmaConfigs.checkin_timeout),
  };
}
