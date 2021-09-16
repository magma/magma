from magma.radio_controller.services.radio_controller.strategies.get_cbsd_id import (
    registration_get_cbsd_id,
    simple_get_cbsd_id,
)

get_cbsd_id_strategies = {
    "registrationRequest": registration_get_cbsd_id,
    "spectrumInquiryRequest": simple_get_cbsd_id,
    "grantRequest": simple_get_cbsd_id,
    "heartbeatRequest": simple_get_cbsd_id,
    "relinquishmentRequest": simple_get_cbsd_id,
    "deregistrationRequest": simple_get_cbsd_id,
}
