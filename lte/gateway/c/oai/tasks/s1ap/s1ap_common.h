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

/** @defgroup _s1ap_impl_ S1AP Layer Reference Implementation
 * @ingroup _ref_implementation_
 * @{
 */

#include <stdint.h>
#include <sys/types.h>

#include "asn_internal.h"
#include "constr_TYPE.h"

#if HAVE_CONFIG_H_
#include "config.h"
#endif

#ifndef FILE_S1AP_COMMON_SEEN
#define FILE_S1AP_COMMON_SEEN

#include "bstrlib.h"

/* Defined in asn_internal.h */
extern int asn_debug;

#if defined(EMIT_ASN_DEBUG_EXTERN)
inline void ASN_DEBUG(const char* fmt, ...);
#endif

#include "S1ap_InitiatingMessage.h"
#include "S1ap_ProtocolExtensionContainer.h"
#include "S1ap_ProtocolExtensionField.h"
#include "S1ap_ProtocolIE-ContainerPair.h"
#include "S1ap_ProtocolIE-Field.h"
#include "S1ap_ProtocolIE-FieldPair.h"
#include "S1ap_S1AP-PDU.h"
#include "S1ap_SuccessfulOutcome.h"
#include "S1ap_UnsuccessfulOutcome.h"
#include "S1ap_asn_constant.h"

#include "S1ap_AllocationAndRetentionPriority.h"
#include "S1ap_BPLMNs.h"
#include "S1ap_Bearers-SubjectToStatusTransfer-Item.h"
#include "S1ap_Bearers-SubjectToStatusTransferList.h"
#include "S1ap_BitRate.h"
#include "S1ap_BroadcastCompletedAreaList.h"
#include "S1ap_CGI.h"
#include "S1ap_CI.h"
#include "S1ap_CNDomain.h"
#include "S1ap_COUNTvalue.h"
#include "S1ap_CSFallbackIndicator.h"
#include "S1ap_CSG-Id.h"
#include "S1ap_CSG-IdList-Item.h"
#include "S1ap_CSG-IdList.h"
#include "S1ap_Cause.h"
#include "S1ap_CauseMisc.h"
#include "S1ap_CauseNas.h"
#include "S1ap_CauseProtocol.h"
#include "S1ap_CauseRadioNetwork.h"
#include "S1ap_CauseTransport.h"
#include "S1ap_Cdma2000HORequiredIndication.h"
#include "S1ap_Cdma2000HOStatus.h"
#include "S1ap_Cdma2000OneXMEID.h"
#include "S1ap_Cdma2000OneXMSI.h"
#include "S1ap_Cdma2000OneXPilot.h"
#include "S1ap_Cdma2000OneXRAND.h"
#include "S1ap_Cdma2000OneXSRVCCInfo.h"
#include "S1ap_Cdma2000PDU.h"
#include "S1ap_Cdma2000RATType.h"
#include "S1ap_Cdma2000SectorID.h"
#include "S1ap_Cell-Size.h"
#include "S1ap_CellID-Broadcast-Item.h"
#include "S1ap_CellID-Broadcast.h"
#include "S1ap_CellIdentity.h"
#include "S1ap_CellTrafficTrace.h"
#include "S1ap_CellType.h"
#include "S1ap_CompletedCellinEAI-Item.h"
#include "S1ap_CompletedCellinEAI.h"
#include "S1ap_CompletedCellinTAI-Item.h"
#include "S1ap_CompletedCellinTAI.h"
#include "S1ap_CriticalityDiagnostics-IE-Item.h"
#include "S1ap_CriticalityDiagnostics-IE-List.h"
#include "S1ap_CriticalityDiagnostics.h"
#include "S1ap_DL-Forwarding.h"
#include "S1ap_DataCodingScheme.h"
#include "S1ap_DeactivateTrace.h"
#include "S1ap_Direct-Forwarding-Path-Availability.h"
#include "S1ap_DownlinkNASTransport.h"
#include "S1ap_DownlinkS1cdma2000tunnelling.h"
#include "S1ap_E-RAB-ID.h"
#include "S1ap_E-RABAdmittedItem.h"
#include "S1ap_E-RABAdmittedList.h"
#include "S1ap_E-RABDataForwardingItem.h"
#include "S1ap_E-RABFailedToSetupItemHOReqAck.h"
#include "S1ap_E-RABFailedtoSetupListHOReqAck.h"
#include "S1ap_E-RABInformationList.h"
#include "S1ap_E-RABInformationListItem.h"
#include "S1ap_E-RABItem.h"
#include "S1ap_E-RABLevelQoSParameters.h"
#include "S1ap_E-RABList.h"
#include "S1ap_E-RABModificationIndication.h"
#include "S1ap_E-RABModifyItemBearerModRes.h"
#include "S1ap_E-RABModifyListBearerModRes.h"
#include "S1ap_E-RABModifyRequest.h"
#include "S1ap_E-RABModifyResponse.h"
#include "S1ap_E-RABReleaseCommand.h"
#include "S1ap_E-RABReleaseIndication.h"
#include "S1ap_E-RABReleaseItemBearerRelComp.h"
#include "S1ap_E-RABReleaseListBearerRelComp.h"
#include "S1ap_E-RABReleaseResponse.h"
#include "S1ap_E-RABSetupItemBearerSURes.h"
#include "S1ap_E-RABSetupItemCtxtSURes.h"
#include "S1ap_E-RABSetupListBearerSURes.h"
#include "S1ap_E-RABSetupListCtxtSURes.h"
#include "S1ap_E-RABSetupRequest.h"
#include "S1ap_E-RABSetupResponse.h"
#include "S1ap_E-RABSubjecttoDataForwardingList.h"
#include "S1ap_E-RABToBeModifiedItemBearerModReq.h"
#include "S1ap_E-RABToBeModifiedListBearerModReq.h"
#include "S1ap_E-RABToBeSetupItemBearerSUReq.h"
#include "S1ap_E-RABToBeSetupItemCtxtSUReq.h"
#include "S1ap_E-RABToBeSetupItemHOReq.h"
#include "S1ap_E-RABToBeSetupListBearerSUReq.h"
#include "S1ap_E-RABToBeSetupListCtxtSUReq.h"
#include "S1ap_E-RABToBeSetupListHOReq.h"
#include "S1ap_E-RABToBeSwitchedDLItem.h"
#include "S1ap_E-RABToBeSwitchedDLList.h"
#include "S1ap_E-RABToBeSwitchedULItem.h"
#include "S1ap_E-RABToBeSwitchedULList.h"
#include "S1ap_E-UTRAN-Trace-ID.h"
#include "S1ap_ECGIList.h"
#include "S1ap_ENB-ID.h"
#include "S1ap_ENB-StatusTransfer-TransparentContainer.h"
#include "S1ap_ENB-UE-S1AP-ID.h"
#include "S1ap_ENBConfigurationTransfer.h"
#include "S1ap_ENBConfigurationUpdate.h"
#include "S1ap_ENBConfigurationUpdateAcknowledge.h"
#include "S1ap_ENBConfigurationUpdateFailure.h"
#include "S1ap_ENBDirectInformationTransfer.h"
#include "S1ap_ENBStatusTransfer.h"
#include "S1ap_ENBX2TLAs.h"
#include "S1ap_ENBname.h"
#include "S1ap_EPLMNs.h"
#include "S1ap_EUTRAN-CGI.h"
#include "S1ap_EmergencyAreaID-Broadcast-Item.h"
#include "S1ap_EmergencyAreaID-Broadcast.h"
#include "S1ap_EmergencyAreaID.h"
#include "S1ap_EmergencyAreaIDList.h"
#include "S1ap_EncryptionAlgorithms.h"
#include "S1ap_ErrorIndication.h"
#include "S1ap_EventType.h"
#include "S1ap_ExtendedRNC-ID.h"
#include "S1ap_ForbiddenInterRATs.h"
#include "S1ap_ForbiddenLACs.h"
#include "S1ap_ForbiddenLAs-Item.h"
#include "S1ap_ForbiddenLAs.h"
#include "S1ap_ForbiddenTACs.h"
#include "S1ap_ForbiddenTAs-Item.h"
#include "S1ap_ForbiddenTAs.h"
#include "S1ap_GBR-QosInformation.h"
#include "S1ap_GERAN-Cell-ID.h"
#include "S1ap_GTP-TEID.h"
#include "S1ap_GUMMEI.h"
#include "S1ap_GUMMEIType.h"
#include "S1ap_Global-ENB-ID.h"
#include "S1ap_HFN.h"
#include "S1ap_HandoverCancel.h"
#include "S1ap_HandoverCancelAcknowledge.h"
#include "S1ap_HandoverCommand.h"
#include "S1ap_HandoverFailure.h"
#include "S1ap_HandoverNotify.h"
#include "S1ap_HandoverPreparationFailure.h"
#include "S1ap_HandoverRequest.h"
#include "S1ap_HandoverRequestAcknowledge.h"
#include "S1ap_HandoverRequired.h"
#include "S1ap_HandoverRestrictionList.h"
#include "S1ap_HandoverType.h"
#include "S1ap_IMSI.h"
#include "S1ap_InitialContextSetupFailure.h"
#include "S1ap_InitialContextSetupRequest.h"
#include "S1ap_InitialContextSetupResponse.h"
#include "S1ap_InitialUEMessage.h"
#include "S1ap_InitiatingMessage.h"
#include "S1ap_IntegrityProtectionAlgorithms.h"
#include "S1ap_Inter-SystemInformationTransferType.h"
#include "S1ap_InterfacesToTrace.h"
#include "S1ap_L3-Information.h"
#include "S1ap_LAC.h"
#include "S1ap_LAI.h"
#include "S1ap_LastVisitedCell-Item.h"
#include "S1ap_LastVisitedEUTRANCellInformation.h"
#include "S1ap_LastVisitedGERANCellInformation.h"
#include "S1ap_LastVisitedUTRANCellInformation.h"
#include "S1ap_LocationReport.h"
#include "S1ap_LocationReportingControl.h"
#include "S1ap_LocationReportingFailureIndication.h"
#include "S1ap_M-TMSI.h"
#include "S1ap_MME-Code.h"
#include "S1ap_MME-Group-ID.h"
#include "S1ap_MME-UE-S1AP-ID.h"
#include "S1ap_MMEConfigurationTransfer.h"
#include "S1ap_MMEConfigurationUpdate.h"
#include "S1ap_MMEConfigurationUpdateAcknowledge.h"
#include "S1ap_MMEConfigurationUpdateFailure.h"
#include "S1ap_MMEDirectInformationTransfer.h"
#include "S1ap_MMEStatusTransfer.h"
#include "S1ap_MMEname.h"
#include "S1ap_MSClassmark2.h"
#include "S1ap_MSClassmark3.h"
#include "S1ap_MessageIdentifier.h"
#include "S1ap_NAS-PDU.h"
#include "S1ap_NASNonDeliveryIndication.h"
#include "S1ap_NASSecurityParametersfromE-UTRAN.h"
#include "S1ap_NASSecurityParameterstoE-UTRAN.h"
#include "S1ap_NumberOfBroadcasts.h"
#include "S1ap_NumberofBroadcastRequest.h"
#include "S1ap_OldBSS-ToNewBSS-Information.h"
#include "S1ap_OverloadAction.h"
#include "S1ap_OverloadResponse.h"
#include "S1ap_OverloadStart.h"
#include "S1ap_OverloadStop.h"
#include "S1ap_PDCP-SN.h"
#include "S1ap_PLMNidentity.h"
#include "S1ap_Paging.h"
#include "S1ap_PagingDRX.h"
#include "S1ap_PathSwitchRequest.h"
#include "S1ap_PathSwitchRequestAcknowledge.h"
#include "S1ap_PathSwitchRequestFailure.h"
#include "S1ap_Pre-emptionCapability.h"
#include "S1ap_Pre-emptionVulnerability.h"
#include "S1ap_PriorityLevel.h"
#include "S1ap_PrivateMessage.h"
#include "S1ap_QCI.h"
#include "S1ap_RAC.h"
#include "S1ap_RIMInformation.h"
#include "S1ap_RIMRoutingAddress.h"
#include "S1ap_RIMTransfer.h"
#include "S1ap_RNC-ID.h"
#include "S1ap_RRC-Container.h"
#include "S1ap_RRC-Establishment-Cause.h"
#include "S1ap_ReceiveStatusofULPDCPSDUs.h"
#include "S1ap_RelativeMMECapacity.h"
#include "S1ap_RepetitionPeriod.h"
#include "S1ap_ReportArea.h"
#include "S1ap_RequestType.h"
#include "S1ap_Reset.h"
#include "S1ap_ResetAcknowledge.h"
#include "S1ap_ResetType.h"
#include "S1ap_S-TMSI.h"
#include "S1ap_S1SetupFailure.h"
#include "S1ap_S1SetupRequest.h"
#include "S1ap_S1SetupResponse.h"
#include "S1ap_SONConfigurationTransfer.h"
#include "S1ap_SONInformation.h"
#include "S1ap_SONInformationReply.h"
#include "S1ap_SONInformationRequest.h"
#include "S1ap_SRVCCHOIndication.h"
#include "S1ap_SRVCCOperationPossible.h"
#include "S1ap_SecurityContext.h"
#include "S1ap_SecurityKey.h"
#include "S1ap_SerialNumber.h"
#include "S1ap_ServedGUMMEIs.h"
#include "S1ap_ServedGUMMEIsItem.h"
#include "S1ap_ServedGroupIDs.h"
#include "S1ap_ServedMMECs.h"
#include "S1ap_ServedPLMNs.h"
#include "S1ap_Source-ToTarget-TransparentContainer.h"
#include "S1ap_SourceBSS-ToTargetBSS-TransparentContainer.h"
#include "S1ap_SourceRNC-ToTargetRNC-TransparentContainer.h"
#include "S1ap_SourceeNB-ID.h"
#include "S1ap_SourceeNB-ToTargeteNB-TransparentContainer.h"
#include "S1ap_SubscriberProfileIDforRFP.h"
#include "S1ap_SuccessfulOutcome.h"
#include "S1ap_SupportedTAs-Item.h"
#include "S1ap_SupportedTAs.h"
#include "S1ap_TAC.h"
#include "S1ap_TAI-Broadcast-Item.h"
#include "S1ap_TAI-Broadcast.h"
#include "S1ap_TAI.h"
#include "S1ap_TAIItem.h"
#include "S1ap_TAIList.h"
#include "S1ap_TAIListforWarning.h"
#include "S1ap_TBCD-STRING.h"
#include "S1ap_Target-ToSource-TransparentContainer.h"
#include "S1ap_TargetBSS-ToSourceBSS-TransparentContainer.h"
#include "S1ap_TargetID.h"
#include "S1ap_TargetRNC-ID.h"
#include "S1ap_TargetRNC-ToSourceRNC-TransparentContainer.h"
#include "S1ap_TargeteNB-ID.h"
#include "S1ap_TargeteNB-ToSourceeNB-TransparentContainer.h"
#include "S1ap_Time-UE-StayedInCell.h"
#include "S1ap_TimeToWait.h"
#include "S1ap_TraceActivation.h"
#include "S1ap_TraceDepth.h"
#include "S1ap_TraceFailureIndication.h"
#include "S1ap_TraceStart.h"
#include "S1ap_TransportLayerAddress.h"
#include "S1ap_TypeOfError.h"
#include "S1ap_UE-HistoryInformation.h"
#include "S1ap_UE-S1AP-ID-pair.h"
#include "S1ap_UE-S1AP-IDs.h"
#include "S1ap_UE-associatedLogicalS1-ConnectionItem.h"
#include "S1ap_UE-associatedLogicalS1-ConnectionListResAck.h"
#include "S1ap_UEAggregateMaximumBitrate.h"
#include "S1ap_UECapabilityInfoIndication.h"
#include "S1ap_UEContextModificationFailure.h"
#include "S1ap_UEContextModificationRequest.h"
#include "S1ap_UEContextModificationResponse.h"
#include "S1ap_UEContextReleaseCommand.h"
#include "S1ap_UEContextReleaseComplete.h"
#include "S1ap_UEContextReleaseRequest.h"
#include "S1ap_UEIdentityIndexValue.h"
#include "S1ap_UEPagingID.h"
#include "S1ap_UERadioCapability.h"
#include "S1ap_UESecurityCapabilities.h"
#include "S1ap_UnsuccessfulOutcome.h"
#include "S1ap_UplinkNASTransport.h"
#include "S1ap_UplinkS1cdma2000tunnelling.h"
#include "S1ap_WarningAreaList.h"
#include "S1ap_WarningMessageContents.h"
#include "S1ap_WarningSecurityInfo.h"
#include "S1ap_WarningType.h"
#include "S1ap_WriteReplaceWarningRequest.h"
#include "S1ap_WriteReplaceWarningResponse.h"
#include "S1ap_X2TNLConfigurationInfo.h"

// UPDATE RELEASE 9
#include "S1ap_BroadcastCancelledAreaList.h"
#include "S1ap_CSGMembershipStatus.h"
#include "S1ap_CancelledCellinEAI-Item.h"
#include "S1ap_CancelledCellinEAI.h"
#include "S1ap_CancelledCellinTAI-Item.h"
#include "S1ap_CancelledCellinTAI.h"
#include "S1ap_CellAccessMode.h"
#include "S1ap_CellID-Cancelled-Item.h"
#include "S1ap_CellID-Cancelled.h"
#include "S1ap_ConcurrentWarningMessageIndicator.h"
#include "S1ap_Data-Forwarding-Not-Possible.h"
#include "S1ap_DownlinkNonUEAssociatedLPPaTransport.h"
#include "S1ap_DownlinkUEAssociatedLPPaTransport.h"
#include "S1ap_E-RABList.h"
#include "S1ap_EUTRANRoundTripDelayEstimationInfo.h"
#include "S1ap_EmergencyAreaID-Cancelled-Item.h"
#include "S1ap_EmergencyAreaID-Cancelled.h"
#include "S1ap_ExtendedRepetitionPeriod.h"
#include "S1ap_KillRequest.h"
#include "S1ap_KillResponse.h"
#include "S1ap_LPPa-PDU.h"
#include "S1ap_PS-ServiceNotAvailable.h"
#include "S1ap_Routing-ID.h"
#include "S1ap_StratumLevel.h"
#include "S1ap_SynchronisationStatus.h"
#include "S1ap_TAI-Cancelled-Item.h"
#include "S1ap_TAI-Cancelled.h"
#include "S1ap_TimeSynchronisationInfo.h"
#include "S1ap_UplinkNonUEAssociatedLPPaTransport.h"
#include "S1ap_UplinkUEAssociatedLPPaTransport.h"

// UPDATE RELEASE 10
#include "S1ap_GUMMEIList.h"
#include "S1ap_GWContextReleaseIndication.h"
#include "S1ap_MMERelaySupportIndicator.h"
#include "S1ap_ManagementBasedMDTAllowed.h"
#include "S1ap_PagingPriority.h"
#include "S1ap_PrivacyIndicator.h"
#include "S1ap_RelayNode-Indicator.h"
#include "S1ap_TrafficLoadReductionIndication.h"

/* Checking version of ASN1C compiler */
#if (ASN1C_ENVIRONMENT_VERSION < ASN1C_MINIMUM_VERSION)
#error "You are compiling s1ap with the wrong version of ASN1C"
#endif

extern int asn_debug;
extern int asn1_xer_print;

#include <stdbool.h>

#include "mme_default_values.h"
#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "3gpp_33.401.h"
#include "security_types.h"
#include "common_types.h"
#include "s1ap_state.h"

#define S1AP_FIND_PROTOCOLIE_BY_ID(IE_TYPE, ie, container, IE_ID, mandatory)   \
  do {                                                                         \
    IE_TYPE** ptr;                                                             \
    ie = NULL;                                                                 \
    for (ptr = container->protocolIEs.list.array;                              \
         ptr < &container->protocolIEs.list                                    \
                    .array[container->protocolIEs.list.count];                 \
         ptr++) {                                                              \
      if ((*ptr)->id == IE_ID) {                                               \
        ie = *ptr;                                                             \
        break;                                                                 \
      }                                                                        \
    }                                                                          \
    if (ie == NULL) {                                                          \
      if (mandatory)                                                           \
        OAILOG_ERROR(                                                          \
            LOG_S1AP,                                                          \
            "S1AP_FIND_PROTOCOLIE_BY_ID %ld: %s %d: Mandatory ie is NULL\n",   \
            IE_ID, __FILE__, __LINE__);                                        \
      else                                                                     \
        OAILOG_DEBUG(                                                          \
            LOG_S1AP,                                                          \
            "S1AP_FIND_PROTOCOLIE_BY_ID %ld: %s %d: Optional ie is NULL\n",    \
            IE_ID, __FILE__, __LINE__);                                        \
    }                                                                          \
                                                                               \
  } while (0)

/** \brief Function callback prototype.
 **/
typedef int (*s1ap_message_handler_t)(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu);

/** \brief Handle criticality
 \param criticality Criticality of the IE
 @returns void
 **/
void s1ap_handle_criticality(S1ap_Criticality_t criticality);

#endif /* FILE_S1AP_COMMON_SEEN */
