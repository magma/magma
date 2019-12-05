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

import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import React from 'react';

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

function Interface({iface}: {iface: InterfaceType}) {
  const ip =
    iface?.subinterfaces?.subinterface?.[0]?.['openconfig-if-ip:ipv4']
      ?.addresses?.address?.[0]?.ip || '';

  return (
    <div>
      <DeviceStatusCircle
        isGrey={false}
        isActive={(iface.state || iface)['oper-status'] === 'UP'}
      />
      {iface.name || iface.state?.name || ''}
      {ip && ` (${ip})`}
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
  const {device} = props;
  const interfaces: ?Array<InterfaceType> =
    device?.status?.['openconfig-interfaces:interfaces']?.interface;
  const latencies: ?LatenciesStateModel =
    device?.status?.['fbc-symphony-device:system']?.['latencies'];

  if (!interfaces && !latencies) {
    return <div>{'<No state reported>'}</div>;
  }

  const interfaceRows = (interfaces || []).map((iface, i) => (
    <Interface key={i} iface={iface} />
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
