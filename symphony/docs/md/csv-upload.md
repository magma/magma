---
id: csv-upload
title: Upload Inventory Data via CSV    
---

![](https://s3.amazonaws.com/purpleheadband.images/wiki/uploadmenu2.png)

## Intro

Here you can find a guide for uploading your CSV data into Inventory.

## 1. Before importing you data
* Location types, Equipment types, Port types etc should be defined manually before upload.
* Verify that the locations types are defined by their order:
Under the location types tab, drag and drop the locations types so that it's arranged from big ones to small ones.
* Export the relevant data to CSV (which will be the template for the upload), [here's how you do it](export.md)
* It's important to note that there are four different upload flows for exported data:
	* Exported Equipment (for add and edit)
	* Exported Ports (edit only)
	* Exported Links (add & edit)
	* Exported Locations (edit only)
	* Exported Services (add & edit) 
* CSV file should be of utf8 format.

## 2. Importing the Data

* Once the CSV file is exported, it can be modified and be used in order to upload new records or edit existing ones.
### 2.1 Equipment
#### 2.1.1 Editing Existing Equipment
* As long as the value of the "Equipment ID" column is not empty, a row will be treated as "to be edited".
* Possible fields to be edited:
   * "Equipment Name" - for renaming an equipment instance.
      * If an equipment is being renamed and it has references to it down the CSV (as parent equipment for example) - edit those as well to the new name.
   * "External ID"
   * Every property
      * Property is editable as long as the corresponding equipment-type supports this property (can be verified on Inventory, under "configure"-> "equipment-types")

#### 2.1.2 Adding New Equipment

* As long as the value of the "Equipment ID" column is empty, a row will be treated as "to be added".
* Fields:
   * "Equipment Name"
   * "Equipment Type" - must be an existing name
   * "External ID" - (optional) an ID from a third-party system 
   * List of locations, from big to small - (will be added on the fly if not exists, but with no location properties)
   * "Parent Equipment (3)" - if exists
   *  "Position (3)" - if exists
   * "Parent Equipment (2)" - if exists
   * "Position (2)" - if exists
   * "Parent Equipment" - if exists
   * "Parent Position" - if exists
   * List of properties for this equipment - if exist
* Equipment positions won't be added on the fly - they should exist and be free in advance of the new import run.

#### 2.1.3 Uploading the modified CSV


![](https://s3.amazonaws.com/purpleheadband.images/wiki/exported_data_for_upload.png)

* Now that we have our CSV ready - 
   *  Inventory
   * Click the '+' sign and a dialog will be opened.
   * Click the "Bulk Upload" tab
   * "Upload Exported Equipment"
   * Choose the edited file.


### 2.2 Ports
#### 2.2.1 Editing Existing Ports
* As long as the value of the "Port ID" column is not empty, a row will be treated as "to be edited".
* Possible fields to be edited:
   * Every port property
      * Property is editable as long as the corresponding port-type supports this property (can be verified on Inventory, under "configure"-> "port-types"-> "port-properties")

#### 2.2.2 Uploading the modified CSV
* Now that we have our CSV ready - 
   *  Inventory
   * Click the '+' sign and a dialog will be opened.
   * Click the "Bulk Upload" tab
   * "Upload Exported Ports"
   * Choose the edited file.
   
### 2.3 Links
#### 2.3.1 Editing Existing Link
* As long as the value of the "Link ID" column is not empty, a row will be treated as "to be edited".
* Possible fields to be edited:
   * Every link property
      * Property is editable as long as the corresponding port-type supports this property (can be verified on Inventory, under "configure"-> "port-types"-> "link-properties")

#### 2.3.2 Adding New Links

* As long as the value of the "Link ID" column is empty, a row will be treated as "to be added".
* Fields for each one of the ports (wrote "A" but the same behavior for "B"):
   * "Port A Name" - must be a valid port name from the equipment type.
   * "Equipment A Name" - will be added on the fly if not exists.
   * "Equipment A Type" - must be an existing name.
   * List of locations, from big to small - (will be added on the fly if not exists, but with no location properties)
   * "Parent Equipment (3) A" - if exists
   * "Position (3) A" - if exists
   * "Parent Equipment (2) A" - if exists
   * "Position (2) A" - if  exists
   * "Parent Equipment A" - if exists
   * "Parent Position A" - if exists
   * {*Exact same columns for  port B*}
   * List of link properties taken from both ports - if exist
* Equipment positions won't be added on the fly - they should exist and be free in advance of the new import run.

#### 2.3.3 Uploading the modified CSV
* Now that we have our CSV ready - 
   *  Inventory
   * Click the '+' sign and a dialog will be opened.
   * Click the "Bulk Upload" tab
   * "Upload Exported Links"
   * Choose the edited file.
   
   
### 2.4 Locations
#### 2.4.1 Editing Existing Locations
* As long as the value of the "Location ID" column is not empty, a row will be treated as "to be edited".
* Possible fields to be edited:
   * "External ID"
   * "Latitude"
   * "Longitude"

#### 2.1.2 Adding New Location

* As long as the value of the "Location ID" column is empty, a row will be treated as "to be added".
* Fields:

   * List of locations, from big to small - (will be added on the fly if not exists, but with no location properties)
   * "External ID" - (optional) an ID from a third-party system 
   * "Latitude" - to specify the coordinates of the location
   * "Longitude" - to specify the coordinates of the location


## 3. [DEPRECATED] Locations upload - not via export 

Location files should be of the following form:

![](https://s3.amazonaws.com/purpleheadband.images/wiki/full_location.png)

* First row is the location types (in order, specified on the "configure->location types page) followed by are the properties, 'external ID'and latitude/longitude.
* the location  types should be at the beginning, but after  that the order does not matter
* Properties and 'external ID' will be added to the smallest location of that row (on the prev example - (2nd row) 1st floor on '2392 S Wayside D' Building is of size - 200 sq ft)
* 'external ID' is an optional column name to add to each location and it later can be searched upon.
* Add 'latitude' and 'longitude' columns to specify the coordinates of the location

* If there is not one hierarchy - multiple files can be uploaded.
      * for example, if there can be both "Buildings" and "Rooms" under cities - have one file with the hierarchy of ...City => Building => Floor etc ..., and one with the ....City=>Room

