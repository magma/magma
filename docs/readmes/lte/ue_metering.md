---
id: ue_metering
title: UE Usage Metering
hide_title: true
---
# UE Usage Metering

Magma currently supports basic usage metering. This allows for real-time
monitoring of data usage specified with the following labels:
- `IMSI`
- `session_id`
- `traffic direction`

This information is available through our metrics REST endpoint.
The metric name used is `ue_traffic`.

![Swagger REST API Endpoint](assets/ue_metering.png)