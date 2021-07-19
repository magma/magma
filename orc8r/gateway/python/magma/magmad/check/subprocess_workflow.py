"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Util module for the general network health check workflow of executing process
calls (e.g. `ping`, `traceroute`) via either `subprocess` or `asyncio` and
parsing the results.
"""

import asyncio
import subprocess


def exec_and_parse_subprocesses(params, arg_list_func, result_parser_func):
    """
    Execute processes via `subprocess`, block while waiting for results, then
    parse the output with the supplied parsing callback.

    Use exec_and_parse_subprocesses_async for a non-blocking async workflow.

    Args:
        params (iterable):
            params for the process calls as an iterable
        arg_list_func (param -> [str]):
            Factory function to convert a param into a list of args for the
            process call
        result_parser_func (stdout (str) -> stderr (str) -> param -> result):
            Callback to invoke to parse the output of the command. The return
            values of these callback invocations will be the return value of
            this function.

    Returns:
        [result]:
            A collection of results parsed from the output of the
            executed subprocesses
    """
    subprocs = [
        subprocess.Popen(
            arg_list_func(param),
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        ) for param in params
    ]
    outputs = [subproc.communicate() for subproc in subprocs]
    return _parse_results(params, outputs, result_parser_func)


@asyncio.coroutine
def exec_and_parse_subprocesses_async(
    params, arg_list_func, result_parser_func,
    loop=None,
):
    """
    Asynchronously execute and parse results from subprocesses. NOTE: This
    workflow can only be used from the main thread!

    Use exec_and_parse_subprocesses if you are not running in the main thread
    or if you need a non-async blocking workflow.

    See exec_and_parse_subprocess for additional args and return documentation.

    Args:
        loop: event loop to execute within (optional)
    """
    loop = loop or asyncio.get_event_loop()
    futures = [
        asyncio.create_subprocess_exec(
            *arg_list_func(param),
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        ) for param in params
    ]
    subprocs = yield from asyncio.gather(*futures)
    outputs = yield from asyncio.gather(
        *[subproc.communicate() for subproc in subprocs],
    )
    return _parse_results(params, outputs, result_parser_func)


def _parse_results(params, outputs, result_parser_func):
    param_output_zip = zip(params, outputs)
    return map(
        lambda param_and_outputs: result_parser_func(
            param_and_outputs[1][0], param_and_outputs[1][1],
            param_and_outputs[0],
        ),
        param_output_zip,
    )
