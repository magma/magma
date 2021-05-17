---
id: build_install_magma_pkg_in_agw
title: Build and install a magma package in AGW
hide_title: true
---
# Build and install a magma package in AGW

**Description:** Purpose of this document is describe how to build and install a magma package in AGW.

**Environment:** AGW

**Steps:**

1. **Clone or Update  your local repository**. In any machine (i.e. your computer) clone the magma repo  and checkout the version from where you would like to build your package. For example, for v1.1 you can run:

```
git clone https://github.com/magma/magma.git
git checkout v1.1
```

2. **Install prerequisites**. Make sure you have installed all the tools specified in the prerequisites https://magma.github.io/magma/docs/basics/prerequisites#prerequisites

3. **Build and create deb package**. In your local magma repo, go to the path `magma/lte/gateway` and run the command `fab dev package:vcs=git`

This command will create a vagrant magma machine, then build and create a deb package.

4. **Locate the packages**. Once the above command finished. You need to enter the VM to verify the deb packages are there.

```
vagrant ssh magma
cd ~/magma-packages/
```
You will need only the ones that say `magma_1.1.XXX` and `magma-sctpd_1.1.XXX` (for v1.1 versions)

5. **Download the package**. You can download the files to your computer from the vagrant machine. To do so, you can install a vagrant plugin in your computer and then download the package from the VM to your computer with the following commands:

```
vagrant plugin install vagrant-scp
vagrant scp magma: ~/magma-packages/<deb_package>
```

6. **Upload the package to AGW** that you would like to install.

7. **Install the package**. In order to install the new deb package in AGW, you can run

`sudo apt -f install <magma package>`

8. **Restart the magma services**
```
sudo service magma@* stop
sudo service magma@magmad restart
```
9. You can **verify the installed version** with

`apt show magma`
