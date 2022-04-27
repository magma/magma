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
 */

import React from 'react';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  header: {
    backgroundColor: theme.palette.common.white,
    borderBottom: `1px solid ${theme.palette.divider}`,
  },
  paper: {
    margin: theme.spacing(3),
  },
  searchBar: {
    marginLeft: theme.spacing(1),
  },
}));

export default function Logs() {
  const classes = useStyles();
  return (
    <>
      <div className={classes.header}>
        <TextField
          id="outlined-search"
          placeholder="Search logs"
          type="search"
          margin="normal"
          variant="outlined"
          className={classes.searchBar}
        />
      </div>
      <div className={classes.paper} />
    </>
  );
}
