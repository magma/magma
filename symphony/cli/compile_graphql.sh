#!/usr/bin/env bash

python3 -m graphql_compiler.gql.cli ../graph/graphql/schema/symphony.graphql pyinventory/graphql/
python3 ./extract_graphql_deprecations.py ../graph/graphql/schema/symphony.graphql ../docs/md/graphql-breaking-changes.md
