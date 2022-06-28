# Orc8r REST API v1

This directory contains v1 of the Orc8r REST API.

We do not provide any guarantees about the stability of the REST API. This will likely change in the future.

The generated portions of this directory are written by the `build.py --generate` command.

## `swagger.yml` file

This file contains the full specification of the Orc8r REST API.

The spec file is automatically generated as an aggregation of the `swagger.v1.yml` files throughout the codebase.

## Client stubs

Per-language client bindings are generated based on the `swagger.yml` file. See `orc8r/cloud/api/` for more details.
