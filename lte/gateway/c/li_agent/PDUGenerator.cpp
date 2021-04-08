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
#include <cassert>
#include <iostream>
#include <string>

#include <netinet/in.h>
#include <netinet/ip.h>
#include <net/if.h>
#include <net/ethernet.h>
#include <arpa/inet.h>

#include <tins/tins.h>
#include <tins/ip.h>

#include "PDUGenerator.h"

const PDU::PDUType LIX3PDU::pdu_flag = PDU::USER_DEFINED_PDU;
#define LI_X3_LINK_TYPE 0x08ae
#define BASE_HEADER_LEN 40  // 40 octets is the min header length

void print_bytes(void* ptr, int size) {
  unsigned char* p = (unsigned char*) ptr;
  int i;
  for (i = 0; i < size; i++) {
    printf("%02hhX ", p[i]);
  }
  printf("\n");
}

void append_to_vector(std::vector<uint8_t>* x, std::vector<uint8_t>* y) {
  x->reserve(x->size() + distance(y->begin(), y->end()));
  x->insert(x->end(), y->begin(), y->end());
}

void set_xid(struct pdu_info* pdu, uint32_t) {}

namespace magma {
namespace lte {

using namespace Tins;

PDUGenerator::PDUGenerator(
    std::shared_ptr<ProxyConnector> proxy_connector,
    const std::string& pkt_dst_mac, const std::string& pkt_src_mac)
    : pkt_dst_mac_(pkt_dst_mac),
      pkt_src_mac_(pkt_src_mac),
      proxy_connector_(proxy_connector) {
  MLOG(MINFO) << "Using interface for pkt generation";
  // iface_ = NetworkInterface("testy1");
  // Don't know why but changing ethernet-> IP produces an error, resolve
  Allocators::register_allocator<EthernetII, LIX3PDU>(LI_X3_LINK_TYPE);
  directoryd_client_ = std::make_shared<magma::AsyncDirectorydClient>();
}

std::vector<uint8_t> PDUGenerator::get_conditional_attr(void) {
  struct conditional_attributes attributes;
  attributes.timestamp = 678;

  std::vector<uint8_t> conditional_attr(
      (uint8_t*) &attributes,
      (uint8_t*) &attributes + sizeof(struct conditional_attributes));

  return conditional_attr;
}

bool extract_ip_addr(
    const u_char* packet, std::string* src_ip, std::string* dst_ip) {
  const struct ether_header* ethernetHeader;
  const struct ip* ipHeader;
  // char sourceIP[INET_ADDRSTRLEN];
  // char destIP[INET_ADDRSTRLEN];

  ethernetHeader = (struct ether_header*) packet;
  if (ntohs(ethernetHeader->ether_type) != ETHERTYPE_IP) {
    ipHeader = (struct ip*) (packet + sizeof(struct ether_header));
    inet_ntop(AF_INET, &(ipHeader->ip_src), (char*) src_ip, INET_ADDRSTRLEN);
    inet_ntop(AF_INET, &(ipHeader->ip_dst), (char*) dst_ip, INET_ADDRSTRLEN);
    // src_ip.assign(sourceIP, INET_ADDRSTRLEN);
    // dst_ip.assign(destIP, INET_ADDRSTRLEN);
    return true;
  } else {
    return false;
  }
}

bool PDUGenerator::send_packet(
    const struct pcap_pkthdr* phdr, const u_char* pdata) {
  PacketSender sender;

  struct pdu_info pdu;

  pdu.version  = 2;
  pdu.pdu_type = 2;
  pdu.header_length =
      sizeof(struct pdu_info) + sizeof(struct conditional_attributes);
  // pdu.header_length = BASE_HEADER_LEN + sizeof(struct
  // conditional_attributes);
  pdu.payload_length    = phdr->len;
  pdu.payload_format    = 5;  // 5 for ipv4 packets
  pdu.payload_direction = 1;  // TODO set 2/3, 2 is to target, 3 is from target

  std::vector<uint8_t> data(
      (uint8_t*) &pdu, (uint8_t*) &pdu + sizeof(struct pdu_info));

  // TODO extract both IPs, check both in directoryd
  std::string src_ip;
  std::string dst_ip;
  extract_ip_addr(pdata, &src_ip, &dst_ip);

  auto request = directoryd_client_->get_directoryd_xid_field(
      src_ip, [this, src_ip](Status status, DirectoryField resp) {
        if (!status.ok()) {
          MLOG(MERROR) << "Could not fetch subscriber with ip - " << src_ip;
        } else {
          printf("REsp value %s", resp.value().c_str());
        }
      });
  if (!request) {
    MLOG(MERROR) << "Could not query directoryd for ip - " << src_ip;
  }

  /* Append the conditional attributes to the header */
  std::vector<uint8_t> cond = get_conditional_attr();
  append_to_vector(&data, &cond);

  /* Append the packet to the header */
  std::vector<uint8_t> pkt_data((uint8_t*) pdata, (uint8_t*) pdata + phdr->len);

  append_to_vector(&data, &pkt_data);
  printf("len is %lu\n", sizeof(struct pdu_info));
  printf("cond is %lu\n", sizeof(struct conditional_attributes));
  printf("pdu is %u\n", pdu.payload_length);
  printf("data is %lu\n", data.size());

  print_bytes(data.data(), sizeof(struct pdu_info));
  // Random mac header
  // EthernetII eth_(pkt_dst_mac_, pkt_src_mac_);
  // eth_ /= LIX3PDU(data.data(),  pdu.header_length + pdu.payload_length);

  // sender.send(eth_, iface_);
  proxy_connector_->SendData(
      data.data(), pdu.header_length + pdu.payload_length);

  return true;
}

}  // namespace lte
}  // namespace magma