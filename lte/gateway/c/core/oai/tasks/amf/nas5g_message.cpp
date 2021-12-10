#include <string.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <netinet/in.h>
#include <sstream>
#include <iostream>
#include <iomanip>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/secu/secu_defs.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/amf/amf_data.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.301.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_as.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_fsm.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GDLNASTransport.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_common.h"
namespace magma5g {

#define NAS5G_MESSAGE_SECURITY_HEADER_SIZE 7

/* Functions used to decode layer 3 NAS messages */
int nas5g_message_header_decode(
    const unsigned char* const buffer, amf_msg_header* const header,
    const uint32_t length, amf_nas_message_decode_status_t* const status,
    bool* const is_sr);

static int _nas5g_message_plain_decode(
    const unsigned char* buffer, const amf_msg_header* header,
    nas_message_plain_t* msg, uint32_t length);

static int _nas5g_message_protected_decode(
    unsigned char* const buffer, amf_msg_header* header,
    nas_message_plain_t* msg, uint32_t length,
    amf_security_context_t* const amf_security_context,
    amf_nas_message_decode_status_t* status);

/* Functions used to encode layer 3 NAS messages */
static int _nas5g_message_header_encode(
    unsigned char* buffer, const amf_msg_header* header, uint32_t length);

int _nas5g_message_plain_encode(
    unsigned char* buffer, const amf_msg_header* header,
    const nas_message_plain_t* msg, uint32_t length);

static int _nas5g_message_protected_encode(
    unsigned char* buffer, const nas_message_security_protected_t* msg,
    uint32_t length, void* security);

/* Functions used to decrypt and encrypt layer 3 NAS messages */
static int _nas5g_message_decrypt(
    unsigned char* const dest, unsigned char* const src, uint8_t type,
    uint32_t code, uint8_t seq, uint32_t length,
    amf_security_context_t* const amf_security_context,
    amf_nas_message_decode_status_t* status);

static int _nas5g_message_encrypt(
    unsigned char* dest, const unsigned char* src, uint8_t type, uint32_t code,
    uint8_t seq, int const direction, uint32_t length,
    amf_security_context_t* const amf_security_context);

/* Functions used for integrity protection of layer 3 NAS messages */
static uint32_t _nas5g_message_get_mac(
    const unsigned char* const buffer, uint32_t const length,
    int const direction, amf_security_context_t* const amf_security_context);

std::string get_message_type_str(uint8_t type);

/****************************************************************************
 *                                                                           *
 *   Name:  nas5g_message_decode()                                           *
 *                                                                           *
 *   Description: Decode layer 3 NAS message                                 *
 *                                                                           *
 *   Inputs:  buffer:  Pointer to the buffer containing layer 3              *
 *                     NAS message data                                      *
 *            length:  Number of bytes that should be decoded                *
 *            security:  security context                                    *
 *            Others:  None                                                  *
 *                                                                           *
 *   Outputs: msg:   L3 NAS message structure to be filled                   *
 *            Return:  The number of bytes in the buffer if the              *
 *                     data have been successfully decoded;                  *
 *                     A negative error code otherwise.                      *
 *            Others:  Return the computed mac if security context is        *
 *                     established                                           *
 *                                                                           *
 ****************************************************************************/
int nas5g_message_decode(
    const unsigned char* const buffer, amf_nas_message_t* msg, uint32_t length,
    void* security, amf_nas_message_decode_status_t* status) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_security_context_t* amf_security_context =
      (amf_security_context_t*) security;
  int bytes    = 0;
  uint32_t mac = 0;
  int size     = 0;
  bool is_sr   = false;
  /*
   * Decode the header
   */
  if (amf_security_context) {
    status->security_context_available = 1;
  }
  size =
      nas5g_message_header_decode(buffer, &msg->header, length, status, &is_sr);
  if (size < 0) {
    OAILOG_ERROR(LOG_AMF_APP, "NAS Header decode failed");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, TLV_BUFFER_TOO_SHORT);
  }
  if (size > 2) {
    // found security header
    /*
     * Compute offset of the sequence number field
     */
    int offset = size - sizeof(uint8_t);
    if (amf_security_context) {
      status->security_context_available = 1;
      if (SECU_DIRECTION_UPLINK == amf_security_context->direction_decode) {
        if (amf_security_context->ul_count.seq_num >
            msg->header.sequence_number) {
          amf_security_context->ul_count.overflow += 1;
        }

        amf_security_context->ul_count.seq_num = msg->header.sequence_number;
      } else {
        if (amf_security_context->dl_count.seq_num >
            msg->header.sequence_number) {
          amf_security_context->dl_count.overflow += 1;
        }

        amf_security_context->dl_count.seq_num = msg->header.sequence_number;
      }

      /*
       * Compute the NAS message authentication code, return 0 if no security
       * context
       */
      mac = _nas5g_message_get_mac(
          buffer + offset, length - offset,
          amf_security_context->direction_decode, amf_security_context);
      /*
       * Check NAS message integrity
       */
      if (mac == msg->header.message_authentication_code) {
        status->mac_matched = 1;
      } else {
        OAILOG_DEBUG(
            LOG_AMF_APP,
            "msg->header.message_authentication_code = %04X != computed = "
            "%04X\n",
            msg->header.message_authentication_code, mac);
      }
    }

    /*
     * Decode security protected NAS message
     */
    bytes = _nas5g_message_protected_decode(
        (unsigned char* const)(buffer + size), &msg->header, &msg->plain,
        length - size, amf_security_context, status);
  } else {
    /*
     * Decode plain NAS message
     */
    bytes =
        _nas5g_message_plain_decode(buffer, &msg->header, &msg->plain, length);
  }

  if (bytes < 0) {
    OAILOG_ERROR(LOG_AMF_APP, "NAS Decode failed");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, bytes);
  }
  OAILOG_DEBUG(
      LOG_AMF_APP, "[%s] Msg plain decode bytes[0-%d]\n%s",
      get_message_type_str(msg->plain.amf.header.message_type).c_str(), bytes,
      uint8_to_hex_string(buffer, bytes).c_str());

  OAILOG_INFO(
      LOG_AMF_APP, "Decoded msg(nas5g) id: [%x]-name [%s]",
      msg->plain.amf.header.message_type,
      get_message_type_str(msg->plain.amf.header.message_type).c_str());

  if (size > 1) {
    // Security Protected NAS message decoded
    OAILOG_FUNC_RETURN(LOG_AMF_APP, size + bytes);
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  nas5g_message_encode()                                          **
 **                                                                        **
 ** Description: Encode layer 3 NAS message                                **
 **                                                                        **
 ** Inputs  msg:   L3 NAS message structure to encode                      **
 **         length:  Maximal capacity of the output buffer                 **
 **         Others:  None                                                  **
 **                                                                        **
 ** Outputs:   buffer:  Pointer to the encoded data buffer                 **
 **            Return:  The number of bytes in the buffer if the           **
 **                     data have been successfully encoded;               **
 **                     A negative error code otherwise.                   **
 **            Others:  None                                               **
 **                                                                        **
 ***************************************************************************/
int nas5g_message_encode(
    unsigned char* buffer, const amf_nas_message_t* const msg, uint32_t length,
    void* security) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_security_context_t* amf_security_context =
      (amf_security_context_t*) security;
  int bytes;

  /*
   * Encode the header
   */
  int size = _nas5g_message_header_encode(buffer, &msg->header, length);
  if (size < 0) {
    OAILOG_ERROR(LOG_AMF_APP, "NAS Header encode failed");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, TLV_BUFFER_TOO_SHORT);
  } else if (size > 2) {
    /*
     * Encode security protected NAS message
     */
    bytes = _nas5g_message_protected_encode(
        buffer + size, &msg->security_protected, length - size,
        amf_security_context);
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
          LOG_AMF_APP,
          "Offset %d = %d - %lu, hdr encode = %d, length = %" PRIu32
          "bytes = %d\n",
          offset, size, sizeof(uint8_t), size, length, bytes);
      uint32_t mac = _nas5g_message_get_mac(
          buffer + offset, bytes + size - offset,
          amf_security_context->direction_encode, amf_security_context);
      /*
       * Set the message authentication code of the NAS message
       */
      *(uint32_t*) (buffer + sizeof(uint16_t)) = htonl(mac);

      if (amf_security_context) {
        /*
         * TS 124.301, section 4.4.3.1
         * * * * The NAS sequence number part of the NAS COUNT shall be
         * * * * exchanged between the UE and the MME as part of the
         * * * * NAS signaling. After each new or retransmitted outbound
         * * * * security protected NAS message, the sender shall increase
         * * * * the NAS COUNT number by one. Specifically, on the sender
         * * * * side, the NAS sequence number shall be increased by one,
         * * * * and if the result is zero (due to wrap around), the NAS
         * * * * overflow counter shall also be incremented by one (see
         * * * * subclause 4.4.3.5).
         */
        if (SECU_DIRECTION_DOWNLINK == amf_security_context->direction_encode) {
          amf_security_context->dl_count.seq_num += 1;

          if (!amf_security_context->dl_count.seq_num) {
            amf_security_context->dl_count.overflow += 1;
          }
          OAILOG_DEBUG(
              LOG_AMF_APP,
              "Incremented amf_security_context.dl_count.seq_num -> %u\n",
              amf_security_context->dl_count.seq_num);
        } else {
          amf_security_context->ul_count.seq_num += 1;

          if (!amf_security_context->ul_count.seq_num) {
            amf_security_context->ul_count.overflow += 1;
          }
          OAILOG_DEBUG(
              LOG_AMF_APP,
              "Incremented amf_security_context.ul_count.seq_num -> %u\n",
              amf_security_context->ul_count.seq_num);
        }
      } else {
        OAILOG_DEBUG(
            LOG_AMF_APP,
            "Did not increment amf_security_context.xl_count.seq_num because "
            "no "
            "security context\n");
      }
    }
    OAILOG_INFO(
        LOG_AMF_APP, "Encoded msg(nas5g) id: [%x]-name [%s]",
        msg->security_protected.plain.amf.header.message_type,
        get_message_type_str(
            msg->security_protected.plain.amf.header.message_type)
            .c_str());
  } else {
    /*
     * Encode plain NAS message
     */
    bytes =
        _nas5g_message_plain_encode(buffer, &msg->header, &msg->plain, length);

    OAILOG_INFO(
        LOG_AMF_APP, "Encoded msg(nas5g) id: [%x]-name [%s]",
        msg->plain.amf.header.message_type,
        get_message_type_str(msg->plain.amf.header.message_type).c_str());
  }

  if (bytes < 0) {
    OAILOG_ERROR(LOG_AMF_APP, "NAS Encode failed");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, bytes);
  }

  if (size > 2) {
    // Security Protected NAS message encoded
    OAILOG_FUNC_RETURN(LOG_AMF_APP, size + bytes);
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  nas5g_message_header_decode()                                   **
 **                                                                        **
 ** Description: Decode header of a security protected NAS message         **
 **                                                                        **
 ** Inputs:  buffer:  Pointer to the buffer containing layer 3             **
 **                   message data                                         **
 **          length:  Number of bytes that should be decoded               **
 **          Others:  None                                                 **
 **                                                                        **
 ** Outputs:   header:  Security header structure to be filled             **
 **            Return:  The size in bytes of the security header           **
 **                     if data have been successfully decoded;            **
 **                     1, if the header is not a security header          **
 **                     (header of plain NAS message);                     **
 **                     -1 otherwise.                                      **
 **            Others:  None                                               **
 **                                                                        **
 ***************************************************************************/
int nas5g_message_header_decode(
    const unsigned char* const buffer, amf_msg_header* const header,
    const uint32_t length, amf_nas_message_decode_status_t* const status,
    bool* const is_sr) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  int size = 0;

  /*
   * Decode the first octet of the header
   */
  header->extended_protocol_discriminator = *buffer;
  size++;
  header->security_header_type = *(buffer + 1) & 0xf;
  size++;
  *is_sr = false;
  if (header->extended_protocol_discriminator ==
      M5GS_MOBILITY_MANAGEMENT_MESSAGE) {
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
            OAILOG_FUNC_RETURN(LOG_AMF_APP, size);
            break;
          default:
            OAILOG_WARNING(LOG_AMF_APP, "Unknown security header type");
        }
      }
      if (*is_sr == false) {
        if (length < NAS5G_MESSAGE_SECURITY_HEADER_SIZE) {
          /*
           * The buffer is not big enough to contain security header
           */
          OAILOG_WARNING(
              LOG_AMF_APP,
              "NET-API   - The size of the header (%u) "
              "exceeds the buffer length %" PRIu32 "\n",
              NAS5G_MESSAGE_SECURITY_HEADER_SIZE, length);
          OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
        }
        // Decode the message authentication code
        DECODE_U32(buffer + size, header->message_authentication_code, size);
        // Decode the sequence number
        DECODE_U8(buffer + size, header->sequence_number, size);
        OAILOG_DEBUG(
            LOG_AMF_APP,
            "epd:%x, security_header_type:%x, "
            "message_authentication_code:%x, sequence_number:%x",
            header->extended_protocol_discriminator,
            header->security_header_type, header->message_authentication_code,
            header->sequence_number);
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:  _nas5g_message_plain_decode()                                   **
 **                                                                        **
 ** Description: Decode plain NAS message                                  **
 **                                                                        **
 ** Inputs:  buffer:  Pointer to the buffer containing plain NAS           **
 **                   message data                                         **
 **          header:  Header of the plain NAS message                      **
 **          length:  Number of bytes that should be decoded               **
 **          Others:  None                                                 **
 **                                                                        **
 ** Outputs:   msg:   Decoded NAS message                                  **
 **            Return:  The number of bytes in the buffer if the           **
 **                     data have been successfully decoded;               **
 **                     A negative error code otherwise.                   **
 **            Others:  None                                               **
 **                                                                        **
 ***************************************************************************/
static int _nas5g_message_plain_decode(
    const unsigned char* buffer, const amf_msg_header* header,
    nas_message_plain_t* msg, uint32_t length) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  int bytes = TLV_PROTOCOL_NOT_SUPPORTED;
  AmfMsg amf_msg;
  if (header->extended_protocol_discriminator ==
      M5GS_MOBILITY_MANAGEMENT_MESSAGE) {
    /*
     * Decode Mobility Management L3
     */
    bytes = amf_msg.M5gNasMessageDecodeMsg(
        (AmfMsg*) &msg->amf, (uint8_t*) buffer, length);
  } else {
    /*
     * Discard L3 messages with not supported protocol discriminator
     */
    OAILOG_WARNING(
        LOG_AMF_APP,
        "NET-API   - Protocol discriminator 0x%x is "
        "not supported\n",
        header->extended_protocol_discriminator);
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  _nas5g_message_protected_decode()                               **
 **                                                                        **
 ** Description: Decode security protected NAS message                     **
 **                                                                        **
 ** Inputs:  buffer:  Pointer to the buffer containing the secu-           **
 **                   rity protected NAS message data                      **
 **          header:  Header of the security protected NAS message         **
 **          length:  Number of bytes that should be decoded               **
 **                   amf_security_context: security context               **
 **          Others:  None                                                 **
 **                                                                        **
 ** Outputs:   msg:     Decoded NAS message                                **
 **            Return:  The number of bytes in the buffer if the           **
 **                     data have been successfully decoded;               **
 **                     A negative error code otherwise.                   **
 **            Others:  None                                               **
 **                                                                        **
 ***************************************************************************/
static int _nas5g_message_protected_decode(
    unsigned char* const buffer, amf_msg_header* header,
    nas_message_plain_t* msg, uint32_t length,
    amf_security_context_t* const amf_security_context,
    amf_nas_message_decode_status_t* const status) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  int bytes                      = TLV_BUFFER_TOO_SHORT;
  unsigned char* const plain_msg = (unsigned char*) calloc(1, length);

  if (plain_msg) {
    /*
     * Decrypt the security protected NAS message
     */
    header->extended_protocol_discriminator = _nas5g_message_decrypt(
        plain_msg, buffer, header->security_header_type,
        header->message_authentication_code, header->sequence_number, length,
        amf_security_context, status);
    /*
     * Decode the decrypted message as plain NAS message
     */
    bytes = _nas5g_message_plain_decode(plain_msg, header, msg, length);

    free(plain_msg);
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, bytes);
}

/*
   -----------------------------------------------------------------------------
      Functions used to encode layer 3 NAS messages
   -----------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:  _nas5g_message_header_encode()                                  **
 **                                                                        **
 ** Description: Encode header of a security protected NAS message         **
 **                                                                        **
 ** Inputs   header:  Security header structure to encode                  **
 **          length:  Maximal capacity of the output buffer                **
 **          Others:  None                                                 **
 **                                                                        **
 ** Outputs:   buffer:  Pointer to the encoded data buffer                 **
 **            Return:  The number of bytes in the buffer if the           **
 **                     data have been successfully encoded;               **
 **                     1, if the header is not a security header          **
 **                     (header of plain NAS message);                     **
 **                     -1 otherwise.                                      **
 **            Others:  None                                               **
 **                                                                        **
 ***************************************************************************/
static int _nas5g_message_header_encode(
    unsigned char* buffer, const amf_msg_header* header, uint32_t length) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  int size = 0;

  /*
   * Encode the first octet of the header
   */
  *buffer = header->extended_protocol_discriminator;
  size++;
  *(buffer + size) = header->security_header_type & 0xf;
  size++;
  if (header->extended_protocol_discriminator ==
      M5GS_MOBILITY_MANAGEMENT_MESSAGE) {
    if (header->security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) {
      if (length < NAS5G_MESSAGE_SECURITY_HEADER_SIZE) {
        /*
         * The buffer is not big enough to contain security header
         */
        OAILOG_WARNING(
            LOG_AMF_APP,
            "NET-API   - The size of the header (%u) "
            "exceeds the buffer length %" PRIu32 "\n",
            NAS5G_MESSAGE_SECURITY_HEADER_SIZE, length);
        OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
      }

      /*
       * Encode the message authentication code
       */
      ENCODE_U32(buffer + size, header->message_authentication_code, size);
      /*
       * Encode the sequence number
       */
      ENCODE_U8(buffer + size, header->sequence_number, size);
    }
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:  _nas5g_message_plain_encode()                                   **
 **                                                                        **
 ** Description: Encode plain NAS message                                  **
 **                                                                        **
 ** Inputs:    pd:    Protocol discriminator of the NAS message            **
 *  *                 to encode                                            **
 **            msg:   Plain NAS message structure to encode                **
 **          length:  Maximal capacity of the output buffer                **
 **          Others:  None                                                 **
 **                                                                        **
 ** Outputs:   buffer:  Pointer to the encoded data buffer                 **
 **            Return:  The number of bytes in the buffer if the           **
 **                     data have been successfully encoded;               **
 **                     A negative error code otherwise.                   **
 **            Others:  None                                               **
 **                                                                        **
 ***************************************************************************/
int _nas5g_message_plain_encode(
    unsigned char* buffer, const amf_msg_header* header,
    const nas_message_plain_t* msg, uint32_t length) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  int bytes = TLV_PROTOCOL_NOT_SUPPORTED;
  AmfMsg amf_msg_test;
  if (M5GS_MOBILITY_MANAGEMENT_MESSAGE ==
      msg->amf.header.extended_protocol_discriminator) {
    /*
     * Encode Mobility Management L3 message
     */
    bytes = amf_msg_test.M5gNasMessageEncodeMsg(
        (AmfMsg*) &msg->amf, (uint8_t*) buffer, length);

    if (bytes < 0) {
      OAILOG_WARNING(LOG_AMF_APP, "Encoding Message Failed");
      OAILOG_FUNC_RETURN(LOG_AMF_APP, bytes);
    }
    OAILOG_DEBUG(
        LOG_AMF_APP, "[%s] Msg plain encode bytes[0-%d]\n%s",
        get_message_type_str(msg->amf.header.message_type).c_str(), bytes,
        uint8_to_hex_string(buffer, bytes).c_str());
  } else {
    /*
     * Discard L3 messages with not supported protocol discriminator
     */
    OAILOG_WARNING(
        LOG_AMF_APP,
        "NET-API   - Protocol discriminator 0x%x is "
        "not supported\n",
        header->extended_protocol_discriminator);
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  _nas5g_message_protected_encode()                               **
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
static int _nas5g_message_protected_encode(
    unsigned char* buffer, const nas_message_security_protected_t* msg,
    uint32_t length, void* security) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_security_context_t* amf_security_context =
      (amf_security_context_t*) security;
  int bytes                = TLV_BUFFER_TOO_SHORT;
  unsigned char* plain_msg = (unsigned char*) calloc(1, length);

  if (plain_msg) {
    /*
     * Encode the security protected NAS message as plain NAS message
     */
    int size = _nas5g_message_plain_encode(
        plain_msg, &msg->header, &msg->plain, length);
    if (size > 0 && security) {
      /*
       * Encrypt the encoded plain NAS message
       */
      bytes = _nas5g_message_encrypt(
          buffer, plain_msg, msg->header.security_header_type,
          msg->header.message_authentication_code, msg->header.sequence_number,
          amf_security_context->direction_encode, size, amf_security_context);
    }
  }
  free(plain_msg);

  OAILOG_FUNC_RETURN(LOG_AMF_APP, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  _nas5g_message_decrypt()                                        **
 **                                                                        **
 ** Description: Decrypt security protected NAS message                    **
 **                                                                        **
 ** Inputs   src:   Pointer to the encrypted data buffer                   **
 **          security_header_type:    The security header type             **
 **          code:    The message authentication code                      **
 **          seq:   The sequence number                                    **
 **          length:  Maximal capacity of the output buffer                **
 **          Others:  None                                                 **
 **                                                                        **
 ** Outputs:   dest:    Pointer to the decrypted data buffer               **
 **            Return:  The protocol discriminator of the message          **
 **                     that has been decrypted;                           **
 **            Others:  None                                               **
 **                                                                        **
 ***************************************************************************/
static int _nas5g_message_decrypt(
    unsigned char* const dest, unsigned char* const src,
    uint8_t security_header_type, uint32_t code, uint8_t seq, uint32_t length,
    amf_security_context_t* const amf_security_context,
    amf_nas_message_decode_status_t* status) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  uint8_t direction     = SECU_DIRECTION_UPLINK;
  int size              = 0;
  amf_msg_header header = {0};
  switch (security_header_type) {
    case SECURITY_HEADER_TYPE_NOT_PROTECTED:
    case SECURITY_HEADER_TYPE_SERVICE_REQUEST:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_NEW:
      OAILOG_DEBUG(
          LOG_AMF_APP,
          "No decryption of message length %" PRIu32
          "according to security header "
          "type "
          "0x%02x\n",
          length, security_header_type);

      memset(dest, 0, length);
      memcpy(dest, src, length);

      DECODE_U8(dest, *(uint8_t*) (&header), size);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, header.extended_protocol_discriminator);
      break;

    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED_NEW:
      if (amf_security_context) {
        direction = amf_security_context->direction_decode;
        switch (amf_security_context->selected_algorithms.encryption) {
          case M5G_NAS_SECURITY_ALGORITHMS_5G_EA0:
            OAILOG_DEBUG(
                LOG_AMF_APP,
                "M5G_NAS_SECURITY_ALGORITHMS_5G_EA0 dir %d ul_count.seq_num %d "
                "dl_count.seq_num %d\n",
                direction, amf_security_context->ul_count.seq_num,
                amf_security_context->dl_count.seq_num);
            memset(dest, 0, length);
            memcpy(dest, src, length);
            /*
             * Decode the first octet (security header type or EPS bearer
             * identity,
             * * * * and protocol discriminator)
             */
            DECODE_U8(dest, *(uint8_t*) (&header), size);
            OAILOG_FUNC_RETURN(
                LOG_AMF_APP, header.extended_protocol_discriminator);
            break;

          default:
            OAILOG_ERROR(
                LOG_AMF_APP, "Unsupported Cyphering protection algorithm %d\n",
                amf_security_context->selected_algorithms.encryption);

            memset(dest, 0, length);
            memcpy(dest, src, length);
            /*
             * Decode the first octet (security header type or EPS bearer
             * identity,
             * * * * and protocol discriminator)
             */
            DECODE_U8(dest, *(uint8_t*) (&header), size);
            OAILOG_FUNC_RETURN(
                LOG_AMF_APP, header.extended_protocol_discriminator);
            break;
        }
      } else {
        OAILOG_ERROR(LOG_AMF_APP, "No security context\n");
        memset(dest, 0, length);
        memcpy(dest, src, length);

        DECODE_U8(dest, *(uint8_t*) (&header), size);
        OAILOG_FUNC_RETURN(LOG_AMF_APP, header.extended_protocol_discriminator);
      }

      break;

    default:
      OAILOG_ERROR(
          LOG_AMF_APP, "Unknown security header type %u", security_header_type);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, 0);
  };
}

/****************************************************************************
 **                                                                        **
 ** Name:  _nas5g_message_encrypt()                                        **
 **                                                                        **
 ** Description: Encrypt plain NAS message                                 **
 **                                                                        **
 ** Inputs   src:   Pointer to the decrypted data buffer                   **
 **          security_header_type:    The security header type             **
 **          code:    The message authentication code                      **
 **          seq:   The sequence number                                    **
 **          direction: The sequence number                                **
 **          length:  Maximal capacity of the output buffer                **
 **          Others:  None                                                 **
 **                                                                        **
 ** Outputs:   dest:    Pointer to the encrypted data buffer               **
 **            Return:  The number of bytes in the output buffer           **
 **                     if data have been successfully encrypted;          **
 **                     RETURNerror otherwise.                             **
 **            Others:  None                                               **
 **                                                                        **
 ***************************************************************************/
static int _nas5g_message_encrypt(
    unsigned char* dest, const unsigned char* src, uint8_t security_header_type,
    uint32_t code, uint8_t seq, int const direction, uint32_t length,
    amf_security_context_t* const amf_security_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (!amf_security_context) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "No security context set for encryption protection algorithm\n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, 0);
  }

  switch (security_header_type) {
    case SECURITY_HEADER_TYPE_NOT_PROTECTED:
    case SECURITY_HEADER_TYPE_SERVICE_REQUEST:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_NEW:
      OAILOG_DEBUG(
          LOG_AMF_APP,
          "No encryption of message according to security header type 0x%02x\n",
          security_header_type);
      memcpy(dest, src, length);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, length);
      break;

    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED:
    case SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED_NEW:
      switch (amf_security_context->selected_algorithms.encryption) {
        case M5G_NAS_SECURITY_ALGORITHMS_5G_EA0:
          OAILOG_DEBUG(
              LOG_AMF_APP,
              "M5G_NAS_SECURITY_ALGORITHMS_5G_EA0 dir %d ul_count.seq_num %d "
              "dl_count.seq_num %d\n",
              direction, amf_security_context->ul_count.seq_num,
              amf_security_context->dl_count.seq_num);
          memcpy(dest, src, length);
          OAILOG_FUNC_RETURN(LOG_AMF_APP, length);
          break;

        default:
          OAILOG_ERROR(
              LOG_AMF_APP, "Unsupported Cyphering protection algorithm %d\n",
              amf_security_context->selected_algorithms.encryption);
          break;
      }

      break;

    default:
      OAILOG_ERROR(
          LOG_AMF_APP, "Unknown security header type %u\n",
          security_header_type);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, 0);
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, length);
}

/*
   -----------------------------------------------------------------------------
    Functions used for integrity protection of layer 3 NAS messages
   -----------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:  _nas5g_message_get_mac()                                        **
 **                                                                        **
 ** Description: Run integrity algorithm onto cyphered or uncyphered NAS   **
 **    message encoded in the input buffer and return the compu-           **
 **    ted message authentication code                                     **
 **                                                                        **
 ** Inputs   buffer:  Pointer to the integrity protected data              **
 **       buffer                                                           **
 **    count:   Value of the uplink NAS counter                            **
 **    length:  Length of the input buffer                                 **
 **      direction                                                         **
 **    Others:  None                                                       **
 **                                                                        **
 ** Outputs:   None                                                        **
 **      Return:  The message authentication code                          **
 **    Others:  None                                                       **
 **                                                                        **
 ***************************************************************************/
static uint32_t _nas5g_message_get_mac(
    const unsigned char* const buffer, uint32_t const length,
    int const direction, amf_security_context_t* const amf_security_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);

  if (!amf_security_context) {
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "No security context set for integrity protection algorithm\n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, 0);
  }

  switch (amf_security_context->selected_algorithms.integrity) {
    case M5G_NAS_SECURITY_ALGORITHMS_128_5G_IA1: {
      uint8_t mac[4];
      nas_stream_cipher_t stream_cipher;
      uint32_t count;
      uint32_t* mac32;

      if (direction == SECU_DIRECTION_UPLINK) {
        count = 0x00000000 |
                ((amf_security_context->ul_count.overflow & 0x0000FFFF) << 8) |
                (amf_security_context->ul_count.seq_num & 0x000000FF);
      } else {
        count = 0x00000000 |
                ((amf_security_context->dl_count.overflow & 0x0000FFFF) << 8) |
                (amf_security_context->dl_count.seq_num & 0x000000FF);
      }

      OAILOG_INFO(
          LOG_AMF_APP,
          "M5G_NAS_SECURITY_ALGORITHMS_128_5G_IA1 %s count.seq_num %u count "
          "%u\n",
          (direction == SECU_DIRECTION_UPLINK) ? "UPLINK" : "DOWNLINK",
          (direction == SECU_DIRECTION_UPLINK) ?
              amf_security_context->ul_count.seq_num :
              amf_security_context->dl_count.seq_num,
          count);
      stream_cipher.key        = amf_security_context->knas_int;
      stream_cipher.key_length = AUTH_KNAS_INT_SIZE;
      stream_cipher.count      = count;
      stream_cipher.bearer     = 0x01;  // 33.401 section 8.1.1
      stream_cipher.direction  = direction;
      stream_cipher.message    = const_cast<uint8_t*>(buffer);
      /*
       *        * length in bits
       *               */
      stream_cipher.blength = length << 3;
      nas_stream_encrypt_eia1(&stream_cipher, mac);
      OAILOG_INFO(
          LOG_AMF_APP,
          "M5G_NAS_SECURITY_ALGORITHMS_128_5G_IA1 returned MAC %x.%x.%x.%x(%u) "
          "for "
          "length "
          "%" PRIu32 ", direction %d, count %d\n",
          mac[0], mac[1], mac[2], mac[3], *(reinterpret_cast<uint32_t*>(&mac)),
          length, direction, count);
      mac32 = reinterpret_cast<uint32_t*>(&mac);
      OAILOG_FUNC_RETURN(LOG_NAS, ntohl(*mac32));
    } break;

    case M5G_NAS_SECURITY_ALGORITHMS_128_5G_IA2: {
      uint8_t mac[4];
      nas_stream_cipher_t stream_cipher;
      uint32_t count;
      uint32_t* mac32;

      if (direction == SECU_DIRECTION_UPLINK) {
        count = 0x00000000 |
                ((amf_security_context->ul_count.overflow & 0x0000FFFF) << 8) |
                (amf_security_context->ul_count.seq_num & 0x000000FF);
      } else {
        count = 0x00000000 |
                ((amf_security_context->dl_count.overflow & 0x0000FFFF) << 8) |
                (amf_security_context->dl_count.seq_num & 0x000000FF);
      }

      OAILOG_DEBUG(
          LOG_AMF_APP,
          "M5G_NAS_SECURITY_ALGORITHMS_5G_IA2 dir %s count.seq_num %u count "
          "%u\n",
          (direction == SECU_DIRECTION_UPLINK) ? "UPLINK" : "DOWNLINK",
          (direction == SECU_DIRECTION_UPLINK) ?
              amf_security_context->ul_count.seq_num :
              amf_security_context->dl_count.seq_num,
          count);

      stream_cipher.key        = amf_security_context->knas_int;
      stream_cipher.key_length = AUTH_KNAS_INT_SIZE;
      stream_cipher.count      = count;
      stream_cipher.bearer     = 0x01;  // 33.401 section 8.1.1
      stream_cipher.direction  = direction;
      stream_cipher.message    = (uint8_t*) buffer;
      /*
       * length in bits
       */
      stream_cipher.blength = length << 3;
      nas_stream_encrypt_eia2(&stream_cipher, mac);
      OAILOG_DEBUG(
          LOG_AMF_APP,
          "M5G_NAS_SECURITY_ALGORITHMS_5G_IA2 returned MAC %x.%x.%x.%x(%u) for "
          "length "
          "%" PRIu32 "direction %d, count %d\n",
          mac[0], mac[1], mac[2], mac[3], *((uint32_t*) &mac), length,
          direction, count);
      mac32 = (uint32_t*) &mac;
      OAILOG_FUNC_RETURN(LOG_AMF_APP, ntohl(*mac32));
    } break;
    case M5G_NAS_SECURITY_ALGORITHMS_5G_IA0:
      OAILOG_DEBUG(
          LOG_AMF_APP,
          "M5G_NAS_SECURITY_ALGORITHMS_5G_IA0 dir %s count.seq_num %u\n",
          (direction == SECU_DIRECTION_UPLINK) ? "UPLINK" : "DOWNLINK",
          (direction == SECU_DIRECTION_UPLINK) ?
              amf_security_context->ul_count.seq_num :
              amf_security_context->dl_count.seq_num);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, 0);
      break;

    default:
      OAILOG_ERROR(
          LOG_AMF_APP, "Unsupported integrity protection algorithm %d\n",
          amf_security_context->selected_algorithms.integrity);
      break;
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, 0);
}

std::string get_message_type_str(uint8_t type) {
  std::string msgtype_str;
  switch (type) {
    case REG_REQUEST:
      msgtype_str = "REG_REQUEST";
      break;
    case REG_ACCEPT:
      msgtype_str = "REG_ACCEPT";
      break;
    case REG_COMPLETE:
      msgtype_str = "REG_COMPLETE";
      break;
    case M5G_SERVICE_REQUEST:
      msgtype_str = "M5G_SERVICE_REQUEST";
      break;
    case M5G_SERVICE_ACCEPT:
      msgtype_str = "M5G_SERVICE_ACCEPT";
      break;
    case M5G_SERVICE_REJECT:
      msgtype_str = "M5G_SERVICE_REJECT";
      break;
    case M5G_IDENTITY_REQUEST:
      msgtype_str = "M5G_IDENTITY_REQUEST";
      break;
    case M5G_IDENTITY_RESPONSE:
      msgtype_str = "M5G_IDENTITY_RESPONSE";
      break;
    case AUTH_REQUEST:
      msgtype_str = "AUTH_REQUEST";
      break;
    case AUTH_RESPONSE:
      msgtype_str = "AUTH_RESPONSE";
      break;
    case AUTH_FAILURE:
      msgtype_str = "AUTH_FAILURE";
      break;
    case SEC_MODE_COMMAND:
      msgtype_str = "SEC_MODE_COMMAND";
      break;
    case SEC_MODE_COMPLETE:
      msgtype_str = "SEC_MODE_COMPLETE";
      break;
    case DE_REG_REQUEST_UE_ORIGIN:
      msgtype_str = "DE_REG_REQUEST_UE_ORIGIN";
      break;
    case DE_REG_ACCEPT_UE_ORIGIN:
      msgtype_str = "DE_REG_ACCEPT_UE_ORIGIN";
      break;
    case ULNASTRANSPORT:
      msgtype_str = "ULNASTRANSPORT";
      break;
    case DLNASTRANSPORT:
      msgtype_str = "DLNASTRANSPORT";
      break;
    default:
      msgtype_str = "UNKNOWN";
      break;
  }
  return msgtype_str;
}

std::string uint8_to_hex_string(const uint8_t* v, const size_t s) {
  std::stringstream ss;

  if (!v || (0 == s)) return ss.str();

  ss << std::hex << std::setfill('0');
  for (unsigned int i = 0; i < s; i++) {
    ss << std::hex << std::setw(2) << static_cast<int>(v[i]) << " ";
    if ((i + 1) % 8 == 0) ss << " ";
    if ((i + 1) % 16 == 0) ss << "\n";
  }
  return ss.str();
}
}  // namespace magma5g
