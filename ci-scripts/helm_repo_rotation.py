"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import os
import re
import time

import requests
from artifactory import ArtifactoryPath

# The number of versions kept in the helm repo
NUMBER_OF_CHARTS = 20
CHART_NAMES = ('cwf-orc8r', 'feg-orc8r', 'fbinternal-orc8r', 'lte-orc8r', 'orc8r')


def exit_environment_not_set(variable_name):
    """Exit due to environment variable not set"""
    print(variable_name + " is not set.")
    exit(1)


def delete_old_artifacts(list_of_artifacts, chart_name):
    """Extract list of artifacts for this chart and only keep the last NUMBER_OF_CHARTS versions uploaded"""
    artifacts = list(filter(lambda x: chart_name == x[0], list_of_artifacts))
    if len(artifacts) < NUMBER_OF_CHARTS:
        return
    artifacts.sort(key=lambda date: time.strptime(date[1], '%d-%b-%Y %H:%M'))
    for artifact in artifacts[:-(NUMBER_OF_CHARTS + 1)]:
        artifact_path = ArtifactoryPath(
            full_artifactory_url + artifact[2],
            auth=credentials,
        )
        if artifact_path.exists():
            print("Deleting artifact " + artifact[2])
            artifact_path.unlink()


# Verify existence of environment variables
username = os.getenv('HELM_CHART_MUSEUM_USERNAME')
if not username:
    exit_environment_not_set('HELM_CHART_MUSEUM_USERNAME')
password = os.getenv('HELM_CHART_MUSEUM_TOKEN')
if not password:
    exit_environment_not_set('HELM_CHART_MUSEUM_TOKEN')
artifactory_url = os.getenv('HELM_CHART_ARTIFACTORY_URL')
if not artifactory_url:
    exit_environment_not_set('HELM_CHART_ARTIFACTORY_URL')
helm_repo = os.getenv('HELM_CHART_MUSEUM_REPO')
if not helm_repo:
    exit_environment_not_set('HELM_CHART_MUSEUM_REPO')

# Login to artifactory
credentials = (username, password)

full_artifactory_url = artifactory_url + helm_repo + '/'
artifactory_path = ArtifactoryPath(
    full_artifactory_url,
    auth=credentials,
)
artifactory_path.touch()

# List from html website as GET API does not provide creation timestamp
r = requests.get(full_artifactory_url)

# Split response body per artifact
artifact_list = re.search("<a href(.*)</pre>", str(r.content)).group(0).split('<a href')
artifact_list_trimmed = list(filter(lambda x: x != "", artifact_list))

# Extract artifact name and date
clean_artifact_map = map(
    lambda p: re.search(">(.*)</a>[ ]+(.*)  ", str(p)).group(1, 2),
    artifact_list_trimmed,
)
# Extract chart name for sorting
clean_artifacts_and_charts = list(
    map(lambda x: ['-'.join(x[0].split('-')[:-1]), x[1], x[0]], clean_artifact_map),
)

# Make sure to only keep 20 artifacts per chart
for chart_name in CHART_NAMES:
    delete_old_artifacts(clean_artifacts_and_charts, chart_name)
