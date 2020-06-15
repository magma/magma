#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import json
import logging
import queue
import ssl
import threading
import uuid
from base64 import b64encode
from types import TracebackType
from typing import Any, Callable, Dict, Optional, Type

import websocket


CallbackType = Callable[[str, Dict[str, Any]], None]

GQL_WS_SUBPROTOCOL = "graphql-ws"

# all the message types
GQL_CONNECTION_INIT = "connection_init"
GQL_START = "start"
GQL_STOP = "stop"
GQL_CONNECTION_TERMINATE = "connection_terminate"
GQL_CONNECTION_ERROR = "connection_error"
GQL_CONNECTION_ACK = "connection_ack"
GQL_DATA = "data"
GQL_ERROR = "error"
GQL_COMPLETE = "complete"
GQL_CONNECTION_KEEP_ALIVE = "ka"

logging.basicConfig()


class ConnectionException(Exception):
    """Exception thrown during connection errors to the GraphQL server"""


class InvalidPayloadException(Exception):
    """Exception thrown if payload recived from server is mal-formed or cannot be parsed """


class SubscriptionClient:
    def __init__(self, url: str, username: str, password: str) -> None:
        self.ws_url: str = url
        self._connection_init_done = False
        # cache of the headers for a session
        self._headers: Optional[Dict[str, str]] = None
        # map of subscriber id to a callback function
        self._subscriber_callbacks: Dict[str, CallbackType] = {}
        # our general receive queue
        self._queue: queue.Queue[Dict[str, Any]] = queue.Queue()
        # map of queues for each subscriber
        self._subscriber_queues: Dict[str, queue.Queue[Dict[str, Any]]] = {}
        self._shutdown_receiver = False
        self.connect(username, password)

    def connect(self, username: str, password: str) -> None:
        """
        Initializes a connection with the server.
        """
        user_and_pass = b64encode(f"{username}:{password}".encode("utf-8")).decode(
            "ascii"
        )
        header = {"Authorization": "Basic %s" % user_and_pass}
        self._connection = websocket.create_connection(
            self.ws_url,
            subprotocols=[GQL_WS_SUBPROTOCOL],
            header=header,
            sslopt={"cert_reqs": ssl.CERT_NONE},
        )
        # start the receiver thread
        self._recevier_thread = threading.Thread(target=self._receiver_task)
        self._recevier_thread.start()

    def __dump_queues(self) -> None:
        logging.debug("[GQL_CLIENT] => Dump of all the internal queues")
        dumps = list(map(lambda q: (q[0], q[1].queue), self._subscriber_queues.items()))
        logging.debug(f"[GQL_CLIENT] => Operation queues: \n {dumps}")

    # wait for any valid message, while ignoring GQL_CONNECTION_KEEP_ALIVE
    def _receiver_task(self) -> None:
        """the receive function of the client. Which validates response from the
        server and queues data """
        while not self._shutdown_receiver:
            self.__dump_queues()
            res = self._connection.recv()
            try:
                msg = json.loads(res)
            except json.JSONDecodeError as err:
                logging.warning(
                    f"Ignoring. Server sent invalid JSON data: {res} \n {err}"
                )

            # ignore messages which are GQL_CONNECTION_KEEP_ALIVE
            if msg["type"] != GQL_CONNECTION_KEEP_ALIVE:

                # check all GQL_DATA and GQL_COMPLETE should have 'id'.
                # Otherwise, server is sending malformed responses, error out!
                if msg["type"] in [GQL_DATA, GQL_COMPLETE] and "id" not in msg:
                    # TODO: main thread can't catch this exception; setup
                    # exception queues. but this scenario will only happen with
                    # servers having glaring holes in implementing the protocol
                    # correctly, which is rare. hence this is not very urgent
                    err = f'Protocol Violation.\nExpected "id" in {msg}, but could not find.'
                    raise InvalidPayloadException(err)

                # if the message has an id, it is meant for a particular operation
                if "id" in msg:
                    op_id = msg["id"]

                    # put it in the correct operation/subscriber queue
                    if op_id not in self._subscriber_queues:
                        self._subscriber_queues[op_id] = queue.Queue()
                    self._subscriber_queues[op_id].put(msg)

                    # if a callback fn exists with the id, call it
                    if op_id in self._subscriber_callbacks:
                        user_fn = self._subscriber_callbacks[op_id]
                        user_fn(op_id, msg)

                # if it doesn't have an id, put in the global queue
                else:
                    self._queue.put(msg)

    def _insert_subscriber(self, op_id: str, callback_fn: CallbackType) -> None:
        self._subscriber_callbacks[op_id] = callback_fn

    def _remove_subscriber(self, op_id: str) -> None:
        del self._subscriber_callbacks[op_id]

    def _create_operation_queue(self, op_id: str) -> None:
        self._subscriber_queues[op_id] = queue.Queue()

    def _remove_operation_queue(self, op_id: str) -> None:
        if op_id in self._subscriber_queues:
            del self._subscriber_queues[op_id]

    def _get_operation_result(self, op_id: str) -> Dict[str, Any]:
        return self._subscriber_queues[op_id].get()

    def _connection_init(self, headers: Optional[Dict[str, Any]] = None) -> None:
        # if we have already initialized and the passed headers are same as
        # prev headers, then do nothing and return
        if self._connection_init_done and headers == self._headers:
            return

        self._headers = headers
        # send the `connection_init` message with the payload
        payload = {"type": GQL_CONNECTION_INIT, "payload": {"headers": headers}}
        self._connection.send(json.dumps(payload))

        res = self._queue.get()

        if res["type"] == GQL_CONNECTION_ERROR:
            err = res.get("payload", "unknown error")
            raise ConnectionException(err)
        if res["type"] == GQL_CONNECTION_ACK:
            self._connection_init_done = True
            return

        err_msg = (
            "Unknown message from server, this client did not understand. "
            + "Original message: "
            + res["type"]
        )
        raise ConnectionException(err_msg)

    def _start(
        self, payload: Dict[str, Any], callback: Optional[CallbackType] = None
    ) -> str:
        """ pass a callback function only if this is a subscription """
        op_id = uuid.uuid4().hex
        frame = {"id": op_id, "type": GQL_START, "payload": payload}
        self._create_operation_queue(op_id)
        if callback:
            self._insert_subscriber(op_id, callback)
        self._connection.send(json.dumps(frame))
        return op_id

    def _stop(self, op_id: str) -> None:
        payload = {"id": op_id, "type": GQL_STOP}
        self._connection.send(json.dumps(payload))

    def query(
        self,
        query: str,
        variables: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Run a GraphQL query or mutation. The `query` argument is a GraphQL query
        string. You can pass optional variables and headers.
        PS: To run a subscription, see the `subscribe` method.
        """
        self._connection_init(headers)
        payload = {"headers": headers, "query": query, "variables": variables}
        op_id = self._start(payload)
        res = self._get_operation_result(op_id)
        self._stop(op_id)
        ack = self._get_operation_result(op_id)
        if ack["type"] != GQL_COMPLETE:
            logging.warning(f"Expected to receive complete, but received: {ack}")
        self._remove_operation_queue(op_id)
        return res

    def subscribe(
        self,
        query: str,
        variables: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, Any]] = None,
        callback: Optional[CallbackType] = None,
    ) -> str:
        """
        Run a GraphQL subscription.
        Parameters:
        query (str): the GraphQL query string
        callback (function): a callback function. This is mandatory.
        This callback function is called, everytime there is new data from the
        subscription.
        variables (dict): (optional) GraphQL variables
        headers (dict): (optional) a dictionary of headers for the session
        Returns:
        op_id (str): The operation id (a UUIDv4) for this subscription operation
        """

        # sanity check that the user passed a valid function
        if not callback and not callable(callback):
            raise TypeError(
                "the argument `callback` is mandatory and it should be a function"
            )

        self._connection_init(headers)
        payload = {"headers": headers, "query": query, "variables": variables}
        op_id = self._start(payload, callback)
        return op_id

    def stop_subscribe(self, op_id: str) -> None:
        """
        Stop a subscription. Takes an operation ID (`op_id`) and stops the
        subscription.
        """
        self._stop(op_id)
        self._remove_subscriber(op_id)
        self._remove_operation_queue(op_id)

    def close(self) -> None:
        """
        Close the connection with the server. To reconnect, use the `connect`
        method.
        """
        self._shutdown_receiver = True
        self._connection.close()
        self._recevier_thread.join()

    def __enter__(self) -> "SubscriptionClient":
        """ enter method for context manager """
        return self

    def __exit__(
        self,
        exc_type: Optional[Type[BaseException]],
        exc_value: Optional[BaseException],
        exc_traceback: Optional[TracebackType],
    ) -> None:
        """ exit method for context manager """
        self.close()
