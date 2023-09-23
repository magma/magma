#!/usr/bin/env python3
# Note: Python version >= 3.7 is needed to run this script
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
import argparse
import os
import subprocess  # noqa: S404 ignore security warning about subprocess
from typing import List

if os.getenv('MAGMA_ROOT'):
    MAGMA_ROOT = os.environ["MAGMA_ROOT"]
else:
    raise KeyError("'MAGMA_ROOT' needs to be set and point to the Magma root directory!")

LINT_DOCKER_PATH = os.path.join(
    MAGMA_ROOT,
    'lte/gateway/docker/python-precommit/',
)
IMAGE_NAME = 'magma/py-lint'
ORC8R_PYTHON_PATH = 'orc8r/gateway/python/magma'
LTE_PYTHON_PATH = 'lte/gateway/python/magma'
GITHUB_IMAGE_NAME = 'ghcr.io/magma/magma/python-precommit:latest'


def main() -> None:
    """Provide command-line options to format/lint Magma's Python codebase"""
    print("Magma root is " + MAGMA_ROOT)
    args = _parse_args()

    _build_docker_image()

    # If no paths are specified, default to magma services
    if args.diff:
        args.paths = _get_diff_against_master()
    if not args.paths:
        print("Please specify at least one path for format/lint!")
        return
    if args.format:
        _format_diff(args.paths)
    if args.lint:
        _run_flake8(args.paths)


def _build_docker_image():
    print("Building the py-lint docker image... This may take a minute or two")
    cmd = [
        'docker', 'build', '-t', IMAGE_NAME,
        '-f', os.path.join(LINT_DOCKER_PATH, 'Dockerfile'),
        MAGMA_ROOT,
    ]
    _run(cmd)


def _format_diff(paths: List[str]):
    for path in paths:
        # when changing any of these commands,
        # make sure to change the corresponding github action
        _run_docker_cmd(['isort', path])
        _run_add_trailing_comma(path)
        # This should be consistent with .github/workflows/python-workflow.yml
        autopep8_checks = 'W191,W291,W292,W293,W391,E131,E1,E2,E3'
        _run_docker_cmd(['autopep8', '--select', autopep8_checks, '-r', '--in-place', path])


def _run_add_trailing_comma(path: str):
    abs_path = os.path.join(os.path.abspath(MAGMA_ROOT), path)
    if os.path.isfile(abs_path):
        # TODO upgrade to --py36-plus eventually
        _run_docker_cmd(
            [
                'add-trailing-comma', '--py35-plus',
                '--exit-zero-even-if-changed', path,
            ],
        )


def _run_flake8(paths: List[str]):
    for path in paths:
        _run_docker_cmd(['flake8', '--exit-zero', path])


def _run_docker_cmd(commands: List[str]):
    volume_cmd = ['-v', os.path.abspath(MAGMA_ROOT) + ':/code']
    docker_image = IMAGE_NAME + ':latest'
    cmd_prefix = 'docker run -it -u 0'.split(' ')
    cmd = cmd_prefix + volume_cmd + [docker_image] + commands
    _run(cmd)


def _run(cmd: List[str]) -> None:
    print("Running '%s'..." % ' '.join(cmd))
    try:
        subprocess.run(cmd, check=True)  # noqa: S603
    except subprocess.CalledProcessError as err:
        print(err)
        exit(err.returncode)


def _get_diff_against_master() -> List[str]:
    files = "git diff --name-only --diff-filter=ACMRT master HEAD"
    changed_files_in_commit = subprocess.run(files.split(), capture_output=True).stdout.decode().split('\n')  # noqa: S603
    changed_py_files = []
    for item in changed_files_in_commit:
        if item.endswith('.py'):
            changed_py_files.append(item)
    print("Changed files since master: " + str(changed_py_files))
    return changed_py_files


def _parse_args() -> argparse.Namespace:
    """Parse the command line args

    Returns:
        argparse.Namespace: the created parser
    """
    parser = argparse.ArgumentParser(description='Python lint/format tool')

    parser.add_argument(
        '--format', '-f',
        action='store_true',
        help='Run formatting commands',
    )
    parser.add_argument(
        '--lint', '-l',
        action='store_true',
        help='Run flake8',
    )
    parser.add_argument(
        '--paths', '-p',
        nargs='+',
        help=(
            'Paths (relative from repo root) to run the linter/formatter '
            + 'against.'
            + 'Can specify multiple paths by running "-p path1 path2"'
        ),
    )
    parser.add_argument(
        '--diff', '-d',
        action='store_true',
        help=(
            'Run the command on all changed files against master. '
            + ' (equivalent to files given by '
            + '"git diff --name-only --diff-filter=ACMRT  master HEAD")'
        ),
    )

    return parser.parse_args()


if __name__ == '__main__':
    main()
