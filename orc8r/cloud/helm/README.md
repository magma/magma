Orchestrator Helm Charts
===

If you're making changes to the chart under orc8r/, you MUST:


1. Bump the versions of any updated subcharts in their Chart.yaml files

2. Update the requirements.yaml and requirements.lock files in the main chart
to match.

3. Bump the version of the main chart in Chart.yaml

4. Run `helm dependency update` to pull in any `orc8rlib` changes

5. `helm package orc8r` from this directory and deploy the .tgz file to the
chart repository. You can do this with the JFrog web UI or using the JFrog
CLI (if you've configured it): `jfrog rt upload orc8r-X.tgz orc8r-charts`

If you're making changes to the chart under orc8rlib/, you need to run
`helm dependency update` from the directory of the helm chart (e.g.
`/magma/orc8r/cloud/helm/orc8r`) that uses this library chart in order to
pull in the changes.