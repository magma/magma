from flask import Blueprint, request
from flask_json import as_json

from dp.cloud.python.protocol_controller.grpc_communication.get_common_rc_response import get_common_bulk_rc_response
from dp.cloud.python.protocol_controller.plugins.cbsd_sas.validators.registration_request import RegistrationRequestSchema

registration_page = Blueprint("registration", __name__)


@registration_page.route('/registration', methods=('POST', ))
@as_json
def registration():
    return get_common_bulk_rc_response(request, "registrationResponse", RegistrationRequestSchema)
