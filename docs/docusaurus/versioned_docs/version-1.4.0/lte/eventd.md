---
id: version-1.4.0-eventd
title: Eventd
hide_title: true
original_id: eventd
---
# Eventd
### Overview

Eventd is a structured logging service that unifies event-logs from Magma services.
We use [Swagger 2.0](https://swagger.io/specification/) specifications to define events that services can emit, to keep track of a state machine or process.

Events are placed into `streams`, which are like logical buckets for events. They can be organized using a `tag`, and must conform to a structure defined by its `event_type`.

### Infrastructure

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

### Gateway eventd gRPC Interface

To publish gateway events on a gateway, the gRPC interface must be used.

```
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
```
{
  "stream_name": "sessiond",
  "event_type": "session_created",
  "tag": "IMSI001010000000099",
  "value": {"apn": "oai.ipv4", "imei": "", "imsi": "IMSI001010000000099", "ip_addr": "192.168.128.96", "mac_addr": "", "msisdn": "", "pdp_start_time": 1598803879, "session_id": "IMSI001010000000099-736956", "spgw_ip": "10.0.2.1"
}
```
As can be seen with this example, `value` holds the event specific data.
An event must be defined before they can be send to `eventd` service.

### Dependencies

- [Fluentd](https://www.fluentd.org/) is run on the orc8r for event aggregation, and moves them to elasticsearch
- [Fluent Bit](https://fluentbit.io/) (run by the `td-agent-bit` service) is required to forward event-logs. As such, eventd only runs on LTE for now.
- [Elasticsearch](https://www.elastic.co/elastic-stack)
- [Swagger-codegen](https://github.com/swagger-api/swagger-codegen) is used to generate python classes from the swagger definitions. This is currently used to generate convenience classes to emit events
- [Bravado-core](https://github.com/Yelp/bravado-core) is used type-check an event based on its `event_type`. This is a python dependency, and helps to simplify the processes of validating an event object against a schema

### How-to: Create an event

- Create a spec under `<plugin_name>/swagger` e.g. `lte/swagger/mock_events.v1.yml`, or add an event to one of the specs
- Register the `event_type` and the location of the swagger file under `eventd.yml`'s event registry.
- Make an RPC call to eventd's `log_event` from your service, using the appropriate client API.
  - (Python-only) Use `make build` under `lte/gateway` to generate swagger models into the `$PYTHON_BUILD` directory. e.g. Use model `ue_added` with the import `<plugin_name>.swagger.models.ue_added`

