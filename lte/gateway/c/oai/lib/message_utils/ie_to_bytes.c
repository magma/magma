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

#include <stdbool.h>

#include "ie_to_bytes.h"

// 18.4.24 in 3GPP TS 29.018
void tmsi_status_to_bytes(const bool* tmsi_status, char* byte_arr) {
  if (*tmsi_status) {
    byte_arr[0] = 0x01;
  } else {
    byte_arr[0] = 0x00;
  }
  return;
}

// 9.9.3.32 in 3GPP TS 24.301
void tai_to_bytes(const tai_t* tai, char* byte_arr) {
  byte_arr[0] = (tai->plmn.mcc_digit2 << 4) | (tai->plmn.mcc_digit1 & 0x0F);
  byte_arr[1] = (tai->plmn.mnc_digit3 << 4) | (tai->plmn.mcc_digit3 & 0x0F);
  byte_arr[2] = (tai->plmn.mnc_digit2 << 4) | (tai->plmn.mnc_digit1 & 0x0F);
  byte_arr[3] = tai->tac >> 8;
  byte_arr[4] = tai->tac;
  return;
}

// 10.5.1.3 in 3GPP TS 24.008
void lai_to_bytes(const lai_t* lai, char* byte_arr) {
  /*byte_arr[0] = (lai->plmn.mcc_digit2 << 4) | (lai->plmn.mcc_digit1 & 0x0F);
  byte_arr[1] = (lai->plmn.mnc_digit3 << 4) | (lai->plmn.mcc_digit3 & 0x0F);
  byte_arr[2] = (lai->plmn.mnc_digit2 << 4) | (lai->plmn.mnc_digit1 & 0x0F);*/
  byte_arr[0] = (lai->mccdigit2 << 4) | (lai->mccdigit1 & 0x0F);
  byte_arr[1] = (lai->mncdigit3 << 4) | (lai->mccdigit3 & 0x0F);
  byte_arr[2] = (lai->mncdigit2 << 4) | (lai->mncdigit1 & 0x0F);
  byte_arr[3] = lai->lac >> 8;
  byte_arr[4] = lai->lac;
  return;
}

// 8.21.5 in 3GPP TS 29.274
void ecgi_to_bytes(const ecgi_t* ecgi, char* byte_arr) {
  byte_arr[0] = (ecgi->plmn.mcc_digit2 << 4) | (ecgi->plmn.mcc_digit1 & 0x0F);
  byte_arr[1] = (ecgi->plmn.mnc_digit3 << 4) | (ecgi->plmn.mcc_digit3 & 0x0F);
  byte_arr[2] = (ecgi->plmn.mnc_digit2 << 4) | (ecgi->plmn.mnc_digit1 & 0x0F);
  byte_arr[3] = ecgi->cell_identity.enb_id >> 16 & 0x0F;
  byte_arr[4] = ecgi->cell_identity.enb_id >> 8;
  byte_arr[5] = ecgi->cell_identity.enb_id;
  byte_arr[6] = ecgi->cell_identity.cell_id;
  return;
}

// 10.5.1.6 in 3GPP TS 24.008
void mobile_station_classmark2_to_bytes(
    const MobileStationClassmark2_t* mscm2, char* byte_arr) {
  byte_arr[0] = (mscm2->revisionlevel << 5) | (mscm2->esind << 4) |
                (mscm2->a51 << 3) | mscm2->rfpowercapability;
  byte_arr[1] = (mscm2->pscapability << 6) | (mscm2->ssscreenindicator << 4) |
                (mscm2->smcapability << 3) | (mscm2->vbs << 2) |
                (mscm2->vgcs << 1) | mscm2->fc;
  byte_arr[2] = (mscm2->cm3 << 7) | (mscm2->lcsvacap << 5) |
                (mscm2->ucs2 << 4) | (mscm2->solsa << 3) | (mscm2->cmsp << 2) |
                (mscm2->a53 << 1) | mscm2->a52;
  return;
}

// 7.3.9 in 3GPP TS 29.272
void plmn_to_bytes(const plmn_t* plmn, char* byte_arr) {
  byte_arr[0] = (plmn->mcc_digit2 << 4) | (plmn->mcc_digit1 & 0x0F);
  byte_arr[1] = (plmn->mnc_digit3 << 4) | (plmn->mcc_digit3 & 0x0F);
  byte_arr[2] = (plmn->mnc_digit2 << 4) | (plmn->mnc_digit1 & 0x0F);
  return;
}
