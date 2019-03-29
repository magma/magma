"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.

Util module for executing multiple `dpkg` commands via subprocess.
"""

from collections import namedtuple

import asyncio
import re

from magma.magmad.check import subprocess_workflow


DpkgCommandParams = namedtuple('DpkgCommandParams', [])
DpkgCommandResult = namedtuple('DpkgCommandResult',
                               ['error', 'kernel_versions_installed'])


def get_kernel_versions():
    """
    Execute dpkg commands via subprocess. Blocks while waiting for output.

    Returns:
        [DpkgCommandResult]: stats from the executed dpkg commands
    """
    return subprocess_workflow.exec_and_parse_subprocesses(
        [DpkgCommandParams()],
        _get_dpkg_command_args_list,
        parse_dpkg_output,
    )


@asyncio.coroutine
def get_kernel_versions_async(loop=None):
    """
    Execute dpkg commands asynchronously.

    Args:
        loop: asyncio event loop (optional)

    Returns:
        [DpkgCommandResult]: stats from the executed dpkg commands
    """
    return subprocess_workflow.exec_and_parse_subprocesses_async(
        [DpkgCommandParams()],
        _get_dpkg_command_args_list,
        parse_dpkg_output,
        loop,
    )


def _get_dpkg_command_args_list(_):
    return ['dpkg', '--list']


def parse_dpkg_output(stdout, stderr, _):
    """
    Parse stdout output from a dpkg command.
    """
    if stderr:
        return DpkgCommandResult(
            kernel_versions_installed=None,
            error=str(stderr),
        )
    else:
        installed = re.findall(r'\S*linux-image\S*', str(stdout))
        return DpkgCommandResult(
            kernel_versions_installed=installed,
            error=None
        )
