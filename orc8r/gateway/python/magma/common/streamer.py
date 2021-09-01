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

import abc
import logging
import threading
import time
from typing import Any, List

import grpc
import snowflake
from google.protobuf import any_pb2
from magma.common import serialization_utils
from magma.common.metrics import STREAMER_RESPONSES
from magma.common.service_registry import ServiceRegistry
from magma.configuration.service_configs import get_service_config_value
from orc8r.protos.streamer_pb2 import DataUpdate, StreamRequest
from orc8r.protos.streamer_pb2_grpc import StreamerStub


class StreamerClient(threading.Thread):
    """
    StreamerClient provides an interface to communicate with the Streamer
    service in the cloud to get updates for a stream.

    The StreamerClient spawns a thread which listens to updates and
    schedules a callback in the asyncio event loop when an update
    is received from the cloud.

    If the connection to the cloud gets terminated, the StreamerClient
    would retry (TBD: with exponential backoff) to connect back to the cloud.
    """

    class Callback:

        @abc.abstractmethod
        def get_request_args(self, stream_name: str) -> Any:
            """
            This is called before every stream request to collect any extra
            arguments to send up to the cloud streamer service.

            Args:
                stream_name:
                Name of the stream that the request arg will be sent to

            Returns: A protobuf message
            """
            pass

        @abc.abstractmethod
        def process_update(
            self, stream_name: str, updates: List[DataUpdate],
            resync: bool,
        ):
            """
            Called when we get an update from the cloud. This method will
            be called in the event loop provided to the StreamerClient.

            Args:
                stream_name: Name of the stream
                updates: Array of updates
                resync: if true, the application can clear the
                    contents before applying the updates
            """
            raise NotImplementedError()

    def __init__(self, stream_callbacks, loop):
        """
        Args:
            stream_callbacks ({string: Callback}): Mapping of stream names to
            callbacks to subscribe to.
            loop: asyncio event loop to schedule the callback
        """
        threading.Thread.__init__(self)
        self._stream_callbacks = stream_callbacks
        self._loop = loop
        # Set this thread as daemon thread. We can kill this background
        # thread abruptly since we handle all updates (and database
        # transactions) in the asyncio event loop.
        self.daemon = True

        # Don't allow stream update rate faster than every 5 seconds
        self._reconnect_pause = get_service_config_value(
            'streamer', 'reconnect_sec', 60,
        )
        self._reconnect_pause = max(5, self._reconnect_pause)
        logging.info("Streamer reconnect pause: %d", self._reconnect_pause)
        self._stream_timeout = get_service_config_value(
            'streamer', 'stream_timeout', 150,
        )
        logging.info("Streamer timeout: %d", self._stream_timeout)

    def run(self):
        while True:
            try:
                channel = ServiceRegistry.get_rpc_channel(
                    'streamer', ServiceRegistry.CLOUD,
                )
                client = StreamerStub(channel)
                self.process_all_streams(client)
            except Exception as exp:  # pylint: disable=broad-except
                logging.error("Error with streamer: %s", exp)

            # If the connection is terminated, wait for a period of time
            # before connecting back to the cloud.
            # TODO: make this more intelligent (exponential backoffs, etc.)
            time.sleep(self._reconnect_pause)

    def process_all_streams(self, client):
        for stream_name, callback in self._stream_callbacks.items():
            try:
                self.process_stream_updates(client, stream_name, callback)

                STREAMER_RESPONSES.labels(result='Success').inc()
            except grpc.RpcError as err:
                logging.error(
                    "Error! Streaming from the cloud failed! [%s] %s",
                    err.code(), err.details(),
                )
                STREAMER_RESPONSES.labels(result='RpcError').inc()
            except ValueError as err:
                logging.error("Error! Streaming from cloud failed! %s", err)
                STREAMER_RESPONSES.labels(result='ValueError').inc()

    def process_stream_updates(self, client, stream_name, callback):
        extra_args = self._get_extra_args_any(callback, stream_name)
        request = StreamRequest(
            gatewayId=snowflake.snowflake(),
            stream_name=stream_name,
            extra_args=extra_args,
        )
        for update_batch in client.GetUpdates(
                request, timeout=self._stream_timeout,
        ):
            self._loop.call_soon_threadsafe(
                callback.process_update,
                stream_name,
                update_batch.updates,
                update_batch.resync,
            )

    @staticmethod
    def _get_extra_args_any(callback, stream_name):
        extra_args = callback.get_request_args(stream_name)
        if extra_args is None:
            return None
        else:
            extra_any = any_pb2.Any()
            extra_any.Pack(extra_args)
            return extra_any


def get_stream_serialize_filename(stream_name):
    return '/var/opt/magma/streams/{}'.format(stream_name)


class SerializingStreamCallback(StreamerClient.Callback):
    """
    Streamer client callback which decodes stream update as a string and writes
    it to a file, overwriting the previous contents of that file. The file
    location is defined by get_stream_serialize_filename.

    This callback will only save the newest update, with each successive update
    overwriting the previous.
    """

    def get_request_args(self, stream_name: str) -> Any:
        return None

    def process_update(self, stream_name, updates, resync):
        if not updates:
            return
        # For now, we only care about the last (newest) update
        for update in updates[:-1]:
            logging.info('Ignoring update %s', update.key)

        logging.info('Serializing stream update %s', updates[-1].key)
        filename = get_stream_serialize_filename(stream_name)
        serialization_utils.write_to_file_atomically(
            filename,
            updates[-1].value.decode(),
        )
