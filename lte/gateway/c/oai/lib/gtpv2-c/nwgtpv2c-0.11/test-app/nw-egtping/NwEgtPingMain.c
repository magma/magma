/*----------------------------------------------------------------------------*
                      Copyright (C) 2010 Amit Chawre.
  ----------------------------------------------------------------------------*/

/**
   @file NwEgtPingMain.c
   @brief This is a program demostrating usage of nw-gtpv2c library for eGTP
   ping.
*/

#include <stdio.h>
#include <assert.h>
#include <signal.h>
#include <string.h>
#include "NwEvt.h"
#include "NwLog.h"
#include "NwGtpv2c.h"

#include "NwMiniLogMgrEntity.h"
#include "NwMiniTmrMgrEntity.h"
#include "NwMiniUdpEntity.h"
#include "NwMiniUlpEntity.h"

#ifndef NW_ASSERT
#define NW_ASSERT assert
#endif

static NwCharT* gLogLevelStr[] = {"EMER", "ALER", "CRIT", "ERRO",
                                  "WARN", "NOTI", "INFO", "DEBG"};

typedef struct NwEgtPingS {
  uint8_t localIpStr[20];
  uint8_t targetIpStr[20];
  uint32_t pingInterval;
  uint32_t pingCount;
  NwGtpv2cNodeUlpT ulpObj;
  NwGtpv2cNodeUdpT udpObj;
} NwEgtPingT;

static NwGtpv2cNodeUlpT ulpObj;
static NwGtpv2cNodeUdpT udpObj;

static NwEgtPingT egtPing;

void nwEgtPingHandleSignal(int sigNum) {
  printf(
      "\n--- %s (" NW_IPV4_ADDR ") EGTPING statistics --- ",
      egtPing.targetIpStr, NW_IPV4_ADDR_FORMAT(inet_addr(egtPing.targetIpStr)));
  printf(
      "\n%u requests sent, %u response received, %d%% packet loss \n\n",
      udpObj.packetsSent, udpObj.packetsRcvd,
      (udpObj.packetsSent ? 100 * (udpObj.packetsSent - udpObj.packetsRcvd) /
                                udpObj.packetsSent :
                            0));
  exit(sigNum);
}

nw_rc_t nwEgtPingHelp() {
  printf("Usage: egtping [-i interval] [-c count] [-l local-ip] ");
  printf("\n               [-t3 t3-time] [-n3 n3-count] destination");
  printf("\n");
  printf(
      "\n       -i <interval>     : Interval between two echo request "
      "messages. "
      "(Default: 1 sec)");
  printf(
      "\n       -c <count>        : Stop after sending count pings. (Default: "
      "Infinite)");
  printf("\n       -t <t3-time>      : GTP T3 timeout value. (Default: 2 sec)");
  printf("\n       -n <n3-count>     : GTP N3 count value. (Default: 2 sec)");
  printf(
      "\n       -l <local-ip>     : Local IP adddress to use. (Default: All "
      "local IPs)");
  printf("\n       -h                : Show this message.");
  printf("\n");
  printf("\n");
}

nw_rc_t nwEgtPingParseCmdLineOpts(int argc, char* argv[]) {
  nw_rc_t rc = NW_OK;
  int i      = 0;

  i++;
  egtPing.pingInterval = 1;
  egtPing.pingCount    = 0xffffffff;

  if (argc < 2) return NW_FAILURE;

  if ((argc == 2) &&
      ((strcmp("--help", argv[i]) == 0) || (strcmp(argv[i], "-h") == 0)))
    return NW_FAILURE;

  while (i < argc - 1) {
    NW_LOG(NW_LOG_LEVEL_DEBG, "Processing cmdline arg %s", argv[i]);

    if ((strcmp("--local-ip", argv[i]) == 0) || (strcmp(argv[i], "-l") == 0)) {
      i++;

      if (i >= (argc - 1)) return NW_FAILURE;

      strcpy(egtPing.localIpStr, (argv[i]));
    } else if (
        (strcmp("--interval", argv[i]) == 0) || (strcmp(argv[i], "-i") == 0)) {
      i++;

      if (i >= (argc - 1)) return NW_FAILURE;

      egtPing.pingInterval = atoi(argv[i]);
    } else if (
        (strcmp("--count", argv[i]) == 0) || (strcmp(argv[i], "-c") == 0)) {
      i++;

      if (i >= (argc - 1)) return NW_FAILURE;

      egtPing.pingCount = atoi(argv[i]);
    } else if (
        (strcmp("--help", argv[i]) == 0) || (strcmp(argv[i], "-h") == 0)) {
      rc = NW_FAILURE;
    } else {
      return NW_FAILURE;
    }

    i++;
  }

  strcpy(egtPing.targetIpStr, (argv[i]));
  return rc;
}

/*---------------------------------------------------------------------------
                  T H E      M A I N      F U N C T I O N
  --------------------------------------------------------------------------*/

int main(int argc, char* argv[]) {
  nw_rc_t rc;
  uint32_t logLevel;
  uint8_t* logLevelStr;
  nw_gtpv2c_StackHandleT hGtpv2cStack = 0;
  nw_gtpv2c_ulp_entity_t ulp;
  nw_gtpv2c_udp_entity_t udp;
  nw_gtpv2c_timer_mgr_entity_t tmrMgr;
  nw_gtpv2c_log_mgr_entity_t logMgr;

  printf("EGTPING 0.1, Copyright (C) 2011 Amit Chawre.\n");
  rc = nwEgtPingParseCmdLineOpts(argc, argv);

  if (rc != NW_OK) {
    rc = nwEgtPingHelp();
    exit(rc);
  }

  logLevelStr = getenv("NW_LOG_LEVEL");

  if (logLevelStr == NULL) {
    logLevel = NW_LOG_LEVEL_INFO;
  } else {
    if (strncmp(logLevelStr, "EMER", 4) == 0)
      logLevel = NW_LOG_LEVEL_EMER;
    else if (strncmp(logLevelStr, "ALER", 4) == 0)
      logLevel = NW_LOG_LEVEL_ALER;
    else if (strncmp(logLevelStr, "CRIT", 4) == 0)
      logLevel = NW_LOG_LEVEL_CRIT;
    else if (strncmp(logLevelStr, "ERRO", 4) == 0)
      logLevel = NW_LOG_LEVEL_ERRO;
    else if (strncmp(logLevelStr, "WARN", 4) == 0)
      logLevel = NW_LOG_LEVEL_WARN;
    else if (strncmp(logLevelStr, "NOTI", 4) == 0)
      logLevel = NW_LOG_LEVEL_NOTI;
    else if (strncmp(logLevelStr, "INFO", 4) == 0)
      logLevel = NW_LOG_LEVEL_INFO;
    else if (strncmp(logLevelStr, "DEBG", 4) == 0)
      logLevel = NW_LOG_LEVEL_DEBG;
  }

  /*---------------------------------------------------------------------------
      Initialize event library
    --------------------------------------------------------------------------*/
  NW_EVT_INIT();
  /*---------------------------------------------------------------------------
      Initialize Log Manager
    --------------------------------------------------------------------------*/
  nwMiniLogMgrInit(nwMiniLogMgrGetInstance(), logLevel);
  /*---------------------------------------------------------------------------
      Initialize Gtpv2c Stack Instance
    --------------------------------------------------------------------------*/
  rc = nwGtpv2cInitialize(&hGtpv2cStack);

  if (rc != NW_OK) {
    NW_LOG(
        NW_LOG_LEVEL_ERRO,
        "Failed to create gtpv2c stack instance. Error '%u' occured", rc);
    exit(1);
  }

  rc = nwGtpv2cSetLogLevel(hGtpv2cStack, logLevel);
  /*---------------------------------------------------------------------------
     Set up Ulp Entity
    --------------------------------------------------------------------------*/
  rc = nwGtpv2cUlpInit(&ulpObj, hGtpv2cStack, egtPing.localIpStr);
  NW_ASSERT(NW_OK == rc);
  ulp.hUlp           = (nw_gtpv2c_UlpHandleT) &ulpObj;
  ulp.ulpReqCallback = nwGtpv2cUlpProcessStackReqCallback;
  rc                 = nwGtpv2cSetUlpEntity(hGtpv2cStack, &ulp);
  NW_ASSERT(NW_OK == rc);
  /*---------------------------------------------------------------------------
     Set up Udp Entity
    --------------------------------------------------------------------------*/
  rc = nwGtpv2cUdpInit(&udpObj, hGtpv2cStack, egtPing.localIpStr);
  NW_ASSERT(NW_OK == rc);
  udp.hUdp               = (nw_gtpv2c_UdpHandleT) &udpObj;
  udp.udpDataReqCallback = nwGtpv2cUdpDataReq;
  rc                     = nwGtpv2cSetUdpEntity(hGtpv2cStack, &udp);
  NW_ASSERT(NW_OK == rc);
  /*---------------------------------------------------------------------------
     Set up Log Entity
    --------------------------------------------------------------------------*/
  tmrMgr.tmrMgrHandle     = 0;
  tmrMgr.tmrStartCallback = nwTimerStart;
  tmrMgr.tmrStopCallback  = nwTimerStop;
  rc                      = nwGtpv2cSetTimerMgrEntity(hGtpv2cStack, &tmrMgr);
  NW_ASSERT(NW_OK == rc);
  /*---------------------------------------------------------------------------
     Set up Log Entity
    --------------------------------------------------------------------------*/
  logMgr.logMgrHandle   = (nw_gtpv2c_LogMgrHandleT) nwMiniLogMgrGetInstance();
  logMgr.logReqCallback = nwMiniLogMgrLogRequest;
  rc                    = nwGtpv2cSetLogMgrEntity(hGtpv2cStack, &logMgr);
  NW_ASSERT(NW_OK == rc);
  /*---------------------------------------------------------------------------
      Send Message Request to Gtpv2c Stack Instance
    --------------------------------------------------------------------------*/
  NW_LOG(
      NW_LOG_LEVEL_NOTI, "EGTPING %s (" NW_IPV4_ADDR ")", egtPing.targetIpStr,
      NW_IPV4_ADDR_FORMAT(inet_addr(egtPing.targetIpStr)));
  rc = nwGtpv2cUlpPing(
      &ulpObj, inet_addr(egtPing.targetIpStr), egtPing.pingCount,
      egtPing.pingInterval, 2, 3);
  NW_ASSERT(NW_OK == rc);
  /*---------------------------------------------------------------------------
     Install signal handler
    --------------------------------------------------------------------------*/
  signal(SIGINT, nwEgtPingHandleSignal);
  /*---------------------------------------------------------------------------
     Event loop
    --------------------------------------------------------------------------*/
  NW_EVT_LOOP();
  NW_LOG(NW_LOG_LEVEL_ERRO, "Exit from eventloop, no events to process!");
  /*---------------------------------------------------------------------------
      Destroy Gtpv2c Stack Instance
    --------------------------------------------------------------------------*/
  rc = nwGtpv2cFinalize(hGtpv2cStack);

  if (rc != NW_OK) {
    NW_LOG(
        NW_LOG_LEVEL_ERRO,
        "Failed to finalize gtpv2c stack instance. Error '%u' occured", rc);
  }

  return rc;
}
