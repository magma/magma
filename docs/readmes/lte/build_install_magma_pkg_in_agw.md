---
id: build_install_magma_pkg_in_agw
title: Build and install a magma package in AGW
hide_title: true
---
# Build and install a magma package in AGW

## Description

The purpose of this document is to describe how to build and install a magma package in AGW.

## Environment

AGW

## Steps

Note that the following is supported starting with magma v1.9. For older versions see respective documentation.

1. **Clone or Update your local repository**.
    In any machine (i.e. your computer) clone the magma repo (and switch to the branch you want to build).

    ```bash
    git clone https://github.com/magma/magma.git
    ```

2. **Install prerequisites**.
    Make sure you have installed all the tools specified in the [prerequisites](https://magma.github.io/magma/docs/basics/prerequisites#prerequisites)

3. **Build and create deb package**.
    To build an AGW package spin up a vagrant machine and then build and create a deb package.

    From `$MAGMA_ROOT/lte/gateway` on your host machine run:

    ```bash
    vagrant up magma
    vagrant ssh magma
    ```

    In the VM from `$MAGMA_ROOT` run:

    ```bash
    bazel run //lte/gateway/release:release_build --config=production
    ```

    To create a package for development or testing, run

    ```bash
    bazel run //lte/gateway/release:release_build
    ```

    Omitting the `--config=production` flag will compile all C++ services with `Debug` compiler flags and enable ASAN. This is recommended for testing only as it will impact performance. In contrast, the production package has C++ services built with `RelWithDebInfo` compiler flags.

4. **Locate the packages**.
    Once the above command finished you can find the packages inside the VM:

    ```bash
    cd /tmp/packages
    ```

    There should be two packages named `magma_1.9.XXX` and `magma-sctpd_1.9.XXX` (for v1.9 versions).

5. **Copy the packages to the target machine**.

6. **Install the package**.
    In order to install the new deb package in AGW, you can run

    ```bash
    sudo apt -f install MAGMA_PACKAGE
    ```

7. **Restart the magma services**

    ```bash
    sudo service magma@* stop
    sudo service magma@magmad restart
    ```

8. You can **verify the installed version** with

    ```bash
    apt show magma
    ```
