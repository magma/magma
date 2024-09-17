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

/*! \file gtpv1_u_messages_def.hpp
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

// "#pragma once" will not work for this file, because this file is included
// in include/messages_def.h, which is included more than once within enum
// and structure in the file intertask_interface_types.h
// See comment in "lte/gateway/c/core/oai/include/messages_def.h" for details

MESSAGE_DEF(GTPV1U_CREATE_TUNNEL_REQ, Gtpv1uCreateTunnelReq,
            gtpv1uCreateTunnelReq)
MESSAGE_DEF(GTPV1U_CREATE_TUNNEL_RESP, Gtpv1uCreateTunnelResp,
            gtpv1uCreateTunnelResp)
MESSAGE_DEF(GTPV1U_UPDATE_TUNNEL_REQ, Gtpv1uUpdateTunnelReq,
            gtpv1uUpdateTunnelReq)
MESSAGE_DEF(GTPV1U_UPDATE_TUNNEL_RESP, Gtpv1uUpdateTunnelResp,
            gtpv1uUpdateTunnelResp)
MESSAGE_DEF(GTPV1U_DELETE_TUNNEL_REQ, Gtpv1uDeleteTunnelReq,
            gtpv1uDeleteTunnelReq)
MESSAGE_DEF(GTPV1U_DELETE_TUNNEL_RESP, Gtpv1uDeleteTunnelResp,
            gtpv1uDeleteTunnelResp)
MESSAGE_DEF(GTPV1U_TUNNEL_DATA_IND, Gtpv1uTunnelDataInd, gtpv1uTunnelDataInd)
MESSAGE_DEF(GTPV1U_TUNNEL_DATA_REQ, Gtpv1uTunnelDataReq, gtpv1uTunnelDataReq)
