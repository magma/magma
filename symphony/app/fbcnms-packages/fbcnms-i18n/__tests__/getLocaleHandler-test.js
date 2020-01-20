/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import express from 'express';
import i18nextMiddleware from 'i18next-express-middleware';
import request from 'supertest';
import {getLocaleHandler} from '..';

test('if requested locale is loaded, serve it', async () => {
  const i18n = getI18n({
    resources: {
      en: {
        translation: {
          hello: 'hello',
        },
      },
      es: {
        translation: {
          hello: 'hola',
        },
      },
    },
  });

  const app = makeApp(i18n);
  let response = await request(app)
    .get('/?locale=en')
    .expect(200);
  expect(response.body.hello).toBe('hello');

  response = await request(app)
    .get('/?locale=es')
    .expect(200);
  expect(response.body.hello).toBe('hola');
});

test('if requested locale does not exist, return a 404 error', async () => {
  const i18n = getI18n({
    resources: {
      en: {
        translation: {
          hello: 'hello',
        },
      },
      es: {
        translation: {
          hello: 'hola',
        },
      },
    },
  });

  const app = makeApp(i18n);
  const response = await request(app)
    .get('/?locale=fr')
    .expect(404);
  expect(response.body.hello).toBe(undefined);
  expect(response.body.message).toBe('error loading language');
});

function makeApp(i18n: any) {
  const app = express();
  // $FlowFixMe - flow can't handle express middleware arrays
  app.get('/', getLocaleHandler(i18n));
  return app;
}

function getI18n(config = {}) {
  const i18n = require('i18next');
  const initConfig: any = {
    fallbackLng: false,
    detection: {
      order: ['querystring'],
      lookupQuerystring: 'locale',
    },
    ...config,
  };
  i18n.use(i18nextMiddleware.LanguageDetector);
  i18n.init(initConfig);
  return i18n;
}
