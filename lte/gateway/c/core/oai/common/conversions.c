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

/*! \file conversions.c
  \brief
  \author Sebastien ROUX
  \company Eurecom
*/

#include <stdint.h>
#include <ctype.h>

#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"

static const char hex_to_ascii_table[16] = {
    '0', '1', '2', '3', '4', '5', '6', '7',
    '8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
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

void hexa_to_ascii(uint8_t* from, char* to, size_t length) {
  size_t i;

  for (i = 0; i < length; i++) {
    uint8_t upper = (from[i] & 0xf0) >> 4;
    uint8_t lower = from[i] & 0x0f;

    to[2 * i] = hex_to_ascii_table[upper];
    to[2 * i + 1] = hex_to_ascii_table[lower];
  }
}

int ascii_to_hex(uint8_t* dst, const char* h) {
  const unsigned char* hex = (const unsigned char*)h;
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
imsi64_t imsi_to_imsi64(const imsi_t* const imsi) {
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

//-----------------------------------------------------------------------------
void imsi_string_to_3gpp_imsi(const Imsi_t* Imsi, imsi_t* imsi) {
  memset(imsi->u.value, 0xff, IMSI_BCD8_SIZE);
  imsi->u.num.digit1 = Imsi->digit[0] - 0x30;
  imsi->u.num.digit2 = Imsi->digit[1] - 0x30;
  imsi->u.num.digit3 = Imsi->digit[2] - 0x30;
  imsi->u.num.digit4 = Imsi->digit[3] - 0x30;
  imsi->u.num.digit5 = Imsi->digit[4] - 0x30;
  imsi->u.num.digit6 = Imsi->digit[5] - 0x30;
  if (Imsi->length >= 7) {
    imsi->u.num.digit7 = Imsi->digit[6] - 0x30;
    if (Imsi->length >= 8) {
      imsi->u.num.digit8 = Imsi->digit[7] - 0x30;
      if (Imsi->length >= 9) {
        imsi->u.num.digit9 = Imsi->digit[8] - 0x30;
        if (Imsi->length >= 10) {
          imsi->u.num.digit10 = Imsi->digit[9] - 0x30;
          if (Imsi->length >= 11) {
            imsi->u.num.digit11 = Imsi->digit[10] - 0x30;
            if (Imsi->length >= 12) {
              imsi->u.num.digit12 = Imsi->digit[11] - 0x30;
              if (Imsi->length >= 13) {
                imsi->u.num.digit13 = Imsi->digit[12] - 0x30;
                if (Imsi->length >= 14) {
                  imsi->u.num.digit14 = Imsi->digit[13] - 0x30;
                  if (Imsi->length >= 15) {
                    imsi->u.num.digit15 = Imsi->digit[14] - 0x30;
                  }
                }
              }
            }
          }
        }
      }
    }
  }
  imsi->length = Imsi->length;
}

//------------------------------------------------------------------------------
imsi64_t amf_imsi_to_imsi64(const imsi_t* const imsi) {
  imsi64_t imsi64 = INVALID_IMSI64;
  bool skip_last_digit = false;

  if (imsi) {
    imsi64 = 0;
    for (int i = 0; i < IMSI_BCD8_SIZE; i++) {
      /*Bring 2 digits to LSB and calculate
       * each digit/nibel range would be 0 to 9
       */
      uint8_t d2 = imsi->u.value[i];
      uint8_t d1 = (d2 & 0xf0) >> 4;
      d2 = d2 & 0x0f;
      if (d1 < 10) {
        imsi64 = imsi64 * 10 + d1;
      } else {
        skip_last_digit = true;
      }

      if (d2 < 10) {
        imsi64 = imsi64 * 10 + d2;
      } else {
        skip_last_digit = true;
      }
    }
  }

  // As per latest sim behavior "2224560000000010" is received
  // This implies first 15 digits are valid and last is  not considered
  if (skip_last_digit == false) {
    imsi64 = imsi64 / 10;
  }

  return imsi64;
}

/* Clear GUTI without free it */
void clear_guti(guti_t* const guti) {
  memset(guti, 0, sizeof(guti_t));
  guti->m_tmsi = INVALID_TMSI;
}
/* Clear IMSI without free it */
void clear_imsi(imsi_t* const imsi) { memset(imsi, 0, sizeof(imsi_t)); }
/* Clear IMEI without free it */
void clear_imei(imei_t* const imei) { memset(imei, 0, sizeof(imei_t)); }
/* Clear IMEISV without free it */
void clear_imeisv(imeisv_t* const imeisv) {
  memset(imeisv, 0, sizeof(imeisv_t));
}

//------------------------------------------------------------------------------
bstring fteid_ip_address_to_bstring(const struct fteid_s* const fteid) {
  bstring bstr = NULL;
  if (fteid->ipv4 && fteid->ipv6) {
    bstring ipv6 = NULL;
    bstr = blk2bstr(&fteid->ipv4_address.s_addr, 4);
    ipv6 = blk2bstr(&fteid->ipv6_address, 16);
    bconcat(bstr, ipv6);
    bdestroy(ipv6);
  } else if (fteid->ipv4) {
    bstr = blk2bstr(&fteid->ipv4_address.s_addr, 4);
  } else if (fteid->ipv6) {
    bstr = blk2bstr(&fteid->ipv6_address, 16);
  } else {
    char avoid_seg_fault[4] = {0, 0, 0, 0};
    bstr = blk2bstr(avoid_seg_fault, 4);
  }
  return bstr;
}

void get_fteid_ip_address(const struct fteid_s* const fteid,
                          ip_address_t* const ip_address) {
  if (fteid->ipv4) {
    ip_address->pdn_type = IPv4;
    memcpy(&ip_address->address.ipv4_address, &fteid->ipv4_address, 4);
  }

  if (fteid->ipv6) {
    ip_address->pdn_type = IPv6;
    memcpy(&ip_address->address.ipv6_address, &fteid->ipv6_address, 16);
  }

  if (fteid->ipv4 && fteid->ipv6) {
    ip_address->pdn_type = IPv4_AND_v6;
  }
}

//------------------------------------------------------------------------------
bstring ip_address_to_bstring(const ip_address_t* ip_address) {
  bstring bstr = NULL;
  switch (ip_address->pdn_type) {
    case IPv4:
      bstr = blk2bstr(&ip_address->address.ipv4_address.s_addr, 4);
      break;
    case IPv6:
      bstr = blk2bstr(&ip_address->address.ipv6_address, 16);
      break;
    case IPv4_AND_v6:
      bstr = blk2bstr(&ip_address->address.ipv4_address.s_addr, 4);
      bcatblk(bstr, &ip_address->address.ipv6_address, 16);
      break;
    case IPv4_OR_v6:
      // do it like that now, TODO
      bstr = blk2bstr(&ip_address->address.ipv4_address.s_addr, 4);
      break;
    default: {
    }
  }
  return bstr;
}

//------------------------------------------------------------------------------
void bstring_to_ip_address(bstring const bstr, ip_address_t* const ip_address) {
  if (bstr) {
    switch (blength(bstr)) {
      case 4:
        ip_address->pdn_type = IPv4;
        memcpy(&ip_address->address.ipv4_address, bstr->data, blength(bstr));
        break;
      case 16:
        ip_address->pdn_type = IPv6;
        memcpy(&ip_address->address.ipv6_address, bstr->data, blength(bstr));
        break;
        break;
      case 20:
        ip_address->pdn_type = IPv4_AND_v6;
        memcpy(&ip_address->address.ipv4_address, bstr->data, 4);
        memcpy(&ip_address->address.ipv6_address, &bstr->data[4], 16);
        break;
      default: {
      }
    }
  }
}

//------------------------------------------------------------------------------
void copy_paa(paa_t* paa_dst, paa_t* paa_src) {
  memcpy(paa_dst, paa_src, sizeof(paa_t));
}

//------------------------------------------------------------------------------
bstring paa_to_bstring(const paa_t* paa) {
  bstring bstr = NULL;
  switch (paa->pdn_type) {
    case IPv4:
      bstr = blk2bstr(&paa->ipv4_address.s_addr, 4);
      break;
    case IPv6:
      if (paa->ipv6_prefix_length == IPV6_PREFIX_LEN) {
        bstr = blk2bstr(&paa->ipv6_address, paa->ipv6_prefix_length / 8);
      } else {
        OAILOG_ERROR(LOG_COMMON, "Invalid ipv6_prefix_length : %u\n",
                     paa->ipv6_prefix_length);
      }
      break;
    case IPv4_AND_v6:
      if (paa->ipv6_prefix_length == IPV6_PREFIX_LEN) {
        bstr = blk2bstr(&paa->ipv6_address, paa->ipv6_prefix_length / 8);
        bcatblk(bstr, &paa->ipv4_address, 4);
      } else {
        OAILOG_ERROR(LOG_COMMON, "Invalid ipv6_prefix_length : %u\n",
                     paa->ipv6_prefix_length);
      }
      break;
    case IPv4_OR_v6:
      // do it like that now, TODO
      bstr = blk2bstr(&paa->ipv4_address.s_addr, 4);
      break;
    default: {
    }
  }
  return bstr;
}

void bstring_to_paa(const bstring bstr, paa_t* paa) {
  if (bstr) {
    switch (blength(bstr)) {
      case 4:
        paa->pdn_type = IPv4;
        memcpy(&paa->ipv4_address, bstr->data, blength(bstr));
        break;
      case 8:
        paa->pdn_type = IPv6;
        memcpy(&paa->ipv6_address, bstr->data, blength(bstr));
        break;
      case 12:
        paa->pdn_type = IPv4_AND_v6;
        memcpy(&paa->ipv4_address, bstr->data, 4);
        memcpy(&paa->ipv6_address, &bstr->data[4], 8);
        break;
      default:
        break;
    }
  }
}

// Return the hex representation of a char array
char* bytes_to_hex(char* byte_array, int length, char* hex_array) {
  int i;
  for (i = 0; i < length; i++) {
    snprintf(hex_array + i * 3, 3, " %02x", (unsigned char)byte_array[i]);
  }
  return hex_array;
}

void convert_guti_to_string(const guti_t* guti_p,
                            char (*guti_str)[GUTI_STRING_LEN]) {
  snprintf(*guti_str, GUTI_STRING_LEN, "%x%x%x%x%x%x%04x%02x%08x",
           guti_p->gummei.plmn.mcc_digit1, guti_p->gummei.plmn.mcc_digit2,
           guti_p->gummei.plmn.mcc_digit3, guti_p->gummei.plmn.mnc_digit1,
           guti_p->gummei.plmn.mnc_digit2, guti_p->gummei.plmn.mnc_digit3,
           guti_p->gummei.mme_gid, guti_p->gummei.mme_code, guti_p->m_tmsi);
}
