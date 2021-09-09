CBSD_SERIAL_NR = "cbsdSerialNumber"
FCC_ID = "fccId"
CBSD_ID = "cbsdId"


def registration_get_cbsd_id(request_payload):
    return f'{request_payload[FCC_ID]}/{request_payload[CBSD_SERIAL_NR]}'


def simple_get_cbsd_id(request_payload):
    return request_payload[CBSD_ID]
