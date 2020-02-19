/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AddEditProjectTypeCard_editingProjectType} from './__generated__/AddEditProjectTypeCard_editingProjectType.graphql';
import type {AddProjectTypeMutationVariables} from '../../mutations/__generated__/AddWorkOrderTypeMutation.graphql';
import type {EditProjectTypeInput} from './mutations/__generated__/EditProjectTypeMutation.graphql';
import type {EditProjectTypeMutationVariables} from '../../mutations/__generated__/EditWorkOrderTypeMutation.graphql';
import type {ProjectTypeWorkOrderTemplatesPanel_workOrderTypes} from './__generated__/ProjectTypeWorkOrderTemplatesPanel_workOrderTypes.graphql';

import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import Button from '@fbcnms/ui/components/design-system/Button';
import CreateProjectTypeMutation from './mutations/CreateProjectTypeMutation';
import EditProjectTypeMutation from './mutations/EditProjectTypeMutation';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import NameDescriptionSection from '@fbcnms/ui/components/NameDescriptionSection';
import ProjectTypeWorkOrderTemplatesPanel from './ProjectTypeWorkOrderTemplatesPanel';
import PropertyTypeTable from '../form/PropertyTypeTable';
import React, {useCallback, useMemo, useState} from 'react';
import update from 'immutability-helper';
import {ConnectionHandler} from 'relay-runtime';
import {FormValidationContextProvider} from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {sortByIndex} from '../draggable/DraggableUtils';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

const useStyles = makeStyles(() => ({
  root: {
    padding: '24px 16px',
    maxHeight: '100%',
    flexGrow: 1,
    overflow: 'hidden',
    display: 'flex',
    flexDirection: 'column',
  },
  header: {
    display: 'flex',
    paddingBottom: '24px',
  },
  body: {
    overflowY: 'auto',
  },
  buttons: {
    display: 'flex',
  },
  cancelButton: {
    marginRight: '8px',
  },
}));

type Props = {
  workOrderTypes: ProjectTypeWorkOrderTemplatesPanel_workOrderTypes,
  editingProjectType: ?AddEditProjectTypeCard_editingProjectType,
  onCancelClicked: () => void,
  onProjectTypeSaved: () => void,
};

const AddEditProjectTypeCard = (props: Props) => {
  const {
    editingProjectType,
    onCancelClicked,
    onProjectTypeSaved,
    workOrderTypes,
  } = props;
  const enqueueSnackbar = useEnqueueSnackbar();
  const classes = useStyles();
  const initialProjectTypeInput: EditProjectTypeInput = useMemo(
    () => ({
      id: editingProjectType?.id ?? '',
      name: editingProjectType?.name ?? '',
      description: editingProjectType?.description ?? undefined,
      workOrders: (editingProjectType?.workOrders ?? [])
        .map(wo => wo?.type)
        .filter(Boolean)
        .map(woType => ({type: woType.id})),
      properties: (editingProjectType?.properties ?? []: any)
        .slice()
        .sort(sortByIndex),
    }),
    [editingProjectType],
  );

  const [projectTypeInput, setProjectTypeInput] = useState(
    initialProjectTypeInput,
  );

  const deleteTempId = (definition: {id: ?string}) => {
    const newDef = {...definition};
    if (definition.id && definition.id.includes('@tmp')) {
      newDef['id'] = undefined;
    }
    return newDef;
  };

  const onAdd = useCallback(() => {
    const variables: AddProjectTypeMutationVariables = {
      input: {
        name: projectTypeInput.name,
        description: projectTypeInput.description ?? undefined,
        workOrders: projectTypeInput.workOrders,
        properties: (projectTypeInput.properties ?? [])
          .filter(propType => !!propType.name)
          .map(deleteTempId),
      },
    };

    const updater = store => {
      const rootQuery = store.getRoot();
      const newNode = store.getRootField('createProjectType');
      if (!newNode) {
        return;
      }
      const types = ConnectionHandler.getConnection(
        rootQuery,
        'WorkOrderProjectTypesQuery_projectTypes',
      );
      const edge = ConnectionHandler.createEdge(
        store,
        types,
        newNode,
        'ProjectTypesEdge',
      );
      ConnectionHandler.insertEdgeAfter(types, edge);
    };

    const callbacks = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          enqueueSnackbar(errors[0].message, {
            variant: 'error',
          });
        } else {
          onProjectTypeSaved();
        }
      },
    };

    CreateProjectTypeMutation(variables, callbacks, updater);
  }, [projectTypeInput, onProjectTypeSaved, enqueueSnackbar]);

  const onEdit = useCallback(() => {
    const woDefsMap = new Map();
    (editingProjectType?.workOrders ?? []).forEach(wo =>
      woDefsMap.set(wo?.type?.id, wo?.id),
    );
    const variables: EditProjectTypeMutationVariables = {
      input: {
        id: projectTypeInput.id,
        name: projectTypeInput.name,
        description: projectTypeInput.description ?? undefined,
        workOrders: (projectTypeInput.workOrders ?? []).map(x => ({
          id: woDefsMap.get(x?.type),
          type: x?.type,
        })),
        properties: (projectTypeInput.properties ?? [])
          .filter(propType => !!propType.name)
          .map(deleteTempId),
      },
    };
    const callbacks = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          enqueueSnackbar(errors[0].message, {
            variant: 'error',
          });
        } else {
          onProjectTypeSaved();
        }
      },
    };

    EditProjectTypeMutation(variables, callbacks);
  }, [
    editingProjectType,
    projectTypeInput.id,
    projectTypeInput.name,
    projectTypeInput.description,
    projectTypeInput.workOrders,
    projectTypeInput.properties,
    enqueueSnackbar,
    onProjectTypeSaved,
  ]);

  const onSave = useCallback(() => {
    if (editingProjectType !== null) {
      return onEdit();
    }
    return onAdd();
  }, [editingProjectType, onAdd, onEdit]);

  return (
    <FormValidationContextProvider>
      <div className={classes.root}>
        <div className={classes.header}>
          <Breadcrumbs
            breadcrumbs={[
              {
                id: 'project_templates',
                name: 'Project Templates',
                onClick: onCancelClicked,
              },
              editingProjectType
                ? {
                    id: editingProjectType.id,
                    name: editingProjectType.name,
                  }
                : {
                    id: 'new_project_type',
                    name: 'New Project Template',
                  },
            ]}
            size="large"
          />
          <div className={classes.buttons}>
            <Button
              className={classes.cancelButton}
              onClick={onCancelClicked}
              skin="regular">
              Cancel
            </Button>
            <FormAction>
              <Button onClick={onSave} disabled={!projectTypeInput.name}>
                Save
              </Button>
            </FormAction>
          </div>
        </div>
        <div className={classes.body}>
          <ExpandingPanel title="Details">
            <NameDescriptionSection
              title="Project Name"
              name={projectTypeInput.name}
              description={projectTypeInput.description}
              descriptionPlaceholder="Describe the project"
              onNameChange={value =>
                setProjectTypeInput(
                  update(projectTypeInput, {name: {$set: value}}),
                )
              }
              onDescriptionChange={value =>
                setProjectTypeInput(
                  update(projectTypeInput, {description: {$set: value}}),
                )
              }
            />
          </ExpandingPanel>
          <ProjectTypeWorkOrderTemplatesPanel
            selectedWorkOrderTypeIds={(projectTypeInput.workOrders ?? []).map(
              wo => wo.type,
            )}
            workOrderTypes={workOrderTypes}
            onWorkOrderTypesSelected={ids => {
              setProjectTypeInput(
                update(projectTypeInput, {
                  workOrders: {$set: ids.map(id => ({type: id}))},
                }),
              );
            }}
          />
          <ExpandingPanel title="Properties">
            <PropertyTypeTable
              supportDelete={true}
              // eslint-disable-next-line flowtype/no-weak-types
              propertyTypes={(projectTypeInput.properties ?? []: any)}
              onPropertiesChanged={properties => {
                setProjectTypeInput(
                  update(projectTypeInput, {
                    properties: {$set: properties},
                  }),
                );
              }}
            />
          </ExpandingPanel>
        </div>
      </div>
    </FormValidationContextProvider>
  );
};

export default createFragmentContainer(AddEditProjectTypeCard, {
  editingProjectType: graphql`
    fragment AddEditProjectTypeCard_editingProjectType on ProjectType {
      id
      name
      description
      workOrders {
        id
        type {
          id
          name
        }
      }
      properties {
        id
        name
        type
        index
        stringValue
        intValue
        booleanValue
        floatValue
        latitudeValue
        longitudeValue
        rangeFromValue
        rangeToValue
        isEditable
        isMandatory
        isInstanceProperty
        isDeleted
      }
    }
  `,
});
