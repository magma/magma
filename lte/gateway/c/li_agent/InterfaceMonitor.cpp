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

#include "InterfaceMonitor.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <iostream>

#include <libmnl/libmnl.h>
#include <linux/netfilter/nfnetlink.h>
#include <linux/netfilter/nfnetlink_conntrack.h>

#include <linux/if_packet.h>
#include <string.h>
#include <sys/ioctl.h>
#include <sys/socket.h>
#include <net/if.h>
#include <netinet/ether.h>
#include <linux/ip.h>
#include <memory>
#include <pcap.h>

#include "magma_logging.h"

namespace magma {

InterfaceMonitor::InterfaceMonitor(
    const std::string& iface_name, std::unique_ptr<PDUGenerator> pkt_gen)
    : iface_name_(iface_name), pkt_gen_(std::move(pkt_gen)) {}

static void packet_handler(
    u_char* user, const struct pcap_pkthdr* phdr, const u_char* pdata) {
  reinterpret_cast<PDUGenerator*>(user)->send_packet(phdr, pdata);
}

int InterfaceMonitor::init_iface_pcap_monitor() {
  char errbuf[PCAP_ERRBUF_SIZE];
  pcap_t* pcap;
  int ret;

  pcap = pcap_open_live(
      iface_name_.c_str(), MAX_PKT_SIZE, PROMISCUOUS_MODE,
      PKT_BUF_READ_TIMEOUT_MS, errbuf);
  if (pcap == nullptr) {
    MLOG(MFATAL) << "Could not capture packets on " << iface_name_
                 << ", exiting";
    return -1;
  }
  MLOG(MINFO) << "Successfully started live pcap sniffing";

  ret = pcap_loop(
      pcap, -1, packet_handler, reinterpret_cast<u_char*>(pkt_gen_.get()));

  if (ret == -1) {
    MLOG(MERROR) << "Could not capture packets";
    if (pcap != nullptr) {
      pcap_close(pcap);
    }
    return -1;
  }

  return 0;
}

}  // namespace magma
