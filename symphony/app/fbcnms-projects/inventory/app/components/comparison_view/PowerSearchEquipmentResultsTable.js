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
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import EquipmentBreadcrumbs from '../equipment/EquipmentBreadcrumbs';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {AutoSizer, Column, Table} from 'react-virtualized';
import {capitalize} from '@fbcnms/util/strings';
import {createFragmentContainer, graphql} from 'react-relay';
import {lowerCase} from 'lodash';
import {withRouter} from 'react-router-dom';
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
});

type Props = WithAlert &
  WithStyles<typeof styles> &
  ContextRouter & {
    equipment: PowerSearchEquipmentResultsTable_equipment,
    onEquipmentSelected: (equipment: Equipment) => void,
    onWorkOrderSelected: (workOrderId: string) => void,
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
    const {classes, history} = this.props;
    let content = null;

    if (dataKey === 'name') {
      content = (
        <Button
          variant="text"
          onClick={() => this.props.onEquipmentSelected(rowData)}>
          {cellData}
        </Button>
      );
    } else if (dataKey === 'status' && rowData.futureState) {
      content = (
        <Button
          variant="text"
          onClick={() => this.props.onWorkOrderSelected(rowData.workOrder.id)}>
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
          size="small"
        />
      );
    } else {
      content = <Text className={classes.cellText}>{cellData}</Text>;
    }

    return <div className={classes.cell}>{content}</div>;
  };

  render() {
    const {classes, equipment} = this.props;
    if (equipment.length === 0) {
      return null;
    }
    const equipmetStatusEnabled = this.context.isFeatureEnabled(
      'planned_equipment',
    );
    const externalIDEnabled = this.context.isFeatureEnabled('external_id');

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
