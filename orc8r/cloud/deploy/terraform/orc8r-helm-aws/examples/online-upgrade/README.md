Online v1.1.0 Orchestrator Upgrade
===

IMPORTANT: Read the "Upgrading from 1.0" Orchestrator documentation for details
on how to use this module to upgrade your deployment. If you just
`terraform apply` here things will probably go very poorly.

The files in this directory can be copied and used as-is to perform an online
upgrade of the Orchestrator application from v1.0.x to v1.1.x. Please read the
descriptions of all the variables very carefully, and double-check the output
of `terraform apply` before applying changes. A misconfiguration may result in
significant downtime.

A lot of variables in this module have been set with defaults equal to the
legacy Terraform root module's defaults. If you copied that as-is to deploy
your v1.0.x application, you can leave these default values alone. If something
does not match up with your old Terraform configuration, you will see an
unexpected planned step in the output of `terraform apply` before it asks
for confirmation.

Detailed upgrade steps can be found in the Orchestrator documentation at
https://facebookincubator.github.io/magma.
