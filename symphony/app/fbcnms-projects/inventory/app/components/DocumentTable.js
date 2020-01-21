/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {DocumentTable_files} from './__generated__/DocumentTable_files.graphql';
import type {DocumentTable_hyperlinks} from './__generated__/DocumentTable_hyperlinks.graphql';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import FileAttachment from './FileAttachment';
import HyperlinkTableRow from './HyperlinkTableRow';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import {createFragmentContainer, graphql} from 'react-relay';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  table: {
    minWidth: 70,
    marginBottom: '12px',
  },
});

type Props = WithStyles<typeof styles> & {
  files: DocumentTable_files,
  hyperlinks: DocumentTable_hyperlinks,
  onDocumentDeleted: (file: DocumentTable_files) => void,
};

class DocumentTable extends React.Component<Props> {
  static contextType = AppContext;
  context: AppContextType;

  render() {
    const {classes, onDocumentDeleted} = this.props;
    const files = [...this.props.files].filter(Boolean);
    const hyperlinks = [...this.props.hyperlinks].filter(Boolean);
    const categoriesEnabled = this.context.isFeatureEnabled('file_categories');
    let sortedFiles = files;
    if (categoriesEnabled) {
      sortedFiles = files.sort((fileA, fileB) =>
        sortLexicographically(fileA.category ?? '', fileB.category ?? ''),
      );
    } else {
      sortedFiles = files.sort((fileA, fileB) =>
        sortLexicographically(fileA.fileName, fileB.fileName),
      );
    }
    return files.length > 0 ? (
      <Table className={classes.table}>
        <TableBody>
          {sortedFiles.map(file => (
            <FileAttachment
              key={file.id}
              file={file}
              onDocumentDeleted={onDocumentDeleted}
            />
          ))}
          {hyperlinks.map(hyperlink => (
            <HyperlinkTableRow key={hyperlink.id} hyperlink={hyperlink} />
          ))}
        </TableBody>
      </Table>
    ) : null;
  }
}

export default withStyles(styles)(
  createFragmentContainer(DocumentTable, {
    files: graphql`
      fragment DocumentTable_files on File @relay(plural: true) {
        id
        fileName
        category
        ...FileAttachment_file
      }
    `,
    hyperlinks: graphql`
      fragment DocumentTable_hyperlinks on Hyperlink @relay(plural: true) {
        id
        category
        ...HyperlinkTableRow_hyperlink
      }
    `,
  }),
);
