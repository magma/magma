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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {FeatureFlag} from './FeatureFlagsDialog';

import EditIcon from '@material-ui/icons/Edit';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import FeatureFlagsDialog from './FeatureFlagsDialog';
import IconButton from '@material-ui/core/IconButton';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../../components/LoadingFiller';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import axios from 'axios';

// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import renderList from '../../../app/util/renderList';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';

const useStyles = makeStyles(_ => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: '10px',
  },
}));

function EditFeatureFlagsDialog(props: {
  featureFlags: FeatureFlag[],
  setFeatureFlags: (featureFlags: FeatureFlag[]) => void,
}) {
  const params = useParams();
  const navigate = useNavigate();

  return (
    <FeatureFlagsDialog
      featureFlag={nullthrows(props.featureFlags.find(f => f.id === params.id))}
      onClose={() => navigate('..')}
      onSave={flag => {
        const newFeatureFlags = [...props.featureFlags];
        for (let i = 0; i < newFeatureFlags.length; i++) {
          if (newFeatureFlags[i].id === flag.id) {
            newFeatureFlags[i] = flag;
            break;
          }
        }
        props.setFeatureFlags(newFeatureFlags);
        navigate('..');
      }}
    />
  );
}

export default function Features() {
  const classes = useStyles();
  const navigate = useNavigate();

  const [featureFlags, setFeatureFlags] = useState<?(FeatureFlag[])>(null);
  useEffect(() => {
    axios.get('/host/feature/async').then(({data}) => setFeatureFlags(data));
  }, []);

  if (!featureFlags) {
    return <LoadingFiller />;
  }

  const rows = featureFlags.map(row => (
    <TableRow key={row.id}>
      <TableCell>{row.title}</TableCell>
      <TableCell>{row.enabledByDefault ? 'Yes' : 'No'}</TableCell>
      <TableCell>
        {renderList(
          // $FlowIgnore
          Object.keys(row.config).filter(org => row.config[org].enabled),
        )}
      </TableCell>
      <TableCell>
        {renderList(
          // $FlowIgnore
          Object.keys(row.config).filter(org => !row.config[org].enabled),
        )}
      </TableCell>
      <TableCell>
        <IconButton onClick={() => navigate(`edit/${row.id}`)}>
          <EditIcon />
        </IconButton>
      </TableCell>
    </TableRow>
  ));

  return (
    <div className={classes.paper}>
      <Paper elevation={2}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Feature</TableCell>
              <TableCell>Enabled By Default</TableCell>
              <TableCell>Enabled For</TableCell>
              <TableCell>Disabled For</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
        </Table>
      </Paper>
      <Routes>
        <Route
          path="/edit/:id"
          element={
            <EditFeatureFlagsDialog
              featureFlags={featureFlags}
              setFeatureFlags={setFeatureFlags}
            />
          }
        />
      </Routes>
    </div>
  );
}
