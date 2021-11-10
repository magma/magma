from magma.mappings.types import CbsdStates, GrantStates
from dp.protos.active_mode_pb2 import Authorized, Granted, Registered, Unregistered

cbsd_state_mapping = {
    CbsdStates.UNREGISTERED.value: Unregistered,
    CbsdStates.REGISTERED.value: Registered,
}
grant_state_mapping = {
    GrantStates.GRANTED.value: Granted,
    GrantStates.AUTHORIZED.value: Authorized,
}
