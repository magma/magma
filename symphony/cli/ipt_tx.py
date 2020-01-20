#!/usr/bin/env python3

import sys
from collections import namedtuple

import unicodecsv as csv
from pyinventory import InventoryClient
from pyinventory.reporter import FailedOperationException, InventoryReporter


def import_tx_row(client, data_identifier, data):
    try:
        location = client.add_location(
            data.Latitud,
            data.Longitud,
            {"Codigo Unico Estacion": data.CodigoUnico, "Direccion": data.DIRECCION},
            ("Departamento", data.Departamento),
            ("Provincia", data.Provincia),
            ("Distrito", data.Distrito),
            ("Centro Poblado", data.CentroPoblado),
            ("Estacion", data.Estacion),
        )
        vsat_equipment = None
        mw_equipment = None
        if data.VSATModel != "":
            equipment_type = "VSAT {}".format(data.VSATModel)
            vsat_equipment = client.add_equipment(
                data.VSATModel,
                equipment_type,
                location,
                {"SAT Band": data.SATBand, "Satellite": data.Satellite},
            )
        if data.MWModel != "":
            equipment_type = "MW{}".format(data.MWModel)
            mw_equipment = client.add_equipment(data.MWModel, equipment_type, location, {})
        if vsat_equipment is not None and mw_equipment is not None:
            client.add_link(vsat_equipment, "Port A", mw_equipment, "Port A")
    except FailedOperationException as e:
        print(e)
        e.log_failed_operation(data_identifier, data)
    except Exception as e:
        print(e)
        print(data)
        raise


def import_tx(email, password, csvPath, csvOutPath, csvErrPath):
    with open(csvPath, mode="rb") as infile:
        reader = csv.reader(infile, delimiter=",", encoding="utf-8")
        columns_row = next(reader)
        columns_row = ["ITEM"] + columns_row[1:]
        Data = namedtuple("Data", columns_row)
        reporter = InventoryReporter(csvOutPath, csvErrPath)
        client = InventoryClient(email, password, "ipt", reporter=reporter)
        for i, data in enumerate(map(Data._make, reader)):
            import_tx_row(client, "{}:{}".format(csvPath, i), data)


if __name__ == "__main__":
    if len(sys.argv) != 6:
        # flake8: noqa: E999
        print(
            "Usage: ipt_tx.py {email} {password} {csv_path} {csv_out_path} {csv_err_path}"
        )
        sys.exit(1)
    import_tx(sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4], sys.argv[5])
    sys.exit(0)
