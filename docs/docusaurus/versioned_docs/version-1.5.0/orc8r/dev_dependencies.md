---
id: version-1.5.0-dev_dependencies
title: Module Dependencies
hide_title: true
original_id: dev_dependencies
---

# Module Dependencies

Within Magma, each module (orc8r, lte, cwf, feg, wifi, fbinternal)
is in charge of generating and defining their own API objects, specifically
defined as Swagger objects. A module's Swagger objects may depend on those
from other modules. This results in the following module dependency tree.

![Module Dependencies](assets/orc8r/gen_construct_dependencies.png)
