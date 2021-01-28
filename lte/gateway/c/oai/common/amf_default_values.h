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

#ifndef FILE_AMF_DEFAULT_VALUES_SEEN
#define FILE_AMF_DEFAULT_VALUES_SEEN

/*******************************************************************************
 * Timer Constants
 ******************************************************************************/
#define AMF_STATISTIC_TIMER_S (60)

/*******************************************************************************
 * NGAP Constants
 ******************************************************************************/

#define NGAP_PORT_NUMBER (38412)  ///< NGAP SCTP IANA ASSIGNED Port Number
#define NGAP_SCTP_PPID (60)  ///< NGAP SCTP Payload Protocol Identifier (PPID)

#define NGAP_OUTCOME_TIMER_DEFAULT (5)  ///< NGAP Outcome drop timer (s)

/*******************************************************************************
 * SCTP Constants
 ******************************************************************************/

#define SCTP_RECV_BUFFER_SIZE (1 << 16)
#define SCTP_OUT_STREAMS (32)
#define SCTP_IN_STREAMS (32)
#define SCTP_MAX_ATTEMPTS (5)

/*******************************************************************************
 * AMF global definitions
 ******************************************************************************/

#define AMFC (1)
#define AMFGID (1)
#define AMFPOINTER (1)
#define PLMN_MCC (208)
#define PLMN_MNC (34)
#define PLMN_MNC_LEN (2)
#define PLMN_TAC (1)

#define RELATIVE_CAPACITY (15)

/*******************************************************************************
 * GRPC Service Constants
 ******************************************************************************/
//#define GRPCSERVICES_SERVER_ADDRESS "127.0.0.1:50073"

#endif /* FILE_AMF_DEFAULT_VALUES_SEEN */
