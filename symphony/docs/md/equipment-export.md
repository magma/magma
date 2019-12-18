---
id: equipment-export
title: Exporting your equipment data
---

### 1. Before Exporting your data

* Verify that the locations types are defined by their order:
   * Under the location types tab, drag and drop the locations types so that it's arranged from big ones to small ones.

### 2. Exporting the data

* Go to the 'Equipment Search' tab (magnifying glass icon).
* Use the filters bar to get a subset of your result (optional).
* Click the "Export" button on the top right corner.

 
![](https://s3.amazonaws.com/purpleheadband.images/wiki/exportdata.png)



* A CSV file containing the filtered equipment list will be downloaded.
* Every row represents an equipment , and it will be of the following form and in the following order:
   * "Equipment ID"
   * "Equipment Name"
   * "Equipment Type"
   * List of locations, from big to small
   * "Parent Equipment (3)" grand grand parent of direct parent of the equipment
   * "Parent Equipment (2)"  grand parent of the equipment
   * "Parent Equipment" direct parent of the equipment
   * "Parent Position" position name under the direct parent (e.g. `slot #1`)
   * List of properties for this equipment.
* For example -  

![](https://s3.amazonaws.com/purpleheadband.images/wiki/exported.png)    


### 3. Import new/edited data using the same template
* [Check the upload wiki](csv-upload.md#importing-exported-data)

