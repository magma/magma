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

/*! \file 3gpp_33.401.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_3GPP_33_401_SEEN
#define FILE_3GPP_33_401_SEEN

//------------------------------------------------------------------------------
// 5.1.3.2 Algorithm Identifier Values
//------------------------------------------------------------------------------
#define EEA0_ALG_ID 0b000
#define EEA1_128_ALG_ID 0b001
#define EEA2_128_ALG_ID 0b010

//------------------------------------------------------------------------------
// 5.1.4.2 Algorithm Identifier Values
//------------------------------------------------------------------------------
#define EIA0_ALG_ID 0b000
#define EIA1_128_ALG_ID 0b001
#define EIA2_128_ALG_ID 0b010

//------------------------------------------------------------------------------
// 6.1.2 Distribution of authentication data from HSS to serving network
//------------------------------------------------------------------------------
/* NOTE 2: It is recommended that the MME fetch only one EPS authentication
 * vector at a time as the need to perform AKA runs has been reduced in EPS
 * through the use of a more elaborate key hierarchy. In particular, service
 * requests can be authenticated using a stored K ASME without the need to
 * perform AKA. Furthermore, the sequence number management schemes in
 * TS 33.102, Annex C [4], designed to avoid re-synchronisation problems caused
 * by interleaving use of batches of authentication vectors, are only optional.
 * Re-synchronisation problems in EPS can be avoided, independently of the
 * sequence number management scheme, by immediately using an authentication
 * vector retrieved from the HSS in an authentication procedure between UE and
 * MME.
 */
#define MAX_EPS_AUTH_VECTORS 1

#endif /* FILE_3GPP_33_401_SEEN */
