/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const path = require('path');
const fs = require('fs');
const i18next = require('i18next');
const i18nextMiddleware = require('i18next-express-middleware');
const FilesystemBackend = require('i18next-node-fs-backend');

const DEVELOPMENT = process.env.NODE_ENV !== 'production';
const LOCALE_PARAM = 'locale';

const defaultConfig = {
  preload: ['en_US'],
  fallbackLng: 'en_US',
  backend: {},
  detection: {
    order: ['querystring', 'path', 'cookie'],
    // keys to lookup in the http path (req.params)
    lookupPath: LOCALE_PARAM,
    lookupQuerystring: LOCALE_PARAM,
    nsLookupPath: 'ns',
    // cache the user's language in a cookie
    caches: ['cookie'],
    /**
     * Hack to fix path language detection:
     * By default, if req.params is not present, language detector will attempt
     * to split the url and find locale and ns using positions in the url.
     * This fails for routes which do not include the language in the path.
     * Here's an example:
     * app.use('/heres/my/route', i18nextMiddleware.handle())
     * will detect 'my' as the language, but instead we should detect nothing
     * and fallback to english
     **/
    lookupFromPathIndex: NaN,
  },
};

/**
 * Builds a default working i18n instance.
 * - Uses language detection to pull the user's language from the url path
 * - Stores translations on the local filesystem and serves them asyncronously.
 *
 * @param config can be an object to override default config options or a
 * function. If a function is passed, the default config will be passed and the
 * function's return value will be used as the config.
 **/
export function i18nBuilder(
  config: void | Object | ((defaultConfig: Object) => Object),
) {
  const initConfig =
    typeof config === 'function'
      ? config(defaultConfig)
      : typeof config === 'object'
      ? Object.assign(defaultConfig, config)
      : defaultConfig;
  if (!initConfig.backend) {
    initConfig.backend = initFsBackendOptions();
  }
  i18next
    .use(i18nextMiddleware.LanguageDetector)
    .use(FilesystemBackend)
    .init(initConfig);

  makeLocaleDirectory(
    {locale: i18next.options.fallbackLng[0], namespace: 'translation'},
    i18next,
  );

  return i18next;
}

export function initFsBackendOptions(
  options?: {localesDir: string} = {localesDir: './locales'},
) {
  const LOCALES_DIR = options.localesDir || './locales';
  const fsBackendOptions = {
    // path where resources get loaded from
    loadPath: path.join(LOCALES_DIR, '/{{lng}}/{{ns}}.json'),
    // path to post missing resources
    addPath: path.join(LOCALES_DIR, '/{{lng}}/{{ns}}.json'),
    // jsonIndent to use when storing json files
    jsonIndent: 2,
  };
  return fsBackendOptions;
}

/**
 * Gets all translations for a specific locale
 * Language is detected
 * Namespace is pulled from the url path
 */
export function getLocaleHandler(i18NextInstance: any) {
  return [
    i18nextMiddleware.handle(i18NextInstance),
    (req: any, res: any) => {
      // this comes from language detection
      const locale = req.i18n.language;
      let ns = req.params[req.i18n.options.detection.nsLookupPath || 'ns'];
      if (!ns) {
        ns = 'translation';
      }
      const getIsLangLoaded = () => req.i18n.hasResourceBundle(locale, ns);
      return new Promise((resolve, reject) => {
        if (getIsLangLoaded()) {
          resolve(req.i18n.getResourceBundle(locale, ns));
        } else {
          req.i18n.loadLanguages(locale, err => {
            const isLangLoaded = getIsLangLoaded();
            if (err || !isLangLoaded) {
              return reject(err);
            } else if (isLangLoaded) {
              return resolve(req.i18n.getResourceBundle(locale, ns));
            } else {
              reject(
                new Error('could not load requested language or a fallback'),
              );
            }
          });
        }
      })
        .then(resources => {
          res.send(resources);
        })
        .catch(_err => {
          return res.status(404).send({message: 'error loading language'});
        });
    },
  ];
}

/**
 * Creates a directory for a locale based on the loadPath parameter provided to
 * fs backend options.
 **/
export function makeLocaleDirectory(
  {locale, namespace}: {locale: string, namespace?: string},
  i18NextInstance: any,
) {
  if (typeof namespace === 'undefined' || namespace.trim() === '') {
    namespace = 'translation';
  }
  const addPath = i18NextInstance.services.interpolator.interpolate(
    i18next.options.backend.addPath,
    {
      lng: locale,
      ns: namespace,
    },
  );
  const addPathDirectory = path.dirname(addPath);
  if (!fs.existsSync(addPathDirectory)) {
    mkDirByPathSync(addPathDirectory);
  }
}

/**
 * DEVELOPMENT ONLY
 * Source key extraction - saves untranslated keys to the configured backend
 **/
export function addMissingKeysHandler(i18NextInstance: any) {
  return [
    developmentOnly(
      'translation key extraction is only enabled in development',
    ),
    i18nextMiddleware.handle(i18NextInstance),
    i18nextMiddleware.missingKeyHandler(i18NextInstance, {
      lngParam: LOCALE_PARAM,
    }),
  ];
}

function developmentOnly(message) {
  return (req: any, res: any, next: any) => {
    if (!DEVELOPMENT) {
      return res
        .status(403)
        .send({error: message || 'Development only feature'});
    }
    return next();
  };
}

/**
 * node v6 does not support mkdir -p via the recursive flag
 * https://stackoverflow.com/a/40686853/2188014
 **/
function mkDirByPathSync(targetDir) {
  const sep = path.sep;
  const initDir = path.isAbsolute(targetDir) ? sep : '';

  return targetDir.split(sep).reduce((parentDir, childDir) => {
    const curDir = path.resolve('.', parentDir, childDir);
    try {
      fs.mkdirSync(curDir);
    } catch (err) {
      if (err.code === 'EEXIST') {
        // curDir already exists!
        return curDir;
      }

      /*
       * To avoid `EISDIR` error on Mac and
       * `EACCES`-->`ENOENT` and `EPERM` on Windows.
       */
      if (err.code === 'ENOENT') {
        // Throw the original parentDir error on curDir `ENOENT` failure.
        throw new Error(`EACCES: permission denied, mkdir '${parentDir}'`);
      }

      const caughtErr = ['EACCES', 'EPERM', 'EISDIR'].indexOf(err.code) > -1;
      if (!caughtErr || (caughtErr && curDir === path.resolve(targetDir))) {
        throw err; // Throw if it's just the last created dir.
      }
    }

    return curDir;
  }, initDir);
}
