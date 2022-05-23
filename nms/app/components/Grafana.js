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
 * @flow
 * @format
 */

import React from 'react';

// $FlowFixMe migrated to typescript
import LoadingFiller from './LoadingFiller';

import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    height: '100%',
    flexGrow: 1,
  },
  dashboardsIframe: {
    width: '100%',
    border: 0,
  },
}));

type Props = {
  grafanaURL: string,
};

export default function GrafanaDashboards(props: Props) {
  const classes = useStyles();
  const [isLoading, setIsLoading] = useState(true);
  return (
    <>
      {isLoading ? <LoadingFiller /> : null}
      <div className={classes.root}>
        <iframe
          src={props.grafanaURL}
          onLoad={() => setIsLoading(false)}
          className={classes.dashboardsIframe}
        />
      </div>
    </>
  );
}
