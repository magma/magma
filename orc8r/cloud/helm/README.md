Orchestrator Helm Charts
===

If you're making changes to the chart under orc8r/, you MUST:


1. Bump the versions of any updated subcharts in their Chart.yaml files

2. Update the requirements.yaml and requirements.lock files in the main chart
to match.

3. Bump the version of the main chart in Chart.yaml

4. `helm package orc8r` from this directory and deploy the .tgz file to the
chart repository. You can do this with the JFrog web UI or using the JFrog
CLI (if you've configured it): `jfrog rt upload orc8r-X.tgz orc8r-charts`
