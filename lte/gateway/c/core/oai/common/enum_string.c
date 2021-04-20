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

/*! \file enum_string.c
  \brief
  \author Sebastien ROUX
  \company Eurecom
*/

#include <stdio.h>
#include <stdlib.h>

#include "common_types.h"
#include "enum_string.h"

enum_to_string_t rat_to_string[NUMBER_OF_RAT_TYPE] = {
    {RAT_WLAN, "WLAN"},      {RAT_VIRTUAL, "VIRUTAL"},
    {RAT_UTRAN, "UTRAN"},    {RAT_GERAN, "GERAN"},
    {RAT_GAN, "GAN"},        {RAT_HSPA_EVOLUTION, "HSPA_EVOLUTION"},
    {RAT_EUTRAN, "E-UTRAN"}, {RAT_CDMA2000_1X, "CDMA2000_1X"},
    {RAT_HRPD, "HRPD"},      {RAT_UMB, "UMB"},
    {RAT_EHRPD, "EHRPD"},
};

enum_to_string_t network_access_mode_to_string[NAM_MAX] = {
    {NAM_PACKET_AND_CIRCUIT, "PACKET AND CIRCUIT"},
    {NAM_RESERVED, "RESERVED"},
    {NAM_ONLY_PACKET, "ONLY PACKET"},
};

enum_to_string_t all_apn_conf_ind_to_string[ALL_APN_MAX] = {
    {ALL_APN_CONFIGURATIONS_INCLUDED, "ALL APN CONFIGURATIONS INCLUDED"},
    {MODIFIED_ADDED_APN_CONFIGURATIONS_INCLUDED,
     "MODIFIED ADDED APN CONFIGURATIONS INCLUDED"},
};

enum_to_string_t pdn_type_to_string[IP_MAX] = {
    {IPv4, "IPv4"},
    {IPv6, "IPv6"},
    {IPv4_AND_v6, "IPv4 and IPv6"},
    {IPv4_OR_v6, "IPv4 or IPv6"},
};

static int compare_values(const void* m1, const void* m2) {
  enum_to_string_t* mi1 = (enum_to_string_t*) m1;
  enum_to_string_t* mi2 = (enum_to_string_t*) m2;

  return (mi1->enum_value - mi2->enum_value);
}

char* enum_to_string(
    int enum_val, enum_to_string_t* string_table, int nb_element) {
  enum_to_string_t* res;
  enum_to_string_t temp;

  temp.enum_value = enum_val;
  res             = bsearch(
      &temp, string_table, nb_element, sizeof(enum_to_string_t),
      compare_values);

  if (res == NULL) {
    return "UNKNOWN";
  }

  return res->enum_value_name;
}
