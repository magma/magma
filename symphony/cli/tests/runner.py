#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import argparse
import sys
from unittest import TestLoader, TextTestRunner

import pyinventory_tests
from pyinventory_tests.utils.constant import (
    TESTS_PATTERN,
    XML_OUTPUT_DIRECTORY,
    TestMode,
)
from xmlrunner import XMLTestRunner


if __name__ == "__main__":

    parser = argparse.ArgumentParser(description="test")
    parser.add_argument(
        "-p",
        "--pattern",
        help="Filter of tests to run. Example: '*TestLocation*'",
        default=TESTS_PATTERN,
    )
    parser.add_argument(
        "-o",
        "--output",
        help="Output path where to store xml result files",
        default=XML_OUTPUT_DIRECTORY,
    )
    parser.add_argument(
        "-l",
        "--local",
        help="Run against which tenant in local environment. Default: fb-test",
        type=str,
        const="fb-test",
        default=None,
        nargs="?",
    )
    parser.add_argument(
        "-r",
        "--remote",
        help="Run against which tenant in production staging environment",
        type=str,
        default=None,
    )

    args = parser.parse_args()

    if args.local is not None:
        pyinventory_tests.utils.TEST_MODE = TestMode.LOCAL
        pyinventory_tests.utils.TENANT = args.local
    elif args.remote is not None:
        pyinventory_tests.utils.TEST_MODE = TestMode.REMOTE
        pyinventory_tests.utils.TENANT = args.remote

    loader = TestLoader()
    loader.testNamePatterns = [args.pattern]
    suite = loader.loadTestsFromModule(pyinventory_tests)

    if args.output:
        runner = XMLTestRunner(output=args.output, verbosity=2)
    else:
        runner = TextTestRunner(verbosity=2)
    result = runner.run(suite)
    if len(result.errors) != 0 or len(result.failures) != 0:
        sys.exit(1)
    sys.exit(0)
