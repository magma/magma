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
import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import FileUpload from './FileUpload';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {withSnackbar} from 'notistack';

type Props = {
  entityId: ?string,
  entityType: ImageEntity,
} & WithSnackbarProps;

type State = {
  isLoadingDocument: boolean,
  anchorEl: ?HTMLElement,
};

const FileTypeEnum = {
  IMAGE: 'IMAGE',
  FILE: 'FILE',
};

// TODO: We should make categories configurable and dynamic.
// we are testing this for now
const categories = [
  'Archivos de Estudios Pre-instalación',
  'Archivos de Contratos',
  'Archivos de TSS',
  'DataFills',
  'ATP',
  'Topología',
  'Archivos Simulación',
  'Reportes de Mantenimiento',
  'Fotos',
];

class DocumentsAddButton extends React.Component<Props, State> {
  static contextType = AppContext;
  context: AppContextType;

  state = {
    isLoadingDocument: false,
    anchorEl: null,
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

    if (!categoriesEnabled) {
      return (
        <FileUpload
          button={
            <Button
              color="primary"
              variant="contained"
              onClick={() =>
                ServerLogger.info(LogEvents.LOCATION_CARD_UPLOAD_FILE_CLICKED)
              }>
              Upload
            </Button>
          }
          onFileUploaded={this.onDocumentUploaded(null)}
          onProgress={() => this.setState({isLoadingDocument: true})}
        />
      );
    }
    return (
      <>
        <Button color="primary" variant="contained" onClick={this.handleClick}>
          Upload
        </Button>
        <Menu
          id="simple-menu"
          anchorEl={this.state.anchorEl}
          keepMounted
          open={Boolean(this.state.anchorEl)}
          onClose={this.handleClose}>
          {categories.map(category => (
            <FileUpload
              key={category}
              button={
                <MenuItem onClick={this.handleClose}>{category}</MenuItem>
              }
              onFileUploaded={this.onDocumentUploaded(category)}
              onProgress={() => this.setState({isLoadingDocument: true})}
            />
          ))}
        </Menu>
      </>
    );
  }

  handleClick = event => {
    ServerLogger.info(LogEvents.LOCATION_CARD_UPLOAD_FILE_CLICKED);
    this.setState({anchorEl: event.currentTarget});
  };

  handleClose = () => {
    this.setState({anchorEl: null});
  };

  onDocumentUploaded = (category: ?string) => (file, key) => {
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
