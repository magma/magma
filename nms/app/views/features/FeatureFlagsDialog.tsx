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
import type {
  Organization,
  OrganizationName,
} from '../organizations/Organizations';

import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormLabel from '@material-ui/core/FormLabel';
import LoadingFillerBackdrop from '../../components/LoadingFillerBackdrop';
import React from 'react';
import axios from 'axios';

import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {useAxios} from '../../hooks';
import {useEffect, useState} from 'react';

const useStyles = makeStyles({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
});

export type FeatureFlag = {
  id: string;
  title: string;
  config: Record<
    OrganizationName,
    {
      id: string;
      enabled: boolean;
    }
  >;
  enabledByDefault: boolean;
};

type Props = {
  onClose: () => void;
  onSave: (featureFlag: FeatureFlag) => void;
  featureFlag: FeatureFlag;
};

type Status = 'enabled' | 'disabled' | 'default';
type FeatureFlagStatus = {
  [orgName: string]: Status;
};

export default function (props: Props) {
  const classes = useStyles();
  const {error, isLoading, response} = useAxios<{
    organizations: Array<Organization>;
  }>({
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
    const response = await axios.post<FeatureFlag>(
      `/host/feature/async/${props.featureFlag.id}`,
      createPayload(featureFlagStatus, props.featureFlag),
    );
    props.onSave(response.data);
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>{props.featureFlag.title}</DialogTitle>
      <DialogContent>
        {error && <FormLabel error>{getErrorMessage(error)}</FormLabel>}
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
                  let nextStatus: Status = 'default';
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
        <Button
          onClick={() => void onSave()}
          variant="contained"
          color="primary">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}

function buildStatusObject(
  organizations: Array<Organization>,
  featureFlag: FeatureFlag,
): FeatureFlagStatus {
  const status: FeatureFlagStatus = {};
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

function createPayload(status: FeatureFlagStatus, featureFlag: FeatureFlag) {
  const toCreate: Array<{organization: string; enabled: boolean}> = [];
  const toDelete: Record<string, null> = {};
  const toUpdate: Record<string, {enabled: boolean}> = {};

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
