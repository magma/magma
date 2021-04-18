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

#include <tins/tins.h>

#include "magma_logging.h"
#include "DirectorydClient.h"
#include "ProxyConnector.h"

#define XID_LENGTH 16
#define LI_X3_LINK_TYPE 0x08ae
#define PDU_VERSION 2
#define PDU_TYPE 2
#define IP_PAYLOAD_FORMAT 5
#define DIRECTION_TO_TARGET 2
#define DIRECTION_FROM_TARGET 3

struct ip_extraction_pair {
  bool successful;
  std::string src_ip;
  std::string dst_ip;
};

struct pdu_info {
  uint16_t version;
  uint16_t pdu_type;
  uint32_t header_length;
  uint32_t payload_length;
  uint16_t payload_format;
  uint16_t payload_direction;
  uint8_t xid[XID_LENGTH];
};

struct conditional_attributes {
  uint64_t timestamp;
};

namespace magma {

class PDUGenerator {
 public:
  PDUGenerator(
      std::unique_ptr<ProxyConnector> proxy_connector,
      std::unique_ptr<DirectorydClient> directoryd_client,
      const std::string& pkt_dst_mac, const std::string& pkt_src_mac);

  /**
   * Send packet
   * @param phdr - packet header
   * @param pdata - packet data
   * @return true if the operation was successful
   */
  bool send_packet(const struct pcap_pkthdr* phdr, const u_char* pdata);

 private:
  std::string iface_name_;
  std::string pkt_dst_mac_;
  std::string pkt_src_mac_;
  Tins::NetworkInterface iface_;
  std::unique_ptr<DirectorydClient> directoryd_client_;
  std::unique_ptr<ProxyConnector> proxy_connector_;

  void set_conditional_attr(
      const struct pcap_pkthdr* phdr,
      struct conditional_attributes* attributes);
  void* generate_pkt(const struct pcap_pkthdr* phdr, const u_char* pdata);
  void handle_ip_lookup_callback(
      std::string src_ip, void* data, struct pdu_info* pdu, Status status,
      DirectoryField resp);
};

}  // namespace magma
