## Security

The Magma Core Foundation takes the security of our software products and services seriously, which includes all source code repositories managed through our [Repositories on Github](https://github.com/magma).

If you believe you have found a security vulnerability in any Magma Core Foundation repository please report it to us as described below. If we accept your report, we may add your name to the acknowledgments section of our site. At your option, we may use an alias, link to your home page, or not have any such acknowledgment.

## Risks

The Magma materials are provided in accordance with the licenses made available in the LICENSE file. Prior to using the materials, it is highly recommended that you test and verify that the materials meet your specific requirements, including, without limitation, any and all security and performance requirements.

## Reporting Security Issues

Report vulnerabilities by sending email to [security@lists.magmacore.org](mailto:security@lists.magmacore.org). _Please **do not report security vulnerabilities using Github Issues**._  While the Magma Core Foundation is committed to operating openly and transparently, reporting security vulnerabilities publicly could place users of the Foundation's software at risk of having their systems exploited.

In order to help us understand the scope of the issue and to triage your report more quickly please include the information listed below (as much as you can provide):

- Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
- Step-by-step instructions to reproduce the issue
- Any special configuration required to reproduce the issue
- Impact of the issue, including how an attacker might exploit the issue
- Proof-of-concept or exploit code (if possible)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL).

We prefer all communications to be in English.

## Response Process

Certain members of the Magma team have been designated as the "vulnerability management team" and will receive any e-mails sent to `security@lists.magmacore.org`. When receiving such an e-mail, they will:

1. Reply to the e-mail acknowledging its receipt, cc'ing `security@lists.magmacore.org` so that the other members of the team are aware that they are handling the issue.  

   - If you are a submitter and do not receive a response within 72 hours, please follow up via email to ensure we received your original message.

2. Inform the #security channel in [the Magma Slack](https://magmacore.slack.com)

2. Evaluate severity using the [Scoring Rubric Template](https://github.com/magma/security/blob/main/Scoring%20rubric%20for%20Magma%20weaknesses.xlsx).

3. Create [an issue in the security repo](https://github.com/magma/security/issues).

4. Research engineering solutions, possibly in consultation with codeowners.

5. Decide whether the report calls for a formal advisory in addition to an issue.

   1. If it does not call for an advisory, inform the discloser, taking care not to reveal potential weaknesses, then prioritize and fix the issue according to ordinary workflows.

   2. If the report does merit an advisory, the vulnerability management team will execute the Advisory Procedure below.
  
6. Add the discloser to the security acknowledgements section in README.md, if the discloser wishes.

### Advisory Procedure

1. Assign a priority as follows:

   - Critical Risk: shall be treated as "P0" issue, with a fix developed and released as soon as possible even if it requires an out-or-sequence release.
   - High Risk: shall be treated as "P1" issue, with a fix delivered no later than the next scheduled "intermediate" release.
   - Medium Risk: shall be treated a "P1" issue, with a fix delivered no later than the next scheduled "intermediate" release.
   - Low Risk: shall be treated as a "P2" issue, with a fix delivered no later than the next scheduled "major" release

2. Use GitHub to create a [new security advisory](https://github.com/magma/magma/security/advisories/new). For more information on security advisories, see [the GitHub documentation](https://docs.github.com/en/free-pro-team@latest/github/managing-security-vulnerabilities/managing-security-vulnerabilities-in-your-project).

   - One must be one of the repo admins to do this. Vulnerability management team members who are not also a repo admin will reach out to the repo admins until they find one who can create the advisory.

3. [Add the reporter](https://docs.github.com/en/free-pro-team@latest/github/managing-security-vulnerabilities/adding-a-collaborator-to-a-security-advisory) to the security advisory so that they can get updates.

4. Assign an engineering resource to the ticket to implement the fix. Identify and inform the relevant Magma Codeowners, adding them to the security advisory.

5. As the fix is being developed, reach out to the reporter to ask them if they would like to be involved and whether they would like to be credited. For credit, the GitHub security advisory UI has a field that allows contributors to be credited.

6. When the issue is resolved, contact the Magma release team and Magma Core Foundation's Outreach Committee to coordinate the publication of the security advisory.

## Security Vulnerability Team Members

The Security Vulnerability team members shall include the members of the TSC, current Codeowners of the Magma Project as approved by the TSC, and additional contributors responsible for vulnerability management.

## Policy

The Magma Core Foundation follows the principle of [The CERT Guide to Coordinated Vulnerability Disclosure](https://resources.sei.cmu.edu/asset_files/SpecialReport/2017_003_001_503340.pdf)
