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

/*! \file 3gpp_24.008.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_3GPP_24_008_SEEN
#define FILE_3GPP_24_008_SEEN

#include <stdbool.h>
#include <stdint.h>

#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"

#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"

//#warning "Set it to max size of message"
#define IE_UNDEFINED_MAX_LENGTH 1024

//******************************************************************************
// 10.5.1 Common information elements
//******************************************************************************

typedef enum common_ie_e {
  C_MOBILE_STATION_CLASSMARK_2_IEI            = 0x11, /* 0x11 = 17   */
  C_LOCATION_AREA_IDENTIFICATION_IEI          = 0x13, /* 0x13 = 19   */
  C_MOBILE_STATION_CLASSMARK_3_IEI            = 0x20, /* 0x20 = 32   */
  C_MOBILE_IDENTITY_IEI                       = 0x23, /* 0x23 = 35   */
  C_PLMN_LIST_IEI                             = 0x4A, /* 0x4A = 74   */
  C_CIPHERING_KEY_SEQUENCE_NUMBER_IEI         = 0x80, /* 0x80 = 128  */
  C_MS_NETWORK_FEATURE_SUPPORT_IEI            = 0xC0, /* 0xC- = 192- */
  C_NETWORK_RESOURCE_IDENTIFIER_CONTAINER_IEI = 0x10, /* 0x10 = 16   */
} common_ie_t;

//------------------------------------------------------------------------------
// 10.5.1.2 Ciphering Key Sequence Number
//------------------------------------------------------------------------------

#define CIPHERING_KEY_SEQUENCE_NUMBER_IE_TYPE 1
#define CIPHERING_KEY_SEQUENCE_NUMBER_IE_MIN_LENGTH 1
#define CIPHERING_KEY_SEQUENCE_NUMBER_IE_MAX_LENGTH 1

typedef uint8_t ciphering_key_sequence_number_t;

int encode_ciphering_key_sequence_number_ie(
    ciphering_key_sequence_number_t* cipheringkeysequencenumber,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_ciphering_key_sequence_number_ie(
    ciphering_key_sequence_number_t* cipheringkeysequencenumber,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.1.3 Location Area Identification
//------------------------------------------------------------------------------
#define LOCATION_AREA_IDENTIFICATION_IE_TYPE 3
#define LOCATION_AREA_IDENTIFICATION_IE_MIN_LENGTH 6
#define LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH 6

#define INVALID_LAC_0000                                                       \
  (uint16_t) 0x0000 /*!< \brief  This LAC can be coded using a full            \
                       hexadecimal representation except for the following     \
                       reserved hexadecimal values: 0000, and FFFE.   */

#define INVALID_LAC_FFFE                                                       \
  (uint16_t) 0xFFFE /*!< \brief  This LAC can be coded using a full            \
                       hexadecimal representation except for the following     \
                       reserved hexadecimal values: 0000, and FFFE.   */
typedef uint16_t
    lac_t; /*!< \brief  Location Area Code (LAC) is a fixed length code (of 2
              octets) identifying a location area within a PLMN */

typedef struct location_area_identification_s {
  uint8_t mccdigit2 : 4;
  uint8_t mccdigit1 : 4;
  uint8_t mncdigit3 : 4;
  uint8_t mccdigit3 : 4;
  uint8_t mncdigit2 : 4;
  uint8_t mncdigit1 : 4;
  lac_t lac;
} location_area_identification_t;

typedef location_area_identification_t lai_t;

int encode_location_area_identification_ie(
    location_area_identification_t* locationareaidentification,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_location_area_identification_ie(
    location_area_identification_t* locationareaidentification,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.1.4 Mobile Identity
//------------------------------------------------------------------------------
#define MOBILE_IDENTITY_IE_TYPE 4
#define MOBILE_IDENTITY_IE_MIN_LENGTH 3
#define MOBILE_IDENTITY_IE_MAX_LENGTH 11

#define MOBILE_IDENTITY_NOT_AVAILABLE_GSM_LENGTH 1
#define MOBILE_IDENTITY_NOT_AVAILABLE_GPRS_LENGTH 3
#define MOBILE_IDENTITY_NOT_AVAILABLE_LTE_LENGTH 3
#define MOBILE_IDENTITY_IE_IMSI_LENGTH 8
#define MOBILE_IDENTITY_IE_IMEI_LENGTH 8
#define MOBILE_IDENTITY_IE_IMEISV_LENGTH 9
#define MOBILE_IDENTITY_IE_TMGI_LENGTH 9
#define MOBILE_IDENTITY_IE_TMSI_LENGTH 5

typedef struct imsi_mobile_identity_s {
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
} imsi_mobile_identity_t;

typedef struct imei_mobile_identity_s {
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
} imei_mobile_identity_t;

typedef struct imeisv_mobile_identity_s {
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
} imeisv_mobile_identity_t;

typedef struct tmgi_mobile_identity_s {
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
} tmgi_mobile_identity_t;

typedef struct tmsi_mobile_identity_s {
  uint8_t f : 4;
  uint8_t oddeven : 1;
  uint8_t typeofidentity : 3;
  uint8_t tmsi[4];
} tmsi_mobile_identity_t;

typedef struct no_mobile_identity_s {
  uint8_t digit1 : 4;
  uint8_t oddeven : 1;
  uint8_t typeofidentity : 3;
  uint8_t digit2 : 4;
  uint8_t digit3 : 4;
  uint8_t digit4 : 4;
  uint8_t digit5 : 4;
} no_mobile_identity_t;

typedef union mobile_identity_s {
#define MOBILE_IDENTITY_IMSI 0b001
#define MOBILE_IDENTITY_IMEI 0b010
#define MOBILE_IDENTITY_IMEISV 0b011
#define MOBILE_IDENTITY_TMSI 0b100
#define MOBILE_IDENTITY_TMGI 0b101
#define MOBILE_IDENTITY_NOT_AVAILABLE 0b000
  imsi_mobile_identity_t imsi;
  imei_mobile_identity_t imei;
  imeisv_mobile_identity_t imeisv;
  tmsi_mobile_identity_t tmsi;
  tmgi_mobile_identity_t tmgi;
  no_mobile_identity_t no_id;
} mobile_identity_t;

int encode_mobile_identity_ie(
    mobile_identity_t* mobileidentity, const bool iei_present, uint8_t* buffer,
    const uint32_t len);
int decode_mobile_identity_ie(
    mobile_identity_t* mobileidentity, const bool iei_present, uint8_t* buffer,
    const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.1.6 Mobile Station Classmark 2
//------------------------------------------------------------------------------
#define MOBILE_STATION_CLASSMARK_2_IE_TYPE 4
#define MOBILE_STATION_CLASSMARK_2_IE_MIN_LENGTH 5
#define MOBILE_STATION_CLASSMARK_2_IE_MAX_LENGTH 5

typedef struct mobile_station_classmark2_s {
  uint8_t revisionlevel : 2;
  uint8_t esind : 1;
  uint8_t a51 : 1;
  uint8_t rfpowercapability : 3;
  uint8_t pscapability : 1;
  uint8_t ssscreenindicator : 2;
  uint8_t smcapability : 1;
  uint8_t vbs : 1;
  uint8_t vgcs : 1;
  uint8_t fc : 1;
  uint8_t cm3 : 1;
  uint8_t lcsvacap : 1;
  uint8_t ucs2 : 1;
  uint8_t solsa : 1;
  uint8_t cmsp : 1;
  uint8_t a53 : 1;
  uint8_t a52 : 1;
} mobile_station_classmark2_t;

int encode_mobile_station_classmark_2_ie(
    mobile_station_classmark2_t* mobilestationclassmark2,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_mobile_station_classmark_2_ie(
    mobile_station_classmark2_t* mobilestationclassmark2,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.1.7 Mobile Station Classmark 3
//------------------------------------------------------------------------------
#define MOBILE_STATION_CLASSMARK_3_IE_TYPE 4
#define MOBILE_STATION_CLASSMARK_3_IE_MIN_LENGTH 34  // TODO
#define MOBILE_STATION_CLASSMARK_3_IE_MAX_LENGTH 34

typedef struct mobile_station_classmark3_s {
  uint8_t byte[32];  // TODO
} mobile_station_classmark3_t;

int encode_mobile_station_classmark_3_ie(
    mobile_station_classmark3_t* mobilestationclassmark3,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_mobile_station_classmark_3_ie(
    mobile_station_classmark3_t* mobilestationclassmark3,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.1.13 PLMN list
//------------------------------------------------------------------------------
#define PLMN_LIST_IE_TYPE 4
#define PLMN_LIST_IE_MIN_LENGTH 5
#define PLMN_LIST_IE_MAX_LENGTH 47
#define PLMN_LIST_IE_MAX_PLMN 15

typedef struct plmn_list_s {
  plmn_t plmn[PLMN_LIST_IE_MAX_PLMN];
  uint8_t num_plmn;
} plmn_list_t;

int encode_plmn_list_ie(
    plmn_list_t* plmnlist, const bool iei_present, uint8_t* buffer,
    const uint32_t len);
int decode_plmn_list_ie(
    plmn_list_t* plmnlist, const bool iei_present, uint8_t* buffer,
    const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.1.15 MS network feature support
//------------------------------------------------------------------------------
#define MS_NETWORK_FEATURE_SUPPORT_IE_TYPE 1
#define MS_NETWORK_FEATURE_SUPPORT_IE_MIN_LENGTH 1
#define MS_NETWORK_FEATURE_SUPPORT_IE_MAX_LENGTH 1

typedef struct ms_network_feature_support_s {
  uint8_t spare_bits : 3;
  uint8_t extended_periodic_timers : 1;
} ms_network_feature_support_t;

int encode_ms_network_feature_support_ie(
    ms_network_feature_support_t* msnetworkfeaturesupport,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_ms_network_feature_support_ie(
    ms_network_feature_support_t* msnetworkfeaturesupport,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
//------------------------------------------------------------------------------
// 10.5.5.31 Network Resource Identifier Container
//------------------------------------------------------------------------------
#define NETWORK_RESOURCE_IDENTIFIER_CONTAINER_IE_TYPE 4
#define NETWORK_RESOURCE_IDENTIFIER_CONTAINER_IE_MIN_LENGTH 4  // TODO
#define NETWORK_RESOURCE_IDENTIFIER_CONTAINER_IE_MAX_LENGTH 4

typedef struct network_resource_identifier_container_s {
  uint8_t byte[32];  // TODO
} network_resource_identifier_container_t;

int encode_network_resource_identifier_container_ie(
    network_resource_identifier_container_t* networkresourceidentifiercontainer,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_network_resource_identifier_container_ie(
    network_resource_identifier_container_t* networkresourceidentifiercontainer,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//******************************************************************************
// 10.5.3 Mobility management information elements.
//******************************************************************************

typedef enum mobility_managenent_ie_e {
  MM_AUTHENTICATION_PARAMETER_RAND_IEI     = 0x21, /* 0x21 = 33 */
  MM_AUTHENTICATION_PARAMETER_AUTN_IEI     = 0x28, /* 0x28 = 40 */
  MM_AUTHENTICATION_RESPONSE_PARAMETER_IEI = 0x29, /* 0x28 = 41 */
  MM_AUTHENTICATION_FAILURE_PARAMETER_IEI  = 0x30, /* 0x30 = 48 */
  MM_EMERGENCY_NUMBER_LIST_IEI             = 0x34, /* 0x34 = 52 */
  MM_FULL_NETWORK_NAME_IEI                 = 0x43, /* 0x43 = 67 */
  MM_SHORT_NETWORK_NAME_IEI                = 0x45, /* 0x45 = 69 */
  MM_TIME_ZONE_IEI                         = 0x46, /* 0x46 = 70 */
  MM_TIME_ZONE_AND_TIME_IEI                = 0x47, /* 0x47 = 71 */
  MM_DAYLIGHT_SAVING_TIME_IEI              = 0x49, /* 0x49 = 73 */
} mobility_managenent_ie_t;

//------------------------------------------------------------------------------
// 10.5.3.1 Authentication parameter RAND
//------------------------------------------------------------------------------
#define AUTHENTICATION_PARAMETER_RAND_IE_TYPE 3
#define AUTHENTICATION_PARAMETER_RAND_IE_MIN_LENGTH 17
#define AUTHENTICATION_PARAMETER_RAND_IE_MAX_LENGTH 17

typedef bstring authentication_parameter_rand_t;

int encode_authentication_parameter_rand_ie(
    authentication_parameter_rand_t authenticationparameterrand,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_authentication_parameter_rand_ie(
    authentication_parameter_rand_t* authenticationparameterrand,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.3.1.1 Authentication Parameter AUTN (UMTS and EPS authentication
// challenge)
//------------------------------------------------------------------------------
#define AUTHENTICATION_PARAMETER_AUTN_IE_TYPE 4
#define AUTHENTICATION_PARAMETER_AUTN_IE_MIN_LENGTH 18
#define AUTHENTICATION_PARAMETER_AUTN_IE_MAX_LENGTH 18

typedef bstring authentication_parameter_autn_t;

int encode_authentication_parameter_autn_ie(
    authentication_parameter_autn_t authenticationparameterautn,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_authentication_parameter_autn_ie(
    authentication_parameter_autn_t* authenticationparameterautn,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.3.2 Authentication Response parameter
//------------------------------------------------------------------------------
#define AUTHENTICATION_RESPONSE_PARAMETER_IE_TYPE 3
#define AUTHENTICATION_RESPONSE_PARAMETER_IE_MIN_LENGTH 5
#define AUTHENTICATION_RESPONSE_PARAMETER_IE_MAX_LENGTH 5

typedef bstring authentication_response_parameter_t;

int encode_authentication_response_parameter_ie(
    authentication_response_parameter_t authenticationresponseparameter,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_authentication_response_parameter_ie(
    authentication_response_parameter_t* authenticationresponseparameter,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.3.2.2 Authentication Failure parameter (UMTS and EPS authentication
// challenge)
//------------------------------------------------------------------------------
#define AUTHENTICATION_FAILURE_PARAMETER_IE_TYPE 4
#define AUTHENTICATION_FAILURE_PARAMETER_IE_MIN_LENGTH 16
#define AUTHENTICATION_FAILURE_PARAMETER_IE_MAX_LENGTH 16

typedef bstring authentication_failure_parameter_t;

int encode_authentication_failure_parameter_ie(
    authentication_failure_parameter_t authenticationfailureparameter,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_authentication_failure_parameter_ie(
    authentication_failure_parameter_t* authenticationfailureparameter,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.3.5a Network Name
//------------------------------------------------------------------------------
#define NETWORK_NAME_IE_TYPE 4
#define NETWORK_NAME_IE_MIN_LENGTH 4
// TODO
#define NETWORK_NAME_IE_MAX_LENGTH 255

typedef struct network_name_s {
  uint8_t codingscheme : 3;
  uint8_t addci : 1;
  uint8_t numberofsparebitsinlastoctet : 3;
  bstring textstring;
} network_name_t;

int encode_network_name_ie(
    network_name_t* networkname, const uint8_t iei, uint8_t* buffer,
    const uint32_t len);
int decode_network_name_ie(
    network_name_t* networkname, const uint8_t iei, uint8_t* buffer,
    const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.3.8 Time Zone
//------------------------------------------------------------------------------
#define TIME_ZONE_IE_TYPE 3
#define TIME_ZONE_IE_MIN_LENGTH 2
#define TIME_ZONE_IE_MAX_LENGTH 2

typedef uint8_t time_zone_t;

int encode_time_zone(
    time_zone_t* timezone, const bool iei_present, uint8_t* buffer,
    const uint32_t len);
int decode_time_zone(
    time_zone_t* timezone, const bool iei_present, uint8_t* buffer,
    const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.3.9 Time Zone and Time
//------------------------------------------------------------------------------
#define TIME_ZONE_AND_TIME_IE_TYPE 3
#define TIME_ZONE_AND_TIME_MIN_LENGTH 8
#define TIME_ZONE_AND_TIME_MAX_LENGTH 8

typedef struct time_zone_and_time_s {
  uint8_t year;
  uint8_t month;
  uint8_t day;
  uint8_t hour;
  uint8_t minute;
  uint8_t second;
  time_zone_t timezone;
} time_zone_and_time_t;

int encode_time_zone_and_time(
    time_zone_and_time_t* timezoneandtime, const bool iei_present,
    uint8_t* buffer, const uint32_t len);
int decode_time_zone_and_time(
    time_zone_and_time_t* timezoneandtime, const bool iei_present,
    uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.3.12 Daylight Saving Time
//------------------------------------------------------------------------------
#define DAYLIGHT_SAVING_TIME_IE_TYPE 4
#define DAYLIGHT_SAVING_TIME_IE_MIN_LENGTH 3
#define DAYLIGHT_SAVING_TIME_IE_MAX_LENGTH 3

typedef uint8_t daylight_saving_time_t;

int encode_daylight_saving_time_ie(
    daylight_saving_time_t* daylightsavingtime, const bool iei_present,
    uint8_t* buffer, const uint32_t len);
int decode_daylight_saving_time_ie(
    daylight_saving_time_t* daylightsavingtime, const bool iei_present,
    uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.3.13 Emergency Number List
//------------------------------------------------------------------------------
#define EMERGENCY_NUMBER_LIST_IE_TYPE 4
#define EMERGENCY_NUMBER_LIST_IE_MIN_LENGTH 5
#define EMERGENCY_NUMBER_LIST_IE_MAX_LENGTH 50

typedef struct emergency_number_list_s {
  uint8_t lengthofemergencynumberinformation;
  uint8_t spare : 3;
  uint8_t emergencyservicecategoryvalue : 5;
#define EMERGENCY_NUMBER_MAX_DIGITS (MSISDN_DIGIT_SIZE * 2)
  uint8_t number_digit[EMERGENCY_NUMBER_MAX_DIGITS];
  struct emergency_number_list_s* next;
} emergency_number_list_t;

int encode_emergency_number_list_ie(
    emergency_number_list_t* emergencynumberlist, const bool iei_present,
    uint8_t* buffer, const uint32_t len);
int decode_emergency_number_list_ie(
    emergency_number_list_t* emergencynumberlist, const bool iei_present,
    uint8_t* buffer, const uint32_t len);

//******************************************************************************
// 10.5.4 Call control information elements
//******************************************************************************
typedef enum call_control_ie_e {
  CC_SUPPORTED_CODEC_LIST_IE = 0x40, /* 0x40 = 64  */
} call_control_ie_t;

//------------------------------------------------------------------------------
// 10.5.4.32 Supported codec list
//------------------------------------------------------------------------------
#define SUPPORTED_CODEC_LIST_IE_TYPE 4
#define SUPPORTED_CODEC_LIST_IE_MIN_LENGTH 5
#define SUPPORTED_CODEC_LIST_IE_MAX_LENGTH IE_UNDEFINED_MAX_LENGTH

typedef bstring supported_codec_list_t;

int encode_supported_codec_list_ie(
    supported_codec_list_t* supportedcodeclist, const bool iei_present,
    uint8_t* buffer, const uint32_t len);
int decode_supported_codec_list_ie(
    supported_codec_list_t* supportedcodeclist, const bool iei_present,
    uint8_t* buffer, const uint32_t len);

//******************************************************************************
// 10.5.5 GPRS mobility management information elements
//******************************************************************************

typedef enum gprs_mobility_managenent_ie_e {
  GMM_PTMSI_SIGNATURE_IEI                              = 0x19, /* 0x19 = 25  */
  GMM_MS_NETWORK_CAPABILITY_IEI                        = 0x31, /* 0x31 = 49  */
  GMM_DRX_PARAMETER_IEI                                = 0x5C, /* 0x5C = 92  */
  GMM_VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_IEI = 0x5D, /* 0x5D = 93  */
  GMM_TMSI_STATUS_IEI    = 0x90, /* 0x90 = 144 (shifted by 4)*/
  GMM_IMEISV_REQUEST_IEI = 0xC0, /* 0xC0 = 192 (shifted by 4)*/
} gprs_mobility_managenent_ie_t;

//------------------------------------------------------------------------------
// 10.5.5.4 TMSI status
//------------------------------------------------------------------------------
#define TMSI_STATUS_IE_TYPE 1
#define TMSI_STATUS_IE_MIN_LENGTH 1
#define TMSI_STATUS_IE_MAX_LENGTH 1

typedef uint8_t tmsi_status_t;

int encode_tmsi_status(
    tmsi_status_t* tmsistatus, const bool iei_present, uint8_t* buffer,
    const uint32_t len);
int decode_tmsi_status(
    tmsi_status_t* tmsistatus, const bool iei_present, uint8_t* buffer,
    const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.5.6 DRX parameter
//------------------------------------------------------------------------------
#define DRX_PARAMETER_IE_TYPE 3
#define DRX_PARAMETER_IE_MIN_LENGTH 3
#define DRX_PARAMETER_IE_MAX_LENGTH 3

typedef struct drx_parameter_s {
  uint8_t splitpgcyclecode;
  uint8_t cnspecificdrxcyclelengthcoefficientanddrxvaluefors1mode : 4;
  uint8_t splitonccch : 1;
  uint8_t nondrxtimer : 3;
} drx_parameter_t;

int encode_drx_parameter_ie(
    drx_parameter_t* drxparameter, const bool iei_present, uint8_t* buffer,
    const uint32_t len);
int decode_drx_parameter_ie(
    drx_parameter_t* drxparameter, const bool iei_present, uint8_t* buffer,
    const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.5.8 P-TMSI signature
//------------------------------------------------------------------------------
#define P_TMSI_SIGNATURE_IE_TYPE 3
#define P_TMSI_SIGNATURE_IE_MIN_LENGTH 4
#define P_TMSI_SIGNATURE_IE_MAX_LENGTH 4

typedef uint32_t p_tmsi_signature_t;

int encode_p_tmsi_signature_ie(
    p_tmsi_signature_t ptmsisignature, const bool iei_present, uint8_t* buffer,
    const uint32_t len);
int decode_p_tmsi_signature_ie(
    p_tmsi_signature_t* ptmsisignature, const bool iei_present, uint8_t* buffer,
    const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.5.9 Identity type 2
//------------------------------------------------------------------------------
#define IDENTITY_TYPE_2_TYPE 1
#define IDENTITY_TYPE_2_IE_MIN_LENGTH 1
#define IDENTITY_TYPE_2_IE_MAX_LENGTH 1

#define IDENTITY_TYPE_2_IMSI 0b001
#define IDENTITY_TYPE_2_IMEI 0b010
#define IDENTITY_TYPE_2_IMEISV 0b011
#define IDENTITY_TYPE_2_TMSI 0b100

typedef uint8_t identity_type2_t;

int encode_identity_type_2_ie(
    identity_type2_t* identitytype2, bool is_ie_present, uint8_t* buffer,
    const uint32_t len);
int decode_identity_type_2_ie(
    identity_type2_t* identitytype2, bool is_ie_present, uint8_t* buffer,
    const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.5.10 IMEISV request
//------------------------------------------------------------------------------
#define IMEISV_REQUEST_IE_TYPE 1
#define IMEISV_REQUEST_IE_MIN_LENGTH 1
#define IMEISV_REQUEST_IE_MAX_LENGTH 1

typedef uint8_t imeisv_request_t;

#define IMEISV_NOT_REQUESTED 0b000
#define IMEISV_REQUESTED 0b001

int encode_imeisv_request_ie(
    imeisv_request_t* imeisvrequest, bool is_ie_present, uint8_t* buffer,
    const uint32_t len);
int decode_imeisv_request_ie(
    imeisv_request_t* imeisvrequest, bool is_ie_present, uint8_t* buffer,
    const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.5.12 MS network capability
//------------------------------------------------------------------------------
#define MS_NETWORK_CAPABILITY_IE_TYPE 4
// TODO
#define MS_NETWORK_CAPABILITY_IE_MIN_LENGTH 5
#define MS_NETWORK_CAPABILITY_IE_MAX_LENGTH 10

typedef struct ms_network_capability_s {
#define MS_NETWORK_CAPABILITY_GEA1 0b10000000
  uint8_t gea1 : 1;
#define MS_NETWORK_CAPABILITY_SM_CAP_VIA_DEDICATED_CHANNELS 0b01000000
  uint8_t smdc : 1;
#define MS_NETWORK_CAPABILITY_SM_CAP_VIA_GPRS_CHANNELS 0b00100000
  uint8_t smgc : 1;
#define MS_NETWORK_CAPABILITY_UCS2_SUPPORT 0b00010000
  uint8_t ucs2 : 1;
#define MS_NETWORK_CAPABILITY_SS_SCREENING_INDICATOR 0b00001100
  uint8_t sssi : 2;
#define MS_NETWORK_CAPABILITY_SOLSA 0b00000010
  uint8_t solsa : 1;
#define MS_NETWORK_CAPABILITY_REVISION_LEVEL_INDICATOR 0b00000001
  uint8_t revli : 1;

#define MS_NETWORK_CAPABILITY_PFC_FEATURE_MODE 0b10000000
  uint8_t pfc : 1;
#define MS_NETWORK_CAPABILITY_GEA2 0b01000000
#define MS_NETWORK_CAPABILITY_GEA3 0b00100000
#define MS_NETWORK_CAPABILITY_GEA4 0b00010000
#define MS_NETWORK_CAPABILITY_GEA5 0b00001000
#define MS_NETWORK_CAPABILITY_GEA6 0b00000100
#define MS_NETWORK_CAPABILITY_GEA7 0b00000010
  uint8_t egea : 6;
#define MS_NETWORK_CAPABILITY_LCS_VA 0b00000001
  uint8_t lcs : 1;

#define MS_NETWORK_CAPABILITY_PS_INTER_RAT_HO_GERAN_TO_UTRAN_IU 0b10000000
  uint8_t ps_ho_utran : 1;
#define MS_NETWORK_CAPABILITY_PS_INTER_RAT_HO_GERAN_TO_EUTRAN_S1 0b01000000
  uint8_t ps_ho_eutran : 1;
#define MS_NETWORK_CAPABILITY_EMM_COMBINED_PROCEDURE 0b00100000
  uint8_t emm_cpc : 1;
#define MS_NETWORK_CAPABILITY_ISR 0b00010000
  uint8_t isr : 1;
#define MS_NETWORK_CAPABILITY_SRVCC 0b00001000
  uint8_t srvcc : 1;
#define MS_NETWORK_CAPABILITY_EPC 0b00000100
  uint8_t epc_cap : 1;
#define MS_NETWORK_CAPABILITY_NOTIFICATION 0b00000010
  uint8_t nf_cap : 1;
#define MS_NETWORK_CAPABILITY_GERAN_NETWORK_SHARING 0b00000001
  uint8_t geran_ns : 1;

#define MS_NETWORK_CAPABILITY_USER_PLANE_INTEGRITY_PROTECTION_SUPPORT 0b10000000
  uint8_t up_integ_prot_support : 1;
#define MS_NETWORK_CAPABILITY_GIA4 0b01000000
  uint8_t gia4 : 1;
#define MS_NETWORK_CAPABILITY_GIA5 0b00100000
  uint8_t gia5 : 1;
#define MS_NETWORK_CAPABILITY_GIA6 0b00010000
  uint8_t gia6 : 1;
#define MS_NETWORK_CAPABILITY_GIA7 0b00001000
  uint8_t gia7 : 1;
#define MS_NETWORK_CAPABILITY_EPCO_IE_INDICATOR 0b00000100
  uint8_t epco_ie_ind : 1;
#define MS_NETWORK_CAPABILITY_RESTRICTION_ON_USE_OF_ENHANCED_COVERAGE_CAPABILITY \
  0b00000010
  uint8_t rest_use_enhanc_cov_cap : 1;
#define MS_NETWORK_CAPABILITY_DUAL_CONNECTIVITY_EUTRA_NR_CAPABILITY 0b00000001
  uint8_t en_dc : 1;
} ms_network_capability_t;

int encode_ms_network_capability_ie(
    ms_network_capability_t* msnetworkcapability, const bool iei_present,
    uint8_t* buffer, const uint32_t len) __attribute__((unused));
int decode_ms_network_capability_ie(
    ms_network_capability_t* msnetworkcapability, const bool iei_present,
    uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.5.15 Routing area identification
//------------------------------------------------------------------------------

typedef uint8_t routing_area_code_t; /*!< \brief  Routing Area Code    */

/*! \struct  rai_t
 * \brief Routing Area Identification (RAI).
 */
typedef struct routing_area_identification_s {
  location_area_identification_t
      lai; /*!< \brief  See 3GPP TS 23.003 4.2 Composition of the Routing Area
              Identification (RAI)    */
  routing_area_code_t rac; /*!< \brief  Routing Area Code    */
} routing_area_identification_t;

//------------------------------------------------------------------------------
// 10.5.5.28 Voice domain preference and UE's usage setting
//------------------------------------------------------------------------------
#define VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_TYPE 4
#define VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_MINIMUM_LENGTH 3
#define VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_MAXIMUM_LENGTH 3

typedef struct voice_domain_preference_and_ue_usage_setting_s {
  uint8_t spare : 5;
#define UE_USAGE_SETTING_VOICE_CENTRIC 0b0
#define UE_USAGE_SETTING_DATA_CENTRIC 0b1
  uint8_t ue_usage_setting : 1;
#define VOICE_DOMAIN_PREFERENCE_CS_VOICE_ONLY 0b00
#define VOICE_DOMAIN_PREFERENCE_IMS_PS_VOICE_ONLY 0b01
#define VOICE_DOMAIN_PREFERENCE_CS_VOICE_PREFERRED_IMS_PS_VOICE_AS_SECONDARY   \
  0b10
#define VOICE_DOMAIN_PREFERENCE_IMS_PS_VOICE_PREFERRED_CS_VOICE_AS_SECONDARY   \
  0b11
  uint8_t voice_domain_for_eutran : 2;
} voice_domain_preference_and_ue_usage_setting_t;

int encode_voice_domain_preference_and_ue_usage_setting(
    voice_domain_preference_and_ue_usage_setting_t*
        voicedomainpreferenceandueusagesetting,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_voice_domain_preference_and_ue_usage_setting(
    voice_domain_preference_and_ue_usage_setting_t*
        voicedomainpreferenceandueusagesetting,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//******************************************************************************
// 10.5.6 Session management information elements
//******************************************************************************

typedef enum session_managenent_ie_e {
  SM_PROTOCOL_CONFIGURATION_OPTIONS_IEI      = 0x27, /* 0x27 = 39 */
  SM_ACCESS_POINT_NAME_IEI                   = 0x28, /* 0x28 = 40 */
  SM_QUALITY_OF_SERVICE_IEI                  = 0x30, /* 0x30 = 48 */
  SM_LLC_SERVICE_ACCESS_POINT_IDENTIFIER_IEI = 0x32, /* 0x32 = 50 */
  SM_PACKET_FLOW_IDENTIFIER_IEI              = 0x34, /* 0x34 = 52 */
  SM_TRAFFIC_FLOW_TEMPLATE_IEI               = 0x36, /* 0x36 = 54 */
  SM_LINKED_TI_IEI                           = 0x5D, /* 0x5D = 93 */

} session_managenent_ie_t;

//------------------------------------------------------------------------------
// 10.5.6.1 Access Point Name
//------------------------------------------------------------------------------
#define ACCESS_POINT_NAME_IE_TYPE 4
#define ACCESS_POINT_NAME_IE_MIN_LENGTH 3
#define ACCESS_POINT_NAME_IE_MAX_LENGTH 102
#define ACCESS_POINT_NAME_MAX_LENGTH 100

typedef bstring access_point_name_t;

int encode_access_point_name_ie(
    access_point_name_t accesspointname, const bool iei_present,
    uint8_t* buffer, const uint32_t len);
int decode_access_point_name_ie(
    access_point_name_t* accesspointname, const bool iei_present,
    uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.6.3 Protocol configuration options
//------------------------------------------------------------------------------
#define PROTOCOL_CONFIGURATION_OPTIONS_IE_TYPE 4
#define PROTOCOL_CONFIGURATION_OPTIONS_IE_MIN_LENGTH 3
#define PROTOCOL_CONFIGURATION_OPTIONS_IE_MAX_LENGTH 253
#define PCO_CONFIGURATION_PROTOCOL_PPP_FOR_USE_WITH_IP_PDP_TYPE_OR_IP_PDN_TYPE \
  0b000

// Protocol identifiers defined in RFC 3232
#define PCO_PI_LCP (0xC021)
#define PCO_PI_PAP (0xC023)
#define PCO_PI_CHAP (0xC223)
#define PCO_PI_IPCP (0x8021)

/* CONTAINER IDENTIFIER MS to network direction:*/
#define PCO_CI_P_CSCF_IPV6_ADDRESS_REQUEST (0x0001)
#define PCO_CI_DNS_SERVER_IPV6_ADDRESS_REQUEST (0x0003)
// NOT SUPPORTED                                                        (0x0004)
#define PCO_CI_MS_SUPPORT_OF_NETWORK_REQUESTED_BEARER_CONTROL_INDICATOR (0x0005)
// RESERVED                                                             (0x0006)
#define PCO_CI_DSMIPV6_HOME_AGENT_ADDRESS_REQUEST (0x0007)
#define PCO_CI_DSMIPV6_HOME_NETWORK_PREFIX_REQUEST (0x0008)
#define PCO_CI_DSMIPV6_IPV4_HOME_AGENT_ADDRESS_REQUEST (0x0009)
#define PCO_CI_IP_ADDRESS_ALLOCATION_VIA_NAS_SIGNALLING (0x000A)
#define PCO_CI_IPV4_ADDRESS_ALLOCATION_VIA_DHCPV4 (0x000B)
#define PCO_CI_P_CSCF_IPV4_ADDRESS_REQUEST (0x000C)
#define PCO_CI_DNS_SERVER_IPV4_ADDRESS_REQUEST (0x000D)
#define PCO_CI_MSISDN_REQUEST (0x000E)
#define PCO_CI_IFOM_SUPPORT_REQUEST (0x000F)
#define PCO_CI_IPV4_LINK_MTU_REQUEST (0x0010)
// RESERVED (0xFF00..FFFF)

/* CONTAINER IDENTIFIER Network to MS direction:*/
#define PCO_CI_P_CSCF_IPV6_ADDRESS (0x0001)
#define PCO_CI_DNS_SERVER_IPV6_ADDRESS (0x0003)
#define PCO_CI_POLICY_CONTROL_REJECTION_CODE (0x0004)
#define PCO_CI_SELECTED_BEARER_CONTROL_MODE (0x0005)
// RESERVED                                                             (0x0006)
#define PCO_CI_DSMIPV6_HOME_AGENT_ADDRESS (0x0007)
#define PCO_CI_DSMIPV6_HOME_NETWORK_PREFIX (0x0008)
#define PCO_CI_DSMIPV6_IPV4_HOME_AGENT_ADDRESS (0x0009)
// RESERVED                                                             (0x000A)
// RESERVED                                                             (0x000B)
#define PCO_CI_P_CSCF_IPV4_ADDRESS (0x000C)
#define PCO_CI_DNS_SERVER_IPV4_ADDRESS (0x000D)
#define PCO_CI_MSISDN (0x000E)
#define PCO_CI_IFOM_SUPPORT (0x000F)
#define PCO_CI_IPV4_LINK_MTU (0x0010)
// RESERVED (0xFF00..FFFF)

/* Both directions:*/
#define PCO_CI_IM_CN_SUBSYSTEM_SIGNALING_FLAG (0x0002)

typedef struct pco_protocol_or_container_id_s {
  uint16_t id;
  uint8_t length;
  bstring contents;
} pco_protocol_or_container_id_t;

typedef struct protocol_configuration_options_s {
  uint8_t ext : 1;
  uint8_t spare : 4;
  uint8_t configuration_protocol : 3;
  uint8_t num_protocol_or_container_id;
  // arbitrary value, can be greater than defined (250/3)
  // Setting this value to 30 to support maximum possible number of protocol id
  // or container id defined in 24.008 release 13
#define PCO_UNSPEC_MAXIMUM_PROTOCOL_ID_OR_CONTAINER_ID 30
  pco_protocol_or_container_id_t
      protocol_or_container_ids[PCO_UNSPEC_MAXIMUM_PROTOCOL_ID_OR_CONTAINER_ID];
} protocol_configuration_options_t;

typedef struct TimeZoneAndTime_s {
  uint8_t year;
  uint8_t month;
  uint8_t day;
  uint8_t hour;
  uint8_t minute;
  uint8_t second;
  uint8_t timezone;
} TimeZoneAndTime_t;

void copy_protocol_configuration_options(
    protocol_configuration_options_t* const pco_dst,
    const protocol_configuration_options_t* const pco_src);
void clear_protocol_configuration_options(
    protocol_configuration_options_t* const pco);
void free_protocol_configuration_options(
    protocol_configuration_options_t** const protocol_configuration_options);

int decode_protocol_configuration_options(
    protocol_configuration_options_t* protocolconfigurationoptions,
    const uint8_t* const buffer, const uint32_t len);

int decode_protocol_configuration_options_ie(
    protocol_configuration_options_t* protocolconfigurationoptions,
    const bool iei_present, const uint8_t* const buffer, const uint32_t len);

int encode_protocol_configuration_options(
    const protocol_configuration_options_t* const protocolconfigurationoptions,
    uint8_t* buffer, const uint32_t len);

int encode_protocol_configuration_options_ie(
    const protocol_configuration_options_t* const protocolconfigurationoptions,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.6.5 Quality of service
//------------------------------------------------------------------------------
#define QUALITY_OF_SERVICE_IE_TYPE 4
#define QUALITY_OF_SERVICE_IE_MIN_LENGTH 14
#define QUALITY_OF_SERVICE_IE_MAX_LENGTH 18

typedef struct quality_of_service_s {
  uint8_t delayclass : 3;
  uint8_t reliabilityclass : 3;
  uint8_t peakthroughput : 4;
  uint8_t precedenceclass : 3;
  uint8_t meanthroughput : 5;
  uint8_t trafficclass : 3;
  uint8_t deliveryorder : 2;
  uint8_t deliveryoferroneoussdu : 3;
  uint8_t maximumsdusize;
  uint8_t maximumbitrateuplink;
  uint8_t maximumbitratedownlink;
  uint8_t residualber : 4;
  uint8_t sduratioerror : 4;
  uint8_t transferdelay : 6;
  uint8_t traffichandlingpriority : 2;
  uint8_t guaranteedbitrateuplink;
  uint8_t guaranteedbitratedownlink;
  uint8_t signalingindication : 1;
  uint8_t sourcestatisticsdescriptor : 4;
} quality_of_service_t;

int encode_quality_of_service_ie(
    quality_of_service_t* qualityofservice, const bool iei_present,
    uint8_t* buffer, const uint32_t len);
int decode_quality_of_service_ie(
    quality_of_service_t* qualityofservice, const bool iei_present,
    uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.6.7 Linked TI
//------------------------------------------------------------------------------
#define LINKED_TI_IE_TYPE 4
#define LINKED_TI_IE_MIN_LENGTH 3
#define LINKED_TI_IE_MAX_LENGTH 4
typedef struct linked_ti_s {
  uint8_t tiflag : 3;
  uint8_t tivalue : 3;
  uint8_t spare : 4;
  uint8_t ext : 1;
  uint8_t tivalue_cont : 5;
} linked_ti_t;

int encode_linked_ti_ie(
    linked_ti_t* llcserviceaccesspointidentifier, const bool iei_present,
    uint8_t* buffer, const uint32_t len);
int decode_linked_ti_ie(
    linked_ti_t* llcserviceaccesspointidentifier, const bool iei_present,
    uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.6.9 LLC service access point identifier
//------------------------------------------------------------------------------
#define LLC_SERVICE_ACCESS_POINT_IDENTIFIER_IE_TYPE 3
#define LLC_SERVICE_ACCESS_POINT_IDENTIFIER_IE_MIN_LENGTH 2
#define LLC_SERVICE_ACCESS_POINT_IDENTIFIER_IE_MAX_LENGTH 2

typedef uint8_t llc_service_access_point_identifier_t;

int encode_llc_service_access_point_identifier_ie(
    llc_service_access_point_identifier_t* llcserviceaccesspointidentifier,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_llc_service_access_point_identifier_ie(
    llc_service_access_point_identifier_t* llcserviceaccesspointidentifier,
    const bool iei_present, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.6.11 Packet Flow Identifier
//------------------------------------------------------------------------------
#define PACKET_FLOW_IDENTIFIER_IE_TYPE 4
#define PACKET_FLOW_IDENTIFIER_IE_MIN_LENGTH 3
#define PACKET_FLOW_IDENTIFIER_IE_MAX_LENGTH 3

typedef uint8_t packet_flow_identifier_t;

int encode_packet_flow_identifier_ie(
    packet_flow_identifier_t* packetflowidentifier, const bool iei_present,
    uint8_t* buffer, const uint32_t len);
int decode_packet_flow_identifier_ie(
    packet_flow_identifier_t* packetflowidentifier, const bool iei_present,
    uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
// 10.5.6.12 Traffic Flow Template
//------------------------------------------------------------------------------
#define TRAFFIC_FLOW_TEMPLATE_MINIMUM_LENGTH 3
#define TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH 257

/*
 * ----------------------------------------------------------------------------
 *        Packet filter list
 * ----------------------------------------------------------------------------
 */

#define TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR 0b00010000
#define TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR 0b00100000
#define TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER 0b00110000
#define TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT 0b01000000
#define TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE 0b01000001
#define TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT 0b01010000
#define TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE 0b01010001
#define TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX 0b01100000
#define TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS 0b01110000
#define TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL 0b10000000

/*
 * Port range
 * -----------
 */
typedef struct port_range_s {
  uint16_t lowlimit;
  uint16_t highlimit;
} port_range_t;

/*
 * Packet filter content
 * ---------------------
 */
typedef struct packet_filter_contents_s {
#define TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG (1 << 0)
#define TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG (1 << 1)
#define TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG (1 << 2)
#define TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG (1 << 3)
#define TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE_FLAG (1 << 4)
#define TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG (1 << 5)
#define TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE_FLAG (1 << 6)
#define TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX_FLAG (1 << 7)
#define TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS_FLAG (1 << 8)
#define TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG (1 << 9)
  uint16_t flags;
#define TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE 4
  struct {
    uint8_t addr;
    uint8_t mask;
  } ipv4remoteaddr[TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE];
#define TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE 16
  struct {
    uint8_t addr;
    uint8_t mask;
  } ipv6remoteaddr[TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE];
  uint8_t protocolidentifier_nextheader;
  uint16_t singlelocalport;
  port_range_t localportrange;
  uint16_t singleremoteport;
  port_range_t remoteportrange;
  uint32_t securityparameterindex;
  struct {
    uint8_t value;
    uint8_t mask;
  } typdeofservice_trafficclass;
  uint32_t flowlabel;
} packet_filter_contents_t;

/*
 * Packet filter list when the TFP operation is "delete existing TFT"
 * and "no TFT operation" shall be empty.
 * ---------------------------------------------------------------
 * The empty struct generates -Wextern-c-compat warnings, as the sizeof
 * an empty struct is zero in c and one byte in c++.
 * As the struct use of this empty struct is not intended to be byte-
 * compatible with 3GPP standard, we can add a placeholder uint8_t here
 * to ensure safe behavior in c/c++ mixed codebase.
 */
typedef struct {
  uint8_t unused;
} no_packet_filter_t;

typedef no_packet_filter_t delete_existing_tft_t;
typedef no_packet_filter_t no_tft_operation_t;

/*
 * Packet filter list when the TFT operation is "delete existing TFT"
 * shall contain a variable number of packet filter identifiers.
 * ------------------------------------------------------------------
 */
#define TRAFFIC_FLOW_TEMPLATE_PACKET_IDENTIFIER_MAX 16
typedef struct packet_filter_identifier_s {
  uint8_t identifier;
} packet_filter_identifier_t;

typedef packet_filter_identifier_t delete_packet_filter_t;

/*
 * Packet filter list when the TFT operation is "create new TFT",
 * "add packet filters to existing TFT" and "replace packet filters
 * in existing TFT" shall contain a variable number of packet filters
 * ------------------------------------------------------------------
 */
#define TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX 4
typedef struct packet_filter_s {
  uint8_t spare : 2;
#define TRAFFIC_FLOW_TEMPLATE_PRE_REL7_TFT_FILTER 0b00
#define TRAFFIC_FLOW_TEMPLATE_DOWNLINK_ONLY 0b01
#define TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY 0b10
#define TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL 0b11
  uint8_t direction : 2;
  uint8_t identifier : 4;
  uint8_t eval_precedence;
  uint8_t length;
  packet_filter_contents_t packetfiltercontents;
} packet_filter_t;

typedef packet_filter_t create_new_tft_t;
typedef packet_filter_t add_packet_filter_t;
typedef packet_filter_t replace_packet_filter_t;
/*
 * Packet filter list
 * ------------------
 */
typedef union {
  create_new_tft_t createnewtft[TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX];
  add_packet_filter_t
      addpacketfilter[TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX];
  replace_packet_filter_t
      replacepacketfilter[TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX];
  delete_packet_filter_t
      deletepacketfilter[TRAFFIC_FLOW_TEMPLATE_PACKET_IDENTIFIER_MAX];
  delete_existing_tft_t deleteexistingtft;
  no_tft_operation_t notftoperation;
} packet_filter_list_t;

typedef struct parameter_s {
#define PARAMETER_IDENTIFIER_AUTHORIZATION_TOKEN 0x01  // Authorization Token
#define PARAMETER_IDENTIFIER_FLOW_IDENTIFIER 0x02      // Flow Identifier
#define PARAMETER_IDENTIFIER_PACKET_FILTER_IDENTIFIER                          \
  0x03  // Packet Filter Identifier
  uint8_t parameteridentifier;
  uint8_t length;
  bstring contents;
} parameter_t;

typedef struct parameters_list_s {
  uint8_t num_parameters;
#define TRAFFIC_FLOW_TEMPLATE_NB_PARAMETERS_MAX 16  // TODO or may use []
  parameter_t parameter[TRAFFIC_FLOW_TEMPLATE_NB_PARAMETERS_MAX];
} parameters_list_t;

typedef struct traffic_flow_template_s {
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_SPARE 0b000
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT 0b001
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_EXISTING_TFT 0b010
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_ADD_PACKET_FILTER_TO_EXISTING_TFT 0b011
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_REPLACE_PACKET_FILTERS_IN_EXISTING_TFT    \
  0b100
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_PACKET_FILTERS_FROM_EXISTING_TFT   \
  0b101
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_NO_TFT_OPERATION 0b110
#define TRAFFIC_FLOW_TEMPLATE_OPCODE_RESERVED 0b111
  uint8_t tftoperationcode : 3;
#define TRAFFIC_FLOW_TEMPLATE_PARAMETER_LIST_IS_NOT_INCLUDED 0
#define TRAFFIC_FLOW_TEMPLATE_PARAMETER_LIST_IS_INCLUDED 1
  uint8_t ebit : 1;
  uint8_t numberofpacketfilters : 4;
  packet_filter_list_t packetfilterlist;
  parameters_list_t parameterslist;  // The parameters list contains a variable
                                     // number of parameters that may be
  // transferred. If the parameters list is included, the E bit is set to 1;
  // otherwise, the E bit is set to 0.
} traffic_flow_template_t;

#define TFT_ENCODE_IEI_TRUE true
#define TFT_ENCODE_IEI_FALSE false
#define TFT_ENCODE_LENGTH_TRUE true
#define TFT_ENCODE_LENGTH_FALSE false

int encode_traffic_flow_template(
    const traffic_flow_template_t* trafficflowtemplate, uint8_t* buffer,
    const uint32_t len);
int encode_traffic_flow_template_ie(
    const traffic_flow_template_t* const trafficflowtemplate,
    const bool iei_present, uint8_t* buffer, const uint32_t len);
int decode_traffic_flow_template(
    traffic_flow_template_t* trafficflowtemplate, const uint8_t* const buffer,
    const uint32_t len);
int decode_traffic_flow_template_ie(
    traffic_flow_template_t* trafficflowtemplate, const bool iei_present,
    const uint8_t* const buffer, const uint32_t len);
void copy_traffic_flow_template(
    traffic_flow_template_t* const tft_dst,
    const traffic_flow_template_t* const tft_src);
void free_traffic_flow_template(traffic_flow_template_t** tft);

//******************************************************************************
// 10.5.7 GPRS Common information elements
//******************************************************************************

typedef enum gprs_common_ie_e {
  GPRS_C_TIMER_3402_VALUE_IEI          = 0x17, /* 0x17 = 23 */
  GPRS_C_TIMER_3423_VALUE_IEI          = 0x59, /* 0x59 = 89 */
  GPRS_C_TIMER_3412_VALUE_IEI          = 0x5A, /* 0x5A = 90 */
  GPRS_C_TIMER_3412_EXTENDED_VALUE_IEI = 0x5E, /* 0x5E = 94 */
} gprs_common_ie_t;

//------------------------------------------------------------------------------
// 10.5.7.3 GPRS Timer
//------------------------------------------------------------------------------
#define GPRS_TIMER_IE_TYPE 3
#define GPRS_TIMER_IE_MIN_LENGTH 2
#define GPRS_TIMER_IE_MAX_LENGTH 2

typedef struct gprs_timer_s {
#define GPRS_TIMER_UNIT_2S 0b000   /* 2 seconds  */
#define GPRS_TIMER_UNIT_60S 0b001  /* 1 minute */
#define GPRS_TIMER_UNIT_360S 0b010 /* decihours  */
#define GPRS_TIMER_UNIT_0S 0b111   /* deactivated  */
  uint8_t unit : 3;
  uint8_t timervalue : 5;
} gprs_timer_t;

int encode_gprs_timer_ie(
    gprs_timer_t* gprstimer, uint8_t iei, uint8_t* buffer, const uint32_t len);
int decode_gprs_timer_ie(
    gprs_timer_t* gprstimer, uint8_t iei, uint8_t* buffer, const uint32_t len);
long gprs_timer_value(gprs_timer_t* gprstimer);

#endif /* FILE_3GPP_24_008_SEEN */
