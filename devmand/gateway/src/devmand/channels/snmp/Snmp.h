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

#include <string>

#include <devmand/channels/snmp/Response.h>

namespace devmand {
namespace channels {
namespace snmp {

// name or address of peer (may include transport specifier and/or port number)
using Peer = std::string;
using Community = std::string;
using Version = std::string;
using SecurityLevel = std::string;

} // namespace snmp
} // namespace channels
} // namespace devmand
