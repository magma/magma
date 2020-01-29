/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import ListAltIcon from '@material-ui/icons/ListAlt';
import MapButtonGroup from '@fbcnms/ui/components/map/MapButtonGroup';
import MapIcon from '@material-ui/icons/Map';
import Text from '@fbcnms/ui/components/design-system/Text';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  bar: {
    display: 'flex',
    flexDirection: 'column',
    padding: '16px 24px',
    paddingBottom: '0',
  },
  barRow: {
    display: 'flex',
    flexDirection: 'row',
    '&:not(:first-child)': {
      paddingTop: '8px',
    },
  },
  expandedBarPart: {
    flexGrow: '1',
  },
  groupButtons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  buttonContent: {
    paddingTop: '4px',
  },
  actionButtons: {
    display: 'flex',
    flexDirection: 'row',
  },
  actionButton: {
    '&:not(:first-child)': {
      paddingleft: '8px',
    },
  },
});

export type DisplayOptionTypes = 'table' | 'map';
export const DisplayOptions = {
  table: 'table',
  map: 'map',
};

type ActionButtonProps = {
  title: string,
  action: () => void,
};

type Props = {
  title: string,
  onViewToggleClicked?: (id: string) => void,
  actionButtons?: Array<ActionButtonProps>,
  searchBar?: React.Node,
};

const InventoryViewHeader = (props: Props) => {
  const classes = useStyles();

  return (
    <div className={classes.bar}>
      <div className={classes.barRow}>
        <Text className={classes.expandedBarPart} variant="h6">
          {props.title}
        </Text>
        {!!props.onViewToggleClicked && (
          <MapButtonGroup
            onIconClicked={props.onViewToggleClicked}
            buttons={[
              {
                item: <ListAltIcon className={classes.buttonContent} />,
                id: DisplayOptions.table,
              },
              {
                item: <MapIcon className={classes.buttonContent} />,
                id: DisplayOptions.map,
              },
            ]}
          />
        )}
      </div>
      <div className={classes.barRow}>
        <div className={classes.expandedBarPart}>{props.searchBar}</div>
        {!!props.actionButtons && (
          <div className={classes.actionButtons}>
            {props.actionButtons.map(actionButton => (
              <Button
                className={classes.actionButton}
                onClick={actionButton.action}>
                {actionButton.title}
              </Button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default InventoryViewHeader;
