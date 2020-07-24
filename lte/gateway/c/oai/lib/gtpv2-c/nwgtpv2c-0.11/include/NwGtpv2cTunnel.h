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

#ifndef __NW_GTPV2C_TUNNEL_H__
#define __NW_GTPV2C_TUNNEL_H__

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "tree.h"
#include "NwTypes.h"
#include "NwUtils.h"
#include "NwError.h"
#include "NwGtpv2c.h"

#ifdef __cplusplus
extern "C" {
#endif

struct nw_gtpv2c_stack_s;

typedef struct nw_gtpv2c_tunnel_s {
  uint32_t teid;
  union {
    struct sockaddr_in ipv4_addr;
    struct sockaddr_in6 ipv6_addr;
  } ipAddrRemote;

  nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel;
  RB_ENTRY(nw_gtpv2c_tunnel_s)
  tunnelMapRbtNode; /**< RB Tree Data Structure Node        */
  struct nw_gtpv2c_tunnel_s* next;
} nw_gtpv2c_tunnel_t;

nw_gtpv2c_tunnel_t* nwGtpv2cTunnelNew(
    struct nw_gtpv2c_stack_s* hStack, uint32_t teid,
    struct sockaddr* ipAddrRemote, nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel);

nw_rc_t nwGtpv2cTunnelDelete(
    struct nw_gtpv2c_stack_s* pStack, nw_gtpv2c_tunnel_t* thiz);

nw_rc_t nwGtpv2cTunnelGetUlpTunnelHandle(
    nw_gtpv2c_tunnel_t* thiz, nw_gtpv2c_ulp_tunnel_handle_t* phUlpTunnel);

#ifdef __cplusplus
}
#endif

#endif

/*--------------------------------------------------------------------------*
 *                      E N D     O F    F I L E                            *
 *--------------------------------------------------------------------------*/
