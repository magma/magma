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

grant_requests = [
    {
        "grantRequest": [
            {
                "cbsdId": "foo",
                "operationParam": {
                    "maxEirp": 100,
                    "operationFrequencyRange":
                        {
                            "lowFrequency": 3560000000,
                            "highFrequency": 3561000000,
                        },
                },
            },
        ],
    },
    {
        "grantRequest": [
            {
                "cbsdId": "bar",
                "operationParam": {
                    "maxEirp": 100,
                    "operationFrequencyRange":
                        {
                            "lowFrequency": 3560000000,
                            "highFrequency": 3561000000,
                        },
                },
            },
        ],
    },
]
