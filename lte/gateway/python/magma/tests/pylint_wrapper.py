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

import warnings

PYLINT_AVAILABLE = False
PYLINT_IMPORT_PROBLEM = 'Error importing pylint'

try:
    from pylint import lint, reporters
    PYLINT_AVAILABLE = True
except (NotADirectoryError, ImportError, ModuleNotFoundError) as e:
    PYLINT_IMPORT_PROBLEM = e


class PyLintWrapper():

    def __init__(
        self, ignored_modules=None, ignored_classes=None,
        show_categories=None, disable_ids=None,
    ):
        def default(x, y):
            return x if x is not None else y

        self.ignored_modules = default(ignored_modules, [])
        self.ignored_classes = default(ignored_classes, [])
        self.show_categories = default(show_categories, [])
        self.disable_ids = default(disable_ids, ['error', 'fatal'])

    def filter(self, message):
        """
        Check if we should show this error.
        Override in child class for custom filtering.
        Return True to add message to report.
        """
        return message.category in self.show_categories and \
            message.module not in self.ignored_modules

    def run_pylint(self, path):
        with warnings.catch_warnings():
            # suppress this warnings from output
            warnings.filterwarnings(
                "ignore",
                category=PendingDeprecationWarning,
            )
            warnings.filterwarnings(
                "ignore",
                category=DeprecationWarning,
            )
            warnings.filterwarnings(
                "ignore",
                category=ImportWarning,
            )

            linter = lint.PyLinter()
            # Register standard checkers.
            linter.load_default_plugins()
            linter.set_reporter(reporters.CollectingReporter())
            # we can simply use linter.error_mode(),
            # but I prefer to filter errors later.
            linter.global_set_option("ignored-modules", self.ignored_modules)
            linter.global_set_option("ignored-classes", self.ignored_classes)
            linter.global_set_option("disable", self.disable_ids)
            linter.check(path)

            return linter.reporter.messages

    def assertNoLintErrors(self, path):
        msg = "PyLint found errors:\n"

        problems = self.run_pylint(path)
        errors_per_module = {}
        for message in problems:
            if self.filter(message):
                group = errors_per_module.setdefault(message.module, [])
                group.append(message)
        for module, errors in errors_per_module.items():
            msgs_for_module = [
                "{}: {}, {}: {} ({})".format(
                    message.msg_id, message.line, message.column,
                    message.msg, message.symbol,
                ) for message in errors
            ]
            msg += "************* Module {}\n".format(module)
            msg += "\n".join(msgs_for_module) + "\n"
        if errors_per_module:
            raise AssertionError(msg)
