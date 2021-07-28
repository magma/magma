---
id: p014_proposal_process
title: Magma Proposals
hide_title: true
---

# Magma Proposals

> **NOTE**: you actually want [How to Open a Proposal](../contributing/contribute_proposals.md) for the current process. This document describes the process as it was originally proposed.

As an open source project, it is important that all contributors are able to
observe substantitive Proposals, comment on them, and be able to discover
historical discussions.
This preserves important historical context for decisions, reduces load on
repeat Proposals, and enables the broad community to suggest and execute
large scale changes.

It is important that this Proposal process be appropriately light-weight,
clearly documented, and exeucted within a reasonable and advertised timeline.

Proposals may take all sorts of shapes, including:

- Design Proposals
- Repository change Proposals
- Process Proposals

## Magma Proposal Process

Magma has elected to base its Proposal Process off the
[Golang Design Proposals Process](https://github.com/golang/proposal#proposing-changes-to-go).

The proposal process is the process for reviewing a proposal and reaching
a decision about whether to accept or decline the proposal.

1. The proposal author
   [creates a brief issue](https://github.com/magma/magma/issues/new)
   describing the proposal.\
   Note: There is no need for a design document at this point.\
   Note: A non-proposal issue can be turned into a proposal by simply adding
         the `type: proposal` label.

2. A discussion on the issue tracker aims to triage the proposal into one of
   three outcomes:
     - Accept proposal, or
     - decline proposal, or
     - ask for a design doc.

   If the proposal is accepted or declined, the process is done.
   Otherwise the discussion is expected to identify concerns that
   should be addressed in a more detailed design.

3. The proposal author writes a [design doc](#design-documents) to work out
   details of the proposed design and address the concerns raised in the
   initial issue discussion.

4. Once comments and revisions on the design doc wind down, there is a final
   discussion on the issue, to reach one of two outcomes:
    - Accept proposal or
    - decline proposal.

After the proposal is accepted or declined (whether after step 2 or step 4),
implementation work proceeds in the same way as any other contribution.

## Detail

### Goals

- Make sure that proposals get a proper, fair, timely, recorded evaluation
  with a clear answer.
- Make past proposals easy to find, to avoid duplicated effort.
- If a design doc is needed, make sure contributors know how to write a good
  one.

### Definitions

- A **proposal** is a suggestion filed as a GitHub issue, identified by having
  the `type: proposal` label.
- A **design doc** is the expanded form of a proposal, written when the
  proposal needs more careful explanation and consideration.

### Scope

The proposal process should be used for any notable change or addition to the
language, libraries, tools, processes.

Since proposals begin (and will often end) with the filing of an issue, even
small changes can go through the proposal process if appropriate.

**If in doubt, file a proposal**.

### Design Documents

As noted above, some (but not all) proposals need to be elaborated in a design
document.

- The design doc should be checked in to
  [magma/docs/readmes/proposals/](https://github.com/magma/magma/tree/master/docs/readmes/proposals)
  as `pNNN_short-name.md`, where `NNN` is the next unused proposal number in
  the directory and `short-name` is a short name (a few dash-separated words
  at most).
  Follow the Magma Github contribution process from the
  [Magma Contributing Conventions](https://docs.magmacore.org/docs/next/contributing/contribute_conventions).

- The design doc should follow [the template](TEMPLATE.md).

- The design doc should address any specific concerns raised during the
  initial issue discussion.

- It is expected that the design doc may go through multiple checked-in
  revisions.
  New design doc authors may be paired with a design doc "shepherd" to help
  work on the doc.

- For ease of change review, design documents should be wrapped around the
  80 column mark.
  [Each sentence should start on a new line](http://rhodesmill.org/brandon/2012/one-sentence-per-line/)
  so that comments can be made accurately and the diff kept shorter.
  **For example, see this document source**.
  - In Emacs, loading `fill.el` from this directory will make
    `fill-paragraph` format text this way.

- Comments on Github Design Doc PRs should be restricted to grammar, spelling,
or procedural errors related to the preparation of the proposal itself.
All other comments should be addressed to the related GitHub issue.

### Proposal Review

The principal goal of the review meeting is to make sure that proposals
are receiving attention from the right people, by cc'ing relevant developers,
raising important questions, pinging lapsed discussions, and generally trying
to guide discussion toward agreement about the outcome.
The discussion itself is expected to happen on the issue tracker,
so that anyone can take part.

A group of Magma team members holds “proposal review meetings”
approximately on a weekly basis to review pending proposals.

The proposal review meetings also identify issues where
consensus has been reached and the process can be
advanced to the next step (by marking the proposal accepted
or declined or by asking for a design doc).

Minutes are posted to [0_proposal-minutes.md](0_proposal-minutes.md)
after the conclusion of the weekly meeting, so that anyone
interested in which proposals are under active consideration
can follow that issue.

The state of Proposal issues are tracked by Github label:
- `proposal:needs:design doc`
  - Proposal has been deemed to merit a more detailed design doc.
  - Please ensure it Addresses the design questions raised in this Github
    Issue.
- `proposal:status:In Review`
  - This proposal is in active review, all discussion will be within this
    Github Issue.
- `proposal:status:Accepted`
  - This proposal has been accepted by general consensus
- `proposal:status:Rejected`
  - This proposal has been rejected by general consensus

## Alternatives Considered

The following open source / public Proposal processes were surveyed and considered.

- [Chromium Design Doc Process](https://chromium.googlesource.com/chromium/src/+/master/docs/contributing.md#design-documents)
- [IETF RFC Process](https://www.ietf.org/standards/process/informal/)
- [Python PEP Process](https://www.python.org/dev/peps/pep-0001/#:~:text=A%20Process%20PEP%20describes%20a,an%20event%20in)%20a%20process.&text=Examples%20include%20procedures%2C%20guidelines%2C%20changes,also%20considered%20a%20Process%20PEP.)
