---
id: events_monitoring
title: Events Monitoring
hide_title: true
---
# Events Monitoring
### Overview

On top of our logging and metrics, Magma has a metrics monitoring system.
Events capture rich, searchable data pertaining to the gateways,
and entities tracked by the gateways.
Events are a work in progress, but has been designed to track the following:
- S1AP message tracing
- Subscriber session tracing
- Subscriber data usage, bearer creation
- Gateway connection status
- eNodeB tracing
- New gateway configurations pushed

If Fluentd, Fluent Bit, and Elasticsearch have been correctly configured on
your setup, events should automatically be available.

### Viewing Events

We recommend two primary ways to view or query for events: either through our
NMS UI, or through our REST API.

#### NMS

![Magma events table](assets/lte/events_table.png?raw=true "Magma Events Table")

The easiest way to view Magma's events is through the NMS UI.
Events can be seen on the main network dashboard.
Gateway and subscriber specific events tables are also provided.

#### REST API

We provide a REST API for querying events.
The Magma NMS itself uses this API to provide its interface.
We use OpenAPI, so more details can be found in our various Swagger files.
`swagger.yml`
`swagger.v1.yml`

Currently we provide three GET endpoints for our REST API:

```
/events/{network_id}:
/events/{network_id}/{stream_name}:
/events/{network_id}/about/count:
```

### Event Streams

Events are divided into the following streams, that can be filtered on either
the NMS or the REST API

(This may not be an exhaustive list)

- `magmad`: Events related to general access gateway function
- `health` Events related to federation gateway health
- `sessions`: Events related to subscriber sessions
- `mme` Events related to MME function, S1 issues, and subscriber attach/detach

### Available Events

Events are defined in swagger yaml files.
This list may not be exhaustive.

#### Access Gateway Events

`magmad_events.v1.yml`

*deleted_stored_mconfig*: The stored mconfig was deleted

*updated_stored_mconfig*: The stored mconfig was updated

*processed_updates*: Stream updates were successfully processed
```
Properties
- updates (new mconfig received)
```

*restarted_services*: Services were restarted
```
Properties
- services
```

*established_sync_rpc_stream*: SyncRPC connection was established

*disconnected_sync_rpc_stream*: SyncRPC stream was disconnected
    
#### Federation Gateway Events

`health_events.v1.yml`

*gateway_promotion_succeeded*: Gateway successfully promoted to active

*gateway_promotion_failed*: Gateway promotion to active failed
```
Properties
- failure reason
```

*gateway_demotion_succeded*: Gateway successfully demoted to standby

*gateway_demotion_failed*: Gateway demotion to standby failed
```
Properties
- failure reason
```
        
#### AAA Events

`aaa_server_events.v1.yml`

*authentication_succeeded*: Used to track successful subscriber authentications
```
Properties
- IMSI
- session ID
- MAC address
- APN
```

*authentication_failed*: Used to track failed subscriber authentications
```
Properties
- IMSI
- MAC address
- APN
- failure reason
```

*session_termination_succeeded*: Used to track AAA server produced session terminations
```
Properties
- IMSI
- session ID
- MAC address
- APN
- reason for termination
```

*session_termination_failed*: Used to track AAA server produced session termination failures
```
Properties
- IMSI
- session ID
- MAC address
- APN
- reason for termination
- failure reason
```

#### Subscriber Events

`session_manager_events.v1.yml`

`mme_events.v1.yml`

*session_created*: Used to track when a session was created
```
Properties
- imsi:
- ip_addr:
- msisdn:
- imei:
- spgw_ip:
- session_id:
- apn:
- mac_addr:
- pdp_start_time:
```

*session_create_failure*: Used to track when a session creation failed
```
Properties
- imsi:
- apn:
- mac_addr:
- failure_reason:
```

*session_updated*: Used to track when a session update is reported
```
Properties
- imsi:
- apn:
- mac_addr:
- ip_addr:
```

*session_update_failure*: Used to track when a session update has failed
```
Properties
- imsi:
- apn:
- mac_addr:
- ip_addr:
- failure_reason:
```

*session_terminated*: Used to track total session metrics
```
Properties
- imsi:
- apn:
- mac_addr:
- ip_addr:
- msisdn:
- imei:
- spgw_ip:
- session_id:
- total_tx:
- total_rx:
- charging_tx:
- charging_rx:
- monitoring_tx:
- monitoring_rx:
- pdp_start_time:
- pdp_end_time:
```

*attach_success*: Used to track when UE attaches successfully
```
Properties
- imsi:
```

*detach_success*: Used to track when UE detaches successfully
```
Properties
- imsi:
- action:
```

*s1_setup_success*: Used to track establishment of S1 connection
```
Properties
- enb_name:
- enb_id:
```

### How-to: Create an event

View under our 'Eventd' doc.
