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
import type {PowerSearchEquipmentResultsTable_equipment} from './__generated__/PowerSearchEquipmentResultsTable_equipment.graphql';
import type {Theme} from '@material-ui/core';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import EquipmentBreadcrumbs from '../equipment/EquipmentBreadcrumbs';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {AutoSizer, Column, Table} from 'react-virtualized';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {capitalize} from '@fbcnms/util/strings';
import {createFragmentContainer, graphql} from 'react-relay';
import {lowerCase} from 'lodash';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

import 'react-virtualized/styles.css';

const styles = (theme: Theme) => ({
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
  },
  cell: {
    padding: '16px 17px',
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
  headerText: {
    fontSize: '12px',
    lineHeight: '16px',
    color: 'rgba(0, 0, 0, 0.54)',
    textTransform: 'none',
  },
  cellText: {
    fontSize: '13px',
    lineHeight: '16px',
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
});

type Props = WithAlert &
  WithStyles<typeof styles> &
  ContextRouter & {
    equipment: PowerSearchEquipmentResultsTable_equipment,
    selectedEquipment?: ?Equipment,
    onEquipmentSelected?: (equipment: Equipment) => void,
    onWorkOrderSelected?: (workOrderId: string) => void,
    onRowSelected?: (equipment: Equipment) => void,
  };

class PowerSearchEquipmentResultsTable extends React.Component<Props> {
  static contextType = AppContext;
  context: AppContextType;

  _headerRenderer = ({label}) => {
    const {classes} = this.props;
    return (
      <div className={classes.cell}>
        <Text className={classes.headerText}>{label}</Text>
      </div>
    );
  };

  _cellRenderer = ({dataKey, rowData, cellData}) => {
    const {
      classes,
      history,
      onEquipmentSelected,
      onWorkOrderSelected,
      onRowSelected,
    } = this.props;
    let content = null;

    if (dataKey === 'name') {
      if (onEquipmentSelected) {
        content = (
          <Button variant="text" onClick={() => onEquipmentSelected(rowData)}>
            {cellData}
          </Button>
        );
      } else {
        content = (
          <Text color="primary" variant="body2">
            {cellData}
          </Text>
        );
      }
    } else if (dataKey === 'status' && rowData.futureState) {
      if (onWorkOrderSelected) {
        content = (
          <Button
            variant="text"
            onClick={() => onWorkOrderSelected(rowData.workOrder.id)}>
            {cellData}
          </Button>
        );
      }
    } else if (dataKey === 'location') {
      content = (
        <EquipmentBreadcrumbs
          equipment={rowData}
          showSelfEquipment={false}
          onParentLocationClicked={
            onRowSelected
              ? null
              : (locationId: string) =>
                  history.push(InventoryAPIUrls.location(locationId))
          }
          onEquipmentClicked={
            onRowSelected
              ? null
              : (equipmentId: string) =>
                  history.push(InventoryAPIUrls.equipment(equipmentId))
          }
          size="small"
        />
      );
    } else {
      content = (
        <Text className={classes.cellText} variant="body2">
          {cellData}
        </Text>
      );
    }

    return <div className={classes.cell}>{content}</div>;
  };

  render() {
    const {classes, equipment, onRowSelected, selectedEquipment} = this.props;
    if (equipment.length === 0) {
      return null;
    }
    const equipmetStatusEnabled = this.context.isFeatureEnabled(
      'planned_equipment',
    );
    const externalIDEnabled = this.context.isFeatureEnabled('external_id');

    return equipment.length > 0 ? (
      <AutoSizer>
        {(height: number, width: number) => (
          <Table
            className={classes.table}
            height={height}
            width={width}
            headerHeight={50}
            rowHeight={50}
            rowCount={equipment.length}
            rowGetter={(index: number) => equipment[index]}
            gridClassName={classes.table}
            rowClassName={({index}) =>
              classNames({
                [classes.header]: index === -1,
                [classes.row]: index !== -1,
                [classes.clickableRow]: onRowSelected != null,
                [classes.checked]:
                  selectedEquipment &&
                  index !== -1 &&
                  equipment[index].id === selectedEquipment.id,
              })
            }
            onRowClick={({_event, _index, rowData}) =>
              onRowSelected && onRowSelected(rowData)
            }>
            <Column
              label="Name"
              dataKey="name"
              width={250}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
            />
            {externalIDEnabled && (
              <Column
                label="External ID"
                dataKey="id"
                width={150}
                flexGrow={1}
                headerRenderer={this._headerRenderer}
                cellRenderer={this._cellRenderer}
                cellDataGetter={({rowData}) => rowData.externalId}
              />
            )}
            <Column
              label="Location"
              dataKey="location"
              width={250}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
            />
            {equipmetStatusEnabled && (
              <Column
                label="Status"
                dataKey="status"
                cellDataGetter={({rowData}) =>
                  rowData.futureState && rowData.workOrder
                    ? `${capitalize(
                        lowerCase(rowData.workOrder.status),
                      )} ${lowerCase(rowData.futureState)}`
                    : 'Installed'
                }
                width={250}
                flexGrow={1}
                headerRenderer={this._headerRenderer}
                cellRenderer={this._cellRenderer}
              />
            )}
            <Column
              label="Type"
              dataKey="type"
              cellDataGetter={({rowData}) => rowData.equipmentType.name}
              width={250}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
            />
          </Table>
        )}
      </AutoSizer>
    ) : null;
  }
}

export default withRouter(
  withAlert(
    withStyles(styles)(
      createFragmentContainer(PowerSearchEquipmentResultsTable, {
        equipment: graphql`
          fragment PowerSearchEquipmentResultsTable_equipment on Equipment
            @relay(plural: true) {
            id
            name
            futureState
            externalId
            equipmentType {
              id
              name
            }
            workOrder {
              id
              status
            }
            ...EquipmentBreadcrumbs_equipment
          }
        `,
      }),
    ),
  ),
);
