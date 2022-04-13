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
import sys
from typing import Dict, Iterator, List

import github
from fabric.api import env, lcd, local
from github.PullRequest import PullRequest

RELEASE_BRANCHES = (
    'v1.1',
    'v1.2',
    'v1.3',
    'v1.4',
    'v1.5',
    'v1.6',
    'v1.7',
    'magma-5G',
)


def find_release_commits(  # noqa: S107
    repo_name: str = 'magma/magma',
    token_file: str = '~/.magma/github_access_token',
    lookback: str = '100',
):
    """
    find_release_commits goes though the last `lookback` commits on master and
    looks for PRs with a `apply-<RELEASE_BRANCH>` label and without a
    `backported-<RELEASE_BRANCH>` label. These commits are they cherry picked
    to the RELEASE_BRANCH, pushed to origin and a `backported-<RELEASE_BRANCH>`
    label is applied to the PR.

    Args:
        repo_name: Which github repository to look at.
        token_file: Path to file containing a github personal access token.
        lookback: How many commits to look back.
    """
    if not os.path.isfile(os.path.expanduser(token_file)):
        print(
            f'Create a file `{token_file}` with a Github '
            'access token in the contents and nothing else.\n'
            'Create a new Personal Access Token at '
            'https://github.com/settings/tokens and grant it all "repo" '
            "permissions if you don't already have one.",
        )
        sys.exit(1)
    lookback = int(lookback)

    with open(os.path.expanduser(token_file), 'r') as f:
        token = f.read().strip()
    g = github.Github(login_or_token=token)
    repo = g.get_repo(repo_name)

    # Query repo for all closed PR's (getting the first 100 is fine)
    print('Checking Github API for recently merged PRs...\n')
    pulls = repo.get_pulls(
        state='closed', base='master', sort='created',
        direction='desc',
    )
    pulls = pulls[:lookback]
    commits_by_release = {}     # type: Dict[str, List[PullRequest]]
    for pr in pulls:
        labelset = {lab.name for lab in pr.labels}
        for rel in RELEASE_BRANCHES:
            needs_pick = f'apply-{rel}' in labelset and \
                         f'backported-{rel}' not in labelset and \
                         pr.merged
            if needs_pick:
                if rel not in commits_by_release:
                    commits_by_release[rel] = []
                commits_by_release[rel].append(pr)

    if not commits_by_release:
        print('No commits to cherry-pick!')
        return

    print('Found commits to cherry-pick:')
    for rel in RELEASE_BRANCHES:
        if rel not in commits_by_release:
            continue
        num_commits = len(commits_by_release[rel])
        print(f'\t{rel}: {num_commits} commits')
    print('')

    print('Checking out an ephemeral clean Magma copy in .gitclone...\n')
    local('rm -rf .gitclone')
    local('mkdir -p .gitclone/magma')
    local('git config --global user.email "tim@magmacore.org"')
    local('git config --global user.name "backport-bot"')
    local(f'git clone git@github.com:{repo_name} .gitclone/magma')

    for rel in RELEASE_BRANCHES:
        if rel not in commits_by_release:
            continue
        print(f'\n\n===== Beginning cherry-pick procedures for branch {rel} =====\n\n')
        _pick_commits(repo_name, rel, reversed(commits_by_release[rel]))


def _pick_commits(repo: str, rel: str, pulls: Iterator[PullRequest]):
    env.warn_only = True
    with lcd('.gitclone/magma'):
        local('git checkout -- .')
        local(f'git checkout {rel}')

        for pr in pulls:
            sha = pr.merge_commit_sha[:8]
            print('')
            print(f'--- Cherry-picking PR#{pr.number} ({sha}) onto {rel} ---')

            pick_status = local(f'git cherry-pick {sha}')

            print('')
            if not pick_status.succeeded:
                print('')
                print(
                    'Aborting the automated cherry-pick procedure '
                    f'for branch {rel}. Perform the following steps after '
                    'this script finishes execution:',
                )

                print('\t1. cd .gitclone/magma')
                print(
                    f'\t2. Manually cherry-pick commit {sha} onto {rel} '
                    'and resolve conflicts:\n'
                    '\t\tgit status\n'
                    '\t\t<resolve all conflicts manually>\n'
                    '\t\tgit add .\n'
                    '\t\tgit cherry-pick --continue',
                )
                print(
                    '\t3. Push the branch upstream:\n'
                    f'\t\tgit push origin {rel}',
                )
                print(
                    f'\t4. Add the "backported-{rel}" label to '
                    f'PR#{pr.number} at '
                    f'https://github.com/{repo}/pull/{pr.number}',
                )
                print('\t5. Run this fab script again')
                print('')
                sys.exit(1)

            env.warn_only = False
            local(f'git push origin {rel}')
            pr.add_to_labels(f'backported-{rel}')
            print(
                f'Successfully picked PR#{pr.number} ({sha}) onto {rel} and '
                f'marked the PR as backported.',
            )
