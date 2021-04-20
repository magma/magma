/*----------------------------------------------------------------------------*
 *                                                                            *
 *            M I N I M A L I S T I C     U L P     E N T I T Y               *
 *                                                                            *
 *                    Copyright (C) 2010 Amit Chawre.                         *
 *                                                                            *
 *----------------------------------------------------------------------------*/

/**
 * @file NwMiniUlpEntity.h
 * @brief This file contains example of a minimalistic ULP entity.
 */

#include <stdio.h>
#include <assert.h>
#include "NwEvt.h"
#include "NwLog.h"

#ifndef __NW_MINI_ULP_H__
#define __NW_MINI_ULP_H__

typedef struct {
  uint8_t peerIpStr[16];
  uint32_t restartCounter;
  nw_gtpv2c_StackHandleT hGtpv2cStack;
} NwGtpv2cNodeUlpT;

#ifdef __cplusplus
extern "C" {
#endif

nw_rc_t nwGtpv2cUlpInit(
    NwGtpv2cNodeUlpT* thiz, nw_gtpv2c_StackHandleT hGtpv2cStack,
    char* peerIpStr);

nw_rc_t nwGtpv2cUlpDestroy(NwGtpv2cNodeUlpT* thiz);

nw_rc_t nwGtpv2cUlpCreateSessionRequestToPeer(NwGtpv2cNodeUlpT* thiz);

nw_rc_t nwGtpv2cUlpProcessStackReqCallback(
    nw_gtpv2c_UlpHandleT hUlp, nw_gtpv2c_ulp_api_t* pUlpApi);

#ifdef __cplusplus
}
#endif

#endif
