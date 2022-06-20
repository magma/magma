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

import express from 'express';
// $FlowFixMe migrated to typescript
import logging from '../../shared/logging';

import type {Middleware} from 'express';
const logger = logging.getLogger(module);

type WebpackMiddlewareOptions = {|
  distPath: string,
  devWebpackConfig: Object,
|};

type WebpackSmartMiddlewareOptions = {|
  ...WebpackMiddlewareOptions,
  devMode: boolean,
|};

// $FlowIgnore[value-as-type]
function webpackDevMiddleware(options: WebpackMiddlewareOptions): Middleware {
  const {devWebpackConfig} = options;
  const webpack = require('webpack');
  const webpackMiddleware = require('webpack-dev-middleware');
  const webpackHotMiddleware = require('webpack-hot-middleware');
  const compiler = webpack(devWebpackConfig);
  const middleware = webpackMiddleware(compiler, {
    publicPath: devWebpackConfig.output.publicPath,
    contentBase: 'src',
    logger,
    stats: {
      colors: true,
      hash: false,
      timings: true,
      chunks: false,
      chunkModules: false,
      modules: false,
    },
  });

  const router = express.Router();
  router.use(middleware);
  router.use(
    webpackHotMiddleware(compiler, {
      reload: true,
    }),
  );
  return router;
}

export default function webpackSmartMiddleware(
  options: WebpackSmartMiddlewareOptions,
  // $FlowIgnore[value-as-type]
): Middleware {
  const {devMode, devWebpackConfig, distPath} = options;

  const router = express.Router();
  if (process.env.NODE_ENV === 'test') {
    // Do nothing
  } else if (devMode) {
    router.use(webpackDevMiddleware({devWebpackConfig, distPath}));
  } else {
    // serve built resources from static/dist/ folder
    router.use(devWebpackConfig.output.publicPath, express.static(distPath));
  }
  return router;
}
