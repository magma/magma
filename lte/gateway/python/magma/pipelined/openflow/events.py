import copy

from ryu.controller import event


class EventSendPacket(event.EventBase):
    def __init__(self, pkt, imsi=None):
        super().__init__()
        self.packet = copy.copy(pkt)
        self.imsi = imsi
