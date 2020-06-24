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
import ListIcon from '@material-ui/icons/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import ListItemText from '@material-ui/core/ListItemText';
import React from 'react';
import symphony from '@fbcnms/ui/theme/symphony';
import withSuspense from '../../common/withSuspense';
import {makeStyles} from '@material-ui/styles';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {useCallback, useMemo, useState} from 'react';
import {useProjectTemplateNodes} from '../../common/Project';

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
  onSelect: ?(projectTypeId: ?string) => void,
|}>;

function ProjectTypesList(props: Props) {
  const {onSelect} = props;
  const classes = useStyles();

  const projectTypes = useProjectTemplateNodes();

  const [selectedProjectTypeId, setSelectedProjectTypeId] = useState(null);

  const handleListItemClick = useCallback(
    selectedProjectType => {
      const selectedProjectTypeId = selectedProjectType?.id;
      setSelectedProjectTypeId(selectedProjectTypeId);

      if (onSelect != null) {
        onSelect(selectedProjectTypeId);
      }
    },
    [onSelect],
  );

  const listItems = useMemo(
    () =>
      projectTypes
        .slice()
        .sort((projectTypeA, projectTypeB) =>
          sortLexicographically(projectTypeA.name, projectTypeB.name),
        )
        .map(projectType => (
          <>
            <FormActionWithPermissions
              permissions={{
                entity: 'project',
                action: 'create',
                projectTypeId: projectType.id,
              }}>
              <ListItem
                className={classes.listItem}
                button
                key={projectType.id}
                selected={selectedProjectTypeId === projectType.id}
                onClick={() => handleListItemClick(projectType)}>
                <ListItemAvatar className={classes.listAvatar}>
                  <Avatar className={classes.avatar}>
                    <ListIcon />
                  </Avatar>
                </ListItemAvatar>
                <ListItemText primary={projectType.name} />
              </ListItem>
            </FormActionWithPermissions>
          </>
        )),
    [
      classes.avatar,
      classes.listAvatar,
      classes.listItem,
      handleListItemClick,
      projectTypes,
      selectedProjectTypeId,
    ],
  );
  return <List>{listItems}</List>;
}

export default withSuspense(ProjectTypesList);
