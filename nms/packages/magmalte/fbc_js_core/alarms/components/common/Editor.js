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
 *
 * @flow strict-local
 * @format
 *
 * Wrappper component for editors such as AddEditRule, AddEditReceiver, etc
 */

import * as React from 'react';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import {makeStyles} from '@material-ui/styles';

export type Props = {
  children: React.Node,
  onExit: () => void,
  onSave: () => Promise<void> | void,
  isNew: boolean,
  title?: string,
  description?: string,
};
const useStyles = makeStyles(theme => ({
  gridContainer: {
    flexGrow: 1,
  },
  editingSpace: {
    height: '100%',
    padding: theme.spacing(3),
  },
}));
export default function Editor({
  children,
  isNew,
  onExit,
  onSave,
  title,
  description,
  ...props
}: Props) {
  const classes = useStyles();

  return (
    <Grid {...props} className={classes.gridContainer} container spacing={0}>
      <Grid className={classes.editingSpace} item xs>
        <form
          onSubmit={e => {
            e.preventDefault();
            onSave();
          }}
          data-testid="editor-form">
          <Grid container spacing={4} direction="column" wrap="nowrap">
            <Grid container item wrap="nowrap" xs={12}>
              <Grid item xs={6}>
                <Typography variant="h5" noWrap>
                  {title}
                </Typography>
                <Typography variant="body2" color="textSecondary" noWrap>
                  {description}
                </Typography>
              </Grid>
              <Grid
                container
                item
                spacing={1}
                xs={6}
                justify="flex-end"
                alignItems="center">
                <Grid item>
                  <Button variant="outlined" onClick={() => onExit()}>
                    Close
                  </Button>
                </Grid>
                <Grid item>
                  <Button
                    variant="contained"
                    color="primary"
                    type="submit"
                    data-testid="editor-submit-button">
                    {isNew ? 'Add' : 'Save'}
                  </Button>
                </Grid>
              </Grid>
            </Grid>
            <Grid container item spacing={3}>
              <Grid
                container
                item
                direction="column"
                spacing={2}
                wrap="nowrap"
                xs={12}>
                {children}
              </Grid>
            </Grid>
          </Grid>
        </form>
      </Grid>
    </Grid>
  );
}
