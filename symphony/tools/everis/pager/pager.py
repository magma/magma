#!/bin/bin/env python3

import urllib3
import gql
import json
import argparse
import sys

from graphql.language.printer import print_ast
from graphql.execution import ExecutionResult
from gql.transport.requests import HTTPTransport
from requests import post
from requests.auth import HTTPBasicAuth


urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

query = gql.gql(
    """
    query Locations($after: Cursor, $first: Int) {
        locations(after: $after, first: $first) {
            edges {
                node {
                    id
                    externalId
                    name
                    locationType {
                        id
                        name
                        mapType
                        mapType
                    }
                    parentLocation {
                        id
                    }
                    latitude
                    longitude
                    equipments {
                        id
                    }
                    properties {
                        id
                        stringValue
                        propertyType {
                            name
                        }
                    }
                    siteSurveyNeeded
                }
            }
            pageInfo {
                hasNextPage
                endCursor
            }
        }
    }
"""
)


class Transport(HTTPTransport):
    def __init__(
        self, url, auth=None, use_json=False, timeout=None, verify_ssl=None, **kwargs
    ):
        super().__init__(url, **kwargs)
        self.auth = auth
        self.default_timeout = timeout
        self.use_json = use_json
        self.verify_ssl = verify_ssl

    def execute(self, document, variable_values=None, timeout=None):
        query_str = print_ast(document)
        payload = {"query": query_str, "variables": variable_values or {}}

        data_key = "json" if self.use_json else "data"
        post_args = {
            "headers": self.headers,
            "auth": self.auth,
            "cookies": self.cookies,
            "timeout": timeout or self.default_timeout,
            "verify": self.verify_ssl,
            data_key: payload,
        }
        request = post(self.url, **post_args)
        request.raise_for_status()

        result = request.json()
        assert (
            "errors" in result or "data" in result
        ), 'Received non-compatible response "{}"'.format(result)
        return ExecutionResult(errors=result.get("errors"), data=result.get("data"))


class Client(gql.Client):
    def __init__(self, url, username, password, **kwargs):
        auth = HTTPBasicAuth(username, password)
        transport = Transport(url, auth=auth, use_json=True, verify_ssl=False)
        super().__init__(
            transport=transport, fetch_schema_from_transport=True, **kwargs
        )


def paginate(client, query, step=100):
    values = {"first": step}
    aggr = rsp = client.execute(query, variable_values=values)
    while rsp["locations"]["pageInfo"]["hasNextPage"]:
        values["after"] = rsp["locations"]["pageInfo"]["endCursor"]
        rsp = client.execute(query, variable_values=values)
        aggr["locations"]["edges"].append(rsp["locations"]["edges"])
    return aggr


def main():
    parser = argparse.ArgumentParser(description="Loads locations from inventory")
    parser.add_argument(
        "-u", "--username", type=str, required=True, help="login username"
    )
    parser.add_argument(
        "-p", "--password", type=str, required=True, help="login password"
    )
    parser.add_argument(
        "-e", "--endpoint", type=str, required=True, help="graphql endpoint"
    )
    parser.add_argument("--step", type=int, help="pagination step")
    args = parser.parse_args()

    client = Client(url=args.endpoint, username=args.username, password=args.password)
    data = paginate(client, query)
    json.dump(data, sys.stdout)


if __name__ == "__main__":
    main()
