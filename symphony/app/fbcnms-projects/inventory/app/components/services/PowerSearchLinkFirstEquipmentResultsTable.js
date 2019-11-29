/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {ContextRouter} from 'react-router-dom';
import type {Equipment} from '../../common/Equipment';
import type {PowerSearchLinkFirstEquipmentResultsTable_equipment} from './__generated__/PowerSearchLinkFirstEquipmentResultsTable_equipment.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import EquipmentBreadcrumbs from '../equipment/EquipmentBreadcrumbs';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {AutoSizer, Column, Table} from 'react-virtualized';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

import 'react-virtualized/styles.css';

const styles = theme => ({
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
});

type Props = WithAlert &
  WithStyles<typeof styles> &
  ContextRouter & {
    equipment: PowerSearchLinkFirstEquipmentResultsTable_equipment,
    onEquipmentSelected: (equipment: Equipment) => void,
  };

class PowerSearchLinkFirstEquipmentResultsTable extends React.Component<Props> {
  static contextType = AppContext;
  context: AppContextType;

  _headerRenderer = ({label}) => {
    const {classes} = this.props;
    return (
      <div className={classes.cell}>
        <Text variant="subtitle2">{label}</Text>
      </div>
    );
  };

  _cellRenderer = ({dataKey, rowData, cellData}) => {
    const {classes, history, onEquipmentSelected} = this.props;
    let content = null;

    if (dataKey === 'name') {
      content = (
        <Button variant="text" onClick={() => onEquipmentSelected(rowData)}>
          {cellData}
        </Button>
      );
    } else if (dataKey === 'location') {
      content = (
        <EquipmentBreadcrumbs
          equipment={rowData}
          showSelfEquipment={false}
          onParentLocationClicked={locationId =>
            history.push(
              `inventory/` + (locationId ? `?location=${locationId}` : ''),
            )
          }
          onEquipmentClicked={equipmentId =>
            history.push(
              `inventory/` + (equipmentId ? `?equipment=${equipmentId}` : ''),
            )
          }
          variant="body2"
        />
      );
    } else {
      content = <Text variant="body2">{cellData}</Text>;
    }

    return <div className={classes.cell}>{content}</div>;
  };

  render() {
    const {classes, equipment} = this.props;
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
            rowClassName={({index}) => (index === -1 ? classes.header : '')}>
            <Column
              label="Equipment A Name"
              dataKey="name"
              width={250}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              headerClassName={classes.column}
              className={classes.column}
            />
            <Column
              label="Equipment Type"
              dataKey="type"
              cellDataGetter={({rowData}) => rowData.equipmentType.name}
              width={250}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              headerClassName={classes.column}
              className={classes.column}
            />
            <Column
              label="Location"
              dataKey="location"
              width={250}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              headerClassName={classes.column}
              className={classes.column}
            />
          </Table>
        )}
      </AutoSizer>
    ) : null;
  }
}

export default withAlert(
  withStyles(styles)(
    createFragmentContainer(PowerSearchLinkFirstEquipmentResultsTable, {
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
    }),
  ),
);
