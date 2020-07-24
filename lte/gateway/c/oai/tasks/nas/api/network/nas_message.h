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

Source      nas_message.h

Version     0.1

Date        2012/26/09

Product     NAS stack

Subsystem   Application Programming Interface

Author      Frederic Maurel

Description Defines the layer 3 messages supported by the NAS sublayer
        protocol and functions used to encode and decode

*****************************************************************************/
#ifndef FILE_NAS_MESSAGE_SEEN
#define FILE_NAS_MESSAGE_SEEN

#include <linux/byteorder/little_endian.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "nas/commonDef.h"
#include "emm_msg.h"
#include "emm_data.h"
#include "esm_msg.h"
#include "3gpp_24.007.h"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

#define NAS_MESSAGE_SECURITY_HEADER_SIZE 6
#define NAS_MESSAGE_SERVICE_REQUEST_SECURITY_HEADER_SIZE 4
/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/* Structure of security protected header */
typedef struct nas_message_security_header_s {
#ifdef __LITTLE_ENDIAN_BITFIELD
  eps_protocol_discriminator_t protocol_discriminator : 4;
  uint8_t security_header_type : 4;
#endif
#ifdef __BIG_ENDIAN_BITFIELD
  uint8_t security_header_type : 4;
  uint8_t protocol_discriminator : 4;
#endif
  uint32_t message_authentication_code;
  uint8_t sequence_number;
} nas_message_security_header_t;

/* Structure of plain NAS message */
typedef union {
  EMM_msg emm; /* EPS Mobility Management messages */
  ESM_msg esm; /* EPS Session Management messages  */
} nas_message_plain_t;

/* Structure of security protected NAS message */
typedef struct nas_message_security_protected_s {
  nas_message_security_header_t header;
  nas_message_plain_t plain;
} nas_message_security_protected_t;

/*
 * Structure of a layer 3 NAS message
 */
typedef union {
  nas_message_security_header_t header;
  nas_message_security_protected_t security_protected;
  nas_message_plain_t plain;
} nas_message_t;

typedef struct nas_message_decode_status_s {
  uint8_t integrity_protected_message : 1;
  uint8_t ciphered_message : 1;
  uint8_t mac_matched : 1;
  uint8_t security_context_available : 1;
  int emm_cause;
} nas_message_decode_status_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/
int nas_message_header_decode(
    const unsigned char* const buffer,
    nas_message_security_header_t* const header, const size_t length,
    nas_message_decode_status_t* const status, bool* const is_sr);

int nas_message_encrypt(
    const unsigned char* inbuf, unsigned char* outbuf,
    const nas_message_security_header_t* header, size_t length, void* security);

int nas_message_decrypt(
    const unsigned char* const inbuf, unsigned char* const outbuf,
    nas_message_security_header_t* header, size_t length, void* security,
    nas_message_decode_status_t* status);

int nas_message_decode(
    const unsigned char* const buffer, nas_message_t* msg, size_t length,
    void* security, nas_message_decode_status_t* status);

int nas_message_encode(
    unsigned char* buffer, const nas_message_t* const msg, size_t length,
    void* security);

#endif /* FILE_NAS_MESSAGE_SEEN*/
