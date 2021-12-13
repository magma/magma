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
import sys
from time import sleep
from typing import List, Optional

import requests
from fabric.api import run
from fabric.context_managers import cd, settings
from fabric.contrib.files import exists
from fabric.exceptions import CommandTimeout
from fabric.operations import get, local, put
from fabric.state import env

LEASE_MAX_RETRIES = 5

LTE_STACK = 'lte'
CWF_STACK = 'cwf'
STACKS = [
    LTE_STACK,
    CWF_STACK,
]

CWF_IMAGES = [
    'cwag_go',
    'gateway_go',
    'gateway_python',
    'gateway_sessiond',
    'gateway_pipelined',
]

DEFAULT_DOCKER_REG = 'facebookconnectivity-magma-docker.jfrog.io'


class FabricException(Exception):
    pass


class NodeLease:
    def __init__(self, node_id: str, lease_id: str, vpn_ip: str):
        self.node_id = node_id
        self.lease_id = lease_id
        self.vpn_ip = vpn_ip


def cwf():
    env.stack = CWF_STACK


def lte():
    env.stack = LTE_STACK


def integ_test(
    repo: str = 'git@github.com:magma/magma.git',
    branch: str = '', sha1: str = '', tag: str = '',
    pr_num: str = '',
    magma_root: str = '',
    node_ssh_key: str = 'ci_node.pem',
    api_url: str = 'https://api-staging.magma.etagecom.io',
    cert_file: str = 'ci_operator.pem',
    cert_key_file: str = 'ci_operator.key.pem',
    run_integ_test: str = 'True',
    build_package: str = 'False',
    deploy_artifacts: str = 'False',
    package_cert: str = 'rootCA.pem',
    package_control_proxy: str = 'control_proxy.yml',
):
    if env.stack not in STACKS:
        raise ValueError(f'Stack {env.stack} is not a valid stack.')
    should_test = run_integ_test == 'True'
    should_build = build_package == 'True'
    should_deploy = deploy_artifacts == 'True'

    lease = _acquire_lease_with_retry(api_url, cert_file, cert_key_file)
    if lease is None:
        print('Did not acquire a node lease in a reasonable time frame.')
        with open(f'{env.stack}_test_status', 'w+') as f:
            f.write('SKIP')
        raise Exception("Didn't get a node")
    _write_lease_to_disk(lease)

    try:
        _set_host_for_lease(lease, node_ssh_key)
        _checkout_code(repo, branch, sha1, tag, pr_num, magma_root)
        # Try to destroy all vm with vagrant destroy and workarround locked virtualbox vm
        run(
            "vagrant global-status 2>/dev/null | "
            "awk '/virtualbox/{print $1}' | "
            "xargs -I {} vagrant destroy -f {}"
            "&> >(grep -oP '(?<=\"unregistervm\", \").*(?=\",)') | "
            "xargs -I {} vboxmanage startvm {} --type emergencystop",
        )
        # Destroy all running vagrant VMs. If we use the same node to run integ
        # tests on more than one repo, Vagrant will complain about colliding
        # VM names.
        # Pipe stderr to /dev/null to silence annoying vagrant update prompts
        run(
            "vagrant global-status 2>/dev/null | "
            "awk '/virtualbox/{print $1}' | "
            "xargs -I {} vagrant destroy -f {}",
        )

        if should_test:
            if env.stack == LTE_STACK:
                _run_remote_lte_integ_test(repo, magma_root)
            elif env.stack == CWF_STACK:
                _run_remote_cwf_integ_test(repo, magma_root)

        if should_build and env.stack == LTE_STACK:
            # Destroy the VM if we didn't use it to run the integ tests in
            # this same job
            should_destroy_vm = not should_test
            _run_remote_lte_package(
                repo, magma_root,
                package_cert, package_control_proxy,
                should_destroy_vm,
            )

        if should_deploy:
            if env.stack == LTE_STACK:
                _deploy_lte_packages(repo, magma_root)
            elif env.stack == CWF_STACK:
                raise Exception('CWAG artifacts should be built on Circle')
    except CommandTimeout as e:
        print(
            'Remote run timed out, this probably indicates a problem '
            'with the job',
        )
        raise e
    finally:
        _release_node_lease(
            api_url, lease.node_id, lease.lease_id,
            cert_file, cert_key_file,
        )


def _write_lease_to_disk(lease: NodeLease):
    """
    Write the lease to disk so circleCI job can release it if the fabric step
    fails or times out
    """
    with open('lease_node.out', 'w+') as f:
        f.write(lease.node_id)
    with open('lease_id.out', 'w+') as f:
        f.write(lease.lease_id)
    with open(f'{env.stack}_test_status', 'w+') as f:
        f.write('RUN')


def _set_host_for_lease(lease: NodeLease, node_ssh_key: str):
    env.host_string = f'magma@{lease.vpn_ip}'
    env.hosts = [env.host_string]
    env.key_filename = node_ssh_key
    env.disable_known_hosts = True


def _acquire_lease_with_retry(
    api_url: str,
    cert_file: str,
    cert_key_file: str,
) -> Optional[NodeLease]:
    lease_retries = 0
    while lease_retries < LEASE_MAX_RETRIES:
        lease = _acquire_node_lease(api_url, cert_file, cert_key_file)
        if lease is not None:
            print(
                f'Acquired lease for {lease.node_id}, '
                f'lease ID {lease.lease_id}',
            )
            return lease
        lease_retries += 1
        print('No nodes found, trying again after 5 minutes...')
        sleep(300)
    return None


def _acquire_node_lease(
    api_url: str,
    cert_file: str,
    cert_key_file: str,
) -> Optional[NodeLease]:
    resource = f'{api_url}/magma/v1/ci/reserve'
    resp = requests.post(resource, cert=(cert_file, cert_key_file))
    if resp.status_code != 200:
        print(
            f'Received status code {resp.status_code} from node lease '
            f'request',
        )
        return None
    resp_obj = resp.json()
    return NodeLease(resp_obj['id'], resp_obj['lease_id'], resp_obj['vpn_ip'])


def _checkout_code(
    repo: str, branch: str, sha1: str, tag: str, pr_num: str,
    magma_root: str,
):
    repo_name = _get_repo_name(repo)
    if not exists(repo_name):
        _run_git(f'git clone {repo}')
    else:
        with cd(repo_name):
            _run_git(f'git remote set-url origin "{repo}"', warn_only=True)

    # This logic comes from the CircleCI `checkout` step
    # TODO: allow PR builds - or does circle set branch env var to
    #  PR branch already?
    branch = branch or 'master'
    with cd(f'{repo_name}/{magma_root}'):
        _run_git('git clean -d -f')
        if tag:
            _run_git(f'git fetch --force origin "refs/tags/{tag}"')
            _run_git(f'git reset --hard {sha1}')
            _run_git(f'git checkout -q {tag}')
        else:
            _run_git(
                f'git fetch --force origin '
                f'"{branch}:remotes/origin/{branch}"',
            )
            _run_git(f'git reset --hard {sha1}')
            _run_git(f'git checkout -q -B {branch}')

        _run_git(f'git reset --hard {sha1}')


def _get_test_summaries_and_logs(test_result_code: int):
    # Copy from node
    local('mkdir -p test-results')
    with settings(warn_only=True):
        get('test-results', 'test-results')
    # Copy to the directory CircleCI expects
    local('sudo mkdir -p /tmp/test-results/')
    if len(os.listdir('test-results')):
        local('sudo mv test-results/* /tmp/test-results/')

    # On failure, transfer logs from all 3 VMs and copy to the log
    # directory. This will get stored as an artifact in the CircleCI
    # config.
    if test_result_code:
        tar_file_name = "lte-test-logs.tar.gz"
        # On failure, transfer logs into current directory
        log_path = './' + tar_file_name
        run(f'fab get_test_logs:dst_path="{log_path}"', warn_only=True)
        # Copy the log files out from the node
        local('mkdir lte-artifacts')
        if exists(log_path):
            get(tar_file_name, 'lte-artifacts')
        local('sudo mkdir -p /tmp/logs/')
        if os.listdir('lte-artifacts'):
            local('sudo mv lte-artifacts/* /tmp/logs/')


def _run_remote_lte_integ_test(repo: str, magma_root: str):
    repo_name = _get_repo_name(repo)
    with cd(f'{repo_name}/{magma_root}/lte/gateway'):
        test_result = run('fab integ_test', timeout=180 * 60, warn_only=True)

        # Transfer test summaries into current directory
        run('fab get_test_summaries:dst_path="test-results"', warn_only=True)

        _get_test_summaries_and_logs(test_result.return_code)

        # Exit with the original test result
        sys.exit(test_result.return_code)


def _run_remote_cwf_integ_test(repo: str, magma_root: str):
    repo_name = _get_repo_name(repo)
    with cd(f'{repo_name}/{magma_root}/cwf/gateway'):
        with cd('docker'):
            # Note: it seems like the --parallel flag makes the fabric output
            # unreadable, with all build outputs interleaved
            run(
                'docker-compose '
                '-f docker-compose.yml '
                '-f docker-compose.override.yml '
                '-f docker-compose.nginx.yml '
                '-f docker-compose.integ-test.yml '
                'build --parallel',
            )
        test_xml = "tests.xml"
        with settings(abort_exception=FabricException):
            try:
                result = run(
                    'fab integ_test:'
                    'destroy_vm=True,'
                    'transfer_images=True,'
                    'test_result_xml=' + test_xml,
                    timeout=110 * 60, warn_only=True,
                )
            except Exception as e:
                _transfer_all_artifacts()
                print(f'Exception while running cwf integ_test\n {e}')
                sys.exit(1)
        # Move JUnit test result to /tmp/test-results directory
        if exists(test_xml):
            local('mkdir cwf-tests-xml')
            get(test_xml, 'cwf-tests-xml')
            local('sudo mkdir -p /tmp/test-results/')
            local('sudo mv cwf-tests-xml/* /tmp/test-results/')
        # On failure, transfer logs of key services from docker containers and
        # copy to the log directory. This will get stored as an artifact in the
        # circleCI config.
        if result.return_code:
            _transfer_all_artifacts()
        sys.exit(result.return_code)


def _transfer_all_artifacts():
    services = "sessiond session_proxy pcrf ocs pipelined ingress"
    run(
        f'fab transfer_artifacts:services="{services}",'
        'get_core_dump=True',
    )
    # Copy log files out from the node
    local('mkdir cwf-artifacts')
    get('*.log', 'cwf-artifacts')
    if exists("coredump.tar.gz"):
        get('coredump.tar.gz', 'cwf-artifacts')
    local('sudo mkdir -p /tmp/logs/')
    local('sudo mv cwf-artifacts/* /tmp/logs/')


def _run_remote_lte_package(
    repo: str, magma_root: str,
    package_cert: str, package_control_proxy: str,
    destroy_vm: bool,
):
    repo_name = _get_repo_name(repo)

    remote_secrets_dir = f'/home/magma/{repo_name}/{magma_root}/fb'
    cert_file = f'{remote_secrets_dir}/rootCA.pem'
    control_proxy_file = f'{remote_secrets_dir}/control_proxy.yml'

    # Upload rootCA, control proxy config to CI node
    run(f'mkdir -p {remote_secrets_dir}')
    put(package_cert, cert_file)
    put(package_control_proxy, control_proxy_file)

    # These map to a different directory on the vagrant VM
    vagrant_cert = '/home/vagrant/magma/fb/rootCA.pem'
    vagrant_cp = '/home/vagrant/magma/fb/control_proxy.yml'

    with cd(f'{repo_name}/{magma_root}/lte/gateway'):
        fab_args = f'vcs=git,all_deps=False,' \
                   f'cert_file={vagrant_cert},' \
                   f'proxy_config={vagrant_cp},' \
                   f'destroy_vm={destroy_vm}'
        run(f'fab release package:{fab_args}')
        # This will create /tmp/packages.tar.gz, /tmp/packages.txt on the
        # remote CI executor node (the current fab host)
        run('fab copy_packages')


def _deploy_lte_packages(repo: str, magma_root: str):
    repo_name = _get_repo_name(repo)

    # Grab all the build artifacts we need from the CI node
    get('/tmp/packages.tar.gz', 'packages.tar.gz')
    get('/tmp/packages.txt', 'packages.txt')
    get('/tmp/magma_version', 'magma_version')
    get(
        f'{repo_name}/{magma_root}/lte/gateway/release/magma.lockfile.debian',
        'magma.lockfile.debian',
    )

    with open('magma_version') as f:
        magma_version = f.readlines()[0].strip()
    s3_path = f's3://magma-images/gateway/{magma_version}'
    local(
        f'aws s3 cp packages.txt {s3_path}.deplist '
        f'--acl bucket-owner-full-control',
    )
    local(
        f'aws s3 cp magma.lockfile.debian {s3_path}.lockfile.debian '
        f'--acl bucket-owner-full-control',
    )
    local(
        f'aws s3 cp packages.tar.gz {s3_path}.deps.tar.gz '
        f'--acl bucket-owner-full-control',
    )


def _destroy_vms(
    repo: str, magma_root: str,
    path: str, vms: List[str],
) -> None:
    repo_name = _get_repo_name(repo)
    with cd(f'{repo_name}/{magma_root}/{path}'), settings(warn_only=True):
        run(f'vagrant destroy -f {" ".join(vms)}')


def _release_node_lease(
    api_url: str, node_id: str, lease_id: str,
    cert_file: str, cert_key_file: str,
) -> None:
    resource = f'{api_url}/magma/v1/ci/nodes/{node_id}/release/{lease_id}'
    resp = requests.post(resource, cert=(cert_file, cert_key_file))
    if resp.status_code != 204:
        raise Exception(
            f'Got response code {resp.status_code} when releasing '
            f'worker node, it may not have been released.',
        )
    print(f'Released node {node_id}')


def _run_git(cmd: str, **kwargs):
    run(f'GIT_SSH_COMMAND="ssh -i ~/.ssh/gh_deploy.pem" {cmd}', **kwargs)


def _get_repo_name(repo: str) -> str:
    return re.match(r'.+/(.+)\.git', repo).group(1)


def _get_repo_org(repo: str) -> str:
    return re.match(r'.+:(.+)/.*\.git', repo).group(1)
