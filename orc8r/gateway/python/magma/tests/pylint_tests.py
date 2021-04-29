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
            ],
            show_categories=["warning", "error", "fatal"],
        )
        # TODO look up directories in magma/orc8r/gateway/python/magma
        directories = [
            'common',
            'configuration',
            'ctraced',
            'directoryd',
            'eventd',
            'magmad',
            'state',
        ]
        parent_path = os.path.dirname(os.path.dirname(__file__))
        for directory in directories:
            path = os.path.join(parent_path, directory)
            py_wrap.assertNoLintErrors(path)
