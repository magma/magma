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

import CircularProgress from '@material-ui/core/CircularProgress';
import React from 'react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  loadingContainer: {
    minHeight: 500,
    paddingTop: 200,
    textAlign: 'center',
  },
});

const LoadingFiller = () => {
  const classes = useStyles();
  return (
    <div className={classes.loadingContainer}>
      <CircularProgress size={50} />
    </div>
  );
};

export default LoadingFiller;
