/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#pragma once

#include "PDUGenerator.h"

namespace magma {
namespace lte {

class InterfaceMonitor {
 public:
  InterfaceMonitor(
      const std::string& iface_name, std::shared_ptr<PDUGenerator> pkt_gen);

  int init_iface_pcap_monitor(void);

 private:
  std::string iface_name_;
  std::shared_ptr<PDUGenerator> pkt_gen_;
};

}  // namespace lte
}  // namespace magma