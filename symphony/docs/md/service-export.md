---
id: service-export
title: Exporting your service data
---

### Exporting the data

* Go to the 'Services' tab (chain icon).
* Use the filters bar to get a subset of your result (optional).
* Click the "Export" button on the top right corner.


* A CSV file containing the filtered service list will be downloaded.
* Every row represents an service , and it will be of the following form and in the following order:
   * "Service ID" - Internal id of the service in inventory. Unique
   * "Service Name" - Unique
   * "Service Type" - Name of the type of the service (defined in 'Configure' tab)
   * "Service External ID" - Service ID used to identify the service in other systems (CRM for example).
                             Unique. Maybe be empty
   * "Customer Name" - Unique. Maybe be empty
   * "Customer External ID" - Customer ID used to identify the customer in other systems (CRM for example).
                              Unique. Maybe be empty.
   * "Status" - 4 options: PENDING, IN_SERVICE, MAINTENANCE, DISCONNECTED
   * List of properties for this service.

