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

from typing import List


class TraceBuildException(Exception):
    pass


class TraceBuilder:
    """Builds commands to run call tracing with TShark.

    TODO(andreilee): Support tracing for subscriber
    TODO(andreilee): Support tracing by 3gpp protocol
    """

    def build_trace_command(
        self,
        interfaces: List[str],
        max_filesize: int,
        timeout: int,
        output_filename: str,
        capture_filters: str,
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
            capture_filters: TShark capture filters.
              Filters are applied while packet capture is occurring.

        Returns:
            A command that can be run by subprocess.Popen to start the trace.

            Example:
            ["tshark", "-i", "eth0", "-a", "filesize:4000", "-w", "out.pcap"]

        Raises:
            TraceBuildException: Could not successfully build a command to
              start a call trace with the specified arguments
        """
        command = ["tshark"]

        # Specify interfaces
        if "any" in interfaces:
            command.extend(["-ni", "any"])
        else:
            for interface in interfaces:
                command.extend(["-i", interface])

        if len(capture_filters) > 0:
            command.extend(["-f", capture_filters])

        # Specify max filesize. tshark will terminate after it is reached.
        if max_filesize != -1:
            command.extend(["-a", "filesize:" + str(max_filesize)])

        command.extend(["-a", "duration:" + str(timeout)])

        # Specify output file
        command.extend(["-w", output_filename])

        return command

    def build_postprocess_command(
        self,
        input_filename: str,
        display_filters: str,
        output_filename: str,
    ) -> List[str]:
        """Builds command to postprocess a trace with display filters.

        Does not execute the command.

        Args:
            input_filename: Input file path
            display_filters: TShark display filters.
            output_filename: Output file path

        Returns:
            A command that can be run by subprocess.Popen to postprocess the
            trace.

            Example:
            ["tshark", "-r", "in.pcap", "-R", "displayfilter",
             "-w", "out.pcap"]
        """
        return [
            "tshark",
            "-r", input_filename,
            "-Y", display_filters,
            "-w", output_filename,
        ]


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
        return TraceBuilder()
    raise TraceBuildException(
        "Failed to create trace builder, "
        "invalid tool name specified: {}"
        .format(tool_name),
    )
