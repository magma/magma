---
id: docusaurus
title: Docusaurus
hide_title: true
---
# Docusaurus
### Generating the Documentation Website

1. Ensure [docker](https://docs.docker.com/install/) is installed
2. From `magma/docs`, run `./docusaurus/create_docusaurus_website.sh`. This copies the generated files to `docs/web` and deletes any old generated files.
3. Navigate to http://127.0.0.1:3000/web/ to view a local version of the site.

### Directory Structure

The documentation website is generated using [docusaurus](https://docusaurus.io/) from
the README files stored in `docs/readmes/`. The generated website files are
stored in `docs/web`. The docusaurus files needed to generate the website are
stored in `docs/docusaurus`.
