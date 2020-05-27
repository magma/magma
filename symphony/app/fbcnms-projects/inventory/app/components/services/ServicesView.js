/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {ServicesView_service} from './__generated__/ServicesView_service.graphql.js';
import type {WithStyles} from '@material-ui/core';

import Link from '@fbcnms/ui/components/Link';
import LocationLink from '../location/LocationLink';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {AutoSizer, Column, Table} from 'react-virtualized';
import {createFragmentContainer, graphql} from 'react-relay';
import {discoveryMethods} from '../../common/Service';
import {serviceStatusToVisibleNames} from '../../common/Service';
import {withStyles} from '@material-ui/core/styles';

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
  linkText: {
    color: 'inherit',
    display: 'inline',
    fontWeight: 'bold',
  },
});

type Props = {
  onServiceSelected: (serviceId: string) => void,
  service: ServicesView_service,
} & WithStyles<typeof styles> &
  ContextRouter;

class ServicesView extends React.Component<Props> {
  _headerRenderer = ({label}) => {
    const {classes} = this.props;
    return (
      <div className={classes.cell}>
        <Text className={classes.headerText}>{label}</Text>
      </div>
    );
  };

  _nameRenderer = ({rowData}) => {
    const {classes, onServiceSelected} = this.props;
    const content = (
      <Link onClick={() => onServiceSelected(rowData.id)}>{rowData.name}</Link>
    );
    return <div className={classes.cell}>{content}</div>;
  };

  _locationRenderer = ({cellData}) => {
    if (cellData == null) return;
    return <LocationLink title={cellData.name} id={cellData.id} />;
  };

  _cellRenderer = ({dataKey, _, cellData}) => {
    const {classes} = this.props;
    let data = cellData ?? '';
    if (dataKey === 'status') {
      data = serviceStatusToVisibleNames[data];
    }
    if (dataKey === 'discovery_method') {
      if (data == '') {
        data = discoveryMethods.MANUAL;
      } else {
        data = discoveryMethods[data];
      }
    }
    const content = <Text className={classes.cellText}>{data}</Text>;
    return <div className={classes.cell}>{content}</div>;
  };

  render() {
    const {classes, service} = this.props;
    if (service.length === 0) {
      return <div />;
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
            rowCount={service.length}
            rowGetter={({index}) => service[index]}
            gridClassName={classes.table}
            rowClassName={({index}) => (index === -1 ? classes.header : '')}>
            <Column
              label="Name"
              dataKey="name"
              width={300}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._nameRenderer}
            />
            <Column
              label="Type"
              dataKey="type"
              cellDataGetter={({rowData}) => rowData.serviceType?.name}
              width={180}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
            />
            <Column
              label="Discovery Method"
              dataKey="discovery_method"
              cellDataGetter={({rowData}) =>
                rowData.serviceType?.discoveryMethod
              }
              width={180}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
            />
            <Column
              label="Service ID"
              dataKey="service_id"
              cellDataGetter={({rowData}) => rowData.externalId}
              width={100}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
            />
            <Column
              label="Customer"
              dataKey="customer"
              cellDataGetter={({rowData}) => rowData.customer?.name}
              width={100}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
            />
            <Column
              label="Status"
              dataKey="status"
              cellDataGetter={({rowData}) => rowData.status}
              width={100}
              flexGrow={1}
              headerRenderer={this._headerRenderer}
              cellRenderer={this._cellRenderer}
            />
          </Table>
        )}
      </AutoSizer>
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(ServicesView, {
    service: graphql`
      fragment ServicesView_service on Service @relay(plural: true) {
        id
        name
        externalId
        status
        customer {
          id
          name
        }
        serviceType {
          id
          name
          discoveryMethod
          propertyTypes {
            ...PropertyTypeFormField_propertyType
            ...DynamicPropertiesGrid_propertyTypes
          }
        }
        properties {
          ...PropertyFormField_property
          ...DynamicPropertiesGrid_properties
        }
      }
    `,
  }),
);
