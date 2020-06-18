#!/usr/bin/env python3
import glob
import json
import os

import requests
from graphql import build_ast_schema, build_client_schema, get_introspection_query
from graphql.language.parser import parse

from .constant import FRAGMENT_DIRNAME


def load_introspection_from_server(url):
    query = get_introspection_query()
    request = requests.post(url, json={"query": query})
    if request.status_code == 200:
        return request.json()["data"]

    raise Exception(
        f"Request failure with {request.status_code} status code for query {query}"
    )


def load_introspection_from_file(filename):
    with open(filename, "r") as fin:
        return json.load(fin)


def load_schema(uri):
    introspection = (
        load_introspection_from_file(uri)
        if os.path.isfile(uri)
        else load_introspection_from_server(uri)
    )
    return build_client_schema(introspection)


def compile_schema_library(schema_library: str):
    full_schema = ""
    # use the following line to use .graphqls files as well
    # os.path.join(schema_library, "**/*.graphql*"), recursive=True
    schema_filepaths = glob.glob(
        os.path.join(schema_library, "**/*.graphql"), recursive=True
    )
    for schema_filepath in schema_filepaths:
        with open(schema_filepath) as schema_file:
            full_schema = full_schema + schema_file.read()
    return build_ast_schema(parse(full_schema))


def read_fragment_queries(graphql_library: str):
    full_fragments = {}
    fragment_filenames = glob.glob(
        os.path.join(graphql_library, f"**/{FRAGMENT_DIRNAME}/*.graphql"),
        recursive=True,
    )
    for fragment_filepath in fragment_filenames:
        with open(fragment_filepath) as fragment_file:
            full_fragments.update({fragment_filepath: fragment_file.read()})
    return full_fragments
