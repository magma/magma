/*----------------------------------------------------------------------------*
 *                                                                            *
                                n w - g t p v 2 c
      G P R S   T u n n e l i n g    P r o t o c o l   v 2 c    S t a c k
 *                                                                            *
 *                                                                            *
   Copyright (c) 2010-2011 Amit Chawre
   All rights reserved.
 *                                                                            *
   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions
   are met:
 *                                                                            *
   1. Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
   2. Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
   3. The name of the author may not be used to endorse or promote products
      derived from this software without specific prior written permission.
 *                                                                            *
   THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
   IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
   OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
   IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,
   INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
   NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
   DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
   THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
   (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
   THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
  ----------------------------------------------------------------------------*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#include "NwTypes.h"
#include "NwUtils.h"
#include "NwError.h"
#include "NwGtpv2cPrivate.h"
#include "NwGtpv2cTunnel.h"

#ifdef __cplusplus
extern "C" {
#endif

static nw_gtpv2c_tunnel_t* gpGtpv2cTunnelPool = NULL;

//------------------------------------------------------------------------------
nw_gtpv2c_tunnel_t* nwGtpv2cTunnelNew(
    struct nw_gtpv2c_stack_s* pStack, uint32_t teid,
    struct sockaddr* ipAddrRemote, nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel) {
  nw_gtpv2c_tunnel_t* thiz;

  if (gpGtpv2cTunnelPool) {
    thiz               = gpGtpv2cTunnelPool;
    gpGtpv2cTunnelPool = gpGtpv2cTunnelPool->next;
  } else {
    NW_GTPV2C_MALLOC(
        pStack, sizeof(nw_gtpv2c_tunnel_t), thiz, nw_gtpv2c_tunnel_t*);
  }

  if (thiz) {
    memset(thiz, 0, sizeof(nw_gtpv2c_tunnel_t));
    thiz->teid = teid;
    memcpy(
        (void*) &thiz->ipAddrRemote, ipAddrRemote,
        ipAddrRemote->sa_family == AF_INET ? sizeof(struct sockaddr_in) :
                                             sizeof(struct sockaddr_in6));
    thiz->hUlpTunnel = hUlpTunnel;
  }
  return thiz;
}

//------------------------------------------------------------------------------
nw_rc_t nwGtpv2cTunnelDelete(
    __attribute__((unused)) struct nw_gtpv2c_stack_s* pStack,
    nw_gtpv2c_tunnel_t* thiz) {
  thiz->next         = gpGtpv2cTunnelPool;
  gpGtpv2cTunnelPool = thiz;
  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t nwGtpv2cTunnelGetUlpTunnelHandle(
    nw_gtpv2c_tunnel_t* thiz, nw_gtpv2c_ulp_tunnel_handle_t* phUlpTunnel) {
  *phUlpTunnel = (thiz ? thiz->hUlpTunnel : 0x00000000);
  return NW_OK;
}

#ifdef __cplusplus
}
#endif

/*--------------------------------------------------------------------------*
                        E N D     O F    F I L E
  --------------------------------------------------------------------------*/
