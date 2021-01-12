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

import errno
import logging
import os
import pathlib
import subprocess
import time
from subprocess import SubprocessError
from .command_builder import get_trace_builder

_TRACE_FILE_NAME = "call_trace"
_TRACE_FILE_EXT = "pcap"
_MAX_FILESIZE = 4000  # ~ 4 MiB for a trace


class TraceManager:
    """
    TraceManager is a wrapper for tshark/tcpdump specifically for starting and
    stopping call/interface/subscriber traces.

    Only a single trace can be captured at a time.
    """
    def __init__(self, config):
        self._is_active = False # is call trace being captured
        self._proc = None
        self._trace_directory = config.get("trace_directory",
                                           "/var/opt/magma/trace")  # type: str
        # Specify southbound interfaces
        self._trace_interfaces = config.get("trace_interfaces",
                                           ["eth0"])  # type: List[str]

        # Should specify absolute path of trace filename if trace is active
        self._trace_filename = ""  # type: str

        tool_name = config.get("trace_tool", "tshark")  # type: str
        self._trace_builder = get_trace_builder(tool_name)

    def start_trace(self) -> bool:
        """Start a call trace.

        Captures all packets across the eth0 interface.

        Returns:
            True if successfully started call trace
        """
        if self._is_active:
            logging.error("Failed to start trace: Trace already active")
            return False

        # Example filename path:
        #   /var/opt/magma/trace/call_trace_1607358641.pcap
        self._trace_filename = "{0}/{1}_{2}.{3}".format(
            self._trace_directory,
            _TRACE_FILE_NAME,
            int(time.time()),
            _TRACE_FILE_EXT)

        command = self._trace_builder.build_trace_command(
            self._trace_interfaces,
            _MAX_FILESIZE,
            self._trace_filename)

        logging.info("Starting trace with tshark, command: [%s]",
                     ' '.join(command))

        self._ensure_trace_directory_exists()

        # TODO(andreilee): Handle edge case where only one instance of the
        #                  process can be running, and may have been started
        #                  by something external as well.
        try:
            self._proc = subprocess.Popen(command)
        except SubprocessError as e:
            logging.error("Failed to start trace: %s", str(e))
            return False

        self._is_active = True
        logging.info("Successfully started trace with tshark")
        return True

    def end_trace(self) -> bytes:
        """Ends call trace, if currently active.

        Returns:
            Call trace file in bytes
        """
        # If trace is active, then stop it
        if self._is_active:
            # If the process has ended, then _proc isn't None
            self._proc.poll()
            if self._proc.returncode is None:
                self._proc.terminate()

        # Read trace data into bytes
        with open(self._trace_filename, "rb") as trace_file:
            data = trace_file.read()  # type: bytes

        # Ensure the tmp trace file is deleted
        self._ensure_tmp_file_deleted()
        self._trace_filename = ""

        self._is_active = False

        # Everything cleaned up, return bytes
        return data

    def _ensure_tmp_file_deleted(self):
        """Ensure that tmp trace file is deleted.

        Uses exception handling rather than a check for file existence to avoid
        TOCTTOU bug
        """
        try:
            os.remove(self._trace_filename)
        except OSError as e:
            if e.errno != errno.ENOENT:
                logging.error("Error when deleting tmp trace file: %s", str(e))

    def _ensure_trace_directory_exists(self) -> None:
        pathlib.Path(self._trace_directory).mkdir(parents=True, exist_ok=True)
