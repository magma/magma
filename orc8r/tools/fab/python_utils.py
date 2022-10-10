"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""


def strtobool(boolstring: str) -> bool:
    """
    Convert a string representation of truth to true (1) or false (0).
    True values are y, yes, t, true, on and 1;
    False values are n, no, f, false, off and 0.
    Raises ValueError if boolstring is anything else.

    This function is a replacement for distutils.util.strtobool,
    which will be deprecated in Python v3.10. The goal of reimplementing
    this function is to future proof the code. For more information see:
    https://peps.python.org/pep-0632/ 
    """
    original = boolstring
    if not isinstance(boolstring, str):
        raise ValueError("The provided value is not a string.")
    boolstring = boolstring.lower()
    if boolstring in ["y", "yes", "t", "true", "on", "1"]:
        return True
    if boolstring in ["n", "no", "f", "false", "off", "0"]:
        return False
    raise ValueError("Value '{}' could not be converted to bool.".format(original))
