/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  AddEquipmentPortTypeMutationResponse,
  AddEquipmentPortTypeMutationVariables,
} from '../../mutations/__generated__/AddEquipmentPortTypeMutation.graphql';
import type {
  EditEquipmentPortTypeMutationResponse,
  EditEquipmentPortTypeMutationVariables,
} from '../../mutations/__generated__/EditEquipmentPortTypeMutation.graphql';
import type {EquipmentPortType} from '../../common/EquipmentType.js';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AddEquipmentPortTypeMutation from '../../mutations/AddEquipmentPortTypeMutation';
import Button from '@fbcnms/ui/components/design-system/Button';
import CardSection from '../CardSection';
import EditEquipmentPortTypeMutation from '../../mutations/EditEquipmentPortTypeMutation';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import Grid from '@material-ui/core/Grid';
import PageFooter from '@fbcnms/ui/components/PageFooter';
import PropertyTypeTable from '../form/PropertyTypeTable';
import React from 'react';
import SectionedCard from '@fbcnms/ui/components/SectionedCard';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import update from 'immutability-helper';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ConnectionHandler} from 'relay-runtime';
import {FormContextProvider} from '../../common/FormContext';
import {createFragmentContainer, graphql} from 'react-relay';
import {getGraphError} from '../../common/EntUtils';
import {getPropertyDefaultValue} from '../../common/PropertyType';
import {sortByIndex} from '../draggable/DraggableUtils';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  root: {
    paddingTop: '24px',
  },
  cardContent: {
    padding: '0px 24px 24px 24px',
  },
  header: {
    marginBottom: '21px',
    paddingBottom: '0px',
  },
  input: {
    display: 'inline-flex',
    width: '305px',
  },
  section: {
    marginBottom: '28px',
  },
  closeButton: {
    marginRight: theme.spacing(),
  },
  headerText: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
  },
  cards: {
    height: 'calc(100% - 60px)',
    padding: '8px 24px',
    overflowY: 'auto',
  },
});

type Props = WithSnackbarProps &
  WithStyles<typeof styles> &
  WithAlert & {
    open: boolean,
    onClose: () => void,
    onSave: (portType: any) => void,
    editingEquipmentPortType?: ?EquipmentPortType,
  };

type State = {
  editingEquipmentPortType: EquipmentPortType,
};

class AddEditEquipmentPortTypeCard extends React.Component<Props, State> {
  state = {
    editingEquipmentPortType: this.getEditingEquipmentPortType(),
  };

  _nameInputRef = React.createRef();

  componentDidMount() {
    this._nameInputRef.current && this._nameInputRef.current.focus();
  }

  render() {
    const {classes, onClose} = this.props;
    const {editingEquipmentPortType} = this.state;
    const propertyTypes = editingEquipmentPortType.propertyTypes
      .slice()
      .sort(sortByIndex);
    const linkPropertyTypes = editingEquipmentPortType.linkPropertyTypes
      .slice()
      .sort(sortByIndex);

    const isOnEdit = !!this.props.editingEquipmentPortType;
    return (
      <FormContextProvider
        permissions={{
          entity: 'location',
          action: isOnEdit ? 'update' : 'create',
        }}>
        <div className={classes.cards}>
          <SectionedCard>
            <div className={classes.header}>
              <Text className={classes.headerText}>
                {this.props.editingEquipmentPortType
                  ? 'Edit Port Type'
                  : 'New Port Type'}
              </Text>
            </div>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <FormField label="Name" required>
                  <TextInput
                    name="name"
                    variant="outlined"
                    className={classes.input}
                    value={editingEquipmentPortType.name}
                    onChange={this.nameChanged}
                    ref={this._nameInputRef}
                  />
                </FormField>
              </Grid>
            </Grid>
          </SectionedCard>
          <SectionedCard>
            <Grid container direction="column" spacing={3}>
              <Grid item xs={12} xl={7}>
                <CardSection
                  className={classes.section}
                  title="Port Properties">
                  <PropertyTypeTable
                    propertyTypes={propertyTypes}
                    onPropertiesChanged={this._propertyChangedHandler}
                    supportMandatory={false}
                  />
                </CardSection>
              </Grid>
            </Grid>
          </SectionedCard>
          <SectionedCard>
            <Grid container direction="column" spacing={3}>
              <Grid item xs={12} xl={7}>
                <CardSection
                  className={classes.section}
                  title="Link Properties">
                  <PropertyTypeTable
                    propertyTypes={linkPropertyTypes}
                    onPropertiesChanged={this._linkPropertyChangedHandler}
                    supportMandatory={false}
                  />
                </CardSection>
              </Grid>
            </Grid>
          </SectionedCard>
        </div>
        <PageFooter>
          <Button
            className={classes.closeButton}
            skin="regular"
            onClick={onClose}>
            Cancel
          </Button>
          <Button disabled={this.isSaveDisabled()} onClick={this.onSave}>
            Save
          </Button>
        </PageFooter>
      </FormContextProvider>
    );
  }

  isSaveDisabled() {
    return (
      !this.state.editingEquipmentPortType.name ||
      !this.state.editingEquipmentPortType.propertyTypes.every(propType => {
        return (
          propType.isInstanceProperty || !!getPropertyDefaultValue(propType)
        );
      }) ||
      !this.state.editingEquipmentPortType.linkPropertyTypes.every(propType => {
        return (
          propType.isInstanceProperty || !!getPropertyDefaultValue(propType)
        );
      })
    );
  }

  buildAddMutationVariables = (): AddEquipmentPortTypeMutationVariables => {
    const {id: _, ...addVars} = this.buildEditMutationVariables().input;
    return {input: {...addVars}};
  };

  buildEditMutationVariables = (): EditEquipmentPortTypeMutationVariables => {
    const {
      id,
      name,
      propertyTypes,
      linkPropertyTypes,
    } = this.state.editingEquipmentPortType;

    const deleteTempId = (definition: Object) => {
      const newDef = {...definition};
      if (definition.id && definition.id.includes('@tmp')) {
        newDef['id'] = undefined;
      }
      return newDef;
    };
    return {
      input: {
        id: id,
        name: name,
        properties: propertyTypes
          .filter(propType => !!propType.name)
          .map(deleteTempId),
        linkProperties: linkPropertyTypes
          .filter(propType => !!propType.name)
          .map(deleteTempId),
      },
    };
  };

  onSave = () => {
    const {name} = this.state.editingEquipmentPortType;
    if (!name) {
      const eror = 'Name cannot be empty';
      this.props.enqueueSnackbar(eror, {
        children: key => (
          <SnackbarItem id={key} message={eror} variant="error" />
        ),
      });
      return;
    }
    if (this.props.editingEquipmentPortType) {
      this.editEquipmentPortType();
    } else {
      this.addNewEquipmentPortType();
    }
  };

  editEquipmentPortType = () => {
    const variables = this.buildEditMutationVariables();
    const callbacks: MutationCallbacks<EditEquipmentPortTypeMutationResponse> = {
      onCompleted: (response, errors) => {
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
        } else {
          this.props.onSave &&
            this.props.onSave(response.editEquipmentPortType);
        }
      },

      onError: (error: Error) => {
        const msg = getGraphError(error);
        this.props.enqueueSnackbar(msg, {
          children: key => (
            <SnackbarItem id={key} message={msg} variant="error" />
          ),
        });
      },
    };

    EditEquipmentPortTypeMutation(variables, callbacks);
  };

  addNewEquipmentPortType = () => {
    const variables = this.buildAddMutationVariables();
    const callbacks: MutationCallbacks<AddEquipmentPortTypeMutationResponse> = {
      onCompleted: (response, errors) => {
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
        } else {
          this.props.onSave && this.props.onSave(response.addEquipmentPortType);
        }
      },
      onError: (error: Error) => {
        const msg = getGraphError(error);
        this.props.enqueueSnackbar(msg, {
          children: key => (
            <SnackbarItem id={key} message={msg} variant="error" />
          ),
        });
      },
    };
    const updater = store => {
      const rootQuery = store.getRoot();
      const newNode = store.getRootField('addEquipmentPortType');
      if (!newNode) {
        return;
      }
      const types = ConnectionHandler.getConnection(
        rootQuery,
        'EquipmentPortTypes_equipmentPortTypes',
      );
      if (types == null) {
        return;
      }

      const edge = ConnectionHandler.createEdge(
        store,
        types,
        newNode,
        'EquipmentPortTypesEdge',
      );
      ConnectionHandler.insertEdgeBefore(types, edge);
    };
    AddEquipmentPortTypeMutation(variables, callbacks, updater);
  };

  fieldChangedHandler = (field: 'name') => event =>
    this.setState({
      editingEquipmentPortType: {
        ...this.state.editingEquipmentPortType,
        [field]: event.target.value,
      },
    });

  nameChanged = this.fieldChangedHandler('name');

  _propertyChangedHandler = properties => {
    this.setState(prevState => {
      return {
        editingEquipmentPortType: update(prevState.editingEquipmentPortType, {
          propertyTypes: {$set: properties},
        }),
      };
    });
  };

  _linkPropertyChangedHandler = properties => {
    this.setState(prevState => {
      return {
        editingEquipmentPortType: update(prevState.editingEquipmentPortType, {
          linkPropertyTypes: {$set: properties},
        }),
      };
    });
  };

  getEditingEquipmentPortType(): EquipmentPortType {
    const {editingEquipmentPortType} = this.props;
    const index = editingEquipmentPortType?.propertyTypes.length ?? 0;
    return {
      id: editingEquipmentPortType?.id ?? 'EquipmentPortType@tmp0',
      name: editingEquipmentPortType?.name ?? '',
      numberOfPortDefinitions:
        editingEquipmentPortType?.numberOfPortDefinitions ?? 0,
      propertyTypes: [
        ...(editingEquipmentPortType?.propertyTypes ?? []),
        {
          id: 'PropertyType@tmp' + index,
          name: '',
          index: index,
          type: 'string',
          nodeType: null,
          booleanValue: false,
          stringValue: null,
          intValue: null,
          floatValue: null,
          latitudeValue: null,
          longitudeValue: null,
          isEditable: true,
          isInstanceProperty: true,
        },
      ],
      linkPropertyTypes: [
        ...(editingEquipmentPortType?.linkPropertyTypes ?? []),
        {
          id: 'LinkPropertyType@tmp' + index,
          name: '',
          index: index,
          type: 'string',
          nodeType: null,
          booleanValue: false,
          stringValue: null,
          intValue: null,
          floatValue: null,
          latitudeValue: null,
          longitudeValue: null,
          isEditable: true,
          isInstanceProperty: true,
        },
      ],
    };
  }
}

export default withStyles(styles)(
  withAlert(
    withSnackbar(
      createFragmentContainer(AddEditEquipmentPortTypeCard, {
        /* eslint-disable relay/unused-fields */
        editingEquipmentPortType: graphql`
          fragment AddEditEquipmentPortTypeCard_editingEquipmentPortType on EquipmentPortType {
            id
            name
            numberOfPortDefinitions
            propertyTypes {
              id
              name
              type
              nodeType
              index
              stringValue
              intValue
              booleanValue
              floatValue
              latitudeValue
              longitudeValue
              isEditable
              isInstanceProperty
            }
            linkPropertyTypes {
              id
              name
              type
              nodeType
              index
              stringValue
              intValue
              booleanValue
              floatValue
              latitudeValue
              longitudeValue
              isEditable
              isInstanceProperty
            }
          }
        `,
      }),
    ),
  ),
);
