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

import asyncio
from xml.etree import ElementTree

from aiohttp import web
from magma.common.misc_utils import get_ip_from_if
from magma.configuration.service_configs import load_service_config
from magma.enodebd.enodeb_status import get_enb_status, update_status_metrics
from magma.enodebd.logger import EnodebdLogger as logger
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_manager import StateMachineManager

from . import metrics


class StatsManager:
    """ HTTP server to receive performance management uploads from eNodeB and
        translate to metrics """
    # Dict to map performance counter names (from eNodeB) to metrics
    # For eNodeB sub-counters, the counter name is shown as
    # '<counter>:<sub-counter>'
    PM_FILE_TO_METRIC_MAP = {
        'RRC.AttConnEstab': metrics.STAT_RRC_ESTAB_ATT,
        'RRC.SuccConnEstab': metrics.STAT_RRC_ESTAB_SUCC,
        'RRC.AttConnReestab': metrics.STAT_RRC_REESTAB_ATT,
        'RRC.AttConnReestab._Cause:RRC.AttConnReestab.RECONF_FAIL':
            metrics.STAT_RRC_REESTAB_ATT_RECONF_FAIL,
        'RRC.AttConnReestab._Cause:RRC.AttConnReestab.HO_FAIL':
            metrics.STAT_RRC_REESTAB_ATT_HO_FAIL,
        'RRC.AttConnReestab._Cause:RRC.AttConnReestab.OTHER':
            metrics.STAT_RRC_REESTAB_ATT_OTHER,
        'RRC.SuccConnReestab': metrics.STAT_RRC_REESTAB_SUCC,
        'ERAB.NbrAttEstab': metrics.STAT_ERAB_ESTAB_ATT,
        'ERAB.NbrSuccEstab': metrics.STAT_ERAB_ESTAB_SUCC,
        'ERAB.NbrFailEstab': metrics.STAT_ERAB_ESTAB_FAIL,
        'ERAB.NbrReqRelEnb': metrics.STAT_ERAB_REL_REQ,
        'ERAB.NbrReqRelEnb.CauseUserInactivity':
            metrics.STAT_ERAB_REL_REQ_USER_INAC,
        'ERAB.NbrReqRelEnb.Normal': metrics.STAT_ERAB_REL_REQ_NORMAL,
        'ERAB.NbrReqRelEnb._Cause:ERAB.NbrReqRelEnb.CauseRADIORESOURCESNOTAVAILABLE':
            metrics.STAT_ERAB_REL_REQ_RES_NOT_AVAIL,
        'ERAB.NbrReqRelEnb._Cause:ERAB.NbrReqRelEnb.CauseREDUCELOADINSERVINGCELL':
            metrics.STAT_ERAB_REL_REQ_REDUCE_LOAD,
        'ERAB.NbrReqRelEnb._Cause:ERAB.NbrReqRelEnb.CauseFAILUREINTHERADIOINTERFACEPROCEDURE':
            metrics.STAT_ERAB_REL_REQ_FAIL_IN_RADIO,
        'ERAB.NbrReqRelEnb._Cause:ERAB.NbrReqRelEnb.CauseRELEASEDUETOEUTRANGENERATEDREASONS':
            metrics.STAT_ERAB_REL_REQ_EUTRAN_REAS,
        'ERAB.NbrReqRelEnb._Cause:ERAB.NbrReqRelEnb.CauseRADIOCONNECTIONWITHUELOST':
            metrics.STAT_ERAB_REL_REQ_RADIO_CONN_LOST,
        'ERAB.NbrReqRelEnb._Cause:ERAB.NbrReqRelEnb.CauseOAMINTERVENTION':
            metrics.STAT_ERAB_REL_REQ_OAM_INTV,
        'PDCP.UpOctUl': metrics.STAT_PDCP_USER_PLANE_BYTES_UL,
        'PDCP.UpOctDl': metrics.STAT_PDCP_USER_PLANE_BYTES_DL,
    }

    # Check if radio transmit is turned on every 10 seconds.
    CHECK_RF_TX_PERIOD = 10

    def __init__(self, enb_acs_manager: StateMachineManager):
        self.enb_manager = enb_acs_manager
        self.loop = asyncio.get_event_loop()
        self._prev_rf_tx = False
        self.mme_timeout_handler = None

    def run(self) -> None:
        """ Create and start HTTP server """
        svc_config = load_service_config("enodebd")

        app = web.Application()
        app.router.add_route('POST', "/{something}", self._post_handler)

        handler = app.make_handler()
        create_server_func = self.loop.create_server(
            handler,
            host=get_ip_from_if(svc_config['tr069']['interface']),
            port=svc_config['tr069']['perf_mgmt_port'],
        )

        self._periodic_check_rf_tx()
        self.loop.run_until_complete(create_server_func)

    def _periodic_check_rf_tx(self) -> None:
        self._check_rf_tx()
        self.mme_timeout_handler = self.loop.call_later(
            self.CHECK_RF_TX_PERIOD,
            self._periodic_check_rf_tx,
        )

    def _check_rf_tx(self) -> None:
        """
        Check if eNodeB should be connected to MME but isn't, and maybe reboot.

        If the eNB doesn't report connection to MME within a timeout period,
        get it to reboot in the hope that it will fix things.

        Usually, enodebd polls the eNodeB for whether it is connected to MME.
        This method checks the last polled MME connection status, and if
        eNodeB should be connected to MME but it isn't.
        """
        # Clear stats when eNodeB stops radiating. This is
        # because eNodeB stops sending performance metrics at this point.
        serial_list = self.enb_manager.get_connected_serial_id_list()
        for enb_serial in serial_list:
            handler = self.enb_manager.get_handler_by_serial(enb_serial)
            self._check_rf_tx_for_handler(handler)

    def _check_rf_tx_for_handler(self, handler: EnodebAcsStateMachine) -> None:
        status = get_enb_status(handler)
        if self._prev_rf_tx and not status.rf_tx_on:
            self._clear_stats()
        self._prev_rf_tx = status.rf_tx_on

        # Update status metrics
        update_status_metrics(status)

    def _get_enb_label_from_request(self, request) -> str:
        label = 'default'
        ip = request.headers.get('X-Forwarded-For')

        if ip is None:
            peername = request.transport.get_extra_info('peername')
            if peername is not None:
                ip, _ = peername

        if ip is None:
            return label

        label = self.enb_manager.get_serial_of_ip(ip)
        if label:
            logger.debug('Found serial %s for ip %s', label, ip)
        else:
            logger.error("Couldn't find serial for ip", ip)
        return label

    @asyncio.coroutine
    def _post_handler(self, request) -> web.Response:
        """ HTTP POST handler """
        # Read request body and convert to XML tree
        body = yield from request.read()

        root = ElementTree.fromstring(body)
        self._parse_pm_xml(self._get_enb_label_from_request(request), root)

        # Return success response
        return web.Response()

    def _parse_pm_xml(self, enb_label, xml_root) -> None:
        """
        Parse performance management XML from eNodeB and populate metrics.
        The schema for this XML document, along with an example, is shown in
        tests/stats_manager_tests.py.
        """
        for measurement in xml_root.findall('Measurements'):
            object_type = measurement.findtext('ObjectType')
            names = measurement.find('PmName')
            data = measurement.find('PmData')
            if object_type == 'EutranCellTdd':
                self._parse_tdd_counters(enb_label, names, data)
            elif object_type == 'ManagedElement':
                # Currently no counters to parse
                pass
            elif object_type == 'SctpAssoc':
                # Currently no counters to parse
                pass

    def _parse_tdd_counters(self, enb_label, names, data):
        """
        Parse eNodeB performance management counters from TDD structure.
        Most of the logic is just to extract the correct counter based on the
        name of the statistic. Each counter is either of type 'V', which is a
        single integer value, or 'CV', which contains multiple integer
        sub-elements, named 'SV', which we add together. E.g:
        <V i="9">0</V>
        <CV i="10">
          <SN>RRC.AttConnReestab.RECONF_FAIL</SN>
          <SV>0</SV>
          <SN>RRC.AttConnReestab.HO_FAIL</SN>
          <SV>0</SV>
          <SN>RRC.AttConnReestab.OTHER</SN>
          <SV>0</SV>
        </CV>
        See tests/stats_manager_tests.py for a more complete example.
        """
        index_data_map = self._build_index_to_data_map(data)
        name_index_map = self._build_name_to_index_map(names)

        # For each performance metric, extract value from XML document and set
        # internal metric to that value.
        for pm_name, metric in self.PM_FILE_TO_METRIC_MAP.items():

            elements = pm_name.split(':')
            counter = elements.pop(0)
            if len(elements) == 0:
                subcounter = None
            else:
                subcounter = elements.pop(0)

            index = name_index_map.get(counter)
            if index is None:
                logger.warning('PM counter %s not found in PmNames', counter)
                continue

            data_el = index_data_map.get(index)
            if data_el is None:
                logger.warning('PM counter %s not found in PmData', counter)
                continue

            if data_el.tag == 'V':
                if subcounter is not None:
                    logger.warning('No subcounter in PM counter %s', counter)
                    continue

                # Data is singular value
                try:
                    value = int(data_el.text)
                except ValueError:
                    logger.info(
                        'PM value (%s) of counter %s not integer',
                        data_el.text, counter,
                    )
                    continue
            elif data_el.tag == 'CV':
                # Check whether we want just one subcounter, or sum them all
                subcounter_index = None
                if subcounter is not None:
                    index = 0
                    for sub_name_el in data_el.findall('SN'):
                        if sub_name_el.text == subcounter:
                            subcounter_index = index
                        index = index + 1

                if subcounter is not None and subcounter_index is None:
                    logger.warning('PM subcounter (%s) not found', subcounter)
                    continue

                # Data is multiple sub-elements. Sum them, or select the one
                # of interest
                value = 0
                try:
                    index = 0
                    for sub_data_el in data_el.findall('SV'):
                        if subcounter_index is None or \
                                subcounter_index == index:
                            value = value + int(sub_data_el.text)
                        index = index + 1
                except ValueError:
                    logger.error(
                        'PM value (%s) of counter %s not integer',
                        sub_data_el.text, pm_name,
                    )
                    continue
            else:
                logger.warning(
                    'Unknown PM data type (%s) of counter %s',
                    data_el.tag, pm_name,
                )
                continue

            # Apply new value to metric
            if pm_name == 'PDCP.UpOctUl' or pm_name == 'PDCP.UpOctDl':
                metric.labels(enb_label).set(value)
            else:
                metric.set(value)

    def _build_index_to_data_map(self, data_etree):
        """
        Parse XML ElementTree and build a dict mapping index to data XML
        element. The relevant part of XML schema being parsed is:
        <xs:element name="PmData">
         <xs:complexType>
          <xs:sequence minOccurs="0" maxOccurs="unbounded">
           <xs:element name="Pm">
            <xs:complexType>
             <xs:choice minOccurs="0" maxOccurs="unbounded">
              <xs:element name="V">
               <xs:complexType>
                <xs:simpleContent>
                 <xs:extension base="xs:string">
                  <xs:attribute name="i" type="xs:integer" use="required"/>
                 </xs:extension>
                </xs:simpleContent>
               </xs:complexType>
              </xs:element>
              <xs:element name="CV">
               <xs:complexType>
                <xs:sequence minOccurs="0" maxOccurs="unbounded">
                 <xs:element name="SN" type="xs:string"/>
                 <xs:element name="SV" type="xs:string"/>
                </xs:sequence>
                <xs:attribute name="i" type="xs:integer" use="required"/>
               </xs:complexType>
              </xs:element>
             </xs:choice>
             <xs:attribute name="Dn" type="xs:string" use="required"/>
             <xs:attribute name="UserLabel" type="xs:string" use="required"/>
            </xs:complexType>
           </xs:element>
          </xs:sequence>
         </xs:complexType>
        </xs:element>

        Inputs:
            - XML elementree element corresponding to 'PmData' in above schema
        Outputs:
            - Dict mapping index ('i' in above schema) to data elementree
              elements ('V' and 'CV' in above schema)
        """
        # Construct map of index to pm_data XML element
        index_data_map = {}
        for pm_el in data_etree.findall('Pm'):
            for data_el in pm_el.findall('V'):
                index = data_el.get('i')
                if index is not None:
                    index_data_map[index] = data_el
            for data_el in pm_el.findall('CV'):
                index = data_el.get('i')
                if index is not None:
                    index_data_map[index] = data_el

        return index_data_map

    def _build_name_to_index_map(self, name_etree):
        """
        Parse XML ElementTree and build a dict mapping name to index. The
        relevant part of XML schema being parsed is:
        <xs:element name="PmName">
         <xs:complexType>
          <xs:sequence minOccurs="0" maxOccurs="unbounded">
           <xs:element name="N">
            <xs:complexType>
             <xs:simpleContent>
              <xs:extension base="xs:string">
               <xs:attribute name="i" type="xs:integer" use="required"/>
              </xs:extension>
             </xs:simpleContent>
            </xs:complexType>
           </xs:element>
          </xs:sequence>
         </xs:complexType>
        </xs:element>

        Inputs:
            - XML elementree element corresponding to 'PmName' in above schema
        Outputs:
            - Dict mapping name ('N' in above schema) to index ('i' in above
              schema)
        """
        # Construct map of pm_name to index
        name_index_map = {}
        for name_el in name_etree.findall('N'):
            name_index_map[name_el.text] = name_el.get('i')

        return name_index_map

    def _clear_stats(self) -> None:
        """
        Clear statistics. Called when eNodeB management plane disconnects
        """
        logger.info('Clearing performance counter statistics')
        # Set all metrics to 0 if eNodeB not connected
        for pm_name, metric in self.PM_FILE_TO_METRIC_MAP:
            # eNB data usage metrics will not be cleared
            if pm_name not in ('PDCP.UpOctUl', 'PDCP.UpOctDl'):
                metric.set(0)
