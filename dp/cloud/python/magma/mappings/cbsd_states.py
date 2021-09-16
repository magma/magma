from magma.mappings.types import CbsdStates, GrantStates, Switch
from dp.protos.active_mode_pb2 import Off, On, Registered, Unregistered
from dp.protos.common_pb2 import Authorized, Granted

cbsd_state_mapping = {
    CbsdStates.UNREGISTERED.value: Unregistered,
    CbsdStates.REGISTERED.value: Registered,
}
grant_state_mapping = {
    GrantStates.GRANTED.value: Granted,
    GrantStates.AUTHORIZED.value: Authorized,
}
switch_mapping = {
    Switch.ON.value: On,
    Switch.OFF.value: Off,
}
