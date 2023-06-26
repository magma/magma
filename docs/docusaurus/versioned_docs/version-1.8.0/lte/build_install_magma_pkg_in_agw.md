---
id: version-1.8.0-build_install_magma_pkg_in_agw
title: Build and install a magma package in AGW
hide_title: true
original_id: build_install_magma_pkg_in_agw
---
# Build and install a magma package in AGW

**Description:** Purpose of this document is describe how to build and install a magma package in AGW.

**Environment:** AGW

**Steps:**

1. **Clone or Update  your local repository**. In any machine (i.e. your computer) clone the magma repo  and checkout the version from where you would like to build your package. For example, for v1.8 you can run:

    ```bash
    git clone https://github.com/magma/magma.git
    git checkout v1.8
    ```

2. **Install prerequisites**. Make sure you have installed all the tools specified in the prerequisites <https://magma.github.io/magma/docs/basics/prerequisites#prerequisites>

3. **Build and create deb package**.
    To build an AGW package, use the script located at `$MAGMA_ROOT/lte/gateway/fabfile.py`. The commands below will create a vagrant machine, then build and create a deb package.

    The following commands are to be run from `$MAGMA_ROOT/lte/gateway` on your host machine.
    To create a package for production, run

    ```bash
    fab release package
    ```

    To create a package for development or testing, run

    ```bash
    fab dev package
    ```

    The `dev` flag will compile all C++ services with `Debug` compiler flags and enable ASAN. This is recommended for testing only as it will impact performance. In contrast, the production package has C++ services built with `RelWithDebInfo` compiler flags.

4. **Locate the packages**. Once the above command finished. You need to enter the VM to verify the deb packages are there.

    ```bash
    vagrant ssh magma
    cd ~/magma-packages/
    ```

    You will need only the ones that say `magma_1.1.XXX` and `magma-sctpd_1.1.XXX` (for v1.1 versions)

5. **Download the package**. You can download the files to your computer from the vagrant machine. To do so, you can install a vagrant plugin in your computer and then download the package from the VM to your computer with the following commands:

    ```bash
    vagrant plugin install vagrant-scp
    vagrant scp magma: ~/magma-packages/<deb_package>
    ```

6. **Upload the package to AGW** that you would like to install.

7. **Install the package**. In order to install the new deb package in AGW, you can run

    ```bash
    sudo apt -f install MAGMA_PACKAGE
    ```

8. **Restart the magma services**

    ```bash
    sudo service magma@* stop
    sudo service magma@magmad restart
    ```

9. You can **verify the installed version** with

    ```bash
    apt show magma
    ```
