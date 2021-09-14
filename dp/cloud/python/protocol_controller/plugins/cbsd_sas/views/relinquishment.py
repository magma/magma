from dp.cloud.python.protocol_controller.grpc_communication.get_common_rc_response import (
    get_common_bulk_rc_response,
)
from dp.cloud.python.protocol_controller.plugins.cbsd_sas.validators.relinquishment_request import (
    RelinquishmentRequestSchema,
)
from flask import Blueprint, request
from flask_json import as_json

relinquishment_page = Blueprint("relinquishment", __name__)


@relinquishment_page.route('/relinquishment', methods=('POST', ))
@as_json
def relinquishment():
    return get_common_bulk_rc_response(request, "relinquishmentResponse", RelinquishmentRequestSchema)
