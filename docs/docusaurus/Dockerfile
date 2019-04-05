FROM node:8.11.4

WORKDIR /app/website

EXPOSE 3000 35729
COPY docusaurus/package.json /app/website/package.json
RUN yarn install
COPY docusaurus /app/website
COPY readmes /app/docs
RUN yarn build

CMD ["yarn", "start"]
