/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ProjectTypeWorkOrderTemplatesPanel_workOrderTypes} from './__generated__/ProjectTypeWorkOrderTemplatesPanel_workOrderTypes.graphql';

import Avatar from '@material-ui/core/Avatar';
import Button from '@fbcnms/ui/components/design-system/Button';
import CheckIcon from '@material-ui/icons/Check';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import React, {useState} from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import WorkIcon from '@material-ui/icons/Work';

import {makeStyles} from '@material-ui/styles';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';

const useStyles = makeStyles(theme => ({
  root: {
    position: 'relative',
  },
  avatar: {
    backgroundColor: '#e4f2ff',
  },
  workOrderName: {
    fontSize: '16px',
    lineHeight: '24px',
    color: theme.palette.blueGrayDark,
    fontWeight: 500,
    flexGrow: 1,
  },
  dialogTitle: {
    padding: '24px',
    paddingBottom: '16px',
  },
  dialogTitleText: {
    fontSize: '20px',
    lineHeight: '24px',
    color: theme.palette.blueGrayDark,
    fontWeight: 500,
  },
  dialogContent: {
    padding: 0,
    height: '400px',
    overflowY: 'scroll',
  },
  list: {
    paddingTop: 0,
    paddingBottom: 0,
  },
  listItem: {
    paddingLeft: '24px',
    paddingRight: '24px',
  },
  listAvatar: {
    minWidth: '52px',
  },
  dialogActions: {
    position: 'absolute',
    padding: '24px',
    bottom: 0,
    display: 'flex',
    justifyContent: 'flex-end',
    width: '100%',
    backgroundColor: 'rgba(255, 255, 255, 0.9)',
    zIndex: 2,
  },
  spacer: {
    height: '80px',
  },
  checkIcon: {
    color: theme.palette.primary.main,
    fontSize: '18px',
  },
}));

type Props = {
  initialSelectedWorkOrderTypeIds: Array<string>,
  workOrderTypes: ProjectTypeWorkOrderTemplatesPanel_workOrderTypes,
  open: boolean,
  onClose: () => void,
  onSaveClicked: (workOrderTypeIds: Array<string>) => void,
};

const ProjectTypeSelectWorkOrdersDialog = ({
  open,
  onClose,
  initialSelectedWorkOrderTypeIds,
  workOrderTypes,
  onSaveClicked,
}: Props) => {
  const classes = useStyles();
  const [selectedWorkOrderTypeIds, setSelectedWorkOrderTypeIds] = useState(
    initialSelectedWorkOrderTypeIds,
  );

  return (
    <Dialog
      onClose={onClose}
      open={open}
      fullWidth={true}
      maxWidth="sm"
      className={classes.root}>
      <DialogTitle className={classes.dialogTitle}>
        <Text className={classes.dialogTitleText}>Select Work Orders</Text>
      </DialogTitle>
      <DialogContent className={classes.dialogContent}>
        <List className={classes.list}>
          {workOrderTypes
            .filter(Boolean)
            .sort((wotA, wotB) => sortLexicographically(wotA.name, wotB.name))
            .map(workOrder => (
              <ListItem
                className={classes.listItem}
                button
                onClick={() =>
                  setSelectedWorkOrderTypeIds(
                    selectedWorkOrderTypeIds.includes(workOrder.id)
                      ? selectedWorkOrderTypeIds.filter(
                          existingId => existingId !== workOrder.id,
                        )
                      : [...selectedWorkOrderTypeIds, workOrder.id],
                  )
                }
                key={workOrder.id}>
                <ListItemAvatar className={classes.listAvatar}>
                  <Avatar className={classes.avatar}>
                    <WorkIcon />
                  </Avatar>
                </ListItemAvatar>
                <Text className={classes.workOrderName}>{workOrder.name}</Text>
                {selectedWorkOrderTypeIds.includes(workOrder.id) && (
                  <CheckIcon className={classes.checkIcon} />
                )}
              </ListItem>
            ))}
        </List>
        <div className={classes.spacer} />
      </DialogContent>
      <DialogActions className={classes.dialogActions}>
        <Button onClick={onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={() => onSaveClicked(selectedWorkOrderTypeIds)}>
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ProjectTypeSelectWorkOrdersDialog;
