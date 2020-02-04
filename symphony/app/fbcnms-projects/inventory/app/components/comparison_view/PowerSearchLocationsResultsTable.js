/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type PowerSearchLocationsResultsTable_locations from './__generated__/PowerSearchLocationsResultsTable_locations.graphql';
import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {ContextRouter} from 'react-router-dom';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import Box from '@material-ui/core/Box';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import Button from '@fbcnms/ui/components/design-system/Button';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {AutoSizer, Column, Table} from 'react-virtualized';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {createFragmentContainer, graphql} from 'react-relay';
import {getPropertyValue} from '../../common/Property';
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
    lineHeight: '100%',
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
    locations: PowerSearchLocationsResultsTable_locations,
  };

class PowerSearchLocationsResultsTable extends React.Component<Props> {
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

  _onLocationClickedCallback = locationId => {
    const {history} = this.props;
    history.replace(InventoryAPIUrls.location(locationId));
  };

  _cellRenderer = ({dataKey, rowData, cellData}) => {
    const {classes} = this.props;
    let content = null;

    if (dataKey === 'name') {
      content = (
        <Button
          variant="text"
          onClick={() => this._onLocationClickedCallback(rowData.id)}>
          {cellData}
        </Button>
      );
    } else if (dataKey === 'breadcrumbs') {
      const breadcrumbs = cellData.map(l => ({
        id: l.id,
        name: l.name,
        subtext: l.locationType.name,
        onClick: () => this._onLocationClickedCallback(l.id),
      }));
      content = <Breadcrumbs breadcrumbs={breadcrumbs} size={'small'} />;
    } else if (dataKey === 'properties') {
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
    return rowData.properties.length > 3
      ? 40 + rowData.properties.length * 10
      : 50;
  };

  render() {
    const {classes, locations} = this.props;
    if (locations.length === 0) {
      return null;
    }
    const externalIDEnabled = this.context.isFeatureEnabled('external_id');

    return locations.length > 0 ? (
      <AutoSizer>
        {({height, width}) => (
          <Table
            className={classes.table}
            height={height}
            width={width}
            headerHeight={50}
            rowHeight={({index}) => this._getRowHeight(locations[index])}
            rowCount={locations.length}
            rowGetter={({index}) => locations[index]}
            gridClassName={classes.table}
            rowClassName={({index}) => (index === -1 ? classes.header : '')}>
            <Column
              label="Location Name"
              dataKey="name"
              width={150}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => rowData.name}
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
              label="Location Type"
              dataKey="type"
              width={150}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => rowData.locationType.name}
            />
            <Column
              label="Location Ancestors"
              dataKey="breadcrumbs"
              width={350}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => rowData.locationHierarchy}
            />
            <Column
              label="Properties"
              dataKey="properties"
              width={500}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
              cellDataGetter={({rowData}) => rowData.properties}
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
      createFragmentContainer(PowerSearchLocationsResultsTable, {
        locations: graphql`
          fragment PowerSearchLocationsResultsTable_locations on Location
            @relay(plural: true) {
            id
            name
            externalId
            locationType {
              id
              name
              propertyTypes {
                id
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
              equipmentValue {
                id
                name
              }
              locationValue {
                id
                name
              }
              serviceValue {
                id
                name
              }
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
            locationHierarchy {
              id
              name
              locationType {
                name
              }
            }
          }
        `,
      }),
    ),
  ),
);
