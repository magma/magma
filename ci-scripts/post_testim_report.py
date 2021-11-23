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

import argparse
import json
from collections import namedtuple

import requests
import xmltodict

# API endpoint for posting a message to Slack
MESSAGE_URL = url = "https://slack.com/api/chat.postMessage"

TestFailure = namedtuple(
    'TestFailure',
    [
        'test_name',
        'failure_msg',
    ],
)


def get_test_data(results_filename):
    """
    Returns dict of test results

    Args:
        results_filename: JUnit XML to be processed
    Returns:
        Dict JSON representation of passed in XML
    """
    file = open(results_filename)
    data = file.read()
    return xmltodict.parse(data)


def get_test_count(test_data):
    """
    Args:
        test_data: json of test data
    Returns:
        int: total test count
    """
    return int(test_data.get("testsuites").get("testsuite").get("@tests"))


def get_test_failures(test_data):
    """
    Args:
        test_data: json of test data
    Returns:
        List[TestFailure]
    """
    test_cases = test_data.get("testsuites").get("testsuite").get("testcase")
    failures = []
    for result in test_cases:
        if "failure" in result:
            name = result.get("@name")
            failure_msg = result.get("failure").get("@message")
            failures.append(
                TestFailure(test_name=name, failure_msg=failure_msg),
            )
    return failures


def blockify(failure):
    """
    Args:
        failure: TestFailure
    Returns:
        Dict - Block Kit representation for Slack posting
    """
    return {
        "type": "section",
        "text": {
            "type": "mrkdwn",
            "text": f'*{failure.test_name}*\n```{failure.failure_msg}```',
        },
    }


def header_blocks(test_count, failures):
    """
    Creates blocks for either success or failure

    Args:
        test_count: int - Total test count
        failures: List[TestFailure]
    Returns:
        List[Dict] - Block Kit representation for Slack posting
    """
    blocks = []
    if len(failures) > 0:
        blocks.append({
            "type": "header",
            "text": {
                "type": "plain_text",
                "text": f'Daily Tests - FAILED {len(failures)} TESTS',
            },
        })
    else:
        blocks.append({
            "type": "header",
            "text": {
                "type": "plain_text",
                "text": 'Daily Tests - PASSED',
            },
        })
    successes = test_count - len(failures)
    blocks.append({
        "type": "context",
        "elements": [
            {
                "type": "mrkdwn",
                "text": f'Passed *{successes}/{test_count}* tests',
            },
        ],
    })
    blocks.append({"type": "divider"})
    return blocks


def post_to_slack(test_count, failures, auth_token, channel_id):
    """
    Post success/failure message to Slack

    Args:
        test_count: int - Total test count
        failures: List[TestFailure]
        auth_token: str -  Slack authorization token
        channel_id: str - Slack channel id
    Returns:
        None
    """
    headers = {
        'content-type': 'application/json',
        'Accept-Charset': 'UTF-8',
        'Authorization': auth_token,
    }
    blocks = header_blocks(test_count, failures)\
        + list(map(blockify, failures))
    payload = {
        "channel": channel_id,
        "blocks": blocks,
    }
    r = requests.post(MESSAGE_URL, data=json.dumps(payload), headers=headers)
    print(r.content)


# Configuring how we parse arguments
parser = argparse.ArgumentParser()
parser.add_argument("file", help="JUnit XML filepath")
parser.add_argument("token", help="Slack Auth token")
parser.add_argument("channel", help="Slack channel id")

# And actually parsing arguments here
args = parser.parse_args()
xml_file = args.file

test_data = get_test_data(xml_file)
test_count = get_test_count(test_data)
failures = get_test_failures(test_data)

post_to_slack(test_count, failures, args.token, args.channel)
