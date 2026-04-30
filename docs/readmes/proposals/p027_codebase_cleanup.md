---
id: p027_codebase_cleanup
title: Codebase Cleanup and Project Revitalization
hide_title: true
---

# Proposal: Codebase Cleanup and Project Revitalization

- *Author: @shahrahman*
- *Last Updated: 02/20/26*
- *Status: In Review*

## Problem

Magma has had zero commits to `master` since April 2025.
Of the 28 listed CODEOWNERS, only 1 remains active.
The project has 94 unreviewed PRs, 390 open issues, 109
unmerged branches, and critical security vulnerabilities in
shipping configurations.
CI pipelines target `ubuntu-20.04` (EOL) and use deprecated
GitHub Actions versions.
Dependencies span multiple years of staleness across Go,
Python, and Node ecosystems.

To the outside world, the project appears abandoned.
Every day this persists, potential contributors, adopters,
and sponsors are lost.

Before Magma can pursue any new capabilities — including
5G SA, CAMARA APIs, or NWDAF — the project must demonstrate
that it is actively maintained, secure, and buildable.

## Goals

1. **Make the project visibly alive.**
   Delete stale branches, triage the PR and issue backlog,
   and publish a "State of Magma" announcement.

2. **Remediate critical security vulnerabilities.**
   Fix shell command injection in magmad, upgrade vulnerable
   NMS authentication libraries, and audit committed key
   material.

3. **Unblock development.**
   Upgrade CI runners, GitHub Actions, and language toolchains
   so that PRs can be tested and merged.

4. **Reduce technical debt.**
   Standardize dependency versions across 20+ Go modules,
   upgrade outdated Python and Node packages, and remove dead
   code.

5. **Rebuild community trust.**
   Update governance documents, CODEOWNERS, and issue
   templates.
   Create a clear pipeline for new contributors.

## Non-goals

- New feature development (5G SA, NEF, NWDAF, etc.) —
  covered by separate proposals.
- Architectural changes to Orc8r, AGW, or FeG.
- NMS redesign or React 18 migration — requires its own
  evaluation.
- Changes to the Magma release process or versioning scheme.

## Proposal

The cleanup is organized into 9 phases, ordered by risk
and impact.
Phases are independent enough to be executed in parallel
where resources allow, but the priority ordering ensures
the highest-value work happens first.

### Phase 1: Zero-Risk Repository Hygiene

Low-risk, high-visibility changes that immediately signal
active maintainership.

**1.1 Delete stale branches (~80+ branches)**

- 2 merged branches safe to delete immediately
- 30+ stale dependabot branches — close associated PRs
  and delete
- 36+ abandoned feature branches (pre-2023, no activity)
- Old release branches (v1.0 through v1.7, v1.8\_backport,
  v1.11.11) — keep only `master`, `v1.8`, `v1.9`,
  `gh-pages`
- 8 recently active branches — contact authors before
  any action (Lucas Gonze, mickymkumar, shashidhar-patil)

**1.2 Triage open PRs (94 open)**

- Close stale PRs (no activity >6 months) with a
  standardized message inviting resubmission
- Batch viable dependabot PRs into consolidated updates
- Prioritize review of recent community PRs and
  mentorship contributions

**1.3 Triage open issues (390 open)**

- Label sweep: ensure all issues have component labels
- Close stale issues (no activity >1 year, no repro steps)
- Consolidate duplicates
- Mark "won't fix" for deprecated components
  (`openwrt`, `experimental`)
- Refresh `good first issue` labels

**1.4 Update GitHub project board**

- Archive stale project boards
- Add issue templates (bug report, feature request,
  support question) — one of only 2 missing community
  health indicators

### Phase 2: Security Remediation

Critical security issues that should be addressed before
any other code changes.

**2.1 Committed keys and secrets**

The `keys/` directory contains cryptographic key material.
Determine whether these are test or production keys.
If production: rotate, scrub from history, add to
`.gitignore`.
If test: add clear documentation.

**2.2 Shell command injection**

The `generic_command_config` in magmad.yml files across
AGW, FeG, and Orc8r gateway configs allows `bash {}`
execution with `allow_params: True`.
Anyone with Orchestrator API access can execute arbitrary
shell commands on gateways.

Actions:
- Remove `bash` from default configs
- Add input sanitization to `shell_command_executor.py`
- Implement a command allowlist

Files affected:
- `orc8r/gateway/configs/magmad.yml`
- `lte/gateway/configs/magmad.yml`
- `feg/gateway/configs/magmad.yml`

**2.3 Subprocess shell=True audit**

Audit all Python `subprocess` calls using `shell=True`
for user-input injection paths.
Replace with `shell=False` and explicit argument lists
where safe.

**2.4 Weak cryptographic algorithms**

`PipelineD_HEConfig` in `lte/protos/mconfig/mconfigs.proto`
includes RC4 and MD5 options, both cryptographically broken.
Deprecate these options and add documentation warnings.

**2.5 NMS authentication library upgrades**

| Package | Current | Issue |
|---------|---------|-------|
| `passport` | 0.4.0 | Released 2014, session fixation |
| `openid-client` | 2.4.5 | Missing OIDC security fixes |
| `axios` | 0.21.2 | CVE-2023-45857, CVE-2023-26159 |
| `express-session` | 1.15.6 | Missing security patches |
| `helmet` | 4.0.0 | Missing modern CSP directives |

These protect user authentication and must be upgraded
as a priority.

### Phase 3: CI/CD Modernization

**3.1 Upgrade GitHub Actions**

All 48+ workflow files use outdated action versions:

| Action | Current | Target |
|--------|---------|--------|
| `actions/checkout` | v3 | v4 |
| `actions/upload-artifact` | mixed v3/v4 | v4 |
| `actions/cache` | v3 | v4 |
| `actions/setup-go` | varies | v5 |
| `actions/setup-python` | varies | v5 |

**3.2 Upgrade runner OS**

All workflows use `ubuntu-20.04` (EOL April 2025).
Migrate to `ubuntu-22.04` or `ubuntu-24.04`.

**3.3 Upgrade language versions in CI**

- Go: currently 1.20-1.21 in CI, target 1.22+
- Python: currently 3.8.10 (EOL) in CI, target 3.11+

**3.4 Fix broken integrations**

- Verify `secrets.SLACK_WEBHOOK` — if Slack workspace is
  dead, remove notification steps
- Consolidate inconsistent clang-format linter versions

**3.5 Docker image modernization**

- Remove `ubuntu:18.04` Dockerfiles (EOL April 2023)
- Upgrade `python:3.9.10-alpine` base images

### Phase 4: Dependency Upgrades

**4.1 Go module standardization**

20+ Go modules have inconsistent dependency versions:

| Dependency | Range | Target |
|------------|-------|--------|
| gRPC | v1.43.0 – v1.64.1 | Standardize v1.62+ |
| `golang.org/x/crypto` | v0.0.0 – v0.23.0 | Standardize latest |
| `golang.org/x/net` | v0.7.0 – v0.25.0 | Standardize latest |
| `go-redis/redis` | v6.14.1 (deprecated) | v9+ |
| `prometheus/client_golang` | v1.12.2 | v1.19+ |
| `Masterminds/squirrel` | v1.1.1-pre (2018) | v1.5+ |

All modules should target Go 1.22+ minimum.

**4.2 Python dependency upgrades**

- `grpcio`: currently 1.37.1-1.53.2, target 1.62+
- `PyYAML`: 5.4.1, target 6.0.2
- `prometheus_client`: 0.3.1, target 0.20+
- `SQLAlchemy`: 1.4.15, target 2.0+
- Minimum Python version: 3.10+

**4.3 Node/NMS dependency upgrades**

Beyond the security-critical packages in Phase 2:

| Package | Current | Target |
|---------|---------|--------|
| `react` | 17.0.2 | 18.x |
| `webpack` | 4.46.0 | 5.x |
| `typescript` | 4.6.4 | 5.x |
| `eslint` | 7.3.2 | 9.x |
| `immutable` | 4.0.0-rc.12 | 4.3+ (pre-release!) |

**4.4 Bazel external dependencies**

| Rule | Current | Target |
|------|---------|--------|
| `rules_python` | v0.5.0 | v0.35+ |
| `rules_go` | v0.28.0 | v0.48+ |
| `bazel-gazelle` | v0.23.0 | v0.37+ |
| Protobuf | v3.19.1 | v25+ |

### Phase 5: Dead Code and Stale Component Removal

**5.1 Evaluate stale top-level directories**

| Directory | Status | Recommendation |
|-----------|--------|----------------|
| `experimental/` | Cloudstrapper, last updated 2021-22 | Archive or delete |
| `openwrt/` | Minimal content, very stale | Delete |
| `show-tech/` | Ansible diagnostics, stale | Archive or delete |
| `cn/` | 5G Core (experimental) | Evaluate, likely delete |
| `hil_testing/` | Requires specific hardware | Evaluate |

**5.2 Clean deprecated protobuf references**

5 DEPRECATED markers in `mconfigs.proto`, 6 deprecated
fields across NAS/NGAP state protos.
Proto fields cannot be removed (wire compatibility), but
referencing code can be cleaned up.

**5.3 Remove legacy NMS components**

`nms/app/views/alarms/legacy/` — evaluate and remove
if no longer used.

**5.4 Clean deprecated Python patterns**

- `@asyncio.coroutine` decorators (issue #15602 filed)
- Unnecessary `__future__` imports

**5.5 Evaluate generated code in version control**

Generated `*.pb.go`, `*_swaggergen.go`, and `nms/generated/`
files are tracked in git.
Evaluate generating these in CI instead to reduce repo size
and eliminate generated-code merge conflicts.

### Phase 6: Code Quality and Modernization

- Run `go mod tidy` and `golangci-lint` across all modules
- Update Python linting (`flake8` 3.9→7.0, `mypy` checks)
- Update TypeScript linting (ESLint 7→9)
- Run `clang-format` on all C/C++ files
- Verify BSD-3-Clause license headers across all source
  files

### Phase 7: Documentation Cleanup

- Replace inaccessible `fb.quip.com` links with public
  equivalents
- Update stale documentation (179 markdown files in
  `docs/readmes/`)
- Update README badges and installation instructions
- Verify Docusaurus site builds and navigation
- Update CODEOWNERS to reflect active maintainers

### Phase 8: Test Infrastructure

- Replace 76 `time.Sleep`/`sleep` calls in CWF integration
  tests with polling/event-driven checks
- Fix tests using `@flaky(max_runs=3)` — root cause, don't
  retry
- Audit `@unittest.skip` / `t.Skip()` / `xit(` for tests
  whose underlying issues may be fixed
- Verify `codecov.yml` is generating and uploading reports
- Evaluate containerizing integration tests (currently
  require 3 Vagrant VMs)

### Phase 9: Governance and Community

- Update `CONTRIBUTING.md` with current processes
- Reconstitute TSC membership with active participants
- Update CODEOWNERS:
  - Add active contributors (Lucas Gonze, Jordan Vrtanoski,
    Lucas Amaral)
  - Remove inactive entries
- Add GitHub issue templates (bug report, feature request,
  support question)
- Review Dependabot configuration for full ecosystem
  coverage
- Verify Slack workspace and GitHub Discussions health
- Publish "State of Magma" announcement:
  new maintainership, cleanup progress, call for
  contributors

## Execution Priority

| Priority | Scope | Risk | Impact |
|----------|-------|------|--------|
| **P0** | Phase 2: Security fixes | Medium | Critical |
| **P1** | Phase 1: Branch/PR/issue cleanup | None | High visibility |
| **P2** | Phase 3: CI/CD modernization | Low | Unblocks development |
| **P2** | Phase 2.5: NMS auth upgrades | Medium | Security |
| **P3** | Phase 4: Dependency upgrades | Medium | Maintainability |
| **P3** | Phase 5: Dead code removal | Low | Reduces confusion |
| **P4** | Phase 6: Code quality | Low | Long-term health |
| **P4** | Phase 7: Documentation | None | Community trust |
| **P5** | Phase 8: Test infrastructure | Low | Reliability |
| **P5** | Phase 9: Governance | None | Community building |

## Estimated Scope

| Phase | Approximate PRs | Files Touched |
|-------|-----------------|---------------|
| Phase 1 | 3-5 | ~0 code (branch/issue mgmt) |
| Phase 2 | 5-8 | ~20 files |
| Phase 3 | 5-10 | ~50 workflow files |
| Phase 4 | 10-15 | ~40 dependency files |
| Phase 5 | 5-8 | Hundreds (deletions) |
| Phase 6 | 10+ | Hundreds (formatting) |
| Phase 7 | 3-5 | ~50 docs files |
| Phase 8 | 5-10 | ~30 test files |
| Phase 9 | 2-3 | ~5 governance files |

## Alternatives Considered

**Do nothing.**
The project continues to lose visibility and relevance.
The security vulnerabilities remain.
Potential contributors see a dead project and move on.
This is not viable if Magma is to have a future.

**Minimal cleanup (security only).**
Fixes the most urgent issues but does not address the
perception problem.
94 unreviewed PRs and 109 stale branches still signal
abandonment.
Rejected because visibility is as important as security
for project revival.

**Full rewrite / fresh start.**
Abandons 400k+ lines of production-tested code.
Loses the existing community, deployment base, and
institutional knowledge.
Rejected because the existing codebase is fundamentally
sound — it needs maintenance, not replacement.

## Good First Issues for Community Contributors

Several cleanup tasks are well-scoped for new contributors:

- Python `@asyncio.coroutine` → `async def` migration
  (#15602)
- C-style casting fixes (#14977)
- `__future__` import cleanup
- Broken documentation link fixes (fb.quip.com references)
- `subprocess` `shell=True` → `shell=False` conversions
- License header consistency checks
- Individual `go mod tidy` runs per module
