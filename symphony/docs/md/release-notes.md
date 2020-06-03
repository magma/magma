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
* 5/22/2020
    * Mobile App
        * **[My Tasks]** **[Checklist]**	Photos upload in checklist:	The user is now able to upload more than 1 photo per checklist item during site survey. The user can choose whether to upload photos from the phone gallery or take photos via the camera.
        * **[My Tasks]**	Offline mode:	The user is now able to complete a work order tasks even if the device does not have internet connection.
    
    * Platform
        * **[Property type]**	Support of "User" property:	The user is able to add "user" as "property type" while creating/editing templates and assign a specific value from a drop down list.
		
    * Inventory Management
    	* **[Services]** Added "discovery method" field in Services:	The user can now choose "inventory" as "discovery method" during the creation of a service type. When a new "service type" is defined and the option "inventory" is selected, all Services  matching that specific "service type" definition will be automatically detected every hour and presented in the list of Services.
        * **[Services]**	Added "export" for Services: The user can now export the list of Services in csv format
        * **[Services]**	Added "filter by location" for Services	The user can now filter Services also by "location type"
    
    * Workforce Management
        * **[Work Orders]** **[UI]** 	Work Order's location opened in a new tab: When opening a Work order page, the user can click on the work order's location and a new tab in the browser will be opened, allowing the user to keep both tabs (Work order and Location's tabs) opened at the same time.  	
        * **[Templates]**	Checklist definition: The user is now able to define checklist categories on a template and easily create work orders that inherit them. 
        * **[Work Orders]** **[Checklist]**	Checklist response view enabled: The user is able to view the responses within the checklist categories even when the work order's status is in "done" status.
     
* 4/30/2020
    * Mobile App
        * **New Release**: New major app release. This release includes work order: view work orders, edit work orders, edit checklist. Download from here: [https://play.google.com/apps/testing/cloud.thesymphony](https://play.google.com/apps/testing/cloud.thesymphony)
    * Inventory and Work Orders
        * **New Property Type**: Added property type of type "Work Order". Now you can create new properties that are linking to specific work orders.
        * **New Property Type**: Added property type of type "User". Now you can create new properties that are linking to specific users.
    * Infra
        * **Permission Model**: New infrastructure for permission model is ready. User roles were replaced with the new model. This is unblocking custom permission definition.
        * **Library Upgrade**: Upgraded gqlgen library to 0.11.3 (Included open census support)
        * **Library Upgrade**: Upgraded open source charts to latest version
    * APIs
        * **Pyinventory**: added get_equipments_by_type, get_equipments_by_location, get_equipment_by_external_id, get_equipment_type_property_type, get_equipment_type_property_type_by_external_id , edit_equipment_type_property_type, get_property_type_id 
        * **Pyinventory**: added external_id to: add_equipment, add_equipment_to_position, copy_equipment_in_position, copy_equipment, get_or_create_equipment, get_or_create_equipment_in_position
        * **Pyinventory**: Performance improvments. get_location and get_location_by_external_id run X3 times faster

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