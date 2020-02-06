/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment, EquipmentPort} from '../common/Equipment';
import type {WithStyles} from '@material-ui/core';

import AvailablePortsTable_ports from './__generated__/AvailablePortsTable_ports.graphql';
import EquipmentBreadcrumbs from './equipment/EquipmentBreadcrumbs';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {AutoSizer, Column, Table} from 'react-virtualized';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

import 'react-virtualized/styles.css';

const styles = {
  noResultsRoot: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: '100px',
  },
  noResultsLabel: {
    color: symphony.palette.D600,
  },
  checked: {
    backgroundColor: symphony.palette.B50,
  },
  row: {
    '&:hover': {
      backgroundColor: symphony.palette.B50,
    },
    '&:focus': {
      outline: 'none',
    },
  },
  table: {
    outline: 'none',
  },
  cell: {
    padding: '14px 16px',
  },
  header: {
    borderBottom: '2px solid #f0f0f0',
    margin: '0px',
  },
  column: {
    '&&': {
      margin: '0px',
      textTransform: 'none',
    },
  },
  checked: {
    backgroundColor: symphony.palette.B50,
  },
  row: {
    '&:hover': {
      backgroundColor: symphony.palette.background,
    },
    '&:focus': {
      outline: 'none',
    },
  },
  clickableRow: {
    cursor: 'pointer',
  },
};

type Props = {
  equipment: Equipment,
  ports: AvailablePortsTable_ports,
  selectedPort: ?EquipmentPort,
  onPortSelected?: (port: EquipmentPort) => void,
} & WithStyles<typeof styles>;

const AvailablePortsTable = (props: Props) => {
  const {equipment, ports, selectedPort, onPortSelected, classes} = props;

  const headerRenderer = ({label}) => {
    return (
      <div className={classes.cell}>
        <Text variant="subtitle2">{label}</Text>
      </div>
    );
  };

  const cellRenderer = ({dataKey, _, cellData}) => {
    let content = null;

    if (dataKey.startsWith('parent_equipment')) {
      content = (
        <EquipmentBreadcrumbs
          equipment={cellData}
          size="small"
          variant="body2"
        />
      );
    } else {
      content = (
        <Text variant={dataKey === 'port_name' ? 'subtitle2' : 'body2'}>
          {cellData}
        </Text>
      );
    }
    return <div className={classes.cell}>{content}</div>;
  };

  const onRowClicked = ({_event, _index, rowData}) => {
    onPortSelected && onPortSelected(rowData);
  };

  if (ports.length === 0) {
    return (
      <div className={classes.noResultsRoot}>
        <Text variant="h6" className={classes.noResultsLabel}>
          {`${fbt(
            'No ports for ' +
              fbt.param('equipment type name', equipment.equipmentType.name) +
              ' ' +
              fbt.param('equipment name', equipment.name),
            'Message when no ports found are for a chosen equipment',
          )}
          `}
        </Text>
      </div>
    );
  }

  return (
    <AutoSizer>
      {({height, width}) => (
        <Table
          className={classes.table}
          height={height}
          width={width}
          headerHeight={50}
          rowHeight={50}
          rowCount={ports.length}
          rowGetter={({index}) => ports[index]}
          gridClassName={classes.table}
          rowClassName={({index}) =>
            classNames({
              [classes.header]: index === -1,
              [classes.row]: index !== -1,
              [classes.clickableRow]: onRowClicked != null,
              [classes.checked]:
                selectedPort &&
                index !== -1 &&
                ports[index].id === selectedPort.id,
            })
          }
          onRowClick={onRowClicked}>
          <Column
            label="Port Name"
            dataKey="port_name"
            width={250}
            flexGrow={1}
            cellDataGetter={({rowData}) => rowData.definition.name}
            headerRenderer={headerRenderer}
            cellRenderer={cellRenderer}
            headerClassName={classes.column}
            className={classes.column}
          />
          <Column
            label="Port Type"
            dataKey="port_type"
            width={250}
            flexGrow={1}
            cellDataGetter={({rowData}) => rowData.definition.portType?.name}
            headerRenderer={headerRenderer}
            cellRenderer={cellRenderer}
            headerClassName={classes.column}
            className={classes.column}
          />
          <Column
            label="Parent Equipment"
            dataKey="parent_equipment"
            width={250}
            flexGrow={1}
            cellDataGetter={({rowData}) => rowData.parentEquipment}
            headerRenderer={headerRenderer}
            cellRenderer={cellRenderer}
            headerClassName={classes.column}
            className={classes.column}
          />
          <Column
            label="Visible Label"
            dataKey="visible_label"
            width={250}
            flexGrow={1}
            cellDataGetter={({rowData}) => rowData.definition.visibleLabel}
            headerRenderer={headerRenderer}
            cellRenderer={cellRenderer}
            headerClassName={classes.column}
            className={classes.column}
          />
        </Table>
      )}
    </AutoSizer>
  );
};

export default withStyles(styles)(
  createFragmentContainer(AvailablePortsTable, {
    ports: graphql`
      fragment AvailablePortsTable_ports on EquipmentPort @relay(plural: true) {
        id
        parentEquipment {
          id
          name
          ...EquipmentBreadcrumbs_equipment
        }
        definition {
          id
          name
          portType {
            name
          }
          visibleLabel
        }
      }
    `,
  }),
);
