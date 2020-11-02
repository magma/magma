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
/*! \file sctp_messages_def.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

// WARNING: Do not include this header directly. Use intertask_interface.h
// instead.

MESSAGE_DEF(SCTP_INIT_MSG, sctp_init_t, sctpInit)
MESSAGE_DEF(SCTP_DATA_REQ, sctp_data_req_t, sctp_data_req)
MESSAGE_DEF(SCTP_DATA_IND, sctp_data_ind_t, sctp_data_ind)
MESSAGE_DEF(SCTP_DATA_CNF, sctp_data_cnf_t, sctp_data_cnf)
MESSAGE_DEF(SCTP_NEW_ASSOCIATION, sctp_new_peer_t, sctp_new_peer)
MESSAGE_DEF(
    SCTP_CLOSE_ASSOCIATION, sctp_close_association_t, sctp_close_association)
MESSAGE_DEF(
    SCTP_MME_SERVER_INITIALIZED, sctp_mme_server_initialized_t,
    sctp_mme_server_initialized)
MESSAGE_DEF(
    SCTP_AMF_SERVER_INITIALIZED, sctp_amf_server_initialized_t,
    sctp_amf_server_initialized)
