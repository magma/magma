This directory contains an SBOM for Magmacore.

It was produced using scancode-toolkit (https://scancode-toolkit.readthedocs.io) using the following CLI recipe:

`./scancode --license --copyright --json-pp sbom.json ../../my-fork/magma/nms`

The generated file being over GitHub's 50 MB limit, the json file was then compressed using this recipe:
`tar cfz sbom-json.tgz sbom.json`


