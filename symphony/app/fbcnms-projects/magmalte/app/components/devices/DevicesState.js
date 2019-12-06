/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {FullDevice} from './DevicesUtils';

import CheckBoxIcon from '@material-ui/icons/CheckBox';
import CheckBoxOutlineBlankIcon from '@material-ui/icons/CheckBoxOutlineBlank';
import Checkbox from '@material-ui/core/Checkbox';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React, {useState} from 'react';

import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(() => ({
  checkbox: {
    padding: '4px',
  },
}));

type Props = {
  device: FullDevice,
};

// TODO: complete device interface model
type InterfaceType = {
  'oper-status'?: string,
  name?: string,

  state: {
    ifindex?: number,
    name?: string,
    'oper-status'?: string,
  },

  config: {
    enabled: boolean,
  },

  subinterfaces?: {
    subinterface: Array<{
      'openconfig-if-ip:ipv4': {
        addresses: {
          address: Array<{ip: string}>,
        },
      },
    }>,
  },
};

function Interface({
  iface,
  device,
  isLoading,
  setIsLoading,
}: {
  iface: InterfaceType,
  device: FullDevice,
  isLoading: boolean,
  setIsLoading: boolean => void,
}) {
  const classes = useStyles();
  const {networkId} = useRouter().match.params;

  const ip =
    iface?.subinterfaces?.subinterface?.[0]?.['openconfig-if-ip:ipv4']
      ?.addresses?.address?.[0]?.ip || '';

  const onChange = async event => {
    setIsLoading(true);
    event.persist();
    const newDevice = await MagmaV1API.getSymphonyByNetworkIdDevicesByDeviceId({
      networkId,
      deviceId: device.id,
    });

    const managedDevices = JSON.parse(newDevice.config?.device_config || '{}');
    const index = findIndex(
      managedDevices['openconfig-interfaces:interfaces'].interface,
      i => iface.name === i.name,
    );
    managedDevices['openconfig-interfaces:interfaces'].interface[
      index
    ].config.enabled = event.target.checked;

    MagmaV1API.putSymphonyByNetworkIdDevicesByDeviceId({
      networkId,
      deviceId: device.id,
      symphonyDevice: {
        ...newDevice,
        config: {
          ...newDevice.config,
          device_config: JSON.stringify(managedDevices),
        },
      },
    });

    setIsLoading(false);
  };

  return (
    <div>
      <DeviceStatusCircle
        isGrey={false}
        isActive={(iface.state || iface)['oper-status'] === 'UP'}
      />
      {iface.name || iface.state?.name || ''}
      {ip && ` (${ip})`}
      <Checkbox
        disabled={isLoading}
        className={classes.checkbox}
        defaultChecked={iface?.config?.enabled}
        onChange={onChange}
        color="primary"
        icon={<CheckBoxOutlineBlankIcon fontSize="small" />}
        checkedIcon={<CheckBoxIcon fontSize="small" />}
      />
    </div>
  );
}

type LatenciesStateModel = {
  latency?: Array<{type: string, src: string, dst: string, rtt: number}>,
};

function renderLatencies(
  state: ?LatenciesStateModel,
): Array<React$Element<'div'>> {
  // if no state, then display nothing (different from empty list)
  if (!state) {
    return [];
  }

  const info = [];
  if (!(state.latency?.length == 0)) {
    info.push(
      ...(state.latency || []).map((latency, i) => {
        const key = `latencies_${i}`;
        const rtt = latency.rtt > 0 ? `${latency.rtt / 1000} ms` : 'timeout';
        return (
          <div key={key}>
            {latency.src} -> {latency.dst} ({latency.type}): {rtt}
          </div>
        );
      }),
    );
  }

  if (info.length == 0) {
    info.push(<div key="latencies_none">No latencies reported</div>);
  }

  return info;
}

export default function DevicesState(props: Props) {
  const [isLoading, setIsLoading] = useState(false);

  const {device} = props;
  const interfaces: ?Array<InterfaceType> =
    device?.status?.['openconfig-interfaces:interfaces']?.interface;
  const latencies: ?LatenciesStateModel =
    device?.status?.['fbc-symphony-device:system']?.['latencies'];

  if (!interfaces && !latencies) {
    return <div>{'<No state reported>'}</div>;
  }

  const interfaceRows = (interfaces || []).map((iface, i) => (
    <Interface
      key={i}
      iface={iface}
      device={device}
      isLoading={isLoading}
      setIsLoading={setIsLoading}
    />
  ));

  if (interfaceRows.length === 0) {
    interfaceRows.push(<div key="interfaces_none">No interfaces reported</div>);
  }

  return (
    <>
      {interfaceRows}
      {renderLatencies(latencies)}
    </>
  );
}
