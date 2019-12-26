/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment} from '../../common/Equipment';
import type {PowerSearchLinkFirstEquipmentResultsTable_equipment} from './__generated__/PowerSearchLinkFirstEquipmentResultsTable_equipment.graphql';

import EquipmentBreadcrumbs from '../equipment/EquipmentBreadcrumbs';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {AutoSizer, Column, Table} from 'react-virtualized';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';

import 'react-virtualized/styles.css';

const useStyles = makeStyles(theme => ({
  root: {
    width: '100%',
    marginTop: theme.spacing(3),
    overflowX: 'auto',
  },
  table: {
    outline: 'none',
  },
  header: {
    borderBottom: '2px solid #f0f0f0',
    margin: '0px',
  },
  cell: {
    padding: '14px 16px',
  },
  addButton: {
    paddingLeft: '16px',
    paddingRight: '16px',
  },
  futureState: {
    textTransform: 'capitalize',
    maxWidth: '50px',
  },
  icon: {
    padding: '0px',
    marginLeft: theme.spacing(),
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
      backgroundColor: symphony.palette.B50,
    },
    '&:focus': {
      outline: 'none',
    },
  },
}));

type Props = {
  equipment: PowerSearchLinkFirstEquipmentResultsTable_equipment,
  selectedEquipment: ?Equipment,
  onEquipmentSelected: (equipment: Equipment) => void,
};

const PowerSearchLinkFirstEquipmentResultsTable = (props: Props) => {
  const classes = useStyles();
  const {equipment, selectedEquipment, onEquipmentSelected} = props;

  const headerRenderer = ({label}) => {
    return (
      <div className={classes.cell}>
        <Text variant="subtitle2">{label}</Text>
      </div>
    );
  };

  const cellRenderer = ({dataKey, rowData, cellData}) => {
    let content = null;

    if (dataKey === 'location') {
      content = (
        <EquipmentBreadcrumbs
          equipment={rowData}
          showSelfEquipment={false}
          variant="body2"
        />
      );
    } else {
      content = (
        <Text variant={dataKey === 'name' ? 'subtitle2' : 'body2'}>
          {cellData}
        </Text>
      );
    }

    return <div className={classes.cell}>{content}</div>;
  };

  const onRowClicked = ({_event, _index, rowData}) => {
    onEquipmentSelected(rowData);
  };

  if (equipment.length === 0) {
    return null;
  }

  return equipment.length > 0 ? (
    <AutoSizer>
      {({height, width}) => (
        <Table
          className={classes.table}
          height={height}
          width={width}
          headerHeight={50}
          rowHeight={50}
          rowCount={equipment.length}
          rowGetter={({index}) => equipment[index]}
          gridClassName={classes.table}
          rowClassName={({index}) =>
            classNames({
              [classes.header]: index === -1,
              [classes.row]: index !== -1,
              [classes.checked]:
                selectedEquipment &&
                index !== -1 &&
                equipment[index].id === selectedEquipment.id,
            })
          }
          onRowClick={onRowClicked}>
          <Column
            label="Equipment A Name"
            dataKey="name"
            width={250}
            flexGrow={1}
            headerRenderer={headerRenderer}
            cellRenderer={cellRenderer}
            headerClassName={classes.column}
            className={classes.column}
          />
          <Column
            label="Equipment Type"
            dataKey="type"
            cellDataGetter={({rowData}) => rowData.equipmentType.name}
            width={250}
            flexGrow={1}
            headerRenderer={headerRenderer}
            cellRenderer={cellRenderer}
            headerClassName={classes.column}
            className={classes.column}
          />
          <Column
            label="Location"
            dataKey="location"
            width={250}
            flexGrow={1}
            headerRenderer={headerRenderer}
            cellRenderer={cellRenderer}
            headerClassName={classes.column}
            className={classes.column}
          />
        </Table>
      )}
    </AutoSizer>
  ) : null;
};

export default createFragmentContainer(
  PowerSearchLinkFirstEquipmentResultsTable,
  {
    equipment: graphql`
      fragment PowerSearchLinkFirstEquipmentResultsTable_equipment on Equipment
        @relay(plural: true) {
        id
        name
        futureState
        equipmentType {
          id
          name
        }
        ...EquipmentBreadcrumbs_equipment
      }
    `,
  },
);
