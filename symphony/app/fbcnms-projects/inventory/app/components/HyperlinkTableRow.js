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
import type {HyperlinkTableRow_hyperlink} from './__generated__/HyperlinkTableRow_hyperlink.graphql';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import HyperlinkTableMenu from './HyperlinkTableMenu';
import InsertLinkIcon from '@material-ui/icons/InsertLink';
import React from 'react';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import symphony from '@fbcnms/ui/theme/symphony';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

const styles = () => ({
  cell: {
    height: '48px',
  },
  nameCell: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    ...symphony.typography.caption,
  },
  thumbnail: {
    marginRight: '20px',
    display: 'flex',
    alignItems: 'center',
  },
  icon: {
    fontSize: '24px',
    width: '32px',
  },
  displayName: {
    ...symphony.typography.caption,
  },
  moreIcon: {
    fill: symphony.palette.D400,
  },
});

type Props = {
  entityId: string,
  hyperlink: HyperlinkTableRow_hyperlink,
} & WithStyles<typeof styles>;

type State = {
  isImageDialogOpen: boolean,
};

class HyperlinkTableRow extends React.Component<Props, State> {
  static contextType = AppContext;
  context: AppContextType;

  render() {
    const categoriesEnabled = this.context.isFeatureEnabled('file_categories');
    const {classes, hyperlink, entityId} = this.props;
    if (hyperlink === null) {
      return null;
    }
    return (
      <TableRow key={hyperlink.id} hover={false}>
        {categoriesEnabled && (
          <TableCell
            padding="none"
            component="th"
            scope="row"
            className={classes.cell}>
            {hyperlink.category}
          </TableCell>
        )}
        <TableCell
          padding="none"
          component="th"
          scope="row"
          className={classes.cell}>
          <a
            className={classes.nameCell}
            href={hyperlink.url}
            target="_blank"
            title={hyperlink.url}>
            <div className={classes.thumbnail}>
              <InsertLinkIcon color="primary" className={classes.icon} />
            </div>
            <div className={classes.displayName}>
              {hyperlink.displayName || hyperlink.url}
            </div>
          </a>
        </TableCell>
        <TableCell className={classes.cell} />
        <TableCell className={classes.cell} />
        <TableCell className={classes.cell} />
        <TableCell
          padding="none"
          className={classes.cell}
          scope="row"
          align="right"
          component="th">
          <HyperlinkTableMenu entityId={entityId} hyperlink={hyperlink} />
        </TableCell>
      </TableRow>
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(HyperlinkTableRow, {
    hyperlink: graphql`
      fragment HyperlinkTableRow_hyperlink on Hyperlink {
        id
        category
        url
        displayName
      }
    `,
  }),
);
