/*----------------------------------------------------------------------------*
 *                                                                            *
           M I N I M A L I S T I C    T M R M G R     E N T I T Y
 *                                                                            *
                      Copyright (C) 2010 Amit Chawre.
 *                                                                            *
  ----------------------------------------------------------------------------*/

/**
   @file NwMiniTmrMgrEntity.c
   @brief This file ontains example of a minimalistic timer manager entity.
*/

#include <stdio.h>
#include <assert.h>
#include "NwEvt.h"
#include "NwGtpv2c.h"
#include "NwMiniLogMgrEntity.h"
#include "NwMiniTmrMgrEntity.h"

#ifndef NW_ASSERT
#define NW_ASSERT assert
#endif

#ifdef __cplusplus
extern "C" {
#endif

static NwCharT* gLogLevelStr[] = {"EMER", "ALER", "CRIT", "ERRO",
                                  "WARN", "NOTI", "INFO", "DEBG"};

/*---------------------------------------------------------------------------
   Private functions
  --------------------------------------------------------------------------*/

static void NW_TMR_CALLBACK(nwGtpv2cNodeHandleStackTimerTimeout) {
  nw_rc_t rc;
  NwGtpv2cNodeTmrT* pTmr = (NwGtpv2cNodeTmrT*) arg;

  /*
   * Send Timeout Request to Gtpv2c Stack Instance
   */
  rc = nwGtpv2cProcessTimeout(pTmr->timeoutArg);
  NW_ASSERT(NW_OK == rc);
  free(pTmr);
  return;
}

/*---------------------------------------------------------------------------
   Public functions
  --------------------------------------------------------------------------*/

nw_rc_t nwTimerStart(
    nw_gtpv2c_TimerMgrHandleT tmrMgrHandle, uint32_t timeoutSec,
    uint32_t timeoutUsec, uint32_t tmrType, void* timeoutArg,
    nw_gtpv2c_TimerHandleT* hTmr) {
  nw_rc_t rc = NW_OK;
  NwGtpv2cNodeTmrT* pTmr;
  struct timeval tv;

  pTmr = (NwGtpv2cNodeTmrT*) malloc(sizeof(NwGtpv2cNodeTmrT));
  /*
   * set the timevalues
   */
  timerclear(&tv);
  tv.tv_sec        = timeoutSec;
  tv.tv_usec       = timeoutUsec;
  pTmr->timeoutArg = timeoutArg;
  evtimer_set(&pTmr->ev, nwGtpv2cNodeHandleStackTimerTimeout, pTmr);
  /*
   * add event
   */
  event_add(&(pTmr->ev), &tv);
  *hTmr = (nw_gtpv2c_TimerHandleT) pTmr;
  return rc;
}

nw_rc_t nwTimerStop(
    nw_gtpv2c_TimerMgrHandleT tmrMgrHandle, nw_gtpv2c_TimerHandleT hTmr) {
  evtimer_del(&(((NwGtpv2cNodeTmrT*) hTmr)->ev));
  free((void*) hTmr);
  return NW_OK;
}

#ifdef __cplusplus
}
#endif
