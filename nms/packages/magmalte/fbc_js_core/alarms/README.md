# FBCNMS Alarms

This package provides UI elements for configuring the alerting system with Prometheus and Alertmanager via Magma's prometheus configmanager. To see an implementation, take a look at the Magma [NMS](https://github.com/magma/magma/tree/master/nms).

![Screenshot](https://raw.githubusercontent.com/magma/fbc-js-core/main/fbcnms-packages/fbcnms-alarms/create_alert_screenshot.png)

To install:

Yarn

```
yarn add @fbcnms/alarms
```

npm

```
npm install @fbcnms/alarms
```

## License

[BSD-2-Clause](https://opensource.org/licenses/BSD-2-Clause)

# Development

The easiest way to do development on FBCNMS Alarms is to use a workflow similar to npm / yarn link. Symlinking this package into your own app's node_modules folder and running it that way. This can cause a number of problems though, particularly when it comes to dependency resolution and Webpack/Babel.

We use [yalc](https://github.com/wclr/yalc) to resolve the afformentioned issues with using npm / yarn link.

First install yalc globally

Yarn:

```
yarn global add yalc
```

NPM:

```
npm i yalc -g
```

Next, start the `dev` yarn script to watch source files for changes and publish to the local yalc repo.

```
yarn run dev
```

Next, cd to your project. This should be the same project which has a dependency on @fbcnms/alarms in its package.json.

```
yalc link @fbcnms/alarms
```

Your project is now able to resolve @fbcnms/alarms.

## Setting up webpack/babel

In webpack.config.js using babel-loader:

First enumerate all the @fbcnms/ packages and their paths

```
const path = require('path');
const packageJson = require('./package.json');

const fbcnmsPackages = Object.keys(packageJson.dependencies)
  .filter(key => key.includes('@fbcnms'))
  .map(pkg =>
    path.join(
      path.resolve(require.resolve(path.join(pkg, 'package.json'))),
      '../',
    ),
  );
```

Next, add them to your babel-loader setup:

```
{
    test: /\.(js|jsx|mjs)$/,
    include: [
        'your app dir',
        fbcnmsPackages,
    ],
    loader: require.resolve('babel-loader'),
}
```

In the future this may not be necessary, but some projects like to import the untransformed sources since we don't currently publish flow-defs.
