"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import re
import subprocess  # noqa: S404
import sys

cmd = 'git diff --name-only `git merge-base origin/master HEAD`'
list_files = subprocess.check_output(cmd, shell=True, universal_newlines=True)  # noqa: S602
mme_is_impacted = False

for line in list_files.split('\n'):
    res = re.search('lte/gateway/Makefile', line)
    if res is not None:
        mme_is_impacted = True
    res = re.search('lte/gateway/c/core/oai', line)
    if res is not None:
        mme_is_impacted = True
    res = re.search('lte/gateway/c/sctpd', line)
    if res is not None:
        mme_is_impacted = True
    res = re.search('lte/gateway/docker/mme', line)
    if res is not None:
        mme_is_impacted = True
    res = re.search('JenkinsFile-OAI-Container-GitHub', line)
    if res is not None:
        mme_is_impacted = True
    res = re.search('generateHtmlReport-OAI-pipeline.py', line)
    if res is not None:
        mme_is_impacted = True
    res = re.search('check_pr_modified_files_for_oai_pipeline.py', line)
    if res is not None:
        mme_is_impacted = True
    res = re.search('ci-scripts/docker/Dockerfile.mme.ci.ubuntu18', line)
    if res is not None:
        mme_is_impacted = True
    res = re.search('ci-scripts/docker/Dockerfile.mme.ci.rhel8', line)
    if res is not None:
        mme_is_impacted = True
    res = re.search('feg/cloud/go/protos/s6a_proxy.pb.go', line)
    if res is not None:
        mme_is_impacted = True
    res = re.search('feg/protos/s6a_proxy.proto', line)
    if res is not None:
        mme_is_impacted = True

if mme_is_impacted:
    sys.exit(-1)
else:
    sys.exit(0)
