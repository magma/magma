/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */

/*******************************************************************************
 * This file had been created by asn1tostruct.py script v1.0.2
 * Please do not modify this file but regenerate it via script.
 * Created on: 2020-08-11 20:47:34.759997 by calsoft
 * from ['NGAP-PDU-Contents.asn']
 ******************************************************************************/
#include "ngap_common.h"
#include "ngap_ies_defs.h"
#include "log.h"

int ngap_decode_initialcontextsetupresponseies(
    InitialContextSetupResponseIEs_t *initialContextSetupResponseIEs,
    ANY_t *any_p) {

    InitialContextSetupResponse_t  initialContextSetupResponse;
    InitialContextSetupResponse_t *initialContextSetupResponse_p = &initialContextSetupResponse;
    int i, decoded = 0;
    int tempDecoded = 0;
    assert(any_p != NULL);
    assert(initialContextSetupResponseIEs != NULL);

    memset(initialContextSetupResponseIEs, 0, sizeof(InitialContextSetupResponseIEs_t));
   OAILOG_DEBUG (LOG_NGAP, "Decoding message InitialContextSetupResponseIEs (%s:%d)\n", __FILE__, __LINE__);

    ANY_to_type_aper(any_p, &asn_DEF_InitialContextSetupResponse, (void**)&initialContextSetupResponse_p);

    for (i = 0; i < initialContextSetupResponse_p->initialContextSetupResponse_ies.list.count; i++) {
        Ngap_IE_t *ie_p;
        ie_p = initialContextSetupResponse_p->initialContextSetupResponse_ies.list.array[i];
        switch(ie_p->id) {
            case Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID:
            {
                AMF_UE_NGAP_ID_t *amfuengapid_p = NULL;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_AMF_UE_NGAP_ID, (void**)&amfuengapid_p);
                if (tempDecoded < 0 || amfuengapid_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE amf_ue_ngap_id failed\n");
                    if (amfuengapid_p)
                        ASN_STRUCT_FREE(asn_DEF_AMF_UE_NGAP_ID, amfuengapid_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_AMF_UE_NGAP_ID, amfuengapid_p);
                memcpy(&initialContextSetupResponseIEs->amf_ue_ngap_id, amfuengapid_p, sizeof(AMF_UE_NGAP_ID_t));
                FREEMEM(amfuengapid_p);
                amfuengapid_p = NULL;
            } break;
            case Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID:
            {
                RAN_UE_NGAP_ID_t *ranuengapid_p = NULL;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_RAN_UE_NGAP_ID, (void**)&ranuengapid_p);
                if (tempDecoded < 0 || ranuengapid_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE ran_ue_ngap_id failed\n");
                    if (ranuengapid_p)
                        ASN_STRUCT_FREE(asn_DEF_RAN_UE_NGAP_ID, ranuengapid_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_RAN_UE_NGAP_ID, ranuengapid_p);
                memcpy(&initialContextSetupResponseIEs->ran_ue_ngap_id, ranuengapid_p, sizeof(RAN_UE_NGAP_ID_t));
                FREEMEM(ranuengapid_p);
                ranuengapid_p = NULL;
            } break;
            default:
               OAILOG_ERROR (LOG_NGAP, "Unknown protocol IE id (%d) for message initialcontextsetupresponseies\n", (int)ie_p->id);
        }
    }
    ASN_STRUCT_FREE(asn_DEF_InitialContextSetupResponse, initialContextSetupResponse_p);
    return decoded;
}

int ngap_decode_initialcontextsetuprequesties(
    InitialContextSetupRequestIEs_t *initialContextSetupRequestIEs,
    ANY_t *any_p) {

    InitialContextSetupRequest_t  initialContextSetupRequest;
    InitialContextSetupRequest_t *initialContextSetupRequest_p = &initialContextSetupRequest;
    int i, decoded = 0;
    int tempDecoded = 0;
    assert(any_p != NULL);
    assert(initialContextSetupRequestIEs != NULL);

    memset(initialContextSetupRequestIEs, 0, sizeof(InitialContextSetupRequestIEs_t));
   OAILOG_DEBUG (LOG_NGAP, "Decoding message InitialContextSetupRequestIEs (%s:%d)\n", __FILE__, __LINE__);

    ANY_to_type_aper(any_p, &asn_DEF_InitialContextSetupRequest, (void**)&initialContextSetupRequest_p);

    for (i = 0; i < initialContextSetupRequest_p->initialContextSetupRequest_ies.list.count; i++) {
        Ngap_IE_t *ie_p;
        ie_p = initialContextSetupRequest_p->initialContextSetupRequest_ies.list.array[i];
        switch(ie_p->id) {
            case Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID:
            {
                AMF_UE_NGAP_ID_t *amfuengapid_p = NULL;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_AMF_UE_NGAP_ID, (void**)&amfuengapid_p);
                if (tempDecoded < 0 || amfuengapid_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE amf_ue_ngap_id failed\n");
                    if (amfuengapid_p)
                        ASN_STRUCT_FREE(asn_DEF_AMF_UE_NGAP_ID, amfuengapid_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_AMF_UE_NGAP_ID, amfuengapid_p);
                memcpy(&initialContextSetupRequestIEs->amf_ue_ngap_id, amfuengapid_p, sizeof(AMF_UE_NGAP_ID_t));
                FREEMEM(amfuengapid_p);
                amfuengapid_p = NULL;
            } break;
            case Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID:
            {
                RAN_UE_NGAP_ID_t *ranuengapid_p = NULL;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_RAN_UE_NGAP_ID, (void**)&ranuengapid_p);
                if (tempDecoded < 0 || ranuengapid_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE ran_ue_ngap_id failed\n");
                    if (ranuengapid_p)
                        ASN_STRUCT_FREE(asn_DEF_RAN_UE_NGAP_ID, ranuengapid_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_RAN_UE_NGAP_ID, ranuengapid_p);
                memcpy(&initialContextSetupRequestIEs->ran_ue_ngap_id, ranuengapid_p, sizeof(RAN_UE_NGAP_ID_t));
                FREEMEM(ranuengapid_p);
                ranuengapid_p = NULL;
            } break;
            case Ngap_ProtocolIE_ID_id_GUAMI:
            {
                GUAMI_t *guami_p = NULL;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_GUAMI, (void**)&guami_p);
                if (tempDecoded < 0 || guami_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE guami failed\n");
                    if (guami_p)
                        ASN_STRUCT_FREE(asn_DEF_GUAMI, guami_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_GUAMI, guami_p);
                memcpy(&initialContextSetupRequestIEs->guami, guami_p, sizeof(GUAMI_t));
                FREEMEM(guami_p);
                guami_p = NULL;
            } break;
            /* Optional field */
            case Ngap_ProtocolIE_ID_id_MaskedIMEISV:
            {
                MaskedIMEISV_t *maskedIMEISV_p = NULL;
                initialContextSetupRequestIEs->presenceMask |= INITIALCONTEXTSETUPREQUESTIES_MASKEDIMEISV_PRESENT;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_MaskedIMEISV, (void**)&maskedIMEISV_p);
                if (tempDecoded < 0 || maskedIMEISV_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE maskedIMEISV failed\n");
                    if (maskedIMEISV_p)
                        ASN_STRUCT_FREE(asn_DEF_MaskedIMEISV, maskedIMEISV_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_MaskedIMEISV, maskedIMEISV_p);
                memcpy(&initialContextSetupRequestIEs->maskedIMEISV, maskedIMEISV_p, sizeof(MaskedIMEISV_t));
                FREEMEM(maskedIMEISV_p);
                maskedIMEISV_p = NULL;
            } break;
            /* Optional field */
            case Ngap_ProtocolIE_ID_id_NAS_PDU:
            {
                NAS_PDU_t *naspdu_p = NULL;
                initialContextSetupRequestIEs->presenceMask |= INITIALCONTEXTSETUPREQUESTIES_NAS_PDU_PRESENT;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_NAS_PDU, (void**)&naspdu_p);
                if (tempDecoded < 0 || naspdu_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE nas_pdu failed\n");
                    if (naspdu_p)
                        ASN_STRUCT_FREE(asn_DEF_NAS_PDU, naspdu_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_NAS_PDU, naspdu_p);
                memcpy(&initialContextSetupRequestIEs->nas_pdu, naspdu_p, sizeof(NAS_PDU_t));
                FREEMEM(naspdu_p);
                naspdu_p = NULL;
            } break;
            default:
               OAILOG_ERROR (LOG_NGAP, "Unknown protocol IE id (%d) for message initialcontextsetuprequesties\n", (int)ie_p->id);
        }
    }
    ASN_STRUCT_FREE(asn_DEF_InitialContextSetupRequest, initialContextSetupRequest_p);
    return decoded;
}

int ngap_decode_initialuemessage_ies(
    InitialUEMessage_IEs_t *initialUEMessage_IEs,
    ANY_t *any_p) {

    InitialUEMessage_t  initialUEMessage;
    InitialUEMessage_t *initialUEMessage_p = &initialUEMessage;
    int i, decoded = 0;
    int tempDecoded = 0;
    assert(any_p != NULL);
    assert(initialUEMessage_IEs != NULL);

    memset(initialUEMessage_IEs, 0, sizeof(InitialUEMessage_IEs_t));
   OAILOG_DEBUG (LOG_NGAP, "Decoding message InitialUEMessage_IEs (%s:%d)\n", __FILE__, __LINE__);

    ANY_to_type_aper(any_p, &asn_DEF_InitialUEMessage, (void**)&initialUEMessage_p);

    for (i = 0; i < initialUEMessage_p->initialUEMessage_ies.list.count; i++) {
        Ngap_IE_t *ie_p;
        ie_p = initialUEMessage_p->initialUEMessage_ies.list.array[i];
        switch(ie_p->id) {
            case Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID:
            {
                RAN_UE_NGAP_ID_t *ranuengapid_p = NULL;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_RAN_UE_NGAP_ID, (void**)&ranuengapid_p);
                if (tempDecoded < 0 || ranuengapid_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE ran_ue_ngap_id failed\n");
                    if (ranuengapid_p)
                        ASN_STRUCT_FREE(asn_DEF_RAN_UE_NGAP_ID, ranuengapid_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_RAN_UE_NGAP_ID, ranuengapid_p);
                memcpy(&initialUEMessage_IEs->ran_ue_ngap_id, ranuengapid_p, sizeof(RAN_UE_NGAP_ID_t));
                FREEMEM(ranuengapid_p);
                ranuengapid_p = NULL;
            } break;
            case Ngap_ProtocolIE_ID_id_NAS_PDU:
            {
                NAS_PDU_t *naspdu_p = NULL;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_NAS_PDU, (void**)&naspdu_p);
                if (tempDecoded < 0 || naspdu_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE nas_pdu failed\n");
                    if (naspdu_p)
                        ASN_STRUCT_FREE(asn_DEF_NAS_PDU, naspdu_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_NAS_PDU, naspdu_p);
                memcpy(&initialUEMessage_IEs->nas_pdu, naspdu_p, sizeof(NAS_PDU_t));
                FREEMEM(naspdu_p);
                naspdu_p = NULL;
            } break;
            case Ngap_ProtocolIE_ID_id_UserLocationInformation:
            {
                UserLocationInformation_t *userLocationInformation_p = NULL;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_UserLocationInformation, (void**)&userLocationInformation_p);
                if (tempDecoded < 0 || userLocationInformation_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE userLocationInformation failed\n");
                    if (userLocationInformation_p)
                        ASN_STRUCT_FREE(asn_DEF_UserLocationInformation, userLocationInformation_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_UserLocationInformation, userLocationInformation_p);
                memcpy(&initialUEMessage_IEs->userLocationInformation, userLocationInformation_p, sizeof(UserLocationInformation_t));
                FREEMEM(userLocationInformation_p);
                userLocationInformation_p = NULL;
            } break;
            case Ngap_ProtocolIE_ID_id_RRCEstablishmentCause:
            {
                RRCEstablishmentCause_t *rrcEstablishmentCause_p = NULL;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_RRCEstablishmentCause, (void**)&rrcEstablishmentCause_p);
                if (tempDecoded < 0 || rrcEstablishmentCause_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE rrcEstablishmentCause failed\n");
                    if (rrcEstablishmentCause_p)
                        ASN_STRUCT_FREE(asn_DEF_RRCEstablishmentCause, rrcEstablishmentCause_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_RRCEstablishmentCause, rrcEstablishmentCause_p);
                memcpy(&initialUEMessage_IEs->rrcEstablishmentCause, rrcEstablishmentCause_p, sizeof(RRCEstablishmentCause_t));
                FREEMEM(rrcEstablishmentCause_p);
                rrcEstablishmentCause_p = NULL;
            } break;
            /* Optional field */
            case Ngap_ProtocolIE_ID_id_UEContextRequest:
            {
                UEContextRequest_t *ueContextRequest_p = NULL;
                initialUEMessage_IEs->presenceMask |= INITIALUEMESSAGE_IES_UECONTEXTREQUEST_PRESENT;
                tempDecoded = ANY_to_type_aper(&ie_p->value, &asn_DEF_UEContextRequest, (void**)&ueContextRequest_p);
                if (tempDecoded < 0 || ueContextRequest_p == NULL) {
                   OAILOG_ERROR (LOG_NGAP, "Decoding of IE ueContextRequest failed\n");
                    if (ueContextRequest_p)
                        ASN_STRUCT_FREE(asn_DEF_UEContextRequest, ueContextRequest_p);
                    return -1;
                }
                decoded += tempDecoded;
                if (asn1_xer_print)
                    xer_fprint(stdout, &asn_DEF_UEContextRequest, ueContextRequest_p);
                memcpy(&initialUEMessage_IEs->ueContextRequest, ueContextRequest_p, sizeof(UEContextRequest_t));
                FREEMEM(ueContextRequest_p);
                ueContextRequest_p = NULL;
            } break;
            default:
               OAILOG_ERROR (LOG_NGAP, "Unknown protocol IE id (%d) for message initialuemessage_ies\n", (int)ie_p->id);
        }
    }
    ASN_STRUCT_FREE(asn_DEF_InitialUEMessage, initialUEMessage_p);
    return decoded;
}

int free_initialcontextsetupresponse(
    InitialContextSetupResponseIEs_t *initialContextSetupResponseIEs) {

    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_AMF_UE_NGAP_ID, &initialContextSetupResponseIEs->amf_ue_ngap_id);
    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_RAN_UE_NGAP_ID, &initialContextSetupResponseIEs->ran_ue_ngap_id);
    return 0;
}

int free_initialcontextsetuprequest(
    InitialContextSetupRequestIEs_t *initialContextSetupRequestIEs) {

    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_AMF_UE_NGAP_ID, &initialContextSetupRequestIEs->amf_ue_ngap_id);
    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_RAN_UE_NGAP_ID, &initialContextSetupRequestIEs->ran_ue_ngap_id);
    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_GUAMI, &initialContextSetupRequestIEs->guami);
    /* Optional field */
    if ((initialContextSetupRequestIEs->presenceMask & INITIALCONTEXTSETUPREQUESTIES_MASKEDIMEISV_PRESENT)
        == INITIALCONTEXTSETUPREQUESTIES_MASKEDIMEISV_PRESENT) 
        ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_MaskedIMEISV, &initialContextSetupRequestIEs->maskedIMEISV);
    /* Optional field */
    if ((initialContextSetupRequestIEs->presenceMask & INITIALCONTEXTSETUPREQUESTIES_NAS_PDU_PRESENT)
        == INITIALCONTEXTSETUPREQUESTIES_NAS_PDU_PRESENT) 
        ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_NAS_PDU, &initialContextSetupRequestIEs->nas_pdu);
    return 0;
}

int free_initialuemessage_(
    InitialUEMessage_IEs_t *initialUEMessage_IEs) {

    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_RAN_UE_NGAP_ID, &initialUEMessage_IEs->ran_ue_ngap_id);
    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_NAS_PDU, &initialUEMessage_IEs->nas_pdu);
    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_UserLocationInformation, &initialUEMessage_IEs->userLocationInformation);
    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_RRCEstablishmentCause, &initialUEMessage_IEs->rrcEstablishmentCause);
    /* Optional field */
    if ((initialUEMessage_IEs->presenceMask & INITIALUEMESSAGE_IES_UECONTEXTREQUEST_PRESENT)
        == INITIALUEMESSAGE_IES_UECONTEXTREQUEST_PRESENT) 
        ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_UEContextRequest, &initialUEMessage_IEs->ueContextRequest);
    return 0;
}

