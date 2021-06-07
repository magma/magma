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

/*! \file asn1_conversions.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_ASN1_CONVERSIONS_SEEN
#define FILE_ASN1_CONVERSIONS_SEEN

#include "BIT_STRING.h"
#include "assertions.h"

//-----------------------begin func -------------------

/*! \fn uint8_t BIT_STRING_to_uint8(BIT_STRING_t *)
 * \brief This function extract at most a 8 bits value from a BIT_STRING_t
 * object, the exact bits number depend on the BIT_STRING_t contents.
 * \param[in] pointer to the BIT_STRING_t object.
 * \return the extracted value.
 */
static inline uint8_t BIT_STRING_to_uint8(BIT_STRING_t* asn) {
  DevCheck((asn->size == 1), asn->size, 0, 0);

  return asn->buf[0] >> asn->bits_unused;
}

/*! \fn uint16_t BIT_STRING_to_uint16(BIT_STRING_t *)
 * \brief This function extract at most a 16 bits value from a BIT_STRING_t
 * object, the exact bits number depend on the BIT_STRING_t contents.
 * \param[in] pointer to the BIT_STRING_t object.
 * \return the extracted value.
 */
static inline uint16_t BIT_STRING_to_uint16(BIT_STRING_t* asn) {
  uint16_t result = 0;
  int index       = 0;

  DevCheck((asn->size > 0) && (asn->size <= 2), asn->size, 0, 0);

  switch (asn->size) {
    case 2:
      result |= asn->buf[index++] << (8 - asn->bits_unused);

    case 1:
      result |= asn->buf[index] >> asn->bits_unused;
      break;

    default:
      break;
  }

  return result;
}

/*! \fn uint32_t BIT_STRING_to_uint32(BIT_STRING_t *)
 * \brief  This function extract at most a 32 bits value from a BIT_STRING_t
 * object, the exact bits number depend on the BIT_STRING_t contents.
 * \param[in] pointer to the BIT_STRING_t object.
 * \return the extracted value.
 */
static inline uint32_t BIT_STRING_to_uint32(BIT_STRING_t* asn) {
  uint32_t result = 0;
  int index;
  int shift;

  DevCheck((asn->size > 0) && (asn->size <= 4), asn->size, 0, 0);

  shift = ((asn->size - 1) * 8) - asn->bits_unused;
  for (index = 0; index < (asn->size - 1); index++) {
    result |= asn->buf[index] << shift;
    shift -= 8;
  }

  result |= asn->buf[index] >> asn->bits_unused;

  return result;
}

/*! \fn uint64_t BIT_STRING_to_uint64(BIT_STRING_t *)
 * \brief  This function extract at most a 64 bits value from a BIT_STRING_t
 * object, the exact bits number depend on the BIT_STRING_t contents.
 * \param[in] pointer to the BIT_STRING_t object.
 * \return the extracted value.
 */
static inline uint64_t BIT_STRING_to_uint64(BIT_STRING_t* asn) {
  uint64_t result = 0;
  int index;
  int shift;

  DevCheck((asn->size > 0) && (asn->size <= 8), asn->size, 0, 0);

  shift = ((asn->size - 1) * 8) - asn->bits_unused;
  for (index = 0; index < (asn->size - 1); index++) {
    result |= asn->buf[index] << shift;
    shift -= 8;
  }

  result |= asn->buf[index] >> asn->bits_unused;

  return result;
}

#endif /* FILE_ASN1_CONVERSIONS_SEEN */
