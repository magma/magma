# Pipelined

### Diagram

```
    GTP port            Local Port
     Uplink              Downlink
        |                   |
        |                   |
        V                   V
    -------------------------------
    |            Table 0          |
    |         GTP APP (OAI)       |
    |- sets IMSI metadata         |
    |- sets tunnel id on downlink |
    |- sets eth src/dst on uplink |
    -------------------------------
                  |
                  V
    -------------------------------
    |          Table 1            |
    |           inout             |
    |- sets direction bit         |
    -------------------------------
                  |
                  V
    -------------------------------
    |          Table 2            |
    |            ARP              |
    |- Forwards non-ARP traffic   |
    |- Responds to ARP requests w/| ---> Arp traffic - LOCAL
    |  ovs bridge MAC             |
    -------------------------------
                  |
                  V
    -------------------------------
    |          Table 3            |
    |       access control        |
    |- Forwards normal traffic    |
    |- Drops traffic with ip      |
    |  address that matches the   |
    |  ip blacklist               |
    -------------------------------
                  |
                  V
   Configurable apps managed by cloud
            (Tables 4-19)
                  |
                  V
    -------------------------------
    |          Table 20           |
    |           inout             |
    |- Forwards uplink traffic to |
    |  LOCAL port                 |
    |- Forwards downlink traffic  |
    |  to GTP port                |
    -------------------------------
        |                   |
        |                   |
        V                   V
    GTP port            Local Port
    downlink              uplink

```

### Configurable apps

```
    ------------------------------- 
    |          Table X            |
    |          metering           |
    |- Assigns unique flow id to  |
    |  IP traffic                 |
    |- Receives flow stats from   |
    |  OVS and forwards to cloud  |
    -------------------------------
    
    -------------------------------
    |          Table X            |
    |            DPI              |
    |- Assigns App ID to each new |
    |  IP tuple encountered       |
    |- Optional, requires separate|
    |  DPI engine                 |
    -------------------------------

    -------------------------------
    |          Table X            |
    |        enforcement          |
    |- Activates/deactivates rules|
    |  for a subscriber           |
    |- Keeps track of flow stats  |
    |  and sends to sessiond      |
    -------------------------------   
```

### Reserved registers
```
+----------+------------+----------------------+-----------------------------+
| Register |    Type    |         Use          |           Set by            |
+----------+------------+----------------------+-----------------------------+
| metadata | Write-once | Stores IMSI          | Table 0 (GTP application)   |
| reg0     | Scratch    | Temporary Arithmetic | Any                         |
| reg1     | Global     | Direction bit        | Table 1 (inout application) |
| reg2     | Local      | Policy number        | enforcement app             |
| reg3     | Local      | App ID               | DPI app                     |
+----------+------------+----------------------+-----------------------------+
```
