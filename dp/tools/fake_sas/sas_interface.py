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
"""SAS interfaces for testing."""

from __future__ import absolute_import, division, print_function

import abc

import six


class SasInterface(six.with_metaclass(abc.ABCMeta, object)):
    """WinnForum standardized interfaces.

    Includes SAS-CBSD interface and (will include) SAS-SAS interface.

    """

    @abc.abstractmethod
    def Registration(self, request, ssl_cert=None, ssl_key=None):
        """SAS-CBSD Registration interface.

        Registers CBSDs.

        Request and response are both lists of dictionaries. Each dictionary
        contains all fields of a single request/response.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "registrationRequest" and the value is a list of individual CBSD
            registration requests (each of which is itself a dictionary).
          ssl_cert: Path to SSL cert file, if None, will use default cert file.
          ssl_key: Path to SSL key file, if None, will use default key file.
        Returns:
          A dictionary with a single key-value pair where the key is
          "registrationResponse" and the value is a list of individual CBSD
          registration responses (each of which is itself a dictionary).
        """
        pass

    @abc.abstractmethod
    def SpectrumInquiry(self, request, ssl_cert=None, ssl_key=None):
        """SAS-CBSD SpectrumInquiry interface.

        Performs spectrum inquiry for CBSDs.

        Request and response are both lists of dictionaries. Each dictionary
        contains all fields of a single request/response.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "spectrumInquiryRequest" and the value is a list of individual CBSD
            spectrum inquiry requests (each of which is itself a dictionary).
          ssl_cert: Path to SSL cert file, if None, will use default cert file.
          ssl_key: Path to SSL key file, if None, will use default key file.
        Returns:
          A dictionary with a single key-value pair where the key is
          "spectrumInquiryResponse" and the value is a list of individual CBSD
          spectrum inquiry responses (each of which is itself a dictionary).
        """
        pass

    @abc.abstractmethod
    def Grant(self, request, ssl_cert=None, ssl_key=None):
        """SAS-CBSD Grant interface.

        Request and response are both lists of dictionaries. Each dictionary
        contains all fields of a single request/response.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "grantRequest" and the value is a list of individual CBSD
            grant requests (each of which is itself a dictionary).
          ssl_cert: Path to SSL cert file, if None, will use default cert file.
          ssl_key: Path to SSL key file, if None, will use default key file.
        Returns:
          A dictionary with a single key-value pair where the key is
          "grantResponse" and the value is a list of individual CBSD
          grant responses (each of which is itself a dictionary).
        """
        pass

    @abc.abstractmethod
    def Heartbeat(self, request, ssl_cert=None, ssl_key=None):
        """SAS-CBSD Heartbeat interface.

        Requests heartbeat for a grant for CBSDs.

        Request and response are both lists of dictionaries. Each dictionary
        contains all fields of a single request/response.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "heartbeatRequest" and the value is a list of individual CBSD
            heartbeat requests (each of which is itself a dictionary).
          ssl_cert: Path to SSL cert file, if None, will use default cert file.
          ssl_key: Path to SSL key file, if None, will use default key file.
        Returns:
          A dictionary with a single key-value pair where the key is
          "heartbeatResponse" and the value is a list of individual CBSD
          heartbeat responses (each of which is itself a dictionary).
        """
        pass

    @abc.abstractmethod
    def Relinquishment(self, request, ssl_cert=None, ssl_key=None):
        """SAS-CBSD Relinquishment interface.

        Relinquishes grant for CBSDs.

        Request and response are both lists of dictionaries. Each dictionary
        contains all fields of a single request/response.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "relinquishmentRequest" and the value is a list of individual CBSD
            relinquishment requests (each of which is itself a dictionary).
          ssl_cert: Path to SSL cert file, if None, will use default cert file.
          ssl_key: Path to SSL key file, if None, will use default key file.
        Returns:
          A dictionary with a single key-value pair where the key is
          "relinquishmentResponse" and the value is a list of individual CBSD
          relinquishment responses (each of which is itself a dictionary).
        """
        pass

    @abc.abstractmethod
    def Deregistration(self, request, ssl_cert=None, ssl_key=None):
        """SAS-CBSD Deregistration interface.

        Deregisters CBSDs.

        Request and response are both lists of dictionaries. Each dictionary
        contains all fields of a single request/response.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "deregistrationRequest" and the value is a list of individual CBSD
            deregistration requests (each of which is itself a dictionary).
          ssl_cert: Path to SSL cert file, if None, will use default cert file.
          ssl_key: Path to SSL key file, if None, will use default key file.
        Returns:
          A dictionary with a single key-value pair where the key is
          "deregistrationResponse" and the value is a list of individual CBSD
          deregistration responses (each of which is itself a dictionary).
        """
        pass

    @abc.abstractmethod
    def GetEscSensorRecord(self, request, ssl_cert=None, ssl_key=None):
        """SAS-SAS ESC Sensor Record Exchange interface

        Requests a Pull Command to get the ESC Sensor Data Message

        Args:
          request: A string containing Esc Sensor Record Id
          ssl_cert: Path to SSL cert file, if None, will use default cert file.
          ssl_key: Path to SSL key file, if None, will use default key file.
        Returns:
          A dictionary of Esc Sensor Data Message object specified in
          WINNF-16-S-0096
        """
        pass

    @abc.abstractmethod
    def GetFullActivityDump(self, ssl_cert=None, ssl_key=None):
        """SAS-SAS Full Activity Dump interface.

        Requests a Pull Command to get Full Activity Dump Message.

        Args:
          ssl_cert: Path to SSL cert file, if None, will use default cert file.
          ssl_key: Path to SSL key file, if None, will use default key file.
        Returns:
          A dictionary containing the FullActivityDump object specified in WINNF-16-S-0096
        """
        pass

    @abc.abstractmethod
    def DownloadFile(self, url, ssl_cert=None, ssl_key=None):
        """SAS-SAS Get data from json files after generate the
         Full Activity Dump Message
        Returns:
         the message as an "json data" object specified in WINNF-16-S-0096
        """
        pass


class SasAdminInterface(six.with_metaclass(abc.ABCMeta, object)):
    """Minimal test control interface for the SAS under test."""

    @abc.abstractmethod
    def Reset(self):
        """SAS admin interface to reset the SAS between test cases."""
        pass

    @abc.abstractmethod
    def InjectFccId(self, request):
        """SAS admin interface to inject fcc id information into SAS under test.

        Args:
          request: A dictionary with the following key-value pairs:
            "fccId": (string) valid fccId to be injected into SAS under test
            "fccMaxEirp": (double) optional; default value of 47 dBm/10 MHz
        """
        pass

    @abc.abstractmethod
    def InjectUserId(self, request):
        """SAS admin interface to whitelist a user ID in the SAS under test.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "userId" and the value is a string of valid userId to be whitelisted by
            the SAS under test.
        """
        pass

    @abc.abstractmethod
    def BlacklistByFccId(self, request):
        """Inject an FCC ID which will be blacklisted by the SAS under test.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "fccId" and the value is the FCC ID (string) to be blacklisted.
        """
        pass

    @abc.abstractmethod
    def BlacklistByFccIdAndSerialNumber(self, request):
        """Inject an (FCC ID, serial number) pair which will be blacklisted by the
           SAS under test.

        Args:
          request: A dictionary with the following key-value pairs:
            "fccId": (string) blacklisted FCC ID
            "serialNumber": (string) blacklisted serial number
        """
        pass

    @abc.abstractmethod
    def PreloadRegistrationData(self, request):
        """SAS admin interface to preload registration data into SAS under test.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "registrationData" and the value is a list of individual CBSD
            registration data which need to be preloaded into SAS (each of which is
            itself a dictionary). The dictionary is a RegistrationRequest object,
            the fccId and cbsdSerialNumber fields are required, other fields are
            optional.
        """
        pass

    @abc.abstractmethod
    def InjectExclusionZone(self, request):
        """Inject exclusion zone information into SAS under test.

        Args:
          request: A dictionary with the following key-value pairs:
            "zone": A GeoJSON object defining the exclusion zone to be injected to SAS UUT.
            "frequencyRanges": A list of frequency ranges for the exclusion zone.
        """
        pass

    @abc.abstractmethod
    def InjectZoneData(self, request):
        """Inject PPA or NTIA zone information into SAS under test.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "record" and the value is ZoneData object to be injected into
            SAS under test. For more information about ZoneData please see
            the SAS-SAS TS (WINNF-16-S-0096).
        """
        pass

    @abc.abstractmethod
    def InjectPalDatabaseRecord(self, request):
        """Inject a PAL Database record into the SAS under test.

        Args:
          request:
          For the contents of this request, please refer to the PAL Database TS
          (WINNF-16-S-0245).
        """
        pass

    @abc.abstractmethod
    def InjectFss(self, request):
        """SAS admin interface to inject FSS information into SAS under test.

        Args:
            request: A dictionary with a single key-value pair where the key is
            "record" and the value is a fixed satellite service object
            (which is itself a dictionary). The dictionary is an
            IncumbentProtectionData object (specified in SAS-SAS TS).
        """
        pass

    @abc.abstractmethod
    def InjectWisp(self, request):
        """SAS admin interface to inject WISP information into SAS under test.

        Args:
          request: A dictionary with two key-value pairs where the keys are
            "record" and "zone" with the values IncumbentProtectionData
            object (specified in SAS-SAS TS) and a GeoJSON Object respectively
        Note: Required Field in IncumbentProtectionData are id, type,
        deploymentParam->operationParam->operationFrequencyRange->
        lowFrequency, highFrequency
        """
        pass

    @abc.abstractmethod
    def InjectSasAdministratorRecord(self, request):
        """SAS admin interface to inject SAS Administrator Record into SAS under test.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "record" and the value is a SAS Administrator information (which is
            itself a dictionary). The dictionary is an SASAdministrator object
            (Specified in SAS-SAS TS WINNF-16-S-0096)
        """
        pass

    @abc.abstractmethod
    def InjectEscSensorDataRecord(self, request):
        """SAS admin interface to inject ESC Sensor Data Record into SAS under test.

        Args:
          request: A dictionary with a single key-value pair where the key is
            "record" and the value is a EscSensorData object (which is
            itself a dictionary specified in SAS-SAS TS WINNF-16-S-0096)
        Behavior: SAS should act as if it is connected to an ESC sensor with
        the provided parameters.
        """
        pass

    @abc.abstractmethod
    def TriggerMeasurementReportRegistration(self):
        """SAS admin interface to trigger measurement report request for all subsequent
        registration request

        Note: The SAS should request a measurement report in the RegistrationResponse
        (if status == 0)
        """
        pass

    @abc.abstractmethod
    def TriggerMeasurementReportHeartbeat(self):
        """SAS admin interface to trigger measurement report request for all subsequent
        heartbeat request

        Note: The SAS should request a measurement report in the HeartbeatResponse
        (if status == 0)
        """
        pass

    @abc.abstractmethod
    def TriggerPpaCreation(self, request):
        """SAS admin interface to trigger PPA creation based on the CBSD Ids,
        Pal Ids and Provided Contour

        Args:
          request: A dictionary with multiple key-value pairs where the keys are
            cbsdIds: array of string containing CBSD Id
            palIds: array of string containing PAL Id
            providedContour(optional): GeoJSON Object

        Returns:
          PPA Id in string format
        """
        pass

    @abc.abstractmethod
    def TriggerDailyActivitiesImmediately(self):
        """SAS admin interface to trigger daily activities immediately which will
        execute the following activities:
          1. Pull from all External Database and other SASes (URLs will be injected to
          SAS UUT using another RPC Call)
          2. Run IAP and DPA Calculations
          3. Apply EIRP updates to devices
        """
        pass

    @abc.abstractmethod
    def TriggerEnableNtiaExclusionZones(self):
        """SAS admin interface to trigger enforcement of the NTIA exclusion zones 
        """
        pass

    @abc.abstractmethod
    def GetDailyActivitiesStatus(self):
        """SAS admin interface to get the daily activities status
        Returns:
          A dictionary with a single key-value pair where the key is "completed" and the
          value is a boolean with value as true if the daily activities is completed and
          false if the daily activities is running/failing.
        """
        pass

    @abc.abstractmethod
    def TriggerEnableScheduledDailyActivities(self):
        """SAS admin interface to trigger the daily activities according to the 
           schedule agreed upon by SAS admins.                         
        """
        pass

    @abc.abstractmethod
    def InjectCpiUser(self, request):
        """SAS admin interface to add a CPI User as if it came directly from the CPI database.

        Args:
          request: A dictionary with the following key-value pairs:
            "cpiId": (string) valid cpiId to be injected into SAS under test
            "cpiName": (string) valid name for cpi user to be injected into SAS under test
            "cpiPublicKey": (string) public key value for cpi user to be injected into SAS under test
        """
        pass

    @abc.abstractmethod
    def TriggerLoadDpas(self):
        """SAS admin interface to load all ESC-monitored DPAs and immediately activate all of them.
        """
        pass

    @abc.abstractmethod
    def TriggerBulkDpaActivation(self, request):
        """SAS admin interface to bulk DPA activation/deactivation
         Args:
          request: A dictionary with the following key-value pairs:
            "activate": (boolean) if True, activate all ESC-monitored DPAs on all channels
                else deactivate all ESC-monitored DPAs on all channels
        """
        pass

    @abc.abstractmethod
    def TriggerDpaActivation(self, request):
        """SAS admin interface to activate specific DPA on specific channel
         Args:
          request: A dictionary with the following key-value pairs:
            "dpaId": (string) it represents the field "name" in the kml file of DPAs
            "frequencyRange": frequencyRange of DPA Channel with lowFrequency, highFrequency

        """
        pass

    @abc.abstractmethod
    def TriggerDpaDeactivation(self, request):
        """SAS admin interface to deactivate specific DPA on specific channel
         Args:
          request: A dictionary with the following key-value pairs:
            "dpaId": (string) it represents the field "name" in the kml file of DPAs
            "frequencyRange": frequencyRange of DPA Channel with lowFrequency, highFrequency
        """
        pass

    @abc.abstractmethod
    def TriggerEscDisconnect(self):
        """Simulates the ESC (ESC-DE) being disconnected from the SAS UUT."""
        pass

    @abc.abstractmethod
    def TriggerFullActivityDump(self):
        """SAS admin interface to trigger generation of a Full Activity Dump.

        Note: SAS does not need to complete generation before returning HTTP 200.
        See the testing API specification for more details.
        """
        pass

    @abc.abstractmethod
    def InjectPeerSas(self, request):
        """SAS admin interface to inject a peer SAS into the SAS UUT.

        Args:
          request: A dictionary with the following key-value pairs:
            "certificateHash": the sha1 fingerprint of the certificate
            "url": base URL of the peer SAS.
        """
        pass

    @abc.abstractmethod
    def QueryPropagationAndAntennaModel(self, request):
        """SAS admin interface to query propagation and antenna gains for CBSD and FSS  or Provided PPA Contour

        Args:
          request: A dictionary with multiple key-value pairs where the keys are
            reliabilityLevel: (permitted values: -1, 0.05, 0.95)
            cbsd: dictionary defining cbsd
            fss(optional): dictionary defining fss
            ppa(optional): GeoJSON Object

        Returns:
          double pathlossDb (pathloss in dB)
          double txAntennaGainDbi (transmitter antenna gain in dBi in the direction of the receiver)
          double rxAntennaGainDbi (optional) (receiver antenna gain in dBi in the direction of the transmitter)

        """
        pass

    @abc.abstractmethod
    def GetPpaCreationStatus(self):
        """SAS admin interface to get the most recent PPA creation status
        Returns:
          A dictionary with a two key-value pairs where the keys are "completed" and
          "withError". The values are of boolean type. The value for "completed" flag
          set to True if the ppa creation(for the most recent ppa creation) has completed or
          to False if the PPA creation is still in progress. The value for "withError" flag is
          is set to True if the PPA creation has completed with error(s) else it is set to False.
        """
        pass

    @abc.abstractmethod
    def InjectDatabaseUrl(self, request):
        """Inject the Database URL into SAS.

        Args:
          request: Contains database url to be injected.
        """
        pass
