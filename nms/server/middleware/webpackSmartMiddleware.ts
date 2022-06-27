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
 */

import express from 'express';
import logging from '../../shared/logging';
import webpack from 'webpack';
import webpackHotMiddleware, {MiddlewareOptions} from 'webpack-hot-middleware';
import webpackMiddleware from 'webpack-dev-middleware';

import {RequestHandler} from 'express';
const logger = (logging.getLogger(
  module,
) as unknown) as webpackMiddleware.Options['logger'];

type WebpackMiddlewareOptions = {
  distPath: string;
  devWebpackConfig: webpack.Configuration;
};

type WebpackSmartMiddlewareOptions = WebpackMiddlewareOptions & {
  devMode: boolean;
};

function webpackDevMiddleware(
  options: WebpackMiddlewareOptions,
): RequestHandler {
  const {devWebpackConfig} = options;
  const compiler = webpack(devWebpackConfig);
  const middleware = webpackMiddleware(compiler, {
    publicPath: devWebpackConfig.output?.publicPath,
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
    webpackHotMiddleware(compiler, ({
      reload: true,
    } as unknown) as MiddlewareOptions),
  );
  return router;
}

export default function webpackSmartMiddleware(
  options: WebpackSmartMiddlewareOptions,
): RequestHandler {
  const {devMode, devWebpackConfig, distPath} = options;

  const router = express.Router();
  if (process.env.NODE_ENV === 'test') {
    // Do nothing
  } else if (devMode) {
    router.use(webpackDevMiddleware({devWebpackConfig, distPath}));
  } else {
    // serve built resources from static/dist/ folder
    router.use(
      devWebpackConfig.output?.publicPath || '',
      express.static(distPath),
    );
  }
  return router;
}
