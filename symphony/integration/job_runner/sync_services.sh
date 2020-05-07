#!/usr/bin/env bash

for tenant in $(/grpcurl -plaintext graph:443 graph.TenantService.List | jq -r '.tenants[] | .name')
do
    curl http://graph/jobs/sync_services -H "x-auth-organization: $tenant" -H "x-auth-automation-name: job_runner" -H "x-auth-user-role: OWNER"
done
