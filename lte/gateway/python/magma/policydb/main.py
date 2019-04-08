"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import logging

from magma.common.service import MagmaService
from magma.common.streamer import StreamerClient
from .streamer_callback import PolicyDBStreamerCallback


def main():
    """ main() for subscriberdb """
    service = MagmaService('policydb')
    # Start a background thread to stream updates from the cloud
    if service.config['enable_streaming']:
        callback = PolicyDBStreamerCallback(service.loop)
        stream = StreamerClient({"policydb": callback}, service.loop)
        stream.start()
    else:
        logging.info('enable_streaming set to False. Streamer disabled!')

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
