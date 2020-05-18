/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions} from './__generated__/ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions.graphql';

import CardSection from '../CardSection';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import inventoryTheme from '../../common/theme';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {sortByIndex} from '../draggable/DraggableUtils';

const useStyles = makeStyles(_theme => ({
  table: inventoryTheme.table,
  cell: {
    paddingLeft: '0px',
  },
}));

type Props = {
  serviceEndpointDefinitions: ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions,
};

const ServiceEndpointDefinitionStaticTable = (props: Props) => {
  const {serviceEndpointDefinitions} = props;
  const classes = useStyles();

  if (serviceEndpointDefinitions.length === 0) {
    return null;
  }
  return (
    <CardSection title="Endpoint Types">
      <Table component="div" className={classes.table}>
        <TableHead component="div">
          <TableRow component="div">
            <TableCell component="div" className={classes.cell}>
              <fbt desc="">Name</fbt>
            </TableCell>
            <TableCell component="div" className={classes.cell}>
              <fbt desc="">Endpoint Function</fbt>
            </TableCell>
            <TableCell component="div" className={classes.cell}>
              <fbt desc="">Equipment Type</fbt>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody component="div">
          {serviceEndpointDefinitions
            .slice()
            .sort(sortByIndex)
            .map(serviceEndpointDefinition => (
              <TableRow component="div" key={serviceEndpointDefinition.id}>
                <TableCell className={classes.cell} component="div" scope="row">
                  <Text variant="body2">{serviceEndpointDefinition.name}</Text>
                </TableCell>
                <TableCell className={classes.cell} component="div" scope="row">
                  <Text variant="body2">{serviceEndpointDefinition.role}</Text>
                </TableCell>
                <TableCell className={classes.cell} component="div" scope="row">
                  <Text variant="body2">
                    {serviceEndpointDefinition.equipmentType.name}
                  </Text>
                </TableCell>
              </TableRow>
            ))}
        </TableBody>
      </Table>
    </CardSection>
  );
};

export default createFragmentContainer(ServiceEndpointDefinitionStaticTable, {
  serviceEndpointDefinitions: graphql`
    fragment ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions on ServiceEndpointDefinition
      @relay(plural: true) {
      id
      name
      role
      index
      equipmentType {
        id
        name
      }
    }
  `,
});
