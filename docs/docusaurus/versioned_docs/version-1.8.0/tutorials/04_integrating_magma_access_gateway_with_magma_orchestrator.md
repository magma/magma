---
id: version-1.8.0-04_integrating_agw_with_orc8r
title: 4. Integrating Magma Access Gateway with Magma Orchestrator
hide_title: true
original_id: 04_integrating_agw_with_orc8r
---

# 4. Integrating Magma Access Gateway with Magma Orchestrator

## Integrate Magma Access Gateway with Magma Orchestrator

Offer an application endpoint from Orchestrator:

```console
juju offer orc8r.orc8r-nginx:orchestrator
juju consume orc8r.orc8r-nginx
```

Relate Magma Access Gateway with Orchestrator:

```console
juju relate orc8r-nginx:orchestrator magma-access-gateway-operator
```

Wait for the application to go back to `Active-Idle`:

```console
ubuntu@host:~$ juju status
Model  Controller     Cloud/Region   Version  SLA          Timestamp
edge   aws-us-east-2  aws/us-east-2  2.9.42   unsupported  16:09:01-05:00

App                            Version  Status  Scale  Charm                          Channel  Rev  Exposed  Message
magma-access-gateway-operator           active      1  magma-access-gateway-operator  stable    29  no

Unit                              Workload  Agent  Machine  Public address  Ports  Message
magma-access-gateway-operator/0*  active    idle   0        18.189.227.182

Machine  State    Address         Inst id                Series  AZ  Message
0        started  18.189.227.182  manual:18.189.227.182  focal       Manually provisioned machine
```

Fetch the Access Gateway's `Hardware ID` and `Challenge Key` and note those values:

```console
juju run-action magma-access-gateway-operator/0 get-access-gateway-secrets --wait
```

The output should look like:

```console
ubuntu@host:~$ juju run-action magma-access-gateway-operator/0 get-access-gateway-secrets --wait
unit-magma-access-gateway-operator-0:
  UnitId: magma-access-gateway-operator/0
  id: "22"
  results:
    challenge-key: MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE4bFCLDcHSi0fmESrejkTdJlBk/Mi/z/30VoV3dYTwWmOo1+xBjUjnMMBpWWlUbmdyOaSk32xg4/Pa9gq6gBj37INrB2zbgBfi5kdHbyFzbuIjak919/m5739tIb3NCYR
    hardware-id: 26236e99-f06d-4686-a888-696c7f2910c9
  status: completed
  timing:
    completed: 2023-03-17 11:47:24 +0000 UTC
    enqueued: 2023-03-17 11:47:20 +0000 UTC
    started: 2023-03-17 11:47:23 +0000 UTC
```

## Create a network in Magma Orchestrator

### Create a user in the `magma-test` organization

1. Login to the `host` Orchestrator organization at this address: `https://host.nms.<your domain>`.
2. Click on the :fontawesome-solid-user-plus: icon next to the `magma-test` organization
3. Add a user with the following attributes:
   - email: `admin@juju.com`
   - password: `password123`
   - role: `Super User`

### Create a network in the `magma-test` organization

1. Login to the `magma-test` organization at this address: `https://magma-test.nms.<your domain>`.
Use the credentials from the previous step.
2. On the left pane, click on "Networks"
3. Click on "Add Network"
4. Fill in the following values:
   - Network ID: `my-network`
   - Network Name: `my-network`
   - Description: `my-network`
   - Network Type: `lte`
5. Refresh the page. You should now see your network dashboard

### Change the Network configuration

1. Click on the "Networks" tab on the left pane
2. Next to the "EPC" box, click on "Edit"
3. Change the following values:
   - MCC: `001`
   - MNC: `01`
   - TAC: `7`
4. Click on "Save"

### Add the Access Gateway to the network

1. Navigate to "Equipment" via the left pane
2. Click on "Add New"
3. Fill in the following values:
   - Gateway Name: `my-gateway`
   - Gateway ID: `my-gateway`
   - Hardware UUID: `<Access Gateway Hardware ID>`
   - Gateway Description: `my-gateway`
   - Challenge Key: `<Access Gateway Challenge Key>`
4. Click on "Save and Continue". You should ignore the next tabs and continue clicking on "Save and Continue".
5. Click on "my-gateway"
6. You should see your gateway's health go to "Good" after a few minutes
