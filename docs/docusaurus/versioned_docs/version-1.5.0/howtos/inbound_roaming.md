---
id: version-1.5.0-inbound_roaming
title: Inbound Roaming
hide_title: true
original_id: inbound_roaming
---

*Last Updated: 6/21/2021*

# Inbound Roaming
Inbound Roaming allows a Magma operator to provide service for subscribers
belonging to other operators (roaming subscribers)

Inbound Roaming requires Magma operator to reach agreements with
other operators to have direct connectivity to their HSS through `S6a`
(diameter) and PGW through `S8` (gtp) interfaces. VPN is suggested to reach
those roaming services but is out of scope of Magma.

At the bottom of this document you have
[configuration example](#configuration-example)

## Architecture
Currently, we support two Architectures:

- **Local Non-Federated + Roaming**: where local subscribers are stored in
  SubscriberDB and roaming subscribers use a remote Federated Gateway to reach
  HSS/PGW.
- **Local Federated + Roaming**: where local Subscribers use local Federated
  Gateway to reach HSS/PCRF/OCS and roaming Subscribers use a remote Federated
  Gateway to reach HSS/PGW.

Any roaming architecture will be composed of an LTE Federated Network and Feg
Network that will serve the local subscribers. Then we will have as many extra
FEG Network as roaming agreements.

The example in the picture blow shows a possible architecture for a
Local Non-Federated + Roaming case (one roaming agreement). As you can see in
the picture below, VPN is suggested to reach to the user plane at the remote
PGW.

![Magma events table](assets/feg/inbound_roaming_architecture_non_federated.png?raw=true "Non-Federated Inbound Roaming")

Other configurations like SubscriberDB with PCRF/OCS to handle local
subscribers may also work but are not tested yet.

Configurations with SubscriberDB, local HSS and remote HSS all at the same
time is not supported yet.

## Prerequisites
Before starting to configure roaming setup, first you need to bring up a
setup to handle your own/local subscribers. So before configuring Inbound
Roaming you need:
- Install [Or8cr](https://docs.magmacore.org/docs/orc8r/architecture_overview),
- Install [Federatetion Gateway](https://docs.magmacore.org/docs/feg/deploy_intro) and,
- Install [Access Gateway](https://docs.magmacoreorg/docs/lte/setup_deb).
- Create a Federate Deployment (see [below](#Create a Federated Deployment)).
- Make sure your setup is able to serve calls with your local subscribers

Once you are done you should either:
- Local Non-Federated case:  `a Federated LTE Network with an LTE Gateway (AGW)` + `a Federation Network WITHOUT Federation Gateway (FeG)`
- Local Federated case: `a Federated LTE Network with an LTE Gateway (AGW)`+
  `a Federation Network WITH Federation Gateway (FeG)`

We will refer to them as `local` Networks and Gateways to differentiate them
from the `roaming` Networks and Gateways we will create in the next step.

### Create a Federated Deployment
As mentioned, Inbound Roaming requires of a FeG gateway to reach the roaming
network. That is why Federated Deployment is required. Please, configure it
using this guide for [Federated Deployment](https://docs.magmacore.org/docs/feg/federated_FWA_setup_guide).

All architectures requiere a Local FeG Network to exist. However depending on
your architecture, you may not need to create a local FeG Gateway inside that
FeG Network. Please check the table below which indicates what **gateways**
are required depending on the architecture

| Network Type ->       | `feg` roaming | `feg` local | `feg_lte` (AGW) |
| :-------------------: | :-----------: | :----------:| :--------------:|
| Local - Non Federated | Yes           | No          | Yes             |
| Local - Federated     | Yes           | Yes         | Yes             |


For both, Federated and Non-Federated mode, on `lte_feg` network under `epc`
you should set `hss_relay_enabled` to `true`. Selection of HSS or
subscriberDb will be done by Inbound Roaming mapping configuration.

In case you are not using PCRF and OCS, then your
not need to create a FeG Gateway from the previous gide on your Federated
Deployment. Also remember in your Federated LTE network, under `epc`, you
will have to set `gx_gy_relay_enabled` to `false`, so the request are sent to
internal policy entity, not to the PCRF or OCS.

## Inbound Roaming configuration
The following instructions use Orc8r Swagger API to configure Inbound Roaming.
You can do the same using Swagger API or NMS JSON editor.

In these instructions we will mainly use GET and PUT methods to read and write
from Swagger. We will use GET to see the content of Network/Gateway,
we will copy and paste that into the result of GET into the PUT method to
modify parameters.

Below are the steps to add Inbound Roaming to your current setup:
- Create roaming Federated Networks and Gateways.
- Configure local Access Gateway Network routing based on PLMN.
- Configure local Federation Network routing based on PLMN.
- Configure roaming Federation Network served networks
- Configure FeG Gateways
- Check connectivity

### 1. Create Roaming Federated Networks and Federated Gateways
Inbound Roaming needs as many FeG Networks as roaming agreements. Don't
forget to create a Federated Gateway per FeG Network.

Roaming FeG Networks do not need to serve any Federated LTE Network (only the
FeG Network created on Pre Requisites needs to serve a Federated LTE Network)

Those roaming Federated Gateways will need `S6a` and `S8` interfaces
configured (make sure you have the other operators HSS and PGW
parameters). To configure those interfaces go to **Swagger API**:
- Go to`Federation Gateways` GET method `Get a specific
  federation gateway` and search the configuration for one of those roaming
  FeG Networks.
- Copy/paste the response into the PUT method `Update an entire federation
  gateway record`
- Edit the `6a` and `S8` fields (check the example from the GET method to
  see any missing parameter)
- Hit Execute (check no errors are show on the Swagger Responses)
- Run the GET method again to see the changes.

Note that you need a routable IP to the roaming HSS and PGW to configure those
interfaces.

### 2. Configure Local Access Gateway Network routing.
When a request gets to Access Gateway, this will have to be routed to
the proper call flow to either use `SubscriberDB` or `S6a/S5` or`S6a/S8`.
based on subscriber PLMN.

To enable that routing you will have to configure it in your local Access
Gateway Network we created in Pre Requisites. On **Swagger API**:
- Go to `Federated LTE Networks` and search using GET method `Describe a
  federated LTE network` your local (non roaming) Federated LTE Network
- Copy/paste the response into PUT method `Update an entire Federated LTE
  network`
- Find a key `federation`. If you completed Pre Requisites properly, you
  should have there `feg_network_id` pointing to your FeG Network Modify/add
- Add the routing dictionary following the example, adding an entry per each
PLMN
```text
  "federation": {
    "federated_modes_mapping": {
      "enabled": true,
      "mapping": [
        {
          "apn": "",
          "imsi_range": "",
          "mode": "local_subscriber",
          "plmn": "123456"
        },
        {
          "apn": "",
          "imsi_range": "",
          "mode": "s8_subscriber",
          "plmn": "9999"
        }
      ]
    },
    "feg_network_id": "example_feg_network"
  },
```
- Field `mode` will indicate the flow the subscriber will take. Use
  `local_subscriber` for PLMN served by your own SubscriberDB/HSS. Use
  `s8_subscriber` to use roam HSS and roam PGW. Leave `apn` and `imsi_range`
  blank since it is not supported yet.

- Note `hss_relay_enabled` must be enabled. The decision to send it to HSS or
  not will be taken by `federated_modes_mapping`. If you disable, s8_subscribers
  will not be sent to the FeG to get the HSS

- Flag `gx_gy_relay_enabled` can be enabled or disabled depending if your
  network works with local policy db or with OCS and PCRF (gx/gy). If your
  local subscribers authenticate with HSS but use GX/GY, then you will have to
  leave it as `True`.

### 3. Configure Local Federation Network routing
When a request gets to the Orc8r, this will have to be routed to the proper
FeG Network which serves that PLMN.

To enable that routing you will have to configure it in your local FeG Network
we created in Pre Requisites. On **Swagger API**:
- Go to `Federation Networks` and search using GET method `Describe a
  federation network` your local (non roaming) FeG Network
- Copy/paste the response into PUT method `Update an entire federation network`
- Modify/add the `nh_routes` (see the example on Swagger API if it is
  missing from your configuration). On the map match the PLMN, and the name
  of the roam FeG Network.
```
    "nh_routes": {
      "00102": "inbound_feg",
      "9999": "feg_roaming_network_1"
    },
```
- Hit Execute (check no errors are show on the Swagger Responses)
- Run the GET method again to see the changes.

### 4. Configure Roaming Federation Network served networks
Roaming Federation Networks will need a last configuration in order to match
them with their serving Local Federation Network. To do that, add to the
Roaming Federation Networks configuration the following key
```
    "served_nh_ids": [
      "example_feg_network"
    ],
```

That means that the Inbound FeG Network will be served by the local network
(in this case called `example_feg_network`)

### 5. Configure Roaming FeG Gateway
Configure Roaming FeG gateway serving roaming subscribers, but just configure
`s6a` and `s8`. Configure local GTP port to match with your PGW GTP-U port.
`apn_operator_sufix` is optional and will just add a suffix to the APN sent by
the UE.
```
   "s8": {
      "apn_operator_suffix": ".operator.com",
      "local_address": "foo.bar.com:5555",
      "pgw_address": "foo.bar.com:5555"
    },
```
Note you don't need to define the local IP, you can just use :port
`"local_address": ":5555"`

### 6. Check connectivity
- From Access Gateway make sure your Access Gateway is able to reach PGW-U
  IP.
```
    ping -I sgi_interface pgw_u_ip
```
- Fom Federated Gateway make sure you can reach PGW-C IP and test your
Federated Gateway can reach PGW using this command
```
    cd /var/opt/magma/docker
    sudo docker-compose exec s8_proxy /var/opt/magma/bin/s8_cli cs -server
    192.168.32.118:2123 123456789012345
    # where 192.168.32.118:2123 is the ip and port of the PGW-C
    # where 123456789012345 is a valid imsi (if you use a not valid imsi
    you can still check the connectivity, but you will get a GTP error back
    from PGW
    # Add -use_builtincli flag if you don't have a FeG setup properly yet
```

## Test and troubleshooting
It is recommendable that before running the tests, enable some extra
logging capabilities in both Access Gateway, and Federated Gateway to
trace the call.

For better details in Access Gateway logs:
- Enable `log_level: DEBUG` in `mme.yml` and `subscriberdb.yml`
- Enable `print_grpc_payload: True` on `subscriberdb.yml`
- Restart magma, so the changes are taken
- See the logs using `sudo journalctl -fu magma@mme` or sudo `journalctl -fu
  magma@subscriberdb`

For better details Federated Gateway logs:
- Add GRPC printing in the following services  `s6a_proxy`,
  `s8_proxy` adding `MAGMA_PRINT_GRPC_PAYLOAD: 1`. For example for s6a_proxy
```
  s6a_proxy:
    <<: *goservice
    container_name: s6a_proxy
    command: envdir /var/opt/magma/envdir /var/opt/magma/bin/s6a_proxy -logtostderr=true -v=0
    environment:
      MAGMA_PRINT_GRPC_PAYLOAD: 1
```
- Restart docker process, so the vars are taken `sudo docker-compose down` and
  `sudo docker-compose up -d`
- Display the logs using for example`sudo docker-compose logs -f s8_proxy`

### Test with s6a_cli and s8_cli
FeG has a couple of clients to run an HSS Authentication Request (s6a) and
Create Session Request (s8) without the need of having a UE. You can run them
either on FeG or AGW.

- Run From FeG
```
# Use FeG s6a_proxy
sudo docker-compose exec s8_proxy /var/opt/magma/bin/s6a_cli air -remote_s6a 001002000000810
# Use s6a_porxy that runs on the cli
sudo docker-compose exec s8_proxy /var/opt/magma/bin/s6a_cli air -use_builtincli false  -remote_s6a 001002000000810

# use FeG s8_proxy
sudo docker-compose exec s8_proxy /var/opt/magma/bin/s8_cli cs -server 192.168.32.118:2123 -delete 3
# use s8_porxy that runs on the cli
sudo docker-compose exec s8_proxy /var/opt/magma/bin/s8_cli cs -server 192.168.32.118:2123 -use_builtincli false -delete 3
```

- Run from AGW
```
# Extract the binaries from docker container from FeG, and move them to AGW
sudo docker cp s6a_proxy:/var/opt/magma/bin/s6a_cli .
sudo docker cp s8_proxy:/var/opt/magma/bin/s8_cli .
```

```
# Execute from AGW
./s6a_cli air -remote_s6a 001002000000810
./s8_cli cs -server 192.168.32.118:2123 -delete 3 -apn inet -remote_s8 001002000000810
```

## Configuration example

Attached you can find the configuration that handle local subscribers with
both subscriber db and HSS and roaming subscribers:

- PLMN 88888: uses subscriber DB to authenticate and Gx/Gy for accounting.
  That is why we have `gx_gy_relay_enabled` as True. Those subscribers are
  never sent to the FeG.
- PLMN 00102: MME sends those subscribers to be authenticated through the FeG.
  When the request reaches Orc8r (in Feg Relay service) and using `nh_routes`
  configured on the local FeG network, those subscribers are forwarded to
  `inbound_feg` network
- Rest of PLMN: MME forwards any other PLMN to be authenticated through the
  FeG. In the orc8r they are forwarded to the local FeG network `terravm_feg_network`

[[inbound_roaming_sample.zip]](assets/feg/inbound_roaming_sample.zip)
