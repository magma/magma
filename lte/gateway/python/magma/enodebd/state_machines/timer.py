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

from datetime import datetime, timedelta


class StateMachineTimer():
    def __init__(self, seconds_remaining: int) -> None:
        self.start_time = datetime.now()
        self.seconds = seconds_remaining

    def is_done(self) -> bool:
        time_elapsed = datetime.now() - self.start_time
        if time_elapsed > timedelta(seconds=self.seconds):
            return True
        return False

    def seconds_elapsed(self) -> int:
        time_elapsed = datetime.now() - self.start_time
        return int(time_elapsed.total_seconds())

    def seconds_remaining(self) -> int:
        return max(0, self.seconds - self.seconds_elapsed())
