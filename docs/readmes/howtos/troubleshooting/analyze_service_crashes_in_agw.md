---
id: analyze_service_crashes_in_agw
title: Analyze Service crashes in AGW
hide_title: true
---
# Analyze Service crashes in AGW

**Description:** Purpose of this guide is to identify if an outage or service restart (mme, sessiond) was caused due
to an AGW service crash.

**Environment:** AGW deployed on bare metal

**Components:** AGW

**Triaging steps:**

**1.  Identify and quantify the impact.** Verify if there is any temporary or permanent drop in metrics related to
service:
  - Number of Connected eNBs (Grafana -> Dashboards -> Networks)

  - Network of Connected UE (Grafana -> Dashboards -> Networks)

  - Network of Registered UE (Grafana -> Dashboards -> Networks)

  - Attach/ Reg attempts (Grafana -> Dashboards -> Networks)

  - Attach Success Rate (Grafana -> Dashboards -> Networks)

  - S6a Authentication Success Rate (Grafana -> Dashboards -> Networks)

  - Service Request Success Rate (Grafana -> Dashboards -> Networks)

  - Session Create Success Rate (Grafana -> Dashboards -> Networks)

  - Upload/Download Throughput (Grafana -> Dashboards -> Gateway)

    Note: Number of sites(enodeb) down, users affected, and outage duration are key indicators of service impact.

**2. Use metrics to verify if there was a recent service restart** and if this correlates with the metrics degradation. Get an estimated timestamp from the metrics. You can use the following metrics:

  - unexpected_service_restarts

  - service_restart_status

**3. Verify if the restart was intentionally triggered** by a user (AGW reboot)

  - Use `last reboot` to list the last logged in users and system last reboot time and date.

    Use this timestamp and compare with the timestamp in the metrics degradation to confirm both events are related(Consider the time zone difference between Orc8r and AGW.)

    If the service was not intentionally restarted. Follow below steps to confirm the outage was caused due to a service crash.

**4. Capture the service crash syslogs and coredumps**. Use the approximate time in metrics and look for the events in both syslogs and coredumps

  - In syslogs located in `/var/log`, look for service terminating events and its previous logs. For example, in below mme service crash there is a segfault reported in mme service before its being terminated.

  ```
  Dec 5 22:25:55 magma kernel: [266759.489500] ITTI 3[13887]: segfault at 1d6d80 ip 000055b0080da0c2 sp 00007f529e6c0310 error 4 in mme[55b0077bd000+e79000]
  Dec 5 22:25:59 magma systemd[1]: magma@mme.service: Main process exited, code=killed, status=11/SEGV
  ```

  - Service crashes with a segmentation fault will create coredumps in `/var/core/` folder. Verify if coredumps have been created and obtain the coredump that matches the time of the outage/crash. Depending on the type of service crash the name of the coredump will vary. More detail in https://magma.github.io/magma/docs/lte/dev_notes#analyzing-coredumps

**5. Get the backtrace using the coredumps**. To analyze the coredumps, you need 3 requirements.


  - Coredump file

  - service binary, (ie. for mme `/usr/local/bin/mme`)

  - gdb package installed

You can read the coredumps in any AGW or machine with these requirements. Example:  `gdb mme core-1607217955-ITTI`

Within the gdb shell, `bt` command will display the backtrace for the segmentation fault. You should expect
an output like this:

```
203	/home/vagrant/magma/lte/gateway/c/oai/tasks/nas/nas_procedures.c: No such file or directory.
[Current thread is 1 (process 13887)]
(gdb) bt
#0  get_nas_specific_procedure_attach (ctxt=ctxt@entry=0x6210000d64b0) at /home/vagrant/magma/lte/gateway/c/oai/tasks/nas/nas_procedures.c:203
#1  0x000055b0080c2aa8 in emm_proc_attach_request (ue_id=ue_id@entry=2819, is_mm_ctx_new=is_mm_ctx_new@entry=true, ies=<optimized out>, ies@entry=0x608000054280)
    at /home/vagrant/magma/lte/gateway/c/oai/tasks/nas/emm/Attach.c:318
#2  0x000055b0080e5e1f in emm_recv_attach_request (ue_id=<optimized out>, originating_tai=originating_tai@entry=0x7f529e6c24d2, originating_ecgi=originating_ecgi@entry=0x7f529e6c2660,
    msg=msg@entry=0x7f529e6c1e10, is_initial=<optimized out>, is_mm_ctx_new=<optimized out>, emm_cause=<optimized out>, decode_status=<optimized out>)
    at /home/vagrant/magma/lte/gateway/c/oai/tasks/nas/emm/sap/emm_recv.c:384
#3  0x000055b0080e27b4 in _emm_as_establish_req (msg=msg@entry=0x7f529e6c25e0, emm_cause=emm_cause@entry=0x7f529e6c257c)
    at /home/vagrant/magma/lte/gateway/c/oai/tasks/nas/emm/sap/emm_as.c:802
#4  0x000055b0080e4460 in emm_as_send (msg=msg@entry=0x7f529e6c25d8) at /home/vagrant/magma/lte/gateway/c/oai/tasks/nas/emm/sap/emm_as.c:185
#5  0x000055b0080d1720 in emm_sap_send (msg=msg@entry=0x7f529e6c25d0) at /home/vagrant/magma/lte/gateway/c/oai/tasks/nas/emm/sap/emm_sap.c:105
#6  0x000055b0080b7037 in nas_proc_establish_ind (ue_id=ue_id@entry=2819, is_mm_ctx_new=<optimized out>, originating_tai=..., ecgi=..., as_cause=<optimized out>, s_tmsi=...,
    s_tmsi@entry=..., msg=0x629000302036) at /home/vagrant/magma/lte/gateway/c/oai/tasks/nas/nas_proc.c:185
#7  0x000055b007ebbf87 in mme_app_handle_initial_ue_message (mme_app_desc_p=mme_app_desc_p@entry=0x60800000b500, initial_pP=initial_pP@entry=0x629000302026)
    at /home/vagrant/magma/lte/gateway/c/oai/tasks/mme_app/mme_app_bearer.c:727
#8  0x000055b007ebab98 in handle_message (loop=<optimized out>, reader=<optimized out>, arg=<optimized out>) at /home/vagrant/magma/lte/gateway/c/oai/tasks/mme_app/mme_app_main.c:182
```

**6. Obtain event that triggered the crash**. Every time a service restarts it will generate a log file (i.e. mme.log). Inside the coredump folder you will find the log (i.e. mme.log) that was generated just before the crash. In order to understand what was the event that triggered the crash, get the last event (Attach Request, Detach, timer expiring, etc.) in the log file.

Note: If you can't find the timestamp of the crash in syslogs, you can use the last log generated in the log found in the coredump to get the exact timestamp of that crash.

**7. Put together all information and isolate different root causes**. If you have multiple crashes, make sure you perform these steps for each of these events. This will help to confirm the crashes are happening due a single or multiple issue. For each event, you should collect the following information

  - Detail of Impact
  - mme binary and the magma version deployed
  - Log in the syslog of the crash
  - Backtrace, reading the coredump
  - Get the event that triggered the crash in mme log

**8. Investigate or seek for help**. Use the collected information to look into previous Github issues/bugs and confirm if this is a known issue or bug that has been fixed in a later version. Otherwise, open a new report with the information collected.
