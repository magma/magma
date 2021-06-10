#!/usr/bin/env python3

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

from prometheus_client import start_http_server, Counter
import random
import time
import sys
import argparse

# Create a metric to track time spent and requests made.
UE_TRAFFIC = Counter('ue_traffic', 'SIM Description of counter', ['IMSI', 'direction', 'session_id', 'service_sim'])
UE_REPORTED_USAGE = Counter('ue_reported_usage', 'SIM Description of counter', ['IMSI', 'direction', 'session_id', 'service_sim'])
UE_DROPPED_USAGE = Counter('ue_dropped_usage', 'SIM Description of counter', ['IMSI', 'direction', 'session_id', 'service_sim'])


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("port", help="metrics port", type=int)
    parser.add_argument("prefix", help="imsi prefix")
    parser.add_argument("total", help="total number of subscribers", type=int)
    args = parser.parse_args()

    # Start up the server to expose the metrics.
    # Generate some requests.
    start_http_server(args.port)

    while True:
        for i in range(args.total):
            imsi = "IMSI%s%d" % (args.prefix, i)
            for d in ['up', 'down']:
                UE_TRAFFIC.labels(IMSI=imsi, direction=d, service_sim='test', session_id='session_test').inc()
                UE_REPORTED_USAGE.labels(IMSI=imsi, direction=d, service_sim='test', session_id='session_test').inc()
                UE_DROPPED_USAGE.labels(IMSI=imsi, direction=d, service_sim='test', session_id='session_test').inc()
        time.sleep(1)