#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import sys
from argparse import ArgumentParser, Namespace
from unittest import TestLoader, TestSuite, TextTestRunner

from xmlrunner import XMLTestRunner

from . import pyinventory_tests, utils
from .utils.constant import TESTS_PATTERN, XML_OUTPUT_DIRECTORY, TestMode


if __name__ == "__main__":

    parser = ArgumentParser(description="test")
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
    # TODO(T64902729): Restore after support for cleaning production between tests
    # parser.add_argument(
    #     "-r",
    #     "--remote",
    #     help="Run against which tenant in production staging environment",
    #     type=str,
    #     default=None,
    # )

    args: Namespace = parser.parse_args()

    if args.local is not None:
        utils.TEST_MODE = TestMode.LOCAL
        utils.TENANT = args.local
    # TODO(T64902729): Restore after support for cleaning production between tests
    # elif args.remote is not None:
    #     utils.TEST_MODE = TestMode.REMOTE
    #     utils.TENANT = args.remote

    loader = TestLoader()
    loader.testNamePatterns = [args.pattern]
    suite: TestSuite = loader.loadTestsFromModule(pyinventory_tests)

    if args.output:
        runner: TextTestRunner = XMLTestRunner(output=args.output, verbosity=2)
    else:
        runner: TextTestRunner = TextTestRunner(verbosity=2)
    result = runner.run(suite)
    if len(result.errors) != 0 or len(result.failures) != 0:
        sys.exit(1)
    sys.exit(0)
