// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <algorithm>
#include <experimental/filesystem>

#include <folly/Format.h>
#include <folly/GLog.h>

#include <devmand/DhcpdConfig.h>
#include <devmand/FileUtils.h>

namespace devmand {

const char* dhcpSubnet = "10.22.18.0";
const char* dhcpSubnetMask = "255.255.255.128";
const char* dhcpServerIp = "10.22.18.10";
const char* dhcpDefaultGwIp = "10.22.18.1";

const char* installerFileName =
    "EC_AS7316_26XB-OcNOS_SP-1.0.0.324.81a-SP_CSR-S0-P0-installer";

const char* dhcpdConfigFile = "/etc/dhcp/dhcpd.conf";
const char* dhcpdConfigTemplate = R"template(
option ocnos-license-url code 251 = text;
option ocnos-provision-url code 250 = text;

subnet {0} netmask {1} {{
     option routers {2};
     option subnet-mask {1};
     default-lease-time 5555;
     max-lease-time 77777;
}}

{3}
)template";

const char* dhcpdConfigHostTemplate = R"template(
host {0} {{
  hardware ethernet {1};
  fixed-address {2};
  option subnet-mask {3};
  option routers {4};
  option ocnos-license-url = "http://{5}/IPI-80A2354DC6FA.bin";
  option default-url = "http://{5}/{6}";
  option ocnos-provision-url = "http://{5}/{0}.conf";
}}
)template";

DhcpdConfig::DhcpdConfig() {
  FileUtils::mkdir(
      std::experimental::filesystem::path(dhcpdConfigFile).parent_path());
}

void DhcpdConfig::add(Host& host) {
  if (hosts.emplace(host).second) {
    // TODO We are just rewriting every time we get a new host which could be a
    // lot better but it works.
    rewrite();
  } else {
    LOG(ERROR) << "Failed to add host";
  }
}

void DhcpdConfig::remove(Host& host) {
  if (hosts.erase(host) != 1) {
    LOG(ERROR) << "Failed to delete host " << host.name;
  } else {
    rewrite();
  }
}

std::string DhcpdConfig::getHostSection(const Host& host) {
  // auto mac = host.mac.toString();
  // mac.erase(std::remove(mac.begin(), mac.end(), ':'), mac.end());
  return folly::sformat(
      dhcpdConfigHostTemplate,
      host.name,
      host.mac.toString(),
      host.ip.str(),
      dhcpSubnetMask,
      dhcpDefaultGwIp,
      dhcpServerIp,
      installerFileName);
}

void DhcpdConfig::rewrite() {
  std::string hostSections;
  for (auto& host : hosts) {
    hostSections += getHostSection(host);
  }
  std::string dhcpConfigFileContents = folly::sformat(
      dhcpdConfigTemplate,
      dhcpSubnet,
      dhcpSubnetMask,
      dhcpServerIp,
      hostSections);
  FileUtils::write(dhcpdConfigFile, dhcpConfigFileContents);
}

} // namespace devmand
