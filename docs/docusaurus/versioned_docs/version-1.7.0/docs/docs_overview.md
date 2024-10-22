---
id: version-1.7.0-docs_overview
title: Documentation Overview
hide_title: true
original_id: docs_overview
---

# Documentation Overview

This directory contains the higher-level documentation for the Magma project.

The documentation is structured as a series of READMEs, which are then organized into a Docusaurus site for easy consumption.

## About

### Docusaurus

[Docusaurus](https://docusaurus.io/) is a framework for static site generation, focusing on documentation-based sites.

The [`sidebars.json`](https://v1.docusaurus.io/docs/en/navigation) and [`siteConfig.js`](https://v1.docusaurus.io/docs/en/site-config) are the main entrypoints for updating our documentation site. The former determines which README docs are displayed, and the latter sets site-wide config values.

There are two categories of documentation: "current" and "versioned". The current docs are stored under `docs/readmes/`. When cutting a specific release, we also capture a snapshot of the current docs and store them as versioned docs under `docs/docusaurus/versioned_docs/`. The sidebars are similarly versioned, with the versioned sidebars stored under `docs/docusaurus/versioned_sidebars`.

### Conventions

Follow the existing conventions when naming and placing new READMEs. Notably,

- Use short, concise verbs as section and document titles (e.g. "Upgrade", "Deploy", "Debug")
- Guide names and format should match across Magma components (i.e. Orc8r's "Upgrade to v1.4" and AGW's "Upgrade to v1.4" should flow together logically)
- A document's ID should match its filename

Some examples of proper naming

- Upgrade guide (upgrade to v1.4)
    - Title: `Upgrade to v1.4`
    - ID: `upgrade_1_4`
    - Filename: `lte/upgrade_1_4.md`, `orc8r/upgrade_1_4.md`
- Deploy guide (install)
    - Title: `Install Orchestrator`, `Install Access Gateway`
    - ID: `deploy_install`
    - Filename: `lte/deploy_install.md`, `orc8r/deploy_install.md`

## Howto

### Make changes

- Add a doc: add the doc to the appropriate path under `docs/readmes/`, then update `sidebars.json` to include the added doc
- Update a current doc: edit the relevant doc under `docs/readmes/`
- Update a versioned doc: first update the corresponding current doc, then edit the relevant doc under `docs/docusaurus/versioned_docs`

### Precommit

Before pushing, verify your Markdown passes the linter. This is enforced by the `markdown-lint` [CI check](../contributing/contribute_ci_checks.md). See [Markdownlint rules](https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md) and [Magma's Markdownlint config](https://github.com/magma/magma/blob/master/docs/readmes/.markdownlint.yaml) for more info

```bash
# Check if changes pass Markdown lints
cd ${MAGMA_ROOT}/docs && make precommit

# Automatically fix some of the lints
cd ${MAGMA_ROOT}/docs && make precommit_fix
```

### View local changes

Generate and run a local server serving your changes

```bash
cd ${MAGMA_ROOT}/docs && make dev
```

### Cut Docusaurus versioned docs

If you want to cut a new versioned documentation release, bring up the Docusaurus container

```bash
cd ${MAGMA_ROOT}/docs/docusaurus
docker-compose down
docker build --tag magma_docusaurus .
docker-compose -f docker-compose.yml -f docker-compose.publish.yml up --detach
docker-compose exec docusaurus bash
```

From inside the container, build the site

```bash
yarn install
yarn build
```

Now you can create a new versioned release

```bash
yarn run version X.Y.0  # e.g. version 1.5.0
```

Commit all the new generated files and tweak the sidebars if you need to. Run `make dev` to preview your changes.
