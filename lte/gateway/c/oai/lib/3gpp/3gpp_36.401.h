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

/*! \file 3gpp_36.401.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_3GPP_36_401_SEEN
#define FILE_3GPP_36_401_SEEN

#include <stdint.h>

//------------------------------------------------------------------------------
// 6.2 E-UTRAN Identifiers
//------------------------------------------------------------------------------

typedef uint32_t
    enb_ue_s1ap_id_t; /*!< \brief  An eNB UE S1AP ID shall be allocated so as to
                         uniquely identify the UE over the S1 interface within
                         an eNB. When an MME receives an eNB UE S1AP ID it shall
                         store it for the duration of the UE-associated logical
                         S1-connection for this UE. Once known to an MME this IE
                         is included in all UE associated S1-AP signalling.
                                                                          The
                         eNB UE S1AP ID shall be unique within the eNB logical
                         node. */

typedef uint32_t
    mme_ue_s1ap_id_t; /*!< \brief  A MME UE S1AP ID shall be allocated so as to
                         uniquely identify the UE over the S1 interface within
                         the MME. When an eNB receives MME UE S1AP ID it shall
                         store it for the duration of the UE-associated logical
                         S1-connection for this UE. Once known to an eNB this IE
                         is included in all UE associated S1-AP signalling.
                                                                          The
                         MME UE S1AP ID shall be unique within the MME logical
                         node.*/

#endif /* FILE_3GPP_36_401_SEEN */
