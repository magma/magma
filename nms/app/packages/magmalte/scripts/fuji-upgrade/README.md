# Magma Pre 1.5 Upgrade NMS DB Data Migration

Magma versions 1.4 and prior use a separate database for the NMS and the orc8r.
This provides a process to combine these two databases together, a necessary
process before completing the rest of the upgrade from 1.4 to 1.5.

## Usage

Run on a machine with access to the kubernetes cluster you are using to run Magma.

*Paste in shell prompt:*

`wget https://raw.githubusercontent.com/magma/magma/master/nms/app/packages/magmalte/scripts/fuji-upgrade/pre-upgrade-migration.sh && chmod +x pre-upgrade-migration.sh && ./pre-upgrade-migration.sh`
