## Security

The Magma Core Foundation takes the security of our software products and services seriously, which includes all source code repositories managed through our [Repositories on Github](https://github.com/magma).

If you believe you have found a security vulnerability in any Magma Core Foundation repository please report it to us as described below.

## Reporting Security Issues

`Please do not report security vulnerabilities using Github Issues.  While the Magma Core Foundation is committed to operating openly and transparently reporting security vulnerabilities publicly could place users of the Foundation's software at risk of having their systems exploited.`

Instead, please report them via email to security@magmacore.org.

In order to help us understand the scope of the issue and to triage your report more quickly please include the information listed below (as much as you can provide):

  * Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
  * Step-by-step instructions to reproduce the issue
  * Any special configuration required to reproduce the issue
  * Impact of the issue, including how an attacker might exploit the issue
  * Proof-of-concept or exploit code (if possible)
  * Full paths of source file(s) related to the manifestation of the issue
  * The location of the affected source code (tag/branch/commit or direct URL).

## Response Process

Certain members of the Magma team have been designated as the "vulnerability management team" and will receive any e-mails sent to `security@magmacore.org`. When receiving such an e-mail, they will:

1. Reply to the e-mail acknowledging its receipt, cc'ing `security@magmacore.org` so that the other members of the team are aware that they are handling the issue.  You should receive this initial response within 72 hours, usually sooner. If for some reason you do not, please follow up via email to ensure we received your original message.

2. Security issues will be evaluated by the Security Vulnerability team and assigned severity level using the Common Vulnerability Scoring System as defined by [First.org](https://first.org/cvss/user-guide).  Development of remediation and fixes shall be prioritized as follows:

   - Critical Risk: shall be treated as "P0" issue, with a fix developed and released as soon as possible even if it requires an out-or-sequence release.
   - High Risk: shall be treated as "P1" issue, with a fix delivered no later than the next scheduled "intermediate" release.
   - Medium Risk: shall be treated a "P1" issue, with a fix delivered no later than the next scheduled "intermediate" release.
   - Low Risk: shall be treated as a "P2" issue, with a fix delivered no later than the next scheduled "major" release

3. Create a new [security advisory](https://github.com/magma/magma/security/advisories/new).
   One must be one of the repo admins to do this. Vulnerability management team members who are not
   also a repo admin will reach out to the repo admins until they find one who can create the advisory.
   The repo admins who are also vulnerability management team members are @tbd.
4. [Add the reporter](https://docs.github.com/en/free-pro-team@latest/github/managing-security-vulnerabilities/adding-a-collaborator-to-a-security-advisory)
   to the security advisory so that they can get updates.
5. Inform the relevant Magma Codeowners, adding them to the security advisory.

As the fix is being developed, they will then reach out to the reporter to ask them if they would like to be involved and whether they would like to be credited. For credit, the GitHub security advisory UI has a field that allows contributors to be credited.

When the issue is resolved, they will contact the Magma release team and Magma Core Foundation's Outreach Committee to coordinate the publication of the security advisory.


For more information on security advisories, see [the GitHub documentation](https://docs.github.com/en/free-pro-team@latest/github/managing-security-vulnerabilities/managing-security-vulnerabilities-in-your-project).

## Security Vulnerabiltiy Team Members

The Security Vulneratily team members shall include the members of the TSC and other current Codeowners of the Magma Project as approved by the TSC.

## Preferred Languages

We prefer all communications to be in English.

## Policy

The Magma Core Foundation follows the principle of [The CERT Guide to Coordinated Vulnerability Disclosure](https://resources.sei.cmu.edu/asset_files/SpecialReport/2017_003_001_503340.pdf)
