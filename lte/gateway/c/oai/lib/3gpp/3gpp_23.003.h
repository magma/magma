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

/*! \file 3gpp_23.003.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_3GPP_23_003_SEEN
#define FILE_3GPP_23_003_SEEN

//==============================================================================
// 12  Identification of PLMN, RNC, Service Area, CN domain and Shared Network
// Area
//==============================================================================
//------------------------------------------------------------------------------
// 12.1  PLMN Identifier
//------------------------------------------------------------------------------
/*
 * A Public Land Mobile Network is uniquely identified by its PLMN identifier.
 * PLMN-Id consists of Mobile Country Code (MCC) and Mobile Network Code (MNC).
 */
/*
 * Public Land Mobile Network identifier
 * PLMN = BCD encoding (Mobile Country Code + Mobile Network Code)
 */
/*! \struct  plmn_t
 * \brief Public Land Mobile Network identifier.
 *        PLMN = BCD encoding (Mobile Country Code + Mobile Network Code).
 */
typedef struct plmn_s {
  uint8_t mcc_digit2 : 4;
  uint8_t mcc_digit1 : 4;
  uint8_t mnc_digit3 : 4;
  uint8_t mcc_digit3 : 4;
  uint8_t mnc_digit2 : 4;
  uint8_t mnc_digit1 : 4;
} plmn_t;

//------------------------------------------------------------------------------
// 12.2  CN Domain Identifier
// 12.3  CN Identifier
// 12.4  RNC Identifier
// 12.5  Service Area Identifier
// 12.6 Shared Network Area Identifier
//------------------------------------------------------------------------------

//==============================================================================
// 2 Identification of mobile subscribers
//==============================================================================

//------------------------------------------------------------------------------
// 2.2 Composition of IMSI
//------------------------------------------------------------------------------

/*! \struct  imsi_t
 * \brief Structure containing an IMSI, BCD structure.
 */
typedef struct imsi_s {
  union {
    struct {
      uint8_t digit2 : 4;
      uint8_t digit1 : 4;
      uint8_t digit4 : 4;
      uint8_t digit3 : 4;
      uint8_t digit6 : 4;
      uint8_t digit5 : 4;
      uint8_t digit8 : 4;
      uint8_t digit7 : 4;
      uint8_t digit10 : 4;
      uint8_t digit9 : 4;
      uint8_t digit12 : 4;
      uint8_t digit11 : 4;
      uint8_t digit14 : 4;
      uint8_t digit13 : 4;
#define EVEN_PARITY 0
#define ODD_PARITY 1
      uint8_t parity : 4;
      uint8_t digit15 : 4;
    } num; /*!< \brief  IMSI shall consist of decimal digits (0 through 9)
              only.*/
#define IMSI_BCD8_SIZE                                                         \
  8 /*!< \brief  The number of digits in IMSI shall not exceed 15.       */
    uint8_t value[IMSI_BCD8_SIZE];
  } u;
  uint8_t length;
} imsi_t;

#define IMSI_BCD_DIGITS_MAX 15
typedef struct {
  uint8_t digit[IMSI_BCD_DIGITS_MAX + 1];  // +1 for '\0` macro sprintf changed
                                           // in snprintf
  uint8_t length;

} Imsi_t;
//------------------------------------------------------------------------------
// 2.4 Structure of TMSI
//------------------------------------------------------------------------------

#define TMSI_SIZE 4
typedef uint32_t
    tmsi_t; /*!< \brief  Since the TMSI has only local significance (i.e. within
               a VLR and the area controlled by a VLR, or within an SGSN and the
               area controlled by an SGSN, or within an MME and the area
               controlled by an MME), the structure and coding of it can be
               chosen by agreement between operator and manufacturer in order to
               meet local needs.
                                                                          The
               TMSI consists of 4 octets. It can be coded using a full
               hexadecimal representation. */

#define INVALID_TMSI                                                           \
  UINT32_MAX /*!< \brief  The network shall not allocate a TMSI with all 32    \
                bits equal to 1 (this is because the TMSI must be stored in    \
                the SIM, and the SIM uses 4 octets with all bits               \
                                                                        equal  \
                to 1 to indicate that no valid TMSI is available).  */

// 2.5 Structure of LMSI
// 2.6  Structure of TLLI
// 2.7 Structure of P-TMSI Signature

//------------------------------------------------------------------------------
// 2.8 Globally Unique Temporary UE Identity (GUTI)
//------------------------------------------------------------------------------

#define INVALID_M_TMSI                                                         \
  UINT32_MAX /*!< \brief  The network shall not allocate a TMSI with all 32    \
                bits equal to 1 (this is because the TMSI must be stored in    \
                the SIM, and the SIM uses 4 octets with all bits               \
                                                                        equal  \
                to 1 to indicate that no valid TMSI is available).  */

typedef uint16_t
    mme_gid_t; /*!< \brief  MME Group ID shall be of 16 bits length. */
typedef uint8_t
    mme_code_t; /*!< \brief  MME Code shall be of 8 bits length.      */

/*! \struct  gummei_t
 * \brief Structure containing the Globally Unique MME Identity.
 */
typedef struct gummei_s {
  plmn_t plmn;         /*!< \brief  GUMMEI               */
  mme_gid_t mme_gid;   /*!< \brief  MME group identifier */
  mme_code_t mme_code; /*!< \brief  MME code             */
} gummei_t;

/*! \struct  guti_t
 * \brief Structure containing the Globally Unique Temporary UE Identity.
 */
typedef struct guti_s {
  gummei_t gummei; /*!< \brief  Globally Unique MME Identity             */
  tmsi_t m_tmsi;   /*!< \brief  M-Temporary Mobile Subscriber Identity   */
} guti_t;

// 2.9 Structure of the S-Temporary Mobile Subscriber Identity (S-TMSI)

/*! \struct  s_tmsi_t
 * \brief Structure of the S-Temporary Mobile Subscriber Identity (S-TMSI).
 */
typedef struct s_tmsi_s {
  mme_code_t mme_code; /* MME code that allocated the GUTI     */
  tmsi_t m_tmsi;       /* M-Temporary Mobile Subscriber Identity   */
} s_tmsi_t;
//==================================================================================
//----------------- 5G Globally Unique Temporary UE Identity (GUTI)-------------
typedef uint16_t
    amf_gid_t; /*!< \brief  AMF Group ID shall be of 16 bits length. */
typedef uint8_t
    amf_code_t; /*!< \brief  AMF Code shall be of 8 bits length.      */
typedef uint8_t amf_Pointer_t;  // 9.3.3.19 AMF Pointer is used to identify one
                                // or more AMF(s) within the AMF Set.
/*! \struct  guamfi_t
 * \brief Structure containing the Globally Unique AMF Identity.
 */
typedef struct guamfi_s {
  plmn_t plmn;         /*!< \brief  GUAMFI               */
  amf_gid_t amf_gid;   /*!< \brief  AMF group identifier */
  amf_code_t amf_code; /*!< \brief  AMF code             */
  amf_Pointer_t amf_Pointer;
} guamfi_t;
typedef struct guti_m5_s {
  guamfi_t guamfi; /*!< \brief  Globally Unique AMF Identity             */
  tmsi_t m_tmsi;   /*!< \brief  M-Temporary Mobile Subscriber Identity   */
} guti_m5_t;
typedef struct s_tmsi_m5_s {
  amf_code_t amf_code; /* AMF code that allocated the GUTI     */
  tmsi_t m_tmsi;       /* M-Temporary Mobile Subscriber Identity   */
} s_tmsi_m5_t;

//==============================================================================
// 3 Numbering plan for mobile stations
//==============================================================================

//------------------------------------------------------------------------------
// 3.3 Structure of MS international PSTN/ISDN number (MSISDN)
//------------------------------------------------------------------------------

/*! \struct  msisdn_t
 * \brief MS international PSTN/ISDN number (MSISDN).
 */
typedef struct msisdn_s {
  uint8_t ext : 1;
  /* Type Of Number           */
#define MSISDN_TON_UNKNOWKN 0b000
#define MSISDN_TON_INTERNATIONAL 0b001
#define MSISDN_TON_NATIONAL 0b010
#define MSISDN_TON_NETWORK 0b011
#define MSISDN_TON_SUBCRIBER 0b100
#define MSISDN_TON_ABBREVIATED 0b110
#define MSISDN_TON_RESERVED 0b111
  uint8_t ton : 3;
  /* Numbering Plan Identification    */
#define MSISDN_NPI_UNKNOWN 0b0000
#define MSISDN_NPI_ISDN_TELEPHONY 0b0001
#define MSISDN_NPI_GENERIC 0b0010
#define MSISDN_NPI_DATA 0b0011
#define MSISDN_NPI_TELEX 0b0100
#define MSISDN_NPI_MARITIME_MOBILE 0b0101
#define MSISDN_NPI_LAND_MOBILE 0b0110
#define MSISDN_NPI_ISDN_MOBILE 0b0111
#define MSISDN_NPI_PRIVATE 0b1110
#define MSISDN_NPI_RESERVED 0b1111
  uint8_t npi : 4;
  /* Dialing Number           */
  struct {
    uint8_t lsb : 4;
    uint8_t msb : 4;
#define MSISDN_DIGIT_SIZE 10
  } digit[MSISDN_DIGIT_SIZE];
} msisdn_t;

// 3.4 Mobile Station Roaming Number (MSRN) for PSTN/ISDN routeing
// 3.5  Structure of Mobile Station International Data Number
// 3.6  Handover Number
// 3.7  Structure of an IP v4 address
// 3.8  Structure of an IP v6 address

//==============================================================================
// 4 Identification of location areas and base stations
//==============================================================================

//------------------------------------------------------------------------------
// 4.1 Composition of the Location Area Identification (LAI)
//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
// 4.2 Composition of the Routing Area Identification (RAI)
//------------------------------------------------------------------------------

// 4.3 Base station identification
// 4.3.1 Cell Identity (CI) and Cell Global Identification (CGI)
// 4.3.2 Base Station Identify Code (BSIC)
// 4.4  Regional Subscription Zone Identity (RSZI)
// 4.5  Location Number
// 4.6  Composition of the Service Area Identification (SAI)

//------------------------------------------------------------------------------
// 4.7  Closed Subscriber Group
//------------------------------------------------------------------------------
typedef uint32_t
    csg_id_t; /*!< \brief  The CSG‑ID shall be fix length 27 bit value. */

// 4.8 HNB Name
// 4.9 CSG Type
// 4.10  HNB Unique Identity
// 5 Identification of MSCs, GSNs and location registers

//==============================================================================
// 6 International Mobile Station Equipment Identity and Software Version Number
//==============================================================================

//------------------------------------------------------------------------------
// 6.2.1 Composition of IMEI
//------------------------------------------------------------------------------

/*! \struct  imei_t
 * \brief Structure containing an International Mobile station Equipment
 * Identity, BCD structure. The IMEI is composed of the following elements (each
 * element shall consist of decimal digits only):
 *        - Type Allocation Code (TAC). Its length is 8 digits;
 *        - Serial Number (SNR) is an individual serial number uniquely
 * identifying each equipment within the TAC. Its length is 6 digits;
 *        -  Check Digit (CD) / Spare Digit (SD): If this is the Check Digit see
 * paragraph below; if this digit is Spare Digit it shall be set to zero, when
 * transmitted by the MS.
 *
 *        The IMEI (14 digits) is complemented by a Check  Digit (CD). The Check
 * Digit is not part of the digits transmitted when the IMEI is checked, as
 * described below. The Check Digit is intended to avoid manual transmission
 * errors, e.g. when customers register stolen MEs at the operator's customer
 * care desk. The Check Digit is defined according to the Luhn formula, as
 * defined in annex B.
 *
 *        NOTE: The Check Digit is not applied to the Software Version Number.
 */
typedef struct imei_s {
  uint8_t length;
  union {
    struct {
      uint8_t tac2 : 4;
      uint8_t tac1 : 4;
      uint8_t tac4 : 4;
      uint8_t tac3 : 4;
      uint8_t tac6 : 4;
      uint8_t tac5 : 4;
      uint8_t tac8 : 4;
      uint8_t tac7 : 4;
      uint8_t snr2 : 4;
      uint8_t snr1 : 4;
      uint8_t snr4 : 4;
      uint8_t snr3 : 4;
      uint8_t snr6 : 4;
      uint8_t snr5 : 4;
#define EVEN_PARITY 0
#define ODD_PARITY 1
      uint8_t parity : 4;
      uint8_t cdsd : 4;
    } num;
#define IMEI_BCD8_SIZE 8
    uint8_t value[IMEI_BCD8_SIZE];
  } u;
} imei_t;

//------------------------------------------------------------------------------
// 6.2.2 Composition of IMEISV
//------------------------------------------------------------------------------

/*! \struct  imeisv_t
 * \brief Structure containing an International Mobile station Equipment
 * Identity, BCD structure. The IMEISV is composed of the following elements
 * (each element shall consist of decimal digits only):
 *        - Type Allocation Code (TAC). Its length is 8 digits;
 *        - Serial Number (SNR) is an individual serial number uniquely
 * identifying each equipment within each TAC. Its length is 6 digits;
 *        - Software Version Number (SVN) identifies the software version number
 * of the mobile equipment. Its length is 2 digits. Regarding updates of the
 * IMEISV: The security requirements of 3GPP TS 22.016 [32] apply only to the
 * TAC and SNR, but not to the SVN part of the IMEISV.
 */
typedef struct imeisv_s {
  uint8_t length;
  union {
    struct {
      uint8_t tac2 : 4;
      uint8_t tac1 : 4;
      uint8_t tac4 : 4;
      uint8_t tac3 : 4;
      uint8_t tac6 : 4;
      uint8_t tac5 : 4;
      uint8_t tac8 : 4;
      uint8_t tac7 : 4;
      uint8_t snr2 : 4;
      uint8_t snr1 : 4;
      uint8_t snr4 : 4;
      uint8_t snr3 : 4;
      uint8_t snr6 : 4;
      uint8_t snr5 : 4;
      uint8_t svn2 : 4;
      uint8_t svn1 : 4;
#define EVEN_PARITY 0
#define ODD_PARITY 1
      uint8_t parity : 4;
    } num;
#define IMEISV_BCD8_SIZE 9
    uint8_t value[IMEISV_BCD8_SIZE];
  } u;
} imeisv_t;

// 7 Identification of Voice Group Call and Voice Broadcast Call Entities
// 8 SCCP subsystem numbers

//==============================================================================
// 9 Definition of Access Point Name
//==============================================================================
// TODO

// 10  Identification of the Cordless Telephony System entities
// 11  Identification of Localised Service Area
// 12 -> section moved to the top of this file.
// 13  Numbering, addressing and identification within the IP multimedia core
// network subsystem 14  Numbering, addressing and identification for 3GPP
// System to WLAN Interworking 15 Identification of Multimedia
// Broadcast/Multicast Service 15.2  Structure of TMGI
// TODO NAS ?
// 15.3  Structure of MBMS SAI
// 15.4  Home Network Realm
// 16  Numbering, addressing and identification within the GAA subsystem
// 17  Numbering, addressing and identification within the Generic Access
// Network 18  Addressing and Identification for IMS Service Continuity and
// Single-Radio Voice Call Continuity

//==============================================================================
// 19  Numbering, addressing and identification for the Evolved Packet Core
// (EPC)
//==============================================================================

//------------------------------------------------------------------------------
// 19.4  Identifiers for Domain Name System procedures
//------------------------------------------------------------------------------

//..............................................................................
// 19.4.2  Fully Qualified Domain Names (FQDNs)
//..............................................................................
// 19.4.2.2 Access Point Name FQDN (APN-FQDN)

// 19.4.2.3  Tracking Area Identity (TAI)

// 19.4.2.4  Mobility Management Entity (MME)
// 19.4.2.5  Routing Area Identity (RAI) - EPC

//------------------------------------------------------------------------------
// 19.6  E-UTRAN Cell Identity (ECI) and E-UTRAN Cell Global Identification
// (ECGI)
//------------------------------------------------------------------------------
/*! \struct  eci_t
 * \brief The ECI shall be of fixed length of 28 bits and shall be coded using
 *        full hexadecimal representation. The exact coding of the ECI is the
 * responsibility of each PLMN operator. */
typedef struct eci_s {
  uint32_t enb_id : 20;
  /* Anoop - This is correct only when eNB type is macro. In case eNB type is
   * Home eNB then all the 28 bits are used for eNB id . This needs
   * correction since MME uses combination of enb_id and "eNB S1AP UEid" for the
   * key to UE context,this may not work if MME is connected to many HeNBs -
   * which is not critical now.*/

  uint32_t cell_id : 8;
  uint32_t empty : 4;
} eci_t;

/*! \struct  ecgi_t
 * \brief The E-UTRAN Cell Global Identification (ECGI) shall be composed of the
 * concatenation of the PLMN Identifier (PLMN-Id) and the E-UTRAN Cell Identity
 * (ECI) .
 */
typedef struct ecgi_s {
  plmn_t plmn;
  eci_t
      cell_identity; /*!< \brief  The ECI shall be of fixed length of 28 bits */
} ecgi_t;

// 20  Addressing and Identification for IMS Centralized Services
// 21 Addressing and Identification for Dual Stack Mobile IPv6 (DSMIPv6)
// 22 Addressing and identification for ANDSF
// 23 Numbering, addressing and identification for the Relay Node OAM System

/* Clear GUTI without free it */
void clear_guti(guti_t* const guti);
/* Clear IMSI without free it */
void clear_imsi(imsi_t* const imsi);
/* Clear IMEI without free it */
void clear_imei(imei_t* const imei);
/* Clear IMEISV without free it */
void clear_imeisv(imeisv_t* const imeisv);

#endif /* FILE_3GPP_23_003_SEEN */
