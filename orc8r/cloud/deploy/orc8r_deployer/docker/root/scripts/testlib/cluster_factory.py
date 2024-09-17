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


from testlib.awscluster import AWSClusterFactory
from testlib.cluster_definitions import ClusterTemplate, ClusterType


class ClusterFactory():
    def create_cluster(
            self, typ, template: str = "", **kwargs
    ):
        if typ == ClusterType.AWS:
            return AWSClusterFactory().create_cluster(template, **kwargs)
        elif typ == ClusterType.LOCAL:
            return LocalClusterFactory().create_cluster(template, **kwargs)


class LocalClusterFactory():
    def create_cluster(self, template: ClusterTemplate = None):
        pass
