#!/usr/bin/env node
/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

require('@fbcnms/babel-register');

const {main} = require('relay-compiler/lib/bin/RelayCompilerMain');
const yargs = require('yargs');
const glob = require('glob');
const prependFile = require('prepend-file');
const fs = require('fs');

const HEADER = `/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 `;

function markFilesGenerated(argv) {
  glob
    .sync('**/__generated__/*.graphql.js', {
      root: argv.src,
    })
    .forEach(file => {
      fs.readFile(file, (err, data) => {
        if (err) throw err;
        if (data.indexOf(HEADER) < 0) {
          prependFile(file, HEADER);
        }
      });
    });
}

let RelayConfig;
try {
  // eslint-disable-next-line no-eval
  RelayConfig = eval('require')('relay-config');
} catch (_) {}

const options = {
  schema: {
    describe: 'Path to schema.graphql or schema.json',
    demandOption: true,
    type: 'string',
    array: false,
  },
  src: {
    describe: 'Root directory of application code',
    demandOption: true,
    type: 'string',
    array: false,
  },
  watch: {
    describe: 'If specified, watches files and regenerates on changes',
    type: 'boolean',
    default: false,
  },
};

// Load external config
const config = RelayConfig && RelayConfig.loadConfig();

// Parse CLI args
const argv = yargs
  .usage('Create Relay generated files\n\n$0 --schema <path> --src <path>')
  .options(options)
  .config(config)
  .help().argv;

const compilerConfig = {
  src: argv.src,
  schema: argv.schema,
  watch: argv.watch,
  validate: false,
  noFutureProofEnums: false,
  language: 'javascript',
  include: ['**'],
  exclude: ['**/node_modules/**', '**/__mocks__/**', '**/__generated__/**'],
  verbose: false,
  quiet: false,
  watchman: true,
  validate: false,
};

// Start the application
main(compilerConfig)
  .then(_ => markFilesGenerated(argv))
  .catch(error => {
    console.error(String(error.stack || error));
    process.exit(1);
  });
