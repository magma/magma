
# TeraVM script for CI
TeraVM is a test environment that is being used for checking Magma regression.
This script allows CI to update Magma gateways on different environments 
and run TeraVM (NG40) tests

## Configuration
### Fill in config file
Please gather all the information required on `fabfile_teravm_setup.json` file. 

On `general` you will find a set of variables that define APT repo and 
location of some docker components as well as user names to access VMs and the 
default list of tests that should be run.

On the `setups` section you may define IP configuration of as many setups as you 
need as long as they have the same `general` configuration.

### Password-les ssh
In order for this script to work, fab needs to be able to access ssh without using 
passowrd. So ssh password-less (private/public key based) is required.

For that reason you will need to generate a private/public ssh key on 
the machine where this script is running (if you haven't yet). 
```
ssh-keygen -t rsa -b 4096 -C "your_email@domain.com"
```

Then append the public key (.pub) to AGW, FEG and NG40 VMs on ` ~/.ssh/authorized_keys` 
file 

Test from the VM running this script using
```
ssh remote_username@server_ip_address
```

## Update 
### Update AGW
To update AGW run the following command 
```
fab upgrade_teravm_agw:setup_1,9cdb9470
```
Where `setup_1` is the setups key you use on config json file 
Where `9cdb9470` is the github hash. 

### Update FEG
To update FEG run the following command 
```
fab upgrade_teravm_feg:setup_1,9cdb9470
```
Where `setup_1` is the setups key you use on config json file 
Where `9cdb9470` is the github hash. 

## Run Test
To run the test on NG40 you can use 
```
fab run_3gpp_tests:setup_1
```
Where `setup_1` is the setups key you use on config json file 

## Update and Run Test
You can also do all 3 steps at the same time using 
```
fab upgrade_and_run_3gpp_tests:setup_1,9cdb9470
```
Where `setup_1` is the setups key you use on config json file 
Where `9cdb9470` is the github hash. 
