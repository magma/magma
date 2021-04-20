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

/*! \file 3gpp_requirements_24.301.h
   \brief
   \author  Lionel GAUTHIER
   \date
   \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_3GPP_REQUIREMENTS_24_301_SEEN
#define FILE_3GPP_REQUIREMENTS_24_301_SEEN

#include "3gpp_requirements.h"
#include "log.h"

#define REQUIREMENT_3GPP_24_301(rElEaSe_sEcTiOn__OaImark)                      \
  REQUIREMENT_3GPP_SPEC(                                                       \
      LOG_NAS, "Hit 3GPP TS 24_301" #rElEaSe_sEcTiOn__OaImark                  \
               " : " rElEaSe_sEcTiOn__OaImark##_BRIEF "\n")
#define NO_REQUIREMENT_3GPP_24_301(rElEaSe_sEcTiOn__OaImark)                   \
  REQUIREMENT_3GPP_SPEC(                                                       \
      LOG_NAS, "#NOT IMPLEMENTED 3GPP TS 24_301" #rElEaSe_sEcTiOn__OaImark     \
               " : " rElEaSe_sEcTiOn__OaImark##_BRIEF "\n")
#define NOT_REQUIREMENT_3GPP_24_301(rElEaSe_sEcTiOn__OaImark)                  \
  REQUIREMENT_3GPP_SPEC(                                                       \
      LOG_NAS, "#NOT ASSERTED 3GPP TS 24_301" #rElEaSe_sEcTiOn__OaImark        \
               " : " rElEaSe_sEcTiOn__OaImark##_BRIEF "\n")

//-----------------------------------------------------------------------------------------------------------------------

#define R10_4_4_4_3__1                                                         \
  "MME24.301R10_4.4.4.3_1: Integrity checking of NAS signalling messages in the MME\
                                                                                                                        \
    Except the messages listed below, no NAS signalling messages shall be processed by the receiving EMM entity in the  \
    MME or forwarded to the ESM entity, unless the secure exchange of NAS messages has been established for the NAS     \
    signalling connection:                                                                                              \
    - EMM messages:                                                                                                     \
        - ATTACH REQUEST;                                                                                               \
        - IDENTITY RESPONSE (if requested identification parameter is IMSI);                                            \
        - AUTHENTICATION RESPONSE;                                                                                      \
        - AUTHENTICATION FAILURE;                                                                                       \
        - SECURITY MODE REJECT;                                                                                         \
        - DETACH REQUEST;                                                                                               \
        - DETACH ACCEPT;                                                                                                \
        - TRACKING AREA UPDATE REQUEST.                                                                                 \
    NOTE 1: The TRACKING AREA UPDATE REQUEST message is sent by the UE without integrity protection, if                 \
    the tracking area updating procedure is initiated due to an inter-system change in idle mode and no                 \
    current EPS security context is available in the UE. The other messages are accepted by the MME                     \
    without integrity protection, as in certain situations they are sent by the UE before security can be               \
    activated.                                                                                                          \
                                                                                                                        \
    All ESM messages are integrity protected except a PDN CONNECTIVITY REQUEST message if it is sent                    \
    piggybacked in ATTACH REQUEST message and NAS security is not activated.                                           "
#define R10_4_4_4_3__1_BRIEF                                                   \
  "Integrity checking of NAS signalling messages exception in the MME"

#define R10_4_4_4_3__2                                                         \
  "MME24.301R10_4.4.4.3_2: Integrity checking of NAS signalling messages in the MME                \
                                                                                                                        \
    Once a current EPS security context exists, until the secure exchange of NAS messages has been established for the  \
    NAS signalling connection, the receiving EMM entity in the MME shall process the following NAS signalling           \
    messages, even if the MAC included in the message fails the integrity check or cannot be verified, as the EPS security  \
    context is not available in the network:                                                                            \
    - ATTACH REQUEST;                                                                                                   \
    - IDENTITY RESPONSE (if requested identification parameter is IMSI);                                                \
    - AUTHENTICATION RESPONSE;                                                                                          \
    - AUTHENTICATION FAILURE;                                                                                           \
    - SECURITY MODE REJECT;                                                                                             \
    - DETACH REQUEST (if sent before security has been activated);                                                      \
    - DETACH ACCEPT;                                                                                                    \
    - TRACKING AREA UPDATE REQUEST;                                                                                     "
#define R10_4_4_4_3__2_BRIEF                                                   \
  "Process NAS signalling message in the MME, even if it fails the integrity " \
  "check or MAC cannot be verified"

//-----------------------------------------------------------------------------------------------------------------------
// GUTI REALLOCATION
//-----------------------------------------------------------------------------------------------------------------------

#define R10_5_4_1_6_c                                                          \
  "GUTI reallocation and attach procedure collision                                                 \
    If the network receives an ATTACH REQUEST message before the ongoing GUTI reallocation procedure has                \
been completed the network shall proceed with the attach procedure after deletion of the EMM context."
#define R10_5_4_1_6_c_BRIEF "GUTI reallocation and attach procedure collision"

//-----------------------------------------------------------------------------------------------------------------------
// AUTHENTICATION
//-----------------------------------------------------------------------------------------------------------------------

#define R10_5_4_2_2                                                            \
  "Authentication initiation by the network                                                           \
    When a NAS signalling connection exists, the network can initiate an authentication procedure at any time. For      \
    restrictions applicable after handover or inter-system handover to S1 mode see subclause 5.5.3.2.3.                 \
    The network initiates the authentication procedure by sending an AUTHENTICATION REQUEST message to the UE           \
    and starting the timer T3460 (see example in figure 5.4.2.2.1). The AUTHENTICATION REQUEST message contains         \
    the parameters necessary to calculate the authentication response (see 3GPP TS 33.401 [19])."
#define R10_5_4_2_2_BRIEF "Authentication initiation by the network"

//------------------------------
#define R10_5_4_2_4__1                                                         \
  "Authentication completion by the network                                                        \
    Upon receipt of an AUTHENTICATION RESPONSE message, the network stops the timer T3460 and checks the                \
    correctness of RES (see 3GPP TS 33.401 [19])."
#define R10_5_4_2_4__1_BRIEF                                                   \
  "AUTHENTICATION RESPONSE received, stop T3460, check RES"

#define R10_5_4_2_4__2                                                         \
  "Authentication completion by the network                                                        \
    If the authentication procedure has been completed successfully and the related eKSI is stored in the EPS security  \
    context of the network, the network shall include a different eKSI value in the AUTHENTICATION REQUEST              \
    message when it initiates a new authentication procedure."
#define R10_5_4_2_4__2_BRIEF                                                   \
  "authentication procedure is success, new eKSI for new authentication "      \
  "procedure"

#define R10_5_4_2_4__3                                                         \
  "Authentication completion by the network                                                        \
    Upon receipt of an AUTHENTICATION FAILURE message, the network stops the timer T3460. In the case where the         \
    EMM cause #21 \"synch failure\" is received, the core network may renegotiate with the HSS/AuC and provide the UE   \
    with new authentication parameters."
#define R10_5_4_2_4__3_BRIEF                                                   \
  "AUTHENTICATION FAILURE received with EMM cause sync failure, renegociate "  \
  "with HSS."

//------------------------------
#define R10_5_4_2_5__1                                                         \
  "Authentication not accepted by the network                                                       \
    If the authentication response returned by the UE is not valid, the network response depends upon the type of identity \
    used by the UE in the initial NAS message, that is:                                                                 \
    - if the GUTI was used; or                                                                                          \
    - if the IMSI was used.                                                                                             \
    If the GUTI was used, the network should initiate an identification procedure. If the IMSI given by the UE during the\
    identification procedure differs from the IMSI the network had associated with the GUTI, the authentication should be\
    restarted with the correct parameters. Otherwise, if the IMSI provided by the UE is the same as the IMSI stored in the\
    network (i.e. authentication has really failed), the network should proceed as described below."
#define R10_5_4_2_5__1_BRIEF "AUTHENTICATION RESPONSE not accepted"

#define R10_5_4_2_5__2                                                         \
  "Authentication not accepted by the network                                                       \
    If the IMSI was used for identification in the initial NAS message, or the network decides not to initiate the       \
    identification procedure after an unsuccessful authentication procedure, the network should send an                  \
    AUTHENTICATION REJECT message to the UE."
#define R10_5_4_2_5__2_BRIEF                                                   \
  "AUTHENTICATION RESPONSE not accepted send AUTHENTICATION REJECT"

// Upon receipt of an AUTHENTICATION REJECT message, the UE shall set the update
// status to EU3 ROAMING NOT ALLOWED, delete the stored GUTI, TAI list, last
// visited registered TAI and KSI ASME . The USIM shall be considered invalid
// until switching off the UE or the UICC containing the USIM is removed.
//
// If A/Gb or Iu mode is supported by the UE, the UE shall in addition handle
// the GMM parameters GMM state, GPRS update status, P-TMSI, P-TMSI signature,
// RAI and GPRS ciphering key sequence number and the MM parameters update
// status, TMSI, LAI and ciphering key sequence number as specified in 3GPP
// TS 24.008 [13] for the case when the authentication and ciphering procedure
// is not accepted by the network.
//
// If the AUTHENTICATION REJECT message is received by the UE, the UE shall
// abort any EMM signalling procedure, stop any of the timers T3410, T3417 or
// T3430 (if running) and enter state EMM-DEREGISTERED.
//
// Depending on local requirements or operator preference for emergency bearer
// services, if the UE has a PDN connection for emergency bearer services
// established or is establishing a PDN connection for emergency bearer
// services, the MME need not follow the procedures specified for the
// authentication failure in the present subclause. The MME may continue a
// current EMM specific procedure or PDN connectivity request procedure. Upon
// completion of the authentication procedure, if not initiated as part of
// another procedure, or upon completion of the EMM procedure or PDN
// connectivity request procedure, the MME shall deactivate all non-emergency
// EPS bearers, if any, by initiating an EPS bearer context deactivation
// procedure. The network shall consider the UE to be attached for emergency
// bearer services only.

//------------------------------
#define R10_5_4_2_7_a                                                          \
  "Abnormal cases                                                                                   \
    Lower layer failure:                                                                                                \
    Upon detection of lower layer failure before the AUTHENTICATION RESPONSE is received, the network                   \
    shall abort the procedure."
#define R10_5_4_2_7_a_BRIEF "Lower layer failure"

#define R10_5_4_2_7_b                                                          \
  "Abnormal cases                                                                                   \
    Expiry of timer T3460:                                                                                              \
    The network shall, on the first expiry of the timer T3460, retransmit the AUTHENTICATION REQUEST                    \
    message and shall reset and start timer T3460. This retransmission is repeated four times, i.e. on the fifth expiry \
    of timer T3460, the network shall abort the authentication procedure and any ongoing EMM specific procedure         \
    and release the NAS signalling connection."
#define R10_5_4_2_7_b_BRIEF "Expiry of timer T3460"

/*#define R10_5_4_2_7_c__1 "Abnormal cases \
    Authentication failure (EMM cause #20 \"MAC failure\"): \
    The UE shall send an AUTHENTICATION FAILURE message, with EMM cause #20
   \"MAC failure\" according                   \
    to subclause 5.4.2.6, to the network and start timer T3418 (see example in
   figure 5.4.2.7.1). Furthermore, the UE   \ shall stop any of the
   retransmission timers that are running (e.g. T3410, T3417, T3421 or
   T3430)."*/

#define R10_5_4_2_7_c__2                                                       \
  "Abnormal cases                                                                                \
    Authentication failure (EMM cause #20 \"MAC failure\"):                                                             \
    Upon the first receipt of an AUTHENTICATION FAILURE message from the UE with EMM cause #20 \"MAC failure\", the     \
    network may initiate the identification procedure described in subclause 5.4.4. This is to allow the network to     \
    obtain the IMSI from the UE. The network may then check that the GUTI originally used in the authentication         \
    challenge corresponded to the correct IMSI. Upon receipt of the IDENTITY REQUEST message from the                   \
    network, the UE shall send the IDENTITY RESPONSE message."
#define R10_5_4_2_7_c__2_BRIEF                                                 \
  "AUTHENTICATION FAILURE (EMM cause #20 \"MAC failure\")"

#define R10_5_4_2_7_c__NOTE1                                                   \
  "Abnormal cases                                                                            \
    Upon receipt of an AUTHENTICATION FAILURE message from the UE with EMM cause #20 \"MAC failure\", the network may   \
    also terminate the authentication procedure (see subclause 5.4.2.5)."
#define R10_5_4_2_7_c__NOTE1_BRIEF "Terminate the authentication procedure"

/* Case implemented but have to think how to trace it
#define R10_5_4_2_7_c__3 "Abnormal cases \
    If the GUTI/IMSI mapping in the network was incorrect, the network should
respond by sending a new                  \
    AUTHENTICATION REQUEST message to the UE. Upon receiving the new
AUTHENTICATION REQUEST                             \
    message from the network, the UE shall stop the timer T3418, if running, and
then process the challenge             \ information as normal."*/

// If the network is validated successfully (an AUTHENTICATION REQUEST that
// contains a valid SQN and MAC is received), the UE shall send the
// AUTHENTICATION RESPONSE message to the network and shall start any
// retransmission timers (e.g. T3410, T3417, T3421 or T3430) if they were
// running and stopped when the UE received the first failed AUTHENTICATION
// REQUEST message.

// If the UE receives the second AUTHENTICATION REQUEST while T3418 is running,
// and the MAC value cannot be resolved, the UE shall follow the procedure
// specified in this subclause, item c, starting again from the beginning, or if
// the message contains a UMTS authentication challenge, the UE shall follow the
// procedure specified in item d. If the SQN is invalid, the UE shall proceed as
// specified in item e.

// It can be assumed that the source of the authentication challenge is not
// genuine (authentication not accepted by the UE) if any of the following
// occur:
//- after sending the AUTHENTICATION FAILURE message with the EMM cause #20 "MAC
// failure" the timer T3418 expires;
//- the UE detects any combination of the authentication failures: EMM causes
//#20 "MAC failure" and #21 "synch failure", during three consecutive
// authentication challenges. The authentication challenges shall be considered
// as consecutive only, if the authentication challenges causing the second and
// third authentication failure are received by the UE, while the timer T3418 or
// T3420 started after the previous authentication failure is running. When it
// has been deemed by the UE that the source of the authentication challenge is
// not genuine (i.e. authentication not accepted by the UE), the UE shall
// proceed as described in item f.

#define R10_5_4_2_7_d__1                                                       \
  "Abnormal cases                                                                                \
    Authentication failure (EMM cause #26 \"non-EPS authentication unacceptable\"):                                     \
    The UE shall send an AUTHENTICATION FAILURE message, with EMM cause #26 \"non-EPS authentication                    \
    unacceptable\", to the network and start the timer T3418 (see example in figure 5.4.2.7.1). Furthermore, the UE     \
    shall stop any of the retransmission timers that are running (e.g. T3410, T3417, T3421 or T3430). Upon the first    \
    receipt of an AUTHENTICATION FAILURE message from the UE with EMM cause #26 \"non-EPS                               \
    authentication unacceptable\", the network may initiate the identification procedure described in subclause 5.4.4.  \
    This is to allow the network to obtain the IMSI from the UE. The network may then check that the GUTI               \
    originally used in the authentication challenge corresponded to the correct IMSI. Upon receipt of the IDENTITY      \
    REQUEST message from the network, the UE shall send the IDENTITY RESPONSE message."
#define R10_5_4_2_7_d__1_BRIEF                                                 \
  "Authentication failure (EMM cause #26 \"non-EPS authentication "            \
  "unacceptable\")"

#define R10_5_4_2_7_d__NOTE2                                                   \
  "Abnormal cases                                                                            \
    Upon receipt of an AUTHENTICATION FAILURE message from the UE with EMM cause #26 \"non-                             \
    EPS authentication unacceptable\", the network may also terminate the authentication procedure (see                 \
    subclause 5.4.2.5)."
#define R10_5_4_2_7_d__NOTE2_BRIEF "Terminate the authentication procedure"

/* Case implemented but have to think how to trace it
#define R10_5_4_2_7_d__2 "Abnormal cases \
If the GUTI/IMSI mapping in the network was incorrect, the network should
respond by sending a new                      \
    AUTHENTICATION REQUEST message to the UE. Upon receiving the new
AUTHENTICATION REQUEST                             \
    message from the network, the UE shall stop the timer T3418, if running, and
then process the challenge             \
    information as normal. If the GUTI/IMSI mapping in the network was correct,
the network terminates the              \ authentication procedure (see
subclause 5.4.2.5)."*/

/*#define R10_5_4_2_7_e__1 "Abnormal cases \
    Authentication failure (EMM cause #21 \"synch failure\"): \
    The UE shall send an AUTHENTICATION FAILURE message, with EMM cause #21
   \"synch failure\", to the                   \
    network and start the timer T3420 (see example in figure 5.4.2.7.2).
   Furthermore, the UE shall stop any of the      \ retransmission timers that
   are running (e.g. T3410, T3417, T3421 or T3430)."*/

#define R10_5_4_2_7_e__2                                                       \
  "Abnormal cases                                                                                \
    Upon the first receipt of an AUTHENTICATION FAILURE message from the UE with the EMM cause #21 \"synch failure\",   \
    the network shall use the returned AUTS parameter from the authentication failure parameter IE in the AUTHENTICATION\
    FAILURE message, to re-synchronise."
#define R10_5_4_2_7_e__2_BRIEF "Re-synchronise with AUTS parameter"

#define R10_5_4_2_7_e__3                                                       \
  "Abnormal cases                                                                                \
    The re-synchronisation procedure requires the MME to delete all unused                                              \
    authentication vectors for that IMSI and obtain new vectors from the HSS. When re-synchronisation is complete,      \
    the network shall initiate the authentication procedure. Upon receipt of the AUTHENTICATION REQUEST                 \
    message, the UE shall stop the timer T3420, if running."
#define R10_5_4_2_7_e__3_BRIEF "Re-synchronisation, new vectors"

#define R10_5_4_2_7_e__NOTE3                                                   \
  "Abnormal cases                                                                            \
    Upon receipt of two consecutive AUTHENTICATION FAILURE messages from the UE with EMM                                \
    cause #21 \"synch failure\", the network may terminate the authentication procedure by sending an                   \
    AUTHENTICATION REJECT message."
#define R10_5_4_2_7_e__NOTE3_BRIEF                                             \
  "Two consecutive AUTHENTICATION FAILURE messages (synch failure): "          \
  "AUTHENTICATION REJECT"

// If the network is validated successfully (a new AUTHENTICATION REQUEST is
// received which contains a valid SQN and MAC) while T3420 is running, the UE
// shall send the AUTHENTICATION RESPONSE message to the network and shall start
// any retransmission timers (e.g. T3410, T3417, T3421 or T3430), if they were
// running and stopped when the UE received the first failed AUTHENTICATION
// REQUEST message.

// If the UE receives the second AUTHENTICATION REQUEST while T3420 is running,
// and the MAC value cannot be resolved, the UE shall follow the procedure
// specified in item c or if the message contains a UMTS authentication
// challenge, the UE shall proceed as specified in item d; if the SQN is
// invalid, the UE shall follow the procedure specified in this subclause, item
// e, starting again from the beginning.
//
// The UE shall deem that the network has failed the authentication check and
// proceed as described in item f if any of the following occurs:
//
//- the timer T3420 expires;
//- the UE detects any combination of the authentication failures: EMM cause #20
//"MAC failure", #21 "synch failure", or #26 "non-EPS authentication
// unacceptable", during three consecutive authentication challenges. The
// authentication challenges shall be considered as consecutive only if the
// authentication challenges causing the second and third authentication failure
// are received by the UE while the timer T3418 or T3420 started after the
// previous authentication failure is running.

// f) Network failing the authentication check:
// If the UE deems that the network has failed the authentication check, then it
// shall request RRC to locally release the RRC connection and treat the active
// cell as barred (see 3GPP TS 36.331 [22]). The UE shall start any
// retransmission timers (e.g. T3410, T3417, T3421 or T3430), if they were
// running and stopped when the UE received the first AUTHENTICATION REQUEST
// message containing an invalid MAC or SQN.

// g) Transmission failure of AUTHENTICATION RESPONSE message or AUTHENTICATION
// FAILURE message indication from lower layers (if the authentication procedure
// is triggered by a tracking area updating procedure) The UE shall re-initiate
// the tracking area updating procedure.

// h) Transmission failure of AUTHENTICATION RESPONSE message or AUTHENTICATION
// FAILURE message indication with TAI change from lower layers (if the
// authentication procedure is triggered by a service request procedure) If the
// current TAI is not in the TAI list, the authentication procedure shall be
// aborted and a tracking area updating procedure shall be initiated. If the
// current TAI is still part of the TAI list, it is up to the UE implementation
// how to re-run the ongoing procedure that triggered the authentication
// procedure.

// i) Transmission failure of AUTHENTICATION RESPONSE message or AUTHENTICATION
// FAILURE message indication without TAI change from lower layers (if the
// authentication procedure is triggered by a service request procedure) It is
// up to the UE implementation how to re-run the ongoing procedure that
// triggered the authentication procedure.

#define R10_5_4_2_7_j                                                          \
  "Abnormal cases                                                                            \
    Lower layers indication of non-delivered NAS PDU due to handover                                             \
    If the AUTHENTICATION REQUEST message could not be delivered due to an intra MME handover and the            \
    target TA is included in the TAI list, then upon successful completion of the intra MME handover the MME     \
    shall retransmit the AUTHENTICATION REQUEST message. If a failure of handover procedure is reported by       \
    the lower layer and the S1 signalling connection exists, the MME shall retransmit the AUTHENTICATION         \
    REQUEST message."
#define R10_5_4_2_7_j_BRIEF                                                    \
  "Lower layers indication of non-delivered NAS PDU due to handover"

// For items c, d, and e:
// Depending on local requirements or operator preference for emergency bearer
// services, if the UE has a PDN connection for emergency bearer services
// established or is establishing a PDN connection for emergency bearer
// services, the MME need not follow the procedures specified for the
// authentication failure specified in the present subclause. The MME may
// respond to the AUTHENTICATION FAILURE message by initiating the security mode
// control procedure selecting the "null integrity protection algorithm" EIA0,
// null ciphering algorithm or may abort the authentication procedure and
// continue using the current security context, if any. The MME shall deactivate
// all non-emergency EPS bearer contexts, if any, by initiating an EPS bearer
// context deactivation procedure. If there is an ongoing PDN connectivity
// procedure, the MME shall deactivate all non-emergency EPS bearer contexts
// upon completion of the PDN connectivity procedure. The network shall consider
// the UE to be attached for emergency bearer services only. If a UE has a PDN
// connection for emergency bearer services established or is establishing a PDN
// connection for emergency bearer services and sends an AUTHENTICATION FAILURE
// message to the MME with the EMM cause appropriate for these cases (#20, #21,
// or #26, respectively) and receives the SECURITY MODE COMMAND message before
// the timeout of timer T3418 or T3420, the UE shall deem that the network has
// passed the authentication check successfully, stop timer T3418 or T3420,
// respectively, and execute the security mode control procedure. If a UE has a
// PDN connection for emergency bearer services established or is establishing a
// PDN connection for emergency bearer services when timer T3418 or T3420
// expires, the UE shall not deem that the network has failed the authentication
// check and not behave as described in item f. Instead the UE shall continue
// using the current security context, if any, deactivate all non-emergency EPS
// bearer contexts, if any, by initiating UE requested PDN disconnect procedure.
// If there is an ongoing PDN connectivity procedure, the UE shall deactivate
// all non-emergency EPS bearer contexts upon completion of the PDN connectivity
// procedure. The UE shall consider itself to be attached for emergency bearer
// services only.

//-----------------------------------------------------------------------------------------------------------------------
// SECURITY MODE CONTROL
//-----------------------------------------------------------------------------------------------------------------------

#define R10_5_4_3_2__1                                                         \
  "NAS security mode control initiation by the network                                             \
    The MME initiates the NAS security mode control procedure by sending a SECURITY MODE COMMAND message                \
    to the UE and starting timer T3460 (see example in figure 5.4.3.2.1)."
#define R10_5_4_3_2__1_BRIEF "SMC initiation, start T3460"

#define R10_5_4_3_2__2                                                         \
  "NAS security mode control initiation by the network                                             \
    The MME shall reset the downlink NAS COUNT counter and use it to integrity protect the initial SECURITY MODE        \
    COMMAND message if the security mode control procedure is initiated:                                                \
    - to take into use the EPS security context created after a successful execution of the EPS authentication          \
      procedure;"
#define R10_5_4_3_2__2_BRIEF "SMC initiation, reset DL NAS count, use SC"

//#define R10_5_4_3_2__3 "NAS security mode control initiation by the network
//    The MME shall reset the downlink NAS COUNT counter and use it to integrity
//    protect the initial SECURITY MODE COMMAND message if the security mode
//    control procedure is initiated:
//    - upon receipt of TRACKING AREA UPDATE REQUEST message including a GPRS
//    ciphering key sequence
//      number IE, if the MME wishes to create a mapped EPS security context
//      (i.e. the type of security context flag is set to \"mapped security
//      context\" in the NAS key set identifier IE included in the SECURITY MODE
//      COMMAND message)."

// Done for KASME (not K'ASME) but have to trace it
//#define R10_5_4_3_2__4 "NAS security mode control initiation by the network
//    The MME shall send the SECURITY MODE COMMAND message unciphered, but shall
//    integrity protect the message with the NAS integrity key based on K ASME
//    or mapped K' ASME indicated by the eKSI included in the message. The MME
//    shall set the security header type of the message to \"integrity protected
//    with new EPS security context\"."

// The MME shall create a locally generated K ASME and send the SECURITY MODE
// COMMAND message including a KSI value in the NAS key set identifier IE set to
// "000" and EIA0 and EEA0 as the selected NAS security algorithms when the
// security mode control procedure is initiated:
//- during an attach procedure for emergency bearer services if no shared EPS
// security context is available;
//- during a tracking area updating procedure for a UE that has a PDN connection
// for emergency bearer services if no shared EPS security context is available;
// or
//- after a failed authentication procedure for a UE that has a PDN connection
// for emergency bearer services if continued usage of a shared security context
// is not possible.

// The UE shall process a SECURITY MODE COMMAND message including a KSI value in
// the NAS key set identifier IE set to "000" and EIA0 and EEA0 as the selected
// NAS security algorithms and, if accepted, create a locally generated K ASME
// when the security mode control procedure is initiated:
//- during an attach procedure for emergency bearer services;
//- during a tracking area updating procedure when the UE has a PDN connection
// for emergency bearer services; or
//- after an authentication procedure when the UE has a PDN connection for
// emergency bearer services. NOTE 1: The process for creation of the locally
// generated K ASME by the MME and the UE is implementation dependent.

// Upon receipt of a TRACKING AREA UPDATE REQUEST message including a GPRS
// ciphering key sequence number IE, if the MME does not have the valid current
// EPS security context indicated by the UE, the MME shall either:
//- indicate the use of the new mapped EPS security context to the UE by setting
// the type of security context flag in the NAS key set identifier IE to "mapped
// security context" and the KSI value related to the security context of the
// source system; or
//- set the KSI value "000" in the NAS key set identifier IE if the MME sets
// EIA0 and EEA0 as the selected NAS security algorithms if the UE has a PDN
// connection for emergency bearer services. While having a current mapped EPS
// security context with the UE, if the MME wants to take the native EPS
// security context into use, the MME shall include the eKSI that indicates the
// native EPS security context in the SECURITY MODE COMMAND message.

#define R10_5_4_3_2__14                                                        \
  "NAS security mode control initiation by the network                                            \
    The MME shall include the replayed security capabilities of the UE (including the security capabilities with regard to\
    NAS, RRC and UP (user plane) ciphering as well as NAS, RRC integrity, and other possible target network security    \
    capabilities, i.e. UTRAN/GERAN if UE included them in the message to network), the replayed nonce UE if the UE      \
    included it in the message to the network, the selected NAS ciphering and integrity algorithms and the Key Set      \
    Identifier (eKSI)."
#define R10_5_4_3_2__14_BRIEF                                                  \
  "SMC initiation, include replayed security capabilities, eKSI, algos"

// The MME shall include both the nonce MME and the nonce UE when creating a
// mapped EPS security context during inter- system change from A/Gb mode to S1
// mode or Iu mode to S1 mode in EMM-IDLE mode.

// The MME may initiate a SECURITY MODE COMMAND in order to change the NAS
// security algorithms for a current EPS security context already in use. The
// MME re-derives the NAS keys from K ASME with the new NAS algorithm identities
// as input and provides the new NAS algorithm identities within the SECURITY
// MODE COMMAND message.

// Additionally, the MME may request the UE to include its IMEISV in the
// SECURITY MODE COMPLETE message. NOTE 2: The AS and NAS security capabilities
// will be the same, i.e. if the UE supports one algorithm for NAS it is also be
// supported for AS.

//------------------------------

#define R10_5_4_3_4__1                                                         \
  "NAS security mode control completion by the network                                             \
    The MME shall, upon receipt of the SECURITY MODE COMPLETE message, stop timer T3460."
#define R10_5_4_3_4__1_BRIEF "SMC completion, stop T3460, "

#define R10_5_4_3_4__2                                                         \
  "NAS security mode control completion by the network                                             \
    From this time onward the MME shall integrity protect and encipher all signalling messages with the selected NAS    \
    integrity and ciphering algorithms."
#define R10_5_4_3_4__2_BRIEF "SMC completion, integ. cipher. all messages"

//------------------------------

// 5.4.3.5
// NAS security mode command not accepted by the UE
// If the security mode command cannot be accepted, the UE shall send a SECURITY
// MODE REJECT message. The SECURITY MODE REJECT message contains an EMM cause
// that typically indicates one of the following cause values: #23: UE security
// capabilities mismatch; #24: security mode rejected, unspecified.

#define R10_5_4_3_5__2                                                         \
  "NAS security mode command not accepted by the UE                                                \
    Upon receipt of the SECURITY MODE REJECT message, the MME shall stop timer T3460. The MME shall also abort          \
    the ongoing procedure that triggered the initiation of the NAS security mode control procedure."
#define R10_5_4_3_5__2_BRIEF                                                   \
  "SECURITY MODE REJECT received, stop T3460, abort procedure"

#define R10_5_4_3_5__3                                                         \
  "NAS security mode command not accepted by the UE                                                \
    Both the UE and the MME shall apply the EPS security context in use before the initiation of the security mode control\
    procedure, if any, to protect the SECURITY MODE REJECT message and any other subsequent messages according to       \
    the rules in subclauses 4.4.4 and 4.4.5."
#define R10_5_4_3_5__3_BRIEF "SECURITY MODE REJECT received, apply previous SC"

//------------------------------

#define R10_5_4_3_7_a                                                          \
  "Abnormal cases on the network side                                                               \
    Lower layer failure before the SECURITY MODE COMPLETE or SECURITY MODE REJECT message is received                   \
    The network shall abort the procedure."
#define R10_5_4_3_7_a_BRIEF                                                    \
  "Lower layer failure before SECURITY MODE COMPLETE or SECURITY MODE REJECT " \
  "is received"

#define R10_5_4_3_7_b__1                                                       \
  "Abnormal cases on the network side                                                            \
    Expiry of timer T3460                                                                                               \
    The network shall, on the first expiry of the timer T3460, retransmit the SECURITY MODE COMMAND                     \
    message and shall reset and start timer T3460. This retransmission is repeated four times, i.e."
#define R10_5_4_3_7_b__1_BRIEF "SMC, Expiry of timer T3460"

#define R10_5_4_3_7_b__2                                                       \
  "Abnormal cases on the network side                                                            \
    Expiry of timer T3460                                                                                               \
    on the fifth expiry of timer T3460, the procedure shall be aborted."
#define R10_5_4_3_7_b__2_BRIEF                                                 \
  "Expiry of timer T3460, procedure shall be aborted"

// NOTE:
// If the SECURITY MODE COMMAND message was sent to create a mapped EPS security
// context during inter-system change from A/Gb mode to S1 mode or Iu mode to S1
// mode, then the network does not generate new values for the nonce MME and the
// nonce UE , but includes the same values in the SECURITY MODE COMMAND message
// (see the subclause 7.2.4.4 in 3GPP TS 33.401 [19]).

#define R10_5_4_3_7_c                                                          \
  "Abnormal cases on the network side                                                               \
    Collision between security mode control procedure and attach, service request, tracking area updating procedure     \
    or detach procedure not indicating switch off                                                                       \
    The network shall abort the security mode control procedure and proceed with the UE initiated procedure."
#define R10_5_4_3_7_c_BRIEF                                                    \
  "Collision between SMC procedure and attach, SR, TAU or detach not "         \
  "indicating switch off"

#define R10_5_4_3_7_d                                                          \
  "Abnormal cases on the network side                                                               \
    Collision between security mode control procedure and other EMM procedures than in item c                           \
    The network shall progress both procedures."
#define R10_5_4_3_7_d_BRIEF                                                    \
  "Collision between security mode control procedure and other EMM procedures"

#define R10_5_4_3_7_e                                                          \
  "Abnormal cases on the network side                                                               \
    Lower layers indication of non-delivered NAS PDU due to handover                                                    \
    If the SECURITY MODE COMMAND message could not be delivered due to an intra MME handover and the                    \
    target TA is included in the TAI list, then upon successful completion of the intra MME handover the MME            \
    shall retransmit the SECURITY MODE COMMAND message. If a failure of the handover procedure is reported              \
    by the lower layer and the S1 signalling connection exists, the MME shall retransmit the SECURITY MODE              \
    COMMAND message."
#define R10_5_4_3_7_e_BRIEF                                                    \
  "Lower layers indication of non-delivered NAS PDU due to handover"

//-----------------------------------------------------------------------------------------------------------------------
// IDENTIFICATION
//-----------------------------------------------------------------------------------------------------------------------

#define R10_5_4_4_1                                                            \
  "General                                                                                            \
    The identification procedure is used by the network to request a particular UE to provide specific identification   \
    parameters, e.g. the International Mobile Subscriber Identity (IMSI) or the International Mobile Equipment Identity \
    (IMEI). IMEI and IMSI definition and structure are specified in 3GPP TS 23.003 [2].                                 \
    For mobile device supporting both 3GPP access and cdma2000 ® access a single IMEI is used to identify the device as \
    specified in 3GPP TS 22.278 [1C]."
#define R10_5_4_4_1_BRIEF "Identification procedure"

#define R10_5_4_4_2                                                            \
  "Identification initiation by the network                                                           \
    The network initiates the identification procedure by sending an IDENTITY REQUEST message to the UE and starting    \
    the timer T3470 (see example in figure 5.4.4.2.1). The IDENTITY REQUEST message specifies the requested             \
    identification parameters in the Identity type information element."
#define R10_5_4_4_2_BRIEF "Identification initiation by the network"

#define R10_5_4_4_4                                                            \
  "Identification completion by the network                                                           \
    Upon receipt of the IDENTITY RESPONSE the network shall stop the timer T3470."
#define R10_5_4_4_4_BRIEF "Identification completion by the network"

#define R10_5_4_4_6_a                                                          \
  "Abnormal cases on the network side                                                               \
    Lower layer failure                                                                                                 \
    Upon detection of a lower layer failure before the IDENTITY RESPONSE is received, the network shall abort           \
    any ongoing EMM procedure."
#define R10_5_4_4_6_a_BRIEF                                                    \
  "Lower layer failure detected, abort any EMM procedure"

#define R10_5_4_4_6_b__1                                                       \
  "Abnormal cases on the network side                                                            \
    Expiry of timer T3470                                                                                               \
    The identification procedure is supervised by the network by the timer T3470. The network shall, on the first       \
    expiry of the timer T3470, retransmit the IDENTITY REQUEST message and reset and restart the timer T3470."
#define R10_5_4_4_6_b__1_BRIEF "Expiry of timer T3470"

#define R10_5_4_4_6_b__2                                                       \
  "Abnormal cases on the network side                                                            \
    Expiry of timer T3470                                                                                               \
    This retransmission is repeated four times, i.e. on the fifth expiry of timer T3470, the network shall abort the    \
    identification procedure and any ongoing EMM procedure."
#define R10_5_4_4_6_b__2_BRIEF "Expiry of timer T3470"

#define R10_5_4_4_6_c                                                          \
  "Abnormal cases on the network side                                                               \
    Collision of an identification procedure with an attach procedure                                                   \
    If the network receives an ATTACH REQUEST message before the ongoing identification procedure has been              \
    completed and no attach procedure is pending on the network (i.e. no ATTACH ACCEPT/REJECT message has               \
    still to be sent as an answer to an ATTACH REQUEST message), the network shall proceed with the attach              \
    procedure."
#define R10_5_4_4_6_c_BRIEF                                                    \
  "Collision of an identification procedure with an attach procedure"
#define R10_5_4_4_6_d                                                          \
  "Abnormal cases on the network side                                                               \
    Collision of an identification procedure with an attach procedure when the identification procedure has been        \
    caused by an attach procedure"
#define R10_5_4_4_6_d_BRIEF                                                    \
  "Collision of an identification procedure caused by attach procedure with "  \
  "an attach procedure"

#define R10_5_4_4_6_d__1                                                       \
  "Abnormal cases on the network side                                                            \
    Collision of an identification procedure with an attach procedure when the identification procedure has been        \
    caused by an attach procedure                                                                                       \
   If the network receives an ATTACH REQUEST message before the ongoing identification procedure has been               \
   completed and an attach procedure is pending (i.e. an ATTACH ACCEPT/REJECT message has to be sent as an              \
   answer to an earlier ATTACH REQUEST message), then:                                                                  \
   - If one or more of the information elements in the ATTACH REQUEST message differ from the ones                      \
     received within the previous ATTACH REQUEST message, the network shall proceed with the new attach                 \
     procedure;"
#define R10_5_4_4_6_d__1_BRIEF                                                 \
  "Collision of an identification procedure with an attach procedure, attach " \
  "IEs changed, new attach go on"
#define R10_5_4_4_6_d__2                                                       \
  "Abnormal cases on the network side                                                            \
    Collision of an identification procedure with an attach procedure when the identification procedure has been        \
    caused by an attach procedure                                                                                       \
   If the network receives an ATTACH REQUEST message before the ongoing identification procedure has been               \
   completed and an attach procedure is pending (i.e. an ATTACH ACCEPT/REJECT message has to be sent as an              \
   answer to an earlier ATTACH REQUEST message), then:                                                                  \
   - If the information elements do not differ, then the network shall not treat any further this new ATTACH            \
     REQUEST."
#define R10_5_4_4_6_d__2_BRIEF                                                 \
  "Collision of an identification procedure with an attach procedure, attach " \
  "IEs unchanged, new attach discarded"

//-----------------------------------------------------------------------------------------------------------------------
// ATTACH
//-----------------------------------------------------------------------------------------------------------------------

#define R10_5_5_1__1                                                           \
  "Attach procedure - General                                                                        \
    The attach procedure is used to attach to an EPC for packet services in EPS.                                        \
    ...                                                                                                                 \
    If the MME does not support an attach for emergency bearer services, the MME shall reject any request to attach with\
    an attach type set to \"EPS emergency attach\"."
#define R10_5_5_1__1_BRIEF ""

//#define R10_5_5_1__2 "Attach procedure - General
//    With a successful attach procedure, a context is established for the UE in
//    the MME, and a default bearer is established between the UE and the PDN
//    GW, thus enabling always-on IP connectivity to the UE. The network may
//    also initiate the activation of dedicated bearers as part of the attach
//    procedure."

//...

//------------------------------

#define R10_5_5_1_2_3__1                                                       \
  "Attach procedure - EMM common procedure initiation                                            \
    The network may initiate EMM common procedures, e.g. the identification, authentication and security mode control   \
    procedures during the attach procedure, depending on the information received in the ATTACH REQUEST message         \
    (e.g. IMSI, GUTI and KSI)."
#define R10_5_5_1_2_3__1_BRIEF                                                 \
  "EMM common procedure initiation during attach procedure"

//#define R10_5_5_1_2_3__2 "Attach procedure - EMM common procedure initiation
//    If the network receives an ATTACH REQUEST message containing the Old GUTI
//    type IE and the EPS mobile identity IE with type of identity indicating
//    \"GUTI\", and the network does not follow the use of the most significant
//    bit of the <MME group id> as specified in 3GPP TS 23.003 [2],
//    subclause 2.8.2.2.2, the network shall use the Old GUTI type IE to
//    determine whether the mobile identity included in the EPS mobile identity
//    IE is a native GUTI or a mapped GUTI."

//#define R10_5_5_1_2_3__3 "Attach procedure - EMM common procedure initiation
//    During an attach for emergency bearer services, the MME may choose to skip
//    the authentication procedure even if no EPS security context is available
//    and proceed directly to the execution of the security mode control
//    procedure as specified in subclause 5.4.3."

//------------------------------

#define R10_5_5_1_2_4__1                                                       \
  "Attach accepted by the network                                                                \
    During an attach for emergency bearer services, if not restricted by local regulations, the MME shall not check for \
    mobility and access restrictions, regional restrictions, subscription restrictions, or perform CSG access control when\
    processing the ATTACH REQUEST message. The network shall not apply subscribed APN based congestion control          \
    during an attach procedure for emergency bearer services."
#define R10_5_5_1_2_4__1_BRIEF ""

#define R10_5_5_1_2_4__2                                                       \
  "Attach accepted by the network                                                                \
    If the attach request is accepted by the network, the MME shall send an ATTACH ACCEPT message to the UE and         \
    start timer T3450. The MME shall send the ATTACH ACCEPT message together with an ACTIVATE DEFAULT EPS               \
    BEARER CONTEXT REQUEST message contained in the ESM message container information element to activate the           \
    default bearer (see subclause 6.4.1). The network may also initiate the activation of dedicated bearers towards the UE\
    by invoking the dedicated EPS bearer context activation procedure (see subclause 6.4.2)."
#define R10_5_5_1_2_4__2_BRIEF ""

#define R10_5_5_1_2_4__3                                                       \
  "Attach accepted by the network                                                                \
    If the attach request is accepted by the network, the MME shall delete the stored UE radio capability information, if\
    any."
#define R10_5_5_1_2_4__3_BRIEF                                                 \
  "Attach accepted by the network, delete the stored UE radio capability "     \
  "information."

#define R10_5_5_1_2_4__4                                                       \
  "Attach accepted by the network                                                                \
    If the UE has included the UE network capability IE or the MS network capability IE or both in the ATTACH           \
    REQUEST message, the MME shall store all octets received from the UE, up to the maximum length defined for the      \
    respective information element.                                                                                     \
    NOTE:                                                                                                               \
      This information is forwarded to the new MME during inter-MME handover or to the new SGSN during                  \
      inter-system handover to A/Gb mode or Iu mode."
#define R10_5_5_1_2_4__4_BRIEF                                                 \
  "Attach accepted by the network, store UE network capability IE or the MS "  \
  "network capability IE or both."

#define R10_5_5_1_2_4__5                                                       \
  "Attach accepted by the network                                                                \
    If the UE specific DRX parameter was included in the DRX Parameter IE in the ATTACH REQUEST message, the            \
    MME shall replace any stored UE specific DRX parameter with the received parameter and use it for the downlink      \
    transfer of signalling and user data."
#define R10_5_5_1_2_4__5_BRIEF                                                 \
  "Attach accepted by the network, use DRX parameter for the downlink "        \
  "transfer of signalling and user data"

#define R10_5_5_1_2_4__6                                                       \
  "Attach accepted by the network                                                                \
    The MME shall assign and include the TAI list the UE is registered to in the ATTACH ACCEPT message. The UE,         \
    upon receiving an ATTACH ACCEPT message, shall delete its old TAI list and store the received TAI list."
#define R10_5_5_1_2_4__6_BRIEF ""

// If the ATTACH ACCEPT message contains a T3412 extended value IE, then the UE
// shall use the value in T3412 extended value IE as periodic tracking area
// update timer (T3412). If the ATTACH ACCEPT message does not contain T3412
// extended value IE, then the UE shall use the value in T3412 value IE as
// periodic tracking area update timer (T3412).

// Upon receiving the ATTACH ACCEPT message, the UE shall stop timer T3410.

#define R10_5_5_1_2_4__9                                                       \
  "Attach accepted by the network                                                                \
    The GUTI reallocation may be part of the attach procedure. When the ATTACH REQUEST message includes the IMSI        \
    or IMEI, or the MME considers the GUTI provided by the UE is invalid, or the GUTI provided by the UE was assigned   \
    by another MME, the MME shall allocate a new GUTI to the UE. The MME shall include in the ATTACH ACCEPT             \
    message the new assigned GUTI together with the assigned TAI list. In this case the MME shall enter state EMM-      \
    COMMON-PROCEDURE-INITIATED as described in subclause 5.4.1."
#define R10_5_5_1_2_4__9_BRIEF ""

#define R10_5_5_1_2_4__10                                                      \
  "Attach accepted by the network                                                               \
    For a shared network, the TAIs included in the TAI list can contain different PLMN identities. The MME indicates the\
    selected core network operator PLMN identity to the UE in the GUTI (see 3GPP TS 23.251 [8B])."
#define R10_5_5_1_2_4__10_BRIEF ""

// If the ATTACH ACCEPT message contains a GUTI, the UE shall use this GUTI as
// the new temporary identity. The UE shall delete its old GUTI and store the
// new assigned GUTI. If no GUTI has been included by the MME in the ATTACH
// ACCEPT message, the old GUTI, if any available, shall be kept.

// If A/Gb mode or Iu mode is supported in the UE, the UE shall set its TIN to
// "GUTI" when receiving the ATTACH ACCEPT message.

#define R10_5_5_1_2_4__13                                                      \
  "Attach accepted by the network                                                               \
    The MME may also include a list of equivalent PLMNs in the ATTACH ACCEPT message. Each entry in the list            \
    contains a PLMN code (MCC+MNC). The UE shall store the list as provided by the network, and if the attach           \
    procedure is not for emergency bearer services, the UE shall remove from the list any PLMN code that is already in the\
    list of forbidden PLMNs. In addition, the UE shall add to the stored list the PLMN code of the registered PLMN that \
    sent the list. The UE shall replace the stored list on each receipt of the ATTACH ACCEPT message. If the ATTACH     \
    ACCEPT message does not contain a list, then the UE shall delete the stored list."
#define R10_5_5_1_2_4__13_BRIEF ""

#define R10_5_5_1_2_4__14                                                      \
  "Attach accepted by the network                                                               \
    The network informs the UE about the support of specific features, such as IMS voice over PS session, location services\
    (EPC-LCS, CS-LCS) or emergency bearer services, in the EPS network feature support information element. In a UE     \
    with IMS voice over PS capability, the IMS voice over PS session indicator and the emergency bearer services indicator\
    shall be provided to the upper layers. The upper layers take the IMS voice over PS session indicator into account as\
    specified in 3GPP TS 23.221 [8A], subclause 7.2a and subclause 7.2b, when selecting the access domain for voice     \
    sessions or calls. When initiating an emergency call, the upper layers also take the emergency bearer services indicator\
    into account for the access domain selection. In a UE with LCS capability, location services indicators (EPC-LCS, CS-\
    LCS) shall be provided to the upper layers. When MO-LR procedure is triggered by the UE's application, those        \
    indicators are taken into account as specified in 3GPP TS 24.171 [13C]."
#define R10_5_5_1_2_4__14_BRIEF ""

// If the UE has initiated the attach procedure due to manual CSG selection and
// receives an ATTACH ACCEPT message; and the UE sent the ATTACH REQUEST message
// in a CSG cell, the UE shall check if the CSG ID and associated PLMN identity
// of the cell are contained in the Allowed CSG list. If not, the UE shall add
// that CSG ID and associated PLMN identity to the Allowed CSG list and the UE
// may add the HNB Name (if provided by lower layers) to the Allowed CSG list if
// the HNB Name is present in neither the Operator CSG list nor the Allowed CSG
// list.

// When the UE receives the ATTACH ACCEPT message combined with the ACTIVATE
// DEFAULT EPS BEARER CONTEXT REQUEST message, it shall forward the ACTIVATE
// DEFAULT EPS BEARER CONTEXT REQUEST message to the ESM sublayer. Upon receipt
// of an indication from the ESM sublayer that the default EPS bearer context
// has been activated, the UE shall send an ATTACH COMPLETE message together
// with an ACTIVATE DEFAULT EPS BEARER CONTEXT ACCEPT message contained in the
// ESM message container information element to the network.

// Additionally, the UE shall reset the attach attempt counter and tracking area
// updating attempt counter, enter state EMM- REGISTERED and set the EPS update
// status to EU1 UPDATED.

// When the UE receives any ACTIVATE DEDICATED EPS BEARER CONTEXT REQUEST
// messages during the attach procedure, the UE shall forward the ACTIVATE
// DEDICATED EPS BEARER CONTEXT REQUEST message(s) to the ESM sublayer. The UE
// shall send a response to the ACTIVATE DEDICATED EPS BEARER CONTEXT REQUEST
// message(s) after successful completion of the attach procedure.

// If the attach procedure was initiated in S101 mode, the lower layers are
// informed about the successful completion of the procedure.

#define R10_5_5_1_2_4__20                                                      \
  "Attach accepted by the network                                                               \
    Upon receiving an ATTACH COMPLETE message, the MME shall stop timer T3450, enter state EMM-REGISTERED               \
    and consider the GUTI sent in the ATTACH ACCEPT message as valid."
#define R10_5_5_1_2_4__20_BRIEF                                                \
  "Attach accepted by the network, ATTACH COMPLETE received, enter state "     \
  "EMM-REGISTERED"

//------------------------------
#define R10_5_5_1_2_7_a                                                        \
  "Abnormal cases on the network side                                                             \
    a) Lower layer failure                                                                                              \
    If a lower layer failure occurs before the message ATTACH COMPLETE has been received from the UE, the               \
    network shall locally abort the attach procedure, enter state EMM-DEREGISTERED and shall not resend the             \
    message ATTACH ACCEPT. If a new GUTI was assigned to the UE in the attach procedure, the MME shall                  \
    consider both the old and the new GUTI as valid until the old GUTI can be considered as invalid by the network      \
    or the EMM context which has been marked as detached in the network is released.                                    \
    If the old GUTI was allocated by an MME other than the current MME, the current MME does not need to retain         \
    the old GUTI. If the old GUTI is used by the UE in a subsequent attach message, the network may use the             \
    identification procedure to request the UE's IMSI.       "
#define R10_5_5_1_2_7_a_BRIEF                                                  \
  "Abnormal cases on the network side: Lower layer failure"

#define R10_5_5_1_2_7_b__1                                                     \
  "Abnormal cases on the network side                                                          \
    b) Protocol error                                                                                                   \
    If the ATTACH REQUEST message is received with a protocol error, the network shall return an ATTACH                 \
    REJECT message with the following EMM cause value: invalid mandatory information       "
#define R10_5_5_1_2_7_b__1_BRIEF                                               \
  "Abnormal cases on the network side: Protocol error invalid mandatory "      \
  "information"
#define R10_5_5_1_2_7_b__2                                                     \
  "Abnormal cases on the network side                                                          \
    b) Protocol error                                                                                                   \
    If the ATTACH REQUEST message is received with a protocol error, the network shall return an ATTACH                 \
    REJECT message with the following EMM cause value: information element non-existent or not implemented"
#define R10_5_5_1_2_7_b__2_BRIEF                                               \
  "Abnormal cases on the network side: Protocol error information element "    \
  "non-existent or not implemented"
#define R10_5_5_1_2_7_b__3                                                     \
  "Abnormal cases on the network side                                                          \
    b) Protocol error                                                                                                   \
    If the ATTACH REQUEST message is received with a protocol error, the network shall return an ATTACH                 \
    REJECT message with the following EMM cause value: conditional IE error"
#define R10_5_5_1_2_7_b__3_BRIEF                                               \
  "Abnormal cases on the network side: Protocol error conditional IE error"
#define R10_5_5_1_2_7_b__4                                                     \
  "Abnormal cases on the network side                                                          \
    b) Protocol error                                                                                                   \
    If the ATTACH REQUEST message is received with a protocol error, the network shall return an ATTACH                 \
    REJECT message with the following EMM cause value: unspecified"
#define R10_5_5_1_2_7_b__4_BRIEF                                               \
  "Abnormal cases on the network side: protocol error, unspecified"

#define R10_5_5_1_2_7_c__1                                                     \
  "Abnormal cases on the network side                                                          \
    c) T3450 time-out                                                                                                   \
    On the first expiry of the timer, the network shall retransmit the ATTACH ACCEPT message and shall reset and        \
    restart timer T3450.                                                                                                \
    This retransmission is repeated four times"
#define R10_5_5_1_2_7_c__1_BRIEF                                               \
  "Abnormal cases on the network side: T3450 time-out"

#define R10_5_5_1_2_7_c__2                                                     \
  "Abnormal cases on the network side                                                          \
    c) T3450 time-out                                                                                                   \
    i.e. on the fifth expiry of timer T3450, the attach procedure shall be                                              \
    aborted and the MME enters state EMM-DEREGISTERED. "
#define R10_5_5_1_2_7_c__2_BRIEF                                               \
  "Abnormal cases on the network side: fifth expiry of timer T3450"

#define R10_5_5_1_2_7_c__3                                                     \
  "Abnormal cases on the network side                                                          \
    c) T3450 time-out                                                                                                   \
    If a new GUTI was allocated in the ATTACH                                                                           \
    ACCEPT message, the network shall consider both the old and the new GUTI as valid until the old GUTI can be         \
    considered as invalid by the network. If the old GUTI was allocated by an MME other than the current MME,           \
    the current MME does not need to retain the old GUTI. or the EMM context which has been marked as detached          \
    in the network is released.                                                                                         \
    If the old GUTI is used by the UE in a subsequent attach message, the network acts as specified for case a above."
#define R10_5_5_1_2_7_c__3_BRIEF ""

#define R10_5_5_1_2_7_d__1                                                     \
  "Abnormal cases on the network side                                                          \
    ATTACH REQUEST received after the ATTACH ACCEPT message has been sent and before the ATTACH                         \
    COMPLETE message is received                                                                                        \
    - If one or more of the information elements in the ATTACH REQUEST message differ from the ones                     \
    received within the previous ATTACH REQUEST message, the previously initiated attach procedure shall                \
    be aborted if the ATTACH COMPLETE message has not been received and the new attach procedure shall                  \
    be progressed; "
#define R10_5_5_1_2_7_d__1_BRIEF                                               \
  "ATTACH REQUEST with changed IEs received after ATTACH ACCEPT sent and "     \
  "before ATTACH COMPLETE received"

#define R10_5_5_1_2_7_d__2                                                     \
  "Abnormal cases on the network side                                                          \
    ATTACH REQUEST received after the ATTACH ACCEPT message has been sent and before the ATTACH                         \
    COMPLETE message is received                                                                                        \
    - if the information elements do not differ, then the ATTACH ACCEPT message shall be resent and the timer           \
    T3450 shall be restarted if an ATTACH COMPLETE message is expected. In that case, the retransmission                \
    counter related to T3450 is not incremented."
#define R10_5_5_1_2_7_d__2_BRIEF                                               \
  "ATTACH REQUEST with same IEs received, ATTACH ACCEPT message shall be "     \
  "resent"
#define R10_5_5_1_2_7_d__2_a_BRIEF                                             \
  "ATTACH REQUEST with same IEs received, T3450 shall be restarted if an "     \
  "ATTACH COMPLETE message is expected"

#define R10_5_5_1_2_7_e__1                                                     \
  "Abnormal cases on the network side                                                          \
    e) More than one ATTACH REQUEST received and no ATTACH ACCEPT or ATTACH REJECT message has                          \
    been sent                                                                                                           \
    - If one or more of the information elements in the ATTACH REQUEST message differs from the ones                    \
    received within the previous ATTACH REQUEST message, the previously initiated attach procedure shall                \
    be aborted and the new attach procedure shall be executed"
#define R10_5_5_1_2_7_e__1_BRIEF ""

#define R10_5_5_1_2_7_e__2                                                     \
  "Abnormal cases on the network side                                                          \
    e) More than one ATTACH REQUEST received and no ATTACH ACCEPT or ATTACH REJECT message has                          \
    been sent                                                                                                           \
    - if the information elements do not differ, then the network shall continue with the previous attach procedure     \
    and shall ignore the second ATTACH REQUEST message."
#define R10_5_5_1_2_7_e__2_BRIEF                                               \
  "Ignore the this duplicate ATTACH REQUEST message, continue with the "       \
  "previous attach procedure"

#define R10_5_5_1_2_7_f                                                        \
  "Abnormal cases on the network side                                                             \
    f) ATTACH REQUEST received in state EMM-REGISTERED                                                                  \
    If an ATTACH REQUEST message is received in state EMM-REGISTERED the network may initiate the                       \
    EMM common procedures; if it turned out that the ATTACH REQUEST message was sent by a UE that has                   \
    already been attached, the EMM context, EPS bearer contexts, if any, are deleted and the new ATTACH                 \
    REQUEST is progressed."
#define R10_5_5_1_2_7_f_BRIEF "ATTACH REQUEST received in state EMM-REGISTERED"

#define R10_5_5_1_2_7_g                                                        \
  "Abnormal cases on the network side                                                             \
    g) TRACKING AREA UPDATE REQUEST message received before ATTACH COMPLETE message.                                    \
    Timer T3450 shall be stopped. The allocated GUTI in the attach procedure shall be considered as valid and the       \
    tracking area updating procedure shall be rejected with the EMM cause #10 \"implicitly detached\" as described in   \
    subclause 5.5.3.2.5."
#define R10_5_5_1_2_7_g_BRIEF                                                  \
  "TRACKING AREA UPDATE REQUEST message received before ATTACH COMPLETE "      \
  "message"

//-----------------------------------------------------------------------------------------------------------------------
// TAU
//-----------------------------------------------------------------------------------------------------------------------

//#define R10_5_5_3_2_3__1 "EMM common procedure initiation
//    If the network receives a TRACKING AREA UPDATE REQUEST message containing
//    the Old GUTI type IE, and the network does not follow the use of the most
//    significant bit of the <MME group id> to distinguish the node type as
//    specified in 3GPP TS 23.003 [2], subclause 2.8.2.2.2, the network shall
//    use the Old GUTI type IE to determine whether the mobile identity included
//    in the Old GUTI IE is a native GUTI or a mapped GUTI."

#define R10_5_5_3_2_3__2                                                       \
  "EMM common procedure initiation                                                               \
    During the tracking area updating procedure, the MME may initiate EMM common procedures, e.g. the EMM               \
    authentication and security mode control procedures."
#define R10_5_5_3_2_3__2_BRIEF ""

//#define R10_5_5_3_2_3__3 "EMM common procedure initiation
//    The MME may be configured to skip the authentication procedure even if no
//    EPS security context is available and proceed directly to the execution of
//    the security mode control procedure as specified in subclause 5.4.3,
//    during a tracking area updating procedure for a UE that has only a PDN
//    connection for emergency bearer services.""

//#define R10_5_5_3_2_3__4 "EMM common procedure initiation
//    The MME shall not initiate an EMM authentication procedure before
//    completion of the tracking area updating procedure, if the following
//    conditions apply: a) the UE initiated the tracking area updating procedure
//    after handover or inter-system handover to S1 mode; b) the target cell is
//    a shared network cell; and
//        -the UE has provided its GUTI in the Old GUTI IE or the Additional
//        GUTI IE in the TRACKING AREA
//         UPDATE REQUEST message, and the PLMN identity included in the GUTI is
//         different from the selected PLMN identity of the target cell; or
//        -the UE has mapped the P-TMSI and RAI into the Old GUTI IE and not
//        included an Additional GUTI IE in
//         the TRACKING AREA UPDATE REQUEST message, and the PLMN identity
//         included in the RAI is different from the selected PLMN identity of
//         the target cell."

//------------------------------

#define R10_5_5_3_2_4__1a                                                      \
  "Normal and periodic tracking area updating procedure accepted by the network                 \
    If the tracking area update request has been accepted by the network, the MME shall send a TRACKING AREA            \
    UPDATE ACCEPT message to the UE."
#define R10_5_5_3_2_4__1a_BRIEF ""

#define R10_5_5_3_2_4__1b                                                      \
  "Normal and periodic tracking area updating procedure accepted by the network                 \
    If the MME assigns a new GUTI for the UE, a GUTI shall be included in the                                           \
    TRACKING AREA UPDATE ACCEPT message. In this case, the MME shall start timer T3450 and enter state EMM-             \
    COMMON-PROCEDURE-INITIATED as described in subclause 5.4.1."
#define R10_5_5_3_2_4__1b_BRIEF ""

#define R10_5_5_3_2_4__1c                                                      \
  "Normal and periodic tracking area updating procedure accepted by the network                 \
    The MME may include a new TAI list for the UE in the TRACKING AREA UPDATE ACCEPT message."
#define R10_5_5_3_2_4__1c_BRIEF ""

#define R10_5_5_3_2_4__2                                                       \
  "Normal and periodic tracking area updating procedure accepted by the network                  \
    If the UE has included the UE network capability IE or the MS network capability IE or both in the TRACKING AREA    \
    UPDATE REQUEST message, the MME shall store all octets received from the UE, up to the maximum length defined       \
    for the respective information element."
#define R10_5_5_3_2_4__2_BRIEF ""

#define R10_5_5_3_2_4__NOTE1                                                   \
  "Normal and periodic tracking area updating procedure accepted by the network              \
    NOTE 1: This information is forwarded to the new MME during inter-MME handover or to the new SGSN during            \
    inter-system handover to A/Gb mode or Iu mode."
#define R10_5_5_3_2_4__NOTE1_BRIEF ""

#define R10_5_5_3_2_4__3                                                       \
  "Normal and periodic tracking area updating procedure accepted by the network                  \
    If a UE radio capability information update needed IE is included in the TRACKING AREA UPDATE REQUEST               \
    message, the MME shall delete the stored UE radio capability information, if any."
#define R10_5_5_3_2_4__3_BRIEF ""

#define R10_5_5_3_2_4__4                                                       \
  "Normal and periodic tracking area updating procedure accepted by the network                  \
    If the UE specific DRX parameter was included in the DRX Parameter IE in the TRACKING AREA UPDATE                   \
    REQUEST message, the network shall replace any stored UE specific DRX parameter with the received parameter and     \
    use it for the downlink transfer of signalling and user data."
#define R10_5_5_3_2_4__4_BRIEF ""

//#define R10_5_5_3_2_4__5 "Normal and periodic tracking area updating procedure
// accepted by the network
//    If an EPS bearer context status IE is included in the TRACKING AREA UPDATE
//    REQUEST message, the MME shall deactivate all those EPS bearer contexts
//    locally (without peer-to-peer signalling between the MME and the UE) which
//    are active on the network side, but are indicated by the UE as being
//    inactive. If a default EPS bearer context is marked as inactive in the EPS
//    bearer context status IE included in the TRACKING AREA UPDATE REQUEST
//    message, and this default bearer is not associated with the last PDN of
//    the user in the MME, the MME shall locally deactivate all EPS bearer
//    contexts associated to the PDN connection with the default EPS bearer
//    context without peer-to-peer ESM signalling to the UE."

//#define R10_5_5_3_2_4__6 "Normal and periodic tracking area updating procedure
// accepted by the network
//    If the EPS bearer context status IE is included in the TRACKING AREA
//    UPDATE REQUEST, the MME shall include an EPS bearer context status IE in
//    the TRACKING AREA UPDATE ACCEPT message, indicating which EPS bearer
//    contexts are active in the MME."

//#define R10_5_5_3_2_4__7 "Normal and periodic tracking area updating procedure
// accepted by the network
//    If the EPS update type IE included in the TRACKING AREA UPDATE REQUEST
//    message indicates \"periodic updating\", and the UE was previously
//    successfully attached for EPS and non-EPS services, subject to operator
//    policies the MME should allocate a TAI list that does not span more than
//    one location area."

//#define R10_5_5_3_2_4__8 "Normal and periodic tracking area updating procedure
// accepted by the network
//    If the TRACKING AREA UPDATE ACCEPT message contains T3412 extended value
//    IE, then the UE shall use the T3412 extended value IE as periodic tracking
//    area update timer (T3412). If the TRACKING AREA UPDATE ACCEPT contains
//    T3412 value IE, but not T3412 extended value IE, then the UE shall use
//    value in T3412 value IE as periodic tracking area update timer (T3412). If
//    neither T3412 value IE nor T3412 extended value IE is included, the UE
//    shall use the value currently stored, e.g. from a prior ATTACH ACCEPT or
//    TRACKING AREA UPDATE ACCEPT message."

//#define R10_5_5_3_2_4__9 "Normal and periodic tracking area updating procedure
// accepted by the network
//    Also during the tracking area updating procedure without "active" flag, if
//    the MME has deactivated EPS bearer context(s) locally for any reason, the
//    MME shall inform the UE of the deactivated EPS bearer context(s) by
//    including the EPS bearer context status IE in the TRACKING AREA UPDATE
//    ACCEPT message."

//#define R10_5_5_3_2_4__10 "Normal and periodic tracking area updating
// procedure accepted by the network
//    If due to regional subscription restrictions or access restrictions the UE
//    is not allowed to access the TA, but it has a PDN connection for emergency
//    bearer services established, the MME may accept the TRACKING AREA UPDATE
//    REQUEST message and deactivate all non-emergency EPS bearer contexts by
//    initiating an EPS bearer context deactivation procedure when the TAU is
//    initiated in EMM-CONNECTED mode. When the TAU is initiated in EMM- IDLE
//    mode, the MME locally deactivates all non-emergency EPS bearer contexts
//    and informs the UE via the EPS bearer context status IE in the TRACKING
//    AREA UPDATE ACCEPT message. The MME shall not deactivate the emergency EPS
//    bearer contexts. The network shall consider the UE to be attached for
//    emergency bearer services only and shall indicate in the EPS update result
//    IE in the TRACKING AREA UPDATE ACCEPT message that ISR is not activated."

//#define R10_5_5_3_2_4__11 "Normal and periodic tracking area updating
// procedure accepted by the network
//    If a TRACKING AREA UPDATE REQUEST message is received from a UE with a
//    LIPA PDN connection, and if:
//    - a GW Transport Layer Address IE value identifying a L-GW is provided by
//    the lower layer together with the TRACKING AREA UPDATE REQUEST message,
//    and the P-GW address included in the EPS bearer context of the LIPA PDN
//    Connection is different from the provided GW Transport Layer Address IE
//    value (see 3GPP TS 36.413 [36]); or
//    - no GW Transport Layer Address is provided together with the tracking
//    area update request by the lower layer, then the MME locally deactivates
//    all EPS bearer contexts associated with the LIPA PDN connection. If active
//    EPS bearer contexts remain for the UE and the TRACKING AREA UPDATE REQUEST
//    request message is accepted, the MME informs the UE via the EPS bearer
//    context status IE in the TRACKING AREA UPDATE ACCEPT message that EPS
//    bearer contexts were locally deactivated."

//#define R10_5_5_3_2_4__12 "Normal and periodic tracking area updating
// procedure accepted by the network
//    For a shared network, the TAIs included in the TAI list can contain
//    different PLMN identities. The MME indicates the selected core network
//    operator PLMN identity to the UE in the GUTI (see 3GPP TS 23.251 [8B])."

//#define R10_5_5_3_2_4__13 "Normal and periodic tracking area updating
// procedure accepted by the network
//    If the "active" flag is included in the TRACKING AREA UPDATE REQUEST
//    message, the MME shall re-establish the radio and S1 bearers for all
//    active EPS bearer contexts."

//#define R10_5_5_3_2_4__14 "Normal and periodic tracking area updating
// procedure accepted by the network
//    If the "active" flag is not included in the TRACKING AREA UPDATE REQUEST
//    message, the MME may also re- establish the radio and S1 bearers for all
//    active EPS bearer contexts due to downlink pending data or downlink
//    pending signalling."

//#define R10_5_5_3_2_4__15 "Normal and periodic tracking area updating
// procedure accepted by the network
//    Upon receiving a TRACKING AREA UPDATE ACCEPT message, the UE shall stop
//    timer T3430, reset the tracking area updating attempt counter, enter state
//    EMM-REGISTERED and set the EPS update status to EU1 UPDATED. If the
//    message contains a GUTI, the UE shall use this GUTI as new temporary
//    identity for EPS services and shall store the new GUTI. If no GUTI was
//    included by the MME in the TRACKING AREA UPDATE ACCEPT message, the old
//    GUTI shall be used. If the UE receives a new TAI list in the TRACKING AREA
//    UPDATE ACCEPT message, the UE shall consider the new TAI list as valid and
//    the old TAI list as invalid; otherwise, the UE shall consider the old TAI
//    list as valid."

//#define R10_5_5_3_2_4__16 "Normal and periodic tracking area updating
// procedure accepted by the network
//    If the UE had initiated the tracking area updating procedure in EMM-IDLE
//    mode to perform an inter-system change from A/Gb mode or Iu mode to S1
//    mode and the nonce UE was included in the TRACKING AREA UPDATE REQUEST
//    message, the UE shall delete the nonce UE upon receipt of the TRACKING
//    AREA UPDATE ACCEPT message."

//#define R10_5_5_3_2_4__17 "Normal and periodic tracking area updating
// procedure accepted by the network
//    If an EPS bearer context status IE is included in the TRACKING AREA UPDATE
//    ACCEPT message, the UE shall deactivate all those EPS bearers contexts
//    locally (without peer-to-peer signalling between the UE and the MME) which
//    are active in the UE, but are indicated by the MME as being inactive. If a
//    default EPS bearer context is marked as inactive in the EPS bearer context
//    status IE included in the TRACKING AREA UPDATE ACCEPT message, and this
//    default bearer is not associated with the last PDN in the UE, the UE shall
//    locally deactivate all EPS bearer contexts associated to the PDN
//    connection with the default EPS bearer context without peer-to-peer ESM
//    signalling to the MME. If only the PDN connection for emergency bearer
//    services remains established, the UE shall consider itself attached for
//    emergency bearer services only."

//#define R10_5_5_3_2_4__18 "Normal and periodic tracking area updating
// procedure accepted by the network
//    The MME may also include of list of equivalent PLMNs in the TRACKING AREA
//    UPDATE ACCEPT message. Each entry in the list contains a PLMN code
//    (MCC+MNC). The UE shall store the list as provided by the network, and if
//    there is no PDN connection for emergency bearer services established, the
//    UE shall remove from the list any PLMN code that is already in the list of
//    forbidden PLMNs. If there is a PDN connection for emergency bearer
//    services established, the UE shall remove from the list of equivalent
//    PLMNs any PLMN code present in the list of forbidden PLMNs when the PDN
//    connection for emergency bearer services is released. In addition, the UE
//    shall add to the stored list the PLMN code of the registered PLMN that
//    sent the list. The UE shall replace the stored list on each receipt of the
//    TRACKING AREA UPDATE ACCEPT message. If the TRACKING AREA UPDATE ACCEPT
//    message does not contain a list, then the UE shall delete the stored
//    list."

//#define R10_5_5_3_2_4__19 "Normal and periodic tracking area updating
// procedure accepted by the network
//    The network may also indicate in the EPS update result IE in the TRACKING
//    AREA UPDATE ACCEPT message that ISR is active. If the UE is attached for
//    emergency bearer services, the network shall indicate in the EPS update
//    result IE in the TRACKING AREA UPDATE ACCEPT message that ISR is not
//    activated. If the TRACKING AREA UPDATE ACCEPT message contains: i) no
//    indication that ISR is activated, the UE shall set the TIN to "GUTI"; ii)
//    an indication that ISR is activated, then:
//      - if the UE is required to perform routing area updating for IMS voice
//      termination as specified in
//        3GPP TS 24.008 [13], annex P.5, the UE shall set the TIN to "GUTI";
//      - if the UE had initiated the tracking area updating procedure due to a
//      change in UE network capability or
//        change in DRX parameters, the UE shall set the TIN to "GUTI"; or
//      - the UE shall regard a previously assigned P-TMSI and RAI as valid and
//      registered with the network. If the
//        TIN currently indicates "P-TMSI" and the periodic routing area update
//        timer T3312 is running, the UE shall set the TIN to "RAT-related
//        TMSI". If the TIN currently indicates "P-TMSI" and the periodic
//        routing area update timer T3312 has already expired, the UE shall set
//        the TIN to "GUTI"."

//#define R10_5_5_3_2_4__20 "Normal and periodic tracking area updating
// procedure accepted by the network
//    The network informs the UE about the support of specific features, such as
//    IMS voice over PS session, location services (EPC-LCS, CS-LCS) or
//    emergency bearer services, in the EPS network feature support information
//    element. In a UE with IMS voice over PS capability, the IMS voice over PS
//    session indicator and the emergency bearer services indicator shall be
//    provided to the upper layers. The upper layers take the IMS voice over PS
//    session indicator into account as specified in 3GPP TS 23.221 [8A],
//    subclause 7.2a and subclause 7.2b, when selecting the access domain for
//    voice sessions or calls. When initiating an emergency call, the upper
//    layers also take the emergency bearer services indicator into account for
//    the access domain selection. When the UE determines via the IMS voice over
//    PS session indicator that the network does not support IMS voice over PS
//    sessions in S1 mode, then the UE shall not locally release any persistent
//    EPS bearer context. When the UE determines via the emergency bearer
//    services indicator that the network does not support emergency bearer
//    services in S1 mode, then the UE shall not locally release any emergency
//    EPS bearer context if there is a radio bearer associated with that
//    context. In a UE with LCS capability, location services indicators
//    (EPC-LCS, CS-LCS) shall be provided to the upper layers. When MO-LR
//    procedure is triggered by the UE's application, those indicators are taken
//    into account as specified in 3GPP TS 24.171 [13C]."

//#define R10_5_5_3_2_4__21 "Normal and periodic tracking area updating
// procedure accepted by the network
//    If the UE has initiated the tracking area updating procedure due to manual
//    CSG selection and receives a TRACKING AREA UPDATE ACCEPT message, and the
//    UE sent the TRACKING AREA UPDATE REQUEST message in a CSG cell, the UE
//    shall check if the CSG ID and associated PLMN identity of the cell where
//    the UE has sent the TRACKING AREA UPDATE REQUEST message are contained in
//    the Allowed CSG list. If not, the UE shall add that CSG ID and associated
//    PLMN identity to the Allowed CSG list and the UE may add the HNB Name (if
//    provided by lower layers) to the Allowed CSG list if the HNB Name is
//    present in neither the Operator CSG list nor the Allowed CSG list."

//#define R10_5_5_3_2_4__22 "Normal and periodic tracking area updating
// procedure accepted by the network
//    If the TRACKING AREA UPDATE ACCEPT message contained a GUTI, the UE shall
//    return a TRACKING AREA UPDATE COMPLETE message to the MME to acknowledge
//    the received GUTI."

//#define R10_5_5_3_2_4__23 "Normal and periodic tracking area updating
// procedure accepted by the network
//    Upon receiving a TRACKING AREA UPDATE COMPLETE message, the MME shall stop
//    timer T3450, and shall consider the GUTI sent in the TRACKING AREA UPDATE
//    ACCEPT message as valid."

//...

//------------------------------

//#define R10_5_5_3_2_5__1 "Normal and periodic tracking area updating procedure
// not accepted by the network
//    If the tracking area updating cannot be accepted by the network, the MME
//    sends a TRACKING AREA UPDATE REJECT message to the UE including an
//    appropriate EMM cause value."

//#define R10_5_5_3_2_5__2 "Normal and periodic tracking area updating procedure
// not accepted by the network
//    If the MME locally deactivates EPS bearer contexts for the UE (see
//    subclause 5.5.3.2.4) and no active EPS bearer contexts remain for the UE,
//    the MME shall send the TRACKING AREA UPDATE REJECT message including the
//    EMM cause value #10 "Implicitly detached"."

//#define R10_5_5_3_2_5__3 "Normal and periodic tracking area updating procedure
// not accepted by the network
//    If the tracking area update request is rejected due to general NAS level
//    mobility management congestion control, the network shall set the EMM
//    cause value to #22 "congestion" and assign a back-off timer T3346."

//#define R10_5_5_3_2_5__4 "Normal and periodic tracking area updating procedure
// not accepted by the network
//    Upon receiving the TRACKING AREA UPDATE REJECT message, the UE shall stop
//    timer T3430, stop any transmission of user data, and take the following
//    actions depending on the EMM cause value received. #3 (Illegal UE); or #6
//    (Illegal ME);
//      The UE shall set the EPS update status to EU3 ROAMING NOT ALLOWED (and
//      shall store it according to subclause 5.1.3.3) and shall delete any
//      GUTI, last visited registered TAI, TAI list and eKSI. The UE shall
//      consider the USIM as invalid for EPS services until switching off or the
//      UICC containing the USIM is removed. The UE shall delete the list of
//      equivalent PLMNs and shall enter the state EMM-DEREGISTERED. If A/Gb
//      mode or Iu mode is supported by the UE, the UE shall handle the GMM
//      parameters GMM state, GPRS update status, P-TMSI, P-TMSI signature, RAI
//      and GPRS ciphering key sequence number and the MM parameters update
//      status, TMSI, LAI and ciphering key sequence number as specified in 3GPP
//      TS 24.008 [13] for the case when the normal routing area updating
//      procedure is rejected with the GMM cause with the same value. The USIM
//      shall be considered as invalid also for non-EPS services until switching
//      off or the UICC containing the USIM is removed.
//    NOTE 1: The possibility to configure a UE so that the radio transceiver
//    for a specific radio access technology is not
//      active, although it is implemented in the UE, is out of scope of the
//      present specification.
//    #7 (EPS services not allowed);
//      The UE shall set the EPS update status to EU3 ROAMING NOT ALLOWED (and
//      shall store it according to subclause 5.1.3.3) and shall delete any
//      GUTI, last visited registered TAI, TAI list and eKSI. The UE shall
//      consider the USIM as invalid for EPS services until switching off or the
//      UICC containing the USIM is removed. The UE shall delete the list of
//      equivalent PLMNs and shall enter the state EMM-DEREGISTERED. If the EPS
//      update type is "periodic updating", a UE in CS/PS mode 1 or CS/PS mode 2
//      of operation is still IMSI attached for non-EPS services. The UE shall
//      select GERAN or UTRAN radio access technology and proceed with
//      appropriate MM specific procedure according to the MM service state. The
//      UE shall not reselect E-UTRAN radio access technology until switching
//      off or the UICC containing the USIM is removed.
// ...

//------------------------------

//#define R10_5_5_3_2_7_a "Abnormal cases on the network side
//    If a lower layer failure occurs before the message TRACKING AREA UPDATE
//    COMPLETE has been received from the UE and a GUTI has been assigned, the
//    network shall abort the procedure and shall consider both, the old and new
//    GUTI as valid until the old GUTI can be considered as invalid by the
//    network (see subclause 5.4.1.4). During this period the network may use
//    the identification procedure followed by a GUTI reallocation procedure if
//    the old GUTI is used by the UE in a subsequent message. The network may
//    page with IMSI if paging with old and new S-TMSI fails. Paging with IMSI
//    causes the UE to re-attach as described in subclause 5.6.2.2.2."

//#define R10_5_5_3_2_7_b "Abnormal cases on the network side
//    Protocol error
//    If the TRACKING AREA UPDATE REQUEST message has been received with a
//    protocol error, the network shall return a TRACKING AREA UPDATE REJECT
//    message with one of the following EMM cause values: #96: invalid mandatory
//    information element error; #99: information element non-existent or not
//    implemented; #100: conditional IE error; or #111: protocol error,
//    unspecified."

//#define R10_5_5_3_2_7_c "Abnormal cases on the network side
//    T3450 time-out
//    On the first expiry of the timer, the network shall retransmit the
//    TRACKING AREA UPDATE ACCEPT message and shall reset and restart timer
//    T3450. The retransmission is performed four times, i.e. on the fifth
//    expiry of timer T3450, the tracking area updating procedure is aborted.
//    Both, the old and the new GUTI shall be considered as valid until the old
//    GUTI can be considered as invalid by the network (see subclause 5.4.1.4).
//    During this period the network acts as described for case a above."

//#define R10_5_5_3_2_7_d "Abnormal cases on the network side
//    TRACKING AREA UPDATE REQUEST received after the TRACKING AREA UPDATE
//    ACCEPT message has been sent and before the TRACKING AREA UPDATE COMPLETE
//    message is received
//    - If one or more of the information elements in the TRACKING AREA UPDATE
//    REQUEST message differ from the ones received within the previous TRACKING
//    AREA UPDATE REQUEST message, the previously initiated tracking area
//    updating procedure shall be aborted if the TRACKING AREA UPDATE COMPLETE
//    message has not been received and the new tracking area updating procedure
//    shall be progressed; or
//    - if the information elements do not differ, then the TRACKING AREA UPDATE
//    ACCEPT message shall be resent and the timer T3450 shall be restarted if
//    an TRACKING AREA UPDATE COMPLETE message is expected. In that case, the
//    retransmission counter related to T3450 is not incremented."

//#define R10_5_5_3_2_7_e "Abnormal cases on the network side
//    More than one TRACKING AREA UPDATE REQUEST received and no TRACKING AREA
//    UPDATE ACCEPT or TRACKING AREA UPDATE REJECT message has been sent
//    - If one or more of the information elements in the TRACKING AREA UPDATE
//    REQUEST message differs from the ones received within the previous
//    TRACKING AREA UPDATE REQUEST message, the previously initiated tracking
//    area updating procedure shall be aborted and the new tracking area
//    updating procedure shall be progressed;
//    - if the information elements do not differ, then the network shall
//    continue with the previous tracking area updating procedure and shall not
//    treat any further this TRACKING AREA UPDATE REQUEST message.

//#define R10_5_5_3_2_7_f "Abnormal cases on the network side
//    Lower layers indication of non-delivered NAS PDU due to handover
//    If the TRACKING AREA UPDATE ACCEPT message or TRACKING AREA UPDATE REJECT
//    message could not be delivered due to handover then the MME shall
//    retransmit the TRACKING AREA UPDATE ACCEPT message or TRACKING AREA UPDATE
//    REJECT message if the failure of handover procedure is reported by the
//    lower layer and the S1 signalling connection exists.
#define R10_9_9_3_7_1__1                                                       \
  "Detach type information element - Type of detach                    \
    All other values are interpreted as 'combined EPS/IMSI detach' in this version of the     \
    protocol."
#define R10_9_9_3_7_1__1_BRIEF                                                 \
  "Forced 'combined EPS/IMSI detach' to Type of detach"

#define R10_9_9_3_11__1                                                        \
  "EPS attach type value                                                \
    All other values are unused and shall be interpreted as 'EPS attach', if received by the  \
    network."
#define R10_9_9_3_11__1_BRIEF "Forced 'EPS attach' to EPS attach type value"

#endif /* FILE_3GPP_REQUIREMENTS_24_301_SEEN */
