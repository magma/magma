from magma.db_service.models import DBCbsd


def get_cbsd_basic_params(cbsd: DBCbsd) -> (str, str, str):
    network_id = ''
    fcc_id = ''
    cbsd_serial_number = ''
    if cbsd:
        network_id = cbsd.network_id or ''
        fcc_id = cbsd.fcc_id or ''
        cbsd_serial_number = cbsd.cbsd_serial_number or ''
    return fcc_id, network_id, cbsd_serial_number
