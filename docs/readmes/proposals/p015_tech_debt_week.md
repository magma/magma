---
id: p015_tech_debt_week
title: Tech Debt Week Processes
hide_title: true
---
# Proposal: Tech Debt Week (TDW) Processes

Author: @electronjoe

Last updated: 04/21/21

## What is Tech Debt Week

Tech debt weeks are an opportunity for all Magma core contributors, plus new and existing comunity stakeholders - to focus on burning down tech debt topics.  The objective is to provide responsive and focused support for the following.

- GH Issue dissemination
- GH Issue clarification and or prototype demonstration
- GH PR early feedback
- GH PR final review
- GH PR merge support

## Getting Help

There are several options for getting help - reach out early and often! Preferred contact will be provided to you based upon the GH Issue you are working, but additionally the following general methods always exist.

- GH user tag (if not at all urgent)
- [Video conference](TBD LINK) monitored by [TDW On-Call](#magma-tdw-on-call)
- Slack channel [#project-tech-debt-week-apr-2021](https://magmacore.slack.com/archives/C0202TW1VRN) monitored by [TDW On-Call](#magma-tdw-on-call)
- Page: email the TDW On-call at tdw@magma.pagerduty.com
  - Provide a preferred contact (e.g. slack id, github id, ...)

## Process for Contributors

1. Review the list of [TDW Available Work]()
   1. What `Unassigned` tasks interest you?
2. Reach out to [TDW On-Call](#tdw-on-call-contact)
   1. Settle on a task
   2. You will be provided contact information for an [Issue Lead](#magma-issue-lead)
   3. You now have a task! Fun!
3. Review the GH issue
   1. Ensure a comment has been filed tagging the [Issue Lead](#magma-issue-lead) GH id for tracking
   2. If possible, assign the GH Issue to yourself (`Assignees` at right side)
      1. Otherwise comment on the GH issue - to make clear you are assigned
   3. See if you understand the GH Issue, prepare questions
4. Reach out to [Issue Lead](#magma-issue-lead) using the provided contact method
   1. Check with the lead to ensure you understand the problem
   2. And the preferred solution approach
   3. If the Issue lacks detail, do not hesitate to ask
      1. For an example of the shape of the solution
      2. To existing committed code that is similar
5. Create an early Draft PR
   1. Work towards a draft PR early, post to GH as a [Draft](https://github.blog/2019-02-14-introducing-draft-pull-requests/) and
      1. Directly notify your [Issue Lead](#magma-issue-lead) over preferred communications
      2. Post it in Slack [#project-tech-debt-week-apr-2021](https://magmacore.slack.com/archives/C0202TW1VRN)
      3. You can convert to draft if you mistakenly create full fledge PR (`Still in progress? Convert to draft`)
   2. Review Timeliness and Tone (for both Draft and Final)
      1. Expect [Timely review](#issues-lead-response-time)
      2. With a [Spheherd menality](https://github.com/magma/magma/blob/master/docs/readmes/contributing/contribute_codeowners.md#shepherds-not-gatekeepers)
      3. Feel free to ask for precise assistance or examples (e.g. a GH gist with example code on an ask)
   3. The [Issue Lead](#magma-issue-lead) will aim to make contact daily if they haven't heard from you
   4. Apply any directional changes advised to draft PR
6. When a final PR is ready for review
   1. Click `Ready for Review` to convert from Draft PR to PR
   2. Directly notify your [Issue Lead](#magma-issue-lead) over preferred communications
   3. Post it in Slack [#project-tech-debt-week-apr-2021](https://magmacore.slack.com/archives/C0202TW1VRN)
   4. See above Review Timeliness and Tone
7. If Review is completed but you are blocked by CI Checks
   1. Reach out to your [Issue Lead](#magma-issue-lead) or the [TDW On-Call](#tdw-on-call-contact) for help
8. When Review is complete and CI Checks are clear
   1. Reach out to [Issue Lead](#magma-issue-lead) or [TDW On-Call](#tdw-on-call-contact) to Merge

## Magma TDW On-Call

### TDW On-call Contact

- [Video conference](TBD LINK) monitored by [TDW On-Call](#magma-tdw-on-call)
- Slack channel [#project-tech-debt-week-apr-2021](https://magmacore.slack.com/archives/C0202TW1VRN) monitored by [TDW On-Call](#magma-tdw-on-call)
- Page: email the TDW On-call at tdw@magma.pagerduty.com
  - Provide a preferred contact (e.g. slack id, github id, ...)

### TDW On-Call Response time

During TDW, We ask that the following interactions and actions are given `highest work priority` - ideally with a response time of `< 1 hour` between the hours of `9AM and 3PM Pacific Time`.  This is to ensure maximum velocity for our external contributors during this tight week work window.

### TDW On-Call Responsibilities

1. Monitor for communications
   1. Pop open the video conference
      1. If empty you can mute yourself and turn off video
      2. Make sure sound is on!
   2. GH Taggs of your user (consider [slack integration](https://slack.github.com/))
   3. Slack channel [#project-tech-debt-week-apr-2021](https://magmacore.slack.com/archives/C0202TW1VRN)
   4. Pages
2. General Assistance
   1. You may be contacted for help with CI failures
   2. You may be contacted for help with the final Merge
3. Contribution life cycle for TDW On-Call
   1. Find appropriate task for contributor
      1. Area of expertise
      2. Complexity
      3. Compare with task scope and provided detail
   2. Reach out to potential [Magma Issue Lead](#magma-issue-lead) and confirm ownership
      1. With confirmation, reply to GH Issue with text `Issue Lead: <username tag>`
      2. Update status in `TDW Available Work` coordinating GH Issue
         1. Unassigned -> Assigned and add tag to GH user assigned
   3. Put contributor in touch with GH Issue lead
      1. Generally the person filing the issue
      2. Otherwise a CODEOWNER
      3. Wait for notifications from GH Issue Lead
         1. E.g. Completed or MIA
      4. Update status in `TDW Available Work` coordinating GH Issue
         1. Assigned -> Merged
4. Handle MIA contributions
   1. If a [Magma Issue Lead](#magma-issue-lead-responsibilities) reaches out about lack of contact with a contributor
      1. Attempt direct outreach for 2 days (on top of 2 days issue lead waited)
      2. Move GH Issue back to available stack
         1. Assigned -> Unassigned

## Magma Issue Lead

A Magma Issue Lead is the designated individual who has filed a GH Issue deemed appropriate for TDW inclusion.  These individuals may be a CODEOWNER for the relevant code base, but at least have a clear understanding of the issue's resolution and has strong communications channels with the CODEOWNERS implicated.  Determining who is Issue Lead will occur between TDW On-Call and relevant Mamga core team members.

### Issues Lead Response time

During TDW, We ask that the following interactions and actions are given `highest work priority` - ideally with a response time of `< 1 hour` between the hours of `9AM and 3PM Pacific Time`.  This is to ensure maximum velocity for our external contributors during this tight week work window.

### Issues Lead Responsibilities

The following describes the contribution life cycle for a Magma Issue Lead.

1. Receive first contact from contributor assigned your GH Issue
   1. This contact will be made by e.g. the [TDW on-call](#magma-tdw-on-call-responsibilities)
   2. Ensure that the Magma TDW On-Call tagged your GH username in the issue
      1. With text `Issue Lead: <username tag>` for clarity to all
2. Validate alignment of contributor withe GH Issue
   1. Understanding of the problem statement and solution
   2. Appropriate matching of skills and knowledge
      1. If possibly a mismatch, check with TDW on-call about re-matching
3. Give contributor preferred communications method
   1. E.g. slack username, GH tagging, etc - your call
4. Ask for an early Draft PR
   1. So that any substantial direction correction can happen early
   2. Check in at least once a day on the progress - even if you haven't heard anything
   3. In providing PR feedback (draft or other)
      1. Apply the [Spheherd menality](https://github.com/magma/magma/blob/master/docs/readmes/contributing/contribute_codeowners.md#shepherds-not-gatekeepers) to PR review
      2. Provide early directional feedback
   4. Once directionally correct - ask for outreach on final PR
   5. Repeat the `reach out any time` mentality and preferred contact method
5. Wait for final PR notification
   1. Apply timely review
   2. Check CI Checks, help contributor resolve any failures
   3. Once all CI checks are clear and review is good, confirm Mereg with contributor
   4. Merge PR and notify `TDW On-Call` of the Merge for accounting
6. If at any stage communications halts for more than two days
   1. Notify [TDW on-call](#magma-tdw-on-call-responsibilities)

## Magma's TDW Planning and Execution Process

**If you are not personally planning another Magma TDW, you can stop reading here.**

The following is an attempt to generate a repeatable recipe for `Tech Debt Weeks` at Magma.

1. Getting the word out via early and pro-active outreach
   1. 1:1 notification of all partners
   2. Announcement in Magma community meetings
   3. Announcement in Slack
   4. Announcement on GH Discussions
   5. Announcement in Magma Conferences
2. Collect diverse set of well scoped, well defined GH Issues
   1. Across all segments of the code base
      1. Continuous Integration
      2. c/c++ AGW
      3. Python AGW
      4. Golang (Orc8r, FEG, CWF)
      5. Documentation
      6. Kubernetes Config
      7. Installation Process
   2. They should be very scoped (Say 4-8 hour efforts max)
   2. They should be very clear in their ask
      1. Ideally even with sample code and precise insertion links to GH files
      2. Or even a prototype example...
3. Set up TDW on-call rotation during business hours
   1. Point members at this document's [TDW on-call](#magma-tdw-on-call)
      1. Ensure they read it!
   2. Have all TDW on-call set up pagerduty accounts
      1. And tie them to the pagingn service at `tdw@magma.pagerduty.com`
4. Set up a video conference meeting for the TDW on-call
   1. And cross post to Slack [#project-tech-debt-week-apr-2021](https://magmacore.slack.com/archives/C0202TW1VRN)
5. Amend this process document to incorporate improvmenets
   1. From lessons learned during the TDW
