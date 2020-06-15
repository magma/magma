#!/usr/bin/env python3

import re
from typing import Any, Dict, Optional

from gql.gql.graphql_client import GraphqlClient
from gql.gql.reporter import DUMMY_REPORTER, Reporter
from requests import Session
from requests.auth import HTTPBasicAuth
from requests.models import Response

from .common.endpoint import (
    LOCALHOST_SYMPHONY,
    SYMPHONY_GRAPHQL,
    SYMPHONY_LOGIN,
    SYMPHONY_STORE_DELETE,
    SYMPHONY_STORE_PUT,
    SYMPHONY_URI,
)


class SymphonyClient(GraphqlClient):
    def __init__(
        self,
        email: str,
        password: str,
        tenant: str = "fb-test",
        app: str = "Pyinventory",
        is_local_host: bool = False,
        is_dev_mode: bool = False,
        reporter: Reporter = DUMMY_REPORTER,
    ) -> None:
        """This is the class to use for working with symphony server.

            The __init__ method uses the credentials to establish session with
            the symphony server. It also consumes graphql schema for
            validations, and validates the client version is compatible with server.

            Args:
                email (str): The email of the user to connect with.
                password (str): The password of the user to connect with.
                tenant (str, optional): The tenant to connect to -
                            should be the beginning of "{}.purpleheadband.cloud"
                            The default is "fb-test" for QA environment
                is_local_host (bool, optional): Used for developers to connect to
                            local symphony. This changes the address and also
                            disable verification of ssl certificate
                is_dev_mode (bool, optional): Used for developers to connect to
                            local symphony from a container. This changes the
                            address and also disable verification of ssl
                            certificate
                reporter (object, optional): Use reporter.Reporter to
                            store reports on all successful and failed mutations
                            in symphony. The default is DummyReporter that
                            discards reports

        """

        self.email = email
        self.password = password
        self.address: str = (
            LOCALHOST_SYMPHONY.format(tenant)
            if is_local_host
            else SYMPHONY_URI.format(tenant)
        )
        graphql_endpoint_address = self.address + SYMPHONY_GRAPHQL
        self.app = app

        self.session: Session = Session()
        auth = HTTPBasicAuth(email, password)
        verify_ssl = not is_local_host and not is_dev_mode
        self.session.verify = verify_ssl
        if not verify_ssl:
            import urllib3

            urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

        self.put_endpoint: str = self.address + SYMPHONY_STORE_PUT
        self.delete_endpoint: str = self.address + SYMPHONY_STORE_DELETE

        super().__init__(graphql_endpoint_address, self.session, app, auth, reporter)

    def store_file(self, file_path: str, file_type: str, is_global: bool) -> str:
        # TODO(T64504906): Remove after basic auth is enabled
        if "x-csrf-token" not in self.session.headers:
            self._login()
        sign_response = self.session.get(
            self.put_endpoint,
            params={"contentType": file_type},
            headers={"Is-Global": str(is_global)},
        )
        sign_response = sign_response.json()
        signed_url = sign_response["URL"]
        with open(file_path, "rb") as f:
            file_data = f.read()
        response = self.session.put(
            signed_url, data=file_data, headers={"Content-Type": file_type}
        )
        response.raise_for_status()
        return sign_response["key"]

    def delete_file(self, key: str, is_global: bool) -> None:
        # TODO(T64504906): Remove after basic auth is enabled
        if "x-csrf-token" not in self.session.headers:
            self._login()
        sign_response = self.session.delete(
            self.delete_endpoint.format(key),
            headers={"Is-Global": str(is_global)},
            allow_redirects=False,
        )
        sign_response.raise_for_status()
        assert sign_response.status_code == 307
        signed_url = sign_response.headers["location"]
        response = self.session.delete(signed_url)
        response.raise_for_status()

    def post(self, url: str, json: Optional[Dict[str, Any]] = None) -> Response:
        # TODO(T64504906): Remove after basic auth is enabled
        if "x-csrf-token" not in self.session.headers:
            self._login()
        return self.session.post(
            "".join([self.address, url]), json=json, headers={"User-Agent": self.app}
        )

    def put(self, url: str, json: Optional[Dict[str, Any]] = None) -> Response:
        # TODO(T64504906): Remove after basic auth is enabled
        if "x-csrf-token" not in self.session.headers:
            self._login()
        return self.session.put(
            "".join([self.address, url]), json=json, headers={"User-Agent": self.app}
        )

    def _login(self) -> None:
        login_endpoint = self.address + SYMPHONY_LOGIN
        response = self.session.get(login_endpoint)
        match = re.search(b'"csrfToken":"([^"]+)"', response.content)
        assert match is not None, "Problem with symphony login"
        csrf_token = match.group(1).decode("ascii")
        login_data = "_csrf={0}&email={1}&password={2}".format(
            csrf_token, self.email, self.password
        ).encode("ascii")
        response = self.session.post(
            login_endpoint,
            data=login_data,
            headers={"Content-type": "application/x-www-form-urlencoded"},
        )
        response.raise_for_status()
        assert (
            re.search('"email":"{}"'.format(self.email).encode(), response.content)
            is not None
        ), "Credentials are incorrect"
        self.session.headers.update({"x-csrf-token": csrf_token})
