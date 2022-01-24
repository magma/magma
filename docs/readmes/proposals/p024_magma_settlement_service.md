# Proposal: Settlement Service for Magma

Author(s): [@arsenii-oganov]

Last updated: 01/14/2022

## 1.0 Objectives

The objective of this work is to extend Magma and the Helium blockchain to add support for settlement interfaces to bridge the gap between traditional roaming settlement and HIP 27s proposal of packet purchasing OUI managed by the helium network operator.
Software built to accomplish this (aka settlement domain service) will be open source under BSD-3-Clause and will reside in the github repository of the DeWi Alliance. Some changes / additions required to the Magma software will be committed to the Magma software repository under the governance of the Linux foundation, such that it can be effectively maintained in the future releases, but will remain openly available to all Helium community participants.
As a result of this effort, any MNO, MVNO or 3rd party Gateway Manufacturer will be able to use the software to operate their own roaming interconnect with a network roaming into Helium.
Out of Scope: The commercial aspects of how an MNO/MVNO purchases data credits on the Helium network through a session purchaser operator is outside the scope of this document and is a commercial arrangement between the clearing house, session purchaser operator and the MNO.

## 2.0 Background

The accounting and settlement of a Helium cellular network is unlike traditional cellular settlement. Unlike traditional cellular networks, a Helium cellular network contains radio and accounting entities (e.g. Miners) owned by private parties not affiliated with the core network operator and/or the subscriber home network provider. Additionally, the subscriber home network operator in a Helium network needs a reward path to pay individual radio/miner owners for the data consumed during subscriber sessions that traverse their deployment. A trust-limited accounting system must be implemented to ensure accounting data integrity and secure the reward path.
When creating this accounting system, the goal has been to avoid significant (or any) deviations from traditional 3GPP interfaces to allow maximum portability and lower the barriers to entry for vendors.
The systems centers around the new usage of a standard component, the Online Charging System (OCS). The OCS traditional role is quota management and realtime session management and accounting. In this system, the signaling of standard OCS is repurposed to provide control and accounting functions to a new Helium entity called the Session Purchaser.
The Session Purchasers job is to account and enforce sessions against a subscribers home network operator account on the helium network. The implementation below describes the role and signaling associated with a Session Purchaser.

## 3.0 Implementation scope

### 3.1 High - level solution architecture

![arch](https://user-images.githubusercontent.com/93994458/149513109-47b80417-7488-49c4-a2f0-ff94e6c57ebf.png)

### 3.2 End to end scenarios

Below will be described e2e scenarios needs to be implemented in scope of this project (Note: scenarios may be decomposed into different use-cases depends on actors involved and/or business case):

- State channel management procedure

The procedure describes Opening, Update and Closure of State channel management call flow.
The detailed call flow for the procedure is specified in section 3.2.1 “State channel management procedure”.

- UE attach procedure

The procedure in which the UE registers to the network, and receives the usage quota to create the EPS Bearer between the UE and the PGW, in order to be able to send and receive data.
The detailed call flow for the procedure is specified in section 3.2.2 “UE attach procedure”.

- UE quota update procedure

The procedure describes a next quota request by Threshold from the previously received quota from SGW to OCS.
The detailed call flow for the procedure is specified in section 3.2.3 “UE quota update procedure”.

- UE initiated session termination procedure

The procedure in which the UE initiates session termination. The procedure includes the last quota usage report from SGW to OCS.
The detailed call flow for the procedure is specified in section 3.2.4 “UE initiated session termination procedure”.

- Network initiated termination procedure

The procedure in which the Network (S/PGW) initiates session termination. The procedure includes the last quota usage report from SGW to OCS.
The detailed call flow for the procedure is specified in section 3.2.5 “Network initiated termination procedure”.

- TAP-Out generation and exporting procedure

The procedure to generate TAP3 files and transfer the files to the DCH (HPLMN) partner.

- RAP processing procedure

The procedure to return rejected TAP files and records to the visited network operator for corrections.

#### 3.2.1 State channel management procedure

Following diagram shows e2e State Channel life cycle management that will be in real life. It consists of existing functionality (marked with black), new functionality to be implemented in magma (marked with blue).

![State channel management procedure (1)](https://user-images.githubusercontent.com/93994458/149513469-0380c35c-d109-4c27-8014-d05b6cbd9a5e.png)

This scenario consists of the steps below:

- Step [1-2] The State channel opening procedure. The Helium Handler sends Open State Channel Request to open the State Channel. State Channel is opened once the positive Response is received from State Channel.
- Step [3] On this Step the Handler sends the Used Data Credits for each Used Data Quota per Miner pub key.
- Step [4-5] The State channel closure procedure. The closure procedure contains two type of triggers - all the Data credits have been used from the State channel and end of an Epoch.

#### 3.2.2 UE attach procedure

Following diagram shows e2e UE attach scenario that will be in real life. It consists of existing functionality (marked with black), new functionality to be implemented in magma (marked with blue).
The attach procedure implies that State Channel is already in open state.

![Attach call flow](https://user-images.githubusercontent.com/93994458/149513666-63ae0f4b-c8be-4cfc-a739-909723fca3aa.png)

This scenario consists of the steps below:

- Step [1 - 2] Initial attach procedure in a mobile network. The UE establishes radio link synchronization with an eNB, the UE creates a connection for data delivery via an Attach Request message sent to the eNB, the eNB then forwards the attach request to an MME/SGW.
- Step [3] The SGW sends CCR-I to the Session Purchaser requesting volume quota.
- Step [4] Session Purchaser checks if the IMSI-prefix (i.e. the HPLMN operator) has balance to answer with granted data units on CCR-I Request which has been received on Step [3].
- Step [5] In response to a CCR-I (Step [3]), the Session Purchaser returns a CCA-I message that indicates success (DIAMETER_SUCCESS) or failure (DIAMETER AUTHORIZATION REJECTED) depending on whether the HPLMN operator has sufficient credit for the requested services. In case of HPLMN operator has sufficient credit the CCA contains the Granted Service Unit in Bytes.
- Step [6] The SGW sends signed Session Init Event to Miner. The message is signed with SGW PubKey.
- Step [7] The Miner sends signed Session Init Event to Helium Handler. The message is signed with Miner PubKey.
- Step [8 - 9] The Create Session Request and Response messages to establish UE requested PDN connectivity through the SGW and Home PGW.
- Step [10-15] Establishing radio link between UE and eNodeB for packet data transfer.
- Step [16-17] Packet data between UE and Internet.

#### 3.2.3 UE quota update procedure

Following diagram shows e2e UE update scenario that will be in real life. It consists of existing functionality (marked with black), new functionality to be implemented in magma and helium (marked with blue).

![Update call flow](https://user-images.githubusercontent.com/93994458/149513805-f78e1c18-bb12-4f1f-abf9-0a0ea5e09d20.png)

This scenario consists of the steps below:

- Step [1-2] Packet data transfer is ongoing between UE and Internet.
- Step [3] The threshold of the quota is configured on SGW (e.g. 95% of granted service unit).
- Step [4] In this Step SGW sends Session Update Request with used service unit.
- Step [5] Session Purchaser checks if the HPLMN operator has balance to answer with granted data units on CCR-U which has been received on Step [4].
- Step [6] In response to a CCR-U (Step [5]), the Session Purchaser returns a CCA-U message that indicates success (DIAMETER_SUCCESS) or failure (DIAMETER AUTHORIZATION REJECTED) depending on whether the HPLMN operator has sufficient credit for the requested services.
- Step [7] The SGW sends signed data usage to Miner. The message is signed with SGW PubKey.
- Step [8] The Miner sends signed data usage to Helium Handler. The message is signed with Miner PubKey.
- Step [9] The Helium Handler updates State Channel with signed data usage received on Step [8]
- Step [10-11] Packet data continuation between UE and Internet.

#### 3.2.4 UE initiated session termination procedure

Following diagram shows e2e UE initiated session termination scenario that will be in real life. It consists of existing functionality (marked with black), new functionality to be implemented in magma and helium (marked with blue).

![UE Terminate call flow](https://user-images.githubusercontent.com/93994458/149513885-75facecc-54ef-4710-b6de-65b4881ddebb.png)

This scenario consists of the steps below:

- Step [1-2] UE detach procedure in a mobile network. The UE closes radio link with an eNB and MME/SGW.
- Step [3] In this Step SGW sends Session Termination Request with last used service unit.
- Step [4] The Session Purchaser update the HPLMN operator balance with last used data units.
- Step [5] The Session Purchaser sends Response on CCR-T (Step [3]) to close the session.
- Step [6] The SGW sends data usage to Miner. The message is signed with SGW PubKey.
- Step [7] The Miner sends data usage to Helium Handler. The message is signed with Miner PubKey.
- Step [8] The Helium Handler updates State Channel with signed data usage received on Step [7]
- Step [9] (Optional) The Helium Handler sends FreedomFi balance update to Session Purchaser.
- This steps checks for discrepancy between FF and Matrixx, helpful early, hopefully long term becomes less needed.
- Step [10-11] The Message between SGW to Home PGW once UE Requested PDN Disconnection to close the PDN connection.
- Step [12-13] MME/eNodeB closes radio link for the UE.

#### 3.2.5 Network initiated termination procedure

Following diagram shows e2e Network initiated termination that will be in real life. It consists of existing functionality (marked with black), new functionality to be implemented in magma and helium (marked with blue).

![Network Terminate call flow](https://user-images.githubusercontent.com/93994458/149514008-8c9ae8e1-1c55-44db-b0da-68972c3cd693.png)

This scenario consists of the steps below:

- Step [1] Delete session request comes from home PGW to close the PDN connection.
- Step [2] In this Step SGW sends Session Termination Request with last used service unit.
- Step [3] The Session Purchaser update the HPLMN operator balance with last used data units.
- Step [4] The Session Purchaser sends Response on CCR-T (Step [2]) to close the session.
- Step [5] The SGW sends data usage to Miner. The message is signed with SGW PubKey.
- Step [6] The Miner sends data usage to Helium Handler. The message is signed with Miner PubKey.
- Step [7] The Helium Handler updates State Channel with signed data usage received on Step [6]
- Step [8] (Optional) The Helium Handler sends FreedomFi balance update to Session Purchaser.
- Step [9] Delete session response from SGW to Home PGW to indicate the PDN connection is closed.
- Step [10-11] MME/eNodeB closes radio link for the UE.

#### 3.2.6 TAP-Out procedure

A roaming agreement between the home network operator and the visited network operator defines the terms that enable each other's customers access to the wireless networks. The visited network operator records the activities performed by the roaming subscriber and then sends the call event details to the home network operator in the format agreed upon in the roaming agreement, usually Transferred Account Procedure (TAP) format. TAP is the process that allows a visited network operator to send call event detail records of roaming subscribers to their respective home network operators to be able to bill for the subscriber's roaming usage.

#### 3.2.8 RAP-In procedure

The home network operator validates the data in the TAP files to ensure that it conforms to the TAP standard and to the terms of the roaming agreement. If the received TAP file contains any errors, the home network operator can reject the entire file or only the incorrect call event detail records. The incorrect file or records are returned to the visited network operator in a Returned Account Procedure (RAP) file.
RAP process is used to return rejected TAP files and records to the visited network operator for corrections. A RAP file contains the rejected TAP file or records and additional data about the error, such as the error code or the error cause. The visited network operator corrects the errors and sends the corrected TAP file back to the home network operator.

### 3.3 Integration interfaces

The following figure shows the integration interfaces between Magma domain and OCS /DCH.

![Integration interfaces](https://user-images.githubusercontent.com/93994458/149514277-062e0834-cb8c-49af-b0e4-867c4664c694.png)

#### 3.3.1 DIAMETER Gy interface specification

Interface transport protocol: TCP/Diameter/DCCA.
The 3GPP based Gy interface between Session Proxy and Session Purchaser to implement real-time quota management for the packet data service for the roaming subscribers. While this interface resembles Gy in almost every way, this is a non-standard use of the Gy protocol defined to allow accounting at the operator level and allow the session purchaser a control point to approve or deny sessions at initiation or at update points throughout a session.
The Session Purchaser is the Diameter Credit Control server, which provides the online charging data to the Session Proxy based on volume quota.
The connection between the Session Proxy (client) and Session Purchaser (server) is TCP based. There are a series of message exchanges to manage quota in real-time.

#### 3.3.2 MME/SGW ↔︎ Helium Miner

Session lifecycle events and usage reports should be delivered from each deployed AGW to the Network-Server (helium-handler), which uses that information to update a state channel. The Helium Miner container does such communication between hotspot(AGW) and Network Server. So, a trusted channel should be formed between MME/SGW and Miner to let a Miner container communicate with the Helium network on behalf of AGW. For having a trusted channel between AGW and Miner new component is introduced - a session-forwarder. Session-forwarder is a mediator between Magma-specific components deployed on AGW and a Helium Miner.
To let Miner trust Session-forwarder, session-forwarder is deployed on the same AGW where software is signed and Linux Secure Boot is used to verify all software components being loaded. So, Miner can trust to all received session-related messages. Miner just does forward all message to helium-handler.

![SGW-Miner](https://user-images.githubusercontent.com/93994458/149515220-2ed6ea77-bcb3-4322-a615-c23abc54f580.png)

1. Sessiond(Magma component) requires new sessions statistics gRPC stream API.
2. When new session is created or SessionState is updated on sessiond side, sessiond streams those updates via session statistics API to session-forwarder.
3. Session-forwarder talks to helium miner via gRPC.

Sessiond API update:

```service SessionStaticsStreamer {
rpc SessionStatistics(orc8r.Void) returns (stream SessionStatisticsUpdate) {}
}```

Miner API update (session-forwarder <-> Miner) :
```service api {
rpc session_usage(session_usage_update) returns (orc8r.Void) {}
}```

#### 3.3.3 Helium Miner ↔︎ Network Server (router/helium-handler)

For each 3gpp session Miner should forward session usage stats from Session-forwarder to Network Server.
```service api {
rpc session_usage(session_usage_update) returns (orc8r.Void) {}
}```

#### 3.3.4 helium-handler ↔︎ State channel

Helium-handler can straightly verify session_usage_update messages from Miner and trust them because they are signed by Session-forwarder using Miner key. So classic offer/purchase/reject semantics can be eliminated.

1. Initial session negotiation is done between miner and network server (helium-handler).
2. When session is in progress, miner forwards session_update messages which helium-handler adds as cbrs_session_usage_commit messages the state channel.
3. When session is terminated the last cbrs_session_usage_commit is added to the state channel for the session.

```message blockchain_state_channel_message_v1 {
oneof msg {
blockchain_state_channel_response_v1 response = 2;
blockchain_state_channel_packet_v1 packet = 4;
blockchain_state_channel_offer_v1 offer = 5;
blockchain_state_channel_purchase_v1 purchase = 6;
blockchain_state_channel_banner_v1 banner = 7; // DEPRECATED
blockchain_state_channel_rejection_v1 reject = 8;
blockchain_state_channel_cbrs_session_usage_commit_v1 cbrs_session_usage_commit = 9;
}
}```

#### 3.3.5 Session Forwarder ↔︎ CDR/TAP/RAP Engine

Session-Forwarder should send all session details to the CDR/TAP/RAP Engine to generate TAP files.

#### 3.3.6 CDR/TAP/RAP Engine ↔︎ DCH TAP handling system

Interface transport protocol: SFTP.

Roaming outcollect processing is using to track and rate activities of subscribers from other wireless networks that roam on FreedomFi network. Outcollect processing allows to rate the visiting subscribers' roaming usage using InterCarrier Tariff rates and generate TAP files consisting of the visiting subscriber's call event detail records, which should be send to roaming partners along with an invoice to bill them for their subscribers' roaming usage.
Roaming outcollect processing involves:

- Rating the visiting subscribers' roaming CDRs using InterCarrier Tariff rates specified in the roaming agreements and generating TAP files for each roaming partner.
- Handling errors in TAP files returned back in RAP files to FreedomFi from a roaming partner for corrections.

To validate the TAP file, the following types of validations are performed on DCH or/and HPLMN side.

**TAP3 fatal error validation:**
TAP3 fatal error validation is performed first to ensure all required data is present and valid. For example, if the TAP file is missing a required block, the entire file is rejected and written to a RAP file.

**TAP3 severe error validation:**
TAP3 severe error validation is performed if fatal error validation is successful. TAP records are validated to check for incorrect or missing reference data or content. For example, if a TAP record is missing a required field, the record is rejected and written to a RAP file, but all other TAP records in the file are processed.

The expected frequency of sending TAP3 files is two tap files in 24 hours.
The expected maximum limit for TAP3 files is 200,000 records, or 50 KB.

TAP3 File Logical Structure:

![TAP3](https://user-images.githubusercontent.com/93994458/149515412-cd140472-e549-4d12-bbcb-eed6860cd9dc.png)

Mandatory - blue
Conditional - orange
Optional - green

Only the GPRS Call will be used in the the Call Event Details section.

Note: A Notification file is sent where the transfer mechanism is electronic file transfer and there is no data available for transfer. (TBD with Tomia).
RAP File Logical Structure:

![RAP](https://user-images.githubusercontent.com/93994458/149515517-80f63415-1606-48c3-b3d0-263d5623a173.png)

Mandatory - blue
Conditional - orange
Optional - green

Note: Only one of the elements grouped at Return Detail level is applicable.

## 4.0 Roadmap and Schedule

**MS1:** Demonstrate producing CDRs via API call for stationary UEs (single miner whole session)

**MS2:** Demonstrate Session Purchaser governed network access

**MS3:** Demonstrate miner signed CDRs with validation and data recorded to State Channel

**MS4:** Demonstrate TAP record generation, export, import, and reconciliation.

**MS5:** Launch first operator with roaming interfaces integrated and recording to blockchain reconciled CDRs in settlement server managed state channels

**MS6:** Upstream code to master
