---
id: version-1.2.0-nms_subscriber
title: NMS Subscriber
hide_title: true
original_id: nms_subscriber
---

# Subscriber Dashboard

## Subscriber Overview
![subscriber_overview1](assets/nms/userguide/subscriber_overview1.png)
The Subscriber Overview page contains a table listing all subscribers in the network. Columns include Name, IMSI, Service configuration(ACTIVE/INACTIVE), Current data usage, Average data usage and Last reported time(displayed if subscriber monitoring is enabled). The table contains an action menu to view, edit or delete an individual subscriber.

## Subscriber Detail
Each Subscriber Detail page comprises of Overview, Events and Config sections.

### Overview
![subscriber_detail1](assets/nms/userguide/subscriber_detail1.png)
* Information section displays Name and IMSI of the subscriber
* Status KPIs display GatewayID, UE latency (if available), connection status (TBD) and eNodeB Serial (TBD)
* Chart displays the upload and download data usage for the subscriber in stacked histograms. The chart shows the upload and download bytes used by the subscriber within the specified interval. The interval is chosen based on the time window selected by the user. The default time window is 3 hours. For a 3-hour window, data usage over a 15 minute interval is displayed as indicated below.
* ![stacked_barchart](assets/nms/userguide/stacked_barchart.png)

* Data usage KPIs to show the data used by the subscriber in the last 1 hour, 1 day, 1 week and 1 month time range.
* Event table displays the subscriber specific events in the last 3 hours.

### Events
![subscriber_detail2](assets/nms/userguide/subscriber_detail2.png)
* Event chart displays the frequency of subscriber specific events in a histogram.
* Event table displays the time, event type and event description. Detailed event description can be viewed by clicking the dropdown on the left. It comes with a search capability to filter events based on type. For example, the user can filter events based on session_creation event
![event_table_session_filtering](assets/nms/userguide/event_table_session_filtering.png)
* Event page comes with a datetime selector which can be used to filter event chart and event table within a selected time window.

### Config contains
![subscriber_detail3](assets/nms/userguide/subscriber_detail3.png)
Subscriber section to view/edit LTE service state, Data plan, Auth Key and Auth OPC and Active APN associated with the subscriber.

# Subscriber Configuration
Subscriber configuration in the overview page provides the ability to do bulk provisioning of subscribers by either
entering the details manually or uploading them through a subscriber csv file.

## Adding new subscriber
Clicking on the 'Add subscriber' button in the subscriber overview page opens the bulk subscriber configuration page.
![add_subscriber_button](assets/nms/userguide/add_subscriber_button.png)

The bulk provisioning page opens up to an empty table. Rows can be added by clicking on the 'add row' button on the right.
Each row contains columns including Subscriber Name, IMSI, Auth Key, Auth OPC, Service state, Data Plan, Active APN.
auth key and auth_opc should be specified as hex. The service state can be either active or inactive. The data plan
and active APN can be chosen based on what is available.

For example,
John, IMSI001010000000079,<auth_key>,<auth_opc>,ACTIVE,default,oai_ipv4

### Provisioning
![subscriber_config1](assets/nms/userguide/subscriber_config1.png)

![subscriber_config2](assets/nms/userguide/subscriber_config2.png)

### Uploading a subscriber csv file
![subscriber_config3](assets/nms/userguide/subscriber_config3.png)

![subscriber_config4](assets/nms/userguide/subscriber_config4.png)

## Editing subscriber
Subscriber information can be edited from the subscriber detail page.

![subscriber_config5](assets/nms/userguide/subscriber_config5.png)

## Deleting subscriber
Subscriber can be deleted from the subscriber overview page as follows.

![subscriber_deletion1](assets/nms/userguide/subscriber_deletion1.png)

![subscriber_deletion1](assets/nms/userguide/subscriber_deletion2.png)
