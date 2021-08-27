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

#include <iostream>
#include "util_nas5g_pkt.h"
#include "rfc_1877.h"
#include "rfc_1332.h"

namespace magma5g {

//  API for testing decode registration request
bool decode_registration_request_msg(
    RegistrationRequestMsg* reg_request, const uint8_t* buffer, uint32_t len) {
  bool decode_success = true;
  uint8_t* decode_reg_buffer =
      const_cast<uint8_t*>(reinterpret_cast<const uint8_t*>(buffer));

  if (reg_request->DecodeRegistrationRequestMsg(
          reg_request, decode_reg_buffer, len) < 0) {
    decode_success = false;
  }

  return (decode_success);
}

//  API for testing encode registration reject
bool encode_registration_reject_msg(
    RegistrationRejectMsg* reg_reject, const uint8_t* buffer, uint32_t len) {
  bool encode_success = true;
  uint8_t* encode_reg_buffer =
      const_cast<uint8_t*>(reinterpret_cast<const uint8_t*>(buffer));

  if (reg_reject->EncodeRegistrationRejectMsg(
          reg_reject, encode_reg_buffer, len) < 0) {
    encode_success = false;
  }

  return (encode_success);
}

//  API for testing decode registration reject
bool decode_registration_reject_msg(
    RegistrationRejectMsg* reg_reject, const uint8_t* buffer, uint32_t len) {
  bool decode_success = true;
  uint8_t* decode_reg_buffer =
      const_cast<uint8_t*>(reinterpret_cast<const uint8_t*>(buffer));

  if (reg_reject->DecodeRegistrationRejectMsg(
          reg_reject, decode_reg_buffer, len) < 0) {
    decode_success = false;
  }

  return (decode_success);
}

void gen_pco_push_protocol_or_container_id(
    protocol_configuration_options_t* const pco,
    pco_protocol_or_container_id_t* const poc_id) {
  pco->protocol_or_container_ids[pco->num_protocol_or_container_id].id =
      poc_id->id;
  pco->protocol_or_container_ids[pco->num_protocol_or_container_id].length =
      poc_id->length;
  pco->protocol_or_container_ids[pco->num_protocol_or_container_id].contents =
      poc_id->contents;
  poc_id->contents = NULL;
  pco->num_protocol_or_container_id += 1;
}

int gen_dns_pco_options(protocol_configuration_options_t* const pco_resp) {
  struct sockaddr_in primary_dns_sa;
  pco_protocol_or_container_id_t poc_id_resp = {0};
  uint8_t dns_array[4];
  in_addr_t ipcp_out_dns_prim_ipv4_addr = INADDR_NONE;

  inet_pton(AF_INET, "192.168.1.100", &(primary_dns_sa.sin_addr));

  ipcp_out_dns_prim_ipv4_addr = primary_dns_sa.sin_addr.s_addr;
  poc_id_resp.id              = PCO_CI_DNS_SERVER_IPV4_ADDRESS;
  poc_id_resp.length          = 4;
  dns_array[0] = (uint8_t)(ipcp_out_dns_prim_ipv4_addr & 0x000000FF);
  dns_array[1] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 8) & 0x000000FF);
  dns_array[2] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 16) & 0x000000FF);
  dns_array[3] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 24) & 0x000000FF);
  poc_id_resp.contents = blk2bstr(dns_array, sizeof(dns_array));

  gen_pco_push_protocol_or_container_id(pco_resp, &poc_id_resp);

  return (poc_id_resp.length + 7);
}

void gen_ipcp_pco_options(protocol_configuration_options_t* const pco_resp) {
  pco_protocol_or_container_id_t poc_id_resp = {0};
  in_addr_t ipcp_out_dns_prim_ipv4_addr      = INADDR_NONE;
  in_addr_t ipcp_out_dns_sec_ipv4_addr       = INADDR_NONE;
  struct sockaddr_in primary_dns_sa;
  struct sockaddr_in secondry_dns_sa;
  uint8_t idp[6]          = {0};
  uint8_t ids[6]          = {0};
  int16_t ipcp_out_length = 0;

  inet_pton(AF_INET, "192.168.1.100", &(primary_dns_sa.sin_addr));
  inet_pton(AF_INET, "8.8.8.8", &(secondry_dns_sa.sin_addr));

  /* Code + Identifier + Length */
  ipcp_out_length = IPCP_CODE_BYTES + IPCP_IDENTIFIER_BYTES + IPC_LENGTH_BYTES;
  poc_id_resp.id  = PCO_PI_IPCP;
  poc_id_resp.length   = 0;
  uint8_t cil[4]       = {0};
  poc_id_resp.contents = blk2bstr(cil, 4);

  pco_resp->ext                          = 1;
  pco_resp->spare                        = 0;
  pco_resp->num_protocol_or_container_id = 0;
  pco_resp->configuration_protocol       = 0;

  ipcp_out_dns_prim_ipv4_addr = primary_dns_sa.sin_addr.s_addr;
  idp[0]                      = IPCP_OPTION_PRIMARY_DNS_SERVER_IP_ADDRESS;
  idp[1]                      = 6;
  idp[2] = (uint8_t)(ipcp_out_dns_prim_ipv4_addr & 0x000000FF);
  idp[3] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 8) & 0x000000FF);
  idp[4] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 16) & 0x000000FF);
  idp[5] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 24) & 0x000000FF);
  ipcp_out_length += 6;
  bcatblk(poc_id_resp.contents, idp, 6);

  ipcp_out_dns_sec_ipv4_addr = secondry_dns_sa.sin_addr.s_addr;
  ids[0]                     = IPCP_OPTION_SECONDARY_DNS_SERVER_IP_ADDRESS;
  ids[1]                     = 6;
  ids[2] = (uint8_t)(ipcp_out_dns_sec_ipv4_addr & 0x000000FF);
  ids[3] = (uint8_t)((ipcp_out_dns_sec_ipv4_addr >> 8) & 0x000000FF);
  ids[4] = (uint8_t)((ipcp_out_dns_sec_ipv4_addr >> 16) & 0x000000FF);
  ids[5] = (uint8_t)((ipcp_out_dns_sec_ipv4_addr >> 24) & 0x000000FF);
  ipcp_out_length += 6;
  bcatblk(poc_id_resp.contents, ids, 6);

  // finally we can fill code, length
  poc_id_resp.length = ipcp_out_length;  // fill value after parsing req
  poc_id_resp.contents->data[0] = IPCP_CODE_CONFIGURE_NACK;
  poc_id_resp.contents->data[1] = 0;
  poc_id_resp.contents->data[2] = (uint8_t)(ipcp_out_length >> 8);
  poc_id_resp.contents->data[3] = (uint8_t)(ipcp_out_length & 0x00FF);

  gen_pco_push_protocol_or_container_id(pco_resp, &poc_id_resp);
}

}  // namespace magma5g

