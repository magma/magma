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
from collections import namedtuple
from subprocess import SubprocessError
from typing import List

from .command_builder import get_trace_builder

_TRACE_FILE_NAME = "call_trace"
_TRACE_FILE_NAME_POSTPROCESSED = "call_trace_postprocessed"
_TRACE_FILE_EXT = "pcapng"
_MAX_FILESIZE = 4000  # ~ 4 MiB for a trace
_POSTPROCESSING_TIMEOUT = 10 # 10 seconds for TShark to apply display filters

EndTraceResult = namedtuple('EndTraceResult', ['success', 'data'])


class TraceManager:
    """
    TraceManager is a wrapper for TShark specifically for starting and
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
        self._trace_filename_postprocessed = ""  # type: str

        # TShark display filters are saved to postprocess packet capture files
        self._display_filters = "" # type: str

        self._tool_name = config.get("trace_tool", "tshark")  # type: str
        self._trace_builder = get_trace_builder(self._tool_name)

    def start_trace(
        self,
        capture_filters: str,
        display_filters: str,
    ) -> bool:
        """Start a call trace.

        Note:
            The output file location is appended to the custom run options,
            matching trace_directory in ctraced.yml

        Args:
            capture_filters: Capture filters for running TShark.
              Equivalent to the -f option of TShark.
              Syntax based on BPF (Berkeley Packet Filter)
            display_filters: Display filters for running TShark.
              Equivalent to the -Y option of TShark.

        Returns:
            True if successfully started call trace
        """
        if self._is_active:
            logging.error("TraceManager: Failed to start trace: "
                          "Trace already active")
            return False

        self._build_trace_filename()

        command = self._trace_builder.build_trace_command(
            self._trace_interfaces,
            _MAX_FILESIZE,
            self._trace_filename,
            capture_filters,
        )

        self._display_filters = display_filters

        return self._execute_start_trace_command(command)

    def end_trace(self) -> EndTraceResult:
        """Ends call trace, if currently active.

        Returns:
            success: True if call trace finished without issue
            data: Call trace file in bytes
        """
        # If trace is active, then stop it
        if self._is_active:
            stopped = self._stop_trace()
            if not stopped:
                return EndTraceResult(False, None)
            while True:
                if self._ensure_trace_file_exists():
                    logging.info("TraceManager: Trace file written!")
                    break
                logging.info("TraceManager: Waiting 1s for trace file to be "
                             "written...")
                time.sleep(1)

        # Perform postprocessing of capture file with TShark display filters
        if len(self._display_filters) > 0:
            succeeded = self._postprocess_trace()
            if not succeeded:
                return EndTraceResult(False, None)

            # Read trace data into bytes
            with open(self._trace_filename_postprocessed, "rb") as trace_file:
                data = trace_file.read()  # type: bytes
        else:
            # Read trace data into bytes
            with open(self._trace_filename, "rb") as trace_file:
                data = trace_file.read()  # type: bytes

        # Ensure the tmp trace file is deleted
        self._ensure_tmp_file_deleted()
        self._trace_filename = ""

        self._is_active = False

        logging.info("TraceManager: Call trace has ended")

        # Everything cleaned up, return bytes
        return EndTraceResult(True, data)

    def _stop_trace(self) -> bool:
        # If the process has ended, then _proc isn't None
        self._proc.poll()
        return_code = self._proc.returncode
        if self._proc.returncode is None:
            logging.info("TraceManager: Ending call trace")
            self._proc.terminate()
        else:
            logging.info("TraceManager: Tracing process return code: %s",
                         return_code)

        logging.debug("TraceManager: Trace logs:")
        logging.debug("<" * 25)
        while True:
            line = self._proc.stdout.readline()
            if not line:
                break
            logging.debug("| %s", str(line.rstrip()))
        logging.debug(">" * 25)

        if return_code is not None:
            self._is_active = False
            return False

        return True

    def _postprocess_trace(self) -> bool:
        command = self._trace_builder.build_postprocess_command(
            self._trace_filename,
            self._display_filters,
            self._trace_filename_postprocessed
        )
        logging.info("TraceManager: Starting postprocess, command: [%s]",
                     ' '.join(command))
        try:
            self._proc = subprocess.Popen(
                command,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT)
        except subprocess.CalledProcessError as e:
            self._is_active = False
            logging.error("TraceManager: Failed to postprocess trace: %s",
                          str(e))
            return False

        self._proc.wait()
        logging.debug("<" * 25)
        while True:
            line = self._proc.stdout.readline()
            if not line:
                break
            logging.debug("| %s", str(line.rstrip()))
        logging.debug(">" * 25)
        logging.info("TraceManager: Finished postprocess")
        return True

    def _build_trace_filename(self):
        # Example filename path:
        #   /var/opt/magma/trace/call_trace_1607358641.pcap
        self._trace_filename = "{0}/{1}_{2}.{3}".format(
            self._trace_directory,
            _TRACE_FILE_NAME,
            int(time.time()),
            _TRACE_FILE_EXT)

        self._trace_filename_postprocessed = "{0}/{1}_{2}.{3}".format(
            self._trace_directory,
            _TRACE_FILE_NAME_POSTPROCESSED,
            int(time.time()),
            _TRACE_FILE_EXT)

    def _execute_start_trace_command(self, command: List[str]) -> bool:
        """Executes a command to start a call trace

        Args:
            command: Shell command with each token ordered in the list.
              example: ["tshark", "-i", "eth0"] would be for "tshark -i eth0"

        Returns:
            True if successfully executed command
        """
        logging.info("TraceManager: Starting trace with %s, command: [%s]",
                     self._tool_name, ' '.join(command))

        self._ensure_trace_directory_exists()

        # TODO(andreilee): Handle edge case where only one instance of the
        #                  process can be running, and may have been started
        #                  by something external as well.
        try:
            self._proc = subprocess.Popen(
                command,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT)
        except SubprocessError as e:
            logging.error("TraceManager: Failed to start trace: %s", str(e))
            return False

        self._is_active = True
        logging.info("TraceManager: Successfully started trace with %s",
                     self._tool_name)
        return True

    def _ensure_trace_file_exists(self) -> bool:
        return os.path.isfile(self._trace_filename)

    def _ensure_tmp_file_deleted(self):
        """Ensure that tmp trace file is deleted.

        Uses exception handling rather than a check for file existence to avoid
        TOCTTOU bug
        """
        try:
            os.remove(self._trace_filename)
            if os.path.isfile(self._trace_filename_postprocessed):
                os.remove(self._trace_filename_postprocessed)
        except OSError as e:
            if e.errno != errno.ENOENT:
                logging.error("TraceManager: Error when deleting tmp trace "
                              "file: %s", str(e))

    def _ensure_trace_directory_exists(self) -> None:
        pathlib.Path(self._trace_directory).mkdir(parents=True, exist_ok=True)
