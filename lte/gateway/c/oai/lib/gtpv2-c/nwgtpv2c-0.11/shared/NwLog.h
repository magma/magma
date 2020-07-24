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
#ifndef __NW_LOG_H__
#define __NW_LOG_H__

#include <libgen.h>

#include "NwTypes.h"

/**
 * @file NwLog.h
 * @brief This header file contains global shared logging definitions.
 */

#ifdef __cplusplus
extern "C" {
#endif

/*---------------------------------------------------------------------------
 * Log Level Definitions
 *--------------------------------------------------------------------------*/

#define NW_LOG_LEVEL_EMER (0) /**< system is unusable              */
#define NW_LOG_LEVEL_ALER (1) /**< action must be taken immediately*/
#define NW_LOG_LEVEL_CRIT (2) /**< critical conditions             */
#define NW_LOG_LEVEL_ERRO (3) /**< error conditions                */
#define NW_LOG_LEVEL_WARN (4) /**< warning conditions              */
#define NW_LOG_LEVEL_NOTI (5) /**< normal but signification condition */
#define NW_LOG_LEVEL_INFO (6) /**< informational                   */
#define NW_LOG_LEVEL_DEBG (7) /**< debug-level messages            */

/*---------------------------------------------------------------------------
 * IPv4 logging macros
 *--------------------------------------------------------------------------*/

#define NW_IPV4_ADDR "%u.%u.%u.%u"
#define NW_IPV4_ADDR_FORMAT(__addr)                                            \
  (uint8_t)((__addr) &0x000000ff), (uint8_t)(((__addr) &0x0000ff00) >> 8),     \
      (uint8_t)(((__addr) &0x00ff0000) >> 16),                                 \
      (uint8_t)(((__addr) &0xff000000) >> 24)

#define NW_IPV4_ADDR_FORMATP(__paddr)                                          \
  (uint8_t)(*((uint8_t*) (__paddr)) & 0x000000ff),                             \
      (uint8_t)(*((uint8_t*) (__paddr + 1)) & 0x000000ff),                     \
      (uint8_t)(*((uint8_t*) (__paddr + 2)) & 0x000000ff),                     \
      (uint8_t)(*((uint8_t*) (__paddr + 3)) & 0x000000ff)

#ifdef __cplusplus
}
#endif

#endif /* __NW_LOG_H__ */
