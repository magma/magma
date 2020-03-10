"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import abc

from contextlib import contextmanager


class BaseStore(metaclass=abc.ABCMeta):
    """
    BaseStore class defines the interfaces that different types of
    the subscriber data store need to expose.

    A well-defined interface would allow us to create different storage
    types based on need (in-memory, on-disk, cloud, etc.) and we can
    chain them for creating hybrid models.

    Implementations of BaseStore should be thread safe.
    """

    @abc.abstractmethod
    def add_subscriber(self, subscriber_data):
        """
        Method that should add the subscriber.

        Args:
            subscriber_data - SubscriberData protobuf message
        Raises:
            DuplicateSubscriberError if the subscriber is already present
        """
        raise NotImplementedError()

    @abc.abstractmethod
    @contextmanager
    def edit_subscriber(self, subscriber_id):
        """
        Context Manager to update the subscriber data.
        Provides the subscriber data as the context, and the underlying
        store guarantees the update to be atomic (by doing the
        necessary locking).

        Args:
            subscriber_id - unique identifier for the subscriber
        Raises:
            SubscriberNotFoundError if the subscriber is not present
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def delete_subscriber(self, subscriber_id):
        """
        Method that should delete a subscriber, if present.

        Args:
            subscriber_id - unique identifier for the subscriber
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def delete_all_subscribers(self):
        """
        Method that should remove all the subscribers from the store
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def resync(self, subscribers):
        """
        Method that should resync the store with the mentioned list of
        subscribers. The resync leaves the current state of subscribers
        intact.

        Args:
            subscribers - list of subscribers to be in the store.
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def get_subscriber_data(self, subscriber_id):
        """
        Method that should return the subscriber data for the subscriber.

        Args:
            subscriber_id - unique identifier for the subscriber
        Returns:
            SubscriberData protobuf message
        Raises:
            SubscriberNotFoundError if the subscriber is not present
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def list_subscribers(self):
        """
        Method that should return the list of subscribers stored

        Returns:
            List of subscriber ids
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def on_ready(self):
        """
        Awaitable interface to block until datastore is
        ready.

        Returns:
            Awaitable
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def add_apn_config(self, apn_config):
        """
        Method that should add the APN configuration.

        Args:
            apn_config - APNConfiguration protobuf message
        Raises:
            DuplicateApnError if the APN is already present
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def edit_apn_config(self, apn_config):
        """
        Method to update the APN configuration.

        Args:
            apn_config - APNConfiguration protobuf message
        Raises:
             ApnNotFoundError if the APN is not present
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def delete_apn_config(self, apn_config):
        """
        Method that should delete an APN, if present.

        Args:
            apn_config - APNConfiguration protobuf message
        """
        raise NotImplementedError()


class SubscriberNotFoundError(Exception):
    """
    Exception thrown when a subscriber is not present in the store
    when a query is requested for that subscriber
    """
    pass


class DuplicateSubscriberError(Exception):
    """
    Exception thrown when a subscriber is requested to be added to the store,
    and the subscriber is already present. The application can choose
    to delete the old subscriber and add, or declare an error.
    """
    pass


class DuplicateApnError(Exception):
    """
    Exception thrown when APN is requested to be added to the store,
    and the APN is already present.
    """

    pass


class ApnNotFoundError(Exception):
    """
    Exception thrown when APN is not present in the store
    when a query is requested for that APN
    """

    pass
