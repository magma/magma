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

/** @brief Intertask Interface common types
 * Contains type definitions used for generating and parsing ITTI messages.
 * @author Laurent Winckel <laurent.winckel@eurecom.fr>
 */

#ifndef _ITTI_TYPES_H_
#define _ITTI_TYPES_H_

#include <stdint.h>

#define CHARS_TO_UINT32(c1, c2, c3, c4)                                        \
  ((((uint32_t) c4) << 24) | (((uint32_t) c3) << 16) |                         \
   (((uint32_t) c2) << 8) | ((uint32_t) c1))

#define MESSAGE_NUMBER_CHAR_FORMAT "%11u"

/* Intertask message types */
enum itti_message_types_e {
  ITTI_DUMP_XML_DEFINITION     = CHARS_TO_UINT32('\n', 'I', 'x', 'd'),
  ITTI_DUMP_XML_DEFINITION_END = CHARS_TO_UINT32('i', 'X', 'D', '\n'),

  ITTI_DUMP_MESSAGE_TYPE     = CHARS_TO_UINT32('\n', 'I', 'm', 's'),
  ITTI_DUMP_MESSAGE_TYPE_END = CHARS_TO_UINT32('i', 'M', 'S', '\n'),

  ITTI_STATISTIC_MESSAGE_TYPE     = CHARS_TO_UINT32('\n', 'I', 's', 't'),
  ITTI_STATISTIC_MESSAGE_TYPE_END = CHARS_TO_UINT32('i', 'S', 'T', '\n'),

  /* This signal is not meant to be used by remote analyzer */
  ITTI_DUMP_EXIT_SIGNAL = CHARS_TO_UINT32('e', 'X', 'I', 'T'),
};

typedef uint32_t itti_message_types_t;

/* Message header is the common part that should never change between
 * remote process and this one.
 */
typedef struct {
  /* The size of this structure */
  uint32_t message_size;
  itti_message_types_t message_type;
} itti_socket_header_t;

typedef struct {
  char message_number_char[12]; /* 9 chars are needed to store an unsigned 32
                                   bits value in decimal, but must be a multiple
                                   of 32 bits to avoid alignment issues */
} itti_signal_header_t;

#define INSTANCE_DEFAULT (UINT16_MAX - 1)
#define INSTANCE_ALL (UINT16_MAX)

typedef uint16_t instance_t;

#endif
