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
#ifndef FILE_SGS_MESSAGES_TYPES_SEEN
#define FILE_SGS_MESSAGES_TYPES_SEEN

#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "common_ies.h"
#include "TrackingAreaIdentity.h"

#define SGSAP_LOCATION_UPDATE_REQ(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.sgsap_location_update_req
#define SGSAP_LOCATION_UPDATE_ACC(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.sgsap_location_update_acc
#define SGSAP_LOCATION_UPDATE_REJ(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.sgsap_location_update_rej
#define SGSAP_SERVICE_REQUEST(mSGpTR) (mSGpTR)->ittiMsg.sgsap_service_request
#define SGSAP_PAGING_REQUEST(mSGpTR) (mSGpTR)->ittiMsg.sgsap_paging_request
#define SGSAP_PAGING_REJECT(mSGpTR) (mSGpTR)->ittiMsg.sgsap_paging_reject
#define SGSAP_UE_UNREACHABLE(mSGpTR) (mSGpTR)->ittiMsg.sgsap_ue_unreachable
#define SGSAP_EPS_DETACH_IND(mSGpTR) (mSGpTR)->ittiMsg.sgsap_eps_detach_ind
#define SGSAP_EPS_DETACH_ACK(mSGpTR) (mSGpTR)->ittiMsg.sgsap_eps_detach_ack
#define SGSAP_IMSI_DETACH_IND(mSGpTR) (mSGpTR)->ittiMsg.sgsap_imsi_detach_ind
#define SGSAP_IMSI_DETACH_ACK(mSGpTR) (mSGpTR)->ittiMsg.sgsap_imsi_detach_ack
#define SGSAP_STATUS(mSGpTR) (mSGpTR)->ittiMsg.sgsap_status
#define SGSAP_TMSI_REALLOC_COMP(mSGpTR)                                        \
  (mSGpTR)->ittiMsg.sgsap_tmsi_realloc_comp
#define SGSAP_MM_INFORMATION_REQ(mSGpTR)                                       \
  (mSGpTR)->ittiMsg.sgsap_mm_information_req
#define SGSAP_ALERT_REQUEST(mSGpTR) (mSGpTR)->ittiMsg.sgsap_alert_request
#define SGSAP_ALERT_ACK(mSGpTR) (mSGpTR)->ittiMsg.sgsap_alert_ack
#define SGSAP_ALERT_REJECT(mSGpTR) (mSGpTR)->ittiMsg.sgsap_alert_reject
#define SGSAP_UPLINK_UNITDATA(mSGpTR) (mSGpTR)->ittiMsg.sgsap_uplink_unitdata
#define SGSAP_DOWNLINK_UNITDATA(mSGpTR)                                        \
  (mSGpTR)->ittiMsg.sgsap_downlink_unitdata
#define SGSAP_RELEASE_REQ(mSGpTR) (mSGpTR)->ittiMsg.sgsap_release_req
#define SGSAP_UE_ACTIVITY_IND(mSGpTR) (mSGpTR)->ittiMsg.sgsap_ue_activity_ind
#define SGSAP_VLR_RESET_INDICATION(mSGpTR)                                     \
  (mSGpTR)->ittiMsg.sgsap_vlr_reset_indication
#define SGSAP_VLR_RESET_ACK(mSGpTR) (mSGpTR)->ittiMsg.sgsap_vlr_reset_ack
#define SGSAP_SERVICE_ABORT_REQ(mSGpTR)                                        \
  (mSGpTR)->ittiMsg.sgsap_service_abort_req

typedef enum SgsCause_e {
  SGS_CAUSE_NORMAL,
  SGS_CAUSE_IMSI_DETACHED_FOR_EPS_SERVICE,
  SGS_CAUSE_IMSI_DETACHED_FOR_EPS_SERVICE_AND_NONEPS_SERVICE,
  SGS_CAUSE_IMSI_UNKNOWN,
  SGS_CAUSE_IMSI_DETACHED_FOR_NONEPS_SERVICE,
  SGS_CAUSE_IMSI_IMPLICITLY_DETACHED_FOR_NONEPS_SERVICE,
  SGS_CAUSE_UE_UNREACHABLE,
  SGS_CAUSE_MESSAGE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE,
  SGS_CAUSE_MISSING_MANDATORY_IE,
  SGS_CAUSE_INVALID_MANDATORY_IE,
  SGS_CAUSE_CONDITIONAL_IE_ERROR,
  SGS_CAUSE_SEMANTICALLY_INCORRECT_IE,
  SGS_CAUSE_MESSAGE_UNKNOWN,
  SGS_CAUSE_MT_CSFB_CALL_REJECTED_BY_USER,
  SGS_CAUSE_UE_TEMPORARILY_UNREACHABLE,
} SgsCause_t;

typedef enum SgsRejectCause_e {
  SGS_INVALID_CAUSE,
  SGS_IMSI_UNKNOWN_IN_HLR                              = 2,
  SGS_ILLEGAL_MS                                       = 3,
  SGS_IMSI_UNKNOWN_IN_VLR                              = 4,
  SGS_IMEI_NOT_ACCEPTED                                = 5,
  SGS_ILLEGAL_UE                                       = 6,
  SGS_PLMN_NOT_ALLOWED                                 = 11,
  SGS_LOCATION_AREA_NOT_ALLOWED                        = 12,
  SGS_ROAMING_NOT_ALLOWED_IN_THIS_LOCATION_AREA        = 13,
  SGS_NO_SUITABLE_CELLS_IN_LOCATION_AREA               = 15,
  SGS_MSC_NOT_REACHABLE                                = 16,
  SGS_NETWORK_FAILURE                                  = 17,
  SGS_MAC_FAILURE                                      = 20,
  SGS_SYNCH_FAILURE                                    = 21,
  SGS_CONGESTION                                       = 22,
  SGS_GSM_AUTHENTICATION_UNACCEPTABLE                  = 23,
  SGS_NOT_AUTHORIZED_FOR_THIS_CSG                      = 25,
  SGS_SERVICE_OPTION_NOT_SUPPORTED                     = 32,
  SGS_REQUESTED_SERVICE_OPTION_NOT_SUBSCRIBED          = 33,
  SGS_SERVICE_OPTION_TEMPORARILY_OUT_OF_ORDER          = 34,
  SGS_CALL_CANNOT_BE_IDENTIFIED                        = 38,
  SGS_RETRY_UPON_ENTRY_INTO_NEW_CELL                   = 48,
  SGS_SEMANTICALLY_INCORRECT_MESSAGE                   = 95,
  SGS_INVALID_MANDATORY_INFORMATION                    = 96,
  SGS_MSG_TYPE_NON_EXISTENT_NOT_IMPLEMENTED            = 97,
  SGS_MSG_TYPE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE      = 98,
  SGS_INFORMATION_ELEMENT_NON_EXISTENT_NOT_IMPLEMENTED = 99,
  SGS_CONDITIONAL_IE_ERROR                             = 100,
  SGS_MSG_NOT_COMPATIBLE_WITH_PROTOCOL_STATE           = 101,
  SGS_PROTOCOL_ERROR_UNSPECIFIED                       = 111
} SgsRejectCause_t;

typedef enum LocationUpdateType_e {
  IMSI_ATTACH            = 1,
  NORMAL_LOCATION_UPDATE = 2
} LocationUpdateType_t;

/*
 * As per specification 29.118 Section 9.4.7,
 * IMSI detach from EPS service type is as below:
 * 0 0 0 0 0 0 0 0 Interpreted as reserved in this version of the protocol
 * 0 0 0 0 0 0 0 1 Network initiated IMSI detach from EPS services
 * 0 0 0 0 0 0 1 0 UE initiated IMSI detach from EPS services
 * 0 0 0 0 0 0 1 1 EPS services not allowed
 */
typedef enum {
  SGS_EPS_DETACH_TYPE_RESERVED = 0,
  SGS_NW_INITIATED_IMSI_DETACH_FROM_EPS,
  SGS_UE_INITIATED_IMSI_DETACH_FROM_EPS,
  SGS_EPS_SERVICES_NOT_ALLOWED,
  SGS_EPS_DETACH_TYPE_MAX
} SgsEpsDetachType_t;

/*
 * As per specification 29.118 Section 9.4.8,
 * IMSI detach from non-EPS service type is as below:
 * 0 0 0 0 0 0 0 0 Interpreted as reserved in this version of the protocol
 * 0 0 0 0 0 0 0 1 Explicit UE initiated IMSI detach from non-EPS services
 * 0 0 0 0 0 0 1 0 Combined UE initiated IMSI detach from EPS and non-EPS
 * services 0 0 0 0 0 0 1 1 Implicit network initiated IMSI detach from EPS and
 * non-EPS services
 */
typedef enum {
  SGS_NONEPS_DETACH_TYPE_RESERVED = 0,
  SGS_EXPLICIT_UE_INITIATED_IMSI_DETACH_FROM_NONEPS,
  SGS_COMBINED_UE_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS,
  SGS_IMPLICIT_NW_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS,
  SGS_NONEPS_DETACH_TYPE_MAX
} SgsNonEpsDetachType_t;

typedef struct SelectedCsDomainOperator_s {
  uint8_t csfbmccdigit2 : 4;
  uint8_t csfbmccdigit1 : 4;
  uint8_t csfbmncdigit3 : 4;
  uint8_t csfbmccdigit3 : 4;
  uint8_t csfbmncdigit2 : 4;
  uint8_t csfbmncdigit1 : 4;
} SelectedCsDomainOperator_t;

typedef struct itti_sgsap_location_update_req_s {
#define SGSAP_OLD_LAI (1 << 0)
#define SGSAP_TMSI_STATUS (1 << 1)
#define SGSAP_IMEISV (1 << 2)
#define SGSAP_TAI (1 << 3)
#define SGSAP_E_CGI (1 << 4)
#define SGSAP_TMSI_NRI_CONTAINER (1 << 5)
#define SGSAP_SELECTED_CS_DOMAIN_OPERATOR (1 << 6)

  uint8_t presencemask;
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  LocationUpdateType_t locationupdatetype;
  lai_t newlaicsfb;  // New LAI
  lai_t oldlaicsfb;  // Old LAI
  bool tmsistatus;
  char imeisv[MAX_IMEISV_SIZE + 1];
  uint8_t imsisv_length;
  tai_t tai;
  ecgi_t ecgi;
  uint16_t tmsinricontainer;
  SelectedCsDomainOperator_t selectedcsdomainop;
} itti_sgsap_location_update_req_t;

typedef struct itti_sgsap_location_update_acc_s {
#define SGSAP_MOBILE_IDENTITY (1 << 0)
  uint8_t presencemask;
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  lai_t laicsfb;
  MobileIdentity_t mobileid;
  uint8_t additional_updt_type;
} itti_sgsap_location_update_acc_t;

typedef struct itti_sgsap_location_update_rej_s {
#define SGSAP_LAI (1 << 0)
  uint8_t presencemask;
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  SgsRejectCause_t cause;
  lai_t laicsfb;
} itti_sgsap_location_update_rej_t;

typedef struct itti_sgsap_paging_request_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  uint8_t service_indicator; /* Indicates SMS or CS call */

  /* If an optional value is present and should be encoded, the corresponding
   * Bit mask should be set to 1.
   */

#define PAGING_REQUEST_TMSI_PARAMETER_PRESENT (1 << 0)
#define PAGING_REQUEST_CLI_PARAMETER_PRESENT (1 << 1)
#define PAGING_REQUEST_LAI_PARAMETER_PRESENT (1 << 2)
  uint8_t presencemask;
  uint32_t opt_tmsi;
  bstring opt_cli;
  lai_t opt_lai;
} itti_sgsap_paging_request_t;

typedef struct itti_sgsap_service_request_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  uint8_t service_indicator; /* Indicates SMS or CS call */

  /* If an optional value is present and should be encoded, the corresponding
   * Bit mask should be set to 1.
   */

#define SERVICE_REQUEST_IMEISV_PARAMETER_PRESENT (1 << 0)
#define SERVICE_REQUEST_UE_TIMEZONE_PARAMETER_PRESENT (1 << 1)
#define SERVICE_REQUEST_MOBILE_STATION_CLASSMARK_2_PARAMETER_PRESENT (1 << 2)
#define SERVICE_REQUEST_TAI_PARAMETER_PRESENT (1 << 3)
#define SERVICE_REQUEST_ECGI_PARAMETER_PRESENT (1 << 4)
#define SERVICE_REQUEST_UE_EMM_MODE_PARAMETER_PRESENT (1 << 5)
  uint8_t presencemask;
  char opt_imeisv[MAX_IMEISV_SIZE + 1];
  uint8_t opt_imeisv_length;
  TimeZone opt_ue_time_zone;
  MobileStationClassmark2_t opt_mobilestationclassmark2;
  tai_t opt_tai;
  ecgi_t opt_ecgi;
  UeEmmMode opt_ue_emm_mode;
} itti_sgsap_service_request_t;

typedef struct itti_sgsap_paging_reject_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  SgsCause_t sgs_cause; /* Paging Reject SGS Cause */
} itti_sgsap_paging_reject_t;

typedef struct itti_sgsap_ue_unreachable_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  SgsCause_t sgs_cause; /* Paging Reject SGS Cause */
  /* TBD optional IEs*/
} itti_sgsap_ue_unreachable_t;

typedef struct itti_sgsap_eps_detach_ind_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  SgsEpsDetachType_t eps_detach_type; /* IMSI DETACH FOR EPS */
} itti_sgsap_eps_detach_ind_t;

typedef struct itti_sgsap_imsi_detach_ind_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  SgsNonEpsDetachType_t noneps_detach_type; /* IMSI DETACH FOR NON-EPS */
} itti_sgsap_imsi_detach_ind_t;

/*
 * Common structure for below messages since it has common parameter - IMSI
 * SGSAP EPS Detach Ack
 * SGSAP IMSI Detach Ack
 * SGSAP TMSI Reallocation Complete
 * SGSAP ALERT Request for non-eps procedure
 * SGSAP ALERT Ack for non-eps procedure
 * SGSAP SERVICE ABORT Req
 */
typedef struct sgsap_imsi_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
} sgsap_imsi_t;

typedef sgsap_imsi_t itti_sgsap_imsi_detach_ack_t;
typedef sgsap_imsi_t itti_sgsap_eps_detach_ack_t;
typedef sgsap_imsi_t itti_sgsap_tmsi_reallocation_comp_t;
typedef sgsap_imsi_t itti_sgsap_alert_request_t;
typedef sgsap_imsi_t itti_sgsap_alert_ack_t;
typedef sgsap_imsi_t itti_sgsap_service_abort_req_t;

typedef struct error_msg_t {
  uint8_t msg_type;
  union {
    itti_sgsap_location_update_acc_t sgsap_location_update_acc;
    itti_sgsap_location_update_rej_t sgsap_location_update_rej;
  } u;
} error_msg_t;

typedef struct itti_sgsap_status_s {
#define SGSAP_IMSI (1 << 0)
#define ERROR_MESSAGE_LEN_MX 255
  uint8_t presencemask;
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  SgsCause_t cause;
  char error_msg_rcvd[ERROR_MESSAGE_LEN_MX];
  error_msg_t error_msg;
} itti_sgsap_status_t;

/* Common structure form sgsap-paging reject and sgsap-alert reject */
typedef itti_sgsap_paging_reject_t itti_sgsap_alert_reject_t;

typedef struct itti_sgsap_mm_information_req_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD8_SIZE + 1];
/* As per spec,No upper length limit is specified except for that given by the
   maximum number of octets in a L3 message max length of L3 message is 251,
   setting to 255 */
#define MAX_NETWORK_NAME_LEN 255
#define MAX_LSA_IDENTIFIER_LEN 5
#define MM_INFORMATION_REQUEST_FULL_NW_NAME_PRESENT (1 << 0)
#define MM_INFORMATION_REQUEST_SHORT_NW_NAME_PRESENT (1 << 1)
#define MM_INFORMATION_REQUEST_LOCAL_TIME_ZONE_PRESENT (1 << 2)
#define MM_INFORMATION_REQUEST_UNIVERSAL_TIME_AND_TIME_ZONE_PRESENT (1 << 3)
#define MM_INFORMATION_REQUEST_LSA_IDENTIFIER_PRESENT (1 << 4)
#define MM_INFORMATION_REQUEST_NW_DAYLIGHT_SAVING_TIME_PRESENT (1 << 5)
  uint8_t presencemask;
  char full_network_name[MAX_NETWORK_NAME_LEN];
  char short_network_name[MAX_NETWORK_NAME_LEN];
  TimeZone localtimezone;
  TimeZoneAndTime_t universaltimeandlocaltimezone;
  uint8_t lsa_identifier[MAX_LSA_IDENTIFIER_LEN];
  uint8_t networkdaylightsavingtime;
} itti_sgsap_mm_information_req_t;

/* SGSAP UPLINK UNITDATA itti message */
typedef struct itti_sgsap_uplink_unitdata_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  bstring nas_msg_container;
/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define UPLINK_UNITDATA_IMEISV_PARAMETER_PRESENT (1 << 0)
#define UPLINK_UNITDATA_UE_TIMEZONE_PARAMETER_PRESENT (1 << 1)
#define UPLINK_UNITDATA_MOBILE_STATION_CLASSMARK_2_PARAMETER_PRESENT (1 << 2)
#define UPLINK_UNITDATA_TAI_PARAMETER_PRESENT (1 << 3)
#define UPLINK_UNITDATA_ECGI_PARAMETER_PRESENT (1 << 4)
  uint8_t presencemask;
  char opt_imeisv[MAX_IMEISV_SIZE + 1];
  uint8_t opt_imeisv_length;
  TimeZone opt_ue_time_zone;
  MobileStationClassmark2_t opt_mobilestationclassmark2;
  tai_t opt_tai;
  ecgi_t opt_ecgi;
} itti_sgsap_uplink_unitdata_t;

/* SGSAP DOWNLINK UNITDATA itti message */
typedef struct itti_sgsap_downlink_unitdata_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  bstring nas_msg_container;
} itti_sgsap_downlink_unitdata_t;

/* SGSAP RELEASE REQUEST itti message */
typedef struct itti_sgsap_release_req_s {
#define RELEASE_REQ_CAUSE_PARAMETER_PRESENT (1 << 0)
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  uint8_t presencemask;
  SgsCause_t opt_cause;
} itti_sgsap_release_req_t;

/* Common structure form sgsap-paging reject and sgsap-alert reject */

/* SGSAP UE Activity Indication */
typedef struct itti_sgsap_ue_activity_ind_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  /*TBD: Do we required the optional parameter 'Maximum UE Available Time'?
   * since it is based on deployment scenario*/
} itti_sgsap_ue_activity_ind_t;

typedef struct itti_sgsap_vlr_reset_indication_s {
  uint8_t length;
  char vlr_name[MAX_VLR_NAME_LENGTH];
} itti_sgsap_vlr_reset_indication_t;

typedef struct itti_sgsap_vlr_reset_ack_s {
  uint8_t length;
  char mme_name[MAX_MME_NAME_LENGTH];
} itti_sgsap_vlr_reset_ack_t;

#endif /* FILE_SGS_MESSAGES_TYPES_SEEN */
