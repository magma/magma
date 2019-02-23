"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
import logging
from xml.etree import ElementTree

from aiohttp import web

from magma.common.misc_utils import get_ip_from_if
from magma.configuration.service_configs import load_service_config
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

    def run(self):
        """ Create and start HTTP server """
        svc_config = load_service_config("enodebd")

        app = web.Application()
        app.router.add_route('POST', "/{something}", self.post_handler)

        loop = asyncio.get_event_loop()
        handler = app.make_handler()
        create_server_func = loop.create_server(
            handler,
            host=get_ip_from_if(svc_config['tr069']['interface']),
            port=svc_config['tr069']['perf_mgmt_port'])

        loop.run_until_complete(create_server_func)

    @asyncio.coroutine
    def post_handler(self, request):
        """ HTTP POST handler """
        # Read request body and convert to XML tree
        body = yield from request.read()

        root = ElementTree.fromstring(body)
        self.parse_pm_xml(root)

        # Return success response
        return web.Response()

    def parse_pm_xml(self, xml_root):
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
                self.parse_tdd_counters(names, data)
            elif object_type == 'ManagedElement':
                # Currently no counters to parse
                pass
            elif object_type == 'SctpAssoc':
                # Currently no counters to parse
                pass

    def parse_tdd_counters(self, names, data):
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
        index_data_map = self.build_index_to_data_map(data)
        name_index_map = self.build_name_to_index_map(names)

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
                logging.info('PM counter %s not found in PmNames', counter)
                continue

            data_el = index_data_map.get(index)
            if data_el is None:
                logging.info('PM counter %s not found in PmData', counter)
                continue

            if data_el.tag == 'V':
                if subcounter is not None:
                    logging.info('No subcounter in PM counter %s', counter)
                    continue

                # Data is singular value
                try:
                    value = int(data_el.text)
                except ValueError:
                    logging.info('PM value (%s) of counter %s not integer',
                                 data_el.text, counter)
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
                    logging.info('PM subcounter (%s) not found', subcounter)
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
                    logging.info('PM value (%s) of counter %s not integer',
                                 sub_data_el.text, pm_name)
                    continue
            else:
                logging.info('Unknown PM data type (%s) of counter %s',
                             data_el.tag, pm_name)
                continue

            # Apply new value to metric
            metric.set(value)

    def build_index_to_data_map(self, data_etree):
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

    def build_name_to_index_map(self, name_etree):
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

    def clear_stats(self):
        """
        Clear statistics. Called when eNodeB management plane disconnects
        """
        logging.info('Clearing statistics')
        # Set all metrics to 0 if eNodeB not connected
        for metric in self.PM_FILE_TO_METRIC_MAP.values():
            metric.set(0)
