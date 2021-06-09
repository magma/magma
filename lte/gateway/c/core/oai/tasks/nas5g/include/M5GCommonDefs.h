#pragma once
#include <arpa/inet.h>
#include "common_defs.h"
#include "magma_logging.h"
#include "glogwrapper/glog_logging.h"

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
    int iLEN = 0;                                                              \
    if (bUFFER != NULL) {                                                      \
      while (iLEN < ((int) lEN)) {                                             \
        MLOG(MDEBUG) << " 0x" << std::hex << int(*(bUFFER + iLEN));            \
        iLEN++;                                                                \
      }                                                                        \
    }                                                                          \
  }
