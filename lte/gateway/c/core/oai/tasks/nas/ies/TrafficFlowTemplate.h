/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#ifndef TRAFFIC_FLOW_TEMPLATE_H_
#define TRAFFIC_FLOW_TEMPLATE_H_
#include <stdint.h>

#define TRAFFIC_FLOW_TEMPLATE_MINIMUM_LENGTH 2
#define TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH 256

/*
 * ----------------------------------------------------------------------------
 *        Packet filter list
 * ----------------------------------------------------------------------------
 */

#define TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR 0b00010000
#define TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR 0b00100000
#define TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER 0b00110000
#define TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT 0b01000000
#define TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE 0b01000001
#define TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT 0b01010000
#define TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE 0b01010001
#define TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX 0b01100000
#define TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS 0b01110000
#define TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL 0b10000000

/*
 * Packet filter content
 * ---------------------
 */
typedef struct {
#define TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG (1 << 0)
#define TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG (1 << 1)
#define TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG (1 << 2)
#define TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG (1 << 3)
#define TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE_FLAG (1 << 4)
#define TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG (1 << 5)
#define TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE_FLAG (1 << 6)
#define TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX_FLAG (1 << 7)
#define TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS_FLAG (1 << 8)
#define TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG (1 << 9)
  uint16_t flags;
#define TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE 4
  struct {
    uint8_t addr;
    uint8_t mask;
  } ipv4remoteaddr[TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE];
#define TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE 16
  struct {
    uint8_t addr;
    uint8_t mask;
  } ipv6remoteaddr[TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE];
  uint8_t protocolidentifier_nextheader;
  uint16_t singlelocalport;
  struct {
    uint16_t lowlimit;
    uint16_t highlimit;
  } localportrange;
  uint16_t singleremoteport;
  struct {
    uint16_t lowlimit;
    uint16_t highlimit;
  } remoteportrange;
  uint32_t securityparameterindex;
  struct {
    uint8_t value;
    uint8_t mask;
  } typdeofservice_trafficclass;
  uint32_t flowlabel;
} PacketFilter;

/*
 * Packet filter list when the TFP operation is "delete existing TFT"
 * and "no TFT operation" shall be empty.
 * ---------------------------------------------------------------
 */
typedef struct {
} NoPacketFilter;

typedef NoPacketFilter DeleteExistingTft;
typedef NoPacketFilter NoTftOperation;

/*
 * Packet filter list when the TFT operation is "delete existing TFT"
 * shall contain a variable number of packet filter identifiers.
 * ------------------------------------------------------------------
 */
#define TRAFFIC_FLOW_TEMPLATE_PACKET_IDENTIFIER_MAX 16
typedef struct {
  uint8_t identifier;
} PacketFilterIdentifiers[TRAFFIC_FLOW_TEMPLATE_PACKET_IDENTIFIER_MAX];

typedef PacketFilterIdentifiers DeletePacketFilter;

/*
 * Packet filter list when the TFT operation is "create new TFT",
 * "add packet filters to existing TFT" and "replace packet filters
 * in existing TFT" shall contain a variable number of packet filters
 * ------------------------------------------------------------------
 */
#define TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX 4
typedef struct {
  uint8_t identifier : 4;
#define TRAFFIC_FLOW_TEMPLATE_PRE_REL7_TFT_FILTER 0b00
#define TRAFFIC_FLOW_TEMPLATE_DOWNLINK_ONLY 0b01
#define TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY 0b10
#define TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL 0b11
  uint8_t direction : 2;
  uint8_t eval_precedence;
  PacketFilter packetfilter;
} PacketFilters[TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX];

typedef PacketFilters CreateNewTft;
typedef PacketFilters AddPacketFilter;
typedef PacketFilters ReplacePacketFilter;

/*
 * Packet filter list
 * ------------------
 */
typedef union {
  CreateNewTft createtft;
  DeleteExistingTft deletetft;
  AddPacketFilter addpacketfilter;
  ReplacePacketFilter replacepacketfilter;
  DeletePacketFilter deletepacketfilter;
  NoTftOperation notftoperation;
} PacketFilterList;

/*
 * ----------------------------------------------------------------------------
 *        Parameters list
 * ----------------------------------------------------------------------------
 */

typedef struct {
  /* TODO */
} ParameterList;

/*
 * ----------------------------------------------------------------------------
 *      Traffic Flow Template information element
 * ----------------------------------------------------------------------------
 */

typedef struct TrafficFlowTemplate_tag {
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_SPARE 0b000
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE 0b001
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE 0b010
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_ADD_PACKET 0b011
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_REPLACE_PACKET 0b100
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_PACKET 0b101
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_NO_OPERATION 0b110
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_RESERVED 0b111
  uint8_t tftoperationcode : 3;
#define TRAFFIC_FLOW_TEMPLATE_PARAMETER_LIST_IS_NOT_INCLUDED 0
#define TRAFFIC_FLOW_TEMPLATE_PARAMETER_LIST_IS_INCLUDED 1
  uint8_t ebit : 1;
  uint8_t numberofpacketfilters : 4;
  PacketFilterList packetfilterlist;
  ParameterList parameterlist;
} TrafficFlowTemplate;

int encode_traffic_flow_template(
    TrafficFlowTemplate* trafficflowtemplate, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_traffic_flow_template(
    TrafficFlowTemplate* trafficflowtemplate, uint8_t iei, uint8_t* buffer,
    uint32_t len);

void dump_traffic_flow_template_xml(
    TrafficFlowTemplate* trafficflowtemplate, uint8_t iei);

#endif /* TRAFFIC FLOW TEMPLATE_H_ */
