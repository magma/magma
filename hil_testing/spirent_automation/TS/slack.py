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
from typing import Dict, List, Union

import config
import requests


class SlackSender:
    BASE_SLACK_URL = "https://hooks.slack.com/services"  # + big long token

    def __init__(self):
        # Use FAKE_VALUE for debugging
        self.webhook_path = config.SLACK.get("slack_webhook_path")

    def _post_message(self, blocks: List[Dict[str, any]]):
        url = f"{self.BASE_SLACK_URL}/{self.webhook_path}"
        if self.webhook_path == "TESTING":
            print(f"## FAKE SLACK POST to {url}")
            print(blocks)
        else:
            response = requests.post(url, json={"blocks": blocks})
            response.raise_for_status()

    def send_report(
        self, test_results: Dict[str, any], sw_ver: str, ts: str, link: str = None,
    ):
        if not self.webhook_path:
            return

        def to_message_block(
            name: str, res: Dict[str, any],
        ) -> Union[None, Dict[str, any]]:
            """Render the test results"""

            if len(res) == 0:
                return None
            elif res.get("result") == "PASS":
                return None

            text = f"*{name}: *{res.get('result')} | {', '.join(map(str, res.get('fail_procs')))}"

            return {"type": "section", "text": {"type": "mrkdwn", "text": text}}

        def generate_message_blocks():
            for name, res in test_results.items():
                yield to_message_block(name, res)

        header_text = '*HIL "{}" Test suite run report on "{}".*'.format(ts, sw_ver)
        if link:
            header_text += f" <{link}|See full details>"

        message_blocks = [
            {"type": "section", "text": {"type": "mrkdwn", "text": header_text}},
        ] + list(filter(lambda b: b is not None, generate_message_blocks()))

        if len(message_blocks) >= 1:
            logging.info("Sending %d blocks to slack", len(message_blocks))
            self._post_message(message_blocks)
        else:
            logging.debug("No results to send to slack")
