/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export const EnodebDeviceClass = Object.freeze({
  BAICELLS_NOVA_233_2_OD_FDD: 'Baicells Nova-233 G2 OD FDD',
  BAICELLS_NOVA_243_OD_TDD: 'Baicells Nova-243 OD TDD',
  BAICELLS_ID: 'Baicells ID TDD/FDD',
  NURAN_CAVIUM_OC_LTE: 'NuRAN Cavium OC-LTE',
});

export const EnodebBandwidthOption = Object.freeze({
  '3': 3,
  '5': 5,
  '10': 10,
  '15': 15,
  '20': 20,
});

export const DEFAULT_ENODEB = Object.freeze({
  device_class: 'Baicells ID TDD/FDD',
  earfcndl: 0,
  pci: 0,
  special_subframe_pattern: 0,
  subframe_assignment: 0,
  bandwidth_mhz: EnodebBandwidthOption['3'],
  tac: 0,
  cell_id: 0,
  transmit_enabled: false,
});

// This is the format in which the REST API understands the config
export type EnodebPayload = {
  device_class: $Values<typeof EnodebDeviceClass>,
  earfcndl: number,
  pci: number,
  subframe_assignment: number,
  special_subframe_pattern: number,
  bandwidth_mhz: $Values<typeof EnodebBandwidthOption>,
  tac: number,
  cell_id: number,
  transmit_enabled: boolean,
};

export type Enodeb = {
  serialId: string,
  deviceClass: $Values<typeof EnodebDeviceClass>,
  earfcndl: number,
  subframeAssignment: number,
  specialSubframePattern: number,
  pci: number,
  bandwidthMhz: $Values<typeof EnodebBandwidthOption>,
  tac: number,
  cellId: number,
  transmitEnabled: boolean,
};
