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

#include "lte/gateway/c/li_agent/src/PDUGenerator.h"
#include "lte/gateway/c/li_agent/src/Utilities.h"

#include <uuid/uuid.h>
#include <netinet/ip.h>
#include <net/ethernet.h>

#include <future>
#include <string>
#include <memory>
#include <utility>

namespace magma {
namespace lte {

#define ETHERNET_HDR_LEN 14
#define MAX_EXPORT_RETRIES 3

#define PDU_TYPE 2
#define PDU_VERSION 2
#define IP_PAYLOAD_FORMAT 5
#define DIRECTION_TO_TARGET 2
#define DIRECTION_FROM_TARGET 3

#define SEQNBR_ATTRID 8
#define TIMESTAMP_ATTRID 9

#define SET_INT64_TLV(tlv, id, value)      \
  do {                                     \
    (tlv)->type = htons(id);               \
    (tlv)->data = htobe64(value);          \
    (tlv)->size = htons(sizeof(uint64_t)); \
  } while (0)

FlowInformation extract_flow_information(const u_char* packet) {
  FlowInformation ret;
  char src[INET_ADDRSTRLEN];
  char dst[INET_ADDRSTRLEN];

  const struct ip* iphdr;
  ret.successful = false;
  const struct ether_header* ethhdr = (struct ether_header*)packet;
  if (ntohs(ethhdr->ether_type) == ETHERTYPE_IP) {
    iphdr = (struct ip*)(packet + sizeof(struct ether_header));
    ret.src_ip = inet_ntop(AF_INET, &(iphdr->ip_src), src, INET_ADDRSTRLEN);
    ret.dst_ip = inet_ntop(AF_INET, &(iphdr->ip_dst), dst, INET_ADDRSTRLEN);
    ret.successful = true;
  }
  return ret;
}

static InterceptState build_new_intercept_state(
    std::string subid, const magma::mconfig::NProbeTask& task) {
  MLOG(MDEBUG) << "Create new intercept state for task " << task.task_id();
  InterceptState state;
  state.target_id = subid;
  state.task_id = task.task_id();
  state.domain_id = task.domain_id();
  state.correlation_id = task.correlation_id();
  state.sequence_number = 0;
  return state;
}

PDUGenerator::PDUGenerator(const std::string& pkt_dst_mac,
                           const std::string& pkt_src_mac, int sync_interval,
                           int inactivity_time,
                           std::unique_ptr<ProxyConnector> proxy_connector,
                           std::unique_ptr<MobilitydClient> mobilityd_client,
                           magma::mconfig::LIAgentD mconfig)
    : pkt_dst_mac_(pkt_dst_mac),
      pkt_src_mac_(pkt_src_mac),
      sync_interval_(sync_interval),
      inactivity_time_(inactivity_time),
      prev_sync_time_(0),
      proxy_connector_(std::move(proxy_connector)),
      mobilityd_client_(std::move(mobilityd_client)),
      mconfig_(mconfig) {}

bool PDUGenerator::process_packet(const struct pcap_pkthdr* phdr,
                                  const u_char* pdata) {
  FlowInformation flow = extract_flow_information(pdata);
  if (!flow.successful) {
    MLOG(MERROR)
        << "Could not extract flow information from the packet, skipping";
    return false;
  }

  auto diff = time_difference_from_now(prev_sync_time_);
  if (diff > static_cast<uint64_t>(sync_interval_)) {
    // load mconfig config to get updated nprobe tasks
    mconfig_ = magma::lte::load_mconfig();
    prev_sync_time_ = get_time_in_sec_since_epoch();
  }

  std::string idx;
  if (get_intercept_state_idx(flow, &idx) == false) {
    MLOG(MERROR) << "Could not find subscriber for src ip - " << flow.src_ip
                 << ", and dst ip - " << flow.dst_ip;
    return false;
  }

  uint32_t rlen;
  uint16_t direction =
      (idx == flow.src_ip) ? DIRECTION_FROM_TARGET : DIRECTION_TO_TARGET;

  void* record = generate_record(phdr, pdata, idx, direction, &rlen);
  if (record == nullptr) {
    return false;
  }

  auto exported = export_record(record, rlen, MAX_EXPORT_RETRIES);
  MLOG(MDEBUG) << "Exported packet " << state_map_[idx].sequence_number
               << " with length " << rlen;
  free(record);
  return exported;
}

void PDUGenerator::delete_inactive_tasks() {
  auto diff = time_difference_from_now(prev_sync_time_);
  if (diff < static_cast<uint64_t>(sync_interval_)) {
    return;
  }

  auto it = state_map_.begin();
  while (it != state_map_.end()) {
    auto inactive = time_difference_from_now(it->second.last_exported);
    if (inactive > static_cast<uint64_t>(inactivity_time_)) {
      MLOG(MDEBUG) << "Delete state for task " << it->second.task_id;
      it = state_map_.erase(it);
    } else {
      it++;
    }
  }
  return;
}

void* PDUGenerator::generate_record(const struct pcap_pkthdr* phdr,
                                    const u_char* pdata, std::string idx,
                                    uint16_t direction, uint32_t* record_len) {
  auto& state = state_map_[idx];
  uint32_t hdr_len = sizeof(X3Header);
  uint32_t pld_len =
      phdr->len -
      ETHERNET_HDR_LEN;  // Skip eth layer as defined in ETSI 103 221-2.

  *record_len = hdr_len + pld_len;
  uint8_t* record = static_cast<uint8_t*>(calloc(1, *record_len));
  if (record == nullptr) {
    MLOG(MERROR) << "Failed to allocate memory " << *record_len;
    *record_len = 0;
    return nullptr;
  }

  X3Header* pdu = reinterpret_cast<X3Header*>(record);
  pdu->version = htons(PDU_VERSION);
  pdu->pdu_type = htons(PDU_TYPE);
  pdu->header_length = htonl(hdr_len);
  pdu->payload_length = htonl(pld_len);
  pdu->payload_format = htons(IP_PAYLOAD_FORMAT);
  pdu->correlation_id = htobe64(state.correlation_id);
  pdu->payload_direction = htons(direction);

  uint64_t tm = (uint64_t)phdr->ts.tv_sec << 32 | phdr->ts.tv_usec;
  SET_INT64_TLV(&pdu->attrs.timestamp, TIMESTAMP_ATTRID, tm);

  SET_INT64_TLV(&pdu->attrs.sequence_number, SEQNBR_ATTRID,
                state.sequence_number);

  auto ret = uuid_parse(state.task_id.c_str(), pdu->xid);
  if (ret != 0) {
    MLOG(MERROR) << "Failed to parse task_id " << state.task_id.c_str();
    free(record);
    *record_len = 0;
    return nullptr;
  }

  memcpy(record + hdr_len, pdata + ETHERNET_HDR_LEN, pld_len);
  state.last_exported = phdr->ts.tv_sec;
  state.sequence_number++;

  return reinterpret_cast<void*>(record);
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

bool PDUGenerator::get_subscriber_id_from_ip(const char* ip_addr,
                                             std::string* subid) {
  struct in_addr addr;
  if (inet_aton(ip_addr, &addr) <= 0) {
    MLOG(MERROR) << "Bad IPv4 address format " << ip_addr;
    return false;
  }

  std::promise<std::string> lookup_res;
  std::future<std::string> lookup_future = lookup_res.get_future();
  mobilityd_client_->get_subscriber_id_from_ip(
      addr,
      [this, &addr, &lookup_res, ip_addr](Status status, SubscriberID resp) {
        if (!status.ok()) {
          MLOG(MDEBUG) << "Could not find subscriber_id for ip " << ip_addr;
          lookup_res.set_value("");
          return;
        }
        MLOG(MDEBUG) << "Found subscriber " << resp.id() << " for ip "
                     << ip_addr;
        lookup_res.set_value(resp.id());
      });

  std::string subid_str = lookup_future.get();
  if (subid_str.empty()) {
    return false;
  }

  if (subid_str.find("IMSI") == std::string::npos) {
    subid_str = "IMSI" + subid_str;
  }
  *subid = subid_str.c_str();
  return true;
}

bool PDUGenerator::get_intercept_state_idx(const FlowInformation& flow,
                                           std::string* idx) {
  if (state_map_.find(flow.src_ip) != state_map_.end()) {
    *idx = flow.src_ip;
  } else if (state_map_.find(flow.dst_ip) != state_map_.end()) {
    *idx = flow.dst_ip;
  }

  if (!idx->empty()) {
    if (is_still_valid_state(*idx)) {
      return true;
    }
    MLOG(MDEBUG) << "Delete invalid state for " << idx;
    state_map_.erase(*idx);
  }
  return create_new_intercept_state(flow, idx);
}

bool PDUGenerator::create_new_intercept_state(const FlowInformation& flow,
                                              std::string* idx) {
  std::string subid;
  if (get_subscriber_id_from_ip(flow.src_ip.c_str(), &subid)) {
    *idx = flow.src_ip;
  } else if (get_subscriber_id_from_ip(flow.dst_ip.c_str(), &subid)) {
    *idx = flow.dst_ip;
  } else {
    return false;
  }

  for (const auto& it : mconfig_.nprobe_tasks()) {
    if (it.target_id() == subid) {
      state_map_[*idx] = build_new_intercept_state(subid, it);
      return true;
    }
  }
  return false;
}

bool PDUGenerator::is_still_valid_state(const std::string& idx) {
  auto diff = time_difference_from_now(state_map_[idx].last_exported);
  if (diff < static_cast<uint64_t>(sync_interval_)) {
    return true;
  }

  auto& state = state_map_[idx];
  for (const auto& task : mconfig_.nprobe_tasks()) {
    if (state.task_id == task.task_id()) {
      MLOG(MDEBUG) << "Found task - " << state.task_id;
      state.correlation_id = task.correlation_id();
      state.domain_id = task.domain_id();
      return true;
    }
  }
  return false;
}

}  // namespace lte
}  // namespace magma
