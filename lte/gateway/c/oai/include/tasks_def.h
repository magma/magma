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
// This task is mandatory and must always be placed in first position
TASK_DEF(TASK_MAIN)

// Other possible tasks in the process

/// FW_IP task
TASK_DEF(TASK_FW_IP)
/// MME Applicative task
TASK_DEF(TASK_MME_APP)
/// S10 task
TASK_DEF(TASK_S10)
/// S11 task
TASK_DEF(TASK_S11)
/// S1AP task
TASK_DEF(TASK_S1AP)
/// S6a task
TASK_DEF(TASK_S6A)
/// SCTP task
TASK_DEF(TASK_SCTP)
/// Serving and Proxy Gateway Application task
TASK_DEF(TASK_SPGW_APP)
/// UDP task
TASK_DEF(TASK_UDP)
// LOGGING TXT TASK
TASK_DEF(TASK_LOG)
// GENERAL PURPOSE SHARED LOGGING TASK
TASK_DEF(TASK_SHARED_TS_LOG)
// UTILITY TASK FOR SYSTEM() CALLS
TASK_DEF(TASK_ASYNC_SYSTEM)
// SERVICE303 TASK
TASK_DEF(TASK_SERVICE303)
TASK_DEF(TASK_SERVICE303_SERVER)
/// SGs task
TASK_DEF(TASK_SGS)
/// SMS_ORC8R task
TASK_DEF(TASK_SMS_ORC8R)
/// GRPC service task for SGs, S6a, SPGW, HA
TASK_DEF(TASK_GRPC_SERVICE)
/// HA task
TASK_DEF(TASK_HA)
/// NGAP task
TASK_DEF(TASK_NGAP)
/// AMF Application task
TASK_DEF(TASK_AMF_APP)
