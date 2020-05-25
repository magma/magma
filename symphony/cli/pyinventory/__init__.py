#!/usr/bin/env python3

from distutils.version import LooseVersion
from typing import Optional, Tuple

from colorama import Fore
from gql.gql.reporter import DUMMY_REPORTER, Reporter
from pysymphony import SymphonyClient

from .api.equipment_type import (
    _populate_equipment_port_types,
    _populate_equipment_types,
)
from .api.location_type import _populate_location_types
from .api.service_type import _populate_service_types
from .common.constant import __version__
from .graphql.query.latest_python_package import LatestPythonPackageQuery


"""Pyinventory is a python package that allows for querying and modifying the
FBC platform inventory using graphql.

This module contains the client that allows to connect to inventory. The client
allows different kind of operations: querying and creating of locations, equipments,
positions and links.

Example of how to connect:
```
from pyinventory import InventoryClient
# since inventory is multi tenant system you will need to insert which
# partner you connect as
client = InventoryClient(email, password, "tenant_name")
location = client.addLocation(-1.22,2.66, ('City', 'Brooklyn'))
client.addEquipment('HW1569', 'Antenna HW', location, {'altitude': 53.5})
```
"""


class InventoryClient(SymphonyClient):

    from .api.file import add_location_image, delete_document, add_file, add_files
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
        get_locations,
    )
    from .api.equipment_type import (
        copy_equipment_type,
        delete_equipment_type_with_equipments,
        _add_equipment_type,
        add_equipment_type,
        get_or_create_equipment_type,
        _edit_equipment_type,
        edit_equipment_type,
        get_equipment_type_property_type,
        get_equipment_type_property_type_by_external_id,
        edit_equipment_type_property_type,
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
        edit_equipment,
        get_equipment_properties,
        get_equipments_by_type,
        get_equipments_by_location,
        get_equipment_by_external_id,
    )
    from .api.link import (
        add_link,
        get_link_in_port_of_equipment,
        get_all_links_and_port_names_of_equipment,
    )
    from .api.service import (
        add_service,
        add_service_endpoint,
        add_service_link,
        get_service,
    )
    from .api.service_type import (
        add_service_type,
        get_service_type,
        edit_service_type,
        delete_service_type,
        delete_service_type_with_services,
    )
    from .api.location_template import (
        apply_location_template_to_location,
        copy_equipment_with_all_attachments,
    )
    from .api.customer import add_customer, delete_customer, get_all_customers
    from .api.port_type import (
        add_equipment_port_type,
        get_equipment_port_type,
        edit_equipment_port_type,
        delete_equipment_port_type,
    )
    from .api.port import get_port, edit_port_properties, edit_link_properties
    from .api.user import (
        add_user,
        get_user,
        edit_user,
        deactivate_user,
        activate_user,
        get_users,
        get_active_users,
    )
    from .api.property_type import get_property_type_id

    def __init__(
        self,
        email: str,
        password: str,
        tenant: str = "fb-test",
        is_local_host: bool = False,
        is_dev_mode: bool = False,
        reporter: Reporter = DUMMY_REPORTER,
    ) -> None:
        """This is the class to use for working with inventory. It contains all
            the functions to query and and edit the inventory.

            The __init__ method populates the different entity types
            for faster run of operations.

            Args:
                email (str): The email of the user to connect with.
                password (str): The password of the user to connect with.
                tenant (str, optional): The tenant to connect to -
                            should be the beginning of "{}.purpleheadband.cloud"
                            The default is "fb-test" for QA environment
                is_local_host (bool, optional): Used for developers to connect to
                            local inventory. This changes the address and also
                            disable verification of ssl certificate
                is_dev_mode (bool, optional): Used for developers to connect to
                            local inventory from a container. This changes the
                            address and also disable verification of ssl
                            certificate
                reporter (object, optional): Use reporter.InventoryReporter to
                            store reports on all successful and failed mutations
                            in inventory. The default is DummyReporter that
                            discards reports

        """
        super().__init__(
            email,
            password,
            tenant,
            f"Pyinventory/{__version__}",
            is_local_host,
            is_dev_mode,
            reporter,
        )
        self._verify_version_is_not_broken()
        self.populate_types()

    def _verify_version_is_not_broken(self) -> None:
        package = self._get_latest_python_package_version()

        latest_version, latest_breaking_version = (
            package if package is not None else (None, None)
        )

        if latest_breaking_version is not None and LooseVersion(
            latest_breaking_version
        ) > LooseVersion(__version__):
            raise Exception(
                "This version of pyinventory is not supported anymore. \
                Please download and install the latest version ({})".format(
                    latest_version
                )
            )

        if latest_version is not None and LooseVersion(latest_version) > LooseVersion(
            __version__
        ):
            print(
                str(Fore.RED)
                + "A newer version of pyinventory exists ({}). \
            It is recommended to download and install it".format(
                    latest_version
                )
            )

    def _get_latest_python_package_version(self) -> Optional[Tuple[str, str]]:

        package = LatestPythonPackageQuery.execute(self)
        if package is not None:
            last_version = package.lastPythonPackage
            last_breaking_version = package.lastBreakingPythonPackage
            if last_version is not None:
                return (
                    last_version.version,
                    last_breaking_version.version
                    if last_breaking_version
                    else last_version.version,
                )
        return None

    def populate_types(self) -> None:
        _populate_location_types(self)
        _populate_equipment_types(self)
        _populate_service_types(self)
        _populate_equipment_port_types(self)
