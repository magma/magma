---
id: version-1.6.X-contribute_proposals
title: Proposals Process
hide_title: true
original_id: contribute_proposals
---

# Proposals Process

This document describes the process by which interested parties may propose changes to the Magma project. Proposals can be any size and cover any topic, e.g. process updates or architectural changes.

To start, consider skimming the existing [Magma project proposals](https://github.com/magma/magma/issues?q=is%3Aissue+label%3A%22type%3A+proposal%22+).

## Tracking proposals

Proposals are tracked as GitHub Issues, using the following labels
- [`type: proposal`](https://github.com/magma/magma/issues?q=is%3Aissue+label%3A%22type%3A+proposal%22+) for all proposals
- [`tsc`](https://github.com/magma/magma/issues?q=is%3Aissue+label%3A%22type%3A+proposal%22+label%3Atsc) for proposals desiring TSC consideration

The `#proposals` Slack channel receives notifications when new proposals are created.

## Submit a proposal

### Submit an Issue

[Submit a GitHub Issue](https://github.com/magma/magma/issues/new/choose) following the "Proposal" template.

The proposal contents should be clear and concise. Aim for a "one-pager" style.

### Receive feedback

The project maintainers (or the TSC, if you applied the `tsc` label) will discuss the proposal and make comments on the Issue. One of four outcomes will be communicated via labels

- [`status: accepted`](https://github.com/magma/magma/labels/status%3A%20accepted) proposal was accepted
- [`status: rejected`](https://github.com/magma/magma/labels/status%3A%20rejected) proposal was rejected
- [`status: withdrawn`](https://github.com/magma/magma/labels/status%3A%20withdrawn) proposal was withdrawn by the author
- [`status: needs design doc`](https://github.com/magma/magma/labels/status%3A%20needs%20design%20doc) proposal needs a design document

If your proposal is accepted, ensure actionable next-steps have been outlined and initiated.

### Optional: needs design doc

If your proposal is labeled as needing a design doc, this means your proposal is nominally accepted, but needs to progress from the one-pager GitHub Issue into a fleshed-out design.

This design doc can start as e.g. a Google Doc or Quip Doc, but needs to eventually make its way into a pull request to add it to the project's Docusaurus.

When writing your design doc, consider following this [standardized design doc template](https://www.industrialempathy.com/posts/design-docs-at-google/).

- Example design doc: [Scaling Orc8r Subscribers Codepath](https://magma.github.io/magma/docs/next/proposals/p010_subscriber_scaling)
- Example design doc PR: [APN Refactoring](https://github.com/magma/magma/pull/7191)
