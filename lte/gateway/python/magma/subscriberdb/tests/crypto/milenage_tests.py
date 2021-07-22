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

import unittest

from magma.subscriberdb.crypto.milenage import Milenage


class MilenageRandomTests(unittest.TestCase):
    """
    Test class RAND method
    """

    def test_rand(self):
        """Tests the random generator"""
        rand = Milenage.generate_rand()
        self.assertEqual(len(rand), 16)


class MilenageTests(unittest.TestCase):
    """
    Test class for the Crypto algorithms
    """

    def setUp(self):
        self.rand = None
        self._old_rand = Milenage.generate_rand
        Milenage.generate_rand = lambda: self.rand

    def tearDown(self):
        Milenage.generate_rand = self._old_rand

    def test_f1(self):
        """
        Test if the f1, f1* algos works as expected.
        This is a test set 1 from 3GPP 35.207 4.3
        """
        # Fake the random
        self.rand = b'#U<\xbe\x967\xa8\x9d!\x8a\xe6M\xaeG\xbf5'
        # Inputs
        k = b'\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc'
        op = b'\xcd\xc2\x02\xd5\x12> \xf6+mgj\xc7,\xb3\x18'
        sqn = b'\xff\x9b\xb4\xd0\xb6\x07'
        amf = b'\xb9\xb9'

        # Outputs
        opc = b'\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf'
        mac_a = b'\x4a\x9f\xfa\xc3\x54\xdf\xaf\xb3'
        mac_s = b'\x01\xcf\xaf\x9e\xc4\xe8\x71\xe9'

        self.assertEqual(Milenage.generate_opc(k, op), opc)
        self.assertEqual(
            Milenage.f1(k, sqn, self.rand, opc, amf),
            (mac_a, mac_s),
        )

    def test_f2_f3_f5(self):
        """ Tests that the f2 and f5 functions work as expected.
        This is a test set 1 from 3GPP 35.207 5.3
        """
        # Fake the random
        self.rand = b'#U<\xbe\x967\xa8\x9d!\x8a\xe6M\xaeG\xbf5'

        # Inputs
        k = b'\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc'
        op = b'\xcd\xc2\x02\xd5\x12> \xf6+mgj\xc7,\xb3\x18'

        # Outputs
        opc = b'\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf'
        f2 = b'\xa5\x42\x11\xd5\xe3\xba\x50\xbf'
        f5 = b'\xaa\x68\x9c\x64\x83\x70'
        f3 = b'\xb4\x0b\xa9\xa3\xc5\x8b\x2a\x05\xbb\xf0\xd9\x87\xb2\x1b\xf8\xcb'

        self.assertEqual(Milenage.generate_opc(k, op), opc)
        xres, ak = Milenage.f2_f5(k, self.rand, opc)
        ck = Milenage.f3(k, self.rand, opc)
        self.assertEqual(xres, f2)
        self.assertEqual(ak, f5)
        self.assertEqual(ck, f3)

    def test_f4(self):
        """ Tests that the f4 function works as expected.
        This is a test set 1 from 3GPP 35.207 6.3
        """
        # Fake the random
        self.rand = b'#U<\xbe\x967\xa8\x9d!\x8a\xe6M\xaeG\xbf5'

        # Inputs
        k = b'\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc'
        op = b'\xcd\xc2\x02\xd5\x12> \xf6+mgj\xc7,\xb3\x18'

        # Outputs
        opc = b'\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf'
        f4 = b'\xf7\x69\xbc\xd7\x51\x04\x46\x04\x12\x76\x72\x71\x1c\x6d\x34\x41'

        self.assertEqual(Milenage.generate_opc(k, op), opc)
        ik = Milenage.f4(k, self.rand, opc)
        self.assertEqual(ik, f4)

    def test_f5_star(self):
        """ Tests that the f5 function works as expected.
        This is a test set 1 from 3GPP 35.207 6.3
        """
        # Fake the random
        self.rand = b'#U<\xbe\x967\xa8\x9d!\x8a\xe6M\xaeG\xbf5'

        # Inputs
        k = b'\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc'
        op = b'\xcd\xc2\x02\xd5\x12> \xf6+mgj\xc7,\xb3\x18'

        # Outputs
        opc = b'\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf'
        ak = b'\x45\x1e\x8b\xec\xa4\x3b'

        self.assertEqual(Milenage.generate_opc(k, op), opc)
        self.assertEqual(Milenage.f5_star(k, self.rand, opc), ak)

    def test_generate_resync(self):
        """ Tests that that we compute the seq and mac_s correctly given
        an auts during re-synchronisation.
        """

        # Inputs
        rand = b'\xcd\x14\xa7S\x97\x7f\xbcq\x8eb\xbd\xdbS]\x88\xf8'
        k = b'\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb'
        op = 16 * b'\x11'
        amf = b'\x00\x00'

        # Outputs
        mac_s = b'\xdb_c`Y\x1f4\xea'
        op_c = b"\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]"
        sqn = 0

        crypto = Milenage(amf)
        self.assertEqual(crypto.generate_opc(k, op), op_c)
        auts = crypto.generate_auts(k, op_c, rand, sqn)
        self.assertEqual(
            crypto.generate_resync(auts, k, op_c, rand),
            (sqn, mac_s),
        )

    def test_eutran_vector(self):
        """Can we compute the vector that OAI generates?"""
        self.rand = (
            b'\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b'
            b'\x0c\r\x0e\x0f'
        )

        # Inputs extracted from OAI logs
        key = b'\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb'
        sqn = 7351
        op = 16 * b'\x11'
        amf = b'\x80\x00'
        plmn = b'\x02\xf8\x59'

        # Outputs extracted from OAI logs
        op_c = b"\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]"
        xres = b'\x2d\xaf\x87\x3d\x73\xf3\x10\xc6'
        autn = b'o\xbf\xa3\x80\x1fW\x80\x00{\xdeY\x88n\x96\xe4\xfe'
        kasme = (
            b'\x87H\xc1\xc0\xa2\x82o\xa4\x05\xb1\xe2~\xa1\x04CJ\xe5V\xc7e'
            b'\xe8\xf0a\xeb\xdb\x8a\xe2\x86\xc4F\x16\xc2'
        )
        crypto = Milenage(amf)
        self.assertEqual(crypto.generate_opc(key, op), op_c)
        (rand_, xres_, autn_, kasme_) = \
            crypto.generate_eutran_vector(key, op_c, sqn, plmn)
        
        self.assertEqual(self.rand, rand_)
        self.assertEqual(xres, xres_)
        self.assertEqual(autn, autn_)
        self.assertEqual(kasme, kasme_)

    def test_m5g_xres_star_vector(self):
        """Can we compute the vector that OAI generates?"""
        self.rand = (b'\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b'
                     b'\x0c\r\x0e\x0f')

        # Inputs extracted from OAI logs
        key = b'F[\\\xe8\xb1\x99\xb4\x9f\xaa_\n.\xe28\xa6\xbc'
        sqn = 20672
        op = 16 * b'\x11'
        amf = b'\x80\x00'
        serving_network = "5G:mnc456.mcc222.3gppnetwork.org"

        # Outputs extracted from OAI logs
        op_c = b'\xc4\xd5\xe49\x91\xb0\xc5Q\xaf\xf8\xb9%<\x131\xab'
        autn = b'\xba\xb4\x99\xbe\xd2"\x80\x00\xfa)\x93+?1\xde\x05'
        
        xres_star = (
                      b'\x99\xd5\x9fA\xdf\xae2\xbd\xcdG\x13\x94\x0e'
                      b'\x11svg4\xc2\x0c\xd9\xb8|";.\x07A\xf42\x07\xb5'
        )
        crypto = Milenage(amf)
        self.assertEqual(crypto.generate_opc(key, op), op_c)
        fiveg_ran_auth_vectors = \
            crypto.generate_m5gran_vector(key, op_c, sqn, serving_network.encode('utf-8'))
        self.assertEqual(self.rand, fiveg_ran_auth_vectors.rand)
        self.assertEqual(xres_star, fiveg_ran_auth_vectors.xres_star)
        self.assertEqual(autn, fiveg_ran_auth_vectors.autn)

    def test_m5g_kausf_vector(self):
        """Can we compute the vector that OAI generates?"""
        self.rand = (b'\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b'
                     b'\x0c\r\x0e\x0f')

        # Inputs extracted from OAI logs
        key = b'F[\\\xe8\xb1\x99\xb4\x9f\xaa_\n.\xe28\xa6\xbc'
        sqn = 20672
        op = 16 * b'\x11'
        amf = b'\x80\x00'
        serving_network = "5G:mnc456.mcc222.3gppnetwork.org"
        snni = serving_network.encode('utf-8')

        # Outputs extracted from OAI logs
        op_c = b'\xc4\xd5\xe49\x91\xb0\xc5Q\xaf\xf8\xb9%<\x131\xab'
        autn = b'\xba\xb4\x99\xbe\xd2"\x80\x00\xfa)\x93+?1\xde\x05'
        kausf= (
                b'\xed\x08\xc3Z\x0b\x93\x88\xdfr\x9a\x9a6\x80e\xd91'
                b'\x9a\x12\x14\x95g\x9c1\xe6\xcd\x14(\xd0W$\x10\xac'
        )        
        crypto = Milenage(amf)
        self.assertEqual(crypto.generate_opc(key, op), op_c)
        fiveg_ran_auth_vectors = \
            crypto.generate_m5gran_vector(key, op_c, sqn, snni)
        self.assertEqual(self.rand, fiveg_ran_auth_vectors.rand)
        self.assertEqual(autn, fiveg_ran_auth_vectors.autn)

        ck = Milenage.f3(key, self.rand, op_c) 
        ik = Milenage.f4(key, self.rand, op_c) 

        kausf_ = crypto.generate_m5g_kausf(ck + ik, snni, autn)
        self.assertEqual(kausf, kausf_)

    def test_m5g_kseaf_vector(self):
        """Can we compute the vector that OAI generates?"""
        self.rand = (b'\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b'
                     b'\x0c\r\x0e\x0f')

        # Inputs extracted from OAI logs
        key = b'F[\\\xe8\xb1\x99\xb4\x9f\xaa_\n.\xe28\xa6\xbc'
        sqn = 20672
        op = 16 * b'\x11'
        amf = b'\x80\x00'
        serving_network = "5G:mnc456.mcc222.3gppnetwork.org"

        # Outputs extracted from OAI logs
        op_c = b'\xc4\xd5\xe49\x91\xb0\xc5Q\xaf\xf8\xb9%<\x131\xab'
        autn = b'\xba\xb4\x99\xbe\xd2"\x80\x00\xfa)\x93+?1\xde\x05'
        kseaf = (
                  b'\x15\xb9\xf0 M\\C\xcbJt\xc9\xa3\xf8\xab\xa5\xafNV'
                  b'\xc4\xc6eG\xf4\x13\xa7\x99\xcc\xf0\xd6[T`'
        )
        crypto = Milenage(amf)
        self.assertEqual(crypto.generate_opc(key, op), op_c)
        fiveg_ran_auth_vectors = \
            crypto.generate_m5gran_vector(key, op_c, sqn, serving_network.encode('utf-8'))
        self.assertEqual(self.rand, fiveg_ran_auth_vectors.rand)
        self.assertEqual(autn, fiveg_ran_auth_vectors.autn)
        self.assertEqual(kseaf, fiveg_ran_auth_vectors.kseaf)

if __name__ == "__main__":
    unittest.main()
