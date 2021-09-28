from magma.configuration_controller.response_processor.strategies.response_processing import (
    process_deregistration_response,
    process_grant_response,
    process_heartbeat_response,
    process_registration_response,
    process_relinquishment_response,
    process_spectrum_inquiry_response,
)

# TODO use enum and constants here
processor_strategies = {
    "registrationRequest": {
        "process_responses": process_registration_response,
    },
    "spectrumInquiryRequest": {
        "process_responses": process_spectrum_inquiry_response,
    },
    "grantRequest": {
        "process_responses": process_grant_response,
    },
    "heartbeatRequest": {
        "process_responses": process_heartbeat_response,
    },
    "relinquishmentRequest": {
        "process_responses": process_relinquishment_response,
    },
    "deregistrationRequest": {
        "process_responses": process_deregistration_response,
    },
}
