/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DocumentMenu_document} from './__generated__/DocumentMenu_document.graphql';

import React from 'react';
import TableRowOptionsButton from './TableRowOptionsButton';
import fbt from 'fbt';
import nullthrows from '@fbcnms/util/nullthrows';
import {DocumentAPIUrls} from '../common/DocumentAPI';
import {createFragmentContainer, graphql} from 'react-relay';

type Props = {
  document: DocumentMenu_document,
  onDocumentDeleted: (document: DocumentMenu_document) => void,
  onDialogOpen: () => void,
};

class DocumentMenu extends React.Component<Props> {
  downloadFileRef: {
    current: null | HTMLAnchorElement,
  } = React.createRef<HTMLAnchorElement>();

  handleDownload = () => {
    if (this.downloadFileRef.current != null) {
      this.downloadFileRef.current.click();
    }
  };

  handlePreview = () => {
    this.props.onDialogOpen();
  };

  handleDelete = () => {
    this.props.onDocumentDeleted(this.props.document);
  };

  render() {
    const {document} = this.props;
    const storeKey = nullthrows(document.storeKey);
    const menuOptions = [
      {
        onClick: this.handlePreview.bind(this),
        caption: fbt(
          'Preview',
          'Caption for menu option for showing image in preview mode',
        ),
        ignorePermissions: true,
      },
      {
        onClick: this.handleDelete.bind(this),
        caption: fbt(
          'Delete',
          'Caption for menu option for deleting a file from files table',
        ),
      },
      {
        onClick: this.handleDownload.bind(this),
        caption: fbt(
          'Download',
          'Caption for menu option for downloading a file from files table',
        ),
        ignorePermissions: true,
      },
    ];
    return (
      <>
        <TableRowOptionsButton options={menuOptions} />
        <a
          href={DocumentAPIUrls.download_url(storeKey, document.fileName)}
          ref={this.downloadFileRef}
          style={{display: 'none'}}
          download
        />
      </>
    );
  }
}

export default createFragmentContainer(DocumentMenu, {
  document: graphql`
    fragment DocumentMenu_document on File {
      id
      fileName
      storeKey
      fileType
    }
  `,
});
