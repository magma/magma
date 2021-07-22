
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

from lte.protos.subscriberdb_pb2 import (
    GSMSubscription,
    LTESubscription,
    SubscriberID,
)
from magma.subscriberdb.sid import SIDUtils

from .crypto.gsm import UnsafePreComputedA3A8
from .crypto.milenage import Milenage
from .crypto.utils import CryptoError


class GSMProcessor(metaclass=abc.ABCMeta):
    """
    Interface for the GSM protocols to interact with other parts
    of subscriberdb.
    """

    @abc.abstractmethod
    def get_gsm_auth_vector(self, imsi):
        """
        Returns the gsm auth tuple for the subsciber by querying the store
        for the secret key.

        Args:
            imsi: the subscriber identifier
        Returns:
            the auth tuple (rand, sres, key) returned by the crypto object
        Raises:
            SubscriberNotFoundError if the subscriber is not present
            CryptoError if the auth tuple couldn't be generated
        """
        raise NotImplementedError()


class LTEProcessor(metaclass=abc.ABCMeta):
    """
    Interface for the LTE protocols to interact with other parts
    of subscriberdb.
    """

    @abc.abstractmethod
    def get_next_lte_auth_seq(self, imsi):
        """
        Returns the sequence number for the next auth operation.

        Args:
            imsi: IMSI string
        Returns:
            the seq number that can be used for the next auth
        Raises:
            SubscriberNotFoundError if the subscriber is not present
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def set_next_lte_auth_seq(self, imsi, seq):
        """
        Updates the LTE auth sequence number.

        Args:
            imsi: IMSI string
            seq: the seq number that will be used during the next auth
        Raises:
            SubscriberNotFoundError if the subscriber is not present
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def generate_lte_auth_vector(self, imsi, plmn):
        """
        Returns the E-UTRAN key vector for the subscriber by querying the
        store for the secret key and generating the vector with the crypto algo
        Args:
            imsi: the subscriber identifier
            plmn (bytes): 24 bit network identifer
        Returns:
            rand (bytes): 128 bit random challenge
            xres (bytes): 128 bit expected result
            autn (bytes): 128 bit authentication token
            kasme (bytes): 256 bit base network authentication code
        Raises:
            SubscriberNotFoundError if the subscriber is not present
            CryptoError if the auth tuple couldn't be generated
        """
        raise NotImplementedError()


class Processor(GSMProcessor, LTEProcessor):
    """
    Core class which glues together all protocols, crypto algorithms and
    subscriber stores.
    """

    def __init__(
        self, store, default_sub_profile,
        sub_profiles, op=None, amf=None,
    ):
        """
        Init the Processor with all the components.

        We use the UnsafePreComputedA3A8 crypto by default for
        GSM authentication. This requires the auth-tuple to be stored directly
        in the store as the key for the subscriber.
        """
        self._store = store
        self._op = op
        self._amf = amf
        self._default_sub_profile = default_sub_profile
        self._sub_profiles = sub_profiles
        if len(op) != 16:
            raise ValueError("OP is invalid len=%d value=%s" % (len(op), op))
        if len(amf) != 2:
            raise ValueError(
                "AMF has invalid length len=%d value=%s" %
                (len(amf), amf),
            )

    def get_sub_profile(self, imsi):
        """
        Returns the subscription profile for subscriber. If the subscriber
        has not profile configured, the default profile is returned.

        Args:
            imsi: IMSI string
        Returns:
            SubscriptionProfile proto struct
        Raises:
            SubscriberNotFoundError if the subscriber is not present
        """
        sid = SIDUtils.to_str(SubscriberID(id=imsi, type=SubscriberID.IMSI))
        subs = self._store.get_subscriber_data(sid)
        return self._sub_profiles.get(
            subs.sub_profile,
            self._default_sub_profile,
        )

    def get_gsm_auth_vector(self, imsi):
        """
        Returns the gsm auth tuple for the subsciber by querying the store
        for the crypto algo and secret keys.
        """
        sid = SIDUtils.to_str(SubscriberID(id=imsi, type=SubscriberID.IMSI))
        subs = self._store.get_subscriber_data(sid)

        if subs.gsm.state != GSMSubscription.ACTIVE:
            raise CryptoError("GSM service not active for %s" % sid)

        # The only GSM crypto algo we support now
        if subs.gsm.auth_algo != GSMSubscription.PRECOMPUTED_AUTH_TUPLES:
            raise CryptoError(
                "Unknown crypto (%s) for %s" %
                (subs.gsm.auth_algo, sid),
            )
        gsm_crypto = UnsafePreComputedA3A8()

        if len(subs.gsm.auth_tuples) == 0:
            raise CryptoError("Auth key not present for %s" % sid)

        return gsm_crypto.generate_auth_tuple(subs.gsm.auth_tuples[0])

    def generate_lte_auth_vector(self, imsi, plmn):
        """
        Returns the lte auth vector for the subscriber by querying the store
        for the crypto algo and secret keys.
        """
        sid = SIDUtils.to_str(SubscriberID(id=imsi, type=SubscriberID.IMSI))
        subs = self._store.get_subscriber_data(sid)

        if subs.lte.state != LTESubscription.ACTIVE:
            raise CryptoError("LTE service not active for %s" % sid)

        if subs.lte.auth_algo != LTESubscription.MILENAGE:
            raise CryptoError(
                "Unknown crypto (%s) for %s" %
                (subs.lte.auth_algo, sid),
            )

        if len(subs.lte.auth_key) != 16:
            raise CryptoError("Subscriber key not valid for %s" % sid)

        if len(subs.lte.auth_opc) == 0:
            opc = Milenage.generate_opc(subs.lte.auth_key, self._op)
        elif len(subs.lte.auth_opc) != 16:
            raise CryptoError("Subscriber OPc is invalid length for %s" % sid)
        else:
            opc = subs.lte.auth_opc

        sqn = self.seq_to_sqn(self.get_next_lte_auth_seq(imsi))
        milenage = Milenage(self._amf)
        return milenage.generate_eutran_vector(
            subs.lte.auth_key,
            opc, sqn, plmn,
        )

    def resync_lte_auth_seq(self, imsi, rand, auts):
        """
        Validates a re-synchronization request and computes the SEQ from
        the AUTS sent by U-SIM
        """
        sid = SIDUtils.to_str(SubscriberID(id=imsi, type=SubscriberID.IMSI))
        subs = self._store.get_subscriber_data(sid)

        if subs.lte.state != LTESubscription.ACTIVE:
            raise CryptoError("LTE service not active for %s" % sid)

        if subs.lte.auth_algo != LTESubscription.MILENAGE:
            raise CryptoError(
                "Unknown crypto (%s) for %s" %
                (subs.lte.auth_algo, sid),
            )

        if len(subs.lte.auth_key) != 16:
            raise CryptoError("Subscriber key not valid for %s" % sid)

        if len(subs.lte.auth_opc) == 0:
            opc = Milenage.generate_opc(subs.lte.auth_key, self._op)
        elif len(subs.lte.auth_opc) != 16:
            raise CryptoError("Subscriber OPc is invalid length for %s" % sid)
        else:
            opc = subs.lte.auth_opc

        dummy_amf = b'\x00\x00'  # Use dummy AMF for re-synchronization
        milenage = Milenage(dummy_amf)
        sqn_ms, mac_s = \
            milenage.generate_resync(auts, subs.lte.auth_key, opc, rand)

        if mac_s != auts[6:]:
            raise CryptoError("Invalid resync authentication code")

        seq_ms = self.sqn_to_seq(sqn_ms)

        # current_seq_number was the seq number the network sent
        # to the mobile station as part of the original auth request.
        current_seq_number = subs.state.lte_auth_next_seq - 1
        if seq_ms >= current_seq_number:
            self.set_next_lte_auth_seq(imsi, seq_ms + 1)
        else:
            seq_delta = current_seq_number - seq_ms
            if seq_delta > (2 ** 28):
                self.set_next_lte_auth_seq(imsi, seq_ms + 1)
            else:
                # This shouldn't have happened
                raise CryptoError(
                    "Re-sync delta in range but UE rejected "
                    "auth: %d" % seq_delta,
                )

    def get_next_lte_auth_seq(self, imsi):
        """
        Returns the sequence number for the next auth operation.
        """
        sid = SIDUtils.to_str(SubscriberID(id=imsi, type=SubscriberID.IMSI))

        # Increment the sequence number.
        # The 3GPP TS 33.102 spec allows wrapping around the maximum value.
        # The re-synchronization mechanism would be used to sync the counter
        # between USIM and HSS when it happens.
        with self._store.edit_subscriber(sid) as subs:
            seq = subs.state.lte_auth_next_seq
            subs.state.lte_auth_next_seq += 1
        return seq

    def set_next_lte_auth_seq(self, imsi, seq):
        """
        Updates the LTE auth sequence number.
        """
        sid = SIDUtils.to_str(SubscriberID(id=imsi, type=SubscriberID.IMSI))

        with self._store.edit_subscriber(sid) as subs:
            subs.state.lte_auth_next_seq = seq

    def get_sub_data(self, imsi):
        """
        Returns the complete subscriber profile for subscriber.
        Args:
            imsi: IMSI string
        Returns:
            SubscriberData proto struct
        """
        sid = SIDUtils.to_str(SubscriberID(id=imsi, type=SubscriberID.IMSI))
        sub_data = self._store.get_subscriber_data(sid)
        return sub_data

    def generate_m5g_auth_vector(self, imsi: str, snni: bytes):
        """
        Returns the m5g auth vector for the subscriber by querying the store
        for the crypto algo and secret keys.
        """
        sid = SIDUtils.to_str(SubscriberID(id=imsi, type=SubscriberID.IMSI))
        subs = self._store.get_subscriber_data(sid)

        if subs.lte.state != LTESubscription.ACTIVE:
            raise CryptoError("LTE service not active for %s" % sid)

        if subs.lte.auth_algo != LTESubscription.MILENAGE:
            raise CryptoError("Unknown crypto (%s) for %s" %
                              (subs.lte.auth_algo, sid))

        if len(subs.lte.auth_key) != 16:
            raise CryptoError("Subscriber key not valid for %s" % sid)

        if len(subs.lte.auth_opc) == 0:
            opc = Milenage.generate_opc(subs.lte.auth_key, self._op)
        elif len(subs.lte.auth_opc) != 16:
            raise CryptoError("Subscriber OPc is invalid length for %s" % sid)
        else:
            opc = subs.lte.auth_opc

        sqn = self.seq_to_sqn(self.get_next_lte_auth_seq(imsi))
        milenage = Milenage(self._amf)
        return milenage.generate_m5gran_vector(subs.lte.auth_key,
                                              opc, sqn, snni)

    @classmethod
    def seq_to_sqn(cls, seq, ind=0):
        """Compute the 48 bit SQN given a seq given the formula defined in
        3GPP TS 33.102 Annex C.3.2. The length of IND is 5 bits.
                                SQN = SEQ || IND

        Args:
            seq (int): the sequence number
            ind (int): the index value
        Returns:
            sqn (int): the 48bit SQN
        """
        return (seq << 5 & 0xFFFFFFFFFFE0) + (ind & 0x1F)

    @classmethod
    def sqn_to_seq(cls, sqn):
        """Compute the SEQ given a 48 bit SQN using the formula defined in
        3GPP TS 33.102 Annex C.3.2. The length of IND is 5 bits.
                                SQN = SEQ || IND

        Args:
            sqn (int): the 48bit SQN
        Returns:
            seq (int): the sequence number
        """
        return sqn >> 5

