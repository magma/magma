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

import type {Organization} from '../organizations/Organizations';

import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormLabel from '@material-ui/core/FormLabel';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import LoadingFillerBackdrop from '../../components/LoadingFillerBackdrop';
import React from 'react';
import axios from 'axios';

import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useAxios} from '../../hooks';
import {useEffect, useState} from 'react';

const useStyles = makeStyles(_ => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

export type FeatureFlag = {
  id: number,
  title: string,
  config: {
    id: string,
    enabled: boolean,
  },
  enabledByDefault: boolean,
};

type Props = {
  onClose: () => void,
  onSave: FeatureFlag => void,
  featureFlag: FeatureFlag,
};

type FeatureFlagStatus = {
  [orgName: string]: 'enabled' | 'disabled' | 'default',
};

export default function (props: Props) {
  const classes = useStyles();
  const {error, isLoading, response} = useAxios({
    method: 'get',
    url: '/host/organization/async',
  });
  const [featureFlagStatus, setFeatureFlagStatus] = useState<FeatureFlagStatus>(
    {},
  );

  useEffect(() => {
    if (response) {
      setFeatureFlagStatus(
        buildStatusObject(response.data.organizations, props.featureFlag),
      );
    }
  }, [response, props.featureFlag]);

  if (error || isLoading || !response) {
    return <LoadingFillerBackdrop />;
  }

  const onSave = async () => {
    const response = await axios.post(
      `/host/feature/async/${props.featureFlag.id}`,
      createPayload(featureFlagStatus, props.featureFlag),
    );
    props.onSave(response.data);
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>{props.featureFlag.title}</DialogTitle>
      <DialogContent>
        {error && <FormLabel error>{error}</FormLabel>}
        {response.data.organizations.map(org => (
          <FormControlLabel
            key={org.name}
            className={classes.input}
            control={
              <Checkbox
                checked={featureFlagStatus[org.name] === 'enabled'}
                indeterminate={featureFlagStatus[org.name] === 'default'}
                color="primary"
                onChange={() => {
                  let nextStatus = 'default';
                  if (featureFlagStatus[org.name] === 'default') {
                    nextStatus = 'enabled';
                  } else if (featureFlagStatus[org.name] === 'enabled') {
                    nextStatus = 'disabled';
                  }
                  setFeatureFlagStatus({
                    ...featureFlagStatus,
                    [org.name]: nextStatus,
                  });
                }}
              />
            }
            label={org.name}
          />
        ))}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}

function buildStatusObject(
  organizations: Organization[],
  featureFlag: FeatureFlag,
): FeatureFlagStatus {
  const status = {};
  organizations.forEach(org => {
    status[org.name] = 'default';
    if (featureFlag.config[org.name]?.enabled) {
      status[org.name] = 'enabled';
    } else if (featureFlag.config[org.name]?.enabled === false) {
      status[org.name] = 'disabled';
    }
  });
  return status;
}

function createPayload(
  status: FeatureFlagStatus,
  featureFlag: FeatureFlag,
): {
  toCreate: {organization: string, enabled: boolean}[],
  toDelete: {[id: number]: null},
  toUpdate: {[id: number]: boolean},
} {
  const toCreate = [];
  const toDelete = {};
  const toUpdate = {};

  Object.keys(status).forEach(orgName => {
    const originalConfig = featureFlag.config[orgName];
    if (status[orgName] !== 'default' && !originalConfig) {
      toCreate.push({
        organization: orgName,
        enabled: status[orgName] === 'enabled',
      });
    } else if (status[orgName] === 'default' && originalConfig) {
      toDelete[originalConfig.id] = null;
    } else if (originalConfig) {
      toUpdate[originalConfig.id] = {
        enabled: status[orgName] === 'enabled',
      };
    }
  });

  return {toCreate, toDelete, toUpdate};
}
