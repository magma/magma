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

#ifndef __NW_GTPV2C_MSG_PARSE_INFO_H__
#define __NW_GTPV2C_MSG_PARSE_INFO_H__

#include "NwTypes.h"
#include "NwGtpv2c.h"

/**
 * @file NwGtpv2cMsgParseInfo.h
 * @brief This file defines APIs for to parse incoming messages.
 */

typedef struct nw_gtpv2c_grouped_ie_parse_info_s {
  uint8_t groupedIeType;
  uint16_t mandatoryIeCount;
  nw_gtpv2c_stack_handle_t hStack;

  struct {
    uint8_t ieMinLength;
    uint8_t iePresence;
  } ieParseInfo[NW_GTPV2C_IE_TYPE_MAXIMUM][NW_GTPV2C_IE_INSTANCE_MAXIMUM];

} nw_gtpv2c_grouped_ie_parse_info_t;

typedef struct nw_gtpv2c_msg_ie_parse_info_s {
  uint16_t msgType;
  uint16_t mandatoryIeCount;
  nw_gtpv2c_stack_handle_t hStack;

  struct {
    uint8_t ieMinLength;
    uint8_t iePresence;
    nw_gtpv2c_grouped_ie_parse_info_t* pGroupedIeInfo;
  } ieParseInfo[NW_GTPV2C_IE_TYPE_MAXIMUM][NW_GTPV2C_IE_INSTANCE_MAXIMUM];

} nw_gtpv2c_msg_ie_parse_info_t;

#ifdef __cplusplus
extern "C" {
#endif

nw_gtpv2c_msg_ie_parse_info_t* nwGtpv2cMsgIeParseInfoNew(
    nw_gtpv2c_stack_handle_t hStack, uint8_t msgType);

nw_rc_t nwGtpv2cMsgIeParseInfoDelete(nw_gtpv2c_msg_ie_parse_info_t* thiz);

nw_rc_t nwGtpv2cMsgIeParse(
    NW_IN nw_gtpv2c_msg_ie_parse_info_t* thiz,
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_INOUT nw_gtpv2c_error_t* pError);

#ifdef __cplusplus
}
#endif

#endif /* __NW_GTPV2C_MSG_PARSE_INFO_H__ */
/*--------------------------------------------------------------------------*
 *                      E N D     O F    F I L E                            *
 *--------------------------------------------------------------------------*/
