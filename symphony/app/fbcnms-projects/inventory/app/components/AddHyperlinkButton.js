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
  AddHyperlinkInput,
  AddHyperlinkMutationResponse,
  AddHyperlinkMutationVariables,
} from './../mutations/__generated__/AddHyperlinkMutation.graphql';
import type {ImageEntity} from '../mutations/__generated__/AddImageMutation.graphql';
import type {MutationCallbacks} from '../mutations/MutationCallbacks.js';
import type {WithSnackbarProps} from 'notistack';

import * as React from 'react';
import AddHyperlinkDialog from './AddHyperlinkDialog';
import AddHyperlinkMutation from '../mutations/AddHyperlinkMutation';
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import PopoverMenu from '@fbcnms/ui/components/design-system/Select/PopoverMenu';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Strings from '../common/CommonStrings';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {useCallback, useContext, useState} from 'react';
import {withSnackbar} from 'notistack';

type addLinkProps = {
  entityId: string,
  entityType: ImageEntity,
};

type Props = addLinkProps & WithSnackbarProps;

const addNewHyperlink = (input: AddHyperlinkInput, onError: string => void) => {
  const variables: AddHyperlinkMutationVariables = {
    input,
  };

  const updater = store => {
    const newNode = store.getRootField('addHyperlink');
    const entityProxy = store.get(input.entityId);
    const hyperlinkNodes = entityProxy.getLinkedRecords('hyperlinks') || [];
    entityProxy.setLinkedRecords([...hyperlinkNodes, newNode], 'hyperlinks');
  };

  const callbacks: MutationCallbacks<AddHyperlinkMutationResponse> = {
    onCompleted: (_, errors) => {
      if (errors && errors[0]) {
        onError(errors[0].message);
      }
    },
    onError: error => onError(error.message),
  };

  AddHyperlinkMutation(variables, callbacks, updater);
};

const AddHyperlinkButton = (props: Props) => {
  const [addHyperlinkDialogOpened, setAddHyperlinkDialogOpened] = useState(
    false,
  );
  const [dialogKey, setDialogKey] = useState(0);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const appContext = useContext(AppContext);
  const categoriesEnabled = appContext.isFeatureEnabled('file_categories');

  const openDialog = useCallback((category: ?string) => {
    setSelectedCategory(category);
    setDialogKey(key => key + 1);
    setAddHyperlinkDialogOpened(true);
    ServerLogger.info(LogEvents.LOCATION_CARD_ADD_HYPERLINK_CLICKED);
  }, []);

  const callAddNewHyperlink = useCallback(
    (url: string, displayName: ?string) => {
      const {entityId, entityType, enqueueSnackbar} = props;
      const onError = errorMessage => {
        enqueueSnackbar(errorMessage, {
          children: key => (
            <SnackbarItem id={key} message={errorMessage} variant="error" />
          ),
        });
      };
      const input = {
        entityId,
        entityType,
        url,
        displayName,
        category: selectedCategory,
      };
      addNewHyperlink(input, onError);
    },
    [props, selectedCategory],
  );

  return (
    <>
      {categoriesEnabled && Strings.documents.categories.length ? (
        <PopoverMenu
          skin="gray"
          menuDockRight={true}
          options={Strings.documents.categories.map(category => ({
            label: category,
            value: category,
          }))}
          onChange={openDialog}>
          {Strings.documents.addLinkButton}
        </PopoverMenu>
      ) : (
        <Button skin="gray" onClick={openDialog}>
          {Strings.documents.addLinkButton}
        </Button>
      )}
      <AddHyperlinkDialog
        key={dialogKey}
        isOpened={addHyperlinkDialogOpened}
        onAdd={callAddNewHyperlink.bind(this)}
        onClose={() => setAddHyperlinkDialogOpened(false)}
        targetCategory={selectedCategory}
      />
    </>
  );
};

export default withSnackbar(AddHyperlinkButton);
