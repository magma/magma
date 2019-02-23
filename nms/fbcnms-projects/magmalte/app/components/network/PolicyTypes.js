/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export type PolicyFlow = {
  action: string,
  match: {
    direction: string,
    ip_proto: string,
    ipv4_src?: string,
    ipv4_dst?: string,
    udp_src?: string,
    udp_dst?: string,
    tcp_src?: string,
    tcp_dst?: string,
  },
};

export type PolicyRule = {
  id: string,
  priority: number,
  flow_list: ?Array<PolicyFlow>,
};

export const ACTION = {
  PERMIT: 'PERMIT',
  DENY: 'DENY',
};

export const DIRECTION = {
  UPLINK: 'UPLINK',
  DOWNLINK: 'DOWNLINK',
};

export const PROTOCOL = {
  IPPROTO_IP: 'IPPROTO_IP',
  IPPROTO_UDP: 'IPPROTO_UDP',
  IPPROTO_TCP: 'IPPROTO_TCP',
  IPPROTO_ICMP: 'IPPROTO_ICMP',
};
