---
id: version-1.1.0-eventd
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

Eventd runs on the gateway under magmad. The files for eventd are under `orc8r/gateway/python`, and it's respective client APIs make RPC calls to the service.

#### Dependencies

- [`FluentBit`](https://fluentbit.io/) (run by the `td-agent-bit` service) is required to forward event-logs. As such, eventd only runs on LTE for now.
- [Swagger-codegen](https://github.com/swagger-api/swagger-codegen) is used to generate python classes from the swagger definitions. This is currently used to generate convenience classes to emit events
- [Bravado-core](https://github.com/Yelp/bravado-core) is used type-check an event based on its `event_type`. This is a python dependency, and helps to simplify the processes of validating an event object against a schema

### How-to: Create an event

- Create a spec under `<plugin_name>/swagger` e.g. `lte/swagger/mock_events.v1.yml`, or add an event to one of the specs
- Register the `event_type` and the location of the swagger file under `eventd.yml`'s event registry.
- Make an RPC call to eventd's `log_event` from your service, using the appropriate client API.
  - (Python-only) Use `make build` under `lte/gateway` to generate swagger models into the `$PYTHON_BUILD` directory. e.g. Use model `ue_added` with the import `<plugin_name>.swagger.models.ue_added`

