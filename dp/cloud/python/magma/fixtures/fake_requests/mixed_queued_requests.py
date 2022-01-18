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

queued_requests = [
    {
        "registrationRequest": [
            {
                "fccId": "321cba",
                "cbsdCategory": "B",
                "callSign": "WSD987",
                "userId": "John Doe",
                "airInterface": {
                    "radioTechnology": "E_UTRA",
                },
                "cbsdSerialNumber": "4321dcba",
                "measCapability": [
                    "RECEIVED_POWER_WITHOUT_GRANT",
                ],
                "installationParam": {
                    "latitude": 37.425056,
                    "longitude": -122.084113,
                    "height": 9.3,
                    "heightType": "AGL",
                    "indoorDeployment": False,
                    "antennaAzimuth": 271,
                    "antennaDowntilt": 3,
                    "antennaGain": 16,
                    "antennaBeamwidth": 30,
                },
                "groupingParam": [
                    {
                        "groupId": "example-group-3",
                        "groupType": "INTERFERENCE_COORDINATION",
                    },
                ],
            },
        ],
    },
    {
        "spectrumInquiryRequest": [
            {
                "fccId": "abc123",
                "cbsdCategory": "A",
                "callSign": "CB987",
                "userId": "John Doe",
                "airInterface": {
                    "radioTechnology": "E_UTRA",
                },
                "cbsdSerialNumber": "abcd1234",
                "measCapability": [
                    "RECEIVED_POWER_WITHOUT_GRANT",
                ],
                "installationParam": {
                    "latitude": 37.419735,
                    "longitude": -122.072205,
                    "height": 6,
                    "heightType": "AGL",
                    "indoorDeployment": True,
                },
                "groupingParam": [
                    {
                        "groupId": "example-group-1",
                        "groupType": "INTERFERENCE_COORDINATION",
                    },
                    {
                        "groupId": "example-group-2",
                        "groupType": "INTERFERENCE_COORDINATION",
                    },
                ],
            },
        ],
    },
]
