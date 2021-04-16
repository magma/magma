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

import abc


class BaseLTEAuthAlgo(metaclass=abc.ABCMeta):
    """
    Abstract class for LTE EUTRAN auth vector
    """

    def __init__(self, amf=b'\x80\x00'):
        """
        Base constructor for auth algos which need to define an operator code
        and authentication management fields.

        Args:
            amf (bytes): 16 bit authentication management field
        """
        self.amf = amf

    @abc.abstractmethod
    def generate_eutran_vector(self, key, opc, sqn, plmn):
        """
        Generate the E-EUTRAN key vector.
        Args:
            key (bytes): 128 bit subscriber key
            opc (bytes): 128 bit operator variant algorithm configuration field
            sqn (int): 48 bit sequence number
            rand (bytes): 128 bit random challenge
            plmn (bytes): 24 bit network identifer
                Octet           Description
                  1      MCC digit 2 | MCC digit 1
                  2      MNC digit 3 | MCC digit 3
                  3      MNC digit 2 | MNC digit 1
        Returns:
            rand (bytes): 128 bit random challenge
            xres (bytes): 128 bit expected result
            autn (bytes): 128 bit authentication token
            kasme (bytes): 256 bit base network authentication code
        """
        pass
