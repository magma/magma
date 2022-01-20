# Python Dependencies managed by Bazel

Requirements.txt holds all Python dependencies which are required by Python-based modules in Magma and have to be built via Bazel. All entries are loadable and, thus part of the Bazel environment if necessary. However, an entry is only loaded via Bazel if a Bazel target applies it as a dependency. 

### How to update Python dependencies:

 1. Insert missing dependencies in requirements.in 

 2. Generate a new version of requirements.txt, including required hashes 
    
       `cd $MAGMA/bazel/external`

       `pip-compile --generate-hashes --output-file=requirements.txt requirements.in`

 The changes are then automatically included in the next Bazel build process.

