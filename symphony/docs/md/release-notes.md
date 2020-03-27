---
id: release-notes
title: Release Notes
---

### Work In Progress

The team is currently working on:
* Saved Searches: Ability to save past searches as bookmarks
* Mobile app: View work orders and site surveys in the mobile app
* Permission model: Customizeable infrastructure to enable roles and policies
* API documentation: Improving documentation for existing APIs
* Check lists for work orders: Enable check list items for work orders
* Enable SSO integration via Keycloak


### Release Notes
* 3/22/2020
    * Inventory
        * **Saved Searches**: Filters can now be saved and named for future use (reports, services, work orders)
        * **Search w/ Breadcrumbs**: Display location/equipment breadcrumbs in Inventory main search results
    * Work Orders
        * **Adding "Close Time"**: Work Orders table now showing a column with the time the WorkOrder status was set to "DONE"
    * Infra
        * **Privacy support**: Add support for allow/deny logic around mutations
        * **Schema for users**: Manage details on users in the DB and connect them to work orders\projects
    * APIs
        * **Pyinventory**: add get_equipment_properties
        * **More Tests**: Improved Pyinventory test coverage
        * **User Managment**: Added GraphQL API for creating\editing users
       
       
* 3/8/2020
    * Bug fixes
        * **Validations on graphql**: several data validation were happening on the UI and not checked when calling directly the GraphQL API. Moved those validations to the GraphQL endpoint.
        * **UI fixes**: Improved User Experience
    * Inventory
        * **Adding warning before deletion**: When an equipment is being deleted, and this equipment has sub-equipments, warn the user that this deletion will delete more than 1 object
    * Work Orders
        * **Export**: Added "Export to CSV" option to Word Orders search view
    * Infra
        * **Subscriptions**: Send notifications via GraphQL subscriptions about changes to WO status
        * **Safe changes to our GraphQL schema**: Block changes to GraphQL that are breaking previous versions from being pushed to production
        * **Adding flow typing**: Improve the Flow coverage in UI files
        * **Enhancing UI Design system**: Icons, Generic View Containers, Radio Groups, Tabs, Different Table views
    * APIs
        * **Pyinventory**: Added: edit equipment & port & link properties