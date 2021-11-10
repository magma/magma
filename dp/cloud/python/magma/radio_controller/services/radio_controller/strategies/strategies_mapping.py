from magma.radio_controller.services.radio_controller.strategies.get_cbsd_filters import (
    registration_get_cbsd_filters,
    simple_get_cbsd_filters,
)

get_cbsd_filter_strategies = {
    "registrationRequest": registration_get_cbsd_filters,
    "spectrumInquiryRequest": simple_get_cbsd_filters,
    "grantRequest": simple_get_cbsd_filters,
    "heartbeatRequest": simple_get_cbsd_filters,
    "relinquishmentRequest": simple_get_cbsd_filters,
    "deregistrationRequest": simple_get_cbsd_filters,
}
