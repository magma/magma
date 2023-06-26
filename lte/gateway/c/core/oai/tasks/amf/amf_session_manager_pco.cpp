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

#include "lte/gateway/c/core/oai/tasks/amf/include/amf_session_manager_pco.hpp"

#include <netinet/in.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
#include "lte/gateway/c/core/oai/common/rfc_1332.h"
#include "lte/gateway/c/core/oai/common/rfc_1877.h"

namespace magma5g {

void sm_clear_protocol_configuration_options(
    protocol_configuration_options_t* const pco) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (pco) {
    for (int i = 0; i < PCO_UNSPEC_MAXIMUM_PROTOCOL_ID_OR_CONTAINER_ID; i++) {
      if (pco->protocol_or_container_ids[i].contents) {
        bdestroy_wrapper(&pco->protocol_or_container_ids[i].contents);
      }
    }
    memset(pco, 0, sizeof(protocol_configuration_options_t));
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void sm_free_protocol_configuration_options(
    protocol_configuration_options_t** const protocol_configuration_options) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  protocol_configuration_options_t* pco = *protocol_configuration_options;
  if (pco) {
    for (int i = 0; i < pco->num_protocol_or_container_id; i++) {
      if (pco->protocol_or_container_ids[i].contents &&
          pco->protocol_or_container_ids[i].length) {
        bdestroy_wrapper(&pco->protocol_or_container_ids[i].contents);
      }
    }
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void sm_copy_protocol_configuration_options(
    protocol_configuration_options_t* const pco_dst,
    const protocol_configuration_options_t* const pco_src) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if ((pco_dst) && (pco_src)) {
    memset(pco_dst, 0, sizeof(protocol_configuration_options_t));

    pco_dst->ext = pco_src->ext;
    pco_dst->spare = pco_src->spare;
    pco_dst->configuration_protocol = pco_src->configuration_protocol;
    pco_dst->num_protocol_or_container_id =
        pco_src->num_protocol_or_container_id;

    for (int i = 0; i < pco_src->num_protocol_or_container_id; i++) {
      pco_dst->protocol_or_container_ids[i].id =
          pco_src->protocol_or_container_ids[i].id;
      pco_dst->protocol_or_container_ids[i].length =
          pco_src->protocol_or_container_ids[i].length;
      if (pco_src->protocol_or_container_ids[i].length) {
        pco_dst->protocol_or_container_ids[i].contents =
            bstrcpy(pco_src->protocol_or_container_ids[i].contents);
      }
    }
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

static void sm_pco_push_protocol_or_container_id(
    protocol_configuration_options_t* const pco,
    pco_protocol_or_container_id_t* const poc_id) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  pco->protocol_or_container_ids[pco->num_protocol_or_container_id].id =
      poc_id->id;
  pco->protocol_or_container_ids[pco->num_protocol_or_container_id].length =
      poc_id->length;
  pco->protocol_or_container_ids[pco->num_protocol_or_container_id].contents =
      poc_id->contents;
  poc_id->contents = NULL;
  pco->num_protocol_or_container_id += 1;
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

uint16_t sm_process_pco_dns_server_request(
    protocol_configuration_options_t* const pco_resp) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  in_addr_t ipcp_out_dns_prim_ipv4_addr = INADDR_NONE;
  pco_protocol_or_container_id_t poc_id_resp = {0};
  uint8_t dns_array[4];

  amf_config_read_lock(&amf_config);
  ipcp_out_dns_prim_ipv4_addr = amf_config.ipv4.default_dns.s_addr;
  amf_config_unlock(&amf_config);

  poc_id_resp.id = PCO_CI_DNS_SERVER_IPV4_ADDRESS;
  poc_id_resp.length = 4;
  dns_array[0] = (uint8_t)(ipcp_out_dns_prim_ipv4_addr & 0x000000FF);
  dns_array[1] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 8) & 0x000000FF);
  dns_array[2] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 16) & 0x000000FF);
  dns_array[3] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 24) & 0x000000FF);
  poc_id_resp.contents = blk2bstr(dns_array, sizeof(dns_array));

  sm_pco_push_protocol_or_container_id(pco_resp, &poc_id_resp);

  OAILOG_FUNC_RETURN(LOG_AMF_APP,
                     (poc_id_resp.length + SM_PCO_IPCP_HDR_LENGTH));
}

uint16_t sm_process_pco_p_cscf_address_request(
    protocol_configuration_options_t* const pco_resp) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  in_addr_t pcscf_prim_ipv4_addr = INADDR_NONE;
  pco_protocol_or_container_id_t poc_id_resp = {0};
  uint8_t pcscf_array[4];

  amf_config_read_lock(&amf_config);
  pcscf_prim_ipv4_addr = amf_config.pcscf_addr.ipv4.s_addr;
  amf_config_unlock(&amf_config);

  poc_id_resp.id = PCO_CI_P_CSCF_IPV4_ADDRESS_REQUEST;
  poc_id_resp.length = 4;
  pcscf_array[0] = (uint8_t)(pcscf_prim_ipv4_addr & 0x000000FF);
  pcscf_array[1] = (uint8_t)((pcscf_prim_ipv4_addr >> 8) & 0x000000FF);
  pcscf_array[2] = (uint8_t)((pcscf_prim_ipv4_addr >> 16) & 0x000000FF);
  pcscf_array[3] = (uint8_t)((pcscf_prim_ipv4_addr >> 24) & 0x000000FF);
  poc_id_resp.contents = blk2bstr(pcscf_array, sizeof(pcscf_array));

  sm_pco_push_protocol_or_container_id(pco_resp, &poc_id_resp);

  OAILOG_FUNC_RETURN(LOG_AMF_APP,
                     (poc_id_resp.length + SM_PCO_IPCP_HDR_LENGTH));
}

uint16_t sm_process_pco_request_ipcp(
    protocol_configuration_options_t* pco_resp,
    const pco_protocol_or_container_id_t* const poc_id) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  pco_protocol_or_container_id_t poc_id_resp = {0};
  in_addr_t ipcp_out_dns_prim_ipv4_addr = INADDR_NONE;
  in_addr_t ipcp_out_dns_sec_ipv4_addr = INADDR_NONE;
  int16_t ipcp_req_remaining_length = poc_id->length;
  size_t pco_in_index = 4;

  uint8_t idp[6] = {0};
  uint8_t ids[6] = {0};
  int16_t ipcp_out_length = 0;

  uint8_t ipcp_req_option = 0;
  int8_t ipcp_req_option_length = 0;

  bool dns_req = false;

  while (ipcp_req_remaining_length >= 2) {
    ipcp_req_option = poc_id->contents->data[pco_in_index];
    ipcp_req_option_length = poc_id->contents->data[pco_in_index + 1];
    ipcp_req_remaining_length =
        ipcp_req_remaining_length - ipcp_req_option_length;

    if ((ipcp_req_option == IPCP_OPTION_PRIMARY_DNS_SERVER_IP_ADDRESS) ||
        (ipcp_req_option == IPCP_OPTION_SECONDARY_DNS_SERVER_IP_ADDRESS)) {
      dns_req = true;
      break;
    }
  }

  if (dns_req == false) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, 0);
  }

  /* Fetch the DNS entries */
  amf_config_read_lock(&amf_config);

  ipcp_out_dns_prim_ipv4_addr = amf_config.ipv4.default_dns.s_addr;
  ipcp_out_dns_sec_ipv4_addr = amf_config.ipv4.default_dns_sec.s_addr;

  amf_config_unlock(&amf_config);

  /* Code + Identifier + Length */
  ipcp_out_length = IPCP_CODE_BYTES + IPCP_IDENTIFIER_BYTES + IPC_LENGTH_BYTES;
  poc_id_resp.id = PCO_PI_IPCP;
  poc_id_resp.length = 0;
  uint8_t cil[4] = {0};
  poc_id_resp.contents = blk2bstr(cil, 4);

  pco_resp->ext = 1;
  pco_resp->spare = 0;
  pco_resp->num_protocol_or_container_id = 0;
  pco_resp->configuration_protocol = 0;

  /* Primary DNS Server IP Address */
  idp[0] = IPCP_OPTION_PRIMARY_DNS_SERVER_IP_ADDRESS;
  idp[1] = 6;
  idp[2] = (uint8_t)(ipcp_out_dns_prim_ipv4_addr & 0x000000FF);
  idp[3] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 8) & 0x000000FF);
  idp[4] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 16) & 0x000000FF);
  idp[5] = (uint8_t)((ipcp_out_dns_prim_ipv4_addr >> 24) & 0x000000FF);
  ipcp_out_length += 6;
  bcatblk(poc_id_resp.contents, idp, 6);
  OAILOG_DEBUG(LOG_AMF_APP,
               "PCO: Protocol identifier IPCP option "
               "PRIMARY_DNS_SERVER_IP_ADDRESS ipcp_out_dns_prim_ipv4_addr "
               "0x%x\n",
               ipcp_out_dns_prim_ipv4_addr);

  /* Secondary DNS Server IP Address */
  ids[0] = IPCP_OPTION_SECONDARY_DNS_SERVER_IP_ADDRESS;
  ids[1] = 6;
  ids[2] = (uint8_t)(ipcp_out_dns_sec_ipv4_addr & 0x000000FF);
  ids[3] = (uint8_t)((ipcp_out_dns_sec_ipv4_addr >> 8) & 0x000000FF);
  ids[4] = (uint8_t)((ipcp_out_dns_sec_ipv4_addr >> 16) & 0x000000FF);
  ids[5] = (uint8_t)((ipcp_out_dns_sec_ipv4_addr >> 24) & 0x000000FF);
  ipcp_out_length += 6;
  bcatblk(poc_id_resp.contents, ids, 6);
  OAILOG_DEBUG(LOG_AMF_APP,
               "PCO: Protocol identifier IPCP option "
               "SECONDARY_DNS_SERVER_IP_ADDRESS ipcp_out_dns_sec_ipv4_addr "
               "0x%x\n",
               ipcp_out_dns_sec_ipv4_addr);
  // finally we can fill code, length
  poc_id_resp.length = ipcp_out_length;  // fill value after parsing req
  poc_id_resp.contents->data[0] = IPCP_CODE_CONFIGURE_NACK;
  poc_id_resp.contents->data[1] = 0;
  poc_id_resp.contents->data[2] = (uint8_t)(ipcp_out_length >> 8);
  poc_id_resp.contents->data[3] = (uint8_t)(ipcp_out_length & 0x00FF);

  sm_pco_push_protocol_or_container_id(pco_resp, &poc_id_resp);

  OAILOG_FUNC_RETURN(LOG_AMF_APP, ipcp_out_length + SM_PCO_IPCP_HDR_LENGTH);
}

uint16_t sm_process_pco_request(protocol_configuration_options_t* pco_req,
                                protocol_configuration_options_t* pco_resp) {
  auto length = 0;
  OAILOG_FUNC_IN(LOG_AMF_APP);
  for (auto id = 0; id < pco_req->num_protocol_or_container_id; id++) {
    switch (pco_req->protocol_or_container_ids[id].id) {
      case PCO_PI_IPCP:
        length += sm_process_pco_request_ipcp(
            pco_resp, &pco_req->protocol_or_container_ids[id]);
        break;

      case PCO_CI_P_CSCF_IPV4_ADDRESS_REQUEST:
        length += sm_process_pco_p_cscf_address_request(pco_resp);
        break;

      case PCO_CI_DNS_SERVER_IPV4_ADDRESS_REQUEST:
        length += sm_process_pco_dns_server_request(pco_resp);
        break;

      default:
        break;
    }
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, length);
}

}  // namespace magma5g
