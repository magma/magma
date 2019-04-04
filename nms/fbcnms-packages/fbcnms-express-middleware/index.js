/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import bodyParser from 'body-parser';
import compression from 'compression';
import cookieParser from 'cookie-parser';
import csrf from 'csurf';
import express from 'express';
import helmet from 'helmet';
import logging from '@fbcnms/logging';
import session from 'express-session';
import webpack from 'webpack';

export type Options = {|
  distPath: string,
  sessionStore: {},
  devMode: boolean,
  sessionStore: session.Session,
  sessionToken: string,
  devWebpackConfig: Object,
|};

import type {
  ExpressResponse,
  ExpressRequest,
  ExpressApplication,
  NextFunction,
} from 'express';

import type {OrganizationMiddlewareRequest} from './organization';

export type FBCNMSMiddleWareRequest = {
  csrfToken: () => string, // from csrf
  body: Object, // from bodyParser
} & OrganizationMiddlewareRequest;

const logger = logging.getLogger(module);

export function middleware(app: ExpressApplication, options: Options) {
  const {
    devMode,
    distPath,
    sessionStore,
    sessionToken,
    devWebpackConfig,
  } = options;

  app.set('trust proxy', 1);
  app.use(helmet());
  app.use(bodyParser.json({limit: '1mb'})); // parse json
  // parse application/x-www-form-urlencoded
  app.use(bodyParser.urlencoded({limit: '1mb', extended: false}));
  app.use(cookieParser());
  app.use(compression());
  app.use(logging.getHttpLogger(module));

  app.use(
    session({
      cookie: {
        secure: !devMode,
      },
      // Used to sign the session cookie
      secret: sessionToken,
      resave: false,
      saveUninitialized: true,
      store: sessionStore,
      unset: 'destroy',
    }),
  );

  sessionStore.sync();

  // Use csrf middleware (uses session, must be declared after)
  app.use(
    csrf({cookie: true, value: (req: ExpressRequest) => req.cookies.csrfToken}),
  );

  if (devMode) {
    // serve developer, non-minified build
    const compiler = webpack(devWebpackConfig);
    const webpackMiddleware = require('webpack-dev-middleware');
    const webpackHotMiddleware = require('webpack-hot-middleware');
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
    app.use(middleware);
    app.use(webpackHotMiddleware(compiler));
  } else {
    // serve built resources from static/dist/ folder
    app.use(devWebpackConfig.output.publicPath, express.static(distPath));
  }

  return function(
    req: FBCNMSMiddleWareRequest,
    res: ExpressResponse,
    next: NextFunction,
  ) {
    res.cookie('csrfToken', req.csrfToken ? req.csrfToken() : '', {
      sameSite: true,
      httpOnly: true,
    });
    next();
  };
}
