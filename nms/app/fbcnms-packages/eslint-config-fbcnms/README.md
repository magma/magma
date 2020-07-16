# eslint-config-fbcnms

This package provides Facebook NMS' eslint config as an extensible shared config. This is only used internally, use at your own risk.

## Usage

Our default export contains all of our ESLint rules, including ECMAScript 6+ and React. It requires `eslint`, `eslint-plugin-relay`, `eslint-plugin-header`, `eslint-plugin-import`, `eslint-plugin-node`, `eslint-plugin-lint`, `eslint-plugin-sort-imports-es6-autofix`, `eslint-config-fb-strict`.

To install:

Yarn
```
yarn add eslint-config-fbcnms eslint-config-fb-strict eslint-plugin-relay eslint-plugin-header eslint-plugin-import eslint-plugin-node eslint-plugin-lint eslint-plugin-sort-imports-es6-autofix --dev
```

npm
```
npm install --save-dev eslint-config-fbcnms eslint-config-fb-strict eslint-plugin-relay eslint-plugin-header eslint-plugin-import eslint-plugin-node eslint-plugin-lint eslint-plugin-sort-imports-es6-autofix
```

- Add "extends": "eslint-config-fbcnms" to your .eslintrc

## License
[BSD-2-Clause](https://opensource.org/licenses/BSD-2-Clause)
