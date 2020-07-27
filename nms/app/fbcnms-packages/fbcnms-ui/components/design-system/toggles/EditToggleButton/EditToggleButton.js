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
import EditIcon from '@material-ui/icons/Edit';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import FormElementContext from '@fbcnms/ui/components/design-system/Form/FormElementContext';
import IconButton from '@material-ui/core/IconButton';
import {makeStyles} from '@material-ui/styles';

type Props = {
  isOnEdit: boolean,
  onChange: (isOnEdit: boolean) => void,
};

const useStyles = makeStyles(() => ({
  disabled: {
    opacity: 0.5,
    pointerEvents: 'none',
  },
}));

const EditToggleButton = (props: Props) => {
  const {isOnEdit, onChange} = props;
  const classes = useStyles();

  return (
    <FormAction>
      <FormElementContext.Consumer>
        {context => (
          <IconButton
            onClick={() => onChange && onChange(!isOnEdit)}
            className={context.disabled ? classes.disabled : ''}
            color="primary">
            {isOnEdit ? <ArrowBackIcon /> : <EditIcon />}
          </IconButton>
        )}
      </FormElementContext.Consumer>
    </FormAction>
  );
};

export default EditToggleButton;
