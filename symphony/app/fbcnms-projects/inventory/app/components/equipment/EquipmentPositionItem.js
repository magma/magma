/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment, EquipmentPosition} from '../../common/Equipment';
import type {EquipmentType} from '../../common/EquipmentType';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {
  RemoveEquipmentFromPositionMutationResponse,
  RemoveEquipmentFromPositionMutationVariables,
} from '../../mutations/__generated__/RemoveEquipmentFromPositionMutation.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import ActionButton from '@fbcnms/ui/components/ActionButton';
import AddToEquipmentDialog from './AddToEquipmentDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import CommonStrings from '@fbcnms/strings/Strings';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import React from 'react';
import RemoveEquipmentFromPositionMutation from '../../mutations/RemoveEquipmentFromPositionMutation';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Text from '@fbcnms/ui/components/design-system/Text';
import Typography from '@material-ui/core/Typography';
import classNames from 'classnames';
import fbt from 'fbt';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {capitalize} from '@fbcnms/util/strings';
import {gray0, gray2} from '@fbcnms/ui/theme/colors';
import {lowerCase} from 'lodash';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  equipment: Equipment,
  position: EquipmentPosition,
  onAttachingEquipmentToPosition: (
    equipmentType: EquipmentType,
    position: EquipmentPosition,
  ) => void,
  onEquipmentPositionClicked: (equipmentId: string) => void,
  workOrderId: ?string,
  onWorkOrderSelected: (workOrderId: string) => void,
} & WithStyles<typeof styles> &
  WithAlert &
  WithSnackbarProps;

type State = {
  isNewEquipmentDialogOpen: boolean,
};

const styles = theme => ({
  root: {
    padding: '10px',
    backgroundColor: gray0,
    borderRadius: '3px',
    width: '172px',
    display: 'flex',
    '&:hover': {
      backgroundColor: theme.palette.grey[50],
      boxShadow: theme.shadows[1],
    },
  },
  equipmentRoot: {
    paddingTop: '5px',
    paddingBottom: '5px',
    alignItems: 'center',
  },
  positionBody: {
    marginRight: theme.spacing(),
    display: 'inline',
    flexGrow: 1,
  },
  equipment: {
    height: '45px',
  },
  equipmentDetails: {
    display: 'flex',
    lineHeight: '1em',
    marginTop: '2px',
    color: gray0,
  },
  equipmentName: {
    whiteSpace: 'nowrap',
    textOverflow: 'ellipsis',
    paddingLeft: '2px',
  },
  equipmentState: {
    marginTop: '2px',
    lineHeight: '1em',
    color: gray2,
  },
  equipmentPositionName: {
    color: gray2,
  },
});

class EquipmentPositionItem extends React.Component<Props, State> {
  state = {
    isNewEquipmentDialogOpen: false,
  };

  render() {
    const {classes, position} = this.props;
    const positionOccupied = position.attachedEquipment != null;
    return (
      <div
        className={classNames({
          [classes.root]: true,
          [classes.equipmentRoot]: true,
        })}>
        <div className={classes.positionBody}>{this.renderEquipment()}</div>
        <FormAction>
          <ActionButton
            action={positionOccupied ? 'remove' : 'add'}
            onClick={() => {
              if (position.attachedEquipment == null) {
                this.setState({isNewEquipmentDialogOpen: true});
                return;
              }
              this.props
                .confirm({
                  title: <fbt desc="">Delete Equipment?</fbt>,
                  message: (
                    <div>
                      <fbt desc="">
                        By removing{' '}
                        <fbt:param name="equipment name">
                          {position.attachedEquipment.name}
                        </fbt:param>{' '}
                        from this position, all information related to this
                        equipment, like links and sub-positions, will be
                        deleted.
                      </fbt>
                      {position.attachedEquipment.services.length > 0 && (
                        <p>
                          <fbt desc="">
                            This attached equipment is used by some services and
                            deleting it can potentially break them.
                          </fbt>
                        </p>
                      )}
                    </div>
                  ),
                  checkboxLabel: <fbt desc="">I understand</fbt>,
                  cancelLabel: CommonStrings.common.cancelButton,
                  confirmLabel: CommonStrings.common.deleteButton,
                  skin: 'red',
                })
                .then(
                  confirmed =>
                    confirmed && this.onDetachEquipmentFromPosition(),
                );
            }}
          />
          <AddToEquipmentDialog
            open={this.state.isNewEquipmentDialogOpen}
            onClose={() => this.setState({isNewEquipmentDialogOpen: false})}
            onEquipmentTypeSelected={equipmentType =>
              this.props.onAttachingEquipmentToPosition(equipmentType, position)
            }
            parentEquipment={this.props.equipment}
            position={position}
          />
        </FormAction>
      </div>
    );
  }

  renderEquipment() {
    const {position, classes} = this.props;
    const equipment = position.attachedEquipment;
    if (equipment === null || equipment === undefined) {
      return (
        <div className={classes.equipment}>
          <div className={classes.equipmentDetails}>
            <Typography
              variant="body2"
              className={classes.equipmentPositionName}>
              {`${position.definition.name}: Available`}
            </Typography>
          </div>
        </div>
      );
    }

    return (
      <div className={classes.equipment}>
        <div className={classes.equipmentDetails}>
          <Text variant="body2" className={classes.equipmentPositionName}>
            {`${position.definition.name}: `}
          </Text>
          <Button
            className={classes.equipmentName}
            variant="text"
            onClick={() => this.props.onEquipmentPositionClicked(equipment.id)}>
            {equipment.name}
          </Button>
        </div>
        {equipment.futureState && (
          <div>
            <Button
              variant="text"
              skin="regular"
              onClick={() =>
                this.props.onWorkOrderSelected(
                  nullthrows(equipment?.workOrder?.id),
                )
              }>
              {`${capitalize(
                lowerCase(equipment?.workOrder?.status),
              )} ${lowerCase(equipment?.futureState)}`}
            </Button>
          </div>
        )}
      </div>
    );
  }

  onDetachEquipmentFromPosition = () => {
    const variables: RemoveEquipmentFromPositionMutationVariables = {
      position_id: this.props.position.id,
      work_order_id: this.props.workOrderId,
    };

    const callbacks: MutationCallbacks<RemoveEquipmentFromPositionMutationResponse> = {
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
      onError: () => {
        this.props.alert('Error removing equipment from position');
      },
    };

    RemoveEquipmentFromPositionMutation(variables, callbacks);
  };
}

export default withStyles(styles)(
  withAlert(withSnackbar(EquipmentPositionItem)),
);
