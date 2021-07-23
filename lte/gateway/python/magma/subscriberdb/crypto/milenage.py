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

import hmac
from Crypto.Cipher import AES
from Crypto.Random import random
from .lte import (
    BaseLTEAuthAlgo,
    FiveGRanAuthVector,
)

class Milenage(BaseLTEAuthAlgo):
    """
    Milenage Algorithm (3GPP TS 35.205, .206, .207, .208)
    """

    def generate_eutran_vector(self, key, opc, sqn, plmn):
        """
        Generate the E-EUTRAN key vector.
        Args:
            key (bytes): 128 bit subscriber key
            opc (bytes): 128 bit operator variant algorithm configuration field
            sqn (int): 48 bit sequence number
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
        sqn_bytes = bytearray.fromhex('{:012x}'.format(sqn))
        rand = Milenage.generate_rand()

        mac_a, _ = Milenage.f1(key, sqn_bytes, rand, opc, self.amf)
        xres, ak = Milenage.f2_f5(key, rand, opc)
        ck = Milenage.f3(key, rand, opc)
        ik = Milenage.f4(key, rand, opc)

        autn = Milenage.generate_autn(sqn_bytes, ak, mac_a, self.amf)
        kasme = Milenage.generate_kasme(ck, ik, plmn, sqn_bytes, ak)
        return rand, xres, autn, kasme

    def generate_m5gran_vector(self, key: bytes, opc: bytes, sqn: int, snni: bytes) -> FiveGRanAuthVector:
        """
        Generate the NGRAN key vector.
        Args:
            key : bytes 
                128 bit subscriber key
            opc : bytes 
                128 bit operator variant algorithm configuration field
            sqn : int 
                48 bit sequence number
            snni : bytes 
                32 bit serving network name consisting of MCC and MNC
        Returns:
            FiveGRanAuthVector : NamedTuple 
                 Consists of (rand, xres_star, autn, kseaf)
        """
        sqn_bytes = bytearray.fromhex('{:012x}'.format(sqn))
        rand = Milenage.generate_rand()

        mac_a, _ = Milenage.f1(key, sqn_bytes, rand, opc, self.amf)
        xres, ak = Milenage.f2_f5(key, rand, opc)
        ck = Milenage.f3(key, rand, opc)
        ik = Milenage.f4(key, rand, opc)

        autn = Milenage.generate_autn(sqn_bytes, ak, mac_a, self.amf)
        xres_star = Milenage.generate_m5g_xres_star(ck + ik, snni, rand, xres)
        kausf = Milenage.generate_m5g_kausf(ck + ik, snni, autn)
        kseaf = Milenage.generate_m5g_kseaf(kausf, snni)

        return FiveGRanAuthVector(rand, xres_star, autn, kseaf)

    def generate_auts(self, key: bytes, opc: bytes, rand: bytes,
                      sqn: int) -> bytes:
        """
        Compute AUTS for re-synchronization using the formula
            AUTS = SQN_MS ^ AK || f1*(SQN_MS || RAND || AMF*)
        Args:
            key (bytes): 128 bit subscriber key
            opc (bytes): 128 bit operator variant algorithm configuration field
            rand (bytes): 128 bit random challenge
            sqn (int), 48 bit sequence number
        Returns:
            auts (bytes): 112 bit authentication token
        """
        ak = self.f5_star(key, rand, opc)
        sqn_bytes = bytearray.fromhex('{:012x}'.format(sqn))
        _, mac_s = self.f1(key, sqn_bytes, rand, opc, self.amf)
        return xor(sqn_bytes, ak) + mac_s

    def generate_resync(self, auts, key, opc, rand):
        """
        Compute SQN_MS and MAC-S from AUTS for re-synchronization
            AUTS = SQN_MS ^ AK || f1*(SQN_MS || RAND || AMF*)
        Args:
            auts (bytes): 112 bit authentication token from client key
            opc (bytes): 128 bit operator variant algorithm configuration field
            key (bytes): 128 bit subscriber key
            rand (bytes): 128 bit random challenge
        Returns:
            sqn_ms (int), 48 bit sequence number from client
            mac_s (bytes), 64 bit resync authentication code
        """
        ak = self.f5_star(key, rand, opc)
        sqn_ms = xor(auts[:6], ak)
        sqn_ms_int = int.from_bytes(sqn_ms, byteorder='big')
        _, mac_s = self.f1(key, sqn_ms, rand, opc, self.amf)
        return sqn_ms_int, mac_s

    @classmethod
    def f1(cls, key, sqn, rand, opc, amf):
        """
        Implementation of f1 and f1*, the network authentication function and
        the re-synchronisation message authentication function according to
        3GPP 35.206 4.1

        Args:
            key (bytes): 128 bit subscriber key
            sqn (bytes): 48 bit sequence number
            rand (bytes): 128 bit random challenge
            opc (bytes): 128 bit computed from OP and subscriber key
            amf (bytes): 16 bit authentication management field
        Returns:
            (64 bit Network auth code, 64 bit Resynch auth code)
        """
        # TEMP = E_K(RAND XOR OP_C)
        temp = cls.encrypt(key, xor(rand, opc))

        # IN1 = SQN || AMF || SQN || AMF
        in1 = (sqn[0:6] + amf[0:2]) * 2

        # Constants from 3GPP 35.206 4.1
        c1 = 16 * b'\x00'  # some constant
        r1 = 8  # rotate by 8 bytes

        # OUT1 = E_K(TEMP XOR rotate(IN1 XOR OP_C, r1) XOR c1) XOR OP_C
        out1_ = cls.encrypt(key, xor(temp, rotate(xor(in1, opc), r1)), c1)
        out1 = xor(opc, out1_)

        #  MAC-A = f1 = OUT1[0] .. OUT1[63]
        #  MAC-S = f1* = OUT1[64] .. OUT1[127]
        return out1[:8], out1[8:]

    @classmethod
    def f2_f5(cls, key, rand, opc):
        """
        Implementation of f2 and f5, the compute anonymity key and response to
        challenge functions according to 3GPP 35.206 4.1

        Args:
            key (bytes): 128 bit subscriber key
            rand (bytes): 128 bit random challenge
            opc (bytes): 128 bit computed from OP and subscriber key
        Returns:
            (xres, ak) = (64 bit response to challenge, 48 bit anonymity key)
        """
        # Constants from 3GPP 35.206 4.1
        c2 = 15 * b'\x00' + b'\x01'  # some constant
        r2 = 0  # rotate by 0 bytes

        # TEMP = E_K(RAND XOR OP_C)
        # OUT2 = E_K(rotate(TEMP XOR OP_C, r2) XOR c2) XOR OP_C
        temp_x_opc = xor(cls.encrypt(key, xor(rand, opc)), opc)
        out2 = xor(cls.encrypt(key, xor(rotate(temp_x_opc, r2), c2)), opc)
        # res = f2 = OUT2[64] ... OUT2[127]
        # ak = f5 = OUT2[0] ... OUT2[47]
        return out2[8:16], out2[0:6]

    @classmethod
    def f3(cls, key, rand, opc):
        """
        Implementation of f3, the compute confidentiality key according
        to 3GPP 35.206 4.1

        Args:
            key (bytes): 128 bit subscriber key
            rand (bytes): 128 bit random challenge
            opc (bytes): 128 bit computed from OP and subscriber key
        Returns:
            ck, 128 bit confidentiality key
        """
        # Constants from 3GPP 35.206 4.1
        c3 = 15 * b'\x00' + b'\x02'  # some constant
        r3 = 4  # rotate by 4 bytes

        # TEMP = E_K(RAND XOR OP_C)
        # OUT3 = E_K(rotate(TEMP XOR OP_C, r3) XOR c3) XOR OP_C
        temp_x_opc = xor(cls.encrypt(key, xor(rand, opc)), opc)
        out3 = xor(cls.encrypt(key, xor(rotate(temp_x_opc, r3), c3)), opc)
        # ck = f3 = OUT3
        return out3

    @classmethod
    def f4(cls, key, rand, opc):
        """
        Implementation of f4, the integrity key according
        to 3GPP 35.206 4.1

        Args:
            key (bytes): 128 bit subscriber key
            rand (bytes): 128 bit random challenge
            opc (bytes): 128 bit computed from OP and subscriber key
        Returns:
            ik, 128 bit integrity key
        """
        # Constants from 3GPP 35.206 4.1
        c4 = 15 * b'\x00' + b'\x04'  # some constant
        r4 = 8  # rotate by 8 bytes

        # TEMP = E_K(RAND XOR OP_C)
        # OUT4 = E_K(rotate(TEMP XOR OP_C, r4) XOR c4) XOR OP_C
        temp_x_opc = xor(cls.encrypt(key, xor(rand, opc)), opc)
        out4 = xor(cls.encrypt(key, xor(rotate(temp_x_opc, r4), c4)), opc)
        # ik = f4 = OUT4
        return out4

    @classmethod
    def f5_star(cls, key, rand, opc):
        """
        Implementation of f5*, the anonymity key according
        to 3GPP 35.206 4.1

        Args:
            key (bytes): 128 bit subscriber key
            rand (bytes): 128 bit random challenge
            opc (bytes): 128 bit computed from OP and subscriber key
        Returns:
            ak, 48 bit anonymity key
        """
        # Constants from 3GPP 35.206 4.1
        c5 = 15 * b'\x00' + b'\x08'  # some constant
        r5 = 12  # rotate by 12 bytes

        # TEMP = E_K(RAND XOR OP_C)
        # OUT5 = E_K(rotate(TEMP XOR OP_C, r5 XOR c5) XOR OP_C
        temp_x_opc = xor(cls.encrypt(key, xor(rand, opc)), opc)
        out5 = xor(cls.encrypt(key, xor(rotate(temp_x_opc, r5), c5)), opc)
        # ak = f5* = OUT5[0] . OUT5[47]
        return out5[:6]

    @classmethod
    def generate_kasme(cls, ck, ik, plmn, sqn, ak):
        """
        KASME derivation function (S_2) according to 3GPP 33.401 Annex A.2.
        This function creates an input string to a key deriviation function.

        The input string to the KDF is composed of 2 input parameters P0, P1
        and their lengths L0, L1 a constant FC which identifies this algorithm.
                        S = FC || P0 || L0 || P1 || L1
        The FC = 0x10 and argument P0 is the 3 octets of the PLMN, and P1 is
        SQN XOR AK. The lengths are in bytes.

        The Kasme is computed by calling the key derivation function with S
        using key CK || IK

        Args:
            ck (bytes): 128 bit confidentiality key
            ik (bytes): 128 bit integrity key
            plmn (bytes): 24 bit network identifer
                Octet           Description
                  1      MCC digit 2 | MCC digit 1
                  2      MNC digit 3 | MCC digit 3
                  3      MNC digit 2 | MNC digit 1
            sqn (bytes): 48 bit sequence number
            ak (bytes): 48 bit anonymity key
        Returns:
            256 bit network base key
        """
        S = b'\x10' + plmn + b'\x00\x03' + xor(sqn, ak) + b'\x00\x06'
        return cls.KDF(ck + ik, S)

    @classmethod
    def generate_rand(cls):
        """
        Generate RAND for Milenage
        Returns:
            (bytes) 128 random bits
        """
        return bytearray.fromhex('{:032x}'.format(random.getrandbits(128)))

    @classmethod
    def generate_opc(cls, key, op):
        """
        Generate the OP_c according to 3GPP 35.205 8.2
        Args:
            key (bytes): 128 bit subscriber key
            op (bytes): 128 bit operator dependent value
        Returns:
            128 bit OP_c
        """
        opc = cls.encrypt(key, op)
        return xor(opc, op)

    @classmethod
    def generate_autn(cls, sqn, ak, mac_a, AMF=b'\x80\x00'):
        """
        Generate network authentication token as defined in 3GPP 25.205 7.2

        Args:
            sqn (bytes): 48 bit sequence number
            ak (bytes): 48 bit anonymity key
            mac_a (bytes): 64 bit network authentication code
            AMF (bytes): 16 bit authentication management field
        Returns:
            autn (bytes): 128 bit authentication token
        """
        return xor(sqn, ak) + AMF + mac_a

    @classmethod
    def KDF(cls, key, buf):
        """
        3GPP Key Derivation Function defined in TS 33.220 to be hmac-sha256

        Args:
            key (bytes): 128 bit secret key
            buf (bytes): the buffer to compute the key from
        Returns:
            258 bit key
        """
        return hmac.new(key, buf, 'sha256').digest()

    @classmethod
    def encrypt(cls, k, buf, IV=16 * b'\x00'):
        """
        Rijndael (AES-128) cipher function used by Milenage

        Args:
            k (bytes): 128 bit encryption key
            buf (bytes): 128 bit buffer to encrypt
            IV (bytes): 128 bit initialization vector
        Returns:
            encrypted output
        """
        aes_cipher = AES.new(k, AES.MODE_CBC, IV)
        return aes_cipher.encrypt(buf)

    @classmethod
    def generate_m5g_xres_star(cls, key: bytes, snni: bytes, rand: bytes,
                               xres: bytes) -> bytes:
        """
        Compute XRES_STAR = K + FC + P0 + L0 + P1 + L1 + P2 + L2
        Args:
            key (bytes): key for xres calculation ck+ik
            snni: Serving network name
            rand (bytes): 128 bit random challenge
            xres (bytes): 64 bit response to challenge
        Returns:
            xres_star (bytes)
        """

        snni_len=len(snni)
        rand_len=len(rand)
        xres_len=len(xres)

        K  =  key
        FC =  bytes.fromhex('6B')
        P0 =  snni
        L0 =  snni_len.to_bytes(2, 'big')
        P1 =  rand
        L1 =  rand_len.to_bytes(2, 'big')
        P2 =  xres
        L2 =  xres_len.to_bytes(2, 'big')

        S = FC + P0 + L0 + P1 + L1 + P2 + L2

        return cls.KDF(K, S)

    @classmethod
    def generate_m5g_kausf(cls, key: bytes, snni: bytes, autn: bytes) -> bytes:
        """
        Compute KAUSF = K + FC + P0 + L0i + P1 + L1
        Args:
            key (bytes) : key combination of ck and ik
            snni (bytes): serving network name
            autn (bytes): 128 bit authentication token
        Returns:
            kausf (bytes) : Key from Authentication server function
        """
        K  =  key
        FC =  bytes.fromhex('6A')
        P0 =  snni
        L0 =  len(snni).to_bytes(2, 'big')
        P1 =  autn[:6]
        L1 =  len(P1).to_bytes(2, 'big')

        S = FC + P0 + L0 + P1 + L1

        return cls.KDF(K, S)

    @classmethod
    def generate_m5g_kseaf(cls, kausf: bytes, snni: bytes) -> bytes:
        """
        Compute KSEAF = K + FC + P0 + L0
        Args:
            key (bytes) : key from authentication server function
            snni (bytes) : serving network name
        Returns:
            kseaf (bytes) : keys from security anchor function
        """
        K  =  kausf
        FC =  bytes.fromhex('6C')
        P0 =  snni
        L0 =  len(snni).to_bytes(2, 'big')

        S = FC + P0 + L0

        return cls.KDF(K, S)

def xor(s1, s2):
    """
    Exclusive-Or of two byte arrays

    Args:
        s1 (bytes): first set of bytes
        s2 (bytes): second set of bytes
    Returns:
        (bytes) s1 ^ s2
    Raises:
        ValueError if s1 and s2 lengths don't match
    """
    if len(s1) != len(s2):
        raise ValueError('Input not equal length: %d %d' % (len(s1), len(s2)))
    return bytes(a ^ b for a, b in zip(s1, s2))


def rotate(input_s, bytes_):
    """
    Rotate a string by a number of bytes

    Args:
        input_s (bytes): the input string
        bytes_ (int): the number of bytes to rotate by
    Returns:
        (bytes) s1 rotated by n bytes
    """
    return bytes(
        input_s[(i + bytes_) % len(input_s)] for i in range(
            len(
            input_s,
            ),
        )
    )
    
