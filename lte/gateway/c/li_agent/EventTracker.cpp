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

#include "EventTracker.h"

#include "magma_logging.h"

namespace magma {
namespace lte {

/* close pcap file */
static void pcap_close(pcap_t* p) {

}

static void packet_handler(
    u_char* user, const struct pcap_pkthdr* phdr, const u_char* pdata) {
  magma::DPIEngine* dpi_instance = reinterpret_cast<magma::DPIEngine*>(user);
  struct qmdpi_result* result;
  int ret;

  dpi_instance->process_packet(pdata, phdr->caplen, &phdr->ts);
}

int init_iface_pcap_monitor(int argc, char* argv[]) {
  char errbuf[PCAP_ERRBUF_SIZE];
  const char* license;
  pcap_t* pcap;
  int ret;

  magma::init_logging(argv[0]);
  magma::set_verbosity(MINFO);

  auto packet_encapsulator = std::make_shared<magma::PacketEncapsulator>();
  // Create response loop for flow response processing
  std::thread packet_encapsulator_thread([&]() {
    std::cout << "Started flow response handler thread\n";
    packet_encapsulator->rpc_response_loop();
  });

  license = getenv("DPI_LICENSE");
  if (license == NULL) {
    MLOG(MFATAL) << "Could not load DPI license env var";
    return 1;
  }
  ret = qmdpi_license_load_from_file(license);
  if (ret < 0) {
    MLOG(MFATAL) << "Could not load DPI license with license file " << license
                 << ", error code - " << ret;
    return -1;
  }

  // Snoop on mon1
  pcap = pcap_open_live("mon1", BUFSIZ, 0, 1000, errbuf);
  if (pcap == nullptr) {
    MLOG(MFATAL) << "Could not capture packets, exiting";
    return 1;
  }
  std::cout << "Successfully started live pcap sniffing" << std::endl;

  magma::service303::MagmaService server(DPID_SERVICE, DPID_VERSION);
  server.Start();

  ret = pcap_loop(
      pcap, -1, packet_handler, reinterpret_cast<u_char*>(&dpi_instance));

  if (ret == -1) {
    MLOG(MINFO) << "Could not capture packets";
    if (p != nullptr) {
        pcap_close(pcap);
    }
    return 1;
  }

  server.Stop();
  packet_encapsulator_thread.join();
  return 0;
}

}  // namespace lte
}  // namespace magma
