/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */

/*! \file conversions.c
  \brief
  \author Sebastien ROUX
  \company Eurecom
*/
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <ctype.h>
#include <inttypes.h>

#include "common_defs.h"
#include "conversions.h"

static const char hex_to_ascii_table[16] = {
  '0',
  '1',
  '2',
  '3',
  '4',
  '5',
  '6',
  '7',
  '8',
  '9',
  'a',
  'b',
  'c',
  'd',
  'e',
  'f',
};

static const signed char ascii_to_hex_table[0x100] = {
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,  1,  2,  3,  4,  5,  6,  7,  8,
  9,  -1, -1, -1, -1, -1, -1, -1, 10, 11, 12, 13, 14, 15, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, 10, 11, 12, 13, 14, 15, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
  -1, -1, -1, -1, -1, -1, -1, -1, -1};

void hexa_to_ascii(uint8_t *from, char *to, size_t length)
{
  size_t i;

  for (i = 0; i < length; i++) {
    uint8_t upper = (from[i] & 0xf0) >> 4;
    uint8_t lower = from[i] & 0x0f;

    to[2 * i] = hex_to_ascii_table[upper];
    to[2 * i + 1] = hex_to_ascii_table[lower];
  }
}

int ascii_to_hex(uint8_t *dst, const char *h)
{
  const unsigned char *hex = (const unsigned char *) h;
  unsigned i = 0;

  for (;;) {
    int high, low;

    while (*hex && isspace(*hex)) hex++;

    if (!*hex) return 1;

    high = ascii_to_hex_table[*hex++];

    if (high < 0) return 0;

    while (*hex && isspace(*hex)) hex++;

    if (!*hex) return 0;

    low = ascii_to_hex_table[*hex++];

    if (low < 0) return 0;

    dst[i++] = (high << 4) | low;
  }
}
//------------------------------------------------------------------------------
imsi64_t imsi_to_imsi64(const imsi_t *const imsi)
{
  imsi64_t imsi64 = INVALID_IMSI64;
  if (imsi) {
    imsi64 = 0;
    for (int i = 0; i < IMSI_BCD8_SIZE; i++) {
      uint8_t d2 = imsi->u.value[i];
      uint8_t d1 = (d2 & 0xf0) >> 4;
      d2 = d2 & 0x0f;
      if (10 > d1) {
        imsi64 = imsi64 * 10 + d1;
        if (10 > d2) {
          imsi64 = imsi64 * 10 + d2;
        } else {
          break;
        }
      } else {
        break;
      }
    }
  }
  return imsi64;
}
