#!/usr/bin/env python3

from .api.equipment_type import _populate_equipment_types
from .api.location_type import _populate_location_types
from .api.service import _populate_service_types
from .consts import DUMMY_REPORTER
from .graphql_client import GraphqlClient
from .reporter import Reporter


"""Pyinventory is a python package that allows for querying and modifying the
FBC platform inventory using graphql.

This module contains the client that allows to connect to inventory. The client
allows different kind of operations: querying and creating of locations, equipments,
positions and links.

Example of how to connect:
    from pyinventory import InventoryClient
    # since inventory is multi tenant system you will need to insert which
    # partner you connect as
    client = InventoryClient(email, password, "ipt")
    location = client.addLocation(-1.22,2.66, ('City', 'Lima'))
    client.addEquipment('HW1569', 'Antenna HW', location, {'altitude': 53.5})

"""


class InventoryClient(GraphqlClient):

    from .api.file import (
        add_location_image,
        add_site_survey_image,
        delete_document,
        delete_site_survey_image,
        add_file,
        add_files,
    )
    from .api.location_type import (
        add_location_type,
        delete_locations_by_location_type,
        delete_location_type_with_locations,
    )
    from .api.location import (
        get_location,
        get_locations_by_external_id,
        get_location_by_external_id,
        get_location_children,
        get_location_documents,
        delete_location,
        add_location,
        edit_location,
        move_location,
    )
    from .api.equipment_type import (
        copy_equipment_type,
        delete_equipment_type_with_equipments,
        _add_equipment_type,
        add_equipment_type,
        get_or_create_equipment_type,
        _edit_equipment_type,
        edit_equipment_type,
    )
    from .api.equipment import (
        add_equipment,
        add_equipment_to_position,
        get_equipment,
        get_equipment_in_position,
        delete_equipment,
        search_for_equipments,
        delete_all_equipments,
        copy_equipment_in_position,
        copy_equipment,
        get_equipment_type_of_equipment,
        get_or_create_equipment,
        get_or_create_equipment_in_position,
    )
    from .api.link import (
        add_link,
        get_link_in_port_of_equipment,
        get_all_links_and_port_names_of_equipment,
        get_port,
    )
    from .api.service import (
        add_service,
        add_service_endpoint,
        add_service_type,
        get_service,
        delete_service_type_with_services,
    )
    from .site_survey import (
        upload_site_survey,
        get_site_surveys,
        delete_site_survey,
        build_site_survey_from_survey_response,
    )
    from .api.location_template import (
        apply_location_template_to_location,
        copy_equipment_with_all_attachments,
    )
    from .api.customer import add_customer, delete_customer, get_all_customers

    def __init__(
        self,
        email: str,
        password: str,
        tenant: str = "fb-test",
        is_local_host: bool = False,
        is_dev_mode: bool = False,
        reporter: Reporter = DUMMY_REPORTER,
    ):
        super().__init__(email, password, tenant, is_local_host, is_dev_mode, reporter)
        _populate_location_types(self)
        _populate_equipment_types(self)
        _populate_service_types(self)
