---
id: config_agw
title: AGW Configuration
sidebar_label: AGW Configuration
hide_title: true
---
# Access Gateway Configuration
### Prerequisites

Before beginning to configure your Magma Access Gateway, you will need to make 
sure that it is running all services without crashing. You will also need a 
working orchestrator setup. The web UI for configuratino is hosted via 
orchestrator.

### Configuration

To set up your python virtualenv, run `magtivate` on your AGW VM.

To get your hardware ID and challenge key, run `show_gateway_info.py`

Configure your access gateway for your desired setup.

Access the Swagger page, which should be running if your Orc8r 
instance is up and running properly.
You can find it at <https://localhost:9443/apidocs> 
or <https://192.169.99.99:9443>

For bare minimum configuration, you'll want to configure the 
following:

* /networks/{network_id}/gateways/{gateway_id}/configs
* /networks/{network_id}/gateways/{gateway_id}/configs/cellular
* /networks/{network_id}/configs/cellular
* /networks/{network_id}/configs/dns

For the most part, configurations can all be default.
You'll want to pay attention to ``/networks/{network_id}/configs/cellular`` as 
the configuration depends on the eNodeB you are using.
