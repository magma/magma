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
import type {Equipment} from '../../common/Equipment';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {
  RemoveEquipmentMutationResponse,
  RemoveEquipmentMutationVariables,
} from '../../mutations/__generated__/RemoveEquipmentMutation.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import CommonStrings from '../../common/CommonStrings';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import React from 'react';
import RemoveEquipmentMutation from '../../mutations/RemoveEquipmentMutation';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import fbt from 'fbt';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {DeleteIcon} from '@fbcnms/ui/components/design-system/Icons';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {capitalize} from '@fbcnms/util/strings';
import {createFragmentContainer, graphql} from 'react-relay';
import {lowerCase} from 'lodash';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  root: {
    width: '100%',
    marginTop: theme.spacing(3),
    overflowX: 'auto',
  },
  table: {
    minWidth: 70,
    marginBottom: '12px',
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
});

type Props = WithSnackbarProps &
  WithAlert &
  WithStyles<typeof styles> & {|
    equipment: Array<Equipment>,
    selectedWorkOrderId: ?string,
    onEquipmentSelected: Equipment => void,
    onWorkOrderSelected: (workOrderId: string) => void,
  |};

class EquipmentTable extends React.Component<Props> {
  static contextType = AppContext;
  context: AppContextType;

  render() {
    const {classes, equipment} = this.props;
    if (equipment.filter(Boolean).length === 0) {
      return null;
    }
    const equipmetStatusEnabled = this.context.isFeatureEnabled(
      'planned_equipment',
    );
    const equipmentLiveStatusEnabled = this.context.isFeatureEnabled(
      'equipment_live_status',
    );

    return equipment.length > 0 ? (
      <Table className={classes.table}>
        <TableHead>
          <TableRow>
            <TableCell>Name</TableCell>
            <TableCell>Type</TableCell>
            <TableCell>Status</TableCell>
            <TableCell />
          </TableRow>
        </TableHead>
        <TableBody>
          {equipment
            .slice()
            .filter(Boolean)
            .sort((x, y) => sortLexicographically(x.name ?? '', y.name ?? ''))
            .map(row => {
              return (
                <TableRow key={row.id}>
                  <TableCell component="th" scope="row">
                    {equipmentLiveStatusEnabled ? (
                      <DeviceStatusCircle
                        isGrey={row.device?.up == null}
                        isActive={row.device?.up ?? false}
                      />
                    ) : null}
                    <Button
                      variant="text"
                      onClick={() => this.props.onEquipmentSelected(row)}>
                      {row.name}
                    </Button>
                  </TableCell>
                  <TableCell component="th" scope="row">
                    {row.equipmentType.name}
                  </TableCell>
                  {equipmetStatusEnabled && (
                    <TableCell>
                      <Button
                        variant="text"
                        onClick={() =>
                          this.props.onWorkOrderSelected(
                            nullthrows(row?.workOrder?.id),
                          )
                        }>
                        {row.futureState
                          ? `${capitalize(
                              lowerCase(row?.workOrder?.status),
                            )} ${lowerCase(row.futureState)}`
                          : ''}
                      </Button>
                    </TableCell>
                  )}
                  <TableCell align="right">
                    <FormAction>
                      <IconButton
                        skin="primary"
                        onClick={() => this.onDelete(row)}
                        icon={DeleteIcon}
                      />
                    </FormAction>
                  </TableCell>
                </TableRow>
              );
            })}
        </TableBody>
      </Table>
    ) : null;
  }

  onDelete(equipment: Equipment) {
    ServerLogger.info(LogEvents.DELETE_EQUIPMENT_CLICKED);
    this.props
      .confirm({
        title: <fbt desc="">Delete Equipment?</fbt>,
        message: (
          <fbt desc="">
            By removing{' '}
            <fbt:param name="equipment name">{equipment.name}</fbt:param> from
            this location, all information related to this equipment, like links
            and sub-positions, will be deleted.
          </fbt>
        ),
        checkboxLabel: <fbt desc="">I understand</fbt>,
        cancelLabel: CommonStrings.common.cancelButton,
        confirmLabel: CommonStrings.common.deleteButton,
        skin: 'red',
      })
      .then(confirmed => {
        if (confirmed) {
          const variables: RemoveEquipmentMutationVariables = {
            id: equipment.id,
            work_order_id: this.props.selectedWorkOrderId,
          };

          const callbacks: MutationCallbacks<RemoveEquipmentMutationResponse> = {
            onCompleted: (_, errors) => {
              if (errors && errors[0]) {
                this.props.enqueueSnackbar(errors[0].message, {
                  children: key => (
                    <SnackbarItem
                      id={key}
                      message={errors[0].message}
                      variant="error"
                    />
                  ),
                });
              }
            },
            onError: (error: any) => {
              this.props.alert('Error: ' + error.source?.errors[0]?.message);
            },
          };

          RemoveEquipmentMutation(variables, callbacks, store => {
            if (!this.props.selectedWorkOrderId) {
              // $FlowFixMe (T62907961) Relay flow types
              store.delete(equipment.id);
            }
          });
        }
      });
  }
}

export default withAlert(
  withStyles(styles)(
    withSnackbar(
      createFragmentContainer(EquipmentTable, {
        equipment: graphql`
          fragment EquipmentTable_equipment on Equipment @relay(plural: true) {
            id
            name
            futureState
            equipmentType {
              id
              name
            }
            workOrder {
              id
              status
            }
            device {
              up
            }
            services {
              id
            }
          }
        `,
      }),
    ),
  ),
);
