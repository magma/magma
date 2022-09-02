/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import * as React from 'react';
import ButtonBase from '@mui/material/ButtonBase';
import Collapse from '@mui/material/Collapse';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import Typography from '@mui/material/Typography';
import {Theme} from '@mui/material/styles';
import {colors} from '../../../../../theme/default';
import {makeStyles} from '@mui/styles';

export type Props = {
  RequiredFields: React.ReactNode;
  OptionalFields?: React.ReactNode;
} & CommonProps;

export type CommonProps = {
  onDelete: () => void;
  onReset: () => void;
  isNew: boolean;
};

/**
 * This component is designed to be composed by a config editor such as the
 * SlackConfigEditor. These are the props required by the editor and shared with
 * ConfigEditor.
 */
export type EditorProps<TConfig> = {
  onUpdate: (config: Partial<TConfig> | TConfig) => void;
  config: TConfig;
} & CommonProps;

const useStyles = makeStyles<Theme>(theme => ({
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
  configEditor: {
    '&:not(:last-of-type)': {
      borderBottom: `1px solid ${colors.primary.mercury}`,
      paddingBottom: theme.spacing(4),
    },
  },
}));

export default function ConfigEditor({
  onDelete,
  onReset,
  RequiredFields,
  OptionalFields,
  isNew,
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
    <Grid
      className={classes.configEditor}
      container
      item
      justifyContent="space-between"
      xs={12}
      alignItems="flex-start">
      <Grid item container spacing={2} direction="column" wrap="nowrap" xs={11}>
        <Grid item xs={12}>
          {RequiredFields}
        </Grid>
        {OptionalFields && (
          <>
            <Grid item>
              <ButtonBase onClick={handleExpandClick} disableTouchRipple>
                <Typography color="primary" variant="body2">
                  {!optionalFieldsExpanded ? 'Show' : 'Hide'} advanced options
                </Typography>
              </ButtonBase>
            </Grid>
            <Grid item xs={12}>
              <Collapse
                in={optionalFieldsExpanded}
                timeout="auto"
                unmountOnExit>
                {OptionalFields}
              </Collapse>
            </Grid>
          </>
        )}
      </Grid>
      <Grid item xs={1} container justifyContent="flex-end">
        <Grid item>
          <EditorMenuButton
            onReset={onReset}
            onDelete={onDelete}
            isNew={isNew}
          />
        </Grid>
      </Grid>
    </Grid>
  );
}

// menu button for top right of card
function EditorMenuButton({
  onReset,
  onDelete,
  isNew,
}: {
  isNew: boolean;
  onReset: () => void;
  onDelete: () => void;
}) {
  const iconRef = React.useRef<HTMLButtonElement>(null);
  const [isMenuOpen, setMenuOpen] = React.useState(false);

  return (
    <>
      <IconButton
        ref={iconRef}
        aria-label="editor-menu"
        edge="end"
        onClick={() => setMenuOpen(true)}
        size="large">
        <MoreVertIcon />
      </IconButton>
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
