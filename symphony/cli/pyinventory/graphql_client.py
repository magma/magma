#!/usr/bin/env python3
# pyre-strict

import re
from distutils.version import LooseVersion
from typing import Any, Dict, Optional, Tuple

import requests
from colorama import Fore
from gql.gql import gql
from gql.gql.client import Client
from graphql.language.ast import DocumentNode

from .consts import (
    DUMMY_REPORTER,
    INVENTORY_ENDPOINT,
    INVENTORY_GRAPHQL_ENDPOINT,
    INVENTORY_LOGIN_ENDPOINT,
    INVENTORY_STORE_DELETE_ENDPOINT,
    INVENTORY_STORE_PUT_ENDPOINT,
    LOCALHOST_INVENTORY_ENDPOINT,
    __version__,
)
from .graphql.latest_python_package_query import LatestPythonPackageQuery
from .reporter import Reporter
from .session import RequestsHTTPSessionTransport


class GraphqlClient:

    locationTypes: Dict[str, Any] = {}
    equipmentTypes: Dict[str, Any] = {}
    serviceTypes: Dict[str, Any] = {}

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

            The __init__ method uses the credentials to establish session with
            the inventory website. It also consumes graphql schema for
            validations, and populate the location types and equipment types
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

        self.email = email
        self.password = password
        self.tenant = tenant
        self.reporter = reporter
        self.address: str = (
            LOCALHOST_INVENTORY_ENDPOINT.format(tenant)
            if is_local_host
            else INVENTORY_ENDPOINT.format(tenant)
        )
        self.login: str = self.address + INVENTORY_LOGIN_ENDPOINT
        self.endpoint: str = self.address + INVENTORY_GRAPHQL_ENDPOINT
        self.put_endpoint: str = self.address + INVENTORY_STORE_PUT_ENDPOINT
        self.delete_endpoint: str = self.address + INVENTORY_STORE_DELETE_ENDPOINT
        self.session = requests.Session()
        self.session.verify = not is_local_host and not is_dev_mode
        self.is_dev_mode = is_dev_mode
        if is_local_host or is_dev_mode:
            import urllib3

            urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
        self._login()
        self.client = Client(
            transport=RequestsHTTPSessionTransport(
                self.session,
                self.endpoint,
                headers={
                    "Accept": "application/json",
                    "Content-Type": "application/json",
                },
            ),
            fetch_schema_from_transport=True,
        )

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

    def _login(self) -> None:
        response = self.session.get(self.login)
        match = re.search(b'"csrfToken":"([^"]+)"', response.content)
        assert match is not None, "Problem with inventory login"
        csrf_token = match.group(1).decode("ascii")
        login_data = "_csrf={0}&email={1}&password={2}".format(
            csrf_token, self.email, self.password
        ).encode("ascii")
        response = self.session.post(
            self.login,
            data=login_data,
            headers={"Content-type": "application/x-www-form-urlencoded"},
        )
        response.raise_for_status()
        assert (
            re.search('"email":"{}"'.format(self.email).encode(), response.content)
            is not None
        ), "Credentials are incorrect"
        self.session.headers.update(
            {"x-csrf-token": csrf_token, "User-Agent": "Pyinventory/" + __version__}
        )

    def _get_latest_python_package_version(self) -> Optional[Tuple[str, str]]:

        package = LatestPythonPackageQuery.execute(self).latestPythonPackage
        if package is not None:
            return (
                package.lastPythonPackage.version,
                package.lastBreakingPythonPackage.version,
            )
        return None

    def call(self, query: str, variables: Dict[str, Any]) -> str:
        return self.client.execute(
            gql(query), variable_values=variables, return_json=False
        )

    def query(
        self, query_name: str, query: DocumentNode, variables: Dict[str, Any]
    ) -> Dict[str, Any]:
        return self.client.execute(query, variable_values=variables)[query_name]
