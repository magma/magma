Docusaurus Pages
===

If you want to cut a new versioned documentation release, delete everything
in the Dockerfile after WORKDIR /app/website, then from this directory:

```bash
docker build --no-cache -f docusaurus/Dockerfile -t docusaurus-doc .
docker run -it -p3000:3000 -v $(pwd)/docusaurus:/app/website -v $(pwd)/readmes:/app/docs docusaurus-doc bash
```

Inside the container,

```bash
yarn install
yarn build
```

Then you can use `yarn` commands to create a new versioned release:

```
yarn run version 1.X.0
```

Commit all the new generated files **except for the `docusaurus/node_modules` 
directory** and tweak the sidebars if you need to. Revert the changes to the 
docusaurus Dockerfile and run `create_docusaurus_website.sh` to preview your
changes.
