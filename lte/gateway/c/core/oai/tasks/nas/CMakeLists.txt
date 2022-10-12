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
include_directories(${CMAKE_SOURCE_DIR})
set(libnas_api_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/api/network/nas_message.cpp
    )

set(libnas_mme_api_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/api/mme/mme_api.cpp
    )

set(libnas_emm_msg_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AttachAccept.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AttachComplete.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AttachReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AttachRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AuthenticationFailure.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AuthenticationReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AuthenticationRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/AuthenticationResponse.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/CsServiceNotification.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/DetachAccept.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/DetachRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/DownlinkNasTransport.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/EmmInformation.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/emm_msg.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/EmmStatus.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/ExtendedServiceRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/GutiReallocationCommand.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/GutiReallocationComplete.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/IdentityRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/IdentityResponse.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/SecurityModeCommand.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/SecurityModeComplete.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/SecurityModeReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/ServiceReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/ServiceRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/TrackingAreaUpdateAccept.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/TrackingAreaUpdateComplete.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/TrackingAreaUpdateReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/TrackingAreaUpdateRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/msg/UplinkNasTransport.cpp
    )

set(libnas_esm_msg_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDedicatedEpsBearerContextAccept.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDedicatedEpsBearerContextReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDedicatedEpsBearerContextRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDefaultEpsBearerContextAccept.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDefaultEpsBearerContextReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ActivateDefaultEpsBearerContextRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/BearerResourceAllocationReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/BearerResourceAllocationRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/BearerResourceModificationReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/BearerResourceModificationRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/DeactivateEpsBearerContextAccept.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/DeactivateEpsBearerContextRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/EsmInformationRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/EsmInformationResponse.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/esm_msg.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/EsmStatus.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ModifyEpsBearerContextAccept.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ModifyEpsBearerContextReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/ModifyEpsBearerContextRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/PdnConnectivityReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/PdnConnectivityRequest.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/PdnDisconnectReject.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/msg/PdnDisconnectRequest.cpp
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/ies)
set(libnas_ies_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/AdditionalUpdateResult.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/AdditionalUpdateType.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/ApnAggregateMaximumBitRate.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/Cli.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/CsfbResponse.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/DetachType.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EmmCause.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsAttachResult.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsAttachType.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsBearerContextStatus.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsMobileIdentity.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsNetworkFeatureSupport.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsQualityOfService.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsUpdateResult.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EpsUpdateType.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EsmCause.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EsmInformationTransferFlag.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/EsmMessageContainer.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/GutiType.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/KsiAndSequenceNumber.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/LcsClientIdentity.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/LcsIndicator.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/LinkedEpsBearerIdentity.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/MessageType.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/NasKeySetIdentifier.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/NasMessageContainer.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/NasRequestType.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/NasSecurityAlgorithms.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/Nonce.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/PagingIdentity.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/PdnAddress.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/PdnType.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/RadioPriority.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/SecurityHeaderType.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/ServiceType.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/ShortMac.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/SsCode.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/SupportedCodecList.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/TrackingAreaIdentity.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/TrackingAreaIdentityList.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/UeAdditionalSecurityCapability.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/UeNetworkCapability.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/UeRadioCapabilityInformationUpdateNeeded.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/UeSecurityCapability.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/ies/NrUeSecurityCapability.cpp
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/util)
set(libnas_utils_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/util/nas_timer.cpp
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/api/mme)
set(libnas_mme_api_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/api/mme/mme_api.cpp
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/emm)
set(libnas_mme_emm_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/Attach.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/Authentication.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/Detach.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/EmmInformation.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/emm_data_ctx.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/emm_main.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/EmmStatusHdl.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/Identification.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/LowerLayer.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/SecurityModeControl.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/ServiceRequestHdl.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/TrackingAreaUpdate.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/EmmCommon.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/NasTransportHdl.cpp
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/emm/sap)
set(libnas_mme_emm_sap_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_as.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_cn.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/EmmCommonProcedureInitiated.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/EmmDeregistered.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/EmmDeregisteredInitiated.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_esm.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_fsm.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_recv.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_reg.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/EmmRegistered.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_sap.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/emm/sap/emm_send.cpp
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/esm)
set(libnas_mme_esm_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/DedicatedEpsBearerContextActivation.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/DefaultEpsBearerContextActivation.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/EpsBearerContextDeactivation.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_data_context.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_ebr.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_ebr_context.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_ip.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_information.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_main.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/esm_pt.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/EsmStatusHdl.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/PdnConnectivity.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/PdnDisconnect.cpp
    )

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/esm/sap)
set(libnas_mme_esm_sap_OBJS
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/sap/esm_recv.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/sap/esm_send.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/esm/sap/esm_sap.cpp
    )

add_library(TASK_NAS
    nas_network.cpp
    nas_proc.cpp
    nas_procedures.cpp
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
