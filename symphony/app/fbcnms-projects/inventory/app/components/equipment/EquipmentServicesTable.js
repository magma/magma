/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Service} from '../../common/Service';
import type {Theme, WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import EquipmentServicesTableMenu from './EquipmentServicesTableMenu';
import IconButton from '@material-ui/core/IconButton';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import Paper from '@material-ui/core/Paper';
import React, {useState} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import {createFragmentContainer, graphql} from 'react-relay';
import {useHistory} from 'react-router';
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
  services: Array<Service>,
  ...$Exact<WithStyles<typeof styles>>,
};

const EquipmentServicesTable = (props: Props) => {
  const [anchorEl, setAnchorEl] = useState<?HTMLElement>(null);
  const [selectedService, setSelectedService] = useState<?Service>(null);
  const history = useHistory();

  const navigateToService = (serviceId: string) => {
    history.push(`/inventory/services?service=${serviceId}`);
  };

  const {services, classes} = props;

  const headCells = [
    {label: 'Name', key: 'name'},
    {label: 'Type', key: 'type'},
    {label: 'Service ID', key: 'service_id'},
    {label: 'Customer', key: 'customer'},
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
                  {cell.label}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {services.slice().map((service, i) => {
              return (
                <TableRow key={`service_${i}`}>
                  <TableCell
                    className={classes.rowFirstCell}
                    component="th"
                    scope="row">
                    <Button
                      variant="text"
                      onClick={() => navigateToService(service.id)}>
                      {service.name}
                    </Button>
                  </TableCell>
                  <TableCell>{service.serviceType.name}</TableCell>
                  <TableCell>{service.externalId}</TableCell>
                  <TableCell>{service.customer?.name}</TableCell>
                  <TableCell>
                    <IconButton
                      onClick={event => {
                        setAnchorEl(event.currentTarget);
                        setSelectedService(service);
                      }}
                      color="secondary">
                      <MoreVertIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </Paper>
      {selectedService && (
        <EquipmentServicesTableMenu
          key={`${selectedService.id}-menu`}
          service={selectedService}
          anchorEl={anchorEl}
          onClose={() => setAnchorEl(null)}
          onViewService={serviceId => navigateToService(serviceId)}
        />
      )}
    </>
  );
};

export default withStyles(styles)(
  createFragmentContainer(EquipmentServicesTable, {
    equipment: graphql`
      fragment EquipmentServicesTable_equipment on Equipment {
        id
        name
        services {
          id
          name
          externalId
          customer {
            name
          }
          serviceType {
            id
            name
          }
        }
      }
    `,
  }),
);
