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

import * as React from 'react';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import ArrowForwardIcon from '@material-ui/icons/ArrowForward';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  iconContainer: {
    border: '1px solid #b3b3b3',
    borderRadius: '100%',
    width: '20px',
    height: '20px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    '&:hover': {
      boxShadow: '0 0 0 5px rgba(53, 120, 229, 0.28)',
      borderColor: theme.palette.primary.main,
      backgroundColor: theme.palette.primary.main,
      '& $icon': {
        color: theme.palette.blueGrayDark,
      },
      cursor: 'pointer',
    },
  },
  icon: {
    color: theme.palette.gray50,
    '&&': {
      fontSize: '20px',
    },
  },
}));

type Props = {
  expanded: boolean,
  onClick: () => void,
};

const ExpandButton = ({expanded, onClick}: Props) => {
  const classes = useStyles();
  return (
    <div className={classes.iconContainer} onClick={onClick}>
      {expanded ? (
        <ArrowBackIcon className={classes.icon} />
      ) : (
        <ArrowForwardIcon className={classes.icon} />
      )}
    </div>
  );
};

export default ExpandButton;
