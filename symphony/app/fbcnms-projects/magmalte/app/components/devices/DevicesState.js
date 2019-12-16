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
  italicUnderline: {
    textDecoration: 'underline',
    fontStyle: 'italic',
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

function interfaceOperStatusIsUp(iface: InterfaceType) {
  return (iface.state || iface)['oper-status'] === 'UP';
}

function Interface({
  iface,
  ifaceIsUp,
  device,
  isLoading,
  setIsLoading,
}: {
  iface: InterfaceType,
  ifaceIsUp: boolean,
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
    let index = findIndex(
      managedDevices['openconfig-interfaces:interfaces'].interface,
      i => iface.name === i.name,
    );

    if (index === -1) {
      const newIface = {
        name: iface.name,
        config: JSON.parse(JSON.stringify(iface.config)),
      };
      managedDevices['openconfig-interfaces:interfaces'].interface.push(
        newIface,
      );
      index =
        managedDevices['openconfig-interfaces:interfaces'].interface.length - 1;
    }
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
      <DeviceStatusCircle isGrey={false} isActive={ifaceIsUp} />
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

function ShowInterfacesList({
  interfaces,
  countUp,
  countDown,
}: {
  interfaces: Array<React$Node>,
  countUp: number,
  countDown: number,
}) {
  const [showsInterfaces, setShowsInterfaces] = useState(false);
  const classes = useStyles();
  return (
    <>
      <div>
        {countUp} interfaces up, {countDown} down&nbsp;&ndash;&nbsp;
        <span
          className={classes.italicUnderline}
          onClick={() => setShowsInterfaces(!showsInterfaces)}>
          {showsInterfaces ? 'Hide List' : 'Show List'}
        </span>
      </div>
      {showsInterfaces && <div key="HiddenInterfaces">{interfaces}</div>}
    </>
  );
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

  let upInterfaceCount = 0;
  let downInterfaceCount = 0;
  const interfaceRows = (interfaces || []).map((iface, i) => {
    const interfaceUp = interfaceOperStatusIsUp(iface);
    if (interfaceUp) {
      upInterfaceCount += 1;
    } else {
      downInterfaceCount += 1;
    }
    return (
      <Interface
        key={i}
        iface={iface}
        ifaceIsUp={interfaceUp}
        device={device}
        isLoading={isLoading}
        setIsLoading={setIsLoading}
      />
    );
  });

  let interfacesDiv;
  if (interfaceRows.length === 0) {
    interfacesDiv = <div>No interfaces reported</div>;
  } else {
    interfacesDiv = (
      <ShowInterfacesList
        interfaces={interfaceRows}
        countUp={upInterfaceCount}
        countDown={downInterfaceCount}
      />
    );
  }
  return (
    <>
      {renderLatencies(latencies)}
      {interfacesDiv}
    </>
  );
}
