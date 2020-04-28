---
id: nms_grafana
title: NMS Grafana
hide_title: true
---

# NMS Grafana Dashboards

Grafana is now replacing the metrics dashboards in the NMS. Grafana provides a much more powerful, configurable, and user-friendly dashboarding solution. Usage should be hassle-free with the new deployment.

### What’s New?

In the metrics page of the nms, there is now a tab called ‘Grafana’
![Grafana homepage](assets/nms/grafana_homepage.png)
You’ll see three dashboards available to you from the start. These replicate the three dashboards in the NMS that are built-in. Go to one of the dashboards, and you’ll now see a Grafana version of the NMS dashboard.
![Grafana variables](assets/nms/grafana_variables.png)
These dashboards contain dropdown selectors to choose which network(s) and gateway(s) you want to look at. In the NMS dashboard you were only able to look at one network at a time, however in Grafana you can look at any collection of networks or gateways you have access to at once. Simply select or deselect the networks/gateways that you want to see and the graphs will be updated. In the top right corner, there is an option to choose the time range that the graphs display. The default is 6 hours.

### Custom Dashboards

With Grafana, you can create your own custom dashboards and populate them with any graphs and queries you want. The simple way is to just click on the “+” icon on the left sidebar, then create a new dashboard. There is ample documentation about grafana dashboards online if you need help creating your dashboard. These dashboards will be accessible by all users of your NMS organization.
![Grafana new dashboard](assets/nms/grafana_new_dashboard.png)

Grafana documentation on creating dashboards: [Grafana Dashboards](https://grafana.com/docs/grafana/latest/features/dashboard/dashboards/)

Prometheus documentation on writing queries: [Prometheus Querying](https://prometheus.io/docs/prometheus/latest/querying/basics/)

If you want to replicate the networkID or gatewayID variables that you find in the preconfigured dashboards, we provide a “template” dashboard to make that easy. Simply open the Template dashboard, and click on the gear icon near the top right. From there, click “Save As” and enter the name you want. Your new dashboard will now have the gatewayID and networkID variables.

An example of how to use these variables in your queries:
![Grafana query](assets/nms/grafana_query.png)
Some technical details: You need to use `=~` when matching label names with these variables in order to see more than one network or gateway at a time. This is because the `=~` operator tells Prometheus to match the value as a regex.
Once you save your new dashboard, it will be visible to all users of your org.

### How to Access

A user must be a Super-User to access the Grafana tab. This is because in Grafana you can query all metrics from the networks in an organization, but in some cases a user in an organization can access only a subset of those networks. To make a user a super-user, have an admin go to the administration page in the NMS and modify their permissions.

Additionally, this feature must be turned on by the NMS administrator through the Master organization’s feature flag page. This page is at [https://master.<nms-hostname>/master/features](https://master.localtest.me/master/features). Look for the feature named “Include tab for Grafana in the Metrics page”. Click the pencil on the right side and enable this for whichever organizations you want to have access to Grafana.
