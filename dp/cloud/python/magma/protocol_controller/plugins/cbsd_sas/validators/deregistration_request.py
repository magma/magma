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

from marshmallow import Schema, fields


class DeregistrationRequestObjectSchema(Schema):
    """
    Deregestration Request object validator class
    """
    cbsdId = fields.String(required=True)  # noqa: N815


class DeregistrationRequestSchema(Schema):
    """
    Deregestration Request validator class
    """
    deregistrationRequest = fields.Nested(  # noqa: N815
        DeregistrationRequestObjectSchema, required=True, many=True, unknown='true',
    )
