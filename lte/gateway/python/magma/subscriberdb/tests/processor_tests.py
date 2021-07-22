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

import tempfile
import unittest

from lte.protos.mconfig.mconfigs_pb2 import SubscriberDB
from lte.protos.subscriberdb_pb2 import (
    GSMSubscription,
    LTESubscription,
    SubscriberData,
    SubscriberState,
)
from magma.subscriberdb import processor
from magma.subscriberdb.crypto.milenage import BaseLTEAuthAlgo, Milenage
from magma.subscriberdb.crypto.utils import CryptoError
from magma.subscriberdb.crypto.lte import FiveGRanAuthVector
from magma.subscriberdb.sid import SIDUtils
from magma.subscriberdb.store.base import SubscriberNotFoundError
from magma.subscriberdb.store.sqlite import SqliteStore

def _dummy_auth_tuple():
    rand = b'ni\x89\xbel\xeeqTT7p\xae\x80\xb1\xef\r'
    sres = b'\xd4\xac\x8bS'
    key = b'\x9f\xf54.\xb9]\x88\x00'
    return (rand, sres, key)


def _dummy_eutran_vector():
    rand = b'\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f'
    xres = b'\x2d\xaf\x87\x3d\x73\xf3\x10\xc6'
    autn = b'o\xbf\xa3\x80\x1fW\x80\x00{\xdeY\x88n\x96\xe4\xfe'
    kasme = (
        b'\x87H\xc1\xc0\xa2\x82o\xa4\x05\xb1\xe2~\xa1\x04CJ\xe5V\xc7e'
        b'\xe8\xf0a\xeb\xdb\x8a\xe2\x86\xc4F\x16\xc2'
    )
    return (rand, xres, autn, kasme)


def _dummy_resync_vector():
    mac_s = 8 * b'\x00'
    sqn = 2**28 + 1 << 5
    return (sqn, mac_s)


def _dummy_opc():
    return b'\x66\xe9\x4b\xd4\xef\x8a\x2c\x3b\x88\x4c\xfa\x59\xca\x34\x2b\x2e'

def _dummy_m5gran_vector():
    rand = b'u\x96=\x9d\xef\xa4\x15\x0e\x95\x852\xcd\xb8$\xb1\xc1'
    xres_star = (
        b'\xe5\xf2u\x80\\M\xaf}\x82P\xfe?\xb7\xd6\x80jkX8a\x8bP'
        b'\x07\x05\xcbY\xdd}]\xf4\xb2%'
    )
    autn = b'Y\xa5o\x867\xff\x80\x00~kI\x8e\xd4\xab\x0f\xee'
    kseaf = (
        b'\t\xc1,\x15\x14%\xbe\xe1/\xe4IT\x7f\xae\xa6\xecT\xcf'
        b'\xacm#\xbbf|\xebu\rG#\x8b\x04\xd3'
    )    
       
    return FiveGRanAuthVector(rand, xres_star, autn, kseaf)

class FakeMilenage(BaseLTEAuthAlgo):
    # pylint:disable=unused-argument
    def generate_eutran_vector(self, key, opc, sqn, plmn):
        return _dummy_eutran_vector()

    def generate_resync(self, auts, key, opc, rand):
        # AMF should be zeros for resync
        assert self.amf == b'\x00\x00'
        return _dummy_resync_vector()

    # pylint:disable=unused-argument
    def generate_m5gran_vector(self, key, opc, sqn, snni):
        return _dummy_m5gran_vector()

    @classmethod
    def generate_opc(cls, key, op):
        return _dummy_opc()


class ProcessorTests(unittest.TestCase):
    """
    Tests for the Processor
    """

    @classmethod
    def setUpClass(cls):
        processor.Milenage = FakeMilenage

    @classmethod
    def tearDownClass(cls):
        processor.Milenage = Milenage

    def setUp(self):
        # Create sqlite3 database for testing
        self._tmpfile = tempfile.TemporaryDirectory()
        store = SqliteStore(self._tmpfile.name + '/')
        op = 16 * b'\x11'
        amf = b'\x80\x00'
        self._sub_profiles = {
            'superfast': SubscriberDB.SubscriptionProfile(
                max_ul_bit_rate=100000, max_dl_bit_rate=50000,
            ),
        }
        self._default_sub_profile = SubscriberDB.SubscriptionProfile(
            max_ul_bit_rate=10000, max_dl_bit_rate=5000,
        )

        self._processor = processor.Processor(
            store, self._default_sub_profile, self._sub_profiles, op, amf,
        )

        # Add some test users
        (rand, sres, gsm_key) = _dummy_auth_tuple()
        gsm = GSMSubscription(
            state=GSMSubscription.ACTIVE,
            auth_tuples=[rand + sres + gsm_key],
        )
        lte_key = 16 * b'\x00'
        lte = LTESubscription(
            state=LTESubscription.ACTIVE,
            auth_key=lte_key,
        )
        lte_opc = LTESubscription(
            state=LTESubscription.ACTIVE,
            auth_key=lte_key,
            auth_opc=Milenage.generate_opc(lte_key, op),
        )
        lte_opc_short = LTESubscription(
            state=LTESubscription.ACTIVE,
            auth_key=lte_key,
            auth_opc=b'\x00',
        )
        state = SubscriberState(lte_auth_next_seq=1)
        sub1 = SubscriberData(
            sid=SIDUtils.to_pb('IMSI11111'), gsm=gsm, lte=lte,
            state=state, sub_profile='superfast',
        )
        sub2 = SubscriberData(
            sid=SIDUtils.to_pb('IMSI22222'),  # No auth keys
            gsm=GSMSubscription(
                state=GSMSubscription.ACTIVE,
            ),
            lte=LTESubscription(state=LTESubscription.ACTIVE),
        )
        sub3 = SubscriberData(
            sid=SIDUtils.to_pb('IMSI33333'),
        )  # No subscribtion
        sub4 = SubscriberData(
            sid=SIDUtils.to_pb('IMSI44444'), lte=lte_opc,
            state=state,
        )
        sub5 = SubscriberData(
            sid=SIDUtils.to_pb('IMSI55555'), lte=lte_opc_short,
            state=state,
        )
        store.add_subscriber(sub1)
        store.add_subscriber(sub2)
        store.add_subscriber(sub3)
        store.add_subscriber(sub4)
        store.add_subscriber(sub5)

    def tearDown(self):
        self._tmpfile.cleanup()

    def test_gsm_auth_success(self):
        """
        Test if we get the correct auth tuple on success
        """
        self.assertEqual(
            self._processor.get_gsm_auth_vector('11111'),
            _dummy_auth_tuple(),
        )

    def test_gsm_auth_imsi_unknown(self):
        """
        Test if we get SubscriberNotFoundError exception
        """
        with self.assertRaises(SubscriberNotFoundError):
            self._processor.get_gsm_auth_vector('12345')

    def test_gsm_auth_key_missing(self):
        """
        Test if we get CryptoError if auth key is missing
        """
        with self.assertRaises(CryptoError):
            self._processor.get_gsm_auth_vector('22222')

    def test_gsm_auth_no_subscription(self):
        """
        Test if we get CryptoError if there is no GSM subscription
        """
        with self.assertRaises(CryptoError):
            self._processor.get_gsm_auth_vector('33333')

    def test_lte_auth_seq(self):
        """
        Test if we can increment the seq number for the LTE auth
        """
        self.assertEqual(self._processor.get_next_lte_auth_seq('11111'), 1)
        self.assertEqual(self._processor.get_next_lte_auth_seq('11111'), 2)
        self.assertEqual(self._processor.get_next_lte_auth_seq('11111'), 3)

        self._processor.set_next_lte_auth_seq('11111', 200)
        self.assertEqual(self._processor.get_next_lte_auth_seq('11111'), 200)
        self.assertEqual(self._processor.get_next_lte_auth_seq('11111'), 201)

        with self.assertRaises(SubscriberNotFoundError):
            self._processor.get_next_lte_auth_seq('12345')

    def test_lte_auth_success(self):
        """
        Test if we get the auth vector
        """
        eutran_vector = _dummy_eutran_vector()
        self.assertEqual(
            self._processor.generate_lte_auth_vector(
                '11111',
                3 * b'\x00',
            ),
            eutran_vector,
        )

    def test_lte_auth_success_opc(self):
        """
        Test if we get the auth vector using passed OPc
        """
        eutran_vector = _dummy_eutran_vector()
        self.assertEqual(
            self._processor.generate_lte_auth_vector(
                '44444',
                3 * b'\x00',
            ),
            eutran_vector,
        )

    def test_lte_auth_fail_opc_short(self):
        """
        Test if we get the a crypto error if the OPc is too short
        """
        with self.assertRaises(CryptoError):
            self._processor.generate_lte_auth_vector('55555', 3 * b'\x00')

    def test_lte_auth_imsi_unknown(self):
        """
        Test if we get SubscriberNotFoundError exception
        """
        with self.assertRaises(SubscriberNotFoundError):
            self._processor.generate_lte_auth_vector('12345', 3 * b'\x00')

    def test_lte_auth_key_missing(self):
        """
        Test if we get CryptoError if auth key is missing
        """
        with self.assertRaises(CryptoError):
            self._processor.generate_lte_auth_vector('22222', 3 * b'\x00')

    def test_lte_auth_no_subscription(self):
        """
        Test if we get CryptoError if there is no LTE subscription
        """
        with self.assertRaises(CryptoError):
            self._processor.generate_lte_auth_vector('33333', 3 * b'\x00')

    def test_lte_resync_seq_overflow_protect(self):
        """
        Test that we update by the sequence number by a number greater than
        2 **28 if the ue_seq number is higher than the network seq number
        """
        self._processor.set_next_lte_auth_seq('11111', 0)
        auts = 14 * b'\x00'
        self._processor.resync_lte_auth_seq('11111', 16 * b'\x00', auts)
        self.assertEqual(
            self._processor.get_next_lte_auth_seq('11111'),
            2 ** 28 + 2,
        )

    def test_lte_resync_seq_smaller(self):
        """
        Test if we correctly update the seq when the SQN_MS is greater than
        next SQN and the MAC_S is correct. We start with a SEQ of 1 and a SEQ_MS
        of 2**28+1 to get a final SEQ of 2**28 + 2
        """
        auts = 14 * b'\x00'
        self._processor.resync_lte_auth_seq('11111', 16 * b'\x00', auts)
        self.assertEqual(
            self._processor.get_next_lte_auth_seq('11111'),
            2**28 + 2,
        )

    def test_lte_resync_seq_equal(self):
        """
        Test if we update the next seq when the current next is what the USIM
        has. We start with a SEQ of 2**28 + 1 and SEQ_MS is equal so we add 1.
        """
        self._processor.set_next_lte_auth_seq('11111', 2**28 + 1)
        auts = 14 * b'\x00'
        self._processor.resync_lte_auth_seq('11111', 16 * b'\x00', auts)
        self.assertEqual(
            self._processor.get_next_lte_auth_seq('11111'),
            2**28 + 2,
        )

    def test_lte_resync_seq_larger(self):
        """
        Test if we ignore the seq update when the SQN_MS isn't larger than
        next SQN and the MAC_S is correct. In this case the SEQ delta is -1 so
        we should ignore.

        NOTE: this test assumes that we've prematurely incremented the
        SEQ on auth failure.
        """
        self._processor.set_next_lte_auth_seq(
            '11111',
            2**28 + 3,
        )
        auts = 14 * b'\x00'
        with self.assertRaises(CryptoError):
            # seq_ms is stubbed to be 2**28 + 1, and we've set seq_network
            # to be 2**28 + 3. seq_ms will be incremented to 2**28 + 2,
            # meaning we still _don't_ allow a resync.
            self._processor.resync_lte_auth_seq('11111', 16 * b'\x00', auts)
        self.assertEqual(
            self._processor.get_next_lte_auth_seq('11111'),
            2**28 + 3,
        )

    def test_lte_resync_bad_mac_s(self):
        """
        Test if we ignore the seq update when the MAC_S is incorrect
        """
        auts = 14 * b'\x01'
        with self.assertRaises(CryptoError):
            self._processor.resync_lte_auth_seq('11111', 16 * b'\x00', auts)
        self.assertEqual(self._processor.get_next_lte_auth_seq('11111'), 1)

    def test_lte_resync_imsi_unknown(self):
        """
        Test if we get SubscriberNotFoundError exception
        """
        with self.assertRaises(SubscriberNotFoundError):
            self._processor.resync_lte_auth_seq(
                '12345', 16 * b'\x00', 14 * b'\x00',
            )

    def test_lte_resync_key_missing(self):
        """
        Test if we get CryptoError if auth key is missing
        """
        with self.assertRaises(CryptoError):
            self._processor.resync_lte_auth_seq(
                '22222', 16 * b'\x00', 14 * b'\x00',
            )

    def test_lte_resync_no_subscription(self):
        """
        Test if we get CryptoError if there is no LTE subscription
        """
        with self.assertRaises(CryptoError):
            self._processor.resync_lte_auth_seq(
                '33333', 16 * b'\x00', 3 * b'\x00',
            )

    def test_sub_profile(self):
        """
        Test if the sub profile is returned as expected
        """
        with self.assertRaises(SubscriberNotFoundError):
            self._processor.get_sub_profile('12345')
        self.assertEqual(
            self._processor.get_sub_profile('11111'),
            self._sub_profiles['superfast'],
        )
        self.assertEqual(
            self._processor.get_sub_profile('33333'),
            self._default_sub_profile,
        )

    def test_m5g_auth_success(self):
        """
        Test if we get the auth vector
        """
        m5Gran_vector = _dummy_m5gran_vector()
        self.assertEqual(
            self._processor.generate_m5g_auth_vector(
                '11111',
                3 * b'\x00',
            ),
            m5Gran_vector,
        )

    def test_m5g_auth_success_opc(self):
        """
        Test if we get the auth vector using passed OPc
        """
        m5Gran_vector = _dummy_m5gran_vector()
        self.assertEqual(
            self._processor.generate_m5g_auth_vector(
                '44444',
                3 * b'\x00',
            ),
            m5Gran_vector,
        )

if __name__ == "__main__":
    unittest.main()

