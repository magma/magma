"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
                    pylint_wrapper.PYLINT_IMPORT_PROBLEM,
                ),
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
                'raise-missing-from',
                'redundant-u-string-prefix',
            ],
            show_categories=["warning", "error", "fatal"],
        )
        excluded_directories = []
        parent_path = os.path.dirname(os.path.dirname(__file__))
        print("Starting pylint ...")
        print(f"Checking all sub-directories of {parent_path}")
        directories = [
            d.name for d in os.scandir(parent_path)
            if d.is_dir() and d.name not in excluded_directories
        ]
        for directory in directories:
            print(f"  Checking {directory}")
            path = os.path.join(parent_path, directory)
            py_wrap.assertNoLintErrors(path)
