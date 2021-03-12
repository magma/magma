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

import os
import re
import subprocess
import sys

cmd = 'git diff --name-only `git merge-base origin/master HEAD`'
listModifiedFiles = subprocess.check_output(cmd, shell=True, universal_newlines=True)
mmeIsImpacted = False

for line in listModifiedFiles.split('\n'):
    res = re.search('lte/gateway/Makefile', line)
    if res is not None:
        mmeIsImpacted = True
    res = re.search('lte/gateway/c/oai', line)
    if res is not None:
        mmeIsImpacted = True
    res = re.search('lte/gateway/c/sctpd', line)
    if res is not None:
        mmeIsImpacted = True
    res = re.search('lte/gateway/docker/mme', line)
    if res is not None:
        mmeIsImpacted = True
    res = re.search('ci-scripts/JenkinsFile-OAI-Container-GitHub|ci-scripts/generateHtmlReport-OAI-pipeline.py|ci-scripts/check_pr_modified_files_for_oai_pipeline.py', line)
    if res is not None:
        mmeIsImpacted = True
    res = re.search('ci-scripts/docker/Dockerfile.mme.ci.ubuntu18', line)
    if res is not None:
        mmeIsImpacted = True

if mmeIsImpacted:
    sys.exit(-1)
else:
    sys.exit(0)
