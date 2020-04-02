#!/usr/bin/env python3

import setuptools


GQL_PACKAGES = ["gql", "gql.*"]
PYHUB_PACKAGES = ["pyhub", "pyhub.*"]


setuptools.setup(
    name="pyhub",
    version="0.0.1",
    author="Facebook Inc.",
    description="Tool for managing the hub",
    packages=setuptools.find_packages(include=GQL_PACKAGES + PYHUB_PACKAGES),
    classifiers=["Programming Language :: Python :: 3.6"],
    include_package_data=True,
    install_requires=["graphql-core-next>=1.0.0", "requests>=2.22.0"],
)
