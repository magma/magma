set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r10.5)

# compile the needed protos
set(NAS_STATE_CPP_PROTOS nas_state spgw_state common_types std_3gpp_types)
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(LTE_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${NAS_STATE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${STATE_PROTO_DIR} ${LTE_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")
include_directories("${S1AP_C_DIR}")
set(libnas_api_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/api/network/nas_message.c
    )

set(libnas_mme_api_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/api/mme/mme_api.c
    )

set(libnas_emm_msg_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AttachAccept.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AttachComplete.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AttachReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AttachRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AuthenticationFailure.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AuthenticationReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AuthenticationRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AuthenticationResponse.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/CsServiceNotification.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/DetachAccept.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/DetachRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/DownlinkNasTransport.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/EmmInformation.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/emm_msg.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/EmmStatus.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/ExtendedServiceRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/GutiReallocationCommand.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/GutiReallocationComplete.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/IdentityRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/IdentityResponse.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/SecurityModeCommand.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/SecurityModeComplete.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/SecurityModeReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/ServiceReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/ServiceRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/TrackingAreaUpdateAccept.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/TrackingAreaUpdateComplete.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/TrackingAreaUpdateReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/TrackingAreaUpdateRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/UplinkNasTransport.c
    )

set(libnas_esm_msg_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDedicatedEpsBearerContextAccept.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDedicatedEpsBearerContextReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDedicatedEpsBearerContextRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDefaultEpsBearerContextAccept.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDefaultEpsBearerContextReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDefaultEpsBearerContextRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/BearerResourceAllocationReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/BearerResourceAllocationRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/BearerResourceModificationReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/BearerResourceModificationRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/DeactivateEpsBearerContextAccept.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/DeactivateEpsBearerContextRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/EsmInformationRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/EsmInformationResponse.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/esm_msg.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/EsmStatus.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ModifyEpsBearerContextAccept.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ModifyEpsBearerContextReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ModifyEpsBearerContextRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/PdnConnectivityReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/PdnConnectivityRequest.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/PdnDisconnectReject.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/PdnDisconnectRequest.c
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/ies)
set(libnas_ies_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/AdditionalUpdateResult.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/AdditionalUpdateType.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/ApnAggregateMaximumBitRate.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/Cli.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/CsfbResponse.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/DetachType.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EmmCause.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsAttachResult.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsAttachType.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsBearerContextStatus.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsMobileIdentity.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsNetworkFeatureSupport.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsQualityOfService.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsUpdateResult.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsUpdateType.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EsmCause.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EsmInformationTransferFlag.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EsmMessageContainer.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/GutiType.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/KsiAndSequenceNumber.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/LcsClientIdentity.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/LcsIndicator.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/LinkedEpsBearerIdentity.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/MessageType.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/NasKeySetIdentifier.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/NasMessageContainer.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/NasRequestType.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/NasSecurityAlgorithms.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/Nonce.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/PagingIdentity.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/PdnAddress.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/PdnType.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/RadioPriority.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/SecurityHeaderType.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/ServiceType.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/ShortMac.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/SsCode.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/SupportedCodecList.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/TrackingAreaIdentity.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/TrackingAreaIdentityList.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/UeAdditionalSecurityCapability.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/UeNetworkCapability.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/UeRadioCapabilityInformationUpdateNeeded.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/UeSecurityCapability.c
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/NrUeSecurityCapability.c
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/util)
set(libnas_utils_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/util/nas_timer.c
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/api/mme)
set(libnas_mme_api_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/api/mme/mme_api.c
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/emm)
set(libnas_mme_emm_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/Attach.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/Authentication.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/Detach.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/EmmInformation.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/emm_data_ctx.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/emm_main.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/EmmStatusHdl.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/Identification.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/LowerLayer.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/SecurityModeControl.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/ServiceRequestHdl.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/TrackingAreaUpdate.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/EmmCommon.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/NasTransportHdl.c
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/emm/sap)
set(libnas_mme_emm_sap_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_as.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_cn.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/EmmCommonProcedureInitiated.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/EmmDeregistered.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/EmmDeregisteredInitiated.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_esm.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_fsm.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_recv.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_reg.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/EmmRegistered.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_sap.c
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_send.c
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/esm)
set(libnas_mme_esm_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/DedicatedEpsBearerContextActivation.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/DefaultEpsBearerContextActivation.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/EpsBearerContextDeactivation.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_data_context.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_ebr.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_ebr_context.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_ip.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_information.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_main.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_pt.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/EsmStatusHdl.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/PdnConnectivity.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/PdnDisconnect.c
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/esm/sap)
set(libnas_mme_esm_sap_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/sap/esm_recv.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/sap/esm_send.c
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/sap/esm_sap.c
    )

add_library(TASK_NAS
    nas_if_s1.c
    nas_if_s6a.c
    nas_network.c
    nas_proc.c
    nas_procedures.c
    nas_state_converter.cpp
    ${libnas_api_OBJS}
    ${libnas_mme_api_OBJS}
    ${libnas_emm_msg_OBJS}
    ${libnas_esm_msg_OBJS}
    ${libnas_ies_OBJS}
    ${libnas_utils_OBJS}
    ${libnas_mme_emm_OBJS}
    ${libnas_mme_emm_sap_OBJS}
    ${libnas_mme_esm_OBJS}
    ${libnas_mme_esm_sap_OBJS}
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(TASK_NAS
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_SECU LIB_3GPP
    TASK_S1AP TASK_MME_APP TASK_SERVICE303 TASK_S6A
    )
if (EMBEDDED_SGW)
  target_link_libraries(TASK_NAS TASK_SGW)
endif (EMBEDDED_SGW)

# TODO leak baby leak
target_include_directories(TASK_NAS PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/util
    ${CMAKE_CURRENT_SOURCE_DIR}/esm
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/sap
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg
    ${CMAKE_CURRENT_SOURCE_DIR}/ies
    ${CMAKE_CURRENT_SOURCE_DIR}/emm
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg
    ${CMAKE_CURRENT_SOURCE_DIR}/api
    ${CMAKE_CURRENT_SOURCE_DIR}/api/mme
    ${CMAKE_CURRENT_SOURCE_DIR}/api/network
    $<TARGET_FILE_DIR:TASK_NAS>
    )
