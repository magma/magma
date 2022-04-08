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
import DataUsageIcon from '@material-ui/icons/DataUsage';
import defaultTheme from '../../theme/default';

<ThemeProvider theme={defaultTheme}>
  <CardTitleRow
    key={'card_title_row'}
    icon={DataUsageIcon}
    label={'Card Title Row'}
  />
</ThemeProvider>
```
