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
  DeleteImageMutationResponse,
  DeleteImageMutationVariables,
  ImageEntity,
} from '../mutations/__generated__/DeleteImageMutation.graphql';
import type {EntityDocumentsTable_files} from './__generated__/EntityDocumentsTable_files.graphql';
import type {MutationCallbacks} from '../mutations/MutationCallbacks.js';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';

import DeleteImageMutation from '../mutations/DeleteImageMutation';
import DocumentTable from './DocumentTable';
import React from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import axios from 'axios';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {DocumentAPIUrls} from '../common/DocumentAPI';
import {createFragmentContainer, graphql} from 'react-relay';
import {withSnackbar} from 'notistack';

type Props = {
  entityType: ImageEntity,
  entityId: string,
  files: EntityDocumentsTable_files,
  className?: string,
} & WithAlert &
  WithSnackbarProps;

class EntityDocumentsTable extends React.Component<Props> {
  constructor(props: Props) {
    super(props);
  }

  render() {
    const {files, className} = this.props;
    return (
      <div className={className}>
        <DocumentTable
          files={files}
          onDocumentDeleted={this.onDocumentDeleted}
        />
      </div>
    );
  }

  onDocumentDeleted = file => {
    this.props
      .confirm(`Are you sure you want to delete "${file.fileName}"?`)
      .then(confirmed => {
        if (confirmed) {
          const variables: DeleteImageMutationVariables = {
            entityType: this.props.entityType,
            entityId: this.props.entityId,
            id: file.id,
          };

          const updater = store => {
            const deletedNode = store.getRootField('deleteImage');
            const proxy = store.get(this.props.entityId);
            const edgeType =
              deletedNode.getValue('fileType') === 'IMAGE' ? 'images' : 'files';

            const currNodes = proxy.getLinkedRecords(edgeType);
            const newNodes = currNodes.filter(file => {
              return file != deletedNode;
            });
            proxy.setLinkedRecords(newNodes, edgeType);
            store.delete(file.id);
            axios.delete(DocumentAPIUrls.delete_url(file.id));
          };

          const callbacks: MutationCallbacks<DeleteImageMutationResponse> = {
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

          DeleteImageMutation(variables, callbacks, updater);
        }
      });
  };
}

export default withAlert(
  withSnackbar(
    createFragmentContainer(EntityDocumentsTable, {
      files: graphql`
        fragment EntityDocumentsTable_files on File @relay(plural: true) {
          ...DocumentTable_files
        }
      `,
    }),
  ),
);
