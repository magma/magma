"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import os
import unittest

from magma.tests import pylint_wrapper


class MagmaPyLintTest(unittest.TestCase):

    def test_pylint(self):
        if not pylint_wrapper.PYLINT_AVAILABLE:
            self.skipTest(
                'Pylint not available, probably because this test is running '
                'under @mode/opt: {}'.format(
                    pylint_wrapper.PYLINT_IMPORT_PROBLEM
                )
            )

        py_wrap = pylint_wrapper.PyLintWrapper(
            disable_ids=[
                'no-member',  # doesn't handle decorators correctly
                'unexpected-keyword-arg',
                # doesn't handle decorators correctly
                'no-value-for-parameter',
                # doesn't handle decorators correctly
                'fixme',  # allow todos
                'unnecessary-pass',  # triggers when pass is ok
            ],
            show_categories=["warning", "error", "fatal"],
        )

        directories = [
            'enodebd',
            # 'mobilityd',
            'pipelined',
            # 'pkt_tester',
            'policydb',
            # 'redirectd',
            'subscriberdb',
        ]
        parent_path = os.path.dirname(os.path.dirname(__file__))
        for directory in directories:
            path = os.path.join(parent_path, directory)
            py_wrap.assertNoLintErrors(path)
