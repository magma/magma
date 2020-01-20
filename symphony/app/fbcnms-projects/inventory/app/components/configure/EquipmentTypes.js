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
import type {EquipmentTypeItem_equipmentType} from './__generated__/EquipmentTypeItem_equipmentType.graphql';
import type {WithStyles} from '@material-ui/core';

import AddEditEquipmentTypeCard from './AddEditEquipmentTypeCard';
import Button from '@fbcnms/ui/components/design-system/Button';
import ConfigueTitle from '@fbcnms/ui/components/ConfigureTitle';
import EquipmentTypeItem from './EquipmentTypeItem';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import withInventoryErrorBoundary from '../../common/withInventoryErrorBoundary';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'relay-runtime';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  root: {
    width: '100%',
    marginTop: '15px',
  },
  paper: {
    flexGrow: 1,
    overflowY: 'hidden',
  },
  content: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-start',
  },
  listItem: {
    marginBottom: theme.spacing(),
  },
  addButton: {
    marginLeft: 'auto',
  },
  addButtonContainer: {
    display: 'flex',
  },
  typesList: {
    padding: '24px',
  },
  title: {
    marginLeft: '10px',
  },
  firstRow: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
});

type Props = ContextRouter &
  WithStyles<typeof styles> & {
    equipmentTypes: Array<EquipmentTypeItem_equipmentType>,
  };

type State = {
  errorMessage: ?string,
  editingEquipmentType: ?EquipmentTypeItem_equipmentType,
  showAddEditCard: boolean,
};

const equipmentTypesQuery = graphql`
  query EquipmentTypesQuery {
    equipmentTypes(first: 500)
      @connection(key: "EquipmentTypes_equipmentTypes") {
      edges {
        node {
          ...EquipmentTypeItem_equipmentType
          ...AddEditEquipmentTypeCard_editingEquipmentType
          id
          name
        }
      }
    }
  }
`;

class EquipmentTypes extends React.Component<Props, State> {
  state = {
    errorMessage: null,
    editingEquipmentType: null,
    showAddEditCard: false,
  };

  render() {
    const {classes} = this.props;
    const {showAddEditCard, editingEquipmentType} = this.state;
    return (
      <InventoryQueryRenderer
        query={equipmentTypesQuery}
        variables={{}}
        render={props => {
          if (showAddEditCard) {
            return (
              <div className={classes.paper}>
                <AddEditEquipmentTypeCard
                  key={'new_equipment_type'}
                  open={showAddEditCard}
                  onClose={this.hideNewEquipmentTypeCard}
                  onSave={this.saveEquipment}
                  editingEquipmentType={editingEquipmentType}
                />
              </div>
            );
          }

          const listItems = props.equipmentTypes.edges
            .map(edge => edge.node)
            .filter(Boolean)
            .sort((eqTypeA, eqTypeB) =>
              sortLexicographically(eqTypeA.name, eqTypeB.name),
            )
            .map(eqType => (
              <div className={classes.listItem} key={`eqType_${eqType.id}`}>
                <EquipmentTypeItem
                  equipmentType={eqType}
                  onEdit={() => this.showAddEditEquipmentTypeCard(eqType)}
                />
              </div>
            ));
          return (
            <div className={classes.typesList}>
              <div className={classes.firstRow}>
                <ConfigueTitle
                  className={classes.title}
                  title={'Equipment Types'}
                  subtitle={'Manage the types of equipment in your inventory'}
                />
                <div className={classes.addButtonContainer}>
                  <Button
                    className={classes.addButton}
                    onClick={() => this.showAddEditEquipmentTypeCard(null)}>
                    Add Equipment Type
                  </Button>
                </div>
              </div>
              <div className={classes.root}>
                <div>{listItems}</div>
              </div>
            </div>
          );
        }}
      />
    );
  }

  showAddEditEquipmentTypeCard = (eqType: ?EquipmentTypeItem_equipmentType) => {
    ServerLogger.info(LogEvents.ADD_EQUIPMENT_TYPE_BUTTON_CLICKED);
    this.setState({editingEquipmentType: eqType, showAddEditCard: true});
  };

  hideNewEquipmentTypeCard = () =>
    this.setState({editingEquipmentType: null, showAddEditCard: false});
  saveEquipment = () => {
    ServerLogger.info(LogEvents.SAVE_EQUIPMENT_TYPE_BUTTON_CLICKED);
    this.setState({
      editingEquipmentType: null,
      showAddEditCard: false,
    });
  };
}

export default withStyles(styles)(
  withRouter(withInventoryErrorBoundary(EquipmentTypes)),
);
