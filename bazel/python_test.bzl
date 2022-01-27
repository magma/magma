# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This file contains parts derived from [1].
# Licensed unter Apache-2.0
# Copyright 2011-2022 Software Freedom Conservancy
# Copyright 2004-2011 Selenium committers
# [1] https://github.com/SeleniumHQ/selenium/blob/1990a43a745669e87acb39331a0742979ecc728a/py/private/pytest.bzl

"""
Enhances Python test execution via Bazel by creating a wrapper on-the-fly that
executes a Python test. The wrapper is used for
* specify a test-framework - here pytest
* use a decorator for creating test coverage data if bazel coverage is used
  (see coverage_decorator.py)
"""

load("@rules_python//python:defs.bzl", "PyInfo", "py_test")

CONTENT_PYTEST_RUNNER = """
import sys
import pytest

from bazel.python_utils.coverage_decorator import coverage_decorator

if __name__ == "__main__":
    with coverage_decorator():
        args =  ["-ra", "-vv"]  + %s + sys.argv[1:] + %s
        sys.exit(pytest.main(args))"""

def _stringify(paths):
    return repr(paths)

def _pytest_runner_impl(ctx):
    if len(ctx.attr.srcs) == 0:
        fail("No test files specified.")

    expanded_args = [ctx.expand_location(arg, ctx.attr.data) for arg in ctx.attr.args]

    runner = ctx.actions.declare_file(ctx.attr.name)
    ctx.actions.write(
        runner,
        CONTENT_PYTEST_RUNNER % (_stringify(expanded_args), _stringify([src.path for src in ctx.files.srcs])),
        is_executable = True,
    )

    return [
        DefaultInfo(
            files = depset([runner]),
            runfiles = ctx.runfiles(ctx.files.data),
            executable = runner,
        ),
    ]

_pytest_runner = rule(
    _pytest_runner_impl,
    attrs = {
        "args": attr.string_list(
            default = [],
        ),
        "data": attr.label_list(
            allow_empty = True,
            allow_files = True,
        ),
        "deps": attr.label_list(
            providers = [
                PyInfo,
            ],
        ),
        "python_version": attr.string(
            values = ["PY2", "PY3"],
            default = "PY3",
        ),
        "srcs": attr.label_list(
            allow_files = [".py"],
        ),
    },
)

def pytest_test(name, srcs, deps = [], args = [], data = [], imports = [], python_version = None, **kwargs):
    runner_target = "%s-runner.py" % name

    _pytest_runner(
        name = runner_target,
        testonly = True,
        srcs = srcs,
        deps = deps,
        args = args,
        data = data,
        python_version = python_version,
    )

    py_test(
        name = name,
        python_version = python_version,
        srcs = srcs + [runner_target],
        deps = deps + ["//bazel/python_utils:coverage_decorator"],
        main = runner_target,
        legacy_create_init = False,
        imports = imports + ["."],
        **kwargs
    )
