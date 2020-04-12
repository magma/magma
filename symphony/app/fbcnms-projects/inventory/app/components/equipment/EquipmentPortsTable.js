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
import type {Theme, WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import Box from '@material-ui/core/Box';
import EquipmentBreadcrumbs from './EquipmentBreadcrumbs';
import EquipmentPortBreadcrumbs from './EquipmentPortBreadcrumbs';
import EquipmentPortsTableMenu from './EquipmentPortsTableMenu';
import Link from '@fbcnms/ui/components/Link';
import Paper from '@material-ui/core/Paper';
import React, {useContext, useState} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import nullthrows from '@fbcnms/util/nullthrows';
import {capitalize} from '@fbcnms/util/strings';
import {createFragmentContainer, graphql} from 'react-relay';
import {find, uniqBy} from 'lodash';
import {
  getInitialPortFromDefinition,
  getNonInstancePortsDefinitions,
} from '../../common/Equipment';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {
  getNonInstancePropertyTypes,
  getPropertyValue,
} from '../../common/Property';
import {lowerCase} from 'lodash';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withStyles} from '@material-ui/core/styles';

const styles = (_theme: Theme) => ({
  rowFirstCell: {
    paddingLeft: '24px',
  },
  futureState: {
    textTransform: 'capitalize',
    maxWidth: '50px',
  },
  paper: {
    width: '100%',
    overflowX: 'auto',
    marginBottom: 0,
    height: '100%',
  },
  iconButtons: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
  },
  header: {
    backgroundColor: 'white',
  },
  headerCell: {
    backgroundColor: 'white',
  },
});

type Props = {
  equipment: Equipment,
  workOrderId: ?string,
  onPortEquipmentClicked: (equipmentId: string) => void,
  onParentLocationClicked: (locationId: string) => void,
  onWorkOrderSelected: (workOrderId: string) => void,
  ...$Exact<WithStyles<typeof styles>>,
};

const findNestedPorts = (
  equipment: ?Equipment,
): Array<EquipmentPortWithBreadcrumbs> => {
  if (!equipment) {
    return [];
  }

  const portsToDisplay = [
    ...(equipment.ports ?? []),
    ...getNonInstancePortsDefinitions(
      equipment.ports ?? [],
      equipment.equipmentType.portDefinitions,
    ).map(portDef =>
      getInitialPortFromDefinition(nullthrows(equipment), portDef),
    ),
  ];
  const directPorts = portsToDisplay.map(port => ({
    ...port,
    breadcrumbs: [equipment],
  }));
  const nestedPorts = (equipment.positions ?? [])
    .map(position =>
      findNestedPorts(position.attachedEquipment).map(portWithBreadcrumbs => ({
        ...portWithBreadcrumbs,
        breadcrumbs: [equipment, ...portWithBreadcrumbs.breadcrumbs],
      })),
    )
    .flatMap(i => i);
  return [...directPorts, ...nestedPorts];
};

const getConnectedPort = (port: EquipmentPort): ?EquipmentPort => {
  if (!port.link || port.link.ports.length < 2) {
    return null;
  }
  if (port.link.ports[0].id != port.id) {
    return port.link.ports[0];
  }
  if (port.link.ports[1].id != port.id) {
    return port.link.ports[1];
  }
};

type SortableColumn = 'parent_eq' | 'connected_eq_name';

const sortRows = (orderBy: SortableColumn, sortDirection: 'asc' | 'desc') => (
  portA,
  portB,
) => {
  let sort = null;
  if (orderBy === 'connected_eq_name') {
    const connectedPortA = getConnectedPort(portA);
    const connectedPortB = getConnectedPort(portB);
    if (connectedPortA === null) {
      sort = -1;
    } else if (connectedPortB === null) {
      sort = 1;
    } else {
      sort = sortLexicographically(
        nullthrows(connectedPortA).parentEquipment.name,
        nullthrows(connectedPortB).parentEquipment.name,
      );
    }
  } else {
    sort = sortLexicographically(
      `${portA.breadcrumbs.map(b => b.name).join('-')}-${
        portA.definition.index
      }`,
      `${portB.breadcrumbs.map(b => b.name).join('-')}-${
        portB.definition.index
      }`,
    );
  }

  return sortDirection === 'asc' ? sort : -sort;
};

const EquipmentPortsTable = (props: Props) => {
  const {isFeatureEnabled} = useContext(AppContext);
  const [orderBy, setOrderBy] = useState<SortableColumn>('parent_eq');
  const [sortDirection, setSortDirection] = useState('asc');

  const {
    equipment,
    workOrderId,
    onPortEquipmentClicked,
    onParentLocationClicked,
    classes,
  } = props;
  const ports = findNestedPorts(equipment);
  if (ports.length === 0) {
    return null;
  }
  const servicesEnabled = isFeatureEnabled('services');
  const linkStatusEnabled = isFeatureEnabled('planned_equipment');
  const propNames = uniqBy(
    ports
      .slice()
      .map(port => port.definition.portType?.propertyTypes ?? [])
      .flatMap(i => i),
    'name',
  )
    .map(propertyTypes => propertyTypes.name)
    .sort(sortLexicographically);

  const headCells = [
    {label: 'Parent Equipment', key: 'parent_eq', sortable: true},
    {label: 'Port name', key: 'port_name'},
    {label: 'Visible Label', key: 'visible_label'},
    {label: 'Type', key: 'type'},
    ...propNames.map(name => ({label: name, key: `prop_${name}`})),
    {
      label: 'Connected Equipment Name',
      key: 'connected_eq_name',
      sortable: true,
    },
    {label: 'Connected Equipment Type', key: 'connected_eq_type'},
    {label: 'Connected Port', key: 'connected_port'},
    {label: 'Link Properties', key: 'link_props'},
    servicesEnabled ? {label: 'Services', key: 'services'} : null,
    linkStatusEnabled ? {label: 'Link Status', key: 'link_status'} : null,
    {label: null, key: 'actions'},
  ].filter(Boolean);

  return (
    <>
      <Paper className={classes.paper}>
        <Table stickyHeader size="small">
          <TableHead className={classes.header}>
            <TableRow>
              {headCells.map(cell => (
                <TableCell key={cell.key} className={classes.headerCell}>
                  {cell.label && cell.sortable ? (
                    <TableSortLabel
                      active={cell.key === orderBy}
                      direction={sortDirection}
                      onClick={() => {
                        if (orderBy === cell.key) {
                          setSortDirection(
                            sortDirection === 'asc' ? 'desc' : 'asc',
                          );
                        } else {
                          setSortDirection('asc');
                        }
                        setOrderBy(cell.key);
                      }}>
                      {cell.label}
                    </TableSortLabel>
                  ) : (
                    cell.label
                  )}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {ports
              .slice()
              .sort(sortRows(orderBy, sortDirection))
              .map((port, i) => {
                const connectedPort = getConnectedPort(port);
                const futureState = port.link?.futureState;
                const relevantPropertyTypes = getNonInstancePropertyTypes(
                  port.properties ?? [],
                  port.definition.portType?.propertyTypes ?? [],
                );
                const allProps = [
                  ...(port.properties ?? []),
                  ...relevantPropertyTypes.map(getInitialPropertyFromType),
                ];
                return (
                  <TableRow key={`port_${i}`}>
                    <TableCell>
                      <EquipmentPortBreadcrumbs
                        port={port}
                        onPortEquipmentClicked={onPortEquipmentClicked}
                      />
                    </TableCell>
                    <TableCell
                      className={classes.rowFirstCell}
                      component="th"
                      scope="row">
                      {port.definition.name}
                    </TableCell>
                    <TableCell>{port.definition.visibleLabel}</TableCell>
                    <TableCell>{port.definition.portType?.name}</TableCell>
                    {propNames.map(name => {
                      const prop = find(
                        allProps,
                        prop => prop.propertyType.name == name,
                      );
                      return (
                        <TableCell key={`property_${name}`}>
                          {prop ? getPropertyValue(prop) : ''}
                        </TableCell>
                      );
                    })}
                    <TableCell>
                      {connectedPort ? (
                        <EquipmentBreadcrumbs
                          onEquipmentClicked={onPortEquipmentClicked}
                          onParentLocationClicked={onParentLocationClicked}
                          equipment={connectedPort.parentEquipment}
                          size="small"
                          showTypes={false}
                        />
                      ) : (
                        'None'
                      )}
                    </TableCell>
                    <TableCell>
                      {connectedPort
                        ? connectedPort.parentEquipment.equipmentType.name
                        : ''}
                    </TableCell>
                    <TableCell>
                      {connectedPort ? connectedPort.definition.name : ''}
                    </TableCell>
                    <TableCell>
                      {port.link && port.link.properties && (
                        <>
                          {port.link.properties.map(property => {
                            const {name} = property.propertyType;
                            const val = getPropertyValue(property) ?? '';
                            return <Box>{`${name}: ${val}`}</Box>;
                          })}
                        </>
                      )}
                    </TableCell>
                    {servicesEnabled && (
                      <TableCell>
                        {port.link &&
                          port.link.services.map(service => (
                            <Box>{service.name}</Box>
                          ))}
                        {port.serviceEndpoints.map(endpoint => (
                          <Box>{`${endpoint.service.name}: ${lowerCase(
                            endpoint.type.role,
                          )}`}</Box>
                        ))}
                      </TableCell>
                    )}
                    {linkStatusEnabled && (
                      <TableCell>
                        <Link
                          onClick={() =>
                            props.onWorkOrderSelected(
                              nullthrows(port.link?.workOrder?.id),
                            )
                          }>
                          {futureState
                            ? `${capitalize(
                                lowerCase(port.link?.workOrder?.status),
                              )} ${lowerCase(futureState)}`
                            : ''}
                        </Link>
                      </TableCell>
                    )}
                    <TableCell>
                      <EquipmentPortsTableMenu
                        key={`${port.id}-menu`}
                        port={port}
                        workOrderId={workOrderId}
                      />
                    </TableCell>
                  </TableRow>
                );
              })}
          </TableBody>
        </Table>
      </Paper>
    </>
  );
};

graphql`
  fragment EquipmentPortsTable_positionAttachedEquipment on Equipment {
    id
    name
    ports {
      ...EquipmentPortsTable_port @relay(mask: false)
    }
    equipmentType {
      portDefinitions {
        id
        name
        visibleLabel
        bandwidth
      }
    }
  }
`;

graphql`
  fragment EquipmentPortsTable_position on EquipmentPosition
    @relay(mask: false) {
    attachedEquipment {
      ...EquipmentPortsTable_positionAttachedEquipment @relay(mask: false)
      positions {
        attachedEquipment {
          ...EquipmentPortsTable_positionAttachedEquipment @relay(mask: false)
          positions {
            attachedEquipment {
              ...EquipmentPortsTable_positionAttachedEquipment
                @relay(mask: false)
              positions {
                attachedEquipment {
                  ...EquipmentPortsTable_positionAttachedEquipment
                    @relay(mask: false)
                }
              }
            }
          }
        }
      }
    }
  }
`;

graphql`
  fragment EquipmentPortsTable_portDefinition on EquipmentPortDefinition {
    id
    name
    index
    visibleLabel
    portType {
      id
      name
      propertyTypes {
        ...PropertyTypeFormField_propertyType @relay(mask: false)
      }
      linkPropertyTypes {
        ...PropertyTypeFormField_propertyType @relay(mask: false)
      }
    }
  }
`;

graphql`
  fragment EquipmentPortsTable_port on EquipmentPort {
    id
    definition {
      ...EquipmentPortsTable_portDefinition @relay(mask: false)
    }
    parentEquipment {
      id
      name
      equipmentType {
        id
        name
        portDefinitions {
          id
          name
          visibleLabel
          portType {
            id
            name
          }
          bandwidth
        }
      }
    }
    link {
      ...EquipmentPortsTable_link @relay(mask: false)
    }
    properties {
      ...PropertyFormField_property @relay(mask: false)
    }
    serviceEndpoints {
      definition {
        role
      }
      service {
        name
      }
    }
  }
`;

graphql`
  fragment EquipmentPortsTable_link on Link {
    id
    futureState
    ports {
      ...EquipmentPortsTable_link_port @relay(mask: false)
    }
    workOrder {
      id
      status
    }
    properties {
      ...PropertyFormField_property @relay(mask: false)
    }
    services {
      id
      name
    }
  }
`;

graphql`
  fragment EquipmentPortsTable_link_port on EquipmentPort {
    id
    definition {
      id
      name
      visibleLabel
      portType {
        linkPropertyTypes {
          ...PropertyTypeFormField_propertyType @relay(mask: false)
        }
      }
    }
    parentEquipment {
      id
      name
      futureState
      equipmentType {
        id
        name
        portDefinitions {
          id
          name
          visibleLabel
          bandwidth
          portType {
            id
            name
          }
        }
      }
      ...EquipmentBreadcrumbs_equipment
    }
    serviceEndpoints {
      definition {
        role
      }
      service {
        name
      }
    }
  }
`;

export default withStyles(styles)(
  createFragmentContainer(EquipmentPortsTable, {
    equipment: graphql`
      fragment EquipmentPortsTable_equipment on Equipment {
        id
        name
        equipmentType {
          id
          name
          portDefinitions {
            id
            ...EquipmentPortsTable_portDefinition @relay(mask: false)
          }
        }
        ports {
          ...EquipmentPortsTable_port @relay(mask: false)
        }
        positions {
          ...EquipmentPortsTable_position @relay(mask: false)
        }
      }
    `,
  }),
);
