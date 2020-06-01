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

const paths = require('./paths');
const webpack = require('webpack');

type Options = {
  hot: boolean,
  projectName?: string,
  extraPaths?: string[],
  entry?: Object,
};

function entry(value: any[], hot: boolean) {
  if (!hot) {
    return value;
  }
  return ['webpack-hot-middleware/client', ...value];
}

function createDevWebpackConfig(options: Options) {
  return {
    mode: 'development',
    devtool: 'source-map',
    entry: Object.assign(
      {
        main: entry([paths.appIndexJs], options.hot),
        login: entry([paths.loginJs], options.hot),
      },
      options.entry || {},
    ),
    externals: [
      {
        xmlhttprequest: '{XMLHttpRequest:XMLHttpRequest}',
      },
    ],
    output: {
      pathinfo: true,
      path: paths.distPath,
      filename: '[name].js',
      chunkFilename: 'static/js/[name].chunk.js',
      publicPath: options.projectName
        ? `/${options.projectName}/static/dist/`
        : '/static/dist/',
    },
    plugins: options.hot ? [new webpack.HotModuleReplacementPlugin()] : [],
    module: {
      rules: [
        {
          // "oneOf" will traverse all following loaders until one will
          // match the requirements. When no loader matches it will fall
          // back to the "file" loader at the end of the loader list.
          oneOf: [
            // "url" loader works like "file" loader except that it embeds
            // assets smaller than specified limit in bytes as data URLs to
            // avoid requests.  A missing `test` is equivalent to a match.
            {
              test: [/\.bmp$/, /\.gif$/, /\.jpe?g$/, /\.png$/],
              loader: require.resolve('url-loader'),
              options: {
                limit: 10000,
                name: 'static/media/[name].[hash:8].[ext]',
              },
            },
            // Process JS with Babel.
            {
              test: /\.(js|jsx|mjs)$/,
              include: [
                paths.appSrc,
                paths.packagesDir,
                ...(options.extraPaths || []),
              ],
              loader: require.resolve('babel-loader'),
              options: {
                configFile: '../../babel.config.js',
                // This is a feature of `babel-loader` for webpack (not Babel
                // itself). It enables caching results in
                // ./node_modules/.cache/babel-loader/ directory for faster
                // rebuilds.
                cacheDirectory: true,
              },
            },
            // "postcss" loader applies autoprefixer to our CSS.
            // "css" loader resolves paths in CSS and adds assets as
            // dependencies.
            // "style" loader turns CSS into JS modules that inject <style>
            // tags.
            // In production, we use a plugin to extract that CSS to a file, but
            // in development "style" loader enables hot editing of CSS.
            {
              test: /\.css$/,
              use: [
                require.resolve('style-loader'),
                {
                  loader: require.resolve('css-loader'),
                  options: {
                    importLoaders: 1,
                  },
                },
                {
                  loader: require.resolve('postcss-loader'),
                  options: {
                    // Necessary for external CSS imports to work
                    // https://github.com/facebookincubator/create-react-app/issues/2677
                    ident: 'postcss',
                    plugins: () => [require('postcss-flexbugs-fixes')],
                  },
                },
              ],
            },
            // "file" loader makes sure those assets get served by
            // WebpackDevServer.
            // When you `import` an asset, you get its (virtual) filename.
            // In production, they would get copied to the `build` folder.
            // This loader doesn't use a "test" so it will catch all modules
            // that fall through the other loaders.
            {
              // Exclude `js` files to keep "css" loader working as it injects
              // its runtime that would otherwise processed through "file"
              // loader. Also exclude `html` and `json` extensions so they get
              // processed by webpacks internal loaders.
              exclude: [/\.(js|jsx|mjs)$/, /\.html$/, /\.json$/],
              loader: require.resolve('file-loader'),
              options: {
                name: 'static/media/[name].[hash:8].[ext]',
              },
            },
          ],
        },
        // ** STOP ** Are you adding a new loader?
        // Make sure to add the new loader(s) before the "file" loader.
      ],
    },
    // Some libraries import Node modules but don't use them in the browser.
    // Tell Webpack to provide empty mocks for them so importing them works.
    node: {
      dgram: 'empty',
      fs: 'empty',
      net: 'empty',
      tls: 'empty',
      child_process: 'empty',
    },
    // Turn off performance hints during development because we don't do any
    // splitting or minification in interest of speed. These warnings become
    // cumbersome.
    performance: {
      hints: false,
    },
    optimization: {
      splitChunks: {
        cacheGroups: {
          vendor: {
            test: /[\\/]node_modules[\\/]/,
            chunks: 'all',
            name: 'vendor',
            filename: 'vendor.js',
          },
        },
      },
    },
    devServer: {
      watchOptions: {
        ignored: /node_modules/,
      },
    },
  };
}

module.exports = {
  createDevWebpackConfig,
};
