// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <folly/GLog.h>

#include <devmand/models/wifi/Model.h>

namespace devmand {
namespace models {
namespace wifi {

void Model::updateRadio(
    folly::dynamic& state,
    int index,
    const YangPath& path,
    const folly::dynamic& value) {
  auto& ap = state["openconfig-access-points:access-points"]["access-point"][0];
  auto& radios = ap["radios"]["radio"];
  auto* radio = radios.get_ptr(index);
  if (radio == nullptr) {
    folly::dynamic rad = folly::dynamic::object;
    rad["id"] = index;
    auto& radc = rad["config"] = folly::dynamic::object;
    radc["id"] = index;
    // clang-format off
  /*
|     |  +--rw enabled?                  boolean
|     |  +--rw transmit-power?           uint8
|     |  +--rw channel?                  uint8
|     |  +--rw channel-width?            uint8
|     |  +--rw dca?                      boolean
|     |  +--rw allowed-channels*         oc-wifi-types:channels-type
|     |  +--rw dtp?                      boolean
|     |  +--rw dtp-min?                  uint8
|     |  +--rw dtp-max?                  uint8
|     |  +--rw antenna-gain?             int8
|     |  +--rw scanning?                 boolean
|     |  +--rw scanning-interval?        uint8
|     |  +--rw scanning-dwell-time?      uint16
|     |  +--rw scanning-defer-clients?   uint8
|     |  +--rw scanning-defer-traffic?   boolean
  */
    auto& rads = rad["state"] = folly::dynamic::object;
    rads["id"] = index;
  /*
|     |  +--ro enabled?                        boolean
|     |  +--ro transmit-power?                 uint8
|     |  +--ro channel?                        uint8
|     |  +--ro channel-width?                  uint8
|     |  +--ro dca?                            boolean
|     |  +--ro allowed-channels*               oc-wifi-types:channels-type
|     |  +--ro dtp?                            boolean
|     |  +--ro dtp-min?                        uint8
|     |  +--ro dtp-max?                        uint8
|     |  +--ro antenna-gain?                   int8
|     |  +--ro scanning?                       boolean
|     |  +--ro scanning-interval?              uint8
|     |  +--ro scanning-dwell-time?            uint16
|     |  +--ro scanning-defer-clients?         uint8
|     |  +--ro scanning-defer-traffic?         boolean
|     |  +--ro base-radio-mac?                 oc-yang:mac-address
|     |  +--ro dfs-hit-time?                   oc-types:timeticks64
|     |  +--ro channel-change-reason?          identityref
|     |  +--ro total-channel-utilization?      oc-types:percentage
|     |  +--ro rx-dot11-channel-utilization?   oc-types:percentage
|     |  +--ro rx-noise-channel-utilization?   oc-types:percentage
|     |  +--ro tx-dot11-channel-utilization?   oc-types:percentage
  */
    rads["counters"] = folly::dynamic::object;
  /*
|     |     +--ro failed-fcs-frames?   oc-yang:counter64
|     |     +--ro noise-floor?         int8
|     +--rw neighbors
|        +--ro neighbor* [bssid]
|           +--ro bssid    -> ../state/bssid
|           +--ro state
|              +--ro bssid?             oc-yang:mac-address
|              +--ro ssid?              string
|              +--ro rssi?              int8
|              +--ro channel?           uint16
|              +--ro primary-channel?   uint16
|              +--ro last-seen?         oc-types:timeticks64
  */
    // clang-format on
    radios.push_back(rad);
  }
  radio = radios.get_ptr(index);

  YangUtils::set(*radio, path, value);
}

void Model::updateSsid(
    folly::dynamic& state,
    int index,
    const YangPath& path,
    const folly::dynamic& value) {
  auto& ap = state["openconfig-access-points:access-points"]["access-point"][0];
  auto& ssids = ap["ssids"]["ssid"];
  auto* ssid = ssids.get_ptr(index);
  if (ssid == nullptr) {
    folly::dynamic s = folly::dynamic::object;
    s["config"] = folly::dynamic::object;
    s["state"] = folly::dynamic::object;
    auto& bssids = s["bssids"] = folly::dynamic::object;
    bssids["bssid"] = folly::dynamic::array;
    ssids.push_back(s);
  }
  ssid = ssids.get_ptr(index);

  YangUtils::set(*ssid, path, value);
}

void Model::updateSsidBssid(
    folly::dynamic& state,
    int indexSsid,
    int indexBssid,
    const YangPath& path,
    const folly::dynamic& value) {
  auto& ap = state["openconfig-access-points:access-points"]["access-point"][0];
  auto& ssids = ap["ssids"]["ssid"];
  auto* ssid = ssids.get_ptr(indexSsid);
  if (ssid == nullptr) {
    folly::dynamic s = folly::dynamic::object;
    s["config"] = folly::dynamic::object;
    s["state"] = folly::dynamic::object;
    auto& bssids = s["bssids"] = folly::dynamic::object;
    bssids["bssid"] = folly::dynamic::array;
    ssids.push_back(s);
  }
  ssid = ssids.get_ptr(indexSsid);

  auto& bssids = (*ssid)["bssids"]["bssid"];
  auto* bssid = bssids.get_ptr(indexBssid);
  if (bssid == nullptr) {
    folly::dynamic b = folly::dynamic::object;
    b["state"] = folly::dynamic::object;
    bssids.push_back(b);
  }
  bssid = bssids.get_ptr(indexBssid);

  YangUtils::set(*bssid, path, value);
}

void Model::init(folly::dynamic& state) {
  // openconfig-ap-manager ####################################################
  auto& papRoot = state["openconfig-ap-manager:provision-aps"] =
      folly::dynamic::object;
  auto& paps = papRoot["provision-ap"] = folly::dynamic::array;
  folly::dynamic pap = folly::dynamic::object;
  pap["config"] = folly::dynamic::object;
  pap["state"] = folly::dynamic::object;
  paps.push_back(pap);

  auto& japRoot = state["openconfig-ap-manager:joined-aps"] =
      folly::dynamic::object;
  auto& japs = japRoot["joined-ap"] = folly::dynamic::array;
  folly::dynamic jap = folly::dynamic::object;
  jap["state"] = folly::dynamic::object;
  japs.push_back(jap);

  // openconfig-access-points #################################################
  auto& apRoot = state["openconfig-access-points:access-points"] =
      folly::dynamic::object;
  auto& aps = apRoot["access-point"] = folly::dynamic::array;
  folly::dynamic ap = folly::dynamic::object;
  //+--rw hostname                oc-inet:domain-name
  auto& raRoot = ap["radios"] = folly::dynamic::object;
  raRoot["radio"] = folly::dynamic::array;
  auto& ssidRoot = ap["ssids"] = folly::dynamic::object;
  ssidRoot["ssid"] = folly::dynamic::array;
  aps.push_back(ap);
}

// clang-format off
/*
+--rw ssids
|  +--rw ssid* [name]
|     +--rw config
|     |  +--rw enabled?                 boolean
|     |  +--rw hidden?                  boolean
|     |  +--rw default-vlan?            oc-vlan-types:vlan-id
|     |  +--rw vlan-list*               oc-vlan-types:vlan-id
|     |  +--rw basic-data-rates*        identityref
|     |  +--rw supported-data-rates*    identityref
|     |  +--rw broadcast-filter?        boolean
|     |  +--rw multicast-filter?        boolean
|     |  +--rw ipv6-ndp-filter?         boolean
|     |  +--rw ipv6-ndp-filter-timer?   uint16
|     |  +--rw station-isolation?       boolean
|     |  +--rw opmode?                  enumeration
|     |  +--rw wpa2-psk?                string
|     |  +--rw server-group?            string
|     |  +--rw dva?                     boolean
|     |  +--rw mobility-domain?         string
|     |  +--rw dhcp-required?           boolean
|     |  +--rw qbss-load?               boolean
|     |  +--rw advertise-apname?        boolean
|     |  +--rw csa?                     boolean
|     |  +--rw ptk-timeout?             uint16
|     |  +--rw gtk-timeout?             uint16
|     |  +--rw dot11k?                  boolean
|     |  +--rw okc?                     boolean
|     +--ro state
|     |  +--ro enabled?                 boolean
|     |  +--ro hidden?                  boolean
|     |  +--ro default-vlan?            oc-vlan-types:vlan-id
|     |  +--ro vlan-list*               oc-vlan-types:vlan-id
|     |  +--ro basic-data-rates*        identityref
|     |  +--ro supported-data-rates*    identityref
|     |  +--ro broadcast-filter?        boolean
|     |  +--ro multicast-filter?        boolean
|     |  +--ro ipv6-ndp-filter?         boolean
|     |  +--ro ipv6-ndp-filter-timer?   uint16
|     |  +--ro station-isolation?       boolean
|     |  +--ro opmode?                  enumeration
|     |  +--ro wpa2-psk?                string
|     |  +--ro server-group?            string
|     |  +--ro dva?                     boolean
|     |  +--ro mobility-domain?         string
|     |  +--ro dhcp-required?           boolean
|     |  +--ro qbss-load?               boolean
|     |  +--ro advertise-apname?        boolean
|     |  +--ro csa?                     boolean
|     |  +--ro ptk-timeout?             uint16
|     |  +--ro gtk-timeout?             uint16
|     |  +--ro dot11k?                  boolean
|     |  +--ro okc?                     boolean
|     +--rw bssids
|     |  +--ro bssid* [radio-id bssid]
|     |     +--ro bssid       -> ../state/bssid
|     |     +--ro radio-id    -> ../state/radio-id
|     |     +--ro state
|     |        +--ro bssid?                    oc-yang:mac-address
|     |        +--ro radio-id?                 uint8
|     |        +--ro num-associated-clients?   uint8
|     |        +--ro counters
|     |           +--ro rx-bss-dot11-channel-utilization?   oc-types:percentage
|     |           +--ro rx-mgmt?                            oc-yang:counter64
|     |           +--ro rx-control?                         oc-yang:counter64
|     |           +--ro rx-data-dist
|     |           |  +--ro rx-0-64?             oc-yang:counter64
|     |           |  +--ro rx-65-128?           oc-yang:counter64
|     |           |  +--ro rx-129-256?          oc-yang:counter64
|     |           |  +--ro rx-257-512?          oc-yang:counter64
|     |           |  +--ro rx-513-1024?         oc-yang:counter64
|     |           |  +--ro rx-1025-2048?        oc-yang:counter64
|     |           |  +--ro rx-2049-4096?        oc-yang:counter64
|     |           |  +--ro rx-4097-8192?        oc-yang:counter64
|     |           |  +--ro rx-8193-16384?       oc-yang:counter64
|     |           |  +--ro rx-16385-32768?      oc-yang:counter64
|     |           |  +--ro rx-32769-65536?      oc-yang:counter64
|     |           |  +--ro rx-65537-131072?     oc-yang:counter64
|     |           |  +--ro rx-131073-262144?    oc-yang:counter64
|     |           |  +--ro rx-262145-524288?    oc-yang:counter64
|     |           |  +--ro rx-524289-1048576?   oc-yang:counter64
|     |           +--ro rx-data-wmm
|     |           |  +--ro vi?   oc-yang:counter64
|     |           |  +--ro vo?   oc-yang:counter64
|     |           |  +--ro be?   oc-yang:counter64
|     |           |  +--ro bk?   oc-yang:counter64
|     |           +--ro rx-mcs
|     |           |  +--ro mcs0?   oc-yang:counter64
|     |           |  +--ro mcs1?   oc-yang:counter64
|     |           |  +--ro mcs2?   oc-yang:counter64
|     |           |  +--ro mcs3?   oc-yang:counter64
|     |           |  +--ro mcs4?   oc-yang:counter64
|     |           |  +--ro mcs5?   oc-yang:counter64
|     |           |  +--ro mcs6?   oc-yang:counter64
|     |           |  +--ro mcs7?   oc-yang:counter64
|     |           |  +--ro mcs8?   oc-yang:counter64
|     |           |  +--ro mcs9?   oc-yang:counter64
|     |           +--ro rx-retries?                         oc-yang:counter64
|     |           +--ro rx-retries-data?                    oc-yang:counter64
|     |           +--ro rx-retries-subframe?                oc-yang:counter64
|     |           +--ro rx-bytes-data?                      oc-yang:counter64
|     |           +--ro tx-bss-dot11-channel-utilization?   oc-types:percentage
|     |           +--ro tx-mgmt?                            oc-yang:counter64
|     |           +--ro tx-control?                         oc-yang:counter64
|     |           +--ro tx-data-dist
|     |           |  +--ro tx-0-64?             oc-yang:counter64
|     |           |  +--ro tx-65-128?           oc-yang:counter64
|     |           |  +--ro tx-129-256?          oc-yang:counter64
|     |           |  +--ro tx-257-512?          oc-yang:counter64
|     |           |  +--ro tx-513-1024?         oc-yang:counter64
|     |           |  +--ro tx-1025-2048?        oc-yang:counter64
|     |           |  +--ro tx-2049-4096?        oc-yang:counter64
|     |           |  +--ro tx-4097-8192?        oc-yang:counter64
|     |           |  +--ro tx-8193-16384?       oc-yang:counter64
|     |           |  +--ro tx-16385-32768?      oc-yang:counter64
|     |           |  +--ro tx-32769-65536?      oc-yang:counter64
|     |           |  +--ro tx-65537-131072?     oc-yang:counter64
|     |           |  +--ro tx-131073-262144?    oc-yang:counter64
|     |           |  +--ro tx-262145-524288?    oc-yang:counter64
|     |           |  +--ro tx-524289-1048576?   oc-yang:counter64
|     |           +--ro tx-data-wmm
|     |           |  +--ro vi?   oc-yang:counter64
|     |           |  +--ro vo?   oc-yang:counter64
|     |           |  +--ro bk?   oc-yang:counter64
|     |           |  +--ro be?   oc-yang:counter64
|     |           +--ro tx-mcs
|     |           |  +--ro mcs0?   oc-yang:counter64
|     |           |  +--ro mcs1?   oc-yang:counter64
|     |           |  +--ro mcs2?   oc-yang:counter64
|     |           |  +--ro mcs3?   oc-yang:counter64
|     |           |  +--ro mcs4?   oc-yang:counter64
|     |           |  +--ro mcs5?   oc-yang:counter64
|     |           |  +--ro mcs6?   oc-yang:counter64
|     |           |  +--ro mcs7?   oc-yang:counter64
|     |           |  +--ro mcs8?   oc-yang:counter64
|     |           |  +--ro mcs9?   oc-yang:counter64
|     |           +--ro tx-retries?                         oc-yang:counter64
|     |           +--ro tx-retries-data?                    oc-yang:counter64
|     |           +--ro tx-retries-subframe?                oc-yang:counter64
|     |           +--ro tx-bytes-data?                      oc-yang:counter64
|     |           +--ro bss-channel-utilization?            oc-types:percentage
|     +--rw wmm
|     |  +--rw config
|     |  |  +--rw trust-dscp?      boolean
|     |  |  +--rw wmm-vo-remark*   uint8
|     |  |  +--rw wmm-vi-remark*   uint8
|     |  |  +--rw wmm-be-remark*   uint8
|     |  |  +--rw wmm-bk-remark*   uint8
|     |  +--ro state
|     |     +--ro trust-dscp?      boolean
|     |     +--ro wmm-vo-remark*   uint8
|     |     +--ro wmm-vi-remark*   uint8
|     |     +--ro wmm-be-remark*   uint8
|     |     +--ro wmm-bk-remark*   uint8
|     +--rw dot11r
|     |  +--rw config
|     |  |  +--rw dot11r?                 boolean
|     |  |  +--rw dot11r-domainid?        uint16
|     |  |  +--rw dot11r-method?          enumeration
|     |  |  +--rw dot11r-r1key-timeout?   uint16
|     |  +--ro state
|     |     +--ro dot11r?                 boolean
|     |     +--ro dot11r-domainid?        uint16
|     |     +--ro dot11r-method?          enumeration
|     |     +--ro dot11r-r1key-timeout?   uint16
|     +--rw dot11v
|     |  +--rw config
|     |  |  +--rw dot11v-dms?               boolean
|     |  |  +--rw dot11v-bssidle?           boolean
|     |  |  +--rw dot11v-bssidle-timeout?   uint16
|     |  |  +--rw dot11v-bsstransition?     boolean
|     |  +--ro state
|     |     +--ro dot11v-dms?               boolean
|     |     +--ro dot11v-bssidle?           boolean
|     |     +--ro dot11v-bssidle-timeout?   uint16
|     |     +--ro dot11v-bsstransition?     boolean
|     +--rw clients
|     |  +--ro client* [mac]
|     |     +--ro mac                    -> ../state/mac
|     |     +--ro state
|     |     |  +--ro mac?        oc-yang:mac-address
|     |     |  +--ro counters
|     |     |     +--ro tx-bytes?     oc-yang:counter64
|     |     |     +--ro rx-bytes?     oc-yang:counter64
|     |     |     +--ro rx-retries?   oc-yang:counter64
|     |     |     +--ro tx-retries?   oc-yang:counter64
|     |     +--ro client-rf
|     |     |  +--ro state
|     |     |     +--ro rssi?              int8
|     |     |     +--ro snr?               uint8
|     |     |     +--ro ss?                uint8
|     |     |     +--ro phy-rate?          uint16
|     |     |     +--ro connection-mode?   enumeration
|     |     |     +--ro frequency?         uint8
|     |     +--ro client-capabilities
|     |     |  +--ro state
|     |     |     +--ro client-capabilities*   identityref
|     |     |     +--ro channel-support*       uint8
|     |     +--ro dot11k-neighbors
|     |     |  +--ro state
|     |     |     +--ro neighbor-bssid?        oc-yang:mac-address
|     |     |     +--ro neighbor-channel?      uint8
|     |     |     +--ro neighbor-rssi?         int8
|     |     |     +--ro neighbor-antenna?      uint8
|     |     |     +--ro channel-load-report?   uint8
|     |     +--ro client-connection
|     |        +--ro state
|     |           +--ro client-state?       identityref
|     |           +--ro connection-time?    uint16
|     |           +--ro username?           string
|     |           +--ro hostname?           string
|     |           +--ro ipv4-address?       oc-inet:ipv4-address
|     |           +--ro ipv6-address?       oc-inet:ipv6-address
|     |           +--ro operating-system?   string
|     +--rw dot1x-timers
|     |  +--rw config
|     |  |  +--rw max-auth-failures?   uint8
|     |  |  +--rw blacklist-time?      uint16
|     |  +--ro state
|     |     +--ro max-auth-failures?   uint8
|     |     +--ro blacklist-time?      uint16
|     +--rw band-steering
|        +--rw config
|        |  +--rw band-steering?   boolean
|        |  +--rw steering-rssi?   int8
|        +--ro state
|           +--ro band-steering?   boolean
|           +--ro steering-rssi?   int8
+--rw assigned-ap-managers
   +--rw ap-manager* [id]
      +--rw id        -> ../config/id
      +--rw config
      |  +--rw id?                        string
      |  +--rw fqdn?                      oc-inet:domain-name
      |  +--rw ap-manager-ipv4-address?   oc-inet:ipv4-address
      |  +--rw ap-manager-ipv6-address*   oc-inet:ipv6-address
      +--ro state
         +--ro id?                        string
         +--ro fqdn?                      oc-inet:domain-name
         +--ro ap-manager-ipv4-address?   oc-inet:ipv4-address
         +--ro ap-manager-ipv6-address*   oc-inet:ipv6-address
         +--ro joined?                    boolean
*/
// clang-format on

} // namespace wifi
} // namespace models
} // namespace devmand
