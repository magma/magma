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

#include <string>
#include <uuid/uuid.h>
#include <tins/tins.h>
#include <unordered_map>

#include "magma_logging.h"
#include "ProxyConnector.h"

namespace magma {
namespace lte {

#define XID_LENGTH 16

typedef struct {
  uint16_t type;  // type
  uint16_t size;  // size of data
  uint64_t data;
} __attribute__((__packed__)) TLV;

typedef struct {
  TLV timestamp;
  TLV sequence_number;
} __attribute__((__packed__)) ConditionalAttributes;

typedef struct {
  uint16_t version;
  uint16_t pdu_type;
  uint32_t header_length;
  uint32_t payload_length;
  uint16_t payload_format;
  uint16_t payload_direction;
  uint8_t xid[XID_LENGTH];
  uint64_t correlation_id;
  ConditionalAttributes attrs;
} __attribute__((__packed__)) X3Header;

typedef struct {
  std::string src_ip;
  std::string dst_ip;
  bool successful;
} FlowInformation;

typedef struct {
  std::string task_id;
  std::string target_id;
  std::string domain_id;
  uint64_t last_exported;
  uint64_t correlation_id;
  uint64_t sequence_number;
} InterceptState;

typedef std::unordered_map<std::string, InterceptState> InterceptStateMap;

class PDUGenerator {
 public:
  PDUGenerator(
      const std::string& pkt_dst_mac, const std::string& pkt_src_mac,
      int sync_interval, int inactivity_time,
      std::unique_ptr<ProxyConnector> proxy_connector);

  /**
   * process_packet retrieves the state of the current interception for
   * this packet by looking in the intercept map or interrogating mobility
   * service and creating a new one. Then it generates the corresponding
   * x3 records and exports it to remote destination over TLS.
   * @param phdr - packet header
   * @param pdata - packet data
   * @return true if the operation was successful
   */
  bool process_packet(const struct pcap_pkthdr* phdr, const u_char* pdata);

 private:
  std::string iface_name_;
  std::string pkt_dst_mac_;
  std::string pkt_src_mac_;
  int sync_interval_;
  int inactivity_time_;
  uint64_t prev_sync_time_;
  Tins::NetworkInterface iface_;
  InterceptStateMap intercept_state_map_;
  std::unique_ptr<ProxyConnector> proxy_connector_;

  /**
   * generate_record builds an x3 record from the current packet as specified
   * in ETSI 103 221-2.
   * @param phdr - packet header
   * @param pdata - packet data
   * @param state - intercept state
   * @param direction - direction of packet
   * @param record_len - output record length
   * @return true if the operation was successful
   */
  void* generate_record(
      const struct pcap_pkthdr* phdr, const u_char* pdata, std::string idx,
      uint16_t direction, uint32_t* record_len);

  /**
   * export_record exports the x3 record over tls to a remote server.
   * @param record- x3 record packet
   * @param size - x3 record length
   * @param retries - number of retries
   * @return true if the operation was successful
   */
  bool export_record(void* record, uint32_t size, int retries);

  /**
   * process_new_task processes a new task when no existing state matches the
   * current packet.
   * @param ipList - source and destination ip address
   * @param idx - ip address of the intercept target
   * @param state - intercept state
   * @return voidl
   */
  bool process_new_task(
      const FlowInformation flow, std::string* idx, InterceptState* state);

  /**
   * set_conditional_attr sets the tlv conditional attributes as defined in ETSI
   * 103 221-2.
   * @param tlv - destination tlv
   * @param type - type of attribute is defined by spec
   * @param value - attribute value
   * @return voidl
   */
  void set_conditional_attr(TLV* tlv, uint16_t type, uint64_t value);

  void refresh_intercept_state_map();

  std::string get_intercept_state_idx(const FlowInformation& flow);
  std::string create_new_intercept_state(const FlowInformation& flow);
};

}  // namespace lte
}  // namespace magma
