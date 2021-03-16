---
id: dev_dependencies
title: Module Dependencies
hide_title: true
---

# Module dependencies on generated constructs

Within Magma, each module(orc8r, lte, cwf, feg, wifi, fbinternal),
is in charge of generating and defining their own constructs. A module's
constructs may have dependencies on those of another module.

![Module Dependencies](assets/orc8r/gen_construct_dependencies.png)
