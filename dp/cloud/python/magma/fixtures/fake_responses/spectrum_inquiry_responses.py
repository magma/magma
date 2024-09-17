"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

zero_channels_for_one_cbsd = {
    "spectrumInquiryResponse": [
        {
            "response": {
                "responseCode": 0,
            },
            "cbsdId": "foo",
            "availableChannel": [],
        },
    ],
}

single_channel_for_one_cbsd = {
    "spectrumInquiryResponse": [
        {
            "response": {
                "responseCode": 0,
            },
            "cbsdId": "foo",
            "availableChannel": [
                {
                    "frequencyRange": {
                        "lowFrequency": 1,
                        "highFrequency": 1,
                    },
                    "channelType": "test",
                    "ruleApplied": "test",
                    "maxEirp": 1,
                },
            ],
        },
    ],
}

single_channel_for_one_cbsd_with_no_max_eirp = {
    "spectrumInquiryResponse": [
        {
            "response": {
                "responseCode": 0,
            },
            "cbsdId": "foo",
            "availableChannel": [
                {
                    "frequencyRange": {
                        "lowFrequency": 1,
                        "highFrequency": 1,
                    },
                    "channelType": "test",
                    "ruleApplied": "test",
                },
            ],
        },
    ],
}

two_channels_for_one_cbsd = {
    "spectrumInquiryResponse": [
        {
            "response": {
                "responseCode": 0,
            },
            "cbsdId": "foo",
            "availableChannel": [
                {
                    "frequencyRange": {
                        "lowFrequency": 1,
                        "highFrequency": 10,
                    },
                    "channelType": "test",
                    "ruleApplied": "test",
                    "maxEirp": 1,
                },
                {
                    "frequencyRange": {
                        "lowFrequency": 20,
                        "highFrequency": 30,
                    },
                    "channelType": "test1",
                    "ruleApplied": "test1",
                    "maxEirp": 2,
                },
            ],
        },
    ],
}
