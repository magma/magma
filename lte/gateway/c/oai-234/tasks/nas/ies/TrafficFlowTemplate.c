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

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "TrafficFlowTemplate.h"

static int decode_traffic_flow_template_delete_packet(
    DeletePacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len);
static int decode_traffic_flow_template_create_tft(
    CreateNewTft* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len);
static int decode_traffic_flow_template_add_packet(
    AddPacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len);
static int decode_traffic_flow_template_replace_packet(
    ReplacePacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len);

static int decode_traffic_flow_template_packet_filter_identifiers(
    PacketFilterIdentifiers* packetfilteridentifiers, uint8_t nbpacketfilters,
    uint8_t* buffer, uint32_t len);
static int decode_traffic_flow_template_packet_filters(
    PacketFilters* packetfilters, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len);

static int encode_traffic_flow_template_delete_packet(
    DeletePacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len);
static int encode_traffic_flow_template_create_tft(
    CreateNewTft* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len);
static int encode_traffic_flow_template_add_packet(
    AddPacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len);
static int encode_traffic_flow_template_replace_packet(
    ReplacePacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len);

static int encode_traffic_flow_template_packet_filter_identifiers(
    PacketFilterIdentifiers* packetfilteridentifiers, uint8_t nbpacketfilters,
    uint8_t* buffer, uint32_t len);
static int encode_traffic_flow_template_packet_filters(
    PacketFilters* packetfilters, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len);

static void dump_traffic_flow_template_packet_filter_identifiers(
    PacketFilterIdentifiers* packetfilteridentifiers, uint8_t nbpacketfilters);
static void dump_traffic_flow_template_packet_filters(
    PacketFilters* packetfilters, uint8_t nbpacketfilters);

int decode_traffic_flow_template(
    TrafficFlowTemplate* trafficflowtemplate, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded        = 0;
  int decoded_result = 0;
  uint8_t ielen      = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  trafficflowtemplate->tftoperationcode      = (*(buffer + decoded) >> 5) & 0x7;
  trafficflowtemplate->ebit                  = (*(buffer + decoded) >> 4) & 0x1;
  trafficflowtemplate->numberofpacketfilters = *(buffer + decoded) & 0xf;
  decoded++;

  /*
   * Decoding packet filter list
   */
  if (trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_PACKET) {
    decoded_result = decode_traffic_flow_template_delete_packet(
        &trafficflowtemplate->packetfilterlist.deletepacketfilter,
        trafficflowtemplate->numberofpacketfilters, (buffer + decoded),
        len - decoded);
  } else if (
      trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE) {
    decoded_result = decode_traffic_flow_template_create_tft(
        &trafficflowtemplate->packetfilterlist.createtft,
        trafficflowtemplate->numberofpacketfilters, (buffer + decoded),
        len - decoded);
  } else if (
      trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_ADD_PACKET) {
    decoded_result = decode_traffic_flow_template_add_packet(
        &trafficflowtemplate->packetfilterlist.addpacketfilter,
        trafficflowtemplate->numberofpacketfilters, (buffer + decoded),
        len - decoded);
  } else if (
      trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_REPLACE_PACKET) {
    decoded_result = decode_traffic_flow_template_replace_packet(
        &trafficflowtemplate->packetfilterlist.replacepacketfilter,
        trafficflowtemplate->numberofpacketfilters, (buffer + decoded),
        len - decoded);
  }
#if NAS_DEBUG
  dump_traffic_flow_template_xml(trafficflowtemplate, iei);
#endif

  if (decoded_result < 0) {
    return decoded_result;
  }

  return (decoded + decoded_result);
}

int encode_traffic_flow_template(
    TrafficFlowTemplate* trafficflowtemplate, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, TRAFFIC_FLOW_TEMPLATE_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_traffic_flow_template_xml(trafficflowtemplate, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = 0x00 |
                        ((trafficflowtemplate->tftoperationcode & 0x7) << 5) |
                        ((trafficflowtemplate->ebit & 0x1) << 4) |
                        (trafficflowtemplate->numberofpacketfilters & 0xf);
  encoded++;

  /*
   * Encoding packet filter list
   */
  if (trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_PACKET) {
    encoded += encode_traffic_flow_template_delete_packet(
        &trafficflowtemplate->packetfilterlist.deletepacketfilter,
        trafficflowtemplate->numberofpacketfilters, (buffer + encoded),
        len - encoded);
  } else if (
      trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE) {
    encoded += encode_traffic_flow_template_create_tft(
        &trafficflowtemplate->packetfilterlist.createtft,
        trafficflowtemplate->numberofpacketfilters, (buffer + encoded),
        len - encoded);
  } else if (
      trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_ADD_PACKET) {
    encoded += encode_traffic_flow_template_add_packet(
        &trafficflowtemplate->packetfilterlist.addpacketfilter,
        trafficflowtemplate->numberofpacketfilters, (buffer + encoded),
        len - encoded);
  } else if (
      trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_REPLACE_PACKET) {
    encoded += encode_traffic_flow_template_replace_packet(
        &trafficflowtemplate->packetfilterlist.replacepacketfilter,
        trafficflowtemplate->numberofpacketfilters, (buffer + encoded),
        len - encoded);
  }

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void dump_traffic_flow_template_xml(
    TrafficFlowTemplate* trafficflowtemplate, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Traffic Flow Template>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(
      LOG_NAS, "    <TFT operation code>%u</TFT operation code>\n",
      trafficflowtemplate->tftoperationcode);
  OAILOG_DEBUG(LOG_NAS, "    <E bit>%u</E bit>\n", trafficflowtemplate->ebit);
  OAILOG_DEBUG(
      LOG_NAS, "    <Number of packet filters>%u</Number of packet filters>\n",
      trafficflowtemplate->numberofpacketfilters);

  if (trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_PACKET) {
    dump_traffic_flow_template_packet_filter_identifiers(
        &trafficflowtemplate->packetfilterlist.deletepacketfilter,
        trafficflowtemplate->numberofpacketfilters);
  } else if (
      trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE) {
    dump_traffic_flow_template_packet_filters(
        &trafficflowtemplate->packetfilterlist.createtft,
        trafficflowtemplate->numberofpacketfilters);
  } else if (
      trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_ADD_PACKET) {
    dump_traffic_flow_template_packet_filters(
        &trafficflowtemplate->packetfilterlist.addpacketfilter,
        trafficflowtemplate->numberofpacketfilters);
  } else if (
      trafficflowtemplate->tftoperationcode ==
      TRAFFIC_FLOW_TEMPLATE_OPCODE_REPLACE_PACKET) {
    dump_traffic_flow_template_packet_filters(
        &trafficflowtemplate->packetfilterlist.replacepacketfilter,
        trafficflowtemplate->numberofpacketfilters);
  }

  OAILOG_DEBUG(LOG_NAS, "</Traffic Flow Template>\n");
}

static int decode_traffic_flow_template_delete_packet(
    DeletePacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len) {
  return decode_traffic_flow_template_packet_filter_identifiers(
      packetfilter, nbpacketfilters, buffer, len);
}

static int decode_traffic_flow_template_create_tft(
    CreateNewTft* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len) {
  return decode_traffic_flow_template_packet_filters(
      packetfilter, nbpacketfilters, buffer, len);
}

static int decode_traffic_flow_template_add_packet(
    AddPacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len) {
  return decode_traffic_flow_template_packet_filters(
      packetfilter, nbpacketfilters, buffer, len);
}

static int decode_traffic_flow_template_replace_packet(
    ReplacePacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len) {
  return decode_traffic_flow_template_packet_filters(
      packetfilter, nbpacketfilters, buffer, len);
}

static int decode_traffic_flow_template_packet_filter_identifiers(
    PacketFilterIdentifiers* packetfilteridentifiers, uint8_t nbpacketfilters,
    uint8_t* buffer, uint32_t len) {
  int decoded = 0, i;

  for (i = 0; (i < nbpacketfilters) && (len - decoded > 0); i++) {
    /*
     * Packet filter identifier
     */
    IES_DECODE_U8(buffer, decoded, (*packetfilteridentifiers)[i].identifier);
  }

  return decoded;
}

static int decode_traffic_flow_template_packet_filters(
    PacketFilters* packetfilters, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0, i, j;

  for (i = 0; (i < nbpacketfilters); i++) {
    if (len - decoded <= 0) {
      /*
       * Mismatch between the number of packet filters subfield,
       * * * * and the number of packet filters in the packet filter list
       */
      return (TLV_VALUE_DOESNT_MATCH);
    }

    /*
     * Initialize the packet filter presence flag indicator
     */
    (*packetfilters)[i].packetfilter.flags = 0;
    /*
     * Packet filter direction
     */
    (*packetfilters)[i].direction = *(buffer + decoded) >> 4;
    /*
     * Packet filter identifier
     */
    (*packetfilters)[i].identifier = *(buffer + decoded) & 0x0f;
    decoded++;
    /*
     * Packet filter evaluation precedence
     */
    IES_DECODE_U8(buffer, decoded, (*packetfilters)[i].eval_precedence);
    /*
     * Length of the Packet filter contents field
     */
    uint8_t pkflen;

    IES_DECODE_U8(buffer, decoded, pkflen);
    /*
     * Packet filter contents
     */
    int pkfstart = decoded;

    while (decoded - pkfstart < pkflen) {
      /*
       * Packet filter component type identifier
       */
      uint8_t component_type;

      IES_DECODE_U8(buffer, decoded, component_type);

      switch (component_type) {
        case TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR:
          /*
           * IPv4 remote address type
           */
          (*packetfilters)[i].packetfilter.flags |=
              TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;

          for (j = 0; j < TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE; j++) {
            (*packetfilters)[i].packetfilter.ipv4remoteaddr[j].addr =
                *(buffer + decoded);
            (*packetfilters)[i].packetfilter.ipv4remoteaddr[j].mask =
                *(buffer + decoded + TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE);
            decoded++;
          }

          decoded += TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE;
          break;

        case TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR:
          /*
           * IPv6 remote address type
           */
          (*packetfilters)[i].packetfilter.flags |=
              TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG;

          for (j = 0; j < TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE; j++) {
            (*packetfilters)[i].packetfilter.ipv6remoteaddr[j].addr =
                *(buffer + decoded);
            (*packetfilters)[i].packetfilter.ipv6remoteaddr[j].mask =
                *(buffer + decoded + TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE);
            decoded++;
          }

          decoded += TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE;
          break;

        case TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER:
          /*
           * Protocol identifier/Next header type
           */
          (*packetfilters)[i].packetfilter.flags |=
              TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG;
          IES_DECODE_U8(
              buffer, decoded,
              (*packetfilters)[i].packetfilter.protocolidentifier_nextheader);
          break;

        case TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT:
          /*
           * Single local port type
           */
          (*packetfilters)[i].packetfilter.flags |=
              TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG;
          IES_DECODE_U16(
              buffer, decoded,
              (*packetfilters)[i].packetfilter.singlelocalport);
          break;

        case TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE:
          /*
           * Local port range type
           */
          (*packetfilters)[i].packetfilter.flags |=
              TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE_FLAG;
          IES_DECODE_U16(
              buffer, decoded,
              (*packetfilters)[i].packetfilter.localportrange.lowlimit);
          IES_DECODE_U16(
              buffer, decoded,
              (*packetfilters)[i].packetfilter.localportrange.highlimit);
          break;

        case TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT:
          /*
           * Single remote port type
           */
          (*packetfilters)[i].packetfilter.flags |=
              TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG;
          IES_DECODE_U16(
              buffer, decoded,
              (*packetfilters)[i].packetfilter.singleremoteport);
          break;

        case TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE:
          /*
           * Remote port range type
           */
          (*packetfilters)[i].packetfilter.flags |=
              TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE_FLAG;
          IES_DECODE_U16(
              buffer, decoded,
              (*packetfilters)[i].packetfilter.remoteportrange.lowlimit);
          IES_DECODE_U16(
              buffer, decoded,
              (*packetfilters)[i].packetfilter.remoteportrange.highlimit);
          break;

        case TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX:
          /*
           * Security parameter index type
           */
          (*packetfilters)[i].packetfilter.flags |=
              TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX_FLAG;
          IES_DECODE_U32(
              buffer, decoded,
              (*packetfilters)[i].packetfilter.securityparameterindex);
          break;

        case TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS:
          /*
           * Type of service/Traffic class type
           */
          (*packetfilters)[i].packetfilter.flags |=
              TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS_FLAG;
          IES_DECODE_U8(
              buffer, decoded,
              (*packetfilters)[i]
                  .packetfilter.typdeofservice_trafficclass.value);
          IES_DECODE_U8(
              buffer, decoded,
              (*packetfilters)[i]
                  .packetfilter.typdeofservice_trafficclass.mask);
          break;

        case TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL:
          /*
           * Flow label type
           */
          (*packetfilters)[i].packetfilter.flags |=
              TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG;
          IES_DECODE_U24(
              buffer, decoded, (*packetfilters)[i].packetfilter.flowlabel);
          break;

        default:
          /*
           * Packet filter component type identifier is not valid
           */
          return (TLV_UNEXPECTED_IEI);
          break;
      }
    }
  }

  if (len - decoded != 0) {
    /*
     * Mismatch between the number of packet filters subfield,
     * * * * and the number of packet filters in the packet filter list
     */
    return (TLV_VALUE_DOESNT_MATCH);
  }

  return decoded;
}

static int encode_traffic_flow_template_delete_packet(
    DeletePacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len) {
  return encode_traffic_flow_template_packet_filter_identifiers(
      packetfilter, nbpacketfilters, buffer, len);
}

static int encode_traffic_flow_template_create_tft(
    CreateNewTft* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len) {
  return encode_traffic_flow_template_packet_filters(
      packetfilter, nbpacketfilters, buffer, len);
}

static int encode_traffic_flow_template_add_packet(
    AddPacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len) {
  return encode_traffic_flow_template_packet_filters(
      packetfilter, nbpacketfilters, buffer, len);
}

static int encode_traffic_flow_template_replace_packet(
    ReplacePacketFilter* packetfilter, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len) {
  return encode_traffic_flow_template_packet_filters(
      packetfilter, nbpacketfilters, buffer, len);
}

static int encode_traffic_flow_template_packet_filter_identifiers(
    PacketFilterIdentifiers* packetfilteridentifiers, uint8_t nbpacketfilters,
    uint8_t* buffer, uint32_t len) {
  int encoded = 0, i;

  for (i = 0; (i < nbpacketfilters) && (len - encoded > 0); i++) {
    /*
     * Packet filter identifier
     */
    IES_ENCODE_U8(buffer, encoded, (*packetfilteridentifiers)[i].identifier);
  }

  return encoded;
}

static int encode_traffic_flow_template_packet_filters(
    PacketFilters* packetfilters, uint8_t nbpacketfilters, uint8_t* buffer,
    uint32_t len) {
  int encoded = 0, i, j;

  for (i = 0; (i < nbpacketfilters) && (len - encoded > 0); i++) {
    if (len - encoded <= 0) {
      /*
       * Mismatch between the number of packet filters subfield,
       * * * * and the number of packet filters in the packet filter list
       */
      return (TLV_VALUE_DOESNT_MATCH);
    }

    /*
     * Packet filter identifier and direction
     */
    IES_ENCODE_U8(
        buffer, encoded,
        (0x00 | ((*packetfilters)[i].direction << 4) |
         ((*packetfilters)[i].identifier)));
    /*
     * Packet filter evaluation precedence
     */
    IES_ENCODE_U8(buffer, encoded, (*packetfilters)[i].eval_precedence);
    /*
     * Save address of the Packet filter contents field length
     */
    uint8_t* pkflenPtr = buffer + encoded;

    encoded++;
    /*
     * Packet filter contents
     */
    int pkfstart  = encoded;
    uint16_t flag = TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;

    while (flag <= TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG) {
      switch ((*packetfilters)[i].packetfilter.flags & flag) {
        case TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG:
          /*
           * IPv4 remote address type
           */
          IES_ENCODE_U8(
              buffer, encoded, TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR);

          for (j = 0; j < TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE; j++) {
            *(buffer + encoded) =
                (*packetfilters)[i].packetfilter.ipv4remoteaddr[j].addr;
            *(buffer + encoded + TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE) =
                (*packetfilters)[i].packetfilter.ipv4remoteaddr[j].mask;
            encoded++;
          }

          encoded += TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE;
          break;

        case TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG:
          /*
           * IPv6 remote address type
           */
          IES_ENCODE_U8(
              buffer, encoded, TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR);

          for (j = 0; j < TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE; j++) {
            *(buffer + encoded) =
                (*packetfilters)[i].packetfilter.ipv6remoteaddr[j].addr;
            *(buffer + encoded + TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE) =
                (*packetfilters)[i].packetfilter.ipv6remoteaddr[j].mask;
            encoded++;
          }

          encoded += TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE;
          break;

        case TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG:
          /*
           * Protocol identifier/Next header type
           */
          IES_ENCODE_U8(
              buffer, encoded, TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER);
          IES_ENCODE_U8(
              buffer, encoded,
              (*packetfilters)[i].packetfilter.protocolidentifier_nextheader);
          break;

        case TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG:
          /*
           * Single local port type
           */
          IES_ENCODE_U8(
              buffer, encoded, TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT);
          IES_ENCODE_U16(
              buffer, encoded,
              (*packetfilters)[i].packetfilter.singlelocalport);
          break;

        case TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE_FLAG:
          /*
           * Local port range type
           */
          IES_ENCODE_U8(
              buffer, encoded, TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE);
          IES_ENCODE_U16(
              buffer, encoded,
              (*packetfilters)[i].packetfilter.localportrange.lowlimit);
          IES_ENCODE_U16(
              buffer, encoded,
              (*packetfilters)[i].packetfilter.localportrange.highlimit);
          break;

        case TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG:
          /*
           * Single remote port type
           */
          IES_ENCODE_U8(
              buffer, encoded, TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT);
          IES_ENCODE_U16(
              buffer, encoded,
              (*packetfilters)[i].packetfilter.singleremoteport);
          break;

        case TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE_FLAG:
          /*
           * Remote port range type
           */
          IES_ENCODE_U8(
              buffer, encoded, TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE);
          IES_ENCODE_U16(
              buffer, encoded,
              (*packetfilters)[i].packetfilter.remoteportrange.lowlimit);
          IES_ENCODE_U16(
              buffer, encoded,
              (*packetfilters)[i].packetfilter.remoteportrange.highlimit);
          break;

        case TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX_FLAG:
          /*
           * Security parameter index type
           */
          IES_ENCODE_U8(
              buffer, encoded, TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX);
          IES_ENCODE_U32(
              buffer, encoded,
              (*packetfilters)[i].packetfilter.securityparameterindex);
          break;

        case TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS_FLAG:
          /*
           * Type of service/Traffic class type
           */
          IES_ENCODE_U8(
              buffer, encoded,
              TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS);
          IES_ENCODE_U8(
              buffer, encoded,
              (*packetfilters)[i]
                  .packetfilter.typdeofservice_trafficclass.value);
          IES_ENCODE_U8(
              buffer, encoded,
              (*packetfilters)[i]
                  .packetfilter.typdeofservice_trafficclass.mask);
          break;

        case TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG:
          /*
           * Flow label type
           */
          IES_ENCODE_U8(buffer, encoded, TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL);
          IES_ENCODE_U24(
              buffer, encoded,
              (*packetfilters)[i].packetfilter.flowlabel & 0x000fffff);
          break;

        default:
          break;
      }

      flag = flag << 1;
    }

    /*
     * Length of the Packet filter contents field
     */
    *pkflenPtr = encoded - pkfstart;
  }

  return encoded;
}

static void dump_traffic_flow_template_packet_filter_identifiers(
    PacketFilterIdentifiers* packetfilteridentifiers, uint8_t nbpacketfilters) {
  int i;

  OAILOG_DEBUG(LOG_NAS, "    <Packet filter list>\n");

  for (i = 0; i < nbpacketfilters; i++) {
    OAILOG_DEBUG(
        LOG_NAS, "        <Identifier>%u</Identifier>\n",
        (*packetfilteridentifiers)[i].identifier);
  }

  OAILOG_DEBUG(LOG_NAS, "    </Packet filter list>\n");
}

static void dump_traffic_flow_template_packet_filters(
    PacketFilters* packetfilters, uint8_t nbpacketfilters) {
  int i;

  OAILOG_DEBUG(LOG_NAS, "    <Packet filter list>\n");

  for (i = 0; i < nbpacketfilters; i++) {
    OAILOG_DEBUG(
        LOG_NAS, "        <Identifier>%u</Identifier>\n",
        (*packetfilters)[i].identifier);
    OAILOG_DEBUG(
        LOG_NAS, "        <Direction>%u</Direction>\n",
        (*packetfilters)[i].direction);
    OAILOG_DEBUG(
        LOG_NAS, "        <Evaluation precedence>%u</Evaluation precedence>\n",
        (*packetfilters)[i].eval_precedence);
    OAILOG_DEBUG(LOG_NAS, "        <Packet filter>\n");

    if ((*packetfilters)[i].packetfilter.flags &
        TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG) {
      OAILOG_DEBUG(
          LOG_NAS,
          "            <IPv4 remote address>%u.%u.%u.%u</IPv4 remote "
          "address>\n",
          (*packetfilters)[i].packetfilter.ipv4remoteaddr[0].addr,
          (*packetfilters)[i].packetfilter.ipv4remoteaddr[1].addr,
          (*packetfilters)[i].packetfilter.ipv4remoteaddr[2].addr,
          (*packetfilters)[i].packetfilter.ipv4remoteaddr[3].addr);
      OAILOG_DEBUG(
          LOG_NAS,
          "            <IPv4 remote address mask>%u.%u.%u.%u</IPv4 remote "
          "address mask>\n",
          (*packetfilters)[i].packetfilter.ipv4remoteaddr[0].mask,
          (*packetfilters)[i].packetfilter.ipv4remoteaddr[1].mask,
          (*packetfilters)[i].packetfilter.ipv4remoteaddr[2].mask,
          (*packetfilters)[i].packetfilter.ipv4remoteaddr[3].mask);
    }

    if ((*packetfilters)[i].packetfilter.flags &
        TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG) {
      OAILOG_DEBUG(
          LOG_NAS,
          "            <Ipv6 remote "
          "address>%x%.2x:%x%.2x:%x%.2x:%x%.2x:%x%.2x:%x%.2x:%x%.2x:%x%.2x</"
          "Ipv6 "
          "remote address>\n",
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[0].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[1].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[2].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[3].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[4].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[5].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[6].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[7].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[8].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[9].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[10].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[11].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[12].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[13].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[14].addr,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[15].addr);
      OAILOG_DEBUG(
          LOG_NAS,
          "            <Ipv6 remote address "
          "mask>%x%.2x:%x%.2x:%x%.2x:%x%.2x:%x%.2x:%x%.2x:%x%.2x:%x%.2x</Ipv6 "
          "remote address mask>\n",
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[0].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[1].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[2].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[3].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[4].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[5].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[6].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[7].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[8].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[9].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[10].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[11].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[12].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[13].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[14].mask,
          (*packetfilters)[i].packetfilter.ipv6remoteaddr[15].mask);
    }

    if ((*packetfilters)[i].packetfilter.flags &
        TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG) {
      OAILOG_DEBUG(
          LOG_NAS,
          "            <Protocol identifier - Next header type>%u</Protocol "
          "identifier - Next header type>\n",
          (*packetfilters)[i].packetfilter.protocolidentifier_nextheader);
    }

    if ((*packetfilters)[i].packetfilter.flags &
        TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG) {
      OAILOG_DEBUG(
          LOG_NAS, "            <Single local port>%u</Single local port>\n",
          (*packetfilters)[i].packetfilter.singlelocalport);
    }

    if ((*packetfilters)[i].packetfilter.flags &
        TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE_FLAG) {
      OAILOG_DEBUG(
          LOG_NAS,
          "            <Local port range low limit>%u</Local port range low "
          "limit>\n",
          (*packetfilters)[i].packetfilter.localportrange.lowlimit);
      OAILOG_DEBUG(
          LOG_NAS,
          "            <Local port range high limit>%u</Local port range high "
          "limit>\n",
          (*packetfilters)[i].packetfilter.localportrange.highlimit);
    }

    if ((*packetfilters)[i].packetfilter.flags &
        TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG) {
      OAILOG_DEBUG(
          LOG_NAS, "            <Single remote port>%u</Single remote port>\n",
          (*packetfilters)[i].packetfilter.singleremoteport);
    }

    if ((*packetfilters)[i].packetfilter.flags &
        TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE_FLAG) {
      OAILOG_DEBUG(
          LOG_NAS,
          "            <Remote port range low limit>%u</Remote port range low "
          "limit>\n",
          (*packetfilters)[i].packetfilter.remoteportrange.lowlimit);
      OAILOG_DEBUG(
          LOG_NAS,
          "            <Remote port range high limit>%u</Remote port range "
          "high "
          "limit>\n",
          (*packetfilters)[i].packetfilter.remoteportrange.highlimit);
    }

    if ((*packetfilters)[i].packetfilter.flags &
        TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX_FLAG) {
      OAILOG_DEBUG(
          LOG_NAS,
          "            <Security parameter index>%u</Security parameter "
          "index>\n",
          (*packetfilters)[i].packetfilter.securityparameterindex);
    }

    if ((*packetfilters)[i].packetfilter.flags &
        TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS_FLAG) {
      OAILOG_DEBUG(
          LOG_NAS,
          "            <Type of service - Traffic class>%u</Type of service - "
          "Traffic class>\n",
          (*packetfilters)[i].packetfilter.typdeofservice_trafficclass.value);
      OAILOG_DEBUG(
          LOG_NAS,
          "            <Type of service - Traffic class mask>%u</Type of "
          "service "
          "- Traffic class mask>\n",
          (*packetfilters)[i].packetfilter.typdeofservice_trafficclass.mask);
    }

    if ((*packetfilters)[i].packetfilter.flags &
        TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG) {
      OAILOG_DEBUG(
          LOG_NAS, "            <Flow label>%u</Flow label>\n",
          (*packetfilters)[i].packetfilter.flowlabel);
    }

    OAILOG_DEBUG(LOG_NAS, "        </Packet filter>\n");
  }

  OAILOG_DEBUG(LOG_NAS, "    </Packet filter list>\n");
}
