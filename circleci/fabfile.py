#  Copyright (c) Facebook, Inc. and its affiliates.
#  All rights reserved.
#
#  This source code is licensed under the BSD-style license found in the
#  LICENSE file in the root directory of this source tree.
import re
from time import sleep
from typing import Optional

import requests
from fabric.api import run
from fabric.context_managers import cd
from fabric.contrib.files import exists
from fabric.state import env

LEASE_MAX_RETRIES = 5


class NodeLease:
    def __init__(self, node_id: str, lease_id: str, vpn_ip: str):
        self.node_id = node_id
        self.lease_id = lease_id
        self.vpn_ip = vpn_ip


def integ_test(repo: str = 'git@github.com:facebookincubator/magma.git',
               branch: str = '', sha1: str = '', tag: str = '',
               pr_num: str = '',
               magma_root: str = '',
               node_ssh_key: str = 'ci_node.pem',
               api_url: str = 'https://api-staging.magma.etagecom.io',
               cert_file: str = 'ci_operator.pem',
               cert_key_file: str = 'ci_operator.key.pem',
               ):
    lease = None
    lease_retries = 0
    while lease is None and lease_retries < LEASE_MAX_RETRIES:
        lease = _acquire_node_lease(api_url, cert_file, cert_key_file)
        if lease is not None:
            print(f'Acquired lease for {lease.node_id}, '
                  f'lease ID {lease.lease_id}')
            break
        lease_retries += 1
        print('No nodes found, trying again after 5 minutes...')
        sleep(300)

    if lease is None:
        print('Did not acquire a node lease in a reasonable time frame.')
        return

    env.host_string = f'magma@{lease.vpn_ip}'
    env.hosts = [env.host_string]
    env.key_filename = node_ssh_key
    env.disable_known_hosts = True
    try:
        _checkout_code(repo, branch, sha1, tag, pr_num)
        _run_remote_integ_test(repo, magma_root)
    finally:
        _release_node_lease(api_url, lease.node_id, lease.lease_id,
                            cert_file, cert_key_file)


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


def _checkout_code(repo: str, branch: str = '', sha1: str = '',
                   tag: str = '', pr_num: str = '',
                   magma_root: str = ''):
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


def _run_remote_integ_test(repo: str, magma_root: str = ''):
    repo_name = _get_repo_name(repo)
    with cd(f'{repo_name}/{magma_root}/lte/gateway'):
        run('fab integ_test')


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
