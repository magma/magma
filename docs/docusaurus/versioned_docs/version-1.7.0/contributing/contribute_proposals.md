---
id: version-1.7.0-contribute_proposals
title: Proposals Process
hide_title: true
original_id: contribute_proposals
---

# Proposals Process

This document describes the process by which interested parties may propose changes to the Magma project. Proposals can be any size and cover any topic, e.g. process updates or architectural changes.

To start, consider skimming the existing [Magma project proposals](https://github.com/magma/magma/issues?q=is%3Aissue+label%3A%22type%3A+proposal%22+).

## Overview

### Acceptance

To be accepted, proposals must receive a majority vote from the TSC.

### Tracking

Proposals are tracked as GitHub Issues, using the [`type: proposal`](https://github.com/magma/magma/issues?q=is%3Aissue+label%3A%22type%3A+proposal%22+) label.

The `#proposals` Slack channel receives notifications when new proposals are created.

## Submit a proposal

### Submit an Issue

[Submit a GitHub Issue](https://github.com/magma/magma/issues/new/choose) following the "Proposal" template. The proposal contents should be clear and concise. Aim for a "one-pager" style.

In many cases, the TSC will delegate their decision to one of the [approvers teams](https://github.com/orgs/magma/teams?query=approvers-). To expedite this process, you can [label your proposal with the applicable component](https://github.com/magma/magma/labels?q=component%3A), which will help the TSC identify domain experts.

### Receive feedback

The TSC, and relevant project maintainers, will discuss the proposal and make comments on the Issue. One of four outcomes will be communicated via labels

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

### Resolution

A proposal Issue will be closed when no further discussion is needed. This occurs after one of the following

- Acceptance, rejection, or withdrawal of the Issue
- Acceptance of a requested design doc

## Conclusion

Please direct process-related questions to the `#governance-tsc` Slack channel.
