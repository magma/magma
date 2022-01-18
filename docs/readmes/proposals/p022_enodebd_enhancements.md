# Proposal: eNodeBd Enhancements for Cell PnP Support

Author(s): [@arsenii-oganov]

Last updated: 11/29/2021

## 1. Objectives

The objective of this work is to add services to address the logistics of scale radio deployments using Magma. **Specifically, plug-n-play onboarding and firmware update via TR-069**.

Plug-n-play onboarding enables auto-discovery and configuration of a new radio connected to a Magma network. This will enable less sophisticated users to add radios to a Magma managed network, as well making it quicker and easier for private cellular deployment users to setup and/or scale their networks.

Firmware update via TR-069 will allow the enodebd ACS to update the firmware on radios connected to the Magma network. Automated and Remote firmware updates are an important requirement for any scale deployment.

Software built to accomplish this will be open source under BSD-3-Clause license and will be committed to the Magma software repository under the governance of the Linux foundation, such that it can be effectively maintained in the future releases. The project will also enable endpoints for more advanced onboarding logic to allow for vendor extensions.

To demonstrate the developed framework, full support for onboarding, configuration, and firmware update will be demonstrated for specific radio models including:

- Baicells Nova 436Q
- Baicells Nova 430
- Baicells Neutrino 430
- Sercomm Englewood

## 2. Background

### 2.1 Terminology

1. **TR-069** refers to an application layer technical specification of the CPE WAN Management Protocol (CWMP) both specified by the Broadband Forum. This application is used for the remote management of configuration, firmware, and metrics in network equipment

2. **ACS** refers to the Automatic Configuration Server which is the TR-069 server in a server client model.

3. Magma **managed** eNBs refer to the radio devices that use TR-069 server communication using the enodebd service. **Externally managed** eNBs, in contrast, do not use TR-069 and are expected to be already configured to match the Cell ID and other configuration parameters so they successfully get connected to the AGW.

### 2.2 Plug-n-Play

The radio onboarding process for Magma currently follows a pre-provisioning workflow. A user with login credentials for the NMS must complete the following steps:

1. Add a new radio element to the NMS with the following fields:

    - Name
    - Serial number
    - Description
    - Externally managed
    - Device Class
    - Bandwidth
    - Cell ID
    - RAN Config (FDD/TDD)
    - EARFCNDL
    - Special Subframe Pattern
    - Subframe Assignment
    - PCI
    - TAC
    - Transmit Enable

2. Assign the radio to an AGW

After this process has been completed the enodebd process in the AGW will now respond to the initial TR-069 Inform message and attempt to configure the radio. Enodebd uses the serial number as the key for matching informs to known radio elements.

If an inform message arrives at enodebd with a serial number that is not known, the TR-069 session is terminated and no further action is taken.

### 2.3 Firmware Update

Magma TR-069 ACS functionality in enodebd currently does not support firmware update procedures.

## 3. Implementation

### 3.1 Plug-n-Play

#### 3.1.1 Block Diagram

![block diagram](https://user-images.githubusercontent.com/93994458/141427317-8e19f9b3-b789-4b8b-9117-c2a4024b9089.png)

#### 3.1.2 Call Flow

![onboarding flow](https://user-images.githubusercontent.com/93994458/141427664-16110fc4-61a4-4bfa-a7b0-b3f379c62150.png)

#### 3.1.3 PnP Scope of Change

##### 3.1.3.1 eNodeBD Service (Modified)

The enodebd service will be updated to support 2 modes of operation, auto-discovery enabled/disabled.

**Auto-Discovery Disabled**
When Auto-Discovery is disabled enodebd will operate as it does today. When a TR-069 inform message is received with a serial number unknown to enodebd, it will terminate the session without further action.

**Auto-Discovery Enabled**
When Auto-discovery is enabled, enodebd will receive a new inform, check it’s local configuration to see if the serial number is known, if it is not, it will forward the serial number and device ID field contents to Control Proxy with the endpoint specified as Onboarding Service.

##### 3.1.3.2 Control Proxy Service (Modified)

The Control Proxy service will be modified to add Onboarding Service onboardd to the service registry. This new endpoint will be used for transporting messages between enodebd and onboardd.

##### 3.1.3.3 Onboardd Service (New)

The Onboarding Service onboardd will be a new service added to support automated radio onboarding. The primary function will be to respond to enodebd new radio messages and using the API, create a new radio element and assign it to the source AGW.

In addition, when the new eNB is a CBRS CBSD, onboardd will also create a new radio element in the Domain Proxy (DP). When the radio is a Cat A indoor with a GPS derived location, the onboardd configuration will be sufficient for the DP to register and obtain a grant for the new eNB.
The onboardd service will only support magma Managed eNodebs, this is due to the requirement of establishing a TR-069 server session to further configure the radio device. Externally managed devices do not setup a TR-069 session.

**_Onboarding Logic_**
Onboardd will have 2 modes for onboarding, default and vendor.

**Default logic** will apply a default configuration parameter set when adding the radio via API. The default configuration will be defined in a configuration file. Default configurations will be matched based on Product Class parameter within the Device ID data. If no matching Product Class configuration file is found, a default configuration can be defined. If the default configuration file is empty, auto-discovery is rejected, otherwise settings within are used.

During the configuration, the Cell ID will need to be randomly generated and assigned to both the radio and the orchestrator radio element.

**_Vendor logic_** will be configurable to query an external service for radio configuration parameters. Onboardd will have a configuration for the endpoint to be used for vendor logic. If the endpoint configuration is populated, all new radio messages will be forwarded to the vendor logic service endpoint. Vendor logic will return radio configuration parameters for radio onboarding. Examples of vendor logic may include allow/block lists of serial numbers or assignment of serial number specific radio configuration parameters.

_Vendor Logic API_
To ensure ease of adoption and integration, a API will be defined for the calls to external vendor logic endpoints. The API will include at least radio_parameters_request and radio_parameters_response endpoints.
Domain Proxy Integration
The new onboardd service will support a set of CBRS configuration elements sufficient for adding a radio to the Domain Proxy. In certain cases, the information will be sufficient to also register a new eNB with the CBRS SAS and allow operation without Certified Professional Installer interaction. In other cases, the eNB will be populated with the majority of the data, making it faster and easier for a CPI to quickly review new radios. CBRS Configuration elements will include:

- SAS URL (TBD, may be global)
- SAS User ID (TBD, may be global)
- CBSD Category
- Location (indoor/outdoor)
- Height Type
- Channel Type (GAA/PAL) (TBD, may be global)
- FCC ID
- Radio Technology (LTE/NR)

This dataset will be supported via static files in default configuration mode and via the vendor logic API.

### 3.2 Firmware Update

#### 3.2.1 Block Diagram

![firmupdate diagram](https://user-images.githubusercontent.com/93994458/141432926-936ef4c8-c282-4717-851d-dab450377cfe.png)

#### 3.2.2 Call Flows

User Configuration

![userupdate flow](https://user-images.githubusercontent.com/93994458/141433768-50794a54-7db8-49e3-9457-cac5e7ecdb8e.png)

FW Update Applied to Radio

![fwupdate flow](https://user-images.githubusercontent.com/93994458/141433791-72954447-00af-4a14-b7a7-79def6f09424.png)

#### 3.2.3 Firmware Update Scope of Changes

##### 3.2.3.1 NMS Service (Modified)

The NMS will be updated to allow creation of eNB upgrade tiers including the selection of applicable Device Classes. NMS will also be extended to allow configuration of upgrade tier parameters including:

- FW image URL
- FW metadata
- Username (default = null)
- Password (default = null)
- FileSize
- TargetFileName
- DelaySeconds (default = 0)
- Md5 (Baicells only)
- RawMode (Baicells only, default = false)

In the case of vendor specific configurations, the NMS will allow user selection to include / exclude these elements in the metadata associated with an upgrade tier.
The NMS will also be updated to display FW version for each managed eNB elements.

##### 3.2.3.2 Configurator/lte Service (Modified)

The orchestrator configurator (and/or lte) service will add support for FW image URL and metadata to be added to the mconfig for LTE AGWs. Metadata will include:

- Username
- Password
- FileSize
- TargetFileName
- DelaySeconds
- Md5
- RawMode

##### 3.2.3.3 enodebd Service (Modified)

The enodebd service will be extended to support the new mconfig format with the FW image URL and metadata elements. Enodebd will also add support for the TR-069 FW update call flow shown above using the Download() method and metadata as shown. This method is used by all target radios. Moreover, this method is the basic FW update method defined in TR-069 protocol and is likely to be supported broadly across more TR-069 capable radios.

eNodebd will be extended to capture and report eNB firmware version reported via TR-069 to the orchestrator database. This will be used as the source for version displayed in the NMS.

Update logic will be implemented as part of Device Classes in enodebd so that any radio specific variations can be scope limited to that Device Class, similar to the FSM and data model today.

## 4. Roadmap & Schedule

**Scope items:**
**MS1:** Demonstrate manually triggered firmware update via enodebd TR-069 ACS (Sercomm and 1 Baicells Radio)

**MS2:** Demonstrate configurator/lte service creation of updated AGW mconfig triggering firmware update of both radios.

**MS3:** Complete NMS updates.

**MS4:** Demonstrate PoC implementation with enodebd talking directly to vendor onboarding server, performing auto-discovery of radio

**MS5:** Demonstrate onboardd addition and deployment inside Magma orchestration, including addition to Control Proxy Service Registry.

**MS6:** Demonstrate message flow through Control Proxy performing Default Logic radio addition using all magma integrated services.

**MS7:** Demonstrate radio auto-discovery for the 4 radio models shown.

**MS8:** Launch Magma integrated PnP in the FreedomFi/Helium network using Vendor Logic extension.

## 5.Schedule

TBD
