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

import AddImageMutation from '../mutations/AddImageMutation';
import AppContext from '@fbcnms/ui/context/AppContext';
import CircularProgress from '@material-ui/core/CircularProgress';
import FileUpload from './FileUpload';
import PopoverMenu from '@fbcnms/ui/components/design-system/Select/PopoverMenu';
import React from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Strings from '../common/CommonStrings';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {withSnackbar} from 'notistack';

type Props = {
  entityId: ?string,
  entityType: ImageEntity,
} & WithSnackbarProps;

type State = {
  isLoadingDocument: boolean,
  isMenuOpened: boolean,
};

const FileTypeEnum = {
  IMAGE: 'IMAGE',
  FILE: 'FILE',
};

class DocumentsAddButton extends React.Component<Props, State> {
  static contextType = AppContext;
  context: AppContextType;
  menuButtonRef = React.createRef();

  state = {
    isLoadingDocument: false,
    isMenuOpened: false,
  };

  render() {
    const {entityId} = this.props;
    const categoriesEnabled = this.context.isFeatureEnabled('file_categories');

    if (!entityId) {
      return null;
    }

    if (this.state.isLoadingDocument) {
      return <CircularProgress />;
    }

    return (
      <>
        {categoriesEnabled && Strings.documents.categories.length ? (
          <PopoverMenu
            skin="primary"
            menuDockRight={true}
            options={Strings.documents.categories.map(category => ({
              label: (
                <FileUpload
                  key={category}
                  button={category}
                  onFileUploaded={this.onDocumentUploaded(category)}
                  onProgress={() => this.setState({isLoadingDocument: true})}
                />
              ),
              value: category,
            }))}>
            {Strings.documents.uploadButton}
          </PopoverMenu>
        ) : (
          <FileUpload
            button={Strings.documents.uploadButton}
            onFileUploaded={this.onDocumentUploaded(null)}
            onProgress={() => this.setState({isLoadingDocument: true})}
          />
        )}
      </>
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
      this.setState({isLoadingDocument: false});
      const newNode = store.getRootField('addImage');
      const fileType = newNode.getValue('fileType');
      const entityProxy = store.get(this.props.entityId);
      if (fileType == FileTypeEnum.IMAGE) {
        const imageNodes = entityProxy.getLinkedRecords('images') || [];
        entityProxy.setLinkedRecords([...imageNodes, newNode], 'images');
      } else {
        const fileNodes = entityProxy.getLinkedRecords('files') || [];
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

export default withSnackbar(DocumentsAddButton);
