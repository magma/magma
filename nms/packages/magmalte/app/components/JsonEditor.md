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

content = {
  "autoupgrade_enabled": true,
  "autoupgrade_poll_interval": 300,
  "checkin_interval": 60,
  "checkin_timeout": 10,
  "dynamic_services": ["td-agent-bit"],
  "feature_flags": {
    "newfeature1": true,
    "newfeature2": false
  },
  "logging": {
    "aggregation": {
      "target_files_by_tag": {
        "enodebd": "/var/log/enodebd.log",
        "mme": "/var/log/mme.log",
        "otherlog": "/var/log/otherlog.log"
      },
      "throttle_interval": "1m",
      "throttle_rate": 1000,
      "throttle_window": 5
    },
    "event_verbosity": 0,
    "log_level": "DEBUG"
  },
  "vpn": {
    "enable_shell": false
  }
};

<ThemeProvider theme={defaultTheme}>
  <JsonEditor
    content={content}
    error={''}
    onSave={state => {content = state;}}
  />
</ThemeProvider>
```
