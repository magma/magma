/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment, EquipmentPort} from '../../common/Equipment';
import type {EquipmentPortWithBreadcrumbs} from './EquipmentPortBreadcrumbs';
import type {WithStyles} from '@material-ui/core';

import EquipmentPortBreadcrumbs from './EquipmentPortBreadcrumbs';
import Link from '@fbcnms/ui/components/Link';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import {
  getInitialPortFromDefinition,
  getNonInstancePortsDefinitions,
} from '../../common/Equipment';
import {isTestEnv} from '@fbcnms/ui/config/RuntimeConfig';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  noResultsRoot: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: '100px',
  },
  noResultsLabel: {
    color: theme.palette.grey[600],
  },
  futureState: {
    textTransform: 'capitalize',
    maxWidth: '50px',
  },
});

type Props = {
  equipment: Equipment,
  sourcePortId: string,
  onPortClicked: (port: EquipmentPort) => void,
} & WithStyles<typeof styles>;

const findNestedAvailablePorts = (
  equipment: ?Equipment,
): Array<EquipmentPortWithBreadcrumbs> => {
  if (!equipment) {
    return [];
  }

  const portsToDisplay = [
    ...equipment.ports,
    ...getNonInstancePortsDefinitions(
      equipment.ports,
      equipment.equipmentType.portDefinitions,
    ).map(portDef =>
      getInitialPortFromDefinition(nullthrows(equipment), portDef),
    ),
  ];
  const directPorts = portsToDisplay
    .map(port => ({
      ...port,
      breadcrumbs: [equipment],
    }))
    .filter(port => !port.link)
    .sort((portA, portB) =>
      sortLexicographically(portA.definition.name, portB.definition.name),
    );
  const nestedPorts = equipment.positions
    .map(position =>
      findNestedAvailablePorts(position.attachedEquipment).map(
        portWithBreadcrumbs => ({
          ...portWithBreadcrumbs,
          breadcrumbs: [equipment, ...portWithBreadcrumbs.breadcrumbs],
        }),
      ),
    )
    // $FlowFixMe: https://github.com/facebook/flow/pull/6948
    .flat();
  return [...directPorts, ...nestedPorts];
};

const AvailablePortsTable = (props: Props) => {
  const {equipment, sourcePortId, onPortClicked, classes} = props;
  const ports = findNestedAvailablePorts(equipment).filter(
    port => port.id != sourcePortId,
  );
  if (ports.length === 0) {
    return (
      <div className={classes.noResultsRoot}>
        <Text variant="h6" className={classes.noResultsLabel}>
          No available ports on
          {` ${equipment.equipmentType.name} ${equipment.name}`}
        </Text>
      </div>
    );
  }

  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>Port Name</TableCell>
          <TableCell>Parent Equipment</TableCell>
          <TableCell>Visible Label</TableCell>
          <TableCell>Type</TableCell>
          {isTestEnv() && <TableCell>Status</TableCell>}
        </TableRow>
      </TableHead>
      <TableBody>
        {ports.map(port => {
          const futureState = port.link?.futureState;
          return (
            <TableRow key={`port_${port.id}`}>
              <TableCell>
                <Link onClick={() => onPortClicked(port)}>
                  {port.definition.name}
                </Link>
              </TableCell>
              <TableCell>
                <EquipmentPortBreadcrumbs port={port} />
              </TableCell>
              <TableCell>{port.definition.visibleLabel}</TableCell>
              {isTestEnv() && (
                <TableCell>
                  <Text variant="caption" className={classes.futureState}>
                    {futureState ? `Planned ${futureState.toLowerCase()}` : ''}
                  </Text>
                </TableCell>
              )}
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
};

export default withStyles(styles)(AvailablePortsTable);
