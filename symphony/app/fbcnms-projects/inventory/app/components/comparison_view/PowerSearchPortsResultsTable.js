/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type PowerSearchPortsResultsTable_ports from './__generated__/PowerSearchPortsResultsTable_ports.graphql';
import type {ContextRouter} from 'react-router-dom';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import Box from '@material-ui/core/Box';
import EquipmentBreadcrumbs from '../equipment/EquipmentBreadcrumbs';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import {AutoSizer, Column, Table} from 'react-virtualized';
import {createFragmentContainer, graphql} from 'react-relay';
import {getPropertyValue} from '../../common/Property';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

import 'react-virtualized/styles.css';

const styles = () => ({
  table: {
    outline: 'none',
  },
  header: {
    borderBottom: '2px solid #f0f0f0',
  },
  cell: {
    padding: '16px 17px',
    lineHeight: '100%',
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
  linkText: {
    color: 'inherit',
    display: 'inline',
    fontWeight: 'bold',
  },
  propsCell: {
    padding: '16px 17px',
    display: 'block',
    overflow: 'auto',
    lineHeight: '100%',
  },
});

type Props = WithAlert &
  WithStyles<typeof styles> &
  ContextRouter & {
    ports: PowerSearchPortsResultsTable_ports,
  };

class PowerSearchPortsResultsTable extends React.Component<Props> {
  _getConnectedPort = (
    port: PowerSearchPortsResultsTable_ports,
  ): ?PowerSearchPortsResultsTable_ports => {
    if (!port || !port.link || port.link.ports.length < 2) {
      return null;
    }
    if (port.link.ports[0].id != port.id) {
      return port.link.ports[0];
    }
    if (port.link.ports[1].id != port.id) {
      return port.link.ports[1];
    }
  };

  _headerRenderer = ({label}) => {
    const {classes} = this.props;
    return (
      <div className={classes.cell}>
        <Text className={classes.headerText}>{label}</Text>
      </div>
    );
  };

  _cellRenderer = ({dataKey, _rowData, cellData}) => {
    const {classes, history} = this.props;
    let content = null;
    if (cellData == null) {
      content = null;
    } else if (dataKey === 'portName' || dataKey === 'equipmentType') {
      content = (
        <Text className={classNames(classes.cellText)}>{cellData}</Text>
      );
    } else if (dataKey === 'equipment') {
      content = (
        <EquipmentBreadcrumbs
          equipment={cellData}
          showSelfEquipment={true}
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
    } else if (dataKey === 'properties' && cellData.length != 0) {
      return (
        <div className={classes.propsCell}>
          {cellData.map(property => {
            const {name} = property.propertyType;
            const val = getPropertyValue(property) ?? '';
            return <Box>{`${name}: ${val}`}</Box>;
          })}
        </div>
      );
    } else {
      content = <Text className={classes.cellText}>{cellData}</Text>;
    }

    return <div className={classes.cell}>{content}</div>;
  };

  _getRowHeight = rowData => {
    return rowData.link?.properties.length > 3
      ? 40 + rowData.link?.properties.length * 10
      : 50;
  };

  render() {
    const {classes, ports} = this.props;
    if (ports.length === 0) {
      return null;
    }
    return ports.length > 0 ? (
      <AutoSizer>
        {({height, width}) => (
          <Table
            className={classes.table}
            height={height}
            width={width}
            headerHeight={50}
            rowHeight={({index}) => this._getRowHeight(ports[index])}
            rowCount={ports.length}
            rowGetter={({index}) => ports[index]}
            gridClassName={classes.table}
            rowClassName={({index}) => (index === -1 ? classes.header : '')}>
            <Column
              label="Equipment"
              dataKey="equipment"
              width={350}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => rowData.parentEquipment}
            />
            <Column
              label="Equipment Type"
              dataKey="equipmentType"
              width={150}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) =>
                rowData.parentEquipment.equipmentType.name
              }
            />
            <Column
              label="Port Name"
              dataKey="portName"
              width={150}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => rowData.definition.name}
            />
            <Column
              label="Properties"
              dataKey="properties"
              width={350}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => rowData.properties}
            />
            <Column
              label="Linked Equipment"
              dataKey="equipment"
              width={350}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => {
                const connectedPort = this._getConnectedPort(rowData);
                return connectedPort?.parentEquipment ?? null;
              }}
            />
            <Column
              label="Linked Equipment Type"
              dataKey="equipmentType"
              width={250}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => {
                const connectedPort = this._getConnectedPort(rowData);
                return connectedPort
                  ? connectedPort.parentEquipment.equipmentType.name
                  : null;
              }}
            />
            <Column
              label="Linked Port Name"
              dataKey="portName"
              width={150}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => {
                const connectedPort = this._getConnectedPort(rowData);
                return connectedPort ? connectedPort.definition.name : null;
              }}
            />
            <Column
              label="Link Properties"
              dataKey="properties"
              width={350}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => rowData.link?.properties}
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
      createFragmentContainer(PowerSearchPortsResultsTable, {
        ports: graphql`
          fragment PowerSearchPortsResultsTable_ports on EquipmentPort
            @relay(plural: true) {
            id
            definition {
              id
              name
            }
            link {
              id
              ports {
                id
                definition {
                  id
                  name
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
                    }
                  }
                  ...EquipmentBreadcrumbs_equipment
                }
              }
              properties {
                id
                stringValue
                intValue
                floatValue
                booleanValue
                latitudeValue
                longitudeValue
                rangeFromValue
                rangeToValue
                propertyType {
                  id
                  name
                  type
                  isEditable
                  isInstanceProperty
                  stringValue
                  intValue
                  floatValue
                  booleanValue
                  latitudeValue
                  longitudeValue
                  rangeFromValue
                  rangeToValue
                }
              }
            }
            parentEquipment {
              id
              name
              equipmentType {
                id
                name
              }
              ...EquipmentBreadcrumbs_equipment
            }
            properties {
              id
              stringValue
              intValue
              floatValue
              booleanValue
              latitudeValue
              longitudeValue
              rangeFromValue
              rangeToValue
              propertyType {
                id
                name
                type
                isEditable
                isInstanceProperty
                stringValue
                intValue
                floatValue
                booleanValue
                latitudeValue
                longitudeValue
                rangeFromValue
                rangeToValue
              }
            }
          }
        `,
      }),
    ),
  ),
);
