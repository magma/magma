# Magma extensions for inbound roaming

Author(s): [@arsenii-oganov]

Last updated: 02/18/2022

## 1 Background and Objective

The objective of this proposal is to extend Magma components to add support for settlement interfaces to bridge the gap between traditional roaming settlement based on Gy interface on Home PGW and new requirements of packet purchasing managed by the Magma core.

Currently, the Magma Access Gateway (AGW) implements a merged SGW+PGW (SPGW) task with support for inbound roaming. The goal of this proposal is to involve sessiond task to support Gy CCR-I/CCR-U messages for roaming subscribers in AGW and FeG.
Software built to accomplish this will be open source under BSD-3-Clause license and will be committed to the Magma software repository under the governance of the Linux foundation, such that it can be effectively maintained in the future releases.

## 2 Implementation Scope

### 2.1 Network Architecture

The planned network architecture to use sessiond for generate an initial CCR. The initial CCR should be generated to check IMSI balance before establish default bearer.

![arch](https://user-images.githubusercontent.com/93994458/154483331-c911ba83-1e9a-4180-b369-a0f4f2574bec.png)

S6a: No changes for S6a are planned in this proposal.
S8 Control Plane: Only one change is planned in this proposal: MME session setup logic should be changed to send Create Session Request after/when Sessiond received quota from OCS. In case of no quota has been received the S1 session should be terminated by the MME.
S8 User Plane: No changes for S8 User Plane are planned in this proposal. We are still going to use MME to install the Table 0 for roaming subscribers.
Sessiond: The main change we are planning to make here is to involve sessiond in the inbound roaming call flow.
The new logic should be added to the MME task:
Before sending Create Session Create request to S8 Proxy the MME should send “Create Session” message to the Sessiond. This message will be trigger to generate CCR-I message towards session-proxy with first volume quota request. When the quota has been received from OCS the sessiond should send "Create Session Response" towards MME to allow S8 crate session procedure. In case of no quota has been received from OCS the sessiond should trigger session termination to the MME.
Once MME received Create Bearer Request and Table0 has been provisioned to the OVS the MME should sent "Final Create Session Request" towards sessiond to start get stats procedure from OVS. After reaching the configurable volume threshold on sessiond the sessiond should generate CCR-U to ask the next quota in OCS. This is an existing procedure for home subscribers but we should add the procedure to the inbound roaming call flow.

### 2.1 Call Flows

The colour scheme in the diagram is as follows:
Black: The existing messages for inbound roaming
Red: The new expected messages
Blue: The existing messages for local users which have to be added to the inbound roaming procedure

#### 2.1.1 UE Attach call flow (successful case)

![attach](https://user-images.githubusercontent.com/93994458/154483466-712c211d-7754-47b2-9993-4dee2da5d052.png)

#### 2.1.2 UE Attach call flow (unsuccessful case)

The behaviour if quota has not been received from OCS.

![negative attach](https://user-images.githubusercontent.com/93994458/154483571-3b7bd11d-51fc-4375-95f6-af3959191fe8.png)

Note: In case if quota has not been received from OCS the MME should terminate the established session using standard GTP-С termination procedure.

#### 2.1.3 Quota Update Call Flow

![update flow](https://user-images.githubusercontent.com/93994458/154483741-db7b1935-2d1d-4b61-9719-3135b0910cf9.png)

#### 2.1.4 Session Termination due to quota exhausted call flow

The session termination procedure should be the same as is for local subscribers.

![termination flow](https://user-images.githubusercontent.com/93994458/154483837-da479f5a-02e6-4147-b8ae-7647917e01a1.png)

## 3 Roadmap and schedule

![roadmap](https://user-images.githubusercontent.com/93994458/154483942-5a96e6e3-8aaa-4465-81a7-0eb860d8019a.png)
