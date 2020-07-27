/*----------------------------------------------------------------------------*
 *                                                                            *
 *                              n w - g t p v 2 c                             *
 *    G P R S   T u n n e l i n g    P r o t o c o l   v 2 c    S t a c k     *
 *                                                                            *
 *                                                                            *
 * Copyright (c) 2010-2011 Amit Chawre                                        *
 * All rights reserved.                                                       *
 *                                                                            *
 * Redistribution and use in source and binary forms, with or without         *
 * modification, are permitted provided that the following conditions         *
 * are met:                                                                   *
 *                                                                            *
 * 1. Redistributions of source code must retain the above copyright          *
 *    notice, this list of conditions and the following disclaimer.           *
 * 2. Redistributions in binary form must reproduce the above copyright       *
 *    notice, this list of conditions and the following disclaimer in the     *
 *    documentation and/or other materials provided with the distribution.    *
 * 3. The name of the author may not be used to endorse or promote products   *
 *    derived from this software without specific prior written permission.   *
 *                                                                            *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR       *
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES  *
 * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.    *
 * IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,           *
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT   *
 * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,  *
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY      *
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT        *
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF   *
 * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.          *
 *----------------------------------------------------------------------------*/

#include <assert.h>

#ifndef __NW_UTILS_H__
#define __NW_UTILS_H__

/**
 * @file NwUtils.h
 * @brief This header file contains utility macro and function definitions.
 */

#define NW_ASSERT assert /**< Assertion */

#define NW_CHK_NULL_PTR(_ptr)                                                  \
  NW_ASSERT(_ptr != NULL) /**< Null pointer check                              \
                           */

#define NW_HTONS(x) ((((x) &0xff00) >> 8) | (((x) &0x00ff) << 8))

#define NW_HTONL(x)                                                            \
  ((((x) &0xff000000) >> 24) | (((x) &0x00ff0000) >> 8) |                      \
   (((x) &0x0000ff00) << 8) | (((x) &0x000000ff) << 24))

#define NW_HTONLL(x)                                                           \
  (((((uint64_t) x) & 0xff00000000000000ULL) >> 56) |                          \
   ((((uint64_t) x) & 0x00ff000000000000ULL) >> 40) |                          \
   ((((uint64_t) x) & 0x0000ff0000000000ULL) >> 24) |                          \
   ((((uint64_t) x) & 0x000000ff00000000ULL) >> 8) |                           \
   ((((uint64_t) x) & 0x000000000000ff00ULL) << 40) |                          \
   ((((uint64_t) x) & 0x00000000000000ffULL) << 56) |                          \
   ((((uint64_t) x) & 0x0000000000ff0000ULL) << 24) |                          \
   ((((uint64_t) x) & 0x00000000ff000000ULL) << 8))

#define NW_NTOHS NW_HTONS
#define NW_NTOHL NW_HTONL
#define NW_NTOHLL NW_HTONLL

#endif /* __NW_UTILS_H__ */
