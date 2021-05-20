---
id: contribute_workflow
title: Development Workflow
hide_title: true
---

# Development Workflow

This document covers the development workflow for contributing to the Magma project.

## General workflow

Magma follows the standard "fork and pull-request" development workflow. For more information on this workflow, consider the following

- [Howto: fork and make pull request](https://guides.github.com/activities/forking/) [(visual version)](https://jarv.is/notes/how-to-pull-request-fork-github/)
- [In-depth: fork and make pull request workflow](https://gist.github.com/Chaser324/ce0505fbed06b947d962)

See the [opinionated workflow](#opinionated-workflow) section below for a low-friction version of these workflows.

Once your PR has been approved and passed all CI checks, use the [`ready2merge`](https://github.com/magma/magma/labels/ready2merge) label to indicate the maintainers can merge it.

## Guidelines

**Required: commits must be signed off.** You must sign-off all commits on the originating branch for a PR, using [the `--signoff` flag in Git](https://stackoverflow.com/questions/1962094). There is a CI check that will fail if any commit in your branch has an unsigned commit. If you've forgotten to sign-off a commit, you can `git commit --amend --signoff`, or `git rebase --signoff` to sign-off an entire branch at once.

**Required: label backward-breaking pull requests.** Use the [`breaking change` label](https://github.com/magma/magma/issues?q=label%3A%22breaking+change%22+). All breaking changes and their mitigation steps will be aggregated in the subsequent release notes. A breaking change fits one or more of the following criteria

1. Will require a manual intervention on the next upgrade (e.g. data migration script)
2. Will break the southbound interfaces between AGW and Orchestrator in a way that will require both components to upgrade in a coordinated way.

**Desired: convincing test plan.** For non-trivial PRs, codeowners will require you to include a test plan detailing any manual verification steps you took.

**Desired: use [imperative mood](https://chris.beams.io/posts/git-commit/) for pull request title and description.** A simple way to check this is to mentally prepend the phrase "This commit will ..." to the title and each bullet point of your description.

## Productivity tools

- Git
    - [`hub`](https://github.com/github/hub): Git CLI wrapper providing extended functionality, e.g. opening a PR
- GitHub
    - [Refined GitHub](https://github.com/sindresorhus/refined-github): browser extension simplifying GitHub interface and adding some useful features
    - [Notifier for GitHub](https://github.com/sindresorhus/notifier-for-github): browser extension to get GitHub notifications sent to your desktop
    - [Sourcegraph](https://docs.sourcegraph.com/integration/browser_extension): browser extension to add [Sourcegraph](https://sourcegraph.com/github.com/magma/magma) functionality directly to GitHub
        - [Codecov](https://sourcegraph.com/extensions/sourcegraph/codecov): Sourcegraph extension to display code coverage metrics
        - [Code ownership](https://sourcegraph.com/extensions/sourcegraph/code-ownership): Sourcegraph extension to show codeowners for each file
        - [Open in IntelliJ](https://sourcegraph.com/extensions/sourcegraph/open-in-intellij): Sourcegraph extension to open a file in your dev machine's IntelliJ instance
- IntelliJ
    - [Native support for GitHub code reviews](https://www.youtube.com/watch?v=MoXxF3aWW8k&ab_channel=IntelliJIDEAbyJetBrains)

## Opinionated workflow

The references above should be sufficient to get started. This section provides an opinionated view on how to efficiently manage this workflow, assuming the following

- up-to-date version of Git
- `hub` wrapper command installed (see above)

### Workflow

This opinionated workflow simplifies the complexities of dealing with Git directly.

It does this by keeping both `your_dev_branch` and `your_dev_branch_base` branches, allowing a straightforward mechanism for rebasing commit stacks.

```bash
# Note: see below for the Git and shell aliases used in this code block

###############
# Get started #
###############

# Fork the Magma repo
# ...

# Clone your forked Magma repo
git clone git@github.com:YOUR_USERNAME/magma.git && cd magma

# Set upstream
git remote add upstream git@github.com:magma/magma.git

###########
# Open PR #
###########

# Checkout master and fast-forward to match upstream
git-update

# Create new dev branch on top of latest master
git-update YOUR_NEW_DEV_BRANCH

# Make changes, and package as a *single* commit on your dev branch
# ...

# Open pull request
git open-pr

# Make requested changes, without committing anything
# ...

# Amend the single commit and force-push to update the PR
git amend-pr

# [optional] Rebase PR onto master
git-rebase

# [optional] Rebase PR onto another PR
git-rebase TARGET_BRANCH

# [note] If there's a merge commit during git-rebase, after completing the
# merge, you need to finish updating the reference branches
git-rebase-finish

###############
# Look around #
###############

# Full commit graph -- super helpful for figuring out what's going on
git graph-all

# Diff between PR and trunk
git diff-base

# Diff between local dev branch and origin dev branch
git diff-origin

############
# Clean up #
############

# When you have no open PRs, you can clear all local branches to clean up
# your git commit graph
git delete-local-branches
```

### Necessary aliases

```gitconfig
# ~/.gitconfig

[alias]
	delete-local-branches = !git branch | grep --invert-match master | xargs git branch --delete
	commit-amend = commit --signoff --amend --no-edit
	diff-base = !git diff $(git branch --show-current)_base
	diff-origin = !git diff origin/$(git branch --show-current)
	graph = graph-all --max-count=30
	graph-all = log --graph --all --format=format:'%C(auto)%h%C(reset) %C(cyan)(%cr)%C(reset)%C(auto)%d%C(reset) %s %C(dim white)- %an%C(reset)'
	amend-pr = !git add --all && git commit --signoff --amend --no-edit && git push --force
	open-pr = !git push-branch && git pull-request --browse
	push-branch = push --set-upstream origin HEAD
```

```bashrc
# shell rc file

# git-update updates master with upstream changes, and optionally creates a feature branch.
function git-update() {
    git checkout master && git pull upstream master && git push origin master

    local br=${1}
    if [[ $br != "" ]]; then
        local br_base=${br}_base
        git branch ${br_base}
        git checkout -b ${br}
    fi
}

# git-rebase rebases current branch on master, or the specified target.
# $1    target
# $2+   args passed to rebase command
#
# Note: if there's a merge conflict, after handling the merge conflict,
# you need to finish by running git-rebase-finish.
function git-rebase() {
    local to=${1:-master}
    local args=${@:2}
    local br=$(git branch --show-current)
    local br_base=${br}_base

    # Save values to file in case rebase fails
    echo "${to} ${br} ${br_base}" > ~/.gitrebase

    git rebase --onto ${to} ${br_base} ${br} ${args} && git checkout ${br_base} && git reset --hard ${to} && git checkout ${br}
}

# git-rebase-finish completes the rebase started by git-rebase.
function git-rebase-finish() {
    local vals
    read -r -a vals < ~/.gitrebase
    local to=${vals[0]}
    local br=${vals[1]}
    local br_base=${vals[2]}
    git checkout ${br_base} && git reset --hard ${to} && git checkout ${br}
}
```
