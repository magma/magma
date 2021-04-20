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

  Source      emm_msg.c

  Version     0.1

  Date        2012/09/27

  Product     NAS stack

  Subsystem   EPS Mobility Management

  Author      Frederic Maurel, Sebastien Roux

  Description Defines EPS Mobility Management messages

*****************************************************************************/
#include <stdint.h>

#include "log.h"
#include "3gpp_24.007.h"
#include "3gpp_24.301.h"
#include "emm_msg.h"
#include "common_defs.h"
#include "common_ies.h"
#include "conversions.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

int emm_msg_encode_header(
    const emm_msg_header_t* header, uint8_t* buffer, uint32_t len);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    emm_msg_decode()                                          **
 **                                                                        **
 ** Description: Decode EPS Mobility Management messages                   **
 **                                                                        **
 ** Inputs:  buffer:    Pointer to the buffer containing the EMM   **
 **             message data                               **
 **          len:       Number of bytes that should be decoded     **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     msg:       The EMM message structure to be filled     **
 **          Return:    The number of bytes in the buffer if data  **
 **             have been successfully decoded;            **
 **             A negative error code otherwise.           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_msg_decode(EMM_msg* msg, uint8_t* buffer, uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int header_result = 0;
  int decode_result = 0;

  /*
   * First decode the EMM message header
   */
  header_result = emm_msg_decode_header(&msg->header, buffer, len);

  if (header_result < 0) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-MSG   - Failed to decode EMM message header "
        "(%d)\n",
        header_result);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, header_result);
  }

  buffer += header_result;
  len -= header_result;
  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMM-MSG   - Message Type 0x%02x\n",
      msg->header.message_type);

  switch (msg->header.message_type) {
    case EMM_INFORMATION:
      decode_result =
          decode_emm_information(&msg->emm_information, buffer, len);
      break;

    case UPLINK_NAS_TRANSPORT:
      decode_result =
          decode_uplink_nas_transport(&msg->uplink_nas_transport, buffer, len);
      break;

    case AUTHENTICATION_REJECT:
      decode_result = decode_authentication_reject(
          &msg->authentication_reject, buffer, len);
      break;

    case AUTHENTICATION_FAILURE:
      decode_result = decode_authentication_failure(
          &msg->authentication_failure, buffer, len);
      break;

    case DETACH_ACCEPT:
      decode_result = decode_detach_accept(&msg->detach_accept, buffer, len);
      break;

    case SERVICE_REJECT:
      decode_result = decode_service_reject(&msg->service_reject, buffer, len);
      break;

    case AUTHENTICATION_REQUEST:
      decode_result = decode_authentication_request(
          &msg->authentication_request, buffer, len);
      break;

    case TRACKING_AREA_UPDATE_REQUEST:
      decode_result = decode_tracking_area_update_request(
          &msg->tracking_area_update_request, buffer, len);
      break;

    case ATTACH_REQUEST:
      decode_result = decode_attach_request(&msg->attach_request, buffer, len);
      break;

    case EMM_STATUS:
      decode_result = decode_emm_status(&msg->emm_status, buffer, len);
      break;

    case IDENTITY_RESPONSE:
      decode_result =
          decode_identity_response(&msg->identity_response, buffer, len);
      break;

    case IDENTITY_REQUEST:
      decode_result =
          decode_identity_request(&msg->identity_request, buffer, len);
      break;

    case GUTI_REALLOCATION_COMMAND:
      decode_result = decode_guti_reallocation_command(
          &msg->guti_reallocation_command, buffer, len);
      break;

    case TRACKING_AREA_UPDATE_REJECT:
      decode_result = decode_tracking_area_update_reject(
          &msg->tracking_area_update_reject, buffer, len);
      break;

    case ATTACH_ACCEPT:
      decode_result = decode_attach_accept(&msg->attach_accept, buffer, len);
      break;

    case SECURITY_MODE_COMPLETE:
      decode_result = decode_security_mode_complete(
          &msg->security_mode_complete, buffer, len);
      // IMEISV is 16 digits. Extra char for null termination.
      char imeisv[MAX_IMEISV_SIZE + 1];
      IMEISV_MOBID_TO_STRING(
          &msg->security_mode_complete.imeisv.imeisv, imeisv,
          MAX_IMEISV_SIZE + 1);
      OAILOG_INFO(LOG_NAS_EMM, "EMM-MSG   - IMEISV: %s", imeisv);
      break;

    case TRACKING_AREA_UPDATE_ACCEPT:
      decode_result = decode_tracking_area_update_accept(
          &msg->tracking_area_update_accept, buffer, len);
      break;

    case ATTACH_REJECT:
      decode_result = decode_attach_reject(&msg->attach_reject, buffer, len);
      break;

    case ATTACH_COMPLETE:
      decode_result =
          decode_attach_complete(&msg->attach_complete, buffer, len);
      break;

    case TRACKING_AREA_UPDATE_COMPLETE:
      decode_result = decode_tracking_area_update_complete(
          &msg->tracking_area_update_complete, buffer, len);
      break;

    case CS_SERVICE_NOTIFICATION:
      decode_result = decode_cs_service_notification(
          &msg->cs_service_notification, buffer, len);
      break;

    case SECURITY_MODE_REJECT:
      decode_result =
          decode_security_mode_reject(&msg->security_mode_reject, buffer, len);
      break;

    case DETACH_REQUEST:
      decode_result = decode_detach_request(&msg->detach_request, buffer, len);
      break;

    case GUTI_REALLOCATION_COMPLETE:
      decode_result = decode_guti_reallocation_complete(
          &msg->guti_reallocation_complete, buffer, len);
      break;

    case SECURITY_MODE_COMMAND:
      decode_result = decode_security_mode_command(
          &msg->security_mode_command, buffer, len);
      break;

    case DOWNLINK_NAS_TRANSPORT:
      decode_result = decode_downlink_nas_transport(
          &msg->downlink_nas_transport, buffer, len);
      break;

    case EXTENDED_SERVICE_REQUEST:
      decode_result = decode_extended_service_request(
          &msg->extended_service_request, buffer, len);
      break;

    case AUTHENTICATION_RESPONSE:
      decode_result = decode_authentication_response(
          &msg->authentication_response, buffer, len);
      break;

    default:
      OAILOG_ERROR(
          LOG_NAS_EMM, "EMM-MSG   - Unexpected message type: 0x%x\n",
          msg->header.message_type);
      decode_result = TLV_WRONG_MESSAGE_TYPE;
      /*
       * TODO: Handle not standard layer 3 messages: SERVICE_REQUEST
       */
  }

  if (decode_result < 0) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-MSG   - Failed to decode L3 EMM message 0x%x "
        "(%d)\n",
        msg->header.message_type, decode_result);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, decode_result);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, header_result + decode_result);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_msg_encode()                                          **
 **                                                                        **
 ** Description: Encode EPS Mobility Management messages                   **
 **                                                                        **
 ** Inputs:  msg:       The EMM message structure to encode        **
 **          length:    Maximal capacity of the output buffer      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     buffer:    Pointer to the encoded data buffer         **
 **          Return:    The number of bytes in the buffer if data  **
 **             have been successfully encoded;            **
 **             A negative error code otherwise.           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_msg_encode(EMM_msg* msg, uint8_t* buffer, uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int header_result;
  int encode_result;

  /*
   * First encode the EMM message header
   */
  header_result = emm_msg_encode_header(&msg->header, buffer, len);

  if (header_result < 0) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-MSG   - Failed to encode EMM message header "
        "(%d)\n",
        header_result);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, header_result);
  }

  buffer += header_result;
  len -= header_result;

  switch (msg->header.message_type) {
    case ATTACH_ACCEPT:
      encode_result = encode_attach_accept(&msg->attach_accept, buffer, len);
      break;

    case ATTACH_COMPLETE:
      encode_result =
          encode_attach_complete(&msg->attach_complete, buffer, len);
      break;

    case ATTACH_REJECT:
      encode_result = encode_attach_reject(&msg->attach_reject, buffer, len);
      break;

    case AUTHENTICATION_FAILURE:
      encode_result = encode_authentication_failure(
          &msg->authentication_failure, buffer, len);
      break;

    case AUTHENTICATION_REJECT:
      encode_result = encode_authentication_reject(
          &msg->authentication_reject, buffer, len);
      break;

    case AUTHENTICATION_REQUEST:
      encode_result = encode_authentication_request(
          &msg->authentication_request, buffer, len);
      break;

    case AUTHENTICATION_RESPONSE:
      encode_result = encode_authentication_response(
          &msg->authentication_response, buffer, len);
      break;

    case CS_SERVICE_NOTIFICATION:
      encode_result = encode_cs_service_notification(
          &msg->cs_service_notification, buffer, len);
      break;

    case DETACH_ACCEPT:
      encode_result = encode_detach_accept(&msg->detach_accept, buffer, len);
      break;

    case DETACH_REQUEST:
      encode_result = encode_detach_request(&msg->detach_request, buffer, len);
      break;

    case DOWNLINK_NAS_TRANSPORT:
      encode_result = encode_downlink_nas_transport(
          &msg->downlink_nas_transport, buffer, len);
      break;

    case EMM_INFORMATION:
      encode_result =
          encode_emm_information(&msg->emm_information, buffer, len);
      break;

    case EMM_STATUS:
      encode_result = encode_emm_status(&msg->emm_status, buffer, len);
      break;

    case EXTENDED_SERVICE_REQUEST:
      encode_result = encode_extended_service_request(
          &msg->extended_service_request, buffer, len);
      break;

    case GUTI_REALLOCATION_COMMAND:
      encode_result = encode_guti_reallocation_command(
          &msg->guti_reallocation_command, buffer, len);
      break;

    case GUTI_REALLOCATION_COMPLETE:
      encode_result = encode_guti_reallocation_complete(
          &msg->guti_reallocation_complete, buffer, len);
      break;

    case IDENTITY_REQUEST:
      encode_result =
          encode_identity_request(&msg->identity_request, buffer, len);
      break;

    case IDENTITY_RESPONSE:
      encode_result =
          encode_identity_response(&msg->identity_response, buffer, len);
      break;

    case SECURITY_MODE_COMMAND:
      encode_result = encode_security_mode_command(
          &msg->security_mode_command, buffer, len);
      break;

    case SECURITY_MODE_COMPLETE:
      encode_result = encode_security_mode_complete(
          &msg->security_mode_complete, buffer, len);
      break;

    case SECURITY_MODE_REJECT:
      encode_result =
          encode_security_mode_reject(&msg->security_mode_reject, buffer, len);
      break;

    case SERVICE_REJECT:
      encode_result = encode_service_reject(&msg->service_reject, buffer, len);
      break;

    case SERVICE_REQUEST:
      encode_result =
          encode_service_request(&msg->service_request, buffer, len);
      break;

    case TRACKING_AREA_UPDATE_ACCEPT:
      encode_result = encode_tracking_area_update_accept(
          &msg->tracking_area_update_accept, buffer, len);
      break;

    case TRACKING_AREA_UPDATE_COMPLETE:
      encode_result = encode_tracking_area_update_complete(
          &msg->tracking_area_update_complete, buffer, len);
      break;

    case TRACKING_AREA_UPDATE_REJECT:
      encode_result = encode_tracking_area_update_reject(
          &msg->tracking_area_update_reject, buffer, len);
      break;

    case TRACKING_AREA_UPDATE_REQUEST:
      encode_result = encode_tracking_area_update_request(
          &msg->tracking_area_update_request, buffer, len);
      break;

    case NW_DETACH_REQUEST:
      encode_result =
          encode_nw_detach_request(&msg->nw_detach_request, buffer, len);
      break;

    case UPLINK_NAS_TRANSPORT:
      encode_result =
          encode_uplink_nas_transport(&msg->uplink_nas_transport, buffer, len);
      break;

    default:
      OAILOG_ERROR(
          LOG_NAS_EMM, "EMM-MSG   - Unexpected message type: 0x%x\n",
          msg->header.message_type);
      encode_result = TLV_WRONG_MESSAGE_TYPE;
      /*
       * TODO: Handle not standard layer 3 messages: SERVICE_REQUEST
       */
  }

  if (encode_result < 0) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-MSG   - Failed to encode L3 EMM message 0x%x "
        "(%d)\n",
        msg->header.message_type, encode_result);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, header_result + encode_result);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    emm_msg_decode_header()                                  **
 **                                                                        **
 ** Description: Decode header of EPS Mobility Management message.         **
 **      The protocol discriminator and the security header type   **
 **      have already been decoded.                                **
 **                                                                        **
 ** Inputs:  buffer:    Pointer to the buffer containing the EMM   **
 **             message                                    **
 **          len:       Number of bytes that should be decoded     **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     header:    The EMM message header to be filled        **
 **          Return:    The size of the header if data have been   **
 **             successfully decoded;                      **
 **             A negative error code otherwise.           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_msg_decode_header(
    emm_msg_header_t* header, const uint8_t* buffer, uint32_t len) {
  int size = 0;

  /*
   * Check the buffer length
   */
  if (len < sizeof(emm_msg_header_t)) {
    return (TLV_BUFFER_TOO_SHORT);
  }

  /*
   * Decode the security header type and the protocol discriminator
   */
  DECODE_U8(buffer + size, *(uint8_t*) (header), size);
  /*
   * Decode the message type
   */
  DECODE_U8(buffer + size, header->message_type, size);

  /*
   * Check the protocol discriminator
   */
  if (header->protocol_discriminator != EPS_MOBILITY_MANAGEMENT_MESSAGE) {
    OAILOG_ERROR(
        LOG_NAS_EMM, "ESM-MSG   - Unexpected protocol discriminator: 0x%x\n",
        header->protocol_discriminator);
    return (TLV_PROTOCOL_NOT_SUPPORTED);
  }

  return (size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_msg_encode_header()                                  **
 **                                                                        **
 **      The protocol discriminator and the security header type   **
 **      have already been encoded.                                **
 **                                                                        **
 ** Inputs:  header:    The EMM message header to encode           **
 **          len:       Maximal capacity of the output buffer      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     buffer:    Pointer to the encoded data buffer         **
 **          Return:    The number of bytes in the buffer if data  **
 **             have been successfully encoded;            **
 **             A negative error code otherwise.           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_msg_encode_header(
    const emm_msg_header_t* header, uint8_t* buffer, uint32_t len) {
  int size = 0;

  /*
   * Check the buffer length
   */
  if (len < sizeof(emm_msg_header_t)) {
    return (TLV_BUFFER_TOO_SHORT);
  }

  /*
   * Check the protocol discriminator
   */
  if (header->protocol_discriminator != EPS_MOBILITY_MANAGEMENT_MESSAGE) {
    OAILOG_ERROR(
        LOG_NAS_EMM, "ESM-MSG   - Unexpected protocol discriminator: 0x%x\n",
        header->protocol_discriminator);
    return (TLV_PROTOCOL_NOT_SUPPORTED);
  }

  /*
   * Encode the security header type and the protocol discriminator
   */
  ENCODE_U8(buffer + size, *(uint8_t*) (header), size);
  /*
   * Encode the message type
   */
  if (header->message_type == NW_DETACH_REQUEST) {
    ENCODE_U8(buffer + size, DETACH_REQUEST, size);
  } else {
    ENCODE_U8(buffer + size, header->message_type, size);
  }
  return (size);
}
