#!/usr/bin/env python3
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
    raise Exception("'MAGMA_ROOT' needs to be set and point to the Magma root directory!")

LINT_DOCKER_PATH = os.path.join(
    MAGMA_ROOT,
    'lte/gateway/docker/python-precommit/',
)
IMAGE_NAME = 'magma/py-lint'
ORC8R_PYTHON_PATH = 'orc8r/gateway/python/magma'
LTE_PYTHON_PATH = 'lte/gateway/python/magma'
GITHUB_IMAGE_NAME = 'ghcr.io/magma/magma/python-precommit:sha-b22d512'


def main() -> None:
    """Provide command-line options to format/lint Magma's Python codebase"""
    print("Magma root is " + MAGMA_ROOT)
    args = _parse_args()

    if args.local:
        print("If you are using the local Dockerfile, please make sure to run `precommit.py --build` to build the image!")

    if args.build_image:
        _build_docker_image()
        return

    # If no paths are specified, default to magma services
    if args.diff:
        args.paths = _get_diff_against_master()
    if not args.paths:
        print("Please specify at least one path for format/lint!")
        return
    if args.format:
        _format_diff(args.paths, args.local)
    if args.lint:
        _run_flake8(args.paths, args.local)


def _build_docker_image():
    print("Building the py-lint docker image... This may take a minute or two")
    cmd = [
        'docker', 'build', '-t', IMAGE_NAME,
        '-f', os.path.join(LINT_DOCKER_PATH, 'Dockerfile'),
        MAGMA_ROOT,
    ]
    _run(cmd)


def _format_diff(paths: List[str], use_local_image: bool):
    for path in paths:
        # when changing any of these commands,
        # make sure to change the corresponding github action
        _run_docker_cmd(['isort', path], use_local_image)
        _run_add_trailing_comma(path, use_local_image)
        # This should be consistent with .github/workflows/python-workflow.yml
        autopep8_checks = 'W191,W291,W292,W293,W391,E131,E1,E2,E3'
        _run_docker_cmd(['autopep8', '--select', autopep8_checks, '-r', '--in-place', path], use_local_image)


def _run_add_trailing_comma(path: str, use_local_image: bool):
    abs_path = os.path.join(os.path.abspath(MAGMA_ROOT), path)
    if os.path.isfile(abs_path):
        # TODO upgrade to --py36-plus eventually
        _run_docker_cmd(
            [
                'add-trailing-comma', '--py35-plus',
                '--exit-zero-even-if-changed', path,
            ], use_local_image,
        )


def _run_flake8(paths: List[str], use_local_image: bool):
    for path in paths:
        _run_docker_cmd(['flake8', '--exit-zero', path], use_local_image)


def _run_docker_cmd(commands: List[str], use_local_image: bool):
    volume_cmd = ['-v', os.path.abspath(MAGMA_ROOT) + ':/code']
    if use_local_image:
        docker_image = IMAGE_NAME + ':latest'
    else:
        docker_image = GITHUB_IMAGE_NAME
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
        '--local',
        action='store_true',
        help='Build the base image from local Dockerfile and use it',
    )
    parser.add_argument(
        '--build_image', '-b',
        action='store_true',
        help='Build the linting Docker image',
    )
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
