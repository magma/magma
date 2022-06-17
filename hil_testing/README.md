# Magma Spirent Testing

NOTE: The code is submited as-is. At the time of writing this README file, the SANITY test suite was fully functional. As long as the architecture  described in this readme is replicated, the expectation is that the test cases will continue to function. 

## Overview

Magma Test case are automated using spirent system.

The current code is sufficient to run the SANITY test suite using Spirent's Landslide platform. This testbed utilized Spirent's C100 M4 chassis. It is licensed with the following components:

Hardware:
  * 8 concurrent test sessions
  * 2x40G NICs
  * 6x1G NICs

Software: 
  * MME Nodal support
  * Sequence Mode support 
 
 For the 4G testcases, Spirent is used to emulate the UEs and the eNBs. MME, SGW, PGW, HSS are provided by Magma and are collectively referred to as the System Under Test (SUT).


## Setup

Clone the repo.

    python3 -m venv .venv
    ln -s .venv/bin/activate
    source activate
    pip3 install -r requirements
    cd spirent_automation

Note that the next time you login, you only need to: 

    source activate
    cd spirent_automation

### Base Spirent system configuration

* The TAS and the TS servers should be setup with up to 8 nics. If the # of licensed executors are not 8, modify `get_ports.py` accordingly.
* Modify `get_ports.py` with the appropriate IP addresses, VLANs, MAC addresses per test port. 
* Setup the uplink data ports as trunk ports with only the required vlans allowed on them. 
* Using the Spirent API, use the POST request and save templates saved in `/hil_testing/TC/spirent_templates`. These templates are required for the SANITY test suite to work.
* Create the DMFs (Data Message Flows) in Spirent Landslide corresponding to the DMFs referenced in TC files located in `/hil_testing/TC/*`.
* Create a list of SUTs in the Spirent Landslide GUI. These SUTs should also be cross refrenced in `config.py` (under Magma/AGW) as well as in the Ansible hosts file located in `hil_testing/Magma_Ansible/inventory/hosts.yaml`.

NOTE: The design for this HWIL testing has been to create a base configuration on the Landslide GUI and then to use the API to create clones of the template and modify critical parameters. 

NOTE: With this approach, every test that is run, will create a new file but will put a "DELETE_ME" tag on it. Automation can be run to delete these files periodically to keep the library clean. 


### Credentials ###

In order to run tests, you need some passwords and secret keys.

* TAS password.
* AWS creds
* AGW Password (Default is `magma`)
* RDS user/pass
* Slack token
You only need the passwords, as tool defaults the usernames in [config.py](config.py)
These passwords will need to be added to your environment (see below). If you need
to save the password in a file, be sure to set the permissions so that
only the owner has access, e.g.
```
chmod 0600 password_file
```

To run the SANITY test suite:

    export TAS_PASSWORD='****'
    export AWS_SECRET_KEY='secret_from_some_vault'
    export MAGMA_PASSWORD='****'
    export RDS_PASS='****'
    export RDS_HOST='****'
    export SLACK_WEBHOOK_PATH='****'
    
    ./main.py run SANITY

This could all be done in a single command line, with passwords from files, e.g.:
    TAS_PASSWORD=$(cat ~/.ssh/cred/.tas) AWS_SECRET_KEY=$(cat ~/.ssh/cred/.aws) MAGMA_PASSWORD=$(cat ~/.ssh/cred/.agw) RDS_PASS=$(cat ~/.ssh/cred/.rdspass) RDS_HOST=$(cat ~/.ssh/cred/.rdshost) SLACK_WEBHOOK_PATH=#(cat ~/.ssh/cred/.slack) ./main.py run SANITY

As an alternative to exporting passwords, you can create a file (for example, .creds.json)
and tell framework to use it:

    cat >.creds.json
    {
        "TAS_PASSWORD":"****",
        "AWS_SECRET_KEY":"secret_from_some_vault",
        "MAGMA_PASSWORD":"****",
        "RDS_PASS":"****",
        "RDS_HOST":"rds_end_point",
        "SLACK_WEBHOOK_PATH": "****"
    }
    
    ./main.py run --credentials-file=./.creds.json SANITY

## Alerting
As of now we alert on Slack private channel [HIL-test](https://magmacore.slack.com/archives/C02164DSGPM) for each test run.


## Usage
 see usage options:

    ./main.py --help
    usage: main.py [-h] {list,run} ...

    Run Hardware in Loop testing

    positional arguments:
        {list,run}  commands
          list      List supported test_suites
          run       Run test suite

    optional arguments:
      -h, --help  show this help message and exit

### listing all test per test suite ###
    ./main.py list sanity (or performance, feature)
    2021-06-02 05:15:53,019 MAGMA_AUTOMATION WARNING Logging set to WARNING
    SANITY                                        :

    TC008_SANITY_data_test
    TC010_SANITY_data_400UE_500k_400sec
    TC011_SANITY_data_600UE_500k_500sec
    TC002_SANITY_control_200UE_12enbs_10rate
    TC009_SANITY_data_200UE_2M_180sec
    TC004_SANITY_control_600UE_12enbs_10rate
    TC003_SANITY_control_400UE_12enbs_10rate
    TC001_SANITY_control_test

### Test Run usage: ###
    main.py run [-h] [--log-level {DEBUG,INFO,WARNING,ERROR,CRITICAL}] [--credentials-file CREDENTIALS_FILE] [--gateway {vagw1,mj_vgw,phy_agw5,phy-u4}] [--no-output-text] [--output-s3]
                   [--upgrade] [--skip-build-check] [--reboot]
                   {SANITY,PERFORMANCE,FEATURE} [only_run [only_run ...]]

    positional arguments:
    {SANITY,PERFORMANCE,FEATURE}
                        Run this group of tests
    only_run              Only run this test (default: None)

    optional arguments:
      -h, --help            show this help message and exit
      --log-level {DEBUG,INFO,WARNING,ERROR,CRITICAL}
                        Log at specified level and higher (default: WARNING)
      --credentials-file CREDENTIALS_FILE, -f CREDENTIALS_FILE
                        Full path to credentials file. JSON format (default: None)
      --gateway {vagw1,mj_vgw,phy_agw5,phy-u4}
                        select gateway (default: mj_vgw)
      --no-output-text      Whether or not to output ascii text tables (default: True)
      --output-s3           Whether or not to send results file to s3 (default: False)
      --upgrade             Whether or not to upgrade SUT (default: False)
      --skip-build-check    Whether or not to run test on same old SUT build (default: False)
      --reboot              Whether or not to reboot SUT before running test (default: False)

### SUT option ###

Framework supports running testsuit on specific AGW from available pool.
by default it will pick reserved for automation.

    ./main.py run SANITY --gateway vagw1

### SUT upgrade ###

Framework supports upgrading SUT to latest magma sw before running test on it. upgrade task executed only if new build available.
if perticular option is given.

    ./main.py run SANITY --gateway vagw1 --upgrade

### Pushing Test results on AWS ###

Framework supports pushing results to AWS s3 web portal.

    ./main.py run sanity --gateway vagw1 --upgrade --output-s3

### Run specific set of test cases for give test suite ###

    ./main.py run --credentials-file=<cred file path> --gateway <sut from pool> SANITY TC001 TC005 --log-level info --output-s3 --upgrade

### Run test with SUT reboot option ###

    ./main.py run --credentials-file=<cred file path> --gateway <sut from pool> SANITY --log-level info --output-s3 --upgrade --reboot

