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
import logging
import time
import uuid
from pathlib import Path
from typing import Any, Dict

from cli.certs import certs_cmd
from cli.cleanup import raw_cleanup, tf_destroy
from cli.configlib import ConfigManager
from cli.install import precheck_cmd, tf_install
from testlib.cluster_definitions import (
    ClusterConfig,
    ClusterCreateError,
    ClusterDestroyError,
    ClusterInternalConfig,
    ClusterTemplate,
    ClusterType,
    GatewayConfig,
)
from utils.ansiblelib import AnsiblePlay, run_playbook
from utils.awslib import (
    check_elastic_role_not_exists,
    get_aws_configs,
    get_bastion_ip,
    get_gateways,
    verify_resources_exist,
)
from utils.common import get_json, init


class AWSCluster(object):
    def __init__(self):
        """Construct a new AWS cluster instance"""
        self.constants = init()

    def upgrade(self):
        raise NotImplementedError()

    def destroy_gateways(self, cluster_config):
        """ destroy the AWS gateways instantiated through test cluster
        Args:
            cluster_config (dict): Test cluster configuration
        """
        if not cluster_config.gateways:
            print("Gateway information not found in cluster configs")
            return

        project_dir = self.constants["project_dir"]
        playbook_dir = self.constants["cloudstrapper_playbooks"]

        # cloudstrapper expects secrets.yaml to be present
        Path(f"{project_dir}/secrets.yaml").touch()

        aws_configs = get_aws_configs()
        cstrap_dict = {
            "testClusterStacks": [gw.gateway_id for gw in cluster_config.gateways],
            "dirLocalInventory": self.constants["project_dir"],
            "idSite": "TestCluster",
            "awsAccessKey": aws_configs["aws_access_key_id"],
            "awsSecretKey": aws_configs["aws_secret_access_key"],
            "dirSecretsLocal": self.constants["secret_dir"],
        }

        cluster_cleanup = AnsiblePlay(
            playbook=f"{playbook_dir}/cluster-provision.yaml",
            tags=['clusterCleanup'],
            extra_vars=cstrap_dict,
        )
        network_cleanup = AnsiblePlay(
            playbook=f"{playbook_dir}/agw-provision.yaml",
            tags=['cleanupBridge', 'cleanupNet'],
            skip_tags=['attachIface'],
            extra_vars=cstrap_dict,
        )

        for playbook in [cluster_cleanup, network_cleanup]:
            print(f"Running playbook {playbook}")
            rc = run_playbook(playbook)
            if rc != 0:
                raise ClusterDestroyError(f"Failed destroying cluster")

    def destroy_orc8r(self):
        rc = tf_destroy(self.constants, warn=False)
        if rc != 0:
            raw_cleanup(self.constants)

    def destroy(self):
        cluster_config_fn = self.constants["test_cluster_config"]
        if not Path(cluster_config_fn).exists():
            print("Cluster config doesn't exist")
            return

        cluster_config = ClusterConfig.from_dict(get_json(cluster_config_fn))

        self.destroy_orc8r()
        self.destroy_gateways(cluster_config)

        # check resources specific to this magma UUID
        verify_resources_exist(cluster_config.uuid)

        # remove cluster config file
        Path(self.constants["test_cluster_config"]).unlink(missing_ok=True)


class AWSClusterFactory():
    @staticmethod
    def generate_cluster_stack(prefix, count):
        return [prefix + str(i) for i in range(count)]

    def create_gateways(
            self,
            constants: dict,
            cluster_uuid: str,
            template: ClusterTemplate = None,
    ) -> Dict[str, Any]:
        """ Create AGW gateways in the test cluster

        Args:
            constants (dict): Constants dictionary
            template (ClusterTemplate, optional): Cluster template definition. Defaults to None.
            skip_certs (bool, optional): Skip certs creation. Defaults to False.
            skip_precheck (bool, optional): Skip prechecks. Defaults to False.

        Raises:
            ClusterCreateError: Exception raised when cluster creation fails

        Returns:
            ClusterConfig: Returns the cluster configuration
        """
        # create edge network
        project_dir = constants["project_dir"]
        playbook_dir = constants["cloudstrapper_playbooks"]

        # cloudstrapper expects secrets.yaml to be present
        Path(f"{project_dir}/secrets.yaml").touch()

        aws_configs = get_aws_configs()
        cluster_stack = AWSClusterFactory.generate_cluster_stack(
            template.gateway.prefix,
            template.gateway.count,
        )

        cstrap_dict = {
            "clusterUuid": cluster_uuid,
            "dirLocalInventory": constants["project_dir"],
            "idSite": "TestCluster",
            "testClusterStacks": cluster_stack,
            "awsAccessKey": aws_configs["aws_access_key_id"],
            "awsSecretKey": aws_configs["aws_secret_access_key"],
            "idGw": "dummy_gateway",
            "dirSecretsLocal": constants["secret_dir"],
            "awsAgwAmi": template.gateway.ami,
            "awsCloudstrapperAmi": template.gateway.cloudstrapper_ami,
            "awsAgwRegion": template.gateway.region,
            "awsAgwAz": template.gateway.az,
            "orc8rDomainName": template.orc8r.infra['orc8r_domain_name'],
        }

        key_create = AnsiblePlay(
            playbook=f"{playbook_dir}/aws-prerequisites.yaml",
            tags=['keyCreate'],
            extra_vars=cstrap_dict,
        )
        bridge_gw_create = AnsiblePlay(
            playbook=f"{playbook_dir}/agw-provision.yaml",
            tags=['createNet', 'createBridge', 'inventory'],
            skip_tags=['attachIface'],
            extra_vars=cstrap_dict,
        )

        # create test instances
        test_inst_create = AnsiblePlay(
            playbook=f"{playbook_dir}/cluster-provision.yaml",
            tags=['clusterStart'],
            extra_vars=cstrap_dict,
        )

        jump_config_dict = {"agws": f"tag_Name_TestClusterBridge"}
        jump_config_dict.update(cstrap_dict)
        test_ssh_configure = AnsiblePlay(
            playbook=f"{playbook_dir}/cluster-provision.yaml",
            tags=['clusterJump'],
            extra_vars=jump_config_dict,
        )

        # configure test instances
        agws_config_dict = {
            "agws": f"tag_Name_{template.gateway.prefix}*",
        }
        agws_config_dict.update(template.gateway.service_config)
        agws_config_dict.update(cstrap_dict)
        test_inst_configure = AnsiblePlay(
            inventory=f"{project_dir}/common_instance_aws_ec2.yaml",
            playbook=f"{playbook_dir}/cluster-configure.yaml",
            tags=['exporter', 'clusterConfigure'],
            extra_vars=agws_config_dict,
        )

        max_retries = 3
        for i in range(max_retries):
            fail = False
            for playbook in [
                    key_create,
                    bridge_gw_create,
                    test_inst_create,
                    test_ssh_configure,
                    test_inst_configure,
            ]:
                print(f"Running playbook {playbook}")
                rc = run_playbook(playbook)
                if rc != 0:
                    fail = True
                    print("Failed creating gateway cluster...trying again")
                    break
                # sleep 10 seconds
                time.sleep(10)

            if not fail:
                break

        # get the newly instantiated gateways
        gateways = []
        for gw_info in get_gateways(template.gateway.prefix):
            (gateway_id, hostname) = gw_info
            gateways.append(
                GatewayConfig(
                    gateway_id=gateway_id,
                    hostname=hostname,
                    hardware_id="",
                ),
            )
        internal_config = ClusterInternalConfig(
            bastion_ip=get_bastion_ip(),
        )
        cluster_config_dict = {
            "uuid": cluster_uuid,
            "internal_config": internal_config,
            "cluster_type": ClusterType.AWS,
            "template": template,
            "gateways": gateways,
        }
        return cluster_config_dict

    def create_orc8r(
            self,
            constants: dict,
            cluster_uuid: str,
            template: ClusterTemplate = None,
            skip_certs=False,
            skip_precheck=False,
    ) -> Dict[str, Any]:
        """ Create an orc8r instance in the test cluster

        Args:
            constants (dict): Constants dictionary
            template (ClusterTemplate, optional): Cluster template definition. Defaults to None.
            skip_certs (bool, optional): Skip certs creation. Defaults to False.
            skip_precheck (bool, optional): Skip prechecks. Defaults to False.

        Raises:
            ClusterCreateError: Exception raised when cluster creation fails
        """
        # configure deployment
        template.orc8r.infra['magma_uuid'] = cluster_uuid
        template.orc8r.infra.update(get_aws_configs())

        # set elastic deploy role based on current state
        k = 'deploy_elasticsearch_service_linked_role'
        template.orc8r.platform[k] = check_elastic_role_not_exists()

        mgr = ConfigManager(constants)
        template_dict = template.orc8r.to_dict()
        for component, configs in template_dict.items():
            for k, v in configs.items():
                mgr.set(component, k, v)
            mgr.commit(component)

        # run playbooks in order
        if not skip_certs:
            logging.debug("Adding self signed and application certs")
            rc = run_playbook(certs_cmd(constants, self_signed=True))
            if rc != 0:
                raise ClusterCreateError(f"Failed running adding certs")

        if not skip_precheck:
            logging.debug("Running installation prechecks")
            rc = run_playbook(precheck_cmd(constants))
            if rc != 0:
                raise ClusterCreateError(f"Failed running prechecks")

        # create the orc8r cluster
        rc = tf_install(constants, warn=False)
        if rc != 0:
            raise ClusterCreateError(f"Failed installing cluster")

        # update dns record for parent domain
        dns_dict = {
            "domain_name": template.orc8r.infra["orc8r_domain_name"],
        }
        dns_dict.update(constants)
        rc = run_playbook(
            AnsiblePlay(
            playbook=f"{constants['playbooks']}/main.yml",
            tags=['update_dns_records'],
            extra_vars=dns_dict,
            ),
        )
        if rc != 0:
            raise ClusterCreateError(
                f"Failed updating dns records for parent domain",
            )

        cluster_config_dict = {
            "uuid": cluster_uuid,
            "cluster_type": ClusterType.AWS,
            "template": template,
        }
        return cluster_config_dict

    def create_cluster(
            self,
            template_fn: str = "",
            cluster_uuid: str = "",
            skip_certs=False,
            skip_precheck=False,
    ):
        """Create AWS based cluster based on provided template

        Args:
            template (ClusterTemplate, optional): Cluster template definition. Defaults to None.
            skip_certs (bool, optional): Skip certs creation. Defaults to False.
            skip_precheck (bool, optional): Skip prechecks. Defaults to False.

        Raises:
            ClusterCreateError: Exception raised when cluster creation fails
        """
        constants = init()
        test_cluster_config_fn = constants["test_cluster_config"]
        if Path(test_cluster_config_fn).exists():
            print("Cluster already exists")
            return

        if not template_fn:
            template_fn = constants['default_template']

        template = ClusterTemplate.from_dict(get_json(template_fn))

        if not cluster_uuid:
            # assign a uuid for this deployment
            cluster_uuid = str(uuid.uuid4())

        cluster_config_dict = {
            "uuid": cluster_uuid,
            "cluster_type": ClusterType.AWS,
            "template": template,
            "internal_config": ClusterInternalConfig(bastion_ip=""),
            "gateways": [],
        }

        try:
            ret = self.create_orc8r(
                constants,
                cluster_uuid,
                template,
                skip_certs,
                skip_precheck,
            )
            cluster_config_dict.update(ret)
            ret = self.create_gateways(constants, cluster_uuid, template)
            cluster_config_dict.update(ret)
        except ClusterCreateError:
            raise
        finally:
            cluster_config = ClusterConfig.from_dict(cluster_config_dict)
            with open(test_cluster_config_fn, "w") as f:
                f.write(cluster_config.to_json())
