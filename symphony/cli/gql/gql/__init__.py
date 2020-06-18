#!/usr/bin/env python3

import warnings

from .client import Client, GraphqlDeprecationWarning, OperationException
from .gql import gql


__all__ = ["gql", "Client", "OperationException"]

warnings.filterwarnings("always", "", GraphqlDeprecationWarning)
