---
id: csv-upload
title: Importing Inventory Data    
---

## Step 1: Manual Definitions

(See Figure 1)

The following must be defined manually before uploading inventory data via CSV:

- Locations Types - using the Catalog (1), verify that the location types are defined by their order, under the LOCATIONS tab (2), drag and drop the location types to arrange them in a hierarchical order - from largest to smallest.
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+1.png' width=800> <br>
Figure 1: Locations 
</p>

- Equipment Types
- Port Types

## Step 2: Create a CSV Template

Export the relevant data to CSV (which will be the template for the upload) according to the steps described in the Exporting Your Data section (add a link TBD).

**note** The CSV file should be of UTF-8 format.

There are five different upload flows for exported data:

- Exported Equipment (for add and edit)
- Exported Ports (edit only)
- Exported Links (add &amp; edit)
- Exported Locations (edit only)
- Exported Services (add &amp; edit)

## Step 3: Importing data

Once the CSV file is exported, it can be modified and used in order to upload new records or edit existing ones.

### Equipment

The following sections describe equipment data upload information and procedures.

#### Step 1: Editing Existing Equipment

(See Figure 2 to Figure 5)
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+2.png' width=1000> <br>
Figure 2: Equipment ID Column
</p>

- Equipment ID (1):
- If the value of the Equipment ID column is not empty, a row will be treated as to be edited.
- As long as the value of the Equipment ID column is empty, a row will be treated as to be added.

- Possible fields to be edited:
- Equipment Name (2) - for renaming an equipment instance.
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+3.png' width=1000> <br>
Figure 3: Equipment Name
</p>

**note** If an equipment is being renamed and it has references to it farther down the CSV (as parent equipment for example) - edit those as well to reflect the new name.

- External ID (3).
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+4.png' width=1000> <br>
Figure 4: External ID Column
</p>

- Every property (4).

A Property is editable as long as the corresponding Equipment Type supports it.

Correspondence can be verified under Catalog/EQUIPMENT/Equipment Types (5).
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+5.png' width=1000> <br>
Figure 5: Properties
</p>

#### Step 2: Adding New Equipment

(See Figure 6)

Fill/edit the following fields in the CSV:

- Equipment Name
- Equipment Type - must be an existing name
- External ID - (optional) an ID from a third-party system
- List of locations, hierarchical order - from largest the smallest.

**note** A new location will be added automatically without location properties.

- Parent Equipment (3) - if exists
- Position (3) - if exists
- Parent Equipment (2) - if exists
- Position (2) - if exists
- Parent Equipment - if exists
- Equipment Position - if exists

**note** Equipment positions will not be added automatically, they must exist and be available in advance of the new import run.
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+6.png' width=1000> <br>
Figure 6: Adding New Equipment
</p>

#### Step 3: Uploading the modified CSV

(See Figure 7 to Figure 11)

In order to upload the modified CSV, perform the following steps:

##### 1) Open Inventory

_Result:_ The **Locations** section appears:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+7.png' width=600> <br>
Figure 7: Locations
</p>
          
##### 2) Click the upper **+** sign (1).

_Result:_ The following window appears **:**
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+8.png' width=600> <br>
Figure 8: BULK Upload Tab
</p>
          
##### 3) Click the **BULK UPLOAD** tab (2).

_Result:_ The **BULK UPLOAD** options appear:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+9.png' width=600> <br>
Figure 9: Upload Exported Equipment
</p>
          
##### 4) Click the **Upload Exported Equipment** option (3).
          
##### 5) Choose a file to upload.

_Result:_ In case of an error, a warning window with a message (4) will appear preventing the file upload:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+10.png' width=600> <br>
Figure 10: Warning Window
</p>

**note** The displayed message (4) will change according to the error.

            1. Open the uploaded file and fix the error.
            2. Repeat steps **1** to **5** as they are described in this section.

- If the file can be uploaded but there are issues, a window with a message (4) and a list of issues (5) will appear:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+11.png' width=600> <br>
Figure 11: List of Issues
</p>

            1. Open the uploaded file and fix the issues (optional).
            2. Repeat steps **1** to **5** as they are described in this section.
            3. Click the **Upload Anyway** button (6).

_Result:_ A summary window appears:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+12.png' width=1000> <br>
Figure 12: Summary Window
</p>

            1. Review what was uploaded (7).
            2. Click the **OK** button (8).

_Result:_ The data is uploaded.

**note** IF there are no errors or issues the summary window will appear and the file will be uploaded.

### Ports

The following sections describe ports data upload information and procedures.

#### Step 1: Editing Existing Ports

(See Figure 13 and Figure 14)
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+13.png' width=1000> <br>
Figure 13: Port ID Column
</p>

- Port ID (1):
- It is mandatory to add a value in the Port ID column.
- Each row with a Port ID value will be treated as to be edited.

- Every port property (2):

Properties are editable as long as the corresponding Port Type supports it.

Correspondence can be verified under Catalog/PORTS/Port Types/Port Properties (3).

**note** Link properties under Port Type are not supported
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+14.png' width=1000> <br>
Figure 14: Properties
</p>

#### Step 2: Uploading the modified CSV

(See Figure 15 to Figure 20)

In order to upload the modified CSV, perform the following steps:

##### 1) Open Inventory

_Result:_ The **Locations** section appears:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+15.png' width=600> <br>
Figure 15: Locations
</p>
          
##### 2) Click the upper **+** sign (1).

_Result:_ The following window appears **:**
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+16.png' width=600> <br>
Figure 16: BULK Upload Tab
</p>

##### 3) Click the **BULK UPLOAD** tab (2).

_Result:_ The **BULK UPLOAD** options appear:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+17.png' width=600> <br>
Figure 17: Upload Exported Ports
</p>

##### 4) Click the **Upload Exported Ports** option (3).

##### 5) Choose a file to upload.

_Result:_ In case of an error, a warning window with a message (4) will appear preventing the file upload:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+18.png' width=600> <br>
Figure 18: Warning Window
</p>

**note** The displayed message (4) will change according to the error.

            1. Open the uploaded file and fix the error.
            2. Repeat steps **1** to **5** as they are described in this section.

- If the file can be uploaded but there are issues, a window with a message (4) and a list of issues (5) will appear:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+19.png' width=600> <br>
Figure 19: List of Issues
</p>

            1. Open the uploaded file and fix the issues (optional).
            2. Repeat steps **1** to **5** as they are described in this section.
            3. Click the **Upload Anyway** button (6).

_Result:_ A summary window appears:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+20.png' width=600> <br>
Figure 20: Summary Window
</p>

            1. Review what was uploaded (7).
            2. Click the **OK** button (8).

_Result:_ The data is uploaded.

**note** If there are no errors or issues the summary window will appear, and the file will be uploaded.

### Links

The following sections describe links data upload information and procedures.

#### Step 1: Editing Existing Links

(See Figure 21 and Figure 22)
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+21.png' width=1000> <br>
Figure 21: Link ID Column
</p>

- Link ID (1):
- If the value of the Link ID column is not empty, a row will be treated as to be edited.
- As long as the value of the Link ID column is empty, a row will be treated as to be added.

- Every Link Property (2):

Properties are editable (2) as long as the corresponding Port Type supports it.

The correspondence can be verified under Catalog/Port Types/Link Properties (3).

**note** Port properties under Port Type are not supported
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+22.png' width=1000> <br>
Figure 22: Properties
</p>

#### Step 2: Adding New Links

(See Figure 23)

Fields for each one of the ports (example for &quot;A&quot; but the same for &quot;B&quot;):

- Port A Name - must be a valid port name from the Equipment Type.
- Equipment A Name

**note** A new Equipment Name will be added automatically.

- Equipment A Type - must be an existing name.
- List of locations, hierarchical order - from largest to smallest.

**note** A new location will be added automatically without location properties.

- Parent Equipment (3) A - if exists.
- Position (3) A - if exists.
- Parent Equipment (2) A - if exists.
- Position (2) A - if exists.
- Parent Equipment A - if exists.
- Equipment Position A - if exists.
- List of link properties taken from both ports for this equipment - if exist.

**note** Equipment positions will not be added automatically, they must exist and be available in advance of the new import run.
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+23.png' width=1000> <br>
Figure 23: Adding New Equipment
</p>

#### Step 3: Uploading the modified CSV

(See Figure 24 to Figure 29)

In order to upload the modified CSV, perform the following steps:

##### 1) Open Inventory

_Result:_ The **Locations** section appears:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+24.png' width=600> <br>
Figure 24: Locations
</p>

##### 2) Click the upper **+** sign (1).

_Result:_ The following window appears **:**
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+25.png' width=600> <br>
Figure 25: BULK Upload Tab
</p>

##### 3) Click the **BULK UPLOAD** tab (2).

_Result:_ The **BULK UPLOAD** options appear:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+26.png' width=600> <br>
Figure 26: Upload Exported Links
</p>

##### 4) Click the **Upload Exported Links** option (3).

##### 5) Choose a file to upload.

_Result:_ In case of an error, a warning window with a message (4) will appear preventing the file upload:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+27.png' width=600> <br>
Figure 27: Warning Window
</p>

**note** The displayed message (4) will change according to the error.

            1. Open the uploaded file and fix the error.
            2. Repeat steps **1** to **5** as they are described in this section.

- If the file can be uploaded but there are issues, a window with a message (4) and a list of issues (5) will appear:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+28.png' width=1000> <br>
Figure 28: List of Issues
</p>

            1. Open the uploaded file and fix the issues (optional).
            2. Repeat steps **1** to **5** as they are described in this section.
            3. Click the **Upload Anyway** button (6).

_Result:_ A summary window appears:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+29.png' width=1000> <br>
Figure 29: Summary Window
</p>

            1. Review what was uploaded (7).
            2. Click the **OK** button (8).

_Result:_ The data is uploaded.

**note** If there are no errors or issues the summary window will appear and the file will be uploaded.

### Locations

The following sections describe locations data upload information and procedures.

#### Step 1: Editing Existing Locations

(See Figure 30 to Figure 31)
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+30.png' width=1000> <br>
Figure 30: Location ID Column
</p>
- Location ID (1):
- If the value of the Location ID column is not empty, a row will be treated as to be edited.
- As long as the value of the Location ID column is empty, a row will be treated as to be added.

- Possible fields to be edited:
- External ID (2).
- Latitude (3).
- Longitude (4).
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+31.png' width=600> <br>
Figure 31: External ID, Latitude, Longitude
</p>

#### Step 2: Adding New Location

(See Figure 32)

Fill/edit the following fields in the CSV:

- List of locations, hierarchical order - from largest to smallest.

**note** A new location will be added automatically without location properties.

- External ID - (optional) an ID from a third-party system
- Latitude - to specify the coordinates of the location.
- Longitude - to specify the coordinates of the location.
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+32.png' width=1000> <br>
Figure 32: Adding New Location
</p>

#### Step 3: Uploading the modified CSV

(See Figure 33 to Figure 38)

In order to upload the modified CSV, perform the following steps:

##### 1) Open Inventory

_Result:_ The **Locations** section appears:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+33.png' width=600> <br>
Figure 33: Locations
</p>

##### 2) Click the upper **+** sign (1).

_Result:_ The following window appears **:**
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+34.png' width=600> <br>
Figure 34: BULK Upload Tab
</p>

##### 3) Click the **BULK UPLOAD** tab (2).

_Result:_ The **BULK UPLOAD** options appear:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+35.png' width=600> <br>
Figure 35: Upload Exported Locations
</p>

##### 4) Click the **Upload Exported Locations** option (3).

##### 5) Choose a file to upload.

_Result:_ In case of an error, a warning window with a message (4) will appear preventing from uploading the file:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+36.png' width=600> <br>
Figure 36: Warning Window
</p>

**note** The displayed message (4) will change according to the error.

            1. Open the uploaded file and fix the error.
            2. Repeat steps **1** to **5** as they are described in this section.

- IF the file can be uploaded but there are issues, a window with a message (4) and a list of issues (5) will appear:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+37.png' width=1000> <br>
Figure 37: List of Issues
</p>

            1. Open the uploaded file and fix the issues (optional).
            1. Repeat steps **1** to **5** as they are described in this section.
            2. Click the **Upload Anyway** button (6).

_Result:_ A summary window appears:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+38.png' width=1000> <br>
Figure 38: Summary Window
</p>

            1. Review what was uploaded (7).
            2. Click the **OK** button (8).

_Result:_ The data is uploaded.

**note** IF there are no errors or issues the summary window will appear and the file will be uploaded.

### Services

The following sections describe services data upload information and procedures.

#### Step 1: Editing Existing Services

(See Figure 39)
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+39.png' width=1000> <br>
Figure 39: Location ID Column
</p>

- Service ID (1):
- If the value of the Service ID column is not empty, a row will be treated as to be edited.
- As long as the value of the Service ID column is empty, a row will be treated as to be added.

#### Step 2: Adding New Service

(See Figure 40)

Fill/edit the following fields in the CSV:

- Service ID - an internal id of the service in inventory, must be unique.
- Service Name - must be unique.
- Service Type - a name of the type of the service (defined in Catalog/SERVICES/Service Type).
- Service External ID - used to identify the service in other systems (CRM for example), unique, can be empty.
- Customer Name - must be unique, can be empty.
- Customer External ID - used to identify the customer in other systems (CRM for example), unique, can be empty.
- Status - four options - PENDING, IN\_SERVICE, MAINTENANCE, DISCONNECTED.
- List of properties for this service.
- List of locations, hierarchical order - from largest to smallest.

**note** A new location will be added automatically without location properties.

<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+40.png' width=1000> <br>
Figure 40: Adding New Service
</p>

#### Step 3: Uploading the modified CSV

(See Figure 41 to Figure 46)

In order to upload the modified CSV, perform the following steps:

##### 1) Open Inventory

_Result:_ The **Locations** section appears:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+41.png' width=600> <br>
Figure 41: Locations
</p>

##### 2) Click the upper **+** sign (1).

_Result:_ The following window appears **:**
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+42.png' width=600> <br>
Figure 42: BULK Upload Tab
</p>

##### 3) Click the **BULK UPLOAD** tab (2).

_Result:_ The **BULK UPLOAD** options appear:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+43.png' width=600> <br>
Figure 43: Upload Exported Service
</p>

##### 4) Click the **Upload Exported Service** option (3).

##### 5) Choose a file to upload.

_Result:_ In case of an error, a warning window with a message (4) will appear preventing from uploading the file:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+44.png' width=600> <br>
Figure 44: Warning Window
</p>

**note** The displayed message (4) will change according to the error.

            1. Open the uploaded file and fix the error.
            2. Repeat steps **1** to **5** as they are described in this section.

- If the file can be uploaded but there are issues, a window with a message (4) and a list of issues (5) will appear:
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+45.png' width=800> <br>
Figure 45: List of Issues
</p>

            1. Open the uploaded file and fix the issues (optional).
            2. Repeat steps **1** to **5** as they are described in this section.
            3. Click the **Upload Anyway** button (6).

_Result:_ A summary window appears:

<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Importing+data/Figure+46.png' width=1000> <br>
Figure 46: Summary Window
</p>

            1. Review what was uploaded (7).
            2. Click the **OK** button (8).

_Result:_ The data is uploaded.

**note** IF there are no errors or issues the summary window will appear and the file will be uploaded.
