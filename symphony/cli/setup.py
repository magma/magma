#!/usr/bin/env python3

import setuptools


GQL_PACKAGES = ["gql", "gql.*"]
COMPILER_PACKAGES = ["graphql_compiler", "graphql_compiler.*"]


setuptools.setup(
    name="gqlc",
    author="Facebook Inc.",
    description="Tool for accessing and modifying FBC Platform Inventory database",
    packages=setuptools.find_packages(include=GQL_PACKAGES + COMPILER_PACKAGES),
    classifiers=["Programming Language :: Python :: 3.7"],
    include_package_data=True,
    install_requires=[
        "graphql-core-next>=1.0.0",
        "unicodecsv>=0.14.1",
        "requests>=2.22.0",
        "dataclasses==0.6",
        "dataclasses-json==0.3.2",
        "deepdiff>=4.3.2",
        "pytest>=5.4.1",
    ],
)
