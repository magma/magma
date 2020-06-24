# Contributing to symphony
We want to make contributing to this project as easy and transparent as
possible.

## Pull Requests

If it's the first time you're create a PR, please follow [this guide](https://kbroman.org/github_tutorial/pages/first_time.html)
to make sure your `git` environment and GitHub account is configured properly.

1. Clone the project and create your branch from `master`:
   - Clone:
     ```console
     git clone git@github.com:facebookincubator/symphony.git
     ```
   - Create a branch:
     ```console
     git checkout -b <branch-name>
     ```
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes and the code lints.
5. For more information about symphony integration environment, see: [integration/README](https://github.com/facebookincubator/symphony/tree/master/integration).
6. Commit your changes.
   - In order to **pull current changes from GitHub** before committing, run: `git pull --rebase --autostash`.
   - Use `git status` and `git diff` to check the working tree status (or install [GitHub for desktop](https://desktop.github.com/) and [GitHub for CLI](https://cli.github.com/)).
   - Use `git add <files...>` to add specific ignored files, or `git add .` to add all files.
   - Use `git commit` to commit your changes (or `git commit -m "commit message` as a shortcut).
   - Push your changes using `git push origin <branch-name>`.
7. [Create a GitHub PR](https://github.com/facebookincubator/symphony/compare).
8. After the CI was passed successfully, request a review from the team members.
   See the [following guide](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/requesting-a-pull-request-review) for more info.
   ![request review](https://help.github.com/assets/images/help/pull_requests/choose-pull-request-reviewer.png)
9. **FB employees:** In the "transition window", we land PRs in Phabricator. In order to land a PR,
   go to [oss symphony](https://www.internalfb.com/intern/opensource/github/repo/2478458042237035/)
   and click "Import" on the relevant PR. This will create a diff with the PR changes (you should approve and land it).
   After the diff was landed, the PR will be closed automatically. 

## Issues
We use GitHub issues to track public bugs. Please ensure your description is
clear and has sufficient instructions to be able to reproduce the issue.

Facebook has a [bounty program](https://www.facebook.com/whitehat/) for the safe
disclosure of security bugs. In those cases, please go through the process
outlined on that page and do not file a public issue.

## License
By contributing to symphony, you agree that your contributions will be licensed
under the LICENSE file in the root directory of this source tree.
