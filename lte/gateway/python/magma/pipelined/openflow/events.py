import copy

from ryu.controller import event


class EventSendPacket(event.EventBase):
    def __init__(self, pkt):
        super().__init__()
        self.packet = copy.copy(pkt)
