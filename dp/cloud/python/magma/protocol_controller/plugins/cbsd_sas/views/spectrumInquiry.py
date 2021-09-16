from magma.protocol_controller.grpc_communication.get_common_rc_response import (
    get_common_bulk_rc_response,
)
from magma.protocol_controller.plugins.cbsd_sas.validators.spectrum_inquiry_request import (
    SpectrumInquiryRequestSchema,
)
from flask import Blueprint, request
from flask_json import as_json

spectrum_inquiry_page = Blueprint("spectrumInquiry", __name__)


@spectrum_inquiry_page.route('/spectrumInquiry', methods=('POST', ))
@as_json
def spectrum_inquiry():
    return get_common_bulk_rc_response(request, "spectrumInquiryResponse", SpectrumInquiryRequestSchema)
