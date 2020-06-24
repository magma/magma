/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Avatar from '@material-ui/core/Avatar';
import FormActionWithPermissions from '../../common/FormActionWithPermissions';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import ListItemText from '@material-ui/core/ListItemText';
import React, {useCallback, useMemo, useState} from 'react';
import WorkIcon from '@material-ui/icons/Work';
import symphony from '@fbcnms/ui/theme/symphony';
import withSuspense from '../../common/withSuspense';
import {makeStyles} from '@material-ui/styles';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {useWorkOrderTemplateNodes} from '../../common/WorkOrder';

const useStyles = makeStyles(() => ({
  avatar: {
    backgroundColor: symphony.palette.B50,
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
}));

type Props = $ReadOnly<{|
  onSelect: ?(workOrderTypeId: ?string) => void,
|}>;

function WorkOrderTypesList(props: Props) {
  const {onSelect} = props;
  const classes = useStyles();
  const workOrderTypes = useWorkOrderTemplateNodes();
  const [selectedWorkOrderTypeId, setSelectedWorkOrderTypeId] = useState(null);

  const handleListItemClick = useCallback(
    clickedWorkOrderType => {
      const selectedWorkOrderTypeId = clickedWorkOrderType?.id;
      setSelectedWorkOrderTypeId(selectedWorkOrderTypeId);
      if (onSelect) {
        onSelect(selectedWorkOrderTypeId);
      }
    },
    [onSelect],
  );

  const listItems = useMemo(
    () =>
      workOrderTypes
        .slice()
        .sort((workOrderTypeA, workOrderTypeB) =>
          sortLexicographically(workOrderTypeA.name, workOrderTypeB.name),
        )
        .map(workOrderType => (
          <FormActionWithPermissions
            permissions={{
              entity: 'workorder',
              action: 'create',
              workOrderTypeId: workOrderType.id,
            }}>
            <ListItem
              className={classes.listItem}
              button
              key={workOrderType.id}
              selected={selectedWorkOrderTypeId === workOrderType.id}
              onClick={() => handleListItemClick(workOrderType)}>
              <ListItemAvatar className={classes.listAvatar}>
                <Avatar className={classes.avatar}>
                  <WorkIcon />
                </Avatar>
              </ListItemAvatar>
              <ListItemText primary={workOrderType.name} />
            </ListItem>
          </FormActionWithPermissions>
        )),
    [
      classes.avatar,
      classes.listAvatar,
      classes.listItem,
      handleListItemClick,
      selectedWorkOrderTypeId,
      workOrderTypes,
    ],
  );

  return <List>{listItems}</List>;
}

export default withSuspense(WorkOrderTypesList);
