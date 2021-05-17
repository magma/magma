---
id: pipelined_tests
title: Pipelined testing framework
hide_title: true
---
# Pipelined testing framework
## Overview
The testing framework aims to isolate pipelined for better testing. This is
achieved by running tests on a different bridge, running only some of the
pipelined apps. Additionally, by inserting OVS *forwarding* flows we isolate
testing only to specific tables.

The framework can also be used with integration testing, using a processing
thread with hub queues, an integ_test flag is provided in the test config.
This means pipelined tests can work with gRPC.

## Functionality breakdown
This is high-level explanation of what happens when running a test. One of the
main principles in designing this framework was making it as component based as
possible so that its easy to add/replace some parts of the process.

### Launch pipelined application directly (not as services)
The first step is to launch pipelined controllers that we want to test.
By launching ryu application directly and avoiding using services we can get
the references to instantiated controllers, thus testing their functionality
directly.

### Isolate the table that is being tested
As we want to run unit tests its necessary to isolate specific tables. This is
done by inserting special *forwarding* flows that will both forward all packets
to the ****table specified as well as setting the required register values that
would have been set by the tables skipped(such as metadata and reg1).

### Insert flow rules if needed (f.e. subscriber policy rules)
Having references to pipelined controllers makes it simple to insert flows into
OVS. The test framework provides an api to add PolicyRules for subscribers.

### Using Scapy insert packets into OVS
For inserting packets the testing framework uses the Scapy library. A wrapper
for easier packet building, packet insertion is provided. After sending packets
its necessary to wait for packets to be received/processed by OVS, example
test files have wait functions to achieve this.

## Testing controller
The testing controller is a ryu app used for instantiating testing
flows(table isolation) and for querying flow stats. The controller is only used
for testing purposes and only runs when invoked in tests.

## API variations
Initially the testing framework was developed using REST and gRPC. The RyuRPC*,
RyuRest* classes are still present but are deprecated. No active tests use them
as they require a running pipelined service to function.
The primary API is RyuDirect*, all active tests use it. Later gRPC calls
will work with this framework after some threading fixes.

## Writing a new test
Example test files are a good place to see the framework in action, also the
`pipelined_test_util.py` file provides convenient functions for easier and
faster test writing.

### Setup that can be used for multiple tests, this should go in `setUpClass`
**Setup a new bridge or use the production bridge**
```
BRIDGE = 'testing_br'
IFACE = 'testing_br'
BridgeTools.create_bridge(BRIDGE, IFACE)
```

**Start ryu apps on a separate thread**
```
# Set the futures for pipelined controller references
enforcement_controller_reference = Future()
testing_controller_reference = Future()

# Build a test_setup for launching ryu apps
test_setup = TestSetup(
    apps=[PipelinedController.Enforcement,
          PipelinedController.Testing],
    references={
        PipelinedController.Enforcement:
            enforcement_controller_reference,
        PipelinedController.Testing:
            testing_controller_reference
    },
    config={
        'bridge_name': cls.BRIDGE,
        'bridge_ip_address': '192.168.128.1',
        'nat_iface': 'eth2',
        'enodeb_iface': 'eth1'
    },
    mconfig=None,
    loop=None
)

# Start the apps from the test_setup config
cls.thread = start_ryu_app_thread(test_setup)

# Wait for apps to start, retrieve references
cls.enforcement_controller = enforcement_controller_reference.result()
cls.testing_controller = testing_controller_reference.result()
```

### Unit test example
**Setup basic information/constants**
```
# Setup subscriber info imsi, ip and a PolicyRule
imsi = 'IMSI010000000088888'
sub_ip = '192.168.128.74'
flow_list1 = [FlowDescription(
    match=FlowMatch(
        ipv4_dst='45.10.0.0/24', direction=FlowMatch.UPLINK),
    action=FlowDescription.PERMIT)
]
policy = PolicyRule(id='simple_match', priority=2, flow_list=flow_list1)

pkts_matched = 256
pkts_sent = 4096
```

**Setup the testing framework classes**
```
# Create a subscriber context, used for adding PolicyRules
sub_context = RyuDirectSubscriberContext(
    imsi, sub_ip, self.enforcement_controller
).add_dynamic_rule(policy)

# Create a table isolator from subscriber context, will set metadata/reg1,
# forward the packets to table 5
isolator = RyuDirectTableIsolator(
    RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                             .build_requests(),
    self.testing_controller
)

# Create a packet sender, an ip packet for testing our PolicyRule
pkt_sender = ScapyPacketInjector(self.IFACE)
packet = IPPacketBuilder()\
    .set_ip_layer('45.10.0.0/20', sub_ip)\
    .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00")\
    .build()

# Generate a flow query for checking the stats of the added rule
flow_query = FlowQuery(
    self.TID, self.testing_controller,
    match=flow_match_to_match(flow_list1[0].match)
)
```

**Test & Verify everything works**
```
# Verify aggregate table stats
# Each FlowTest provides a query and number of packets that query should match
# wait_after_send function ensures all packets are processed before verifying
flow_verifier = FlowVerifier([
    FlowTest(FlowQuery(self.TID, self.testing_controller), pkts_sent),
    FlowTest(flow_query, pkts_matched)
], lambda: wait_after_send(self.testing_controller))

# Initialize all contexts and then send packets
# isolator -  inserts the flow forward rules that forward traffic to table 5,
#             set the required registers (metadata, reg1)
#
# sub_context - adds subscriber flows, activates the PolicyRule provided
#
# flow_verifier - gathers FlowTest stats on when initialized and when exiting
with isolator, sub_context, flow_verifier:
    pkt_sender.send(packet)

# Asserts that conditions were met for each provided FlowTest
flow_verifier.verify()
```
