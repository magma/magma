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
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Collapse from '@material-ui/core/Collapse';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import RootRef from '@material-ui/core/RootRef';
import Tooltip from '@material-ui/core/Tooltip';
import classnames from 'classnames';
import {makeStyles} from '@material-ui/styles';

export type Props = {
  RequiredFields: React.Node,
  OptionalFields?: React.Node,
  ...CommonProps,
};

export type CommonProps = {|
  onDelete: () => void,
  onReset: () => void,
  isNew: boolean,
|};

/**
 * This component is designed to be composed by a config editor such as the
 * SlackConfigEditor. These are the props required by the editor and shared with
 * ConfigEditor.
 */
export type EditorProps<TConfig> = {
  ...CommonProps,
  onUpdate: ($Shape<TConfig> | TConfig) => void,
  config: TConfig,
};

const useStyles = makeStyles(theme => ({
  expand: {
    transform: 'rotate(0deg)',
    marginLeft: 'auto',
    transition: theme.transitions.create('transform', {
      duration: theme.transitions.duration.shortest,
    }),
  },
  expandOpen: {
    transform: 'rotate(180deg)',
  },
}));

export default function ConfigEditor({
  onDelete,
  onReset,
  RequiredFields,
  OptionalFields,
  isNew,
  ...props
}: Props) {
  const classes = useStyles();

  const [optionalFieldsExpanded, setOptionalFieldsExpanded] = React.useState(
    false,
  );
  const handleExpandClick = React.useCallback(
    () => setOptionalFieldsExpanded(x => !x),
    [setOptionalFieldsExpanded],
  );
  return (
    <Card {...props}>
      <CardContent>
        <Grid container justify="flex-end">
          <EditorMenuButton
            onReset={onReset}
            onDelete={onDelete}
            isNew={isNew}
          />
        </Grid>
        <Grid container spacing={2} direction="column" wrap="nowrap">
          {RequiredFields}
        </Grid>
      </CardContent>
      {OptionalFields && (
        <>
          <CardActions disableSpacing>
            <Tooltip title="Advanced" placement="right">
              <IconButton
                className={classnames(classes.expand, {
                  [classes.expandOpen]: optionalFieldsExpanded,
                })}
                onClick={handleExpandClick}
                aria-expanded={optionalFieldsExpanded}
                aria-label="optional fields">
                <ExpandMoreIcon />
              </IconButton>
            </Tooltip>
          </CardActions>
          <Collapse in={optionalFieldsExpanded} timeout="auto" unmountOnExit>
            <CardContent>{OptionalFields}</CardContent>
          </Collapse>
        </>
      )}
    </Card>
  );
}

// menu button for top right of card
function EditorMenuButton({
  onReset,
  onDelete,
  isNew,
}: {
  isNew: boolean,
  onReset: () => void,
  onDelete: () => void,
}) {
  const iconRef = React.useRef<?HTMLElement>();
  const [isMenuOpen, setMenuOpen] = React.useState(false);
  return (
    <>
      <RootRef rootRef={iconRef}>
        <IconButton
          aria-label="editor-menu"
          size="small"
          edge="end"
          onClick={() => setMenuOpen(true)}>
          <MoreVertIcon />
        </IconButton>
      </RootRef>
      <Menu
        anchorEl={iconRef.current}
        open={isMenuOpen}
        onClose={() => setMenuOpen(false)}>
        {!isNew && (
          <MenuItem
            onClick={() => {
              onReset();
              setMenuOpen(false);
            }}>
            Reset
          </MenuItem>
        )}
        <MenuItem
          onClick={() => {
            onDelete();
            setMenuOpen(false);
          }}>
          Delete
        </MenuItem>
      </Menu>
    </>
  );
}
