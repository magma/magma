```jsx
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

/**
 * This file is meant to be viewed using React Styleguidist
 * See the NMS readme.md
 */

import ThemeProvider from '@material-ui/styles/ThemeProvider';
import defaultTheme from '../theme/default';
import {useState} from 'react';

/*
type SelectProps = {
  content: Array<string>,
  defaultValue?: string,
  value: string,
  onChange: string => void,
  testId?: string,
};
*/

const [value, setValue] = useState('');
const content = ['1.0.0', '1.1.0', '1.2.0', '1.2.1'];
<ThemeProvider theme={defaultTheme}>
  Enter version:
  <SelectEdit
    value={value}
    content={content}
    onChange={value => setValue(value)}
  />
</ThemeProvider>
```
