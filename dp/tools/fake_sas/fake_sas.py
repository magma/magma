# flake8: noqa
#    Copyright 2018 SAS Project Authors. All Rights Reserved.
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.


# Some parts of this software was developed by employees of
# the National Institute of Standards and Technology (NIST),
# an agency of the Federal Government.
# Pursuant to title 17 United States Code Section 105, works of NIST employees
# are not subject to copyright protection in the United States and are
# considered to be in the public domain. Permission to freely use, copy,
# modify, and distribute this software and its documentation without fee
# is hereby granted, provided that this notice and disclaimer of warranty
# appears in all copies.

# THE SOFTWARE IS PROVIDED 'AS IS' WITHOUT ANY WARRANTY OF ANY KIND, EITHER
# EXPRESSED, IMPLIED, OR STATUTORY, INCLUDING, BUT NOT LIMITED TO, ANY WARRANTY
# THAT THE SOFTWARE WILL CONFORM TO SPECIFICATIONS, ANY IMPLIED WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND FREEDOM FROM
# INFRINGEMENT, AND ANY WARRANTY THAT THE DOCUMENTATION WILL CONFORM TO THE
# SOFTWARE, OR ANY WARRANTY THAT THE SOFTWARE WILL BE ERROR FREE. IN NO EVENT
# SHALL NIST BE LIABLE FOR ANY DAMAGES, INCLUDING, BUT NOT LIMITED TO, DIRECT,
# INDIRECT, SPECIAL OR CONSEQUENTIAL DAMAGES, ARISING OUT OF, RESULTING FROM,
# OR IN ANY WAY CONNECTED WITH THIS SOFTWARE, WHETHER OR NOT BASED UPON
# WARRANTY, CONTRACT, TORT, OR OTHERWISE, WHETHER OR NOT INJURY WAS SUSTAINED
# BY PERSONS OR PROPERTY OR OTHERWISE, AND WHETHER OR NOT LOSS WAS SUSTAINED
# FROM, OR AROSE OUT OF THE RESULTS OF, OR USE OF, THE SOFTWARE OR SERVICES
# PROVIDED HEREUNDER.

# Distributions of NIST software should also include copyright and licensing
# statements of any third-party software that are legally bundled with the
# code in compliance with the conditions of those licenses.

"""A fake implementation of SasInterface, based on v1.0 of the SAS-CBSD TS.

A local test server could be run by using "python fake_sas.py".

"""
from __future__ import absolute_import, division, print_function

import argparse
import json
import os
import ssl
import sys
import tempfile
import uuid
from datetime import datetime, timedelta

import sas_interface
from OpenSSL import crypto
from six.moves import configparser
from six.moves.BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer

# Fake SAS server configurations.
HOST = '0.0.0.0'
PORT = 9000
CERT_FILE = 'certs/server.cert'
KEY_FILE = 'certs/server.key'
CA_CERT = 'certs/ca.cert'
CIPHERS = [
    'AES128-GCM-SHA256',              # TLS_RSA_WITH_AES_128_GCM_SHA256
    'AES256-GCM-SHA384',              # TLS_RSA_WITH_AES_256_GCM_SHA384
    'ECDHE-RSA-AES128-GCM-SHA256',    # TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
]
ECC_CERT_FILE = 'certs/server-ecc.cert'
ECC_KEY_FILE = 'certs/server-ecc.key'
ECC_CIPHERS = [
    'ECDHE-ECDSA-AES128-GCM-SHA256',  # TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
    'ECDHE-ECDSA-AES256-GCM-SHA384',  # TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
]

MISSING_PARAM = 102
INVALID_PARAM = 103


class FakeSas(sas_interface.SasInterface):
    """A fake implementation of SasInterface.

    Returns success for all requests with plausible fake values for all required
    response fields.

    """

    def __init__(self):
        self.maximum_batch_size = 100
        pass

    def Registration(self, request, ssl_cert=None, ssl_key=None):
        response = {'registrationResponse': []}
        for req in request['registrationRequest']:
            if 'fccId' not in req or 'cbsdSerialNumber' not in req:
                response['registrationResponse'].append({
                    'response': self._GetSuccessResponse(),
                })
                continue
            response['registrationResponse'].append({
                'cbsdId': req['fccId'] + '/' + req['cbsdSerialNumber'],
                'response': self._GetSuccessResponse(),
            })
        return response

    def SpectrumInquiry(self, request, ssl_cert=None, ssl_key=None):
        response = {'spectrumInquiryResponse': []}
        for req in request['spectrumInquiryRequest']:
            response['spectrumInquiryResponse'].append({
                'cbsdId': req['cbsdId'],
                'availableChannel': [{
                    'frequencyRange': {
                        'lowFrequency': 3550000000,
                        'highFrequency': 3700000000,
                    },
                    'channelType': 'GAA',
                    'ruleApplied': 'FCC_PART_96',
                    'maxEirp': 37.0,
                }],
                'response': self._GetSuccessResponse(),
            })
        return response

    def Grant(self, request, ssl_cert=None, ssl_key=None):
        response = {'grantResponse': []}
        for req in request['grantRequest']:
            if ('cbsdId' not in req):
                response['grantResponse'].append({
                    'response': self._GetMissingParamResponse(),
                })
            else:
                if (('highFrequency' not in req['operationParam']['operationFrequencyRange']) or \
                   ('lowFrequency' not in req['operationParam']['operationFrequencyRange'])):
                    response['grantResponse'].append({
                        'cbsdId': req['cbsdId'],
                        'response': self._GetMissingParamResponse(),
                    })
                else:
                    response['grantResponse'].append({
                        'cbsdId': req['cbsdId'],
                        'grantId': 'fake_grant_id_%s' % datetime.utcnow().isoformat(),
                        'channelType': 'GAA',
                        'heartbeatInterval': 60,
                        'response': self._GetSuccessResponse(),
                    })
        return response

    def Heartbeat(self, request, ssl_cert=None, ssl_key=None):
        response = {'heartbeatResponse': []}
        for req in request['heartbeatRequest']:
            transmit_expire_time = datetime.utcnow().replace(
                microsecond=0,
            ) + timedelta(minutes=1)
            response['heartbeatResponse'].append({
                'cbsdId': req['cbsdId'],
                'grantId': req['grantId'],
                'transmitExpireTime': transmit_expire_time.isoformat() + 'Z',
                'response': self._GetSuccessResponse(),
            })
        return response

    def Relinquishment(self, request, ssl_cert=None, ssl_key=None):
        response = {'relinquishmentResponse': []}
        for req in request['relinquishmentRequest']:
            response['relinquishmentResponse'].append({
                'cbsdId': req['cbsdId'],
                'grantId': req['grantId'],
                'response': self._GetSuccessResponse(),
            })
        return response

    def Deregistration(self, request, ssl_cert=None, ssl_key=None):
        response = {'deregistrationResponse': []}
        for req in request['deregistrationRequest']:
            if ('cbsdId' not in req):
                response['deregistrationResponse'].append({
                    'response': self._GetMissingParamResponse(),
                })
            else:
                response['deregistrationResponse'].append({
                    'cbsdId': req['cbsdId'],
                    'response': self._GetSuccessResponse(),
                })
        return response

    def GetEscSensorRecord(self, request, ssl_cert=None, ssl_key=None):
        # Get the Esc Sensor record
        with open(os.path.join('testcases', 'testdata', 'esc_sensor_record_0.json')) as fd:
            esc_sensor_record = json.load(fd)

        if request == esc_sensor_record['id']:
            return esc_sensor_record
        else:
            # Return Empty if invalid Id
            return {}

    def GetFullActivityDump(self, version, ssl_cert=None, ssl_key=None):
        response = json.loads(
            json.dumps({
                'files': [
                    {
                        'url': 'https://raw.githubusercontent.com/Wireless-Innovation-Forum/' +
                        'Spectrum-Access-System/master/schema/empty_activity_dump_file.json',
                        'checksum': 'da39a3ee5e6b4b0d3255bfef95601890afd80709', 'size': 19,
                        'version': version, 'recordType': "cbsd",
                    },
                    {
                        'url': 'https://raw.githubusercontent.com/Wireless-Innovation-Forum/Spectrum-Access-System/master/schema/empty_activity_dump_file.json',
                        'checksum': 'da39a3ee5e6b4b0d3255bfef95601890afd80709', 'size': 19,
                        'version': version, 'recordType': "zone",
                    },
                    {
                        'url': 'https://raw.githubusercontent.com/Wireless-Innovation-Forum/Spectrum-Access-System/master/schema/empty_activity_dump_file.json',
                        'checksum': 'da39a3ee5e6b4b0d3255bfef95601890afd80709', 'size': 19,
                        'version': version, 'recordType': "esc_sensor",
                    },
                    {
                        'url': 'https://raw.githubusercontent.com/Wireless-Innovation-Forum/Spectrum-Access-System/master/schema/empty_activity_dump_file.json',
                        'checksum': 'da39a3ee5e6b4b0d3255bfef95601890afd80709', 'size': 19,
                        'version': version, 'recordType': "coordination",
                    },
                ],
                'generationDateTime': datetime.utcnow().strftime(
                    '%Y-%m-%dT%H:%M:%SZ',
                ),
                'description': "Full activity dump files",
            }),
        )
        return response

    def _GetSuccessResponse(self):
        return {'responseCode': 0}

    def _GetMissingParamResponse(self):
        return {'responseCode': MISSING_PARAM}

    def DownloadFile(self, url, ssl_cert=None, ssl_key=None):
        """SAS-SAS Get data from json files after generate the
         Full Activity Dump Message
        Returns:
         the message as an "json data" object specified in WINNF-16-S-0096
        """
        pass


class FakeSasAdmin(sas_interface.SasAdminInterface):
    """Implementation of SAS Admin for Fake SAS."""

    def Reset(self):
        pass

    def InjectFccId(self, request):
        pass

    def InjectUserId(self, request):
        pass

    def InjectCpiUser(self, request):
        pass

    def BlacklistByFccId(self, request):
        pass

    def BlacklistByFccIdAndSerialNumber(self, request):
        pass

    def PreloadRegistrationData(self, request):
        pass

    def InjectExclusionZone(self, request, ssl_cert=None, ssl_key=None):
        pass

    def InjectZoneData(self, request, ssl_cert=None, ssl_key=None):
        return request['record']['id']

    def InjectPalDatabaseRecord(self, request):
        pass

    def InjectFss(self, request):
        pass

    def InjectWisp(self, request):
        pass

    def InjectSasAdministratorRecord(self, request):
        pass

    def InjectEscSensorDataRecord(self, request):
        pass

    def InjectPeerSas(self, request):
        pass

    def TriggerMeasurementReportRegistration(self):
        pass

    def TriggerMeasurementReportHeartbeat(self):
        pass

    def TriggerPpaCreation(self, request, ssl_cert=None, ssl_key=None):
        return 'zone/ppa/fake_sas/%s/%s' % (
            request['palIds'][0],
            uuid.uuid4().hex,
        )

    def TriggerDailyActivitiesImmediately(self):
        pass

    def TriggerEnableNtiaExclusionZones(self):
        pass

    def TriggerEnableScheduledDailyActivities(self):
        pass

    def QueryPropagationAndAntennaModel(self, request):
        from testcases.WINNF_FT_S_PAT_testcase import (
            computePropagationAntennaModel,
        )
        return computePropagationAntennaModel(request)

    def GetDailyActivitiesStatus(self):
        return {'completed': True}

    def GetPpaCreationStatus(self):
        return {'completed': True, 'withError': False}

    def GetDailyActivitiesStatus(self):
        return {'completed': True}

    def TriggerLoadDpas(self):
        pass

    def TriggerBulkDpaActivation(self, request):
        pass

    def TriggerDpaActivation(self, request):
        pass

    def TriggerFullActivityDump(self):
        pass

    def TriggerDpaDeactivation(self, request):
        pass

    def TriggerEscDisconnect(self):
        pass

    def InjectDatabaseUrl(self, request):
        pass


class FakeSasHandler(BaseHTTPRequestHandler):
    @classmethod
    def SetVersion(cls, cbsd_sas_version, sas_sas_version):
        cls.cbsd_sas_version = cbsd_sas_version
        cls.sas_sas_version = sas_sas_version

    def _parseUrl(self, url):
        """Parse the Url into the path and value."""
        splitted_url = url.split('/')[1:]
        # Returns path and value
        return '/'.join(splitted_url[0:2]), '/'.join(splitted_url[2:])

    def do_POST(self):
        """Handles POST requests."""

        length = int(self.headers.getheader('content-length'))
        if length > 0:
            request = json.loads(self.rfile.read(length))
        if self.path == '/%s/registration' % self.cbsd_sas_version:
            response = FakeSas().Registration(request)
        elif self.path == '/%s/spectrumInquiry' % self.cbsd_sas_version:
            response = FakeSas().SpectrumInquiry(request)
        elif self.path == '/%s/grant' % self.cbsd_sas_version:
            response = FakeSas().Grant(request)
        elif self.path == '/%s/heartbeat' % self.cbsd_sas_version:
            response = FakeSas().Heartbeat(request)
        elif self.path == '/%s/relinquishment' % self.cbsd_sas_version:
            response = FakeSas().Relinquishment(request)
        elif self.path == '/%s/deregistration' % self.cbsd_sas_version:
            response = FakeSas().Deregistration(request)
        elif self.path == '/admin/injectdata/zone':
            response = FakeSasAdmin().InjectZoneData(request)
        elif self.path == '/admin/trigger/create_ppa':
            response = FakeSasAdmin().TriggerPpaCreation(request)
        elif self.path == '/admin/get_daily_activities_status':
            response = FakeSasAdmin().GetDailyActivitiesStatus()
        elif self.path == '/admin/get_daily_activities_status':
            response = FakeSasAdmin().GetDailyActivitiesStatus()
        elif self.path == '/admin/get_ppa_status':
            response = FakeSasAdmin().GetPpaCreationStatus()
        elif self.path == '/admin/query/propagation_and_antenna_model':
            try:
                response = FakeSasAdmin().QueryPropagationAndAntennaModel(request)
            except ValueError:
                self.send_response(400)
            return
        elif self.path in (
            '/admin/reset', '/admin/injectdata/fcc_id',
            '/admin/injectdata/user_id',
            '/admin/injectdata/conditional_registration',
            '/admin/injectdata/blacklist_fcc_id',
            '/admin/injectdata/blacklist_fcc_id_and_serial_number',
            '/admin/injectdata/fss',
            '/admin/injectdata/wisp',
            '/admin/injectdata/peer_sas',
            '/admin/injectdata/pal_database_record',
            '/admin/injectdata/sas_admin',
            '/admin/injectdata/sas_impl',
            '/admin/injectdata/esc_sensor',
            '/admin/injectdata/cpi_user',
            '/admin/trigger/meas_report_in_registration_response',
            '/admin/trigger/meas_report_in_heartbeat_response',
            '/admin/trigger/daily_activities_immediately',
            '/admin/trigger/enable_scheduled_daily_activities',
            '/admin/trigger/load_dpas',
            '/admin/trigger/dpa_activation',
            '/admin/trigger/dpa_deactivation',
            '/admin/trigger/bulk_dpa_activation',
            '/admin/injectdata/exclusion_zone',
            '/admin/trigger/create_full_activity_dump',
            '/admin/injectdata/database_url',
        ):
            response = ''
        else:
            self.send_response(404)
            return
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(response))

    def do_GET(self):
        """Handles GET requests."""
        path, value = self._parseUrl(self.path)
        if path == '%s/esc_sensor' % self.sas_sas_version:
            response = FakeSas().GetEscSensorRecord(value)
        elif path == '%s/dump' % self.sas_sas_version:
            response = FakeSas().GetFullActivityDump(self.sas_sas_version)
        else:
            self.send_response(404)
            return
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(response))


def ParseCrlIndex(index_filename):
    """Returns the list of CRL filenames from a CRL index file."""
    dirname = os.path.dirname(index_filename)
    return [
        os.path.join(dirname, line.rstrip())
        for line in open(index_filename)
    ]


def RunFakeServer(cbsd_sas_version, sas_sas_version, is_ecc, crl_index):
    FakeSasHandler.SetVersion(cbsd_sas_version, sas_sas_version)
    if is_ecc:
        assert ssl.HAS_ECDH
    server = HTTPServer((HOST, PORT), FakeSasHandler)

    ssl_context = ssl.create_default_context(ssl.Purpose.CLIENT_AUTH)
    ssl_context.options |= ssl.CERT_REQUIRED

    # If CRLs were provided, then enable revocation checking.
    if crl_index is not None:
        ssl_context.verify_flags = ssl.VERIFY_CRL_CHECK_CHAIN

        try:
            crl_files = ParseCrlIndex(crl_index)
        except IOError as e:
            print("Failed to parse CRL index file %r: %s" % (crl_index, e))

        # https://tools.ietf.org/html/rfc5280#section-4.2.1.13 specifies that
        # CRLs MUST be DER-encoded, but SSLContext expects the name of a PEM-encoded
        # file, so we must convert it first.
        for f in crl_files:
            try:
                with file(f) as handle:
                    der = handle.read()
                    try:
                        crl = crypto.load_crl(crypto.FILETYPE_ASN1, der)
                    except crypto.Error as e:
                        print("Failed to parse CRL file %r as DER format: %s" % (f, e))
                        return
                    with tempfile.NamedTemporaryFile() as tmp:
                        tmp.write(crypto.dump_crl(crypto.FILETYPE_PEM, crl))
                        tmp.flush()
                        ssl_context.load_verify_locations(cafile=tmp.name)
                print("Loaded CRL file: %r" % f)
            except IOError as e:
                print("Failed to load CRL file %r: %s" % (f, e))
                return

    ssl_context.load_verify_locations(cafile=CA_CERT)
    ssl_context.load_cert_chain(
        certfile=ECC_CERT_FILE if is_ecc else CERT_FILE,
        keyfile=ECC_KEY_FILE if is_ecc else KEY_FILE,
    )
    ssl_context.set_ciphers(':'.join(ECC_CIPHERS if is_ecc else CIPHERS))
    ssl_context.verify_mode = ssl.CERT_REQUIRED
    server.socket = ssl_context.wrap_socket(server.socket, server_side=True)
    print('Will start server at %s:%d, use <Ctrl-C> to stop.' % (HOST, PORT))
    server.serve_forever()


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument(
        '--ecc', help='Use ECDSA certificate', action='store_true',
    )
    parser.add_argument(
        '--crl_index',
        help=(
            'Read a text file containing one DER-encoded CRL file per line, '
            'and enable revocation checking.'
        ),
        dest='crl_index', action='store',
    )
    try:
        args = parser.parse_args()
    except:
        parser.print_help()
        sys.exit(0)
    config_parser = configparser.RawConfigParser()
    config_parser.read(['sas.cfg'])
    cbsd_sas_version = config_parser.get('SasConfig', 'CbsdSasVersion')
    sas_sas_version = config_parser.get('SasConfig', 'SasSasVersion')
    RunFakeServer(cbsd_sas_version, sas_sas_version, args.ecc, args.crl_index)
