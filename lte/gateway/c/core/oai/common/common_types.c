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

/*! \file common_types.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <string.h>
#include <stdlib.h>

#include "bstrlib.h"
#include "assertions.h"
#include "3gpp_23.003.h"
#include "common_types.h"
#include "3gpp_29.274.h"
#include "log.h"
#include "hashtable.h"
#include "mme_config.h"

/* Clear GUTI without free it */
void clear_guti(guti_t* const guti) {
  memset(guti, 0, sizeof(guti_t));
  guti->m_tmsi = INVALID_TMSI;
}
/* Clear IMSI without free it */
void clear_imsi(imsi_t* const imsi) {
  memset(imsi, 0, sizeof(imsi_t));
}
/* Clear IMEI without free it */
void clear_imei(imei_t* const imei) {
  memset(imei, 0, sizeof(imei_t));
}
/* Clear IMEISV without free it */
void clear_imeisv(imeisv_t* const imeisv) {
  memset(imeisv, 0, sizeof(imeisv_t));
}

//------------------------------------------------------------------------------
bstring fteid_ip_address_to_bstring(const struct fteid_s* const fteid) {
  bstring bstr = NULL;
  if (fteid->ipv4) {
    bstr = blk2bstr(&fteid->ipv4_address.s_addr, 4);
  } else if (fteid->ipv6) {
    bstr = blk2bstr(&fteid->ipv6_address, 16);
  } else {
    char avoid_seg_fault[4] = {0, 0, 0, 0};
    bstr                    = blk2bstr(avoid_seg_fault, 4);
  }
  return bstr;
}

void get_fteid_ip_address(
    const struct fteid_s* const fteid, ip_address_t* const ip_address) {
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
    default:;
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
      default:;
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
        OAILOG_ERROR(
            LOG_COMMON, "Invalid ipv6_prefix_length : %u\n",
            paa->ipv6_prefix_length);
      }
      break;
    case IPv4_AND_v6:
      if (paa->ipv6_prefix_length == IPV6_PREFIX_LEN) {
        bstr = blk2bstr(&paa->ipv6_address, paa->ipv6_prefix_length / 8);
        bcatblk(bstr, &paa->ipv4_address, 4);
      } else {
        OAILOG_ERROR(
            LOG_COMMON, "Invalid ipv6_prefix_length : %u\n",
            paa->ipv6_prefix_length);
      }
      break;
    case IPv4_OR_v6:
      // do it like that now, TODO
      bstr = blk2bstr(&paa->ipv4_address.s_addr, 4);
      break;
    default:;
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
