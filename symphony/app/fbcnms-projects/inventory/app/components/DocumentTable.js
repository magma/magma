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

const getHyperlinkSortingValue = (hyperlink, categoriesEnabled) => {
  return `${(categoriesEnabled && hyperlink.category) ||
    ''}${hyperlink.displayName || hyperlink.url}`;
};
const getFileSortingValue = (file, categoriesEnabled) => {
  return `${(categoriesEnabled && file.category) || ''}${file.fileName}`;
};

class DocumentTable extends React.Component<Props> {
  static contextType = AppContext;
  context: AppContextType;

  render() {
    const {classes, onDocumentDeleted} = this.props;
    const categoriesEnabled = this.context.isFeatureEnabled('file_categories');
    const files = this.props.files.map(file => ({
      ...file,
      isFile: true,
      sortingValue: getFileSortingValue(file, categoriesEnabled),
    }));
    const hyperlinks = this.props.hyperlinks.map(hyperlink => ({
      ...hyperlink,
      isHyperlink: true,
      sortingValue: getHyperlinkSortingValue(hyperlink, categoriesEnabled),
    }));
    const allDocuments = [...files, ...hyperlinks].sort((docA, docB) =>
      sortLexicographically(docA.sortingValue, docB.sortingValue),
    );
    return allDocuments.length > 0 ? (
      <Table className={classes.table}>
        <TableBody>
          {allDocuments.map(
            doc =>
              (doc.isFile && (
                <FileAttachment
                  key={doc.id}
                  file={doc}
                  onDocumentDeleted={onDocumentDeleted}
                />
              )) ||
              (doc.isHyperlink && (
                <HyperlinkTableRow key={doc.id} hyperlink={doc} />
              )),
          )}
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
        url
        displayName
        ...HyperlinkTableRow_hyperlink
      }
    `,
  }),
);
