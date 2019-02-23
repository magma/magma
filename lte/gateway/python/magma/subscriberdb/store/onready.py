"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio


class OnDataReady:
    """
    A thread-safe Event mixin interface for triggering a ready event
    when subscribers are added to the data store. Routines can wait on
    the _ready_ event to block until a condition is met:
        1. a subscriber is added
        2. a datastore resync is triggered
    """

    def __init__(self, loop=None):
        self.loop = loop if loop else asyncio.new_event_loop()
        self.event = asyncio.Event(loop=self.loop)

    def add_subscriber(self, _):
        self.loop.call_soon_threadsafe(self.trigger_ready)

    def resync(self, _):
        self.loop.call_soon_threadsafe(self.trigger_ready)

    def trigger_ready(self):
        if not self.event.is_set():
            self.event.set()
