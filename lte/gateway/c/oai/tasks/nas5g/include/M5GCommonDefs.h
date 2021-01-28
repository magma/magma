#pragma once
#include <arpa/inet.h>
#include <stdint.h>
#include <iostream>
#include "M5gNasMessage.h"
#include "magma_logging.h"
#include "glogwrapper/glog_logging.h"

typedef enum {
  TLV_MAC_MISMATCH                  = -14,
  TLV_BUFFER_NULL                   = -13,
  TLV_BUFFER_TOO_SHORT              = -12,
  TLV_PROTOCOL_NOT_SUPPORTED        = -11,
  TLV_WRONG_MESSAGE_TYPE            = -10,
  TLV_OCTET_STRING_TOO_LONG_FOR_IEI = -9,
  TLV_VALUE_DOESNT_MATCH            = -4,
  TLV_MANDATORY_FIELD_NOT_PRESENT   = -3,
  TLV_UNEXPECTED_IEI                = -2,
  RETURN_ERROR                      = -1,
  RETURN_NO_ERROR                   = 0,
  TLV_NO_ERROR                      = RETURN_NO_ERROR,
  TLV_FATAL_ERROR                   = TLV_VALUE_DOESNT_MATCH
} TLVErrorCode;

#define DECODE_U8(bUFFER, vALUE, sIZE)                                         \
  vALUE = *(uint8_t*) (bUFFER);                                                \
  sIZE += sizeof(uint8_t)

#define DECODE_U16(bUFFER, vALUE, sIZE)                                        \
  vALUE = ntohs(*(uint16_t*) (bUFFER));                                        \
  sIZE += sizeof(uint16_t)

#define DECODE_U24(bUFFER, vALUE, sIZE)                                        \
  vALUE = ntohl(*(uint32_t*) (bUFFER)) >> 8;                                   \
  sIZE += sizeof(uint8_t) + sizeof(uint16_t)

#define DECODE_U32(bUFFER, vALUE, sIZE)                                        \
  vALUE = ntohl(*(uint32_t*) (bUFFER));                                        \
  sIZE += sizeof(uint32_t)

#define ENCODE_U8(buffer, value, size)                                         \
  *(uint8_t*) (buffer) = value;                                                \
  size += sizeof(uint8_t)

#define ENCODE_U16(buffer, value, size)                                        \
  *(uint16_t*) (buffer) = htons(value);                                        \
  size += sizeof(uint16_t)

#define ENCODE_U24(buffer, value, size)                                        \
  *(uint32_t*) (buffer) = htonl(value);                                        \
  size += sizeof(uint8_t) + sizeof(uint16_t)

#define ENCODE_U32(buffer, value, size)                                        \
  *(uint32_t*) (buffer) = htonl(value);                                        \
  size += sizeof(uint32_t)

#define CHECK_IEI_DECODER(iEI, bUFFER)                                         \
  if (iEI != bUFFER) {                                                         \
    MLOG(MERROR) << "Error: " << std::dec << TLV_UNEXPECTED_IEI;               \
    return TLV_UNEXPECTED_IEI;                                                 \
  }

#define CHECK_IEI_ENCODER(iEI, TYPEVALUE)                                      \
  if (iEI != TYPEVALUE) {                                                      \
    MLOG(MERROR) << "Error: " << std::dec << TLV_UNEXPECTED_IEI;               \
    return TLV_UNEXPECTED_IEI;                                                 \
  }

#define CHECK_LENGTH_DECODER(bUFFERlENGTH, lENGTH)                             \
  if ((uint32_t) bUFFERlENGTH < (uint32_t) lENGTH) {                           \
    MLOG(MERROR) << "Error: " << std::dec << TLV_BUFFER_TOO_SHORT;             \
    return TLV_BUFFER_TOO_SHORT;                                               \
  }

#define CHECK_PDU_POINTER_AND_LENGTH_ENCODER(bUFFER, mINIMUMlENGTH, lENGTH)    \
  if (bUFFER == NULL) {                                                        \
    MLOG(MERROR) << "Error: " << std::dec << TLV_BUFFER_NULL;                  \
    return TLV_BUFFER_NULL;                                                    \
  }                                                                            \
  if ((uint32_t) lENGTH < (uint32_t) mINIMUMlENGTH) {                          \
    MLOG(MERROR) << "Error: " << std::dec << TLV_BUFFER_TOO_SHORT;             \
    return TLV_BUFFER_TOO_SHORT;                                               \
  }

#define CHECK_PDU_POINTER_AND_LENGTH_DECODER(bUFFER, mINIMUMlENGTH, lENGTH)    \
  if (bUFFER == NULL) {                                                        \
    MLOG(MERROR) << "Error: " << std::dec << TLV_BUFFER_NULL;                  \
    return TLV_BUFFER_NULL;                                                    \
  }                                                                            \
  if ((uint32_t) lENGTH < (uint32_t) mINIMUMlENGTH) {                          \
    MLOG(MERROR) << "Error: " << std::dec << TLV_BUFFER_TOO_SHORT;             \
    return TLV_BUFFER_TOO_SHORT;                                               \
  }

#define IES_ENCODE_U8(buffer, encoded, value)                                  \
  ENCODE_U8(buffer + encoded, value, encoded)

#define IES_ENCODE_U16(buffer, encoded, value)                                 \
  ENCODE_U16(buffer + encoded, value, encoded)

#define IES_ENCODE_U24(buffer, encoded, value)                                 \
  ENCODE_U24(buffer + encoded, value, encoded)

#define IES_ENCODE_U32(buffer, encoded, value)                                 \
  ENCODE_U32(buffer + encoded, value, encoded)

#define IES_DECODE_U8(bUFFER, dECODED, vALUE)                                  \
  DECODE_U8(bUFFER + dECODED, vALUE, dECODED)

#define IES_DECODE_U16(bUFFER, dECODED, vALUE)                                 \
  DECODE_U16(bUFFER + dECODED, vALUE, dECODED)

#define IES_DECODE_U24(bUFFER, dECODED, vALUE)                                 \
  DECODE_U24(bUFFER + dECODED, vALUE, dECODED)

#define IES_DECODE_U32(bUFFER, dECODED, vALUE)                                 \
  DECODE_U32(bUFFER + dECODED, vALUE, dECODED)

#define BUFFER_PRINT_LOG(bUFFER, lEN)                                          \
  {                                                                            \
    uint32_t iLEN = 0;                                                         \
    if (bUFFER != NULL) {                                                      \
      while (iLEN < lEN) {                                                     \
        MLOG(MDEBUG) << " 0x" << hex << int(*(bUFFER + iLEN));                 \
        iLEN++;                                                                \
      }                                                                        \
    }                                                                          \
    MLOG(MDEBUG) << endl;                                                      \
  }
