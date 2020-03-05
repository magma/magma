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
import type {EntityDocumentsTable_hyperlinks} from './__generated__/EntityDocumentsTable_hyperlinks.graphql';
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
  hyperlinks: EntityDocumentsTable_hyperlinks,
  className?: string,
} & WithAlert &
  WithSnackbarProps;

class EntityDocumentsTable extends React.Component<Props> {
  constructor(props: Props) {
    super(props);
  }

  render() {
    const {files, hyperlinks, className, entityId} = this.props;
    return (
      <div className={className}>
        <DocumentTable
          entityId={entityId}
          files={files}
          hyperlinks={hyperlinks}
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
            // $FlowFixMe (T62907961) Relay flow types
            const deletedNode = store.getRootField('deleteImage');
            // $FlowFixMe (T62907961) Relay flow types
            const proxy = store.get(this.props.entityId);
            const edgeType =
              deletedNode.getValue('fileType') === 'IMAGE' ? 'images' : 'files';

            // $FlowFixMe (T62907961) Relay flow types
            const currNodes = proxy.getLinkedRecords(edgeType);
            // $FlowFixMe (T62907961) Relay flow types
            const newNodes = currNodes.filter(file => {
              return file != deletedNode;
            });
            // $FlowFixMe (T62907961) Relay flow types
            proxy.setLinkedRecords(newNodes, edgeType);
            // $FlowFixMe (T62907961) Relay flow types
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
      hyperlinks: graphql`
        fragment EntityDocumentsTable_hyperlinks on Hyperlink
          @relay(plural: true) {
          ...DocumentTable_hyperlinks
        }
      `,
    }),
  ),
);
