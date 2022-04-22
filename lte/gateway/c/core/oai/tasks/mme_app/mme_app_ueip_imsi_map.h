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

#include <unordered_map>
#include <iostream>
#include <string>
#include <vector>

/* Description: ue_ip address is allocated by either roaming PGWs or mobilityd
 * So there is possibility to allocate same ue ip address for different UEs.
 * So defining ue_ip_imsi map with key as ue_ip and value as list of imsis
 * having same ue_ip
 */

typedef std::unordered_map<std::string, std::vector<uint64_t>> UeIpImsiMap;
