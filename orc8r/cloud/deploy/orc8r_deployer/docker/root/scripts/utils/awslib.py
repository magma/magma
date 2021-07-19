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
import subprocess

import boto3
from cli.style import print_error_msg
from utils.common import run_command


def set_aws_configs(params: set, configs: dict):
    """Sets AWS configuration

    Args:
        params (set): List of aws configuration attributes
        configs (dict): Configuration map for a particular component
    """
    for k, v in configs.items():
        if k not in params:
            continue
        cmd = ["aws", "configure", "set", k, v]
        proc_inst = run_command(cmd)
        if proc_inst.returncode != 0:
            print_error_msg(f"Failed configuring aws with {k}")


def get_aws_configs():
    """Gets AWS configuration from environment"""
    env_params_cfg_map = (
        ('AWS_ACCESS_KEY_ID', 'aws_access_key_id'),
        ('AWS_SECRET_ACCESS_KEY', 'aws_secret_access_key'),
        ('AWS_DEFAULT_REGION', 'region'),
    )
    configs = {}
    for env_param, cfg_key in env_params_cfg_map:
        cmd = ["aws", "configure", "get", cfg_key]
        proc_inst = run_command(cmd)
        val = None
        if proc_inst.returncode == 0:
            val = proc_inst.stdout.strip()

        if not val:
            val = os.environ.get(env_param)
        configs[cfg_key] = val

    return configs


def check_elastic_role_not_exists():
    elastic_role = 'AWSServiceRoleForAmazonElasticsearchService'
    client = boto3.client('iam')
    try:
        client.get_role(RoleName=elastic_role)
    except client.exceptions.NoSuchEntityException:
        return True
    return False


def get_gateways(gateway_prefix: str = "agw"):
    client = boto3.client('ec2')
    gateways = []
    try:
        instance_info = client.describe_instances(
            Filters=[
                {
                    'Name': 'tag:Name',
                    'Values': [f'{gateway_prefix}*'],
                },
                {
                    'Name': 'instance-state-name',
                    'Values': ['running'],
                },
            ],
        )
        for reservation in instance_info["Reservations"]:
            gateway_id = ""
            hostname = ""
            for instance in reservation["Instances"]:
                for tags in instance["Tags"]:
                    if tags["Key"] == "Name":
                        gateway_id = tags["Value"]
                        break
                gateway_ip = instance['PrivateIpAddress']
            if gateway_id and gateway_ip:
                gateways.append((gateway_id, gateway_ip))
    except client.exceptions.NoSuchEntityException:
        pass
    return gateways


def get_bastion_ip(gateway_prefix: str = "agw"):
    client = boto3.client('ec2')
    gateways = []
    try:
        instance_info = client.describe_instances(
            Filters=[
                {
                    'Name': 'tag:Name',
                    'Values': ['*Bridge'],
                },
                {
                    'Name': 'instance-state-name',
                    'Values': ['running'],
                },
            ],
        )
        for reservation in instance_info["Reservations"]:
            for instance in reservation["Instances"]:
                return instance['PublicIpAddress']
    except client.exceptions.NoSuchEntityException:
        pass
    return ''


def verify_resources_exist(uuid: str = None):
    """Check if resources with a specific uuid exist

    Args:
        uuid (str, optional): unique id to identify cluster. Defaults to None.
    """
    client = boto3.client('resourcegroupstaggingapi')
    tagFilters = []
    if uuid:
        tag_filters = [{
            'Key': 'magma-uuid',
            'Values': [uuid],
        }]
    resources = client.get_resources(
        TagFilters=tag_filters,
    )
    for resource in resources['ResourceTagMappingList']:
        print(resource["ResourceARN"])
