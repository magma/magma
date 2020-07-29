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

/*! \file enum_string.h
  \brief
  \author Sebastien ROUX
  \company Eurecom
*/
#ifndef FILE_ENUM_STRING_SEEN
#define FILE_ENUM_STRING_SEEN

#include "common_types.h"

typedef struct {
  int enum_value;
  char* enum_value_name;
} enum_to_string_t;

extern enum_to_string_t network_access_mode_to_string[NAM_MAX];
extern enum_to_string_t rat_to_string[NUMBER_OF_RAT_TYPE];
extern enum_to_string_t pdn_type_to_string[IP_MAX];

char* enum_to_string(
    int enum_val, enum_to_string_t* string_table, int nb_element);
#define ACCESS_MODE_TO_STRING(vAL)                                             \
  enum_to_string(                                                              \
      (int) vAL, network_access_mode_to_string,                                \
      sizeof(network_access_mode_to_string) / sizeof(enum_to_string_t))
#define PDN_TYPE_TO_STRING(vAL)                                                \
  enum_to_string(                                                              \
      (int) vAL, pdn_type_to_string,                                           \
      sizeof(pdn_type_to_string) / sizeof(enum_to_string_t))

#endif /* FILE_ENUM_STRING_SEEN */
