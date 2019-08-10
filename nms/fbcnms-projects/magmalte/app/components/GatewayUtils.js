/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  CellularConfig,
  MagmadConfig,
  Record,
} from '../common/MagmaAPIType';
import type {WithStyles} from '@material-ui/core';

import {withStyles} from '@material-ui/core/styles';

import React from 'react';

export const toString = (input: ?number | ?string): string => {
  return input !== null && input !== undefined ? input + '' : '';
};

export type GatewayPayload = {
  gateway_id: string,
  config: ?{
    cellular_gateway: ?CellularConfig,
    magmad_gateway: ?MagmadConfig,
  },
  record: ?(Record & {
    key: {key_type: string},
  }),
  status: ?{
    checkin_time: number,
    version: string,
    vpn_ip: string,
    meta: ?{
      gps_latitude: number,
      gps_longitude: number,
      rf_tx_on: boolean,
      enodeb_connected: number,
      gps_connected: number,
      mme_connected: number,
    },
  },
};

export type Gateway = {
  hwid: string,
  name: string,
  logicalID: string,
  challengeType: string,
  enodebRFTXEnabled: boolean,
  enodebRFTXOn: boolean,
  latLon: {lat: number, lon: number},
  version: string,
  vpnIP: string,
  enodebConnected: boolean,
  gpsConnected: boolean,
  isBackhaulDown: boolean,
  lastCheckin: string,
  mmeConnected: boolean,
  autoupgradePollInterval: ?number,
  checkinInterval: ?number,
  checkinTimeout: ?number,
  tier: ?string,
  autoupgradeEnabled: boolean,
  attachedEnodebSerials: Array<string>,
  ran: {pci: ?number, transmitEnabled: boolean},
  epc: {ipBlock: string, natEnabled: boolean},
  nonEPSService: {
    control: number,
    csfbRAT: number,
    csfbMCC: ?number,
    csfbMNC: ?number,
    lac: ?number,
  },
  rawGateway: GatewayPayload,
};

const styles = {
  status: {
    width: '10px',
    height: '10px',
    borderRadius: '50%',
    display: 'inline-block',
    textAlign: 'center',
    color: 'white',
    fontSize: '10px',
    fontWeight: 'bold',
    marginRight: '5px',
  },
};

const GatewayStatusElement = (
  props: WithStyles<typeof styles> & {isGrey: boolean, isActive: boolean},
) => {
  if (props.isGrey) {
    return (
      <span
        className={props.classes.status}
        style={{backgroundColor: '#bec3c8'}}
      />
    );
  } else if (props.isActive) {
    return (
      <span
        className={props.classes.status}
        style={{backgroundColor: '#05a503'}}
      />
    );
  } else {
    return (
      <span
        className={props.classes.status}
        style={{backgroundColor: '#fa3a3f'}}
      />
    );
  }
};

export const GatewayStatus = withStyles(styles)(GatewayStatusElement);
