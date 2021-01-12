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

from abc import ABC, abstractmethod
from typing import List


class TraceBuildException(Exception):
    pass


class TraceBuilder(ABC):
    """Builds commands to run call tracing.

    TODO(andreilee): Support tracing for subscriber
    TODO(andreilee): Support tracing by 3gpp protocol
    """
    def __init__(self):
        super().__init__()

    @abstractmethod
    def build_trace_command(
            self,
            interfaces: List[str],
            max_filesize: int,
            output_filename: str,
    ) -> List[str]:
        """Builds command to start a trace. Does not execute the command.

        Builds a command to start a trace using the specified arguments.
        It is not guaranteed that the specific implementation of TraceBuilder
        will be able to satisfy all the constraints placed by the specified
        arguments. Check the docstring of the specific implementation for
        additional details.

        Args:
            interfaces: Interfaces to capture the trace on. Specify "any" to
              capture on all interfaces.
            max_filesize: Max capture filesize specified in KiB. Capture will
              stop after reaching the specified size. Value of -1 specifies
              no limit.
            output_filename: Where the output file will be saved.

        Returns:
            A command that can be run by subprocess.Popen to start the trace.

            Example:
            ["tshark", "-i", "eth0", "-a", "filesize:4000", "-w", "out.pcap"]

        Raises:
            TraceBuildException: Could not successfully build a command to
              start a call trace with the specified arguments
        """
        pass


class TSharkTraceBuilder(TraceBuilder):
    """Builds commands to run call tracing with tshark.

    This is the most feature complete TraceBuilder.
    """
    def __init__(self):
        super().__init__()

    def build_trace_command(
            self,
            interfaces: List[str],
            max_filesize: int,
            output_filename: str,
    ) -> List[str]:
        command = ["tshark"]

        # Specify interfaces
        if "any" in interfaces:
            command.extend(["-ni", "any"])
        else:
            for interface in interfaces:
                command.extend(["-i", interface])

        # Specify max filesize. tshark will terminate after it is reached.
        if max_filesize != -1:
            command.extend(["-a", "filesize:" + str(max_filesize)])

        # Specify output file
        command.extend(["-w", output_filename])

        return command


class TcpdumpTraceBuilder(TraceBuilder):
    """Builds commands to run call tracing with tcpdump.

    This is a less feature complete TraceBuilder, as tcpdump has less options
    than a tool such as tshark.
    """
    def __init__(self):
        super().__init__()

    def build_trace_command(
            self,
            interfaces: List[str],
            max_filesize: int,
            output_filename: str,
    ) -> List[str]:
        """Builds command to start a trace using tcpdump. Does not execute it.

        SEE PARENT CLASS FOR DETAILS.

        Details here only specify argument restrictions.

        Argument Restrictions:
            interfaces: Only one interface can be specified.
            max_filesize: Only value of -1 supported (no limit).
            output_filename:

        TODO(andreilee): Support max_filesize.
                         Add postrotate command to stop after max filesize
                         reached.
        TODO(andreilee): Support multiple interfaces.
        """
        command = ["tcpdump"]

        # Specify southbound interfaces
        if len(interfaces) == 1:
            command.extend(["-i", interfaces[0]])
        elif len(interfaces) > 1:
            raise TraceBuildException("""Cannot start trace with tcpdump with
                                         more than one interface specified""")

        # Currently unsupported.
        # While a max filesize can be specified, tcpdump will rotate to new
        # output files after the maximum is reached.
        if max_filesize != -1:
            raise TraceBuildException("""Cannot start trace with tcpdump with
                                         max filesize""")

        # Specify output file
        command.extend(["-w", self._trace_filename])

        return command


def get_trace_builder(tool_name: str) -> TraceBuilder:
    """Factory method for TraceBuilder.

    Args:
      tool_name: Name of tool for capturing a call trace.
        Only options allowed are ["tshark", "tcpdump"].

    Returns:
      The TraceBuilder implementation matching the specified tool name

    Raises:
      TraceBuildException: If an unsupported tool name is specified.
    """
    if tool_name == "tshark":
        return TSharkTraceBuilder()
    elif tool_name == "tcpdump":
        return TcpdumpTraceBuilder()
    raise TraceBuildException("Failed to create trace builder, "
                              "invalid tool name specified: {}"
                              .format(tool_name))
