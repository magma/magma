/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import Avatar from '@material-ui/core/Avatar';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Text from '@fbcnms/ui/components/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import classNames from 'classnames';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  icon: React.Element<any>,
  instanceCount: number,
  instanceNamePlural: string,
  instanceNameSingular: string,
  name: string,
  onDelete: () => void,
  onEdit?: ?() => void,
  allowDelete?: ?boolean,
} & WithStyles<typeof styles>;

const styles = theme => ({
  inline: {
    display: 'flex',
    alignItems: 'center',
    flexGrow: 1,
  },
  root: {
    flexGrow: 1,
  },
  icon: {
    backgroundColor: theme.palette.grey[50],
    color: theme.palette.grey[500],
    marginRight: 15,
  },
  iconButton: {
    display: 'block',
    marginLeft: 'auto',
  },
  boldText: {
    fontWeight: 'bold',
  },
  text: {
    color: '#4d4d4e',
  },
  counter: {
    marginRight: 'auto',
  },
  actionButtons: {
    display: 'flex',
    flexDirection: 'row-reverse',
    alignItems: 'center',
    flexGrow: 1,
  },
  checkbox: {
    margin: '12px',
  },
});

class ConfigureExpansionPanel extends React.Component<Props> {
  render() {
    const {
      classes,
      icon,
      instanceCount,
      instanceNamePlural,
      instanceNameSingular,
      name,
      onEdit,
    } = this.props;
    return (
      <Grid
        container
        className={classes.root}
        direction="row"
        justify="space-between"
        alignItems="center">
        <Grid item xs>
          <div className={classes.inline}>
            <Avatar className={classes.icon}>{icon}</Avatar>
            <Text
              className={classNames(classes.text, classes.boldText)}
              variant="subtitle1">
              {name}
            </Text>
          </div>
        </Grid>
        <Grid item xs>
          <Text
            className={classNames(classes.text, classes.counter)}
            variant="body2">
            {`${instanceCount.toLocaleString()}
                  ${
                    instanceCount == 1
                      ? instanceNameSingular
                      : instanceNamePlural
                  } of this type`}
          </Text>
        </Grid>
        <Grid item xs>
          <div className={classes.actionButtons}>
            {this.deleteButton()}
            {onEdit && (
              <IconButton onClick={onEdit} color="primary">
                <EditIcon />
              </IconButton>
            )}
          </div>
        </Grid>
      </Grid>
    );
  }

  deleteButton = () => {
    const {classes, instanceCount, allowDelete} = this.props;
    const disabled =
      allowDelete !== undefined && allowDelete !== null
        ? !allowDelete
        : instanceCount > 0;
    const deleteButton = (
      <IconButton
        disabled={disabled}
        onClick={this.props.onDelete}
        color="primary"
        className={classes.iconButton}>
        <DeleteIcon />
      </IconButton>
    );
    const tooltip = `Cannot delete a type that is in use`;
    return (
      <div>
        {disabled ? (
          <Tooltip title={tooltip}>
            <div>{deleteButton}</div>
          </Tooltip>
        ) : (
          deleteButton
        )}
      </div>
    );
  };
}

export default withStyles(styles)(ConfigureExpansionPanel);
