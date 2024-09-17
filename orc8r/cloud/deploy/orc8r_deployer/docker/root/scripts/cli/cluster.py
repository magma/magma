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
import click
from testlib.awscluster import AWSCluster
from testlib.cluster_factory import ClusterFactory, ClusterType


@click.group()
@click.version_option()
def cli():
    """Test cluster cli
    """
    pass


@cli.command()
@click.option('--template', default="", help='Location of the template file')
@click.option('--cluster-uuid', default="", help='UUID for the test cluster')
@click.option('--skip-certs', is_flag=True, default=False, help='skip certs')
@click.option(
    '--skip-precheck', is_flag=True,
    default=False, help='skip prechecks',
)
def create(
    template: str, cluster_uuid: str,
    skip_certs: bool, skip_precheck: bool,
):
    """Create a cluster

    Args:
        template ([type], optional): [Cluster template location]. Defaults to None.
    """
    ClusterFactory().create_cluster(
        ClusterType.AWS,
        template,
        cluster_uuid=cluster_uuid,
        skip_certs=skip_certs,
        skip_precheck=skip_precheck,
    )


@cli.command()
def cleanup():
    """Destroy a cluster"""
    AWSCluster().destroy()
