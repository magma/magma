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
import threading
import time
from collections import namedtuple
from subprocess import SubprocessError
from typing import List

import grpc
from magma.ctraced.command_builder import get_trace_builder
from orc8r.protos.ctraced_pb2 import ReportEndedTraceRequest
from orc8r.protos.ctraced_pb2_grpc import CallTraceControllerStub

_TRACE_FILE_NAME = "call_trace"
_TRACE_FILE_NAME_POSTPROCESSED = "call_trace_postprocessed"
_TRACE_FILE_EXT = "pcapng"
_TRACE_FILE_WRITE_TIMEOUT = 20  # 20 seconds for TShark to write a trace to disk
_MAX_FILESIZE = 4000  # ~ 4 MiB for a trace
_POSTPROCESSING_TIMEOUT = 10  # 10 seconds for TShark to apply display filters

EndTraceResult = namedtuple('EndTraceResult', ['success', 'data'])


class TraceManager:
    """
    TraceManager is a wrapper for TShark specifically for starting and
    stopping call/interface/subscriber traces.

    Only a single trace can be captured at a time.
    """

    def __init__(self, config, ctraced_stub: CallTraceControllerStub):
        self._trace_id = ""
        self._is_active = False  # is call trace being captured
        self._is_stopping_trace = False  # is manual stop initiated
        self._proc = None
        self._trace_directory = config.get(
            "trace_directory",
            "/var/opt/magma/trace",
        )  # type: str
        # Specify southbound interfaces
        self._trace_interfaces = config.get(
            "trace_interfaces",
            ["eth0"],
        )  # type: List[str]

        # Should specify absolute path of trace filename if trace is active
        self._trace_filename = ""  # type: str
        self._trace_filename_postprocessed = ""  # type: str

        # TShark display filters are saved to postprocess packet capture files
        self._display_filters = ""  # type: str

        self._tool_name = config.get("trace_tool", "tshark")  # type: str
        self._trace_builder = get_trace_builder(self._tool_name)

        self._ctraced_stub = ctraced_stub

    def start_trace(
        self,
        trace_id: str,
        timeout: int,
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
            logging.error(
                "TraceManager: Failed to start trace: "
                "Trace already active",
            )
            return False

        self._trace_id = trace_id
        self._build_trace_filename()

        command = self._trace_builder.build_trace_command(
            self._trace_interfaces,
            _MAX_FILESIZE,
            timeout,
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
            self._is_stopping_trace = True
            stopped = self._stop_trace()
            if not stopped:
                return EndTraceResult(False, None)
            self._wait_until_trace_file_exists()

        # Perform postprocessing of capture file with TShark display filters
        succeeded = self._conditionally_postprocess_trace()
        if not succeeded:
            return EndTraceResult(False, None)

        data = self._get_final_trace_data()  # type: bytes

        self._cleanup_trace()
        logging.info("TraceManager: Call trace has ended")

        # Everything cleaned up, return bytes
        return EndTraceResult(True, data)

    def _execute_start_trace_command(self, command: List[str]) -> bool:
        """Executes a command to start a call trace

        Args:
            command: Shell command with each token ordered in the list.
              example: ["tshark", "-i", "eth0"] would be for "tshark -i eth0"

        Returns:
            True if successfully executed command
        """
        logging.info(
            "TraceManager: Starting trace with %s, command: [%s]",
            self._tool_name, ' '.join(command),
        )

        self._ensure_trace_directory_exists()

        # TODO(andreilee): Handle edge case where only one instance of the
        #                  process can be running, and may have been started
        #                  by something external as well.

        # TODO(andreilee): Make sure that a fast failure is detected

        def run_tracing_in_thread(on_exit, command: List[str]):
            try:
                self._proc = subprocess.Popen(
                    command,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.STDOUT,
                )
            except SubprocessError as e:
                logging.error(
                    "TraceManager: Failed to start trace: %s",
                    str(e),
                )

            self._is_active = True
            logging.info(
                "TraceManager: Successfully started trace with %s",
                self._tool_name,
            )

            self._proc.wait()
            on_exit()
            return

        thread = threading.Thread(
            target=run_tracing_in_thread,
            args=(self._on_trace_exit, command),
        )
        thread.start()
        return True

    def _on_trace_exit(self):
        if self._is_stopping_trace:
            # Manual stop initiated, skip this automatic process
            return
        logging.info("TraceManager: Call trace timed out")
        if self._proc is None:
            logging.error("TraceManager: Call trace not running")
            self._report_trace_failure()
            return
        self._proc.poll()
        return_code = self._proc.returncode
        logging.debug(
            "TraceManager: Tracing process return code: %s",
            return_code,
        )
        self._dump_trace_logs()

        if return_code != 0:
            self._report_trace_failure()
            return

        self._wait_until_trace_file_exists()

        # Perform postprocessing of capture file with TShark display filters
        succeeded = self._conditionally_postprocess_trace()
        if not succeeded:
            self._report_trace_failure()
            return

        data = self._get_final_trace_data()  # type: bytes

        self._cleanup_trace()
        logging.info("TraceManager: Reporting call trace timeout")

        self._report_trace_success(data)

    def _report_trace_success(self, data: bytes):
        try:
            req = ReportEndedTraceRequest(
                trace_id=self._trace_id,
                success=True,
                trace_content=data,
            )
            self._ctraced_stub.ReportEndedCallTrace(req)
        except grpc.RpcError:
            logging.error(
                'Unable to report successful call trace for %s. ',
                self._trace_id,
            )

    def _report_trace_failure(self):
        try:
            req = ReportEndedTraceRequest(
                trace_id=self._trace_id,
                success=False,
            )
            self._ctraced_stub.ReportEndedCallTrace(req)
        except grpc.RpcError:
            logging.error(
                'Unable to report failed call trace for %s. ',
                self._trace_id,
            )

    def _stop_trace(self) -> bool:
        # If the process has ended, then _proc isn't None
        self._proc.poll()
        return_code = self._proc.returncode
        if self._proc.returncode is None:
            logging.info("TraceManager: Ending call trace")
            self._proc.terminate()
        else:
            logging.debug(
                "TraceManager: Tracing process return code: %s",
                return_code,
            )

        self._dump_trace_logs()

        if return_code is not None:
            self._is_active = False
            return False

        return True

    def _conditionally_postprocess_trace(self) -> bool:
        """ Postprocess trace. Return whether succeeded. """
        if len(self._display_filters) == 0:
            return True

        command = self._trace_builder.build_postprocess_command(
            self._trace_filename,
            self._display_filters,
            self._trace_filename_postprocessed,
        )
        logging.info(
            "TraceManager: Starting postprocess, command: [%s]",
            ' '.join(command),
        )
        try:
            self._proc = subprocess.Popen(
                command,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
            )
        except subprocess.CalledProcessError as e:
            self._is_active = False
            logging.error(
                "TraceManager: Failed to postprocess trace: %s",
                str(e),
            )
            return False

        self._proc.wait()
        if not self._proc.stdout:
            logging.error("TraceManager: Failed to capture STDOUT of tshark")
            return True
        logging.debug("<" * 25)
        while True:
            line = self._proc.stdout.readline()
            if not line:
                break
            logging.debug("| %s", str(line.rstrip()))
        logging.debug(">" * 25)
        logging.info("TraceManager: Finished postprocess")
        return True

    def _get_final_trace_data(self) -> bytes:
        if self._should_postprocess():
            filename = self._trace_filename_postprocessed
        else:
            filename = self._trace_filename

        with open(filename, "rb") as trace_file:
            data = trace_file.read()  # type: bytes
        return data

    def _cleanup_trace(self):
        self._ensure_tmp_file_deleted(self._trace_filename)
        self._ensure_tmp_file_deleted(self._trace_filename_postprocessed)
        self._is_stopping_trace = False
        self._trace_filename = ""
        self._is_active = False

    def _build_trace_filename(self):
        # Example filename path:
        #   /var/opt/magma/trace/call_trace_1607358641.pcap
        self._trace_filename = "{0}/{1}_{2}.{3}".format(
            self._trace_directory,
            _TRACE_FILE_NAME,
            int(time.time()),
            _TRACE_FILE_EXT,
        )

        self._trace_filename_postprocessed = "{0}/{1}_{2}.{3}".format(
            self._trace_directory,
            _TRACE_FILE_NAME_POSTPROCESSED,
            int(time.time()),
            _TRACE_FILE_EXT,
        )

    def _should_postprocess(self) -> bool:
        return len(self._display_filters) > 0

    def _dump_trace_logs(self):
        if not self._proc.stdout:
            logging.error("TraceManager: Failed to capture trace logs")
            return
        logging.debug("TraceManager: Trace logs:")
        logging.debug("<" * 25)
        while True:
            line = self._proc.stdout.readline()
            if not line:
                break
            logging.debug("| %s", str(line.rstrip()))
        logging.debug(">" * 25)

    def _wait_until_trace_file_exists(self):
        time_pending = 0
        while (
            not self._does_trace_file_exist()
            and time_pending < _TRACE_FILE_WRITE_TIMEOUT
        ):
            logging.debug(
                "TraceManager: Waiting 1s for trace file to be "
                "written...",
            )
            time.sleep(1)
            time_pending += 1
        if time_pending >= _TRACE_FILE_WRITE_TIMEOUT:
            self._report_trace_failure()
            return
        logging.debug("TraceManager: Trace file written!")

    def _does_trace_file_exist(self) -> bool:
        return os.path.isfile(self._trace_filename)

    def _ensure_tmp_file_deleted(self, filename: str):
        """Ensure that tmp trace file is deleted.

        Uses exception handling rather than a check for file existence to avoid
        TOCTTOU bug
        """
        try:
            if os.path.isfile(filename):
                os.remove(filename)
        except OSError as e:
            if e.errno != errno.ENOENT:
                logging.error(
                    "TraceManager: Error when deleting tmp trace "
                    "file: %s", str(e),
                )

    def _ensure_trace_directory_exists(self) -> None:
        pathlib.Path(self._trace_directory).mkdir(parents=True, exist_ok=True)
