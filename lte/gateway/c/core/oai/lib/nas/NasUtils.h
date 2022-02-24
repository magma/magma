#include <vector>
#include <string>

#include "NasEnum.h"

namespace nas {

/**** NasUtils****/
class NasUtils {
public:
    static std::vector<uint8_t> HexStringToVector(const std::string &hex)
    {
        if (hex.length() % 2 != 0)
            throw std::runtime_error("hex string has an odd length");

        for (char c : hex)
        {
            if (c >= '0' && c <= '9')
                continue;
            if (c >= 'a' && c <= 'f')
                continue;
            if (c >= 'A' && c <= 'F')
                continue;
            throw std::runtime_error("hex string contains invalid characters");
        }

        std::vector<uint8_t> bytes;
        for (unsigned int i = 0; i < hex.length(); i += 2)
        {
            std::string byteString = hex.substr(i, 2);
            char byte = (char)strtol(byteString.c_str(), nullptr, 16);
            bytes.emplace_back(byte);
        }
        return bytes;
    }
    
    static std::string Enum2String(RegistrationType v)
    {
        std::string str;
        switch (v)
        {
        case RegistrationType::INITIAL_REGISTRATION:
            str =  "Initial Registration";
        case RegistrationType::MOBILITY_REGISTRATION_UPDATING:
            str =   "Mobility Registration";
        case RegistrationType::PERIODIC_REGISTRATION_UPDATING:
            str =   "Periodic Registration";
        case RegistrationType::EMERGENCY_REGISTRATION:
            str =   "Emergency Registration";
        default:
            str =   "?";
        }
        return str;
    }
    static std::string Enum2String(ExtendedProtocolDiscriminator v)
    {
        std::string str;
        switch (v)
        {
        case ExtendedProtocolDiscriminator::MOBILITY_MANAGEMENT_MESSAGES :
            str = "MOBILITY_MANAGEMENT_MESSAGES";
            break;
        case ExtendedProtocolDiscriminator::SESSION_MANAGEMENT_MESSAGES :
            str = "SESSION_MANAGEMENT_MESSAGES";
            break;
        default:
            str =   "?";
            break;
        }
        return str;
    }

    static std::string Enum2String(MessageType msgtype) {
        std::string str;
        switch (msgtype)
        {
        case MessageType::REGISTRATION_REQUEST:
            str ="REGISTRATION_REQUEST";
            break;
        case MessageType::REGISTRATION_ACCEPT:
            str ="REGISTRATION_ACCEPT";
            break;
        case MessageType::REGISTRATION_COMPLETE:
            str ="REGISTRATION_COMPLETE";
            break;
        case MessageType::REGISTRATION_REJECT:
            str ="REGISTRATION_REJECT";
            break;
        case MessageType::DEREGISTRATION_REQUEST_UE_ORIGINATING:
            str ="DEREGISTRATION_REQUEST_UE_ORIGINATING";
            break;
        case MessageType::DEREGISTRATION_ACCEPT_UE_ORIGINATING:
            str ="DEREGISTRATION_ACCEPT_UE_ORIGINATING";
            break;
        case MessageType::DEREGISTRATION_REQUEST_UE_TERMINATED:
            str ="DEREGISTRATION_REQUEST_UE_TERMINATED";
            break;
        case MessageType::DEREGISTRATION_ACCEPT_UE_TERMINATED:
            str ="DEREGISTRATION_ACCEPT_UE_TERMINATED";
            break;
        case MessageType::SERVICE_REQUEST:
            str ="SERVICE_REQUEST";
            break;
        case MessageType::SERVICE_REJECT:
            str ="SERVICE_REJECT";
            break;
        case MessageType::SERVICE_ACCEPT:
            str ="SERVICE_ACCEPT";
            break;
        case MessageType::CONFIGURATION_UPDATE_COMMAND:
            str ="CONFIGURATION_UPDATE_COMMAND";
            break;
        case MessageType::CONFIGURATION_UPDATE_COMPLETE:
            str ="CONFIGURATION_UPDATE_COMPLETE";
            break;
        case MessageType::AUTHENTICATION_REQUEST:
            str ="AuthenticationRequest";
            break;
        case MessageType::AUTHENTICATION_RESPONSE:
            str ="AUTHENTICATION_REQUEST";
            break;
        case MessageType::AUTHENTICATION_REJECT:
            str ="AUTHENTICATION_REJECT";
            break;
        case MessageType::AUTHENTICATION_FAILURE:
            str ="AUTHENTICATION_FAILURE";
            break;
        case MessageType::AUTHENTICATION_RESULT:
            str ="AUTHENTICATION_RESULT";
            break;
        case MessageType::IDENTITY_REQUEST:
            str ="IDENTITY_REQUEST";
            break;
        case MessageType::IDENTITY_RESPONSE:
            str ="IDENTITY_RESPONSE";
            break;
        case MessageType::SECURITY_MODE_COMMAND:
            str ="SECURITY_MODE_COMMAND";
            break;
        case MessageType::SECURITY_MODE_COMPLETE:
            str ="SECURITY_MODE_COMPLETE";
            break;
        case MessageType::SECURITY_MODE_REJECT:
            str ="SECURITY_MODE_REJECT";
            break;
        case MessageType::FIVEG_MM_STATUS:
            str ="FIVEG_MM_STATUS";
            break;
        case MessageType::NOTIFICATION:
            str ="NOTIFICATION";
            break;
        case MessageType::NOTIFICATION_RESPONSE:
            str ="NOTIFICATION_RESPONSE";
            break;
        case MessageType::UL_NAS_TRANSPORT:
            str ="UL_NAS_TRANSPORT";
            break;
        case MessageType::DL_NAS_TRANSPORT:
            str ="DL_NAS_TRANSPORT";
            break;
        case MessageType::PDU_SESSION_ESTABLISHMENT_REQUEST:
            str ="PDU_SESSION_ESTABLISHMENT_REQUEST";
            break;
        case MessageType::PDU_SESSION_ESTABLISHMENT_ACCEPT:
            str ="PDU_SESSION_ESTABLISHMENT_ACCEPT";
            break;
        case MessageType::PDU_SESSION_ESTABLISHMENT_REJECT:
            str ="PDU_SESSION_ESTABLISHMENT_REJECT";
            break;
        case MessageType::PDU_SESSION_AUTHENTICATION_COMMAND:
            str ="PDU_SESSION_AUTHENTICATION_COMMAND";
            break;
        case MessageType::PDU_SESSION_AUTHENTICATION_COMPLETE:
            str ="PDU_SESSION_AUTHENTICATION_COMPLETE";
            break;
        case MessageType::PDU_SESSION_AUTHENTICATION_RESULT:
            str ="PDU_SESSION_AUTHENTICATION_RESULT";
            break;
        case MessageType::PDU_SESSION_MODIFICATION_REQUEST:
            str ="PDU_SESSION_MODIFICATION_REQUEST";
            break;
        case MessageType::PDU_SESSION_MODIFICATION_REJECT:
            str ="PDU_SESSION_MODIFICATION_REJECT";
            break;
        case MessageType::PDU_SESSION_MODIFICATION_COMMAND:
            str ="PDU_SESSION_MODIFICATION_COMMAND";
            break;
        case MessageType::PDU_SESSION_MODIFICATION_COMPLETE:
            str ="PDU_SESSION_MODIFICATION_COMPLETE";
            break;
        case MessageType::PDU_SESSION_MODIFICATION_COMMAND_REJECT:
            str ="PDU_SESSION_MODIFICATION_COMMAND_REJECT";
            break;
        case MessageType::PDU_SESSION_RELEASE_REQUEST:
            str ="PDU_SESSION_RELEASE_REQUEST";
            break;
        case MessageType::PDU_SESSION_RELEASE_REJECT:
            str ="PDU_SESSION_RELEASE_REJECT";
            break;
        case MessageType::PDU_SESSION_RELEASE_COMMAND:
            str ="PDU_SESSION_RELEASE_COMMAND";
            break;
        case MessageType::PDU_SESSION_RELEASE_COMPLETE:
            str ="PDU_SESSION_RELEASE_COMPLETE";
            break;
        case MessageType::FIVEG_SM_STATUS:
            str ="FIVEG_SM_STATUS";
            break;
        default:
            str ="?";
            break;
        }
        return str;
    }
    static std::string Enum2String(InformationElementType ieitype) {
        std::string str;

        switch (ieitype)
        {

        default:
            str ="?";
            break;
        }
        return str;
    }
#if 0
const char *EnumToString(MmCause v)
{
    switch (v)
    {
    case MmCause::ILLEGAL_UE:
        return "ILLEGAL_UE";
    case MmCause::PEI_NOT_ACCEPTED:
        return "PEI_NOT_ACCEPTED";
    case MmCause::ILLEGAL_ME:
        return "ILLEGAL_ME";
    case MmCause::FIVEG_SERVICES_NOT_ALLOWED:
        return "FIVEG_SERVICES_NOT_ALLOWED";
    case MmCause::UE_IDENTITY_CANNOT_BE_DERIVED_FROM_NETWORK:
        return "UE_IDENTITY_CANNOT_BE_DERIVED_FROM_NETWORK";
    case MmCause::IMPLICITY_DEREGISTERED:
        return "IMPLICITY_DEREGISTERED";
    case MmCause::PLMN_NOT_ALLOWED:
        return "PLMN_NOT_ALLOWED";
    case MmCause::TA_NOT_ALLOWED:
        return "TA_NOT_ALLOWED";
    case MmCause::ROAMING_NOT_ALLOWED_IN_TA:
        return "ROAMING_NOT_ALLOWED_IN_TA";
    case MmCause::NO_SUITIBLE_CELLS_IN_TA:
        return "NO_SUITIBLE_CELLS_IN_TA";
    case MmCause::MAC_FAILURE:
        return "MAC_FAILURE";
    case MmCause::SYNCH_FAILURE:
        return "SYNCH_FAILURE";
    case MmCause::CONGESTION:
        return "CONGESTION";
    case MmCause::UE_SECURITY_CAP_MISMATCH:
        return "UE_SECURITY_CAP_MISMATCH";
    case MmCause::SEC_MODE_REJECTED_UNSPECIFIED:
        return "SEC_MODE_REJECTED_UNSPECIFIED";
    case MmCause::NON_5G_AUTHENTICATION_UNACCEPTABLE:
        return "NON_5G_AUTHENTICATION_UNACCEPTABLE";
    case MmCause::N1_MODE_NOT_ALLOWED:
        return "N1_MODE_NOT_ALLOWED";
    case MmCause::RESTRICTED_SERVICE_AREA:
        return "RESTRICTED_SERVICE_AREA";
    case MmCause::LADN_NOT_AVAILABLE:
        return "LADN_NOT_AVAILABLE";
    case MmCause::MAX_PDU_SESSIONS_REACHED:
        return "MAX_PDU_SESSIONS_REACHED";
    case MmCause::INSUFFICIENT_RESOURCES_FOR_SLICE_AND_DNN:
        return "INSUFFICIENT_RESOURCES_FOR_SLICE_AND_DNN";
    case MmCause::INSUFFICIENT_RESOURCES_FOR_SLICE:
        return "INSUFFICIENT_RESOURCES_FOR_SLICE";
    case MmCause::NGKSI_ALREADY_IN_USE:
        return "NGKSI_ALREADY_IN_USE";
    case MmCause::NON_3GPP_ACCESS_TO_CN_NOT_ALLOWED:
        return "NON_3GPP_ACCESS_TO_CN_NOT_ALLOWED";
    case MmCause::SERVING_NETWORK_NOT_AUTHORIZED:
        return "SERVING_NETWORK_NOT_AUTHORIZED";
    case MmCause::PAYLOAD_NOT_FORWARDED:
        return "PAYLOAD_NOT_FORWARDED";
    case MmCause::DNN_NOT_SUPPORTED_OR_NOT_SUBSCRIBED:
        return "DNN_NOT_SUPPORTED_OR_NOT_SUBSCRIBED";
    case MmCause::INSUFFICIENT_USER_PLANE_RESOURCES:
        return "INSUFFICIENT_USER_PLANE_RESOURCES";
    case MmCause::SEMANTICALLY_INCORRECT_MESSAGE:
        return "SEMANTICALLY_INCORRECT_MESSAGE";
    case MmCause::INVALID_MANDATORY_INFORMATION:
        return "INVALID_MANDATORY_INFORMATION";
    case MmCause::MESSAGE_TYPE_NON_EXISTENT_OR_NOT_IMPLEMENTED:
        return "MESSAGE_TYPE_NON_EXISTENT_OR_NOT_IMPLEMENTED";
    case MmCause::MESSAGE_TYPE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE:
        return "MESSAGE_TYPE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE";
    case MmCause::INFORMATION_ELEMENT_NON_EXISTENT_OR_NOT_IMPLEMENTED:
        return "INFORMATION_ELEMENT_NON_EXISTENT_OR_NOT_IMPLEMENTED";
    case MmCause::CONDITIONAL_IE_ERROR:
        return "CONDITIONAL_IE_ERROR";
    case MmCause::MESSAGE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE:
        return "MESSAGE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE";
    case MmCause::UNSPECIFIED_PROTOCOL_ERROR:
        return "UNSPECIFIED_PROTOCOL_ERROR";
    default:
        return "?";
    }
}

const char *EnumToString(eap::ECode v)
{
    switch (v)
    {
    case eap::ECode::REQUEST:
        return "REQUEST";
    case eap::ECode::RESPONSE:
        return "RESPONSE";
    case eap::ECode::SUCCESS:
        return "SUCCESS";
    case eap::ECode::FAILURE:
        return "FAILURE";
    case eap::ECode::INITIATE:
        return "INITIATE";
    case eap::ECode::FINISH:
        return "FINISH";
    default:
        return "?";
    }
}

const char *EnumToString(ESmCause v)
{
    switch (v)
    {
    case ESmCause::INSUFFICIENT_RESOURCES:
        return "INSUFFICIENT_RESOURCES";
    case ESmCause::MISSING_OR_UNKNOWN_DNN:
        return "MISSING_OR_UNKNOWN_DNN";
    case ESmCause::UNKNOWN_PDU_SESSION_TYPE:
        return "UNKNOWN_PDU_SESSION_TYPE";
    case ESmCause::USER_AUTHENTICATION_OR_AUTHORIZATION_FAILED:
        return "USER_AUTHENTICATION_OR_AUTHORIZATION_FAILED";
    case ESmCause::REQUEST_REJECTED_UNSPECIFIED:
        return "REQUEST_REJECTED_UNSPECIFIED";
    case ESmCause::SERVICE_OPTION_TEMPORARILY_OUT_OF_ORDER:
        return "SERVICE_OPTION_TEMPORARILY_OUT_OF_ORDER";
    case ESmCause::PTI_ALREADY_IN_USE:
        return "PTI_ALREADY_IN_USE";
    case ESmCause::REGULAR_DEACTIVATION:
        return "REGULAR_DEACTIVATION";
    case ESmCause::REACTIVATION_REQUESTED:
        return "REACTIVATION_REQUESTED";
    case ESmCause::INVALID_PDU_SESSION_IDENTITY:
        return "INVALID_PDU_SESSION_IDENTITY";
    case ESmCause::SEMANTIC_ERRORS_IN_PACKET_FILTERS:
        return "SEMANTIC_ERRORS_IN_PACKET_FILTERS";
    case ESmCause::SYNTACTICAL_ERROR_IN_PACKET_FILTERS:
        return "SYNTACTICAL_ERROR_IN_PACKET_FILTERS";
    case ESmCause::OUT_OF_LADN_SERVICE_AREA:
        return "OUT_OF_LADN_SERVICE_AREA";
    case ESmCause::PTI_MISMATCH:
        return "PTI_MISMATCH";
    case ESmCause::PDU_SESSION_TYPE_IPV4_ONLY_ALLOWED:
        return "PDU_SESSION_TYPE_IPV4_ONLY_ALLOWED";
    case ESmCause::PDU_SESSION_TYPE_IPV6_ONLY_ALLOWED:
        return "PDU_SESSION_TYPE_IPV6_ONLY_ALLOWED";
    case ESmCause::PDU_SESSION_DOES_NOT_EXIST:
        return "PDU_SESSION_DOES_NOT_EXIST";
    case ESmCause::INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN:
        return "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN";
    case ESmCause::NOT_SUPPORTED_SSC_MODE:
        return "NOT_SUPPORTED_SSC_MODE";
    case ESmCause::INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE:
        return "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE";
    case ESmCause::MISSING_OR_UNKNOWN_DNN_IN_A_SLICE:
        return "MISSING_OR_UNKNOWN_DNN_IN_A_SLICE";
    case ESmCause::INVALID_PTI_VALUE:
        return "INVALID_PTI_VALUE";
    case ESmCause::MAXIMUM_DATA_RATE_PER_UE_FOR_USER_PLANE_INTEGRITY_PROTECTION_IS_TOO_LOW:
        return "MAXIMUM_DATA_RATE_PER_UE_FOR_USER_PLANE_INTEGRITY_PROTECTION_IS_TOO_LOW";
    case ESmCause::SEMANTIC_ERROR_IN_THE_QOS_OPERATION:
        return "SEMANTIC_ERROR_IN_THE_QOS_OPERATION";
    case ESmCause::SYNTACTICAL_ERROR_IN_THE_QOS_OPERATION:
        return "SYNTACTICAL_ERROR_IN_THE_QOS_OPERATION";
    case ESmCause::SEMANTICALLY_INCORRECT_MESSAGE:
        return "SEMANTICALLY_INCORRECT_MESSAGE";
    case ESmCause::INVALID_MANDATORY_INFORMATION:
        return "INVALID_MANDATORY_INFORMATION";
    case ESmCause::MESSAGE_TYPE_NON_EXISTENT_OR_NOT_IMPLEMENTED:
        return "MESSAGE_TYPE_NON_EXISTENT_OR_NOT_IMPLEMENTED";
    case ESmCause::MESSAGE_TYPE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE:
        return "MESSAGE_TYPE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE";
    case ESmCause::INFORMATION_ELEMENT_NON_EXISTENT_OR_NOT_IMPLEMENTED:
        return "INFORMATION_ELEMENT_NON_EXISTENT_OR_NOT_IMPLEMENTED";
    case ESmCause::CONDITIONAL_IE_ERROR:
        return "CONDITIONAL_IE_ERROR";
    case ESmCause::MESSAGE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE:
        return "MESSAGE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE";
    case ESmCause::PROTOCOL_ERROR_UNSPECIFIED:
        return "PROTOCOL_ERROR_UNSPECIFIED";
    default:
        return "?";
    }
}

const char *EnumToString(EPduSessionType v)
{
    switch (v)
    {
    case EPduSessionType::IPV4:
        return "IPV4";
    case EPduSessionType::IPV6:
        return "IPV6";
    case EPduSessionType::IPV4V6:
        return "IPV4V6";
    case EPduSessionType::UNSTRUCTURED:
        return "UNSTRUCTURED";
    case EPduSessionType::ETHERNET:
        return "ETHERNET";
    default:
        return "?";
    }
}
#endif
};

}