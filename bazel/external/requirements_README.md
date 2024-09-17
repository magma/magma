# Python Dependencies managed by Bazel

Requirements.txt holds all Python dependencies which are required by Python-based modules in Magma and have to be built via Bazel. All entries are loadable and, thus part of the Bazel environment if necessary. However, an entry is only loaded via Bazel if a Bazel target applies it as a dependency.

## How to update Python dependencies

 1. Add, remove or modify dependencies in requirements.in

 2. Generate a new version of requirements.txt, including required hashes

       `cd $MAGMA/bazel/external`

       `pip-compile --allow-unsafe --generate-hashes --output-file=requirements.txt requirements.in`

 The changes are then automatically included in the next Bazel build process.

## How to upgrade Python dependencies

In general, existing and thus pinned dependencies in the requirement.txt are not upgraded when the command above is executed. In order to upgrade all dependencies (to the highest possible version) use the following command:

`pip-compile --upgrade --allow-unsafe --generate-hashes --output-file=requirements.txt requirements.in`
