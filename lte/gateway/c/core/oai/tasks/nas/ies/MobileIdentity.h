/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#ifndef MOBILE_IDENTITY_H_
#define MOBILE_IDENTITY_H_
#include <stdint.h>

#define MOBILE_IDENTITY_MINIMUM_LENGTH 3
#define MOBILE_IDENTITY_MAXIMUM_LENGTH 11

#define MOBILE_IDENTITY_NOT_AVAILABLE_GSM_LENGTH 1
#define MOBILE_IDENTITY_NOT_AVAILABLE_GPRS_LENGTH 3
#define MOBILE_IDENTITY_NOT_AVAILABLE_LTE_LENGTH 3

typedef struct {
  uint8_t digit1 : 4;
  uint8_t oddeven : 1;
  uint8_t typeofidentity : 3;
  uint8_t digit2 : 4;
  uint8_t digit3 : 4;
  uint8_t digit4 : 4;
  uint8_t digit5 : 4;
  uint8_t digit6 : 4;
  uint8_t digit7 : 4;
  uint8_t digit8 : 4;
  uint8_t digit9 : 4;
  uint8_t digit10 : 4;
  uint8_t digit11 : 4;
  uint8_t digit12 : 4;
  uint8_t digit13 : 4;
  uint8_t digit14 : 4;
  uint8_t digit15 : 4;
  uint8_t numOfValidImsiDigits : 4;
} ImsiMobileIdentity_t;

typedef struct {
  uint8_t tac1 : 4;
  uint8_t oddeven : 1;
  uint8_t typeofidentity : 3;
  uint8_t tac2 : 4;
  uint8_t tac3 : 4;
  uint8_t tac4 : 4;
  uint8_t tac5 : 4;
  uint8_t tac6 : 4;
  uint8_t tac7 : 4;
  uint8_t tac8 : 4;
  uint8_t snr1 : 4;
  uint8_t snr2 : 4;
  uint8_t snr3 : 4;
  uint8_t snr4 : 4;
  uint8_t snr5 : 4;
  uint8_t snr6 : 4;
  uint8_t cdsd : 4;
} ImeiMobileIdentity_t;

typedef struct {
  uint8_t tac1 : 4;
  uint8_t oddeven : 1;
  uint8_t typeofidentity : 3;
  uint8_t tac2 : 4;
  uint8_t tac3 : 4;
  uint8_t tac4 : 4;
  uint8_t tac5 : 4;
  uint8_t tac6 : 4;
  uint8_t tac7 : 4;
  uint8_t tac8 : 4;
  uint8_t snr1 : 4;
  uint8_t snr2 : 4;
  uint8_t snr3 : 4;
  uint8_t snr4 : 4;
  uint8_t snr5 : 4;
  uint8_t snr6 : 4;
  uint8_t svn1 : 4;
  uint8_t svn2 : 4;
  uint8_t last : 4;
} ImeisvMobileIdentity_t;

typedef struct {
  uint8_t spare : 2;
  uint8_t mbmssessionidindication : 1;
  uint8_t mccmncindication : 1;
#define MOBILE_IDENTITY_EVEN 0
#define MOBILE_IDENTITY_ODD 1
  uint8_t oddeven : 1;
  uint8_t typeofidentity : 3;
  uint32_t mbmsserviceid;
  uint8_t mccdigit2 : 4;
  uint8_t mccdigit1 : 4;
  uint8_t mncdigit3 : 4;
  uint8_t mccdigit3 : 4;
  uint8_t mncdigit2 : 4;
  uint8_t mncdigit1 : 4;
  uint8_t mbmssessionid;
} TmgiMobileIdentity_t;

typedef ImsiMobileIdentity_t TmsiMobileIdentity_t;
typedef ImsiMobileIdentity_t NoMobileIdentity_t;

typedef union MobileIdentity_tag {
#define MOBILE_IDENTITY_IMSI 0b001
#define MOBILE_IDENTITY_IMEI 0b010
#define MOBILE_IDENTITY_IMEISV 0b011
#define MOBILE_IDENTITY_TMSI 0b100
#define MOBILE_IDENTITY_TMGI 0b101
#define MOBILE_IDENTITY_NOT_AVAILABLE 0b000
  ImsiMobileIdentity_t imsi;
  ImeiMobileIdentity_t imei;
  ImeisvMobileIdentity_t imeisv;
  TmsiMobileIdentity_t tmsi;
  TmgiMobileIdentity_t tmgi;
  NoMobileIdentity_t no_id;
} MobileIdentity;

int encode_mobile_identity(
    MobileIdentity* mobileidentity, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_mobile_identity(
    MobileIdentity* mobileidentity, uint8_t iei, uint8_t* buffer, uint32_t len);

void dump_mobile_identity_xml(MobileIdentity* mobileidentity, uint8_t iei);

#endif /* MOBILE IDENTITY_H_ */
