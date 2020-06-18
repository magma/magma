#!/usr/bin/env python3

from graphql.language.ast import DocumentNode
from graphql.language.parser import parse
from graphql.language.source import Source


def gql(request_string: str) -> DocumentNode:
    if isinstance(request_string, str):
        source = Source(request_string, "GraphQL request")
        return parse(source)
    else:
        raise Exception('Received incompatible request "{}".'.format(request_string))
