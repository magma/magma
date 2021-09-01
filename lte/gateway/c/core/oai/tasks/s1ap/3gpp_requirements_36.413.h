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

/*! \file 3gpp_requirements_36.413.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_3GPP_REQUIREMENTS_36_413_SEEN
#define FILE_3GPP_REQUIREMENTS_36_413_SEEN

#include "3gpp_requirements.h"
#include "log.h"

#define REQUIREMENT_3GPP_36_413(rElEaSe_sEcTiOn__OaImark)                      \
  REQUIREMENT_3GPP_SPEC(                                                       \
      LOG_S1AP, "Hit 3GPP TS 36_413" #rElEaSe_sEcTiOn__OaImark                 \
                " : " rElEaSe_sEcTiOn__OaImark##_BRIEF "\n")
#define NO_REQUIREMENT_3GPP_36_413(rElEaSe_sEcTiOn__OaImark)                   \
  REQUIREMENT_3GPP_SPEC(                                                       \
      LOG_S1AP, "#NOT IMPLEMENTED 3GPP TS 36_413" #rElEaSe_sEcTiOn__OaImark    \
                " : " rElEaSe_sEcTiOn__OaImark##_BRIEF "\n")
#define NOT_REQUIREMENT_3GPP_36_413(rElEaSe_sEcTiOn__OaImark)                  \
  REQUIREMENT_3GPP_SPEC(                                                       \
      LOG_S1AP, "#NOT ASSERTED 3GPP TS 36_413" #rElEaSe_sEcTiOn__OaImark       \
                " : " rElEaSe_sEcTiOn__OaImark##_BRIEF "\n")

//-----------------------------------------------------------------------------------------------------------------------
#define R10_8_3_3_2__2                                                         \
  "MME36.413R10_8.3.3.2_2: Successful Operation\
                                                                                                                        \
    The UE CONTEXT RELEASE COMMAND message shall contain the UE S1AP ID pair IE if available, otherwise the             \
    message shall contain the MME UE S1AP ID IE."

#define R10_8_3_3_2__2_BRIEF                                                   \
  "UE CONTEXT RELEASE COMMAND contains UE S1AP ID pair IE or at least MME UE " \
  "S1AP ID IE"

#endif /* FILE_3GPP_REQUIREMENTS_36_413_SEEN */
