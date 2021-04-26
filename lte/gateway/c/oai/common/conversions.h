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

/*! \file conversions.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_CONVERSIONS_SEEN
#define FILE_CONVERSIONS_SEEN
#include <stdio.h>
#include <endian.h>
#include <netinet/in.h>
#include <stdbool.h>
#include <stdint.h>
#include <math.h>

#include "common_types.h"
#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "EpsQualityOfService.h"

/* Endianness conversions for 16 and 32 bits integers from host to network order
 */
#if (BYTE_ORDER == LITTLE_ENDIAN)
#define hton_int32(x)                                                          \
  (((x & 0x000000FF) << 24) | ((x & 0x0000FF00) << 8) |                        \
   ((x & 0x00FF0000) >> 8) | ((x & 0xFF000000) >> 24))

#define hton_int16(x)                                                          \
    (((x & 0x00FF) << 8) | ((x & 0xFF00) >> 8)

#define ntoh_int32_buf(bUF)                                                    \
  ((*(bUF)) << 24) | ((*((bUF) + 1)) << 16) | ((*((bUF) + 2)) << 8) |          \
      (*((bUF) + 3))
#else
#define hton_int32(x) (x)
#define hton_int16(x) (x)
#endif

/* in: X (struct inaddr) in network byte order (big endian)
 * out: bUFF with most significant IP byte in buf[0]
 */
#define IN_ADDR_TO_BUFFER(X, bUFF)                                             \
  do {                                                                         \
    ((uint8_t*) (bUFF))[0] = ((uint8_t*) &((X).s_addr))[0];                    \
    ((uint8_t*) (bUFF))[1] = ((uint8_t*) &((X).s_addr))[1];                    \
    ((uint8_t*) (bUFF))[2] = ((uint8_t*) &((X).s_addr))[2];                    \
    ((uint8_t*) (bUFF))[3] = ((uint8_t*) &((X).s_addr))[3];                    \
  } while (0)

#define BUFFER_TO_IN_ADDR(bUFF, X)                                             \
  do {                                                                         \
    (X).s_addr = (((uint8_t*) (bUFF))[0]) | (((uint8_t*) (bUFF))[1] << 8) |    \
                 (((uint8_t*) (bUFF))[2] << 16) |                              \
                 (((uint8_t*) (bUFF))[3] << 24);                               \
  } while (0)

#define IN6_ADDR_TO_BUFFER(X, bUFF)                                            \
  do {                                                                         \
    ((uint8_t*) (bUFF))[0]  = (X).s6_addr[0];                                  \
    ((uint8_t*) (bUFF))[1]  = (X).s6_addr[1];                                  \
    ((uint8_t*) (bUFF))[2]  = (X).s6_addr[2];                                  \
    ((uint8_t*) (bUFF))[3]  = (X).s6_addr[3];                                  \
    ((uint8_t*) (bUFF))[4]  = (X).s6_addr[4];                                  \
    ((uint8_t*) (bUFF))[5]  = (X).s6_addr[5];                                  \
    ((uint8_t*) (bUFF))[6]  = (X).s6_addr[6];                                  \
    ((uint8_t*) (bUFF))[7]  = (X).s6_addr[7];                                  \
    ((uint8_t*) (bUFF))[8]  = (X).s6_addr[8];                                  \
    ((uint8_t*) (bUFF))[9]  = (X).s6_addr[9];                                  \
    ((uint8_t*) (bUFF))[10] = (X).s6_addr[10];                                 \
    ((uint8_t*) (bUFF))[11] = (X).s6_addr[11];                                 \
    ((uint8_t*) (bUFF))[12] = (X).s6_addr[12];                                 \
    ((uint8_t*) (bUFF))[13] = (X).s6_addr[13];                                 \
    ((uint8_t*) (bUFF))[14] = (X).s6_addr[14];                                 \
    ((uint8_t*) (bUFF))[15] = (X).s6_addr[15];                                 \
  } while (0)

#define BUFFER_TO_INT8(buf, x) (x = ((buf)[0]))

#define INT8_TO_BUFFER(x, buf) ((buf)[0] = (x))

/* Convert an integer on 16 bits to the given bUFFER */
#define INT16_TO_BUFFER(x, buf)                                                \
  do {                                                                         \
    (buf)[0] = (x) >> 8;                                                       \
    (buf)[1] = (x);                                                            \
  } while (0)

/* Convert an integer on 24 bits to the given bUFFER */
#define INT24_TO_BUFFER(x, buf)                                                \
  do {                                                                         \
    (buf)[0] = (x) >> 16;                                                      \
    (buf)[1] = (x) >> 8;                                                       \
    (buf)[2] = (x);                                                            \
  } while (0)

/* Convert an array of char containing vALUE to x */
#define BUFFER_TO_INT16(buf, x)                                                \
  do {                                                                         \
    x = ((buf)[0] << 8) | ((buf)[1]);                                          \
  } while (0)

#define BUFFER_TO_INT24(buf, x)                                                \
  do {                                                                         \
    x = (int32_t)(                                                             \
        ((uint32_t)((buf)[0]) << 16) | ((uint32_t)((buf)[1]) << 8) |           \
        ((uint32_t)((buf)[2])));                                               \
  } while (0)

/* Convert an integer on 32 bits to the given bUFFER */
#define INT32_TO_BUFFER(x, buf)                                                \
  do {                                                                         \
    (buf)[0] = (x) >> 24;                                                      \
    (buf)[1] = (x) >> 16;                                                      \
    (buf)[2] = (x) >> 8;                                                       \
    (buf)[3] = (x);                                                            \
  } while (0)

/* Convert an array of char containing vALUE to x */
#define BUFFER_TO_INT32(buf, x)                                                \
  do {                                                                         \
    x = (int32_t)(                                                             \
        ((uint32_t)((buf)[0]) << 24) | ((uint32_t)((buf)[1]) << 16) |          \
        ((uint32_t)((buf)[2]) << 8) | ((uint32_t)((buf)[3])));                 \
  } while (0)

/* Convert an integer on 32 bits to an octet string from aSN1c tool */
#define INT32_TO_OCTET_STRING(x, aSN)                                          \
  do {                                                                         \
    (aSN)->buf = calloc(4, sizeof(uint8_t));                                   \
    INT32_TO_BUFFER(x, ((aSN)->buf));                                          \
    (aSN)->size = 4;                                                           \
  } while (0)

#define INT32_TO_BIT_STRING(x, aSN)                                            \
  do {                                                                         \
    INT32_TO_OCTET_STRING(x, aSN);                                             \
    (aSN)->bits_unused = 0;                                                    \
  } while (0)

#define UE_ID_INDEX_TO_BIT_STRING(x, aSN)                                      \
  do {                                                                         \
    INT16_TO_OCTET_STRING(x << 6, aSN);                                        \
    (aSN)->bits_unused = 6;                                                    \
  } while (0)

#define AMF_POINTER_TO_BIT_STRING(x, aSN)                                      \
  do {                                                                         \
    INT8_TO_OCTET_STRING(x << 2, aSN);                                         \
    (aSN)->bits_unused = 2;                                                    \
  } while (0)

#define INT24_TO_OCTET_STRING(x, aSN)                                          \
  do {                                                                         \
    (aSN)->buf  = calloc(3, sizeof(uint8_t));                                  \
    (aSN)->size = 3;                                                           \
    INT24_TO_BUFFER(x, (aSN)->buf);                                            \
  } while (0)

#define INT16_TO_OCTET_STRING(x, aSN)                                          \
  do {                                                                         \
    (aSN)->buf  = calloc(2, sizeof(uint8_t));                                  \
    (aSN)->size = 2;                                                           \
    INT16_TO_BUFFER(x, (aSN)->buf);                                            \
  } while (0)

#define INT8_TO_OCTET_STRING(x, aSN)                                           \
  do {                                                                         \
    (aSN)->buf  = calloc(1, sizeof(uint8_t));                                  \
    (aSN)->size = 1;                                                           \
    INT8_TO_BUFFER(x, (aSN)->buf);                                             \
  } while (0)

#define MME_CODE_TO_OCTET_STRING INT8_TO_OCTET_STRING
#define M_TMSI_TO_OCTET_STRING INT32_TO_OCTET_STRING
#define MME_GID_TO_OCTET_STRING INT16_TO_OCTET_STRING

#define OCTET_STRING_TO_INT8(aSN, x)                                           \
  do {                                                                         \
    DevCheck((aSN)->size == 1, (aSN)->size, 0, 0);                             \
    BUFFER_TO_INT8((aSN)->buf, x);                                             \
  } while (0)

#define OCTET_STRING_TO_INT16(aSN, x)                                          \
  do {                                                                         \
    DevCheck((aSN)->size == 2, (aSN)->size, 0, 0);                             \
    BUFFER_TO_INT16((aSN)->buf, x);                                            \
  } while (0)

#define OCTET_STRING_TO_INT24(aSN, x)                                          \
  do {                                                                         \
    DevCheck((aSN)->size == 3, (aSN)->size, 0, 0);                             \
    BUFFER_TO_INT24((aSN)->buf, x);                                            \
  } while (0)

#define OCTET_STRING_TO_INT32(aSN, x)                                          \
  do {                                                                         \
    DevCheck((aSN)->size == 4, (aSN)->size, 0, 0);                             \
    BUFFER_TO_INT32((aSN)->buf, x);                                            \
  } while (0)

#define BIT_STRING_TO_INT32(aSN, x)                                            \
  do {                                                                         \
    DevCheck((aSN)->bits_unused == 0, (aSN)->bits_unused, 0, 0);               \
    OCTET_STRING_TO_INT32(aSN, x);                                             \
  } while (0)

#define BIT_STRING_TO_INT16(aSN, x)                                            \
  do {                                                                         \
    DevCheck((aSN)->bits_unused == 0, (aSN)->bits_unused, 0, 0);               \
    OCTET_STRING_TO_INT16(aSN, x);                                             \
  } while (0)

#define BIT_STRING_TO_CELL_IDENTITY(aSN, vALUE)                                \
  do {                                                                         \
    DevCheck((aSN)->bits_unused == 4, (aSN)->bits_unused, 4, 0);               \
    vALUE.enb_id =                                                             \
        ((aSN)->buf[0] << 12) | ((aSN)->buf[1] << 4) | ((aSN)->buf[2] >> 4);   \
    vALUE.cell_id = ((aSN)->buf[2] << 4) | ((aSN)->buf[3] >> 4);               \
  } while (0)

#define MCC_HUNDREDS(vALUE) ((vALUE) / 100)
/* When MNC is only composed of 2 digits, set the hundreds unit to 0xf */
#define MNC_HUNDREDS(vALUE, mNCdIGITlENGTH)                                    \
  (mNCdIGITlENGTH == 2 ? 15 : (vALUE) / 100)
#define MCC_MNC_DECIMAL(vALUE) (((vALUE) / 10) % 10)
#define MCC_MNC_DIGIT(vALUE) ((vALUE) % 10)

#define MCC_TO_BUFFER(mCC, bUFFER)                                             \
  do {                                                                         \
    DevAssert(bUFFER != NULL);                                                 \
    (bUFFER)[0] = MCC_HUNDREDS(mCC);                                           \
    (bUFFER)[1] = MCC_MNC_DECIMAL(mCC);                                        \
    (bUFFER)[2] = MCC_MNC_DIGIT(mCC);                                          \
  } while (0)

#define MCC_MNC_TO_PLMNID(mCC, mNC, mNCdIGITlENGTH, oCTETsTRING)               \
  do {                                                                         \
    (oCTETsTRING)->buf    = calloc(3, sizeof(uint8_t));                        \
    (oCTETsTRING)->buf[0] = (MCC_MNC_DECIMAL(mCC) << 4) | MCC_HUNDREDS(mCC);   \
    (oCTETsTRING)->buf[1] =                                                    \
        (MNC_HUNDREDS(mNC, mNCdIGITlENGTH) << 4) | MCC_MNC_DIGIT(mCC);         \
    (oCTETsTRING)->buf[2] = (MCC_MNC_DIGIT(mNC) << 4) | MCC_MNC_DECIMAL(mNC);  \
    (oCTETsTRING)->size   = 3;                                                 \
  } while (0)

#define MCC_MNC_TO_TBCD(mCC, mNC, mNCdIGITlENGTH, tBCDsTRING)                  \
  do {                                                                         \
    char _buf[3];                                                              \
    DevAssert((mNCdIGITlENGTH == 3) || (mNCdIGITlENGTH == 2));                 \
    _buf[0] = (MCC_MNC_DECIMAL(mCC) << 4) | MCC_HUNDREDS(mCC);                 \
    _buf[1] = (MNC_HUNDREDS(mNC, mNCdIGITlENGTH) << 4) | MCC_MNC_DIGIT(mCC);   \
    _buf[2] = (MCC_MNC_DIGIT(mNC) << 4) | MCC_MNC_DECIMAL(mNC);                \
    OCTET_STRING_fromBuf(tBCDsTRING, _buf, 3);                                 \
  } while (0)

#define TBCD_TO_MCC_MNC(tBCDsTRING, mCC, mNC, mNCdIGITlENGTH)                  \
  do {                                                                         \
    int mNC_hundred;                                                           \
    DevAssert((tBCDsTRING)->size == 3);                                        \
    mNC_hundred = (((tBCDsTRING)->buf[1] & 0xf0) >> 4);                        \
    if (mNC_hundred == 0xf) {                                                  \
      mNC_hundred    = 0;                                                      \
      mNCdIGITlENGTH = 2;                                                      \
    } else {                                                                   \
      mNCdIGITlENGTH = 3;                                                      \
    }                                                                          \
    mCC = (((((tBCDsTRING)->buf[0]) & 0xf0) >> 4) * 10) +                      \
          ((((tBCDsTRING)->buf[0]) & 0x0f) * 100) +                            \
          (((tBCDsTRING)->buf[1]) & 0x0f);                                     \
    mNC = (mNC_hundred * 100) + ((((tBCDsTRING)->buf[2]) & 0xf0) >> 4) +       \
          ((((tBCDsTRING)->buf[2]) & 0x0f) * 10);                              \
  } while (0)

#define TBCD_TO_PLMN_T(tBCDsTRING, pLMN)                                       \
  do {                                                                         \
    DevAssert((tBCDsTRING)->size == 3);                                        \
    if (((tBCDsTRING)->buf[1] & 0xf0) == 0xf0) {                               \
      (pLMN)->mcc_digit2 = (((tBCDsTRING)->buf[0] & 0xf0) >> 4);               \
      (pLMN)->mcc_digit1 = ((tBCDsTRING)->buf[0] & 0x0f);                      \
      (pLMN)->mcc_digit3 = ((tBCDsTRING)->buf[1] & 0x0f);                      \
      (pLMN)->mnc_digit3 = (((tBCDsTRING)->buf[1] & 0xf0) >> 4);               \
      (pLMN)->mnc_digit1 = ((tBCDsTRING)->buf[2] & 0x0f);                      \
      (pLMN)->mnc_digit2 = (((tBCDsTRING)->buf[2] & 0xf0) >> 4);               \
    } else {                                                                   \
      (pLMN)->mcc_digit2 = (((tBCDsTRING)->buf[0] & 0xf0) >> 4);               \
      (pLMN)->mcc_digit1 = ((tBCDsTRING)->buf[0] & 0x0f);                      \
      (pLMN)->mnc_digit1 = (((tBCDsTRING)->buf[1] & 0xf0) >> 4);               \
      (pLMN)->mcc_digit3 = ((tBCDsTRING)->buf[1] & 0x0f);                      \
      (pLMN)->mnc_digit3 = (((tBCDsTRING)->buf[2] & 0xf0) >> 4);               \
      (pLMN)->mnc_digit2 = ((tBCDsTRING)->buf[2] & 0x0f);                      \
    }                                                                          \
  } while (0)

#define PLMN_T_TO_TBCD(pLMN, tBCDsTRING, mNClENGTH)                            \
  do {                                                                         \
    tBCDsTRING[0] = (pLMN.mcc_digit2 << 4) | pLMN.mcc_digit1;                  \
    /* ambiguous (think about len 2) */                                        \
    if (mNClENGTH == 2) {                                                      \
      tBCDsTRING[1] = (0x0F << 4) | pLMN.mcc_digit3;                           \
      tBCDsTRING[2] = (pLMN.mnc_digit2 << 4) | pLMN.mnc_digit1;                \
    } else {                                                                   \
      tBCDsTRING[1] = (pLMN.mnc_digit3 << 4) | pLMN.mcc_digit3;                \
      tBCDsTRING[2] = (pLMN.mnc_digit2 << 4) | pLMN.mnc_digit1;                \
    }                                                                          \
  } while (0)

#define TBCD_TO_LAI_T(tBCDsTRING, lAI)                                         \
  do {                                                                         \
    DevAssert((tBCDsTRING)->size == 3);                                        \
    (lAI)->mccdigit2 = (((tBCDsTRING)->buf[0] & 0xf0) >> 4);                   \
    (lAI)->mccdigit1 = ((tBCDsTRING)->buf[0] & 0x0f);                          \
    (lAI)->mncdigit3 = (((tBCDsTRING)->buf[1] & 0xf0) >> 4);                   \
    (lAI)->mccdigit3 = ((tBCDsTRING)->buf[1] & 0x0f);                          \
    (lAI)->mncdigit2 = (((tBCDsTRING)->buf[2] & 0xf0) >> 4);                   \
    (lAI)->mncdigit1 = ((tBCDsTRING)->buf[2] & 0x0f);                          \
  } while (0)

#define LAI_T_TO_TBCD(lAI, tBCDsTRING, mNClENGTH)                              \
  do {                                                                         \
    tBCDsTRING[0] = (lAI.mccdigit2 << 4) | lAI.mccdigit1;                      \
    /* ambiguous (think about len 2) */                                        \
    if (mNClENGTH == 2) {                                                      \
      tBCDsTRING[1] = (0x0F << 4) | lAI.mccdigit3;                             \
      tBCDsTRING[2] = (lAI.mncdigit2 << 4) | lAI.mncdigit1;                    \
    } else {                                                                   \
      tBCDsTRING[1] = (lAI.mncdigit3 << 4) | lAI.mccdigit3;                    \
      tBCDsTRING[2] = (lAI.mncdigit2 << 4) | lAI.mncdigit1;                    \
    }                                                                          \
  } while (0)

#define PLMN_T_TO_MCC_MNC(pLMN, mCC, mNC, mNCdIGITlENGTH)                      \
  do {                                                                         \
    mCC = pLMN.mcc_digit1 * 100 + pLMN.mcc_digit2 * 10 + pLMN.mcc_digit3;      \
    mNCdIGITlENGTH = (pLMN.mnc_digit3 == 0xF ? 2 : 3);                         \
    mNC =                                                                      \
        (mNCdIGITlENGTH == 2 ? (pLMN.mnc_digit1 * 10 + pLMN.mnc_digit2) :      \
                               (pLMN.mnc_digit1 * 100) +                       \
                                   pLMN.mnc_digit2 * 10 + pLMN.mnc_digit3);    \
  } while (0)

#define PLMN_T_TO_PLMNID(pLMN, oCTETsTRING)                                    \
  do {                                                                         \
    uint16_t plmn_mcc = 0, plmn_mnc = 0, plmn_mnc_len = 0;                     \
    PLMN_T_TO_MCC_MNC(pLMN, plmn_mcc, plmn_mnc, plmn_mnc_len);                 \
    MCC_MNC_TO_PLMNID(plmn_mcc, plmn_mnc, plmn_mnc_len, oCTETsTRING);          \
  } while (0)
/*
 * TS 36.413 v10.9.0 section 9.2.1.37:
 * Macro eNB ID:
 * Equal to the 20 leftmost bits of the Cell
 * Identity IE contained in the E-UTRAN CGI
 * IE (see subclause 9.2.1.38) of each cell
 * served by the eNB.
 */
#define MACRO_ENB_ID_TO_BIT_STRING(mACRO, bITsTRING)                           \
  do {                                                                         \
    (bITsTRING)->buf         = calloc(3, sizeof(uint8_t));                     \
    (bITsTRING)->buf[0]      = ((mACRO) >> 12);                                \
    (bITsTRING)->buf[1]      = (mACRO) >> 4;                                   \
    (bITsTRING)->buf[2]      = ((mACRO) &0x0f) << 4;                           \
    (bITsTRING)->size        = 3;                                              \
    (bITsTRING)->bits_unused = 4;                                              \
  } while (0)
/*
 * TS 36.413 v10.9.0 section 9.2.1.38:
 * E-UTRAN CGI/Cell Identity
 * The leftmost bits of the Cell
 * Identity correspond to the eNB
 * ID (defined in subclause 9.2.1.37).
 */
#define MACRO_ENB_ID_TO_CELL_IDENTITY(mACRO, cELL_iD, bITsTRING)               \
  do {                                                                         \
    (bITsTRING)->buf         = calloc(4, sizeof(uint8_t));                     \
    (bITsTRING)->buf[0]      = ((mACRO) >> 12);                                \
    (bITsTRING)->buf[1]      = (mACRO) >> 4;                                   \
    (bITsTRING)->buf[2]      = (((mACRO) &0x0f) << 4) | ((cELL_iD) >> 4);      \
    (bITsTRING)->buf[3]      = ((cELL_iD) &0x0f) << 4;                         \
    (bITsTRING)->size        = 4;                                              \
    (bITsTRING)->bits_unused = 4;                                              \
  } while (0)

/* Used to format an uint32_t containing an ipv4 address */
#define IN_ADDR_FMT "%u.%u.%u.%u"
#define PRI_IN_ADDR(aDDRESS)                                                   \
  (uint8_t)((aDDRESS.s_addr) & 0x000000ff),                                    \
      (uint8_t)(((aDDRESS.s_addr) & 0x0000ff00) >> 8),                         \
      (uint8_t)(((aDDRESS.s_addr) & 0x00ff0000) >> 16),                        \
      (uint8_t)(((aDDRESS.s_addr) & 0xff000000) >> 24)

#define IPV4_ADDR_DISPLAY_8(aDDRESS)                                           \
  (aDDRESS)[0], (aDDRESS)[1], (aDDRESS)[2], (aDDRESS)[3]

#define TAC_TO_ASN1 INT16_TO_OCTET_STRING
#define TAC_TO_ASN1_5G INT24_TO_OCTET_STRING
#define GTP_TEID_TO_ASN1 INT32_TO_OCTET_STRING
#define OCTET_STRING_TO_TAC OCTET_STRING_TO_INT16
#define OCTET_STRING_TO_TAC_5G OCTET_STRING_TO_INT24
#define OCTET_STRING_TO_MME_CODE OCTET_STRING_TO_INT8
#define OCTET_STRING_TO_M_TMSI OCTET_STRING_TO_INT32
#define OCTET_STRING_TO_MME_GID OCTET_STRING_TO_INT16
#define OCTET_STRING_TO_CSG_ID OCTET_STRING_TO_INT27

#define OCTET_STRING_TO_AMF_CODE OCTET_STRING_TO_INT8
#define OCTET_STRING_TO_AMF_GID OCTET_STRING_TO_INT16

/* Convert the IMSI contained by a char string NULL terminated to uint64_t */
#define IMSI_STRING_TO_IMSI64(sTRING, iMSI64_pTr)                              \
  sscanf(sTRING, IMSI_64_FMT, iMSI64_pTr)

/* Convert the IMEI contained by a char string NULL terminated to uint64_t */
#define IMEI_STRING_TO_IMEI64(sTRING, iMEI64_pTr)                              \
  sscanf(sTRING, IMEI_64_FMT, iMEI64_pTr)

#define IMSI64_TO_CSFBIMSI(iMsI64_t, cSfBiMsI_t)                               \
  {                                                                            \
    if ((iMsI64_t / 100000000000000) != 0) {                                   \
      cSfBiMsI_t.digit15              = iMsI64_t / 100000000000000;            \
      iMsI64_t                        = iMsI64_t % 100000000000000;            \
      cSfBiMsI_t.digit14              = iMsI64_t / 10000000000000;             \
      iMsI64_t                        = iMsI64_t % 10000000000000;             \
      cSfBiMsI_t.digit13              = iMsI64_t / 1000000000000;              \
      iMsI64_t                        = iMsI64_t % 1000000000000;              \
      cSfBiMsI_t.digit12              = iMsI64_t / 100000000000;               \
      iMsI64_t                        = iMsI64_t % 100000000000;               \
      cSfBiMsI_t.digit11              = iMsI64_t / 10000000000;                \
      iMsI64_t                        = iMsI64_t % 10000000000;                \
      cSfBiMsI_t.digit10              = iMsI64_t / 1000000000;                 \
      iMsI64_t                        = iMsI64_t % 1000000000;                 \
      cSfBiMsI_t.digit9               = iMsI64_t / 100000000;                  \
      iMsI64_t                        = iMsI64_t % 100000000;                  \
      cSfBiMsI_t.digit8               = iMsI64_t / 10000000;                   \
      iMsI64_t                        = iMsI64_t % 10000000;                   \
      cSfBiMsI_t.digit7               = iMsI64_t / 1000000;                    \
      iMsI64_t                        = iMsI64_t % 1000000;                    \
      cSfBiMsI_t.digit6               = iMsI64_t / 100000;                     \
      iMsI64_t                        = iMsI64_t % 100000;                     \
      cSfBiMsI_t.digit5               = iMsI64_t / 10000;                      \
      iMsI64_t                        = iMsI64_t % 10000;                      \
      cSfBiMsI_t.digit4               = iMsI64_t / 1000;                       \
      iMsI64_t                        = iMsI64_t % 1000;                       \
      cSfBiMsI_t.digit3               = iMsI64_t / 100;                        \
      iMsI64_t                        = iMsI64_t % 100;                        \
      cSfBiMsI_t.digit2               = iMsI64_t / 10;                         \
      iMsI64_t                        = iMsI64_t % 10;                         \
      cSfBiMsI_t.digit1               = iMsI64_t;                              \
      cSfBiMsI_t.parity               = 1;                                     \
      cSfBiMsI_t.numOfValidImsiDigits = 15;                                    \
    } else {                                                                   \
      cSfBiMsI_t.digit14              = iMsI64_t / 10000000000000;             \
      iMsI64_t                        = iMsI64_t % 10000000000000;             \
      cSfBiMsI_t.digit13              = iMsI64_t / 1000000000000;              \
      iMsI64_t                        = iMsI64_t % 1000000000000;              \
      cSfBiMsI_t.digit12              = iMsI64_t / 100000000000;               \
      iMsI64_t                        = iMsI64_t % 100000000000;               \
      cSfBiMsI_t.digit11              = iMsI64_t / 10000000000;                \
      iMsI64_t                        = iMsI64_t % 10000000000;                \
      cSfBiMsI_t.digit10              = iMsI64_t / 1000000000;                 \
      iMsI64_t                        = iMsI64_t % 1000000000;                 \
      cSfBiMsI_t.digit9               = iMsI64_t / 100000000;                  \
      iMsI64_t                        = iMsI64_t % 100000000;                  \
      cSfBiMsI_t.digit8               = iMsI64_t / 10000000;                   \
      iMsI64_t                        = iMsI64_t % 10000000;                   \
      cSfBiMsI_t.digit7               = iMsI64_t / 1000000;                    \
      iMsI64_t                        = iMsI64_t % 1000000;                    \
      cSfBiMsI_t.digit6               = iMsI64_t / 100000;                     \
      iMsI64_t                        = iMsI64_t % 100000;                     \
      cSfBiMsI_t.digit5               = iMsI64_t / 10000;                      \
      iMsI64_t                        = iMsI64_t % 10000;                      \
      cSfBiMsI_t.digit4               = iMsI64_t / 1000;                       \
      iMsI64_t                        = iMsI64_t % 1000;                       \
      cSfBiMsI_t.digit3               = iMsI64_t / 100;                        \
      iMsI64_t                        = iMsI64_t % 100;                        \
      cSfBiMsI_t.digit2               = iMsI64_t / 10;                         \
      iMsI64_t                        = iMsI64_t % 10;                         \
      cSfBiMsI_t.digit1               = iMsI64_t;                              \
      cSfBiMsI_t.parity               = 0;                                     \
      cSfBiMsI_t.numOfValidImsiDigits = 14;                                    \
    }                                                                          \
  }

#define IMSI64_TO_STRING(iMSI64, sTRING, _imsi_len)                            \
  snprintf(                                                                    \
      sTRING, IMSI_BCD_DIGITS_MAX + 1, IMSI_64_FMT_DYN_LEN, _imsi_len, iMSI64)
imsi64_t imsi_to_imsi64(const imsi_t* const imsi);
imsi64_t amf_imsi_to_imsi64(const imsi_t* const imsi);

#define IMSI_TO_STRING(iMsI_t_PtR, iMsI_sTr, MaXlEn)                           \
  do {                                                                         \
    int l_i = 0;                                                               \
    int l_j = 0;                                                               \
    while ((l_i < IMSI_BCD8_SIZE) && (l_j < MaXlEn - 1)) {                     \
      if ((((iMsI_t_PtR)->u.value[l_i] & 0xf0) >> 4) > 9) break;               \
      snprintf(                                                                \
          ((iMsI_sTr) + l_j), (MaXlEn - l_j), "%u",                            \
          (((iMsI_t_PtR)->u.value[l_i] & 0xf0) >> 4));                         \
      l_j++;                                                                   \
      if (((iMsI_t_PtR)->u.value[l_i] & 0xf) > 9 || (l_j >= MaXlEn - 1))       \
        break;                                                                 \
      snprintf(                                                                \
          ((iMsI_sTr) + l_j), (MaXlEn - l_j), "%u",                            \
          ((iMsI_t_PtR)->u.value[l_i] & 0xf));                                 \
      l_j++;                                                                   \
      l_i++;                                                                   \
    }                                                                          \
    for (; l_j < MaXlEn; l_j++) (iMsI_sTr)[l_j] = '\0';                        \
  } while (0);

#define IMEI_TO_STRING(iMeI_t_PtR, iMeI_sTr, MaXlEn)                           \
  {                                                                            \
    int l_offset = 0;                                                          \
    int l_ret    = 0;                                                          \
    l_ret        = snprintf(                                                   \
        iMeI_sTr + l_offset, MaXlEn - l_offset, "%u%u%u%u%u%u%u%u",     \
        (iMeI_t_PtR)->u.num.tac1, (iMeI_t_PtR)->u.num.tac2,             \
        (iMeI_t_PtR)->u.num.tac3, (iMeI_t_PtR)->u.num.tac4,             \
        (iMeI_t_PtR)->u.num.tac5, (iMeI_t_PtR)->u.num.tac6,             \
        (iMeI_t_PtR)->u.num.tac7, (iMeI_t_PtR)->u.num.tac8);            \
    if (l_ret > 0) {                                                           \
      l_offset += l_ret;                                                       \
      l_ret = snprintf(                                                        \
          iMeI_sTr + l_offset, MaXlEn - l_offset, "%u%u%u%u%u%u",              \
          (iMeI_t_PtR)->u.num.snr1, (iMeI_t_PtR)->u.num.snr2,                  \
          (iMeI_t_PtR)->u.num.snr3, (iMeI_t_PtR)->u.num.snr4,                  \
          (iMeI_t_PtR)->u.num.snr5, (iMeI_t_PtR)->u.num.snr6);                 \
    }                                                                          \
    if (((iMeI_t_PtR)->u.num.parity != 0x0) && (l_ret > 0)) {                  \
      l_offset += l_ret;                                                       \
      l_ret = snprintf(                                                        \
          iMeI_sTr + l_offset, MaXlEn - l_offset, "%u",                        \
          (iMeI_t_PtR)->u.num.cdsd);                                           \
    }                                                                          \
  }

#define IMEI_MOBID_TO_IMEI64(iMeI_t_PtR, iMEI64)                               \
  {                                                                            \
    (*iMEI64) = (uint64_t)((iMeI_t_PtR)->u.num.tac1);                          \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac2));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac3));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac4));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac5));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac6));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac7));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac8));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.snr1));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.snr2));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.snr3));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.snr4));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.snr5));     \
    (*iMEI64) = (10 * (*iMEI64)) + ((uint64_t)((iMeI_t_PtR)->u.num.snr6));     \
  }

#define IMEI_MOBID_TO_IMEI_TAC64(iMeI_t_PtR, tAc_PtR)                          \
  {                                                                            \
    (*tAc_PtR) = (uint64_t)((iMeI_t_PtR)->u.num.tac1);                         \
    (*tAc_PtR) = (10 * (*tAc_PtR)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac2));   \
    (*tAc_PtR) = (10 * (*tAc_PtR)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac3));   \
    (*tAc_PtR) = (10 * (*tAc_PtR)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac4));   \
    (*tAc_PtR) = (10 * (*tAc_PtR)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac5));   \
    (*tAc_PtR) = (10 * (*tAc_PtR)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac6));   \
    (*tAc_PtR) = (10 * (*tAc_PtR)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac7));   \
    (*tAc_PtR) = (10 * (*tAc_PtR)) + ((uint64_t)((iMeI_t_PtR)->u.num.tac8));   \
  }

#define IMSI_TO_OCTET_STRING(iMsI_sTr, iMsI_len, aSN)                          \
  do {                                                                         \
    int len = 0;                                                               \
    int idx = 0;                                                               \
    if ((iMsI_len % 2) != 0) {                                                 \
      len = (iMsI_len / 2) + 1;                                                \
    } else {                                                                   \
      len = (iMsI_len / 2);                                                    \
    }                                                                          \
    (aSN)->buf = calloc(len, sizeof(uint8_t));                                 \
    for (idx = 0; idx < (len); idx++) {                                        \
      ((aSN)->buf)[idx] = (iMsI_sTr[2 * idx] & 0x0f) |                         \
                          ((iMsI_sTr[(2 * idx) + 1] << 4) & 0xf0);             \
    }                                                                          \
    if ((iMsI_len % 2) != 0) {                                                 \
      ((aSN)->buf)[idx - 1] |= 0xf0;                                           \
    }                                                                          \
    (aSN)->size = len;                                                         \
  } while (0)
#define IMEISV_TO_STRING(iMeIsV_t_PtR, iMeIsV_sTr, MaXlEn)                     \
  {                                                                            \
    int l_offset = 0;                                                          \
    int l_ret    = 0;                                                          \
    l_ret        = snprintf(                                                   \
        iMeIsV_sTr + l_offset, MaXlEn - l_offset, "%u%u%u%u%u%u%u%u",   \
        (iMeIsV_t_PtR)->u.num.tac1, (iMeIsV_t_PtR)->u.num.tac2,         \
        (iMeIsV_t_PtR)->u.num.tac3, (iMeIsV_t_PtR)->u.num.tac4,         \
        (iMeIsV_t_PtR)->u.num.tac5, (iMeIsV_t_PtR)->u.num.tac6,         \
        (iMeIsV_t_PtR)->u.num.tac7, (iMeIsV_t_PtR)->u.num.tac8);        \
    if (l_ret > 0) {                                                           \
      l_offset += l_ret;                                                       \
      l_ret = snprintf(                                                        \
          iMeIsV_sTr + l_offset, MaXlEn - l_offset, "%u%u%u%u%u%u%u%u",        \
          (iMeIsV_t_PtR)->u.num.snr1, (iMeIsV_t_PtR)->u.num.snr2,              \
          (iMeIsV_t_PtR)->u.num.snr3, (iMeIsV_t_PtR)->u.num.snr4,              \
          (iMeIsV_t_PtR)->u.num.snr5, (iMeIsV_t_PtR)->u.num.snr6,              \
          (iMeIsV_t_PtR)->u.num.svn1, (iMeIsV_t_PtR)->u.num.svn2);             \
    }                                                                          \
  }

#define IMEISV_MOBID_TO_STRING(iMeIsV_t_PtR, iMeIsV_sTr, MaXlEn)                 \
  {                                                                              \
    int l_offset = 0;                                                            \
    int l_ret    = 0;                                                            \
    l_ret        = snprintf(                                                     \
        iMeIsV_sTr + l_offset, MaXlEn - l_offset, "%u%u%u%u%u%u%u%u",     \
        (iMeIsV_t_PtR)->tac1, (iMeIsV_t_PtR)->tac2, (iMeIsV_t_PtR)->tac3, \
        (iMeIsV_t_PtR)->tac4, (iMeIsV_t_PtR)->tac5, (iMeIsV_t_PtR)->tac6, \
        (iMeIsV_t_PtR)->tac7, (iMeIsV_t_PtR)->tac8);                      \
    if (l_ret > 0) {                                                             \
      l_offset += l_ret;                                                         \
      l_ret = snprintf(                                                          \
          iMeIsV_sTr + l_offset, MaXlEn - l_offset, "%u%u%u%u%u%u%u%u",          \
          (iMeIsV_t_PtR)->snr1, (iMeIsV_t_PtR)->snr2, (iMeIsV_t_PtR)->snr3,      \
          (iMeIsV_t_PtR)->snr4, (iMeIsV_t_PtR)->snr5, (iMeIsV_t_PtR)->snr6,      \
          (iMeIsV_t_PtR)->svn1, (iMeIsV_t_PtR)->svn2);                           \
    }                                                                            \
  }

/*Used to convert char* IMSI/TMSI Mobile Identity to MobileIdentity(digit)
 * format*/
#define MOBILE_ID_CHAR_TO_MOBILE_ID_IMSI_NAS(mObId_ChAr, mObId_PtR, iMsI_LeN)  \
  {                                                                            \
    uint32_t idx = 0;                                                          \
    if ((iMsI_LeN % 2) == 0) {                                                 \
      mObId_PtR->oddeven = false;                                              \
    } else {                                                                   \
      mObId_PtR->oddeven = true;                                               \
    }                                                                          \
    mObId_PtR->digit1 = (*(mObId_ChAr + idx) >> 4) & 0xf;                      \
    mObId_PtR->numOfValidImsiDigits++;                                         \
    idx++;                                                                     \
    if (idx <= iMsI_LeN) {                                                     \
      mObId_PtR->digit2 = *(mObId_ChAr + idx) & 0xf;                           \
      mObId_PtR->digit3 = (*(mObId_ChAr + idx) >> 4) & 0xf;                    \
      idx++;                                                                   \
      mObId_PtR->numOfValidImsiDigits += 2;                                    \
    }                                                                          \
    if (idx <= iMsI_LeN) {                                                     \
      mObId_PtR->digit4 = *(mObId_ChAr + idx) & 0xf;                           \
      mObId_PtR->digit5 = (*(mObId_ChAr + idx) >> 4) & 0xf;                    \
      idx++;                                                                   \
      mObId_PtR->numOfValidImsiDigits += 2;                                    \
    }                                                                          \
    if (idx <= iMsI_LeN) {                                                     \
      mObId_PtR->digit6 = *(mObId_ChAr + idx) & 0xf;                           \
      mObId_PtR->digit7 = (*(mObId_ChAr + idx) >> 4) & 0xf;                    \
      idx++;                                                                   \
      mObId_PtR->numOfValidImsiDigits += 2;                                    \
    }                                                                          \
    if (idx <= iMsI_LeN) {                                                     \
      mObId_PtR->digit8 = *(mObId_ChAr + idx) & 0xf;                           \
      mObId_PtR->digit9 = (*(mObId_ChAr + idx) >> 4) & 0xf;                    \
      idx++;                                                                   \
      mObId_PtR->numOfValidImsiDigits += 2;                                    \
    }                                                                          \
    if (idx <= iMsI_LeN) {                                                     \
      mObId_PtR->digit10 = *(mObId_ChAr + idx) & 0xf;                          \
      mObId_PtR->digit11 = (*(mObId_ChAr + idx) >> 4) & 0xf;                   \
      idx++;                                                                   \
      mObId_PtR->numOfValidImsiDigits += 2;                                    \
    }                                                                          \
    if (idx <= iMsI_LeN) {                                                     \
      mObId_PtR->digit12 = *(mObId_ChAr + idx) & 0xf;                          \
      mObId_PtR->digit13 = (*(mObId_ChAr + idx) >> 4) & 0xf;                   \
      idx++;                                                                   \
      mObId_PtR->numOfValidImsiDigits += 2;                                    \
    }                                                                          \
    if (idx <= iMsI_LeN) {                                                     \
      mObId_PtR->digit14 = *(mObId_ChAr + idx) & 0xf;                          \
      mObId_PtR->digit15 = (*(mObId_ChAr + idx) >> 4) & 0xf;                   \
      idx++;                                                                   \
      mObId_PtR->numOfValidImsiDigits += 2;                                    \
    }                                                                          \
    if (mObId_PtR->oddeven == false) {                                         \
      mObId_PtR->numOfValidImsiDigits--;                                       \
    }                                                                          \
  }

void imsi_string_to_3gpp_imsi(const Imsi_t* Imsi, imsi_t* imsi);
/*Used to convert char* IMSI/TMSI Mobile Identity to MobileIdentity(digit)
 * format*/
#define MOBILE_ID_CHAR_TO_MOBILE_ID_TMSI_NAS(mObId_ChAr, mObId_PtR, tMsI_LeN)  \
  {                                                                            \
    uint32_t idx = 0;                                                          \
    if ((tMsI_LeN % 2) == 0) {                                                 \
      mObId_PtR->oddeven = false;                                              \
    } else {                                                                   \
      mObId_PtR->oddeven = true;                                               \
    }                                                                          \
    for (idx = 0; idx < 4; idx++) {                                            \
      (mObId_PtR->tmsi)[idx] = mObId_ChAr[idx];                                \
    }                                                                          \
  }

void hexa_to_ascii(uint8_t* from, char* to, size_t length);

int ascii_to_hex(uint8_t* dst, const char* h);
#define UINT8_TO_BINARY_FMT "%c%c%c%c%c%c%c%c"
#define UINT8_TO_BINARY_ARG(bYtE)                                              \
  ((bYtE) &0x80 ? '1' : '0'), ((bYtE) &0x40 ? '1' : '0'),                      \
      ((bYtE) &0x20 ? '1' : '0'), ((bYtE) &0x10 ? '1' : '0'),                  \
      ((bYtE) &0x08 ? '1' : '0'), ((bYtE) &0x04 ? '1' : '0'),                  \
      ((bYtE) &0x02 ? '1' : '0'), ((bYtE) &0x01 ? '1' : '0')

int get_time_zone(void);
#endif /* FILE_CONVERSIONS_SEEN */
