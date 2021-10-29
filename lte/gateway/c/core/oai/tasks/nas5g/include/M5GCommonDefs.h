#pragma once
#include <arpa/inet.h>
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "orc8r/gateway/c/common/logging/magma_logging.h"
#include "lte/gateway/c/core/oai/common/glogwrapper/glog_logging.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GNasEnums.h"

// AMF_TEST scheme output  nibbles needs to be reversed
#define REV_NIBBLE(bUFFER, sIZE)                                               \
  for (int i = 0; i < sIZE; i++) {                                             \
    bUFFER[i] = (((bUFFER[i] & 0xf0) >> 4) | ((bUFFER[i] & 0x0f) << 4));       \
  }

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
      while (iLEN < (uint32_t) lEN) {                                          \
        MLOG(MDEBUG) << " 0x" << std::hex << int(*(bUFFER + iLEN));            \
        iLEN++;                                                                \
      }                                                                        \
    }                                                                          \
  }
