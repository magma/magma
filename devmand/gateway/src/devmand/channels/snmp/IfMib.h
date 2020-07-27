/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once

#include <functional>

#include <folly/futures/Future.h>

#include <devmand/channels/snmp/Channel.h>

namespace devmand {
namespace channels {
namespace snmp {

using Location = std::string;
using Contact = std::string;

struct InterfacePair {
  int index;
  std::string value;
};

using InterfacePairs = std::vector<InterfacePair>;
using InterfaceIndicies = std::vector<int>;

class IfMib {
 public:
  IfMib() = delete;
  ~IfMib() = delete;
  IfMib(const IfMib&) = delete;
  IfMib& operator=(const IfMib&) = delete;
  IfMib(IfMib&&) = delete;
  IfMib& operator=(IfMib&&) = delete;

 public:
  static folly::Future<int> getNumberOfInterfaces(
      channels::snmp::Channel& channel);
  static folly::Future<std::string> getSystemName(
      channels::snmp::Channel& channel);
  static folly::Future<Contact> getSystemContact(
      channels::snmp::Channel& channel);
  static folly::Future<Location> getSystemLocation(
      channels::snmp::Channel& channel);
  static folly::Future<InterfaceIndicies> getInterfaceIndicies(
      channels::snmp::Channel& channel);
  static folly::Future<InterfacePairs> getInterfaceNames(
      channels::snmp::Channel& channel,
      const InterfaceIndicies& indices);
  static folly::Future<InterfacePairs> getInterfaceOperStatuses(
      channels::snmp::Channel& channel,
      const InterfaceIndicies& indices);
  static folly::Future<InterfacePairs> getInterfaceAdminStatuses(
      channels::snmp::Channel& channel,
      const InterfaceIndicies& indices);
  static folly::Future<InterfacePairs> getInterfaceMtus(
      channels::snmp::Channel& channel,
      const InterfaceIndicies& indices);
  static folly::Future<InterfacePairs> getInterfaceTypes(
      channels::snmp::Channel& channel,
      const InterfaceIndicies& indices);
  static folly::Future<InterfacePairs> getInterfaceDescriptions(
      channels::snmp::Channel& channel,
      const InterfaceIndicies& indices);
  static folly::Future<InterfacePairs> getInterfaceLastChange(
      channels::snmp::Channel& channel,
      const InterfaceIndicies& indices);

  static folly::Future<InterfacePairs> getInterfaceField(
      channels::snmp::Channel& channel,
      const std::string& oid,
      const std::function<std::string(std::string)>& formatter = nullptr);

  static folly::Future<InterfacePairs> getInterfaceField(
      channels::snmp::Channel& channel,
      const InterfaceIndicies& indices,
      const std::string& oid,
      const std::function<std::string(std::string)>& formatter = nullptr);
};

} // namespace snmp
} // namespace channels
} // namespace devmand
