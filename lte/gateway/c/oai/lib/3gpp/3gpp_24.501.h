/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once
#include <stdint.h>
#include <stdbool.h>

//==============================================================================
// 9 General message format and information elements coding
//==============================================================================

//------------------------------------------------------------------------------
// 9.2 Protocol discriminator
//------------------------------------------------------------------------------

// 9.3.1 Security header type
#define SECURITY_HEADER_TYPE_NOT_PROTECTED 0b0000
#define SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED 0b0001
#define SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED 0b0010
#define SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_NEW 0b0011
#define SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED_NEW 0b0100
#define SECURITY_HEADER_TYPE_SERVICE_REQUEST 0b1100
#define SECURITY_HEADER_TYPE_RESERVED1 0b1101
#define SECURITY_HEADER_TYPE_RESERVED2 0b1110
#define SECURITY_HEADER_TYPE_RESERVED3 0b1111

//------------------------------------------------------------------------------
// 9.8 Message type
//------------------------------------------------------------------------------
// Table 9.8.1: Message types for EPS mobility management
/* Message identifiers for EPS Mobility Management Registration  REGISTRATION */
//#define REGISTRATION_REQUEST   0b01000001                  /* 65 = 0x41 */
//#define REGISTRATION_ACCEPT    0b01000010                  /* 66 = 0x42 */
//#define REGISTRATION_COMPLETE  0b01000011                  /* 67 = 0x43 */
//#define REGISTRATION_REJECT    0b01000100                  /* 68 = 0x44 */
//#define DEREGISTRATION_REQUEST_UE_INIT 0b01000101                  /* 69 =
// 0x45 */ #define DEREGISTRATION_ACCEPT_UE_INIT  0b01000110                  /*
// 70 = 0x46 */
#define DEREGISTRATION_REQUEST_NW_INIT 0b01000111 /* 71 = 0x47 */
#define DEREGISTRATION_ACCEPT_NW_INIT 0b01001000  /* 72 = 0x48 */
#define DEREGISTRATION_ACCEPT_UE_INIT 0b01000110  /* 70 = 0x46 */

//..............................................................................
//  Table 10.2.1: Timers of 5GS mobility management – UE side
//..............................................................................

#define T3502_DEFAULT_VALUE 720
#define T3510_DEFAULT_VALUE 15
#define T3511_DEFAULT_VALUE 10
#define T3512_DEFAULT_VALUE 3240
#define T3516_DEFAULT_VALUE 30
#define T3517_DEFAULT_VALUE 5
#define T3517_EXT_DEFAULT_VALUE 10
#define T3520_DEFAULT_VALUE 15
#define T3521_DEFAULT_VALUE 15
#define T3523_DEFAULT_VALUE 0  // value provided by network
#define T3540_DEFAULT_VALUE 10
#define T3542_DEFAULT_VALUE 0  // value provided by network
#define T3522_DEFAULT_VALUE 6
#define T3550_DEFAULT_VALUE 6
#define T3560_DEFAULT_VALUE 6
#define T3570_DEFAULT_VALUE 6
#define T3585_DEFAULT_VALUE 0  // Value provided by NW.
#define T3586_DEFAULT_VALUE 0  // TODO-RECHECK
#define T3589_DEFAULT_VALUE 0  // TODO-RECHECK
#define T3595_DEFAULT_VALUE 0

/*
 Registration request
Registration accept
Registration complete
Registration reject
Deregistration request (UE originating)
Deregistration accept (UE originating)
Deregistration request (UE terminated)
Deregistration accept (UE terminated) */
//-------------------------------------------------------------------------------
// A.3 Causes related to PLMN specific network failures and
//	   congestion/authentication failures
//-------------------------------------------------------------------------------
#define AMF_CAUSE_SMF_FAILURE 19  // need to check
#define AMF_CAUSE_MAC_FAILURE 20
#define AMF_CAUSE_SYNCH_FAILURE 21
#define AMF_CAUSE_CONGESTION 22
#define AMF_UE_SECURITY_CAPABILITIES_MISMATCH                                  \
  23  // UE security capabilities mismatch
#define AMF_SECURITY_MODE_REJECT 24

/* 5G Timer structure */
typedef struct nas5g_timer_s {
  long int id;  /* The timer identifier                 */
  uint32_t sec; /* The timer interval value in seconds  */
} nas5g_timer_t;

#endif
