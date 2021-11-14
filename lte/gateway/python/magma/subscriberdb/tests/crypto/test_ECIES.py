# −*− coding: UTF−8 −*−
#/**
# * Software Name : CryptoMobile 
# * Version : 0.3
# *
# * Copyright 2020. Benoit Michau. P1Sec.
# *
# * This program is free software: you can redistribute it and/or modify
# * it under the terms of the GNU General Public License version 2 as published
# * by the Free Software Foundation. 
# *
# * This program is distributed in the hope that it will be useful,
# * but WITHOUT ANY WARRANTY; without even the implied warranty of
# * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# * GNU General Public License for more details. 
# *
# * You will find a copy of the terms and conditions of the GNU General Public
# * License version 2 in the "license.txt" file or
# * see http://www.gnu.org/licenses/ or write to the Free Software Foundation,
# * Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301 USA
# *
# *--------------------------------------------------------
# * File Name : test/test_ECIES.py
# * Created : 2020-01-22
# * Authors : Benoit Michau 
# *--------------------------------------------------------
#*/

########################################################
# CryptoMobile python toolkit
#
# ECIES computation
# as defined in 3GPP TS 33.501, annex C
#######################################################

import unittest

from time     import time
from binascii import unhexlify

from cryptography.hazmat.primitives.asymmetric.x25519   import X25519PrivateKey
from magma.subscriberdb.crypto.EC import *
from magma.subscriberdb.crypto.ECIES import *


# annex C.4.3, ECIES Profile A test data
def test_profileA():
    hn_privkey  = unhexlify('c53c22208b61860b06c62e5406a7b330c2b577aa5558981510d128247d38bd1d')
    hn_pubkey   = unhexlify('5a8d38864820197c3394b92613b20b91633cbd897119273bf8e4a6f4eec0a650')
    eph_privkey = unhexlify('c80949f13ebe61af4ebdbd293ea4f942696b9e815d7e8f0096bbf6ed7de62256')
    eph_pubkey  = unhexlify('b2e92f836055a255837debf850b528997ce0201cb82adfe4be1f587d07d8457d')
    shared_key  = unhexlify('028ddf890ec83cdf163947ce45f6ec1a0e3070ea5fe57e2b1f05139f3e82422a')
    #
    plaintext   = unhexlify('00012080f6')
    ciphertext  = unhexlify('cb02352410')
    mactag      = unhexlify('cddd9e730ef3fa87')
    
    x1 = X25519(eph_privkey)
    x2 = X25519(hn_privkey)
    
    ue = ECIES_UE(profile='A')
    ue.EC.PrivKey = X25519PrivateKey.from_private_bytes(eph_privkey)
    hn = ECIES_HN(profile='A', hn_priv_key=hn_privkey)
    ue.generate_sharedkey(hn_pubkey, fresh=False)
    ue_pk, ue_ct, ue_mac = ue.protect(plaintext)
    hn_ct = hn.unprotect(ue_pk, ue_ct, ue_mac)
    
    return x1.get_pubkey() == eph_pubkey and \
    x1.generate_sharedkey(hn_pubkey) == shared_key and \
    x2.get_pubkey() == hn_pubkey and \
    x2.generate_sharedkey(eph_pubkey) == shared_key and \
    ue_ct == ciphertext and ue_mac == mactag and hn_ct == plaintext


# annex C.4.4, ECIES Profile A test data
def test_profileB():
    hn_pubkey   = unhexlify('0272DA71976234CE833A6907425867B82E074D44EF907DFB4B3E21C1C2256EBCD1') # compressed
    hn_privkey  = unhexlify('F1AB1074477EBCC7F554EA1C5FC368B1616730155E0041AC447D6301975FECDA')
    eph_pubkey  = unhexlify('039AAB8376597021E855679A9778EA0B67396E68C66DF32C0F41E9ACCA2DA9B9D1') # compressed
    eph_privkey = unhexlify('99798858A1DC6A2C68637149A4B1DBFD1FDFF5ADDD62A2142F06699ED7602529')
    shared_key  = unhexlify('6C7E6518980025B982FBB2FF746E3C2E85A196D252099A7AD23EA7B4C0959CAE')
    #
    plaintext   = unhexlify('00012080F6')
    ciphertext  = unhexlify('46A33FC271')
    mactag      = unhexlify('6AC7DAE96AA30A4D')
    
    x1 = ECDH_SECP256R1(_raw_keypair=(eph_privkey, eph_pubkey))
    x2 = ECDH_SECP256R1(_raw_keypair=(hn_privkey, hn_pubkey))
    
    ue = ECIES_UE(profile='B')
    ue.EC._load_raw_keypair(eph_privkey, eph_pubkey)
    hn = ECIES_HN(None, profile='B', _raw_keypair=(hn_privkey, hn_pubkey))
    ue.generate_sharedkey(hn_pubkey, fresh=False)
    ue_pk, ue_ct, ue_mac = ue.protect(plaintext)
    hn_ct = hn.unprotect(ue_pk, ue_ct, ue_mac)
    
    return x1.generate_sharedkey(hn_pubkey) == shared_key and \
    x2.generate_sharedkey(eph_pubkey) == shared_key and \
    ue_ct == ciphertext and ue_mac == mactag and hn_ct == plaintext


def testall():
    return test_profileA() & test_profileB()


def testperf():
    a = None
    T0 = time()
    for i in range(1000):
        a = testall()
    print('1000 full testsets in %.3f seconds' % (time()-T0, ))


def test_ECIES():
    assert( testall() )


if __name__ == '__main__':
    testperf()
