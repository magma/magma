---
id: version-1.5.0-inbound_roaming
title: Inbound Roaming
hide_title: true
original_id: inbound_roaming
---

*Last Updated: 4/08/2021*

# Inbound Roaming

Inbound Roaming allows a Magma operator to provide service for subscribers
belonging to other operators (roaming subscribers)

Inbound Roaming requires Magma operator to reach agreements with
other operators to have direct connectivity to their HSS through `S6a`
(diameter) and PGW through `S8` (gtp) interfaces. VPN is suggested to reach
those roaming services but is out of scope of Magma.

## Architecture

Currently, we support two Architectures:
- **Non-Federated + Roaming**: where local subscribers are stored in
  SubscriberDB and roaming subscribers use a remote Federated Gateway to reach HSS/PGW
- **Federated + Roaming**: where local Subscribers use local Federated Gateway to reach
  HSS/PCRF/OCS and roam Subscribers use a remote Federated Gateway to reach HSS/PGW.

The example below shows a possible architecture for a Non-Federated +
roaming case. As you can see in the picture below, VPN is suggested to
reach to the user plane at the remote PGW.

![Magma events table](assets/feg/inbound_roaming_architecture_non_federated.png?raw=true "Non-Federated Inbound Roaming")

Other configurations like SubscriberDB with PCRF/OCS to handle local
subscribers may also work but are not tested yet.

Configurations with SubscriberDB, local HSS and remote HSS all at the same
time is not supported yet.

## Pre Requisites

Before starting to configure roaming setups, first you need to bring up a
setup to handle your own/local subscribers.

Even one of the architectures mentions `Non-Federated`, Inbound Roaming
requires a `Federated` setup to work. So the to configure Inbound Roaming
you need:
- Functioning [Or8cr](https://docs.magmacore.org/docs/orc8r/architecture_overview),
- Install [Federatetion Gateway](https://docs.magmacore.org/docs/feg/deploy_intro) and,
- Install [Access Gateway](https://docs.magmacoreorg/docs/lte/setup_deb).
- Configure it as a [Federated Deployment](https://docs.magmacore.org/docs/feg/federated_FWA_setup_guide)
- Make sure your setup is able to serve calls with your local subscribers
  (in case of using subscriber DB you will need to set `"hss_relay_enabled":
  false` on Federated LTE Network temporally to test this. Set it back to
  `true` once tested!!)

At the end of this you will have:
- A Federated LTE Network with an LTE Gateway (AGW)
- A Federation Network with a Federation Gateway (FEG)

We will refer to them as `local` Networks and Gateways to differentiate them
from the `roaming`Networks and Gateways we will create in the next step.

*Note that in case of **Non-Federated + roaming** you can skip the creation
of the Federated Gateway, but you still need the Federated Gateway Network that will serve your
Federated LTE Network.*

## Inbound Roaming Configuration
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
- Check connectivity

### 1. Create Roaming Federated Networks and Federated Gateways
Inbound Roaming needs as many FEG Networks as roaming agreements. Don't
forget to create a Federated Gateway per FEG Network.

Roaming FEG Networks do not need to serve any Federated LTE Network (only the
FEG Network created on Pre Requisites needs to serve a Federated LTE Network)

Those roaming Federated Gateways will need `S6a` and `S8` interfaces
configured (make sure you have the other operators HSS and PGW
parameters). To configure those interfaces go to **Swagger API**:
- Go to`Federation Gateways` GET method `Get a specific
  federation gateway` and search the configuration for one of those roaming
  FEG Networks.
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
  should have there `feg_network_id` pointing to your FEG Network Modify/add
- Add the routing dictionary following the example, adding an entry per each
PLMN
```text
  "federation": {
    "federated_modes_mapping": {
      "enabled": false,
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

### 3. Configure Local Federation Network routing
When a request gets to the Orc8r, this will have to be routed to the proper
FEG Network which serves that PLMN.

To enable that routing you will have to configure it in your local FEG Network
we created in Pre Requisites. On **Swagger API**:
- Go to `Federation Networks` and search using GET method `Describe a
  federation network` your local (non roaming) FEG Network
- Copy/paste the response into PUT method `Update an entire federation network`
- Modify/add the `nh_routes` (see the example on Swagger API if it is
  missing from your configuration). On the map match the PLMN, and the name
  of the roam FEG Network.
```
  "nh_routes": {
      "00102": "inbound_feg",
      "9999": "feg_roaming_network_1"
    },
```
- Hit Execute (check no errors are show on the Swagger Responses)
- Run the GET method again to see the changes.

### 4. Check Connectivity
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
    # Add -use_builtincli flag if you don't have a FEG setup properly yet
```

## Test and Troubleshooting
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
