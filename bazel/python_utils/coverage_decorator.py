# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Inspired by https://github.com/bazelbuild/bazel/issues/10660#issuecomment-924223507 by @joshua-cannon-techlabs
# Contributed to https://github.com/bazelbuild/bazel under Apache-2.0

"""
Decorator used for generating test coverage with Bazel.
"""

import contextlib
import os
import pathlib
from typing import List

import coverage
from coverage_lcov import converter

# The following environment variables are created by bazel during 'bazel coverage' calls.
# Generally we should rely on the values set by bazel and do not change them.

# environment variable COVERAGE - should be 1 iff bazel coverage was used
COVERAGE_ENV = "COVERAGE"
# environment variable COVERAGE_DIR - bazel sandbox folder where coverage of traget is created
COVERAGE_DIR_ENV = "COVERAGE_DIR"
# environment variable COVERAGE_MANIFEST - bazel instrumented files, created via --instrumentation_filter
COVERAGE_MANIFEST_ENV = "COVERAGE_MANIFEST"
# environment variable COVERAGE_OUTPUT_FILE - location of generated lcov file
COVERAGE_OUTPUT_FILE_ENV = "COVERAGE_OUTPUT_FILE"

# default name of coverage data file for target
# note: bazel expects this file name - do not change
COVERAGE_DEFAULT_FILE = ".coverage"


@contextlib.contextmanager
def coverage_decorator():
    """
    Create lcov coverage generation via Bazel for the coverage verb. Various environment variables
    are used that are created by Bazel during runtime - see documentation of global variables.
    Coverage is calculated for a single target (usually one Python test file) in .coverage.
    After bazel coverage is finished, the coverage files are merged to the Bazel
    standard coverage folder
        $MAGMA_ROOT/bazel-out/_coverage/_coverage_report.dat
    Note: this decorator should be removed when Bazel is able to handle coverage for Python
        see: https://github.com/bazelbuild/rules_python/issues/43

    Yields:
        The test that is executed.
    """

    if os.getenv(COVERAGE_ENV, None) == "1":
        coverage_file = _get_coverage_file()
        cov = _create_coverage_context(coverage_file)
        cov.start()
    try:
        yield
    finally:
        if os.getenv(COVERAGE_ENV, None) == "1":
            cov.stop()
            cov.save()
            lcov_converter = _create_lcov_converter(coverage_file)
            lcov_converter.create_lcov(os.getenv(COVERAGE_OUTPUT_FILE_ENV))


def _get_coverage_file() -> pathlib.Path:
    coverage_dir = pathlib.Path(os.getenv(COVERAGE_DIR_ENV))
    return coverage_dir / COVERAGE_DEFAULT_FILE


def _get_coverage_sources() -> List[str]:
    coverage_manifest = pathlib.Path(os.getenv(COVERAGE_MANIFEST_ENV))
    return coverage_manifest.read_text().splitlines()


def _create_coverage_context(coverage_file: str) -> coverage.Coverage:
    coverage_sources = _get_coverage_sources()
    return coverage.Coverage(data_file=str(coverage_file), include=coverage_sources)


def _create_lcov_converter(coverage_file: str) -> converter.Converter:
    return converter.Converter(
        relative_path=True,
        config_file=False,
        data_file_path=str(coverage_file),
    )
