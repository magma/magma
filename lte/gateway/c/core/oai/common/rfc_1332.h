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

/*! \file rfc_1332.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_RFC_1332_SEEN
#define FILE_RFC_1332_SEEN

// 2 A PPP Network Control Protocol (NCP) for IP

// Data Link Layer Protocol Field
// Exactly one IPCP packet is encapsulated in the Information field
// of PPP Data Link Layer frames where the Protocol field indicates
// type hex 8021 (IP Control Protocol)

// Code field
// Only Codes 1 through 7 (Configure-Request, Configure-Ack,
// Configure-Nak, Configure-Reject, Terminate-Request, Terminate-Ack
// and Code-Reject) are used. Other Codes should be treated as
// unrecognized and should result in Code-Rejects.
#define IPCP_CODE_CONFIGURE_REQUEST (0x01)
#define IPCP_CODE_CONFIGURE_ACK (0x02)
#define IPCP_CODE_CONFIGURE_NACK (0x03)
#define IPCP_CODE_CONFIGURE_REJECT (0x04)
#define IPCP_CODE_TERMINATE_REQUEST (0x05)
#define IPCP_CODE_TERMINATE_ACK (0x06)
#define IPCP_CODE_REJECT (0x07)

#endif
