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

Then you can use `yarn` commands to create a new versioned release.
