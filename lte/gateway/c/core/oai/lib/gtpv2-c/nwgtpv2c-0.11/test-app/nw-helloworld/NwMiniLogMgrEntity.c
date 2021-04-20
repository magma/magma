/*----------------------------------------------------------------------------*
 *                                                                            *
           M I N I M A L I S T I C    L O G M G R     E N T I T Y
 *                                                                            *
                      Copyright (C) 2010 Amit Chawre.
 *                                                                            *
  ----------------------------------------------------------------------------*/

/**
   @file hello-world.c
   @brief This file contains example of a minimalistic log manager entity.
*/

#include <stdio.h>
#include <assert.h>
#include "NwEvt.h"
#include "NwTypes.h"
#include "NwError.h"
#include "NwLog.h"
#include "NwGtpv2c.h"

#include "NwMiniLogMgrEntity.h"

#ifdef __cplusplus
extern "C" {
#endif

static NwCharT* gLogLevelStr[] = {"EMER", "ALER", "CRIT", "ERRO",
                                  "WARN", "NOTI", "INFO", "DEBG"};

NwMiniLogMgrT __gLogMgr;

/*---------------------------------------------------------------------------
   Public functions
  --------------------------------------------------------------------------*/

NwMiniLogMgrT* nwMiniLogMgrGetInstance() {
  return &(__gLogMgr);
}

nw_rc_t nwMiniLogMgrInit(NwMiniLogMgrT* thiz, uint32_t logLevel) {
  thiz->logLevel = logLevel;
  return NW_OK;
}

nw_rc_t nwMiniLogMgrSetLogLevel(NwMiniLogMgrT* thiz, uint32_t logLevel) {
  thiz->logLevel = logLevel;
}

nw_rc_t nwMiniLogMgrLogRequest(
    nw_gtpv2c_LogMgrHandleT hLogMgr, uint32_t logLevel, NwCharT* file,
    uint32_t line, NwCharT* logStr) {
  NwMiniLogMgrT* thiz = (NwMiniLogMgrT*) hLogMgr;

  if (thiz->logLevel >= logLevel)
    printf(
        "NWGTPV2C-STK  %s - %s <%s,%u>\n", gLogLevelStr[logLevel], logStr,
        basename(file), line);

  return NW_OK;
}

#ifdef __cplusplus
}
#endif
