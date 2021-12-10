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
        """list_subscribers - method that should return the list of subscribers stored"""
        raise NotImplementedError()

    @abc.abstractmethod
    def upsert_subscriber(self, subscriber_data):
        """
        Check if the given subscriber exists in store. If so, update subscriber
        data; otherwise, add subscriber.

        Args:
            subscriber_data: the data of the subscriber to be upserted.
        """
        raise NotImplementedError()

    @abc.abstractmethod
    async def on_ready(self):
        """
        Awaitable interface to block until datastore is
        ready.

        Returns:
            Awaitable
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


class SuciProfileNotFoundError(Exception):
    """
    Exception thrown when a suciprofile is not present in the store
    when a query is requested for that suciprofile
    """
    pass                    # noqa: WPS604


class DuplicateSuciProfileError(Exception):
    """
    Exception thrown when a suciprofile is requested to be added to the store,
    and the subscriber is already present. The application can choose
    to delete the old suciprofile and add, or declare an error.
    """
    pass                    # noqa: WPS604
