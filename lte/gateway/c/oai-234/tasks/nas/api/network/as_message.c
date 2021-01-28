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

/*****************************************************************************

  Source    as_message.c

  Version   0.1

  Date    2012/11/06

  Product   NAS stack

  Subsystem Application Programming Interface

  Author    Frederic Maurel, Lionel GAUTHIER

  Description Defines the messages supported by the Access Stratum sublayer
    protocol (usually RRC and S1AP for E-UTRAN) and functions used
    to encode and decode

*****************************************************************************/
#include <string.h>  // memcpy
#include <stdlib.h>  // free
#include <stdint.h>
#include <stdbool.h>

#include "bstrlib.h"

#include "log.h"
#include "nas/as_message.h"
#include "common_types.h"
#include "dynamic_memory_check.h"
#include "common_defs.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:  as_message_decode()                                       **
 **                                                                        **
 ** Description: Decode AS message and accordingly fills data structure    **
 **                                                                        **
 ** Inputs:  buffer:  Pointer to the buffer containing the       **
 **       message                                    **
 **      length:  Number of bytes that should be decoded     **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   msg:   AS message structure to be filled          **
 **      Return:  The AS message identifier when the buffer  **
 **       has been successfully decoded;             **
 **       RETURNerror otherwise.                     **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
int as_message_decode(const char* buffer, as_message_t* msg, size_t length) {
  OAILOG_FUNC_IN(LOG_NAS);
  int bytes;
  uint8_t** data = NULL;

  /*
   * Get the message type
   */
  msg->msg_id = *(uint16_t*) (buffer);
  bytes       = sizeof(uint16_t);

  switch (msg->msg_id) {
    case AS_NAS_ESTABLISH_REQ:
      /*
       * NAS signalling connection establish request
       */
      bytes += sizeof(nas_establish_req_t) - sizeof(uint8_t*);
      data = &msg->msg.nas_establish_req.initial_nas_msg->data;
      break;

    case AS_NAS_ESTABLISH_IND:
      /*
       * NAS signalling connection establishment indication
       */
      bytes += sizeof(nas_establish_ind_t) - sizeof(uint8_t*);
      data = &msg->msg.nas_establish_ind.initial_nas_msg->data;
      break;

    case AS_NAS_ESTABLISH_RSP:
      /*
       * NAS signalling connection establishment response
       */
      bytes += sizeof(nas_establish_rsp_t) - sizeof(uint8_t*);
      data = &msg->msg.nas_establish_rsp.nas_msg->data;
      break;

    case AS_NAS_ESTABLISH_CNF:
      /*
       * NAS signalling connection establishment confirm
       */
      bytes += sizeof(nas_establish_cnf_t) - sizeof(uint8_t*);
      data = &msg->msg.nas_establish_cnf.nas_msg->data;
      break;

    case AS_UL_INFO_TRANSFER_REQ:
      /*
       * Uplink L3 data transfer request
       */
      bytes += sizeof(ul_info_transfer_req_t) - sizeof(uint8_t*);
      data = &msg->msg.ul_info_transfer_req.nas_msg->data;
      break;

    case AS_UL_INFO_TRANSFER_IND:
      /*
       * Uplink L3 data transfer indication
       */
      bytes += sizeof(ul_info_transfer_ind_t) - sizeof(uint8_t*);
      data = &msg->msg.ul_info_transfer_ind.nas_msg->data;
      break;

    case AS_DL_INFO_TRANSFER_REQ:
      /*
       * Downlink L3 data transfer request
       */
      bytes += sizeof(dl_info_transfer_req_t) - sizeof(uint8_t*);
      data = &msg->msg.dl_info_transfer_req.nas_msg->data;
      break;

    case AS_DL_INFO_TRANSFER_IND:
      /*
       * Downlink L3 data transfer indication
       */
      bytes += sizeof(dl_info_transfer_ind_t) - sizeof(uint8_t*);
      data = &msg->msg.dl_info_transfer_ind.nas_msg->data;
      break;

    case AS_ACTIVATE_BEARER_CONTEXT_REQ:
      bytes += sizeof(dl_info_transfer_req_t) - sizeof(uint8_t*);
      data = &msg->msg.activate_bearer_context_req.nas_msg->data;
      break;

    case AS_BROADCAST_INFO_IND:
    case AS_PAGING_REQ:
    case AS_NAS_RELEASE_REQ:
    case AS_UL_INFO_TRANSFER_CNF:
    case AS_DL_INFO_TRANSFER_CNF:
    case AS_NAS_RELEASE_IND:
    case AS_ERAB_SETUP_IND:
    case AS_ERAB_SETUP_RSP:
    case AS_ERAB_SETUP_CNF:
    case AS_RAB_RELEASE_REQ:
    case AS_RAB_RELEASE_IND:
      /*
       * Messages without dedicated NAS information
       */
      bytes = length;
      break;

    default:
      bytes = 0;
      OAILOG_WARNING(
          LOG_NAS, "NET-API   - AS message 0x%x is not valid", msg->msg_id);
      break;
  }

  if (bytes > 0) {
    if (data) {
      /*
       * Set the pointer to dedicated NAS information
       */
      *data = (uint8_t*) (buffer + bytes);
    }

    /*
     * Decode the message
     */
    memcpy(msg, (as_message_t*) buffer, bytes);
    OAILOG_FUNC_RETURN(LOG_NAS, msg->msg_id);
  }

  OAILOG_WARNING(
      LOG_NAS, "NET-API   - Failed to decode AS message 0x%x", msg->msg_id);
  OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
}

/****************************************************************************
 **                                                                        **
 ** Name:  as_message_encode()                                       **
 **                                                                        **
 ** Description: Encode AS message                                         **
 **                                                                        **
 ** Inputs:  msg:   AS message structure to encode             **
 **      length:  Maximal capacity of the output buffer      **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   buffer:  Pointer to the encoded data buffer         **
 **      Return:  The number of characters in the buffer     **
 **       when data have been successfully encoded;  **
 **       RETURNerror otherwise.                     **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
// int
// as_message_encode (
//  char *buffer,
//  as_message_t * msg,
//  size_t length)
//{
//  OAILOG_FUNC_IN (LOG_NAS);
//  int                                     bytes = sizeof (msg->msg_id);
//  bstring                                 nas_msg = NULL;
//
//  switch (msg->msg_id) {
//  case AS_BROADCAST_INFO_IND:
//    /*
//     * Broadcast information
//     */
//    bytes += sizeof (broadcast_info_ind_t);
//    break;
//
//  case AS_PAGING_REQ:
//    /*
//     * Paging information request
//     */
//    bytes += sizeof (paging_req_t);
//    break;
//
//  case AS_NAS_ESTABLISH_REQ:
//    /*
//     * NAS signalling connection establish request
//     */
//    bytes += sizeof (nas_establish_req_t) - sizeof (uint8_t *);
//    nas_msg = msg->msg.nas_establish_req.initial_nas_msg;
//    msg->msg.nas_establish_req.initial_nas_msg = NULL;
//    break;
//
//  case AS_NAS_ESTABLISH_IND:
//    /*
//     * NAS signalling connection establish indication
//     */
//    bytes += sizeof (nas_establish_ind_t) - sizeof (uint8_t *);
//    nas_msg = msg->msg.nas_establish_ind.initial_nas_msg;
//    msg->msg.nas_establish_ind.initial_nas_msg = NULL;
//    break;
//
//  case AS_NAS_ESTABLISH_RSP:
//    /*
//     * NAS signalling connection establish response
//     */
//    bytes += sizeof (nas_establish_rsp_t) - sizeof (uint8_t *);
//    nas_msg = msg->msg.nas_establish_rsp.nas_msg;
//    msg->msg.nas_establish_rsp.nas_msg = NULL;
//    break;
//
//  case AS_NAS_ESTABLISH_CNF:
//    /*
//     * NAS signalling connection establish confirm
//     */
//    bytes += sizeof (nas_establish_cnf_t) - sizeof (uint8_t *);
//    nas_msg = msg->msg.nas_establish_cnf.nas_msg;
//    msg->msg.nas_establish_cnf.nas_msg = NULL;
//    break;
//
//  case AS_NAS_RELEASE_REQ:
//    /*
//     * NAS signalling connection release request
//     */
//    bytes += sizeof (nas_release_req_t);
//    break;
//
//  case AS_NAS_RELEASE_IND:
//    /*
//     * NAS signalling connection release indication
//     */
//    bytes += sizeof (nas_release_ind_t);
//    break;
//
//  case AS_UL_INFO_TRANSFER_REQ:
//    /*
//     * Uplink L3 data transfer request
//     */
//    bytes += sizeof (ul_info_transfer_req_t) - sizeof (uint8_t *);
//    nas_msg = msg->msg.ul_info_transfer_req.nas_msg;
//    msg->msg.ul_info_transfer_req.nas_msg = NULL;
//    break;
//
//  case AS_UL_INFO_TRANSFER_CNF:
//    /*
//     * Uplink L3 data transfer confirm
//     */
//    bytes += sizeof (ul_info_transfer_cnf_t);
//    break;
//
//  case AS_UL_INFO_TRANSFER_IND:
//    /*
//     * Uplink L3 data transfer indication
//     */
//    bytes += sizeof (ul_info_transfer_ind_t) - sizeof (uint8_t *);
//    nas_msg = msg->msg.ul_info_transfer_ind.nas_msg;
//    msg->msg.ul_info_transfer_ind.nas_msg = NULL;
//    break;
//
//  case AS_DL_INFO_TRANSFER_REQ:
//    /*
//     * Downlink L3 data transfer
//     */
//    bytes += sizeof (dl_info_transfer_req_t) - sizeof (uint8_t *);
//    nas_msg = msg->msg.dl_info_transfer_req.nas_msg;
//    msg->msg.dl_info_transfer_req.nas_msg = NULL;
//    break;
//
//  case AS_DL_INFO_TRANSFER_CNF:
//    /*
//     * Downlink L3 data transfer confirm
//     */
//    bytes += sizeof (dl_info_transfer_cnf_t);
//    break;
//
//  case AS_DL_INFO_TRANSFER_IND:
//    /*
//     * Downlink L3 data transfer indication
//     */
//    bytes += sizeof (dl_info_transfer_ind_t) - sizeof (uint8_t *);
//    nas_msg = msg->msg.dl_info_transfer_ind.nas_msg;
//    msg->msg.dl_info_transfer_ind.nas_msg = NULL;
//    break;
//
//  case AS_RAB_ESTABLISH_REQ:
//    /*
//     * Radio Access Bearer establishment request
//     */
//    bytes += sizeof (rab_establish_req_t);
//    break;
//
//  case AS_RAB_ESTABLISH_IND:
//    /*
//     * Radio Access Bearer establishment indication
//     */
//    bytes += sizeof (rab_establish_ind_t);
//    break;
//
//  case AS_RAB_ESTABLISH_RSP:
//    /*
//     * Radio Access Bearer establishment response
//     */
//    bytes += sizeof (rab_establish_rsp_t);
//    break;
//
//  case AS_RAB_ESTABLISH_CNF:
//    /*
//     * Radio Access Bearer establishment confirm
//     */
//    bytes += sizeof (rab_establish_cnf_t);
//    break;
//
//  case AS_RAB_RELEASE_REQ:
//    /*
//     * Radio Access Bearer release request
//     */
//    bytes += sizeof (rab_release_req_t);
//    break;
//
//  case AS_RAB_RELEASE_IND:
//    /*
//     * Radio Access Bearer release indication
//     */
//    bytes += sizeof (rab_release_ind_t);
//    break;
//
//  default:
//    OAILOG_WARNING(LOG_NAS, "NET-API   - AS message 0x%x is not valid",
//    msg->msg_id); bytes = length; break;
//  }
//
//  if (length > bytes) {
//    /*
//     * Encode the AS message
//     */
//    memcpy (buffer, (unsigned char *)msg, bytes);
//
//    if (nas_msg && (nas_msg->length > 0)) {
//      /*
//       * Copy the NAS message
//       */
//      memcpy (buffer + bytes, nas_msg->data, nas_msg->length);
//      bytes += nas_msg->length;
//      /*
//       * Release NAS message memory
//       */
//      free_wrapper (nas_msg->data);
//      nas_msg->length = 0;
//      nas_msg->data = NULL;
//    }
//
//    OAILOG_FUNC_RETURN (LOG_NAS, bytes);
//  }
//
//  OAILOG_WARNING(LOG_NAS, "NET-API   - Failed to encode AS message 0x%x",
//  msg->msg_id); OAILOG_FUNC_RETURN (LOG_NAS, RETURNerror);
//}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
