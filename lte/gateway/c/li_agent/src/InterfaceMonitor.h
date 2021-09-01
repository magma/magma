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

#include <pcap.h>

#include <string>
#include <memory>

#include "PDUGenerator.h"

namespace magma {
namespace lte {

#define MAX_PKT_SIZE 2048
#define PROMISCUOUS_MODE 0
#define PKT_BUF_READ_TIMEOUT_MS 1000

class InterfaceMonitor {
 public:
  InterfaceMonitor(
      const std::string& iface_name, std::unique_ptr<PDUGenerator> pkt_gen);

  /**
   * init_iface_pcap_monitor starts a live pcap sniffing for an interface
   * provided in service configuration.
   * @return return positif integer if interface monitoring starts successfully.
   */
  int init_interface_monitor();
  int start_capture();

 private:
  pcap_t* pcap_;
  std::string iface_name_;
  std::unique_ptr<PDUGenerator> pkt_gen_;
};

}  // namespace lte
}  // namespace magma
