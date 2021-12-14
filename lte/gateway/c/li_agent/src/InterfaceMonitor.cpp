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

#include <stdio.h>
#include <unistd.h>
#include <utility>

#include "lte/gateway/c/li_agent/src/InterfaceMonitor.h"
#include "orc8r/gateway/c/common/logging/magma_logging.h"

namespace magma {
namespace lte {

InterfaceMonitor::InterfaceMonitor(const std::string& iface_name,
                                   std::unique_ptr<PDUGenerator> pkt_gen)
    : pcap_(nullptr), iface_name_(iface_name), pkt_gen_(std::move(pkt_gen)) {}

static void packet_handler(u_char* user, const struct pcap_pkthdr* phdr,
                           const u_char* pdata) {
  reinterpret_cast<PDUGenerator*>(user)->process_packet(phdr, pdata);
}

int InterfaceMonitor::init_interface_monitor() {
  char errbuf[PCAP_ERRBUF_SIZE];
  int ret;

  pcap_ = pcap_open_live(iface_name_.c_str(), MAX_PKT_SIZE, PROMISCUOUS_MODE,
                         PKT_BUF_READ_TIMEOUT_MS, errbuf);
  if (pcap_ == nullptr) {
    MLOG(MFATAL) << "Could not capture packets on " << iface_name_
                 << ", exiting";
    return -1;
  }

  ret = pcap_setnonblock(pcap_, 1, errbuf);
  if (ret == -1) {
    MLOG(MFATAL) << "Could not set non-blocking mode " << ret << ", exiting";
    return -1;
  }
  MLOG(MINFO) << "Successfully started live pcap sniffing";
  return 0;
}

int InterfaceMonitor::start_capture() {
  int ret;
  while (true) {
    ret = pcap_dispatch(pcap_, -1, packet_handler,
                        reinterpret_cast<u_char*>(pkt_gen_.get()));
    if (ret == -1) {
      MLOG(MERROR) << "Could not capture packets";
      if (pcap_ != nullptr) {
        pcap_close(pcap_);
        pcap_ = nullptr;
      }
      return -1;
    } else if (ret == 0) {
      pkt_gen_->delete_inactive_tasks();
      usleep(100);
    }
  }
  return 0;
}

}  // namespace lte
}  // namespace magma
