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

/*! \file 3gpp_29.018.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_3GPP_29_018_SEEN
#define FILE_3GPP_29_018_SEEN

//..............................................................................
// 9.2 SGSAP Message types
//..............................................................................

typedef enum sgsap_message_types_e {
  SGS_PAGING_REQUEST             = 0x01,
  SGS_PAGING_REJECT              = 0x02,
  SGS_SERVICE_REQUEST            = 0x06,
  SGS_DOWNLINK_UNIT_DATA         = 0x07,
  SGS_UPLINK_UNIT_DATA           = 0x08,
  SGS_LOCATION_UPDATE_REQUEST    = 0x09,
  SGS_LOCATION_UPDATE_ACCEPT     = 0x0A,
  SGS_LOCATION_UPDATE_REJECT     = 0x0B,
  SGS_TMSI_REALLOCATION_COMPLETE = 0x0C,
  SGS_ALERT_REQUEST              = 0x0D,
  SGS_ALERT_ACK                  = 0x0E,
  SGS_ALERT_REJECT               = 0x0F,
  SGS_UE_ACTIVITY_INDICATION     = 0x10,
  SGS_EPS_DETACH_INDICATION      = 0x11,
  SGS_EPS_DETACH_ACK             = 0x12,
  SGS_IMSI_DETACH_INDICATION     = 0x13,
  SGS_IMSI_EPS_DETACH_ACK        = 0x14,
  SGS_RESET_INDICATION           = 0x15,
  SGS_RESET_ACK                  = 0x16,
  SGS_SERVICE_ABORT_REQUEST      = 0x17,
  SGS_MO_CSFB_INDICATION         = 0x18,
  SGS_MM_INFORMATION_REQUEST     = 0x1A,
  SGS_RELEASE_REQUEST            = 0x1B,
  SGS_STATUS                     = 0x1D,
  SGS_UE_UNREACHABLE             = 0x1F,
} sgsap_message_types_t;

#endif /* FILE_3GPP_29_018_SEEN */
