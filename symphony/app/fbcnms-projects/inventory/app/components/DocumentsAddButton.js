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
  AddImageMutationResponse,
  AddImageMutationVariables,
  ImageEntity,
} from '../mutations/__generated__/AddImageMutation.graphql';
import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {MutationCallbacks} from '../mutations/MutationCallbacks.js';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AddImageMutation from '../mutations/AddImageMutation';
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import FileUploadButton from './FileUpload/FileUploadButton';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import PopoverMenu from '@fbcnms/ui/components/design-system/Select/PopoverMenu';
import React from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Strings from '../common/CommonStrings';
import Text from '@fbcnms/ui/components/design-system/Text';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  uploadCategory: {
    padding: '0px',
  },
  uploadCategoryButton: {
    display: 'block',
    padding: '4px',
    textOverflow: 'ellipsis',
    overflow: 'hidden',
    width: '100%',
  },
};

type Props = {
  entityId: ?string,
  entityType: ImageEntity,
} & WithSnackbarProps &
  WithStyles<typeof styles>;

type State = {
  isMenuOpened: boolean,
};

const FileTypeEnum = {
  IMAGE: 'IMAGE',
  FILE: 'FILE',
};

class DocumentsAddButton extends React.Component<Props, State> {
  static contextType = AppContext;
  context: AppContextType;

  state = {
    isMenuOpened: false,
  };

  render() {
    const {entityId, classes} = this.props;
    const categoriesEnabled = this.context.isFeatureEnabled('file_categories');

    if (!entityId) {
      return null;
    }

    return (
      <FormAction>
        {categoriesEnabled && Strings.documents.categories.length ? (
          <PopoverMenu
            skin="primary"
            menuDockRight={true}
            options={Strings.documents.categories.map(category => ({
              key: category,
              label: (
                <FileUploadButton
                  key={category}
                  onFileUploaded={this.onDocumentUploaded(category)}>
                  {openFileUploadDialog => (
                    <Text
                      className={classes.uploadCategoryButton}
                      variant="body2"
                      onClick={openFileUploadDialog}>
                      {category}
                    </Text>
                  )}
                </FileUploadButton>
              ),
              value: category,
              className: classes.uploadCategory,
            }))}>
            {Strings.documents.uploadButton}
          </PopoverMenu>
        ) : (
          <FileUploadButton onFileUploaded={this.onDocumentUploaded(null)}>
            {openFileUploadDialog => (
              <Button skin="primary" onClick={openFileUploadDialog}>
                {Strings.documents.uploadButton}
              </Button>
            )}
          </FileUploadButton>
        )}
      </FormAction>
    );
  }

  onDocumentUploaded = (category: ?string) => (file, key) => {
    ServerLogger.info(LogEvents.LOCATION_CARD_UPLOAD_FILE_CLICKED);
    if (this.props.entityId == null) {
      return;
    }
    const variables: AddImageMutationVariables = {
      input: {
        entityType: this.props.entityType,
        entityId: this.props.entityId,
        imgKey: key,
        fileName: file.name,
        fileSize: file.size,
        modified: new Date(file.lastModified).toISOString(),
        contentType: file.type,
        category: category,
      },
    };

    const updater = store => {
      // $FlowFixMe (T62907961) Relay flow types
      const newNode = store.getRootField('addImage');
      const fileType = newNode.getValue('fileType');
      // $FlowFixMe (T62907961) Relay flow types
      const entityProxy = store.get(this.props.entityId);
      if (fileType == FileTypeEnum.IMAGE) {
        // $FlowFixMe (T62907961) Relay flow types
        const imageNodes = entityProxy.getLinkedRecords('images') || [];
        // $FlowFixMe (T62907961) Relay flow types
        entityProxy.setLinkedRecords([...imageNodes, newNode], 'images');
      } else {
        // $FlowFixMe (T62907961) Relay flow types
        const fileNodes = entityProxy.getLinkedRecords('files') || [];
        // $FlowFixMe (T62907961) Relay flow types
        entityProxy.setLinkedRecords([...fileNodes, newNode], 'files');
      }
    };

    const callbacks: MutationCallbacks<AddImageMutationResponse> = {
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
      onError: () => {},
    };

    AddImageMutation(variables, callbacks, updater);
  };
}

export default withStyles(styles)(withSnackbar(DocumentsAddButton));
