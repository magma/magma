---
id: eventd
title: Event Reporting
hide_title: true
---

# Event Reporting

This document describes how gateway events are reported via the *eventd* service.

## Overview

Eventd is a structured logging service that unifies event-logs from Magma services.
We use [Swagger 2.0](https://swagger.io/specification/) specifications to define events that services can emit, to keep track of a state machine or process.

Events are placed into `streams`, which are like logical buckets for events. They can be organized using a `tag`, and must conform to a structure defined by its `event_type`.

## Infrastructure

![Magma events architecture diagram](assets/lte/events_architecture.png?raw=true "Magma Events Architecture")
The events pipeline is set up with functional pieces across the access gateway, orc8r, and NMS.

First, on the gateway, events are aggregated by the `eventd` service.
To fire events, the gRPC interface on `eventd` must be called.
The `eventd` service is responsible for sending these events to Fluent Bit,
labelled as `td-agent-bit` on the diagram above.

Fluent Bit is a lightweight log forwarder, which we use to aggregate our events
to the orc8r.
(Fluent Bit is also used to aggregate logs, but that is not the topic of this document)
From Fluent Bit, events are forwarded to Fluentd, running on the orc8r.

Fluentd is used on the orc8r to collect events from all gateways.
Fluentd is configured to store all collected events in Elasticsearch,
under an `eventd` index prefix.

For network operators to get visibility into the Magma events, an API is provided.

## Gateway eventd gRPC Interface

To publish gateway events on a gateway, the gRPC interface must be used.

```grpc
service EventService {
  rpc LogEvent (Event) returns (Void) {}
}

message Event {
  string stream_name = 1;
  string event_type = 2;
  string tag = 3;
  string value = 4;
}
```

Example event data:

```json
{
  "stream_name": "sessiond",
  "event_type": "session_created",
  "tag": "IMSI001010000000099",
  "value": {"apn": "oai.ipv4", "imei": "", "imsi": "IMSI001010000000099", "ip_addr": "192.168.128.96", "mac_addr": "", "msisdn": "", "pdp_start_time": 1598803879, "session_id": "IMSI001010000000099-736956", "spgw_ip": "10.0.2.1"
}
```

As can be seen with this example, `value` holds the event specific data.
An event must be defined before they can be send to `eventd` service.

## Dependencies

- [Fluentd](https://www.fluentd.org/) is run on the orc8r for event aggregation, and moves them to elasticsearch
- [Fluent Bit](https://fluentbit.io/) (run by the `td-agent-bit` service) is required to forward event-logs. As such, eventd only runs on LTE for now.
- [Elasticsearch](https://www.elastic.co/elastic-stack)
- [Swagger-codegen](https://github.com/swagger-api/swagger-codegen) is used to generate python classes from the swagger definitions. This is currently used to generate convenience classes to emit events
- [Bravado-core](https://github.com/Yelp/bravado-core) is used type-check an event based on its `event_type`. This is a python dependency, and helps to simplify the processes of validating an event object against a schema

## How-to: Create an event

- Create a spec under `<plugin_name>/swagger` e.g. `lte/swagger/mock_events.v1.yml`, or add an event to one of the specs
- Register the `event_type` and the location of the swagger file under `eventd.yml`'s event registry.
- Make an RPC call to eventd's `log_event` from your service, using the appropriate client API.

### Example

Create new file `$MAGMA_ROOT/lte/swagger/my_events.v1.yml` with content:

```yml
---
swagger: '2.0'

info:
  title: MY event definitions
  description: These events occur in MY task
  version: 1.0.0

definitions:
  my_success:
    type: object
    description: My success
    properties:
      msg:
        type: string
```

Add the following to `$MAGMA_ROOT/lte/gateway/configs/eventd.yml`:

```yml
  my_success:
    module: lte
    filename: my_events.v1.yml
```

Use the following code in a Python script to emit the event:

```python
import snowflake
from magma.eventd.eventd_client import log_event
from orc8r.protos.eventd_pb2 import Event

log_event(
    Event(
        stream_name="magmad",
        event_type="my_success",
        tag=snowflake.snowflake(),
        value="{}",
    ),
)
```

Add the following to the BUILD.bazel file in the directory, where your script is:

```python
MAGMA_ROOT = "RELATIVE_PATH_TO_MAGMA_ROOT"
ORC8R_ROOT = "{}orc8r/gateway/python".format(MAGMA_ROOT)

py_binary(
    name = "your_script_name",
    srcs = ["your_script_name.py"],
    imports = [ORC8R_ROOT],
    legacy_create_init = False,
    deps = [
        "//orc8r/gateway/python/magma/eventd:eventd_client",
        "//orc8r/protos:mconfig_python_proto",
        requirement("snowflake"),
    ],
)
```

Finally call the script with Bazel:

`bazel run //PATH/TO/YOUR/SCRIPT:your_script_name`

The event should show up in the eventd logs.
