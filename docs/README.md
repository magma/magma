# Magma documentation

This directory contains the higher-level documentation for the Magma project.

The documentation is structured as a series of READMEs, which are then
organized into a Docusaurus site for easy consumption.

## About

### Docusaurus

[Docusaurus](https://docusaurus.io/) is a framework for static site generation,
focusing on documentation-based sites.

The `sidebars.json` and `siteConfig.js` are the main entrypoints for updating
our documentation site. The former determines which README docs are displayed,
and the latter sets site-wide config values.

There are two categories of documentation: "current" and "versioned". The
current docs are stored under `docs/readmes/`. When cutting a specific release,
we also capture a snapshot of the current docs and store them as versioned
docs under `docs/docusaurus/versioned_docs/`. The sidebars are similarly
versioned, with the versioned sidebars stored under
`docs/docusaurus/versioned_sidebars`.

## Howto

### Make changes

- Add a doc: add the doc to the appropriate path under `docs/readmes/`, then
update `sidebars.json` to include the added doc
- Update a current doc: edit the relevant doc under `docs/readmes/`
- Update a versioned doc: first update the corresponding current doc, then edit
the relevant doc under `docs/docusaurus/versioned_docs`

### View local changes

Use the provided `create_docusaurus_website.sh` script to generate and run a
local server serving your local changes.

### Cut Docusaurus versioned docs

If you want to cut a new versioned documentation release, delete everything
in the Dockerfile after WORKDIR /app/website, then from this directory:

```bash
docker-compose down
docker build --tag magma_docusaurus .
docker-compose up --detach
docker-compose exec docusaurus bash
```

Inside the container,

```bash
yarn install
yarn build
```

Then you can use `yarn` commands to create a new versioned release:

```bash
yarn run version 1.X.0
```

Commit all the new generated files **except for the `docusaurus/node_modules`
directory** and tweak the sidebars if you need to. Run
`create_docusaurus_website.sh` to preview your changes.
