---
id: csv-upload
title: Upload Inventory Data via CSV    
---

![](https://s3.amazonaws.com/purpleheadband.images/wiki/uploadmenu.png)

## Intro

Here you can find a guide for uploading your CSV data into Inventory.

## 1. Prerequisites

* Upload the location types and equipment types (i.e the schema) including property types via the UI
* CSV file should be of utf8 format.

## 2. Locations upload

Location files should be of the following form:

![](https://s3.amazonaws.com/purpleheadband.images/wiki/full_location.png)

* First row is the location types (in order, specified on the "configure->location types page) followed by are the properties, 'external ID'and latitude/longitude.
* the location  types should be at the beginning, but after  that the order does not matter
* Properties and 'external ID' will be added to the smallest location of that row (on the prev example - (2nd row) 1st floor on '2392 S Wayside D' Building is of size - 200 sq ft)
* 'external ID' is an optional column name to add to each location and it later can be searched upon.
* Add 'latitude' and 'longitude' columns to specify the coordinates of the location

* If there is not one hierarchy - multiple files can be uploaded.
      * for example, if there can be both "Buildings" and "Rooms" under cities - have one file with the hierarchy of ...City => Building => Floor etc ..., and one with the ....City=>Room

## 3. Position Definition

* Have the following columns on the first row:
   * "Equipment_Type" - equipment type name to attach to - should be created in advance
   * "position_name" - name of the 'position' to be added
   * "Position_Visible_Label" - visible label  of position

## 4. Port Definition

![](https://s3.amazonaws.com/purpleheadband.images/wiki/portDef.png)

* Have **all** the following columns:
   * "Equipment_type" - equipment type name to attach to - mandatory
   * "Port_Type"  - value can be empty
   * "Port_Bandwidth" - value can be empty
   * "Port_Visible_Label" - value can be empty
   * "Port_ID" - port name - mandatory

## 5. Equipment upload

 Equipment files should be of the following form:

> Every row is one equipment

![](https://s3.amazonaws.com/purpleheadband.images/wiki/equipfull.png)

* First row will include (in this order):
   * Location types (all the hierarchy) - similar to how it's present on locations upload, means where will this equipment sit
   * “equipment type”  - the name of the equipment type - mandatory
   * “equipment name”  - equipment instance name - mandatory
   * properties for that type - on each row - only the relevant properties shpuld be filled


## 6. Port Connection

* Have the following columns:
   * A_Equipment - Equipment A name
   * A_Port - Port A name
   * B_Equipment - Equipment name
   * B_Port - Port B name
* The script will create a link between the two ports


# Importing Equipment Exported Data


### 1. Before importing you data
* Verify that the locations types are defined by their order:
Under the location types tab, drag and drop the locations types so that it's arranged from big ones to small ones.
* Export the equipment data to CSV, [here's how you do it](equipment-export.md)
### 2. Importing the Data

* Ones the CSV file is exported, it can be modified and be used for uploading new equipment data or editing existing equipment data.
* Let's first edit the file, and on the next step we'll re-upload it.

#### 2.1 Editing Existing Equipment

* As long as the value of the "Equipment ID" column is not empty, a row will be treated as "to be edited".
* Possible fields to be edited:
   * "Equipment Name" - for renaming an equipment instance.
      * If an equipment is being renamed and it has references to it down the CSV (as parent equipment for example) - edit those as well to the new name.
   * Every property
      * Property is editable as long as the corresponding equipment-type supports this property (can be verified on Inventory, under "configure"-> "equipment-types")

#### 2.2 Adding New Equipment

* As long as the value of the "Equipment ID" column is empty, a row will be treated as "to be added".
* Fields:
   * "Equipment Name"
   * "Equipment Type" - must be an existing name
   * List of locations, from big to small - (will be added on the fly if not exist)
   * "Parent Equipment (3)" - if exist
   * "Parent Equipment (2)" - if exist
   * "Parent Equipment" - if exist
   * "Parent Position" - if exist
   * List of properties for this equipment - if exist
* Equipment positions won't be added on the fly - they should exist and be free in advance of the new import run.

#### 2.3 Uploading the modified CSV



![](https://s3.amazonaws.com/purpleheadband.images/wiki/exported_data_for_upload.png)

* Now that we have our CSV ready - 
   *  Inventory
   * Click the '+' sign and a dialog will be opened.
   * Click the "Bulk Upload" tab
   * "Upload Exported Equipment"
   * Choose the edited file.