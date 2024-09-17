---
id: debug_user_control_plane
title: User Control Plane trace CLI
hide_title: true
---

# User Control Plane trace CLI

User control plane trace cli should facilitate filtering user control sessions in mme logs.
The cli has 3 functionalities(list_imsi, list_ue_id, session_trace) for a given IMSI, mme ue id and mme log.

## Usage

CLI outputs user control signaling for a supplied IMSI and mme user id.

```user_trace_cli.py [-h] [-p PATH] {list_imsi,list_ue_id,session_trace}```

Positional arguments:

```bash
{list_imsi,list_ue_id,session_trace}
  list_imsi           List imsi and number of ocurrences in logs
  list_ue_id          List mme and enb ue id pairs for a given imsi
  session_trace       Dump session trace for given user id and imsi

optional arguments:
  -h, --help            show this help message and exit
  -p PATH, --path PATH  Path to the data directory (default: /var/log/mme.log)
```

In general, you can use below flow for this cli.

1. Grab one imsi from list_imsi
2. From the IMSI selected, get the ue_id found from list_ue_id
3. Output with session_trace using imsi + mme_id found in Step 1 and 2.

For example:

```bash
user_trace_cli.py -p mme.log list_imsi

IMSI            Occurrences
------------------------------
000000000000110 331
000000000000111 2
000000000000112 2
000000000000113 2
000000000000114 2
000000000000115 2
000000000000116 2
000000000000117 2
000000000000118 2
000000000000119 2
```

```sh
user_trace_cli.py -p mme.log list_ue_id -i 000000000000110

IMSI: 000000000000110

mme_id          enodeb_id
-------------------------
25 0x19         728 0x2d8
26 0x1a         729 0x2d9
27 0x1b         730 0x2da
28 0x1c         731 0x2db
29 0x1d         732 0x2dc
```

```sh
user_trace_cli.py -p mme.log session_trace -i 200010001018110 -m 25

...
...

000418 Mon Aug 16 16:59:41 2021 7F271E3DF700 INFO  MME-AP tasks/mme_app/mme_app_context.c :0121   [000000000000110] Deleted UE location from directoryd

000422 Mon Aug 16 16:59:41 2021 7F271E3DF700 WARNI MME-AP tasks/mme_app/mme_app_context.c :0335   [000000000000110]  No IMSI hashtable for this IMSI

****************************************************************************************************

Error: Wrong APN configuration

Suggestion: Verify the APN has been provisioned correctly in the phone/CPE and/or Orc8r/HSS

Log: 000401 Mon Aug 16 16:59:41 2021 7F271E3DF700 ERROR NAS-ES tasks/nas/emm/sap/emm_cn.c      :0280    No suitable APN found ue_id=0x00000019)

****************************************************************************************************
```

## What is this command doing?

This command is using regex to match different patterns of IMSI and mme user id in mme logs in order to filter specific session from a given IMSI.

## How to read the output?

- 'list_imsi'
Provide a list of imsi and number of occurrences found in logs. (No input required)

- 'list_ue_id'
List mme and enb ue id pairs for a given IMSI. (IMSI input required)

- 'session_trace'
Output session trace for given user id(mme_id) and IMSI (IMSI and mme_id required). If is a known issue, it will provide the error and suggestions. Two known errors have been added(APN error and IE not supported).
