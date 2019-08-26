"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
from prometheus_client import Counter
STREAMER_RESPONSES = Counter('streamer_responses',
                             'The number of responses by label',
                             ['result'])

SERVICE_ERRORS = Counter('service_errors',
                         'The number of errors logged')
