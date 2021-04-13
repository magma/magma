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
#include <tins/pdu.h>

#include "magma_logging.h"
#include "DirectorydClient.h"
#include "ProxyConnector.h"

struct pdu_info {
  uint16_t version;
  uint16_t pdu_type;
  uint32_t header_length;
  uint32_t payload_length;
  uint16_t payload_format;
  uint16_t payload_direction;
  uint8_t xid[16];
};

struct conditional_attributes {
  uint64_t timestamp;
};

using namespace Tins;
class LIX3PDU : public Tins::PDU {
 public:
  /*
   * Unique protocol identifier. For user-defined PDUs, you **must**
   * use values greater or equal to PDU::USER_DEFINED_PDU;
   */
  static const PDU::PDUType pdu_flag;

  /*
   * Constructor from buffer. This constructor will be called while
   * sniffing packets, whenever a PDU of this type is found.
   *
   * The "data" parameter points to a buffer of length "sz".
   */
  LIX3PDU(const uint8_t* data, uint32_t sz) : buffer_(data, data + sz) {}

  /*
   * Clones the PDU. This method is used when copying PDUs.
   */
  LIX3PDU* clone() const { return new LIX3PDU(*this); }

  /*
   * Retrieves the size of this PDU.
   */
  uint32_t header_size() const { return buffer_.size(); }

  /*
   * This method must return pdu_flag.
   */
  PDUType pdu_type() const { return pdu_flag; }

  /*
   * Serializes the PDU. The serialization output should be written
   * to the buffer pointed to by "data", which is of size "sz". The
   * "sz" parameter will be equal to the value returned by
   * LIX3PDU::header_size.
   */
  void write_serialization(uint8_t* data, uint32_t sz) {
    std::memcpy(data, buffer_.data(), sz);
  }

  /*
   * For libtins 4.0 and lower
   */
  void write_serialization(
      uint8_t* buffer, uint32_t total_sz, const PDU* parent) {
    std::memcpy(buffer, buffer_.data(), total_sz - parent->header_size());
  }

  // This is just a getter to retrieve the buffer member.
  const std::vector<uint8_t>& get_buffer() const { return buffer_; }

 private:
  std::vector<uint8_t> buffer_;
};

namespace magma {
namespace lte {

class PDUGenerator {
 public:
  PDUGenerator(
      std::shared_ptr<ProxyConnector> proxy_connector,
      std::shared_ptr<AsyncDirectorydClient> directoryd_client,
      const std::string& pkt_dst_mac, const std::string& pkt_src_mac);

  std::vector<uint8_t> get_conditional_attr(void);

  /**
   * Send packet
   * @param flow_information - flow_information
   * @return true if the operation was successful
   */
  bool send_packet(const struct pcap_pkthdr* phdr, const u_char* pdata);

 private:
  std::string iface_name_;
  std::string pkt_dst_mac_;
  std::string pkt_src_mac_;
  Tins::NetworkInterface iface_;
  std::shared_ptr<AsyncDirectorydClient> directoryd_client_;
  std::shared_ptr<ProxyConnector> proxy_connector_;
};

}  // namespace lte
}  // namespace magma
