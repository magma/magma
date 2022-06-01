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
 * The add button at the bottom right of the tables
 */

import * as React from 'react';
import AddIcon from '@material-ui/icons/Add';
import Fab from '@material-ui/core/Fab';
import {Theme} from '@material-ui/core/styles';
import {makeStyles} from '@material-ui/styles';

type Props = {
  label: string;
  onClick: () => void;
};

const useStyles = makeStyles<Theme>(theme => ({
  addButton: {
    position: 'fixed',
    bottom: 0,
    right: 0,
    margin: theme.spacing(2),
  },
}));
export default function TableAddButton({label, onClick, ...props}: Props) {
  const classes = useStyles();
  return (
    <Fab
      {...props}
      className={classes.addButton}
      color="primary"
      onClick={onClick}
      aria-label={label}>
      <AddIcon />
    </Fab>
  );
}
