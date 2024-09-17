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


class MsgCounterHandler(logging.Handler):
    """ Register this handler to logging to count the logs by level """

    count_by_level = None

    def __init__(self, *args, **kwargs):
        super(MsgCounterHandler, self).__init__(*args, **kwargs)
        self.count_by_level = {}

    def emit(self, record: logging.LogRecord):
        level = record.levelname
        if (level not in self.count_by_level):
            self.count_by_level[level] = 0
        self.count_by_level[level] += 1

    def pop_error_count(self) -> int:
        error_count = self.count_by_level.get('ERROR', 0)
        self.count_by_level['ERROR'] = 0
        return error_count
