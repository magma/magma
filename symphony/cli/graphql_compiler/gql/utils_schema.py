#!/usr/bin/env python3
import os
import json
import requests
from graphql import get_introspection_query, build_client_schema


def load_introspection_from_server(url):
    query = get_introspection_query()
    request = requests.post(url, json={'query': query})
    if request.status_code == 200:
        return request.json()['data']

    raise Exception('Query failed to run by returning code of '
        f'{request.status_code}. {query}')


def load_introspection_from_file(filename):
    with open(filename, 'r') as fin:
        return json.load(fin)


def load_schema(uri):
    introspection = load_introspection_from_file(uri) if os.path.isfile(uri) \
        else load_introspection_from_server(uri)
    return build_client_schema(introspection)
