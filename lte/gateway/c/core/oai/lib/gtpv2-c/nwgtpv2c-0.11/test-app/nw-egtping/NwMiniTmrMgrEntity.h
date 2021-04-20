/*----------------------------------------------------------------------------*
 *                                                                            *
 *         M I N I M A L I S T I C    T M R M G R     E N T I T Y             *
 *                                                                            *
 *                    Copyright (C) 2010 Amit Chawre.                         *
 *                                                                            *
 *----------------------------------------------------------------------------*/

/**
 * @file NwMiniTmrMgrEntity.c
 * @brief This file ontains example of a minimalistic timer manager entity.
 */

#include <stdio.h>
#include <assert.h>
#include "NwEvt.h"
#include "NwLog.h"

#ifndef NW_ASSERT
#define NW_ASSERT assert
#endif

#ifndef __NW_MINI_TMR_MGR_H__
#define __NW_MINI_TMR_MGR_H__

typedef struct {
  uint8_t _tbd;
} NwGtpv2cNodeTmrMgrT;

typedef struct {
  NwEventT ev;
  void* timeoutArg;
} NwGtpv2cNodeTmrT;

#ifdef __cplusplus
extern "C" {
#endif

nw_rc_t nwTimerStart(
    nw_gtpv2c_TimerMgrHandleT tmrMgrHandle, uint32_t timeoutSec,
    uint32_t timeoutUsec, uint32_t tmrType, void* timeoutArg,
    nw_gtpv2c_TimerHandleT* hTmr);

nw_rc_t nwTimerStop(
    nw_gtpv2c_TimerMgrHandleT tmrMgrHandle, nw_gtpv2c_TimerHandleT hTmr);

#ifdef __cplusplus
}
#endif

#endif
