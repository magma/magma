---
id: config_apn
title: APN Configuration
hide_title: true
---
# Access Point Name (APN) Configuration
UEs can successfully attach and get connected to the Magma AGWs if they have a valid APN configuration in their subscription profiles on the network side. Typically, UEs send APN information explicitly in their connection requests. Magma AGW pulls APN information from the subscription data to verify that UEs have indeed subscription for the requested APN. If APN information is missing from the connection request, AGW picks the first APN in the subscriber profile as the default APN and establishes a connection session according to that default APN.

## Defining APN Configurations
The first step in APN Configuration is to make sure that the desired APN profiles are already defined for the network:
- Navigate to your NMS instance and on the sidebar click on "Configure" button.
- In the newly opened page, on the top bar select "APN CONFIGURATION". If there are already APNs defined, it would show up on this page (e.g., see the screenshot below).
- You can edit or delete any of the existing APN configurations. Note that the updates and deletions would be reflected automatically in subscriber profiles and new attach as well as PDN connection requests would be impacted by these changes.
- You can also add a new APN configuration by clicking on the "Add APN" button and filling up the requested fields. After saving these changes, the page should refresh with the new list of APNs and their configurations.

![Creating an APN Configuration](assets/nms/add_apnconfig.png)

## Adding APN Configurations to Subscriber Profiles
The next step is to add one or more APN configurations to the subscriber profiles so that UEs can start consuming network services based on their APNs:

- For an existing subscriber, to update its subscription profile, simple click on the edit field and perform a multi-select under the "Access Point Names".
- For a new subscriber, similarly fill up the fields including the "Access Point Names" field (screenshot below shows the view after clicking on the "Add Subscriber" button). Once you save the updated or new subscriber information, the APNs added to the subscriber profile would be refreshed and shown on the page.

![Adding subscriber with APN](assets/nms/add_apn2subscriber.png)

### Notes
- The first APN listed under "Active APNs" for each subscriber becomes the default APN that would be used if UE omits the APN information in its connection requests.
- The subscriber data is streamed down to AGWs periodically and the new configs should be reflected on AGW with some lag.
- To check if AGW is already updated, on the AGW, one can run the following command to retrieve the subscriber data that includes the subscribed APN profiles:

`subscriber_cli.py get IMSI<15 digit IMSI>`

An example output for a hypothetical user with IMSI 001010000000001 and APNs "internet", "ims" is shown below:

```bash
sid {
  id: "001010000000001"
}
lte {
  state: ACTIVE
  auth_key: "..." # not shown in this example
}
network_id {
  id: "my_network"
}
state {
}
sub_profile: "default"
non_3gpp {
  apn_config {
    service_selection: "ims"
    qos_profile {
      class_id: 5
      priority_level: 9
    }
    ambr {
      max_bandwidth_ul: 100000
      max_bandwidth_dl: 100000
    }
  }
  apn_config {
    service_selection: "internet"
    qos_profile {
      class_id: 9
      priority_level: 15
    }
    ambr {
      max_bandwidth_ul: 100000000
      max_bandwidth_dl: 200000000
    }
  }
}
```
