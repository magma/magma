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

#include <string.h>  // memcpy
#include <stdlib.h>  // free
#include <stdint.h>
#include <stdbool.h>
#include <netinet/in.h>

#include "log.h"
#include "nas_message.h"
#include "emm_data.h"
#include "secu_defs.h"
#include "dynamic_memory_check.h"
#include "3gpp_24.301.h"
#include "KsiAndSequenceNumber.h"
#include "NasSecurityAlgorithms.h"
#include "ServiceRequest.h"
#include "common_defs.h"
#include "emm_msgDef.h"
#include "nas/securityDef.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

#define SR_MAC_SIZE_BYTES 2

/* Functions used to decode layer 3 NAS messages */

static int nas_message_plain_decode(
    const unsigned char* buffer, const nas_message_security_header_t* header,
    nas_message_plain_t* msg, size_t length);

static int nas_message_protected_decode(
    unsigned char* const buffer, nas_message_security_header_t* header,
    nas_message_plain_t* msg, size_t length,
    emm_security_context_t* const emm_security_context,
    nas_message_decode_status_t* status);

/* Functions used to encode layer 3 NAS messages */
static int nas_message_header_encode(
    unsigned char* buffer, const nas_message_security_header_t* header,
    size_t length);

static int nas_message_plain_encode(
    unsigned char* buffer, const nas_message_security_header_t* header,
    const nas_message_plain_t* msg, size_t length);

static int nas_message_protected_encode(
    unsigned char* buffer, const nas_message_security_protected_t* msg,
    size_t length, void* security);

/* Functions used to decrypt and encrypt layer 3 NAS messages */
static int nas_message_decrypt_a(
    unsigned char* const dest, unsigned char* const src, uint8_t type,
    uint32_t code, uint8_t seq, size_t length,
    emm_security_context_t* const emm_security_context,
    nas_message_decode_status_t* status);

static int nas_message_encrypt_a(
    unsigned char* dest, const unsigned char* src, uint8_t type, uint32_t code,
    uint8_t seq, size_t length,
    emm_security_context_t* const emm_security_context);

/* Functions used for integrity protection of layer 3 NAS messages */
static uint32_t nas_message_get_mac(
    const unsigned char* const buffer, size_t const length, int const direction,
    emm_security_context_t* const emm_security_context);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
 *
 * Name:  nas_message_encrypt()
 *
 * Description: Encripts plain NAS message into security protected NAS
 *    message
 *
 * Inputs:
 *    inbuf:   Input buffer containing plain NAS message
 *    header:  Security protected header to be applied
 *    length:  Number of bytes that should be encrypted
 *    Others:  None
 *
 * Outputs:
 *      outbuf:  Output buffer containing security protected message
 *      Return:  The number of bytes in the output buffer if the input buffer
 * has been successfully encrypted; Negative error code otherwise. Others:  None
 *
 */
int nas_message_encrypt(
    const unsigned char* inbuf, unsigned char* outbuf,
    const nas_message_security_header_t* header, size_t length,
    void* security) {
  OAILOG_FUNC_IN(LOG_NAS);
  emm_security_context_t* emm_security_context =
      (emm_security_context_t*) security;
  int bytes = length;

  /*
   * Encode the header
   */
  int size = nas_message_header_encode(outbuf, header, length);

  if (size < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS, TLV_BUFFER_TOO_SHORT);
  } else if (size > 1) {
    /*
     * Encrypt the plain NAS message.
     * bytes is zero if emm_security_context is null.
     */
    bytes = nas_message_encrypt_a(
        outbuf + size, inbuf, header->security_header_type,
        header->message_authentication_code, header->sequence_number,
        length - size, emm_security_context);

    /*
     * Integrity protected the NAS message
     */
    if (bytes > 0) {
      /*
       * Compute offset of the sequence number field
       */
      int offset = size - sizeof(uint8_t);

      /*
       * Compute the NAS message authentication code
       */
      uint32_t mac = nas_message_get_mac(
          outbuf + offset, bytes + size - offset,
          emm_security_context->direction_encode, emm_security_context);

      /*
       * Set the message authentication code of the NAS message
       */
      uint32_t network_mac = htonl(mac);
      memcpy(
          outbuf + sizeof(uint8_t), (unsigned char*) &network_mac,
          sizeof(uint32_t));
    }
  } else {
    /*
     * The input buffer does not need to be encrypted
     */
    memcpy(outbuf, inbuf, length);
  }

  /*
   * TS 124.301, section 4.4.3.1
   * * * * The NAS sequence number part of the NAS COUNT shall be
   * * * * exchanged between the UE and the MME as part of the
   * * * * NAS signalling. After each new or retransmitted outbound
   * * * * security protected NAS message, the sender shall increase
   * * * * the NAS COUNT number by one. Specifically, on the sender
   * * * * side, the NAS sequence number shall be increased by one,
   * * * * and if the result is zero (due to wrap around), the NAS
   * * * * overflow counter shall also be incremented by one (see
   * * * * subclause 4.4.3.5).
   */
  if (emm_security_context) {
    if (SECU_DIRECTION_DOWNLINK == emm_security_context->direction_encode) {
      emm_security_context->dl_count.seq_num += 1;

      if (!emm_security_context->dl_count.seq_num) {
        emm_security_context->dl_count.overflow += 1;
      }

      OAILOG_DEBUG(
          LOG_NAS, "Incremented emm_security_context.dl_count.seq_num -> %u\n",
          emm_security_context->dl_count.seq_num);
    } else {
      emm_security_context->ul_count.seq_num += 1;

      if (!emm_security_context->ul_count.seq_num) {
        emm_security_context->ul_count.overflow += 1;
      }

      OAILOG_DEBUG(
          LOG_NAS, "Incremented emm_security_context.ul_count.seq_num -> %u\n",
          emm_security_context->ul_count.seq_num);
    }
  }

  if (bytes < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS, bytes);
  }

  if (size > 1) {
    OAILOG_FUNC_RETURN(LOG_NAS, size + bytes);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  nas_message_decrypt()                                     **
 **                                                                        **
 ** Description: Decripts security protected NAS message into plain NAS    **
 **    message                                                   **
 **                                                                        **
 ** Inputs:  inbuf:   Input buffer containing security protected **
 **       NAS message                                **
 **      length:  Number of bytes that should be decrypted   **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   outbuf:  Output buffer containing plain NAS message **
 **    header:  Security protected header applied          **
 **      Return:  The number of bytes in the output buffer   **
 **       if the input buffer has been successfully  **
 **       decrypted; Negative error code otherwise.  **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
int nas_message_decrypt(
    const unsigned char* const inbuf, unsigned char* const outbuf,
    nas_message_security_header_t* header, size_t length, void* security,
    nas_message_decode_status_t* status) {
  OAILOG_FUNC_IN(LOG_NAS);
  emm_security_context_t* emm_security_context =
      (emm_security_context_t*) security;
  int bytes  = length;
  bool is_sr = false;
  /*
   * Decode the header
   */
  int size = nas_message_header_decode(inbuf, header, length, status, &is_sr);

  if (size < 0) {
    OAILOG_DEBUG(LOG_NAS, "MESSAGE TOO SHORT\n");
    OAILOG_FUNC_RETURN(LOG_NAS, TLV_BUFFER_TOO_SHORT);
  } else if (size > 1) {
    if (emm_security_context) {
      status->security_context_available = 1;

      if (emm_security_context->ul_count.seq_num > header->sequence_number) {
        emm_security_context->ul_count.overflow += 1;
      }

      emm_security_context->ul_count.seq_num = header->sequence_number;
    }

    /*
     * Compute offset of the sequence number field
     */
    int offset = size - sizeof(uint8_t);

    /*
     * Compute the NAS message authentication code
     */
    uint32_t mac = nas_message_get_mac(
        inbuf + offset, length - offset, SECU_DIRECTION_UPLINK,
        emm_security_context);

    /*
     * Check NAS message integrity
     */
    if (mac == header->message_authentication_code) {
      status->mac_matched = 1;
      OAILOG_DEBUG(LOG_NAS, "Integrity: MAC Success\n");
    } else {
      OAILOG_CRITICAL(
          LOG_NAS,
          "MAC Failure MSG:%08X(%u) <> INT ALGO:%08X(%u) Type of security "
          "context %u\n",
          header->message_authentication_code,
          header->message_authentication_code, mac, mac,
          (emm_security_context) ? emm_security_context->sc_type : 0);
      // LG: Do not return now (out of spec but we need that with only one MME)
      // LOG_FUNC_RETURN (LOG_NAS, TLV_MAC_MISMATCH);
    }

    /*
     * Decrypt the security protected NAS message
     */
    // OAI_GCC_DIAG_OFF(discarded-qualifiers);
    header->protocol_discriminator = nas_message_decrypt_a(
        outbuf, (unsigned char* const)(inbuf + size),
        header->security_header_type, header->message_authentication_code,
        header->sequence_number, length - size, emm_security_context, status);
    // OAI_GCC_DIAG_ON(discarded-qualifiers);

    bytes = length - size;
  } else {
    OAILOG_DEBUG(LOG_NAS, "Plain NAS message found\n");
    /*
     * The input buffer contains a plain NAS message
     */
    memcpy(outbuf, inbuf, length);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}

/*

   Name:  nas_message_decode()

   Description: Decode layer 3 NAS message

   Inputs:  buffer:  Pointer to the buffer containing layer 3
       NAS message data
       length:  Number of bytes that should be decoded
       security:  security context
       Others:  None

   Outputs:   msg:   L3 NAS message structure to be filled
       Return:  The number of bytes in the buffer if the
         data have been successfully decoded;
         A negative error code otherwise.
       Others:  Return the computed mac if security context is established

*/
int nas_message_decode(
    const unsigned char* const buffer, nas_message_t* msg, size_t length,
    void* security, nas_message_decode_status_t* status) {
  OAILOG_FUNC_IN(LOG_NAS);
  emm_security_context_t* emm_security_context =
      (emm_security_context_t*) security;
  int bytes                    = 0;
  uint32_t mac                 = 0;
  uint16_t short_mac           = 0;
  int size                     = 0;
  bool is_sr                   = false;
  uint8_t sequence_number      = 0;
  uint8_t temp_sequence_number = 0;
  /*
   * Decode the header
   */
  OAILOG_STREAM_HEX(
      OAILOG_LEVEL_DEBUG, LOG_NAS, "Incoming NAS message: ", buffer, length);
  if (emm_security_context) {
    status->security_context_available = 1;
  }
  size =
      nas_message_header_decode(buffer, &msg->header, length, status, &is_sr);

  OAILOG_DEBUG(LOG_NAS, "nas_message_header_decode returned size %d\n", size);

  if (size < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS, TLV_BUFFER_TOO_SHORT);
  }
  if (is_sr) {
    if (length < NAS_MESSAGE_SERVICE_REQUEST_SECURITY_HEADER_SIZE) {
      /*
       * The buffer is not big enough to contain security header
       */
      OAILOG_WARNING(
          LOG_NAS, "Message header %lu bytes is too short %u bytes\n", length,
          NAS_MESSAGE_SECURITY_HEADER_SIZE);
      OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
    }
    // Decode Service Request message. It has different format than any other
    // standard layer-3 message
    DECODE_U8(buffer + size, sequence_number, size);
    DECODE_U16(buffer + size, short_mac, size);

    // shortcut
    msg->plain.emm.header.message_type = SERVICE_REQUEST;
    msg->plain.emm.service_request.ksiandsequencenumber.ksi =
        sequence_number >> 5;
    msg->plain.emm.service_request.ksiandsequencenumber.sequencenumber =
        sequence_number & 0x1F;
    msg->plain.emm.service_request.messageauthenticationcode = short_mac;
    msg->plain.emm.service_request.protocoldiscriminator =
        EPS_MOBILITY_MANAGEMENT_MESSAGE;
    msg->plain.emm.service_request.securityheadertype =
        SECURITY_HEADER_TYPE_SERVICE_REQUEST;
    msg->plain.emm.service_request.messagetype = SERVICE_REQUEST;

    if (emm_security_context == NULL) {
      /*
       * This implies UE context is not present. Send Service Reject with Cause-
       * "UE identity cannot be derived by the network" so that UE can do fresh
       * attach
       */
      status->mac_matched = 0;
      OAILOG_FUNC_RETURN(LOG_NAS, size);
    }

    /*
     * Compute offset of the sequence number field
     */
    // remove ksi
    sequence_number = sequence_number & 0x1F;
    // Estimate 8bit sequence number from 5 bit sequence number
    temp_sequence_number = (emm_security_context->ul_count.seq_num & 0xE0) >> 5;
    if ((emm_security_context->ul_count.seq_num & 0x1F) > sequence_number) {
      temp_sequence_number += 1;
    }
    sequence_number =
        ((temp_sequence_number & 0x07) << 5) | (sequence_number & 0x1F);

    if (emm_security_context->ul_count.seq_num > sequence_number) {
      emm_security_context->ul_count.overflow += 1;
    }
    emm_security_context->ul_count.seq_num = sequence_number;

    /*
     * Compute the NAS message authentication code, return 0 if no security
     * context
     */
    mac = nas_message_get_mac(
        buffer, SR_MAC_SIZE_BYTES, SECU_DIRECTION_UPLINK, emm_security_context);

    /*
     * Check NAS message integrity
     */

    // Compare last 2 LSB bytes for SR
    short_mac = mac & 0x0000FFFF;
    if (short_mac == msg->plain.emm.service_request.messageauthenticationcode) {
      status->mac_matched = 1;
      OAILOG_DEBUG(
          LOG_NAS, "Service Request: message MAC = %04X == computed = %04X\n",
          msg->plain.emm.service_request.messageauthenticationcode, short_mac);
    } else {
      OAILOG_DEBUG(
          LOG_NAS, "Service Request: message MAC = %04X != computed = %04X\n",
          msg->plain.emm.service_request.messageauthenticationcode, short_mac);
    }

    OAILOG_FUNC_RETURN(LOG_NAS, size + bytes);
  }
  if (size > 1) {
    // found security header
    /*
     * Compute offset of the sequence number field
     */
    int offset = size - sizeof(uint8_t);
    if (emm_security_context) {
      status->security_context_available = 1;
      if (SECU_DIRECTION_UPLINK == emm_security_context->direction_decode) {
        if (emm_security_context->ul_count.seq_num >
            msg->header.sequence_number) {
          emm_security_context->ul_count.overflow += 1;
        }

        emm_security_context->ul_count.seq_num = msg->header.sequence_number;
      } else {
        if (emm_security_context->dl_count.seq_num >
            msg->header.sequence_number) {
          emm_security_context->dl_count.overflow += 1;
        }

        emm_security_context->dl_count.seq_num = msg->header.sequence_number;
      }

      /*
       * Compute the NAS message authentication code, return 0 if no security
       * context
       */
      mac = nas_message_get_mac(
          buffer + offset, length - offset,
          emm_security_context->direction_decode, emm_security_context);
      /*
       * Check NAS message integrity
       */
      if (mac == msg->header.message_authentication_code) {
        status->mac_matched = 1;
      } else {
        OAILOG_DEBUG(
            LOG_NAS,
            "msg->header.message_authentication_code = %04X != computed = "
            "%04X\n",
            msg->header.message_authentication_code, mac);
      }
    }

    /*
     * Decode security protected NAS message
     */
    // LG WARNING  msg->plain versus msg->security.plain.
    bytes = nas_message_protected_decode(
        (unsigned char* const)(buffer + size), &msg->header, &msg->plain,
        length - size, emm_security_context, status);
  } else {
    /*
     * Decode plain NAS message
     */
    bytes = nas_message_plain_decode(buffer, &msg->header, &msg->plain, length);
  }

  if (bytes < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS, bytes);
  }

  if (size > 1) {
    OAILOG_FUNC_RETURN(LOG_NAS, size + bytes);
  }
  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  nas_message_encode()                                      **
 **                                                                        **
 ** Description: Encode layer 3 NAS message                                **
 **                                                                        **
 ** Inputs   msg:   L3 NAS message structure to encode         **
 **    length:  Maximal capacity of the output buffer      **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   buffer:  Pointer to the encoded data buffer         **
 **      Return:  The number of bytes in the buffer if the   **
 **       data have been successfully encoded;       **
 **       A negative error code otherwise.           **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
int nas_message_encode(
    unsigned char* buffer, const nas_message_t* const msg, size_t length,
    void* security) {
  OAILOG_FUNC_IN(LOG_NAS);
  emm_security_context_t* emm_security_context =
      (emm_security_context_t*) security;
  int bytes;

  /*
   * Encode the header
   */
  int size = nas_message_header_encode(buffer, &msg->header, length);

  if (size < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS, TLV_BUFFER_TOO_SHORT);
  } else if (size > 1) {
    /*
     * Encode security protected NAS message
     */
    bytes = nas_message_protected_encode(
        buffer + size, &msg->security_protected, length - size,
        emm_security_context);

    /*
     * Integrity protect the NAS message
     */
    if (bytes > 0) {
      /*
       * Compute offset of the sequence number field
       */
      int offset = size - sizeof(uint8_t);

      /*
       * Compute the NAS message authentication code
       */
      OAILOG_DEBUG(
          LOG_NAS,
          "offset %d = %d - %lu, hdr encode = %d, length = %lu bytes = %d\n",
          offset, size, sizeof(uint8_t), size, length, bytes);
      uint32_t mac = nas_message_get_mac(
          buffer + offset, bytes + size - offset,
          emm_security_context->direction_encode, emm_security_context);

      /*
       * Set the message authentication code of the NAS message
       */
      uint32_t network_mac = htonl(mac);
      memcpy(
          buffer + sizeof(uint8_t), (unsigned char*) &network_mac,
          sizeof(uint32_t));

      if (emm_security_context) {
        /*
         * TS 124.301, section 4.4.3.1
         * * * * The NAS sequence number part of the NAS COUNT shall be
         * * * * exchanged between the UE and the MME as part of the
         * * * * NAS signalling. After each new or retransmitted outbound
         * * * * security protected NAS message, the sender shall increase
         * * * * the NAS COUNT number by one. Specifically, on the sender
         * * * * side, the NAS sequence number shall be increased by one,
         * * * * and if the result is zero (due to wrap around), the NAS
         * * * * overflow counter shall also be incremented by one (see
         * * * * subclause 4.4.3.5).
         */
        if (SECU_DIRECTION_DOWNLINK == emm_security_context->direction_encode) {
          emm_security_context->dl_count.seq_num += 1;

          if (!emm_security_context->dl_count.seq_num) {
            emm_security_context->dl_count.overflow += 1;
          }
          OAILOG_DEBUG(
              LOG_NAS,
              "Incremented emm_security_context.dl_count.seq_num -> %u\n",
              emm_security_context->dl_count.seq_num);
        } else {
          emm_security_context->ul_count.seq_num += 1;

          if (!emm_security_context->ul_count.seq_num) {
            emm_security_context->ul_count.overflow += 1;
          }
          OAILOG_DEBUG(
              LOG_NAS,
              "Incremented emm_security_context.ul_count.seq_num -> %u\n",
              emm_security_context->ul_count.seq_num);
        }
      } else {
        OAILOG_DEBUG(
            LOG_NAS,
            "Did not increment emm_security_context.xl_count.seq_num because "
            "no "
            "security context\n");
      }
    }
    /*
     * Log message header
     */

  } else {
    /*
     * Encode plain NAS message
     */
    bytes = nas_message_plain_encode(buffer, &msg->header, &msg->plain, length);
  }

  if (bytes < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS, bytes);
  }

  if (size > 1) {
    OAILOG_FUNC_RETURN(LOG_NAS, size + bytes);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

/*
   -----------------------------------------------------------------------------
      Functions used to decode layer 3 NAS messages
   -----------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:  nas_message_header_decode()                              **
 **                                                                        **
 ** Description: Decode header of a security protected NAS message         **
 **                                                                        **
 ** Inputs:  buffer:  Pointer to the buffer containing layer 3   **
 **       message data                               **
 **      length:  Number of bytes that should be decoded     **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   header:  Security header structure to be filled     **
 **      Return:  The size in bytes of the security header   **
 **       if data have been successfully decoded;    **
 **       1, if the header is not a security header  **
 **       (header of plain NAS message);             **
 **       -1 otherwise.                              **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
int nas_message_header_decode(
    const unsigned char* const buffer,
    nas_message_security_header_t* const header, const size_t length,
    nas_message_decode_status_t* const status, bool* const is_sr) {
  OAILOG_FUNC_IN(LOG_NAS);
  int size = 0;

  /*
   * Decode the first octet of the header (security header type or EPS bearer
   * * * * identity, and protocol discriminator)
   */
  DECODE_U8(buffer, *(uint8_t*) (header), size);

  *is_sr = false;
  if (header->protocol_discriminator == EPS_MOBILITY_MANAGEMENT_MESSAGE) {
    if (header->security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) {
      if (status) {
        switch (header->security_header_type) {
          case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED:
          case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_NEW:
            status->integrity_protected_message = 1;
            break;
          case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED:
          case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED_NEW:
            status->integrity_protected_message = 1;
            status->ciphered_message            = 1;
            break;
          case SECURITY_HEADER_TYPE_SERVICE_REQUEST:
            *is_sr                              = true;
            status->integrity_protected_message = 1;
            /*
             * Current Scope - Service Request message that comes as Initial
             * NAS Message is supported. Note - Service reqeust message which is
             * sent in connected mode due to CSFB comes ciphered as well and is
             * not handled currently. CSFB is not a critical feature from data
             * only service pov and it is not supported.
             */

            OAILOG_FUNC_RETURN(LOG_NAS, size);
            break;
          default:;
        }
      }
      if (*is_sr == false) {
        if (length < NAS_MESSAGE_SECURITY_HEADER_SIZE) {
          /*
           * The buffer is not big enough to contain security header
           */
          OAILOG_WARNING(
              LOG_NAS,
              "NET-API   - The size of the header (%u) "
              "exceeds the buffer length (%lu)\n",
              NAS_MESSAGE_SECURITY_HEADER_SIZE, length);
          OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
        }
        // Decode the message authentication code
        DECODE_U32(buffer + size, header->message_authentication_code, size);
        // Decode the sequence number
        DECODE_U8(buffer + size, header->sequence_number, size);
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:  _nas_message_plain_decode()                               **
 **                                                                        **
 ** Description: Decode plain NAS message                                  **
 **                                                                        **
 ** Inputs:  buffer:  Pointer to the buffer containing plain NAS **
 **       message data                               **
 **    header:  Header of the plain NAS message            **
 **      length:  Number of bytes that should be decoded     **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   msg:   Decoded NAS message                        **
 **      Return:  The number of bytes in the buffer if the   **
 **       data have been successfully decoded;       **
 **       A negative error code otherwise.           **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
static int nas_message_plain_decode(
    const unsigned char* buffer, const nas_message_security_header_t* header,
    nas_message_plain_t* msg, size_t length) {
  OAILOG_FUNC_IN(LOG_NAS);
  int bytes = TLV_PROTOCOL_NOT_SUPPORTED;

  if (header->protocol_discriminator == EPS_MOBILITY_MANAGEMENT_MESSAGE) {
    /*
     * Decode EPS Mobility Management L3 message
     */
    bytes = emm_msg_decode(&msg->emm, (uint8_t*) buffer, length);
  } else if (header->protocol_discriminator == EPS_SESSION_MANAGEMENT_MESSAGE) {
    /*
     * Decode EPS Session Management L3 message
     */
    bytes = esm_msg_decode(&msg->esm, (uint8_t*) buffer, length);
  } else {
    /*
     * Discard L3 messages with not supported protocol discriminator
     */
    OAILOG_WARNING(
        LOG_NAS,
        "NET-API   - Protocol discriminator 0x%x is "
        "not supported\n",
        header->protocol_discriminator);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  _nas_message_protected_decode()                               **
 **                                                                        **
 ** Description: Decode security protected NAS message                     **
 **                                                                        **
 ** Inputs:  buffer:  Pointer to the buffer containing the secu-           **
 **                     rity protected NAS message data                    **
 **          header:  Header of the security protected NAS message       **
 **      length:  Number of bytes that should be decoded             **
 **      emm_security_context: security context                       **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   msg:   Decoded NAS message                        **
 **      Return:  The number of bytes in the buffer if the   **
 **       data have been successfully decoded;       **
 **       A negative error code otherwise.           **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
static int nas_message_protected_decode(
    unsigned char* const buffer, nas_message_security_header_t* header,
    nas_message_plain_t* msg, size_t length,
    emm_security_context_t* const emm_security_context,
    nas_message_decode_status_t* const status) {
  OAILOG_FUNC_IN(LOG_NAS);
  int bytes                      = TLV_BUFFER_TOO_SHORT;
  unsigned char* const plain_msg = (unsigned char*) calloc(1, length);

  if (plain_msg) {
    /*
     * Decrypt the security protected NAS message
     */
    header->protocol_discriminator = nas_message_decrypt_a(
        plain_msg, buffer, header->security_header_type,
        header->message_authentication_code, header->sequence_number, length,
        emm_security_context, status);
    /*
     * Decode the decrypted message as plain NAS message
     */
    bytes = nas_message_plain_decode(plain_msg, header, msg, length);
    free_wrapper((void**) &plain_msg);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}

/*
   -----------------------------------------------------------------------------
      Functions used to encode layer 3 NAS messages
   -----------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:  _nas_message_header_encode()                              **
 **                                                                        **
 ** Description: Encode header of a security protected NAS message         **
 **                                                                        **
 ** Inputs   header:  Security header structure to encode        **
 **    length:  Maximal capacity of the output buffer      **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   buffer:  Pointer to the encoded data buffer         **
 **      Return:  The number of bytes in the buffer if the   **
 **       data have been successfully encoded;       **
 **       1, if the header is not a security header  **
 **       (header of plain NAS message);             **
 **       -1 otherwise.                              **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
static int nas_message_header_encode(
    unsigned char* buffer, const nas_message_security_header_t* header,
    size_t length) {
  OAILOG_FUNC_IN(LOG_NAS);
  int size = 0;

  /*
   * Encode the first octet of the header (security header type or EPS bearer
   * * * * identity, and protocol discriminator)
   */
  ENCODE_U8(buffer, *(uint8_t*) (header), size);

  if (header->protocol_discriminator == EPS_MOBILITY_MANAGEMENT_MESSAGE) {
    if (header->security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) {
      // static uint8_t seq = 0;
      if (length < NAS_MESSAGE_SECURITY_HEADER_SIZE) {
        /*
         * The buffer is not big enough to contain security header
         */
        OAILOG_WARNING(
            LOG_NAS,
            "NET-API   - The size of the header (%u) "
            "exceeds the buffer length (%lu)\n",
            NAS_MESSAGE_SECURITY_HEADER_SIZE, length);
        OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
      }

      /*
       * Encode the message authentication code
       */
      ENCODE_U32(buffer + size, header->message_authentication_code, size);
      /*
       * Encode the sequence number
       */
      ENCODE_U8(buffer + size, header->sequence_number, size);
      // ENCODE_U8(buffer+size, seq, size);
      // seq++;
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:  _nas_message_plain_encode()                               **
 **                                                                        **
 ** Description: Encode plain NAS message                                  **
 **                                                                        **
 ** Inputs:  pd:    Protocol discriminator of the NAS message  **
 **       to encode                                  **
 **      msg:   Plain NAS message structure to encode      **
 **    length:  Maximal capacity of the output buffer      **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   buffer:  Pointer to the encoded data buffer         **
 **      Return:  The number of bytes in the buffer if the   **
 **       data have been successfully encoded;       **
 **       A negative error code otherwise.           **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
static int nas_message_plain_encode(
    unsigned char* buffer, const nas_message_security_header_t* header,
    const nas_message_plain_t* msg, size_t length) {
  OAILOG_FUNC_IN(LOG_NAS);
  int bytes = TLV_PROTOCOL_NOT_SUPPORTED;

  if (EPS_MOBILITY_MANAGEMENT_MESSAGE ==
      msg->emm.header.protocol_discriminator) {
    /*
     * Encode EPS Mobility Management L3 message
     */
    bytes = emm_msg_encode((EMM_msg*) (&msg->emm), (uint8_t*) buffer, length);
  } else if (
      EPS_SESSION_MANAGEMENT_MESSAGE ==
      msg->emm.header.protocol_discriminator) {
    /*
     * Encode EPS Session Management L3 message
     */
    bytes = esm_msg_encode((ESM_msg*) (&msg->esm), (uint8_t*) buffer, length);
  } else {
    /*
     * Discard L3 messages with not supported protocol discriminator
     */
    OAILOG_WARNING(
        LOG_NAS,
        "NET-API   - Protocol discriminator 0x%x is "
        "not supported\n",
        header->protocol_discriminator);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  _nas_message_protected_encode()                               **
 **                                                                        **
 ** Description: Encode security protected NAS message                     **
 **                                                                        **
 ** Inputs    msg:     Security protected NAS message structure            **
 **                    to encode                                           **
 **           length:  Maximal capacity of the output buffer               **
 **           Others:  None                                                **
 **                                                                        **
 ** Outputs:  buffer:  Pointer to the encoded data buffer                  **
 ** Return:   The number of bytes in the buffer if the                     **
 **           data have been successfully encoded;                         **
 **           A negative error code otherwise.                             **
 ** Others:   None                                                         **
 **                                                                        **
 ***************************************************************************/
static int nas_message_protected_encode(
    unsigned char* buffer, const nas_message_security_protected_t* msg,
    size_t length, void* security) {
  OAILOG_FUNC_IN(LOG_NAS);
  emm_security_context_t* emm_security_context =
      (emm_security_context_t*) security;
  int bytes                = TLV_BUFFER_TOO_SHORT;
  unsigned char* plain_msg = (unsigned char*) calloc(1, length);

  if (plain_msg) {
    /*
     * Encode the security protected NAS message as plain NAS message
     */
    int size =
        nas_message_plain_encode(plain_msg, &msg->header, &msg->plain, length);

    if (size > 0) {
      // static uint8_t seq = 0;
      /*
       * Encrypt the encoded plain NAS message
       */
      bytes = nas_message_encrypt_a(
          buffer, plain_msg, msg->header.security_header_type,
          msg->header.message_authentication_code, msg->header.sequence_number,
          size, emm_security_context);
      // seq, size);
      // seq ++;
    }

    free_wrapper((void**) &plain_msg);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}

/*
   -----------------------------------------------------------------------------
        Functions used to decrypt and encrypt layer 3 NAS messages
   -----------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:  nas_message_decrypt_a()                                    **
 **                                                                        **
 ** Description: Decrypt security protected NAS message                    **
 **                                                                        **
 ** Inputs   src:   Pointer to the encrypted data buffer       **
 **    security_header_type:    The security header type                   **
 **    code:    The message authentication code            **
 **    seq:   The sequence number                        **
 **    length:  Maximal capacity of the output buffer      **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   dest:    Pointer to the decrypted data buffer       **
 **      Return:  The protocol discriminator of the message  **
 **       that has been decrypted;                   **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
static int nas_message_decrypt_a(
    unsigned char* const dest, unsigned char* const src,
    uint8_t security_header_type, uint32_t code, uint8_t seq, size_t length,
    emm_security_context_t* const emm_security_context,
    nas_message_decode_status_t* status) {
  OAILOG_FUNC_IN(LOG_NAS);
  nas_stream_cipher_t stream_cipher    = {0};
  uint32_t count                       = 0;
  uint32_t len                         = 0;
  uint8_t direction                    = SECU_DIRECTION_UPLINK;
  int size                             = 0;
  nas_message_security_header_t header = {0};

  switch (security_header_type) {
    case SECURITY_HEADER_TYPE_NOT_PROTECTED:
    case SECURITY_HEADER_TYPE_SERVICE_REQUEST:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_NEW:
      OAILOG_DEBUG(
          LOG_NAS,
          "No decryption of message length %lu according to security header "
          "type "
          "0x%02x\n",
          length, security_header_type);
      len = sizeof(dest);
      memset(dest, 0, len);
      memcpy(dest, src, length);
      DECODE_U8(dest, *(uint8_t*) (&header), size);
      OAILOG_FUNC_RETURN(LOG_NAS, header.protocol_discriminator);
      // LOG_FUNC_RETURN (LOG_NAS, length);
      break;

    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED_NEW:
      if (emm_security_context) {
        direction = emm_security_context->direction_decode;
        switch (emm_security_context->selected_algorithms.encryption) {
          case NAS_SECURITY_ALGORITHMS_EEA1: {
            if (0 == status->mac_matched) {
              OAILOG_ERROR(LOG_NAS, "MAC integrity failed\n");
              OAILOG_FUNC_RETURN(LOG_NAS, 0);
            }
            if (direction == SECU_DIRECTION_UPLINK) {
              count = 0x00000000 |
                      ((emm_security_context->ul_count.overflow & 0x0000FFFF)
                       << 8) |
                      (emm_security_context->ul_count.seq_num & 0x000000FF);
            } else {
              count = 0x00000000 |
                      ((emm_security_context->dl_count.overflow & 0x0000FFFF)
                       << 8) |
                      (emm_security_context->dl_count.seq_num & 0x000000FF);
            }

            OAILOG_DEBUG(
                LOG_NAS,
                "NAS_SECURITY_ALGORITHMS_EEA1 dir %s count.seq_num %u count "
                "%u\n",
                (direction == SECU_DIRECTION_UPLINK) ? "UPLINK" : "DOWNLINK",
                (direction == SECU_DIRECTION_UPLINK) ?
                    emm_security_context->ul_count.seq_num :
                    emm_security_context->dl_count.seq_num,
                count);
            stream_cipher.key        = emm_security_context->knas_enc;
            stream_cipher.key_length = AUTH_KNAS_ENC_SIZE;
            stream_cipher.count      = count;
            stream_cipher.bearer     = 0x00;  // 33.401 section 8.1.1
            stream_cipher.direction  = direction;
            stream_cipher.message    = (uint8_t*) src;
            /*
             * length in bits
             */
            stream_cipher.blength = length << 3;
            nas_stream_encrypt_eea1(&stream_cipher, (uint8_t*) dest);
            /*
             * Decode the first octet (security header type or EPS bearer
             * identity,
             * * * * and protocol discriminator)
             */
            DECODE_U8(dest, *(uint8_t*) (&header), size);
            OAILOG_FUNC_RETURN(LOG_NAS, header.protocol_discriminator);
          } break;

          case NAS_SECURITY_ALGORITHMS_EEA2: {
            if (0 == status->mac_matched) {
              OAILOG_ERROR(LOG_NAS, "MAC integrity failed\n");
              OAILOG_FUNC_RETURN(LOG_NAS, 0);
            }
            if (direction == SECU_DIRECTION_UPLINK) {
              count = 0x00000000 |
                      ((emm_security_context->ul_count.overflow & 0x0000FFFF)
                       << 8) |
                      (emm_security_context->ul_count.seq_num & 0x000000FF);
            } else {
              count = 0x00000000 |
                      ((emm_security_context->dl_count.overflow & 0x0000FFFF)
                       << 8) |
                      (emm_security_context->dl_count.seq_num & 0x000000FF);
            }

            OAILOG_DEBUG(
                LOG_NAS,
                "NAS_SECURITY_ALGORITHMS_EEA2 dir %s count.seq_num %u count "
                "%u\n",
                (direction == SECU_DIRECTION_UPLINK) ? "UPLINK" : "DOWNLINK",
                (direction == SECU_DIRECTION_UPLINK) ?
                    emm_security_context->ul_count.seq_num :
                    emm_security_context->dl_count.seq_num,
                count);
            stream_cipher.key        = emm_security_context->knas_enc;
            stream_cipher.key_length = AUTH_KNAS_ENC_SIZE;
            stream_cipher.count      = count;
            stream_cipher.bearer     = 0x00;  // 33.401 section 8.1.1
            stream_cipher.direction  = direction;
            stream_cipher.message    = (uint8_t*) src;
            /*
             * length in bits
             */
            stream_cipher.blength = length << 3;
            nas_stream_encrypt_eea2(&stream_cipher, (uint8_t*) dest);
            /*
             * Decode the first octet (security header type or EPS bearer
             * identity,
             * * * * and protocol discriminator)
             */
            DECODE_U8(dest, *(uint8_t*) (&header), size);
            OAILOG_FUNC_RETURN(LOG_NAS, header.protocol_discriminator);
          } break;

          case NAS_SECURITY_ALGORITHMS_EEA0:
            OAILOG_DEBUG(
                LOG_NAS,
                "NAS_SECURITY_ALGORITHMS_EEA0 dir %d ul_count.seq_num %d "
                "dl_count.seq_num %d\n",
                direction, emm_security_context->ul_count.seq_num,
                emm_security_context->dl_count.seq_num);
            len = sizeof(dest);
            memset(dest, 0, len);
            memcpy(dest, src, length);
            /*
             * Decode the first octet (security header type or EPS bearer
             * identity,
             * * * * and protocol discriminator)
             */
            DECODE_U8(dest, *(uint8_t*) (&header), size);
            OAILOG_FUNC_RETURN(LOG_NAS, header.protocol_discriminator);
            break;

          default:
            OAILOG_ERROR(
                LOG_NAS, "Unknown Cyphering protection algorithm %d\n",
                emm_security_context->selected_algorithms.encryption);
            len = sizeof(dest);
            memset(dest, 0, len);
            memcpy(dest, src, length);
            /*
             * Decode the first octet (security header type or EPS bearer
             * identity,
             * * * * and protocol discriminator)
             */
            DECODE_U8(dest, *(uint8_t*) (&header), size);
            OAILOG_FUNC_RETURN(LOG_NAS, header.protocol_discriminator);
            break;
        }
      } else {
        OAILOG_ERROR(LOG_NAS, "No security context\n");
        OAILOG_FUNC_RETURN(LOG_NAS, 0);
      }

      break;

    default:
      OAILOG_ERROR(
          LOG_NAS, "Unknown security header type %u", security_header_type);
      OAILOG_FUNC_RETURN(LOG_NAS, 0);
  };
}

/****************************************************************************
 **                                                                        **
 ** Name:  nas_message_encrypt_a()                                    **
 **                                                                        **
 ** Description: Encrypt plain NAS message                                 **
 **                                                                        **
 ** Inputs   src:   Pointer to the decrypted data buffer       **
 **    security_header_type:    The security header type                   **
 **    code:    The message authentication code            **
 **    seq:   The sequence number                        **
 **    direction: The sequence number                        **
 **    length:  Maximal capacity of the output buffer      **
 **    Others:  None                                       **
 **                                                                        **
 ** Outputs:   dest:    Pointer to the encrypted data buffer       **
 **      Return:  The number of bytes in the output buffer   **
 **       if data have been successfully encrypted;  **
 **       RETURNerror otherwise.                     **
 **    Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
static int nas_message_encrypt_a(
    unsigned char* dest, const unsigned char* src, uint8_t security_header_type,
    uint32_t code, uint8_t seq, size_t length,
    emm_security_context_t* const emm_security_context) {
  nas_stream_cipher_t stream_cipher = {0};
  uint32_t count                    = 0;

  OAILOG_FUNC_IN(LOG_NAS);

  if (!emm_security_context) {
    OAILOG_ERROR(
        LOG_NAS,
        "No security context set for encryption protection algorithm\n");
    OAILOG_FUNC_RETURN(LOG_NAS, 0);
  }

  int const direction = emm_security_context->direction_encode;

  switch (security_header_type) {
    case SECURITY_HEADER_TYPE_NOT_PROTECTED:
    case SECURITY_HEADER_TYPE_SERVICE_REQUEST:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_NEW:
      OAILOG_DEBUG(
          LOG_NAS,
          "No encryption of message according to security header type 0x%02x\n",
          security_header_type);
      memcpy(dest, src, length);
      OAILOG_FUNC_RETURN(LOG_NAS, length);
      break;

    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED_NEW:
      switch (emm_security_context->selected_algorithms.encryption) {
        case NAS_SECURITY_ALGORITHMS_EEA1: {
          if (direction == SECU_DIRECTION_UPLINK) {
            count =
                0x00000000 |
                ((emm_security_context->ul_count.overflow & 0x0000FFFF) << 8) |
                (emm_security_context->ul_count.seq_num & 0x000000FF);
          } else {
            count =
                0x00000000 |
                ((emm_security_context->dl_count.overflow & 0x0000FFFF) << 8) |
                (emm_security_context->dl_count.seq_num & 0x000000FF);
          }

          OAILOG_DEBUG(
              LOG_NAS,
              "NAS_SECURITY_ALGORITHMS_EEA1 dir %s count.seq_num %u count %u\n",
              (direction == SECU_DIRECTION_UPLINK) ? "UPLINK" : "DOWNLINK",
              (direction == SECU_DIRECTION_UPLINK) ?
                  emm_security_context->ul_count.seq_num :
                  emm_security_context->dl_count.seq_num,
              count);
          stream_cipher.key        = emm_security_context->knas_enc;
          stream_cipher.key_length = AUTH_KNAS_ENC_SIZE;
          stream_cipher.count      = count;
          stream_cipher.bearer     = 0x00;  // 33.401 section 8.1.1
          stream_cipher.direction  = direction;
          stream_cipher.message    = (uint8_t*) src;
          /*
           * length in bits
           */
          stream_cipher.blength = length << 3;
          nas_stream_encrypt_eea1(&stream_cipher, (uint8_t*) dest);
          OAILOG_FUNC_RETURN(LOG_NAS, length);
        } break;

        case NAS_SECURITY_ALGORITHMS_EEA2: {
          if (direction == SECU_DIRECTION_UPLINK) {
            count =
                0x00000000 |
                ((emm_security_context->ul_count.overflow & 0x0000FFFF) << 8) |
                (emm_security_context->ul_count.seq_num & 0x000000FF);
          } else {
            count =
                0x00000000 |
                ((emm_security_context->dl_count.overflow & 0x0000FFFF) << 8) |
                (emm_security_context->dl_count.seq_num & 0x000000FF);
          }

          OAILOG_DEBUG(
              LOG_NAS,
              "NAS_SECURITY_ALGORITHMS_EEA2 dir %s count.seq_num %u count %u\n",
              (direction == SECU_DIRECTION_UPLINK) ? "UPLINK" : "DOWNLINK",
              (direction == SECU_DIRECTION_UPLINK) ?
                  emm_security_context->ul_count.seq_num :
                  emm_security_context->dl_count.seq_num,
              count);
          stream_cipher.key        = emm_security_context->knas_enc;
          stream_cipher.key_length = AUTH_KNAS_ENC_SIZE;
          stream_cipher.count      = count;
          stream_cipher.bearer     = 0x00;  // 33.401 section 8.1.1
          stream_cipher.direction  = direction;
          stream_cipher.message    = (uint8_t*) src;
          /*
           * length in bits
           */
          stream_cipher.blength = length << 3;
          nas_stream_encrypt_eea2(&stream_cipher, (uint8_t*) dest);
          OAILOG_FUNC_RETURN(LOG_NAS, length);
        } break;

        case NAS_SECURITY_ALGORITHMS_EEA0:
          OAILOG_DEBUG(
              LOG_NAS,
              "NAS_SECURITY_ALGORITHMS_EEA0 dir %d ul_count.seq_num %d "
              "dl_count.seq_num %d\n",
              direction, emm_security_context->ul_count.seq_num,
              emm_security_context->dl_count.seq_num);
          memcpy(dest, src, length);
          OAILOG_FUNC_RETURN(LOG_NAS, length);
          break;

        default:
          OAILOG_ERROR(
              LOG_NAS, "Unknown Cyphering protection algorithm %d\n",
              emm_security_context->selected_algorithms.encryption);
          break;
      }

      break;

    default:
      OAILOG_ERROR(
          LOG_NAS, "Unknown security header type %u\n", security_header_type);
      OAILOG_FUNC_RETURN(LOG_NAS, 0);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, length);
}

/*
   -----------------------------------------------------------------------------
    Functions used for integrity protection of layer 3 NAS messages
   -----------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:  _nas_message_get_mac()                                        **
 **                                                                        **
 ** Description: Run integrity algorithm onto cyphered or uncyphered NAS   **
 **    message encoded in the input buffer and return the compu-         **
 **    ted message authentication code                                   **
 **                                                                        **
 ** Inputs   buffer:  Pointer to the integrity protected data            **
 **       buffer                                                     **
 **    count:   Value of the uplink NAS counter                        **
 **    length:  Length of the input buffer                             **
 **      direction                                                         **
 **    Others:  None                                                   **
 **                                                                        **
 ** Outputs:   None                                                      **
 **      Return:  The message authentication code                    **
 **    Others:  None                                                   **
 **                                                                        **
 ***************************************************************************/
static uint32_t nas_message_get_mac(
    const unsigned char* const buffer, size_t const length, int const direction,
    emm_security_context_t* const emm_security_context) {
  OAILOG_FUNC_IN(LOG_NAS);

  if (!emm_security_context) {
    OAILOG_DEBUG(
        LOG_NAS,
        "No security context set for integrity protection algorithm\n");
    OAILOG_FUNC_RETURN(LOG_NAS, 0);
  }

  switch (emm_security_context->selected_algorithms.integrity) {
    case NAS_SECURITY_ALGORITHMS_EIA1: {
      uint8_t mac[4];
      nas_stream_cipher_t stream_cipher;
      uint32_t count;
      uint32_t* mac32;

      if (direction == SECU_DIRECTION_UPLINK) {
        count = 0x00000000 |
                ((emm_security_context->ul_count.overflow & 0x0000FFFF) << 8) |
                (emm_security_context->ul_count.seq_num & 0x000000FF);
      } else {
        count = 0x00000000 |
                ((emm_security_context->dl_count.overflow & 0x0000FFFF) << 8) |
                (emm_security_context->dl_count.seq_num & 0x000000FF);
      }

      OAILOG_DEBUG(
          LOG_NAS,
          "NAS_SECURITY_ALGORITHMS_EIA1 dir %s count.seq_num %u count %u\n",
          (direction == SECU_DIRECTION_UPLINK) ? "UPLINK" : "DOWNLINK",
          (direction == SECU_DIRECTION_UPLINK) ?
              emm_security_context->ul_count.seq_num :
              emm_security_context->dl_count.seq_num,
          count);
      stream_cipher.key        = emm_security_context->knas_int;
      stream_cipher.key_length = AUTH_KNAS_INT_SIZE;
      stream_cipher.count      = count;
      stream_cipher.bearer     = 0x00;  // 33.401 section 8.1.1
      stream_cipher.direction  = direction;
      stream_cipher.message    = (uint8_t*) buffer;
      /*
       * length in bits
       */
      stream_cipher.blength = length << 3;
      nas_stream_encrypt_eia1(&stream_cipher, mac);
      OAILOG_DEBUG(
          LOG_NAS,
          "NAS_SECURITY_ALGORITHMS_EIA1 returned MAC %x.%x.%x.%x(%u) for "
          "length "
          "%lu direction %d, count %d\n",
          mac[0], mac[1], mac[2], mac[3], *((uint32_t*) &mac), length,
          direction, count);
      mac32 = (uint32_t*) &mac;
      OAILOG_FUNC_RETURN(LOG_NAS, ntohl(*mac32));
    } break;

    case NAS_SECURITY_ALGORITHMS_EIA2: {
      uint8_t mac[4];
      nas_stream_cipher_t stream_cipher;
      uint32_t count;
      uint32_t* mac32;

      if (direction == SECU_DIRECTION_UPLINK) {
        count = 0x00000000 |
                ((emm_security_context->ul_count.overflow & 0x0000FFFF) << 8) |
                (emm_security_context->ul_count.seq_num & 0x000000FF);
      } else {
        count = 0x00000000 |
                ((emm_security_context->dl_count.overflow & 0x0000FFFF) << 8) |
                (emm_security_context->dl_count.seq_num & 0x000000FF);
      }

      OAILOG_DEBUG(
          LOG_NAS,
          "NAS_SECURITY_ALGORITHMS_EIA2 dir %s count.seq_num %u count %u\n",
          (direction == SECU_DIRECTION_UPLINK) ? "UPLINK" : "DOWNLINK",
          (direction == SECU_DIRECTION_UPLINK) ?
              emm_security_context->ul_count.seq_num :
              emm_security_context->dl_count.seq_num,
          count);
      stream_cipher.key        = emm_security_context->knas_int;
      stream_cipher.key_length = AUTH_KNAS_INT_SIZE;
      stream_cipher.count      = count;
      stream_cipher.bearer     = 0x00;  // 33.401 section 8.1.1
      stream_cipher.direction  = direction;
      stream_cipher.message    = (uint8_t*) buffer;
      /*
       * length in bits
       */
      stream_cipher.blength = length << 3;
      nas_stream_encrypt_eia2(&stream_cipher, mac);
      OAILOG_DEBUG(
          LOG_NAS,
          "NAS_SECURITY_ALGORITHMS_EIA2 returned MAC %x.%x.%x.%x(%u) for "
          "length "
          "%lu direction %d, count %d\n",
          mac[0], mac[1], mac[2], mac[3], *((uint32_t*) &mac), length,
          direction, count);
      mac32 = (uint32_t*) &mac;
      OAILOG_FUNC_RETURN(LOG_NAS, ntohl(*mac32));
    } break;

    case NAS_SECURITY_ALGORITHMS_EIA0:
      OAILOG_DEBUG(
          LOG_NAS, "NAS_SECURITY_ALGORITHMS_EIA0 dir %s count.seq_num %u\n",
          (direction == SECU_DIRECTION_UPLINK) ? "UPLINK" : "DOWNLINK",
          (direction == SECU_DIRECTION_UPLINK) ?
              emm_security_context->ul_count.seq_num :
              emm_security_context->dl_count.seq_num);
      OAILOG_FUNC_RETURN(LOG_NAS, 0);
      break;

    default:
      OAILOG_ERROR(
          LOG_NAS, "Unknown integrity protection algorithm %d\n",
          emm_security_context->selected_algorithms.integrity);
      break;
  }

  OAILOG_FUNC_RETURN(LOG_NAS, 0);
}
