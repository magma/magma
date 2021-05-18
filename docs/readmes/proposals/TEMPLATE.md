---
id: proposal_template
title: Title of Design Doc
hide_title: true
---

[This is a template for Magma's change proposal process, documented
[here](README.md).]

# Proposal: [Title]

Author(s): [Author Name, Co-Author Name]

Last updated: [Date]

Discussion at
[https://github.com/magma/magma/issues/4888](https://github.com/magma/magma/issues/4888).

## Context & scope

[Give the reader a very rough overview of the landscape in which the new system 
is being built.  This isn’t a requirements doc; keep it succinct! Bring the 
reader up to speed. Focus on objective background facts.]

### Goals

[Bulleted points; for example, “ACID compliance” for a database.]

### Non-goals

[Bulleted points; for example, “ACID compliance” for a database.]

## Proposal

[Start with an overview of your design, and then go into details. **Focus on 
the trade-offs that you made in the design** in order to make sure this doc has 
long-term value. Given the context (facts), and goals vs non-goals, this is the 
place to suggest a solution and substantiate why it best satisfies those goals.]

## Alternatives considered

[This section lists alternative designs that might have reasonably achieved 
considerable outcomes.  The focus of each should be on the trade-offs that each 
alternative makes and how those trade-offs led to the decision to select the 
primary design. **This is probably the most important section; it shows very 
explicitly why the selected solution is the best given the goals.** This is 
important, since your reader is likely wondering about one or more alternative 
solutions.]

## Cross-cutting concerns

[These are relatively short-sections where the template forces the author to 
write 1-2 sentences to demonstrate consideration of how the design impacts a 
concern -- and how it is addressed.  Teams should standardize on the set of 
concerns.]

### Compatibility

[A discussion of the change with regard to backward / forward compatibility.]

### Observability and Debug

[A description, of how issues with this design would be observed and debugged
in various stages from development through production.]

### Security & privacy

[A description of the security and/or privacy impact of the change (if any).]

## Open issues (if applicable)

[A discussion of issues relating to this proposal for which the author does not
know the solution. This section may be omitted if there are none.]
