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

'use strict';

const paths = require('./paths');
const webpack = require('webpack');
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer')
  .BundleAnalyzerPlugin;
const ManifestPlugin = require('webpack-manifest-plugin');
const TerserPlugin = require('terser-webpack-plugin');

module.exports = {
  mode: 'production',
  devtool: 'cheap-module-eval-source-map',
  entry: {
    main: [paths.appIndexJs],
    login: [paths.loginJs],
    host: [paths.hostJs],
  },
  resolve: {
    extensions: ['.ts', '.tsx', '.js', '.json']
  },
  externals: [
    {
      xmlhttprequest: '{XMLHttpRequest:XMLHttpRequest}',
    },
  ],
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
            loader: require.resolve('url-loader'),
            options: {
              limit: 10000,
              name: 'static/media/[name].[hash:8].[ext]',
            },
            test: [/\.bmp$/, /\.gif$/, /\.jpe?g$/, /\.png$/],
          },
          // Process JS with Babel.
          {
            loader: require.resolve('babel-loader'),
            options: {
              rootMode: 'upward',
              // This is a feature of `babel-loader` for webpack (not Babel
              // itself). It enables caching results in
              // ./node_modules/.cache/babel-loader/ directory for faster
              // rebuilds.
              cacheDirectory: true,
            },
            test: /\.(js|jsx|mjs|ts|tsx)$/,
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
            exclude: [/\.(js|jsx|mjs|ts|tsx)$/, /\.html$/, /\.json$/],
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
    child_process: 'empty',
    dgram: 'empty',
    fs: 'empty',
    net: 'empty',
    tls: 'empty',
  },
  output: {
    chunkFilename: '[name].[chunkhash].chunk.js',
    filename: '[name].[chunkhash].js',
    path: paths.distPath,
    pathinfo: true,
    publicPath: '/nms/static/dist/',
  },
  plugins: [
    new ManifestPlugin(),
    new BundleAnalyzerPlugin({
      analyzerMode: 'static',
      reportFilename: 'report.html',
    }),
    // remove excess locales in moment bloating the bundle
    new webpack.ContextReplacementPlugin(/moment[/\\]locale$/, /en/),
  ],
  optimization: {
    minimizer: [
      new TerserPlugin({
        chunkFilter: chunk => chunk.name !== 'vendor',
        parallel: true,
      }),
    ],
    splitChunks: {
      cacheGroups: {
        vendor: {
          test: /[\\/]node_modules[\\/]/,
          chunks: 'initial',
          name: 'vendor',
          priority: 10,
          enforce: true,
        },
      },
    },
  },
};
