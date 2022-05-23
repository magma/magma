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

// $FlowFixMe migrated to typescript
import LoadingFiller from './LoadingFiller';
import React from 'react';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  backdrop: {
    alignItems: 'center',
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    bottom: 0,
    display: 'flex',
    justifyContent: 'center',
    left: 0,
    position: 'fixed',
    right: 0,
    top: 0,
    zIndex: '13000',
  },
}));

export default function LoadingFillerBackdrop() {
  const classes = useStyles();
  return (
    <div className={classes.backdrop}>
      <LoadingFiller />
    </div>
  );
}
