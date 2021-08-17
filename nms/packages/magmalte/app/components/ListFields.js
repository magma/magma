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

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import IconButton from '@material-ui/core/IconButton';
import React from 'react';
import RemoveCircleOutline from '@material-ui/icons/RemoveCircleOutline';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  container: {
    display: 'block',
    margin: '5px 0',
    whiteSpace: 'nowrap',
    width: '100%',
  },
  input: {
    width: '500px',
    paddingRight: '10px',
  },
  icon: {
    width: '30px',
    height: '30px',
    verticalAlign: 'bottom',
  },
}));

type Props = {
  itemList: Array<string>,
  onChange: (Array<string>) => void,
};

export default function ListFields(props: Props) {
  const classes = useStyles();

  const onChange = (index, value) => {
    const itemList = [...props.itemList];
    itemList[index] = value;
    props.onChange(itemList);
  };

  const removeField = index => {
    const itemList = [...props.itemList];
    itemList.splice(index, 1);
    props.onChange(itemList);
  };

  const addField = () => {
    props.onChange([...props.itemList, '']);
  };

  return (
    <>
      {props.itemList.map((item, index) => (
        <div className={classes.container} key={index}>
          <TextField
            label="Item"
            margin="none"
            value={item}
            onChange={({target}) => onChange(index, target.value)}
            className={classes.input}
          />
          {props.itemList.length !== 1 && (
            <IconButton
              onClick={() => removeField(index)}
              className={classes.icon}>
              <RemoveCircleOutline />
            </IconButton>
          )}
          {index === props.itemList.length - 1 && (
            <IconButton onClick={addField} className={classes.icon}>
              <AddCircleOutline />
            </IconButton>
          )}
        </div>
      ))}
    </>
  );
}
