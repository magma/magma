import logging
import time

from eventlet import queue
from magma.pipelined.app.base import MagmaController
from magma.pipelined.app.inout import EGRESS
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.events import EventSendPacket
from magma.pipelined.openflow.flows import DROP_PRIORITY
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import TestPacket, TEST_PACKET_REG, \
    IMSI_REG
from ryu.app.ofctl.api import get_datapath
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER
from ryu.controller.handler import set_ev_cls
from ryu.lib import hub
from ryu.lib.packet.packet import Packet


class PacketTracingController(MagmaController):
    APP_NAME = "packet_tracer"
    _EVENTS = [EventSendPacket]

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._datapath = None
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME
        )
        self.dropped_table = {}
        self.drop_flows_installed = set()

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self.logger.debug('Tracer connected (dp.id): %d', datapath.id)

        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                                             match=MagmaMatch(),
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

    def install_drop_flows(self):
        """
        Install flows that send test-packets to the controller instead of
        dropping them if no match is found
        """
        drop_tables = set()
        for tables in self._service_manager \
                .get_all_table_assignments() \
                .values():
            drop_tables.update(tables.scratch_tables + [tables.main_table])

        drop_tables.discard(0)
        drop_tables = drop_tables.difference(self.drop_flows_installed)
        for table in drop_tables:
            self.logger.debug('Installing drop flow for %d', table)
            flows.add_trace_packet_output_flow(datapath=self._datapath,
                                               table=table,
                                               match=MagmaMatch(),
                                               instructions=[],
                                               priority=DROP_PRIORITY)
            self.drop_flows_installed.add(table)

    def trace_packet(self, packet, imsi, timeout=2):
        """
        Send a packet and wait until it is processed and the dropped_table dict
        shows which table caused a drop.

        Important: trace_packet initiates a packet-send through
        self.send_event_to_observers() because trace_packet could be called
        from a different thread than the one Ryu app is running on, which will
        cause the datapath to reconnect. To avoid that we initiate an
        EventSendPacket which will be handled on the Ryu event-loop thread.

        :param packet: bytes of the packet
        :param imsi: IMSI string (like 001010000000013 or IMSI001010000000013)
        :param timeout: timeout after which to stop waiting for a packet
        :return: table_id which caused the drop or -1 if it wasn't dropped
        """
        assert isinstance(packet, bytes)
        self.logger.debug('Trace packet: %s', str(Packet(packet)))
        self.send_event_to_observers(EventSendPacket(pkt=packet, imsi=imsi))

        start = time.time()
        while packet not in self.dropped_table:
            if time.time() > start + timeout:
                raise Exception('Timeout while waiting for the packet')
            time.sleep(0.1)
        table_id = self.dropped_table[packet]
        del self.dropped_table[packet]

        # If the table was returned to the user space after the last table =>
        # It wasn't dropped in the process
        if table_id == self._service_manager.get_table_num(EGRESS):
            return -1
        return table_id

    def cleanup_on_disconnect(self, datapath):
        self.logger.debug('Tracer disconnected (dp.id): %d', datapath.id)
        assert self._datapath.id == datapath.id
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    # pylint: disable=try-except-raise
    def _event_loop(self):
        """
        Override the event loop because it was getting stuck with no timeout
        The only difference from the base _event_loop is that
        we set a timeout on receiving an event and continue
        if the queue is empty
        """
        LOG = logging.getLogger('ryu.base.app_manager')
        while self.is_active or not self.events.empty():
            try:
                ev, state = self.events.get(timeout=0.1)
            except queue.Empty:
                continue

            self._events_sem.release()
            if ev == self._event_stop:
                continue
            handlers = self.get_handlers(ev, state)
            for handler in handlers:
                # noinspection PyBroadException
                # pylint: disable=broad-except
                try:
                    handler(ev)
                except hub.TaskExit:
                    # Normal exit.
                    # Propagate upwards, so we leave the event loop.
                    raise
                except Exception:
                    LOG.exception(
                        '%s: Exception occurred during handler processing. '
                        'Backtrace from offending handler '
                        '[%s] servicing event [%s] follows.',
                        self.name, handler.__name__, ev.__class__.__name__)

    @set_ev_cls(EventSendPacket)
    def _send_packet(self, ev):
        """
        First install drop flows if they are not there yet.
        Then send the packet through the switch
        :param ev: EventSendPacket
        """
        self.install_drop_flows()

        pkt = ev.packet
        imsi = ev.imsi
        if isinstance(pkt, (bytes, bytearray)):
            data = bytearray(pkt)
        elif isinstance(pkt, Packet):
            pkt.serialize()
            data = pkt.data
        else:
            raise ValueError('Could not handle packet of type: '
                             '{}'.format(type(pkt)))

        self.logger.debug('Tracer sending packet: %s', str(Packet(data)))
        datapath = get_datapath(self, dpid=self._datapath.id)
        ofp = datapath.ofproto
        ofp_parser = datapath.ofproto_parser
        actions = [
            # Turn on test-packet as we're just tracing it
            ofp_parser.NXActionRegLoad2(dst=TEST_PACKET_REG,
                                        value=TestPacket.ON.value),
            # Add IMSI metadata
            ofp_parser.NXActionRegLoad2(dst=IMSI_REG,
                                        value=encode_imsi(imsi)),
            # Submit to table=0 because otherwise the packet will be dropped!
            ofp_parser.NXActionResubmitTable(table_id=0),
        ]
        datapath.send_packet_out(in_port=ofp.OFPP_LOCAL,
                                 actions=actions,
                                 data=data)

    @set_ev_cls(ofp_event.EventOFPPacketIn, MAIN_DISPATCHER)
    def _handle_packet_in_user_space(self, ev):
        """
        Receive the table_id which caused the packet to be dropped
        """
        if ev.msg.match[TEST_PACKET_REG] != TestPacket.ON.value:
            return
        pkt = Packet(data=ev.msg.data)
        self.logger.debug('Tracer received packet %s', str(pkt))
        self.dropped_table[pkt.data] = ev.msg.table_id
