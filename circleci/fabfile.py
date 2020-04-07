#  Copyright (c) Facebook, Inc. and its affiliates.
#  All rights reserved.
#
#  This source code is licensed under the BSD-style license found in the
#  LICENSE file in the root directory of this source tree.
import re
from typing import List, Optional

import dateutil.parser
import requests
from fabric.api import run
from fabric.context_managers import cd, lcd
from fabric.contrib.files import exists
from fabric.exceptions import CommandTimeout
from fabric.operations import get, local, put
from fabric.state import env
from time import sleep

LEASE_MAX_RETRIES = 5

LTE_STACK = 'lte'
CWF_STACK = 'cwf'
STACKS = [
    LTE_STACK,
    CWF_STACK,
]

CWF_IMAGES = [
    'gateway_go',
    'gateway_python',
    'gateway_radius',
    'gateway_sessiond',
    'gateway_pipelined',
]

DEFAULT_DOCKER_REG = 'facebookconnectivity-magma-docker.jfrog.io'


class NodeLease:
    def __init__(self, node_id: str, lease_id: str, vpn_ip: str):
        self.node_id = node_id
        self.lease_id = lease_id
        self.vpn_ip = vpn_ip


def cwf():
    env.stack = CWF_STACK


def lte():
    env.stack = LTE_STACK


def integ_test(repo: str = 'git@github.com:facebookincubator/magma.git',
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
               docker_registry: str = DEFAULT_DOCKER_REG,
               docker_user: str = 'magmaci-bot',
               jfrog_key: str = 'jfrog_key',
               build_number: str = '0',
               workflow_names: str = 'lte-integ-test;cwf-integ-test',
               circle_key: str = 'circleci_key'):
    if env.stack not in STACKS:
        raise ValueError(f'Stack {env.stack} is not a valid stack.')
    should_test = run_integ_test == 'True'
    should_build = build_package == 'True'
    should_deploy = deploy_artifacts == 'True'

    workflow_list = list(filter(lambda x: x, workflow_names.split(';')))
    circle_key_val = local(f'cat {circle_key}', capture=True)
    lease = _acquire_lease_with_retry(repo, int(build_number), workflow_list,
                                      circle_key_val,
                                      api_url, cert_file, cert_key_file)
    if lease is None:
        print('Did not acquire a node lease in a reasonable time frame.')
        with open(f'{env.stack}_test_status', 'w+') as f:
            f.write('SKIP')
        raise Exception("Didn't get a node")
    _write_lease_to_disk(lease)

    try:
        _set_host_for_lease(lease, node_ssh_key)
        _checkout_code(repo, branch, sha1, tag, pr_num, magma_root)
        if should_test:
            if env.stack == LTE_STACK:
                _run_remote_lte_integ_test(repo, magma_root)
            elif env.stack == CWF_STACK:
                _run_remote_cwf_integ_test(repo, magma_root)

        if should_build and env.stack == LTE_STACK:
            # Destroy the VM if we didn't use it to run the integ tests in
            # this same job
            should_destroy_vm = not should_test
            _run_remote_lte_package(repo, magma_root,
                                    package_cert, package_control_proxy,
                                    should_destroy_vm)

        if should_deploy:
            if env.stack == LTE_STACK:
                _deploy_lte_packages(repo, magma_root)
            elif env.stack == CWF_STACK:
                append_oss_hash = 'facebookincubator/magma' not in repo
                # Rebuild the containers if we didn't run the integ tests in
                # this same job
                should_rebuild = not should_test
                _deploy_cwf_images(repo, magma_root,
                                   docker_registry, docker_user, jfrog_key,
                                   append_oss_hash, should_rebuild)
    except CommandTimeout as e:
        print('Remote run timed out, this probably indicates a problem '
              'with the job')
        raise e
    finally:
        if env.stack == LTE_STACK:
            _destroy_vms(repo, magma_root, 'lte/gateway',
                         ['magma', 'magma_test', 'magma_trfserver'])
        elif env.stack == CWF_STACK:
            _destroy_vms(repo, magma_root, 'cwf/gateway',
                         ['cwag', 'cwag_test'])
        _release_node_lease(api_url, lease.node_id, lease.lease_id,
                            cert_file, cert_key_file)


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


def _acquire_lease_with_retry(repo: str,
                              build_number: int,
                              workflow_names: List[str],
                              circle_key: str,
                              api_url: str,
                              cert_file: str,
                              cert_key_file: str) -> Optional[NodeLease]:
    lease_retries = 0
    while lease_retries < LEASE_MAX_RETRIES:
        lease = _acquire_node_lease(api_url, cert_file, cert_key_file)
        if lease is not None:
            print(f'Acquired lease for {lease.node_id}, '
                  f'lease ID {lease.lease_id}')
            return lease
        lease_retries += 1
        print('No nodes found, trying again after 5 minutes...')
        should_cancel = _do_newer_running_workflows_exist(
            repo,
            build_number, workflow_names, circle_key,
        )
        if should_cancel:
            # TODO: If we `circleci step halt` an integ test job here, how
            #  do we also skip the build/deploy job later in the workflow?
            print('Newer running workflows exist in the queue, I should '
                  'cancel myself when my developer implements this '
                  'functionality')
        sleep(300)
    return None


def _acquire_node_lease(api_url: str,
                        cert_file: str,
                        cert_key_file: str) -> Optional[NodeLease]:
    resource = f'{api_url}/magma/v1/ci/reserve'
    resp = requests.post(resource, cert=(cert_file, cert_key_file))
    if resp.status_code != 200:
        print(f'Received status code {resp.status_code} from node lease '
              f'request')
        return None
    resp_obj = resp.json()
    return NodeLease(resp_obj['id'], resp_obj['lease_id'], resp_obj['vpn_ip'])


def _checkout_code(repo: str, branch: str, sha1: str, tag: str, pr_num: str,
                   magma_root: str):
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
            _run_git(f'git fetch --force origin '
                     f'"{branch}:remotes/origin/{branch}"')
            _run_git(f'git reset --hard {sha1}')
            _run_git(f'git checkout -q -B {branch}')

        _run_git(f'git reset --hard {sha1}')


def _run_remote_lte_integ_test(repo: str, magma_root: str):
    repo_name = _get_repo_name(repo)
    with cd(f'{repo_name}/{magma_root}/lte/gateway'):
        run('fab integ_test', timeout=90*60)


def _run_remote_cwf_integ_test(repo: str, magma_root: str):
    repo_name = _get_repo_name(repo)
    with cd(f'{repo_name}/{magma_root}/cwf/gateway'):
        with cd('docker'):
            # Note: it seems like the --parallel flag makes the fabric output
            # unreadable, with all build outputs interleaved
            run('docker-compose '
                '-f docker-compose.yml '
                '-f docker-compose.override.yml '
                '-f docker-compose.nginx.yml '
                '-f docker-compose.integ-test.yml '
                'build --parallel')
        run('fab integ_test:destroy_vm=True,transfer_images=True',
            timeout=110*60)


def _run_remote_lte_package(repo: str, magma_root: str,
                            package_cert: str, package_control_proxy: str,
                            destroy_vm: bool):
    repo_name = _get_repo_name(repo)

    # We upload to the magma/fb directory on the CI node, but that maps to
    # /home/vagrant/magma/fb on the magma VM
    def secpath(user, file):
        return f'/home/{user}/{repo_name}/{magma_root}/fb/{file}'

    remote_secrets_dir = secpath('magma', '')
    cert_file = secpath('magma', 'rootCA.pem')
    control_proxy_file = secpath('magma', 'control_proxy.yml')

    # Upload rootCA, control proxy config
    run(f'mkdir -p {remote_secrets_dir}')
    put(package_cert, cert_file)
    put(package_control_proxy, control_proxy_file)

    with cd(f'{repo_name}/{magma_root}/lte/gateway'):
        fab_args = f'vcs=git,all_deps=False,' \
                   f'cert_file={secpath("vagrant", "rootCA.pem")},' \
                   f'proxy_config={secpath("vagrant", "control_proxy.yml")},' \
                   f'destroy_vm={destroy_vm}'
        run(f'fab test package:{fab_args}')
        # This will create /tmp/packages.tar.gz, /tmp/packages.txt on the
        # remote CI executor node (the current fab host)
        run('fab copy_packages')


def _deploy_lte_packages(repo: str, magma_root: str):
    repo_name = _get_repo_name(repo)

    # Grab all the build artifacts we need from the CI node
    get('/tmp/packages.tar.gz', 'packages.tar.gz')
    get('/tmp/packages.txt', 'packages.txt')
    get('/tmp/magma_version', 'magma_version')
    get(f'{repo_name}/{magma_root}/lte/gateway/release/magma.lockfile',
        'magma.lockfile')

    with open('magma_version') as f:
        magma_version = f.readlines()[0].strip()
    s3_path = f's3://magma-images/gateway/{magma_version}'
    local(f'aws s3 cp packages.txt {s3_path}.deplist')
    local(f'aws s3 cp magma.lockfile {s3_path}.lockfile')
    local(f'aws s3 cp packages.tar.gz {s3_path}.deps.tar.gz')


def _deploy_cwf_images(repo: str, magma_root: str,
                       docker_registry: str, user: str, jfrog_key: str,
                       append_oss_hash: bool, rebuild: bool):
    repo_name = _get_repo_name(repo)

    if rebuild:
        with cd(f'{repo_name}/{magma_root}/cwf/gateway/docker'):
            # Note: it seems like the --parallel flag makes the fabric output
            # unreadable, with all build outputs interleaved
            run('docker-compose '
                '-f docker-compose.yml '
                '-f docker-compose.override.yml '
                'build --parallel')

    local_hash = local('git rev-parse HEAD', capture=True)
    container_version = local_hash[:8]

    put(jfrog_key, '/tmp/jfrog_key')
    with cd(f'{repo_name}/{magma_root}/cwf/gateway/docker'):
        for img in CWF_IMAGES:
            run(f'../../../orc8r/tools/docker/publish.sh '
                f'-r {docker_registry} -i {img} '
                f'-u {user} -p /tmp/jfrog_key -v {container_version}')

    if append_oss_hash:
        oss_hash = _find_matching_opensource_commit(magma_root)
        container_version = f'{container_version}|{oss_hash}'
    local(f'echo "{container_version}" > cwag_version')


def _find_matching_opensource_commit(
        magma_root: str,
        oss_repo: str = 'https://github.com/facebookincubator/magma.git ',
) -> str:
    # Find corresponding hash in opensource repo by grabbing the message of the
    # latest commit to the magma root directory of the current repository then
    # searching for it in the open source repo
    commit_subj = local(f'git --no-pager log --oneline --pretty=format:"%s" '
                        f'-- {magma_root} | head -n 1', capture=True)
    local('rm -rf /tmp/ossmagma')
    local('mkdir -p /tmp/ossmagma')
    local(f'git clone {oss_repo} /tmp/ossmagma/magma')
    with lcd('/tmp/ossmagma/magma'):
        oss_hash = local(f'git --no-pager log --oneline --pretty=format:"%h" '
                         f'--grep=\'{commit_subj}\' | head -n 1', capture=True)
        return oss_hash


def _do_newer_running_workflows_exist(repo: str,
                                      build_number: int,
                                      workflow_names: List[str],
                                      circle_key: str) -> bool:
    repo_org, repo_name = _get_repo_org(repo), _get_repo_name(repo)
    circle_url = 'https://circleci.com/api/v1.1/project/gh'
    resp = requests.get(f'{circle_url}/{repo_org}/{repo_name}?'
                        f'circle-token={circle_key}&'
                        f'filter=running')
    if resp.status_code != 200:
        print(f'Got status code {resp.status_code} from CircleCI API, will '
              f'not perform auto-cancellation for this job.')
        return False
    builds = resp.json()
    match = list(filter(lambda b: b['build_num'] == build_number, builds))
    if not match:
        print('Did not find the current build in the list of recent builds, '
              'assuming that this is super out of date. Will report that '
              'this build should be canceled.')
        return True
    this_build = match[0]
    this_build_start = dateutil.parser.isoparse(this_build['start_time'])
    for b in builds:
        # Ignore ourselves
        if b['build_num'] == build_number:
            continue
        if b['workflows']['workflow_name'] not in workflow_names:
            continue

        other_start = dateutil.parser.isoparse(b['start_time'])
        if other_start > this_build_start:
            print(f'Found build f{b["build_num"]} with start time after the '
                  f'current build, will stop this run now.')
            return True
    return False


def _destroy_vms(repo: str, magma_root: str,
                 path: str, vms: List[str]) -> None:
    try:
        repo_name = _get_repo_name(repo)
        with cd(f'{repo_name}/{magma_root}/{path}'):
            run(f'vagrant destroy -f {" ".join(vms)}')
    except Exception as e:
        print(f'Caught exception from destroying VMs: {e}')


def _release_node_lease(api_url: str, node_id: str, lease_id: str,
                        cert_file: str, cert_key_file: str) -> None:
    resource = f'{api_url}/magma/v1/ci/nodes/{node_id}/release/{lease_id}'
    resp = requests.post(resource, cert=(cert_file, cert_key_file))
    if resp.status_code != 204:
        raise Exception(f'Got response code {resp.status_code} when releasing '
                        f'worker node, it may not have been released.')
    print(f'Released node {node_id}')


def _run_git(cmd: str, **kwargs):
    run(f'GIT_SSH_COMMAND="ssh -i ~/.ssh/gh_deploy.pem" {cmd}', **kwargs)


def _get_repo_name(repo: str) -> str:
    return re.match(r'.+/(.+)\.git', repo).group(1)


def _get_repo_org(repo: str) -> str:
    return re.match(r'.+:(.+)/.*\.git', repo).group(1)
