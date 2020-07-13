---
id: workforce-intro
title: Workforce Management
---
## Introduction
Workforce Management is a web application (included in Symphony platform) that allows Internet Service Providers (ISP) and Mobile Network Operators (MNO) to easily create, manage and track projects as well as to schedule tasks (work orders) and assign them to the field force. All this via a user-friendly web UI that enables the user to perform complex operations at the touch of a button. 
The user can access to Workforce Management via clicking on the “hamburger” button on the bottom left corner and selecting “Workforce Management” as displayed in the figure 1 below:  
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Workforce+management/figure+1.png' width=300> <br>
Fig.1: Workforce Management tab <br><br><br>
</p> 


On the left side there is a black vertical bar with three main tabs, which are “Work Orders”, “Projects”, and “Templates” (see fig.2). By default, the first tab (Work Orders) is selected. 
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Workforce+management/figure+2.PNG' width=400> <br>
Fig.2: Work Orders, Projects, and Templates tab <br><br><br>
</p> 

## Work Orders <br>
**Definition:**
"In Workforce Management, a “work order” (WO) is a task (or a set of tasks) that can be assigned to a worker (for example a field engineer or a technician) to perform a specific activity. Examples of WOs are “installing a router”, “verifying the correct alignment of a satellite antenna dish”, “replacing a faulty device”, etc."
<br><br>

### Web user interface
Looking at figure 3 below, Work Orders page includes:
- A list of all the already created work orders (1). This list is formatted as a set of rows that can be sorted alphabetically by the column fields below:    
  - Name
  - Template
  - Project
  - Owner
  - Status
  - Creation Time
  - Due Date
  - Location
  - Assignee
  - Close time

- A horizontal “filter bar” positioned on the top of the page (2) that allows the user to search for WOs matching a specific filter criteria. User can filter by 
  - WO data (name, status, owner, etc…)
  - Location Type (Country, State, City, etc…)
  - A pre-defined set of criteria (status Pending & Planned, Name & Owner, etc…)

- The “export” button that allows the user to export work orders data to a .CSV file (3). 
- A “view mode” button that displays the Work Orders as a table (4) (default) or as a map (5) based on their locations.
- On the top right corner, “Create Work Order” button allows the user to create a new work order (6).
- An existing Work Order can be Edited by clicking on the Work order’s name in the list (7)
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Workforce+management/figure+3.PNG' width=1280> <br>
Fig. 3:  Work Orders page <br><br><br>
</p> 

## Projects <br>
**Definition:**
"In Workforce Management, a “Project” is defined as a set of activities (or WOs) that are related to the same project. Examples of Project are “3G to 4G Network Upgrade in Uganda” (that consists, for instance, of “eNodeB installation” WOs in 400 different locations + “acceptance tests” WOs in 10 different cities + other types of WOs)"
<br><br>
### Web user interface
As displayed in fig.4 , the Projects page consists of:
- A list of all the already created projects (1). This list is formatted as a set of rows that can be sorted alphabetically by the fields indicated below:
  - Name
  - Template
  - Location
  - Owner
- A “view mode” button that displays the Projects as a table (2) (default) or as a map (3) based on their location.
- On the top right corner, “Create Project” button allows the user to create a new project (4).
- An existing project can be Edited by clicking on Project’s name in the list (5)
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Workforce+management/figure+4.PNG' width=1280> <br>
Fig. 4:  Projects page <br><br><br>
</p>

## Templates <br>
**Definition:**
"To create a new work order, the user must provide a set of data that depends on the nature of the task itself. The data structure (or schema) used by the user to define a WO is called “work order template”. In other words, a work order template consists of a set of data fields that can be re-used to create WOs of the same type. The user can create customized work order templates for a variety of different activities (or work orders)."

**Definition:**
"Like the “work order template”, a “Project template” consists of a set of data (including several types of WOs) that can be re-used in projects of the same type. For example, the same “3G to 4G Network Upgrade” project template can be used for both projects “3G to 4G Network Upgrade in Uganda” and “3G to 4G Network Upgrade in Kenya” (assuming the network vendor and the network operator is the same)."   
<br><br>
Creating a Work Order or a Project starts with selecting a work order template or a project template, respectively. In both cases, the selected template will include the required data structure for the Work Order or Project we want to create. Once a Work Order or a Project is created, it is no longer dependent on the template used upon creation, although its structure is exactly the same as inherited from its template at time of creation. 
If a WO is created using a certain template, once created, it is not possible to edit its data structure (add new properties, other fields). However, the user can still modify the values of the existing fields. <br><br><br>

### Templates Tab
Fig.5 shows the Templates tab that allow the user to create new or edit existing work order templates and project templates. Project templates enable to quickly create similar projects containing the same work order types.
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Workforce+management/figure+5.PNG' width=1280> <br>
Fig.5 Templates tab <br><br><br>
</p> 

### Work order template
If “Work Orders” is selected (as in Fig.6), the list of all the existing work order templates (1) is presented. This list is formatted as a set of rows that can be sorted alphabetically by the fields indicated below:
- Work order template’s name
- Description
New work order templates can be created by clicking on “Create Work Order Template” button (2).
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Workforce+management/figure+6.PNG' width=1280>
Fig.6 Work Order Templates page <br><br><br>
</p> 

### Project template
If “Projects” is selected (as in Fig.7), all existing project templates are presented (1). Each Project template includes the following information:
- Project template name
- Project template description
- Number of attached work orders <br>

New project templates can be created by clicking the Create Project Templates button (2). Existing project templates can be edited by clicking on the “Edit” button (3). Only projects templates with “no Work Orders” can be deleted (4).
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Workforce+management/figure+7.PNG' width=1280> <br>
Fig.7 Project Template page <br><br><br>
</p> 

