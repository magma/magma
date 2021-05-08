## vars.yml

vars.yml is used by orcl tool to get the user input for orc8r deployment. This file should contain all the input variables in orc8r and orc8r-app module which require user input. The input variables in orc8r and orc8r-app are specified in variables.tf file in their respective modules.

### Note:

* Some input variables in orc8r-app derives its value from orc8r module output, these input variables needn't be included in this file.
* We needn't copy all default values from variables.tf in here. We can only copy those values which might be useful to have pre checks against. For e.g. we set the defaults for chart versions, so that we can verify if helm repo contains the associated chart versions

### Check if the file is in sync

check_vars_sync.py can be used to check if this file is in sync with variables in terraform modules and add any missing variables in here.


