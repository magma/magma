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

import {GatewayStatus} from '@fbcnms/magmalte/app/components/GatewayUtils';

type Props = {
  device: FullDevice,
};

export default function DevicesState(props: Props) {
  const {device} = props;

  let info = '<no status reported>';

  const config = device.status;
  if (config) {
    if (!config['openconfig-interfaces:interfaces']) {
      info = `No interfaces reported`;
    } else {
      const interfaces = config['openconfig-interfaces:interfaces'].interface;
      if (!interfaces) {
        info = `No interfaces reported`;
      } else if (interfaces.length == 0) {
        info = `No interfaces reported`;
      } else {
        info = interfaces.map((iface, i) => {
          const ip =
            iface?.subinterfaces?.subinterface?.[0]?.['openconfig-if-ip:ipv4']
              ?.addresses?.address?.[0]?.ip || '';
          return (
            <div key={i}>
              <GatewayStatus
                isGrey={false}
                isActive={(iface.state || iface)['oper-status'] === 'UP'}
              />
              {iface.name || iface.state?.name || ''}
              {ip && <> ({ip})</>}
            </div>
          );
        });
      }
    }
  }

  return info;
}
