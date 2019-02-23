"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
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
