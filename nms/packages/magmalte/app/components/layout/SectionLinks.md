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
import NetworkContext from '../context/NetworkContext';
import defaultTheme from '../../theme/default';
import {getLteSectionsV2} from '../lte/LteSections';
import {
  BrowserRouter as Router,
} from "react-router-dom";

const contextValue = {
  networkId: 'test_network',
  networkType: 'lte',
};

const [_, sections] = getLteSectionsV2(false);


<ThemeProvider theme={defaultTheme}>
  <Router>
    <NetworkContext.Provider value={contextValue}>
      <SectionLinks sections={sections}/>
    </NetworkContext.Provider>
  </Router>
</ThemeProvider>
```
