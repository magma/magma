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
  LocationTypesQuery,
  LocationTypesQueryResponse,
} from './__generated__/LocationTypesQuery.graphql';

import AddEditLocationTypeCard from './AddEditLocationTypeCard';
import Button from '@fbcnms/ui/components/design-system/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import ConfigueTitle from '@fbcnms/ui/components/ConfigureTitle';
import DroppableTableBody from '../draggable/DroppableTableBody';
import FormActionWithPermissions from '../../common/FormActionWithPermissions';
import LocationTypeItem from './LocationTypeItem';
import React, {useState} from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import {FormContextProvider} from '../../common/FormContext';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';
import {reorder, sortByIndex} from '../draggable/DraggableUtils';
import {saveLocationTypeIndexes} from '../../mutations/EditLocationTypesIndexMutation';
import {useEnqueueSnackbar} from '@fbcnms/alarms/hooks/useSnackbar';
import {useLazyLoadQuery} from 'react-relay/hooks';

const useStyles = makeStyles(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  root: {
    display: 'flex',
    width: '100%',
    flexDirection: 'column',
  },
  table: {
    width: '100%',
    marginTop: '15px',
  },
  paper: {
    flexGrow: 1,
    overflowY: 'hidden',
  },
  typesList: {
    padding: '24px',
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
  progress: {
    alignSelf: 'center',
  },
  title: {
    marginLeft: '10px',
  },
  firstRow: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
}));

type ResponseLocationType = $NonMaybeType<
  $ElementType<
    $ElementType<
      $ElementType<
        $NonMaybeType<
          $ElementType<LocationTypesQueryResponse, 'locationTypes'>,
        >,
        'edges',
      >,
      number,
    >,
    'node',
  >,
>;

const locationTypesQuery = graphql`
  query LocationTypesQuery {
    locationTypes(first: 500) @connection(key: "Catalog_locationTypes") {
      edges {
        node {
          ...LocationTypeItem_locationType
          ...AddEditLocationTypeCard_editingLocationType
          id
          name
          index
        }
      }
    }
  }
`;

const LocationTypes = () => {
  const classes = useStyles();
  const {
    locationTypes,
  }: LocationTypesQueryResponse = useLazyLoadQuery<LocationTypesQuery>(
    locationTypesQuery,
  );

  const locationTypesData: Array<ResponseLocationType> =
    locationTypes?.edges.map(edge => edge.node).filter(Boolean) ?? [];

  const [
    editingLocationType,
    setEditingLocationType,
  ] = useState<?ResponseLocationType>(null);
  const [showAddEditCard, setShowAddEditCard] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const enqueueSnackbar = useEnqueueSnackbar();

  const showAddEditLocationTypeCard = (locType: ?ResponseLocationType) => {
    ServerLogger.info(LogEvents.ADD_LOCATION_TYPE_BUTTON_CLICKED);
    setEditingLocationType(locType);
    setShowAddEditCard(true);
  };

  const hideNewLocationTypeCard = () => {
    setEditingLocationType(null);
    setShowAddEditCard(false);
  };

  const saveLocation = () => {
    ServerLogger.info(LogEvents.SAVE_LOCATION_TYPE_BUTTON_CLICKED);
    setEditingLocationType(null);
    setShowAddEditCard(false);
  };

  if (showAddEditCard) {
    return (
      <div className={classes.paper}>
        <AddEditLocationTypeCard
          open={showAddEditCard}
          onClose={hideNewLocationTypeCard}
          onSave={saveLocation}
          editingLocationType={editingLocationType}
        />
      </div>
    );
  }

  const buildMutationVariables = (newItems: Array<ResponseLocationType>) => {
    return newItems
      .map(item => {
        if (item.index == null) {
          return null;
        }
        return {
          locationTypeID: item.id,
          index: item.index,
        };
      })
      .filter(Boolean);
  };

  const saveOrder = newItems => {
    setIsSaving(true);
    saveLocationTypeIndexes(buildMutationVariables(newItems))
      .catch((errorMessage: string) =>
        enqueueSnackbar(errorMessage, {
          children: (key: string) => (
            <SnackbarItem id={key} message={errorMessage} variant="error" />
          ),
        }),
      )
      .finally(() => setIsSaving(false));
  };

  const onDragEnd = ({source, destination}) => {
    if (destination == null) {
      return;
    }

    ServerLogger.info(LogEvents.LOCATION_TYPE_REORDERED);
    const items = reorder(locationTypesData, source.index, destination.index);
    const newItems = items.map((locTyp: ResponseLocationType, i) => ({
      ...locTyp,
      index: i,
    }));
    saveOrder(newItems);
  };

  return (
    <FormContextProvider
      permissions={{
        entity: 'locationType',
      }}>
      <div className={classes.typesList}>
        <div className={classes.firstRow}>
          <ConfigueTitle
            className={classes.title}
            title={'Location Types'}
            subtitle={
              'Drag and drop location types to arrange them by size, from largest to smallest'
            }
          />
          <div className={classes.addButtonContainer}>
            {isSaving ? (
              <CircularProgress className={classes.progress} />
            ) : null}
            <FormActionWithPermissions
              permissions={{entity: 'locationType', action: 'create'}}>
              <Button
                className={classes.addButton}
                onClick={() => showAddEditLocationTypeCard(null)}>
                Add Location Type
              </Button>
            </FormActionWithPermissions>
          </div>
        </div>
        <div className={classes.root}>
          <DroppableTableBody
            isDragDisabled={isSaving}
            className={classes.table}
            onDragEnd={onDragEnd}>
            {locationTypesData.sort(sortByIndex).map((locType, i) => {
              return (
                <div className={classes.listItem} key={`${locType.id}_${i}`}>
                  <LocationTypeItem
                    locationType={locType}
                    position={i}
                    onEdit={() => showAddEditLocationTypeCard(locType)}
                  />
                </div>
              );
            })}
          </DroppableTableBody>
        </div>
      </div>
    </FormContextProvider>
  );
};

export default LocationTypes;
