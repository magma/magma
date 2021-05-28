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

#ifndef FILE_IE_TO_BYTES_SEEN
#define FILE_IE_TO_BYTES_SEEN

#include <stdbool.h>

#include "common_ies.h"
#include "common_types.h"
#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "TrackingAreaIdentity.h"

// value field length of IEs (IEI and length indicator are excluded)
#define IE_LENGTH_EPS_LOCATION_UPDATE_TYPE 1
#define IE_LENGTH_ECGI 7
#define IE_LENGTH_IMSI_DETACH_FROM_EPS_SERVICE_TYPE 1
#define IE_LENGTH_IMSI_DETACH_FROM_NON_EPS_SERVICE_TYPE 1
#define IE_LENGTH_LAI 5
#define IE_LENGTH_MME_NAME 55
#define IE_LENGTH_MOBILE_STATION_CLASSMARK2 3
#define IE_LENGTH_SERVICE_INDICATOR 1
#define IE_LENGTH_SGS_CAUSE 1
#define IE_LENGTH_TAI 5
#define IE_LENGTH_TMSI_STATUS 1
#define IE_LENGTH_UE_EMM_MODE 1
#define IE_LENGTH_UE_TIMEZONE 1
#define IE_LENGTH_PLMN 3

void tmsi_status_to_bytes(const bool* tmsi_status, char* byte_arr);
void tai_to_bytes(const tai_t* tai, char* byte_arr);
void lai_to_bytes(const lai_t* lai, char* byte_arr);
void ecgi_to_bytes(const ecgi_t* ecgi, char* byte_arr);
void mobile_station_classmark2_to_bytes(
    const MobileStationClassmark2_t* mscm2, char* byte_arr);
void plmn_to_bytes(const plmn_t* plmn, char* byte_arr);

#endif /* FILE_IE_TO_BYTES_SEEN */
