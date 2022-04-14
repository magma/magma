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

import logging
from ipaddress import IPv4Address, IPv4Network, IPv6Address, IPv6Network
from typing import Union

IPAddress = Union[IPv4Address, IPv6Address]
IPNetwork = Union[IPv4Network, IPv6Network]


def log_error_and_raise(exception_class: type, message_template: str, *args):
    """Create an error message by rendering the message template with *args.
    The template is passed to logging.error, and the rendered message is
    passed to the exception constructor.

    Args:
        exception_class (type): [The exception to raise]
        message_template (str): [The template for the log and exception message]
        *args: Further arguments are used for rendering the message from the template

    Raises:
        exception_class: [The exception that is passed to this function]
    """
    logging.error(message_template, *args)
    rendered_message = message_template % args
    raise exception_class(rendered_message)
