import subprocess
import time

from prometheus_client import Gauge

services = ["magma@control_proxy", "magma@directoryd", "magma@enodebd", "magma@health", "magma@mme", "magma@pipelined", "magma@redis", "magma@smsd", "magma@subscriberdb", "magma@ctraced", "magma@dnsd", "magma@eventd", "magma@magmad", "magma@mobilityd", "magma@policydb", "magma@sessiond", "magma@state"]

control_proxy_state = Gauge('controlproxystate', 'maintains state for control proxy')
directoryd_state = Gauge('directorydstate', 'maintains state for directoryd')
enodebd_state = Gauge('enodebdstate', 'maintains state for enodebd')
health_state = Gauge('heathservicestate', 'maintains state for health service')
mme_state = Gauge('mmeservicestate', 'maintains state for mme service')
pipelined_state = Gauge('pipelinedservicestate', 'maintains state for pipelined service')
redis_state = Gauge('redis_state', 'maintains state for redis service')
smsd_state = Gauge('smsd_servicestate', 'maintains state for smsd service')
subscriberdb_state = Gauge('subscriberdbservicestate', 'maintains state for subscriberdb service')
ctraced_state = Gauge('ctraced_service_state', 'maintains state for ctraced service')
dnsd_state = Gauge('dnsd_service_state', 'maintains state for dnsd service')
eventd_state = Gauge('eventd_service_state', 'maintains state for eventd service')
magmad_state = Gauge('magmad_service_state', 'maintains state for magmad service')
mobilityd_state = Gauge('mobilityd_service_state', 'maintains state for mobilityd service')
policydb_state = Gauge('policydb_service_state', 'maintains state for policydb service')
sessiond_state = Gauge('sessiond_service_state', 'maintains state for sessiond service')
state_state = Gauge('state_service_state', 'maintains state for state service')

state_map = {"magma@control_proxy": control_proxy_state, "magma@directoryd": directoryd_state, "magma@enodebd": enodebd_state, "magma@health": health_state, "magma@mme": mme_state, "magma@pipelined": pipelined_state, "magma@redis": redis_state, "magma@smsd": smsd_state, "magma@subscriberdb": subscriberdb_state, "magma@ctraced": ctraced_state, "magma@dnsd": dnsd_state, "magma@eventd": eventd_state, "magma@magmad": magmad_state, "magma@mobilityd": mobilityd_state, "magma@policydb": policydb_state, "magma@sessiond": sessiond_state, "magma@state": state_state}


def check_service_running():
    for service in services:
        p = subprocess.Popen(["systemctl", "is-active", service], stdout=subprocess.PIPE)
        (output, err) = p.communicate()
        output = output.decode('utf-8')
        g = state_map[service]
        if output == "active\n":
          g.set(1)
        else:
          g.set(0)


if __name__ == '__main__':
    while True:
        print("Checking magma service state")
        check_service_running()
        time.sleep(5)
