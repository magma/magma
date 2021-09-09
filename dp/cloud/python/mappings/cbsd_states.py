from dp.protos.active_mode_pb2 import Unregistered, Registered, Granted, Authorized, On, Off
from dp.cloud.python.mappings.types import CbsdStates, GrantStates, Switch

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
