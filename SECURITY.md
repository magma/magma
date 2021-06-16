## Security

The Magma Core Foundations takes the security of our software products and services seriously, which includes all source code repositories managed through our [Repositories on Github](https://github.com/magma).

If you believe you have found a security vulnerability in any Magma Core Foundation repository please report it to us as described below.

## Reporting Security Issues

`Please do not report security vulnerabilities using Github Issues.  While the Magma Core Foundation is committed to operating openly and transparently reporting security vulnerabilities publicly could place users of the Foundation's software at risk of having their systems exploited.`

Instead, please report them via email to security@magmacore.org.  If possible, please encrypt your message using our PGP key (_need to provide a key and describe how to do this..._).

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

0. Reply to the e-mail acknowledging its receipt, cc'ing `security@magmacore.org` so that the other members of the team are aware that they are handling the issue.  You should receive this initial response within 72 hours, usually sooner. If for some reason you do not, please follow up via email to ensure we received your original message.

1. Create a new [security advisory](https://github.com/magma/magma/security/advisories/new).
   One must be one of the repo admins to do this. Vulnerability management team members who are not
   also a repo admin will reach out to the repo admins until they find one who can create the advisory.
   The repo admins who are also vulnerability management team members are @tbd.
2. [Add the reporter](https://docs.github.com/en/free-pro-team@latest/github/managing-security-vulnerabilities/adding-a-collaborator-to-a-security-advisory)
   to the security advisory so that they can get updates.
3. Inform the relevant Magma Codeowners, adding them to the security advisory.

As the fix is being developed, they will then reach out to the reporter to ask them if they would like to be involved and whether they would like to be credited. For credit, the GitHub security advisory UI has a field that allows contributors to be credited.

When the issue is resolved, they will contact the Magma release team and Magma Core Foundation's Outreach Committee to coordinate the publication of the security advisory.

Security issues have the equivalent of a P0 priority level but are
not tracked explicitly in the issue database. This means that we attempt to fix them as quickly as possible.

For more information on security advisories, see [the GitHub documentation](https://docs.github.com/en/free-pro-team@latest/github/managing-security-vulnerabilities/managing-security-vulnerabilities-in-your-project).

## Preferred Languages

We prefer all communications to be in English.

## Policy

The Magma Core Foundation follows the principle of [The CERT Guide to Coordinated Vulnerability Disclosure](https://resources.sei.cmu.edu/asset_files/SpecialReport/2017_003_001_503340.pdf)
