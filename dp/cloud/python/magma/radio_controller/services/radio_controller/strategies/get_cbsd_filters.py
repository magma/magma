from magma.db_service.models import DBCbsd

CBSD_SERIAL_NR = "cbsdSerialNumber"
FCC_ID = "fccId"
USER_ID = "userId"
CBSD_ID = "cbsdId"


def registration_get_cbsd_filters(request_payload):
    return [
        DBCbsd.fcc_id == request_payload.get(FCC_ID),
        DBCbsd.user_id == request_payload.get(USER_ID),
        DBCbsd.cbsd_serial_number == request_payload.get(CBSD_SERIAL_NR),
    ]


def simple_get_cbsd_filters(request_payload):
    return [
        DBCbsd.cbsd_id == request_payload.get(CBSD_ID),
    ]
