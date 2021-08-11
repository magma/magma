# Copyright (C) 2014 Nippon Telegraph and Telephone Corporation.
# Copyright (C) 2014 YAMAMOTO Takashi <yamamoto at valinux co jp>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ofctl service

import numbers

from ryu.base import app_manager

from ryu.controller import ofp_event
from ryu.controller.handler import CONFIG_DISPATCHER, MAIN_DISPATCHER,\
    DEAD_DISPATCHER
from ryu.controller.handler import set_ev_cls

from . import event
from . import exception


class _SwitchInfo(object):
    def __init__(self, datapath):
        self.datapath = datapath
        self.xids = {}
        self.barriers = {}
        self.results = {}


class OfctlService(app_manager.RyuApp):
    def __init__(self, *args, **kwargs):
        super(OfctlService, self).__init__(*args, **kwargs)
        self.name = 'ofctl_service'
        self._switches = {}
        self._observing_events = {}

    def _observe_msg(self, msg_cls):
        assert msg_cls is not None
        ev_cls = ofp_event.ofp_msg_to_ev_cls(msg_cls)
        self._observing_events.setdefault(ev_cls, 0)
        if self._observing_events[ev_cls] == 0:
            self.logger.debug('ofctl: start observing %s', ev_cls)
            self.register_handler(ev_cls, self._handle_reply)
            self.observe_event(ev_cls)
        self._observing_events[ev_cls] += 1

    def _unobserve_msg(self, msg_cls):
        assert msg_cls is not None
        ev_cls = ofp_event.ofp_msg_to_ev_cls(msg_cls)
        assert self._observing_events[ev_cls] > 0
        self._observing_events[ev_cls] -= 1
        if self._observing_events[ev_cls] == 0:
            self.unregister_handler(ev_cls, self._handle_reply)
            self.unobserve_event(ev_cls)
            self.logger.debug('ofctl: stop observing %s', ev_cls)

    def _cancel(self, info, barrier_xid, exception):
        xid = info.barriers.pop(barrier_xid)
        req = info.xids.pop(xid)
        msg = req.msg
        datapath = msg.datapath
        parser = datapath.ofproto_parser
        is_barrier = isinstance(msg, parser.OFPBarrierRequest)

        info.results.pop(xid)

        if not is_barrier and req.reply_cls is not None:
            self._unobserve_msg(req.reply_cls)

        self.logger.error('failed to send message <%s>', req.msg)
        self.reply_to_request(req, event.Reply(exception=exception))

    @staticmethod
    def _is_error(msg):
        return (ofp_event.ofp_msg_to_ev_cls(type(msg)) ==
                ofp_event.EventOFPErrorMsg)

    @set_ev_cls(ofp_event.EventOFPSwitchFeatures, CONFIG_DISPATCHER)
    def _switch_features_handler(self, ev):
        datapath = ev.msg.datapath
        id = datapath.id
        assert isinstance(id, numbers.Integral)
        old_info = self._switches.get(id, None)
        new_info = _SwitchInfo(datapath=datapath)
        self.logger.debug('add dpid %s datapath %s new_info %s old_info %s',
                          id, datapath, new_info, old_info)
        self._switches[id] = new_info
        if old_info:
            old_info.datapath.close()
            for xid in list(old_info.barriers):
                self._cancel(
                    old_info, xid, exception.InvalidDatapath(result=id))

    @set_ev_cls(ofp_event.EventOFPStateChange, DEAD_DISPATCHER)
    def _handle_dead(self, ev):
        datapath = ev.datapath
        id = datapath.id
        self.logger.debug('del dpid %s datapath %s', id, datapath)
        if id is None:
            return
        try:
            info = self._switches[id]
        except KeyError:
            return
        if info.datapath is datapath:
            self.logger.debug('forget info %s', info)
            self._switches.pop(id)
            for xid in list(info.barriers):
                self._cancel(info, xid, exception.InvalidDatapath(result=id))

    @set_ev_cls(event.GetDatapathRequest, MAIN_DISPATCHER)
    def _handle_get_datapath(self, req):
        result = None
        if req.dpid is None:
            result = [v.datapath for v in self._switches.values()]
        else:
            if req.dpid in self._switches:
                result = self._switches[req.dpid].datapath
        self.reply_to_request(req, event.Reply(result=result))

    @set_ev_cls(event.SendMsgRequest, MAIN_DISPATCHER)
    def _handle_send_msg(self, req):
        msg = req.msg
        datapath = msg.datapath
        parser = datapath.ofproto_parser
        is_barrier = isinstance(msg, parser.OFPBarrierRequest)

        try:
            si = self._switches[datapath.id]
        except KeyError:
            self.logger.error('unknown dpid %s' % (datapath.id,))
            rep = event.Reply(exception=exception.
                              InvalidDatapath(result=datapath.id))
            self.reply_to_request(req, rep)
            return

        def _store_xid(xid, barrier_xid):
            assert xid not in si.results
            assert xid not in si.xids
            assert barrier_xid not in si.barriers
            si.results[xid] = []
            si.xids[xid] = req
            si.barriers[barrier_xid] = xid

        if is_barrier:
            barrier = msg
            datapath.set_xid(barrier)
            _store_xid(barrier.xid, barrier.xid)
        else:
            if req.reply_cls is not None:
                self._observe_msg(req.reply_cls)
            datapath.set_xid(msg)
            barrier = datapath.ofproto_parser.OFPBarrierRequest(datapath)
            datapath.set_xid(barrier)
            _store_xid(msg.xid, barrier.xid)
            if not datapath.send_msg(msg):
                return self._cancel(
                    si, barrier.xid,
                    exception.InvalidDatapath(result=datapath.id))

        if not datapath.send_msg(barrier):
            return self._cancel(
                si, barrier.xid,
                exception.InvalidDatapath(result=datapath.id))

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    def _handle_barrier(self, ev):
        msg = ev.msg
        datapath = msg.datapath
        parser = datapath.ofproto_parser
        try:
            si = self._switches[datapath.id]
        except KeyError:
            self.logger.error('unknown dpid %s', datapath.id)
            return
        try:
            xid = si.barriers.pop(msg.xid)
        except KeyError:
            self.logger.error('unknown barrier xid %s', msg.xid)
            return
        result = si.results.pop(xid)
        req = si.xids.pop(xid)
        is_barrier = isinstance(req.msg, parser.OFPBarrierRequest)
        if req.reply_cls is not None and not is_barrier:
            self._unobserve_msg(req.reply_cls)
        if is_barrier and req.reply_cls == parser.OFPBarrierReply:
            rep = event.Reply(result=ev.msg)
        elif any(self._is_error(r) for r in result):
            rep = event.Reply(exception=exception.OFError(result=result))
        elif req.reply_multi:
            rep = event.Reply(result=result)
        elif len(result) == 0:
            rep = event.Reply()
        elif len(result) == 1:
            rep = event.Reply(result=result[0])
        else:
            rep = event.Reply(exception=exception.
                              UnexpectedMultiReply(result=result))
        self.reply_to_request(req, rep)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    def _handle_reply(self, ev):
        msg = ev.msg
        datapath = msg.datapath
        try:
            si = self._switches[datapath.id]
        except KeyError:
            self.logger.error('unknown dpid %s', datapath.id)
            return
        try:
            req = si.xids[msg.xid]
        except KeyError:
            self.logger.error('unknown error xid %s', msg.xid)
            return
        if ((not isinstance(ev, ofp_event.EventOFPErrorMsg)) and
                (req.reply_cls is None or not isinstance(ev.msg, req.reply_cls))):
            self.logger.error('unexpected reply %s for xid %s', ev, msg.xid)
            return
        try:
            si.results[msg.xid].append(ev.msg)
        except KeyError:
            self.logger.error('unknown error xid %s', msg.xid)
