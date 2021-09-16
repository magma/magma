from magma.configuration_controller.response_processor.strategies.map_keys_generation import (
    generate_compound_request_map_key,
    generate_compound_response_map_key,
    generate_registration_request_map_key,
    generate_simple_request_map_key,
    generate_simple_response_map_key,
)
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
        "request_map_key": generate_registration_request_map_key,
        "response_map_key": generate_simple_response_map_key,
        "process_responses": process_registration_response,
    },
    "spectrumInquiryRequest": {
        "request_map_key": generate_simple_request_map_key,
        "response_map_key": generate_simple_response_map_key,
        "process_responses": process_spectrum_inquiry_response,
    },
    "grantRequest": {
        "request_map_key": generate_simple_request_map_key,
        "response_map_key": generate_simple_response_map_key,
        "process_responses": process_grant_response,
    },
    "heartbeatRequest": {
        "request_map_key": generate_compound_request_map_key,
        "response_map_key": generate_compound_response_map_key,
        "process_responses": process_heartbeat_response,
    },
    "relinquishmentRequest": {
        "request_map_key": generate_compound_request_map_key,
        "response_map_key": generate_compound_response_map_key,
        "process_responses": process_relinquishment_response,
    },
    "deregistrationRequest": {
        "request_map_key": generate_simple_request_map_key,
        "response_map_key": generate_simple_response_map_key,
        "process_responses": process_deregistration_response,
    },
}
