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

import React from 'react';

import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';

type Props = {
  device: ?FullDevice,
};

type InterfaceStateModel = {
  // TODO: complete device interface model
  interface?: ?Array<{
    'oper-status'?: string,
    name?: string,

    state: {
      ifindex?: number,
      name?: string,
      'oper-status'?: string,
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
  }>,
};

function GenInterfaces(
  interfaceState: ?InterfaceStateModel,
): Array<React$Element<'div'>> {
  // if no interface state, then display nothing (different from empty list)
  if (!interfaceState) {
    return [];
  }

  const info = [];
  if (!(interfaceState.interface?.length == 0)) {
    info.push(
      ...(interfaceState.interface || []).map((iface, i) => {
        const ip =
          iface?.subinterfaces?.subinterface?.[0]?.['openconfig-if-ip:ipv4']
            ?.addresses?.address?.[0]?.ip || '';
        const key = `interfaces_${i}`;
        return (
          <div key={key}>
            <DeviceStatusCircle
              isGrey={false}
              isActive={(iface.state || iface)['oper-status'] === 'UP'}
            />
            {iface.name || iface.state?.name || ''}
            {ip && <> ({ip})</>}
          </div>
        );
      }),
    );
  }

  if (info.length == 0) {
    info.push(<div key="interfaces_none">No interfaces reported</div>);
  }

  return info;
}

type latenciesStateModel = {
  latency?: Array<{type: string, src: string, dst: string, rtt: number}>,
};

function GenLatencies(
  state: ?latenciesStateModel,
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

export default function DevicesState(
  props: Props,
): Array<React$Element<'div'>> {
  const {device} = props;

  const info = [];

  info.push(
    ...GenInterfaces(device?.status?.['openconfig-interfaces:interfaces']),
    ...GenLatencies(
      device?.status?.['fbc-symphony-device:system']?.['latencies'],
    ),
  );

  if (info.length == 0) {
    info.push(<div key="nostate">&lt;No state reported&gt;</div>);
  }

  return info;
}
