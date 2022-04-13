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

import enum
import pickle

'''
TrafficServerInstance and TrafficTestInstance are payloads used to coordinate
trfgen testing configurations, e.g. IPs and ports.

TrafficMessage and its subclasses are the message containers that wrap around
information about the message type and its payload, as well as an optional
identifier.

The sequence of events for these messages is:


    traffic_server.py                                traffic_util.py
           <---------------TrafficRequest[TEST]--------------
           --------------TrafficResponse[SERVER]------------>
           <--------------TrafficRequest[START]--------------
           -------------TrafficResponse[STARTED]------------>
           -------------TrafficResponse[RESULTS]------------>


TrafficTest sends TrafficRequest (TEST) with a tuple(TrafficTestInstance) as
    the payload to TrafficTestServer.
TrafficTestServer creates TrafficTestDriver, which sends
    TrafficResponse(SERVER) with a tuple(TrafficServerInstance) as the payload
    to TrafficTest. This tuple is indexed to be associated with the received
    TrafficTestInstance objects.
TrafficTest sends TrafficRequest (START) after starting its downlink iperfs.
TrafficTestDriver sends TrafficResponse (STARTED) after starting all its iperf
    instances.
TrafficTest starts its uplink iperfs, and awaits a response from
    TrafficTestDriver for results.
TrafficTestDriver sends TrafficResponse (RESULTS) after tests have completed.
'''


class TrafficServerInstance(object):
    ''' Information about the server instance for a single uplink/downlink
    traffic channel '''

    def __init__(self, ip, port, mac):
        ''' Create a traffic server instance with the given values

        Args:
            ip (ipaddress.ip_address): the IP of the test server
            port (int): the port number of the test server
            mac (str): the MAC address of the test server
        '''
        self.ip = ip
        self.mac = mac
        self.port = port

    def __repr__(self):
        ''' String representation of this test server instance '''
        return ' '.join((
            '%s:' % type(self).__name__,
            '%s:%d' % (self.ip.exploded, self.port),
            'on device',
            self.mac,
        ))


class TrafficTestInstance(object):
    ''' Information about the test instance for a single uplink/downlink
    traffic channel '''

    def __init__(self, is_uplink, is_udp, duration, ip, port):
        ''' Create a traffic test instance with the given values

        Args:
            is_uplink (bool): whether the test is uplink (else downlink)
            is_udp (bool): whether the test is UDP (else TCP)
            duration (int): the duration of the test, in seconds
            ip (ipaddress.ip_address): the IP of the test device (UE)
            port (int): the port number of the test device (UE)
        '''
        self.duration = duration
        self.ip = ip
        self.is_udp = is_udp
        self.is_uplink = is_uplink
        self.port = port

    def __repr__(self):
        ''' String representation of this test instance '''
        return ' '.join((
            '%s:' % type(self).__name__,
            'UPLINK' if self.is_uplink else 'DOWNLINK',
            'UDP' if self.is_udp else 'TCP',
            'test,',
            '%d seconds' % self.duration,
            'for test device at',
            '%s:%d' % (self.ip.exploded, self.port),
        ))


class TrafficMessage(object):
    ''' Message superclass between client and server '''

    def __init__(self, message, identifier, payload):
        ''' Create a TrafficMessage of the given type with the specified
        payload

        Args:
            message: the message type
            identifier (int): a unique identifier used to refer to a single
                test driver
            payload (object): the payload to incorporate
        '''
        self.message = message
        self.id = identifier
        self.payload = payload

    def __repr__(self):
        ''' String representation of this message '''
        payload_str = repr(self.payload)
        return ' '.join((
            '%s' % type(self).__name__,
            '(%s, id %s):' % (self.message.name, str(self.id)),
            payload_str,
        ))

    @staticmethod
    def recv(stream):
        ''' Retrieve a TrafficMessage from the stream and return it

        Args:
            stream (object): a readable binary file-like object

        Returns a TrafficMessage, the message received, or None if a
            non-message was received
        '''
        assert (not stream.closed) and stream.readable

        length = stream.readline()
        if len(length) and length is not b'0':
            length = int(length.decode())
            line = stream.read(length)
            return pickle.loads(line)
        return None

    def send(self, stream):
        ''' Send this TrafficMessage through the given stream

        Args:
            stream (object): a writable binary file-like object
        '''
        assert (not stream.closed) and stream.writable

        pstr = pickle.dumps(self)
        pstrlen = str(len(pstr)).encode()
        stream.write(pstrlen + b'\n' + pstr)
        stream.flush()


# Enumerated type for TrafficRequest; module-level for pickling purposes
TrafficRequestType = enum.unique(
    enum.Enum(
        'TrafficRequestType', 'EXIT SHUTDOWN START TEST',
    ),
)


class TrafficRequest(TrafficMessage):
    ''' Request object sent from client to server '''

    def __init__(self, message, identifier=None, payload=None):
        ''' Create a TrafficRequest of the given type with the specified
        payload

        Args:
            message (TrafficRequestType): the message type
            identifier (int): a unique identifier used to refer to a single
                test driver
            payload (object): the payload to incorporate; defaults to None
        '''
        assert isinstance(message, TrafficRequestType)
        super(TrafficRequest, self).__init__(message, identifier, payload)


# Enumerated type for TrafficResponse; module-level for pickling purposes
TrafficResponseType = enum.unique(
    enum.Enum(
        'TrafficResponseType', 'INFO RESULTS SERVER STARTED',
    ),
)


class TrafficResponse(TrafficMessage):
    ''' Response object sent from server to client '''

    def __init__(self, message, identifier=None, payload=None):
        ''' Create a TrafficResponse of the given type with the specified
        payload

        Args:
            message (TrafficResponseType): the message type
            identifier (int): a unique identifier used to refer to a single
                test driver
            payload (object): the payload to incorporate; defaults to None
        '''
        assert isinstance(message, TrafficResponseType)
        super(TrafficResponse, self).__init__(message, identifier, payload)
