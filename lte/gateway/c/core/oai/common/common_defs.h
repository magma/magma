/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

/*! \file common_defs.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_COMMON_DEFS_SEEN
#define FILE_COMMON_DEFS_SEEN

#include <string.h>  // memcpy
#include <arpa/inet.h>
//------------------------------------------------------------------------------
#define STOLEN_REF
#define CLONE_REF

#define OFFSET_OF(TyPe, MeMBeR) ((size_t) & ((TyPe*) 0)->MeMBeR)
// https://stackoverflow.com/questions/4415524/common-array-length-macro-for-c
#define COUNT_OF(x)                                                            \
  ((sizeof(x) / sizeof(0 [x])) / ((size_t)(!(sizeof(x) % sizeof(0 [x])))))

#define PARENT_STRUCT(cOnTaiNeD, TyPe, MeMBeR)                                 \
  ({                                                                           \
    const typeof(((TyPe*) 0)->MeMBeR)* __MemBeR_ptr = (cOnTaiNeD);             \
    (TyPe*) ((char*) __MemBeR_ptr - OFFSET_OF(TyPe, MeMBeR));                  \
  })

#define OAI_MAX(a, b)                                                          \
  ({                                                                           \
    __typeof__(a) _a = (a);                                                    \
    __typeof__(b) _b = (b);                                                    \
    _a > _b ? _a : _b;                                                         \
  })

#define OAI_MIN(a, b)                                                          \
  ({                                                                           \
    __typeof__(a) _a = (a);                                                    \
    __typeof__(b) _b = (b);                                                    \
    _a < _b ? _a : _b;                                                         \
  })
//------------------------------------------------------------------------------

#define INTERNAL_FLAG_NULL (0)
#define INTERNAL_FLAG_X2_HANDOVER (1 << 0)
#define INTERNAL_FLAG_SKIP_RESPONSE (1 << 2)
#define INTERNAL_FLAG_TRIGGERED_ACK (1 << 3)
#define INTERNAL_LATE_RESPONS_IND (1 << 4)

//------------------------------------------------------------------------------

typedef enum {
  /* Fatal errors - received message should not be processed */
  TLV_MAC_MISMATCH                  = -14,
  TLV_BUFFER_NULL                   = -13,
  TLV_BUFFER_TOO_SHORT              = -12,
  TLV_PROTOCOL_NOT_SUPPORTED        = -11,
  TLV_WRONG_MESSAGE_TYPE            = -10,
  TLV_OCTET_STRING_TOO_LONG_FOR_IEI = -9,

  TLV_VALUE_DOESNT_MATCH          = -4,
  TLV_MANDATORY_FIELD_NOT_PRESENT = -3,
  TLV_UNEXPECTED_IEI              = -2,

  TLV_ERROR_OK = 0,
  /* Defines error code limit below which received message should be discarded
   * because it cannot be further processed */
  TLV_FATAL_ERROR = TLV_VALUE_DOESNT_MATCH

} tlv_error_code_e;

/* Intended to be the primary mechanism to gracefully handle errors across
 * API boundaries.
 * Most functions which can produce a recoverable error should be designed to
 * return this.
 * Inspired by Google's absl::Status
 */
typedef enum {
  RETURNerror = -1,  // An internal invariant is broken
  RETURNok    = 0,   // Ok
} status_code_e;

/* This enum should match with the ModeMapItem_FederatedMode enum
 * defined in mconfigs.proto
 */
typedef enum {
  SPGW_SUBSCRIBER  = 0,  // default mode is HSS + spgw_task
  LOCAL_SUBSCRIBER = 1,  // will use subscriberDb + spgw_task
  S8_SUBSCRIBER    = 2,  // will use federated HSS + s8_task
} mode_map;

//------------------------------------------------------------------------------
#define DECODE_U8(bUFFER, vALUE, sIZE)                                         \
  vALUE = *(uint8_t*) (bUFFER);                                                \
  sIZE += sizeof(uint8_t)

#define DECODE_U16(bUFFER, vALUE, sIZE)                                        \
  memcpy((unsigned char*) &vALUE, bUFFER, sizeof(uint16_t));                   \
  vALUE = ntohs(vALUE);                                                        \
  sIZE += sizeof(uint16_t)

#define DECODE_U24(bUFFER, vALUE, sIZE)                                        \
  memcpy((unsigned char*) &vALUE, bUFFER, sizeof(uint32_t));                   \
  vALUE = ntohl(vALUE) >> 8;                                                   \
  sIZE += sizeof(uint8_t) + sizeof(uint16_t)

#define DECODE_U32(bUFFER, vALUE, sIZE)                                        \
  memcpy((unsigned char*) &vALUE, bUFFER, sizeof(uint32_t));                   \
  vALUE = ntohl(vALUE);                                                        \
  sIZE += sizeof(uint32_t)

#if (BYTE_ORDER == LITTLE_ENDIAN)
#define DECODE_LENGTH_U16(bUFFER, vALUE, sIZE)                                 \
  vALUE = ((*(bUFFER)) << 8) | (*((bUFFER) + 1));                              \
  sIZE += sizeof(uint16_t)
#else
#define DECODE_LENGTH_U16(bUFFER, vALUE, sIZE)                                 \
  vALUE = (*(bUFFER)) | (*((bUFFER) + 1) << 8);                                \
  sIZE += sizeof(uint16_t)
#endif

#define ENCODE_U8(buffer, value, size)                                         \
  *(uint8_t*) (buffer) = value;                                                \
  size += sizeof(uint8_t)

#define ENCODE_U16(buffer, value, size)                                        \
  do {                                                                         \
    uint16_t n_value = htons(value);                                           \
    memcpy(buffer, (unsigned char*) &n_value, sizeof(uint16_t));               \
    size += sizeof(uint16_t);                                                  \
  } while (0)

#define ENCODE_U24(buffer, value, size)                                        \
  do {                                                                         \
    uint32_t n_value = htonl(value);                                           \
    memcpy(buffer, (unsigned char*) &n_value, sizeof(uint32_t));               \
    size += sizeof(uint8_t) + sizeof(uint16_t);                                \
  } while (0)

#define ENCODE_U32(buffer, value, size)                                        \
  do {                                                                         \
    uint32_t n_value = htonl(value);                                           \
    memcpy(buffer, (unsigned char*) &n_value, sizeof(uint32_t));               \
    size += sizeof(uint32_t);                                                  \
  } while (0)

#define IPV4_STR_ADDR_TO_INADDR(AdDr_StR, InAdDr, MeSsAgE)                     \
  do {                                                                         \
    if (inet_aton(AdDr_StR, &InAdDr) <= 0) {                                   \
      Fatal(MeSsAgE);                                                          \
    }                                                                          \
  } while (0)

#define IPV6_STR_ADDR_TO_INADDR(AdDr_StR, InAdDr, MeSsAgE)                     \
  do {                                                                         \
    if (inet_pton(AF_INET6, AdDr_StR, &InAdDr) <= 0) {                         \
      Fatal(MeSsAgE);                                                          \
    }                                                                          \
  } while (0)

#define NIPADDR(addr)                                                          \
  (uint8_t)(addr & 0x000000FF), (uint8_t)((addr & 0x0000FF00) >> 8),           \
      (uint8_t)((addr & 0x00FF0000) >> 16),                                    \
      (uint8_t)((addr & 0xFF000000) >> 24)

#define HIPADDR(addr)                                                          \
  (uint8_t)((addr & 0xFF000000) >> 24), (uint8_t)((addr & 0x00FF0000) >> 16),  \
      (uint8_t)((addr & 0x0000FF00) >> 8), (uint8_t)(addr & 0x000000FF)

#define NIP6ADDR(addr)                                                         \
  ntohs((addr)->s6_addr16[0]), ntohs((addr)->s6_addr16[1]),                    \
      ntohs((addr)->s6_addr16[2]), ntohs((addr)->s6_addr16[3]),                \
      ntohs((addr)->s6_addr16[4]), ntohs((addr)->s6_addr16[5]),                \
      ntohs((addr)->s6_addr16[6]), ntohs((addr)->s6_addr16[7])

#define IN6_ARE_ADDR_MASKED_EQUAL(a, b, m)                                     \
  (((((__const uint32_t*) (a))[0] & (((__const uint32_t*) (m))[0])) ==         \
    (((__const uint32_t*) (b))[0] & (((__const uint32_t*) (m))[0]))) &&        \
   ((((__const uint32_t*) (a))[1] & (((__const uint32_t*) (m))[1])) ==         \
    (((__const uint32_t*) (b))[1] & (((__const uint32_t*) (m))[1]))) &&        \
   ((((__const uint32_t*) (a))[2] & (((__const uint32_t*) (m))[2])) ==         \
    (((__const uint32_t*) (b))[2] & (((__const uint32_t*) (m))[2]))) &&        \
   ((((__const uint32_t*) (a))[3] & (((__const uint32_t*) (m))[3])) ==         \
    (((__const uint32_t*) (b))[3] & (((__const uint32_t*) (m))[3]))))

#define EBI_TO_INDEX(eBi) (eBi - 5)
#define INDEX_TO_EBI(iNdEx) (iNdEx + 5)

#ifndef UNUSED
#define UNUSED(x) (void) (x)
#endif

#endif /* FILE_COMMON_DEFS_SEEN */
