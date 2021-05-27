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

#include <string>
#include <memory>
#include <vector>
#include <uuid/uuid.h>
#include <netinet/ip.h>
#include <net/ethernet.h>

#include "PDUGenerator.h"
#include "MConfigLoader.h"
#include "Utilities.h"

namespace magma {
namespace lte {

#define ETHERNET_HDR_LEN 14
#define MAX_EXPORT_RETRIES 3

#define PDU_TYPE 2
#define PDU_VERSION 2
#define IP_PAYLOAD_FORMAT 5
#define DIRECTION_TO_TARGET 2
#define DIRECTION_FROM_TARGET 3

#define SEQUENCE_NUMBER_ATTRID 8
#define TIMESTAMP_ATTRID 9

#define SET_INT64_TLV(tlv, id, value)                                          \
  do {                                                                         \
    (tlv)->type = htons(id);                                                   \
    (tlv)->data = htobe64(value);                                              \
    (tlv)->size = htons(sizeof(uint64_t));                                     \
  } while (0)

FlowInformation extract_flow_information(const u_char* packet) {
  FlowInformation ret;
  char src[INET_ADDRSTRLEN];
  char dst[INET_ADDRSTRLEN];

  const struct ip* iphdr;
  ret.successful                    = false;
  const struct ether_header* ethhdr = (struct ether_header*) packet;
  if (ntohs(ethhdr->ether_type) == ETHERTYPE_IP) {
    iphdr          = (struct ip*) (packet + sizeof(struct ether_header));
    ret.src_ip     = inet_ntop(AF_INET, &(iphdr->ip_src), src, INET_ADDRSTRLEN);
    ret.dst_ip     = inet_ntop(AF_INET, &(iphdr->ip_dst), dst, INET_ADDRSTRLEN);
    ret.successful = true;
  }
  return ret;
}

PDUGenerator::PDUGenerator(
    const std::string& pkt_dst_mac, const std::string& pkt_src_mac,
    int sync_interval, int inactivity_time,
    std::unique_ptr<ProxyConnector> proxy_connector,
    std::unique_ptr<MobilitydClient> mobilityd_client)
    : pkt_dst_mac_(pkt_dst_mac),
      pkt_src_mac_(pkt_src_mac),
      sync_interval_(sync_interval),
      inactivity_time_(inactivity_time),
      proxy_connector_(std::move(proxy_connector)),
      mobilityd_client_(std::move(mobilityd_client)) {}

bool PDUGenerator::get_subscriber_id_from_ip(
    const char* ip_addr, std::string* subid) {
  struct in_addr addr;
  if (inet_aton(ip_addr, &addr) < 0) {
    return false;
  }

  std::string subid_str;
  int status = mobilityd_client_->GetSubscriberIDFromIP(addr, &subid_str);
  if (subid_str.empty()) {
    return false;
  }

  if (subid_str.find("IMSI") == std::string::npos) {
    subid_str = "IMSI" + subid_str;
  }
  *subid = strdup(subid_str.c_str());
  return true;
}

bool PDUGenerator::is_still_valid_state(std::string idx) {
  auto& state  = intercept_state_map_[idx];
  auto mconfig = magma::lte::load_mconfig();
  for (const auto& task : mconfig.nprobe_tasks()) {
    if (state.task_id == task.task_id()) {
      MLOG(MDEBUG) << "Found task - " << state.task_id;
      state.correlation_id = task.correlation_id();
      state.domain_id      = task.domain_id();
      return true;
    }
  }
  return false;
}

std::string PDUGenerator::get_intercept_state_idx(const FlowInformation& flow) {
  std::string idx;
  if (intercept_state_map_.find(flow.src_ip) != intercept_state_map_.end()) {
    idx = flow.src_ip;
  } else if (
      intercept_state_map_.find(flow.dst_ip) != intercept_state_map_.end()) {
    idx = flow.dst_ip;
  }

  if (!idx.empty()) {
    auto now = get_time_in_sec_since_epoch();
    if (now - intercept_state_map_[idx].last_exported < sync_interval_)
      if (is_still_valid_state(idx)) {
        return idx;
      }
    MLOG(MDEBUG) << "Delete task " << idx;
    intercept_state_map_.erase(idx);
  }
  return create_new_intercept_state(flow);
}

std::string PDUGenerator::create_new_intercept_state(
    const FlowInformation& flow) {
  std::string subid, idx;
  if (get_subscriber_id_from_ip(flow.src_ip.c_str(), &subid)) {
    idx = flow.src_ip;
  } else if (get_subscriber_id_from_ip(flow.dst_ip.c_str(), &subid)) {
    idx = flow.dst_ip;
  } else {
    MLOG(MERROR) << "Could not find subscriber_id for src ip - " << flow.src_ip
                 << ", and dst ip - " << flow.dst_ip;
    return idx;
  }

  auto mconfig = magma::lte::load_mconfig();
  for (const auto& it : mconfig.nprobe_tasks()) {
    if (it.target_id() == subid) {
      MLOG(MDEBUG) << "Create new task " << it.task_id();
      InterceptState state;
      state.target_id           = subid;
      state.task_id             = it.task_id();
      state.domain_id           = it.domain_id();
      state.correlation_id      = it.correlation_id();
      intercept_state_map_[idx] = state;
      break;
    }
  }
  return idx;
}

void* PDUGenerator::generate_record(
    const struct pcap_pkthdr* phdr, const u_char* pdata, std::string idx,
    uint16_t direction, uint32_t* record_len) {
  auto& state      = intercept_state_map_[idx];
  uint32_t hdr_len = sizeof(X3Header);
  uint32_t pld_len =
      phdr->len -
      ETHERNET_HDR_LEN;  // Skip eth layer as defined in ETSI 103 221-2.

  *record_len   = hdr_len + pld_len;
  uint8_t* data = static_cast<uint8_t*>(calloc(1, *record_len));

  X3Header* pdu          = reinterpret_cast<X3Header*>(data);
  pdu->version           = htons(PDU_VERSION);
  pdu->pdu_type          = htons(PDU_TYPE);
  pdu->header_length     = htonl(hdr_len);
  pdu->payload_length    = htonl(pld_len);
  pdu->payload_format    = htons(IP_PAYLOAD_FORMAT);
  pdu->correlation_id    = htobe64(state.correlation_id);
  pdu->payload_direction = htons(direction);

  uuid_parse(state.task_id.c_str(), pdu->xid);
  SET_INT64_TLV(&pdu->attrs.timestamp, TIMESTAMP_ATTRID, phdr->ts.tv_sec);
  SET_INT64_TLV(
      &pdu->attrs.sequence_number, SEQUENCE_NUMBER_ATTRID,
      state.sequence_number);

  memcpy(data + hdr_len, pdata + ETHERNET_HDR_LEN, pld_len);

  state.last_exported = phdr->ts.tv_sec;
  state.sequence_number++;
  return (void*) data;
}

bool PDUGenerator::process_packet(
    const struct pcap_pkthdr* phdr, const u_char* pdata) {
  FlowInformation flow = extract_flow_information(pdata);
  if (!flow.successful) {
    MLOG(MERROR)
        << "Could not extract flow information from the packet, skipping";
    return false;
  }

  auto idx = get_intercept_state_idx(flow);
  if (idx.empty()) return false;

  uint16_t direction =
      (idx == flow.src_ip) ? DIRECTION_FROM_TARGET : DIRECTION_TO_TARGET;

  uint32_t record_len;
  void* record = generate_record(phdr, pdata, idx, direction, &record_len);
  auto ret     = export_record(record, record_len, MAX_EXPORT_RETRIES);

  MLOG(MDEBUG) << "Generated packet "
               << intercept_state_map_[idx].sequence_number
               << " length with length " << record_len;

  free(record);
  return ret;
}

bool PDUGenerator::export_record(void* record, uint32_t size, int retries) {
  for (auto i = 0; i < retries; i++) {
    int ret = proxy_connector_->send_data(record, size);
    if (ret > 0) {
      break;
    } else {
      proxy_connector_->cleanup();
      if (proxy_connector_->setup_proxy_socket() < 0) {
        return false;
      }
    }
  }
  return true;
}

void PDUGenerator::cleanup_inactive_tasks() {
  if (time_difference_from_now(prev_sync_time_) < sync_interval_) {
    return;
  }

  auto it = intercept_state_map_.begin();
  while (it != intercept_state_map_.end()) {
    if (time_difference_from_now(it->second.last_exported) > inactivity_time_) {
      MLOG(MDEBUG) << "Delete task " << it->second.task_id;
      it = intercept_state_map_.erase(it);
    } else {
      it++;
    }
  }
  return;
}

}  // namespace lte
}  // namespace magma
