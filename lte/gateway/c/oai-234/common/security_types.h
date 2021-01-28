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
/*! \file security_types.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdint.h>

#if HAVE_CONFIG_H
#include "config.h"
#endif

#ifndef FILE_SECURITY_TYPES_SEEN
#define FILE_SECURITY_TYPES_SEEN

#define RAND_LENGTH_BITS (128)
#define RAND_LENGTH_OCTETS (RAND_LENGTH_BITS / 8)
#define AUTH_KEY_LENGTH_BITS (128)
#define AUTH_KEY_LENGTH_OCTETS (AUTH_KEY_LENGTH_BITS / 8)
#define KASME_LENGTH_BITS (256)
#define KASME_LENGTH_OCTETS (KASME_LENGTH_BITS / 8)
/* In OCTETS */
#define XRES_LENGTH_MIN (4)
#define XRES_LENGTH_MAX (16)
#define AUTN_LENGTH_BITS (128)
#define AUTN_LENGTH_OCTETS (AUTN_LENGTH_BITS / 8)

/* Some methods to convert a string to an int64_t */
/*
#define STRING_TO_64BITS(sTRING, cONTAINER)    \
    sscanf(sTRING, "%" SCN64, cONTAINER)
#define STRING_TO_U64BITS(sTRING, cONTAINER)    \
    sscanf(sTRING, "%" SCNu64, cONTAINER)
*/

/* Converts a string to 128 bits gmplib integer holder */
#define STRING_TO_XBITS(sTRING, lENGTH, cONTAINER, rET)                        \
  do {                                                                         \
    memcpy(cONTAINER, sTRING, lENGTH);                                         \
    rET = 0;                                                                   \
  } while (0)

#define STRING_TO_128BITS(sTRING, cONTAINER, rET)                              \
  STRING_TO_XBITS(sTRING, 16, cONTAINER, rET)

#define STRING_TO_256BITS(sTRING, cONTAINER, rET)                              \
  STRING_TO_XBITS(sTRING, 32, cONTAINER, rET)

#define STRING_TO_RAND STRING_TO_128BITS
#define STRING_TO_AUTH_KEY STRING_TO_128BITS
#define STRING_TO_AUTH_RES STRING_TO_128BITS
#define STRING_TO_AUTN STRING_TO_128BITS
#define STRING_TO_KASME STRING_TO_256BITS
#define STRING_TO_XRES(sTRING, lENGTH, cONTAINER, rET)                         \
  do {                                                                         \
    STRING_TO_XBITS(sTRING, lENGTH, (cONTAINER)->data, rET);                   \
    if (rET != -1) (cONTAINER)->size = lENGTH;                                 \
  } while (0)

/* RES amd XRES can have a variable length of 4-16 octets */
typedef struct {
  uint8_t size;
  uint8_t data[XRES_LENGTH_MAX];
} res_t;

#define FORMAT_128BITS                                                         \
  "%02x,%02x,%02x,%02x,%02x,%02x,%02x,%02x,%02x,%02x,%02x,%02x,%02x,%02x,%"    \
  "02x,%02x"
#define XRES_FORMAT FORMAT_128BITS
#define RAND_FORMAT FORMAT_128BITS
#define AUTN_FORMAT FORMAT_128BITS
#define KASME_FORMAT FORMAT_128BITS

#define DISPLAY_128BITS(bUFFER)                                                \
  bUFFER[0], bUFFER[1], bUFFER[2], bUFFER[3], bUFFER[4], bUFFER[5], bUFFER[6], \
      bUFFER[7], bUFFER[8], bUFFER[9], bUFFER[10], bUFFER[11], bUFFER[12],     \
      bUFFER[13], bUFFER[14], bUFFER[15]
#define DISPLAY_128BITS_2(bUFFER)                                              \
  bUFFER[16], bUFFER[17], bUFFER[18], bUFFER[19], bUFFER[20], bUFFER[21],      \
      bUFFER[22], bUFFER[23], bUFFER[24], bUFFER[25], bUFFER[26], bUFFER[27],  \
      bUFFER[28], bUFFER[29], bUFFER[30], bUFFER[31]

#define XRES_DISPLAY(bUFFER) DISPLAY_128BITS(bUFFER)
#define RAND_DISPLAY(bUFFER) DISPLAY_128BITS(bUFFER)
#define AUTN_DISPLAY(bUFFER) DISPLAY_128BITS(bUFFER)
/* Display only first 128 bits of KASME */
#define KASME_DISPLAY_1(bUFFER) DISPLAY_128BITS(bUFFER)
#define KASME_DISPLAY_2(bUFFER) DISPLAY_128BITS_2(bUFFER)

/* Holds an E-UTRAN authentication vector */
typedef struct eutran_vector_s {
  uint8_t rand[RAND_LENGTH_OCTETS];
  res_t xres;
  uint8_t autn[AUTN_LENGTH_OCTETS];
  uint8_t kasme[KASME_LENGTH_OCTETS];
} eutran_vector_t;

#define FC_KASME (0x10)
#define FC_KENB (0x11)
#define FC_NH (0x12)
#define FC_KENB_STAR (0x13)
/* 33401 #A.7 Algorithm for key derivation function.
 * This FC should be used for:
 * - NAS Encryption algorithm
 * - NAS Integrity algorithm
 * - RRC Encryption algorithm
 * - RRC Integrity algorithm
 * - User Plane Encryption algorithm
 */
#define FC_ALG_KEY_DER (0x15)
#define FC_KASME_TO_CK (0x16)

typedef enum {
  NAS_ENC_ALG = 0x01,
  NAS_INT_ALG = 0x02,
  RRC_ENC_ALG = 0x03,
  RRC_INT_ALG = 0x04,
  UP_ENC_ALG  = 0x05,
  UP_INT_ALG  = 0x06
} algorithm_type_dist_t;

#endif /* FILE_SECURITY_TYPES_SEEN */
